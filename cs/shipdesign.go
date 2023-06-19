package cs

import (
	"math"
	"time"

	"github.com/google/uuid"
)

type ShipDesign struct {
	ID            int64             `json:"id"`
	GameID        int64             `json:"gameId"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	PlayerNum     int               `json:"playerNum"`
	Dirty         bool              `json:"-"`
	UUID          uuid.UUID         `json:"uuid,omitempty"`
	Name          string            `json:"name"`
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
}

type ShipDesignSpec struct {
	IdealSpeed                int                   `json:"idealSpeed,omitempty"`
	Engine                    string                `json:"engine,omitempty"`
	FuelUsage                 [11]int               `json:"fuelUsage,omitempty"`
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
	ImmuneToOwnDetonation     bool                  `json:"immuneToOwnDetonation"`
	MineLayingRateByMineType  map[MineFieldType]int `json:"mineLayingRateByMineType"`
	Shield                    int                   `json:"shield,omitempty"`
	Colonizer                 bool                  `json:"colonizer,omitempty"`
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
	WeaponSlots               []ShipDesignSlot      `json:"weaponSlots"`
}

type MineLayingRateByMineType struct {
}

type ShipDesignPurpose string

const (
	ShipDesignPurposeNone              ShipDesignPurpose = ""
	ShipDesignPurposeScout             ShipDesignPurpose = "Scout"
	ShipDesignPurposeArmedScout        ShipDesignPurpose = "ArmedScout"
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

func NewShipDesign(player *Player) *ShipDesign {
	return &ShipDesign{GameID: player.GameID, PlayerNum: player.Num, UUID: uuid.New(), Dirty: true, Slots: []ShipDesignSlot{}}
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
	sd.Spec = computeShipDesignSpec(rules, player, sd)
	return sd
}

func computeShipDesignSpec(rules *Rules, player *Player, design *ShipDesign) ShipDesignSpec {
	hull := rules.techs.GetHull(design.Hull)
	spec := ShipDesignSpec{
		Mass:                    hull.Mass,
		Armor:                   hull.Armor,
		FuelCapacity:            hull.FuelCapacity,
		Cost:                    hull.GetPlayerCost(player),
		CargoCapacity:           hull.CargoCapacity,
		CloakUnits:              player.Race.Spec.BuiltInCloakUnits,
		Initiative:              hull.Initiative,
		TorpedoInaccuracyFactor: 1,
		ImmuneToOwnDetonation:   hull.ImmuneToOwnDetonation,
		RepairBonus:             hull.RepairBonus,
		ScanRange:               NoScanner,
		ScanRangePen:            NoScanner,
		SpaceDock:               hull.SpaceDock,
	}

	numTachyonDetectors := 0

	for _, slot := range design.Slots {
		if slot.Quantity > 0 {
			component := rules.techs.GetHullComponent(slot.HullComponent)
			hullSlot := hull.Slots[slot.HullSlotIndex-1]

			// record engine details
			if hullSlot.Type == HullSlotTypeEngine {
				engine := rules.techs.GetEngine(slot.HullComponent)
				spec.Engine = engine.Name
				spec.IdealSpeed = engine.IdealSpeed
				spec.NumEngines = slot.Quantity
				spec.FuelUsage = engine.FuelUsage
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
			spec.Cost = spec.Cost.Add(component.Tech.GetPlayerCost(player).MultiplyInt(slot.Quantity))

			spec.Mass += component.Mass * slot.Quantity
			spec.Armor += component.Armor * slot.Quantity
			spec.Shield += component.Shield * slot.Quantity
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
				if _, ok := spec.MineLayingRateByMineType[component.MineFieldType]; ok {
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
		}

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
			spec.Movement = clamp((spec.IdealSpeed-2)-spec.Mass/70/spec.NumEngines+spec.Movement+player.Race.Spec.MovementBonus, 2, 10)
		} else {
			spec.Movement = 0
		}
	}

	spec.computeScanRanges(rules, player, design, hull)

	return spec
}

// Compute the scan ranges for this ship design The formula is: (scanner1**4 + scanner2**4 + ...
// + scannerN**4)**(.25)
func (spec *ShipDesignSpec) computeScanRanges(rules *Rules, player *Player, design *ShipDesign, hull *TechHull) {
	spec.ScanRange = NoScanner
	spec.ScanRangePen = NoScanner

	// compute scanner as a built in JoaT scanner if it's built in
	builtInScannerMultiplier := player.Race.Spec.BuiltInScannerMultiplier
	if builtInScannerMultiplier > 0 && hull.BuiltInScanner {
		spec.ScanRange = player.TechLevels.Electronics * builtInScannerMultiplier
		if !player.Race.Spec.NoAdvancedScanners {
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
		spec.ScanRange = int(float64(spec.ScanRange) * player.Race.Spec.ScanRangeFactor)
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

func designShip(techStore *TechStore, hull *TechHull, name string, player *Player, hullSetNumber int, purpose ShipDesignPurpose) *ShipDesign {

	design := NewShipDesign(player).WithName(name).WithHull(hull.Name)
	design.Purpose = purpose

	engine := techStore.GetBestEngine(player)
	scanner := techStore.GetBestScanner(player)
	fuelTank := techStore.GetBestFuelTank(player)
	colonizationModule := techStore.GetBestColonizationModule(player)

	for i, hullSlot := range hull.Slots {
		slot := ShipDesignSlot{HullSlotIndex: i + 1}
		slot.Quantity = hullSlot.Capacity

		switch hullSlot.Type {
		case HullSlotTypeEngine:
			slot.HullComponent = engine.Name
		case HullSlotTypeScanner:
			slot.HullComponent = scanner.Name
		case HullSlotTypeMechanical:
			fallthrough
		case HullSlotTypeArmorScannerElectricalMechanical:
			switch purpose {
			case ShipDesignPurposeColonizer:
				fallthrough
			case ShipDesignPurposeColonistFreighter:
				if colonizationModule != nil {
					slot.HullComponent = colonizationModule.Name
					slot.Quantity = 1 // we only need 1 colonization module
				}
			}

		case HullSlotTypeGeneral:
			switch purpose {
			default:
				slot.HullComponent = fuelTank.Name
			}
		}

		// we filled it, add it
		if slot.HullComponent != "" {
			design.Slots = append(design.Slots, slot)
		}
	}

	return design
}
