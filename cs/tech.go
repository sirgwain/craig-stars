package cs

import (
	"math"
	"time"
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
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	TechStoreID  int64            `json:"techStoreId"`
	Name         string           `json:"name"`
	Cost         Cost             `json:"cost"`
	Requirements TechRequirements `json:"requirements" `
	Ranking      int              `json:"ranking,omitempty"`
	Category     TechCategory     `json:"category,omitempty"`
}

type TechRequirements struct {
	TechLevel
	PRTDenied    PRT `json:"prtDenied,omitempty"`
	LRTsRequired LRT `json:"lrtsRequired,omitempty"`
	LRTsDenied   LRT `json:"lrtsDenied,omitempty"`
	PRTRequired  PRT `json:"prtRequired,omitempty"`
}

type TechHullComponent struct {
	Tech
	HullSlotType              HullSlotType  `json:"hullSlotType"`
	Mass                      int           `json:"mass,omitempty"`
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
	FuelUsage  [11]int `json:"fuelUsage,omitempty"`
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
	Slots                    []TechHullSlot `json:"slots,omitempty"`
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

func (t *Tech) String() string                 { return t.Name }
func (t *TechHull) String() string             { return t.Name }
func (t *TechHullComponent) String() string    { return t.Name }
func (t *TechEngine) String() string           { return t.Name }
func (t *TechPlanetaryScanner) String() string { return t.Name }
func (t *TechDefense) String() string          { return t.Name }
func (t *TechTerraform) String() string        { return t.Name }

func (t *Tech) GetPlayerCost(player *Player) Cost {
	// figure out miniaturization
	// this is 4% per level above the required tech we have.
	// We count the smallest diff, i.e. if you have
	// tech level 10 energy, 12 bio and the tech costs 9 energy, 4 bio
	// the smallest level difference you have is 1 energy level (not 8 bio levels)

	levelDiff := TechLevel{-1, -1, -1, -1, -1, -1}

	// From the diff between the player level and the requirements, find the lowest difference
	// i.e. 1 energey level in the example above
	numTechLevelsAboveRequired := math.MaxInt
	if t.Requirements.Energy > 0 {
		levelDiff.Energy = player.TechLevels.Energy - t.Requirements.Energy
		numTechLevelsAboveRequired = minInt(levelDiff.Energy, numTechLevelsAboveRequired)
	}
	if t.Requirements.Weapons > 0 {
		levelDiff.Weapons = player.TechLevels.Weapons - t.Requirements.Weapons
		numTechLevelsAboveRequired = minInt(levelDiff.Weapons, numTechLevelsAboveRequired)
	}
	if t.Requirements.Propulsion > 0 {
		levelDiff.Propulsion = player.TechLevels.Propulsion - t.Requirements.Propulsion
		numTechLevelsAboveRequired = minInt(levelDiff.Propulsion, numTechLevelsAboveRequired)
	}
	if t.Requirements.Construction > 0 {
		levelDiff.Construction = player.TechLevels.Construction - t.Requirements.Construction
		numTechLevelsAboveRequired = minInt(levelDiff.Construction, numTechLevelsAboveRequired)
	}
	if t.Requirements.Electronics > 0 {
		levelDiff.Electronics = player.TechLevels.Electronics - t.Requirements.Electronics
		numTechLevelsAboveRequired = minInt(levelDiff.Electronics, numTechLevelsAboveRequired)
	}
	if t.Requirements.Biotechnology > 0 {
		levelDiff.Biotechnology = player.TechLevels.Biotechnology - t.Requirements.Biotechnology
		numTechLevelsAboveRequired = minInt(levelDiff.Biotechnology, numTechLevelsAboveRequired)
	}

	// for starter techs, they are all 0 requirements, so just use our lowest field
	if numTechLevelsAboveRequired == math.MaxInt {
		numTechLevelsAboveRequired = player.TechLevels.Min()
	}

	// As we learn techs, they get cheaper. We start off with full priced techs, but every additional level of research we learn makes
	// techs cost a little less, maxing out at some discount (i.e. 75% or 80% for races with BET)
	miniaturization := math.Min(player.Race.Spec.MiniaturizationMax, player.Race.Spec.MiniaturizationPerLevel*float64(numTechLevelsAboveRequired))
	// New techs cost BET races 2x
	// new techs will have 0 for miniaturization.
	miniaturizationFactor := player.Race.Spec.NewTechCostFactor
	if numTechLevelsAboveRequired > 0 {
		miniaturizationFactor = 1 - miniaturization
	}

	return Cost{
		int(math.Ceil(float64(t.Cost.Ironium) * miniaturizationFactor)),
		int(math.Ceil(float64(t.Cost.Boranium) * miniaturizationFactor)),
		int(math.Ceil(float64(t.Cost.Germanium) * miniaturizationFactor)),
		int(math.Ceil(float64(t.Cost.Resources) * miniaturizationFactor)),
	}
	// if we are at level 26, a beginner tech would cost (26 * .04)
	// return cost * (1 - Math.Min(.75, .04 * lowestRequiredDiff));
}
