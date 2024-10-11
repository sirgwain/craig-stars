package cs

import (
	"fmt"
	"math"

	"github.com/rs/zerolog"
)

// warpspeed for using a stargate vs moving with warp drive
const StargateWarpSpeed = 11

// time period to perform a task, like patrol
const Indefinite = 0

// use automatic warp speed for patrols
const PatrolWarpSpeedAutomatic = 0

// target fleets in any range when patrolling
const PatrolRangeInfinite = 0

// no target planet, player, etc
const None = 0

// Fleets are player owned collections of ships. Fleets have a slice of Tokens, one for each unique design
// in the fleet. Fleets also have orders that can be updated by the player, in the form of waypoints and the battle plan.
// Fleets are one of the commandable MapObjects in the game.
type Fleet struct {
	MapObject
	FleetOrders
	PlanetNum         int         `json:"planetNum"` // for starbase fleets that are owned by a planet
	BaseName          string      `json:"baseName"`
	Cargo             Cargo       `json:"cargo,omitempty"`
	Fuel              int         `json:"fuel"`
	Age               int         `json:"age"`
	Tokens            []ShipToken `json:"tokens"`
	Heading           Vector      `json:"heading,omitempty"`
	WarpSpeed         int         `json:"warpSpeed,omitempty"`
	PreviousPosition  *Vector     `json:"previousPosition,omitempty"`
	OrbitingPlanetNum int         `json:"orbitingPlanetNum,omitempty"`
	Starbase          bool        `json:"starbase,omitempty"`
	Spec              FleetSpec   `json:"spec,omitempty"`
	battlePlan        *BattlePlan
	struckMineField   bool
	remoteMined       bool
}

type FleetOrders struct {
	Waypoints     []Waypoint   `json:"waypoints"`
	RepeatOrders  bool         `json:"repeatOrders,omitempty"`
	BattlePlanNum int          `json:"battlePlanNum,omitempty"`
	Purpose       FleetPurpose `json:"purpose,omitempty"`
}

type FleetSpec struct {
	ShipDesignSpec
	BaseCloakedCargo int                        `json:"baseCloakedCargo"`
	BasePacketSpeed  int                        `json:"basePacketSpeed"`
	HasMassDriver    bool                       `json:"hasMassDriver,omitempty"`
	HasStargate      bool                       `json:"hasStargate,omitempty"`
	MassDriver       string                     `json:"massDriver,omitempty"`
	MassEmpty        int                        `json:"massEmpty"`
	MaxHullMass      int                        `json:"maxHullMass,omitempty"`
	MaxRange         int                        `json:"maxRange,omitempty"`
	Purposes         map[ShipDesignPurpose]bool `json:"purposes"`
	SafeHullMass     int                        `json:"safeHullMass,omitempty"`
	SafeRange        int                        `json:"safeRange,omitempty"`
	Stargate         string                     `json:"stargate,omitempty"`
	TotalShips       int                        `json:"totalShips"`
}

type Waypoint struct {
	Position             Vector                 `json:"position"`
	WarpSpeed            int                    `json:"warpSpeed"`
	EstFuelUsage         int                    `json:"estFuelUsage,omitempty"`
	Task                 WaypointTask           `json:"task"`
	TransportTasks       WaypointTransportTasks `json:"transportTasks,omitempty"`
	WaitAtWaypoint       bool                   `json:"waitAtWaypoint,omitempty"`
	LayMineFieldDuration int                    `json:"layMineFieldDuration,omitempty"`
	PatrolRange          int                    `json:"patrolRange,omitempty"`
	PatrolWarpSpeed      int                    `json:"patrolWarpSpeed,omitempty"`
	TargetType           MapObjectType          `json:"targetType,omitempty"`
	TargetNum            int                    `json:"targetNum,omitempty"`
	TargetPlayerNum      int                    `json:"targetPlayerNum,omitempty"`
	TargetName           string                 `json:"targetName,omitempty"`
	TransferToPlayer     int                    `json:"transferToPlayer,omitempty"`
	PartiallyComplete    bool                   `json:"partiallyComplete,omitempty"`
	processed            bool                   `json:"-"`
}

type WaypointTask string

const (
	WaypointTaskNone           = ""
	WaypointTaskTransport      = "Transport"
	WaypointTaskColonize       = "Colonize"
	WaypointTaskRemoteMining   = "RemoteMining"
	WaypointTaskMergeWithFleet = "MergeWithFleet"
	WaypointTaskScrapFleet     = "ScrapFleet"
	WaypointTaskLayMineField   = "LayMineField"
	WaypointTaskPatrol         = "Patrol"
	WaypointTaskRoute          = "Route"
	WaypointTaskTransferFleet  = "TransferFleet"
)

type WaypointTransportTasks struct {
	Fuel      WaypointTransportTask `json:"fuel"`
	Ironium   WaypointTransportTask `json:"ironium"`
	Boranium  WaypointTransportTask `json:"boranium"`
	Germanium WaypointTransportTask `json:"germanium"`
	Colonists WaypointTransportTask `json:"colonists"`
}

type WaypointTransportTask struct {
	Amount int                         `json:"amount,omitempty"`
	Action WaypointTaskTransportAction `json:"action,omitempty"`
}

type WaypointTaskTransportAction string

type transportTaskByType map[CargoType]WaypointTransportTask

const (
	// No transport task for the specified cargo.
	TransportActionNone WaypointTaskTransportAction = ""

	// (fuel only) Load or unload fuel until the fleet carries only the exact amount
	// needed to reach the next waypoint. You can use this task to send a fleet
	// loaded with fuel to rescue a stranded fleet. The rescue fleet will transfer
	// only the amount of fuel it can spare without stranding itself.
	TransportActionLoadOptimal WaypointTaskTransportAction = "LoadOptimal"

	// Load as much of the specified cargo as the fleet can hold.
	TransportActionLoadAll WaypointTaskTransportAction = "LoadAll"

	// Unload all the specified cargo at the waypoint.
	TransportActionUnloadAll WaypointTaskTransportAction = "UnloadAll"

	// Load the amount specified only if there is room in the hold.
	TransportActionLoadAmount WaypointTaskTransportAction = "LoadAmount"

	// Unload the amount specified only if the fleet is carrying that amount.
	TransportActionUnloadAmount WaypointTaskTransportAction = "UnloadAmount"

	// Loads up to the specified portion of the cargo hold subject to amount available at waypoint and room left in hold.
	TransportActionFillPercent WaypointTaskTransportAction = "FillPercent"

	// Remain at the waypoint until exactly X % of the hold is filled.
	TransportActionWaitForPercent WaypointTaskTransportAction = "WaitForPercent"

	// (minerals and colonists only) This command waits until all other loads and unloads are complete,
	// then loads as many colonists or amount of a mineral as will fit in the remaining space. For example,
	// setting Load All Germanium, Load Dunnage Ironium, will load all the Germanium that is available,
	// then as much Ironium as possible. If more than one dunnage cargo is specified, they are loaded in
	// the order of Ironium, Boranium, Germanium, and Colonists.
	TransportActionLoadDunnage WaypointTaskTransportAction = "LoadDunnage"

	// Load or unload the cargo until the amount on board is the amount specified.
	// If less than the specified cargo is available, the fleet will not move on.
	TransportActionSetAmountTo WaypointTaskTransportAction = "SetAmountTo"

	// Load or unload the cargo until the amount at the waypoint is the amount specified.
	// This order is always carried out to the best of the fleet’s ability that turn but does not prevent the fleet from moving on.
	TransportActionSetWaypointTo WaypointTaskTransportAction = "SetWaypointTo"
)

type FleetPurpose string

