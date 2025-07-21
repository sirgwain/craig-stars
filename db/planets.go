package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	generated "github.com/sirgwain/craig-stars/db/generated"
)

// get a planet by id
func (c *client) GetPlanet(ctx context.Context, id int64) (*cs.Planet, error) {
	item, err := c.reader.GetPlanet(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertPlanet(item), nil
}

func (c *client) GetPlanetsForGame(ctx context.Context, gameID int64) ([]*cs.Planet, error) {

	items, err := c.reader.GetPlanetsForGame(ctx, gameID)

	if err == sql.ErrNoRows {
		return []*cs.Planet{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertPlanets(items), nil
}

func (c *client) GetPlanetsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.Planet, error) {

	items, err := c.reader.GetPlanetsForPlayer(ctx, generated.GetPlanetsForPlayerParams{
		GameID:    gameID,
		PlayerNum: int64(playerNum),
	})

	if err == sql.ErrNoRows {
		return []*cs.Planet{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertPlanets(items), nil
}

func (c *client) GetPlanetByNum(ctx context.Context, gameID int64, num int) (*cs.Planet, error) {

	rows, err := c.reader.GetPlanetByNum(ctx, generated.GetPlanetByNumParams{
		GameID: gameID,
		Num:    int64(num),
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	planet := c.converter.ConvertPlanet(rows[0].Planet)
	for _, row := range rows {
		if row.FleetID.Valid {
			fleet := generated.Fleet{
				ID:                row.FleetID.Int64,
				CreatedAt:         row.FleetCreatedAt.Time,
				UpdatedAt:         row.FleetUpdatedAt.Time,
				GameID:            row.FleetGameID.Int64,
				BattlePlanNum:     row.FleetBattlePlanNum,
				X:                 row.FleetX.Float64,
				Y:                 row.FleetY.Float64,
				Name:              row.FleetName.String,
				Num:               row.FleetNum,
				PlayerNum:         row.FleetPlayerNum,
				Tokens:            row.FleetTokens,
				Waypoints:         row.FleetWaypoints,
				RepeatOrders:      row.FleetRepeatOrders.Bool,
				PlanetNum:         row.FleetPlanetNum,
				BaseName:          row.FleetBaseName.String,
				Ironium:           row.FleetIronium,
				Boranium:          row.FleetBoranium,
				Germanium:         row.FleetGermanium,
				Colonists:         row.FleetColonists,
				Fuel:              row.FleetFuel,
				Age:               row.FleetAge,
				HeadingX:          row.FleetHeadingX.Float64,
				HeadingY:          row.FleetHeadingY.Float64,
				WarpSpeed:         row.FleetWarpSpeed,
				PreviousPositionX: row.FleetPreviousPositionX,
				PreviousPositionY: row.FleetPreviousPositionY,
				OrbitingPlanetNum: row.FleetOrbitingPlanetNum,
				Starbase:          row.FleetStarbase.Bool,
				Spec:              row.FleetSpec,
				Purpose:           row.FleetPurpose,
				Tags:              row.FleetTags,
			}
			planet.Starbase = c.converter.ConvertFleet(fleet)
		}
	}

	return planet, nil

}

func (c *client) SavePlanet(ctx context.Context, planet *cs.Planet) error {
	if planet.ID == 0 {
		result, err := c.writer.CreatePlanet(ctx, c.converter.ConvertGamePlanetToCreateParams(planet))
		if err != nil {
			return err
		}
		planet.ID = result
	} else {
		_, err := c.writer.UpdatePlanet(ctx, c.converter.ConvertGamePlanetToUpdateParams(planet))
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdatePlanetSpec updates only a planets spec field
func (c *client) UpdatePlanetSpec(ctx context.Context, planet *cs.Planet) error {
	_, err := c.writer.UpdatePlanetSpec(ctx, generated.UpdatePlanetSpecParams{
		ID:   planet.ID,
		Spec: (*generated.PlanetSpec)(&planet.Spec),
	})
	if err != nil {
		return err
	}

	return nil
}
