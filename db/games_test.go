package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateGame(t *testing.T) {

	type args struct {
		c    *client
		game *cs.Game
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create", args{connectTestDB(), &cs.Game{HostID: 1, Name: "test"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := *tt.args.game
			got, err := tt.args.c.CreateGame(t.Context(), tt.args.game)

			// id is automatically added
			want.DBObject = got.DBObject
			want.Rules = got.Rules // TODO: remove when rules are actually loaded/saved
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("CreateGame() did not return error when expected")
				} else {
					t.Fatalf("CreateGame() errored unexpectedly; err = \n%v", err)
				}
			}
			test.CompareAsJSON(t, got, want)
		})
	}
}

func TestUpdateGame(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game, err := c.CreateGame(t.Context(), &cs.Game{HostID: 1, Name: "Test"})
	if err != nil {
		t.Errorf("create game %s", err)
		return
	}

	game.Name = "Test2"
	if err := c.UpdateGame(t.Context(), game); err != nil {
		t.Errorf("update game %s", err)
		return
	}

	updated, err := c.GetGame(t.Context(), game.ID)

	if err != nil {
		t.Errorf("get game %s", err)
		return
	}

	assert.Equal(t, game.Name, updated.Name)

}

func TestGetGame(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	game := cs.NewGame().WithSettings(*cs.NewGameSettings().WithHost(cs.Humanoids()).WithName("test"))
	game.Area = cs.Vector{X: 1, Y: 2}
	game, err := c.CreateGame(t.Context(), game)
	if err != nil {
		t.Errorf("create game %s", err)
		return
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *cs.Game
		wantErr bool
	}{
		{"No results", args{id: 0}, nil, false},
		{"Got game", args{id: game.ID}, game, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetGame(t.Context(), tt.args.id)
			// GetGame returns a GameWithPlayers so we need the empty slice for comparison
			var want *cs.GameWithPlayers
			if tt.want != nil && got != nil {
				want = &cs.GameWithPlayers{Game: got.Game}
				want.UpdatedAt = got.UpdatedAt
				want.CreatedAt = got.CreatedAt
			}

			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("GetGame() did not return error when expected")
				} else {
					t.Fatalf("GetGame() returned errored unexpectedly; err = \n%v", err)
				}
			}
			test.CompareAsJSON(t, got, want)
		})
	}
}

func TestGetGames(t *testing.T) {
	c := connectTestDB()

	// start with 1 game from connectTestDB
	result, err := c.GetGames(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	_, err = c.CreateGame(t.Context(), &cs.Game{HostID: 1, Name: "Test"})
	if err != nil {
		t.Errorf("create game %s", err)
		return
	}

	result, err = c.GetGames(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestGetOpenGames(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	// start with 1 game from connectTestDB
	result, err := c.GetOpenGames(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	_, err = c.CreateGame(t.Context(), &cs.Game{HostID: 2, Name: "Test", State: cs.GameStateSetup, OpenPlayerSlots: 1, Public: true})
	if err != nil {
		t.Errorf("create game %s", err)
		return
	}

	//
	result, err = c.GetOpenGames(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	// create a second closed game
	_, err = c.CreateGame(t.Context(), &cs.Game{HostID: 2, Name: "Test", State: cs.GameStateSetup, OpenPlayerSlots: 0})
	if err != nil {
		t.Errorf("create game %s", err)
		return
	}

	result, err = c.GetOpenGames(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

}

func TestGetGameWithPlayersStatus(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	fg := c.createTestFullGame(t.Context())

	// make sure we don't see our own games
	game, err := c.GetGame(t.Context(), fg.ID)
	assert.Nil(t, err)
	assert.Equal(t, fg.ID, game.ID)
	assert.Equal(t, 1, len(game.Players))

}

func TestDeleteGames(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	result, err := c.GetGames(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))

	game, err := c.CreateGame(t.Context(), &cs.Game{HostID: 1, Name: "Test"})
	if err != nil {
		t.Errorf("create game %s", err)
		return
	}

	// should have our game in the db
	result, err = c.GetGames(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))

	if err := c.DeleteGame(t.Context(), game.ID); err != nil {
		t.Errorf("delete game %s", err)
		return
	}

	// should be no games left in db
	result, err = c.GetGames(t.Context())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestUpdateFullGame(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	fg := c.createTestFullGame(t.Context())

	if err := c.UpdateFullGame(t.Context(), fg); err != nil {
		t.Errorf("update full game %s", err)
		return
	}

	updated, err := c.GetFullGame(t.Context(), fg.ID)

	if err != nil {
		t.Errorf("get full game %s", err)
		return
	}

	assert.Equal(t, fg.Name, updated.Name)

}