const (
	FleetPurposeNone              FleetPurpose = ""
	FleetPurposeScout             FleetPurpose = "Scout"
	FleetPurposeColonizer         FleetPurpose = "Colonizer"
	FleetPurposeBomber            FleetPurpose = "Bomber"
	FleetPurposeFighter           FleetPurpose = "Fighter"
	FleetPurposeCapitalShip       FleetPurpose = "CapitalShip"
	FleetPurposeFreighter         FleetPurpose = "Freighter"
	FleetPurposeColonistFreighter FleetPurpose = "ColonistFreighter"
	FleetPurposeArmedFreighter    FleetPurpose = "ArmedFreighter"
	FleetPurposeMineLayer         FleetPurpose = "MineLayer"
	FleetPurposeMiner             FleetPurpose = "Miner"
	FleetPurposeTerraformer       FleetPurpose = "Terraformer"
	FleetPurposeInvader           FleetPurpose = "Invader"
)

func FleetPurposeFromShipDesignPurpose(purpose ShipDesignPurpose) FleetPurpose {
	switch purpose {
	case ShipDesignPurposeScout:
		return FleetPurposeScout
	case ShipDesignPurposeColonizer:
		return FleetPurposeColonizer

	case ShipDesignPurposeBomber:
		fallthrough
	case ShipDesignPurposeSmartBomber:
		fallthrough
	case ShipDesignPurposeStructureBomber:
		return FleetPurposeBomber

	case ShipDesignPurposeFighter:
		fallthrough
	case ShipDesignPurposeFighterScout:
		return FleetPurposeFighter
	case ShipDesignPurposeCapitalShip:
		return FleetPurposeCapitalShip

	case ShipDesignPurposeFreighter:
		fallthrough
	case ShipDesignPurposeFuelFreighter:
		return FleetPurposeFreighter

	case ShipDesignPurposeMultiPurposeFreighter:
		fallthrough
	case ShipDesignPurposeColonistFreighter:
		return FleetPurposeColonistFreighter

	case ShipDesignPurposeArmedFreighter:
		return FleetPurposeArmedFreighter
	case ShipDesignPurposeMiner:
		return FleetPurposeMiner
	case ShipDesignPurposeTerraformer:
		return FleetPurposeTerraformer
	case ShipDesignPurposeDamageMineLayer:
		fallthrough
	case ShipDesignPurposeSpeedMineLayer:
		return FleetPurposeMineLayer
	}
	return FleetPurposeNone
}

type fleetMoveInterruptedReason int

const (
	fleetMoveInterruptedEngineFailure = iota
	fleetMoveInterruptedHitMineField
)

type fleetMoveInterrupted struct {
	reason    fleetMoveInterruptedReason
	mineField *MineField
}

func newFleet(player *Player, num int, name string, waypoints []Waypoint) Fleet {
	return Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
			Num:       num,
			Name:      fmt.Sprintf("%s #%d", name, num),
			Position:  waypoints[0].Position,
		},
		BaseName: name,
		Tokens:   []ShipToken{},
		FleetOrders: FleetOrders{
			Waypoints: waypoints,
		},
		OrbitingPlanetNum: None,
		battlePlan:        &player.BattlePlans[0],
	}
}

// create a new fleet with a design
func newFleetForDesign(player *Player, design *ShipDesign, quantity, num int, name string, waypoints []Waypoint) Fleet {
	return Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
			Num:       num,
			Name:      fmt.Sprintf("%s #%d", name, num),
			Position:  waypoints[0].Position,
		},
		BaseName: name,
		Tokens: []ShipToken{
			{design: design, DesignNum: design.Num, Quantity: quantity},
		},
		FleetOrders: FleetOrders{
			Waypoints: waypoints,
		},
		OrbitingPlanetNum: None,
		battlePlan:        &player.BattlePlans[0],
	}
}

func newFleetForToken(player *Player, num int, token ShipToken, waypoints []Waypoint) Fleet {
	return Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
			Num:       num,
			Name:      fmt.Sprintf("%s #%d", token.design.Name, num),
			Position:  waypoints[0].Position,
		},
		BaseName: token.design.Name,
		Tokens:   []ShipToken{token},
		FleetOrders: FleetOrders{
			Waypoints: waypoints,
		},
		OrbitingPlanetNum: None,
	}
}

// create a new fleet that is a starbase
func newStarbase(player *Player, planet *Planet, design *ShipDesign, name string) Fleet {
	fleet := newFleetForDesign(player, design, 1, 0, name, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, 1)})
	fleet.PlanetNum = planet.Num
	fleet.Starbase = true

	return fleet
}

func (f *Fleet) String() string {
	return fmt.Sprintf("Fleet %s #%d", f.BaseName, f.Num)
}

func (f *Fleet) withPlayerNum(playerNum int) *Fleet {
	f.PlayerNum = playerNum
	return f
}

func (f *Fleet) withNum(num int) *Fleet {
	f.Num = num
	return f
}

func (f *Fleet) withCargo(cargo Cargo) *Fleet {
	f.Cargo = cargo
	return f
}

func (f *Fleet) withFuel(fuel int) *Fleet {
	f.Fuel = fuel
	return f
}

func (f *Fleet) withPosition(position Vector) *Fleet {
	f.Position = position
	// todo: should we set waypoints in a builder?
	f.Waypoints = []Waypoint{{Position: position}}
	return f
}

func (f *Fleet) withWaypoints(waypoints ...Waypoint) *Fleet {
	f.Waypoints = waypoints
	return f
}

func (f *Fleet) withOrbitingPlanetNum(num int) *Fleet {
	f.OrbitingPlanetNum = num
	return f
}

// get a pointer to a ShipToken a design, or nil if it's not present
func (f *Fleet) getTokenByDesign(designNum int) *ShipToken {
	for i, token := range f.Tokens {
		if token.DesignNum == designNum {
			return &f.Tokens[i]
		}
	}
	return nil
}

func (f *Fleet) Orbiting() bool {
	return f.OrbitingPlanetNum != None
}

func (f *Fleet) Idle() bool {
	return len(f.Waypoints) == 1 && f.Waypoints[0].Task == WaypointTaskNone
}

func (f *Fleet) Rename(name string) {
	f.BaseName = name
	f.Name = fmt.Sprintf("%s #%d", f.BaseName, f.Num)
}

func NewPlanetWaypoint(position Vector, num int, name string, warpSpeed int) Waypoint {
	return Waypoint{
		Position:        position,
		TargetType:      MapObjectTypePlanet,
		TargetNum:       num,
		TargetName:      name,
		TargetPlayerNum: None,
		WarpSpeed:       warpSpeed,
	}
}

func NewFleetWaypoint(position Vector, num int, playerNum int, name string, warpSpeed int) Waypoint {
	return Waypoint{
		Position:        position,
		TargetType:      MapObjectTypeFleet,
		TargetNum:       num,
		TargetPlayerNum: playerNum,
		TargetName:      name,
		WarpSpeed:       warpSpeed,
	}
}

func NewMysteryTraderWaypoint(mt *MysteryTrader, warpSpeed int) Waypoint {
	return Waypoint{
		Position:   mt.Position,
		TargetType: mt.Type,
		TargetNum:  mt.Num,
		TargetName: mt.Name,
		WarpSpeed:  warpSpeed,
	}
}

func NewPositionWaypoint(position Vector, warpSpeed int) Waypoint {
	return Waypoint{
		Position:        position,
		WarpSpeed:       warpSpeed,
		TargetNum:       None,
		TargetPlayerNum: None,
	}
}

func (wp *Waypoint) clearTarget() {
	wp.TargetName = ""
	wp.TargetNum = None
	wp.TargetPlayerNum = None
	wp.TargetType = MapObjectTypeNone
}

func (wp *Waypoint) targetPlanet(planet *Planet) {
	wp.TargetType = MapObjectTypePlanet
	wp.TargetName = planet.Name
	wp.TargetNum = planet.Num
}

func (wp Waypoint) WithTask(task WaypointTask) Waypoint {
	wp.Task = task
	return wp
}

func (wp Waypoint) WithTransportTasks(transportTasks WaypointTransportTasks) Waypoint {
	wp.TransportTasks = transportTasks
	return wp
}

