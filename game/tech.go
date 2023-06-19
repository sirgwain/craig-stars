package game

import (
	"time"

	"gorm.io/gorm"
)

type TechField string

const (
	TechFieldEnergy        TechField = "Energy"
	TechFieldWeapons       TechField = "Weapons"
	TechFieldPropulsion    TechField = "Propulsion"
	TechFieldConstruction  TechField = "Construction"
	TechFieldElectronics   TechField = "Electronics"
	TechFieldBiotechnology TechField = "Biotechnology"
)

type TechCategory string

const (
	TechCategoryArmor            TechCategory = "Armor"
	TechCategoryBeamWeapon       TechCategory = "BeamWeapon"
	TechCategoryBomb             TechCategory = "Bomb"
	TechCategoryElectrical       TechCategory = "Electrical"
	TechCategoryEngine           TechCategory = "Engine"
	TechCategoryMechanical       TechCategory = "Mechanical"
	TechCategoryMineLayer        TechCategory = "MineLayer"
	TechCategoryMineRobot        TechCategory = "MineRobot"
	TechCategoryOrbital          TechCategory = "Orbital"
	TechCategoryPlanetaryScanner TechCategory = "PlanetaryScanner"
	TechCategoryPlanetaryDefense TechCategory = "PlanetaryDefense"
	TechCategoryScanner          TechCategory = "Scanner"
	TechCategoryShield           TechCategory = "Shield"
	TechCategoryShipHull         TechCategory = "ShipHull"
	TechCategoryStarbaseHull     TechCategory = "StarbaseHull"
	TechCategoryTerraforming     TechCategory = "Terraforming"
	TechCategoryTorpedo          TechCategory = "Torpedo"
)

type Tech struct {
	ID           uint             `gorm:"primaryKey" json:"id" header:"Username"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedat"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"deletedAt"`
	TechStoreID  uint             `json:"techStoreId"`
	Name         string           `json:"name"`
	Cost         Cost             `json:"cost" gorm:"embedded;embeddedPrefix:cost_"`
	Requirements TechRequirements `json:"requirements"  gorm:"embedded;embeddedPrefix:requirements_"`
	Ranking      int              `json:"ranking,omitempty"`
	Category     TechCategory     `json:"category,omitempty"`
}

type TechRequirements struct {
	TechLevel
	PRTDenied    PRT `json:"prtDenied,omitempty"`
	LRTsRequired LRT `json:"lrtsRequired"`
	LRTsDenied   LRT `json:"lrtsDenied"`
	PRTRequired  PRT `json:"prtRequired,omitempty"`
}

type TechHullComponent struct {
	Tech
	HullSlotType              HullSlotType  `json:"hullSlotType"`
	Mass                      int           `json:"mass,omitempty"`
	ScanRange                 int           `json:"scanRange"`
	ScanRangePen              int           `json:"scanRangePen"`
	SafeHullMass              int           `json:"safeHullMass"`
	SafeRange                 int           `json:"safeRange"`
	MaxHullMass               int           `json:"maxHullMass"`
	MaxRange                  int           `json:"maxRange"`
	Radiating                 bool          `json:"radiating,omitempty"`
	PacketSpeed               int           `json:"packetSpeed,omitempty"`
	CloakUnits                int           `json:"cloakUnits,omitempty"`
	TerraformRate             int           `json:"terraformRate,omitempty"`
	MiningRate                int           `json:"miningRate,omitempty"`
	KillRate                  float64       `json:"killRate,omitempty"`
	MinKillRate               int           `json:"minKillRate,omitempty"`
	StructureDestroyRate      float64       `json:"structureDestroyRate,omitempty"`
	UnterraformRate           int           `json:"unterraformRate,omitempty"`
	Smart                     bool          `json:"smart,omitempty"`
	CanStealFleetCargo        bool          `json:"canStealFleetCargo,omitempty"`
	CanStealPlanetCargo       bool          `json:"canStealPlanetCargo,omitempty"`
	Armor                     int           `json:"armor,omitempty"`
	Shield                    int           `json:"shield,omitempty"`
	TorpedoBonus              float64       `json:"torpedoBonus,omitempty"`
	InitiativeBonus           int           `json:"initiativeBonus,omitempty"`
	BeamBonus                 float64       `json:"beamBonus,omitempty"`
	ReduceMovement            int           `json:"reduceMovement,omitempty"`
	TorpedoJamming            float64       `json:"torpedoJamming,omitempty"`
	ReduceCloaking            bool          `json:"reduceCloaking,omitempty"`
	CloakUnarmedOnly          bool          `json:"cloakUnarmedOnly,omitempty"`
	MineFieldType             MineFieldType `json:"mineFieldType,omitempty"`
	MineLayingRate            int           `json:"mineLayingRate,omitempty"`
	BeamDefense               int           `json:"beamDefense,omitempty"`
	CargoBonus                int           `json:"cargoBonus,omitempty"`
	ColonizationModule        bool          `json:"colonizationModule,omitempty"`
	FuelBonus                 int           `json:"fuelBonus,omitempty"`
	MovementBonus             int           `json:"movementBonus,omitempty"`
	OrbitalConstructionModule bool          `json:"orbitalConstructionModule,omitempty"`
	Power                     int           `json:"power,omitempty"`
	Range                     int           `json:"range,omitempty"`
	Initiative                int           `json:"initiative,omitempty"`
	Gattling                  bool          `json:"gattling,omitempty"`
	HitsAllTargets            bool          `json:"hitsAllTargets,omitempty"`
	DamageShieldsOnly         bool          `json:"damageShieldsOnly,omitempty"`
	FuelRegenerationRate      int           `json:"fuelRegenerationRate,omitempty"`
	Accuracy                  int           `json:"accuracy,omitempty"`
	CapitalShipMissile        bool          `json:"capitalShipMissile,omitempty"`
}

