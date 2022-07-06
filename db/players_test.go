package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/game"
	"github.com/sirgwain/craig-stars/test"
)

func TestDB_FindPlayerByGameId(t *testing.T) {

	// uncomment this for some database debug logging
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)

	db := connectDB()

	g := game.NewGame()
	g.AddPlayer(game.NewPlayer(1, game.NewRace()))
	if err := db.CreateGame(g); err != nil {
		t.Error(err)
	}

	if err := g.GenerateUniverse(); err != nil {
		t.Error(err)
	}

	if err := db.SaveGame(g); err != nil {
		t.Error(err)
	}

	player := &g.Players[0]

	tests := []struct {
		name    string
		gameId  uint
		userId  uint
		want    *game.Player
		wantErr bool
	}{
		{"No Player", 1, 2, nil, false},
		{"Find Player", 1, 1, player, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.FindPlayerByGameId(tt.gameId, tt.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.FindGameById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil || tt.want != nil {
				if !test.CompareAsJSON(t, got, tt.want) {
					t.Errorf("DB.FindGameById() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
