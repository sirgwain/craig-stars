package db

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/game"
	"github.com/stretchr/testify/assert"
)

func connectDB() *client {
	c := &client{}
	cfg := &config.Config{}
	// cfg.Database.Filename = "../tmp/test-data.db"
	cfg.Database.Filename = ":memory:"
	c.Connect(cfg)
	c.MigrateAll()
	return c
}

func newRandomGame() *game.FullGame {
	client := game.NewClient()
	g := client.CreateGame(1, *game.NewGameSettings())
	player := client.NewPlayer(1, game.Humanoids(), &g.Rules)
	players := []*game.Player{player}
	universe, err := client.GenerateUniverse(&g, players)
	if err != nil {
		panic(err)
	}

	fg := game.FullGame{
		Game:     &g,
		Players:  players,
		Universe: universe,
	}

	return &fg
}

func TestDB_GetGames(t *testing.T) {
	c := connectDB()
	tests := []struct {
		name string
		want []*game.Game
	}{
		{"No games", []*game.Game{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := c.GetGames(); !reflect.DeepEqual(got, tt.want) || err != nil {
				t.Errorf("c.GetGames() = %v, want %v, err %v", got, tt.want, err)
			}
		})
	}
}

func TestDB_CreateGame(t *testing.T) {
	c := connectDB()
	type args struct {
		game *game.Game
	}
	tests := []struct {
		name string
		args args
	}{
		{"Create Game", args{&game.Game{Rules: game.NewRules()}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.CreateGame(tt.args.game)
			if err != nil {
				t.Error(err)
			}
			assert.NotZero(t, tt.args.game.ID)
		})
	}
}

func TestDB_SaveGame(t *testing.T) {
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)

	c := connectDB()
	g := newRandomGame()
	err := c.SaveGame(g)
	if err != nil {
		t.Error(err)
	}

	g, _ = c.FindGameById(g.ID)

	// set a pop value
	pop := 100_000
	g.Planets[0].SetPopulation(pop)
	g.Planets[0].Dirty = true
	g.Players[0].PlanetIntels[0].Population = uint(pop)
	g.Players[0].PlanetIntels[0].Dirty = true

	// add a production queue item
	previousItemCount := len(g.Planets[0].ProductionQueue)
	g.Planets[0].ProductionQueue = append(g.Planets[0].ProductionQueue, game.ProductionQueueItem{Type: game.QueueItemTypeMine, Quantity: 5})

	err = c.SaveGame(g)
	if err != nil {
		t.Error(err)
	}

	loaded, err := c.FindGameById(g.ID)
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
		userID int64
	}
	tests := []struct {
		name string
		db   *client
		args args
		want []game.Game
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{}
			if got, err := c.GetGamesByUser(tt.args.userID); !reflect.DeepEqual(got, tt.want) || err != nil {
				t.Errorf("c.GetGamesByUser() = %v, want %v, err %v", got, tt.want, err)
			}
		})
	}
}

func TestDB_DeleteGameById(t *testing.T) {
	type args struct {
		gameID int64
	}
	tests := []struct {
		name string
		db   *client
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{}
			c.DeleteGameById(tt.args.gameID)
		})
	}
}

func TestDB_FindGameById(t *testing.T) {

	c := connectDB()

	g := newRandomGame()
	if err := c.SaveGame(g); err != nil {
		t.Error(err)
	}

	tests := []struct {
		name    string
		id      int64
		want    *game.FullGame
		wantErr bool
	}{
		{"No Game", 2, nil, false},
		{"Find Game", 1, g, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.FindGameById(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("c.FindGameById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil || tt.want != nil {
				// TODO: figure out a better way to test equivalence
				// this is fragile because the DB modifies the data on save
				// if !test.CompareAsJSON(t, got, tt.want) {
				// 	t.Errorf("c.FindGameById() = %v, want %v", got, tt.want)
				// }
			}
		})
	}
}
