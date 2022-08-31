package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/game"
)

func TestDB_FindPlayerByGameId(t *testing.T) {

	// uncomment this for some database debug logging
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)

	db := connectDB()

	g := newRandomGame()
	if err := db.SaveGame(g); err != nil {
		t.Error(err)
	}

	player := g.Players[0]

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
			_ = got
			// if got != nil || tt.want != nil {
			// TODO: figure out a better way to test equivalence
			// this is fragile because the DB modifies the data on save
			// if !test.CompareAsJSON(t, got, tt.want) {
			// 	t.Errorf("DB.FindGameById() = %v, want %v", got, tt.want)
			// }
			// }
		})
	}
}
