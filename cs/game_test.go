//go:build !wasi && !wasm

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

func TestGameWithPlayers_IsSinglePlayer(t *testing.T) {
	tests := []struct {
		name    string
		players []GamePlayer
		want    bool
	}{
		{"host vs ai", []GamePlayer{{UserID: 1}, {AIControlled: true}}, true},
		{"hot seat", []GamePlayer{{UserID: 1}, {UserID: 1}, {AIControlled: true}}, true},
		{"two users", []GamePlayer{{UserID: 1}, {UserID: 2}}, false},
		{"hot seat with open slot", []GamePlayer{{UserID: 1}, {UserID: 1}, {}}, false},
		{"two open slots", []GamePlayer{{}, {}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GameWithPlayers{Players: tt.players}
			if got := g.IsSinglePlayer(); got != tt.want {
				t.Errorf("IsSinglePlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameSettings_IsSinglePlayer(t *testing.T) {
	tests := []struct {
		name    string
		players []NewGamePlayer
		want    bool
	}{
		{"host vs ai", []NewGamePlayer{{Type: NewGamePlayerTypeHost}, {Type: NewGamePlayerTypeAI}}, true},
		{"hot seat", []NewGamePlayer{{Type: NewGamePlayerTypeHost}, {Type: NewGamePlayerTypeHost}}, true},
		{"host and open", []NewGamePlayer{{Type: NewGamePlayerTypeHost}, {Type: NewGamePlayerTypeOpen}}, false},
		{"host and guest", []NewGamePlayer{{Type: NewGamePlayerTypeHost}, {Type: NewGamePlayerTypeGuest}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := &GameSettings{Players: tt.players}
			if got := settings.IsSinglePlayer(); got != tt.want {
				t.Errorf("IsSinglePlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}
