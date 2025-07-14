package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db/generated"
)

func (c *client) GetShipDesignsForGame(ctx context.Context, gameID int64) ([]*cs.ShipDesign, error) {

	items, err := c.reader.GetShipDesignsForGame(ctx, gameID)

	if err == sql.ErrNoRows {
		return []*cs.ShipDesign{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertShipDesigns(items), nil
}

func (c *client) GetShipDesignsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.ShipDesign, error) {

	items, err := c.reader.GetShipDesignsForPlayer(ctx, generated.GetShipDesignsForPlayerParams{
		Gameid:    gameID,
		Playernum: int64(playerNum),
	})

	if err == sql.ErrNoRows {
		return []*cs.ShipDesign{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertShipDesigns(items), nil
}

// get a shipDesign by id
func (c *client) GetShipDesign(ctx context.Context, id int64) (*cs.ShipDesign, error) {
	item, err := c.reader.GetShipDesign(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertShipDesign(item), nil
}

// get a shipDesign by id
func (c *client) GetShipDesignByNum(ctx context.Context, gameID int64, playerNum, num int) (*cs.ShipDesign, error) {
	item, err := c.reader.GetShipDesignByNum(ctx, generated.GetShipDesignByNumParams{
		Gameid:    gameID,
		Playernum: int64(playerNum),
		Num:       int64(num),
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertShipDesign(item), nil
}

func (c *client) SaveShipDesign(ctx context.Context, shipDesign *cs.ShipDesign) error {
	if shipDesign.ID == 0 {
		result, err := c.writer.CreateShipDesign(ctx, c.converter.ConvertGameShipDesignToCreateParams(shipDesign))
		if err != nil {
			return err
		}
		shipDesign.ID = result.ID
		shipDesign.CreatedAt = result.Createdat
		shipDesign.UpdatedAt = result.Updatedat
	} else {
		result, err := c.writer.UpdateShipDesign(ctx, c.converter.ConvertGameShipDesignToUpdateParams(shipDesign))
		if err != nil {
			return err
		}
		shipDesign.UpdatedAt = result
	}

	return nil
}

// delete a shipDesign by id
func (c *client) DeleteShipDesign(ctx context.Context, id int64) error {
	return c.writer.DeleteShipDesign(ctx, id)
}
