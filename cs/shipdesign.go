package cs

import (
	"fmt"
	"math"
	"strings"

	"slices"
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
	MovementBonus             int                   `json:"movementBonus,omitempty"`
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

type MineLayingRateByMineType struct {
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

// Compute the spec for this ShipDesign. This function is mostly for universe generation and tests
func (sd *ShipDesign) WithSpec(rules *Rules, player *Player) *ShipDesign {
	var err error
	sd.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, sd)
	if err != nil {
		panic(fmt.Sprintf("failed to ComputeShipDesignSpec %v", err))
	}
	return sd
}

// validate that this ship design is available to the player
func (sd *ShipDesign) Validate(rules *Rules, player *Player) error {
	if strings.TrimSpace(sd.Name) == "" {
		return fmt.Errorf("design has no name")
	}
	hull := rules.techs.GetHull(sd.Hull)
	if hull == nil {
		return fmt.Errorf("hull %s not found", sd.Hull)
	}
	if !player.HasTech(&hull.Tech) {
		return fmt.Errorf("hull %s is not available to player", hull.Name)
	}

	for _, slot := range sd.Slots {
		if slot.HullSlotIndex < 1 || slot.HullSlotIndex > len(hull.Slots) {
			return fmt.Errorf("hull component index %d out of range", slot.HullSlotIndex)
		}
		hullSlot := hull.Slots[slot.HullSlotIndex-1]
		if slot.Quantity < 0 || slot.Quantity > hullSlot.Capacity {
			return fmt.Errorf("hull component quantity %d out of range", slot.Quantity)
		}
		if hullSlot.Required && hullSlot.Capacity != slot.Quantity {
			return fmt.Errorf("hull component required but quantity %d != capacity %d", slot.Quantity, hullSlot.Capacity)
		}

		// if we have a hull component, check it
		if slot.HullComponent != "" {
			hc := rules.techs.GetHullComponent(slot.HullComponent)
			if hc == nil {
				return fmt.Errorf("hull component %s not found", slot.HullComponent)
			}

			if hullSlot.Type&hc.HullSlotType == 0 {
				return fmt.Errorf("hull component %s won't work in slot %v", hc.Name, hullSlot.Type)
			}

			if len(hc.Requirements.HullsAllowed) > 0 && slices.IndexFunc(hc.Requirements.HullsAllowed, func(h string) bool { return hull.Name == h }) == -1 {
				return fmt.Errorf("hull component %s is not mountable on the %s hull", hc.Name, sd.Hull)
			}

			if len(hc.Requirements.HullsDenied) > 0 && slices.IndexFunc(hc.Requirements.HullsDenied, func(h string) bool { return hull.Name == h }) != -1 {
				return fmt.Errorf("hull component %s is not mountable on the %s hull", hc.Name, sd.Hull)
			}

			if !player.HasTech(&hc.Tech) {
				return fmt.Errorf("hull component %s is not available to player", hc.Name)
			}
		}

	}

	for i, hullSlot := range hull.Slots {
		if hullSlot.Required {
			found := false
			for _, slot := range sd.Slots {
				if slot.HullSlotIndex-1 == i && slot.Quantity == hullSlot.Capacity {
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

// return true if this ship's job requires it to be light
func (p ShipDesignPurpose) IsLightShip() bool {
	// all peacetime ships should not be using armor
	// too heavy and they're gonna die anyways
	return (!p.IsWarship() && p != ShipDesignPurposeStartingFighter) || p.IsBeamShip()
}

// return true if this ship's job is to be some kind of warship
// (starting fighter not included)
func (p ShipDesignPurpose) IsWarship() bool {
	return p.IsBeamShip() || p.IsTorpedoShip()
}

// return true if this ship's job is to use beams
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
		p == ShipDesignPurposeStarbaseQuarter
}

// get the movement for this ship design, based on cargoMass
func (d *ShipDesign) getMovement(cargoMass int) int {
	return getBattleMovement(d.Spec.Engine.IdealSpeed, d.Spec.MovementBonus, d.Spec.Mass+cargoMass, d.Spec.NumEngines)
}

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
		ScanRange:                0, // by default, all ships non-pen scan ships in their radius
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
		return ShipDesignSpec{}, fmt.Errorf("failed to get design cost %w", err)
	}

	// count the number of each type of battle component we have
	torpedoBonusesByCount := map[float64]int{}
	torpedoJammersByCount := map[float64]int{}
	beamBoostersByCount := map[float64]int{}
	beamDeflectorsByCount := map[float64]int{}

	numTachyonDetectors := 0

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
			if hullSlot.Type == HullSlotTypeEngine {
				engine := rules.techs.GetEngine(slot.HullComponent)
				spec.Engine = engine.Engine
				spec.NumEngines = slot.Quantity
			}

			if component.Category == TechCategoryBeamWeapon && component.Power > 0 && (component.Range+hull.RangeBonus) > 0 {
				// mine sweep is power * (range)^2
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
			spec.Armor += int(a)
			spec.Shields += int(s)
			spec.CargoCapacity += component.CargoBonus * slot.Quantity
			spec.FuelCapacity += component.FuelBonus * slot.Quantity
			spec.FuelGeneration += component.FuelGeneration * slot.Quantity
			spec.Colonizer = spec.Colonizer || component.ColonizationModule || component.OrbitalConstructionModule
			spec.Initiative += component.InitiativeBonus * slot.Quantity
			spec.MovementBonus += component.MovementBonus * slot.Quantity
			spec.ReduceMovement = MaxInt(spec.ReduceMovement, component.ReduceMovement) // these don't stack
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
				spec.BasePacketSpeed = MaxInt(spec.BasePacketSpeed, component.PacketSpeed)
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
			rules.JammerCap[hull.Starbase])*rules.JammerMulti[hull.Starbase], 4)

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
		// Movement = IdealEngineSpeed - 2 - Mass / 70 / NumEngines + NumManeuveringJets + 2*NumOverThrusters
		// we added any MovementBonus components above
		// we round up the slightest bit, and we can't go below 2, or above 10
		spec.Movement = getBattleMovement(spec.Engine.IdealSpeed, spec.MovementBonus, spec.Mass, spec.NumEngines)
		spec.MovementFull = getBattleMovement(spec.Engine.IdealSpeed, spec.MovementBonus, spec.Mass+spec.CargoCapacity, spec.NumEngines)
	} else {
		spec.Movement = 0
		spec.MovementFull = 0
	}

	beamPower = int(float64(beamPower) * (spec.BeamBonus))
	if beamPower > 0 {
		// starbases don't move, but for the beam power calcs
		// assume they have a movement of "2" which is the lowest possible
		movement := Clamp(spec.Movement, 2, 10)

		// a movement of 1 1/2 in the UI (i.e. 6) doesn't impact your beam
		// power rating. Anything less reduces your beam power, anything higher increases it
		beamPower += (beamPower * (movement - 6)) / 10
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

// Compute the scan ranges for this ship design
// Formula: (scanner1^4 + scanner2^4 + ...
// + scannerN^4)^(.25)
func (spec *ShipDesignSpec) computeScanRanges(rules *Rules, scannerSpec ScannerSpec, techLevels TechLevel, design *ShipDesign, hull *TechHull) {
	spec.ScanRange = 0
	spec.ScanRangePen = NoScanner

	// compute scanner as a built in JoaT scanner if it's built in
	builtInScannerMultiplier := scannerSpec.BuiltInScannerMultiplier
	if builtInScannerMultiplier > 0 && hull.BuiltInScanner {
		spec.ScanRange = techLevels.Electronics * builtInScannerMultiplier
		spec.ScanRangePen = int(math.Pow(float64(spec.ScanRange)/2, 4))
		spec.ScanRange = int(math.Pow(float64(spec.ScanRange), 4))
	}

	for _, slot := range design.Slots {
		if slot.Quantity == 0 {
			continue
		}

		component := rules.techs.GetHullComponent(slot.HullComponent)
		if !component.Scanner {
			continue
		}

		// bat scanners have 0 range
		if component.ScanRange != NoScanner {
			spec.ScanRange += int(math.Pow(float64(component.ScanRange), 4) * float64(slot.Quantity))
		}

		if component.ScanRangePen != NoScanner {
			if spec.ScanRangePen == NoScanner {
				spec.ScanRangePen = int((math.Pow(float64(component.ScanRangePen), 4)) * float64(slot.Quantity))
			} else {
				spec.ScanRangePen += int((math.Pow(float64(component.ScanRangePen), 4)) * float64(slot.Quantity))
			}
		}
	}

	// now quad root it
	if spec.ScanRange > 0 {
		spec.ScanRange = int(math.Pow(float64(spec.ScanRange), .25) + .5)
		spec.ScanRange = int(float64(spec.ScanRange) * scannerSpec.ScanRangeFactor)
	}

	if spec.ScanRangePen > 0 {
		spec.ScanRangePen = int(math.Pow(float64(spec.ScanRangePen), .25) + .5)
	}

	// true if we have any scanning capability (all fleets should be able to scan at 0, but not pen scan)
	spec.Scanner = spec.ScanRange != NoScanner || spec.ScanRangePen != NoScanner
}

// design a ship/starbase for the AI or as a starting fleet using the best parts available to us
func DesignShip(rules *Rules, hull *TechHull, name string, player *Player, num int, hullSetNumber int, purpose ShipDesignPurpose, fleetPurpose FleetPurpose) (*ShipDesign, error) {

	techStore := rules.techs
	design := NewShipDesign(player, num).WithName(name).WithHull(hull.Name)
	design.Purpose = purpose
	tc := NewTechComparer(rules, player)

	// fuel depots & starter colonies are empty
	if purpose == ShipDesignPurposeFuelDepot || purpose == ShipDesignPurposeStarterColony {
		return design, nil
	} else if purpose == ShipDesignPurposeBeamFighter || purpose == ShipDesignPurposeTorpedoFighter || purpose == ShipDesignPurposeFighterScout || purpose == ShipDesignPurposeStarbase || purpose == ShipDesignPurposeStarbaseHalf || purpose == ShipDesignPurposeStarbaseQuarter {
		// warships get their own separate function for reasons
		design, err := DesignWarship(rules, hull, name, player, num, hullSetNumber, purpose)
		if err != nil {
			return &ShipDesign{}, fmt.Errorf("DesignWarship returned error %w", err)
		} else {
			return design, nil
		}
	}

	techTagsToCheck := []TechTag{
		TechTagBeamWeapon,
		TechTagBomb,
		TechTagSmartBomb,
		TechTagStructureBomb,
		TechTagTorpedo,
		TechTagTorpedoBonus,
		TechTagArmor,
		TechTagShield,
		TechTagMineLayer,
		TechTagHeavyMineLayer,
		TechTagSpeedMineLayer,
		TechTagStargate,
		TechTagMassDriver,
		TechTagTerraformingRobot,
		TechTagMiningRobot,
		TechTagColonyModule,
		TechTagFuelTank,
		TechTagCargoPod,
		TechTagCloak,
		TechTagScanner,
	}

	bestPartsBySlot := map[HullSlotType]map[TechTag]*TechHullComponent{} // represents if we've already checked this hull slot type
	hullSlotsByFlexibility := map[int][]int{}                            // lists all the hull slots in our ship sorted by flexibility
	engine := techStore.GetBestEngine(player, hull, fleetPurpose)

	numColonizationModules := 0
	numFuelTanks := 0
	numCargoPods := 0
	numScanners := 0
	numPacketThrowers := 0
	numStargates := 0
	numBeamWeapons := 0
	numTorpedos := 0

	maxNum := math.MinInt
	for i, hullSlot := range hull.Slots {
		hst := hullSlot.Type
		// first, we loop around once to make our maps

		b := Bitmask(hst).countBits()
		if b > maxNum {
			maxNum = b
		}
		hullSlotsByFlexibility[b] = append(hullSlotsByFlexibility[b], i) // add list index of the hull slot to our slice

		if bestPartsBySlot[hst] == nil && hst != HullSlotTypeEngine { // prevents double counting
			bestPartsBySlot[hst] = map[TechTag]*TechHullComponent{}
			for _, tag := range techTagsToCheck {
				bestPartsBySlot[hst][tag] = tc.GetBestComponentWithTag(design, hst, tag)
			}
		}
	}

	// loop through a second time to check our hull slots
	for i := 1; i <= maxNum; i++ {
		list := hullSlotsByFlexibility[i]
		if list == nil {
			// no slots with this many different part types; skip
			continue
		}
		for _, sn := range list {
			hullSlot := hull.Slots[sn]
			hst := hullSlot.Type
			slot := ShipDesignSlot{HullSlotIndex: sn + 1} // list index 0 gets slot no. 1
			slot.Quantity = hullSlot.Capacity

			scanner := bestPartsBySlot[hst][TechTagScanner]
			beamWeapon := bestPartsBySlot[hst][TechTagBeamWeapon]
			torpedo := bestPartsBySlot[hst][TechTagTorpedo]
			bomb := bestPartsBySlot[hst][TechTagBomb]
			smartBomb := bestPartsBySlot[hst][TechTagSmartBomb]
			structureBomb := bestPartsBySlot[hst][TechTagStructureBomb]
			armor := bestPartsBySlot[hst][TechTagArmor]
			shield := bestPartsBySlot[hst][TechTagShield]
			cargoPod := bestPartsBySlot[hst][TechTagCargoPod]
			fuelTank := bestPartsBySlot[hst][TechTagFuelTank]
			colonizationModule := bestPartsBySlot[hst][TechTagColonyModule]
			battleComputer := bestPartsBySlot[hst][TechTagTorpedoBonus]
			miningRobot := bestPartsBySlot[hst][TechTagMiningRobot]
			terraformRobot := bestPartsBySlot[hst][TechTagTerraformingRobot]
			standardMineLayer := bestPartsBySlot[hst][TechTagMineLayer]
			heavyMineLayer := bestPartsBySlot[hst][TechTagHeavyMineLayer]
			speedMineLayer := bestPartsBySlot[hst][TechTagSpeedMineLayer]
			packetThrower := bestPartsBySlot[hst][TechTagMassDriver]
			stargate := bestPartsBySlot[hst][TechTagStargate]

			if hst&HullSlotTypeEngine != 0 {
				slot.HullComponent = engine.Name
				design.Slots = append(design.Slots, slot)
				continue
			}
			switch purpose {
			case ShipDesignPurposeScout:
				if numScanners == 0 && scanner != nil {
					slot.HullComponent = scanner.Name
					numScanners += slot.Quantity
				}
			case ShipDesignPurposeStartingFighter: // everyone's favorite rinky dinky starter ships
				if numScanners == 0 && scanner != nil {
					slot.HullComponent = scanner.Name
					numScanners += slot.Quantity
				} else if torpedo != nil && beamWeapon != nil {
					if numTorpedos > numBeamWeapons {
						slot.HullComponent = beamWeapon.Name
						numBeamWeapons += slot.Quantity
					} else {
						slot.HullComponent = torpedo.Name
						numTorpedos += slot.Quantity
					}
				} else if battleComputer != nil {
					slot.HullComponent = battleComputer.Name
				} else if armor != nil {
					slot.HullComponent = armor.Name
				}

			// fill the bomb slot based on the type of bomber we want
			// or leave it blank
			case ShipDesignPurposeSmartBomber:
				if smartBomb != nil {
					slot.HullComponent = smartBomb.Name
				}
			case ShipDesignPurposeStructureBomber:
				if structureBomb != nil {
					slot.HullComponent = structureBomb.Name
				}
			case ShipDesignPurposeBomber:
				if bomb != nil {
					slot.HullComponent = bomb.Name
				}
			case ShipDesignPurposeFuelFreighter:
			// nothing happens; our default case is to tack on fuel pods in spare slots
			case ShipDesignPurposeColonizer:
				if colonizationModule != nil && numColonizationModules == 0 {
					slot.HullComponent = colonizationModule.Name
					slot.Quantity = 1 // we only need 1 colonization module
					numColonizationModules++
					break
				}
				fallthrough
			case ShipDesignPurposeArmedFreighter:
				// TODO: Add purpose for cloaked ships and add cloaks accordingly
				fallthrough
			case ShipDesignPurposeFreighter, ShipDesignPurposeColonistFreighter:
				// add cargo pods or fuel pods
				if cargoPod != nil && numCargoPods < numFuelTanks {
					slot.HullComponent = cargoPod.Name
					numCargoPods += slot.Quantity
				}
			case ShipDesignPurposeTerraformer:
				if terraformRobot != nil {
					slot.HullComponent = terraformRobot.Name
				}
			case ShipDesignPurposeMiner:
				if miningRobot != nil {
					slot.HullComponent = miningRobot.Name
				}
			case ShipDesignPurposeSpeedMineLayer:
				if speedMineLayer != nil {
					slot.HullComponent = speedMineLayer.Name
					break
				}
				fallthrough
			case ShipDesignPurposeDamageMineLayer:
				// TODO: Ensure AI uses SD detonating minefield layers if applicable
				if heavyMineLayer != nil {
					slot.HullComponent = heavyMineLayer.Name
				} else if standardMineLayer != nil {
					slot.HullComponent = standardMineLayer.Name
				}
			case ShipDesignPurposePacketThrower:
				if packetThrower != nil {
					slot.HullComponent = packetThrower.Name
					numPacketThrowers++
					break
				}
				fallthrough
			case ShipDesignPurposeStargater:
				if stargate != nil {
					slot.HullComponent = stargate.Name
					numStargates++
					break
				}
				// if in doubt, add an extra packet thrower for lulz
				if packetThrower != nil {
					slot.HullComponent = packetThrower.Name
					numPacketThrowers++
				}
			}

			if slot.HullComponent == "" {
				if fuelTank != nil { // when in doubt, add fuel tanks to empty slots
					slot.HullComponent = fuelTank.Name
					numFuelTanks += slot.Quantity
				} else if shield != nil { // also add shields to freighters so they don't die as much against minefields
					slot.HullComponent = shield.Name
				} else if scanner != nil {
					slot.HullComponent = scanner.Name
				}
			}

			// if we filled the slot, add it to the design's slots
			if slot.HullComponent != "" {
				design.Slots = append(design.Slots, slot)
			}
		}
	}
	return design, nil
}

// Design a warship or starbase based on available parts to fit a specified goal
func DesignWarship(rules *Rules, hull *TechHull, name string, player *Player, num int, hullSetNumber int, purpose ShipDesignPurpose) (*ShipDesign, error) {

	//* DISCLAIMER FOR CODE (RE)VIEWERS: THIS IS A *VERY LONG FUNCTION*. Use the hashtags (#) to jump between sections.
	techStore := rules.techs
	design := NewShipDesign(player, num).WithName(name).WithHull(hull.Name).WithPurpose(purpose)
	tc := NewTechComparer(rules, player)

	// (#) COUNTERS & CONSTANTS

	// A list of hull slots by flexibility (how many different item types you can put in it).
	//
	// We use this to determine what order to loop through things
	var hullSlotsByFlexibility = map[int][]int{}
	capacitorSlots := []int{} // contains all the hull slots we are reserving for beam caps (checked last due to hardcap)
	jammerSlots := []int{}    // contains all the hull slots we are reserving for jammers (checked last due to hardcap)
	engine := techStore.GetBestBattleEngine(player, hull)
	numWeapons := 0
	numSappers := 0
	numPacketThrowers := 0
	numStargates := 0
	var err error
	var hasScanner bool // whether we have a scanner or not
	minMove := 4        // min movement to trigger placing emergency jets
	minWeapons := 6     // min amount of weapons for us to dedicate our GP slots to
	if hull.Type == TechHullTypeFighter {
		minWeapons = int(float64(minWeapons) * 1.5)
		// smaller ships tend to have higher baseline costs for hull/shields/components,
		// so increasing the min weapon count helps add much needed firepower
	}

	// (#) LOGIC AND PART RETRIEVAL
	maxNum := math.MinInt
	for i, hullSlot := range hull.Slots {
		b := Bitmask(hullSlot.Type).countBits()
		if b > maxNum {
			maxNum = b
		}
		hullSlotsByFlexibility[b] = append(hullSlotsByFlexibility[b], i) // add the slot index to our slice
	}

	// loop through hull slots in order of increasing flexibility (single item type first, then 2, then 3...)
	// then in decreasing slot quantity (so bigger slots get used up first)
	// ensures weapons get put on larger slots
	for i := 1; i <= maxNum; i++ {
		list := hullSlotsByFlexibility[i]
		if list == nil {
			// if we have no slots with this many different HullSlotTypes, skip
			continue
		}
		if len(list) > 1 {
			slices.SortStableFunc(list, func(m, n int) int {
				return hull.Slots[n].Capacity - hull.Slots[m].Capacity // reversing m & n puts list in descending order
			})
		}

		var jetCounter int = 1 // counter used to decide whether to add jets or not; ensures we don't get 12 jets and nothing else

		// extract slot numbers from list so we can loop through them
		for _, slotNum := range list {
			hullSlot := hull.Slots[slotNum]
			hst := hullSlot.Type
			designSlot := ShipDesignSlot{HullSlotIndex: slotNum + 1} // list index 0 gets slot no. 1
			designSlot.Quantity = hullSlot.Capacity
			itemToPlace := &TechHullComponent{}

			var weapon, sapper, driver, stargate, scanner, jet *TechHullComponent
			if design.Purpose.IsTorpedoShip() {
				// add a big long ~~cucumber~~ missile
				weapon = tc.GetBestComponentWithTag(design, hst, TechTagTorpedo)
			} else {
				// add beems to our beem sheep
				weapon = tc.GetBestComponentWithTag(design, hst, TechTagBeamWeapon)
			}
			driver = tc.GetBestComponentWithTag(design, hst, TechTagMassDriver)
			stargate = tc.GetBestComponentWithTag(design, hst, TechTagStargate)
			scanner = tc.GetBestComponentWithTag(design, hst, TechTagScanner)
			jet = tc.GetBestComponentWithTag(design, hst, TechTagManeuveringJet)

			switch {
			case engine != nil && hst&HullSlotTypeEngine != 0:
				itemToPlace = &engine.TechHullComponent
				goto partPlacement
			case weapon != nil && ((numWeapons+numSappers) < minWeapons) || hst == HullSlotTypeWeapon:
				// first, if we don't have many weapons already, dedicate some of our flex slots to them
				sapper = tc.GetBestComponentWithTag(design, hst, TechTagShieldSapper)

				var shouldUseSapper bool
				// decide on whether to use sappers or not
				// TODO: Rework this once armor damage multi becomes a thing
				shouldUseSapper = sapper != nil && // have a sapper to use
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
			case (driver != nil || (stargate != nil && numStargates == 0)) && hull.Starbase:
				// add orbital items to starbases
				if driver != nil && numPacketThrowers == 0 {
					itemToPlace = driver
					numPacketThrowers++
				} else if stargate != nil && numStargates == 0 {
					itemToPlace = stargate
					numStargates++
				} else if driver != nil {
					itemToPlace = driver
					numPacketThrowers++
				}

			case scanner != nil && design.Purpose == ShipDesignPurposeFighterScout &&
				!hasScanner:
				// add scanners to armed scouts if they don't have them already
				itemToPlace = scanner
				designSlot.Quantity = 1 // only need 1 scanner, otherwise unnecessary

			case jet != nil && !hull.Starbase && ((jetCounter == 1 && design.Spec.Movement < 10) ||
				design.Spec.Movement < minMove):
				// add on jets if the counter is positioned right OR if we
				itemToPlace = jet
			default:
				// add whatever we need the most
				// TODO: Somehow work in estimated torpedo hitrate into computer cost analysis
				// to incentivize the AI to use them more over shields or jammers
				mostNeededItem, err := tc.getMostNeededComponent(design, hst, designSlot.Quantity)
				if err != nil {
					return &ShipDesign{}, fmt.Errorf("getMostNeededComponent failed to get parts for tag %v, error: %w", TechTagManeuveringJet, err)
				}

				itemToPlace = mostNeededItem
			}

			// reset the counter so we don't get 20 million jets
			jetCounter += 1
			if jetCounter > 3 {
				jetCounter = 1
			}

			// (#) SPEC RECOMPUTATION

			// however we happened to fill the slot, tack it on and recompute spec fields
		partPlacement:
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

				// if all this item does is add beam bonus and/or jamming, tack it on a separate "reserved" list and
				// don't count it in our prevXXX totals
				// (since they're used later to track our starting computing/jamming scores before adding stuff)
				// Parts stay in design slots so as to not interfere with part placement logic,
				// and will be overwritten
				isPureCapacitor := (itemToPlace.Tech.Tags.CountTags() == 1 && itemToPlace.Tags[TechTagBeamCapacitor]) ||
					(itemToPlace.Tech.Tags.CountTags() == 2 && (itemToPlace.Tags[TechTagBeamCapacitor] && itemToPlace.Tags[TechTagTorpedoJammer]))
				isPureJammer := (itemToPlace.Tech.Tags.CountTags() == 1 && itemToPlace.Tags[TechTagTorpedoJammer]) ||
					(itemToPlace.Tech.Tags.CountTags() == 2 && (itemToPlace.Tags[TechTagBeamCapacitor] && itemToPlace.Tags[TechTagTorpedoJammer]))
				if isPureJammer {
					jammerSlots = append(jammerSlots, slotNum)
				}
				if isPureCapacitor {
					capacitorSlots = append(capacitorSlots, slotNum)
				} else {
					designSlot.HullComponent = itemToPlace.Name
				}
				design.Slots = append(design.Slots, designSlot)
				hasScanner = hasScanner || (itemToPlace.ScanRange >= 0 || itemToPlace.ScanRangePen >= 0)

				design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
				if err != nil {
					return &ShipDesign{}, fmt.Errorf("computeShipDesignSpec errored during warship part allocation, error: %w", err)
				}

			}
		}
	}

	// (#) JAMMERS, CAPACITORS & WRAP-UP
	// add on our long lost capacitor & jammer friends
	if len(capacitorSlots) > 0 {
		// remove "pure" capacitor items from the design for now
		design.Slots = slices.DeleteFunc(design.Slots, func(sd ShipDesignSlot) bool {
			item := rules.techs.GetHullComponent(sd.HullComponent)
			return item != nil && (item.Tech.Tags.CountTags() == 1 && item.Tags[TechTagBeamCapacitor] ||
				(item.Tech.Tags.CountTags() == 2 && (item.Tags[TechTagBeamCapacitor] && item.Tags[TechTagTorpedoJammer])))
		})
		design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
		if err != nil {
			return &ShipDesign{}, fmt.Errorf("computeShipDesignSpec errored during warship part allocation, error: %w", err)
		}
		prevCapacitating := design.Spec.BeamBonus
	capLoop:
		for i, id := range capacitorSlots {
			// place our best capacitor into the slot
			hullSlot := hull.Slots[id]
			capacitor := tc.GetBestComponentWithTag(design, hullSlot.Type, TechTagBeamCapacitor)
			slot := ShipDesignSlot{HullComponent: capacitor.Name, HullSlotIndex: id + 1}

			// add them one by one to make sure we don't go overboard
			for slot.Quantity = 1; slot.Quantity < hullSlot.Capacity; slot.Quantity++ {
				prevCapacitating *= 1 + capacitor.BeamBonus
				if prevCapacitating > rules.BeamBonusCap {
					// we hit the beam bonus cap; no more capacitors needed
					design.Slots = append(design.Slots, slot)
					capacitorSlots[i] = 0
					break capLoop
				}
			}
			// add the finished item to the hullSlot and remove it from the list
			design.Slots = append(design.Slots, slot)
			capacitorSlots[i] = 0
		}
	}

	if len(jammerSlots) > 0 {
		// remove pure jammers from slots and re-compute design spec to figure out how much stat we have
		design.Slots = slices.DeleteFunc(design.Slots, func(sd ShipDesignSlot) bool {
			item := rules.techs.GetHullComponent(sd.HullComponent)
			return item != nil && (item.Tech.Tags.CountTags() == 1 && item.Tags[TechTagTorpedoJammer] ||
				(item.Tech.Tags.CountTags() == 2 && (item.Tags[TechTagBeamCapacitor] && item.Tags[TechTagTorpedoJammer])))
		})
		design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
		if err != nil {
			return &ShipDesign{}, fmt.Errorf("computeShipDesignSpec errored during warship part allocation, error: %w", err)
		}
		prevJamming := design.Spec.TorpedoJamming
	jamLoop:
		for i, id := range jammerSlots {
			// place our best jammer into the slot
			hullSlot := hull.Slots[id]
			qty := hullSlot.Capacity
			if design.Purpose == ShipDesignPurposeStarbaseHalf {
				qty /= 2
			} else if design.Purpose == ShipDesignPurposeStarbaseQuarter {
				qty /= 4
			}
			jammer := tc.GetBestComponentWithTag(design, hullSlot.Type, TechTagTorpedoJammer)
			slot := ShipDesignSlot{HullComponent: jammer.Name, HullSlotIndex: id + 1}

			// add them one by one to make sure we don't go overboard
			for slot.Quantity = 1; slot.Quantity < hullSlot.Capacity; slot.Quantity++ {
				prevJamming = roundFloat(getNewJamming(prevJamming, jammer.TorpedoJamming, rules.JammerMulti[hull.Starbase], 1), 4)
				if prevJamming >= rules.JammerCap[hull.Starbase]*rules.JammerMulti[hull.Starbase] {
					// we hit the jamming cap; no more jammers needed
					design.Slots = append(design.Slots, slot)
					jammerSlots[i] = 0
					break jamLoop
				}
			}

			// add the finished item to the design and zero it out
			design.Slots = append(design.Slots, slot)
			jammerSlots[i] = 0
		}
	}

	// finally, fix all the various temporary jank we did to the ship
	design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
	if err != nil {
		return &ShipDesign{}, fmt.Errorf("computeShipDesignSpec errored during warship part allocation, error: %w", err)
	}

	// re-sort hull slots by ascending slot index
	slices.SortFunc(slices.Clip(design.Slots), func(m, n ShipDesignSlot) int {
		return m.HullSlotIndex - n.HullSlotIndex
	})

	return design, nil
}
