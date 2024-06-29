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
	Dirty             bool              `json:"-"`
	Version           int               `json:"version"`
	Hull              string            `json:"hull"`
	HullSetNumber     int               `json:"hullSetNumber"`
	CannotDelete      bool              `json:"cannotDelete,omitempty"`
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
	HullType                  TechHullType          `json:"hullType,omitempty"`
	Engine                    Engine                `json:"engine,omitempty"`
	NumEngines                int                   `json:"numEngines,omitempty"`
	Cost                      Cost                  `json:"cost,omitempty"`
	TechLevel                 TechLevel             `json:"techLevel,omitempty"`
	Mass                      int                   `json:"mass,omitempty"`
	Armor                     int                   `json:"armor,omitempty"`
	FuelCapacity              int                   `json:"fuelCapacity,omitempty"`
	FuelGeneration            int                   `json:"fuelGeneration,omitempty"`
	CargoCapacity             int                   `json:"cargoCapacity,omitempty"`
	CloakUnits                int                   `json:"cloakUnits,omitempty"`
	ScanRange                 int                   `json:"scanRange,omitempty"`
	ScanRangePen              int                   `json:"scanRangePen,omitempty"`
	InnateScanRangePenFactor  float64               `json:"innateScanRangePenFactor,omitempty"`
	RepairBonus               float64               `json:"repairBonus,omitempty"`
	TorpedoInaccuracyFactor   float64               `json:"torpedoInaccuracyFactor,omitempty"`
	TorpedoJamming            float64               `json:"torpedoJamming,omitempty"`
	BeamBonus                 float64               `json:"beamBonus,omitempty"`
	BeamDefense               float64               `json:"beamDefense,omitempty"`
	Initiative                int                   `json:"initiative,omitempty"`
	Movement                  int                   `json:"movement,omitempty"`
	ReduceMovement            int                   `json:"reduceMovement,omitempty"`
	PowerRating               int                   `json:"powerRating,omitempty"`
	Bomber                    bool                  `json:"bomber,omitempty"`
	Bombs                     []Bomb                `json:"bombs,omitempty"`
	SmartBombs                []Bomb                `json:"smartBombs,omitempty"`
	RetroBombs                []Bomb                `json:"retroBombs,omitempty"`
	Scanner                   bool                  `json:"scanner,omitempty"`
	ImmuneToOwnDetonation     bool                  `json:"immuneToOwnDetonation,omitempty"`
	MineLayingRateByMineType  map[MineFieldType]int `json:"mineLayingRateByMineType,omitempty"`
	Shields                   int                   `json:"shields,omitempty"`
	Colonizer                 bool                  `json:"colonizer,omitempty"`
	Starbase                  bool                  `json:"starbase,omitempty"`
	CanLayMines               bool                  `json:"canLayMines,omitempty"`
	SpaceDock                 int                   `json:"spaceDock,omitempty"`
	MiningRate                int                   `json:"miningRate,omitempty"`
	TerraformRate             int                   `json:"terraformRate,omitempty"`
	MineSweep                 int                   `json:"mineSweep,omitempty"`
	CloakPercent              int                   `json:"cloakPercent,omitempty"`
	CloakPercentFullCargo     int                   `json:"cloakPercentFullCargo,omitempty"`
	ReduceCloaking            float64               `json:"reduceCloaking,omitempty"`
	CanStealFleetCargo        bool                  `json:"canStealFleetCargo,omitempty"`
	CanStealPlanetCargo       bool                  `json:"canStealPlanetCargo,omitempty"`
	OrbitalConstructionModule bool                  `json:"orbitalConstructionModule,omitempty"`
	HasWeapons                bool                  `json:"hasWeapons,omitempty"`
	WeaponSlots               []ShipDesignSlot      `json:"weaponSlots,omitempty"`
	Stargate                  string                `json:"stargate,omitempty"`
	SafeHullMass              int                   `json:"safeHullMass,omitempty"`
	SafeRange                 int                   `json:"safeRange,omitempty"`
	MaxHullMass               int                   `json:"maxHullMass,omitempty"`
	MaxRange                  int                   `json:"maxRange,omitempty"`
	MassDriver                string                `json:"massDriver,omitempty"`
	SafePacketSpeed           int                   `json:"safePacketSpeed,omitempty"`
	BasePacketSpeed           int                   `json:"basePacketSpeed,omitempty"`
	AdditionalMassDrivers     int                   `json:"additionalMassDrivers,omitempty"`
	MaxPopulation             int                   `json:"maxPopulation,omitempty"`
	Radiating                 bool                  `json:"radiating,omitempty"`
	NumInstances              int                   `json:"numInstances,omitempty"`
	NumBuilt                  int                   `json:"numBuilt,omitempty"`
	EstimatedRange            int                   `json:"estimatedRange,omitempty"`
	EstimatedRangeFull        int                   `json:"estimatedRangeFull,omitempty"`
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
	return &ShipDesign{PlayerNum: player.Num, Num: num, Dirty: true, Slots: []ShipDesignSlot{}}
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
	sd.Spec = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, sd)
	return sd
}

