package server

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_playerNoteService(t *testing.T) {
	ctx := t.Context()
	dbConn := db.NewConn()
	cfg := config.Config{}
	cfg.Database.Filename = ":memory:"
	require.NoError(t, dbConn.Connect(ctx, &cfg))
	t.Cleanup(func() { dbConn.Close() })
	client := dbConn.NewReadWriteClient()
	user, err := client.CreateUser(ctx, cs.NewUser("host", "", "", cs.RoleUser))
	require.NoError(t, err)
	gr := NewGameRunner(dbConn, cfg)
	fullGame, err := gr.HostGame(user.ID, cs.NewGameSettings().WithHost(cs.Humanoids()).WithAIPlayer(cs.AIDifficultyNormal, 0))
	require.NoError(t, err)

	game, err := client.GetGame(ctx, fullGame.ID)
	require.NoError(t, err)
	require.Len(t, game.Players, 2)
	ctx = context.WithValue(ctx, keyDbRead, dbConn.NewReadClient())
	ctx = context.WithValue(ctx, keyDbWrite, client)
	ctx = context.WithValue(ctx, keyGame, game)
	otherCtx := context.WithValue(ctx, keyGamePlayer, &game.Players[1])
	ctx = context.WithValue(ctx, keyGamePlayer, &game.Players[0])

	service := NewPlayerNoteServiceHandler()
	save := func(ctx context.Context, note *craig_starsv1.PlayerNote) (*craig_starsv1.PlayerNote, error) {
		response, err := service.SavePlayerNote(ctx, connect.NewRequest(&craig_starsv1.SavePlayerNoteRequest{GameId: game.ID, Note: note}))
		if err != nil {
			return nil, err
		}
		return response.Msg.Note, nil
	}
	notes := func(ctx context.Context) []*craig_starsv1.PlayerNote {
		response, err := service.GetPlayerNotes(ctx, connect.NewRequest(&craig_starsv1.GetPlayerNotesRequest{GameId: game.ID}))
		require.NoError(t, err)
		return response.Msg.Notes
	}

	// a note on a planet
	planetNote, err := save(ctx, &craig_starsv1.PlayerNote{
		TargetType: craig_starsv1.MapObjectType_MAP_OBJECT_TYPE_PLANET,
		TargetNum:  3,
		Text:       " colonize next ",
	})
	require.NoError(t, err)
	assert.Equal(t, "colonize next", planetNote.Text)
	assert.Equal(t, int32(game.Players[0].Num), planetNote.PlayerNum)
	assert.NotZero(t, planetNote.GameDbObject.Id)

	// saving another note for the same planet replaces the first
	replaced, err := save(ctx, &craig_starsv1.PlayerNote{
		TargetType: craig_starsv1.MapObjectType_MAP_OBJECT_TYPE_PLANET,
		TargetNum:  3,
		Text:       "too cold",
		ShowOnMap:  true,
	})
	require.NoError(t, err)
	assert.Equal(t, planetNote.GameDbObject.Id, replaced.GameDbObject.Id)

	// a label on the map, which can be moved
	label, err := save(ctx, &craig_starsv1.PlayerNote{Position: &craig_starsv1.Vector{X: 10, Y: 20}, Text: "the front"})
	require.NoError(t, err)
	assert.True(t, label.ShowOnMap)
	label.Position = &craig_starsv1.Vector{X: 30, Y: 40}
	label.Text = "the new front"
	label, err = save(ctx, label)
	require.NoError(t, err)

	got := notes(ctx)
	require.Len(t, got, 2)
	assert.Equal(t, "too cold", got[0].Text)
	assert.True(t, got[0].ShowOnMap)
	assert.Equal(t, "the new front", got[1].Text)
	assert.Equal(t, int32(30), got[1].Position.X)
	assert.Equal(t, int32(40), got[1].Position.Y)

	// notes are private, other players can't see, change or delete them
	assert.Empty(t, notes(otherCtx))
	_, err = save(otherCtx, label)
	assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
	_, err = service.DeletePlayerNote(otherCtx, connect.NewRequest(&craig_starsv1.DeletePlayerNoteRequest{GameId: game.ID, NoteId: label.GameDbObject.Id}))
	assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))

	// another player can have their own note on the same planet
	_, err = save(otherCtx, &craig_starsv1.PlayerNote{
		TargetType: craig_starsv1.MapObjectType_MAP_OBJECT_TYPE_PLANET,
		TargetNum:  3,
		Text:       "mine",
	})
	require.NoError(t, err)
	assert.Len(t, notes(otherCtx), 1)
	assert.Len(t, notes(ctx), 2)

	_, err = save(ctx, &craig_starsv1.PlayerNote{Position: &craig_starsv1.Vector{X: 10, Y: 20}, Text: " "})
	assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))

	_, err = service.DeletePlayerNote(ctx, connect.NewRequest(&craig_starsv1.DeletePlayerNoteRequest{GameId: game.ID, NoteId: label.GameDbObject.Id}))
	require.NoError(t, err)
	assert.Len(t, notes(ctx), 1)

	// tags are saved with the note, and a planet can be tagged without writing anything
	replaced.Tags = []string{" gate  hub ", "Target", "target"}
	replaced, err = save(ctx, replaced)
	require.NoError(t, err)
	assert.Equal(t, []string{"gate hub", "Target"}, replaced.Tags)
	assert.Equal(t, []string{"gate hub", "Target"}, notes(ctx)[0].Tags)
	tagged, err := save(ctx, &craig_starsv1.PlayerNote{
		TargetType: craig_starsv1.MapObjectType_MAP_OBJECT_TYPE_PLANET,
		TargetNum:  4,
		Tags:       []string{"Target"},
	})
	require.NoError(t, err)
	assert.Empty(t, tagged.Text)
	assert.Len(t, notes(ctx), 2)

	// a player can put a note on their own fleets
	playerNum := game.Players[0].Num
	fleets, err := client.GetFleetsForPlayer(ctx, game.ID, playerNum)
	require.NoError(t, err)
	var fleet *cs.Fleet
	for _, f := range fleets {
		if !f.Starbase {
			fleet = f
			break
		}
	}
	require.NotNil(t, fleet)
	fleetNote := &craig_starsv1.PlayerNote{
		TargetType:      craig_starsv1.MapObjectType_MAP_OBJECT_TYPE_FLEET,
		TargetNum:       int32(fleet.Num),
		TargetPlayerNum: int32(playerNum),
		Text:            "decoy",
		Tags:            []string{"scouts"},
	}
	_, err = save(ctx, fleetNote)
	require.NoError(t, err)
	assert.Len(t, notes(ctx), 3)

	// but not on a fleet that doesn't exist, or another player's
	_, err = save(ctx, &craig_starsv1.PlayerNote{
		TargetType:      craig_starsv1.MapObjectType_MAP_OBJECT_TYPE_FLEET,
		TargetNum:       9999,
		TargetPlayerNum: int32(playerNum),
		Text:            "ghost",
	})
	assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
	_, err = save(otherCtx, fleetNote)
	assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))

	// fleet nums are reused, so the note goes away with the fleet
	require.NoError(t, client.DeleteFleet(ctx, fleet.ID))
	got = notes(ctx)
	require.Len(t, got, 2)
	for _, note := range got {
		assert.Equal(t, craig_starsv1.MapObjectType_MAP_OBJECT_TYPE_PLANET, note.TargetType)
	}
}