// get a list of transport tasks keyed by cargotype
func (wp Waypoint) getTransportTasks() transportTaskByType {
	tasks := transportTaskByType{}
	if wp.TransportTasks.Fuel.Action != TransportActionNone {
		tasks[Fuel] = wp.TransportTasks.Fuel
	}
	if wp.TransportTasks.Ironium.Action != TransportActionNone {
		tasks[Ironium] = wp.TransportTasks.Ironium
	}
	if wp.TransportTasks.Boranium.Action != TransportActionNone {
		tasks[Boranium] = wp.TransportTasks.Boranium
	}
	if wp.TransportTasks.Germanium.Action != TransportActionNone {
		tasks[Germanium] = wp.TransportTasks.Germanium
	}
	if wp.TransportTasks.Colonists.Action != TransportActionNone {
		tasks[Colonists] = wp.TransportTasks.Colonists
	}

	return tasks
}

// inject designs into tokens so all the various Compute* functions work
func (f *Fleet) InjectDesigns(designs []*ShipDesign) error {

	designsByNum := make(map[int]*ShipDesign, len(designs))
	for i := range designs {
		design := designs[i]
		designsByNum[design.Num] = design
	}

	// inject the design into this
	for i := range f.Tokens {
		token := &f.Tokens[i]
		token.design = designsByNum[token.DesignNum]
		if token.design == nil {
			return fmt.Errorf("unable to find design %d for fleet %s", token.DesignNum, f.Name)
		}
	}

	return nil
}

// compute all the computable values of this fleet (cargo capacity, armor, mass)
func ComputeFleetSpec(rules *Rules, player *Player, fleet *Fleet) FleetSpec {
	spec := FleetSpec{
		ShipDesignSpec: ShipDesignSpec{
			ScanRangePen:   NoScanner,
			SpaceDock:      UnlimitedSpaceDock,
			ReduceCloaking: 1,
		},
		Purposes: map[ShipDesignPurpose]bool{},
	}
	spec.Mass = fleet.Cargo.Total()

	for _, token := range fleet.Tokens {

		// update our total ship count
		spec.TotalShips += token.Quantity

		if token.design.Purpose != ShipDesignPurposeNone {
			spec.Purposes[token.design.Purpose] = true
		}

		if token.design.Spec.Starbase {
			spec.Starbase = true
		}
		spec.MaxPopulation = MaxInt(spec.MaxPopulation, token.design.Spec.MaxPopulation)
		spec.InnateScanRangePenFactor = math.Max(spec.InnateScanRangePenFactor, token.design.Spec.InnateScanRangePenFactor) // Ultra Station and Death Stars have pen scanning

		// use the lowest ideal speed for this fleet
		// if we have multiple engines
		if (token.design.Spec.Engine != Engine{}) {
			if spec.Engine.IdealSpeed == 0 {
				spec.Engine.IdealSpeed = token.design.Spec.Engine.IdealSpeed
				spec.Engine.FreeSpeed = token.design.Spec.Engine.FreeSpeed
				spec.Engine.MaxSafeSpeed = token.design.Spec.Engine.MaxSafeSpeed
			} else {
				spec.Engine.IdealSpeed = MinInt(spec.Engine.IdealSpeed, token.design.Spec.Engine.IdealSpeed)
				spec.Engine.FreeSpeed = token.design.Spec.Engine.FreeSpeed
				spec.Engine.MaxSafeSpeed = MinInt(spec.Engine.MaxSafeSpeed, token.design.Spec.Engine.MaxSafeSpeed)
			}
		}
		// cost
		spec.Cost = token.design.Spec.Cost.MultiplyInt(token.Quantity)

		// mass
		spec.Mass += token.design.Spec.Mass * token.Quantity
		spec.MassEmpty += token.design.Spec.Mass * token.Quantity

		// armor
		spec.Armor += token.design.Spec.Armor * token.Quantity

		// shield
		spec.Shields += token.design.Spec.Shields * token.Quantity

		// cargo
		spec.CargoCapacity += token.design.Spec.CargoCapacity * token.Quantity

		// fuel
		spec.FuelCapacity += token.design.Spec.FuelCapacity * token.Quantity
		spec.FuelGeneration += token.design.Spec.FuelGeneration * token.Quantity

		// minesweep
		spec.MineSweep += token.design.Spec.MineSweep * token.Quantity

		// remote mining
		spec.MiningRate += token.design.Spec.MiningRate * token.Quantity

		// remote terraforming
		spec.TerraformRate += token.design.Spec.TerraformRate * token.Quantity

		// colonization
		spec.Colonizer = spec.Colonizer || token.design.Spec.Colonizer
		spec.OrbitalConstructionModule = spec.OrbitalConstructionModule || token.design.Spec.OrbitalConstructionModule

		// radiating parts
		spec.Radiating = spec.Radiating || token.design.Spec.Radiating

		// spec all mine layers in the fleet
		if token.design.Spec.CanLayMines {
			spec.CanLayMines = true
			if spec.MineLayingRateByMineType == nil {
				spec.MineLayingRateByMineType = make(map[MineFieldType]int)
			}
			for key := range token.design.Spec.MineLayingRateByMineType {
				if _, ok := spec.MineLayingRateByMineType[key]; ok {
					spec.MineLayingRateByMineType[key] = 0
				}
				spec.MineLayingRateByMineType[key] += token.design.Spec.MineLayingRateByMineType[key] * token.Quantity
			}
		}

		// We should only have one ship stack with spacdock capabilities, but for this logic just go with the max
		spec.SpaceDock = MaxInt(spec.SpaceDock, token.design.Spec.SpaceDock)

		// sadly, the fleet only gets the best repair bonus from one design
		spec.RepairBonus = math.Max(spec.RepairBonus, token.design.Spec.RepairBonus)

		spec.ScanRange = MaxInt(spec.ScanRange, token.design.Spec.ScanRange)
		spec.ScanRangePen = MaxInt(spec.ScanRangePen, token.design.Spec.ScanRangePen)
		if token.design.Spec.Scanner {
			spec.Scanner = true
		}

		// add bombs
		if token.design.Spec.Bomber {
			spec.Bomber = true
			for _, bomb := range token.design.Spec.Bombs {
				bomb.Quantity *= token.Quantity
				spec.Bombs = append(spec.Bombs, bomb)
			}
			for _, bomb := range token.design.Spec.SmartBombs {
				bomb.Quantity *= token.Quantity
				spec.SmartBombs = append(spec.SmartBombs, bomb)
			}
			for _, bomb := range token.design.Spec.RetroBombs {
				bomb.Quantity *= token.Quantity
				spec.RetroBombs = append(spec.RetroBombs, bomb)
			}
		}

		// check if any tokens have weapons
		// we process weapon slots per stack, so we don't need to spec all
		// weapons in a fleet
		if token.design.Spec.HasWeapons {
			spec.HasWeapons = true
		}

		if token.design.Spec.CloakUnits > 0 {
			// calculate the cloak units for this token based on the design's cloak units (i.e. 70 cloak units / kT for a stealh cloak)
			spec.CloakUnits += token.design.Spec.CloakUnits
		} else {
			// if this ship doesn't have cloaking, it counts as cargo (except for races with free cargo cloaking)
			if !player.Race.Spec.FreeCargoCloaking {
				spec.BaseCloakedCargo += token.design.Spec.Mass * token.Quantity
			}
		}

		// choose the best tachyon detector ship
		spec.ReduceCloaking = math.Min(spec.ReduceCloaking, token.design.Spec.ReduceCloaking)

		spec.CanJump = spec.CanJump || token.design.Spec.CanJump
		spec.CanStealFleetCargo = spec.CanStealFleetCargo || token.design.Spec.CanStealFleetCargo
		spec.CanStealPlanetCargo = spec.CanStealPlanetCargo || token.design.Spec.CanStealPlanetCargo

		// stargate fields
		if token.design.Spec.SafeHullMass != 0 {
			spec.Stargate = token.design.Spec.Stargate
			spec.HasStargate = true
			spec.SafeHullMass = token.design.Spec.SafeHullMass
		}
		if token.design.Spec.MaxHullMass != 0 {
			spec.MaxHullMass = token.design.Spec.MaxHullMass
		}
		if token.design.Spec.SafeRange != 0 {
			spec.SafeRange = token.design.Spec.SafeRange
		}
		if token.design.Spec.MaxRange != 0 {
			spec.MaxRange = token.design.Spec.MaxRange
		}

		if token.design.Spec.SafePacketSpeed != 0 {
			spec.HasMassDriver = true
			spec.MassDriver = token.design.Spec.MassDriver
			spec.SafePacketSpeed = token.design.Spec.SafePacketSpeed
			spec.BasePacketSpeed = token.design.Spec.BasePacketSpeed
			spec.AdditionalMassDrivers = token.design.Spec.AdditionalMassDrivers
		}

	}

	// compute the cloaking based on the cloak units and cargo
	spec.CloakPercent = computeFleetCloakPercent(&spec, fleet.Cargo.Total() + spec.BaseCloakedCargo, player.Race.Spec.FreeCargoCloaking)

	if !spec.Starbase {
		spec.EstimatedRange = fleet.getEstimatedRange(player, spec.Engine.IdealSpeed, spec.CargoCapacity)
	}

	return spec
}

