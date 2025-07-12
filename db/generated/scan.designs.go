package generated

import (
	"database/sql/driver"

	"github.com/sirgwain/craig-stars/cs"
)

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
