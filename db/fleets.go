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
		GameID:    gameID,
		PlayerNum: int64(playerNum),
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

	int64Nums := make([]int64, 0, len(nums))
	for _, n := range nums {
		int64Nums = append(int64Nums, int64(n))
	}

	items, err := c.reader.GetFleetsByNums(ctx, generated.GetFleetsByNumsParams{
		GameID:    gameID,
		PlayerNum: int64(playerNum),
		Nums:      int64Nums,
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
		fleet.ID = result
	} else {
		_, err := c.writer.UpdateFleet(ctx, c.converter.ConvertGameFleetToUpdateParams(fleet))
		if err != nil {
			return err
		}
	}

	return nil
}

// delete a fleet by id
func (c *client) DeleteFleet(ctx context.Context, id int64) error {
	return c.writer.DeleteFleet(ctx, id)
}