// compute fuel usage for each waypoint
func (f *Fleet) computeFuelUsage(player *Player) {
	for i := range f.Waypoints {
		wp := &f.Waypoints[i]
		if i > 0 && wp.WarpSpeed < StargateWarpSpeed {
			wpPrevious := f.Waypoints[i-1]
			fuelUsage := f.GetFuelCost(player, wp.WarpSpeed, wp.Position.DistanceTo(wpPrevious.Position))
			wp.EstFuelUsage = fuelUsage
		} else {
			wp.EstFuelUsage = 0
		}
	}
}

// compute a fleet's cloak percent based on its current cargo/mass/cloak units
func computeFleetCloakPercent(spec *FleetSpec, cargoTotal int, freeCargoCloaking bool) int {
	cloakUnits := spec.CloakUnits

	// starbases have no mass or cargo, but fleet cloaking is adjusted for it
	if spec.Mass > 0 {
		// figure out how much cargo we are cloaking
		cloakedCargo := 0
		if !freeCargoCloaking {
			cloakedCargo = cargoTotal
		}

		cloakUnits = int(math.Round(float64(cloakUnits) * float64(spec.MassEmpty) / float64(spec.MassEmpty+cloakedCargo)))
	}
	return getCloakPercentForCloakUnits(cloakUnits)
}

// make sure we don't overflow our fuel. After a battle, we might have more fuel than our fleet can hold
func (fleet *Fleet) reduceFuelToMax() {
	fleet.Fuel = MinInt(fleet.Spec.FuelCapacity, fleet.Fuel)
}

// make sure we don't overflow our cargo. After a battle, we might have more cargo than our fleet can hold
// return any dropped cargo
func (fleet *Fleet) reduceCargoToMax() Cargo {
	capacity := fleet.Spec.CargoCapacity
	cargo := fleet.Cargo

	// no capacity, no cargo
	if capacity == 0 {
		fleet.Cargo = Cargo{}
		return cargo
	}

	// check if we have more cargo than we can hold
	total := fleet.Cargo.Total()
	if total > capacity {

		// save the people first!
		if fleet.Cargo.Colonists > 0 {
			fleet.Cargo.Colonists = MinInt(fleet.Cargo.Colonists, capacity)
		}

		// if we have 110kT of space and 10kT is taken up by colonists, we have 100kT remaining capacity
		// if we have 200kT of minerals left, we keep half of each
		minerals := fleet.Cargo.ToMineral()
		remainingCapacity := MaxInt(0, capacity-fleet.Cargo.Colonists)

		// if we have no capacity left, drop all minerals and
		if remainingCapacity == 0 {
			fleet.Cargo = Cargo{Colonists: fleet.Cargo.Colonists}
			return cargo.Subtract(fleet.Cargo)
		}
		totalMinerals := minerals.Total()

		// reduce each mineral by a percent
		percentToKeep := 1 / (float64(totalMinerals) / float64(remainingCapacity))
		minerals = minerals.MultiplyFloat64(percentToKeep)
		fleet.Cargo = Cargo{
			minerals.Ironium,
			minerals.Boranium,
			minerals.Germanium,
			fleet.Cargo.Colonists,
		}
	}

	return cargo.Subtract(fleet.Cargo)
}

// return true if this fleet would attack another player's fleet, planet, minefield, etc.
func (fleet *Fleet) willAttack(fleetPlayer *Player, otherPlayerNum int) bool {
	switch fleet.battlePlan.AttackWho {
	case BattleAttackWhoEnemies:
		return fleetPlayer.IsEnemy(otherPlayerNum)
	case BattleAttackWhoEnemiesAndNeutrals:
		return fleetPlayer.IsEnemy(otherPlayerNum) || fleetPlayer.IsNeutral(otherPlayerNum)
	case BattleAttackWhoEveryone:
		return true
	default:
		return false
	}

}

func (f *Fleet) availableCargoSpace() int {
	return Clamp(f.Spec.CargoCapacity-f.Cargo.Total(), 0, f.Spec.CargoCapacity)
}

func (f *Fleet) availableFuelSpace() int {
	return Clamp(f.Spec.FuelCapacity-f.Fuel, 0, f.Spec.FuelCapacity)
}

// transfer cargo from a fleet to a cargo holder
func (f *Fleet) transferToDest(dest cargoHolder, cargoType CargoType, transferAmount int) error {
	destCargo := dest.getCargo()

	if transferAmount > 0 && !f.Cargo.CanTransferAmount(cargoType, transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %d to %s, there is not enough in the fleet to transfer", f.Name, transferAmount, dest.getMapObject().Name)
	}

	if transferAmount < 0 && f.availableCargoSpace() < -transferAmount {
		return fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", f.Name, f.availableCargoSpace(), transferAmount, dest.getMapObject().Name)
	}

	if transferAmount < 0 && !destCargo.CanTransferAmount(cargoType, -transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %d from %s, there is not enough to transfer", f.Name, transferAmount, dest.getMapObject().Name)
	}

	if transferAmount > 0 && dest.getCargoCapacity() != Unlimited && (dest.getCargoCapacity()-destCargo.Total()) < transferAmount {
		return fmt.Errorf("fleet %s cannot transfer %d to %s, there is not enough to space to hold the cargo", f.Name, transferAmount, dest.getMapObject().Name)
	}

	// transfer the cargo
	f.Cargo.SubtractAmount(cargoType, transferAmount)
	destCargo.AddAmount(cargoType, transferAmount)

	return nil
}

// remove any empty tokens that were destroyed (by minefields, overgating, battle... it's a dangerous universe)
func (fleet *Fleet) removeEmptyTokens() {
	updatedTokens := make([]ShipToken, 0, len(fleet.Tokens))
	for _, token := range fleet.Tokens {
		// keep this token
		if token.Quantity > 0 {
			updatedTokens = append(updatedTokens, token)
		}
	}
	fleet.Tokens = updatedTokens
}

