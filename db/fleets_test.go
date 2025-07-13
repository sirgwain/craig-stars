package db

import (
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
			GameDBObject: cs.GameDBObject{GameID: 1},
			MapObject:    cs.MapObject{Type: cs.MapObjectTypeFleet, Name: "test"},
		},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			g, player := tt.args.c.createTestGameWithPlayer(t.Context())
			tt.args.fleet.GameID = g.ID
			tt.args.fleet.PlayerNum = player.Num

			want := *tt.args.fleet
			err := tt.args.c.SaveFleet(t.Context(), tt.args.fleet)

			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("CreateFleet() did not return error when expected")
				} else {
					t.Fatalf("CreateFleet() errored unexpectedly; err = \n%v", err)
				}
			}

			got := tt.args.fleet
			// DBObject is returned
			want.GameDBObject = got.GameDBObject
			test.CompareAsJSON(t, got, want)
		})
	}
}

func TestGetFleet(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	g, player := c.createTestGameWithPlayer(t.Context())

	design := c.createTestShipDesign(t.Context(), player, cs.NewShipDesign(player.Num, 1).WithHull(cs.Scout.Name))

	fleet := &cs.Fleet{
		GameDBObject: cs.GameDBObject{GameID: g.ID},
		MapObject:    cs.MapObject{PlayerNum: player.Num, Name: "name", Type: cs.MapObjectTypeFleet},
		Tokens: []cs.ShipToken{
			{Quantity: 1, DesignNum: design.Num},
		},
		FleetOrders: cs.FleetOrders{
			Waypoints: []cs.Waypoint{
				cs.NewPositionWaypoint(cs.Vector{X: 2, Y: 3}, 4),
			},
		},
	}
	if err := c.SaveFleet(t.Context(), fleet); err != nil {
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
		{"Got fleet", args{id: fleet.ID}, fleet, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetFleet(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetFleet() did not return error when expected")
				} else {
					t.Fatalf("GetFleet() errored unexpectedly; err = \n%v", err)
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

func TestGetFleets(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	g, player := c.createTestGameWithPlayer(t.Context())

	// start with 1 planet from connectTestDB
	result, err := c.GetFleetsForGame(t.Context(), g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	if err = c.SaveFleet(t.Context(), &cs.Fleet{GameDBObject: cs.GameDBObject{GameID: g.ID}, MapObject: cs.MapObject{PlayerNum: player.Num}}); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	result, err = c.GetFleetsForGame(t.Context(), g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestUpdateFleet(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	g, player := c.createTestGameWithPlayer(t.Context())
	fleet := &cs.Fleet{GameDBObject: cs.GameDBObject{GameID: g.ID}, MapObject: cs.MapObject{PlayerNum: player.Num}}

	if err := c.SaveFleet(t.Context(), fleet); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	fleet.Name = "Test2"
	if err := c.SaveFleet(t.Context(), fleet); err != nil {
		t.Errorf("update planet %s", err)
		return
	}

	updated, err := c.GetFleet(t.Context(), fleet.ID)

	if err != nil {
		t.Errorf("get planet %s", err)
		return
	}

	assert.Equal(t, fleet.Name, updated.Name)

}
