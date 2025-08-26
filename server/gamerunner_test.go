package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/sirgwain/craig-stars/ai"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/stretchr/testify/assert"
)

func createTestGameRunner(_ context.Context) GameRunner {
	dbConn := db.NewConn()
	cfg := &config.Config{}
	// cfg.Database.Filename = "../data/sqlx.db"
	cfg.Database.Filename = ":memory:"
	cfg.Database.DebugLogging = true
	if err := dbConn.Connect(cfg); err != nil {
		panic(fmt.Errorf("error connecting to test database: \n%w", err))
	}

	return NewGameRunner(dbConn, *cfg)
}

func Test_gameRunner_HostGame(t *testing.T) {

	gr := createTestGameRunner(t.Context())

	fullGame, err := gr.HostGame(1, cs.NewGameSettings().WithHost(cs.Humanoids()).WithAIPlayer(cs.AIDifficultyNormal, 0))

	if err != nil {
		t.Errorf("host game returned error \n%v", err)
		return
	}

	// make sure we generate some universes
	assert.Greater(t, len(fullGame.Planets), 0)
	assert.Greater(t, len(fullGame.Players), 0)
}

func Test_gameRunner_GenerateTurns(t *testing.T) {

	dbConn := db.NewConn()
	cfg := &config.Config{}
	cfg.Database.Filename = ":memory:"
	if err := dbConn.Connect(cfg); err != nil {
		panic(fmt.Errorf("error connecting to test database: \n%w", err))
	}

	// create a race per PRT
	for _, prt := range cs.PRTs {
		race := cs.NewRace()
		race.PRT = prt
		race.Name = fmt.Sprintf("%v", prt)
		race.PluralName = fmt.Sprintf("%vs", prt)
	}

	gr := NewGameRunner(dbConn, *cfg)

	// create a game with AI players for each PRT
	fullGame, err := gr.HostGame(1, cs.NewGameSettings().
		WithName("All Races Test").
		WithSize(cs.SizeMedium).
		WithAIPlayerRace(ai.Races[0], cs.AIDifficultyNormal, 0).
		WithAIPlayerRace(ai.Races[1], cs.AIDifficultyNormal, 1).
		WithAIPlayerRace(ai.Races[2], cs.AIDifficultyNormal, 2).
		WithAIPlayerRace(ai.Races[3], cs.AIDifficultyNormal, 3).
		WithAIPlayerRace(ai.Races[4], cs.AIDifficultyNormal, 0).
		WithAIPlayerRace(ai.Races[5], cs.AIDifficultyNormal, 1).
		WithAIPlayerRace(ai.Races[6], cs.AIDifficultyNormal, 2).
		WithAIPlayerRace(ai.Races[7], cs.AIDifficultyNormal, 3).
		WithAIPlayerRace(ai.Races[8], cs.AIDifficultyNormal, 0).
		WithAIPlayerRace(ai.Races[9], cs.AIDifficultyNormal, 1))

	if err != nil {
		t.Errorf("host game %v", err)
	}

	// generate 100 turns
	for i := 0; i < 100; i++ {
		gr := NewGameRunner(dbConn, *cfg)

		if _, err := gr.GenerateTurn(fullGame.ID); err != nil {
			t.Errorf("GenerateTurn failed on year %d: \n%v", fullGame.Game.Year, err)
		}
	}
}

func Test_gameRunner_getGuestNum(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     int
		wantErr  bool
	}{
		{"1", "guest-1-1", 1, false},
		{"20", "guest-29-20", 20, false},
		{"fail", "bob", 0, true},
		{"fail 2", "bob-1-bob", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &gameRunner{}
			u := cs.User{Username: tt.username}
			got, err := gr.getGuestNum(&u)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("gameRunner.getGuestNum() did not return error when expected")
				} else {
					t.Fatalf("gameRunner.getGuestNum() errored unexpectedly; err = \n%v", err)
				}
			}
			if got != tt.want {
				t.Errorf("gameRunner.getGuestNum() = %v, want %v", got, tt.want)
			}
		})
	}
}
