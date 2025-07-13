package db

import (
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
			_, player := tt.args.c.createTestGameWithPlayer(t.Context())
			tt.args.shipDesign.GameID = player.GameID
			tt.args.shipDesign.PlayerNum = player.Num
			err := tt.args.c.SaveShipDesign(t.Context(), tt.args.shipDesign)

			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("CreateShipDesign() did not return error when expected")
				} else {
					t.Fatalf("CreateShipDesign() errored unexpectedly; err = \n%v", err)
				}
			}

			got := tt.args.shipDesign
			// DBObject is returned
			want.PlayerNum = got.PlayerNum
			want.GameDBObject = got.GameDBObject
			test.CompareAsJSON(t, got, want)
		})
	}
}

func TestGetShipDesign(t *testing.T) {
	rules := cs.NewRules()
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game, player := c.createTestGameWithPlayer(t.Context())
	shipDesign := cs.NewShipDesign(player.Num, 1).WithHull(cs.Scout.Name).WithSpec(&rules, player)
	shipDesign.GameID = game.ID
	if err := c.SaveShipDesign(t.Context(), shipDesign); err != nil {
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
			got, err := c.GetShipDesign(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetShipDesign() did not return error when expected")
				} else {
					t.Fatalf("GetShipDesign() errored unexpectedly; err = \n%v", err)
				}
			}
			if got != nil {
				tt.want.GameDBObject = got.GameDBObject
			}

			test.CompareAsJSON(t, got, tt.want)
		})
	}
}

func TestGetShipDesigns(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game, player := c.createTestGameWithPlayer(t.Context())

	// start with 1 shipDesign from connectTestDB
	result, err := c.GetShipDesignsForPlayer(t.Context(), player.GameID, player.Num)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	if err := c.SaveShipDesign(t.Context(), &cs.ShipDesign{GameDBObject: cs.GameDBObject{GameID: game.ID}, Num: 1, PlayerNum: player.Num, Name: "name"}); err != nil {
		t.Errorf("create shipDesign %s", err)
		return
	}

	result, err = c.GetShipDesignsForPlayer(t.Context(), player.GameID, player.Num)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestDeleteShipDesigns(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game, player := c.createTestGameWithPlayer(t.Context())

	result, err := c.GetShipDesignsForPlayer(t.Context(), player.GameID, player.Num)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	design := &cs.ShipDesign{GameDBObject: cs.GameDBObject{GameID: game.ID}, Num: 1, PlayerNum: player.Num, Name: "name"}
	if err := c.SaveShipDesign(t.Context(), design); err != nil {
		t.Errorf("create shipDesign %s", err)
		return
	}

	// should have our shipDesign in the db
	result, err = c.GetShipDesignsForPlayer(t.Context(), player.GameID, player.Num)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	if err := c.DeleteShipDesign(t.Context(), design.ID); err != nil {
		t.Errorf("delete shipDesign %s", err)
		return
	}

	// should be no shipDesigns left in db
	result, err = c.GetShipDesignsForPlayer(t.Context(), player.GameID, player.Num)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}
