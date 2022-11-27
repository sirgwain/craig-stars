package dbsqlx

import (
	"fmt"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/game"
)

func connectTestDB() *client {

	c := client{
		converter: &GameConverter{},
	}
	cfg := &config.Config{}
	// cfg.Database.Filename = "../data/sqlx.db"
	cfg.Database.Filename = ":memory:"
	cfg.Database.Recreate = true
	cfg.Database.DebugLogging = true
	if err := c.Connect(cfg); err != nil {
		panic(fmt.Errorf("connect to test database, %w", err))
	}
	c.ExecSchema("../schema.sql")

	// create a test user
	if err := c.CreateUser(game.NewUser("admin", "admin", game.RoleAdmin)); err != nil {
		panic(fmt.Errorf("create test database user, %w", err))
	}

	return &c
}

// create a new game
func (c *client) createTestGame() *game.Game {

	g := game.NewGame()
	g.HostID = 1
	if err := c.CreateGame(g); err != nil {
		panic(fmt.Errorf("create test database game, %w", err))
	}

	return g
}

// create a simple game with one player
func (c *client) createTestGameWithPlayer() (*game.Game, *game.Player) {

	g := game.NewGame()
	g.HostID = 1
	if err := c.CreateGame(g); err != nil {
		panic(fmt.Errorf("create test database game, %w", err))
	}

	player := game.NewPlayer(1, game.NewRace().WithSpec(&g.Rules))
	player.GameID = g.ID

	if err := c.CreatePlayer(player); err != nil {
		panic(fmt.Errorf("create test database game player %w", err))
	}

	return g, player
}

func (c *client) createTestShipDesign(player *game.Player, design *game.ShipDesign) {
	design.PlayerID = player.ID
	if err := c.CreateShipDesign(design); err != nil {
		panic(fmt.Errorf("create test design %w", err))
	}
}

func (c *client) createTestFullGame() *game.FullGame {
	gameClient := game.NewClient()
	g, player := c.createTestGameWithPlayer()

	players := []*game.Player{player}
	universe, err := gameClient.GenerateUniverse(g, players)
	if err != nil {
		panic(err)
	}

	fg := game.FullGame{
		Game:     g,
		Players:  players,
		Universe: universe,
	}

	return &fg
}
