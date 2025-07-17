package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	generated "github.com/sirgwain/craig-stars/db/generated"
)

// get a minefield by id
func (c *client) GetMinefield(ctx context.Context, id int64) (*cs.Minefield, error) {
	item, err := c.reader.GetMinefield(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMinefield(item), nil
}

func (c *client) GetMinefieldByNum(ctx context.Context, gameID int64, playerNum int, num int) (*cs.Minefield, error) {

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

	return c.converter.ConvertMinefield(item), nil

}

func (c *client) getMinefieldsForGame(ctx context.Context, gameID int64) ([]*cs.Minefield, error) {
	items, err := c.reader.GetMinefieldsForGame(ctx, gameID)

	if err == sql.ErrNoRows {
		return []*cs.Minefield{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMinefields(items), nil
}

func (c *client) GetMinefieldsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.Minefield, error) {
	items, err := c.reader.GetMinefieldsForPlayer(ctx, generated.GetMinefieldsForPlayerParams{
		GameID:    gameID,
		PlayerNum: int64(playerNum),
	})

	if err == sql.ErrNoRows {
		return []*cs.Minefield{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMinefields(items), nil
}

func (c *client) SaveMinefield(ctx context.Context, minefield *cs.Minefield) error {
	if minefield.ID == 0 {
		result, err := c.writer.CreateMinefield(ctx, c.converter.ConvertGameMinefieldToCreateParams(minefield))
		if err != nil {
			return err
		}
		minefield.ID = result
	} else {
		_, err := c.writer.UpdateMinefield(ctx, c.converter.ConvertGameMinefieldToUpdateParams(minefield))
		if err != nil {
			return err
		}
	}

	return nil
}

// delete a minefield by id
func (c *client) DeleteMinefield(ctx context.Context, id int64) error {
	_, err := c.writer.DeleteMinefield(ctx, id)
	return err
}
