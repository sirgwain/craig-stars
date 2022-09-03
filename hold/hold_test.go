package hold

import (
	"testing"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/game"
	"github.com/sirgwain/craig-stars/test"
)

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

func TestDB_Connect(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Should connect"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := DB{}
			// dir, err := ioutil.TempDir("../tmp", "hold_testdb_connect")
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }
			cfg := &config.Config{}
			// cfg.Database.Filename = fmt.Sprintf("%s/bolt.db", dir)
			cfg.Database.Filename = "../tmp/bolt.db"
			db.Connect(cfg)

			g := newRandomGame()
			if err := db.CreateGame(g.Game); err != nil {
				t.Error(err)
				return
			}
			if err := db.SaveGame(g); err != nil {
				t.Error(err)
				return
			}

			loaded, err := db.FindGameById(g.ID)
			if err != nil {
				t.Error(err)
				return
			}

			if !test.CompareAsJSON(t, loaded, g) {
				t.Errorf("games don't match = %v, want %v", loaded, g)
			}

		})
	}
}
