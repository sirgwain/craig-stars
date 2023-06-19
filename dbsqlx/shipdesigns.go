package dbsqlx

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirgwain/craig-stars/game"
)

type ShipDesign struct {
	ID            int64                  `json:"id,omitempty"`
	CreatedAt     time.Time              `json:"createdAt,omitempty"`
	UpdatedAt     time.Time              `json:"updatedAt,omitempty"`
	PlayerID      int64                  `json:"playerId,omitempty"`
	PlayerNum     int                    `json:"playerNum,omitempty"`
	UUID          uuid.UUID              `json:"uuid,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Version       int                    `json:"version,omitempty"`
	Hull          string                 `json:"hull,omitempty"`
	HullSetNumber int                    `json:"hullSetNumber,omitempty"`
	CanDelete     bool                   `json:"canDelete,omitempty"`
	Slots         ShipDesignSlots        `json:"slots,omitempty"`
	Purpose       game.ShipDesignPurpose `json:"purpose,omitempty"`
	Spec          *ShipDesignSpec        `json:"spec,omitempty"`
}

type ShipDesignSlots []game.ShipDesignSlot
type ShipDesignSpec game.ShipDesignSpec

// db serializer to serialize this to JSON
func (item ShipDesignSlots) Value() (driver.Value, error) {
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
func (item ShipDesignSlots) Scan(src interface{}) error {
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
func (item *ShipDesignSpec) Value() (driver.Value, error) {
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
func (item *ShipDesignSpec) Scan(src interface{}) error {
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

func (c *client) GetShipDesignsForPlayer(playerID int64) ([]game.ShipDesign, error) {

	items := []ShipDesign{}
	if err := c.db.Select(&items, `SELECT * FROM shipDesigns WHERE playerId = ?`, playerID); err != nil {
		if err == sql.ErrNoRows {
			return []game.ShipDesign{}, nil
		}
		return nil, err
	}

	result := make([]game.ShipDesign, len(items))

	for i := range items {
		result[i] = *c.converter.ConvertShipDesign(&items[i])
	}

	return result, nil
}

// get a shipDesign by id
func (c *client) GetShipDesign(id int64) (*game.ShipDesign, error) {
	item := ShipDesign{}
	if err := c.db.Get(&item, "SELECT * FROM shipDesigns WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return c.converter.ConvertShipDesign(&item), nil
}

// create a new shipDesign given something that can execute NamedExec (either a DB or )
func (c *client) CreateShipDesign(shipDesign *game.ShipDesign, tx NamedExecer) error {
	item := c.converter.ConvertGameShipDesign(shipDesign)

	result, err := tx.NamedExec(`
	INSERT INTO shipDesigns (
		createdAt,
		updatedAt,
		playerId,
		playerNum,
		uuid,
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
		:playerId,
		:playerNum,
		:uuid,
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

func (c *client) UpdateShipDesign(shipDesign *game.ShipDesign) error {
	return c.updateShipDesign(shipDesign, c.db)
}

// update an existing shipDesign
func (c *client) updateShipDesign(shipDesign *game.ShipDesign, tx NamedExecer) error {

	item := c.converter.ConvertGameShipDesign(shipDesign)

	if _, err := tx.NamedExec(`
	UPDATE shipDesigns SET
		updatedAt = CURRENT_TIMESTAMP,
		playerId = :playerId,
		playerNum = :playerNum,
		uuid = :uuid,
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
func (c *client) DeleteShipDesign(id int64) error {
	if _, err := c.db.Exec("DELETE FROM shipDesigns WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}
