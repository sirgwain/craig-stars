package generated

import (
	"database/sql/driver"
)

// we json serialize these types with custom Scan/Value methods
type PlayerNoteTags []string

// db serializer to serialize this to JSON
func (item *PlayerNoteTags) Value() (driver.Value, error) {
	if item == nil || *item == nil {
		return "[]", nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *PlayerNoteTags) Scan(src interface{}) error {
	return scanJSON(src, item)
}
