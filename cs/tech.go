package cs

import (
	"fmt"
)

type TechCategory string

const (
	TechCategoryNone             TechCategory = ""
	TechCategoryArmor            TechCategory = "Armor"
	TechCategoryBeamWeapon       TechCategory = "BeamWeapon"
	TechCategoryBomb             TechCategory = "Bomb"
	TechCategoryElectrical       TechCategory = "Electrical"
	TechCategoryEngine           TechCategory = "Engine"
	TechCategoryMechanical       TechCategory = "Mechanical"
	TechCategoryMineLayer        TechCategory = "MineLayer"
	TechCategoryMineRobot        TechCategory = "MineRobot"
	TechCategoryOrbital          TechCategory = "Orbital"
	TechCategoryPlanetary        TechCategory = "Planetary"
	TechCategoryPlanetaryScanner TechCategory = "PlanetaryScanner"
	TechCategoryPlanetaryDefense TechCategory = "PlanetaryDefense"
	TechCategoryScanner          TechCategory = "Scanner"
	TechCategoryShield           TechCategory = "Shield"
	TechCategoryShipHull         TechCategory = "ShipHull"
	TechCategoryStarbaseHull     TechCategory = "StarbaseHull"
	TechCategoryTerraforming     TechCategory = "Terraforming"
	TechCategoryTorpedo          TechCategory = "Torpedo"
)

var TechCategories = []TechCategory{
	TechCategoryArmor,
	TechCategoryBeamWeapon,
	TechCategoryBomb,
	TechCategoryElectrical,
	TechCategoryEngine,
	TechCategoryMechanical,
	TechCategoryMineLayer,
	TechCategoryMineRobot,
	TechCategoryOrbital,
	TechCategoryPlanetary,
	TechCategoryPlanetaryScanner,
	TechCategoryPlanetaryDefense,
	TechCategoryScanner,
	TechCategoryShield,
	TechCategoryShipHull,
	TechCategoryStarbaseHull,
	TechCategoryTerraforming,
	TechCategoryTorpedo,
}

// The basic skeleton of a Tech item, containing name, cost and other essential info.
type Tech struct {
	Name         string           `json:"name"`
	Cost         Cost             `json:"cost"`
	Requirements TechRequirements `json:"requirements" `
	Ranking      int              `json:"ranking"`
	Category     TechCategory     `json:"category"`
	Origin       TechOrigin       `json:"origin,omitempty"`
	Tags         TechTags         `json:"tags,omitempty"`
}

type TechOrigin string

const (
	OriginNone          TechOrigin = ""
	OriginMysteryTrader TechOrigin = "MysteryTrader"
)

type TechRequirements struct {
	TechLevel    `tstype:",extends"`
	PRTsDenied   []PRT    `json:"prtsDenied,omitempty"`
	LRTsRequired LRT      `json:"lrtsRequired,omitempty"`
	LRTsDenied   LRT      `json:"lrtsDenied,omitempty"`
	PRTsRequired []PRT    `json:"prtsRequired,omitempty"`
	HullsAllowed []string `json:"hullsAllowed,omitempty"`
	HullsDenied  []string `json:"hullsDenied,omitempty"`
	Acquirable   bool     `json:"acquirable,omitempty"`
}

type TechHullComponent struct {
	Tech                      `tstype:",extends"`
	HullSlotType              HullSlotType  `json:"hullSlotType"`
	Mass                      int           `json:"mass"`
	Scanner                   bool          `json:"scanner,omitempty"`
	ScanRange                 int           `json:"scanRange,omitempty"`
	ScanRangePen              int           `json:"scanRangePen,omitempty"`
	SafeHullMass              int           `json:"safeHullMass,omitempty"`
	SafeRange                 int           `json:"safeRange,omitempty"`
	MaxHullMass               int           `json:"maxHullMass,omitempty"`
	MaxRange                  int           `json:"maxRange,omitempty"`
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
	BeamDefense               float64       `json:"beamDefense,omitempty"`
	CargoBonus                int           `json:"cargoBonus,omitempty"`
	ColonizationModule        bool          `json:"colonizationModule,omitempty"`
	FuelBonus                 int           `json:"fuelBonus,omitempty"`
	FuelGeneration            int           `json:"fuelGeneration,omitempty"`
	MovementBonus             float64       `json:"movementBonus,omitempty"`
	OrbitalConstructionModule bool          `json:"orbitalConstructionModule,omitempty"`
	Power                     int           `json:"power,omitempty"`
	Range                     int           `json:"range,omitempty"`
	Initiative                int           `json:"initiative,omitempty"`
	Gatling                   bool          `json:"gatling,omitempty"`
	HitsAllTargets            bool          `json:"hitsAllTargets,omitempty"`
	DamageShieldsOnly         bool          `json:"damageShieldsOnly,omitempty"`
	Accuracy                  int           `json:"accuracy,omitempty"`
	CapitalShipMissile        bool          `json:"capitalShipMissile,omitempty"`
	CanJump                   bool          `json:"canJump,omitempty"`
}

