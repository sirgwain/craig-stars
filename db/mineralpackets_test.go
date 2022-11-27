package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/game"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateMineralPacket(t *testing.T) {
	type args struct {
		c             *client
		mineralPacket *game.MineralPacket
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &game.MineralPacket{
			MapObject: game.MapObject{GameID: 1, Name: "test"},
		},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			g, player := tt.args.c.createTestGameWithPlayer()
			tt.args.mineralPacket.GameID = g.ID
			tt.args.mineralPacket.PlayerID = player.ID

			want := *tt.args.mineralPacket
			err := tt.args.c.createMineralPacket(tt.args.mineralPacket, tt.args.c.db)

			// id is automatically added
			want.ID = tt.args.mineralPacket.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMineralPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !test.CompareAsJSON(t, tt.args.mineralPacket, &want) {
				t.Errorf("CreateMineralPacket() = \n%v, want \n%v", tt.args.mineralPacket, want)
			}
		})
	}
}

func TestGetMineralPacket(t *testing.T) {
	c := connectTestDB()
	g, player := c.createTestGameWithPlayer()
	design := game.NewShipDesign(player).WithHull(game.Scout.Name)
	c.createTestShipDesign(player, design)

	mineralPacket := game.MineralPacket{
		MapObject: game.MapObject{GameID: g.ID, PlayerID: player.ID, Name: "name", Type: game.MapObjectTypeMineralPacket},
	}
	if err := c.createMineralPacket(&mineralPacket, c.db); err != nil {
		t.Errorf("create mineralPacket %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *game.MineralPacket
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got mineralPacket", args{id: mineralPacket.ID}, &mineralPacket, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetMineralPacket(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMineralPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("GetMineralPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMineralPackets(t *testing.T) {
	c := connectTestDB()
	g, player := c.createTestGameWithPlayer()

	// start with 1 planet from connectTestDB
	result, err := c.getMineralPacketsForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, []*game.MineralPacket{}, result)

	mineralPacket := game.MineralPacket{MapObject: game.MapObject{GameID: g.ID, PlayerID: player.ID}}
	if err := c.createMineralPacket(&mineralPacket, c.db); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	result, err = c.getMineralPacketsForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestUpdateMineralPacket(t *testing.T) {
	c := connectTestDB()
	g, player := c.createTestGameWithPlayer()
	planet := game.MineralPacket{MapObject: game.MapObject{GameID: g.ID, PlayerID: player.ID}}
	if err := c.createMineralPacket(&planet, c.db); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	planet.Name = "Test2"
	if err := c.UpdateMineralPacket(&planet); err != nil {
		t.Errorf("update planet %s", err)
		return
	}

	updated, err := c.GetMineralPacket(planet.ID)

	if err != nil {
		t.Errorf("get planet %s", err)
		return
	}

	assert.Equal(t, planet.Name, updated.Name)
	assert.Less(t, planet.UpdatedAt, updated.UpdatedAt)

}
