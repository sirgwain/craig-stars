package server

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/game"
)

type TurnGenerationCheckResult uint

const (
	TurnNotGenerated TurnGenerationCheckResult = iota
	TurnGenerated
)

type DBClient db.Client

type GameRunner interface {
	HostGame(hostID int64, settings *game.GameSettings) (*game.FullGame, error)
	AddPlayer(gameID int64, userID int64, race *game.Race) error
	GenerateUniverse(gameID int64) (*game.Universe, error)
	LoadGame(gameID int64) (*game.FullGame, error)
	LoadPlayerGame(gameID int64, userID int64) (*game.Game, *game.FullPlayer, error)
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
func (gr *gameRunner) HostGame(hostID int64, settings *game.GameSettings) (*game.FullGame, error) {
	client := game.NewClient()
	g := client.CreateGame(hostID, *settings)

	// create the game in the database
	err := gr.db.CreateGame(&g)
	if err != nil {
		return nil, err
	}

	players := make([]*game.Player, 0, len(settings.Players))

	for i, player := range settings.Players {
		if player.Type == game.NewGamePlayerTypeHost {

			// add the host player with the host race
			race, err := gr.db.GetRace(player.RaceID)
			if err != nil {
				return nil, err
			}
			if race.UserID != hostID {
				return nil, fmt.Errorf("user %d does not own Race %d", hostID, race.ID)
			}
			log.Debug().Int64("hostID", hostID).Msgf("Adding host to game")
			player := client.NewPlayer(hostID, *race, &g.Rules)
			player.GameID = g.ID
			player.Num = i + 1
			players = append(players, player)
		} else if player.Type == game.NewGamePlayerTypeAI {
			// g.AddPlayer(game.NewAIPlayer())
		} else if player.Type == game.NewGamePlayerTypeOpen {
			g.OpenPlayerSlots++
			log.Debug().Uint("openPlayerSlots", g.OpenPlayerSlots).Msgf("Added open player slot")
		}
	}

	fg := &game.FullGame{
		Game:     &g,
		Players:  players,
		Universe: &game.Universe{}, // todo: populate
	}

	gr.db.UpdateFullGame(fg)

	// generate the universe if this game is done
	if g.OpenPlayerSlots == 0 {
		fg.Universe, err = gr.GenerateUniverse(g.ID)
		if err != nil {
			return nil, err
		}
	}

	return fg, nil
}

// add a player to an existing game
func (gr *gameRunner) AddPlayer(gameID int64, userID int64, race *game.Race) error {

	g, err := gr.db.GetFullGame(gameID)
	if err != nil {
		return fmt.Errorf("unable to load game %d: %w", gameID, err)
	}

	client := game.NewClient()
	player := client.NewPlayer(userID, *race, &g.Rules)
	player.GameID = g.ID
	if err := gr.db.CreatePlayer(player); err != nil {
		return fmt.Errorf("save player %s for game %d: %w", player, gameID, err)
	}
	if err := gr.db.UpdateFullGame(g); err != nil {
		return fmt.Errorf("save game %d: %w", gameID, err)
	}

	// generate the universe if this game is ready
	if g.OpenPlayerSlots == 0 {
		if _, err := gr.GenerateUniverse(g.ID); err != nil {
			return fmt.Errorf("generate universe: %d: %w", gameID, err)
		}
	}

	return nil
}

func (gr *gameRunner) GenerateUniverse(gameID int64) (*game.Universe, error) {
	defer timeTrack(time.Now(), "GenerateUniverse")
	g, err := gr.LoadGame(gameID)
	if err != nil {
		return nil, fmt.Errorf("load game %d: %w", gameID, err)
	}

	client := game.NewClient()
	universe, err := client.GenerateUniverse(g.Game, g.Players)
	if err != nil {
		return nil, fmt.Errorf("generate universe for game %d: %w", gameID, err)
	}

	g.State = game.GameStateWaitingForPlayers
	g.Universe = universe

	err = gr.db.UpdateFullGame(g)
	if err != nil {
		return nil, fmt.Errorf("save game %d: %w", gameID, err)
	}

	return g.Universe, nil
}

// load a full game
func (gr *gameRunner) LoadGame(gameID int64) (*game.FullGame, error) {
	defer timeTrack(time.Now(), "LoadGame")

	g, err := gr.db.GetFullGame(gameID)

	if err != nil {
		return nil, fmt.Errorf("load game %d: %w", gameID, err)
	}

	return g, nil
}

// load a player and the light version of the player game
func (gr *gameRunner) LoadPlayerGame(gameID int64, userID int64) (*game.Game, *game.FullPlayer, error) {

	g, err := gr.db.GetGame(gameID)

	if err != nil {
		return nil, nil, err
	}

	if g.Rules.TechsID == 0 {
		g.Rules.WithTechStore(&game.StaticTechStore)
	} else {
		techs, err := gr.db.GetTechStore(g.Rules.TechsID)
		if err != nil {
			return nil, nil, err
		}
		g.Rules.WithTechStore(techs)
	}

	player, err := gr.db.GetFullPlayerForGame(gameID, userID)

	if err != nil {
		return nil, nil, err
	}

	return g, player, nil
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
	g, err := gr.LoadGame(gameID)

	if err != nil {
		return TurnNotGenerated, err
	}

	// if everyone submitted their turn, generate a new turn
	client := game.NewClient()
	if client.CheckAllPlayersSubmitted(g.Players) {
		if err := client.GenerateTurn(g.Game, g.Universe, g.Players); err != nil {
			return TurnNotGenerated, fmt.Errorf("generate turn -> %w", err)
		}
		if err := gr.db.UpdateFullGame(g); err != nil {
			return TurnNotGenerated, fmt.Errorf("save game after turn generation -> %w", err)
		}
		return TurnGenerated, nil
	}

	return TurnNotGenerated, nil
}
