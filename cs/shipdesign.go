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
	ShipDesignPurposeFighter               ShipDesignPurpose = "Fighter"
	ShipDesignPurposeFighterScout          ShipDesignPurpose = "FighterScout"
	ShipDesignPurposeCapitalShip           ShipDesignPurpose = "CapitalShip"
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
				gattlingMultiplier := 1
				if component.Gattling {
					// gattlings are 4x more mine-sweepery (all gatlings have range of 2)
					// lol, 4x, get it?
					gattlingMultiplier = component.Range * component.Range
				}
				spec.MineSweep += slot.Quantity * component.Power * ((component.Range + hull.RangeBonus) * component.Range) * gattlingMultiplier
			}

			spec.TechLevel = spec.TechLevel.Max(component.Requirements.TechLevel)

			spec.Mass += component.Mass * slot.Quantity
			spec.Armor += int(float64(component.Armor)*raceSpec.ArmorStrengthFactor) * slot.Quantity
			spec.Shields += int(float64(component.Shield)*raceSpec.ShieldStrengthFactor) * slot.Quantity
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
		spec.ReduceCloaking = math.Min(math.Pow((100.0-float64(rules.TachyonCloakReduction))/100, math.Sqrt(float64(numTachyonDetectors))), float64(rules.TachyonMaxCloakReduction)/100)
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
		spec.TorpedoJamming = math.Min(.95, roundFloat(spec.TorpedoJamming, 4))
	}

	// beam bonus defaults to 1
	spec.BeamBonus = 1
	if len(beamBoostersByCount) > 0 {
		for beamBonus, count := range beamBoostersByCount {
			// for 3 flux caps, this calc is 1-(1.2^3) for 1.728x beam damage
			bonus := math.Pow(1+beamBonus, float64(count))

			// multiple beam boosters stack multiplicatively
			spec.BeamBonus *= bonus
		}

		// Return final % bonus, rounded to 4 decimal places and capped at 155% base damage
		spec.BeamBonus = math.Min(roundFloat(spec.BeamBonus, 4), 2.55)
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

		// a movement of 1 1/2 int the UI (i.e. 6) doesn't impact your beam
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

// Compute the scan ranges for this ship design The formula is: (scanner1**4 + scanner2**4 + ...
// + scannerN**4)**(.25)
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

func DesignShip(techStore *TechStore, hull *TechHull, name string, player *Player, num int, hullSetNumber int, purpose ShipDesignPurpose, fleetPurpose FleetPurpose) *ShipDesign {

	design := NewShipDesign(player, num).WithName(name).WithHull(hull.Name)
	design.Purpose = purpose

	// fuel depots are empty
	if purpose == ShipDesignPurposeFuelDepot {
		return design
	}

	battleEngine := techStore.GetBestBattleEngine(player, hull)
	engine := techStore.GetBestEngine(player, hull, fleetPurpose)
	scanner := techStore.GetBestScanner(player)
	fuelTank := techStore.GetBestFuelTank(player)
	cargoPod := techStore.GetBestCargoPod(player)
	beamWeapon := techStore.GetBestBeamWeapon(player)
	torpedo := techStore.GetBestTorpedo(player)
	bomb := techStore.GetBestBomb(player)
	smartBomb := techStore.GetBestSmartBomb(player)
	structureBomb := techStore.GetBestStructureBomb(player)
	shield := techStore.GetBestShield(player)
	armor := techStore.GetBestArmor(player)
	colonizationModule := techStore.GetBestColonizationModule(player)
	battleComputer := techStore.GetBestBattleComputer(player)
	miningRobot := techStore.GetBestMiningRobot(player)
	terraformRobot := techStore.GetBestTerraformRobot(player)
	standardMineLayer := techStore.GetBestMineLayer(player, MineFieldTypeStandard)
	heavyMineLayer := techStore.GetBestMineLayer(player, MineFieldTypeHeavy)
	speedMineLayer := techStore.GetBestMineLayer(player, MineFieldTypeSpeedBump)
	packetThrower := techStore.GetBestPacketThrower(player)
	stargate := techStore.GetBestStargate(player)

	numColonizationModules := 0
	numScanners := 0
	numBeamWeapons := 0
	numTorpedos := 0
	numArmors := 0
	numShields := 0
	numFuelTanks := 0
	numCargoPods := 0
	numPacketThrowers := 0
	numStargates := 0

	for i, hullSlot := range hull.Slots {
		slot := ShipDesignSlot{HullSlotIndex: i + 1}
		slot.Quantity = hullSlot.Capacity

		// reduce quantity of armor and weapons when designing starbases for defense
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
		}

		switch hullSlot.Type {
		case HullSlotTypeEngine:
			if purpose == ShipDesignPurposeFighter {
				// need them battleships to be speedy!
				slot.HullComponent = battleEngine.Name
			} else {
				slot.HullComponent = engine.Name
			}
		case HullSlotTypeScanner:
			slot.HullComponent = scanner.Name
			numScanners++
		case HullSlotTypeWeapon:
			if purpose == ShipDesignPurposeFighterScout {
				slot.HullComponent = beamWeapon.Name
				numBeamWeapons++
			} else {
				if numTorpedos > numBeamWeapons {
					slot.HullComponent = beamWeapon.Name
					numBeamWeapons++ // TODO: Split fighters into 2 classes
				} else {
					slot.HullComponent = torpedo.Name
					numTorpedos++
				}
			}
		case HullSlotTypeBomb:
			// fill the bomb slot based on the type of bomber we want
			// or leave it blank
			switch purpose {
			case ShipDesignPurposeSmartBomber:
				if smartBomb != nil {
					slot.HullComponent = smartBomb.Name
				}
			case ShipDesignPurposeStructureBomber:
				if structureBomb != nil {
					slot.HullComponent = structureBomb.Name
				}
			default:
				if bomb != nil {
					slot.HullComponent = bomb.Name
				}
			}
		case HullSlotTypeShieldArmor:
			// fuel freighters stay fast and loose
			if purpose == ShipDesignPurposeFuelFreighter {
				continue
			}

			// if we are choosing shield or armor, pick shield first, then armor
			if numShields > numArmors {
				slot.HullComponent = armor.Name
				numArmors += slot.Quantity
			} else {
				slot.HullComponent = shield.Name
				numShields += slot.Quantity
			}
		case HullSlotTypeArmor:
			// fuel freighters stay fast and loose
			if purpose == ShipDesignPurposeFuelFreighter {
				continue
			}
			slot.HullComponent = armor.Name
			numArmors++
		case HullSlotTypeShield:
			// fuel freighters stay fast and loose
			if purpose == ShipDesignPurposeFuelFreighter {
				continue
			}
			slot.HullComponent = shield.Name
			numShields++
		case HullSlotTypeMining:
			if purpose == ShipDesignPurposeTerraformer {
				if terraformRobot != nil {
					slot.HullComponent = terraformRobot.Name
				}
			} else if purpose == ShipDesignPurposeMiner {
				if miningRobot != nil {
					slot.HullComponent = miningRobot.Name
				}
			}
		case HullSlotTypeMineLayer:
			switch purpose {
			case ShipDesignPurposeSpeedMineLayer:
				slot.HullComponent = speedMineLayer.Name
			default:
				if heavyMineLayer != nil {
					slot.HullComponent = heavyMineLayer.Name
				} else {
					slot.HullComponent = standardMineLayer.Name
				}
			}
		case HullSlotTypeOrbital:
			fallthrough
		case HullSlotTypeOrbitalElectrical:
			// if this starbase is designed for stargates or packet throwers, fill those
			// first. By default add packet throwers, then stargates, then electrical items

			switch purpose {
			case ShipDesignPurposePacketThrower:
				if packetThrower != nil {
					slot.HullComponent = packetThrower.Name
					numPacketThrowers++
					break
				}
			case ShipDesignPurposeStargater:
				if stargate != nil {
					slot.HullComponent = stargate.Name
					numStargates++
					break
				}
			default:
				// packet throwers for defense, then stargates
				if numPacketThrowers == 0 && packetThrower != nil {
					slot.HullComponent = packetThrower.Name
					numPacketThrowers++
					break
				} else if numStargates == 0 && stargate != nil {
					slot.HullComponent = stargate.Name
					numStargates++
					break
				}
			}
			// spare orbital slots left; use electrical items instead
			fallthrough
		case HullSlotTypeElectrical:
			// TODO: add in jammers, stealth, etc
			switch purpose {
			case ShipDesignPurposeCapitalShip:
				fallthrough
			case ShipDesignPurposeFighter:
				fallthrough
			case ShipDesignPurposeFighterScout:
				fallthrough
			default:
				slot.HullComponent = battleComputer.Name
			}
		case HullSlotTypeMechanical:
			switch purpose {
			case ShipDesignPurposeCapitalShip, ShipDesignPurposeFighter, ShipDesignPurposeFighterScout:
				fallthrough
			case ShipDesignPurposeFuelFreighter:
				slot.HullComponent = fuelTank.Name
				numFuelTanks += slot.Quantity
			case ShipDesignPurposeFreighter:
				fallthrough
			case ShipDesignPurposeColonistFreighter:
				// add cargo pods to freighters if we have a ramscoop
				if engine.FreeSpeed > 1 && cargoPod != nil {
					slot.HullComponent = cargoPod.Name
					numCargoPods += slot.Quantity
				}
			case ShipDesignPurposeColonizer:
				if colonizationModule != nil && numColonizationModules == 0 {
					numColonizationModules++
					slot.HullComponent = colonizationModule.Name
					slot.Quantity = 1 // we only need 1 colonization module
				} else {
					// balance fuel and cargo, fuel firsts
					if numFuelTanks > numCargoPods {
						slot.HullComponent = cargoPod.Name
						numCargoPods += slot.Quantity
					} else {
						slot.HullComponent = fuelTank.Name
						numFuelTanks += slot.Quantity
					}
				}
			default:
				slot.HullComponent = fuelTank.Name
				numFuelTanks += slot.Quantity
			}
		case HullSlotTypeElectricalMechanical:
			switch purpose {
			case ShipDesignPurposeFreighter:
				fallthrough
			case ShipDesignPurposeColonistFreighter:
				// add cargo pods to freighters if we have a ramscoop
				// up to 2 more than fuel tanks (because we still need _some_ fuel)
				if engine.FreeSpeed > 1 && cargoPod != nil && numCargoPods+2 > numFuelTanks {
					slot.HullComponent = cargoPod.Name
					numCargoPods += slot.Quantity
				} else {
					slot.HullComponent = fuelTank.Name
					numFuelTanks += slot.Quantity
				}
			case ShipDesignPurposeColonizer:
				if colonizationModule != nil && numColonizationModules == 0 {
					numColonizationModules++
					slot.HullComponent = colonizationModule.Name
					slot.Quantity = 1 // we only need 1 colonization module
				} else {
					// balance fuel and cargo, fuel first
					if numFuelTanks > numCargoPods && cargoPod != nil {
						slot.HullComponent = cargoPod.Name
						numCargoPods += slot.Quantity
					} else {
						slot.HullComponent = fuelTank.Name
						numFuelTanks += slot.Quantity
					}
				}
			default:
				// can always use more fuel
				slot.HullComponent = fuelTank.Name
				numFuelTanks += slot.Quantity
			}
		case HullSlotTypeScannerElectricalMechanical:
			switch purpose {
			case ShipDesignPurposeFuelFreighter:
				slot.HullComponent = fuelTank.Name
				numFuelTanks += slot.Quantity
			case ShipDesignPurposeFreighter:
				fallthrough
			case ShipDesignPurposeColonistFreighter:
				// add cargo pods to freighters if we have a ramscoop
				// up to 2 more than fuel tanks (because we still need _some_ fuel)
				if engine.FreeSpeed > 1 && cargoPod != nil && numCargoPods+2 > numFuelTanks {
					slot.HullComponent = cargoPod.Name
					numCargoPods += slot.Quantity
				} else {
					slot.HullComponent = fuelTank.Name
					numFuelTanks += slot.Quantity
				}
			case ShipDesignPurposeColonizer:
				if colonizationModule != nil && numColonizationModules == 0 {
					numColonizationModules++
					slot.HullComponent = colonizationModule.Name
					slot.Quantity = 1 // we only need 1 colonization module
				} else {
					// balance fuel and cargo, fuel first
					if numFuelTanks > numCargoPods && cargoPod != nil {
						slot.HullComponent = cargoPod.Name
						numCargoPods += slot.Quantity
					} else {
						slot.HullComponent = fuelTank.Name
						numFuelTanks += slot.Quantity
					}
				}
			default:
				if numScanners == 0 {
					slot.HullComponent = scanner.Name
					numScanners++
				} else {
					// can always use more fuel
					slot.HullComponent = fuelTank.Name
					numFuelTanks += slot.Quantity
				}
			}
		case HullSlotTypeArmorScannerElectricalMechanical:
			switch purpose {
			case ShipDesignPurposeFuelFreighter:
				slot.HullComponent = fuelTank.Name
				numFuelTanks += slot.Quantity
			case ShipDesignPurposeColonizer:
				if colonizationModule != nil && numColonizationModules == 0 {
					numColonizationModules++
					slot.HullComponent = colonizationModule.Name
					slot.Quantity = 1 // we only need 1 colonization module
				} else { // balance fuel and cargo, fuel firsts
					if numFuelTanks > numCargoPods && cargoPod != nil {
						slot.HullComponent = cargoPod.Name
						numCargoPods += slot.Quantity
					} else {
						slot.HullComponent = fuelTank.Name
						numFuelTanks += slot.Quantity
					}
				}
			default:
				if numScanners == 0 {
					slot.HullComponent = scanner.Name
					numScanners++
				} else {
					slot.HullComponent = fuelTank.Name
					numFuelTanks += slot.Quantity
				}
			}

		case HullSlotTypeGeneral:
			switch purpose {
			case ShipDesignPurposeFuelFreighter:
				slot.HullComponent = fuelTank.Name
				numFuelTanks += slot.Quantity
			case ShipDesignPurposeColonizer:
				// balance fuel and cargo, fuel firsts
				if numFuelTanks > numCargoPods && cargoPod != nil {
					slot.HullComponent = cargoPod.Name
					numCargoPods += slot.Quantity
				} else {
					slot.HullComponent = fuelTank.Name
					numFuelTanks += slot.Quantity
				}
			case ShipDesignPurposeFighter:
				fallthrough
			case ShipDesignPurposeFighterScout:
				if numScanners == 0 {
					slot.HullComponent = scanner.Name
					numScanners++
				} else {
					slot.HullComponent = beamWeapon.Name
				}
			default:
				if numScanners == 0 {
					slot.HullComponent = scanner.Name
					numScanners++
				} else {
					slot.HullComponent = fuelTank.Name
					numFuelTanks += slot.Quantity
				}
			}
		}

		// we filled it, add it
		if slot.HullComponent != "" {
			design.Slots = append(design.Slots, slot)
		}
	}

	return design
}
