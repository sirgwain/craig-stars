package dbsqlx

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/game"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreatePlanetIntel(t *testing.T) {

	type args struct {
		c           *client
		planetIntel *game.PlanetIntel
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &game.PlanetIntel{MapObjectIntel: game.MapObjectIntel{Intel: game.Intel{Name: "name"}}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := *tt.args.planetIntel
			_, player := tt.args.c.createTestGameWithPlayer()
			tt.args.planetIntel.PlayerID = player.ID
			err := tt.args.c.CreatePlanetIntel(tt.args.planetIntel)

			// id is automatically added
			want.PlayerID = player.ID
			want.ID = tt.args.planetIntel.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePlanetIntel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.planetIntel, &want) {
				t.Errorf("CreatePlanetIntel() = \n%v, want \n%v", tt.args.planetIntel, want)
			}
		})
	}
}

func TestGetPlanetIntel(t *testing.T) {
	c := connectTestDB()
	_, player := c.createTestGameWithPlayer()
	planetIntel := game.PlanetIntel{MapObjectIntel: game.MapObjectIntel{Type: game.MapObjectTypePlanet, Intel: game.Intel{Name: "name"}}}
	planetIntel.PlayerID = player.ID
	if err := c.CreatePlanetIntel(&planetIntel); err != nil {
		t.Errorf("failed to create planetIntel %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *game.PlanetIntel
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got planetIntel", args{id: planetIntel.ID}, &planetIntel, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetPlanetIntel(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPlanetIntel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("GetPlanetIntel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPlanetIntels(t *testing.T) {
	c := connectTestDB()
	_, player := c.createTestGameWithPlayer()

	// start with 1 planetIntel from connectTestDB
	result, err := c.GetPlanetIntelsForPlayer(player.ID)
	assert.Nil(t, err)
	assert.Equal(t, []game.PlanetIntel{}, result)

	planetIntel := game.PlanetIntel{MapObjectIntel: game.MapObjectIntel{Intel: game.Intel{PlayerID: player.ID, Name: "name"}}}
	if err := c.CreatePlanetIntel(&planetIntel); err != nil {
		t.Errorf("failed to create planetIntel %s", err)
		return
	}

	result, err = c.GetPlanetIntelsForPlayer(player.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestDeletePlanetIntels(t *testing.T) {
	c := connectTestDB()
	_, player := c.createTestGameWithPlayer()

	result, err := c.GetPlanetIntelsForPlayer(player.ID)
	assert.Nil(t, err)
	assert.Equal(t, []game.PlanetIntel{}, result)

	planetIntel := game.PlanetIntel{MapObjectIntel: game.MapObjectIntel{Intel: game.Intel{PlayerID: player.ID, Name: "name"}}}
	if err := c.CreatePlanetIntel(&planetIntel); err != nil {
		t.Errorf("failed to create planetIntel %s", err)
		return
	}

	// should have our planetIntel in the db
	result, err = c.GetPlanetIntelsForPlayer(player.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	if err := c.DeletePlanetIntel(planetIntel.ID); err != nil {
		t.Errorf("failed to delete planetIntel %s", err)
		return
	}

	// should be no planetIntels left in db
	result, err = c.GetPlanetIntelsForPlayer(player.ID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}
