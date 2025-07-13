package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestSaveMineField(t *testing.T) {
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
			GameDBObject: cs.GameDBObject{GameID: 1},
			MapObject:    cs.MapObject{Type: cs.MapObjectTypeMineField, Name: "test"},
		},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			g, player := tt.args.c.createTestGameWithPlayer(t.Context())
			tt.args.mineField.GameID = g.ID
			tt.args.mineField.PlayerNum = player.Num

			want := *tt.args.mineField
			err := tt.args.c.SaveMineField(t.Context(), tt.args.mineField)

			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("SaveMineField() did not return error when expected")
				} else {
					t.Fatalf("SaveMineField() errored unexpectedly; err = \n%v", err)
				}
			}

			got := tt.args.mineField
			// DBObject is returned
			want.GameDBObject = got.GameDBObject
			test.CompareAsJSON(t, got, &want)
		})
	}
}

func TestGetMineField(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	g, player := c.createTestGameWithPlayer(t.Context())

	mineField := &cs.MineField{
		GameDBObject:  cs.GameDBObject{GameID: g.ID},
		MapObject:     cs.MapObject{PlayerNum: player.Num, Name: "name", Type: cs.MapObjectTypeMineField},
		MineFieldType: cs.MineFieldTypeStandard,
	}

	if err := c.SaveMineField(t.Context(), mineField); err != nil {
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
		{"Got mineField", args{id: mineField.ID}, mineField, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetMineField(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetMineField() did not return error when expected")
				} else {
					t.Fatalf("GetMineField() errored unexpectedly; err = \n%v", err)
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

func TestGetMineFields(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	g, player := c.createTestGameWithPlayer(t.Context())

	// start with 1 planet from connectTestDB
	result, err := c.getMineFieldsForGame(t.Context(), g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	if err := c.SaveMineField(t.Context(), &cs.MineField{GameDBObject: cs.GameDBObject{GameID: g.ID}, MapObject: cs.MapObject{PlayerNum: player.Num}}); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	result, err = c.getMineFieldsForGame(t.Context(), g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestUpdateMineField(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	g, player := c.createTestGameWithPlayer(t.Context())
	mineField := &cs.MineField{GameDBObject: cs.GameDBObject{GameID: g.ID}, MapObject: cs.MapObject{PlayerNum: player.Num}}
	if err := c.SaveMineField(t.Context(), mineField); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	mineField.Name = "Test2"
	if err := c.SaveMineField(t.Context(), mineField); err != nil {
		t.Errorf("update planet %s", err)
		return
	}

	updated, err := c.GetMineField(t.Context(), mineField.ID)

	if err != nil {
		t.Errorf("get planet %s", err)
		return
	}

	assert.Equal(t, mineField.Name, updated.Name)

}
