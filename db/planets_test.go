package db

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/game"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreatePlanet(t *testing.T) {
	type args struct {
		c      *client
		planet *game.Planet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &game.Planet{
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
			err := tt.args.c.createPlanet(tt.args.planet, tt.args.c.db)

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
	g := c.createTestGame()

	// start with 1 planet from connectTestDB
	result, err := c.getPlanetsForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, []*game.Planet{}, result)

	planet := game.Planet{MapObject: game.MapObject{GameID: g.ID}}
	if err := c.createPlanet(&planet, c.db); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	result, err = c.getPlanetsForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestGetPlanet(t *testing.T) {
	c := connectTestDB()
	g := c.createTestGame()
	planet := game.Planet{MapObject: game.MapObject{GameID: g.ID, Name: "name", Type: game.MapObjectTypePlanet}}
	if err := c.createPlanet(&planet, c.db); err != nil {
		t.Errorf("create planet %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *game.Planet
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
	g := c.createTestGame()
	planet := game.Planet{MapObject: game.MapObject{GameID: g.ID}}
	if err := c.createPlanet(&planet, c.db); err != nil {
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
