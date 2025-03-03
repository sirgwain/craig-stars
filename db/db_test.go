package db

import (
	"fmt"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
)

func connectTestDB() *client {

	dbConn := dbConn{}
	cfg := &config.Config{}
	// cfg.Database.Filename = "../data/sqlx.db"
	cfg.Database.Filename = ":memory:"
	cfg.Database.DebugLogging = true
	cfg.Database.SkipUpgrade = true
	if err := dbConn.Connect(cfg); err != nil {
		panic(fmt.Errorf("error while connecting to test database: \n%w", err))
	}

	// create a test user
	user, err := cs.NewUser("admin", "admin", "admin@craig-stars.net", cs.RoleAdmin)
	if err != nil {
		panic(fmt.Errorf("error generating test user: \n%w", err))
	}

	if err := dbConn.WrapInTransaction(func(c Client) error {
		if err := c.CreateUser(user); err != nil {
			return fmt.Errorf("error creating test database user: \n%w", err)
		}
		return nil
	}); err != nil {
		panic(fmt.Errorf("error creating test user in db: \n%w", err))
	}

	// create a new c from a transaction
	c, err := dbConn.BeginTransaction()
	if err != nil {
		panic(fmt.Errorf("error beginning test transaction: \n%w", err))
	}

	return c.(*client)
}

func closeTestDB(c *client) {
	if err := c.commit(); err != nil {
		panic(fmt.Errorf("error commiting test transaction: \n%w", err))
	}
}

// create a new game
func (c *client) createTestGame() *cs.Game {

	game := cs.NewGame()
	game.HostID = 1
	if err := c.CreateGame(game); err != nil {
		panic(fmt.Errorf("error creating test database game: \n%w", err))
	}

	return game
}

// create a simple game with one player
func (c *client) createTestGameWithPlayer() (*cs.Game, *cs.Player) {

	gameClient := cs.NewGamer()
	game := gameClient.CreateGame(1, *cs.NewGameSettings())
	if err := c.CreateGame(game); err != nil {
		panic(fmt.Errorf("error creating test database game: \n%w", err))
	}

	player := gameClient.NewPlayer(1, cs.Humanoids(), &game.Rules)
	player.Num = 1
	player.GameID = game.ID

	if err := c.CreatePlayer(player); err != nil {
		panic(fmt.Errorf("error creating test database game player: \n%w", err))
	}

	return game, player
}

func (c *client) createTestShipDesign(player *cs.Player, design *cs.ShipDesign) {
	design.PlayerNum = player.Num
	design.GameID = player.GameID
	if err := c.CreateShipDesign(design); err != nil {
		panic(fmt.Errorf("error creating test design: \n%w", err))
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
