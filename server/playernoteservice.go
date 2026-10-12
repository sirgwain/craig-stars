package server

import (
	"context"
	"fmt"

	"log/slog"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
)

func NewPlayerNoteServiceHandler() craig_starsv1connect.PlayerNoteServiceHandler {
	return &playerNoteService{}
}

type playerNoteService struct{}

// GetPlayerNotes returns all the notes the current player has written in this game
func (s *playerNoteService) GetPlayerNotes(ctx context.Context, req *connect.Request[craig_starsv1.GetPlayerNotesRequest]) (*connect.Response[craig_starsv1.GetPlayerNotesResponse], error) {
	dbClient := contextDb(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	notes, err := dbClient.GetPlayerNotesForPlayer(ctx, game.ID, gamePlayer.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get notes from db"))
	}

	return connect.NewResponse(&craig_starsv1.GetPlayerNotesResponse{
		Notes: converter.C.ConvertCSPlayerNotes(notes),
	}), nil
}

// SavePlayerNote creates a note, or updates the text and tags of an existing one
func (s *playerNoteService) SavePlayerNote(ctx context.Context, req *connect.Request[craig_starsv1.SavePlayerNoteRequest]) (*connect.Response[craig_starsv1.SavePlayerNoteResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	if req.Msg.Note == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing note"))
	}

	note := converter.C.ConvertPlayerNote(req.Msg.Note)
	note.GameID = game.ID
	note.PlayerNum = gamePlayer.Num

	var existing *cs.PlayerNote
	var err error
	if note.ID != 0 {
		existing, err = dbClient.GetPlayerNote(ctx, note.ID)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get note from db"))
		}
		if existing == nil || existing.GameID != game.ID || existing.PlayerNum != gamePlayer.Num {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("note not found"))
		}
		// a note stays attached to whatever it was created for
		note.TargetType = existing.TargetType
		note.TargetNum = existing.TargetNum
		note.TargetPlayerNum = existing.TargetPlayerNum
	}

	if err := note.Validate(game.Area); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if existing == nil && note.TargetType == cs.MapObjectTypeFleet {
		// a note for a fleet that doesn't exist would end up on the next fleet to get that num
		fleet, err := dbClient.GetFleetByNum(ctx, game.ID, gamePlayer.Num, note.TargetNum)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get fleet from db"))
		}
		if fleet == nil {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("fleet not found"))
		}
	}

	if existing == nil && note.HasTarget() {
		// a player has one note per map object, replace it if they already have one
		existing, err = dbClient.GetPlayerNoteForTarget(ctx, game.ID, gamePlayer.Num, cs.MapObjectTarget{
			TargetType:      note.TargetType,
			TargetNum:       note.TargetNum,
			TargetPlayerNum: note.TargetPlayerNum,
		})
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get note from db"))
		}
	}

	if existing != nil {
		note.ID = existing.ID
		note.CreatedAt = existing.CreatedAt
	} else {
		notes, err := dbClient.GetPlayerNotesForPlayer(ctx, game.ID, gamePlayer.Num)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get notes from db"))
		}
		if len(notes) >= cs.PlayerNoteMaxCount {
			return nil, connect.NewError(connect.CodeResourceExhausted, fmt.Errorf("a player can have at most %d notes", cs.PlayerNoteMaxCount))
		}
	}

	if err := dbWriteClient.SavePlayerNote(ctx, note); err != nil {
		slog.Error("save player note in database", slog.Any("error", err), slog.Int64("ID", note.ID))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.SavePlayerNoteResponse{
		Note: converter.C.ConvertCSPlayerNote(note),
	}), nil
}

// DeletePlayerNote deletes one of the current player's notes
func (s *playerNoteService) DeletePlayerNote(ctx context.Context, req *connect.Request[craig_starsv1.DeletePlayerNoteRequest]) (*connect.Response[craig_starsv1.DeletePlayerNoteResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	note, err := dbClient.GetPlayerNote(ctx, req.Msg.NoteId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get note from db"))
	}
	if note == nil || note.GameID != game.ID || note.PlayerNum != gamePlayer.Num {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("note not found"))
	}

	if err := dbWriteClient.DeletePlayerNote(ctx, note.ID); err != nil {
		slog.Error("delete player note from database", slog.Any("error", err), slog.Int64("ID", note.ID))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.DeletePlayerNoteResponse{}), nil
}
