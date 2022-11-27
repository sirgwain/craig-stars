package dbsqlx

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
	Waypoints         *Waypoints `json:"waypoints,omitempty"`
	RepeatOrders      bool       `json:"repeatOrders,omitempty"`
	PlanetNum         int        `json:"planetNum,omitempty"`
	BaseName          string     `json:"baseName,omitempty"`
	Ironium           int        `json:"ironium,omitempty"`
	Boranium          int        `json:"boranium,omitempty"`
	Germanium         int        `json:"germanium,omitempty"`
	Colonists         int        `json:"colonists,omitempty"`
	Fuel              int        `json:"fuel,omitempty"`
	Damage            int        `json:"damage,omitempty"`
	BattlePlanName    string     `json:"battlePlanName,omitempty"`
	HeadingX          float64    `json:"headingX,omitempty"`
	HeadingY          float64    `json:"headingY,omitempty"`
	WarpSpeed         int        `json:"warpSpeed,omitempty"`
	PreviousPositionX *float64   `json:"previousPositionX,omitempty"`
	PreviousPositionY *float64   `json:"previousPositionY,omitempty"`
	OrbitingPlanetNum int        `json:"orbitingPlanetNum,omitempty"`
	Starbase          bool       `json:"starbase,omitempty"`
	Spec              *FleetSpec `json:"spec,omitempty"`
}

type ShipToken struct {
	ID              int64        `json:"id"`
	CreatedAt       sql.NullTime `json:"createdAt"`
	UpdatedAt       sql.NullTime `json:"updatedAt"`
	FleetID         int64        `json:"fleetId"`
	DesignUUID      uuid.UUID    `json:"designUuid,omitempty"`
	Quantity        int          `json:"quantity"`
	Damage          float64      `json:"damage"`
	QuantityDamaged int          `json:"quantityDamaged"`
}

// we json serialize these types with custom Scan/Value methods
type Waypoints []game.Waypoint
type FleetSpec game.FleetSpec

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

type fleetJoin struct {
	Fleet     `json:"fleet,omitempty"`
	ShipToken `json:"fleetShipToken,omitempty"`
}

const fleetJoinSelect = `
SELECT 
	f.id AS 'fleet.id',
	f.createdAt AS 'fleet.createdAt',
	f.updatedAt AS 'fleet.updatedAt',
	f.gameId AS 'fleet.gameId',
	f.playerId AS 'fleet.playerId',
	f.battlePlanName AS 'fleet.battlePlanName',
	f.x AS 'fleet.x',
	f.y AS 'fleet.y',
	f.name AS 'fleet.name',
	f.num AS 'fleet.num',
	f.playerNum AS 'fleet.playerNum',
	f.waypoints AS 'fleet.waypoints',
	f.repeatOrders AS 'fleet.repeatOrders',
	f.planetNum AS 'fleet.planetNum',
	f.baseName AS 'fleet.baseName',
	f.ironium AS 'fleet.ironium',
	f.boranium AS 'fleet.boranium',
	f.germanium AS 'fleet.germanium',
	f.colonists AS 'fleet.colonists',
	f.fuel AS 'fleet.fuel',
	f.damage AS 'fleet.damage',
	f.headingX AS 'fleet.headingX',
	f.headingY AS 'fleet.headingY',
	f.warpSpeed AS 'fleet.warpSpeed',
	f.previousPositionX AS 'fleet.previousPositionX',
	f.previousPositionY AS 'fleet.previousPositionY',
	f.orbitingPlanetNum AS 'fleet.orbitingPlanetNum',
	f.starbase AS 'fleet.starbase',
	f.spec AS 'fleet.spec',

	COALESCE(t.id, 0) AS 'fleetShipToken.id',
	t.createdAt AS 'fleetShipToken.createdAt',
	t.updatedAt AS 'fleetShipToken.updatedAt',
	COALESCE(t.fleetId, 0) AS 'fleetShipToken.fleetId',
	COALESCE(t.designUuid, '') AS 'fleetShipToken.designUuid',
	COALESCE(t.quantity, 0) AS 'fleetShipToken.quantity',
	COALESCE(t.damage, 0) AS 'fleetShipToken.damage',
	COALESCE(t.quantityDamaged, 0) AS 'fleetShipToken.quantityDamaged'
	FROM fleets f
	LEFT JOIN shipTokens t
		ON f.id = t.fleetId
`

// scan through rows from a fleet join and return a list of fleets and their tokens
func (c *client) scanFleetJoin(rows []fleetJoin) []*game.Fleet {
	fleets := []*game.Fleet{}

	fleet := &game.Fleet{}
	for _, row := range rows {
		// found a new fleet
		if row.Fleet.ID != 0 {
			fleet = c.converter.ConvertFleet(&row.Fleet)
			fleet.Tokens = []game.ShipToken{}
			fleets = append(fleets, fleet)
		}

		if row.ShipToken.ID != 0 {
			fleet.Tokens = append(fleet.Tokens, c.converter.ConvertShipToken(row.ShipToken))
		}
	}

	return fleets
}

