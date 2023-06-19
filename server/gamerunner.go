package server

import (
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

type GameRunner interface {
	HostGame(hostID int64, settings *cs.GameSettings) (*cs.FullGame, error)
	AddPlayer(gameID int64, userID int64, race *cs.Race) error
	GenerateUniverse(gameID int64) (*cs.Universe, error)
	LoadGame(gameID int64) (*cs.FullGame, error)
	LoadPlayerGame(gameID int64, userID int64) (*cs.Game, *cs.FullPlayer, error)
	SubmitTurn(gameID int64, userID int64) error
	CheckAndGenerateTurn(gameID int64) (TurnGenerationCheckResult, error)
}

type gameRunner struct {
	db DBClient
}

func NewGameRunner(db DBClient) GameRunner {
	return &gameRunner{db}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// host a new game
func (gr *gameRunner) HostGame(hostID int64, settings *cs.GameSettings) (*cs.FullGame, error) {
	client := cs.NewClient()
	game := client.CreateGame(hostID, *settings)

	// create the game in the database
	err := gr.db.CreateGame(&game)
	if err != nil {
		return nil, err
	}

	players := make([]*cs.Player, 0, len(settings.Players))

	for i, player := range settings.Players {
		if player.Type == cs.NewGamePlayerTypeHost {

			// add the host player with the host race
			race, err := gr.db.GetRace(player.RaceID)
			if err != nil {
				return nil, err
			}
			if race.UserID != hostID {
				return nil, fmt.Errorf("user %d does not own Race %d", hostID, race.ID)
			}
			log.Debug().Int64("hostID", hostID).Msgf("Adding host to game")
			player := client.NewPlayer(hostID, *race, &game.Rules)
			player.GameID = game.ID
			player.Num = i + 1
			players = append(players, player)
		} else if player.Type == cs.NewGamePlayerTypeAI {
			// g.AddPlayer(game.NewAIPlayer())
		} else if player.Type == cs.NewGamePlayerTypeOpen {
			game.OpenPlayerSlots++
			log.Debug().Uint("openPlayerSlots", game.OpenPlayerSlots).Msgf("Added open player slot")
		}
	}

	fullGame := &cs.FullGame{
		Game:     &game,
		Players:  players,
		Universe: &cs.Universe{}, // todo: populate
	}

	gr.db.UpdateFullGame(fullGame)

	// generate the universe if this game is done
	if game.OpenPlayerSlots == 0 {
		fullGame.Universe, err = gr.GenerateUniverse(game.ID)
		if err != nil {
			return nil, err
		}
	}

	return fullGame, nil
}

// add a player to an existing game
func (gr *gameRunner) AddPlayer(gameID int64, userID int64, race *cs.Race) error {

	fullGame, err := gr.LoadGame(gameID)
	if err != nil {
		return fmt.Errorf("unable to load game %d: %w", gameID, err)
	}

	client := cs.NewClient()
	player := client.NewPlayer(userID, *race, &fullGame.Rules)
	player.GameID = fullGame.ID
	if err := gr.db.CreatePlayer(player); err != nil {
		return fmt.Errorf("save player %s for game %d: %w", player, gameID, err)
	}
	if err := gr.db.UpdateFullGame(fullGame); err != nil {
		return fmt.Errorf("save game %d: %w", gameID, err)
	}

	// generate the universe if this game is ready
	if fullGame.OpenPlayerSlots == 0 {
		if _, err := gr.GenerateUniverse(fullGame.ID); err != nil {
			return fmt.Errorf("generate universe: %d: %w", gameID, err)
		}
	}

	return nil
}

// process an ai player's turn
func (gr *gameRunner) processAITurns(universe *cs.Universe, player *cs.Player) {
	pmo := universe.GetPlayerMapObjects(player.Num)
	ai := ai.NewAIPlayer(player, pmo)
	ai.ProcessTurn()
}

func (gr *gameRunner) GenerateUniverse(gameID int64) (*cs.Universe, error) {
	defer timeTrack(time.Now(), "GenerateUniverse")
	fullGame, err := gr.LoadGame(gameID)
	if err != nil {
		return nil, fmt.Errorf("load game %d: %w", gameID, err)
	}

	client := cs.NewClient()
	universe, err := client.GenerateUniverse(fullGame.Game, fullGame.Players)
	if err != nil {
		return nil, fmt.Errorf("generate universe for game %d: %w", gameID, err)
	}

	// process ai players
	for _, player := range fullGame.Players {
		// TODO: ai only ai processing
		// if player.AIControlled {
		gr.processAITurns(universe, player)
		// }
	}

	fullGame.State = cs.GameStateWaitingForPlayers
	fullGame.Universe = universe

	err = gr.db.UpdateFullGame(fullGame)
	if err != nil {
		return nil, fmt.Errorf("save game %d: %w", gameID, err)
	}

	return fullGame.Universe, nil
}

// load a full game
func (gr *gameRunner) LoadGame(gameID int64) (*cs.FullGame, error) {
	defer timeTrack(time.Now(), "LoadGame")

	fullGame, err := gr.db.GetFullGame(gameID)

	if err != nil {
		return nil, fmt.Errorf("load game %d: %w", gameID, err)
	}

	return fullGame, nil
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

// check a game for every player submitted and generate a turn
func (gr *gameRunner) CheckAndGenerateTurn(gameID int64) (TurnGenerationCheckResult, error) {
	defer timeTrack(time.Now(), "CheckAndGenerateTurn")
	fullGame, err := gr.LoadGame(gameID)

	if err != nil {
		return TurnNotGenerated, err
	}

	// if everyone submitted their turn, generate a new turn
	client := cs.NewClient()
	if client.CheckAllPlayersSubmitted(fullGame.Players) {
		if err := client.GenerateTurn(fullGame.Game, fullGame.Universe, fullGame.Players); err != nil {
			return TurnNotGenerated, fmt.Errorf("generate turn -> %w", err)
		}

		// process ai players
		for _, player := range fullGame.Players {
			// TODO: ai only ai processing
			// if player.AIControlled {
			gr.processAITurns(fullGame.Universe, player)
			// }
		}

		if err := gr.db.UpdateFullGame(fullGame); err != nil {
			return TurnNotGenerated, fmt.Errorf("save game after turn generation -> %w", err)
		}
		return TurnGenerated, nil
	}

	return TurnNotGenerated, nil
}
