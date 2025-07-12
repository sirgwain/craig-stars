package generated

import (
	"database/sql/driver"

	"github.com/sirgwain/craig-stars/cs"
)

// we json serialize these types with custom Scan/Value methods
type WormholeSpec cs.WormholeSpec

// db serializer to serialize this to JSON
func (item *WormholeSpec) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *WormholeSpec) Scan(src interface{}) error {
	return scanJSON(src, item)
}
