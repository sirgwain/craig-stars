package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestSaveWormhole(t *testing.T) {
	type args struct {
		c        *client
		wormhole *cs.Wormhole
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.Wormhole{
			GameDBObject: cs.GameDBObject{GameID: 1},
			MapObject:    cs.MapObject{Type: cs.MapObjectTypeWormhole, Name: "test"}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			game := tt.args.c.createTestGame(t.Context())
			tt.args.wormhole.GameID = game.ID

			want := *tt.args.wormhole
			err := tt.args.c.SaveWormhole(t.Context(), tt.args.wormhole)

			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("SaveWormhole() did not return error when expected")
				} else {
					t.Fatalf("SaveWormhole() errored unexpectedly; err = \n%v", err)
				}
			}

			got := tt.args.wormhole
			// DBObject is returned
			want.GameDBObject = got.GameDBObject
			test.CompareAsJSON(t, got, want)
		})
	}
}

func TestGetWormholes(t *testing.T) {
	c := connectTestDB()

	game := c.createTestGame(t.Context())

	// start with 1 wormhole from connectTestDB
	result, err := c.GetWormholesForGame(t.Context(), game.ID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	if err := c.SaveWormhole(t.Context(), &cs.Wormhole{GameDBObject: cs.GameDBObject{GameID: game.ID}, MapObject: cs.MapObject{}}); err != nil {
		t.Errorf("create wormhole %s", err)
		return
	}

	result, err = c.GetWormholesForGame(t.Context(), game.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestGetWormhole(t *testing.T) {
	c := connectTestDB()

	game := c.createTestGame(t.Context())
	wormhole := &cs.Wormhole{GameDBObject: cs.GameDBObject{GameID: game.ID}, MapObject: cs.MapObject{Type: cs.MapObjectTypeWormhole}}
	if err := c.SaveWormhole(t.Context(), wormhole); err != nil {
		t.Errorf("create wormhole %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.Wormhole
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got wormhole", args{id: wormhole.ID}, wormhole, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetWormhole(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetWormhole() did not return error when expected")
				} else {
					t.Fatalf("GetWormhole() errored unexpectedly; err = \n%v", err)
				}
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}

			test.CompareAsJSON(t, got, tt.want)
		})
	}
}

func TestUpdateWormhole(t *testing.T) {
	c := connectTestDB()

	game := c.createTestGame(t.Context())
	wormhole := &cs.Wormhole{GameDBObject: cs.GameDBObject{GameID: game.ID}, MapObject: cs.MapObject{}}
	if err := c.SaveWormhole(t.Context(), wormhole); err != nil {
		t.Errorf("create wormhole %s", err)
		return
	}

	wormhole.Name = "Test2"
	if err := c.SaveWormhole(t.Context(), wormhole); err != nil {
		t.Errorf("update wormhole %s", err)
		return
	}

	updated, err := c.GetWormhole(t.Context(), wormhole.ID)

	if err != nil {
		t.Errorf("get wormhole %s", err)
		return
	}

	assert.Equal(t, wormhole.Name, updated.Name)

}
