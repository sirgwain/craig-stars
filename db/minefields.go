package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	generated "github.com/sirgwain/craig-stars/db/generated"
)

// get a minefield by id
func (c *client) GetMineField(ctx context.Context, id int64) (*cs.MineField, error) {
	item, err := c.reader.GetMineField(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMineField(item), nil
}

func (c *client) GetMineFieldByNum(ctx context.Context, gameID int64, playerNum int, num int) (*cs.MineField, error) {

	item, err := c.reader.GetMineFieldByNum(ctx, generated.GetMineFieldByNumParams{
		Gameid:    gameID,
		Playernum: sql.NullInt64{Valid: true, Int64: int64(playerNum)},
		Num:       sql.NullInt64{Valid: true, Int64: int64(num)},
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMineField(item), nil

}

func (c *client) getMineFieldsForGame(ctx context.Context, gameID int64) ([]*cs.MineField, error) {
	items, err := c.reader.GetMineFieldsForGame(ctx, gameID)

	if err == sql.ErrNoRows {
		return []*cs.MineField{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMineFields(items), nil
}

func (c *client) GetMineFieldsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.MineField, error) {
	items, err := c.reader.GetMineFieldsForPlayer(ctx, generated.GetMineFieldsForPlayerParams{
		Gameid:    gameID,
		Playernum: sql.NullInt64{Valid: true, Int64: int64(playerNum)},
	})

	if err == sql.ErrNoRows {
		return []*cs.MineField{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMineFields(items), nil
}

func (c *client) CreateMineField(ctx context.Context, minefield *cs.MineField) (*cs.MineField, error) {
	result, err := c.writer.CreateMineField(ctx, c.converter.ConvertGameMineFieldToCreateParams(minefield))

	if err != nil {
		return nil, err
	}

	minefield.ID = result.ID
	minefield.CreatedAt = result.Createdat
	minefield.UpdatedAt = result.Updatedat
	return minefield, nil
}

// update an existing minefield
func (c *client) UpdateMineField(ctx context.Context, minefield *cs.MineField) error {

	result, err := c.writer.UpdateMineField(ctx, c.converter.ConvertGameMineFieldToUpdateParams(minefield))
	if err != nil {
		return err
	}

	minefield.UpdatedAt = result
	return nil
}

// delete a minefield by id
func (c *client) DeleteMineField(ctx context.Context, id int64) error {
	return c.writer.DeleteMineField(ctx, id)
}
