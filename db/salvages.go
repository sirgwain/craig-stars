package db

import (
	"context"
	"database/sql"
	"fmt"

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

func (c *client) CreateSalvage(ctx context.Context, salvage *cs.Salvage) (*cs.Salvage, error) {
	result, err := c.writer.CreateSalvage(ctx, c.converter.ConvertGameSalvageToCreateParams(salvage))

	if err != nil {
		return nil, err
	}

	created := c.converter.ConvertSalvage(result)
	return created, nil
}

// update an existing salvage
func (c *client) UpdateSalvage(ctx context.Context, salvage *cs.Salvage) error {

	result, err := c.writer.UpdateSalvage(ctx, c.converter.ConvertGameSalvageToUpdateParams(salvage))
	if err != nil {
		return err
	}

	salvage.UpdatedAt = result.Updatedat
	return nil
}

// This should always be wrapped in a transaction
func (c *client) CreateUpdateOrDeleteSalvages(ctx context.Context, gameID int64, salvages []*cs.Salvage) error {

	// create/update salvages
	for i, salvage := range salvages {
		if salvage.ID == 0 {
			salvage.GameID = gameID
			var err error
			created, err := c.CreateSalvage(ctx, salvage)
			if err != nil {
				return fmt.Errorf("create salvage: %w", err)
			}
			// log.Debug().Int64("GameID", salvage.GameID).Int64("ID", salvage.ID).Msgf("Created salvage %s", salvage.Name)
			salvages[i] = created
		} else if salvage.Delete {
			if err := c.DeleteSalvage(ctx, salvage.ID); err != nil {
				return fmt.Errorf("delete salvage: %w", err)
			}
			// log.Debug().Int64("GameID", salvage.GameID).Int64("ID", salvage.ID).Msgf("Deleted salvage %s", salvage.Name)
		} else {
			if err := c.UpdateSalvage(ctx, salvage); err != nil {
				return fmt.Errorf("update salvage: %w", err)
			}
			// log.Debug().Int64("GameID", salvage.GameID).Int64("ID", salvage.ID).Msgf("Updated salvage %s", salvage.Name)
		}
	}

	return nil
}

// delete a salvage by id
func (c *client) DeleteSalvage(ctx context.Context, id int64) error {
	return c.writer.DeleteSalvage(ctx, id)
}