// get actual armor/shield values for a tech item given its base shield/armor amounts and the multipliers for each
func getArmorShieldAmounts(baseArmor, baseShield float64, qty int, raceSpec RaceSpec, isArmor bool) (armor, shield float64) {
	// TODO: Fix RS shield effect in a less janky way
	if isArmor {
		return baseArmor * raceSpec.ArmorStrengthFactor * float64(qty), baseShield * raceSpec.ShieldStrengthFactor * float64(qty)
	} else {
		return baseArmor * float64(qty), baseShield * raceSpec.ShieldStrengthFactor * float64(qty)
	}
}

type Engine struct {
	IdealSpeed   int     `json:"idealSpeed"`
	FreeSpeed    int     `json:"freeSpeed"`
	MaxSafeSpeed int     `json:"maxSafeSpeed"`
	FuelUsage    [11]int `json:"fuelUsage"`
}

type TechEngine struct {
	TechHullComponent `tstype:",extends"`
	Engine            `tstype:",extends"`
}

type TechHull struct {
	Tech                     `tstype:",extends"`
	Type                     TechHullType   `json:"type"`
	Mass                     int            `json:"mass"`
	Armor                    int            `json:"armor"`
	Shield                   int            `json:"shield,omitempty"`
	FuelCapacity             int            `json:"fuelCapacity"`
	FuelGeneration           int            `json:"fuelGeneration,omitempty"`
	CargoCapacity            int            `json:"cargoCapacity,omitempty"`
	CargoSlotPosition        Vector         `json:"cargoSlotPosition,omitempty"`
	CargoSlotSize            Vector         `json:"cargoSlotSize,omitempty"`
	CargoSlotCircle          bool           `json:"cargoSlotCircle,omitempty"`
	SpaceDock                int            `json:"spaceDock,omitempty"`
	SpaceDockSlotPosition    Vector         `json:"spaceDockSlotPosition,omitempty"`
	SpaceDockSlotSize        Vector         `json:"spaceDockSlotSize,omitempty"`
	SpaceDockSlotCircle      bool           `json:"spaceDockSlotCircle,omitempty"`
	MineLayingBonus          float64        `json:"mineLayingBonus,omitempty"`
	Initiative               int            `json:"initiative"`
	RepairBonus              float64        `json:"repairBonus,omitempty"`
	ImmuneToOwnDetonation    bool           `json:"immuneToOwnDetonation,omitempty"`
	RangeBonus               int            `json:"rangeBonus,omitempty"`
	Starbase                 bool           `json:"starbase,omitempty"`
	OrbitalConstructionHull  bool           `json:"orbitalConstructionHull,omitempty"`
	BuiltInScanner           bool           `json:"builtInScanner,omitempty"`
	DoubleMineEfficiency     bool           `json:"doubleMineEfficiency,omitempty"`
	MaxPopulation            int            `json:"maxPopulation,omitempty"`
	InnateScanRangePenFactor float64        `json:"innateScanRangePenFactor,omitempty"`
	Slots                    []TechHullSlot `json:"slots"`
}

type TechHullSlot struct {
	Type     HullSlotType `json:"type"`
	Capacity int          `json:"capacity"`
	Required bool         `json:"required,omitempty"`
	Position Vector       `json:"position"`
}

type TechHullType string

const (
	TechHullTypeBomber                TechHullType = "Bomber"
	TechHullTypeColonizer             TechHullType = "Colonizer"
	TechHullTypeCapitalShip           TechHullType = "CapitalShip" // big, bulky capital ships
	TechHullTypeFighter               TechHullType = "Fighter"
	TechHullTypeFreighter             TechHullType = "Freighter"
	TechHullTypeFuelTransport         TechHullType = "FuelTransport"
	TechHullTypeMiner                 TechHullType = "Miner"
	TechHullTypeMineLayer             TechHullType = "MineLayer"
	TechHullTypeMultiPurposeFreighter TechHullType = "MultiPurposeFreighter"
	TechHullTypeOrbitalFort           TechHullType = "OrbitalFort"
	TechHullTypeScout                 TechHullType = "Scout"
	TechHullTypeStarbase              TechHullType = "Starbase"
)

