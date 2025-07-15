package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	generated "github.com/sirgwain/craig-stars/db/generated"
)

// get a wormhole by id
func (c *client) GetWormhole(ctx context.Context, id int64) (*cs.Wormhole, error) {
	item, err := c.reader.GetWormhole(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertWormhole(item), nil
}

func (c *client) GetWormholeByNum(ctx context.Context, gameID int64, num int) (*cs.Wormhole, error) {

	item, err := c.reader.GetWormholeByNum(ctx, generated.GetWormholeByNumParams{
		GameID: gameID,
		Num:    int64(num),
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertWormhole(item), nil

}

func (c *client) GetWormholesForGame(ctx context.Context, gameID int64) ([]*cs.Wormhole, error) {
	items, err := c.reader.GetWormholesForGame(ctx, gameID)

	if err == sql.ErrNoRows {
		return []*cs.Wormhole{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertWormholes(items), nil
}

func (c *client) SaveWormhole(ctx context.Context, wormhole *cs.Wormhole) error {
	if wormhole.ID == 0 {
		result, err := c.writer.CreateWormhole(ctx, c.converter.ConvertGameWormholeToCreateParams(wormhole))
		if err != nil {
			return err
		}
		wormhole.ID = result
	} else {
		_, err := c.writer.UpdateWormhole(ctx, c.converter.ConvertGameWormholeToUpdateParams(wormhole))
		if err != nil {
			return err
		}
	}

	return nil
}

// delete a wormhole by id
func (c *client) DeleteWormhole(ctx context.Context, id int64) error {
	_, err := c.writer.DeleteWormhole(ctx, id)
	return err
}
