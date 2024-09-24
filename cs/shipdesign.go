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
	Delete            bool              // used by the AI to mark a design for deletion
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
	ShipDesignPurposeFighterScout          ShipDesignPurpose = "FighterScout"    // armed scouts
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

// get the movement for this ship design, based on cargoMass
func (d *ShipDesign) getMovement(cargoMass int) int {
	return getBattleMovement(d.Spec.Engine.IdealSpeed, d.Spec.MovementBonus, d.Spec.Mass+cargoMass, d.Spec.NumEngines)
}

// return amount this jammer boosts our relative torpedo defense
func (spec ShipDesignSpec) getJamIncrease(rules *Rules, jamAmount float64, qty int) float64 {
	return math.Min(1-(1-spec.TorpedoJamming)*math.Pow(1-jamAmount, float64(qty)), rules.JammerCap[spec.Starbase]) / spec.TorpedoJamming
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

	for _, slot := range design.Slots {
		if slot.Quantity > 0 {
			component := rules.techs.GetHullComponent(slot.HullComponent)
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
					// gattlings are 4x more mine-sweepery (all gatlings have range of 2)
					// lol, 4x, get it?
					gatlingMultiplier = component.Range * component.Range
				}
				spec.MineSweep += slot.Quantity * component.Power * ((component.Range + hull.RangeBonus) * component.Range) * gatlingMultiplier
			}

			spec.TechLevel = spec.TechLevel.Max(component.Requirements.TechLevel)

			spec.Mass += component.Mass * slot.Quantity
			a, s := getActualArmorAmount(float64(component.Armor), float64(component.Shield), raceSpec.ArmorStrengthFactor, raceSpec.ShieldStrengthFactor, component.Category == TechCategoryArmor)
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
		// 95% ^ (SQRT(#_of_detectors) = Reduction factor for other player's cloaking (Capped at 81% or 17TDs)
		spec.ReduceCloaking = math.Min(math.Pow((100.0-float64(rules.TachyonCloakReduction))/100, math.Sqrt(float64(numTachyonDetectors))), 0.81)
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

		// the final jammer is the above sum inverted
		spec.TorpedoJamming = 1 - spec.TorpedoJamming

		// golang, why you be like this? nobody wants 1-.2^1 to be .199999994
		spec.TorpedoJamming = math.Min(roundFloat(spec.TorpedoJamming, 4), rules.JammerCap[hull.Starbase])
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
	design.PlayerNum = player.Num
	design.OriginalPlayerNum = player.Num

	// fuel depots & starter colonies are empty
	if purpose == ShipDesignPurposeFuelDepot || purpose == ShipDesignPurposeStarterColony {
		return design, nil
	} else if purpose == ShipDesignPurposeBeamFighter || purpose == ShipDesignPurposeTorpedoFighter || purpose == ShipDesignPurposeFighterScout || purpose == ShipDesignPurposeStarbase || purpose == ShipDesignPurposeStarbaseHalf || purpose == ShipDesignPurposeStarbaseQuarter {
		// warships get their own separate function for reasons
		design, err := DesignWarship(rules, hull, name, player, num, hullSetNumber, purpose)
		if err != nil {
			return &ShipDesign{}, fmt.Errorf("error %w in DesignWarship", err)
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
	engine := techStore.GetBestEngine(player, hull, fleetPurpose)

	numColonizationModules := 0
	numFuelTanks := 0
	numCargoPods := 0
	numScanners := 0
	numPacketThrowers := 0
	numStargates := 0
	numBeamWeapons := 0
	numTorpedos := 0

	for i, hullSlot := range hull.Slots {
		slot := ShipDesignSlot{HullSlotIndex: i + 1}
		slot.Quantity = hullSlot.Capacity
		hst := hullSlot.Type

		// first, we loop around once to catalog the best items per slot type in the design
		if bestPartsBySlot[hst] == nil && hst != HullSlotTypeEngine { // prevents double counting0
			bestPartsBySlot[hst] = map[TechTag]*TechHullComponent{}
			for _, tag := range techTagsToCheck {
				bestPartsBySlot[hst][tag] = techStore.GetBestComponentWithTags(player, hull, hst, false, false, tag)
			}
		}

		scanner := bestPartsBySlot[hst][TechTagScanner]
		beamWeapon := bestPartsBySlot[hst][TechTagBeamWeapon]
		torpedo := bestPartsBySlot[hst][TechTagTorpedo]
		bomb := bestPartsBySlot[hst][TechTagBomb]
		smartBomb := bestPartsBySlot[hst][TechTagSmartBomb]
		structureBomb := bestPartsBySlot[hst][TechTagStructureBomb]
		armor := bestPartsBySlot[hst][TechTagArmor]
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

		if hst&HullSlotTypeEngine > 0 {
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
		case ShipDesignPurposeStartingFighter: // everyone's favorite rinky dinky starter ship
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
			}
			fallthrough
		case ShipDesignPurposeFreighter, ShipDesignPurposeColonistFreighter:
			// tack on up to 1 extra cargo pod here if we have a ramscoop; otherwise slap on a fuel tank
			// TODO: Add purpose for cloaked ships and add cloaks accordingly
			if cargoPod != nil && engine.FreeSpeed > 1 && numCargoPods < numFuelTanks+1 {
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
			}
		case ShipDesignPurposeDamageMineLayer:
			if heavyMineLayer != nil {
				slot.HullComponent = heavyMineLayer.Name
			} else if standardMineLayer != nil {
				slot.HullComponent = standardMineLayer.Name
			}
		case ShipDesignPurposePacketThrower:
			if packetThrower != nil {
				slot.HullComponent = packetThrower.Name
				numPacketThrowers++
			}
			fallthrough
		case ShipDesignPurposeStargater:
			if stargate != nil {
				slot.HullComponent = stargate.Name
				numStargates++
			}
			// packet throwers for defense, then stargates
			if packetThrower != nil {
				slot.HullComponent = packetThrower.Name
				numPacketThrowers++
			} else if stargate != nil {
				slot.HullComponent = stargate.Name
				numStargates++
			}
		}

		if slot.HullComponent == "" && fuelTank != nil { // when in doubt, add fuel tanks to empty slots
			slot.HullComponent = fuelTank.Name
			numFuelTanks += slot.Quantity
		}

		// if we filled the slot, add it
		if slot.HullComponent != "" {
			design.Slots = append(design.Slots, slot)
		}
	}

	return design, nil
}

// get the best design for a warship or starbase based on available parts
func DesignWarship(rules *Rules, hull *TechHull, name string, player *Player, num int, hullSetNumber int, purpose ShipDesignPurpose) (*ShipDesign, error) {

	//* DISCLAIMER FOR CODE VIEWERS: THIS IS A *VERY LONG FUNCTION* DUE TO THE COMPLEXITY OF MAKING GOOD SHIP DESIGNS.
	//* Use the hashtags (#) to jump between sections.
	techStore := rules.techs
	design := NewShipDesign(player, num).WithName(name).WithHull(hull.Name).WithPurpose(purpose)

	// (#) HELPER VARIABLES AND LISTS/COUNTERS

	var qtyMulti int = 1 // quantity multiplier for starbases
	if purpose == ShipDesignPurposeStarbaseHalf {
		qtyMulti = 2
	} else if purpose == ShipDesignPurposeStarbaseQuarter {
		qtyMulti = 4
	}

	// whether this ship should use torpedoes or not (and battle computers by extension)
	var torpedoShip bool = purpose == ShipDesignPurposeTorpedoFighter || purpose == ShipDesignPurposeStarbase || purpose == ShipDesignPurposeStarbaseHalf || purpose == ShipDesignPurposeStarbaseQuarter
	// subset of tags to check through so we don't have to loop over mining robots and non-combat implements
	// exception is fuel pods which we use as a fallback if all else fails
	techTagsToCheck := []TechTag{
		TechTagBeamWeapon,
		TechTagShieldSapper,
		TechTagBeamCapacitor,
		TechTagBeamDeflector,
		TechTagTorpedo,
		TechTagFuelTank,
		TechTagTorpedoBonus,
		TechTagTorpedoJammer,
		TechTagStargate,
		TechTagMassDriver,
		TechTagArmor,
		TechTagShield,
		TechTagManeuveringJet,
		TechTagScanner,
	}

	// A sorted list of hull slots by flexibility (how many different item types you can put in it).
	//
	// We use this to determine what order to loop through things
	var hullSlotsByFlexibility = map[int][]int{}
	var bestPartsBySlot = map[HullSlotType]map[TechTag]*TechHullComponent{} // a map matching each hull slot to the best parts of each kind available for said slot
	var capacitorSlots = []int{}                                            // contains all the hull slots we are reserving for beam caps (these get checked last due to hardcap)
	var jammerSlots = []int{}                                               // contains all the hull slots we are reserving for jammers (also checked last due to hardcap)
	engine := techStore.GetBestBattleEngine(player, hull)
	numWeapons := 0
	numSappers := 0
	numPacketThrowers := 0
	numStargates := 0
	scanners := false

	// (#) CONSTANTS USED IN EXECUTION
	// These values are used in the code to influence part choice
	// Modifying these will change behavior accordingly

	targetMove := 9
	targetJamming := 0.8
	targetComputing := 0.8
	minWeapons := 9

	// (#) LOGIC AND PART RETRIEVAL
	var err error
	// Compute & initialize our design spec
	design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
	if err != nil {
		return &ShipDesign{}, fmt.Errorf("error %w in ComputeShipDesignSpec when calculating stats for warship part allocation", err)
	}

	maxNum := math.MinInt
	for i, hullSlot := range hull.Slots {
		b := Bitmask(hullSlot.Type).countBits()
		if b > maxNum {
			maxNum = b
		}
		hullSlotsByFlexibility[b] = append(hullSlotsByFlexibility[b], i+1)
		// while we're at it, catalog the best items per hull slot type in our design
		if bestPartsBySlot[hullSlot.Type] != nil { // prevents double counting
			for _, tag := range techTagsToCheck {
				bestPartsBySlot[hullSlot.Type][tag] = techStore.GetBestComponentWithTags(player, hull, hullSlot.Type, false, purpose != ShipDesignPurposeTorpedoFighter, tag)
			}
		}
	}

	// loop through items in order of increasing flexibility
	for i := 1; i <= maxNum; i++ {
		list := hullSlotsByFlexibility[i]
		if list == nil {
			// speed it up if empty
			continue
		}
		var elecCounter int = 1 // counter used to decide which parts to add
		var mechCounter int = 1 // counter used to decide which parts to add

	slotLoop:
		// extract slot numbers from list so we can loop through them
		for _, sn := range list {
			hullSlot := hull.Slots[sn-1]
			hst := hullSlot.Type
			slot := ShipDesignSlot{HullSlotIndex: sn} // list index 0 gets slot no. 1
			slot.Quantity = hullSlot.Capacity

			armor := bestPartsBySlot[hst][TechTagArmor]
			shield := bestPartsBySlot[hst][TechTagShield]
			torpedo := bestPartsBySlot[hst][TechTagTorpedo]
			beamWeapon := bestPartsBySlot[hst][TechTagBeamWeapon]
			sapper := bestPartsBySlot[hst][TechTagShieldSapper]
			capacitor := bestPartsBySlot[hst][TechTagBeamCapacitor]
			stargate := bestPartsBySlot[hst][TechTagStargate]
			driver := bestPartsBySlot[hst][TechTagMassDriver]
			deflector := bestPartsBySlot[hst][TechTagBeamDeflector]
			fuelTank := bestPartsBySlot[hst][TechTagFuelTank]
			computer := bestPartsBySlot[hst][TechTagTorpedoBonus]
			jammer := bestPartsBySlot[hst][TechTagTorpedoJammer]
			jet := bestPartsBySlot[hst][TechTagManeuveringJet]
			scanner := bestPartsBySlot[hst][TechTagScanner]
			var itemToPlace *TechHullComponent

			// (#) PART PLACEMENT LOGIC
			if i == 1 {
				// only 1 part type means we can just check the slot type
				switch {
				case hst == HullSlotTypeEngine: // we should ALWAYS HAVE AN ENGINE TO USE ON OUR SHIPS. IF NOT WE'RE SCREWED
					slot.HullComponent = engine.Name
					design.Slots = append(design.Slots, slot)
					design.Spec.Engine = engine.Engine
					design.Spec.NumEngines = slot.Quantity
					design.Spec.Mass += itemToPlace.Mass * slot.Quantity
					continue slotLoop
				case hst == HullSlotTypeArmor && armor != nil:
					itemToPlace = armor
					slot.Quantity /= qtyMulti
				case hst == HullSlotTypeShield && shield != nil:
					itemToPlace = shield
					slot.Quantity /= qtyMulti
				case hst == HullSlotTypeWeapon:
					if torpedoShip && torpedo != nil {
						// add our best torp/missile
						itemToPlace = torpedo
						slot.Quantity /= qtyMulti
						numWeapons += slot.Quantity
					} else {
						// add beems to our beem sheep

						// decide on whether to use sappers or not
						shouldUseSapper := sapper != nil && sapper.Range == beamWeapon.Range &&
							numWeapons*3 > numSappers // 3:1 beam:sapper ratio
						if shouldUseSapper {
							itemToPlace = sapper
							slot.Quantity /= qtyMulti
							numSappers += slot.Quantity
						} else if beamWeapon != nil {
							itemToPlace = beamWeapon
							slot.Quantity /= qtyMulti
							numWeapons += slot.Quantity
						}
					}
				case hst == HullSlotTypeScanner && scanner != nil:
					itemToPlace = scanner
					scanners = true
				case hst == HullSlotTypeOrbital:
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
				case hst == HullSlotTypeElectrical:
					// elec stuff (jammers/comps)
					switch elecCounter {
					case 1:
						// add beam caps and computers
						if torpedoShip && design.Spec.TorpedoBonus < targetComputing &&
							computer != nil {
							itemToPlace = computer
							break
						} else if !torpedoShip && capacitor != nil &&
							design.Spec.BeamBonus < rules.BeamBonusCap {
							itemToPlace = capacitor
							break
						}
						fallthrough
					case -1:
						// add jammers
						if jammer != nil && design.Spec.TorpedoJamming < targetJamming {
							itemToPlace = jammer
							break
						}
						fallthrough
					default:
						if torpedoShip && computer != nil &&
							design.Spec.TorpedoBonus < design.Spec.TorpedoJamming {
							itemToPlace = computer
						} else if jammer != nil && design.Spec.TorpedoJamming < rules.JammerCap[hull.Starbase] {
							itemToPlace = jammer
						}
					}

					if itemToPlace != nil {
						break
					}
					fallthrough
				case hst == HullSlotTypeMechanical && (deflector != nil || (jet != nil && !design.Spec.Starbase && design.Spec.Movement < 10)):
					if mechCounter == 1 && jet != nil && (!design.Spec.Starbase || design.Spec.Movement < targetMove) {
						itemToPlace = jet
					} else if deflector != nil {
						itemToPlace = deflector
					} else if jet != nil && !design.Spec.Starbase && design.Spec.Movement < 10 {
						itemToPlace = jet
					} else if fuelTank != nil {
						itemToPlace = fuelTank
					}
				}

				// toggle our counters
				elecCounter = -elecCounter
				mechCounter = -mechCounter
			} else {
				// here comes the fun part: figuring out how to rank stuff!

				// first, if we don't have many weapons already, dedicate our flex slots to them
				if (numWeapons + numSappers) < minWeapons {
					// add our best torp/missile
					if torpedoShip && torpedo != nil {
						itemToPlace = torpedo
						slot.Quantity /= qtyMulti
						numWeapons += slot.Quantity
					} else {
						// add beems to our beem sheep
						shouldUseSapper := sapper != nil && sapper.Range == beamWeapon.Range &&
							numWeapons*3 > numSappers // 3:1 beam:sapper ratio
						if shouldUseSapper {
							itemToPlace = sapper
							slot.Quantity /= qtyMulti
							numSappers += slot.Quantity
						} else if beamWeapon != nil {
							itemToPlace = beamWeapon
							slot.Quantity /= qtyMulti
							numWeapons += slot.Quantity
						}
					}
					break
				}

				// add orbital items to starbases
				if design.Spec.Starbase {
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
				}

				switch elecCounter {
				// add on items based on what we need the most
				case 1:
					// jets
					if jet != nil && !design.Spec.Starbase && design.Spec.Movement < targetMove {
						itemToPlace = jet
						break
					}
					fallthrough
				case 2:
					// computers/capacitors
					if torpedoShip && computer != nil &&
						design.Spec.TorpedoBonus <= min(targetComputing, design.Spec.BeamDefense, design.Spec.TorpedoJamming) {
						// our computing score is the lowest out of jamming/deflecting (ie we need it the most); add one on
						itemToPlace = computer
						break
					} else if purpose == ShipDesignPurposeFighterScout && scanner != nil && !scanners {
						// scouts need scanners to function
						itemToPlace = scanner
						scanners = true
						break
					} else if !torpedoShip && capacitor != nil &&
						// the estimated relative dmg increase of adding beam caps to this slot
						math.Min(math.Pow(capacitor.BeamBonus, float64(slot.Quantity)), rules.BeamBonusCap/design.Spec.BeamBonus) >=
							// the estimated relative damage taken decrease of adding jammers or deflectors to this slot
							math.Max(math.Pow(1+deflector.BeamDefense, float64(slot.Quantity)), design.Spec.getJamIncrease(rules, jammer.TorpedoJamming, slot.Quantity)) {
						itemToPlace = capacitor // add it to list; we'll come back for it later
						break
					}
					fallthrough
				case 3:
					// jammers
					if jammer != nil && design.Spec.TorpedoJamming < min(targetJamming, design.Spec.BeamDefense) {
						itemToPlace = jammer
						break
					}
					fallthrough
				case 4:
					// deflectors
					if deflector != nil {
						itemToPlace = deflector
						break
					}
					fallthrough
				case 5:
					// armors
					if design.Spec.Armor <= design.Spec.Shields && armor != nil {
						itemToPlace = armor
						slot.Quantity /= qtyMulti
						break
					}
					fallthrough
				case 6:
					// shields
					if design.Spec.Shields < design.Spec.Armor && shield != nil {
						itemToPlace = shield
						slot.Quantity /= qtyMulti
						break
					}
					fallthrough
				default:
					// everything above us failed; slap on either comps/caps, more jammers/deflectors, extra extra shields or a fuel tank as a last resort
					if torpedoShip && computer != nil &&
						design.Spec.TorpedoBonus <= min(design.Spec.BeamDefense, design.Spec.TorpedoJamming) {
						itemToPlace = computer
						break
					} else if purpose == ShipDesignPurposeFighterScout && scanner != nil && !scanners {
						itemToPlace = scanner
						scanners = true
						break
					} else if !torpedoShip && capacitor != nil &&
						math.Max(math.Pow(capacitor.BeamBonus, float64(slot.Quantity)), rules.BeamBonusCap/design.Spec.BeamBonus) >=
							math.Max(math.Pow(1+deflector.BeamDefense, float64(slot.Quantity)), design.Spec.getJamIncrease(rules, jammer.TorpedoJamming, slot.Quantity)) {
						itemToPlace = capacitor
						break
					} else if jammer != nil && design.Spec.TorpedoJamming <= math.Max(design.Spec.TorpedoJamming, design.Spec.BeamDefense) {
						itemToPlace = jammer
						break
					} else if deflector != nil {
						itemToPlace = deflector
					} else if shield != nil {
						itemToPlace = shield
						slot.Quantity /= qtyMulti
					} else if fuelTank != nil {
						itemToPlace = fuelTank
					}
				}

				elecCounter += 1
				if elecCounter > 6 {
					elecCounter = 0
				}
			}

			// we filled the slot, add it and recompute important fields in spec
			if itemToPlace != nil {
				design.Spec.Mass += itemToPlace.Mass * slot.Quantity
				if !design.Spec.Starbase {
					design.Spec.MovementBonus += itemToPlace.MovementBonus * slot.Quantity
					design.Spec.Movement = design.getMovement(0)
				}
				design.Spec.TorpedoBonus = 1 - (1-design.Spec.TorpedoBonus)*math.Pow(1-itemToPlace.TorpedoBonus, float64(slot.Quantity))
				design.Spec.TorpedoJamming = 1 - (1-design.Spec.TorpedoJamming)*math.Pow(1-itemToPlace.TorpedoJamming, float64(slot.Quantity))
				design.Spec.BeamDefense = 1 - (1-design.Spec.BeamDefense)*math.Pow(1-itemToPlace.BeamDefense, float64(slot.Quantity))
				design.Spec.BeamBonus = 1 - math.Min(design.Spec.BeamBonus*math.Pow(1+itemToPlace.BeamBonus, float64(slot.Quantity)), rules.BeamBonusCap) // failsafe so we don't get 0 power beams
				a, s := getActualArmorAmount(float64(itemToPlace.Armor), float64(itemToPlace.Shield), player.Race.Spec.ArmorStrengthFactor, player.Race.Spec.ShieldStrengthFactor, itemToPlace.Category == TechCategoryArmor)
				design.Spec.Armor += int(a)
				design.Spec.Shields += int(s)
				if itemToPlace.Tech.Tags.hasTag(TechTagBeamCapacitor) && itemToPlace.Power == 0 { // specifically excludes weapons
					capacitorSlots = append(capacitorSlots, sn) // tack it on the cap list; we'll circle back to it later
				} else if itemToPlace.Tech.Tags.hasTag(TechTagTorpedoJammer) && itemToPlace.Power == 0 {
					jammerSlots = append(jammerSlots, sn) // tack it on the jammer list; we'll circle back to it later
				} else {
					slot.HullComponent = itemToPlace.Name
					design.Slots = append(design.Slots, slot)
				}
			}
		}
	}

	// add on our long lost capacitor & jammer friends
	if len(capacitorSlots) > 0 {
		for _, id := range capacitorSlots {
			// place on our best capacitor into the slot
			hullSlot := hull.Slots[id]
			capacitor := bestPartsBySlot[hullSlot.Type][TechTagBeamCapacitor]
			slot := ShipDesignSlot{HullComponent: capacitor.Name, HullSlotIndex: id}

			// add them one by one to make sure we don't go overboard
			for slot.Quantity = 1; slot.Quantity <= hullSlot.Capacity; slot.Quantity++ {
				design.Spec.BeamBonus *= (1 + capacitor.BeamBonus)
				if design.Spec.BeamBonus > rules.BeamBonusCap {
					break
				}
			}
			// add the finished item to the hullSlot
			design.Slots[id] = slot
		}
	}

	if len(jammerSlots) > 0 {
		for _, id := range jammerSlots {
			// place our best jammer into the slot
			hullSlot := hull.Slots[id]
			jammer := bestPartsBySlot[hullSlot.Type][TechTagTorpedoJammer]
			slot := ShipDesignSlot{HullComponent: jammer.Name, HullSlotIndex: id}

			// add them one by one to make sure we don't go overboard
			for slot.Quantity = 1; slot.Quantity <= hullSlot.Capacity; slot.Quantity++ {
				design.Spec.TorpedoJamming = 1 - (1-design.Spec.TorpedoJamming)*jammer.TorpedoJamming
				if design.Spec.TorpedoJamming > rules.JammerCap[hull.Starbase] {
					break
				}
			}
			// add the finished item to the hullSlot
			design.Slots[id] = slot
		}
	}

	// do 1 last failsafe spec recalculation to fix all the jank we did to the ship
	design.Spec, err = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
	if err != nil {
		return &ShipDesign{}, fmt.Errorf("error %w in ComputeShipDesignSpec when calculating stats for warship part allocation", err)
	}

	return design, nil
}

/* 		// reduce quantity of armor and weapons when designing starbases for defense
if hullSlot.Type == HullSlotTypeArmor ||
	hullSlot.Type == HullSlotTypeShield ||
	hullSlot.Type == HullSlotTypeShieldArmor ||
	hullSlot.Type == HullSlotTypeWeaponShield ||
	hullSlot.Type == HullSlotTypeWeapon {
	if purpose == ShipDesignPurposeStarbaseQuarter {
		slot.Quantity = MaxInt(1, hullSlot.Capacity/4)
	} else if purpose == ShipDesignPurposeStarbaseHalf {
		slot.Quantity = MaxInt(1, hullSlot.Capacity/2)
	}
} */
