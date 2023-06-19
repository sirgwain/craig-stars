package dbsqlx

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/game"
	"github.com/stretchr/testify/assert"
)

func TestCreatePlayerMessage(t *testing.T) {

	type args struct {
		c             *client
		playerMessage *game.PlayerMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &game.PlayerMessage{Type: game.PlayerMessageInfo, Text: "message info"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := *tt.args.playerMessage
			_, player := tt.args.c.createTestGameWithPlayer()
			tt.args.playerMessage.PlayerID = player.ID
			err := tt.args.c.CreatePlayerMessage(tt.args.playerMessage, tt.args.c.db)
			
			// id is automatically added
			want.PlayerID = player.ID
			want.ID = tt.args.playerMessage.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePlayerMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.playerMessage, &want) {
				t.Errorf("CreatePlayerMessage() = \n%v, want \n%v", tt.args.playerMessage, want)
			}
		})
	}
}

func TestGetPlayerMessage(t *testing.T) {
	c := connectTestDB()
	_, player := c.createTestGameWithPlayer()
	playerMessage := game.PlayerMessage{Type: game.PlayerMessageInfo, Text: "message"}
	playerMessage.PlayerID = player.ID
	if err := c.CreatePlayerMessage(&playerMessage, c.db); err != nil {
		t.Errorf("failed to create playerMessage %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *game.PlayerMessage
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got playerMessage", args{id: int64(playerMessage.ID)}, &playerMessage, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetPlayerMessage(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPlayerMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPlayerMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPlayerMessages(t *testing.T) {
	c := connectTestDB()
	_, player := c.createTestGameWithPlayer()

	// start with 1 playerMessage from connectTestDB
	result, err := c.GetPlayerMessages()
	assert.Nil(t, err)
	assert.Equal(t, []game.PlayerMessage{}, result)

	playerMessage := game.PlayerMessage{PlayerID: player.ID, Type: game.PlayerMessageInfo, Text: "Test"}
	if err := c.CreatePlayerMessage(&playerMessage, c.db); err != nil {
		t.Errorf("failed to create playerMessage %s", err)
		return
	}

	result, err = c.GetPlayerMessages()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestDeletePlayerMessages(t *testing.T) {
	c := connectTestDB()
	_, player := c.createTestGameWithPlayer()

	result, err := c.GetPlayerMessages()
	assert.Nil(t, err)
	assert.Equal(t, []game.PlayerMessage{}, result)

	playerMessage := game.PlayerMessage{PlayerID: player.ID, Type: game.PlayerMessageInfo, Text: "Test"}
	if err := c.CreatePlayerMessage(&playerMessage, c.db); err != nil {
		t.Errorf("failed to create playerMessage %s", err)
		return
	}

	// should have our playerMessage in the db
	result, err = c.GetPlayerMessages()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	if err := c.DeletePlayerMessage(playerMessage.ID); err != nil {
		t.Errorf("failed to delete playerMessage %s", err)
		return
	}

	// should be no playerMessages left in db
	result, err = c.GetPlayerMessages()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}
