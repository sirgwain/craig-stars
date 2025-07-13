package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	generated "github.com/sirgwain/craig-stars/db/generated"
)

// get a fleet by id
func (c *client) GetFleet(ctx context.Context, id int64) (*cs.Fleet, error) {
	item, err := c.reader.GetFleet(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertFleet(item), nil
}

func (c *client) GetFleetByNum(ctx context.Context, gameID int64, playerNum int, num int) (*cs.Fleet, error) {

	item, err := c.reader.GetFleetByNum(ctx, generated.GetFleetByNumParams{
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

	return c.converter.ConvertFleet(item), nil

}

func (c *client) GetFleetsForGame(ctx context.Context, gameID int64) ([]*cs.Fleet, error) {
	items, err := c.reader.GetFleetsForGame(ctx, gameID)

	if err == sql.ErrNoRows {
		return []*cs.Fleet{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertFleets(items), nil
}

func (c *client) GetFleetsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.Fleet, error) {
	items, err := c.reader.GetFleetsForPlayer(ctx, generated.GetFleetsForPlayerParams{
		Gameid:    gameID,
		Playernum: sql.NullInt64{Valid: true, Int64: int64(playerNum)},
	})

	if err == sql.ErrNoRows {
		return []*cs.Fleet{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertFleets(items), nil
}

func (c *client) GetFleetsByNums(ctx context.Context, gameID int64, playerNum int, nums []int) ([]*cs.Fleet, error) {

	numsInt64 := make([]sql.NullInt64, len(nums))
	for i, n := range nums {
		numsInt64[i] = sql.NullInt64{Valid: true, Int64: int64(n)}
	}
	items, err := c.reader.GetFleetsByNums(ctx, generated.GetFleetsByNumsParams{
		Gameid:    gameID,
		Playernum: sql.NullInt64{Valid: true, Int64: int64(playerNum)},
		Nums:      numsInt64,
	})

	if err == sql.ErrNoRows {
		return []*cs.Fleet{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertFleets(items), nil
}

func (c *client) SaveFleet(ctx context.Context, fleet *cs.Fleet) error {
	if fleet.ID == 0 {
		result, err := c.writer.CreateFleet(ctx, c.converter.ConvertGameFleetToCreateParams(fleet))
		if err != nil {
			return err
		}
		fleet.ID = result.ID
		fleet.CreatedAt = result.Createdat
		fleet.UpdatedAt = result.Updatedat
	} else {
		result, err := c.writer.UpdateFleet(ctx, c.converter.ConvertGameFleetToUpdateParams(fleet))
		if err != nil {
			return err
		}
		fleet.UpdatedAt = result
	}

	return nil
}

// delete a fleet by id
func (c *client) DeleteFleet(ctx context.Context, id int64) error {
	return c.writer.DeleteFleet(ctx, id)
}
