package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	generated "github.com/sirgwain/craig-stars/db/generated"
)

// get a mysterytrader by id
func (c *client) GetMysteryTrader(ctx context.Context, id int64) (*cs.MysteryTrader, error) {
	item, err := c.reader.GetMysteryTrader(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMysteryTrader(item), nil
}

func (c *client) GetMysteryTraderByNum(ctx context.Context, gameID int64, num int) (*cs.MysteryTrader, error) {

	item, err := c.reader.GetMysteryTraderByNum(ctx, generated.GetMysteryTraderByNumParams{
		Gameid: gameID,
		Num:    sql.NullInt64{Valid: true, Int64: int64(num)},
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMysteryTrader(item), nil

}

func (c *client) GetMysteryTradersForGame(ctx context.Context, gameID int64) ([]*cs.MysteryTrader, error) {
	items, err := c.reader.GetMysteryTradersForGame(ctx, gameID)

	if err == sql.ErrNoRows {
		return []*cs.MysteryTrader{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMysteryTraders(items), nil
}

func (c *client) SaveMysteryTrader(ctx context.Context, mysteryTrader *cs.MysteryTrader) error {
	if mysteryTrader.ID == 0 {
		result, err := c.writer.CreateMysteryTrader(ctx, c.converter.ConvertGameMysteryTraderToCreateParams(mysteryTrader))
		if err != nil {
			return err
		}
		mysteryTrader.ID = result.ID
		mysteryTrader.CreatedAt = result.Createdat
		mysteryTrader.UpdatedAt = result.Updatedat
	} else {
		result, err := c.writer.UpdateMysteryTrader(ctx, c.converter.ConvertGameMysteryTraderToUpdateParams(mysteryTrader))
		if err != nil {
			return err
		}
		mysteryTrader.UpdatedAt = result
	}

	return nil
}

// delete a mysterytrader by id
func (c *client) DeleteMysteryTrader(ctx context.Context, id int64) error {
	return c.writer.DeleteMysteryTrader(ctx, id)
}
