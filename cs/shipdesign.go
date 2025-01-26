package cs

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

// Fleets are made up of ships, and each ship has a design. Players start with designs created
// during universe generation, and they can add new designs in the UI.
// Deleting a design deletes all fleets associated with it.
type ShipDesign struct {
	GameDBObject
	Num               int               `json:"num,omitempty"`
	PlayerNum         int               `json:"playerNum"`
	OriginalPlayerNum int               `json:"originalPlayerNum"`
	Name              string            `json:"name"`
	Version           int               `json:"version"`
	Hull              string            `json:"hull"`
	HullSetNumber     int               `json:"hullSetNumber"`
	CannotDelete      bool              `json:"cannotDelete,omitempty"`
	MysteryTrader     bool              `json:"mysteryTrader,omitempty"`
	Slots             []ShipDesignSlot  `json:"slots"`
	Purpose           ShipDesignPurpose `json:"purpose,omitempty"`
	Spec              ShipDesignSpec    `json:"spec"`
	Delete            bool              `json:"-"` // used by the AI to mark a design for deletion
}

type ShipDesignSlot struct {
	HullComponent string `json:"hullComponent"`
	HullSlotIndex int    `json:"hullSlotIndex"`
	Quantity      int    `json:"quantity"`
}

type ShipDesignSpec struct {
	AdditionalMassDrivers     int                   `json:"additionalMassDrivers,omitempty"`
	Armor                     int                   `json:"armor,omitempty"`
	BasePacketSpeed           int                   `json:"basePacketSpeed,omitempty"`
	BeamBonus                 float64               `json:"beamBonus,omitempty"`
	BeamDefense               float64               `json:"beamDefense,omitempty"`
	Bomber                    bool                  `json:"bomber,omitempty"`
	Bombs                     []Bomb                `json:"bombs,omitempty"`
	CanJump                   bool                  `json:"canJump,omitempty"`
	CanLayMines               bool                  `json:"canLayMines,omitempty"`
	CanStealFleetCargo        bool                  `json:"canStealFleetCargo,omitempty"`
	CanStealPlanetCargo       bool                  `json:"canStealPlanetCargo,omitempty"`
	CargoCapacity             int                   `json:"cargoCapacity,omitempty"`
	CloakPercent              int                   `json:"cloakPercent,omitempty"`
	CloakPercentFullCargo     int                   `json:"cloakPercentFullCargo,omitempty"`
	CloakUnits                int                   `json:"cloakUnits,omitempty"`
	Colonizer                 bool                  `json:"colonizer,omitempty"`
	Cost                      Cost                  `json:"cost,omitempty"`
	Engine                    Engine                `json:"engine,omitempty"`
	EstimatedRange            int                   `json:"estimatedRange,omitempty"`
	EstimatedRangeFull        int                   `json:"estimatedRangeFull,omitempty"`
	FuelCapacity              int                   `json:"fuelCapacity,omitempty"`
	FuelGeneration            int                   `json:"fuelGeneration,omitempty"`
	HasWeapons                bool                  `json:"hasWeapons,omitempty"`
	HullType                  TechHullType          `json:"hullType,omitempty"`
	ImmuneToOwnDetonation     bool                  `json:"immuneToOwnDetonation,omitempty"`
	Initiative                int                   `json:"initiative,omitempty"`
	InnateScanRangePenFactor  float64               `json:"innateScanRangePenFactor,omitempty"`
	Mass                      int                   `json:"mass,omitempty"`
	MassDriver                string                `json:"massDriver,omitempty"`
	MaxHullMass               int                   `json:"maxHullMass,omitempty"`
	MaxPopulation             int                   `json:"maxPopulation,omitempty"`
	MaxRange                  int                   `json:"maxRange,omitempty"`
	MineLayingRateByMineType  map[MineFieldType]int `json:"mineLayingRateByMineType,omitempty"`
	MineSweep                 int                   `json:"mineSweep,omitempty"`
	MiningRate                int                   `json:"miningRate,omitempty"`
	Movement                  int                   `json:"movement,omitempty"`
	MovementBonus             float64               `json:"movementBonus,omitempty"`
	MovementFull              int                   `json:"movementFull,omitempty"`
	NumBuilt                  int                   `json:"numBuilt,omitempty"`
	NumEngines                int                   `json:"numEngines,omitempty"`
	NumInstances              int                   `json:"numInstances,omitempty"`
	OrbitalConstructionModule bool                  `json:"orbitalConstructionModule,omitempty"`
	PowerRating               int                   `json:"powerRating,omitempty"`
	Radiating                 bool                  `json:"radiating,omitempty"`
	ReduceCloaking            float64               `json:"reduceCloaking,omitempty"`
	ReduceMovement            int                   `json:"reduceMovement,omitempty"`
	RepairBonus               float64               `json:"repairBonus,omitempty"`
	RetroBombs                []Bomb                `json:"retroBombs,omitempty"`
	SafeHullMass              int                   `json:"safeHullMass,omitempty"`
	SafePacketSpeed           int                   `json:"safePacketSpeed,omitempty"`
	SafeRange                 int                   `json:"safeRange,omitempty"`
	Scanner                   bool                  `json:"scanner,omitempty"`
	ScanRange                 int                   `json:"scanRange,omitempty"`
	ScanRangePen              int                   `json:"scanRangePen,omitempty"`
	Shields                   int                   `json:"shields,omitempty"`
	SmartBombs                []Bomb                `json:"smartBombs,omitempty"`
	SpaceDock                 int                   `json:"spaceDock,omitempty"`
	Starbase                  bool                  `json:"starbase,omitempty"`
	Stargate                  string                `json:"stargate,omitempty"`
	TechLevel                 TechLevel             `json:"techLevel,omitempty"`
	TerraformRate             int                   `json:"terraformRate,omitempty"`
	TorpedoBonus              float64               `json:"torpedoBonus,omitempty"`
	TorpedoJamming            float64               `json:"torpedoJamming,omitempty"`
	WeaponSlots               []ShipDesignSlot      `json:"weaponSlots,omitempty"`
}

type ShipDesignPurpose string

const (
	ShipDesignPurposeNone                  ShipDesignPurpose = ""
	ShipDesignPurposeScout                 ShipDesignPurpose = "Scout"
	ShipDesignPurposeColonizer             ShipDesignPurpose = "Colonizer"
	ShipDesignPurposeBomber                ShipDesignPurpose = "Bomber"
	ShipDesignPurposeStructureBomber       ShipDesignPurpose = "StructureBomber"
	ShipDesignPurposeSmartBomber           ShipDesignPurpose = "SmartBomber"
	ShipDesignPurposeStartingFighter       ShipDesignPurpose = "StartingFighter" // only used for starting designs
	ShipDesignPurposeFighterScout          ShipDesignPurpose = "FighterScout"    // armed beam scouts
	ShipDesignPurposeTorpedoFighter        ShipDesignPurpose = "TorpedoFighter"  // torpedo/missile boats
	ShipDesignPurposeBeamFighter           ShipDesignPurpose = "BeamFighter"     // beam/sapper boats
	ShipDesignPurposeFreighter             ShipDesignPurpose = "Freighter"
	ShipDesignPurposeColonistFreighter     ShipDesignPurpose = "ColonistFreighter"
	ShipDesignPurposeFuelFreighter         ShipDesignPurpose = "FuelFreighter"
	ShipDesignPurposeMultiPurposeFreighter ShipDesignPurpose = "MultiPurposeFreighter"
	ShipDesignPurposeArmedFreighter        ShipDesignPurpose = "ArmedFreighter"
	ShipDesignPurposeMiner                 ShipDesignPurpose = "Miner"
	ShipDesignPurposeTerraformer           ShipDesignPurpose = "Terraformer"
	ShipDesignPurposeDamageMineLayer       ShipDesignPurpose = "DamageMineLayer"
	ShipDesignPurposeSpeedMineLayer        ShipDesignPurpose = "SpeedMineLayer"
	ShipDesignPurposeStarbase              ShipDesignPurpose = "Starbase"
	ShipDesignPurposeStarbaseUnarmed       ShipDesignPurpose = "StarbaseUnarmed"
	ShipDesignPurposeFuelDepot             ShipDesignPurpose = "FuelDepot"
	ShipDesignPurposeStarbaseQuarter       ShipDesignPurpose = "StarbaseQuarter"
	ShipDesignPurposeStarbaseHalf          ShipDesignPurpose = "StarbaseHalf"
	ShipDesignPurposePacketThrower         ShipDesignPurpose = "PacketThrower"
	ShipDesignPurposeStargater             ShipDesignPurpose = "Stargater"
	ShipDesignPurposeFort                  ShipDesignPurpose = "Fort"
	ShipDesignPurposeStarterColony         ShipDesignPurpose = "StarterColony"
)

