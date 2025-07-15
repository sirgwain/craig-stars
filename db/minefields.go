package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	generated "github.com/sirgwain/craig-stars/db/generated"
)

// get a minefield by id
func (c *client) GetMinefield(ctx context.Context, id int64) (*cs.MineField, error) {
	item, err := c.reader.GetMinefield(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMineField(item), nil
}

func (c *client) GetMinefieldByNum(ctx context.Context, gameID int64, playerNum int, num int) (*cs.MineField, error) {

	item, err := c.reader.GetMinefieldByNum(ctx, generated.GetMinefieldByNumParams{
		GameID:    gameID,
		PlayerNum: int64(playerNum),
		Num:       int64(num),
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
	items, err := c.reader.GetMinefieldsForGame(ctx, gameID)

	if err == sql.ErrNoRows {
		return []*cs.MineField{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMineFields(items), nil
}

func (c *client) GetMinefieldsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.MineField, error) {
	items, err := c.reader.GetMinefieldsForPlayer(ctx, generated.GetMinefieldsForPlayerParams{
		GameID:    gameID,
		PlayerNum: int64(playerNum),
	})

	if err == sql.ErrNoRows {
		return []*cs.MineField{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMineFields(items), nil
}

func (c *client) SaveMinefield(ctx context.Context, mineField *cs.MineField) error {
	if mineField.ID == 0 {
		result, err := c.writer.CreateMinefield(ctx, c.converter.ConvertGameMineFieldToCreateParams(mineField))
		if err != nil {
			return err
		}
		mineField.ID = result
	} else {
		_, err := c.writer.UpdateMinefield(ctx, c.converter.ConvertGameMineFieldToUpdateParams(mineField))
		if err != nil {
			return err
		}
	}

	return nil
}

// delete a minefield by id
func (c *client) DeleteMineField(ctx context.Context, id int64) error {
	_, err := c.writer.DeleteMinefield(ctx, id)
	return err
}
