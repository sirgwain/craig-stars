package generated

import (
	"database/sql/driver"

	"github.com/sirgwain/craig-stars/cs"
)

// we json serialize these types with custom Scan/Value methods
type ProductionQueueItems []cs.ProductionQueueItem
type PlanetSpec cs.PlanetSpec

// db serializer to serialize this to JSON
func (item *ProductionQueueItems) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *ProductionQueueItems) Scan(src interface{}) error {
	return scanJSON(src, &item)

}

// db serializer to serialize this to JSON
func (item *PlanetSpec) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *PlanetSpec) Scan(src interface{}) error {
	return scanJSON(src, &item)
}
