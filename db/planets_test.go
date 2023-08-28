package db

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreatePlanet(t *testing.T) {
	type args struct {
		c      *txClient
		planet *cs.Planet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.Planet{
			MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: 1}, Name: "test"}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			game := tt.args.c.createTestGame()
			tt.args.planet.GameID = game.ID

			want := *tt.args.planet
			err := tt.args.c.createPlanet(tt.args.planet)

			// id is automatically added
			want.ID = tt.args.planet.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePlanet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.planet, &want) {
				t.Errorf("CreatePlanet() = \n%v, want \n%v", tt.args.planet, want)
			}
		})
	}
}

func TestGetPlanets(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame()

	// start with 1 planet from connectTestDB
	result, err := c.getPlanetsForGame(game.ID)
	assert.Nil(t, err)
	assert.Equal(t, []*cs.Planet{}, result)

	planet := cs.Planet{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: game.ID}}}
	if err := c.createPlanet(&planet); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	result, err = c.getPlanetsForGame(game.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestGetPlanet(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame()
	planet := cs.Planet{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: game.ID}, Name: "name", Type: cs.MapObjectTypePlanet}}
	if err := c.createPlanet(&planet); err != nil {
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
		{"Got planet", args{id: planet.ID}, &planet, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetPlanet(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPlanet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("GetPlanet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdatePlanet(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := c.createTestGame()
	planet := cs.Planet{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: game.ID}}}
	if err := c.createPlanet(&planet); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	planet.Name = "Test2"
	if err := c.UpdatePlanet(&planet); err != nil {
		t.Errorf("update planet %s", err)
		return
	}

	updated, err := c.GetPlanet(planet.ID)

	if err != nil {
		t.Errorf("get planet %s", err)
		return
	}

	assert.Equal(t, planet.Name, updated.Name)
	assert.Less(t, planet.UpdatedAt, updated.UpdatedAt)

}

func TestGetPlanetByNum(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	g, player := c.createTestGameWithPlayer()
	planet1 := cs.Planet{MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: g.ID}, Name: "name", Num: 1, Type: cs.MapObjectTypePlanet}}
	if err := c.createPlanet(&planet1); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	planet2 := cs.Planet{
		MapObject: cs.MapObject{GameDBObject: cs.GameDBObject{GameID: g.ID}, Name: "name", PlayerNum: player.Num, Num: 2, Type: cs.MapObjectTypePlanet},
	}
	if err := c.createPlanet(&planet2); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	design := cs.NewShipDesign(player, 1).WithHull(cs.SpaceStation.Name)
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
		PlanetNum: planet2.Num,
	}
	if err := c.createFleet(&fleet); err != nil {
		t.Errorf("create fleet %s", err)
		return
	}
	planet2.Starbase = &fleet

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
		{"Got planet Without Starbase", args{gameID: planet1.GameID, num: planet1.Num}, &planet1, false},
		{"Got planet With Starbase", args{gameID: planet2.GameID, num: planet2.Num}, &planet2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetPlanetByNum(tt.args.gameID, tt.args.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPlanet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("GetPlanetByNum() = %v, want %v", got, tt.want)
			}

			if tt.want != nil && tt.want.Starbase != nil {
				if got.Starbase != nil {
					tt.want.Starbase.UpdatedAt = got.Starbase.UpdatedAt
					tt.want.Starbase.CreatedAt = got.Starbase.CreatedAt
				}
				if !test.CompareAsJSON(t, got.Starbase, tt.want.Starbase) {
					t.Errorf("GetPlanetByNum() Starbase = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
