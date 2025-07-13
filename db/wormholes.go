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
		Gameid: gameID,
		Num:    sql.NullInt64{Valid: true, Int64: int64(num)},
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

func (c *client) CreateWormhole(ctx context.Context, wormhole *cs.Wormhole) (*cs.Wormhole, error) {
	result, err := c.writer.CreateWormhole(ctx, c.converter.ConvertGameWormholeToCreateParams(wormhole))

	if err != nil {
		return nil, err
	}

	wormhole.ID = result.ID
	wormhole.CreatedAt = result.Createdat
	wormhole.UpdatedAt = result.Updatedat
	return wormhole, nil
}

// update an existing wormhole
func (c *client) UpdateWormhole(ctx context.Context, wormhole *cs.Wormhole) error {

	result, err := c.writer.UpdateWormhole(ctx, c.converter.ConvertGameWormholeToUpdateParams(wormhole))
	if err != nil {
		return err
	}

	wormhole.UpdatedAt = result
	return nil
}

// delete a wormhole by id
func (c *client) DeleteWormhole(ctx context.Context, id int64) error {
	return c.writer.DeleteWormhole(ctx, id)
}
