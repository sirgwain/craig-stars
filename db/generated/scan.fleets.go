package generated

import (
	"database/sql/driver"

	"github.com/sirgwain/craig-stars/cs"
)

// we json serialize these types with custom Scan/Value methods
type ShipTokens []cs.ShipToken
type Waypoints []cs.Waypoint
type FleetSpec cs.FleetSpec
type Tags cs.Tags

// db serializer to serialize this to JSON
func (item *ShipTokens) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *ShipTokens) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// db serializer to serialize this to JSON
func (item *Waypoints) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *Waypoints) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// db serializer to serialize this to JSON
func (item *FleetSpec) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *FleetSpec) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// db serializer to serialize this to JSON
func (item *Tags) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *Tags) Scan(src interface{}) error {
	return scanJSON(src, item)
}
