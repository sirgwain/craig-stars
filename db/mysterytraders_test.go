package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateMysteryTrader(t *testing.T) {
	type args struct {
		c             *client
		mysteryTrader *cs.MysteryTrader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.MysteryTrader{
			GameDBObject: cs.GameDBObject{GameID: 1},
			MapObject:    cs.MapObject{Type: cs.MapObjectTypeMysteryTrader, Name: "test"}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			game := tt.args.c.createTestGame(t.Context())
			tt.args.mysteryTrader.GameID = game.ID

			want := *tt.args.mysteryTrader
			got, err := tt.args.c.CreateMysteryTrader(t.Context(), tt.args.mysteryTrader)

			// id is automatically added
			want.GameDBObject = got.GameDBObject
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("CreateMysteryTrader() did not return error when expected")
				} else {
					t.Fatalf("CreateMysteryTrader() errored unexpectedly; err = \n%v", err)
				}
			}
			test.CompareAsJSON(t, got, want)
		})
	}
}

func TestGetMysteryTraders(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame(t.Context())

	// start with 1 mysteryTrader from connectTestDB
	result, err := c.GetMysteryTradersForGame(t.Context(), game.ID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	_, err = c.CreateMysteryTrader(t.Context(), &cs.MysteryTrader{GameDBObject: cs.GameDBObject{GameID: game.ID}, MapObject: cs.MapObject{}})
	if err != nil {
		t.Errorf("create mysteryTrader %s", err)
		return
	}

	result, err = c.GetMysteryTradersForGame(t.Context(), game.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestGetMysteryTrader(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame(t.Context())
	mysteryTrader, err := c.CreateMysteryTrader(t.Context(), &cs.MysteryTrader{GameDBObject: cs.GameDBObject{GameID: game.ID}, MapObject: cs.MapObject{Name: "name", Type: cs.MapObjectTypeMysteryTrader}})
	if err != nil {
		t.Errorf("create mysteryTrader %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.MysteryTrader
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got mysteryTrader", args{id: mysteryTrader.ID}, mysteryTrader, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetMysteryTrader(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetMysteryTrader() did not return error when expected")
				} else {
					t.Fatalf("GetMysteryTrader() errored unexpectedly; err = \n%v", err)
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

func TestUpdateMysteryTrader(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame(t.Context())
	mysteryTrader, err := c.CreateMysteryTrader(t.Context(), &cs.MysteryTrader{GameDBObject: cs.GameDBObject{GameID: game.ID}, MapObject: cs.MapObject{}})
	if err != nil {
		t.Errorf("create mysteryTrader %s", err)
		return
	}

	mysteryTrader.Name = "Test2"
	if err := c.UpdateMysteryTrader(t.Context(), mysteryTrader); err != nil {
		t.Errorf("update mysteryTrader %s", err)
		return
	}

	updated, err := c.GetMysteryTrader(t.Context(), mysteryTrader.ID)

	if err != nil {
		t.Errorf("get mysteryTrader %s", err)
		return
	}

	assert.Equal(t, mysteryTrader.Name, updated.Name)

}
