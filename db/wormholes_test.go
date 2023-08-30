package db

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateWormhole(t *testing.T) {
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
			MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: 1}, Name: "test"}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			game := tt.args.c.createTestGame()
			tt.args.wormhole.GameID = game.ID

			want := *tt.args.wormhole
			err := tt.args.c.createWormhole(tt.args.wormhole)

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
	defer func() { closeTestDB(c) }()

	game := c.createTestGame()

	// start with 1 wormhole from connectTestDB
	result, err := c.getWormholesForGame(game.ID)
	assert.Nil(t, err)
	assert.Equal(t, []*cs.Wormhole{}, result)

	wormhole := cs.Wormhole{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: game.ID}}}
	if err := c.createWormhole(&wormhole); err != nil {
		t.Errorf("create wormhole %s", err)
		return
	}

	result, err = c.getWormholesForGame(game.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestGetWormhole(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame()
	wormhole := cs.Wormhole{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: game.ID}, Name: "name", Type: cs.MapObjectTypeWormhole}}
	if err := c.createWormhole(&wormhole); err != nil {
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
	defer func() { closeTestDB(c) }()

	game := c.createTestGame()
	wormhole := cs.Wormhole{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: game.ID}}}
	if err := c.createWormhole(&wormhole); err != nil {
		t.Errorf("create wormhole %s", err)
		return
	}

	wormhole.Name = "Test2"
	if err := c.updateWormhole(&wormhole); err != nil {
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
