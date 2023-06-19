package cs

import (
	"fmt"
	"math"
	"strings"
)

type ShipDesign struct {
	GameDBObject
	Num           int               `json:"num,omitempty"`
	PlayerNum     int               `json:"playerNum"`
	Name          string            `json:"name"`
	Dirty         bool              `json:"-"`
	Version       int               `json:"version"`
	Hull          string            `json:"hull"`
	HullSetNumber int               `json:"hullSetNumber"`
	CanDelete     bool              `json:"canDelete,omitempty"`
	Slots         []ShipDesignSlot  `json:"slots"`
	Purpose       ShipDesignPurpose `json:"purpose,omitempty"`
	Spec          ShipDesignSpec    `json:"spec"`
}

type ShipDesignSlot struct {
	HullComponent string `json:"hullComponent"`
	HullSlotIndex int    `json:"hullSlotIndex"`
	Quantity      int    `json:"quantity"`
	hc            *TechHullComponent
}

type ShipDesignSpec struct {
	HullType                  TechHullType          `json:"hullType,omitempty"`
	Engine                    Engine                `json:"engine,omitempty"`
	NumEngines                int                   `json:"numEngines,omitempty"`
	Cost                      Cost                  `json:"cost,omitempty"`
	Mass                      int                   `json:"mass,omitempty"`
	Armor                     int                   `json:"armor,omitempty"`
	FuelCapacity              int                   `json:"fuelCapacity,omitempty"`
	CargoCapacity             int                   `json:"cargoCapacity,omitempty"`
	CloakUnits                int                   `json:"cloakUnits,omitempty"`
	ScanRange                 int                   `json:"scanRange,omitempty"`
	ScanRangePen              int                   `json:"scanRangePen,omitempty"`
	RepairBonus               float64               `json:"repairBonus,omitempty"`
	TorpedoInaccuracyFactor   float64               `json:"torpedoInaccuracyFactor,omitempty"`
	Initiative                int                   `json:"initiative,omitempty"`
	Movement                  int                   `json:"movement,omitempty"`
	PowerRating               int                   `json:"powerRating,omitempty"`
	Bomber                    bool                  `json:"bomber,omitempty"`
	Bombs                     []Bomb                `json:"bombs,omitempty"`
	SmartBombs                []Bomb                `json:"smartBombs,omitempty"`
	RetroBombs                []Bomb                `json:"retroBombs,omitempty"`
	Scanner                   bool                  `json:"scanner,omitempty"`
	ImmuneToOwnDetonation     bool                  `json:"immuneToOwnDetonation,omitempty"`
	MineLayingRateByMineType  map[MineFieldType]int `json:"mineLayingRateByMineType,omitempty"`
	Shields                   int                   `json:"shield,omitempty"`
	Colonizer                 bool                  `json:"colonizer,omitempty"`
	Starbase                  bool                  `json:"starbase,omitempty"`
	CanLayMines               bool                  `json:"canLayMines,omitempty"`
	SpaceDock                 int                   `json:"spaceDock,omitempty"`
	MiningRate                int                   `json:"miningRate,omitempty"`
	TerraformRate             int                   `json:"terraformRate,omitempty"`
	MineSweep                 int                   `json:"mineSweep,omitempty"`
	CloakPercent              int                   `json:"cloakPercent,omitempty"`
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
	NumInstances              int                   `json:"numInstances,omitempty"`
	NumBuilt                  int                   `json:"numBuilt,omitempty"`
}

type MineLayingRateByMineType struct {
}

type ShipDesignPurpose string