// move a fleet through space, check for minefields, use fuel, etc
func (fleet *Fleet) moveFleet(rules *Rules, mapObjectGetter mapObjectGetter, playerGetter playerGetter) (interrupted *fleetMoveInterrupted) {
	player := playerGetter.getPlayer(fleet.PlayerNum)
	wp0 := fleet.Waypoints[0]
	wp1 := fleet.Waypoints[1]
	totalDist := fleet.Position.DistanceTo(wp1.Position)
	fleet.PreviousPosition = &Vector{fleet.Position.X, fleet.Position.Y}

	dist := float64(wp1.WarpSpeed * wp1.WarpSpeed)
	// round up, if we are <1 away, i.e. the target is 81.9 ly away, warp 9 (81 ly travel) should be able to make it there
	if dist < totalDist && totalDist-dist < 1 {
		dist = math.Ceil(totalDist)
	}

	// make sure we end up at a whole number
	vectorTravelled := wp1.Position.Subtract(fleet.Position).Normalized().Scale(dist)
	dist = vectorTravelled.Length()
	// don't overshoot
	dist = math.Min(totalDist, dist)

	// check for CE engine failure
	if player.Race.Spec.EngineFailureRate > 0 && wp1.WarpSpeed > player.Race.Spec.EngineReliableSpeed && rules.random.Float64() <= player.Race.Spec.EngineFailureRate {
		messager.fleetEngineFailure(player, fleet)
		return &fleetMoveInterrupted{reason: fleetMoveInterruptedEngineFailure}
	}

	// get the cost for the fleet
	fuelCost := fleet.GetFuelCost(player, wp1.WarpSpeed, dist)
	var fuelGenerated int = 0
	if fuelCost > fleet.Fuel {
		// we will run out of fuel
		// if this distance would have cost us 10 fuel but we have 6 left, only travel 60% of the distance.
		distanceFactor := float64(fleet.Fuel) / float64(fuelCost)
		dist = dist * distanceFactor
		fuelCost = fleet.Fuel

		// collide with minefields on route, but don't hit a minefield if we run out of fuel beforehand
		hitMineField, actualDist := checkForMineFieldCollision(rules, playerGetter, mapObjectGetter, fleet, wp1, dist)
		if hitMineField != nil {
			interrupted = &fleetMoveInterrupted{reason: fleetMoveInterruptedHitMineField, mineField: hitMineField}
		}

		// we hit a minefield before we ran out of fuel
		if actualDist != dist {
			dist = actualDist
			fuelCost = fleet.GetFuelCost(player, wp1.WarpSpeed, dist)
		} else {
			wp1.WarpSpeed = fleet.Spec.Engine.FreeSpeed
			fleet.Waypoints[1] = wp1
			messager.fleetOutOfFuel(player, fleet, wp1.WarpSpeed)
			// if we ran out of fuel 60% of the way to our normal distance, the remaining 40% of our time
			// was spent travelling at fuel generation speeds:
			remainingDistanceTravelled := (1 - distanceFactor) * float64(wp1.WarpSpeed*wp1.WarpSpeed)
			dist += remainingDistanceTravelled
			fuelGenerated = fleet.getFuelGeneration(wp1.WarpSpeed, remainingDistanceTravelled)
		}

		fleet.Fuel -= fuelCost
	} else {
		// collide with minefields on route, but don't hit a minefield if we run out of fuel beforehand
		hitMineField, actualDist := checkForMineFieldCollision(rules, playerGetter, mapObjectGetter, fleet, wp1, dist)
		if hitMineField != nil {
			interrupted = &fleetMoveInterrupted{reason: fleetMoveInterruptedHitMineField, mineField: hitMineField}
		}

		if actualDist != dist {
			dist = actualDist
			fuelCost = fleet.GetFuelCost(player, wp1.WarpSpeed, dist)
		}

		fleet.Fuel -= fuelCost
		fuelGenerated = fleet.getFuelGeneration(wp1.WarpSpeed, dist)
	}

	// message the player about fuel generation
	fuelGenerated = MinInt(fuelGenerated, fleet.Spec.FuelCapacity-fleet.Fuel)
	if fuelGenerated > 0 {
		fleet.Fuel += fuelGenerated
		messager.fleetGeneratedFuel(player, fleet, fuelGenerated)
	}

	// assuming we move at all, make sure we are no longer orbiting any planets
	if dist > 0 && fleet.Orbiting() {
		fleet.OrbitingPlanetNum = None
	}

	if totalDist == dist {
		fleet.completeMove(mapObjectGetter, player, wp0, wp1)
	} else {
		// update what other people see for this fleet's speed and direction
		fleet.WarpSpeed = wp1.WarpSpeed
		fleet.Heading = (wp1.Position.Subtract(fleet.Position)).Normalized()

		// move this fleet closer to the next waypoint
		wp0.TargetType = MapObjectTypeNone
		wp0.TargetNum = None
		wp0.TargetPlayerNum = None
		wp0.TargetName = ""
		wp0.PartiallyComplete = true

		fleet.Position = fleet.Position.Add(fleet.Heading.Scale(dist))
		fleet.Position = fleet.Position.Round()
		wp0.Position = fleet.Position

		if fleet.struckMineField {
			fleet.WarpSpeed = 0
			fleet.Heading = Vector{}
		}

		// don't do any transport in mid space, reset this
		if wp0.Task == WaypointTaskTransport {
			wp0.Task = WaypointTaskNone
			wp0.TransportTasks = WaypointTransportTasks{}
		}

		fleet.Waypoints[0] = wp0
	}

	// if we ended up at a planet, make sure we are orbiting it
	if fleet.OrbitingPlanetNum == None {
		for _, mo := range mapObjectGetter.getMapObjectsAtPosition(fleet.Position) {
			if planet, ok := mo.(*Planet); ok {
				fleet.OrbitingPlanetNum = planet.Num
			}
		}
	}
	return interrupted
}

// GateFleet moves the fleet the cool way, with stargates!
func (fleet *Fleet) gateFleet(rules *Rules, mapObjectGetter mapObjectGetter, playerGetter playerGetter) {
	player := playerGetter.getPlayer(fleet.PlayerNum)
	wp0 := fleet.Waypoints[0]
	wp1 := fleet.Waypoints[1]
	totalDist := fleet.Position.DistanceTo(wp1.Position)
	fleet.PreviousPosition = &Vector{fleet.Position.X, fleet.Position.Y}

	var sourceStargate, destStargate PlanetStarbaseSpec

	// if we got here, both source and dest have stargates (unless we're using a jumpgate)
	destPlanet := mapObjectGetter.getPlanet(wp1.TargetNum)

	if destPlanet == nil || !destPlanet.Spec.HasStargate {
		messager.fleetStargateInvalidDest(player, fleet, wp0, wp1)
		return
	}

	destPlanetPlayer := playerGetter.getPlayer(destPlanet.PlayerNum)

	if !destPlanetPlayer.IsFriend(player.Num) {
		messager.fleetStargateInvalidDestOwner(player, fleet, wp0, wp1)
		return
	}

	destStargate = destPlanet.Spec.PlanetStarbaseSpec

	// jumpgate fleets don't use the source planet, only the dest
	var sourcePlanet *Planet
	if fleet.Spec.CanJump {
		sourceStargate = destStargate
	} else {
		sourcePlanet = mapObjectGetter.getPlanet(fleet.OrbitingPlanetNum)

		if sourcePlanet == nil || !sourcePlanet.Spec.HasStargate {
			messager.fleetStargateInvalidSource(player, fleet, wp0)
			return
		}

		sourcePlanetPlayer := playerGetter.getPlayer(sourcePlanet.PlayerNum)
		if sourcePlanetPlayer != nil && !sourcePlanetPlayer.IsFriend(player.Num) {
			messager.fleetStargateInvalidSourceOwner(player, fleet, wp0, wp1)
			return
		}

		sourceStargate = sourcePlanet.Spec.PlanetStarbaseSpec
	}

	// can't gate colonists unless we're IT
	// can't dump colonists into space or on a world we don't own
	// ships with jump gates can gate cargo (amazing)
	if !fleet.Spec.CanJump && !player.Race.Spec.CanGateCargo && fleet.Cargo.Colonists > 0 && (sourcePlanet == nil || !sourcePlanet.OwnedBy(player.Num)) {
		messager.fleetStargateInvalidColonists(player, fleet, wp0, wp1)
		return
	}

	// only the source gate matters for range
	minSafeRange := sourceStargate.SafeRange
	minSafeHullMass := MinInt(sourceStargate.SafeHullMass, destStargate.SafeHullMass)

	// check if we are exceeding the max distance
	if totalDist > float64(minSafeRange*rules.StargateMaxRangeFactor) {
		messager.fleetStargateInvalidRange(player, fleet, wp0, wp1, totalDist)
		return
	}

	// check if any ships exceed the max mass allowed
	for _, token := range fleet.Tokens {
		if token.design.Spec.Mass > minSafeHullMass*rules.StargateMaxHullMassFactor {
			messager.fleetStargateInvalidMass(player, fleet, wp0, wp1)
			return
		}
	}

	// dump cargo if we aren't IT or using a jump gate
	if !fleet.Spec.CanJump && fleet.Cargo.Total() > 0 && !player.Race.Spec.CanGateCargo {
		messager.fleetStargateDumpedCargo(player, fleet, wp0, wp1, fleet.Cargo)
		sourcePlanet.Cargo = sourcePlanet.Cargo.Add(fleet.Cargo)
		fleet.Cargo = Cargo{}
	}

	// apply overgate damage and delete tokens (and possibly the fleet)
	// also vanish tokens for non IT races
	fleet.applyOvergatePenalty(player, rules, totalDist, wp0, wp1, sourceStargate, destStargate)

	// if the fleet is gone, we're done
	if len(fleet.Tokens) == 0 {
		return
	}

	// we survived, warp it!
	fleet.completeMove(mapObjectGetter, player, wp0, wp1)
}

