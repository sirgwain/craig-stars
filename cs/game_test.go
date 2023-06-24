package cs

import "testing"

type testPlayerGetter struct {
	players []*Player
}

func newTestPlayerGetter(players ...*Player) playerGetter {
	return &testPlayerGetter{players}
}

func (pg *testPlayerGetter) getPlayer(num int) *Player {
	for _, p := range pg.players {
		if p.Num == num {
			return p
		}
	}
	return nil
}

func TestGame_GenerateHash(t *testing.T) {
	tests := []struct {
		name string
		id   int64
		salt string
		want string
	}{
		{"hash for salt", 1, "salt", "2c7ece3537e3af629566c1604fdb30"},
		{"hash for salt", 2, "salt", "cff02c8df0b1c9c986e005dc1c1ca9"},
		{"hash for salt", 2, "salt2", "4e042167d14ca83767029c176782e9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{}
			g.ID = tt.id
			if got := g.GenerateHash(tt.salt); got != tt.want {
				t.Errorf("Game.GenerateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
