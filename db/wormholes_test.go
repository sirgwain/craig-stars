package db

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/game"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateWormhole(t *testing.T) {
	type args struct {
		c        *client
		wormhole *game.Wormhole
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &game.Wormhole{
			MapObject: game.MapObject{GameID: 1, Name: "test"},
		},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			g := tt.args.c.createTestGame()
			tt.args.wormhole.GameID = g.ID

			want := *tt.args.wormhole
			err := tt.args.c.createWormhole(tt.args.wormhole, tt.args.c.db)

			// id is automatically added
			want.ID = tt.args.wormhole.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWormhole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.wormhole, &want) {
				t.Errorf("CreateWormhole() = \n%v, want \n%v", tt.args.wormhole, want)
			}
		})
	}
}

func TestGetWormholes(t *testing.T) {
	c := connectTestDB()
	g := c.createTestGame()

	// start with 1 wormhole from connectTestDB
	result, err := c.getWormholesForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, []*game.Wormhole{}, result)

	wormhole := game.Wormhole{MapObject: game.MapObject{GameID: g.ID}}
	if err := c.createWormhole(&wormhole, c.db); err != nil {
		t.Errorf("create wormhole %s", err)
		return
	}

	result, err = c.getWormholesForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestGetWormhole(t *testing.T) {
	c := connectTestDB()
	g := c.createTestGame()
	wormhole := game.Wormhole{MapObject: game.MapObject{GameID: g.ID, Name: "name", Type: game.MapObjectTypeWormhole}}
	if err := c.createWormhole(&wormhole, c.db); err != nil {
		t.Errorf("create wormhole %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *game.Wormhole
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got wormhole", args{id: wormhole.ID}, &wormhole, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetWormhole(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWormhole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("GetWormhole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateWormhole(t *testing.T) {
	c := connectTestDB()
	g := c.createTestGame()
	wormhole := game.Wormhole{MapObject: game.MapObject{GameID: g.ID}}
	if err := c.createWormhole(&wormhole, c.db); err != nil {
		t.Errorf("create wormhole %s", err)
		return
	}

	wormhole.Name = "Test2"
	if err := c.UpdateWormhole(&wormhole); err != nil {
		t.Errorf("update wormhole %s", err)
		return
	}

	updated, err := c.GetWormhole(wormhole.ID)

	if err != nil {
		t.Errorf("get wormhole %s", err)
		return
	}

	assert.Equal(t, wormhole.Name, updated.Name)
	assert.Less(t, wormhole.UpdatedAt, updated.UpdatedAt)

}
