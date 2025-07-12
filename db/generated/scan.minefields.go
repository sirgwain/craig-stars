package generated

import (
	"database/sql/driver"

	"github.com/sirgwain/craig-stars/cs"
)

type MineFieldSpec cs.MineFieldSpec

// db serializer to serialize this to JSON
func (item *MineFieldSpec) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *MineFieldSpec) Scan(src interface{}) error {
	return scanJSON(src, item)
}