type TechEngine struct {
	TechHullComponent
	IdealSpeed int     `json:"idealSpeed,omitempty"`
	FreeSpeed  int     `json:"freeSpeed,omitempty"`
	FuelUsage  [11]int `json:"fuelUsage,omitempty" gorm:"serializer:json"`
}

type TechHull struct {
	Tech
	Type                     TechHullType   `json:"type,omitempty"`
	Mass                     int            `json:"mass,omitempty"`
	Armor                    int            `json:"armor,omitempty"`
	FuelCapacity             int            `json:"fuelCapacity,omitempty"`
	CargoCapacity            int            `json:"cargoCapacity,omitempty"`
	SpaceDock                int            `json:"spaceDock,omitempty"`
	MineLayingFactor         int            `json:"mineLayingFactor,omitempty"`
	BuiltInScanner           bool           `json:"builtInScanner,omitempty"`
	Initiative               int            `json:"initiative,omitempty"`
	RepairBonus              float64        `json:"repairBonus,omitempty"`
	ImmuneToOwnDetonation    bool           `json:"immuneToOwnDetonation,omitempty"`
	RangeBonus               int            `json:"rangeBonus,omitempty"`
	Starbase                 bool           `json:"starbase,omitempty"`
	OrbitalConstructionHull  bool           `json:"orbitalConstructionHull,omitempty"`
	DoubleMineEfficiency     bool           `json:"doubleMineEfficiency,omitempty"`
	InnateScanRangePenFactor float64        `json:"innateScanRangePenFactor,omitempty"`
	Slots                    []TechHullSlot `json:"slots,omitempty" gorm:"serializer:json"`
}

type TechHullSlot struct {
	Type     HullSlotType `json:"type"`
	Capacity int          `json:"capacity"`
	Required bool         `json:"required,omitempty"`
}

type TechHullType string

const (
	TechHullTypeScout          TechHullType = "Scout"
	TechHullTypeColonizer      TechHullType = "Colonizer"
	TechHullTypeBomber         TechHullType = "Bomber"
	TechHullTypeFighter        TechHullType = "Fighter"
	TechHullTypeCapitalShip    TechHullType = "CapitalShip"
	TechHullTypeFreighter      TechHullType = "Freighter"
	TechHullTypeArmedFreighter TechHullType = "ArmedFreighter"
	TechHullTypeFuelTransport  TechHullType = "FuelTransport"
	TechHullTypeMiner          TechHullType = "Miner"
	TechHullTypeMineLayer      TechHullType = "MineLayer"
	TechHullTypeStarbase       TechHullType = "Starbase"
)

type HullSlotType Bitmask

const (
	HullSlotTypeArmor HullSlotType = 1 << iota
	HullSlotTypeBomb
	HullSlotTypeElectrical
	HullSlotTypeEngine
	HullSlotTypeMechanical
	HullSlotTypeMineLayer
	HullSlotTypeOrbital
	HullSlotTypeScanner
	HullSlotTypeShield
	HullSlotTypeMining
	HullSlotTypeWeapon

	HullSlotTypeOrbitalElectrical                = HullSlotTypeOrbital | HullSlotTypeElectrical
	HullSlotTypeShieldElectricalMechanical       = HullSlotTypeShield | HullSlotTypeElectrical | HullSlotTypeMechanical
	HullSlotTypeScannerElectricalMechanical      = HullSlotTypeScanner | HullSlotTypeElectrical | HullSlotTypeMechanical
	HullSlotTypeArmorScannerElectricalMechanical = HullSlotTypeArmor | HullSlotTypeScanner | HullSlotTypeElectrical | HullSlotTypeMechanical
	HullSlotTypeMineElectricalMechanical         = HullSlotTypeMineLayer | HullSlotTypeElectrical | HullSlotTypeMechanical
	HullSlotTypeShieldArmor                      = HullSlotTypeShield | HullSlotTypeArmor
	HullSlotTypeWeaponShield                     = HullSlotTypeShield | HullSlotTypeWeapon
	HullSlotTypeGeneral                          = HullSlotTypeScanner | HullSlotTypeMechanical | HullSlotTypeElectrical | HullSlotTypeShield | HullSlotTypeArmor | HullSlotTypeWeapon | HullSlotTypeMineLayer
)

type TechPlanetaryScanner struct {
	Tech
	ScanRange    int `json:"scanRange,omitempty"`
	ScanRangePen int `json:"scanRangePen,omitempty"`
}

type TechDefense struct {
	Tech
	DefenseCoverage float64 `json:"defenseCoverage,omitempty"`
}

type TechTerraform struct {
	Tech
	Ability int              `json:"ability,omitempty"`
	HabType TerraformHabType `json:"habType,omitempty"`
}

type TerraformHabType string

const (
	TerraformHabTypeGravity     TerraformHabType = "Gravity"
	TerraformHabTypeTemperature TerraformHabType = "Temperature"
	TerraformHabTypeRadiation   TerraformHabType = "Radiation"
	TerraformHabTypeAll         TerraformHabType = "All"
)

func NewTech(name string, cost Cost, requirements TechRequirements, ranking int, category TechCategory) Tech {
	return Tech{
		Name:         name,
		Cost:         cost,
		Requirements: requirements,
		Ranking:      ranking,
		Category:     category,
	}
}