// if the fleet went over safe warp, explode some ships
func (fleet *Fleet) applyOverwarpPenalty(rules *Rules) int {
	if len(fleet.Waypoints) <= 1 {
		return 0
	}
	wp1 := &fleet.Waypoints[1]
	// check for exploded ships
	explodedShips := 0
	for tokenIndex := range fleet.Tokens {
		token := &fleet.Tokens[tokenIndex]
		if wp1.WarpSpeed > token.design.Spec.Engine.MaxSafeSpeed && wp1.WarpSpeed != StargateWarpSpeed {
			// explode some fleets if you go too fast
			for shipIndex := 0; shipIndex < token.Quantity; shipIndex++ {
				if rules.FleetSafeSpeedExplosionChance > rules.random.Float64() {
					explodedShips++
					token.Quantity--
				}
			}
		}
	}

	return explodedShips
}

// applyOvergatePenalty applies damage (if any) to each token that overgated
func (fleet *Fleet) applyOvergatePenalty(player *Player, rules *Rules, distance float64, wp0, wp1 Waypoint, sourceStargate, destStargate PlanetStarbaseSpec) {
	var totalDamage, shipsLostToDamage, shipsLostToTheVoid, startingShips int
	for i := range fleet.Tokens {
		token := &fleet.Tokens[i]
		startingShips += token.Quantity
		// Inner stellar travellers never lose ships to the void, but everyone else does
		if player.Race.Spec.ShipsVanishInVoid {
			rangeVanishChance := token.getStargateRangeVanishingChance(distance, sourceStargate.SafeRange)
			massVanishingChance := token.getStargateMassVanishingChance(sourceStargate.SafeHullMass, rules.StargateMaxHullMassFactor)
			// Combined vanishing chance idea courtesy of ekolis
			vanishingChance := 1 - (1-rangeVanishChance)*(1-massVanishingChance)

			if rangeVanishChance > 0 || massVanishingChance > 0 {
				for i := 0; i < token.Quantity; i++ {
					// check if it vanishes due to range, if not, check if it vanishes due
					// to mass. Each ship can only vanish once
					if vanishingChance > rules.random.Float64() {
						// oh no, we lost a ship!
						shipsLostToTheVoid++
						token.Quantity--
						i--
						if token.QuantityDamaged > 0 {
							// get rid of the damaged ships first
							// if we're out of damaged ships, reset our
							// token damage to 0
							token.QuantityDamaged--
							// can't have damage without damaged ships
							// I don't think this should ever come up
							if token.QuantityDamaged == 0 {
								token.Damage = 0
							}
						}
					}
				}
			}
		}

		// if we didn't lose tokens in
		if token.Quantity > 0 {
			tokenDamage := token.applyOvergateDamage(distance, sourceStargate.SafeRange, sourceStargate.SafeHullMass, destStargate.SafeHullMass, rules.StargateMaxHullMassFactor)

			totalDamage += tokenDamage.damage
			shipsLostToDamage += tokenDamage.shipsDestroyed
		}
	}

	// remove any tokens that were lost completely
	fleet.removeEmptyTokens()

	if len(fleet.Tokens) == 0 {
		messager.fleetStargateDestroyed(player, fleet, wp0, wp1)
	} else {
		if totalDamage > 0 || shipsLostToTheVoid > 0 {
			messager.fleetStargateDamaged(player, fleet, wp0, wp1, totalDamage, startingShips, shipsLostToDamage, shipsLostToTheVoid)
		}
	}
}

// Engine fuel usage calculation courtesy of m.a@stars
func (engine Engine) getFuelCostForEngine(warpSpeed int, mass int, dist float64, ifeFactor float64) int {
	if warpSpeed == 0 {
		return 0
	}
	// 1 mg of fuel will move 200kT of weight 1 LY at a Fuel Usage Number of 100.
	// Number of engines doesn't matter. Neither number of ships with the same engine.

	distanceCeiling := math.Ceil(dist) // rounding to next integer gives best graph fit

	// IFE is applied to drive specifications, just as the helpfile hints.
	// Stars! probably does it outside here once per turn per engine to save time.
	engineEfficiency := math.Ceil(ifeFactor * float64(engine.FuelUsage[warpSpeed]))

	// 20000 = 200*100
	// Safe bet is Stars! does all this with integer math tricks.
	// Subtracting 2000 in a loop would be a way to also get the rounding.
	// Or even bitshift for the 2 and adjust "decimal point" for the 1000
	teorFuel := (math.Floor(float64(mass)*engineEfficiency*distanceCeiling/2000) / 10)
	// using only one decimal introduces another artifact: .0999 gets rounded down to .0

	// The heavier ships will benefit the most from the accuracy
	intFuel := int(math.Ceil(teorFuel))

	// That's all. Nothing really fancy, much less random. Subtle differences in
	// math lib workings might explain the rarer and smaller discrepancies observed
	return intFuel
	// Unrelated to this fuel math are some quirks inside the
	// "negative fuel" watchdog when the remainder of the
	// trip is < 1 ly. Aahh, the joys of rounding! ;o)
}

// Get the Fuel cost for this fleet to travel a certain distance at a certain speed
func (fleet *Fleet) GetFuelCost(player *Player, warpSpeed int, distance float64) int {
	return fleet.getFuelCost(player, warpSpeed, distance, fleet.Spec.CargoCapacity)
}

func (fleet *Fleet) getFuelCost(player *Player, warpSpeed int, distance float64, cargoCapacity int) int {

	// figure out how much fuel we're going to use
	efficiencyFactor := 1 + player.Race.Spec.FuelEfficiencyOffset

	var fuelCost int = 0

	// compute each ship stack separately
	for _, token := range fleet.Tokens {
		// figure out this ship stack's mass as well as its proportion of the cargo
		mass := token.design.Spec.Mass * token.Quantity
		fleetCargo := fleet.Cargo.Total()
		stackCapacity := token.design.Spec.CargoCapacity * token.Quantity

		if cargoCapacity > 0 {
			mass += int(float64(fleetCargo) * (float64(stackCapacity) / float64(cargoCapacity)))
		}

		engine := token.design.Spec.Engine
		fuelCost += engine.getFuelCostForEngine(warpSpeed, mass, distance, efficiencyFactor)
	}

	return fuelCost
}

func (fleet *Fleet) getEstimatedRange(player *Player, warpSpeed int, cargoCapacity int) int {
	fuelCost := fleet.getFuelCost(player, warpSpeed, 1000, cargoCapacity)
	if fuelCost == 0 {
		return Infinite
	}
	return int(float64(fleet.Fuel) / float64(fuelCost) * 1000)
}

// Get the amount of fuel this ship will generate at a given warp
// F = 0 if the engine is running above the highest warp at which it travels for free (i.e. it is using fuel)
// F = D if the engine is running at the highest warp at which it travels for free
// F = 3D if the engine is running 1 warp factor below the highest warp at which it travels for free
// F = 6D if the engine is running 2 warp factors below the highest warp at which it travels for free
// F = 10D if the engine is running 3 or more warp factors below the highest warp at which it travels for free
// Note that the fuel generated is per engine, not per ship; i.e.; a ship with 2, 3, or 4 engines
// produces (or uses) 2, 3, or 4 times as much fuel as a single engine ship.
func (fleet *Fleet) getFuelGeneration(warpSpeed int, distance float64) int {
	fuelGenerated := 0.0
	for _, token := range fleet.Tokens {
		freeSpeed := token.design.Spec.Engine.FreeSpeed
		numEngines := token.design.Spec.NumEngines * token.Quantity
		speedDifference := freeSpeed - warpSpeed
		if speedDifference == 0 {
			fuelGenerated += distance * float64(numEngines)
		} else if speedDifference == 1 {
			fuelGenerated += (3 * distance) * float64(numEngines)
		} else if speedDifference == 2 {
			fuelGenerated += (6 * distance) * float64(numEngines)
		} else if speedDifference >= 3 {
			fuelGenerated += (10 * distance) * float64(numEngines)
		}
	}

	return int(fuelGenerated)
}

