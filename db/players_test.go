//go:build !wasi && !wasm

package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestSavePlayer(t *testing.T) {
	type args struct {
		c      *client
		player *cs.Player
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.Player{UserID: 1, Name: "test"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			game := tt.args.c.createTestGame(t.Context())
			tt.args.player.GameID = game.ID

			want := *tt.args.player
			err := tt.args.c.SavePlayer(t.Context(), tt.args.player)

			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("CreatePlayer() did not return error when expected")
				} else {
					t.Fatalf("CreatePlayer() errored unexpectedly; err = \n%v", err)
				}
			}

			got := tt.args.player
			// DBObject is returned
			want.GameDBObject = got.GameDBObject
			test.CompareAsJSON(t, got, want)
		})
	}
}

func TestUpdatePlayer(t *testing.T) {
	c := connectTestDB()

	c.createTestGame(t.Context())
	player := &cs.Player{GameDBObject: cs.GameDBObject{GameID: 1}, UserID: 1, Name: "Test"}
	if err := c.SavePlayer(t.Context(), player); err != nil {
		t.Errorf("create player %s", err)
		return
	}

	player.Name = "Test2"
	player.Num = 1
	player.Messages = append(player.Messages, cs.PlayerMessage{Type: cs.PlayerMessageInfo, Text: "message1"})
	player.Messages = append(player.Messages, cs.PlayerMessage{Type: cs.PlayerMessageInfo, Text: "message2"})
	if err := c.SavePlayer(t.Context(), player); err != nil {
		t.Errorf("update player %s", err)
		return
	}

	updated, err := c.GetPlayer(t.Context(), player.ID)

	if err != nil {
		t.Errorf("get player %s", err)
		return
	}

	assert.Equal(t, player.Name, updated.Name)
	assert.Equal(t, player.Num, updated.Num)
	assert.Equal(t, 2, len(updated.Messages))

}

func TestGetPlayer(t *testing.T) {
	rules := cs.NewRules()
	c := connectTestDB()

	c.createTestGame(t.Context())
	player := &cs.Player{GameDBObject: cs.GameDBObject{GameID: 1}, UserID: 1, Name: "Test", Race: *cs.NewRace().WithSpec(&rules)}
	if err := c.SavePlayer(t.Context(), player); err != nil {
		t.Errorf("create player %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.Player
		wantErr bool
	}{
		// {"No results", args{id: 0}, nil, false},
		{"Got player", args{id: player.ID}, player, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetPlayer(t.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetPlayer() did not return error when expected")
				} else {
					t.Fatalf("GetPlayer() errored unexpectedly; err = \n%v", err)
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

func Test_GetPlayerForGame(t *testing.T) {
	rules := cs.NewRules()
	c := connectTestDB()

	game := c.createTestGame(t.Context())
	player := &cs.Player{GameDBObject: cs.GameDBObject{GameID: 1}, UserID: 1, Name: "Test", Race: *cs.NewRace().WithSpec(&rules)}
	if err := c.SavePlayer(t.Context(), player); err != nil {
		t.Errorf("create player %s", err)
		return
	}

	// verify it works with no designs
	_, err := c.GetPlayerForGame(t.Context(), game.ID, player.Num)
	if err != nil {
		t.Errorf("GetPlayerForGame %s", err)
		return
	}

	// create a couple designs and join again
	shipDesign1 := &cs.ShipDesign{Num: 1, PlayerNum: player.Num, Name: "name"}
	shipDesign1.GameID = game.ID
	if err := c.SaveShipDesign(t.Context(), shipDesign1); err != nil {
		t.Errorf("create shipDesign %s", err)
		return
	}

	shipDesign2 := &cs.ShipDesign{Num: 2, PlayerNum: player.Num, Name: "name2"}
	shipDesign2.GameID = game.ID
	if err := c.SaveShipDesign(t.Context(), shipDesign2); err != nil {
		t.Errorf("create shipDesign %s", err)
		return
	}

	got, err := c.GetPlayerForGame(t.Context(), game.ID, player.Num)
	if err != nil {
		t.Errorf("GetPlayerForGame %s", err)
		return
	}

	// we expect to have a player with designs
	player.Designs = append(player.Designs, shipDesign1, shipDesign2)

	// clear out the incoming updated/created timestamps, we won't have those

	player.GameDBObject = got.GameDBObject
	for i := range got.Designs {
		player.Designs[i].GameDBObject = got.Designs[i].GameDBObject
	}

	test.CompareAsJSON(t, got, player)
}

func TestGetPlayers(t *testing.T) {
	c := connectTestDB()

	c.createTestGame(t.Context())

	// start with 1 player from connectTestDB
	result, err := c.GetPlayers(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	player := &cs.Player{GameDBObject: cs.GameDBObject{GameID: 1}, UserID: 1, Name: "Test"}
	if err := c.SavePlayer(t.Context(), player); err != nil {
		t.Errorf("create player %s", err)
		return
	}

	result, err = c.GetPlayers(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestDeletePlayers(t *testing.T) {
	c := connectTestDB()

	c.createTestGame(t.Context())

	result, err := c.GetPlayers(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	player := &cs.Player{GameDBObject: cs.GameDBObject{GameID: 1}, UserID: 1, Name: "Test"}
	if err := c.SavePlayer(t.Context(), player); err != nil {
		t.Errorf("create player %s", err)
		return
	}

	// should have our player in the db
	result, err = c.GetPlayers(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	if err := c.DeletePlayer(t.Context(), player.ID); err != nil {
		t.Errorf("delete player %s", err)
		return
	}

	// should be no players left in db
	result, err = c.GetPlayers(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}
