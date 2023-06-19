package db

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/game"
	"github.com/stretchr/testify/assert"
)

func connectDB() *DB {
	db := &DB{}
	cfg := &config.Config{}
	// cfg.Database.Filename = "../tmp/test-data.db"
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
			if got, err := db.GetGames(); !reflect.DeepEqual(got, tt.want) || err != nil {
				t.Errorf("DB.GetGames() = %v, want %v, err %v", got, tt.want, err)
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

	// add a production queue item
	previousItemCount := len(g.Planets[0].ProductionQueue)
	g.Planets[0].ProductionQueue = append(g.Planets[0].ProductionQueue, game.ProductionQueueItem{Type: game.QueueItemTypeMine, Quantity: 5})

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
	assert.Equal(t, previousItemCount+1, len(loaded.Planets[0].ProductionQueue))

	// make sure our production queue item saved
	assert.Equal(t, game.QueueItemTypeMine, loaded.Planets[0].ProductionQueue[len(loaded.Planets[0].ProductionQueue)-1].Type)
	assert.Equal(t, 5, loaded.Planets[0].ProductionQueue[len(loaded.Planets[0].ProductionQueue)-1].Quantity)
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
			if got, err := db.GetGamesByUser(tt.args.userID); !reflect.DeepEqual(got, tt.want) || err != nil {
				t.Errorf("DB.GetGamesByUser() = %v, want %v, err %v", got, tt.want, err)
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
				// TODO: figure out a better way to test equivalence
				// this is fragile because the DB modifies the data on save
				// if !test.CompareAsJSON(t, got, tt.want) {
				// 	t.Errorf("DB.FindGameById() = %v, want %v", got, tt.want)
				// }
			}
		})
	}
}
