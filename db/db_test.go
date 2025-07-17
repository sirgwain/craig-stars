package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sirgwain/craig-stars/ai"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
)

func connectTestDB() *client {

	dbConn := dbConn{}
	cfg := &config.Config{}
	// cfg.Database.Filename = "../data/sqlx.db"
	// cfg.Database.DebugLogging = true
	cfg.Database.Filename = ":memory:"
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

	return dbConn.NewReadWriteClient().(*client)
}

// create a new game
func (c *client) createTestGame(ctx context.Context) *cs.Game {

	gameClient := cs.NewGamer()
	game := gameClient.CreateGame(1, *cs.NewGameSettings())
	err := c.SaveGame(ctx, game)
	if err != nil {
		panic(fmt.Errorf("error creating test database game: \n%w", err))
	}

	return game
}

// create a simple game with one player
func (c *client) createTestGameWithPlayer(ctx context.Context) (*cs.Game, *cs.Player) {

	gameClient := cs.NewGamer()
	game := gameClient.CreateGame(1, *cs.NewGameSettings())
	err := c.SaveGame(ctx, game)
	if err != nil {
		panic(fmt.Errorf("error creating test database game: \n%w", err))
	}

	player := gameClient.NewPlayer(1, cs.Humanoids(), &game.Rules)
	player.Num = 1
	player.GameID = game.ID

	if err := c.SavePlayer(ctx, player); err != nil {
		panic(fmt.Errorf("error creating test database game player: \n%w", err))
	}

	return game, player
}

