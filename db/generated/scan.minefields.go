package generated

import (
	"database/sql/driver"

	"github.com/sirgwain/craig-stars/cs"
)

type MinefieldSpec cs.MinefieldSpec

// db serializer to serialize this to JSON
func (item *MinefieldSpec) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *MinefieldSpec) Scan(src interface{}) error {
	return scanJSON(src, item)
}
