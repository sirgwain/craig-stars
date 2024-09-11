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
	TechTagArmor              TechTag = "Armor"
	TechTagTorpedoBonus       TechTag = "TorpedoBonus"
	TechTagBeamCapacitor      TechTag = "BeamCapacitor"
	TechTagBeamDeflector      TechTag = "BeamDeflector"
	TechTagBeamWeapon         TechTag = "BeamWeapon"
	TechTagBomb               TechTag = "Bomb"
	TechTagCapitalShipMissile TechTag = "CapitalShipMissile"
	TechTagCargoPod           TechTag = "CargoPod"
	TechTagCloak              TechTag = "Cloak"
	TechTagColonyModule       TechTag = "ColonyModule"
	TechTagEngine             TechTag = "Engine"
	TechTagFuelTank           TechTag = "FuelTank"
	TechTagGatlingGun         TechTag = "GatlingGun"
	TechTagHeavyMineLayer     TechTag = "HeavyMineLayer"
	TechTagInitiativeBonus    TechTag = "InitiativeBonus"
	TechTagMassDriver         TechTag = "MassDriver"
	TechTagManeuveringJet     TechTag = "ManeuveringJet"
	TechTagMineLayer          TechTag = "MineLayer"
	TechTagMiningRobot        TechTag = "MiningRobot"
	TechTagPenScanner         TechTag = "PenScanner"
	TechTagRamscoop           TechTag = "Ramscoop"
	TechTagTerraformingRobot  TechTag = "TerraformingRobot"
	TechTagScanner            TechTag = "Scanner"
	TechTagShield             TechTag = "Shield"
	TechTagShieldSapper       TechTag = "ShieldSapper"
	TechTagSmartBomb          TechTag = "SmartBomb"
	TechTagSpeedMineLayer     TechTag = "SpeedMineLayer"
	TechTagStargate           TechTag = "Stargate"
	TechTagStructureBomb      TechTag = "StructureBomb"
	TechTagTerraforming       TechTag = "Terraforming"
	TechTagTorpedo            TechTag = "Torpedo"
	TechTagTorpedoJammer      TechTag = "TorpedoJammer"
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
func (tt TechTags) hasTag(tags ...TechTag) bool {
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

// compare 2 techHullComponents based on a field determined by the specified TechTag
// Precedence is given to other component in case of tie
//
// light denotes to penalize heavy armors/shields in favor of lighter ones
func (hc *TechHullComponent) CompareFieldsByTag(player *Player, other *TechHullComponent, tag TechTag, light bool) bool {
	if other == nil {
		return false
	} else if hc == nil {
		return true
	}
	switch tag {
	case TechTagArmor, TechTagShield:
		var hcArmor, hcShield, otherArmor, otherShield float64
		if tag == TechTagArmor {
			hcArmor = float64(hc.Armor) * player.Race.Spec.ArmorStrengthFactor
			hcShield = float64(hc.Shield) // bonus shield from armor components unaffected by bonuses
			otherArmor = float64(other.Armor) * player.Race.Spec.ArmorStrengthFactor
			otherShield = float64(other.Shield)
		} else {
			hcArmor = float64(hc.Armor) // bonus armor from shield components unaffected by bonuses
			hcShield = float64(hc.Shield) * player.Race.Spec.ShieldStrengthFactor
			otherArmor = float64(other.Armor)
			otherShield = float64(other.Shield) * player.Race.Spec.ShieldStrengthFactor
		}
		if light {
			return (otherArmor+otherShield)/(1+math.Min(float64(other.Mass-30)/10, 0)) /
				(hcArmor+hcShield)/(1+math.Min(float64(hc.Mass-30)/10, 0)) >= 
				getMineralEfficiencyRatio(player, hc, other, true)
		}
		return (otherArmor + otherShield) / (hcArmor + hcShield) >= getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagBeamCapacitor:
		return other.BeamBonus >= hc.BeamBonus
	case TechTagBeamDeflector:
		return other.BeamDefense >= hc.BeamDefense
	case TechTagScanner:
		if hc.ScanRangePen > 0 {
			if other.ScanRangePen > 0 {
				return other.ScanRangePen >= hc.ScanRangePen
			} else {
				// 2nd tech doesn't pen scan; 1st wins by default
				return false
			}
		} else if other.ScanRangePen > 0 {
			// 1st tech doesn't pen scan; 2nd wins by default
			return true
		}
		return other.ScanRange >= hc.ScanRange
	case TechTagInitiativeBonus:
		return float64(other.InitiativeBonus) / float64(hc.InitiativeBonus) >= getMineralEfficiencyRatio(player, hc, other, true)
		// FOR THE RECORD, this works out to be equivalent to comparing unit prices
	case TechTagTorpedoJammer:
		return other.TorpedoJamming >= hc.TorpedoJamming
	case TechTagBeamWeapon, TechTagShieldSapper, TechTagGatlingGun:
		return float64(other.Power*other.Range) / float64(hc.Power*hc.Range) >= getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagTorpedo, TechTagCapitalShipMissile:
		return hc.GetBestTorpedo(player, other)
	case TechTagColonyModule:
		return getMineralEfficiencyRatio(player, hc, other, true) >= 1
	case TechTagCargoPod:
		return float64(other.CargoBonus) / float64(hc.CargoBonus) >= getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagFuelTank:
		return float64(other.FuelBonus+5*other.FuelGeneration) / float64(hc.FuelBonus+5*hc.FuelGeneration) >= 
		getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagMineLayer, TechTagHeavyMineLayer, TechTagSpeedMineLayer:
		return float64(other.MineLayingRate) / float64(hc.MineLayingRate) >= getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagBomb, TechTagSmartBomb:
		return float64(other.KillRate) / float64(hc.KillRate) >= getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagStructureBomb:
		return float64(other.StructureDestroyRate) / float64(hc.StructureDestroyRate) >+ getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagCloak:
		return float64(other.CloakUnits) / float64(hc.CloakUnits) >+ getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagManeuveringJet:
		return float64(other.MovementBonus) / float64(hc.MovementBonus) >+ getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagMassDriver:
		return other.PacketSpeed >= hc.PacketSpeed 
	case TechTagStargate:
		return hc.GetBestStargate(other)
	case TechTagTerraformingRobot:
		return float64(other.TerraformRate) / float64(hc.TerraformRate) >+ getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagMiningRobot:
		return float64(other.MiningRate) / float64(hc.MiningRate) >+ getMineralEfficiencyRatio(player, hc, other, false)
	}
	return false
}

// compare 2 stargates and determine which one is better
// 1st priority mass, 2nd priority distance, 3rd priority ranking
func (hc *TechHullComponent) GetBestStargate(other *TechHullComponent) bool {
	if hc.Tags.hasTag(TechTagStargate) {
		if other.Tags.hasTag(TechTagStargate) && hc != other {
			if hc.SafeHullMass < other.SafeHullMass {
				return false
			} else if hc.SafeHullMass > other.SafeHullMass {
				return true
			} else if hc.SafeRange < other.SafeRange { // same safe mass; check safe distance
				return false
			} else if hc.SafeRange > other.SafeRange {
				return true
			} else if hc.Ranking > other.Ranking { // same distance & range; compare ranking
				return true
			}
		} else if other.Tags.hasTag(TechTagStargate) {
			// first tech not a stargate; 2nd one wins by default
			return true
		}
	}
	return false
}

// return the better of the 2 provided torpedo weapons
// defaults to 1st if both are equal or either one is null
func (hc *TechHullComponent) GetBestTorpedo(player *Player, other *TechHullComponent) bool {
	var hcPower float64
	var otherPower float64
	empty := false
	if hc != nil && hc.Category == TechCategoryTorpedo && hc.Power > 0 && player.HasTech(&hc.Tech) {
		hcPower = float64(hc.Power * hc.Accuracy)
		if hc.CapitalShipMissile {
			hcPower *= 1.5 // cap missiles do more damage than normal torps, but enemies don't always have shields up
		}
	} else {
		empty = true
	}
	if other != nil && other.Category == TechCategoryTorpedo && other.Power > 0 && player.HasTech(&other.Tech) {
		if empty { // first component not a torpedo; 2nd one wins by default
			return true
		}
		otherPower = float64(other.Power * other.Accuracy)
		if other.CapitalShipMissile {
			// TODO: Figure out a value multi for cap missiles that makes sense
			otherPower *= 1.5 // cap missiles do more damage than normal torps, but enemies don't always have shields up
		}
	} else {
		return false // 1st one wins by default
	}

	// if the 2nd torpedo is more cost efficient (in terms of avg damage/minerals spent) than the new weapon, use it
	return otherPower/hcPower >= getMineralEfficiencyRatio(player, hc, other, false)
}

// Returns the ratio of mineral efficiency for the 2 components based on the highest mineral in either one 
// (numerator/denomiator)
//
// resource indicates whether to consider resources in cost analysis 
func getMineralEfficiencyRatio(player *Player, numerator, denominator *TechHullComponent, resource bool) float64 {
	if resource {
		hcCost := numerator.GetPlayerCost(player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset)
		otherCost := denominator.GetPlayerCost(player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset)
		minType := hcCost.Max(otherCost).HighestAmount()
		return float64(hcCost.GetAmount(minType)) / float64(otherCost.GetAmount(minType))
	}

	hcCost := numerator.GetPlayerCost(player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset).ToMineral()
	otherCost := denominator.GetPlayerCost(player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset).ToMineral()
	minType := hcCost.Max(otherCost).HighestType()
	return float64(hcCost.GetAmount(minType)) / float64(otherCost.GetAmount(minType))
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
		if !((purpose == FleetPurposeColonizer || purpose == FleetPurposeColonistFreighter) && tech.Radiating &&
		!(player.Race.ImmuneRad || player.Race.Spec.HabCenter.Rad >= 85)) && 
		&otherTech == nil || tech.Ranking > otherTech.Ranking {
			return hc
		}
	}
	return other
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
