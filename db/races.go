package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
)

func (c *client) GetRaces(ctx context.Context) ([]cs.Race, error) {

	items, err := c.reader.GetRaces(ctx)

	if err == sql.ErrNoRows {
		return []cs.Race{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertRaces(items), nil
}

func (c *client) GetRacesForUser(ctx context.Context, userID int64) ([]cs.Race, error) {

	items, err := c.reader.GetRacesForUser(ctx, userID)

	if err == sql.ErrNoRows {
		return []cs.Race{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertRaces(items), nil

}

// get a race by id
func (c *client) GetRace(ctx context.Context, id int64) (*cs.Race, error) {
	item, err := c.reader.GetRace(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	race := c.converter.ConvertRace(item)
	return &race, nil
}

// create a new race
func (c *client) CreateRace(ctx context.Context, race *cs.Race) (*cs.Race, error) {

	result, err := c.writer.CreateRace(ctx, c.converter.ConvertGameRaceToCreateParams(race))
	if err != nil {
		return nil, err
	}

	created := c.converter.ConvertRace(result)
	return &created, nil
}

// update an existing race
func (c *client) UpdateRace(ctx context.Context, race *cs.Race) error {

	result, err := c.writer.UpdateRace(ctx, c.converter.ConvertGameRaceToUpdateParams(race))
	if err != nil {
		return err
	}

	race.UpdatedAt = result.Updatedat
	return nil
}

// delete a race by id
func (c *client) DeleteRace(ctx context.Context, id int64) error {
	return c.writer.DeleteRace(ctx, id)
}

// delete all races belonging to a user
func (c *client) DeleteUserRaces(ctx context.Context, userID int64) error {
	return c.writer.DeleteUserRaces(ctx, userID)
}