var TechHullTypes = []TechHullType{
	TechHullTypeBomber,
	TechHullTypeColonizer,
	TechHullTypeCapitalShip,
	TechHullTypeFighter,
	TechHullTypeFreighter,
	TechHullTypeFuelTransport,
	TechHullTypeMiner,
	TechHullTypeMineLayer,
	TechHullTypeMultiPurposeFreighter,
	TechHullTypeOrbitalFort,
	TechHullTypeScout,
	TechHullTypeStarbase,
}

func (t TechHullType) IsAttackHull() bool {
	return t == TechHullTypeFighter || t == TechHullTypeCapitalShip || t == TechHullTypeMultiPurposeFreighter
}

func (t TechHullType) IsBomber() bool {
	return t == TechHullTypeBomber
}

type HullSlotType Bitmask

// DO NOT REARRANGE EXISTING VALUES. IT *WILL* BREAK EXISTING GAMES.

const (
	HullSlotTypeNone                = 0
	HullSlotTypeEngine HullSlotType = 1 << iota
	HullSlotTypeScanner
	HullSlotTypeMechanical
	HullSlotTypeBomb
	HullSlotTypeMining
	HullSlotTypeElectrical
	HullSlotTypeShield
	HullSlotTypeArmor
	HullSlotTypeCargo
	HullSlotTypeSpaceDock
	HullSlotTypeWeapon
	HullSlotTypeOrbital
	HullSlotTypeMineLayer

	HullSlotTypeElectricalMechanical             = HullSlotTypeElectrical | HullSlotTypeMechanical
	HullSlotTypeOrbitalElectrical                = HullSlotTypeOrbital | HullSlotTypeElectrical
	HullSlotTypeShieldElectricalMechanical       = HullSlotTypeShield | HullSlotTypeElectrical | HullSlotTypeMechanical
	HullSlotTypeScannerElectricalMechanical      = HullSlotTypeScanner | HullSlotTypeElectrical | HullSlotTypeMechanical
	HullSlotTypeArmorScannerElectricalMechanical = HullSlotTypeArmor | HullSlotTypeScanner | HullSlotTypeElectrical | HullSlotTypeMechanical
	HullSlotTypeMineElectricalMechanical         = HullSlotTypeMineLayer | HullSlotTypeElectrical | HullSlotTypeMechanical
	HullSlotTypeShieldArmor                      = HullSlotTypeShield | HullSlotTypeArmor
	HullSlotTypeWeaponShield                     = HullSlotTypeShield | HullSlotTypeWeapon
	HullSlotTypeGeneral                          = HullSlotTypeScanner | HullSlotTypeMechanical | HullSlotTypeElectrical | HullSlotTypeShield | HullSlotTypeArmor | HullSlotTypeWeapon | HullSlotTypeMineLayer
)

var HullSlotTypes = []HullSlotType{
	HullSlotTypeNone,
	HullSlotTypeEngine,
	HullSlotTypeScanner,
	HullSlotTypeMechanical,
	HullSlotTypeBomb,
	HullSlotTypeMining,
	HullSlotTypeElectrical,
	HullSlotTypeShield,
	HullSlotTypeArmor,
	HullSlotTypeCargo,
	HullSlotTypeSpaceDock,
	HullSlotTypeWeapon,
	HullSlotTypeOrbital,
	HullSlotTypeMineLayer,
	HullSlotTypeElectricalMechanical,
	HullSlotTypeOrbitalElectrical,
	HullSlotTypeShieldElectricalMechanical,
	HullSlotTypeScannerElectricalMechanical,
	HullSlotTypeArmorScannerElectricalMechanical,
	HullSlotTypeMineElectricalMechanical,
	HullSlotTypeShieldArmor,
	HullSlotTypeWeaponShield,
	HullSlotTypeGeneral,
}

