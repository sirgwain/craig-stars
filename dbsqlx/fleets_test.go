package dbsqlx

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/game"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateFleet(t *testing.T) {
	type args struct {
		c      *client
		planet *game.Fleet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &game.Fleet{
			MapObject: game.MapObject{GameID: 1, Name: "test"},
		},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			g := tt.args.c.createTestGame()
			tt.args.planet.GameID = g.ID

			want := *tt.args.planet
			err := tt.args.c.createFleet(tt.args.planet, tt.args.c.db)

			// id is automatically added
			want.ID = tt.args.planet.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFleet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.planet, &want) {
				t.Errorf("CreateFleet() = \n%v, want \n%v", tt.args.planet, want)
			}
		})
	}
}

func TestGetFleet(t *testing.T) {
	c := connectTestDB()
	g, player := c.createTestGameWithPlayer()
	design := game.NewShipDesign(player).WithHull(game.Scout.Name)
	c.createTestShipDesign(player, design)

	fleet := game.Fleet{
		MapObject: game.MapObject{GameID: g.ID, PlayerID: player.ID, Name: "name", Type: game.MapObjectTypeFleet},
		Tokens: []game.ShipToken{
			{Quantity: 1, Design: design},
		},
		FleetOrders: game.FleetOrders{
			Waypoints: []game.Waypoint{
				game.NewPositionWaypoint(game.Vector{X: 2, Y: 3}, 4),
			},
		},
	}
	if err := c.createFleet(&fleet, c.db); err != nil {
		t.Errorf("failed to create fleet %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *game.Fleet
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got fleet", args{id: fleet.ID}, &fleet, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetFleet(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFleet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
				for i := range got.Tokens {
					tt.want.Tokens[i].CreatedAt = got.Tokens[i].CreatedAt
					tt.want.Tokens[i].UpdatedAt = got.Tokens[i].UpdatedAt
				}
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("GetFleet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFleets(t *testing.T) {
	c := connectTestDB()
	g := c.createTestGame()

	// start with 1 planet from connectTestDB
	result, err := c.getFleetsForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, []*game.Fleet{}, result)

	fleet := game.Fleet{MapObject: game.MapObject{GameID: g.ID}}
	if err := c.createFleet(&fleet, c.db); err != nil {
		t.Errorf("failed to create planet %s", err)
		return
	}

	result, err = c.getFleetsForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestUpdateFleet(t *testing.T) {
	c := connectTestDB()
	g := c.createTestGame()
	planet := game.Fleet{MapObject: game.MapObject{GameID: g.ID}}
	if err := c.createFleet(&planet, c.db); err != nil {
		t.Errorf("failed to create planet %s", err)
		return
	}

	planet.Name = "Test2"
	if err := c.UpdateFleet(&planet); err != nil {
		t.Errorf("failed to update planet %s", err)
		return
	}

	updated, err := c.GetFleet(planet.ID)

	if err != nil {
		t.Errorf("failed to get planet %s", err)
		return
	}

	assert.Equal(t, planet.Name, updated.Name)
	assert.Less(t, planet.UpdatedAt, updated.UpdatedAt)

}
