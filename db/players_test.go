package db

import (
	"reflect"
	"testing"
	"time"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreatePlayer(t *testing.T) {
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
			game := tt.args.c.createTestGame()
			tt.args.player.GameID = game.ID

			want := *tt.args.player
			err := tt.args.c.CreatePlayer(tt.args.player)

			// id is automatically added
			want.ID = tt.args.player.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePlayer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.player, &want) {
				t.Errorf("CreatePlayer() = \n%v, want \n%v", tt.args.player, want)
			}
		})
	}
}

func TestUpdatePlayer(t *testing.T) {
	c := connectTestDB()
	c.createTestGame()
	player := cs.Player{GameDBObject: cs.GameDBObject{GameID: 1}, UserID: 1, Name: "Test"}
	if err := c.CreatePlayer(&player); err != nil {
		t.Errorf("create player %s", err)
		return
	}

	player.Name = "Test2"
	player.Num = 1
	if err := c.UpdatePlayer(&player); err != nil {
		t.Errorf("update player %s", err)
		return
	}

	updated, err := c.GetPlayer(player.ID)

	if err != nil {
		t.Errorf("get player %s", err)
		return
	}

	assert.Equal(t, player.Name, updated.Name)
	assert.Equal(t, player.Num, updated.Num)
	assert.Less(t, player.UpdatedAt, updated.UpdatedAt)

}

func TestGetPlayer(t *testing.T) {
	rules := cs.NewRules()
	c := connectTestDB()
	c.createTestGame()
	player := cs.Player{GameDBObject: cs.GameDBObject{GameID: 1}, UserID: 1, Name: "Test", Race: *cs.NewRace().WithSpec(&rules)}
	if err := c.CreatePlayer(&player); err != nil {
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
		{"Got player", args{id: player.ID}, &player, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetPlayer(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPlayer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("GetPlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPlayerWithDesigns(t *testing.T) {
	rules := cs.NewRules()
	c := connectTestDB()
	game := c.createTestGame()
	player := cs.Player{GameDBObject: cs.GameDBObject{GameID: 1}, UserID: 1, Num: 1, Name: "Test", Race: *cs.NewRace().WithSpec(&rules)}
	if err := c.CreatePlayer(&player); err != nil {
		t.Errorf("create player %s", err)
		return
	}

	// verify it works with no designs
	_, err := c.getPlayerWithDesigns("p.gameId = ?", game.ID)
	if err != nil {
		t.Errorf("getPlayerWithDesigns %s", err)
		return
	}

	// create a couple designs and join again
	shipDesign1 := cs.ShipDesign{Num: 1, PlayerNum: player.Num, Name: "name"}
	shipDesign1.GameID = game.ID
	if err := c.CreateShipDesign(&shipDesign1); err != nil {
		t.Errorf("create shipDesign %s", err)
		return
	}

	shipDesign2 := cs.ShipDesign{Num: 2, PlayerNum: player.Num, Name: "name2"}
	shipDesign2.GameID = game.ID
	if err := c.CreateShipDesign(&shipDesign2); err != nil {
		t.Errorf("create shipDesign %s", err)
		return
	}

	got, err := c.getPlayerWithDesigns("p.gameId = ?", game.ID)
	if err != nil {
		t.Errorf("getPlayerWithDesigns %s", err)
		return
	}

	// we expect to have a player with designs
	player.Designs = append(player.Designs, &shipDesign1, &shipDesign2)

	// clear out the incoming updated/created timestamps, we won't have those
	for i := range got {
		p := &got[i]
		p.CreatedAt = time.Time{}
		p.UpdatedAt = time.Time{}
		for _, design := range p.Designs {
			design.CreatedAt = time.Time{}
			design.UpdatedAt = time.Time{}
		}

	}
	if !test.CompareAsJSON(t, got, []*cs.Player{&player}) {
		t.Errorf("getPlayerWithDesigns() = %v, want %v", got, player)
	}
}

func TestGetPlayers(t *testing.T) {
	c := connectTestDB()
	c.createTestGame()

	// start with 1 player from connectTestDB
	result, err := c.GetPlayers()
	assert.Nil(t, err)
	assert.Equal(t, []cs.Player{}, result)

	player := cs.Player{GameDBObject: cs.GameDBObject{GameID: 1}, UserID: 1, Name: "Test"}
	if err := c.CreatePlayer(&player); err != nil {
		t.Errorf("create player %s", err)
		return
	}

	result, err = c.GetPlayers()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestDeletePlayers(t *testing.T) {
	c := connectTestDB()
	c.createTestGame()

	result, err := c.GetPlayers()
	assert.Nil(t, err)
	assert.Equal(t, []cs.Player{}, result)

	player := cs.Player{GameDBObject: cs.GameDBObject{GameID: 1}, UserID: 1, Name: "Test"}
	if err := c.CreatePlayer(&player); err != nil {
		t.Errorf("create player %s", err)
		return
	}

	// should have our player in the db
	result, err = c.GetPlayers()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	if err := c.DeletePlayer(player.ID); err != nil {
		t.Errorf("delete player %s", err)
		return
	}

	// should be no players left in db
	result, err = c.GetPlayers()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestUpdateFullPlayer(t *testing.T) {
	c := connectTestDB()
	game := c.createTestGame()
	player := cs.Player{GameDBObject: cs.GameDBObject{GameID: game.ID}, UserID: 1, Name: "Test"}
	if err := c.CreatePlayer(&player); err != nil {
		t.Errorf("create player %s", err)
		return
	}

	player.Name = "Test2"
	player.Num = 1
	player.Messages = append(player.Messages, cs.PlayerMessage{Type: cs.PlayerMessageInfo, Text: "message1"})
	player.Messages = append(player.Messages, cs.PlayerMessage{Type: cs.PlayerMessageInfo, Text: "message2"})
	if err := c.updateFullPlayer(&player); err != nil {
		t.Errorf("update player %s", err)
		return
	}

	updated, err := c.GetFullPlayerForGame(player.GameID, player.UserID)

	if err != nil {
		t.Errorf("get player %s", err)
		return
	}

	assert.Equal(t, player.Name, updated.Name)
	assert.Equal(t, player.Num, updated.Num)
	assert.Less(t, player.UpdatedAt, updated.UpdatedAt)
	assert.Equal(t, 2, len(updated.Messages))

}
