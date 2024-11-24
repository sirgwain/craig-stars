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
	return !p.IsWarship() || p.IsBeamShip()
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
			a, s := getActualArmorAmount(float64(component.Armor), float64(component.Shield), slot.Quantity, raceSpec, component.Category == TechCategoryArmor)
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

// return the total % amount this TechHullComponent boosts our ship's performance,
// using its TechTags to evaluate individual bonuses
//
// Used to assess parts from different categories to determine
// which one is the best for us to use
func (design *ShipDesign) getWarshipPartBonus(rules *Rules, player *Player, hc *TechHullComponent, qty int) float64 {
	relativeBoost := 1.0
	checkedShield := false

	// check tags individually and tally up the numbers
	for tag := range hc.Tags {
		switch {
		case (tag == TechTagArmor || tag == TechTagShield) && !checkedShield:
			// grab shield and armor stats
			hcArmor, hcShield := getActualArmorAmount(float64(hc.Armor), float64(hc.Shield), qty, player.Race.Spec, hc.Category == TechCategoryArmor)
			oldArmor := float64(design.Spec.Armor)
			oldShield := float64(design.Spec.Shields)
			newArmor := math.Max(hcArmor*float64(qty)+oldArmor, 1) // prevents divide by 0 errors
			newShield := math.Max(hcShield*float64(qty)+oldShield, 1)

			// apply score penalties for having too much armor/shield, up to 3x
			// margin of error is 40% for armor, 30% for shield
			// (if you have too much armor you're more vulnerable to missiles, 
			// but too much shield and you're just dead)
			if newArmor/newShield > 1.4 {
				// reduce our effective armor bonus for adding too much armor
				hcArmor /= math.Min(newArmor/newShield, 3)
				newArmor = hcArmor*float64(qty) + oldArmor
			} else if newShield/newArmor > 1.3 {
				// reduce our effective shield bonus for adding too much shield
				hcShield /= math.Min(newShield/newArmor, 3)
				newShield = hcShield*float64(qty) + oldShield
			}

			// wrap up by adding relative armor & shield boosts to get overall relative bonus
			relativeBoost *= (newArmor + newShield) / (oldArmor + oldShield)
			checkedShield = true // needed to prevent shield/armor parts from double counting themselves
		case tag == TechTagBeamCapacitor:
			if design.Purpose.IsTorpedoShip() {
				break // beam bonus meaningless on missile boats
			}
			// new beam bonus / old beam bonus
			relativeBoost *= math.Min(design.Spec.BeamBonus*math.Pow(1+hc.BeamBonus, float64(qty)),
				rules.BeamBonusCap) / design.Spec.BeamBonus
		case tag == TechTagBeamDeflector:
			relativeBoost *= 1 / math.Pow(1-hc.BeamDefense, float64(qty))
		case tag == TechTagInitiativeBonus:
			// do nothing lol
		case tag == TechTagTorpedoBonus:
			if design.Purpose.IsBeamShip() {
				break // torpedo bonus meaningless on beamships
			}
			fallthrough
		case tag == TechTagTorpedoJammer:
			relativeBoost *= design.Spec.getJamOrComputerIncrease(rules, hc, qty, tag == TechTagTorpedoJammer)

			// all other cases are either stupidly painful to quantify or don't get checked in the actual design function
		}

		if design.Purpose.IsBeamShip() && hc.Mass > 30 {
			// reduce bonus if component is too heavy and we need to go fast
			relativeBoost = 1 + (relativeBoost-1)/(1+float64(hc.Mass-30)/10)
		}
	}
	return roundFloat(relativeBoost, 4)
}

// return relative factor by which a jammer or computer boosts our relative torpedo defense/offense
//
// Formula: 1+oldJamming / 1+newJamming (https://www.desmos.com/calculator/ee0yze4hya)
//
// fieldToCheck determines which stat is being calculated for; true checks jamming whereas false checks computing
func (spec *ShipDesignSpec) getJamOrComputerIncrease(rules *Rules, hc *TechHullComponent, qty int, fieldToCheck bool) float64 {
	var oldBonus, hcBonus, cap, jamMulti float64
	if fieldToCheck {
		jamMulti = rules.JammerMulti[spec.Starbase]
		cap = rules.JammerCap[spec.Starbase] * jamMulti
		oldBonus = spec.TorpedoJamming
		hcBonus = hc.TorpedoJamming
	} else {
		jamMulti = 1
		cap = 1
		oldBonus = spec.TorpedoBonus
		hcBonus = hc.TorpedoBonus
	}

	if oldBonus == cap {
		return 1
	}

	// *I HATE FLOATING POINT ROUNDING ERRORS*
	componentBonus := roundFloat(1-(1-math.Pow(1-hcBonus, float64(qty)))*jamMulti, 7)
	newBonus := roundFloat(math.Min(1-(1-oldBonus)*
		componentBonus, cap), 7)

	// return final score, rounded to 4 decimal places for ease of comparison
	return roundFloat((1+newBonus)/(1+oldBonus), 4)
}

