package cs

import (
	"fmt"
	"math"
	"slices"
	"strings" // only needed for jank temporary solution
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

const (
	OriginNone          string = ""
	OriginMysteryTrader string = "MysteryTrader"
)

type Tech struct {
	Name         string           `json:"name"`
	Cost         Cost             `json:"cost"`
	Requirements TechRequirements `json:"requirements" `
	Ranking      int              `json:"ranking,omitempty"`
	Category     TechCategory     `json:"category,omitempty"`
	Origin       string           `json:"origin,omitempty"`
	Tags         TechTags         `json:"tags,omitempty"`
}

type TechTags map[TechTag]bool

type TechTag string

const (
	TechTagArmor          TechTag = "Armor"
	TechTagTorpedoBonus   TechTag = "Accuracy Bonus"
	TechTagCapacitor      TechTag = "Beam Capacitor"
	TechTagDeflector      TechTag = "Beam Deflector"
	TechTagBeamWeapon     TechTag = "Beam Weapon"
	TechTagMissile        TechTag = "Capital Ship Missile"
	TechTagCargoPod       TechTag = "Cargo Pod"
	TechTagCloak          TechTag = "Cloak"
	TechTagColonyModule   TechTag = "Colony Module"
	TechTagEngine         TechTag = "Engine"
	TechTagFuelPod        TechTag = "Fuel Pod"
	TechTagGatling        TechTag = "Gatling"
	TechTagInitiative     TechTag = "Initiative Bonus"
	TechTagJammer         TechTag = "Torpedo Jammer"
	TechTagMassDriver     TechTag = "Mass Driver"
	TechTagManeuveringJet TechTag = "Maneuvering Jet"
	TechTagMineLayer      TechTag = "Mine Layer"
	TechTagMiningRobot    TechTag = "Remote Mining Robot"
	TechTagRamscoop       TechTag = "Ramscoop"
	TechTagTerraformRobot TechTag = "Orbital Terraforming Module"
	TechTagScanner        TechTag = "Scanner"
	TechTagShield         TechTag = "Shield"
	TechTagSapper         TechTag = "Shield Sapper"
	TechTagStargate       TechTag = "Stargate"
	TechTagTerraforming   TechTag = "Terraforming"
	TechTagTorpedo        TechTag = "Torpedo"
)

// Create a new techTags map from a list of TechTag items
func newTechTags(tags ...TechTag) TechTags {
	var newMap TechTags
	for _, t := range tags {
		newMap[t] = true
	}
	return newMap
}

// returns true if tt has ALL of the specified tags
func (tt TechTags) hasAllTags(tags ...TechTag) bool {
	for _, tag := range tags {
		if !tt[tag] {
			return false
		}
	}
	return true
}

// returns true if tt has AT LEAST 1 of the specified tags
func (tt TechTags) hasOneTag(tags ...TechTag) bool {
	for _, tag := range tags {
		if tt[tag] {
			return true
		}
	}
	return false
}

// returns slice of all tags in tt, sorted alphabetically
func (tt TechTags) GetTags() []string {
	var list []string
	for k, v := range tt {
		if v {
			list = append(list, string(k))
		}
	}
	slices.Sort(list)
	return list
}

type TechRequirements struct {
	TechLevel
	PRTsDenied   []PRT    `json:"prtsDenied,omitempty"`
	LRTsRequired LRT      `json:"lrtsRequired,omitempty"`
	LRTsDenied   LRT      `json:"lrtsDenied,omitempty"`
	PRTsRequired []PRT    `json:"prtsRequired,omitempty"`
	HullsAllowed []string `json:"hullsAllowed,omitempty"`
	HullsDenied  []string `json:"hullsDenied,omitempty"`
	Acquirable   bool     `json:"acquirable,omitempty"`
}

type TechHullComponent struct {
	Tech
	HullSlotType              HullSlotType  `json:"hullSlotType"`
	Mass                      int           `json:"mass,omitempty"`
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
	MovementBonus             int           `json:"movementBonus,omitempty"`
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

type Engine struct {
	IdealSpeed   int     `json:"idealSpeed,omitempty"`
	FreeSpeed    int     `json:"freeSpeed,omitempty"`
	MaxSafeSpeed int     `json:"maxSafeSpeed,omitempty"`
	FuelUsage    [11]int `json:"fuelUsage,omitempty"`
}

type TechEngine struct {
	TechHullComponent
	Engine
}

type TechHull struct {
	Tech
	Type                     TechHullType   `json:"type,omitempty"`
	Mass                     int            `json:"mass,omitempty"`
	Armor                    int            `json:"armor,omitempty"`
	FuelCapacity             int            `json:"fuelCapacity,omitempty"`
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
	BuiltInScanner           bool           `json:"builtInScanner,omitempty"`
	Initiative               int            `json:"initiative,omitempty"`
	RepairBonus              float64        `json:"repairBonus,omitempty"`
	ImmuneToOwnDetonation    bool           `json:"immuneToOwnDetonation,omitempty"`
	RangeBonus               int            `json:"rangeBonus,omitempty"`
	Starbase                 bool           `json:"starbase,omitempty"`
	OrbitalConstructionHull  bool           `json:"orbitalConstructionHull,omitempty"`
	DoubleMineEfficiency     bool           `json:"doubleMineEfficiency,omitempty"`
	MaxPopulation            int            `json:"maxPopulation,omitempty"`
	InnateScanRangePenFactor float64        `json:"innateScanRangePenFactor,omitempty"`
	Slots                    []TechHullSlot `json:"slots,omitempty"`
}

type TechHullSlot struct {
	Type     HullSlotType `json:"type"`
	Capacity int          `json:"capacity"`
	Required bool         `json:"required,omitempty"`
	Position Vector       `json:"position"`
}

type TechHullType string

const (
	TechHullTypeScout                 TechHullType = "Scout"
	TechHullTypeColonizer             TechHullType = "Colonizer"
	TechHullTypeBomber                TechHullType = "Bomber"
	TechHullTypeFighter               TechHullType = "Fighter"
	TechHullTypeCapitalShip           TechHullType = "CapitalShip"
	TechHullTypeFreighter             TechHullType = "Freighter"
	TechHullTypeMultiPurposeFreighter TechHullType = "MultiPurposeFreighter"
	TechHullTypeFuelTransport         TechHullType = "FuelTransport"
	TechHullTypeMiner                 TechHullType = "Miner"
	TechHullTypeMineLayer             TechHullType = "MineLayer"
	TechHullTypeStarbase              TechHullType = "Starbase"
	TechHullTypeOrbitalFort           TechHullType = "OrbitalFort"
)

func (t TechHullType) IsAttackHull() bool {
	return t == TechHullTypeFighter || t == TechHullTypeCapitalShip || t == TechHullTypeMultiPurposeFreighter
}

func (t TechHullType) IsBomber() bool {
	return t == TechHullTypeBomber
}

type HullSlotType Bitmask

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
		return "mine layer"
	case HullSlotTypeOrbitalElectrical:
		return "orbital electrical"
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
	Tech
	ResetPlanet bool `json:"resetPlanet,omitempty"`
}

type TechPlanetaryScanner struct {
	TechPlanetary
	ScanRange    int `json:"scanRange,omitempty"`
	ScanRangePen int `json:"scanRangePen,omitempty"`
}

type Defense struct {
	DefenseCoverage float64 `json:"defenseCoverage,omitempty"`
}

type TechDefense struct {
	TechPlanetary
	Defense
}

type TechTerraform struct {
	Tech
	Ability int              `json:"ability,omitempty"`
	HabType TerraformHabType `json:"habType,omitempty"`
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

func NewTechWithOrigin(name string, cost Cost, requirements TechRequirements, ranking int, category TechCategory, origin string, tags ...TechTag) Tech {
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

// Get baseline cost for this technology given a player's tech levels, minaturization stats & racial cost modifiers
func (t *Tech) GetPlayerCost(techLevels TechLevel, spec MiniaturizationSpec, costOffset TechCostOffset) Cost {
	// figure out miniaturization
	// this is 4% per level above the required tech we have.
	// We count the smallest diff, i.e. if you have
	// tech level 10 energy, 12 bio and the tech costs 9 energy, 4 bio
	// the smallest level difference you have is 1 energy level (not 8 bio levels)

	levelDiff := TechLevel{-1, -1, -1, -1, -1, -1}

	// From the diff between the player level and the requirements, find the lowest difference
	// i.e. 1 energy level in the example above
	numTechLevelsAboveRequired := math.MaxInt
	if t.Requirements.Energy > 0 {
		levelDiff.Energy = techLevels.Energy - t.Requirements.Energy
		numTechLevelsAboveRequired = MinInt(levelDiff.Energy, numTechLevelsAboveRequired)
	}
	if t.Requirements.Weapons > 0 {
		levelDiff.Weapons = techLevels.Weapons - t.Requirements.Weapons
		numTechLevelsAboveRequired = MinInt(levelDiff.Weapons, numTechLevelsAboveRequired)
	}
	if t.Requirements.Propulsion > 0 {
		levelDiff.Propulsion = techLevels.Propulsion - t.Requirements.Propulsion
		numTechLevelsAboveRequired = MinInt(levelDiff.Propulsion, numTechLevelsAboveRequired)
	}
	if t.Requirements.Construction > 0 {
		levelDiff.Construction = techLevels.Construction - t.Requirements.Construction
		numTechLevelsAboveRequired = MinInt(levelDiff.Construction, numTechLevelsAboveRequired)
	}
	if t.Requirements.Electronics > 0 {
		levelDiff.Electronics = techLevels.Electronics - t.Requirements.Electronics
		numTechLevelsAboveRequired = MinInt(levelDiff.Electronics, numTechLevelsAboveRequired)
	}
	if t.Requirements.Biotechnology > 0 {
		levelDiff.Biotechnology = techLevels.Biotechnology - t.Requirements.Biotechnology
		numTechLevelsAboveRequired = MinInt(levelDiff.Biotechnology, numTechLevelsAboveRequired)
	}

	// for starter techs, they are all 0 requirements, so just use our lowest field
	if numTechLevelsAboveRequired == math.MaxInt {
		numTechLevelsAboveRequired = techLevels.Min()
	}

	// As we learn techs, they get cheaper. We start off with full priced techs, but every additional level of research we learn makes
	// techs cost a little less, maxing out at some discount (i.e. 75% or 80% for races with BET)
	miniaturization := math.Min(spec.MiniaturizationMax, spec.MiniaturizationPerLevel*float64(numTechLevelsAboveRequired))
	// New techs cost BET races 2x
	// new techs will have 0 for miniaturization.
	miniaturizationFactor := spec.NewTechCostFactor
	if numTechLevelsAboveRequired > 0 {
		miniaturizationFactor = 1 - miniaturization
	}

	// apply any tech cost offsets
	// TODO: Implement IT 25% gate discount in actually less janky way
	cost := t.Cost
	switch t.Category {
	case TechCategoryEngine:
		cost = cost.MultiplyFloat64(1 + costOffset.Engine)
	case TechCategoryBeamWeapon:
		cost = cost.MultiplyFloat64(1 + costOffset.BeamWeapon)
	case TechCategoryBomb:
		cost = cost.MultiplyFloat64(1 + costOffset.Bomb)
	case TechCategoryTorpedo:
		cost = cost.MultiplyFloat64(1 + costOffset.Torpedo)
	case TechCategoryOrbital:
		if strings.Contains(t.Name, "Stargate") {
			cost = cost.MultiplyFloat64(1 + costOffset.Stargate)
		}
	case TechCategoryTerraforming:
		cost = cost.MultiplyFloat64(1 + costOffset.Terraforming)
	}

	return Cost{
		int(roundHalfDown(float64(cost.Ironium) * miniaturizationFactor)),
		int(roundHalfDown(float64(cost.Boranium) * miniaturizationFactor)),
		int(roundHalfDown(float64(cost.Germanium) * miniaturizationFactor)),
		int(roundHalfDown(float64(cost.Resources) * miniaturizationFactor)),
	}

	// if we are at level 26, a beginner tech would cost (26 * .04)
	// return cost * (1 - Math.Min(.75, .04 * lowestRequiredDiff));
}
