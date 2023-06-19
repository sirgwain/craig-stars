package db

import (
	"fmt"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
)

func connectTestDB() *client {

	c := client{
		converter: &GameConverter{},
	}
	cfg := &config.Config{}
	// cfg.Database.Filename = "../data/sqlx.db"
	cfg.Database.Filename = ":memory:"
	cfg.Database.UsersFilename = ":memory:"
	cfg.Database.Recreate = true
	cfg.Database.DebugLogging = true
	if err := c.Connect(cfg); err != nil {
		panic(fmt.Errorf("connect to test database, %w", err))
	}

	// create a test user
	user, err := cs.NewUser("admin", "admin", "admin@craig-stars.net", cs.RoleAdmin)
	if err != nil {
		panic(fmt.Errorf("generate test user, %w", err))
	}
	if err := c.CreateUser(user); err != nil {
		panic(fmt.Errorf("create test database user, %w", err))
	}

	return &c
}

// create a new game
func (c *client) createTestGame() *cs.Game {

	game := cs.NewGame()
	game.HostID = 1
	if err := c.CreateGame(game); err != nil {
		panic(fmt.Errorf("create test database game, %w", err))
	}

	return game
}

// create a simple game with one player
func (c *client) createTestGameWithPlayer() (*cs.Game, *cs.Player) {

	gameClient := cs.NewGamer()
	game := gameClient.CreateGame(1, *cs.NewGameSettings())
	if err := c.CreateGame(game); err != nil {
		panic(fmt.Errorf("create test database game, %w", err))
	}

	player := gameClient.NewPlayer(1, cs.Humanoids(), &game.Rules)
	player.Num = 1
	player.GameID = game.ID

	if err := c.CreatePlayer(player); err != nil {
		panic(fmt.Errorf("create test database game player %w", err))
	}

	return game, player
}

func (c *client) createTestShipDesign(player *cs.Player, design *cs.ShipDesign) {
	design.PlayerNum = player.Num
	design.GameID = player.GameID
	if err := c.CreateShipDesign(design); err != nil {
		panic(fmt.Errorf("create test design %w", err))
	}
}

func (c *client) createTestFullGame() *cs.FullGame {
	gameClient := cs.NewGamer()
	g, player := c.createTestGameWithPlayer()

	players := []*cs.Player{player}
	universe, err := gameClient.GenerateUniverse(g, players)
	if err != nil {
		panic(err)
	}

	fg := cs.FullGame{
		Game:     g,
		Players:  players,
		Universe: universe,
	}

	return &fg
}
