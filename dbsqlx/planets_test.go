package dbsqlx

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/game"
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
	result, err := c.GetPlanetsForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, []*game.Planet{}, result)

	planet := game.Planet{MapObject: game.MapObject{GameID: g.ID}}
	if err := c.createPlanet(&planet, c.db); err != nil {
		t.Errorf("failed to create planet %s", err)
		return
	}

	result, err = c.GetPlanetsForGame(g.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestUpdatePlanet(t *testing.T) {
	c := connectTestDB()
	g := c.createTestGame()
	planet := game.Planet{MapObject: game.MapObject{GameID: g.ID}}
	if err := c.createPlanet(&planet, c.db); err != nil {
		t.Errorf("failed to create planet %s", err)
		return
	}

	planet.Name = "Test2"
	if err := c.UpdatePlanet(&planet); err != nil {
		t.Errorf("failed to update planet %s", err)
		return
	}

	updated, err := c.GetPlanet(int64(planet.ID))

	if err != nil {
		t.Errorf("failed to get planet %s", err)
		return
	}

	assert.Equal(t, planet.Name, updated.Name)
	assert.Less(t, planet.UpdatedAt, updated.UpdatedAt)

}
