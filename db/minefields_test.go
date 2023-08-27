package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateMineField(t *testing.T) {
	type args struct {
		c         *client
		mineField *cs.MineField
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.MineField{
			MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: 1}, Name: "test"},
		},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			g, player := tt.args.c.createTestGameWithPlayer()
			tt.args.mineField.GameID = g.ID
			tt.args.mineField.PlayerNum = player.Num

			want := *tt.args.mineField
			err := tt.args.c.createMineField(tt.args.mineField, tt.args.c.db)

			// id is automatically added
			want.ID = tt.args.mineField.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMineField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !test.CompareAsJSON(t, tt.args.mineField, &want) {
				t.Errorf("CreateMineField() = \n%v, want \n%v", tt.args.mineField, want)
			}
		})
	}
}

func TestGetMineField(t *testing.T) {
	c := connectTestDB()
	g, player := c.createTestGameWithPlayer()

	mineField := cs.MineField{
		MapObject:     cs.MapObject{GameDBObject: cs.GameDBObject{GameID: g.ID}, PlayerNum: player.Num, Name: "name", Type: cs.MapObjectTypeMineField},
		MineFieldType: cs.MineFieldTypeStandard,
	}
	if err := c.createMineField(&mineField, c.db); err != nil {
		t.Errorf("create mineField %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.MineField
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got mineField", args{id: mineField.ID}, &mineField, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetMineField(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMineField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("GetMineField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMineFields(t *testing.T) {
	c := connectTestDB()
	g, player := c.createTestGameWithPlayer()

	// start with 1 planet from connectTestDB
	result, err := c.getMineFieldsForGame(c.db, g.ID)
	assert.Nil(t, err)
	assert.Equal(t, []*cs.MineField{}, result)

	mineField := cs.MineField{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: g.ID}, PlayerNum: player.Num}}
	if err := c.createMineField(&mineField, c.db); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	result, err = c.getMineFieldsForGame(c.db, g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestUpdateMineField(t *testing.T) {
	c := connectTestDB()
	g, player := c.createTestGameWithPlayer()
	planet := cs.MineField{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: g.ID}, PlayerNum: player.Num}}
	if err := c.createMineField(&planet, c.db); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	planet.Name = "Test2"
	if err := c.UpdateMineField(&planet); err != nil {
		t.Errorf("update planet %s", err)
		return
	}

	updated, err := c.GetMineField(planet.ID)

	if err != nil {
		t.Errorf("get planet %s", err)
		return
	}

	assert.Equal(t, planet.Name, updated.Name)
	assert.Less(t, planet.UpdatedAt, updated.UpdatedAt)

}
