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
		panic(fmt.Errorf("Failed to connect to test database, %w", err))
	}
	c.ExecSchema("../schema.sql")

	// create a test user
	if err := c.CreateUser(game.NewUser("admin", "admin", game.RoleAdmin)); err != nil {
		panic(fmt.Errorf("Failed to create test database user, %w", err))
	}

	return &c
}

// create a new game
func (c *client) createTestGame() *game.Game {

	g := game.NewGame()
	g.HostID = 1
	if err := c.CreateGame(g); err != nil {
		panic(fmt.Errorf("Failed to create test database game, %w", err))
	}

	return g
}

// create a simple game with one player
func (c *client) createTestGameWithPlayer() (*game.Game, *game.Player) {

	g := game.NewGame()
	g.HostID = 1
	if err := c.CreateGame(g); err != nil {
		panic(fmt.Errorf("Failed to create test database game, %w", err))
	}

	player := game.NewPlayer(1, game.NewRace().WithSpec(&g.Rules))
	player.GameID = g.ID

	if err := c.CreatePlayer(player); err != nil {
		panic(fmt.Errorf("Failed to create test database game player, %w", err))
	}

	// add the player's race into the db
	player.Race.PlayerID = &player.ID
	player.Race.UserID = player.UserID
	if err := c.CreateRace(&player.Race); err != nil {
		panic(fmt.Errorf("Failed to create test database game race, %w", err))
	}

	return g, player
}
