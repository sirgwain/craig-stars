package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	generated "github.com/sirgwain/craig-stars/db/generated"
)

// get a salvage by id
func (c *client) GetSalvage(ctx context.Context, id int64) (*cs.Salvage, error) {
	item, err := c.reader.GetSalvage(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertSalvage(item), nil
}

func (c *client) GetSalvageByNum(ctx context.Context, gameID int64, num int) (*cs.Salvage, error) {

	item, err := c.reader.GetSalvageByNum(ctx, generated.GetSalvageByNumParams{
		Gameid: gameID,
		Num:    sql.NullInt64{Valid: true, Int64: int64(num)},
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertSalvage(item), nil

}

func (c *client) GetSalvagesForGame(ctx context.Context, gameID int64) ([]*cs.Salvage, error) {
	items, err := c.reader.GetSalvagesForGame(ctx, gameID)

	if err == sql.ErrNoRows {
		return []*cs.Salvage{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertSalvages(items), nil
}

func (c *client) GetSalvagesForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.Salvage, error) {
	items, err := c.reader.GetSalvagesForPlayer(ctx, generated.GetSalvagesForPlayerParams{
		Gameid:    gameID,
		Playernum: sql.NullInt64{Valid: true, Int64: int64(playerNum)},
	})

	if err == sql.ErrNoRows {
		return []*cs.Salvage{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertSalvages(items), nil
}

func (c *client) SaveSalvage(ctx context.Context, salvage *cs.Salvage) error {
	if salvage.ID == 0 {
		result, err := c.writer.CreateSalvage(ctx, c.converter.ConvertGameSalvageToCreateParams(salvage))
		if err != nil {
			return err
		}
		salvage.ID = result.ID
		salvage.CreatedAt = result.Createdat
		salvage.UpdatedAt = result.Updatedat
	} else {
		result, err := c.writer.UpdateSalvage(ctx, c.converter.ConvertGameSalvageToUpdateParams(salvage))
		if err != nil {
			return err
		}
		salvage.UpdatedAt = result
	}

	return nil
}

// delete a salvage by id
func (c *client) DeleteSalvage(ctx context.Context, id int64) error {
	return c.writer.DeleteSalvage(ctx, id)
}