// get a fleet by id
func (c *client) GetFleet(id int64) (*game.Fleet, error) {
	rows := []fleetJoin{}
	if err := c.db.Select(&rows, fmt.Sprintf("%s WHERE f.id = ?", fleetJoinSelect), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// check if we have any results
	if len(rows) == 0 {
		return nil, nil
	}

	fleets := c.scanFleetJoin(rows)

	if len(fleets) != 1 {
		return nil, fmt.Errorf("scan fleetJoin rows")
	}

	fleet := fleets[0]

	// load a fleet's designs
	designUUIDs := make([]uuid.UUID, 0, len(fleet.Tokens))
	for i := range fleet.Tokens {
		designUUIDs = append(designUUIDs, fleet.Tokens[i].DesignUUID)
	}

	// this might be an error case, or we're in a unit test
	if len(designUUIDs) == 0 {
		log.Warn().Int64("ID", fleet.ID).Msg("fleet has no designs associated with tokens")
		return fleet, nil
	}

	designs, err := c.getShipDesignsByUUIDs(designUUIDs)
	if err != nil {
		return nil, fmt.Errorf("get designs by UUIDs -> %w", err)
	}

	designsByUUIDs := make(map[uuid.UUID]*game.ShipDesign, len(designs))
	for i := range designs {
		design := &designs[i]
		designsByUUIDs[design.UUID] = design
	}

	fleet.InjectDesigns(designsByUUIDs)

	return fleet, nil
}

func (c *client) getFleetsForGame(gameId int64) ([]*game.Fleet, error) {
	rows := []fleetJoin{}
	if err := c.db.Select(&rows, fmt.Sprintf("%s WHERE f.gameId = ?", fleetJoinSelect), gameId); err != nil {
		if err == sql.ErrNoRows {
			return []*game.Fleet{}, nil
		}
		return nil, err
	}

	// check if we have any results
	if len(rows) == 0 {
		return []*game.Fleet{}, nil
	}

	fleets := c.scanFleetJoin(rows)

	if len(fleets) == 0 {
		return nil, fmt.Errorf("scan fleetJoin rows")
	}

	return fleets, nil
}

func (c *client) getFleetsForPlayer(playerId int64) ([]*game.Fleet, error) {
	rows := []fleetJoin{}
	if err := c.db.Select(&rows, fmt.Sprintf("%s WHERE f.playerId = ?", fleetJoinSelect), playerId); err != nil {
		if err == sql.ErrNoRows {
			return []*game.Fleet{}, nil
		}
		return nil, err
	}

	// check if we have any results
	if len(rows) == 0 {
		return []*game.Fleet{}, nil
	}

	fleets := c.scanFleetJoin(rows)

	if len(fleets) == 0 {
		return nil, fmt.Errorf("scan fleetJoin rows")
	}

	return fleets, nil
}

// create a new game
func (c *client) createFleet(fleet *game.Fleet, tx SQLExecer) error {
	item := c.converter.ConvertGameFleet(fleet)
	result, err := tx.NamedExec(`
	INSERT INTO fleets (
		createdAt,
		updatedAt,
		gameId,
		playerId,
		battlePlanName,
		x,
		y,
		name,
		num,
		playerNum,
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
		:playerId,
		:battlePlanName,
		:x,
		:y,
		:name,
		:num,
		:playerNum,
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

	for i := range fleet.Tokens {
		token := &fleet.Tokens[i]
		token.FleetID = fleet.ID
		if err := c.createShipToken(token, tx); err != nil {
			return fmt.Errorf("create ShipToken %w", err)
		}
	}

	return nil
}

func (c *client) createShipToken(token *game.ShipToken, tx SQLExecer) error {
	result, err := tx.NamedExec(`
	INSERT INTO shipTokens (
		createdAt,
		updatedAt,
		fleetId,
		designUuid,
		quantity,
		damage,
		quantityDamaged
	)
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:fleetId,
		:designUuid,
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
func (c *client) updateFleet(fleet *game.Fleet, tx SQLExecer) error {

	item := c.converter.ConvertGameFleet(fleet)

	if _, err := tx.NamedExec(`
	UPDATE fleets SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		playerId = :playerId,
		battlePlanName = :battlePlanName,
		x = :x,
		y = :y,
		name = :name,
		num = :num,
		playerNum = :playerNum,
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

	if _, err := tx.Exec("DELETE FROM shipTokens WHERE fleetId = ?", fleet.ID); err != nil {
		return fmt.Errorf("delete existing shipTokens for fleet %w", err)
	}

	for i := range fleet.Tokens {
		token := &fleet.Tokens[i]
		token.ID = 0
		token.FleetID = fleet.ID
		if err := c.createShipToken(token, tx); err != nil {
			return fmt.Errorf("create ShipToken %w", err)
		}
	}
	return nil
}