// Complete a move from one waypoint to another
func (fleet *Fleet) completeMove(mapObjectGetter mapObjectGetter, player *Player, wp0 Waypoint, wp1 Waypoint) {
	fleet.Position = wp1.Position

	// find out if we arrived at a planet, either by reaching our target fleet
	// or reaching a planet
	if wp1.TargetType == MapObjectTypeFleet && wp1.TargetPlayerNum != None && wp1.TargetNum != None {
		target := mapObjectGetter.getFleet(wp1.TargetPlayerNum, wp1.TargetNum)
		fleet.OrbitingPlanetNum = target.OrbitingPlanetNum
	} else if wp1.TargetType == MapObjectTypePlanet && wp1.TargetNum != None {
		fleet.OrbitingPlanetNum = wp1.TargetNum
	} else if wp1.TargetType == MapObjectTypeWormhole && wp1.TargetNum != None {
		target := mapObjectGetter.getWormhole(wp1.TargetNum)
		dest := mapObjectGetter.getWormhole(target.DestinationNum)
		player.discoverer.discoverWormholeLink(target, dest)
		fleet.Position = dest.Position
		fleet.Waypoints[1] = NewPositionWaypoint(fleet.Position, fleet.Spec.Engine.IdealSpeed)
	}

	// if we wait at a waypoint while unloading, we "complete" our move but don't actually move
	// TODO: this is weird, can we just not complete a move if we are waiting at a waypoint?
	if !wp0.WaitAtWaypoint {
		fleet.Waypoints = fleet.Waypoints[1:]
	}

	// we arrived, process the current task (the previous waypoint)
	if len(fleet.Waypoints) == 1 {
		fleet.WarpSpeed = 0
		fleet.Heading = Vector{}
	} else {
		wp1 = fleet.Waypoints[1]
		fleet.WarpSpeed = wp1.WarpSpeed
		fleet.Heading = (wp1.Position.Subtract(fleet.Position)).Normalized()
	}
}

// colonize a planet
// TODO: return an error and stop colonization
func (fleet *Fleet) colonizePlanet(rules *Rules, player *Player, planet *Planet) {
	planet.PlayerNum = player.Num
	planet.ProductionQueue = []ProductionQueueItem{}
	planet.Cargo = planet.Cargo.Add(fleet.Cargo)
	fleet.Cargo = Cargo{}

	if len(player.ProductionPlans) > 0 {
		plan := player.ProductionPlans[0]
		plan.Apply(planet)
	}

	if player.Race.Spec.InnateMining {
		hab := player.Race.GetPlanetHabitability(planet.Hab)
		maxPop := planet.getMaxPopulation(rules, player, hab)
		planet.Mines = planet.innateMines(player, planet.productivePopulation(planet.population(), maxPop))
	}

	if player.Race.Spec.InnateScanner {
		planet.Scanner = true
	}

	planet.Spec = computePlanetSpec(rules, player, planet)
}

// get the minerals and resources recovered from a scrapped fleet
// from the stars wiki:
// After battle, 1/3 of the mineral cost of the destroyed ships is left as salvage. If the battle took place in orbit, these minerals are deposited on the planet below.
// In deep space, each type of mineral decays 10%, or 10kT per year, whichever is higher. Salvage deposited on planets does not decay.
// Scrapping: (from help file)
//
// A ship scrapped at a starbase deposits 80% of the original minerals on the planet, or 90% of the minerals and 70% of the resources if the LRT 'Ultimate Recycling' is selected.
// A ship scrapped at a planet with no starbase leaves 33% of the original minerals on the planet, or 45% of the minerals if the LRT Ultimate Recycling is selected.
// Wih UR the resources recovered is:
// (resources the ship costs * resources on the planet)/(resources the ship cost + resources on the planet)
// The maximum recoverable resources occurs when the cost of the scrapped ship equals the resources produced at the planet where it is scrapped.
//
// A ship scrapped in space leaves no minerals behind.
// When a ship design is deleted, all such ships vanish leaving nothing behind. (moral: scrap before you delete!)
func (fleet *Fleet) getScrapAmount(rules *Rules, player *Player, planet *Planet) Cost {

	// create a new cargo instance out of our fleet cost
	scrappedCost := fleet.Spec.Cost

	scrapMineralFactor := rules.ScrapMineralAmount
	scrapResourceFactor := rules.ScrapResourceAmount
	extraResources := 0
	planetResources := 0

	if planet != nil && planet.OwnedBy(player.Num) {
		planetResources = planet.Spec.ResourcesPerYear + planet.bonusResources
		// UR races get resources when scrapping
		if planet.Spec.HasStarbase {
			// scrapping over a planet we own with a starbase, calculate bonus minerals and resources
			scrapMineralFactor += player.Race.Spec.ScrapMineralOffsetStarbase
			scrapResourceFactor += player.Race.Spec.ScrapResourcesOffsetStarbase
		} else {
			// scrapping over a planet we own without a starbase, calculate bonus minerals and resources
			scrapMineralFactor += player.Race.Spec.ScrapMineralOffset
			scrapResourceFactor += player.Race.Spec.ScrapResourcesOffset
		}
	}

	// figure out much cargo and resources we get
	scrappedCost = scrappedCost.MultiplyFloat64(scrapMineralFactor)

	if scrapResourceFactor > 0 {
		// Formula for calculating resources: (Current planet production * Extra resources)/(Current planet production + Extra Resources)
		extraResources = int(float64(fleet.Spec.Cost.Resources)*scrapResourceFactor + .5)
		extraResources = int(float64(planetResources*extraResources) / float64(planetResources+extraResources))
		scrappedCost.Resources += extraResources
	} else {
		scrappedCost.Resources = 0
	}

	return scrappedCost
}

