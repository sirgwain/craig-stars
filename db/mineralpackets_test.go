//go:build !wasi && !wasm

package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestSaveMineralPacket(t *testing.T) {
	type args struct {
		c             *client
		mineralPacket *cs.MineralPacket
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.MineralPacket{
			GameDBObject: cs.GameDBObject{GameID: 1},
			MapObject:    cs.MapObject{Type: cs.MapObjectTypeMineralPacket, Name: "test"}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			g, player := tt.args.c.createTestGameWithPlayer(t.Context())
			tt.args.mineralPacket.GameID = g.ID
			tt.args.mineralPacket.PlayerNum = player.Num

			want := *tt.args.mineralPacket
			err := tt.args.c.SaveMineralPacket(t.Context(), tt.args.mineralPacket)

			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("SaveMineralPacket() did not return error when expected")
				} else {
					t.Fatalf("SaveMineralPacket() errored unexpectedly; err = \n%v", err)
				}
			}

			got := tt.args.mineralPacket
			// DBObject is returned
			want.GameDBObject = got.GameDBObject
			test.CompareAsJSON(t, got, &want)
		})
	}
}

func TestGetMineralPacket(t *testing.T) {
	c := connectTestDB()

	g, player := c.createTestGameWithPlayer(t.Context())

	mineralPacket := &cs.MineralPacket{
		GameDBObject: cs.GameDBObject{GameID: g.ID},
		MapObject:    cs.MapObject{PlayerNum: player.Num, Name: "name", Type: cs.MapObjectTypeMineralPacket},
	}
	if err := c.SaveMineralPacket(t.Context(), mineralPacket); err != nil {
		t.Errorf("create mineralPacket %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.MineralPacket
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got mineralPacket", args{id: mineralPacket.ID}, mineralPacket, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetMineralPacket(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetMineralPacket() did not return error when expected")
				} else {
					t.Fatalf("GetMineralPacket() errored unexpectedly; err = \n%v", err)
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

func TestGetMineralPackets(t *testing.T) {
	c := connectTestDB()

	g, player := c.createTestGameWithPlayer(t.Context())

	// start with 1 planet from connectTestDB
	result, err := c.getMineralPacketsForGame(t.Context(), g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	if err := c.SaveMineralPacket(t.Context(), &cs.MineralPacket{GameDBObject: cs.GameDBObject{GameID: g.ID}, MapObject: cs.MapObject{PlayerNum: player.Num}}); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	result, err = c.getMineralPacketsForGame(t.Context(), g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestUpdateMineralPacket(t *testing.T) {
	c := connectTestDB()

	g, player := c.createTestGameWithPlayer(t.Context())
	mineralPacket := &cs.MineralPacket{GameDBObject: cs.GameDBObject{GameID: g.ID}, MapObject: cs.MapObject{PlayerNum: player.Num}}
	if err := c.SaveMineralPacket(t.Context(), mineralPacket); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	mineralPacket.Name = "Test2"
	if err := c.SaveMineralPacket(t.Context(), mineralPacket); err != nil {
		t.Errorf("update planet %s", err)
		return
	}

	updated, err := c.GetMineralPacket(t.Context(), mineralPacket.ID)

	if err != nil {
		t.Errorf("get planet %s", err)
		return
	}

	assert.Equal(t, mineralPacket.Name, updated.Name)

}
