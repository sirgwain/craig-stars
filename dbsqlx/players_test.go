package dbsqlx

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/game"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreatePlayer(t *testing.T) {
	type args struct {
		c      *client
		player *game.Player
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &game.Player{UserID: 1, Name: "test"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a test game
			g := tt.args.c.createTestGame()
			tt.args.player.GameID = g.ID

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
	player := game.Player{UserID: 1, GameID: 1, Name: "Test"}
	if err := c.CreatePlayer(&player); err != nil {
		t.Errorf("failed to create player %s", err)
		return
	}

	player.Name = "Test2"
	player.Num = 1
	if err := c.UpdatePlayer(&player); err != nil {
		t.Errorf("failed to update player %s", err)
		return
	}

	updated, err := c.GetPlayer(player.ID)

	if err != nil {
		t.Errorf("failed to get player %s", err)
		return
	}

	assert.Equal(t, player.Name, updated.Name)
	assert.Equal(t, player.Num, updated.Num)
	assert.Less(t, player.UpdatedAt, updated.UpdatedAt)

}

func TestGetPlayer(t *testing.T) {
	rules := game.NewRules()
	c := connectTestDB()
	c.createTestGame()
	player := game.Player{UserID: 1, GameID: 1, Name: "Test", Race: *game.NewRace().WithSpec(&rules)}
	if err := c.CreatePlayer(&player); err != nil {
		t.Errorf("failed to create player %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *game.Player
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

func TestGetPlayers(t *testing.T) {
	c := connectTestDB()
	c.createTestGame()

	// start with 1 player from connectTestDB
	result, err := c.GetPlayers()
	assert.Nil(t, err)
	assert.Equal(t, []game.Player{}, result)

	player := game.Player{UserID: 1, GameID: 1, Name: "Test"}
	if err := c.CreatePlayer(&player); err != nil {
		t.Errorf("failed to create player %s", err)
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
	assert.Equal(t, []game.Player{}, result)

	player := game.Player{UserID: 1, GameID: 1, Name: "Test"}
	if err := c.CreatePlayer(&player); err != nil {
		t.Errorf("failed to create player %s", err)
		return
	}

	// should have our player in the db
	result, err = c.GetPlayers()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	if err := c.DeletePlayer(player.ID); err != nil {
		t.Errorf("failed to delete player %s", err)
		return
	}

	// should be no players left in db
	result, err = c.GetPlayers()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestUpdateFullPlayer(t *testing.T) {
	c := connectTestDB()
	g := c.createTestGame()
	player := game.Player{UserID: 1, GameID: g.ID, Name: "Test"}
	if err := c.CreatePlayer(&player); err != nil {
		t.Errorf("failed to create player %s", err)
		return
	}

	player.Name = "Test2"
	player.Num = 1
	player.Messages = append(player.Messages, game.PlayerMessage{PlayerID: player.ID, Type: game.PlayerMessageInfo, Text: "message1"})
	player.Messages = append(player.Messages, game.PlayerMessage{PlayerID: player.ID, Type: game.PlayerMessageInfo, Text: "message2"})
	if err := c.updateFullPlayer(&player); err != nil {
		t.Errorf("failed to update player %s", err)
		return
	}

	updated, err := c.GetFullPlayerForGame(player.GameID, player.UserID)

	if err != nil {
		t.Errorf("failed to get player %s", err)
		return
	}

	assert.Equal(t, player.Name, updated.Name)
	assert.Equal(t, player.Num, updated.Num)
	assert.Less(t, player.UpdatedAt, updated.UpdatedAt)
	assert.Equal(t, 2, len(updated.Messages))

}
