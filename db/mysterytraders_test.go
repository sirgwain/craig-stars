package db

import (
	"reflect"
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
			MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: 1}, Name: "test"}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			game := tt.args.c.createTestGame()
			tt.args.mysteryTrader.GameID = game.ID

			want := *tt.args.mysteryTrader
			err := tt.args.c.createMysteryTrader(tt.args.mysteryTrader)

			// id is automatically added
			want.ID = tt.args.mysteryTrader.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMysteryTrader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.mysteryTrader, &want) {
				t.Errorf("CreateMysteryTrader() = \n%v, want \n%v", tt.args.mysteryTrader, want)
			}
		})
	}
}

func TestGetMysteryTraders(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame()

	// start with 1 mysteryTrader from connectTestDB
	result, err := c.getMysteryTradersForGame(game.ID)
	assert.Nil(t, err)
	assert.Equal(t, []*cs.MysteryTrader{}, result)

	mysteryTrader := cs.MysteryTrader{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: game.ID}}}
	if err := c.createMysteryTrader(&mysteryTrader); err != nil {
		t.Errorf("create mysteryTrader %s", err)
		return
	}

	result, err = c.getMysteryTradersForGame(game.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestGetMysteryTrader(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame()
	mysteryTrader := cs.MysteryTrader{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: game.ID}, Name: "name", Type: cs.MapObjectTypeMysteryTrader}}
	if err := c.createMysteryTrader(&mysteryTrader); err != nil {
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
		{"Got mysteryTrader", args{id: mysteryTrader.ID}, &mysteryTrader, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetMysteryTrader(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMysteryTrader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("GetMysteryTrader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateMysteryTrader(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame()
	mysteryTrader := cs.MysteryTrader{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: game.ID}}}
	if err := c.createMysteryTrader(&mysteryTrader); err != nil {
		t.Errorf("create mysteryTrader %s", err)
		return
	}

	mysteryTrader.Name = "Test2"
	if err := c.updateMysteryTrader(&mysteryTrader); err != nil {
		t.Errorf("update mysteryTrader %s", err)
		return
	}

	updated, err := c.GetMysteryTrader(mysteryTrader.ID)

	if err != nil {
		t.Errorf("get mysteryTrader %s", err)
		return
	}

	assert.Equal(t, mysteryTrader.Name, updated.Name)
	assert.Less(t, mysteryTrader.UpdatedAt, updated.UpdatedAt)

}