const (
	ShipDesignPurposeNone              ShipDesignPurpose = ""
	ShipDesignPurposeScout             ShipDesignPurpose = "Scout"
	ShipDesignPurposeColonizer         ShipDesignPurpose = "Colonizer"
	ShipDesignPurposeBomber            ShipDesignPurpose = "Bomber"
	ShipDesignPurposeFighter           ShipDesignPurpose = "Fighter"
	ShipDesignPurposeFighterScout      ShipDesignPurpose = "FighterScout"
	ShipDesignPurposeCapitalShip       ShipDesignPurpose = "CapitalShip"
	ShipDesignPurposeFreighter         ShipDesignPurpose = "Freighter"
	ShipDesignPurposeColonistFreighter ShipDesignPurpose = "ColonistFreighter"
	ShipDesignPurposeFuelFreighter     ShipDesignPurpose = "FuelFreighter"
	ShipDesignPurposeArmedFreighter    ShipDesignPurpose = "ArmedFreighter"
	ShipDesignPurposeMiner             ShipDesignPurpose = "Miner"
	ShipDesignPurposeTerraformer       ShipDesignPurpose = "Terraformer"
	ShipDesignPurposeDamageMineLayer   ShipDesignPurpose = "DamageMineLayer"
	ShipDesignPurposeSpeedMineLayer    ShipDesignPurpose = "SpeedMineLayer"
	ShipDesignPurposeStarbase          ShipDesignPurpose = "Starbase"
	ShipDesignPurposeFort              ShipDesignPurpose = "Fort"
	ShipDesignPurposeStarterColony     ShipDesignPurpose = "StarterColony"
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

func ComputeShipDesignSpec(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design *ShipDesign) ShipDesignSpec {
	hull := rules.techs.GetHull(design.Hull)
	spec := ShipDesignSpec{
		Mass:                    hull.Mass,
		Armor:                   hull.Armor,
		FuelCapacity:            hull.FuelCapacity,
		Cost:                    hull.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec),
		CargoCapacity:           hull.CargoCapacity,
		CloakUnits:              raceSpec.BuiltInCloakUnits,
		Initiative:              hull.Initiative,
		TorpedoInaccuracyFactor: 1,
		ImmuneToOwnDetonation:   hull.ImmuneToOwnDetonation,
		RepairBonus:             hull.RepairBonus,
		ScanRange:               NoScanner,
		ScanRangePen:            NoScanner,
		SpaceDock:               hull.SpaceDock,
		Starbase:                hull.Starbase,
		HullType:                hull.Type,
	}

	numTachyonDetectors := 0

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

			spec.Mass += component.Mass * slot.Quantity
			spec.Armor += component.Armor * slot.Quantity
			spec.Shields += component.Shield * slot.Quantity
			spec.CargoCapacity += component.CargoBonus * slot.Quantity
			spec.FuelCapacity += component.FuelBonus * slot.Quantity
			spec.Colonizer = spec.Colonizer || component.ColonizationModule || component.OrbitalConstructionModule
			spec.Initiative += component.InitiativeBonus
			spec.Movement += component.MovementBonus * slot.Quantity
			spec.MiningRate += component.MiningRate * slot.Quantity
			spec.TerraformRate += component.TerraformRate * slot.Quantity
			spec.OrbitalConstructionModule = spec.OrbitalConstructionModule || component.OrbitalConstructionModule
			spec.CanStealFleetCargo = spec.CanStealFleetCargo || component.CanStealFleetCargo
			spec.CanStealPlanetCargo = spec.CanStealPlanetCargo || component.CanStealPlanetCargo

			// Add this mine type to the layers this design has
			if component.MineLayingRate > 0 {
				spec.CanLayMines = true
				if spec.MineLayingRateByMineType == nil {
					spec.MineLayingRateByMineType = make(map[MineFieldType]int)
				}
				if _, ok := spec.MineLayingRateByMineType[component.MineFieldType]; !ok {
					spec.MineLayingRateByMineType[component.MineFieldType] = 0
				}
				spec.MineLayingRateByMineType[component.MineFieldType] += component.MineLayingRate * slot.Quantity * hull.MineLayingFactor
			}

			// i.e. two .3f battle computers is (1 -.3) * (1 - .3) or (.7 * .7) or it decreases innaccuracy by 49%
			// so a 75% accurate torpedo would be 100 - (100 - 75) * .49 = 100 - 12.25 or 88% accurate
			// a 75% accurate torpedo with two 30% comps and one 50% comp would be
			// 100 - (100 - 75) * .7 * .7 * .5 = 94% accurate
			// if TorpedoInnaccuracyDecrease is 1 (default), it's just 75%
			spec.TorpedoInaccuracyFactor *= float64(math.Pow((1 - float64(component.TorpedoBonus)), float64(slot.Quantity)))

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
			}

			if component.Power > 0 {
				spec.HasWeapons = true
				spec.PowerRating += component.Power * slot.Quantity
				spec.WeaponSlots = append(spec.WeaponSlots, slot)
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
				spec.BasePacketSpeed = maxInt(spec.BasePacketSpeed, component.PacketSpeed)
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

	// determine the safe speed for this design
	spec.SafePacketSpeed = spec.BasePacketSpeed + spec.AdditionalMassDrivers

	// figure out the cloak as a percentage after we specd our cloak units
	spec.CloakPercent = getCloakPercentForCloakUnits(spec.CloakUnits)

	if numTachyonDetectors > 0 {
		// 95% ^ (SQRT(#_of_detectors) = Reduction factor for other player's cloaking (Capped at 81% or 17TDs)
		spec.ReduceCloaking = math.Pow((100.0-float64(rules.TachyonCloakReduction))/100, math.Sqrt(float64(numTachyonDetectors)))
	}

	if spec.NumEngines > 0 {
		// Movement = IdealEngineSpeed - 2 - Mass / 70 / NumEngines + NumManeuveringJets + 2*NumOverThrusters
		// we added any MovementBonus components above
		// we round up the slightest bit, and we can't go below 2, or above 10
		spec.Movement = clamp((spec.Engine.IdealSpeed-2)-spec.Mass/70/spec.NumEngines+spec.Movement+raceSpec.MovementBonus, 2, 10)
	} else {
		spec.Movement = 0
	}

	spec.computeScanRanges(rules, raceSpec.ScannerSpec, techLevels, design, hull)

	return spec
}

