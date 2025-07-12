package db

import (
	"context"
	"database/sql"
	"fmt"

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

func (c *client) CreateFleet(ctx context.Context, fleet *cs.Fleet) (*cs.Fleet, error) {
	result, err := c.writer.CreateFleet(ctx, c.converter.ConvertGameFleetToCreateParams(fleet))

	if err != nil {
		return nil, err
	}

	created := c.converter.ConvertFleet(result)
	return created, nil
}

// update an existing fleet
func (c *client) UpdateFleet(ctx context.Context, fleet *cs.Fleet) error {

	result, err := c.writer.UpdateFleet(ctx, c.converter.ConvertGameFleetToUpdateParams(fleet))
	if err != nil {
		return err
	}

	fleet.UpdatedAt = result.Updatedat
	return nil
}

// This should always be wrapped in a transaction
func (c *client) CreateUpdateOrDeleteFleets(ctx context.Context, gameID int64, fleets []*cs.Fleet) error {

	// create/update fleets
	for i, fleet := range fleets {
		if fleet.ID == 0 {
			fleet.GameID = gameID
			var err error
			created, err := c.CreateFleet(ctx, fleet)
			if err != nil {
				return fmt.Errorf("create fleet: %w", err)
			}
			// log.Debug().Int64("GameID", fleet.GameID).Int64("ID", fleet.ID).Msgf("Created fleet %s", fleet.Name)
			fleets[i] = created
		} else if fleet.Delete {
			if err := c.DeleteFleet(ctx, fleet.ID); err != nil {
				return fmt.Errorf("delete fleet: %w", err)
			}
			// log.Debug().Int64("GameID", fleet.GameID).Int64("ID", fleet.ID).Msgf("Deleted fleet %s", fleet.Name)
		} else {
			if err := c.UpdateFleet(ctx, fleet); err != nil {
				return fmt.Errorf("update fleet: %w", err)
			}
			// log.Debug().Int64("GameID", fleet.GameID).Int64("ID", fleet.ID).Msgf("Updated fleet %s", fleet.Name)
		}
	}

	return nil
}

// delete a fleet by id
func (c *client) DeleteFleet(ctx context.Context, id int64) error {
	return c.writer.DeleteFleet(ctx, id)
}