func NewShipDesign(player *Player, num int) *ShipDesign {
	return &ShipDesign{PlayerNum: player.Num, Num: num, Slots: []ShipDesignSlot{}}
}

func (sd *ShipDesign) WithName(name string) *ShipDesign {
	sd.Name = name
	return sd
}
func (sd *ShipDesign) WithHull(hull string) *ShipDesign {
	sd.Hull = hull
	return sd
}
func (sd *ShipDesign) WithSlots(slots []ShipDesignSlot) *ShipDesign {
	sd.Slots = slots
	return sd
}

func (sd *ShipDesign) WithPurpose(purpose ShipDesignPurpose) *ShipDesign {
	sd.Purpose = purpose
	return sd
}
func (sd *ShipDesign) WithHullSetNumber(num int) *ShipDesign {
	sd.HullSetNumber = num
	return sd
}

// Compute the spec for this ShipDesign. This function is mostly for universe generation and tests.
//
// See [ComputeShipDesignSpec]
func (sd *ShipDesign) WithSpec(rules *Rules, player *Player) *ShipDesign {
	var err error
	sd.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, sd)
	if err != nil {
		panic(fmt.Sprintf("ComputeShipDesignSpec returned error %v", err))
	}
	return sd
}

// validate that this ship design is correct and is avaliable for the given player
func (sd *ShipDesign) Validate(rules *Rules, player *Player) error {
	// Basic design/hull checks
	if strings.TrimSpace(sd.Name) == "" {
		return fmt.Errorf("design has no name")
	}
	hull := rules.techs.GetHull(sd.Hull)
	if hull == nil {
		return fmt.Errorf("hull %q was not found in techStore", sd.Hull)
	}
	if !player.HasTech(&hull.Tech) {
		return fmt.Errorf("hull %q is not available to player", hull)
	}

	// slot index checks
	for index, designSlot := range sd.Slots {
		hullSlot := hull.Slots[Clamp(designSlot.HullSlotIndex, 1, len(hull.Slots))-1] // prevents index out of range for nil errors & lets us use nice switch statement
		switch {
		case designSlot.HullSlotIndex <= 0:
			return fmt.Errorf("design slot #%d's HullSlotIndex is 0 or less (%d)", index, designSlot.HullSlotIndex)
		case designSlot.HullSlotIndex > len(hull.Slots):
			return fmt.Errorf("design slot #%d's HullSlotIndex is out of range (%d vs %d)", index, designSlot.HullSlotIndex, len(hull.Slots))
		case designSlot.Quantity < 0:
			return fmt.Errorf("design slot #%d has a negative number of components (%d)", index, designSlot.Quantity)
		case designSlot.Quantity > hullSlot.Capacity:
			return fmt.Errorf("design slot #%d has more components than the hull slot could hold (%d vs %d)", index, designSlot.Quantity, hullSlot.Capacity)
		case hullSlot.Required && designSlot.Quantity < hullSlot.Capacity:
			return fmt.Errorf("design slot #%d has too few components for a mandatory hull slot (%d vs %d)", index, designSlot.Quantity, hullSlot.Capacity)
		case designSlot.HullComponent != "":
			// if we have a hull component, check it
			hc := rules.techs.GetHullComponent(designSlot.HullComponent)
			if hc == nil {
				return fmt.Errorf("hull component %q was not found in tech store", designSlot.HullComponent)
			}

			if hullSlot.Type&hc.HullSlotType == 0 {
				return fmt.Errorf("hull component %q cannot be placed in %s slot", hc, hullSlot.Type)
			}

			if len(hc.Requirements.HullsAllowed) > 0 && slices.Contains(hc.Requirements.HullsAllowed, hull.Name) {
				return fmt.Errorf("hull component %q is not usable on hull %s", hc, hull)
			}

			if len(hc.Requirements.HullsDenied) > 0 && !slices.Contains(hc.Requirements.HullsDenied, hull.Name) {
				return fmt.Errorf("hull component %q is forbidden on hull %s", hc, hull)
			}

			if !player.HasTech(&hc.Tech) {
				return fmt.Errorf("hull component %s is not available to player", hc)
			}
		}

	}

	// check required slots to make sure they're filled properly
	for i, hullSlot := range hull.Slots {
		if hullSlot.Required {
			found := false
			for _, slot := range sd.Slots {
				if slot.HullSlotIndex-1 == i && slot.HullComponent != "" {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("%d %s required", hullSlot.Capacity, hullSlot.Type.String())
			}
		}
	}

	return nil
}

// compare two ship design's slots and return true if they are equal
func (d *ShipDesign) SlotsEqual(otherSlots []ShipDesignSlot) bool {
	if len(d.Slots) != len(otherSlots) {
		return false
	}
	for i, v := range d.Slots {
		if v != otherSlots[i] {
			return false
		}
	}
	return true
}

// return true if this ship's purpose requires it to be light
// (ie beam warships & peacetime ships)
func (p ShipDesignPurpose) IsLightShip() bool {
	// all peacetime ships should not be using armor
	// too heavy and they're gonna die anyways
	return (!p.IsWarship() && p != ShipDesignPurposeStartingFighter) || p.IsBeamShip()
}

// return true if this ship's purpose is to be some kind of warship
// (starting fighter not included)
func (p ShipDesignPurpose) IsWarship() bool {
	return p.IsBeamShip() || p.IsTorpedoShip()
}

// return true if this ship's purpose is to use beams
func (p ShipDesignPurpose) IsBeamShip() bool {
	return p == ShipDesignPurposeBeamFighter ||
		p == ShipDesignPurposeFighterScout ||
		p == ShipDesignPurposeArmedFreighter
}

// return true if this ship's job is to use torpedoes
func (p ShipDesignPurpose) IsTorpedoShip() bool {
	return p == ShipDesignPurposeTorpedoFighter ||
		p == ShipDesignPurposeStarbase ||
		p == ShipDesignPurposeStarbaseHalf ||
		p == ShipDesignPurposeStarbaseQuarter ||
		p == ShipDesignPurposeFort
}

// get the movement for this ship design, based on cargoMass
func (d *ShipDesign) getMovement(rules *Rules, cargoMass int) int {
	return getBattleMovement(rules.MovementMin, rules.MovementMax, d.Spec.Engine.IdealSpeed, float64(d.Spec.MovementBonus), d.Spec.Mass+cargoMass, d.Spec.NumEngines)
}

// returns the new jamming/computing bonus
func getNewJamming(prevBonus, componentBonus, multi float64, qty int) float64 {
	baseMulti := 1 - prevBonus/multi // undo multi before multiplication
	compMulti := math.Pow(1-componentBonus, float64(qty))
	return (1 - baseMulti*compMulti) * multi
}

// returns the new beam defense factor after adding the given components
func getNewBeamBonus(prevBonus, componentBonus float64, qty int) float64 {
	return prevBonus * math.Pow(1+componentBonus, float64(qty))
}

// Compute a ship design's Spec
func ComputeShipDesignSpec(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design *ShipDesign) (ShipDesignSpec, error) {
	hull := rules.techs.GetHull(design.Hull)
	if hull == nil {
		return ShipDesignSpec{}, fmt.Errorf("failed to find hull %s in techstore", design.Hull)
	}
	c := NewCostCalculator()
	spec := ShipDesignSpec{
		Mass:                     hull.Mass,
		Armor:                    hull.Armor,
		Shields:                  hull.Shield,
		FuelCapacity:             hull.FuelCapacity,
		FuelGeneration:           hull.FuelGeneration,
		Cost:                     Cost{}, // will assign cost later with error handling
		TechLevel:                hull.Requirements.TechLevel,
		CargoCapacity:            hull.CargoCapacity,
		CloakUnits:               raceSpec.BuiltInCloakUnits,
		Initiative:               hull.Initiative,
		ImmuneToOwnDetonation:    hull.ImmuneToOwnDetonation,
		MovementBonus:            raceSpec.MovementBonus,
		RepairBonus:              hull.RepairBonus,
		ScanRange:                0, // by default, all ships non-pen scan ships in their radius (ie at their position)
		ScanRangePen:             NoScanner,
		SpaceDock:                hull.SpaceDock,
		Starbase:                 hull.Starbase,
		MaxPopulation:            hull.MaxPopulation,
		HullType:                 hull.Type,
		InnateScanRangePenFactor: hull.InnateScanRangePenFactor,
	}

	var err error
	spec.Cost, err = c.GetDesignCost(rules, techLevels, raceSpec, design)
	if err != nil {
		return ShipDesignSpec{}, fmt.Errorf("failed to get design cost; error %w", err)
	}

	// count the number of each type of battle component we have
	torpedoBonusesByCount := map[float64]int{}
	torpedoJammersByCount := map[float64]int{}
	beamBoostersByCount := map[float64]int{}
	beamDeflectorsByCount := map[float64]int{}

	numTachyonDetectors := 0

	var armor, shield float64

	// rating calcs
	beamPower := 0
	torpedoPower := 0
	bombsPower := 0

	for i, slot := range design.Slots {
		if slot.Quantity > 0 {
			component := rules.techs.GetHullComponent(slot.HullComponent)
			if component == nil || component.Name == "" {
				// assume slot is empty; cut it out and carry on
				design.Slots = append(design.Slots[:i], design.Slots[i+1:]...)
				continue
			}
			hullSlot := hull.Slots[slot.HullSlotIndex-1]

			// record engine details
			// TODO: Add support for multiple engine "slots"
			// (all would have to share the same engine type)
			if hullSlot.Type == HullSlotTypeEngine {
				engine := rules.techs.GetEngine(slot.HullComponent)
				spec.Engine = engine.Engine
				spec.NumEngines = slot.Quantity
			}

			if component.Category == TechCategoryBeamWeapon && component.Power > 0 && (component.Range+hull.RangeBonus) > 0 {
				// mines swept/yr = power * (range)^2
				gatlingMultiplier := 1
				if component.Gatling {
					// gattlings are 4x more mine-sweepery (all gatlings have range of 2; 2^2=4)
					gatlingMultiplier = component.Range * component.Range
				}
				spec.MineSweep += slot.Quantity * component.Power * ((component.Range + hull.RangeBonus) * component.Range) * gatlingMultiplier
			}

			spec.TechLevel = spec.TechLevel.Max(component.Requirements.TechLevel)

			spec.Mass += component.Mass * slot.Quantity
			a, s := getArmorShieldAmounts(float64(component.Armor), float64(component.Shield), slot.Quantity, raceSpec, component.Category == TechCategoryArmor)
			armor += a
			shield += s
			spec.CargoCapacity += component.CargoBonus * slot.Quantity
			spec.FuelCapacity += component.FuelBonus * slot.Quantity
			spec.FuelGeneration += component.FuelGeneration * slot.Quantity
			spec.Colonizer = spec.Colonizer || component.ColonizationModule || component.OrbitalConstructionModule
			spec.Initiative += component.InitiativeBonus * slot.Quantity
			spec.MovementBonus += component.MovementBonus * float64(slot.Quantity)
			spec.ReduceMovement = Max(spec.ReduceMovement, component.ReduceMovement) // these don't stack
			spec.MiningRate += component.MiningRate * slot.Quantity
			spec.TerraformRate += component.TerraformRate * slot.Quantity
			spec.OrbitalConstructionModule = spec.OrbitalConstructionModule || component.OrbitalConstructionModule
			spec.CanStealFleetCargo = spec.CanStealFleetCargo || component.CanStealFleetCargo
			spec.CanStealPlanetCargo = spec.CanStealPlanetCargo || component.CanStealPlanetCargo
			spec.CanJump = spec.CanJump || component.CanJump
			spec.Radiating = spec.Radiating || component.Radiating

			// Add this mine type to the layers this design has
			if component.MineLayingRate > 0 {
				spec.CanLayMines = true
				if spec.MineLayingRateByMineType == nil {
					spec.MineLayingRateByMineType = make(map[MineFieldType]int)
				}
				if _, ok := spec.MineLayingRateByMineType[component.MineFieldType]; !ok {
					spec.MineLayingRateByMineType[component.MineFieldType] = 0
				}
				spec.MineLayingRateByMineType[component.MineFieldType] += int(float64(component.MineLayingRate) * float64(slot.Quantity) * (1 + hull.MineLayingBonus))
			}

			// count battle computers, jammers, capacitors & deflectors
			if component.TorpedoBonus > 0 {
				torpedoBonusesByCount[component.TorpedoBonus] += slot.Quantity
			}
			if component.TorpedoJamming > 0 {
				torpedoJammersByCount[component.TorpedoJamming] += slot.Quantity
			}
			if component.BeamBonus > 0 {
				beamBoostersByCount[component.BeamBonus] += slot.Quantity
			}
			if component.BeamDefense > 0 {
				beamDeflectorsByCount[component.BeamDefense] += slot.Quantity
			}

			// if this slot has a bomb, this design is a bomber
			if component.HullSlotType == HullSlotTypeBomb || component.MinKillRate > 0 || component.KillRate > 0 || component.StructureDestroyRate > 0 || component.UnterraformRate > 0 {
				spec.Bomber = true
				bomb := Bomb{
					Quantity:             slot.Quantity,
					KillRate:             component.KillRate,
					MinKillRate:          component.MinKillRate,
					StructureDestroyRate: component.StructureDestroyRate,
					UnterraformRate:      component.UnterraformRate,
				}
				if component.UnterraformRate > 0 {
					spec.RetroBombs = append(spec.RetroBombs, bomb)
				} else if component.Smart {
					spec.SmartBombs = append(spec.SmartBombs, bomb)
				} else {
					spec.Bombs = append(spec.Bombs, bomb)
				}

				// bombs add to rating
				bombsPower += int((bomb.KillRate*10 + bomb.StructureDestroyRate)) * slot.Quantity * 2
			}

			if component.Power > 0 {
				spec.HasWeapons = true
				spec.WeaponSlots = append(spec.WeaponSlots, slot)
				switch component.Category {
				case TechCategoryBeamWeapon:
					// beams contribute to the rating based on range, but sappers
					// are 1/3rd rated to compensate for high power
					rating := component.Power * slot.Quantity * (component.Range + 3) / 4
					if component.DamageShieldsOnly {
						rating /= 3
					}
					beamPower += rating
				case TechCategoryTorpedo:
					torpedoPower += component.Power * slot.Quantity * (component.Range - 2) / 2
				}
			}

			// cloaking
			if component.CloakUnits > 0 {
				spec.CloakUnits += component.CloakUnits * slot.Quantity
			}
			if component.ReduceCloaking {
				numTachyonDetectors++
			}
			// cargo and space dock that are built into the hull
			// the space dock assumes that there is only one slot like that
			// it won't add them up

			if hullSlot.Type&HullSlotTypeSpaceDock > 0 {
				spec.SpaceDock = hullSlot.Capacity
			}

			// mass drivers
			if component.PacketSpeed > 0 {
				// if we already have a mass driver at this speed, add an additional mass driver to up
				// our speed
				if spec.BasePacketSpeed == component.PacketSpeed {
					spec.AdditionalMassDrivers++
				}
				spec.BasePacketSpeed = Max(spec.BasePacketSpeed, component.PacketSpeed)
				spec.MassDriver = component.Name
			}

			// stargate fields
			if component.SafeHullMass != 0 {
				spec.Stargate = component.Name
				spec.SafeHullMass = component.SafeHullMass
			}
			if component.MaxHullMass != 0 {
				spec.MaxHullMass = component.MaxHullMass
			}
			if component.SafeRange != 0 {
				spec.SafeRange = component.SafeRange
			}
			if component.MaxRange != 0 {
				spec.MaxRange = component.MaxRange
			}
		}
	}

	spec.Armor += int(armor)
	spec.Shields += int(shield)

	// ISB gives some special starbase bonuses
	// Discount is already handled in cost function
	if hull.Starbase {
		spec.CloakUnits += raceSpec.BuiltInCloakUnits
	}

	// determine the safe speed for this design
	spec.SafePacketSpeed = spec.BasePacketSpeed + spec.AdditionalMassDrivers

	// figure out the cloak as a percentage after we spend our cloak units
	spec.CloakPercent = getCloakPercentForCloakUnits(spec.CloakUnits)
	spec.CloakPercentFullCargo = getCloakPercentForCloakUnits(int(math.Round(float64(spec.CloakUnits) * float64(spec.Mass) / float64(spec.Mass+spec.CargoCapacity))))

	if numTachyonDetectors > 0 {
		// 95% ^ (SQRT(#_of_detectors) = reduction factor for other players' cloaks (capped at 81% or 17TDs)
		spec.ReduceCloaking = math.Min(math.Pow((1-rules.TachyonCloakReduction), math.Sqrt(float64(numTachyonDetectors))), rules.TachyonMaxCloakReduction)
	} else {
		spec.ReduceCloaking = 1
	}

	// Calculate final bonuses for computing, jamming, capacitating & jamming
	// TODO: Benchmark these and swap to new functions if faster
	if len(torpedoBonusesByCount) > 0 {
		spec.TorpedoBonus = 1
		for torpedoBonus, count := range torpedoBonusesByCount {
			// for 3 Battle Computer 30s, this calc is 1-(.7^3) or 65%
			bonus := 1 - math.Pow(1-torpedoBonus, float64(count))

			// if there are multiple battle computer slots all working together, they multiply together
			// 1−((1−BC20Bonus)×(1−BC30Bonus)×(1−BC50Bonus))
			spec.TorpedoBonus *= 1 - bonus
		}

		// the final bonus is the above sum inverted
		spec.TorpedoBonus = 1 - spec.TorpedoBonus

		// golang, why you be like this? nobody wants 1-.2^1 to be .199999994
		spec.TorpedoBonus = roundFloat(spec.TorpedoBonus, 4)
	}

	if len(torpedoJammersByCount) > 0 {
		spec.TorpedoJamming = 1
		for torpedoJammer, count := range torpedoJammersByCount {
			// for 3 Jammer 10s, this calc is 1-(.9^3) or 27.1%
			jammer := 1 - math.Pow(1-torpedoJammer, float64(count))

			// if there are multiple jammer slots all working together, they multiply together
			// 1−((1−Jammer10)×(1−Jammer20)×(1−Jammer30))
			spec.TorpedoJamming *= 1 - jammer
		}

		// the final jam anount is the above sum inverted
		spec.TorpedoJamming = 1 - spec.TorpedoJamming

		// round off answer and apply relevant caps/multipliers
		// golang, why you be like this? nobody wants 1-.2^1 to be .199999994
		spec.TorpedoJamming = roundFloat(math.Min(spec.TorpedoJamming,
			rules.JammerCap.Get(hull.Starbase))*
			rules.JammerMulti.Get(hull.Starbase), 4)

	}

	// beam bonus defaults to 1
	spec.BeamBonus = 1
	if len(beamBoostersByCount) > 0 {
		for beamBonus, count := range beamBoostersByCount {
			// for 3 flux caps, this calc is 1-(1.2^3) for 1.728x beam damage
			bonus := math.Pow(1+beamBonus, float64(count))

			// multiple beam boosters stack multiplicatively
			spec.BeamBonus *= bonus

			if spec.BeamBonus > rules.BeamBonusCap {
				// save a bit of computing power by breaking early if over cap
				break
			}
		}

		// Return final % bonus, rounded to 4 decimal places and capped at 2.55x base damage
		spec.BeamBonus = math.Min(roundFloat(spec.BeamBonus, 4), rules.BeamBonusCap)
	}

	if len(beamDeflectorsByCount) > 0 {
		spec.BeamDefense = 1
		for beamDefense, count := range beamDeflectorsByCount {
			// for 3 deflectors, this calc is 1-(0.9^3) for 0.729x beam damage taken
			bonus := math.Pow(1-beamDefense, float64(count))

			// multiple beam deflectors stack multiplicatively
			spec.BeamDefense *= bonus
		}

		// Return final % dmg reduction, rounded to 4 decimal places
		spec.BeamDefense = roundFloat(spec.BeamDefense, 4)
	}

	if spec.NumEngines > 0 {
		// Movement = (IdealEngineSpeed - 2) - (Mass / (70 * NumEngines)) + Move Bonus
		// move bonus is rounded up before evaluation
		spec.Movement = getBattleMovement(rules.MovementMin, rules.MovementMax, spec.Engine.IdealSpeed, spec.MovementBonus, spec.Mass, spec.NumEngines)
		spec.MovementFull = getBattleMovement(rules.MovementMin, rules.MovementMax, spec.Engine.IdealSpeed, spec.MovementBonus, spec.Mass+spec.CargoCapacity, spec.NumEngines)
	} else {
		spec.Movement = 0
		spec.MovementFull = 0
	}

	beamPower = int(float64(beamPower) * (spec.BeamBonus))
	if beamPower > 0 {
		// starbases don't move, but for the beam power calcs
		// assume they have a movement of "2" which is the lowest possible
		movement := Clamp(spec.Movement, rules.MovementMin, rules.MovementMax)

		// a movement of 1 1/2 in the UI (halfwar between max & min) doesn't impact your beam
		// power rating. Anything less reduces it, anything higher increases it
		beamPower += (beamPower * (movement - (rules.MovementMin+rules.MovementMax)/2)) / rules.MovementMax
	}
	spec.PowerRating = beamPower + torpedoPower + bombsPower

	spec.computeScanRanges(rules, raceSpec.ScannerSpec, techLevels, design, hull)

	// compute the estimated range for this design
	if spec.NumEngines > 0 {
		fuelCostFor1kly := spec.Engine.getFuelCostForEngine(spec.Engine.IdealSpeed, spec.Mass, 1000, 1+raceSpec.FuelEfficiencyOffset)
		fuelCostFor1klyFull := spec.Engine.getFuelCostForEngine(spec.Engine.IdealSpeed, spec.Mass+spec.CargoCapacity, 1000, 1+raceSpec.FuelEfficiencyOffset)

		if fuelCostFor1kly == 0 {
			spec.EstimatedRange = Infinite
			spec.EstimatedRangeFull = Infinite
		} else {
			spec.EstimatedRange = int(float64(spec.FuelCapacity) / float64(fuelCostFor1kly) * 1000)
			spec.EstimatedRangeFull = int(float64(spec.FuelCapacity) / float64(fuelCostFor1klyFull) * 1000)
		}
	}

	return spec, nil
}

// Compute the scan ranges for this ship design.
//
// Formula: (scanner1^4 + scanner2^4 + ...
// + scannerN^4)^(0.25)
func (spec *ShipDesignSpec) computeScanRanges(rules *Rules, scannerSpec ScannerSpec, techLevels TechLevel, design *ShipDesign, hull *TechHull) {
	spec.ScanRange = 0
	spec.ScanRangePen = 0
	hasPenScan := false // counter to track if we have a pen scanner or not

	// compute built in scanner if hull allows for it
	if hull.BuiltInScanner {
		builtInScanner := scannerSpec.BuiltInScanner
		builtInNormal := builtInScanner.NormalMulti.Multiply(techLevels).Total()
		builtInPen := builtInScanner.PenMulti.Multiply(techLevels).Total()
		if builtInNormal > 0 {
			spec.ScanRange = PowInt(builtInNormal, 4)
		}
		if builtInPen > 0 {
			spec.ScanRangePen = PowInt(builtInPen, 4)
			hasPenScan = spec.ScanRangePen > 0
		}
	}

	// loop through slots to add scan ranges up
	for _, slot := range design.Slots {
		if slot.Quantity == 0 {
			continue
		}

		component := rules.techs.GetHullComponent(slot.HullComponent)
		if component == nil || !component.Scanner {
			continue
		}

		// Add (scanrange)^4 to our tally for both normal and pen scans
		if component.ScanRange != NoScanner {
			spec.ScanRange += PowInt(component.ScanRange, 4) * slot.Quantity
		}

		if component.ScanRangePen != NoScanner {
			hasPenScan = true
			spec.ScanRangePen += PowInt(component.ScanRangePen, 4) * slot.Quantity
		}
	}

	// time to quad root everything
	if spec.ScanRange > 0 {
		s := math.Pow(float64(spec.ScanRange), .25) * scannerSpec.ScanRangeFactor
		spec.ScanRange = int(math.Round(s))
	}

	if spec.ScanRangePen > 0 {
		s := math.Pow(float64(spec.ScanRangePen), .25)
		spec.ScanRangePen = int(math.Round(s))
	} else if !hasPenScan {
		spec.ScanRangePen = NoScanner
	}

	// Update scanner field if we have any scanning capabilities whatsoever
	// all fleets should be able to regular scan at range 0 (i.e. see planet occupation status while in orbit), but not pen scan
	spec.Scanner = spec.ScanRange != NoScanner || spec.ScanRangePen != NoScanner
}

// Design a ship/starbase for the AI or as a starting fleet using the best parts available to us
//
// Warship design is handled by (and delegated to) [cs.DesignWarship] instead
func DesignShip(rules *Rules, hull *TechHull, name string, player *Player, num int, hullSetNumber int, purpose ShipDesignPurpose, fleetPurpose FleetPurpose) (*ShipDesign, error) {

	techStore := rules.techs
	design := NewShipDesign(player, num).WithName(name).WithHull(hull.Name).WithHullSetNumber(hullSetNumber).WithPurpose(purpose)
	tc := NewTechComparer(rules, player)

	// fuel depots & starter colonies are empty
	if purpose == ShipDesignPurposeFuelDepot || purpose == ShipDesignPurposeStarterColony {
		return design, nil
	} else if purpose == ShipDesignPurposeBeamFighter ||
		purpose == ShipDesignPurposeTorpedoFighter ||
		purpose == ShipDesignPurposeFighterScout ||
		purpose == ShipDesignPurposeStarbase ||
		purpose == ShipDesignPurposeFort ||
		purpose == ShipDesignPurposeStarbaseHalf ||
		purpose == ShipDesignPurposeStarbaseQuarter {
		// warships & bases get their own separate function for reasons
		design, err := DesignWarship(rules, hull, name, player, num, hullSetNumber, purpose)
		if err != nil {
			return &ShipDesign{}, fmt.Errorf("DesignWarship returned error %w", err)
		} else {
			return design, nil
		}
	}

	bestPartsBySlot := map[HullSlotType]map[TechTag]*TechHullComponent{} // represents if we've already checked this hull slot type
	hullSlotsByFlexibility := map[int][]int{}                            // lists all the hull slots in our ship sorted by flexibility
	engine := techStore.GetBestEngine(player, hull, fleetPurpose)

	numFuelTanks := 0
	numCargoPods := 0
	numBeamWeapons := 0
	numTorpedoes := 0
	var hasGate, hasDriver, hasScanner, hasColonyModule bool

	maxNum := math.MinInt
	for i, hullSlot := range hull.Slots {
		hst := hullSlot.Type
		// first, we loop around once to make our maps

		b := Bitmask(hst).countBits()
		if b > maxNum {
			maxNum = b
		}
		hullSlotsByFlexibility[b] = append(hullSlotsByFlexibility[b], i) // add list index of the hull slot to our slice
		bestPartsBySlot[hst] = map[TechTag]*TechHullComponent{}
	}

	// loop through a second time to check our hull slots
	var hullSlot TechHullSlot
	getPartAndCache := func(t TechTag) *TechHullComponent {
		return UpdateLookupMap(bestPartsBySlot[hullSlot.Type], t, func(tt TechTag) *TechHullComponent {
			return tc.GetBestComponentWithTag(design, hullSlot.Type, hullSlot.Capacity, tt)
		})
	}

	for i := range maxNum {
		list := hullSlotsByFlexibility[i+1]
		if list == nil {
			// no slots with this many different part types; skip
			continue
		}
		for _, sn := range list {
			hullSlot = hull.Slots[sn]
			hst := hullSlot.Type
			slot := ShipDesignSlot{HullSlotIndex: sn + 1} // list index 0 gets slot no. 1
			slot.Quantity = hullSlot.Capacity

			if hst&HullSlotTypeEngine != 0 {
				slot.HullComponent = engine.Name
				design.Slots = append(design.Slots, slot)
				continue
			}

			// assign slots based on purpose
		purposeSwitch:
			switch purpose {
			case ShipDesignPurposeScout:
				scanner := getPartAndCache(TechTagScanner)
				if !hasScanner && scanner != nil {
					slot.HullComponent = scanner.Name
					hasScanner = true
				}
			case ShipDesignPurposeStartingFighter: // everyone's favorite rinky dinky starter ships
				scanner := getPartAndCache(TechTagScanner)
				if scanner != nil && !hasScanner {
					slot.HullComponent = scanner.Name
					hasScanner = true
					break purposeSwitch
				}
				beamWeapon := getPartAndCache(TechTagBeamWeapon)
				torpedo := getPartAndCache(TechTagTorpedo)
				if torpedo != nil && beamWeapon != nil {
					if numTorpedoes > numBeamWeapons {
						slot.HullComponent = beamWeapon.Name
						numBeamWeapons += slot.Quantity
					} else {
						slot.HullComponent = torpedo.Name
						numTorpedoes += slot.Quantity
					}
					break purposeSwitch
				}
				battleComputer := getPartAndCache(TechTagTorpedoBonus)
				if battleComputer != nil {
					slot.HullComponent = battleComputer.Name
					break purposeSwitch
				}
				armor := getPartAndCache(TechTagArmor)
				if armor != nil {
					slot.HullComponent = armor.Name
				}

			// fill the bomb slot based on the type of bomber we want
			// or leave it blank
			case ShipDesignPurposeSmartBomber:
				smartBomb := getPartAndCache(TechTagSmartBomb)
				if smartBomb != nil {
					slot.HullComponent = smartBomb.Name
				}
			case ShipDesignPurposeStructureBomber:
				structureBomb := getPartAndCache(TechTagStructureBomb)
				if structureBomb != nil {
					slot.HullComponent = structureBomb.Name
				}
			case ShipDesignPurposeBomber:
				bomb := getPartAndCache(TechTagBomb)
				if bomb != nil {
					slot.HullComponent = bomb.Name
				}
			case ShipDesignPurposeFuelFreighter:
			// nothing happens; our default case is to tack on fuel pods in spare slots
			case ShipDesignPurposeColonizer:
				colonyModule := getPartAndCache(TechTagColonyModule)
				if colonyModule != nil && !hasColonyModule {
					slot.HullComponent = colonyModule.Name
					slot.Quantity = 1 // we only need 1 colonization module
					hasColonyModule = true
					break purposeSwitch
				}
				fallthrough
			case ShipDesignPurposeArmedFreighter:
				// TODO: Add purpose for cloaked pokey ships and add cloaks accordingly
				fallthrough
			case ShipDesignPurposeFreighter, ShipDesignPurposeColonistFreighter:
				cargoPod := getPartAndCache(TechTagCargoPod)
				// add cargo pods or fuel pods
				if cargoPod != nil && numCargoPods < numFuelTanks {
					slot.HullComponent = cargoPod.Name
					numCargoPods += slot.Quantity
				}
			case ShipDesignPurposeTerraformer:
				terraformRobot := getPartAndCache(TechTagTerraformingRobot)
				if terraformRobot != nil {
					slot.HullComponent = terraformRobot.Name
				}
			case ShipDesignPurposeMiner:
				miningRobot := getPartAndCache(TechTagMiningRobot)
				if miningRobot != nil {
					slot.HullComponent = miningRobot.Name
				}
			case ShipDesignPurposeSpeedMineLayer:
				speedMineLayer := getPartAndCache(TechTagSpeedMineLayer)
				if speedMineLayer != nil {
					slot.HullComponent = speedMineLayer.Name
					break
				}
				fallthrough
			case ShipDesignPurposeDamageMineLayer:
				// TODO: Ensure AI uses SD detonating minefield layers if applicable
				heavyMineLayer := getPartAndCache(TechTagHeavyMineLayer)
				standardMineLayer := getPartAndCache(TechTagMineLayer)
				if heavyMineLayer != nil {
					slot.HullComponent = heavyMineLayer.Name
				} else if standardMineLayer != nil {
					slot.HullComponent = standardMineLayer.Name
				}
			case ShipDesignPurposePacketThrower:
				massDriver := getPartAndCache(TechTagMassDriver)
				if massDriver != nil && !hasDriver {
					slot.HullComponent = massDriver.Name
					hasDriver = true
				}
			case ShipDesignPurposeStargater:
				stargate := getPartAndCache(TechTagStargate)
				if stargate != nil && !hasGate {
					slot.HullComponent = stargate.Name
					hasGate = true
				}
			case ShipDesignPurposeStarbaseUnarmed:
				stargate := getPartAndCache(TechTagStargate)
				if stargate != nil && !hasGate {
					slot.HullComponent = stargate.Name
					hasGate = true
					break purposeSwitch
				}
				massDriver := getPartAndCache(TechTagMassDriver)
				// if have space, add packet throwers as well
				if massDriver != nil {
					slot.HullComponent = massDriver.Name
					hasDriver = true
				}
			}

			if slot.HullComponent == "" && !hull.Starbase {
				fuelTank := getPartAndCache(TechTagFuelTank)
				shield := getPartAndCache(TechTagShield)
				scanner := getPartAndCache(TechTagScanner)
				switch {
				case fuelTank != nil: // when in doubt, add fuel tanks to empty slots
					slot.HullComponent = fuelTank.Name
					numFuelTanks += slot.Quantity
				case shield != nil: // add shields to freighters so they don't die as much against minefields
					slot.HullComponent = shield.Name
				case scanner != nil && !hasScanner: // also add scanners to super fuels in a pinc
					slot.HullComponent = scanner.Name
					hasScanner = true
				}
			}

			// if we filled the slot, add it to the design's slots
			if slot.HullComponent != "" {
				design.Slots = append(design.Slots, slot)
			}
		}
	}

	var err error
	design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
	if err != nil {
		return &ShipDesign{}, fmt.Errorf("computeShipDesignSpec errored during DesignShip, error: %w", err)
	}
	return design, nil
}

// Design a warship or starbase based on available parts to fit a specified goal
func DesignWarship(rules *Rules, hull *TechHull, name string, player *Player, num int, hullSetNumber int, purpose ShipDesignPurpose) (*ShipDesign, error) {

	//* DISCLAIMER FOR CODE (RE)VIEWERS: THIS IS A *VERY LONG FUNCTION*. Use the hashtags (#) to jump between sections.
	techStore := rules.techs
	design := NewShipDesign(player, num).WithName(name).WithHull(hull.Name).WithHullSetNumber(hullSetNumber).WithPurpose(purpose)
	tc := NewTechComparer(rules, player)

	// (#) COUNTERS & CONSTANTS

	// A list of hull slot indices sorted by slot type/capacity
	//
	// We use this to determine what order to loop through things
	var hullSlotNumsSorted = []int{}
	var capacitorSlots = []int{} // contains all the hull slots we are reserving for beam caps (checked last due to hardcap)
	var jammerSlots = []int{}    // contains all the hull slots we are reserving for jammers (checked last due to hardcap)
	var jetSlots = []int{}       // contains all the hull slots we are reserving for jets (checked last due to hardcap)
	var err error
	numWeapons := 0
	numSappers := 0
	var hasDriver, hasScanner, hasStargate bool
	minWeapons := 6 // min amount of weapons for us to dedicate our GP slots to
	if hull.Type == TechHullTypeFighter {
		minWeapons = int(float64(minWeapons) * 1.5)
		// smaller ships tend to have higher baseline costs for hull/shields/components,
		// so increasing the min weapon count helps add much needed firepower
	}
	bestPartsBySlot := map[HullSlotType]map[TechTag]*TechHullComponent{} // used for quick lookups for repeat parts
	var engineSlots = []int{}                                            // contains all our engine slots
	var numEngines int

	// (#) SLOT INDICING/SORTING

	// add the slots to our slice & initialize our lookup map if needed
	for i, hullSlot := range hull.Slots {
		if hullSlot.Type&HullSlotTypeEngine != 0 && !hull.Starbase { // don't add engines to list for starbases
			numEngines += hullSlot.Capacity
			engineSlots = append(engineSlots, i)
			continue
		}
		hullSlotNumsSorted = append(hullSlotNumsSorted, i)
		if bestPartsBySlot[hullSlot.Type] == nil {
			bestPartsBySlot[hullSlot.Type] = map[TechTag]*TechHullComponent{}
		}
	}

	// get our engine slots out of the way
	if len(engineSlots) == 0 && !hull.Starbase {
		return nil, fmt.Errorf("DesignWarship() could not find any engine slots in hull %s", hull)
	} else {
		bestEngine := techStore.GetBestBattleEngine(player, hull, numEngines)
		for _, i := range engineSlots {
			h := hull.Slots[i]
			design.Slots = append(design.Slots, ShipDesignSlot{
				HullComponent: bestEngine.Name, HullSlotIndex: i + 1, Quantity: h.Capacity})
		}

	}

	// sort through hull slots in order of increasing slot type
	// then in decreasing slot quantity (so bigger slots get used up first)
	// ensures weapons get put on larger slots first all else being equal
	hullSlotTypePriority := map[HullSlotType]int{
		HullSlotTypeNone:                             0,
		HullSlotTypeEngine:                           1,
		HullSlotTypeSpaceDock:                        1 << 1,
		HullSlotTypeCargo:                            1 << 2,
		HullSlotTypeBomb:                             1 << 3,
		HullSlotTypeMining:                           1 << 4,
		HullSlotTypeMineLayer:                        1 << 5,
		HullSlotTypeScanner:                          1 << 6,
		HullSlotTypeOrbital:                          1 << 7,
		HullSlotTypeWeapon:                           1 << 8,
		HullSlotTypeShield:                           1 << 9,
		HullSlotTypeArmor:                            1 << 10,
		HullSlotTypeMechanical:                       1 << 11,
		HullSlotTypeElectrical:                       1 << 12,
		HullSlotTypeElectricalMechanical:             1<<11 | 1<<12,
		HullSlotTypeOrbitalElectrical:                1<<7 | 1<<12,
		HullSlotTypeShieldElectricalMechanical:       1<<9 | 1<<11 | 1<<12,
		HullSlotTypeScannerElectricalMechanical:      1<<6 | 1<<11 | 1<<12,
		HullSlotTypeArmorScannerElectricalMechanical: 1<<10 | 1<<6 | 1<<11 | 1<<12,
		HullSlotTypeMineElectricalMechanical:         1<<5 | 1<<11 | 1<<12,
		HullSlotTypeShieldArmor:                      1<<9 | 1<<10,
		HullSlotTypeWeaponShield:                     1<<8 | 1<<9,
		HullSlotTypeGeneral:                          1<<5 | 1<<6 | 1<<8 | 1<<9 | 1<<10 | 1<<11 | 1<<12,
	}
	if len(hullSlotNumsSorted) > 1 {
		slices.SortStableFunc(hullSlotNumsSorted, func(m, n int) int {
			b := hullSlotTypePriority[hull.Slots[m].Type] - hullSlotTypePriority[hull.Slots[n].Type]
			if b != 0 {
				return b
			}
			// reversing m & n puts list in descending quantity order (biggest first)
			return hull.Slots[n].Capacity - hull.Slots[m].Capacity
		})
	}

	var hullSlot TechHullSlot

	// wrapper function to make code prettier
	getPartAndCache := func(t TechTag) *TechHullComponent {
		return UpdateLookupMap(bestPartsBySlot[hullSlot.Type], t, func(tt TechTag) *TechHullComponent {
			return tc.GetBestComponentWithTag(design, hullSlot.Type, hullSlot.Capacity, tt)
		})
	}

	// extract slot numbers from list so we can loop through them
	for _, slotNum := range hullSlotNumsSorted {
		hullSlot = hull.Slots[slotNum]
		designSlot := ShipDesignSlot{HullSlotIndex: slotNum + 1} // list index 0 gets slot no. 1
		designSlot.Quantity = hullSlot.Capacity
		var itemToPlace *TechHullComponent

		// assign hull components, using map lookups to avoid repetition
		var weapon, driver, stargate, scanner *TechHullComponent
		if design.Purpose.IsBeamShip() {
			weapon = getPartAndCache(TechTagBeamWeapon)
		} else {
			weapon = getPartAndCache(TechTagTorpedo)
		}
		driver = getPartAndCache(TechTagMassDriver)
		stargate = getPartAndCache(TechTagStargate)
		scanner = getPartAndCache(TechTagScanner)

		switch {
		case weapon != nil && ((numWeapons+numSappers) < minWeapons) || hullSlot.Type == HullSlotTypeWeapon:
			// if we don't have many weapons already or this is a
			// weapons-only slot, slap on some guns

			// decide on whether to use sappers or not
			// TODO: Rework this once armorDamageMulti becomes a techHullComponent property
			sapper := getPartAndCache(TechTagShieldSapper)

			shouldUseSapper := sapper != nil && // have a sapper to use
				sapper.Range == weapon.Range && // sapper has at least as much range as our main guns
				numWeapons > numSappers*3 && // 3:1 gun:sapper ratio
				tc.compareWeaponPowers(weapon, sapper) // sapper does more damage per hit

			if shouldUseSapper {
				itemToPlace = sapper
				numSappers += designSlot.Quantity
			} else if weapon != nil {
				itemToPlace = weapon
				numWeapons += designSlot.Quantity
			}
		case scanner != nil && design.Purpose == ShipDesignPurposeFighterScout &&
			!hasScanner:
			// add scanners to armed scouts if they don't have them already
			// TODO: Add a way to determine the "least needed" slot rather than tacking a scanner on the first one we find
			itemToPlace = scanner
			if itemToPlace.Tags.Count() == 1 && itemToPlace.Tags.HasTag(TechTagScanner) { // covers for non-useless scanner items
				designSlot.Quantity = 1
			}
			// Note that due to the hull slot sorting done earlier,
			// any "Scanner/XXX" slots will only be checked _after_ every single "Scanner only" slot
		case (driver != nil || (stargate != nil && !hasStargate)) && hull.Starbase &&
			purpose != ShipDesignPurposeFort:
			// add orbital items to starbases if avaliable
			if driver != nil && !hasDriver {
				itemToPlace = driver
				hasDriver = true
			} else if stargate != nil && !hasStargate {
				itemToPlace = stargate
				hasStargate = true
			} else if driver != nil {
				itemToPlace = driver
				hasDriver = true
			}
		case hullSlot.Type == HullSlotTypeShield: // covers for langston shell
			shield := getPartAndCache(TechTagShield)

			itemToPlace = shield
		default:
			// add whatever we need the most
			itemToPlace, err = tc.getMostNeededComponent(design, hullSlot.Type, designSlot.Quantity)
			if err != nil {
				return nil, fmt.Errorf("getMostNeededComponent failed to get parts, error %w", err)
			}
		}

		// (#) SPEC RECOMPUTATION

		// however we happened to fill the slot, tack it on and recompute spec fields
		if itemToPlace != nil {
			designSlot.HullComponent = itemToPlace.Name
			// reduce qty for partially built starbases
			if itemToPlace.HullSlotType&(HullSlotTypeShieldArmor|HullSlotTypeWeapon) != 0 {
				if design.Purpose == ShipDesignPurposeStarbaseHalf {
					designSlot.Quantity /= 2
				} else if design.Purpose == ShipDesignPurposeStarbaseQuarter {
					designSlot.Quantity /= 4
				}
			}

			// if all this item does is add beam bonus and/or jamming, tack it on a separate "reserved" list
			// Parts stay in design slots so as to not interfere with part placement logic,
			// and will be removed and re-added later to prevent overcapping
			isPureJammer := itemToPlace.Tags.hasTags([]TechTag{TechTagTorpedoJammer}, CombatTechTags...)
			isPureCapacitor := itemToPlace.Tags.hasTags([]TechTag{TechTagBeamCapacitor}, CombatTechTags...)
			isPureJet := itemToPlace.Tags.hasTags([]TechTag{TechTagManeuveringJet}, CombatTechTags...)
			if isPureJammer {
				jammerSlots = append(jammerSlots, slotNum)
			} else if isPureCapacitor {
				capacitorSlots = append(capacitorSlots, slotNum)
			} else if isPureJet {
				jetSlots = append(jetSlots, slotNum)
			}
			design.Slots = append(design.Slots, designSlot)
			hasScanner = hasScanner || itemToPlace.Scanner

			design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
			if err != nil {
				return nil, fmt.Errorf("computeShipDesignSpec errored during warship part allocation, error: %w", err)
			}

		}
	}

	// (#) JAMMERS, CAPACITORS & JETS
	// add on our long lost capacitor & jammer friends
	if len(capacitorSlots) > 0 && design.Spec.BeamBonus > rules.BeamBonusCap {
		// remove "pure" capacitor items from the design to get an accurate read of our stats
		design.Slots = slices.DeleteFunc(design.Slots, func(sd ShipDesignSlot) bool {
			item := rules.techs.GetHullComponent(sd.HullComponent)
			return item != nil && item.Tags.hasTags([]TechTag{TechTagBeamCapacitor}, CombatTechTags...)
		})
		design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
		if err != nil {
			return nil, fmt.Errorf("computeShipDesignSpec errored during warship part allocation, error: %w", err)
		}
		prevCapacitating := design.Spec.BeamBonus
	capLoop:
		for _, id := range capacitorSlots {
			// place our best capacitor into the slot
			hullSlot := hull.Slots[id]
			capacitor := tc.GetBestComponentWithTag(design, hullSlot.Type, hullSlot.Capacity, TechTagBeamCapacitor)
			slot := ShipDesignSlot{HullComponent: capacitor.Name, HullSlotIndex: id + 1}
			// add them one by one to make sure we don't go overboard
			for range hullSlot.Capacity {
				slot.Quantity++
				prevCapacitating *= 1 + capacitor.BeamBonus
				if prevCapacitating >= rules.BeamBonusCap {
					// we hit the beam bonus cap; no more capacitors needed
					design.Slots = append(design.Slots, slot)
					break capLoop
				}
			}
			// add the finished item to the hullSlot and remove it from the list
			design.Slots = append(design.Slots, slot)
		}
	}

	if len(jammerSlots) > 0 && design.Spec.TorpedoJamming > rules.JammerCap.Get(design.Spec.Starbase) {
		// remove pure jammers from slots and re-compute design spec to figure out how much stat we have
		design.Slots = slices.DeleteFunc(design.Slots, func(sd ShipDesignSlot) bool {
			item := rules.techs.GetHullComponent(sd.HullComponent)
			return item != nil && item.Tags.hasTags([]TechTag{TechTagTorpedoJammer}, CombatTechTags...)
		})
		design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
		if err != nil {
			return nil, fmt.Errorf("computeShipDesignSpec errored during warship part allocation, error: %w", err)
		}
		prevJamming := design.Spec.TorpedoJamming
	jamLoop:
		for _, id := range jammerSlots {
			// place our best jammer into the slot
			hullSlot := hull.Slots[id]
			jammer := tc.GetBestComponentWithTag(design, hullSlot.Type, hullSlot.Capacity, TechTagTorpedoJammer)
			slot := ShipDesignSlot{HullComponent: jammer.Name, HullSlotIndex: id + 1, Quantity: 0}

			// add them one by one to make sure we don't go overboard
			for range hullSlot.Capacity {
				slot.Quantity++
				prevJamming = getNewJamming(prevJamming, jammer.TorpedoJamming, rules.JammerMulti.Get(hull.Starbase), 1)
				if prevJamming >= rules.JammerCap.Get(hull.Starbase)*rules.JammerMulti.Get(hull.Starbase) {
					// we hit the jamming cap; no more jammers needed
					design.Slots = append(design.Slots, slot)
					break jamLoop
				}
			}

			// add the finished item to the design and zero it out
			design.Slots = append(design.Slots, slot)
		}
	}

	if len(jetSlots) > 0 && design.Spec.Movement >= rules.MovementMax {
		// remove pure jets from slots and re-compute design spec to figure out how sped we are
		design.Slots = slices.DeleteFunc(design.Slots, func(sd ShipDesignSlot) bool {
			item := rules.techs.GetHullComponent(sd.HullComponent)
			return item != nil && item.Tags.hasTags([]TechTag{TechTagManeuveringJet}, CombatTechTags...)
		})
		design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
		if err != nil {
			return nil, fmt.Errorf("computeShipDesignSpec errored during warship part allocation, error: %w", err)
		}
	jetLoop:
		for _, id := range jetSlots {
			// place our best jammer into the slot
			hullSlot := hull.Slots[id]
			jet := tc.GetBestComponentWithTag(design, hullSlot.Type, hullSlot.Capacity, TechTagManeuveringJet)
			slot := ShipDesignSlot{HullComponent: jet.Name, HullSlotIndex: id + 1, Quantity: 0}

			// add them one by one to make sure we don't go overboard
			for range hullSlot.Capacity {
				slot.Quantity++
				prevMovement := getBattleMovement(rules.MovementMin, rules.MovementMax, design.Spec.Engine.IdealSpeed, design.Spec.MovementBonus+jet.MovementBonus*float64(slot.Quantity), design.Spec.Mass+jet.Mass*slot.Quantity, design.Spec.NumEngines)
				if prevMovement >= rules.MovementMax {
					// we are going brrr enough; stop
					design.Slots = append(design.Slots, slot)
					break jetLoop
				}
			}

			// add the finished item to the design and zero it out
			design.Slots = append(design.Slots, slot)
		}
	}

	// (#) FINAL WRAP UP
	// Fix all the various temporary jank we did to the ship
	// and check to make sure it's actually functional
	design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
	if err != nil {
		return nil, fmt.Errorf("computeShipDesignSpec errored during warship part allocation, error: %w", err)
	}

	// re-sort hull slots by ascending slot index and remove unused capacity
	design.Slots = slices.Clip(design.Slots)
	slices.SortFunc(design.Slots, func(m, n ShipDesignSlot) int {
		return m.HullSlotIndex - n.HullSlotIndex
	})

	if len(design.Spec.WeaponSlots) == 0 {
		// our "completed" warship has no actual weapons; we assume the build process failed somehow
		slotList := map[string]int{}
		for _, slot := range design.Slots {
			slotList[slot.HullComponent] += slot.Quantity
		}
		slotsString := MapToStringDelimited(slotList, strings.Compare, ": ", "\n")
		return nil, fmt.Errorf("DesignWarship returned ship with no weapon slots; part tallies: \n%v", slotsString)
	}

	return design, nil
}

// return relative factor by which a jammer/computer/deflector boosts our relative torpedo defense/offense
//
// Formula: 1+oldJamming / 1+newJamming (https://www.desmos.com/calculator/cpcyiloqeg)
//
// fieldToCheck determines which stat is being calculated for (jamming, computing or deflecting);
// panics if incorrect tag is given
func (spec *ShipDesignSpec) getJamOrComputerBonus(rules *Rules, hc *TechHullComponent, qty int, fieldToCheck TechTag) float64 {
	var oldBonus, hcBonus, cap, jamMulti float64
	switch fieldToCheck {
	case TechTagTorpedoJammer:
		jamMulti = rules.JammerMulti.Get(spec.Starbase)
		cap = rules.JammerCap.Get(spec.Starbase) * jamMulti
		oldBonus = spec.TorpedoJamming
		hcBonus = hc.TorpedoJamming
	case TechTagTorpedoBonus:
		jamMulti = 1
		cap = 1
		oldBonus = spec.TorpedoBonus
		hcBonus = hc.TorpedoBonus
	case TechTagBeamDeflector:
		jamMulti = 1
		cap = 1
		// TODO: change this after BeamDefense refactor
		if spec.BeamDefense == 0 {
			oldBonus = 0
		} else {
			oldBonus = 1 - spec.BeamDefense
		}
		hcBonus = hc.BeamDefense
	default:
		panic(fmt.Sprintf("incorrect TechTag %s given to getJamOrComputerBonus", fieldToCheck))
	}

	if oldBonus == cap {
		return 1
	}

	// *I HATE FLOATING POINT ROUNDING ERRORS*
	newBonus := math.Min(getNewJamming(oldBonus, hcBonus, jamMulti, qty), cap)

	return (1 + newBonus) / (1 + oldBonus)
}
