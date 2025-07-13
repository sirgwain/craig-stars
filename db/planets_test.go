package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestSavePlanet(t *testing.T) {
	type args struct {
		c      *client
		planet *cs.Planet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.Planet{
			GameDBObject: cs.GameDBObject{GameID: 1},
			MapObject:    cs.MapObject{Type: cs.MapObjectTypePlanet, Name: "test"}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			game := tt.args.c.createTestGame(t.Context())
			tt.args.planet.GameID = game.ID

			want := *tt.args.planet
			err := tt.args.c.SavePlanet(t.Context(), tt.args.planet)

			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("CreatePlanet() did not return error when expected")
				} else {
					t.Fatalf("CreatePlanet() errored unexpectedly; err = \n%v", err)
				}
			}

			got := tt.args.planet
			// DBObject is returned
			want.GameDBObject = got.GameDBObject
			test.CompareAsJSON(t, got, want)
		})
	}
}

func TestGetPlanets(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame(t.Context())

	// start with 1 planet from connectTestDB
	result, err := c.GetPlanetsForGame(t.Context(), game.ID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	if err := c.SavePlanet(t.Context(), &cs.Planet{GameDBObject: cs.GameDBObject{GameID: game.ID}, MapObject: cs.MapObject{}}); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	result, err = c.GetPlanetsForGame(t.Context(), game.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestGetPlanet(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame(t.Context())
	planet := &cs.Planet{GameDBObject: cs.GameDBObject{GameID: game.ID}, MapObject: cs.MapObject{Name: "name", Type: cs.MapObjectTypePlanet}}
	if err := c.SavePlanet(t.Context(), planet); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.Planet
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got planet", args{id: planet.ID}, planet, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetPlanet(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetPlanet() did not return error when expected")
				} else {
					t.Fatalf("GetPlanet() errored unexpectedly; err = \n%v", err)
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

func TestUpdatePlanet(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame(t.Context())
	planet := &cs.Planet{GameDBObject: cs.GameDBObject{GameID: game.ID}, MapObject: cs.MapObject{}}
	if err := c.SavePlanet(t.Context(), planet); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	planet.Name = "Test2"
	if err := c.SavePlanet(t.Context(), planet); err != nil {
		t.Errorf("update planet %s", err)
		return
	}

	updated, err := c.GetPlanet(t.Context(), planet.ID)

	if err != nil {
		t.Errorf("get planet %s", err)
		return
	}

	assert.Equal(t, planet.Name, updated.Name)

}

func TestGetPlanetByNum(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	g, player := c.createTestGameWithPlayer(t.Context())
	planet1 := &cs.Planet{GameDBObject: cs.GameDBObject{GameID: g.ID}, MapObject: cs.MapObject{Name: "name", Num: 1, Type: cs.MapObjectTypePlanet}}
	if err := c.SavePlanet(t.Context(), planet1); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	planet2 := &cs.Planet{
		GameDBObject: cs.GameDBObject{GameID: g.ID},
		MapObject:    cs.MapObject{Name: "name", PlayerNum: player.Num, Num: 2, Type: cs.MapObjectTypePlanet},
	}
	if err := c.SavePlanet(t.Context(), planet2); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	design := c.createTestShipDesign(t.Context(), player, cs.NewShipDesign(player.Num, 1).WithHull(cs.SpaceStation.Name))

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
		PlanetNum: planet2.Num,
	}
	if err := c.SaveFleet(t.Context(), fleet); err != nil {
		t.Errorf("create fleet %s", err)
		return
	}
	planet2.Starbase = fleet

	type args struct {
		gameID int64
		num    int
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.Planet
		wantErr bool
	}{
		{"No results", args{gameID: 0, num: 0}, nil, false},
		{"Got planet Without Starbase", args{gameID: planet1.GameID, num: planet1.Num}, planet1, false},
		{"Got planet With Starbase", args{gameID: planet2.GameID, num: planet2.Num}, planet2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetPlanetByNum(t.Context(), tt.args.gameID, tt.args.num)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetPlanet() did not return error when expected")
				} else {
					t.Fatalf("GetPlanet() errored unexpectedly; err = \n%v", err)
				}
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}

			test.CompareAsJSON(t, got, tt.want)

			if tt.want != nil && tt.want.Starbase != nil {
				if got.Starbase != nil {
					tt.want.Starbase.UpdatedAt = got.Starbase.UpdatedAt
					tt.want.Starbase.CreatedAt = got.Starbase.CreatedAt
				}

				test.CompareAsJSON(t, got.Starbase, tt.want.Starbase)
			}
		})
	}
}
