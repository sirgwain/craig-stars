package dbsqlx

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/game"
	"github.com/stretchr/testify/assert"
)

func TestCreateGame(t *testing.T) {

	type args struct {
		c    *client
		game *game.Game
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &game.Game{HostID: 1, Name: "test"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := *tt.args.game
			err := tt.args.c.CreateGame(tt.args.game)

			// id is automatically added
			want.ID = tt.args.game.ID
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.game, &want) {
				t.Errorf("CreateGame() = \n%v, want \n%v", tt.args.game, want)
			}
		})
	}
}

func TestUpdateGame(t *testing.T) {
	c := connectTestDB()
	game := game.Game{HostID: 1, Name: "Test"}
	if err := c.CreateGame(&game); err != nil {
		t.Errorf("failed to create game %s", err)
		return
	}

	game.Name = "Test2"
	if err := c.UpdateGame(&game); err != nil {
		t.Errorf("failed to update game %s", err)
		return
	}

	updated, err := c.GetGame(int64(game.ID))

	if err != nil {
		t.Errorf("failed to get game %s", err)
		return
	}

	assert.Equal(t, game.Name, updated.Name)
	assert.Less(t, game.UpdatedAt, updated.UpdatedAt)

}

func TestGetGame(t *testing.T) {
	c := connectTestDB()
	g := game.Game{HostID: 1, Name: "Test"}
	if err := c.CreateGame(&g); err != nil {
		t.Errorf("failed to create game %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *game.Game
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got game", args{id: int64(g.ID)}, &g, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetGame(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				tt.want.UpdatedAt = got.UpdatedAt
				tt.want.CreatedAt = got.CreatedAt
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGames(t *testing.T) {
	c := connectTestDB()

	// start with 1 game from connectTestDB
	result, err := c.GetGames()
	assert.Nil(t, err)
	assert.Equal(t, []game.Game{}, result)

	game := game.Game{HostID: 1, Name: "Test"}
	if err := c.CreateGame(&game); err != nil {
		t.Errorf("failed to create game %s", err)
		return
	}

	result, err = c.GetGames()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestDeleteGames(t *testing.T) {
	c := connectTestDB()

	result, err := c.GetGames()
	assert.Nil(t, err)
	assert.Equal(t, []game.Game{}, result)

	game := game.Game{HostID: 1, Name: "Test"}
	if err := c.CreateGame(&game); err != nil {
		t.Errorf("failed to create game %s", err)
		return
	}

	// should have our game in the db
	result, err = c.GetGames()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	if err := c.DeleteGame(int64(game.ID)); err != nil {
		t.Errorf("failed to delete game %s", err)
		return
	}

	// should be no games left in db
	result, err = c.GetGames()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestUpdateFullGame(t *testing.T) {
	c := connectTestDB()
	fg := c.createTestFullGame()

	if err := c.UpdateFullGame(fg); err != nil {
		t.Errorf("failed to update full game %s", err)
		return
	}

	updated, err := c.GetFullGame(fg.ID)

	if err != nil {
		t.Errorf("failed to get full game %s", err)
		return
	}

	assert.Equal(t, fg.Name, updated.Name)
	assert.Less(t, fg.UpdatedAt, updated.UpdatedAt)

}