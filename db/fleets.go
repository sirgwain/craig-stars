package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type Fleet struct {
	ID                int64       `json:"id,omitempty"`
	GameID            int64       `json:"gameId,omitempty"`
	CreatedAt         time.Time   `json:"createdAt,omitempty"`
	UpdatedAt         time.Time   `json:"updatedAt,omitempty"`
	X                 float64     `json:"x,omitempty"`
	Y                 float64     `json:"y,omitempty"`
	Name              string      `json:"name,omitempty"`
	Num               int         `json:"num,omitempty"`
	PlayerNum         int         `json:"playerNum,omitempty"`
	Tags              Tags        `json:"tags,omitempty"`
	Tokens            *ShipTokens `json:"tokens,omitempty"`
	Waypoints         *Waypoints  `json:"waypoints,omitempty"`
	RepeatOrders      bool        `json:"repeatOrders,omitempty"`
	PlanetNum         int         `json:"planetNum,omitempty"`
	BaseName          string      `json:"baseName,omitempty"`
	Ironium           int         `json:"ironium,omitempty"`
	Boranium          int         `json:"boranium,omitempty"`
	Germanium         int         `json:"germanium,omitempty"`
	Colonists         int         `json:"colonists,omitempty"`
	Fuel              int         `json:"fuel,omitempty"`
	Damage            int         `json:"damage,omitempty"`
	BattlePlanName    string      `json:"battlePlanName,omitempty"`
	HeadingX          float64     `json:"headingX,omitempty"`
	HeadingY          float64     `json:"headingY,omitempty"`
	WarpSpeed         int         `json:"warpSpeed,omitempty"`
	PreviousPositionX *float64    `json:"previousPositionX,omitempty"`
	PreviousPositionY *float64    `json:"previousPositionY,omitempty"`
	OrbitingPlanetNum int         `json:"orbitingPlanetNum,omitempty"`
	Starbase          bool        `json:"starbase,omitempty"`
	Spec              *FleetSpec  `json:"spec,omitempty"`
}

// we json serialize these types with custom Scan/Value methods
type ShipTokens []cs.ShipToken
type Waypoints []cs.Waypoint
type FleetSpec cs.FleetSpec

// db serializer to serialize this to JSON
func (item *ShipTokens) Value() (driver.Value, error) {
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
func (item *ShipTokens) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// db serializer to serialize this to JSON
func (item *Waypoints) Value() (driver.Value, error) {
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
func (item *Waypoints) Scan(src interface{}) error {
	return scanJSON(src, item)
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
	return scanJSON(src, item)
}

// get a fleet by id
func (c *client) GetFleet(id int64) (*cs.Fleet, error) {
	item := Fleet{}
	if err := c.db.Get(&item, "SELECT * FROM fleets WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	fleet := c.converter.ConvertFleet(&item)

	// load a fleet's designs
	designNums := make([]int, 0, len(fleet.Tokens))
	for i := range fleet.Tokens {
		designNums = append(designNums, fleet.Tokens[i].DesignNum)
	}

	// this might be an error case, or we're in a unit test
	if len(designNums) == 0 {
		log.Warn().Int64("ID", fleet.ID).Msg("fleet has no designs associated with tokens")
		return fleet, nil
	}

	designs, err := c.getShipDesignsByNums(designNums)
	if err != nil {
		return nil, fmt.Errorf("get designs by nums -> %w", err)
	}

	fleet.InjectDesigns(designs)

	return fleet, nil
}

func (c *client) getFleetsForGame(gameID int64) ([]*cs.Fleet, error) {

	items := []Fleet{}
	if err := c.db.Select(&items, `SELECT * FROM fleets WHERE gameId = ?`, gameID); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.Fleet{}, nil
		}
		return nil, err
	}

	results := make([]*cs.Fleet, len(items))
	for i := range items {
		results[i] = c.converter.ConvertFleet(&items[i])
	}

	return results, nil
}

func (c *client) GetFleetsForPlayer(gameID int64, playerNum int) ([]*cs.Fleet, error) {

	items := []Fleet{}
	if err := c.db.Select(&items, `SELECT * FROM fleets WHERE gameId = ? AND playerNum = ?`, gameID, playerNum); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.Fleet{}, nil
		}
		return nil, err
	}

	results := make([]*cs.Fleet, len(items))
	for i := range items {
		results[i] = c.converter.ConvertFleet(&items[i])
	}

	return results, nil
}

func (c *client) GetFleetByNum(gameID int64, num int) (*cs.Fleet, error) {

	item := Fleet{}
	if err := c.db.Get(&item, `SELECT * FROM fleets WHERE gameId = ? AND num = ?`, gameID, num); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	fleet := c.converter.ConvertFleet(&item)
	return fleet, nil

}

// create a new game
func (c *client) createFleet(fleet *cs.Fleet, tx SQLExecer) error {
	item := c.converter.ConvertGameFleet(fleet)
	result, err := tx.NamedExec(`
	INSERT INTO fleets (
		createdAt,
		updatedAt,
		gameId,
		battlePlanName,
		x,
		y,
		name,
		num,
		playerNum,
		tokens,
		waypoints,
		repeatOrders,
		planetNum,
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
		:battlePlanName,
		:x,
		:y,
		:name,
		:num,
		:playerNum,
		:tokens,
		:waypoints,
		:repeatOrders,
		:planetNum,
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

func (c *client) UpdateFleet(fleet *cs.Fleet) error {
	return c.updateFleet(fleet, c.db)
}

// update an existing fleet
func (c *client) updateFleet(fleet *cs.Fleet, tx SQLExecer) error {

	item := c.converter.ConvertGameFleet(fleet)

	if _, err := tx.NamedExec(`
	UPDATE fleets SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		battlePlanName = :battlePlanName,
		x = :x,
		y = :y,
		name = :name,
		num = :num,
		playerNum = :playerNum,
		tokens = :tokens,
		waypoints = :waypoints,
		repeatOrders = :repeatOrders,
		planetNum = :planetNum,
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

func (c *client) DeleteFleet(fleetID int64) error {
	if _, err := c.db.Exec("DELETE FROM fleets where id = ?", fleetID); err != nil {
		return fmt.Errorf("delete fleet %d %w", fleetID, err)
	}
	return nil
}

func (c *client) deleteFleet(fleetID int64, tx SQLExecer) error {
	if _, err := tx.Exec("DELETE FROM fleets where id = ?", fleetID); err != nil {
		return fmt.Errorf("delete fleet %d %w", fleetID, err)
	}
	return nil
}
