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

type GameRunner struct {
	db db.Client
}

func NewGameRunner(db db.Client) *GameRunner {
	return &GameRunner{db}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// host a new game
func (gr *GameRunner) HostGame(hostID uint64, settings *game.GameSettings) (*game.FullGame, error) {
	client := game.NewClient()
	g := client.CreateGame(hostID, *settings)

	// create the game in the database
	err := gr.db.CreateGame(&g)
	if err != nil {
		return nil, err
	}

	players := make([]*game.Player, 0, len(settings.Players))

	for _, player := range settings.Players {
		if player.Type == game.NewGamePlayerTypeHost {

			// add the host player with the host race
			race, err := gr.db.FindRaceById(player.RaceID)
			if err != nil {
				return nil, err
			}
			if race.UserID != hostID {
				return nil, fmt.Errorf("user %d does not own Race %d", hostID, race.ID)
			}
			log.Debug().Uint64("hostID", hostID).Msgf("Adding host to game")
			player := client.NewPlayer(hostID, *race, &g.Rules)
			player.GameID = g.ID
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

	gr.db.SaveGame(fg)

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
func (gr *GameRunner) AddPlayer(gameID uint64, userID uint64, race *game.Race) error {

	g, err := gr.db.FindGameById(gameID)
	if err != nil {
		return fmt.Errorf("unable to load game %d: %w", gameID, err)
	}

	client := game.NewClient()
	player := client.NewPlayer(userID, *race, &g.Rules)
	player.GameID = g.ID
	if err := gr.db.SavePlayer(player); err != nil {
		return fmt.Errorf("failed to save player %s for game %d: %w", player, gameID, err)
	}
	if err := gr.db.SaveGame(g); err != nil {
		return fmt.Errorf("failed to save game %d: %w", gameID, err)
	}

	// generate the universe if this game is ready
	if g.OpenPlayerSlots == 0 {
		if _, err := gr.GenerateUniverse(g.ID); err != nil {
			return fmt.Errorf("failed to generate universe: %d: %w", gameID, err)
		}
	}

	return nil
}

func (gr *GameRunner) GenerateUniverse(gameID uint64) (*game.Universe, error) {
	defer timeTrack(time.Now(), "GenerateUniverse")
	g, err := gr.LoadGame(gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to load game %d: %w", gameID, err)
	}

	client := game.NewClient()
	universe, err := client.GenerateUniverse(g.Game, g.Players)
	if err != nil {
		return nil, fmt.Errorf("failed to generate universe for game %d: %w", gameID, err)
	}

	g.State = game.GameStateWaitingForPlayers
	g.Universe = universe

	err = gr.db.SaveGame(g)
	if err != nil {
		return nil, fmt.Errorf("failed to save game %d: %w", gameID, err)
	}

	return g.Universe, nil
}

// load a full game
func (gr *GameRunner) LoadGame(gameID uint64) (*game.FullGame, error) {
	defer timeTrack(time.Now(), "LoadGame")

	g, err := gr.db.FindGameById(gameID)

	if err != nil {
		return nil, fmt.Errorf("failed to load game %d: %w", gameID, err)
	}

	return g, nil
}

// load a player and the light version of the player game
func (gr *GameRunner) LoadPlayerGame(gameID uint64, userID uint64) (*game.Game, *game.FullPlayer, error) {

	g, err := gr.db.FindGameByIdLight(gameID)

	if err != nil {
		return nil, nil, err
	}

	if g.Rules.TechsID == 0 {
		g.Rules.WithTechStore(&game.StaticTechStore)
	} else {
		techs, err := gr.db.FindTechStoreById(g.Rules.TechsID)
		if err != nil {
			return nil, nil, err
		}
		g.Rules.WithTechStore(techs)
	}

	player, err := gr.db.FindPlayerByGameId(gameID, userID)

	if err != nil {
		return nil, nil, err
	}

	return g, player, nil
}

// submit a turn for a player
func (gr *GameRunner) SubmitTurn(gameID uint64, userID uint64) error {
	player, err := gr.db.FindPlayerByGameIdLight(gameID, userID)
	if err != nil {
		return fmt.Errorf("failed to find player for user %d, game %d: %w", userID, gameID, err)
	}

	player.SubmittedTurn = true
	gr.db.SavePlayer(player)
	return nil
}

// check a game for every player submitted and generate a turn
func (gr *GameRunner) CheckAndGenerateTurn(gameID uint64) (TurnGenerationCheckResult, error) {
	defer timeTrack(time.Now(), "CheckAndGenerateTurn")
	g, err := gr.LoadGame(gameID)

	if err != nil {
		return TurnNotGenerated, err
	}

	// if everyone submitted their turn, generate a new turn
	client := game.NewClient()
	if client.CheckAllPlayersSubmitted(g.Players) {
		client.GenerateTurn(g.Game, g.Universe, g.Players)
		gr.db.SaveGame(g)
		return TurnGenerated, nil
	}

	return TurnNotGenerated, nil
}
