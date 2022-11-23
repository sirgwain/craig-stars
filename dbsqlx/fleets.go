package dbsqlx

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/sirgwain/craig-stars/game"
)

type Fleet struct {
	ID                int64      `json:"id,omitempty"`
	GameID            int64      `json:"gameId,omitempty"`
	CreatedAt         time.Time  `json:"createdAt,omitempty"`
	UpdatedAt         time.Time  `json:"updatedAt,omitempty"`
	PlayerID          int64      `json:"playerId,omitempty"`
	X                 float64    `json:"x,omitempty"`
	Y                 float64    `json:"y,omitempty"`
	Name              string     `json:"name,omitempty"`
	Num               int        `json:"num,omitempty"`
	PlayerNum         int        `json:"playerNum,omitempty"`
	Tags              Tags       `json:"tags,omitempty"`
	Waypoints         Waypoints  `json:"waypoints,omitempty"`
	RepeatOrders      bool       `json:"repeatOrders,omitempty"`
	PlanetID          int64      `json:"planetId,omitempty"`
	BaseName          string     `json:"baseName,omitempty"`
	Ironium           int        `json:"ironium,omitempty"`
	Boranium          int        `json:"boranium,omitempty"`
	Germanium         int        `json:"germanium,omitempty"`
	Colonists         int        `json:"colonists,omitempty"`
	Fuel              int        `json:"fuel,omitempty"`
	Damage            int        `json:"damage,omitempty"`
	BattlePlanID      int64      `json:"battlePlanId,omitempty"`
	HeadingX          float64    `json:"headingX,omitempty"`
	HeadingY          float64    `json:"headingY,omitempty"`
	WarpSpeed         int        `json:"warpSpeed,omitempty"`
	PreviousPositionX *float64   `json:"previousPositionX,omitempty"`
	PreviousPositionY *float64   `json:"previousPositionY,omitempty"`
	OrbitingPlanetNum int        `json:"orbitingPlanetNum,omitempty"`
	Starbase          bool       `json:"starbase,omitempty"`
	Spec              *FleetSpec `json:"spec,omitempty"`
}

// we json serialize these types with custom Scan/Value methods
type Waypoints []game.Waypoint
type FleetSpec game.FleetSpec

// db serializer to serialize this to JSON
func (item Waypoints) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// db deserializer to read this from JSON
func (item Waypoints) Scan(src interface{}) error {
	if src == nil {
		// leave empty
		return nil
	}

	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, &item)
	case string:
		return json.Unmarshal([]byte(v), &item)
	}
	return errors.New("type assertion failed")
}

// db serializer to serialize this to JSON
func (item *FleetSpec) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// db deserializer to read this from JSON
func (item *FleetSpec) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, item)
	case string:
		return json.Unmarshal([]byte(v), item)
	}
	return errors.New("type assertion failed")
}

// get a fleet by id
func (c *client) GetFleet(id int64) (*game.Fleet, error) {
	item := Fleet{}
	if err := c.db.Get(&item, "SELECT * FROM fleets WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	fleet := c.converter.ConvertFleet(&item)
	return fleet, nil
}

func (c *client) getFleetsForGame(gameId int64) ([]*game.Fleet, error) {

	// don't include password in bulk select
	items := []Fleet{}
	if err := c.db.Select(&items, `SELECT * FROM fleets WHERE gameId = ?`, gameId); err != nil {
		if err == sql.ErrNoRows {
			return []*game.Fleet{}, nil
		}
		return nil, err
	}

	results := make([]*game.Fleet, len(items))
	for i := range items {
		results[i] = c.converter.ConvertFleet(&items[i])
	}

	return results, nil
}

// create a new game
func (c *client) createFleet(fleet *game.Fleet, tx NamedExecer) error {
	item := c.converter.ConvertGameFleet(fleet)
	result, err := tx.NamedExec(`
	INSERT INTO fleets (
		createdAt,
		updatedAt,
		gameId,
		playerId,
		battlePlanId,
		x,
		y,
		name,
		num,
		playerNum,
		waypoints,
		repeatOrders,
		planetId,
		baseName,
		ironium,
		boranium,
		germanium,
		colonists,
		fuel,
		damage,
		headingX,
		headingY,
		warpSpeed,
		previousPositionX,
		previousPositionY,
		orbitingPlanetNum,
		starbase,
		spec
	)
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:gameId,
		:playerId,
		:battlePlanId,
		:x,
		:y,
		:name,
		:num,
		:playerNum,
		:waypoints,
		:repeatOrders,
		:planetId,
		:baseName,
		:ironium,
		:boranium,
		:germanium,
		:colonists,
		:fuel,
		:damage,
		:headingX,
		:headingY,
		:warpSpeed,
		:previousPositionX,
		:previousPositionY,
		:orbitingPlanetNum,
		:starbase,
		:spec
	)
	`, item)

	if err != nil {
		return err
	}

	// update the id of our passed in game
	fleet.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) createShipToken(token *game.ShipToken, tx NamedExecer) error {
	result, err := tx.NamedExec(`
	INSERT INTO shipTokens (
		createdAt,
		updatedAt,
		fleetId,
		designId,
		quantity,
		damage,
		quantityDamaged
	)
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:fleetId,
		:designId,
		:quantity,
		:damage,
		:quantityDamaged
	)
	`, token)

	if err != nil {
		return err
	}

	// update the id of our passed in game
	token.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) UpdateFleet(fleet *game.Fleet) error {
	return c.updateFleet(fleet, c.db)
}

// update an existing fleet
func (c *client) updateFleet(fleet *game.Fleet, tx NamedExecer) error {

	item := c.converter.ConvertGameFleet(fleet)

	if _, err := tx.NamedExec(`
	UPDATE fleets SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		playerId = :playerId,
		battlePlanId = :battlePlanId,
		x = :x,
		y = :y,
		name = :name,
		num = :num,
		playerNum = :playerNum,
		waypoints = :waypoints,
		repeatOrders = :repeatOrders,
		planetId = :planetId,
		baseName = :baseName,
		ironium = :ironium,
		boranium = :boranium,
		germanium = :germanium,
		colonists = :colonists,
		fuel = :fuel,
		damage = :damage,
		headingX = :headingX,
		headingY = :headingY,
		warpSpeed = :warpSpeed,
		previousPositionX = :previousPositionX,
		previousPositionY = :previousPositionY,
		orbitingPlanetNum = :orbitingPlanetNum,
		starbase = :starbase,
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// update an existing shiptoken
func (c *client) updateShipToken(token *game.ShipToken, tx NamedExecer) error {
	if _, err := tx.NamedExec(`
	UPDATE shipTokens SET
		updatedAt = CURRENT_TIMESTAMP,
		fleetId = :fleetId,
		designId = :designId,
		quantity = :quantity,
		damage = :damage,
		quantityDamaged = :quantityDamaged
	WHERE id = :id
	`, token); err != nil {
		return err
	}

	return nil
}
