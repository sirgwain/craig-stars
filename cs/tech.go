package cs

import (
	"fmt"
	"math"
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

// get actual armor/shield value for a tech item given its shield/armor amounts and the multipliers for each
func getActualArmorAmount(armor, shield float64, qty int, raceSpec RaceSpec, isArmor bool) (float64, float64) {
	// TODO: Fix RS shield effect in a less janky way
	if isArmor {
		return armor * raceSpec.ArmorStrengthFactor * float64(qty), shield * raceSpec.ShieldStrengthFactor * float64(qty)
	} else {
		return armor * float64(qty), shield * raceSpec.ShieldStrengthFactor * float64(qty)
	}
}

// compare 2 mining robots and determine which is more resource efficient (ie higher return on investment)
func (hc *TechHullComponent) getBestMiningRobot(player *Player, other *TechHullComponent) bool {
	return float64(other.MiningRate)/float64(hc.MiningRate) > getResourceEfficiencyRatio(player, other, hc) ||
		(float64(other.MiningRate)/float64(hc.MiningRate) == getResourceEfficiencyRatio(player, other, hc) &&
			other.Ranking > hc.Ranking)
}

// compare 2 stargates and determine which one is better
// 1st priority mass, 2nd priority distance, 3rd priority ranking
func (hc *TechHullComponent) getBestStargate(other *TechHullComponent) bool {
	if hc != other {
		switch {
		case hc.SafeHullMass < other.SafeHullMass:
			return false
		case hc.SafeHullMass > other.SafeHullMass:
			return true
		case hc.SafeRange < other.SafeRange: // same safe mass; check safe distance
			return false
		case hc.SafeRange > other.SafeRange:
			return true
		case hc.Ranking > other.Ranking: // same distance & range; compare ranking
			return true
		}
	}
	return false
}

// return the better of the 2 provided torpedo weapons.
// Defaults to 1st if both are equal or either one is null and breaks ties by item ranking
func (hc *TechHullComponent) getBestTorpedo(player *Player, other *TechHullComponent) bool {
	var hcPower float64
	var otherPower float64
	capMissileMulti := 1.5
	// TODO: Figure out a value multi for capital ship missiles that makes sense
	// missiles do 2x damage on shieldless foes, but enemies don't always have shields down

	hcPower = float64(hc.Power) * float64(hc.Accuracy+10) / 100 // add a bit of accuracy bonus to account for computing vs jamming
	if hc.CapitalShipMissile {
		hcPower *= capMissileMulti
	}
	otherPower = float64(other.Power) * float64(other.Accuracy+10) / 100
	if other.CapitalShipMissile {
		otherPower *= capMissileMulti
	}

	// if the 2nd torpedo is more cost efficient (in terms of avg damage/minerals spent) than the new weapon, use it
	return otherPower/hcPower > getCostEfficiencyRatio(player, hc, other, false) ||
		otherPower/hcPower == getCostEfficiencyRatio(player, hc, other, false) && other.Ranking >= hc.Ranking
}

// Returns the ratio of mineral efficiency for 2 TechHullComponents by totaling the tech's Cost struct
// (numeratorTotal / denominatorTotal)
//
// resource indicates whether to consider resources in cost analysis (if false, only minerals will be tallied)
func getCostEfficiencyRatio(player *Player, numerator, denominator *TechHullComponent, resource bool) float64 {
	hcCost := numerator.GetPlayerCost(player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset)
	otherCost := denominator.GetPlayerCost(player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset)
	if !resource {
		// zero out resource costs cuz we don't need em
		hcCost.Resources = 0
		otherCost.Resources = 0
	}
	return float64(hcCost.Total()) / float64(otherCost.Total())
}

// Returns the ratio of resource expenditure for 2 TechHullComponents
// (numeratorTotal / denominatorTotal)
func getResourceEfficiencyRatio(player *Player, numerator, denominator *TechHullComponent) float64 {
	hcRes := numerator.GetPlayerCost(player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset).Resources
	otherRes := denominator.GetPlayerCost(player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset).Resources
	return float64(hcRes) / float64(otherRes)
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

// returns the better of the 2 engines
func (hc *TechEngine) CompareEngine(player *Player, other *TechEngine, purpose FleetPurpose) *TechEngine {
	tech := hc.TechHullComponent
	otherTech := other.TechHullComponent
	if player.HasTech(&tech.Tech) {
		// colony ships don't want radiating engines if we would lose colonists from it
		if ((purpose == FleetPurposeColonizer || purpose == FleetPurposeColonistFreighter) && tech.Radiating &&
			!(player.Race.ImmuneRad || player.Race.Spec.HabCenter.Rad >= 85)) ||
			otherTech.Ranking > tech.Ranking {
			return other
		}
	}
	return hc
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
		return "mine layer"
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

	// planetary items don't get miniaturization
	if t.Category == TechCategoryPlanetary || t.Category == TechCategoryTerraforming || t.Category == TechCategoryPlanetaryDefense || t.Category == TechCategoryPlanetaryScanner {
		miniaturizationFactor = 1
	}

	// apply any tech cost offsets to our item cost
	cost := t.Cost
	var highestCostMulti float64
	for tag := range t.Tags {
		highestCostMulti = math.Min(1+costOffset[tag], highestCostMulti)
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
