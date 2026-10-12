package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	generated "github.com/sirgwain/craig-stars/db/generated"
)

// get a player note by id
func (c *client) GetPlayerNote(ctx context.Context, id int64) (*cs.PlayerNote, error) {
	item, err := c.reader.GetPlayerNote(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertPlayerNote(item), nil
}

func (c *client) GetPlayerNotesForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.PlayerNote, error) {
	items, err := c.reader.GetPlayerNotesForPlayer(ctx, generated.GetPlayerNotesForPlayerParams{
		GameID:    gameID,
		PlayerNum: int64(playerNum),
	})
	if err == sql.ErrNoRows {
		return []*cs.PlayerNote{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertPlayerNotes(items), nil
}

// get the note a player has attached to a map object, or nil if there isn't one
func (c *client) GetPlayerNoteForTarget(ctx context.Context, gameID int64, playerNum int, target cs.MapObjectTarget) (*cs.PlayerNote, error) {
	item, err := c.reader.GetPlayerNoteForTarget(ctx, generated.GetPlayerNoteForTargetParams{
		GameID:          gameID,
		PlayerNum:       int64(playerNum),
		TargetType:      target.TargetType,
		TargetNum:       int64(target.TargetNum),
		TargetPlayerNum: int64(target.TargetPlayerNum),
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertPlayerNote(item), nil
}

func (c *client) SavePlayerNote(ctx context.Context, note *cs.PlayerNote) error {
	if note.ID == 0 {
		result, err := c.writer.CreatePlayerNote(ctx, c.converter.ConvertGamePlayerNoteToCreateParams(note))
		if err != nil {
			return err
		}
		note.ID = result
	} else {
		_, err := c.writer.UpdatePlayerNote(ctx, c.converter.ConvertGamePlayerNoteToUpdateParams(note))
		if err != nil {
			return err
		}
	}

	return nil
}

// delete a player note by id
func (c *client) DeletePlayerNote(ctx context.Context, id int64) error {
	_, err := c.writer.DeletePlayerNote(ctx, id)
	return err
}