// design a ship/starbase for the AI or as a starting fleet using the best parts available to us
func DesignShip(rules *Rules, hull *TechHull, name string, player *Player, num int, hullSetNumber int, purpose ShipDesignPurpose, fleetPurpose FleetPurpose) (*ShipDesign, error) {

	techStore := rules.techs
	tc := NewTechComparer()
	design := NewShipDesign(player, num).WithName(name).WithHull(hull.Name)
	design.Purpose = purpose

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
				var err error
				bestPartsBySlot[hst][tag], err = tc.GetBestComponentWithTag(rules, player, design, hst, tag)
				if err != nil {
					return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag returned error %w", err)
				}
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

// get the best design for a warship or starbase based on available parts
func DesignWarship(rules *Rules, hull *TechHull, name string, player *Player, num int, hullSetNumber int, purpose ShipDesignPurpose) (*ShipDesign, error) {

	//* DISCLAIMER FOR CODE (RE)VIEWERS: THIS IS A *VERY LONG FUNCTION*. Use the hashtags (#) to jump between sections.
	techStore := rules.techs
	design := NewShipDesign(player, num).WithName(name).WithHull(hull.Name).WithPurpose(purpose)
	tc := NewTechComparer()

	// (#) COUNTERS & CONSTANTS

	// A list of hull slots by flexibility (how many different item types you can put in it).
	//
	// We use this to determine what order to loop through things
	var hullSlotsByFlexibility = map[int][]int{}
	capacitorSlots := []int{} // contains all the hull slots we are reserving for beam caps (checked last due to hardcap)
	jammerSlots := []int{}    // contains all the hull slots we are reserving for jammers (checked last due to hardcap)
	engine := techStore.GetBestBattleEngine(player, hull)
	var prevCapacitating float64 // our design's beamBonus before adding extra beam caps; used in final cap loop
	var prevJamming float64      // our design's jamming before adding extra jammers; used in final jamming loop
	numWeapons := 0
	numSappers := 0
	numPacketThrowers := 0
	numStargates := 0
	var hasScanner bool // whether we have a scanner or not

	// These values are used in the code to influence part choice;
	// modifying these will change behavior accordingly
	targetMove := 9
	targetComputing := 0.8
	targetJamming := 0.8 * rules.JammerMulti[hull.Starbase]
	minWeapons := 6
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


		var elecCounter int = 1 // counter used to decide which parts to add; ensures even spread of parts
		var mechCounter int = 1 // counter used to decide which parts to add

		// extract slot numbers from list so we can loop through them
		for _, slotNum := range list {
			hullSlot := hull.Slots[slotNum]
			hst := hullSlot.Type
			designSlot := ShipDesignSlot{HullSlotIndex: slotNum + 1} // list index 0 gets slot no. 1
			designSlot.Quantity = hullSlot.Capacity

			// (#) SINGLE SLOT ITEM PLACEMENT
			itemToPlace := &TechHullComponent{}
			if i == 1 {
				// only 1 part type means we can check
				// the slot type manually and assign parts on a case-by-case basis that way
				switch hst {
				case HullSlotTypeEngine: // we should ALWAYS HAVE AN ENGINE TO USE ON OUR SHIPS. IF NOT WE'RE SCREWED
					itemToPlace = &engine.TechHullComponent
				case HullSlotTypeShield:
					shield, err := tc.GetBestComponentWithTag(rules, player, design, hst, TechTagShield)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagShield, err)
					}
					itemToPlace = shield
				case HullSlotTypeArmor:
					armor, err := tc.GetBestComponentWithTag(rules, player, design, hst, TechTagArmor)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagArmor, err)
					}
					itemToPlace = armor
				case HullSlotTypeWeapon:
					var weapon, sapper *TechHullComponent
					var err error
					var shouldUseSapper bool 
					if design.Purpose.IsTorpedoShip() {
						// add our best torp/missile
						weapon, err = tc.GetBestComponentWithTag(rules, player, design, hst, TechTagTorpedo)
						if err != nil {
							return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagTorpedo, err)
						}
					} else {
						// add beems to our beem sheep
						weapon, err = tc.GetBestComponentWithTag(rules, player, design, hst, TechTagBeamWeapon)
						if err != nil {
							return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagBeamWeapon, err)
						}
					}

					sapper, err = tc.GetBestComponentWithTag(rules, player, design, hst, TechTagShieldSapper)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagShieldSapper, err)
					}

					// decide on whether to use sappers or not
					// TODO: Rework this once armor damage multi becomes a thing
					shouldUseSapper = sapper != nil && // have a sapper to use
						sapper.Range == weapon.Range && // sapper has at least as much range as our main guns
						numWeapons*3 > numSappers && // 3:1 gun:sapper ratio
						tc.compareWeaponPowers(rules, weapon, sapper) // sapper does more damage per hit
					if shouldUseSapper {
						itemToPlace = sapper
						numSappers += designSlot.Quantity
					} else {
						itemToPlace = weapon
						numWeapons += designSlot.Quantity
					}
				case HullSlotTypeScanner:
					scanner, err := tc.GetBestComponentWithTag(rules, player, design, hst, TechTagScanner)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagScanner, err)
					}
					itemToPlace = scanner
				case HullSlotTypeOrbital:
					driver, err := tc.GetBestComponentWithTag(rules, player, design, hst, TechTagMassDriver)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagMassDriver, err)
					}
					stargate, err := tc.GetBestComponentWithTag(rules, player, design, hst, TechTagStargate)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagStargate, err)
					}

					// place a packet thrower, then a stargate, then another packet thrower
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
				case HullSlotTypeElectrical:
					// elec stuff (jammers/comps/caps)
					computer, err := tc.GetBestComponentWithTag(rules, player, design, hst, TechTagTorpedoBonus)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagTorpedoBonus, err)
					}
					capacitor, err := tc.GetBestComponentWithTag(rules, player, design, hst, TechTagBeamCapacitor)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagBeamCapacitor, err)
					}
					jammer, err := tc.GetBestComponentWithTag(rules, player, design, hst, TechTagTorpedoJammer)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagTorpedoJammer, err)
					}

					switch elecCounter {
					case 1:
						// add beam caps and computers
						if design.Purpose.IsTorpedoShip() && computer != nil &&
							design.Spec.TorpedoBonus < targetComputing {
							itemToPlace = computer
							break
						} else if !design.Purpose.IsTorpedoShip() && capacitor != nil &&
							design.Spec.BeamBonus < rules.BeamBonusCap {
							itemToPlace = capacitor
							break
						}
						fallthrough
					case -1:
						// add jammers until we don't need them
						if jammer != nil && design.Spec.TorpedoJamming < math.Min(targetJamming, rules.JammerCap[hull.Starbase]) {
							itemToPlace = jammer
							break
						}
						fallthrough
					default:
						// more computers
						if design.Purpose.IsTorpedoShip() && computer != nil &&
							design.Spec.TorpedoBonus < design.Spec.TorpedoJamming {
							itemToPlace = computer
							break
						} else if jammer != nil && design.Spec.TorpedoJamming < rules.JammerCap[hull.Starbase] {
							itemToPlace = jammer
							break
						}
					}

					fallthrough // nothing else matched so we might as well check for elec deflectors/jets as a failsafe
				case HullSlotTypeMechanical:
					deflector, err := tc.GetBestComponentWithTag(rules, player, design, hst, TechTagBeamDeflector)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagBeamDeflector, err)
					}
					jet, err := tc.GetBestComponentWithTag(rules, player, design, hst, TechTagManeuveringJet)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagManeuveringJet, err)
					}

					if mechCounter == 1 && jet != nil &&
						!hull.Starbase && design.Spec.Movement < targetMove {
						itemToPlace = jet
					} else if deflector != nil {
						itemToPlace = deflector
					} else if jet != nil && !hull.Starbase && design.Spec.Movement < 10 {
						// TODO: change code once max battle movement is added to BattleRules
						itemToPlace = jet
					}

					if itemToPlace != nil {
						break
					}
					fallthrough // just slap on a fuel tank dammit
				default:
					fuelTank, err := tc.GetBestComponentWithTag(rules, player, design, hst, TechTagFuelTank)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagFuelTank, err)
					}
					itemToPlace = fuelTank
				}
				// toggle our counters
				elecCounter = -elecCounter
				mechCounter = -mechCounter
			} else {
				// (#) MULTI-CATEGORY ITEM PLACEMENT

				var weapon, sapper, driver, stargate, scanner, jet *TechHullComponent
				var err error
				if design.Purpose.IsTorpedoShip() {
					// add a big long ~~cucumber~~ missile
					weapon, err = tc.GetBestComponentWithTag(rules, player, design, hst, TechTagTorpedo)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagTorpedo, err)
					}
				} else {
					// add beems to our beem sheep
					weapon, err = tc.GetBestComponentWithTag(rules, player, design, hst, TechTagBeamWeapon)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagBeamWeapon, err)
					}
				}
				driver, err = tc.GetBestComponentWithTag(rules, player, design, hst, TechTagMassDriver)
				if err != nil {
					return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagMassDriver, err)
				}
				stargate, err = tc.GetBestComponentWithTag(rules, player, design, hst, TechTagStargate)
				if err != nil {
					return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagStargate, err)
				}
				scanner, err = tc.GetBestComponentWithTag(rules, player, design, hst, TechTagScanner)
				if err != nil {
					return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagScanner, err)
				}
				jet, err = tc.GetBestComponentWithTag(rules, player, design, hst, TechTagManeuveringJet)
				if err != nil {
					return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagManeuveringJet, err)
				}

				switch {
				case (numWeapons + numSappers) < minWeapons && weapon != nil:
					// first, if we don't have many weapons already, dedicate some of our flex slots to them
					sapper, err = tc.GetBestComponentWithTag(rules, player, design, hst, TechTagShieldSapper)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagShieldSapper, err)
					}
					
					var shouldUseSapper bool 
					// decide on whether to use sappers or not
					// TODO: Rework this once armor damage multi becomes a thing
					shouldUseSapper = sapper != nil && // have a sapper to use
						sapper.Range == weapon.Range && // sapper has at least as much range as our main guns
						numWeapons > numSappers*3 && // 3:1 gun:sapper ratio
						tc.compareWeaponPowers(rules, weapon, sapper) // sapper does more damage per hit
					if shouldUseSapper {
						itemToPlace = sapper
						numSappers += designSlot.Quantity
					} else if weapon != nil {
						itemToPlace = weapon
						numWeapons += designSlot.Quantity
					}
				case hull.Starbase && (driver != nil || (stargate != nil && numStargates == 0)):
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
				
				case design.Purpose == ShipDesignPurposeFighterScout && 
				!hasScanner && scanner != nil:
					// add scanners to armed scouts if they don't have them already
					itemToPlace = scanner
					designSlot.Quantity = 1 // only need 1 scanner, otherwise too expensive
					
				case AbsInt(elecCounter) == 1 && !hull.Starbase && jet != nil &&
					design.Spec.Movement < 10: // adjust movement target req based on slot #
					// add on jets if the counter is positioned right
					itemToPlace = jet
				default:
					// add whatever we need the most
					mostNeededItem, err := tc.getMostNeededComponent(rules, player, design, hst, designSlot.Quantity)
					if err != nil {
						return &ShipDesign{}, fmt.Errorf("getMostNeededComponent failed to get parts for tag %v, error: %w", TechTagManeuveringJet, err)
					}

					itemToPlace = mostNeededItem
				}

				// reset the counter so we don't get 20 million jets
				elecCounter += 1
				if elecCounter > 4 {
					elecCounter = 1
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

				// if all this item does is add beam bonus and/or jamming, tack it on a separate "reserved" list,
				// remove it from the design's slots temporarily and don't count it in our previous stat totals
				// (since they're used later on to avoid overcaps)
				// Spec fields still get updated as to not interfere with part placement logic
				isPureCapacitor := (itemToPlace.Tech.Tags.CountTags() == 1 && itemToPlace.Tags[TechTagBeamCapacitor]) ||
					(itemToPlace.Tech.Tags.CountTags() == 2 && (itemToPlace.Tags[TechTagBeamCapacitor] && itemToPlace.Tags[TechTagTorpedoJammer]))
				isPureJammer := (itemToPlace.Tech.Tags.CountTags() == 1 && itemToPlace.Tags[TechTagTorpedoJammer]) ||
					(itemToPlace.Tech.Tags.CountTags() == 2 && (itemToPlace.Tags[TechTagBeamCapacitor] && itemToPlace.Tags[TechTagTorpedoJammer]))
				if isPureCapacitor {
					capacitorSlots = append(capacitorSlots, slotNum)
					prevCapacitating = design.Spec.BeamBonus
				}
				if isPureJammer {
					jammerSlots = append(jammerSlots, slotNum)
					prevJamming = design.Spec.TorpedoJamming
				}
				designSlot.HullComponent = itemToPlace.Name
				design.Slots = append(design.Slots, designSlot)
				hasScanner = hasScanner || (itemToPlace.ScanRange >= 0 || itemToPlace.ScanRangePen >= 0)

				var err error
				design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
				if err != nil {
					return &ShipDesign{}, fmt.Errorf("computeShipDesignSpec errored during warship part allocation, error: %w", err)
				}

				if isPureJammer || isPureCapacitor {
					design.Slots = design.Slots[:len(design.Slots)-1] // re-slice w/o last element, functionally "removing" it
				} else {
					prevCapacitating = design.Spec.BeamBonus
					prevJamming = design.Spec.TorpedoJamming
				}
			}
		}
	}

	// (#) JAMMERS, CAPACITORS & WRAP-UP
	// add on our long lost capacitor & jammer friends
	if len(capacitorSlots) > 0 {
	capLoop:
		for i, id := range capacitorSlots {
			// place on our best capacitor into the slot
			hullSlot := hull.Slots[id]
			capacitor, err := tc.GetBestComponentWithTag(rules, player, design, hullSlot.Type, TechTagBeamCapacitor)
			if err != nil {
				return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagBeamCapacitor, err)
			}
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
		// undo jammer penalty temporarily to make calculations easier;
		// we recompute the spec immediately after this anyhow so no permanent harm done
		design.Spec.TorpedoJamming /= rules.JammerMulti[hull.Starbase]
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
			jammer, err := tc.GetBestComponentWithTag(rules, player, design, hullSlot.Type, TechTagTorpedoJammer)
			if err != nil {
				return &ShipDesign{}, fmt.Errorf("GetBestComponentWithTag failed to get parts for tag %v, error: %w", TechTagTorpedoJammer, err)
			}
			slot := ShipDesignSlot{HullComponent: jammer.Name, HullSlotIndex: id + 1}

			// add them one by one to make sure we don't go overboard
			// we don't actually
			for slot.Quantity = 1; slot.Quantity <= hullSlot.Capacity; slot.Quantity++ {
				prevJamming = roundFloat(1-(1-design.Spec.TorpedoJamming)*jammer.TorpedoJamming, 4)
				if prevJamming > rules.JammerCap[hull.Starbase] {
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
	
	// if we happen to have any unfilled slots left over, try to fill em up as a last minute failsafe
	// probably unnecessary but I like to be secure
	leftoverSlots := slices.Clip(slices.DeleteFunc(slices.Clone(design.Slots), func(sd ShipDesignSlot) bool {
		return sd.Quantity != 0 // zero out all previously filled slots
	}))
	if len(leftoverSlots) > 0 {
		for _, designSlot := range leftoverSlots {
			// place our best items into any empty slots
			id := designSlot.HullSlotIndex
			hullSlot := hull.Slots[id-1]
			hst := hullSlot.Type
			qty := hullSlot.Capacity

			var err error
			design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
			if err != nil {
				return &ShipDesign{}, fmt.Errorf("computeShipDesignSpec errored during warship part allocation, error: %w", err)
			}
			mostNeededItem, err := tc.getMostNeededComponent(rules, player, design, hst, qty)
			if err != nil {
				return &ShipDesign{}, fmt.Errorf("getMostNeededComponent failed to get most needed part, error: %w", err)
			}
			if mostNeededItem.HullSlotType&(HullSlotTypeWeaponShield|HullSlotTypeArmor) != 0 {
				if design.Purpose == ShipDesignPurposeStarbaseHalf {
					qty /= 2
				} else if design.Purpose == ShipDesignPurposeStarbaseQuarter {
					qty /= 4
				}
			}
			design.Slots = append(design.Slots, ShipDesignSlot{HullComponent: mostNeededItem.Name, HullSlotIndex: id, Quantity: qty})
		}
	}

	// finally, fix all the various temporary jank we did to the ship
	var err error
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
