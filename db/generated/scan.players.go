package generated

import (
	"database/sql/driver"

	"github.com/sirgwain/craig-stars/cs"
)

// we json serialize these types with custom Scan/Value methods
type CargoTransfers cs.CargoTransfers
type BattlePlans []cs.BattlePlan
type ProductionPlans []cs.ProductionPlan
type TransportPlans []cs.TransportPlan
type PlayerRelationships []cs.PlayerRelationship
type PlayerMessages []cs.PlayerMessage
type PlayerScores []cs.PlayerScore
type AcquiredTechs map[string]bool
type BattleRecords []cs.BattleRecord
type PlayerIntels []cs.PlayerIntel
type ScoreIntels []cs.ScoreIntel
type PlanetIntels []*cs.Planet
type FleetIntels []*cs.Fleet
type ShipDesignIntels []*cs.ShipDesign
type MineralPacketIntels []*cs.MineralPacket
type SalvageIntels []*cs.Salvage
type MinefieldIntels []*cs.Minefield
type MysteryTraderIntels []*cs.MysteryTrader
type WormholeIntels []*cs.Wormhole
type PlayerRace cs.Race
type PlayerSpec cs.PlayerSpec
type PlayerStats cs.PlayerStats

// db serializer to serialize this to JSON
func (item *CargoTransfers) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *CargoTransfers) Scan(src interface{}) error {
	return scanJSON(src, &item)
}

// db serializer to serialize this to JSON
func (item *BattlePlans) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *BattlePlans) Scan(src interface{}) error {
	return scanJSON(src, &item)
}

// db serializer to serialize this to JSON
func (item *ProductionPlans) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *ProductionPlans) Scan(src interface{}) error {
	return scanJSON(src, &item)
}

// db serializer to serialize this to JSON
func (item *TransportPlans) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *TransportPlans) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// db serializer to serialize this to JSON
func (item *PlayerRace) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *PlayerRace) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// db serializer to serialize this to JSON
func (item *PlayerSpec) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *PlayerSpec) Scan(src interface{}) error {
	return scanJSON(src, item)

}

// db serializer to serialize this to JSON
func (item *PlayerStats) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *PlayerStats) Scan(src interface{}) error {
	// stats are weird. They are
	switch v := src.(type) {
	case []byte:
		if len(v) == 4 && string(v) == "null" {
			return nil
		}
	}
	return scanJSON(src, item)
}

func (item *PlayerRelationships) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *PlayerRelationships) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *PlayerMessages) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *PlayerMessages) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *PlayerScores) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *PlayerScores) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *AcquiredTechs) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *AcquiredTechs) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *BattleRecords) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *BattleRecords) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *PlayerIntels) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *PlayerIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *ScoreIntels) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *ScoreIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *PlanetIntels) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *PlanetIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *FleetIntels) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *FleetIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *ShipDesignIntels) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *ShipDesignIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *MineralPacketIntels) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *MineralPacketIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *SalvageIntels) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *SalvageIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *MinefieldIntels) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *MinefieldIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *MysteryTraderIntels) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *MysteryTraderIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *WormholeIntels) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return valueJSON(item)
}

func (item *WormholeIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}