// all single type hull slot types
var BasicHullSlotTypes = []HullSlotType{
	HullSlotTypeNone,
	HullSlotTypeEngine,
	HullSlotTypeScanner,
	HullSlotTypeMechanical,
	HullSlotTypeBomb,
	HullSlotTypeMining,
	HullSlotTypeElectrical,
	HullSlotTypeShield,
	HullSlotTypeArmor,
	HullSlotTypeCargo,
	HullSlotTypeSpaceDock,
	HullSlotTypeWeapon,
	HullSlotTypeOrbital,
	HullSlotTypeMineLayer,
}

func (hst HullSlotType) String() string {
	switch hst {
	case HullSlotTypeNone:
		return "none"
	case HullSlotTypeEngine:
		return "engine"
	case HullSlotTypeScanner:
		return "scanner"
	case HullSlotTypeMechanical:
		return "mechanical"
	case HullSlotTypeBomb:
		return "bomb"
	case HullSlotTypeMining:
		return "mining"
	case HullSlotTypeElectrical:
		return "electrical"
	case HullSlotTypeShield:
		return "shield"
	case HullSlotTypeArmor:
		return "armor"
	case HullSlotTypeCargo:
		return "cargo"
	case HullSlotTypeSpaceDock:
		return "space dock"
	case HullSlotTypeWeapon:
		return "weapon"
	case HullSlotTypeOrbital:
		return "orbital"
	case HullSlotTypeMineLayer:
		return "minelayer"
	case HullSlotTypeOrbitalElectrical:
		return "orbital electrical"
	case HullSlotTypeElectricalMechanical:
		return "electrical mechanical"
	case HullSlotTypeShieldElectricalMechanical:
		return "shield electrical mechanical"
	case HullSlotTypeScannerElectricalMechanical:
		return "scanner electrical mechanical"
	case HullSlotTypeArmorScannerElectricalMechanical:
		return "armor scanner electrical mechanical"
	case HullSlotTypeMineElectricalMechanical:
		return "mine electrical mechanical"
	case HullSlotTypeShieldArmor:
		return "shield armor"
	case HullSlotTypeWeaponShield:
		return "weapon shield"
	case HullSlotTypeGeneral:
		return "general"
	default:
		return fmt.Sprintf("unknown slot type (%d)", hst)
	}
}

type TechPlanetary struct {
	Tech        `tstype:",extends"`
	ResetPlanet bool `json:"resetPlanet,omitempty"`
}

type TechPlanetaryScanner struct {
	TechPlanetary `tstype:",extends"`
	ScanRange     int `json:"scanRange"`
	ScanRangePen  int `json:"scanRangePen"`
}

type Defense struct {
	DefenseCoverage float64 `json:"defenseCoverage"`
}

type TechDefense struct {
	TechPlanetary `tstype:",extends"`
	Defense       `tstype:",extends"`
}

type TechTerraform struct {
	Tech    `tstype:",extends"`
	Ability int              `json:"ability"`
	HabType TerraformHabType `json:"habType"`
}

type TerraformHabType string

const (
	TerraformHabTypeNone TerraformHabType = ""
	TerraformHabTypeGrav TerraformHabType = "Grav"
	TerraformHabTypeTemp TerraformHabType = "Temp"
	TerraformHabTypeRad  TerraformHabType = "Rad"
	TerraformHabTypeAll  TerraformHabType = "All"
)

func FromHabType(habType HabType) TerraformHabType {
	switch habType {
	case Grav:
		return TerraformHabTypeGrav
	case Temp:
		return TerraformHabTypeTemp
	case Rad:
		return TerraformHabTypeRad
	default:
		return TerraformHabTypeNone
	}
}

func NewTech(name string, cost Cost, requirements TechRequirements, ranking int, category TechCategory, tags ...TechTag) Tech {
	return Tech{
		Name:         name,
		Cost:         cost,
		Requirements: requirements,
		Ranking:      ranking,
		Category:     category,
		Tags:         newTechTags(tags...),
	}
}

func NewTechWithOrigin(name string, cost Cost, requirements TechRequirements, ranking int, category TechCategory, origin TechOrigin, tags ...TechTag) Tech {
	t := NewTech(name, cost, requirements, ranking, category, tags...)
	t.Origin = origin
	return t
}

func (t *Tech) String() string                 { return t.Name }
func (t *TechHull) String() string             { return t.Name }
func (t *TechHullComponent) String() string    { return t.Name }
func (t *TechEngine) String() string           { return t.Name }
func (t *TechPlanetaryScanner) String() string { return t.Name }
func (t *TechDefense) String() string          { return t.Name }
func (t *TechTerraform) String() string        { return t.Name }
