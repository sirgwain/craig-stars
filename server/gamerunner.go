package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/ai"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

var errNotFound = errors.New("resource was not found")

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
	"#FF4500",
	"#FF8C00",
	"#008000",
	"#00FA9A",
	"#7FFFD4",
	"#8A2BE2",
	"#FF1493",
	"#D2691E",
	"#F0FFF0",
}

type GameRunner interface {
	HostGame(hostID int64, settings *cs.GameSettings) (*cs.FullGame, error)
	JoinGame(gameID int64, userID int64, raceID int64, color string) error
	GenerateUniverse(game *cs.Game) error
	LoadPlayerGame(gameID int64, userID int64) (*cs.GameWithPlayers, *cs.FullPlayer, error)
	SubmitTurn(gameID int64, userID int64) error
	CheckAndGenerateTurn(gameID int64) (TurnGenerationCheckResult, error)
	GenerateTurn(gameID int64) error
}

type gameRunner struct {
	db     DBClient
	client cs.Gamer
	config config.Config
}

func NewGameRunner(db DBClient, config config.Config) GameRunner {
	return &gameRunner{db, cs.NewGamer(), config}
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

	// generate an invite hash for the game
	game.Hash = game.GenerateHash(gr.config.Game.InviteLinkSalt)
	if err := gr.db.UpdateGame(game); err != nil {
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
			player.DefaultHullSet = playerSetting.DefaultHullSet
			player.Ready = true
			players = append(players, player)
		} else if playerSetting.Type == cs.NewGamePlayerTypeAI {
			log.Debug().Int64("hostID", hostID).Msg("Adding ai player to game")
			race := *cs.NewRace()
			if playerSetting.Race.Name != "" {
				race = playerSetting.Race
			}
			player := gr.client.NewPlayer(0, race, &game.Rules)
			player.GameID = game.ID
			player.Num = i + 1
			player.AIControlled = true
			player.Name = fmt.Sprintf("AI Player %d", player.Num)
			player.Color = playerSetting.Color
			player.DefaultHullSet = playerSetting.DefaultHullSet
			player.Ready = true
			players = append(players, player)
		} else if playerSetting.Type == cs.NewGamePlayerTypeOpen {
			log.Debug().Uint("openPlayerSlots", game.OpenPlayerSlots).Msg("Added open player slot to game")
			race := cs.NewRace()
			player := gr.client.NewPlayer(0, *race, &game.Rules)
			player.GameID = game.ID
			player.Num = i + 1
			player.Name = "Open Slot"
			player.Color = playerSetting.Color
			player.DefaultHullSet = playerSetting.DefaultHullSet
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
				player.Color = fmt.Sprintf("#%s", hex.EncodeToString(color))
			}
		}
	}

	// load the tech store for this game
	techStore, err := gr.loadTechStore(game)
	if err != nil {
		return nil, err
	}

	universe := cs.NewUniverse(&game.Rules)
	fullGame := &cs.FullGame{
		Game:      game,
		Universe:  &universe,
		TechStore: techStore,
		Players:   players,
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
func (gr *gameRunner) JoinGame(gameID int64, userID int64, raceID int64, color string) error {

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
	player.Ready = true
	player.Color = color

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

	// fix duplicate colors
	// TODO: this is fragile
	for i, p := range fullGame.Players {
		if p.Color == color && p.Num != player.Num {
			if player.Num-1 < len(colors) {
				player.Color = colors[player.Num-1]
			} else {
				color := make([]byte, 3)
				rand.Read(color)
				player.Color = fmt.Sprintf("#%s", hex.EncodeToString(color))
			}

			if err := gr.db.UpdatePlayer(player); err != nil {
				return fmt.Errorf("update open slot player %s for game %d: %w", p, gameID, err)
			}

			fullGame.Players[i] = player

		}
	}

	fullGame.OpenPlayerSlots--
	if err := gr.db.UpdateGame(fullGame.Game); err != nil {
		return fmt.Errorf("save game %d: %w", gameID, err)
	}

	log.Info().Int64("GameID", gameID).Int64("UserID", userID).Str("Race", player.Race.Name).Msgf("Joined game %s", fullGame.Name)

	return nil
}

func (gr *gameRunner) GenerateUniverse(game *cs.Game) error {
	defer timeTrack(time.Now(), "GenerateUniverse")
	fullGame, err := gr.loadGame(game.ID)
	if err != nil {
		return fmt.Errorf("load game %d: %w", game.ID, err)
	}

	return gr.generateUniverse(fullGame)
}

// load a player and the light version of the player game
func (gr *gameRunner) LoadPlayerGame(gameID int64, userID int64) (*cs.GameWithPlayers, *cs.FullPlayer, error) {

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

	if err := gr.db.SubmitPlayerTurn(gameID, player.Num, true); err != nil {
		return fmt.Errorf("submitting player turn %w", err)
	}
	return nil
}

// CheckAndGenerateTurn check a game for every player submitted and generate a turn
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
		if !player.AIControlled {
			continue
		}
		// TODO: make this use copies to ensure the ai only updates orders?
		// TODO: ai only ai processing
		pmo := fullGame.Universe.GetPlayerMapObjects(player.Num)
		ai := ai.NewAIPlayer(fullGame.Game, fullGame.TechStore, player, pmo)
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

// load the tech store into a game
func (gr *gameRunner) loadTechStore(game *cs.Game) (*cs.TechStore, error) {
	if game.Rules.TechsID == 0 {
		return &cs.StaticTechStore, nil
	}
	techs, err := gr.db.GetTechStore(game.Rules.TechsID)
	if err != nil {
		return nil, err
	}
	return techs, nil
}
