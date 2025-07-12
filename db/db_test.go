package db

import (
	"context"
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
		if _, err := c.CreateUser(context.Background(), user); err != nil {
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
func (c *client) createTestGame(ctx context.Context) *cs.Game {

	gameClient := cs.NewGamer()
	game, err := c.CreateGame(ctx, gameClient.CreateGame(1, *cs.NewGameSettings()))
	if err != nil {
		panic(fmt.Errorf("error creating test database game: \n%w", err))
	}

	return game
}

// create a simple game with one player
func (c *client) createTestGameWithPlayer(ctx context.Context) (*cs.Game, *cs.Player) {

	gameClient := cs.NewGamer()
	game, err := c.CreateGame(ctx, gameClient.CreateGame(1, *cs.NewGameSettings()))
	if err != nil {
		panic(fmt.Errorf("error creating test database game: \n%w", err))
	}

	player := gameClient.NewPlayer(1, cs.Humanoids(), &game.Rules)
	player.Num = 1
	player.GameID = game.ID

	player, err = c.CreatePlayer(ctx, player)
	if err != nil {
		panic(fmt.Errorf("error creating test database game player: \n%w", err))
	}

	return game, player
}

func (c *client) createTestShipDesign(ctx context.Context, player *cs.Player, design *cs.ShipDesign) *cs.ShipDesign {
	design.PlayerNum = player.Num
	design.GameID = player.GameID
	var err error
	created, err := c.CreateShipDesign(ctx, design)
	if err != nil {
		panic(fmt.Errorf("error creating test design: \n%w", err))
	}
	return created
}

func (c *client) createTestFullGame(ctx context.Context) *cs.FullGame {
	gameClient := cs.NewGamer()
	g, player := c.createTestGameWithPlayer(ctx)

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