// Compute the scan ranges for this ship design The formula is: (scanner1**4 + scanner2**4 + ...
// + scannerN**4)**(.25)
func (spec *ShipDesignSpec) computeScanRanges(rules *Rules, scannerSpec ScannerSpec, techLevels TechLevel, design *ShipDesign, hull *TechHull) {
	spec.ScanRange = NoScanner
	spec.ScanRangePen = NoScanner

	// compute scanner as a built in JoaT scanner if it's built in
	builtInScannerMultiplier := scannerSpec.BuiltInScannerMultiplier
	if builtInScannerMultiplier > 0 && hull.BuiltInScanner {
		spec.ScanRange = techLevels.Electronics * builtInScannerMultiplier
		if !scannerSpec.NoAdvancedScanners {
			spec.ScanRangePen = int(math.Pow(float64(spec.ScanRange)/2, 4))
		}
		spec.ScanRange = int(math.Pow(float64(spec.ScanRange), 4))
	}

	for _, slot := range design.Slots {
		if slot.Quantity > 0 {
			component := rules.techs.GetHullComponent(slot.HullComponent)

			// bat scanners have 0 range
			if component.ScanRange != NoScanner {
				spec.ScanRange += int(math.Pow(float64(component.ScanRange), 4) * float64(slot.Quantity))
			}

			if component.ScanRangePen != NoScanner {
				spec.ScanRangePen += int((math.Pow(float64(component.ScanRangePen), 4)) * float64(slot.Quantity))
			}
		}
	}

	// now quad root it
	if spec.ScanRange != NoScanner {
		spec.ScanRange = int(math.Pow(float64(spec.ScanRange), .25) + .5)
		spec.ScanRange = int(float64(spec.ScanRange) * scannerSpec.ScanRangeFactor)
	}

	if spec.ScanRangePen != NoScanner {
		spec.ScanRangePen = int(math.Pow(float64(spec.ScanRangePen), .25) + .5)
	}

	// if we have no pen scan but we have a regular scan, set the pen scan range to 0
	if spec.ScanRangePen == NoScanner {
		if spec.ScanRange != NoScanner {
			spec.ScanRangePen = 0
		} else {
			spec.ScanRangePen = NoScanner
		}
	}

	// true if we have any scanning capability
	spec.Scanner = spec.ScanRange != NoScanner || spec.ScanRangePen != NoScanner
}

