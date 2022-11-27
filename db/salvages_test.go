package db

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/game"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateSalvage(t *testing.T) {
	type args struct {
		c       *client
		salvage *game.Salvage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &game.Salvage{
			MapObject: game.MapObject{GameID: 1, Name: "test"},
		},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			g, player := tt.args.c.createTestGameWithPlayer()
			tt.args.salvage.GameID = g.ID
			tt.args.salvage.PlayerID = player.ID

			want := *tt.args.salvage
			err := tt.args.c.createSalvage(tt.args.salvage, tt.args.c.db)

			// id is automatically added
			want.ID = tt.args.salvage.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSalvage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.salvage, &want) {
				t.Errorf("CreateSalvage() = \n%v, want \n%v", tt.args.salvage, want)
			}
		})
	}
}

func TestGetSalvage(t *testing.T) {
	c := connectTestDB()
	g, player := c.createTestGameWithPlayer()
	design := game.NewShipDesign(player).WithHull(game.Scout.Name)
	c.createTestShipDesign(player, design)

	salvage := game.Salvage{
		MapObject: game.MapObject{GameID: g.ID, PlayerID: player.ID, Name: "name", Type: game.MapObjectTypeSalvage},
	}
	if err := c.createSalvage(&salvage, c.db); err != nil {
		t.Errorf("create salvage %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *game.Salvage
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got salvage", args{id: salvage.ID}, &salvage, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetSalvage(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSalvage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("GetSalvage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSalvages(t *testing.T) {
	c := connectTestDB()
	g, player := c.createTestGameWithPlayer()

	// start with 1 planet from connectTestDB
	result, err := c.getSalvagesForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, []*game.Salvage{}, result)

	salvage := game.Salvage{MapObject: game.MapObject{GameID: g.ID, PlayerID: player.ID}}
	if err := c.createSalvage(&salvage, c.db); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	result, err = c.getSalvagesForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestUpdateSalvage(t *testing.T) {
	c := connectTestDB()
	g, player := c.createTestGameWithPlayer()
	planet := game.Salvage{MapObject: game.MapObject{GameID: g.ID, PlayerID: player.ID}}
	if err := c.createSalvage(&planet, c.db); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	planet.Name = "Test2"
	if err := c.UpdateSalvage(&planet); err != nil {
		t.Errorf("update planet %s", err)
		return
	}

	updated, err := c.GetSalvage(planet.ID)

	if err != nil {
		t.Errorf("get planet %s", err)
		return
	}

	assert.Equal(t, planet.Name, updated.Name)
	assert.Less(t, planet.UpdatedAt, updated.UpdatedAt)

}