// getTransferAmount gets the amount of cargo to transfer for loading a cargo type from a cargoholder
func (fleet *Fleet) getCargoLoadAmount(dest cargoHolder, cargoType CargoType, task WaypointTransportTask) (transferAmount int, waitAtWaypoint bool) {
	availableCapacity := fleet.Spec.CargoCapacity - fleet.Cargo.Total()
	availableToLoad := dest.getCargo().GetAmount(cargoType)
	currentAmount := fleet.Cargo.GetAmount(cargoType)
	totalCapacity := fleet.Spec.CargoCapacity

	// fuel transfers use different tanks
	if cargoType == Fuel {
		availableCapacity = fleet.Spec.FuelCapacity - fleet.Fuel
		availableToLoad = dest.getFuel()
		currentAmount = fleet.Fuel
		totalCapacity = fleet.Spec.FuelCapacity
	}

	switch task.Action {
	case TransportActionLoadOptimal:
		// fuel only
		// we set our fuel to whatever it takes to finish our waypoints and transfer the rest to the ICargoHolder target.
		// If the target is a planet or starbase (and has infinite fuel capacity), we skip this and don't give them our fuel
		if cargoType == Fuel && dest.getFuelCapacity() != Unlimited {
			fuelRequiredForWaypoints := 0
			for i := 1; i < len(fleet.Waypoints); i++ {
				fuelRequiredForWaypoints += fleet.Waypoints[i].EstFuelUsage
			}
			leftoverFuel := fleet.Fuel - fuelRequiredForWaypoints
			fuelCapacityAvailable := dest.getFuelCapacity() - dest.getFuel()
			if leftoverFuel > 0 && fuelCapacityAvailable > 0 {
				// transfer the lowest of how much fuel capacity they have available or how much we can give
				// this is a bit weird because we are doing a "Load", but it's actually an unload of fuel
				// from us to a dest fleet, so make the transferAmount negative.
				transferAmount = MaxInt(-leftoverFuel, -(dest.getFuelCapacity() - dest.getFuel()))
			}
		}
	case TransportActionLoadAll:
		// load all available, based on our constraints
		transferAmount = MinInt(availableToLoad, availableCapacity)
	case TransportActionLoadAmount:
		transferAmount = MinInt(MinInt(availableToLoad, task.Amount), availableCapacity)
	case TransportActionWaitForPercent:
		fallthrough
	case TransportActionFillPercent:
		// we want a percent of our hold to be filled with some amount, figure out how
		// much that is in kT, i.e. 50% of 100kT would be 50kT of this mineral
		var taskAmountkT = int(float64(task.Amount) / 100 * float64(totalCapacity))

		if currentAmount >= taskAmountkT {
			// no need to transfer any, move on
			return 0, false
		} else {

			// transfer up to our percent specified
			// wait here if we haven't loaded the amount we want
			// but move on if we are out of cargo space (in case the user suffers from innumeracy and said they wanted 50% 50% 50%)
			transferAmount = MinInt(MinInt(availableToLoad, taskAmountkT-currentAmount), availableCapacity)
			if (transferAmount+currentAmount) < taskAmountkT && task.Action == TransportActionWaitForPercent && (availableCapacity-transferAmount) > 0 {
				waitAtWaypoint = true
			}
		}
	case TransportActionSetAmountTo:
		// only transfer the min of what we have, vs what we need, vs the capacity
		transferAmount = MaxInt(0, MinInt(MinInt(availableToLoad, task.Amount-currentAmount), availableCapacity))
		if transferAmount < (task.Amount - currentAmount) {
			waitAtWaypoint = true
		}
	case TransportActionSetWaypointTo:
		// Check how much the destination has of what we want
		// if we SetWaypointTo 100kT germanium and they have 120kT, we load 20kT if we can fit it
		if availableToLoad <= task.Amount {
			// they are below the amount, we won't load (this TransportAction will possibly be used to unload later)
			break
		} else {
			// only transfer down to what we set
			transferAmount = MinInt(MinInt(availableToLoad, availableToLoad-task.Amount), availableCapacity)
		}

	case TransportActionLoadDunnage:
		// (minerals and colonists only) This command waits until all other loads and unloads are
		// complete, then loads as many colonists or amount of a mineral as will fit in the remaining
		// space. For example, setting Load All Germanium, Load Dunnage Ironium, will load all the
		// Germanium that is available, then as much Ironium as possible. If more than one dunnage cargo
		// is specified, they are loaded in the order of Ironium, Boranium, Germanium, and Colonists.
		transferAmount = MinInt(availableToLoad, availableCapacity)
	}

	// let the caller know how much of this cargo we load
	return transferAmount, waitAtWaypoint
}

// getCargoUnloadAmount gets the amount of cargo to transfer for unloading a cargo type from a cargoholder
func (fleet *Fleet) getCargoUnloadAmount(dest cargoHolder, cargoType CargoType, task WaypointTransportTask) (transferAmount int, waitAtWaypoint bool) {

	capacity := dest.getCargoCapacity()
	currentAmount := fleet.Cargo.GetAmount(cargoType)

	var availableToUnload int
	if cargoType == Fuel {
		availableToUnload = fleet.Fuel
		capacity = MaxInt(0, dest.getFuelCapacity()-dest.getFuel())
		currentAmount = fleet.Fuel
	} else {
		availableToUnload = fleet.Cargo.GetAmount(cargoType)
	}
	switch task.Action {
	case TransportActionUnloadAll:
		// unload all available, based on our constraints
		if capacity == Unlimited {
			transferAmount = availableToUnload
		} else {
			transferAmount = MinInt(availableToUnload, capacity)
		}
	case TransportActionUnloadAmount:
		// don't unload more than the task says
		if capacity == Unlimited {
			transferAmount = MinInt(availableToUnload, task.Amount)
		} else {
			transferAmount = MinInt(MinInt(availableToUnload, task.Amount), capacity)
		}
	case TransportActionSetAmountTo:
		// set the amount in our hold to amount, or do nothing if we have under that amount
		transferAmount = MaxInt(0, MinInt(availableToUnload, currentAmount-task.Amount))
	case TransportActionSetWaypointTo:
		// Make sure the waypoint has at least whatever we specified
		var currentAmount = dest.getCargo().GetAmount(cargoType)

		if currentAmount >= task.Amount {
			// no need to transfer any, move on
			break
		} else {
			// only transfer the min of what we have, vs what we need, vs the capacity
			if capacity == Unlimited {
				transferAmount = MinInt(availableToUnload, task.Amount-currentAmount)
			} else {
				transferAmount = MinInt(MinInt(availableToUnload, task.Amount-currentAmount), capacity)
			}
		}
	}
	return transferAmount, waitAtWaypoint
}

// Repair a fleet. This changes based on where the fleet is
func (fleet *Fleet) repairFleet(log zerolog.Logger, rules *Rules, player *Player, orbiting *Planet) {
	needsRepair := false
	for _, token := range fleet.Tokens {
		if token.QuantityDamaged > 0 {
			needsRepair = true
		}
	}
	if !needsRepair {
		return
	}

	rate := RepairRateMoving

	if len(fleet.Waypoints) == 1 {
		if orbiting != nil {
			if fleet.Spec.Bomber && player.IsEnemy(orbiting.PlayerNum) {
				// no repairs while bombing
				rate = RepairRateNone
			} else {
				if orbiting.OwnedBy(player.Num) {
					rate = RepairRateOrbitingOwnPlanet
				} else {
					rate = RepairRateOrbiting
				}
			}
		} else {
			rate = RepairRateStopped
		}
	}

	repairRate := rules.RepairRates[rate]
	if repairRate > 0 {
		// apply any bonuses for the fleet
		repairRate += fleet.Spec.RepairBonus

		if rate == RepairRateOrbitingOwnPlanet && orbiting.Starbase != nil && !orbiting.Starbase.Delete {
			// apply any bonuses for the starbase if we own this planet and it has a starbase
			repairRate += orbiting.Starbase.Spec.RepairBonus
		}

		for i := range fleet.Tokens {
			token := &fleet.Tokens[i]

			// IS races double repair
			// repair some percentage of armor
			// 100dp armor@3% repair over a planet means
			// it repairs 3dp per turn. All damaged tokens repair
			// at the same rate
			repairAmount := MaxInt(1, int(float64(token.design.Spec.Armor)*repairRate*player.Race.Spec.RepairFactor))

			// Remove damage from this fleet by its armor * repairRate
			token.Damage = math.Max(0, token.Damage-float64(repairAmount))
			if token.Damage == 0 {
				token.QuantityDamaged = 0
			}

			log.Debug().
				Int("Player", fleet.PlayerNum).
				Str("Fleet", fleet.Name).
				Str("Token", token.design.Name).
				Int("RepairAmount", repairAmount).
				Int("QuantityDamaged", token.QuantityDamaged).
				Int("Damage", int(token.Damage)).
				Msgf("fleet token repaired")

		}
	}
}

// Repair a starbase
func (fleet *Fleet) repairStarbase(log zerolog.Logger, rules *Rules, player *Player) {
	repairRate := rules.RepairRates[RepairRateStarbase]
	token := &fleet.Tokens[0]

	// IS races repair starbases 1.5x
	repairAmount := MaxInt(1, int(float64(token.design.Spec.Armor)*repairRate*player.Race.Spec.StarbaseRepairFactor))

	// Remove damage from this fleet by its armor * repairRate
	token.Damage = math.Max(0, fleet.Tokens[0].Damage-float64(repairAmount))

	log.Debug().
		Int("Player", fleet.PlayerNum).
		Str("Fleet", fleet.Name).
		Str("Token", token.design.Name).
		Int("RepairAmount", repairAmount).
		Int("QuantityDamaged", token.QuantityDamaged).
		Int("Damage", int(token.Damage)).
		Msgf("starbase repaired")

}
