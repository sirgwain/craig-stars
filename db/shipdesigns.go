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
	CannotDelete  bool                 `json:"cannotDelete,omitempty"`
	Slots         *ShipDesignSlots     `json:"slots,omitempty"`
	Purpose       cs.ShipDesignPurpose `json:"purpose,omitempty"`
	Spec          *ShipDesignSpec      `json:"spec,omitempty"`
	CanDelete     sql.NullBool         `json:"canDelete,omitempty"` // unused
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

func (c *client) GetShipDesignsForPlayer(gameID int64, playerNum int) ([]*cs.ShipDesign, error) {

	items := []*ShipDesign{}
	if err := c.reader.Select(&items, `SELECT * FROM shipDesigns WHERE gameId = ? AND playerNum = ?`, gameID, playerNum); err != nil {
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

func (c *client) getShipDesignsByNums(gameID int64, playerNum int, nums []int) ([]*cs.ShipDesign, error) {

	query, args, err := sqlx.In(`SELECT * FROM shipDesigns WHERE gameId = ? AND playerNum = ? AND num IN (?)`, gameID, playerNum, nums)
	if err != nil {
		return nil, err
	}

	query = c.reader.Rebind(query)
	items := []*ShipDesign{}
	if err := c.reader.Select(&items, query, args...); err != nil {
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
func (c *client) GetShipDesign(id int64) (*cs.ShipDesign, error) {
	item := ShipDesign{}
	if err := c.reader.Get(&item, "SELECT * FROM shipDesigns WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return c.converter.ConvertShipDesign(&item), nil
}

// get a shipDesign by id
func (c *client) GetShipDesignByNum(gameID int64, playerNum, num int) (*cs.ShipDesign, error) {
	item := ShipDesign{}
	if err := c.reader.Get(&item, "SELECT * FROM shipDesigns WHERE gameId = ? AND playerNum = ? AND num = ?", gameID, playerNum, num); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return c.converter.ConvertShipDesign(&item), nil
}

func (c *client) CreateShipDesign(shipDesign *cs.ShipDesign) error {
	item := c.converter.ConvertGameShipDesign(shipDesign)

	result, err := c.writer.NamedExec(`
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
		cannotDelete,
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
		:cannotDelete,
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
func (c *client) UpdateShipDesign(shipDesign *cs.ShipDesign) error {

	item := c.converter.ConvertGameShipDesign(shipDesign)

	if _, err := c.writer.NamedExec(`
	UPDATE shipDesigns SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		num = :num,
		playerNum = :playerNum,
		name = :name,
		version = :version,
		hull = :hull,
		hullSetNumber = :hullSetNumber,
		cannotDelete = :cannotDelete,
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
func (c *client) DeleteShipDesign(id int64) error {
	if _, err := c.writer.Exec("DELETE FROM shipDesigns WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}

