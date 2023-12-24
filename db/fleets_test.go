package db

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateFleet(t *testing.T) {
	type args struct {
		c     *client
		fleet *cs.Fleet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.Fleet{
			MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: 1}, Name: "test"},
		},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			g, player := tt.args.c.createTestGameWithPlayer()
			tt.args.fleet.GameID = g.ID
			tt.args.fleet.PlayerNum = player.Num

			want := *tt.args.fleet
			err := tt.args.c.CreateFleet(tt.args.fleet)

			// id is automatically added
			want.ID = tt.args.fleet.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFleet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.fleet, &want) {
				t.Errorf("CreateFleet() = \n%v, want \n%v", tt.args.fleet, want)
			}
		})
	}
}

func TestGetFleet(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	g, player := c.createTestGameWithPlayer()
	design := cs.NewShipDesign(player, 1).WithHull(cs.Scout.Name)
	c.createTestShipDesign(player, design)

	fleet := cs.Fleet{
		MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: g.ID}, PlayerNum: player.Num, Name: "name", Type: cs.MapObjectTypeFleet},
		Tokens: []cs.ShipToken{
			{Quantity: 1, DesignNum: design.Num},
		},
		FleetOrders: cs.FleetOrders{
			Waypoints: []cs.Waypoint{
				cs.NewPositionWaypoint(cs.Vector{X: 2, Y: 3}, 4),
			},
		},
	}
	if err := c.CreateFleet(&fleet); err != nil {
		t.Errorf("create fleet %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.Fleet
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
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("GetFleet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFleets(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	g, player := c.createTestGameWithPlayer()

	// start with 1 planet from connectTestDB
	result, err := c.getFleetsForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, []*cs.Fleet{}, result)

	fleet := cs.Fleet{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: g.ID}, PlayerNum: player.Num}}
	if err := c.CreateFleet(&fleet); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	result, err = c.getFleetsForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestUpdateFleet(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	g, player := c.createTestGameWithPlayer()
	planet := cs.Fleet{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: g.ID}, PlayerNum: player.Num}}
	if err := c.CreateFleet(&planet); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	planet.Name = "Test2"
	if err := c.UpdateFleet(&planet); err != nil {
		t.Errorf("update planet %s", err)
		return
	}

	updated, err := c.GetFleet(planet.ID)

	if err != nil {
		t.Errorf("get planet %s", err)
		return
	}

	assert.Equal(t, planet.Name, updated.Name)
	assert.Less(t, planet.UpdatedAt, updated.UpdatedAt)

}
