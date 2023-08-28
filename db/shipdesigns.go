package db

import (
	"database/sql"
	"database/sql/driver"

	"github.com/jmoiron/sqlx"
	"github.com/sirgwain/craig-stars/cs"
)

type ShipDesign struct {
	ID            int64                `json:"id,omitempty"`
	GameID        int64                `json:"gameId,omitempty"`
	UpdatedAt     sql.NullTime         `json:"updatedAt,omitempty"`
	CreatedAt     sql.NullTime         `json:"createdAt,omitempty"`
	Num           int                  `json:"num,omitempty"`
	PlayerNum     int                  `json:"playerNum,omitempty"`
	Name          string               `json:"name,omitempty"`
	Version       int                  `json:"version,omitempty"`
	Hull          string               `json:"hull,omitempty"`
	HullSetNumber int                  `json:"hullSetNumber,omitempty"`
	CanDelete     bool                 `json:"canDelete,omitempty"`
	Slots         *ShipDesignSlots     `json:"slots,omitempty"`
	Purpose       cs.ShipDesignPurpose `json:"purpose,omitempty"`
	Spec          *ShipDesignSpec      `json:"spec,omitempty"`
}

type ShipDesignSlots []cs.ShipDesignSlot
type ShipDesignSpec cs.ShipDesignSpec

// db serializer to serialize this to JSON
func (item *ShipDesignSlots) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *ShipDesignSlots) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// db serializer to serialize this to JSON
func (item *ShipDesignSpec) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *ShipDesignSpec) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (c *txClient) GetShipDesignsForPlayer(gameID int64, playerNum int) ([]*cs.ShipDesign, error) {

	items := []*ShipDesign{}
	if err := c.db.Select(&items, `SELECT * FROM shipDesigns WHERE gameId = ? AND playerNum = ?`, gameID, playerNum); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.ShipDesign{}, nil
		}
		return nil, err
	}

	result := make([]*cs.ShipDesign, len(items))

	for i := range items {
		result[i] = c.converter.ConvertShipDesign(items[i])
	}

	return result, nil
}

func (c *txClient) getShipDesignsByNums(gameID int64, playerNum int, nums []int) ([]*cs.ShipDesign, error) {

	query, args, err := sqlx.In(`SELECT * FROM shipDesigns WHERE gameId = ? AND playerNum = ? AND num IN (?)`, gameID, playerNum, nums)
	if err != nil {
		return nil, err
	}

	query = c.db.Rebind(query)
	items := []*ShipDesign{}
	if err := c.db.Select(&items, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.ShipDesign{}, nil
		}
		return nil, err
	}

	result := make([]*cs.ShipDesign, len(items))

	for i := range items {
		result[i] = c.converter.ConvertShipDesign(items[i])
	}

	return result, nil
}

// get a shipDesign by id
func (c *txClient) GetShipDesign(id int64) (*cs.ShipDesign, error) {
	item := ShipDesign{}
	if err := c.db.Get(&item, "SELECT * FROM shipDesigns WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return c.converter.ConvertShipDesign(&item), nil
}

// get a shipDesign by id
func (c *txClient) GetShipDesignByNum(gameID int64, playerNum, num int) (*cs.ShipDesign, error) {
	item := ShipDesign{}
	if err := c.db.Get(&item, "SELECT * FROM shipDesigns WHERE gameId = ? AND playerNum = ? AND num = ?", gameID, playerNum, num); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return c.converter.ConvertShipDesign(&item), nil
}

func (c *txClient) CreateShipDesign(shipDesign *cs.ShipDesign) error {
	item := c.converter.ConvertGameShipDesign(shipDesign)

	result, err := c.db.NamedExec(`
	INSERT INTO shipDesigns (
		createdAt,
		updatedAt,
		gameId,
		num,
		playerNum,
		name,
		version,
		hull,
		hullSetNumber,
		canDelete,
		slots,
		purpose,
		spec
	)
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:gameId,
		:num,
		:playerNum,
		:name,
		:version,
		:hull,
		:hullSetNumber,
		:canDelete,
		:slots,
		:purpose,
		:spec
	)
	`, item)

	if err != nil {
		return err
	}

	// update the id of our passed in shipDesign
	shipDesign.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// update an existing shipDesign
func (c *txClient) UpdateShipDesign(shipDesign *cs.ShipDesign) error {

	item := c.converter.ConvertGameShipDesign(shipDesign)

	if _, err := c.db.NamedExec(`
	UPDATE shipDesigns SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		num = :num,
		playerNum = :playerNum,
		name = :name,
		version = :version,
		hull = :hull,
		hullSetNumber = :hullSetNumber,
		canDelete = :canDelete,
		slots = :slots,
		purpose = :purpose,
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// delete a shipDesign by id
func (c *txClient) DeleteShipDesign(id int64) error {
	if _, err := c.db.Exec("DELETE FROM shipDesigns WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}

// delete a ship design and update/delete fleets associated with the design
// this is a separate function so we can do all the db interactions in a single transaction
func (c *txClient) DeleteShipDesignWithFleets(id int64, fleetsToUpdate, fleetsToDelete []*cs.Fleet) error {

	if _, err := c.db.Exec("DELETE FROM shipDesigns WHERE id = ?", id); err != nil {
		return err
	}

	for _, fleet := range fleetsToUpdate {
		if err := c.UpdateFleet(fleet); err != nil {
			return err
		}
	}

	for _, fleet := range fleetsToDelete {
		if err := c.DeleteFleet(fleet.ID); err != nil {
			return err
		}
	}

	return c.db.Commit()
}