func (sd *ShipDesign) MarkDirty() {
	sd.Dirty = true
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
func (d *ShipDesign) SlotsEqual(other *ShipDesign) bool {
	if len(d.Slots) != len(other.Slots) {
		return false
	}
	for i, v := range d.Slots {
		if v != other.Slots[i] {
			return false
		}
	}
	return true
}

func ComputeShipDesignSpec(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design *ShipDesign) ShipDesignSpec {
	hull := rules.techs.GetHull(design.Hull)
	spec := ShipDesignSpec{
		Mass:                     hull.Mass,
		Armor:                    hull.Armor,
		FuelCapacity:             hull.FuelCapacity,
		FuelGeneration:           hull.FuelGeneration,
		Cost:                     hull.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec),
		TechLevel:                hull.Requirements.TechLevel,
		CargoCapacity:            hull.CargoCapacity,
		CloakUnits:               raceSpec.BuiltInCloakUnits,
		Initiative:               hull.Initiative,
		TorpedoInaccuracyFactor:  1,
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

	torpedoJammingFactor := 1.0
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
			spec.Cost = spec.Cost.Add(component.Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec).MultiplyInt(slot.Quantity))
			spec.TechLevel = spec.TechLevel.Max(component.Requirements.TechLevel)

			spec.Mass += component.Mass * slot.Quantity
			spec.Armor += int(float64(component.Armor)*raceSpec.ArmorStrengthFactor) * slot.Quantity
			spec.Shields += int(float64(component.Shield)*raceSpec.ShieldStrengthFactor) * slot.Quantity
			spec.CargoCapacity += component.CargoBonus * slot.Quantity
			spec.FuelCapacity += component.FuelBonus * slot.Quantity
			spec.FuelGeneration += component.FuelGeneration * slot.Quantity
			spec.Colonizer = spec.Colonizer || component.ColonizationModule || component.OrbitalConstructionModule
			spec.Initiative += component.InitiativeBonus
			spec.Movement += component.MovementBonus * slot.Quantity
			spec.ReduceMovement = MaxInt(spec.ReduceMovement, component.ReduceMovement) // these don't stack
			spec.MiningRate += component.MiningRate * slot.Quantity
			spec.TerraformRate += component.TerraformRate * slot.Quantity
			spec.OrbitalConstructionModule = spec.OrbitalConstructionModule || component.OrbitalConstructionModule
			spec.CanStealFleetCargo = spec.CanStealFleetCargo || component.CanStealFleetCargo
			spec.CanStealPlanetCargo = spec.CanStealPlanetCargo || component.CanStealPlanetCargo
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

			// i.e. two .3f battle computers is (1 -.3) * (1 - .3) or (.7 * .7) or it decreases innaccuracy by 49%
			// so a 75% accurate torpedo would be 100 - (100 - 75) * .49 = 100 - 12.25 or 88% accurate
			// a 75% accurate torpedo with two 30% comps and one 50% comp would be
			// 100 - (100 - 75) * .7 * .7 * .5 = 94% accurate
			// if TorpedoInnaccuracyDecrease is 1 (default), it's just 75%
			spec.TorpedoInaccuracyFactor *= float64(math.Pow((1 - float64(component.TorpedoBonus)), float64(slot.Quantity)))
			torpedoJammingFactor *= float64(math.Pow((1 - float64(component.TorpedoJamming)), float64(slot.Quantity)))

			// beam bonuses
			spec.BeamBonus += component.BeamBonus * float64(slot.Quantity)
			spec.BeamDefense += component.BeamDefense * float64(slot.Quantity)

			// if this slot has a bomb, this design is a bomber
			if component.HullSlotType == HullSlotTypeBomb || component.MinKillRate > 0 {
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
				bombsPower += int((bomb.KillRate+bomb.StructureDestroyRate)*10) * slot.Quantity * 2
			}

			if component.Power > 0 {
				spec.HasWeapons = true
				spec.WeaponSlots = append(spec.WeaponSlots, slot)
				switch component.Category {
				case TechCategoryBeamWeapon:
					// beams contribute to the rating based on range, but sappers
					// are 1/3rd rated
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
			// cargo and space doc that are built into the hull
			// the space dock assumes that there is only one slot like that
			// it won't add them up

			if hullSlot.Type&HullSlotTypeSpaceDock > 0 {
				spec.SpaceDock = hullSlot.Capacity
			}

			// mass drivers
			if component.PacketSpeed > 0 {
				// if we already have a massdriver at this speed, add an additional mass driver to up
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
	if hull.Starbase {
		spec.CloakUnits += raceSpec.BuiltInCloakUnits
		spec.Cost = spec.Cost.MultiplyFloat64(raceSpec.StarbaseCostFactor)
	}

	spec.TorpedoJamming = 1 - torpedoJammingFactor

	// determine the safe speed for this design
	spec.SafePacketSpeed = spec.BasePacketSpeed + spec.AdditionalMassDrivers

	// figure out the cloak as a percentage after we specd our cloak units
	spec.CloakPercent = getCloakPercentForCloakUnits(spec.CloakUnits)
	spec.CloakPercentFullCargo = getCloakPercentForCloakUnits(int(math.Round(float64(spec.CloakUnits) * float64(spec.Mass) / float64(spec.Mass+spec.CargoCapacity))))

	if numTachyonDetectors > 0 {
		// 95% ^ (SQRT(#_of_detectors) = Reduction factor for other player's cloaking (Capped at 81% or 17TDs)
		spec.ReduceCloaking = math.Pow((100.0-float64(rules.TachyonCloakReduction))/100, math.Sqrt(float64(numTachyonDetectors)))
	}

	if spec.NumEngines > 0 {
		// Movement = IdealEngineSpeed - 2 - Mass / 70 / NumEngines + NumManeuveringJets + 2*NumOverThrusters
		// we added any MovementBonus components above
		// we round up the slightest bit, and we can't go below 2, or above 10
		spec.Movement = Clamp((spec.Engine.IdealSpeed-2)-spec.Mass/70/spec.NumEngines+spec.Movement+raceSpec.MovementBonus, 2, 10)
	} else {
		spec.Movement = 0
	}

	beamPower = int(float64(beamPower) * (1 + spec.BeamBonus))
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
	return spec
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
			slot.HullComponent = engine.Name
		case HullSlotTypeScanner:
			numScanners++
			slot.HullComponent = scanner.Name
		case HullSlotTypeWeapon:
			if numTorpedos > numBeamWeapons {
				slot.HullComponent = beamWeapon.Name
				numBeamWeapons++
			} else {
				slot.HullComponent = torpedo.Name
				numTorpedos++
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

			// if we are choosing shield or armor, pick  armor first, then shield
			if numArmors > numShields {
				slot.HullComponent = shield.Name
			} else {
				slot.HullComponent = armor.Name
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
			if miningRobot != nil {
				slot.HullComponent = miningRobot.Name
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
			// first. By default add packet throwers, then stargates, then battle computers
			switch purpose {
			case ShipDesignPurposePacketThrower:
				if packetThrower != nil {
					slot.HullComponent = packetThrower.Name
					numPacketThrowers++
				}
			case ShipDesignPurposeStargater:
				if stargate != nil {
					slot.HullComponent = packetThrower.Name
					numStargates++
				}
			default:
				// packet throwers for defense, then stargates
				if numPacketThrowers == 0 && packetThrower != nil {
					slot.HullComponent = packetThrower.Name
					numPacketThrowers++
				} else if numStargates == 0 && stargate != nil {
					slot.HullComponent = stargate.Name
					numStargates++
				}
				if slot.HullComponent == "" && hullSlot.Type == HullSlotTypeOrbitalElectrical {
					slot.HullComponent = battleComputer.Name
				}
			}
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
			case ShipDesignPurposeFuelFreighter:
				slot.HullComponent = fuelTank.Name
				numFuelTanks++
			case ShipDesignPurposeFreighter:
				fallthrough
			case ShipDesignPurposeColonistFreighter:
				// add cargo pods to freighters if we have a ramscoop
				if engine.FreeSpeed > 1 && cargoPod != nil {
					slot.HullComponent = cargoPod.Name
					numCargoPods++
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
						numCargoPods++
					} else {
						slot.HullComponent = fuelTank.Name
						numFuelTanks++
					}
				}
			default:
				slot.HullComponent = fuelTank.Name
				numFuelTanks++
			}
		case HullSlotTypeScannerElectricalMechanical:
			switch purpose {
			case ShipDesignPurposeFuelFreighter:
				slot.HullComponent = fuelTank.Name
				numFuelTanks++
			case ShipDesignPurposeFreighter:
				fallthrough
			case ShipDesignPurposeColonistFreighter:
				// add cargo pods to freighters if we have a ramscoop
				if engine.FreeSpeed > 1 && cargoPod != nil {
					slot.HullComponent = cargoPod.Name
					numCargoPods++
				} else {
					slot.HullComponent = fuelTank.Name
					numFuelTanks++
				}
			case ShipDesignPurposeColonizer:
				if colonizationModule != nil && numColonizationModules == 0 {
					numColonizationModules++
					slot.HullComponent = colonizationModule.Name
					slot.Quantity = 1 // we only need 1 colonization module
				} else {
					// balance fuel and cargo, fuel firsts
					if numFuelTanks > numCargoPods && cargoPod != nil {
						slot.HullComponent = cargoPod.Name
						numCargoPods++
					} else {
						slot.HullComponent = fuelTank.Name
						numFuelTanks++
					}
				}
			default:
				if numScanners == 0 {
					numScanners++
					slot.HullComponent = scanner.Name
				} else {
					// can always use more
					slot.HullComponent = fuelTank.Name
					numFuelTanks++
				}
			}

		case HullSlotTypeArmorScannerElectricalMechanical:
			switch purpose {
			case ShipDesignPurposeFuelFreighter:
				slot.HullComponent = fuelTank.Name
				numFuelTanks++
			case ShipDesignPurposeColonizer:
				if colonizationModule != nil && numColonizationModules == 0 {
					numColonizationModules++
					slot.HullComponent = colonizationModule.Name
					slot.Quantity = 1 // we only need 1 colonization module
				} else { // balance fuel and cargo, fuel firsts
					if numFuelTanks > numCargoPods && cargoPod != nil {
						slot.HullComponent = cargoPod.Name
						numCargoPods++
					} else {
						slot.HullComponent = fuelTank.Name
						numFuelTanks++
					}
				}
			default:
				if numScanners == 0 {
					slot.HullComponent = scanner.Name
					numScanners++
				} else {
					slot.HullComponent = fuelTank.Name
					numFuelTanks++
				}
			}

		case HullSlotTypeGeneral:
			switch purpose {
			case ShipDesignPurposeFuelFreighter:
				slot.HullComponent = fuelTank.Name
				numFuelTanks++
			case ShipDesignPurposeColonizer:
				// balance fuel and cargo, fuel firsts
				if numFuelTanks > numCargoPods && cargoPod != nil {
					slot.HullComponent = cargoPod.Name
					numCargoPods++
				} else {
					slot.HullComponent = fuelTank.Name
					numFuelTanks++
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
					numFuelTanks++
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
