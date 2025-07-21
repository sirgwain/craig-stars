package generated

import (
	"database/sql/driver"

	"github.com/sirgwain/craig-stars/cs"
)

// we json serialize these types with custom Scan/Value methods
type MysteryTraderSpec cs.MysteryTraderSpec
type MysteryTraderPlayersRewarded map[int]bool

// db serializer to serialize this to JSON
func (item *MysteryTraderSpec) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *MysteryTraderSpec) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// db serializer to serialize this to JSON
func (item *MysteryTraderPlayersRewarded) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *MysteryTraderPlayersRewarded) Scan(src interface{}) error {
	return scanJSON(src, item)
}
