package server

import (
	"fmt"
	"testing"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/stretchr/testify/assert"
)

func createTestGameRunner() GameRunner {
	db := db.NewClient()
	cfg := &config.Config{}
	// cfg.Database.Filename = "../data/sqlx.db"
	cfg.Database.Filename = ":memory:"
	if err := db.Connect(cfg); err != nil {
		panic(fmt.Errorf("connect to test database, %w", err))
	}

	// create a test user with a single race
	user, err := cs.NewUser("admin", "admin", "admin@craig-stars.net", cs.RoleAdmin)
	if (err != nil) {
		panic(fmt.Errorf("generate test user, %w", err))
	}
	if err := db.CreateUser(user); err != nil {
		panic(fmt.Errorf("create test database user, %w", err))
	}

	race := cs.Humanoids()
	race.UserID = 1
	if err := db.CreateRace(&race); err != nil {
		panic(fmt.Errorf("create test user race, %w", err))
	}

	return &gameRunner{
		db:     db,
		client: cs.NewGamer(),
	}
}

func Test_gameRunner_HostGame(t *testing.T) {

	gr := createTestGameRunner()

	fullGame, err := gr.HostGame(1, cs.NewGameSettings().WithHost(1).WithAIPlayer(cs.AIDifficultyNormal, 0))

	if err != nil {
		t.Errorf("host game %v", err)
	}

	// make sure we generate some universes
	assert.Greater(t, len(fullGame.Planets), 0)
	assert.Greater(t, len(fullGame.Players), 0)
}

func Test_gameRunner_GenerateTurns(t *testing.T) {

	gr := createTestGameRunner()

	// create a game with a couple AI players
	fullGame, err := gr.HostGame(1, cs.NewGameSettings().WithAIPlayer(cs.AIDifficultyNormal, 0).WithAIPlayer(cs.AIDifficultyNormal, 1))

	if err != nil {
		t.Errorf("host game %v", err)
	}

	// generate 100 turns
	for i := 0; i < 100; i++ {
		if err := gr.GenerateTurn(fullGame.ID); err != nil {
			t.Errorf("generate turn %v", err)
		}
	}
}
