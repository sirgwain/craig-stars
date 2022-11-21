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
	cfg.Database.Filename = "../data/sqlx.db"
	// cfg.Database.Filename = ":memory:"
	cfg.Database.Recreate = true
	cfg.Database.DebugLogging = true
	if err := c.Connect(cfg); err != nil {
		panic(fmt.Errorf("Failed to connect to test database, %w", err))
	}
	c.ExecSchema("../schema.sql")

	// create a test user
	if err := c.CreateUser(game.NewUser("admin", "admin", game.RoleAdmin)); err != nil {
		panic(fmt.Errorf("Failed to create test database user, %w", err))
	}

	return &c
}

func (c *client) createTestGame() *game.Game {

	g := game.NewGame()
	g.HostID = 1
	if err := c.CreateGame(g); err != nil {
		panic(fmt.Errorf("Failed to create test database game, %w", err))
	}

	return g
}