func DesignShip(techStore *TechStore, hull *TechHull, name string, player *Player, num int, hullSetNumber int, purpose ShipDesignPurpose) *ShipDesign {

	design := NewShipDesign(player, num).WithName(name).WithHull(hull.Name)
	design.Purpose = purpose

	engine := techStore.GetBestEngine(player)
	scanner := techStore.GetBestScanner(player)
	fuelTank := techStore.GetBestFuelTank(player)
	beamWeapon := techStore.GetBestBeamWeapon(player)
	shield := techStore.GetBestShield(player)
	armor := techStore.GetBestArmor(player)
	colonizationModule := techStore.GetBestColonizationModule(player)
	battleComputer := techStore.GetBestBattleComputer(player)
	miningRobot := techStore.GetBestMiningRobot(player)
	standardMineLayer := techStore.GetBestMineLayer(player, MineFieldTypeStandard)
	heavyMineLayer := techStore.GetBestMineLayer(player, MineFieldTypeHeavy)
	speedMineLayer := techStore.GetBestMineLayer(player, MineFieldTypeSpeedBump)

	numColonizationModules := 0
	numScanners := 0

	for i, hullSlot := range hull.Slots {
		slot := ShipDesignSlot{HullSlotIndex: i + 1}
		slot.Quantity = hullSlot.Capacity

		switch hullSlot.Type {
		case HullSlotTypeEngine:
			slot.HullComponent = engine.Name
		case HullSlotTypeScanner:
			numScanners++
			slot.HullComponent = scanner.Name
		case HullSlotTypeWeapon:
			slot.HullComponent = beamWeapon.Name
		case HullSlotTypeShieldArmor:
			fallthrough
		case HullSlotTypeArmor:
			slot.HullComponent = armor.Name
		case HullSlotTypeShield:
			slot.HullComponent = shield.Name
		case HullSlotTypeMining:
			slot.HullComponent = miningRobot.Name
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
		case HullSlotTypeElectrical:
			switch purpose {
			case ShipDesignPurposeCapitalShip:
				fallthrough
			case ShipDesignPurposeFighter:
				fallthrough
			case ShipDesignPurposeFighterScout:
				slot.HullComponent = battleComputer.Name
			default:
			}
		case HullSlotTypeMechanical:
			switch purpose {
			case ShipDesignPurposeColonizer:
				fallthrough
			case ShipDesignPurposeColonistFreighter:
				if colonizationModule != nil && numColonizationModules == 0 {
					numColonizationModules++
					slot.HullComponent = colonizationModule.Name
					slot.Quantity = 1 // we only need 1 colonization module
				} else {
					slot.HullComponent = fuelTank.Name
				}
			default:
				slot.HullComponent = fuelTank.Name
			}
		case HullSlotTypeScannerElectricalMechanical:
			switch purpose {
			case ShipDesignPurposeColonizer:
				fallthrough
			case ShipDesignPurposeColonistFreighter:
				if colonizationModule != nil && numColonizationModules == 0 {
					numColonizationModules++
					slot.HullComponent = colonizationModule.Name
					slot.Quantity = 1 // we only need 1 colonization module
				} else if numScanners == 0 {
					numScanners++
					slot.HullComponent = scanner.Name
				}
			default:
				if numScanners == 0 {
					numScanners++
					slot.HullComponent = scanner.Name
				} else {

				}
			}

		case HullSlotTypeArmorScannerElectricalMechanical:
			switch purpose {
			case ShipDesignPurposeColonizer:
				fallthrough
			case ShipDesignPurposeColonistFreighter:
				if colonizationModule != nil && numColonizationModules == 0 {
					numColonizationModules++
					slot.HullComponent = colonizationModule.Name
					slot.Quantity = 1 // we only need 1 colonization module
				} else if numScanners == 0 {
					slot.HullComponent = scanner.Name
					numScanners++
				} else {
					slot.HullComponent = fuelTank.Name
				}
			default:
				if numScanners == 0 {
					slot.HullComponent = scanner.Name
					numScanners++
				} else {
					slot.HullComponent = fuelTank.Name
				}
			}

		case HullSlotTypeGeneral:
			switch purpose {
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
