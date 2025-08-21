//go:build !wasi && !wasm

package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestSaveSalvage(t *testing.T) {
	type args struct {
		c       *client
		salvage *cs.Salvage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.Salvage{
			GameDBObject: cs.GameDBObject{GameID: 1},
			MapObject:    cs.MapObject{Type: cs.MapObjectTypeSalvage, Name: "test"}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			g, player := tt.args.c.createTestGameWithPlayer(t.Context())
			tt.args.salvage.GameID = g.ID
			tt.args.salvage.PlayerNum = player.Num

			want := *tt.args.salvage
			err := tt.args.c.SaveSalvage(t.Context(), tt.args.salvage)

			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("SaveSalvage() did not return error when expected")
				} else {
					t.Fatalf("SaveSalvage() errored unexpectedly; err = \n%v", err)
				}
			}

			got := tt.args.salvage
			// DBObject is returned
			want.GameDBObject = got.GameDBObject
			test.CompareAsJSON(t, got, want)
		})
	}
}

func TestGetSalvage(t *testing.T) {
	c := connectTestDB()

	g, player := c.createTestGameWithPlayer(t.Context())

	salvage := &cs.Salvage{
		GameDBObject: cs.GameDBObject{GameID: g.ID},
		MapObject:    cs.MapObject{PlayerNum: player.Num, Name: "name", Type: cs.MapObjectTypeSalvage},
	}
	if err := c.SaveSalvage(t.Context(), salvage); err != nil {
		t.Errorf("create salvage %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.Salvage
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got salvage", args{id: salvage.ID}, salvage, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetSalvage(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetSalvage() did not return error when expected")
				} else {
					t.Fatalf("GetSalvage() errored unexpectedly; err = \n%v", err)
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

func TestGetSalvages(t *testing.T) {
	c := connectTestDB()

	g, player := c.createTestGameWithPlayer(t.Context())

	// start with 1 salvage from connectTestDB
	result, err := c.GetSalvagesForGame(t.Context(), g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	if err := c.SaveSalvage(t.Context(), &cs.Salvage{GameDBObject: cs.GameDBObject{GameID: g.ID}, MapObject: cs.MapObject{PlayerNum: player.Num}}); err != nil {
		t.Errorf("create salvage %s", err)
		return
	}

	result, err = c.GetSalvagesForGame(t.Context(), g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestUpdateSalvage(t *testing.T) {
	c := connectTestDB()

	g, player := c.createTestGameWithPlayer(t.Context())
	salvage := &cs.Salvage{GameDBObject: cs.GameDBObject{GameID: g.ID}, MapObject: cs.MapObject{PlayerNum: player.Num}}
	if err := c.SaveSalvage(t.Context(), salvage); err != nil {
		t.Errorf("create salvage %s", err)
		return
	}

	salvage.Name = "Test2"
	if err := c.SaveSalvage(t.Context(), salvage); err != nil {
		t.Errorf("update salvage %s", err)
		return
	}

	updated, err := c.GetSalvage(t.Context(), salvage.ID)

	if err != nil {
		t.Errorf("get salvage %s", err)
		return
	}

	assert.Equal(t, salvage.Name, updated.Name)

}
