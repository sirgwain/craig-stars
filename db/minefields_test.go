package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestSaveMinefield(t *testing.T) {
	type args struct {
		c         *client
		minefield *cs.Minefield
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.Minefield{
			GameDBObject: cs.GameDBObject{GameID: 1},
			MapObject:    cs.MapObject{Type: cs.MapObjectTypeMinefield, Name: "test"},
		},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			g, player := tt.args.c.createTestGameWithPlayer(t.Context())
			tt.args.minefield.GameID = g.ID
			tt.args.minefield.PlayerNum = player.Num

			want := *tt.args.minefield
			err := tt.args.c.SaveMinefield(t.Context(), tt.args.minefield)

			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("SaveMinefield() did not return error when expected")
				} else {
					t.Fatalf("SaveMinefield() errored unexpectedly; err = \n%v", err)
				}
			}

			got := tt.args.minefield
			// DBObject is returned
			want.GameDBObject = got.GameDBObject
			test.CompareAsJSON(t, got, &want)
		})
	}
}

func TestGetMinefield(t *testing.T) {
	c := connectTestDB()

	g, player := c.createTestGameWithPlayer(t.Context())

	minefield := &cs.Minefield{
		GameDBObject:  cs.GameDBObject{GameID: g.ID},
		MapObject:     cs.MapObject{PlayerNum: player.Num, Name: "name", Type: cs.MapObjectTypeMinefield},
		MinefieldType: cs.MinefieldTypeStandard,
	}

	if err := c.SaveMinefield(t.Context(), minefield); err != nil {
		t.Errorf("create minefield %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.Minefield
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got minefield", args{id: minefield.ID}, minefield, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetMinefield(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetMinefield() did not return error when expected")
				} else {
					t.Fatalf("GetMinefield() errored unexpectedly; err = \n%v", err)
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

func TestGetMinefields(t *testing.T) {
	c := connectTestDB()

	g, player := c.createTestGameWithPlayer(t.Context())

	// start with 1 planet from connectTestDB
	result, err := c.getMinefieldsForGame(t.Context(), g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	if err := c.SaveMinefield(t.Context(), &cs.Minefield{GameDBObject: cs.GameDBObject{GameID: g.ID}, MapObject: cs.MapObject{PlayerNum: player.Num}}); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	result, err = c.getMinefieldsForGame(t.Context(), g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestUpdateMinefield(t *testing.T) {
	c := connectTestDB()

	g, player := c.createTestGameWithPlayer(t.Context())
	minefield := &cs.Minefield{GameDBObject: cs.GameDBObject{GameID: g.ID}, MapObject: cs.MapObject{PlayerNum: player.Num}}
	if err := c.SaveMinefield(t.Context(), minefield); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	minefield.Name = "Test2"
	if err := c.SaveMinefield(t.Context(), minefield); err != nil {
		t.Errorf("update planet %s", err)
		return
	}

	updated, err := c.GetMinefield(t.Context(), minefield.ID)

	if err != nil {
		t.Errorf("get planet %s", err)
		return
	}

	assert.Equal(t, minefield.Name, updated.Name)

}
