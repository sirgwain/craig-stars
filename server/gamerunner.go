package server

import (
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

func (gr *GameRunner) NewGame(hostID uint) *game.Game {
	game := game.NewGame()
	game.HostID = hostID

	gr.db.CreateGame(game)

	return game
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
	game, err := gr.db.FindGameById(gameID)

	if err != nil {
		return nil, err
	}

	return game, nil
}

// load a player and the light version of the player game
func (gr *GameRunner) LoadPlayerGame(gameID uint, userID uint) (*game.Game, *game.Player, error) {

	g, err := gr.db.FindGameByIdLight(gameID)

	if err != nil {
		return nil, nil, err
	}

	if g.Rules.TechsID != 0 {
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