func (c *client) createTestShipDesign(ctx context.Context, player *cs.Player, design *cs.ShipDesign) *cs.ShipDesign {
	design.PlayerNum = player.Num
	design.GameID = player.GameID

	if err := c.SaveShipDesign(ctx, design); err != nil {
		panic(fmt.Errorf("error creating test design: \n%w", err))
	}
	return design
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

func BenchmarkUpdateFullGame(b *testing.B) {
	dbConn := dbConn{}
	cfg := &config.Config{}
	// cfg.Database.Filename = "../data/sqlx.db"
	cfg.Database.Filename = ":memory:"
	cfg.Database.DebugLogging = false
	cfg.Database.SkipUpgrade = true

	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	b.Run("Small Game one turn", func(b *testing.B) {
		var err error
		if err = dbConn.Connect(cfg); err != nil {
			b.Fatalf("error while connecting to test database %v", err)
		}
		defer dbConn.Close()

		c := dbConn.NewReadWriteClient()

		gameClient := cs.NewGamer()

		b.ResetTimer()
		for b.Loop() {
			b.StopTimer()

			game := cs.NewGame()
			player := gameClient.NewPlayer(0, cs.Humanoids(), &game.Rules).WithNum(1).WithAIControlled(true)
			players := []*cs.Player{player}
			universe, err := gameClient.GenerateUniverse(game, players)
			if err != nil {
				b.Fatalf("failed to generate universe")
			}

			fullGame := &cs.FullGame{
				Game:      game,
				Universe:  universe,
				TechStore: &cs.StaticTechStore,
				Players:   players,
			}

			err = c.SaveGame(b.Context(), fullGame.Game)
			if err != nil {
				b.Fatalf("failed to create game %v", err)
			}

			gameClient.GenerateTurn(fullGame.Game, fullGame.Universe, fullGame.Players)

			b.StartTimer()
			c.UpdateFullGame(b.Context(), fullGame)
		}
	})

	b.Run("Large Game many turns and players", func(b *testing.B) {
		var err error
		if err = dbConn.Connect(cfg); err != nil {
			b.Fatalf("error while connecting to test database %v", err)
		}
		defer dbConn.Close()

		c := dbConn.NewReadWriteClient()
		gameClient := cs.NewGamer()

		settings := cs.NewGameSettings().WithSize(cs.SizeLarge).WithDensity(cs.DensityPacked)
		game := cs.NewGame().WithSettings(*settings)
		players := make([]*cs.Player, len(ai.Races))
		for i := range ai.Races {
			players[i] = gameClient.NewPlayer(0, ai.Races[i], &game.Rules).WithNum(i + 1).WithAIControlled(true)
		}
		universe, err := gameClient.GenerateUniverse(game, players)
		if err != nil {
			b.Fatalf("failed to generate universe")
		}

		fullGame := &cs.FullGame{
			Game:      game,
			Universe:  universe,
			TechStore: &cs.StaticTechStore,
			Players:   players,
		}

		err = c.SaveGame(b.Context(), fullGame.Game)
		if err != nil {
			b.Fatalf("failed to create game %v", err)
		}

		// generate some turns
		for range 50 {
			gameClient.GenerateTurn(fullGame.Game, fullGame.Universe, fullGame.Players)
		}

		b.ResetTimer()
		for b.Loop() {
			if err := c.UpdateFullGame(b.Context(), fullGame); err != nil {
				b.Fatalf("failed to update game %v", err)
			}
		}
	})
}

func BenchmarkGetFullGame(b *testing.B) {
	dbConn := dbConn{}
	cfg := &config.Config{}
	// cfg.Database.Filename = "../data/sqlx.db"
	cfg.Database.Filename = ":memory:"
	cfg.Database.DebugLogging = false
	cfg.Database.SkipUpgrade = true

	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	b.Run("Small Game one turn", func(b *testing.B) {
		var err error
		if err = dbConn.Connect(cfg); err != nil {
			b.Fatalf("error while connecting to test database %v", err)
		}
		defer dbConn.Close()

		c := dbConn.NewReadWriteClient()

		gameClient := cs.NewGamer()
		game := cs.NewGame()
		player := gameClient.NewPlayer(0, cs.Humanoids(), &game.Rules).WithNum(1).WithAIControlled(true)
		players := []*cs.Player{player}
		universe, err := gameClient.GenerateUniverse(game, players)
		if err != nil {
			b.Fatalf("failed to generate universe")
		}

		fullGame := &cs.FullGame{
			Game:      game,
			Universe:  universe,
			TechStore: &cs.StaticTechStore,
			Players:   players,
		}

		err = c.SaveGame(b.Context(), fullGame.Game)
		if err != nil {
			b.Fatalf("failed to create game %v", err)
		}

		gameClient.GenerateTurn(fullGame.Game, fullGame.Universe, fullGame.Players)
		if err := c.UpdateFullGame(b.Context(), fullGame); err != nil {
			b.Fatalf("failed to update full game %v", err)
		}

		b.ResetTimer()
		for b.Loop() {

			b.StartTimer()
			if _, err := c.GetFullGame(b.Context(), game.ID); err != nil {
				b.Fatalf("failed to update game %v", err)
			}
		}
	})

	b.Run("Large Game many turns and players", func(b *testing.B) {
		var err error
		if err = dbConn.Connect(cfg); err != nil {
			b.Fatalf("error while connecting to test database %v", err)
		}
		defer dbConn.Close()

		c := dbConn.NewReadWriteClient()
		gameClient := cs.NewGamer()

		settings := cs.NewGameSettings().WithSize(cs.SizeLarge).WithDensity(cs.DensityPacked)
		game := cs.NewGame().WithSettings(*settings)
		players := make([]*cs.Player, len(ai.Races))
		for i := range ai.Races {
			players[i] = gameClient.NewPlayer(0, ai.Races[i], &game.Rules).WithNum(i + 1).WithAIControlled(true)
		}
		universe, err := gameClient.GenerateUniverse(game, players)
		if err != nil {
			b.Fatalf("failed to generate universe")
		}

		fullGame := &cs.FullGame{
			Game:      game,
			Universe:  universe,
			TechStore: &cs.StaticTechStore,
			Players:   players,
		}

		err = c.SaveGame(b.Context(), fullGame.Game)
		if err != nil {
			b.Fatalf("failed to create game %v", err)
		}

		// generate some turns
		for range 50 {
			gameClient.GenerateTurn(fullGame.Game, fullGame.Universe, fullGame.Players)
		}

		if err := c.UpdateFullGame(b.Context(), fullGame); err != nil {
			b.Fatalf("failed to update game %v", err)
		}

		b.ResetTimer()
		for b.Loop() {
			if _, err := c.GetFullGame(b.Context(), game.ID); err != nil {
				b.Fatalf("failed to update game %v", err)
			}
		}
	})
}
