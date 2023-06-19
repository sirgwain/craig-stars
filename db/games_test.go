package db

import (
	"os"
	"reflect"
	"testing"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/game"
	"github.com/stretchr/testify/assert"
)

func connectDB() *DB {
	db := &DB{}
	cfg := &config.Config{}
	cfg.Database.Filename = ":memory:"
	db.Connect(cfg)
	db.MigrateAll()
	return db
}

func newRandomGame() *game.Game {
	g := game.NewGame()
	// g.Size = game.SizeTiny
	// g.Density = game.DensitySparse
	g.AddPlayer(game.NewPlayer(1, game.NewRace()))
	g.GenerateUniverse()
	return g
}

func compareAsJSON(t *testing.T, g1 *game.Game, g2 *game.Game) bool {
	if g1 == nil && g2 == nil {
		return true
	} else if g1 == nil && g2 != nil || g1 != nil && g2 == nil {
		return false
	} else {
		json1, err := json.MarshalIndent(g1, "", "  ")
		if err != nil {
			t.Errorf("Failed to compare %s, error = %v", g1, err)
		}
		json2, err := json.MarshalIndent(g2, "", "  ")
		if err != nil {
			t.Errorf("Failed to compare %s, error = %v", g2, err)
		}

		if string(json1) != string(json2) {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			log.Debug().Msgf("\n\njson1: %s\n", string(json1))
			log.Debug().Msgf("\n\njson2: %s\n", string(json2))
			return false
		} else {
			return true
		}
	}
}

func TestDB_GetGames(t *testing.T) {
	db := connectDB()
	tests := []struct {
		name string
		want []game.Game
	}{
		{"No games", []game.Game{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.GetGames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.GetGames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_CreateGame(t *testing.T) {
	db := connectDB()
	type args struct {
		game *game.Game
	}
	tests := []struct {
		name string
		args args
	}{
		{"Create Game", args{newRandomGame()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.CreateGame(tt.args.game)
			if err != nil {
				t.Error(err)
			}
			assert.NotZero(t, tt.args.game.ID)
			assert.NotZero(t, tt.args.game.Players[0].ID)
			assert.NotZero(t, tt.args.game.Players[0].Race.ID)
			assert.NotZero(t, tt.args.game.Planets[0].ID)
		})
	}
}

func TestDB_SaveGame(t *testing.T) {
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)

	db := connectDB()
	g := newRandomGame()
	err := db.CreateGame(g)
	if err != nil {
		t.Error(err)
	}

	g, _ = db.FindGameById(g.ID)

	// set a pop value
	pop := 100_000
	g.Planets[0].SetPopulation(pop)
	g.Planets[0].Dirty = true
	g.Players[0].PlanetIntels[0].Population = uint(pop)
	g.Players[0].PlanetIntels[0].Dirty = true
	err = db.SaveGame(g)
	if err != nil {
		t.Error(err)
	}

	loaded, err := db.FindGameById(g.ID)
	if err != nil {
		t.Error(err)
	}

	// make sure our pop saved
	assert.Equal(t, pop, loaded.Planets[0].Population())
	assert.Equal(t, pop, int(loaded.Players[0].PlanetIntels[0].Population))
}

func TestDB_GetGamesByUser(t *testing.T) {
	type args struct {
		userID uint
	}
	tests := []struct {
		name string
		db   *DB
		args args
		want []game.Game
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{}
			if got := db.GetGamesByUser(tt.args.userID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.GetGamesByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_DeleteGameById(t *testing.T) {
	type args struct {
		id uint
	}
	tests := []struct {
		name string
		db   *DB
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{}
			db.DeleteGameById(tt.args.id)
		})
	}
}

func TestDB_FindGameById(t *testing.T) {

	db := connectDB()
	g := newRandomGame()
	g.Name = "Test Game"
	err := db.CreateGame(g)

	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name    string
		id      uint
		want    *game.Game
		wantErr bool
	}{
		{"No Game", 2, nil, false},
		{"Find Game", 1, g, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.FindGameById(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.FindGameById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil || tt.want != nil {
				if !compareAsJSON(t, got, tt.want) {
					t.Errorf("DB.FindGameById() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
