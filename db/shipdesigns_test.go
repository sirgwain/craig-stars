package db

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateShipDesign(t *testing.T) {

	type args struct {
		c          *client
		shipDesign *cs.ShipDesign
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.ShipDesign{Name: "name", Num: 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := *tt.args.shipDesign
			_, player := tt.args.c.createTestGameWithPlayer()
			tt.args.shipDesign.GameID = player.GameID
			tt.args.shipDesign.PlayerNum = player.Num
			err := tt.args.c.CreateShipDesign(tt.args.shipDesign)

			// id is automatically added
			want.GameID = player.GameID
			want.PlayerNum = player.Num
			want.ID = tt.args.shipDesign.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateShipDesign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.shipDesign, &want) {
				t.Errorf("CreateShipDesign() = \n%v, want \n%v", tt.args.shipDesign, want)
			}
		})
	}
}

func TestGetShipDesign(t *testing.T) {
	rules := cs.NewRules()
	c := connectTestDB()
	_, player := c.createTestGameWithPlayer()
	shipDesign := cs.NewShipDesign(player, 1).WithHull(cs.Scout.Name).WithSpec(&rules, player)
	if err := c.CreateShipDesign(shipDesign); err != nil {
		t.Errorf("create shipDesign %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.ShipDesign
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got shipDesign", args{id: shipDesign.ID}, shipDesign, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetShipDesign(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShipDesign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("GetShipDesign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetShipDesigns(t *testing.T) {
	c := connectTestDB()
	_, player := c.createTestGameWithPlayer()

	// start with 1 shipDesign from connectTestDB
	result, err := c.GetShipDesignsForPlayer(player.GameID, player.Num)
	assert.Nil(t, err)
	assert.Equal(t, []*cs.ShipDesign{}, result)

	shipDesign := cs.ShipDesign{GameID: player.GameID, Num: 1, PlayerNum: player.Num, Name: "name"}
	if err := c.CreateShipDesign(&shipDesign); err != nil {
		t.Errorf("create shipDesign %s", err)
		return
	}

	result, err = c.GetShipDesignsForPlayer(player.GameID, player.Num)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestDeleteShipDesigns(t *testing.T) {
	c := connectTestDB()
	_, player := c.createTestGameWithPlayer()

	result, err := c.GetShipDesignsForPlayer(player.GameID, player.Num)
	assert.Nil(t, err)
	assert.Equal(t, []*cs.ShipDesign{}, result)

	shipDesign := cs.ShipDesign{GameID: player.GameID, Num: 1, PlayerNum: player.Num, Name: "name"}
	if err := c.CreateShipDesign(&shipDesign); err != nil {
		t.Errorf("create shipDesign %s", err)
		return
	}

	// should have our shipDesign in the db
	result, err = c.GetShipDesignsForPlayer(player.GameID, player.Num)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	if err := c.DeleteShipDesign(shipDesign.ID); err != nil {
		t.Errorf("delete shipDesign %s", err)
		return
	}

	// should be no shipDesigns left in db
	result, err = c.GetShipDesignsForPlayer(player.GameID, player.Num)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}
