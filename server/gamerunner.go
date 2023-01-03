package server

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/ai"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

type TurnGenerationCheckResult uint
type DBClient db.Client

const (
	TurnNotGenerated TurnGenerationCheckResult = iota
	TurnGenerated
)

var colors = []string{
	"#C33232",
	"#1F8BA7",
	"#43A43E",
	"#8D29CB",
	"#B88628",
}

type GameRunner interface {
	HostGame(hostID int64, settings *cs.GameSettings) (*cs.FullGame, error)
	JoinGame(gameID int64, userID int64, raceID int64) error
	GenerateUniverse(gameID int64, userID int64) error
	LoadPlayerGame(gameID int64, userID int64) (*cs.Game, *cs.FullPlayer, error)
	SubmitTurn(gameID int64, userID int64) error
	CheckAndGenerateTurn(gameID int64) (TurnGenerationCheckResult, error)
	GenerateTurn(gameID int64) error
}

type gameRunner struct {
	db     DBClient
	client cs.Gamer
}

func NewGameRunner(db DBClient) GameRunner {
	return &gameRunner{db, cs.NewGamer()}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// host a new game
func (gr *gameRunner) HostGame(hostID int64, settings *cs.GameSettings) (*cs.FullGame, error) {
	game := gr.client.CreateGame(hostID, *settings)

	user, err := gr.db.GetUser(hostID)
	if err != nil {
		return nil, err
	}

	// create the game in the database
	err = gr.db.CreateGame(game)
	if err != nil {
		return nil, err
	}

	game.NumPlayers = len(settings.Players)
	players := make([]*cs.Player, 0, len(settings.Players))

	for i, playerSetting := range settings.Players {
		if playerSetting.Type == cs.NewGamePlayerTypeHost {

			// add the host player with the host race
			race, err := gr.db.GetRace(playerSetting.RaceID)
			if err != nil {
				return nil, err
			}
			if race.UserID != hostID {
				return nil, fmt.Errorf("user %d does not own Race %d", hostID, race.ID)
			}
			log.Debug().Int64("hostID", hostID).Msg("Adding host to game")
			player := gr.client.NewPlayer(hostID, *race, &game.Rules)
			player.GameID = game.ID
			player.Num = i + 1
			player.Name = user.Username
			player.Color = playerSetting.Color
			players = append(players, player)
		} else if playerSetting.Type == cs.NewGamePlayerTypeAI {
			log.Debug().Int64("hostID", hostID).Msg("Adding ai player to game")
			race := cs.NewRace()
			race.Name = "Cool AI"
			race.PluralName = "Cool AIs"
			player := gr.client.NewPlayer(hostID, *race, &game.Rules)
			player.GameID = game.ID
			player.Num = i + 1
			player.AIControlled = true
			player.Name = "AI Player"
			player.Color = playerSetting.Color
			players = append(players, player)
		} else if playerSetting.Type == cs.NewGamePlayerTypeOpen {
			log.Debug().Uint("openPlayerSlots", game.OpenPlayerSlots).Msg("Added open player slot to game")
			race := cs.NewRace()
			player := gr.client.NewPlayer(hostID, *race, &game.Rules)
			player.GameID = game.ID
			player.Num = i + 1
			player.Name = "Open Slot"
			player.Color = playerSetting.Color
			players = append(players, player)
			game.OpenPlayerSlots++
		}
	}

	for _, player := range players {
		if player.Color == "" || (player.Num != 1 && player.Color == "#0000FF") {
			if player.Num-1 < len(colors) {
				player.Color = colors[player.Num-1]
			} else {
				color := make([]byte, 3)
				rand.Read(color)
				player.Color = fmt.Sprintf("#%X%X%X", color[0], color[1], color[2])
			}
		}
	}

	universe := cs.NewUniverse(&game.Rules)
	fullGame := &cs.FullGame{
		Game:     game,
		Players:  players,
		Universe: &universe,
	}

	gr.db.UpdateFullGame(fullGame)

	// generate the universe if this game is done
	if game.OpenPlayerSlots == 0 {
		err = gr.generateUniverse(fullGame)
		if err != nil {
			return nil, err
		}
	}

	return fullGame, nil
}

// add a player to an existing game
func (gr *gameRunner) JoinGame(gameID int64, userID int64, raceID int64) error {

	user, err := gr.db.GetUser(userID)
	if err != nil {
		return fmt.Errorf("unable to load user %d: %w", userID, err)
	}
	if user == nil {
		return fmt.Errorf("no user for id %d found. %w", userID, errNotFound)
	}

	fullGame, err := gr.loadGame(gameID)
	if err != nil {
		return fmt.Errorf("unable to load game %d: %w", gameID, err)
	}

	race, err := gr.db.GetRace(raceID)
	if err != nil {
		log.Error().Err(err).Int64("ID", raceID).Msg("get race from database")
		return fmt.Errorf("failed to get race from database")
	}

	// validate
	if race == nil {
		return fmt.Errorf("no race for id %d found. %w", raceID, errNotFound)
	}

	if race.UserID != userID {
		return fmt.Errorf("user %d does not control race %d", userID, raceID)
	}

	if fullGame == nil {
		return fmt.Errorf("no game for id %d found. %w", gameID, errNotFound)
	}

	if fullGame.OpenPlayerSlots == 0 {
		return fmt.Errorf("no slots left in this game")
	}

	player := gr.client.NewPlayer(userID, *race, &fullGame.Rules)
	player.GameID = fullGame.ID
	player.Name = user.Username

	// claim an open slot
	for i, p := range fullGame.Players {
		if p.Name == "Open Slot" {
			// take over this empty player
			player.Num = p.Num
			player.ID = p.ID
			if err := gr.db.UpdatePlayer(player); err != nil {
				return fmt.Errorf("update open slot player %s for game %d: %w", p, gameID, err)
			}

			fullGame.Players[i] = player
			break
		}
	}

	fullGame.OpenPlayerSlots--
	if err := gr.db.UpdateGame(fullGame.Game); err != nil {
		return fmt.Errorf("save game %d: %w", gameID, err)
	}

	log.Info().Int64("GameID", gameID).Int64("UserID", userID).Str("Race", player.Race.Name).Msgf("Joined game %s", fullGame.Name)

	return nil
}

func (gr *gameRunner) GenerateUniverse(gameID int64, userID int64) error {
	defer timeTrack(time.Now(), "GenerateUniverse")
	fullGame, err := gr.loadGame(gameID)
	if err != nil {
		return fmt.Errorf("load game %d: %w", gameID, err)
	}

	if fullGame.HostID != userID {
		return fmt.Errorf("user %d is not the host for %d", userID, gameID)
	}

	err = gr.generateUniverse(fullGame)
	return err
}

// load a player and the light version of the player game
func (gr *gameRunner) LoadPlayerGame(gameID int64, userID int64) (*cs.Game, *cs.FullPlayer, error) {

	game, err := gr.db.GetGame(gameID)

	if err != nil {
		return nil, nil, err
	}

	if game.Rules.TechsID == 0 {
		game.Rules.WithTechStore(&cs.StaticTechStore)
	} else {
		techs, err := gr.db.GetTechStore(game.Rules.TechsID)
		if err != nil {
			return nil, nil, err
		}
		game.Rules.WithTechStore(techs)
	}

	player, err := gr.db.GetFullPlayerForGame(gameID, userID)

	if err != nil {
		return nil, nil, err
	}

	return game, player, nil
}

// submit a turn for a player
func (gr *gameRunner) SubmitTurn(gameID int64, userID int64) error {
	player, err := gr.db.GetLightPlayerForGame(gameID, userID)
	if err != nil {
		return fmt.Errorf("find player for user %d, game %d: %w", userID, gameID, err)
	}

	if player == nil {
		return fmt.Errorf("player for user %d, game %d not found", userID, gameID)
	}

	player.SubmittedTurn = true
	gr.db.UpdateLightPlayer(player)
	return nil
}

//CheckAndGenerateTurn check a game for every player submitted and generate a turn
func (gr *gameRunner) CheckAndGenerateTurn(gameID int64) (TurnGenerationCheckResult, error) {
	defer timeTrack(time.Now(), "CheckAndGenerateTurn")
	fullGame, err := gr.loadGame(gameID)

	if err != nil {
		return TurnNotGenerated, err
	}

	// if everyone submitted their turn, generate a new turn
	if gr.client.CheckAllPlayersSubmitted(fullGame.Players) {
		if err := gr.generateTurn(fullGame); err != nil {
			return TurnNotGenerated, err
		} else {
			return TurnGenerated, nil
		}
	}

	return TurnNotGenerated, nil
}

// GenerateTurn generate a new turn, regardless of whether player's have all submitted their turns
func (gr *gameRunner) GenerateTurn(gameID int64) error {
	defer timeTrack(time.Now(), "GenerateTurn")
	fullGame, err := gr.loadGame(gameID)

	if err != nil {
		return err
	}

	return gr.generateTurn(fullGame)
}

// load a full game
func (gr *gameRunner) loadGame(gameID int64) (*cs.FullGame, error) {
	defer timeTrack(time.Now(), "loadGame")

	fullGame, err := gr.db.GetFullGame(gameID)

	if err != nil {
		return nil, fmt.Errorf("load game %d: %w", gameID, err)
	}

	return fullGame, nil
}

func (gr *gameRunner) generateUniverse(fullGame *cs.FullGame) error {

	universe, err := gr.client.GenerateUniverse(fullGame.Game, fullGame.Players)
	if err != nil {
		return fmt.Errorf("generate universe for game %d: %w", fullGame.ID, err)
	}

	// ai processing
	gr.processAITurns(fullGame)

	fullGame.State = cs.GameStateWaitingForPlayers
	fullGame.Universe = universe

	err = gr.db.UpdateFullGame(fullGame)
	if err != nil {
		return fmt.Errorf("save game %d: %w", fullGame.ID, err)
	}

	return nil
}

// process an the ai player's turns
func (gr *gameRunner) processAITurns(fullGame *cs.FullGame) {
	for _, player := range fullGame.Players {
		// TODO: make this use copies to ensure the ai only updates orders?
		// TODO: ai only ai processing
		pmo := fullGame.Universe.GetPlayerMapObjects(player.Num)
		ai := ai.NewAIPlayer(player, pmo)
		ai.ProcessTurn()

		if player.AIControlled {
			player.SubmittedTurn = true
		}
	}
}

// generate a turn for a game
func (gr *gameRunner) generateTurn(fullGame *cs.FullGame) error {
	// if everyone submitted their turn, generate a new turn
	if err := gr.client.GenerateTurn(fullGame.Game, fullGame.Universe, fullGame.Players); err != nil {
		return fmt.Errorf("generate turn -> %w", err)
	}

	// ai processing
	gr.processAITurns(fullGame)

	if err := gr.db.UpdateFullGame(fullGame); err != nil {
		return fmt.Errorf("save game after turn generation -> %w", err)
	}
	return nil

}
