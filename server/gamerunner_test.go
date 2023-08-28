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
	cfg.Database.UsersFilename = ":memory:"
	if err := db.Connect(cfg); err != nil {
		panic(fmt.Errorf("connect to test database, %w", err))
	}

	client, err := db.BeginTransaction()
	if err != nil {
		panic(fmt.Errorf("connect to test database, %w", err))
	}

	// create a test user with a single race
	user, err := cs.NewUser("admin", "admin", "admin@craig-stars.net", cs.RoleAdmin)
	if err != nil {
		panic(fmt.Errorf("generate test user, %w", err))
	}
	if err := client.CreateUser(user); err != nil {
		panic(fmt.Errorf("create test database user, %w", err))
	}

	race := cs.Humanoids()
	race.UserID = 1
	if err := client.CreateRace(&race); err != nil {
		panic(fmt.Errorf("create test user race, %w", err))
	}

	return &gameRunner{
		db:     client,
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

	dbClient := db.NewClient()
	cfg := &config.Config{}
	cfg.Database.Filename = ":memory:"
	cfg.Database.UsersFilename = ":memory:"
	if err := dbClient.Connect(cfg); err != nil {
		panic(fmt.Errorf("connect to test database, %w", err))
	}

	tx, err := dbClient.BeginTransaction()
	if err != nil {
		panic(fmt.Errorf("connect to test database, %w", err))
	}
	defer func() { dbClient.Rollback(tx) }()

	// create a test user
	user, err := cs.NewUser("admin", "admin", "admin@craig-stars.net", cs.RoleAdmin)
	if err != nil {
		panic(fmt.Errorf("generate test user, %w", err))
	}
	if err := tx.CreateUser(user); err != nil {
		panic(fmt.Errorf("create test database user, %w", err))
	}

	prts := make(map[cs.PRT]cs.Race, len(cs.PRTs))

	// create a race per PRT
	for _, prt := range cs.PRTs {
		race := cs.NewRace()
		race.PRT = prt
		race.Name = fmt.Sprintf("%v", prt)
		race.PluralName = fmt.Sprintf("%vs", prt)
	}

	gr := gameRunner{
		db:     tx,
		client: cs.NewGamer(),
	}

	// create a game with AI players for each PRT
	fullGame, err := gr.HostGame(1, cs.NewGameSettings().
		WithName("All Races Test").
		WithSize(cs.SizeMedium).
		WithAIPlayerRace(prts[cs.HE], cs.AIDifficultyNormal, 0).
		WithAIPlayerRace(prts[cs.SS], cs.AIDifficultyNormal, 1).
		WithAIPlayerRace(prts[cs.WM], cs.AIDifficultyNormal, 2).
		WithAIPlayerRace(prts[cs.CA], cs.AIDifficultyNormal, 3).
		WithAIPlayerRace(prts[cs.IS], cs.AIDifficultyNormal, 0).
		WithAIPlayerRace(prts[cs.SD], cs.AIDifficultyNormal, 1).
		WithAIPlayerRace(prts[cs.PP], cs.AIDifficultyNormal, 2).
		WithAIPlayerRace(prts[cs.IT], cs.AIDifficultyNormal, 3).
		WithAIPlayerRace(prts[cs.AR], cs.AIDifficultyNormal, 0).
		WithAIPlayerRace(prts[cs.JoaT], cs.AIDifficultyNormal, 1))

	if err != nil {
		t.Errorf("host game %v", err)
	}

	// commit the game changes before we gen turns
	err = dbClient.Commit(tx)
	if err != nil {
		t.Errorf("commit %v", err)
	}

	// generate 100 turns
	for i := 0; i < 50; i++ {
		tx, err := dbClient.BeginTransaction()
		if err != nil {
			t.Errorf("begin transaction for turn %d %v", i, err)
			return
		}
		defer func() { dbClient.Rollback(tx) }()

		gr := gameRunner{
			db:     tx,
			client: cs.NewGamer(),
		}

		if _, err := gr.GenerateTurn(fullGame.ID); err != nil {
			t.Errorf("generate turn %v", err)
		}

		if err := dbClient.Commit(tx); err != nil {
			t.Errorf("commit transaction for turn %d %v", i, err)
			return
		}
	}
}
