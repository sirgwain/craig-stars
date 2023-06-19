package server

import (
	"fmt"

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

// host a new game
func (gr *GameRunner) HostGame(hostID uint, settings *game.GameSettings) (*game.FullGame, error) {
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
			log.Debug().Uint("hostID", hostID).Msgf("Adding host to game")
			players = append(players, client.NewPlayer(hostID, *race, &g.Rules))
		} else if player.Type == game.NewGamePlayerTypeAI {
			// g.AddPlayer(game.NewAIPlayer())
		} else if player.Type == game.NewGamePlayerTypeOpen {
			g.OpenPlayerSlots++
			log.Debug().Uint("openPlayerSlots", g.OpenPlayerSlots).Msgf("Added open player slot")
		}
	}

	// generate the universe if this game is done
	if g.OpenPlayerSlots == 0 {
		gr.GenerateUniverse(g.ID)
	}

	return &game.FullGame{
		Game:     &g,
		Players:  players,
		Universe: &game.Universe{}, // todo: populate
	}, nil
}

// add a player to an existing game
func (gr *GameRunner) AddPlayer(gameID uint, userID uint, race *game.Race) error {

	g, err := gr.db.FindGameById(gameID)
	if err != nil {
		return err
	}

	client := game.NewClient()
	player := client.NewPlayer(userID, *race, &g.Rules)
	player.GameID = g.ID
	if err := gr.db.SavePlayer(player); err != nil {
		return err
	}
	if err := gr.db.SaveGame(g); err != nil {
		return err
	}

	// generate the universe if this game is ready
	if g.OpenPlayerSlots == 0 {
		if err := gr.GenerateUniverse(g.ID); err != nil {
			return err
		}
	}

	return nil
}

func (gr *GameRunner) GenerateUniverse(gameID uint) error {
	g, err := gr.LoadGame(gameID)
	if err != nil {
		return err
	}

	client := game.NewClient()
	universe, err := client.GenerateUniverse(g.Game, g.Players)
	if err != nil {
		return err
	}

	g.State = game.GameStateWaitingForPlayers
	g.Universe = universe

	err = gr.db.SaveGame(g)
	if err != nil {
		return err
	}

	return nil
}

// load a full game
func (gr *GameRunner) LoadGame(gameID uint) (*game.FullGame, error) {
	g, err := gr.db.FindGameById(gameID)

	if err != nil {
		return nil, err
	}

	return g, nil
}

// load a player and the light version of the player game
func (gr *GameRunner) LoadPlayerGame(gameID uint, userID uint) (*game.Game, *game.FullPlayer, error) {

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
func (gr *GameRunner) SubmitTurn(gameID uint, userID uint) error {
	player, err := gr.db.FindPlayerByGameIdLight(gameID, userID)
	if err != nil {
		return err
	}

	player.SubmittedTurn = true
	gr.db.SavePlayer(player)
	return nil
}

// check a game for every player submitted and generate a turn
func (gr *GameRunner) CheckAndGenerateTurn(id uint) (TurnGenerationCheckResult, error) {
	g, err := gr.LoadGame(id)

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
