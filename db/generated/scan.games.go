package generated

import (
	"database/sql/driver"

	"github.com/sirgwain/craig-stars/cs"
)

// we json serialize these types with custom Scan/Value methods
type Rules cs.Rules

// db serializer to serialize this to JSON
func (item *Rules) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *Rules) Scan(src interface{}) error {
	return scanJSON(src, item)
}
