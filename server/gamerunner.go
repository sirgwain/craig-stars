package server

import (
	"fmt"

	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/game"
)

type TurnGenerationCheckResult uint

const (
	TurnNotGenerated TurnGenerationCheckResult = iota
	TurnGenerated
)

type GameRunner struct {
	db db.Service
}

func NewGameRunner(db db.Service) *GameRunner {
	return &GameRunner{db}
}

// host a new game
func (gr *GameRunner) HostGame(hostID uint, settings *game.GameSettings) (*game.Game, error) {
	g := game.NewGame().WithSettings(settings)
	g.HostID = hostID

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
			g.AddPlayer(game.NewPlayer(hostID, race))
		}
	}

	err := gr.db.CreateGame(g)
	if err != nil {
		return nil, err
	}

	if err := gr.GenerateUniverse(g.ID); err != nil {
		return nil, err
	}

	return g, nil
}

func (gr *GameRunner) AddPlayer(gameID uint, player *game.Player) error {
	game, err := gr.LoadGame(gameID)
	if err != nil {
		return err
	}

	game.AddPlayer(player)
	gr.db.SaveGame(game)

	return nil
}

func (gr *GameRunner) GenerateUniverse(gameID uint) error {
	game, err := gr.LoadGame(gameID)
	if err != nil {
		return err
	}

	err = game.GenerateUniverse()
	if err != nil {
		return err
	}

	err = gr.db.SaveGame(game)
	if err != nil {
		return err
	}

	return nil
}

// load a full game
func (gr *GameRunner) LoadGame(gameID uint) (*game.Game, error) {
	g, err := gr.db.FindGameById(gameID)

	if err != nil {
		return nil, err
	}

	return g, nil
}

// load a player and the light version of the player game
func (gr *GameRunner) LoadPlayerGame(gameID uint, userID uint) (*game.Game, *game.Player, error) {

	g, err := gr.db.FindGameByIdLight(gameID)

	if err != nil {
		return nil, nil, err
	}

	if g.Rules.TechsID == 0 {
		g.Rules.Techs = &game.StaticTechStore
	} else {
		techs, err := gr.db.FindTechStoreById(g.Rules.TechsID)
		if err != nil {
			return nil, nil, err
		}
		g.Rules.Techs = techs
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
	game, err := gr.LoadGame(id)

	if err != nil {
		return TurnNotGenerated, err
	}

	// if everyone submitted their turn, generate a new turn
	if game.CheckAllPlayersSubmitted() {
		game.GenerateTurn()
		gr.db.SaveGame(game)
		return TurnGenerated, nil
	}

	return TurnNotGenerated, nil
}
