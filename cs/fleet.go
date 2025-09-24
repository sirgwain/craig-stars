package cs

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
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
	GameDBObject
	MapObject
	FleetOrders
	PlanetNum         int         `json:"planetNum"` // for starbase fleets that are owned by a planet
	BaseName          string      `json:"baseName"`
	Cargo             Cargo       `json:"cargo,omitzero"`
	Fuel              int         `json:"fuel"`
	Age               int         `json:"age"`
	Tokens            []ShipToken `json:"tokens"`
	Heading           Vector      `json:"heading"`
	WarpSpeed         int         `json:"warpSpeed,omitempty"`
	PreviousPosition  *Vector     `json:"previousPosition,omitempty"`
	OrbitingPlanetNum int         `json:"orbitingPlanetNum,omitempty"`
	Starbase          bool        `json:"starbase,omitempty"`
	Spec              FleetSpec   `json:"spec"`
	battlePlan        *BattlePlan
	struckMinefield   bool
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
	BaseCloakedCargo int                        `json:"baseCloakedCargo,omitempty"`
	HasMassDriver    bool                       `json:"hasMassDriver,omitempty"`
	HasStargate      bool                       `json:"hasStargate,omitempty"`
	MassEmpty        int                        `json:"massEmpty,omitempty"`
	Purposes         map[ShipDesignPurpose]bool `json:"purposes,omitzero"`
	TotalShips       int                        `json:"totalShips,omitempty"`
}

type Waypoint struct {
	MapObjectTarget
	Position             Vector                 `json:"position"`
	WarpSpeed            int                    `json:"warpSpeed"`
	EstFuelUsage         int                    `json:"estFuelUsage,omitempty"`
	Task                 WaypointTask           `json:"task,omitempty"`
	TransportTasks       WaypointTransportTasks `json:"transportTasks,omitzero"`
	WaitAtWaypoint       bool                   `json:"waitAtWaypoint,omitempty"`
	LayMinefieldDuration int                    `json:"layMinefieldDuration,omitempty"`
	PatrolRange          int                    `json:"patrolRange,omitempty"`
	PatrolWarpSpeed      int                    `json:"patrolWarpSpeed,omitempty"`
	TransferToPlayer     int                    `json:"transferToPlayer,omitempty"`
	PartiallyComplete    bool                   `json:"partiallyComplete,omitempty"`
	processed            bool                   `json:"-"`
}

type WaypointTask string

const (
	WaypointTaskNone           WaypointTask = ""
	WaypointTaskTransport      WaypointTask = "Transport"
	WaypointTaskColonize       WaypointTask = "Colonize"
	WaypointTaskRemoteMining   WaypointTask = "RemoteMining"
	WaypointTaskMergeWithFleet WaypointTask = "MergeWithFleet"
	WaypointTaskScrapFleet     WaypointTask = "ScrapFleet"
	WaypointTaskLayMinefield   WaypointTask = "LayMinefield"
	WaypointTaskPatrol         WaypointTask = "Patrol"
	WaypointTaskRoute          WaypointTask = "Route"
	WaypointTaskTransferFleet  WaypointTask = "TransferFleet"
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

// TODO: Add a "set waypoint to %" command

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

	// Loads up to the specified portion of the cargo hold, subject to amount available at waypoint and room left in hold.
	TransportActionFillPercent WaypointTaskTransportAction = "FillPercent"

	// Remain at the waypoint until exactly X % of the hold is filled.
	TransportActionWaitForPercent WaypointTaskTransportAction = "WaitForPercent"

	// (minerals and colonists only) This command waits until all other loads and unloads are complete,
	// then loads as many colonists or minerals will fit in the remaining space.
	// If more than one dunnage cargo is specified, they are performed in
	// the order of Ironium, Boranium, Germanium, and Colonists.
	TransportActionLoadDunnage WaypointTaskTransportAction = "LoadDunnage"

	// Load or unload the cargo until the amount on board is the amount specified.
	// If less than the specified cargo is available, the fleet will not move on.
	TransportActionSetAmountTo WaypointTaskTransportAction = "SetAmountTo"

	// Load or unload the cargo until the amount at the waypoint is the amount specified.
	// This order is always carried out to the best of the fleet’s ability that turn
	// but does not prevent the fleet from moving on.
	TransportActionSetWaypointTo WaypointTaskTransportAction = "SetWaypointTo"
)

// the purpose for a fleet's existence (ie what it's supposed to be doing),
// exported to allow the AI to plan ship movements
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
	case ShipDesignPurposeBomber, ShipDesignPurposeSmartBomber, ShipDesignPurposeStructureBomber:
		return FleetPurposeBomber
	case ShipDesignPurposeStartingFighter, ShipDesignPurposeFighterScout:
		return FleetPurposeFighter
	case ShipDesignPurposeTorpedoFighter, ShipDesignPurposeBeamFighter:
		return FleetPurposeCapitalShip
	case ShipDesignPurposeFreighter, ShipDesignPurposeFuelFreighter, ShipDesignPurposeMultiPurposeFreighter:
		return FleetPurposeFreighter
	case ShipDesignPurposeColonistFreighter:
		return FleetPurposeColonistFreighter
	case ShipDesignPurposeArmedFreighter:
		return FleetPurposeArmedFreighter
	case ShipDesignPurposeMiner:
		return FleetPurposeMiner
	case ShipDesignPurposeTerraformer:
		return FleetPurposeTerraformer
	case ShipDesignPurposeDamageMineLayer, ShipDesignPurposeSpeedMineLayer:
		return FleetPurposeMineLayer
	}
	return FleetPurposeNone
}

// get fleet purpose from hull type; design purpose used only as backup for niche cases
func FleetPurposeFromTechHullType(hull TechHullType, purpose ShipDesignPurpose) FleetPurpose {
	switch hull {
	case TechHullTypeScout:
		return FleetPurposeScout
	case TechHullTypeColonizer:
		return FleetPurposeColonizer
	case TechHullTypeBomber:
		return FleetPurposeBomber
	case TechHullTypeFighter:
		return FleetPurposeFighter
	case TechHullTypeCapitalShip:
		return FleetPurposeCapitalShip
	case TechHullTypeFreighter, TechHullTypeFuelTransport, TechHullTypeMultiPurposeFreighter:
		return FleetPurposeFreighter
	case TechHullTypeMineLayer:
		return FleetPurposeMineLayer
	case TechHullTypeMiner:
		if purpose == ShipDesignPurposeMiner {
			return FleetPurposeMiner
		}
		return FleetPurposeTerraformer
	}
	return FleetPurposeNone
}

type fleetMoveInterruptedReason int

const (
	fleetMoveInterruptedEngineFailure = iota
	fleetMoveInterruptedHitMinefield
)

type fleetMoveInterrupted struct {
	reason    fleetMoveInterruptedReason
	minefield *Minefield
}

func NewFleet(player *Player, num int, name string, waypoints []Waypoint) *Fleet {
	return &Fleet{
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
	f.Name = fmt.Sprintf("%s #%d", f.BaseName, f.Num)
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
		Position: position,
		MapObjectTarget: MapObjectTarget{
			TargetType:      MapObjectTypePlanet,
			TargetNum:       num,
			TargetName:      name,
			TargetPlayerNum: None,
		},
		WarpSpeed: warpSpeed,
	}
}

func NewFleetWaypoint(position Vector, num int, playerNum int, name string, warpSpeed int) Waypoint {
	return Waypoint{
		Position: position,
		MapObjectTarget: MapObjectTarget{
			TargetType:      MapObjectTypeFleet,
			TargetNum:       num,
			TargetPlayerNum: playerNum,
			TargetName:      name,
		},
		WarpSpeed: warpSpeed,
	}
}

func NewMysteryTraderWaypoint(mt *MysteryTrader, warpSpeed int) Waypoint {
	return Waypoint{
		Position: mt.Position,
		MapObjectTarget: MapObjectTarget{
			TargetType: mt.Type,
			TargetNum:  mt.Num,
			TargetName: mt.Name,
		},
		WarpSpeed: warpSpeed,
	}
}

func NewPositionWaypoint(position Vector, warpSpeed int) Waypoint {
	return Waypoint{
		Position:  position,
		WarpSpeed: warpSpeed,
		MapObjectTarget: MapObjectTarget{
			TargetNum:       None,
			TargetPlayerNum: None,
		},
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
func (tt WaypointTransportTasks) getTransportTasks() transportTaskByType {
	tasks := transportTaskByType{}
	if tt.Fuel.Action != TransportActionNone {
		tasks[Fuel] = tt.Fuel
	}
	if tt.Ironium.Action != TransportActionNone {
		tasks[Ironium] = tt.Ironium
	}
	if tt.Boranium.Action != TransportActionNone {
		tasks[Boranium] = tt.Boranium
	}
	if tt.Germanium.Action != TransportActionNone {
		tasks[Germanium] = tt.Germanium
	}
	if tt.Colonists.Action != TransportActionNone {
		tasks[Colonists] = tt.Colonists
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

// compute all the computable values of this fleet (cargo capacity, armor, mass, etc.)
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
		spec.MaxPopulation = max(spec.MaxPopulation, token.design.Spec.MaxPopulation)
		spec.InnateScanRangePenFactor = max(spec.InnateScanRangePenFactor, token.design.Spec.InnateScanRangePenFactor) // Ultra Station and Death Stars have pen scanning

		// use the lowest ideal speed for this fleet
		// if we have multiple engines
		if len(token.design.Spec.Engine.FuelUsage) > 0 {
			if spec.Engine.IdealSpeed == 0 {
				spec.Engine.IdealSpeed = token.design.Spec.Engine.IdealSpeed
				spec.Engine.FreeSpeed = token.design.Spec.Engine.FreeSpeed
				spec.Engine.MaxSafeSpeed = token.design.Spec.Engine.MaxSafeSpeed
			} else {
				spec.Engine.IdealSpeed = min(spec.Engine.IdealSpeed, token.design.Spec.Engine.IdealSpeed)
				spec.Engine.FreeSpeed = token.design.Spec.Engine.FreeSpeed
				spec.Engine.MaxSafeSpeed = min(spec.Engine.MaxSafeSpeed, token.design.Spec.Engine.MaxSafeSpeed)
			}
		}
		// cost
		spec.Cost = MultiplyCost(token.design.Spec.Cost, token.Quantity)

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
		// TODO: Rework radiating into a custom datatype
		spec.Radiating = spec.Radiating || token.design.Spec.Radiating

		// spec all mine layers in the fleet
		if token.design.Spec.CanLayMines {
			spec.CanLayMines = true
			if spec.MineLayingRateByMineType == nil {
				spec.MineLayingRateByMineType = make(map[MinefieldType]int)
			}
			for key := range token.design.Spec.MineLayingRateByMineType {
				if _, ok := spec.MineLayingRateByMineType[key]; ok {
					spec.MineLayingRateByMineType[key] = 0
				}
				spec.MineLayingRateByMineType[key] += token.design.Spec.MineLayingRateByMineType[key] * token.Quantity
			}
		}

		// We should only have one ship stack with spacdock capabilities, but for this logic just go with the max
		spec.SpaceDock = max(spec.SpaceDock, token.design.Spec.SpaceDock)

		// sadly, the fleet only gets the best repair bonus from one design
		spec.RepairBonus = max(spec.RepairBonus, token.design.Spec.RepairBonus)

		spec.ScanRange = max(spec.ScanRange, token.design.Spec.ScanRange)
		spec.ScanRangePen = max(spec.ScanRangePen, token.design.Spec.ScanRangePen)
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
			// calculate the cloak units for this token based on the design's cloak units (70 cloak units / kT for a stealth cloak)
			spec.CloakUnits += token.design.Spec.CloakUnits
		} else if !player.Race.Spec.FreeCargoCloaking {
			// if this ship doesn't have cloaking, it counts as cargo (except for races with free cargo cloaking)
			spec.BaseCloakedCargo += token.design.Spec.Mass * token.Quantity
		}

		// choose the best tachyon detector ship
		spec.ReduceCloaking = min(spec.ReduceCloaking, token.design.Spec.ReduceCloaking)

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
	spec.CloakPercent = computeFleetCloakPercent(&spec, fleet.Cargo.Total()+spec.BaseCloakedCargo, player.Race.Spec.FreeCargoCloaking)

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
			wp.EstFuelUsage = f.GetFuelCost(player, wp.WarpSpeed, math.Ceil(wp.Position.DistanceTo(wpPrevious.Position)))
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
	fleet.Fuel = min(fleet.Spec.FuelCapacity, fleet.Fuel)
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
			fleet.Cargo.Colonists = min(fleet.Cargo.Colonists, capacity)
		}

		// if we have 110kT of space and 10kT is taken up by colonists, we have 100kT remaining capacity
		// if we have 200kT of minerals left, we keep half of each
		minerals := fleet.Cargo.ToMineral()
		remainingCapacity := max(0, capacity-fleet.Cargo.Colonists)

		// if we have no capacity left, drop all minerals and
		if remainingCapacity == 0 {
			fleet.Cargo = Cargo{Colonists: fleet.Cargo.Colonists}
			return cargo.Subtract(fleet.Cargo)
		}
		totalMinerals := minerals.Total()

		// reduce each mineral by a percent
		percentToKeep := 1 / (float64(totalMinerals) / float64(remainingCapacity))
		minerals = minerals.MultiplyFloat64(percentToKeep, math.Floor)
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
	dist = min(totalDist, dist)

	// check for CE engine failure
	if player.Race.Spec.EngineFailureRate > 0 &&
		wp1.WarpSpeed > player.Race.Spec.EngineReliableSpeed &&
		player.Race.Spec.EngineFailureRate >= rules.random.Float64() {
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
		hitMinefield, actualDist := checkForMinefieldCollision(rules, playerGetter, mapObjectGetter, fleet, wp1, dist)
		if hitMinefield != nil {
			interrupted = &fleetMoveInterrupted{reason: fleetMoveInterruptedHitMinefield, minefield: hitMinefield}
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
		hitMinefield, actualDist := checkForMinefieldCollision(rules, playerGetter, mapObjectGetter, fleet, wp1, dist)
		if hitMinefield != nil {
			interrupted = &fleetMoveInterrupted{reason: fleetMoveInterruptedHitMinefield, minefield: hitMinefield}
		}

		if actualDist != dist {
			dist = actualDist
			fuelCost = fleet.GetFuelCost(player, wp1.WarpSpeed, dist)
		}

		fleet.Fuel -= fuelCost
		fuelGenerated = fleet.getFuelGeneration(wp1.WarpSpeed, dist)
	}

	// message the player about fuel generation
	fuelGenerated = min(fuelGenerated, fleet.Spec.FuelCapacity-fleet.Fuel)
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

		if fleet.struckMinefield {
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
			messager.fleetStargateInvalidSourceOwner(player, fleet, wp0)
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
	minSafeHullMass := min(sourceStargate.SafeHullMass, destStargate.SafeHullMass)

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
	fleet.applyOvergatePenalty(rules, player, totalDist, wp0, wp1, sourceStargate, destStargate)

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
				if rules.FleetSafeSpeedExplosionChance >= rules.random.Float64() {
					explodedShips++
					token.Quantity--
				}
			}
		}
	}

	return explodedShips
}

// applyOvergatePenalty damages and/or vanishes ShipTokens inside overgating fleets based on distance.
func (fleet *Fleet) applyOvergatePenalty(rules *Rules, player *Player, distance float64, wp0, wp1 Waypoint, sourceStargate, destStargate PlanetStarbaseSpec) {
	var totalDamage, shipsLostToDamage, shipsLostToTheVoid, startingShips int
	for i := range fleet.Tokens {
		token := &fleet.Tokens[i]
		startingShips += token.Quantity
		// IT players never lose ships to the void, but everyone else does
		if player.Race.Spec.ShipsVanishInVoid {
			shipsLostToTheVoid += token.applyOvergateVanishing(rules, distance, sourceStargate.SafeRange, sourceStargate.SafeHullMass)
		}

		// damage any remaining tokens if we have any left
		tokenDamage := token.applyOvergateDamage(distance, sourceStargate.SafeRange, sourceStargate.SafeHullMass, destStargate.SafeHullMass, rules.StargateMaxHullMassFactor)
		totalDamage += tokenDamage.damage
		shipsLostToDamage += tokenDamage.shipsDestroyed
	}

	// remove any tokens that were lost completely
	fleet.removeEmptyTokens()

	if len(fleet.Tokens) == 0 {
		messager.fleetStargateDestroyed(player, fleet, wp0, wp1)
	} else if totalDamage > 0 || shipsLostToTheVoid > 0 {
		messager.fleetStargateDamaged(player, fleet, wp0, wp1, totalDamage, shipsLostToDamage, shipsLostToTheVoid)
	}
}

// Engine fuel usage calculation courtesy of m.a@stars
func (engine Engine) getFuelCostForEngine(warpSpeed int, mass int, dist float64, ifeFactor float64) int {
	if warpSpeed == 0 {
		return 0
	}
	// 1 mg of fuel will move 200kT of weight 1 LY at a Fuel Usage Number of 100.
	// Number of engines doesn't matter, nor number of ships with the same engine.

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

// Get the fuel cost for this fleet to travel a certain distance at a certain speed,
// exported to allow for cross package usage.
func (fleet *Fleet) GetFuelCost(player *Player, warpSpeed int, distance float64) int {
	if warpSpeed == StargateWarpSpeed {
		return 0
	}
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
			// @sirgwain: Consider making this allocate cargo optimally for least fuel usage as QoL option
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
		planet.Mines = innateMines(player.Race.Spec.InnateMinesFactor, planet.GetPopulation())
	}

	if player.Race.Spec.InnateScanner {
		planet.Scanner = true
	}

	planet.Spec = ComputePlanetSpec(rules, player, planet)
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
func (fleet *Fleet) getScrapAmount(rules *Rules, player *Player, planet *Planet, colonize bool) Cost {

	// create a new cargo instance out of our fleet cost
	scrappedCost := fleet.Spec.Cost
	scrapMineralFactor := rules.ScrapMineralAmount
	scrapResourceFactor := rules.ScrapResourceAmount
	extraResources := 0
	planetResources := 0

	if planet != nil {
		if colonize {
			scrapMineralFactor = rules.ScrapColonizeAmount
		} else {
			planetResources = planet.Spec.ResourcesPerYear + planet.bonusResources
			// UR races get resources when scrapping (not colonizing)
			if planet.Spec.HasStarbase {
				// scrapping over a planet with a starbase, calculate bonus minerals and resources
				scrapMineralFactor += player.Race.Spec.ScrapMineralOffsetStarbase
				scrapResourceFactor += player.Race.Spec.ScrapResourcesOffsetStarbase
			} else {
				// scrapping over a planet without a starbase, calculate bonus minerals and resources
				scrapMineralFactor += player.Race.Spec.ScrapMineralOffset
				scrapResourceFactor += player.Race.Spec.ScrapResourcesOffset
			}
		}
	}

	// figure out much cargo and resources we get
	scrappedCost = MultiplyCost(scrappedCost, scrapMineralFactor)

	if scrapResourceFactor > 0 {
		// Formula for calculating resources: (Current planet production * Extra resources)/(Current planet production + Extra Resources)
		extraResources = int(float64(fleet.Spec.Cost.Resources)*scrapResourceFactor + .5) // add 0.5 to round up
		extraResources = int(float64(planetResources*extraResources) / float64(planetResources+extraResources))
		scrappedCost.Resources = extraResources
	} else {
		scrappedCost.Resources = 0
	}

	return scrappedCost
}

// Repair a fleet. This changes based on where the fleet is
func (fleet *Fleet) repairFleet(log *slog.Logger, rules *Rules, player *Player, orbiting *Planet) {
	needsRepair := false
	// Check if anything even needs repairing
	for _, token := range fleet.Tokens {
		if token.QuantityDamaged > 0 {
			needsRepair = true
			break
		}
	}

	if !needsRepair {
		return
	}

	var rate RepairRate
	switch {
	case len(fleet.Waypoints) > 1:
		rate = RepairRateMoving
	case orbiting == nil:
		// we're standing still, but not at a planet
		rate = RepairRateStopped
	case fleet.Spec.Bomber && player.IsEnemy(orbiting.PlayerNum):
		// no repairs while bombing
		rate = RepairRateNone
	case orbiting.OwnedBy(player.Num):
		// Confirmed - boosted repairs only occur on _your_ planets (not allies')
		rate = RepairRateOrbitingOwnPlanet
	default:
		rate = RepairRateOrbiting
	}

	repairRate := rules.RepairRates[rate]
	if repairRate <= 0 {
		return
	}

	// apply any bonuses for this fleet
	// TODO: Should this apply to _all_ fleets at the given location?
	// We could probably pre-screen fleets at the same location for global
	// repair bonus the first time around and recycle that value for subsequent ones
	repairRate += fleet.Spec.RepairBonus

	if rate == RepairRateOrbitingOwnPlanet && orbiting.Starbase != nil && !orbiting.Starbase.Delete {
		// apply any bonuses from orbiting our own starbase (if present)
		repairRate += orbiting.Starbase.Spec.RepairBonus
	}

	for i := range fleet.Tokens {
		token := &fleet.Tokens[i]

		// IS races double repair
		// repair some percentage of armor
		// 100dp armor@3% repair over a planet means
		// it repairs 3dp per turn. All damaged tokens repair
		// at the same rate
		repairAmount := max(1, int(float64(token.design.Spec.Armor)*repairRate*player.Race.Spec.RepairFactor))

		// Remove damage from this fleet by its armor * repairRate
		token.Damage = math.Floor(max(0, token.Damage-float64(repairAmount)))
		if token.Damage == 0 {
			token.QuantityDamaged = 0
		}

		log.Debug("fleet token repaired",
			slog.Int("Player", fleet.PlayerNum),
			slog.String("Fleet", fleet.Name),
			slog.String("Token", token.design.Name),
			slog.Int("RepairAmount", repairAmount),
			slog.Int("QuantityDamaged", token.QuantityDamaged),
			slog.Int("Damage", int(token.Damage)))

	}
}

// Repair a starbase
func (fleet *Fleet) repairStarbase(log *slog.Logger, rules *Rules, player *Player) {
	repairRate := rules.RepairRates[RepairRateStarbase]
	token := &fleet.Tokens[0]

	// IS races repair starbases 1.5x
	repairAmount := max(1, int(float64(token.design.Spec.Armor)*repairRate*player.Race.Spec.StarbaseRepairFactor))

	// Remove damage from this fleet by its armor * repairRate
	token.Damage = math.Floor(max(0, fleet.Tokens[0].Damage-float64(repairAmount)))

	log.Debug("starbase repaired",
		slog.Int("Player", fleet.PlayerNum),
		slog.String("Fleet", fleet.Name),
		slog.String("Token", token.design.Name),
		slog.Int("RepairAmount", repairAmount),
		slog.Int("QuantityDamaged", token.QuantityDamaged),
		slog.Int("Damage", int(token.Damage)))

}

type WaypointDest struct {
	MO       MapObject `json:"mo"`
	Position Vector    `json:"position"`
}

// CanColonize returns true if this fleet can colonize the planet
func (f *Fleet) CanColonize(planet *Planet) bool {
	return f.Spec.Colonizer && f.Cargo.Colonists > 0 && planet != nil && !planet.Owned() && planet.Spec.TerraformedHabitability > 0
}

// CanFuel returns true if this fleet will refuel at the planet
func (f *Fleet) CanFuel(player *Player, planet *Planet) bool {
	return planet != nil && planet.Owned() && planet.Spec.DockCapacity != 0 && player.IsFriend(planet.PlayerNum)
}

// CanRemoteMine returns true if this fleet can remote mine the planet
func (f *Fleet) CanRemoteMine(player *Player, planet *Planet) bool {
	return f.Spec.MiningRate > 0 && planet != nil && !planet.Owned() || (player.Race.Spec.CanRemoteMineOwnPlanets && planet.OwnedBy(player.Num))
}

// CanRemoteMine returns true if this fleet can remote mine the planet
func (f *Fleet) CanJump(player *Player, dist float64, orbiting *Planet, target *Planet) bool {
	if target == nil {
		return false
	}

	highestShipMass := 0
	for _, token := range f.Tokens {
		highestShipMass = max(highestShipMass, token.design.Spec.Mass)
	}

	destSafeHullMass := target.Spec.SafeHullMass
	destSafeRange := target.Spec.SafeRange
	destStargateSafe :=
		target.Spec.HasStargate &&
			target.Owned() &&
			player.IsFriend(target.PlayerNum) &&
			float64(destSafeRange) >= dist &&
			highestShipMass <= destSafeHullMass

	if f.Spec.CanJump {
		// we have a jump gate installed in our ship, we only care about the destination gate
		return destStargateSafe
	} else {
		if orbiting == nil || !orbiting.Spec.HasStargate {
			return false
		}
		canGateCargo := player.Race.Spec.CanGateCargo
		sourceSafeHullMass := orbiting.Spec.SafeHullMass
		sourceSafeRange := orbiting.Spec.SafeRange
		sourceStargateSafe :=
			(canGateCargo || f.Cargo.Total() == 0) &&
				orbiting.Owned() &&
				player.IsFriend(target.PlayerNum) &&
				float64(sourceSafeRange) >= dist &&
				highestShipMass <= sourceSafeHullMass
		return destStargateSafe && sourceStargateSafe
	}

}

type waypointInfo struct {
	selectedWaypoint *Waypoint
	nextWaypoint     *Waypoint
	previousWaypoint *Waypoint
	waypointIndex    int
}

func (f *Fleet) getSelectedWaypointInfo(currentSelectedWaypointIndex int) waypointInfo {
	index := currentSelectedWaypointIndex
	if index == -1 || index >= len(f.Waypoints) {
		index = 0
	}

	info := waypointInfo{
		selectedWaypoint: &f.Waypoints[index],
		waypointIndex:    index,
	}

	if index > 0 {
		info.previousWaypoint = &f.Waypoints[index-1]
	}
	if index < len(f.Waypoints)-1 {
		info.nextWaypoint = &f.Waypoints[index+1]
	}

	return info
}

// AddWaypoint adds a new waypoint at a destination after the currentSelectedWaypointIndex.
// This function returns the index of the newly added waypoint, or 0 if no waypoint is added
func (f *Fleet) AddWaypoint(
	player *Player,
	dest WaypointDest,
	currentSelectedWaypointIndex int,
	fastestWaypoint bool,
) int {
	f.computeFuelUsage(player)
	info := f.getSelectedWaypointInfo(currentSelectedWaypointIndex)
	selectedWaypoint := info.selectedWaypoint
	nextWaypoint := info.nextWaypoint
	index := info.waypointIndex

	// the position is either a position or the target's position
	position := dest.Position
	if position == (Vector{}) && dest.MO.Type != MapObjectTypeNone {
		position = dest.MO.Position
	}

	if position == (Vector{}) ||
		position == selectedWaypoint.Position ||
		(nextWaypoint != nil && position == nextWaypoint.Position) {
		slog.Debug("Not adding waypoint position",
			slog.String("position", position.String()),
			slog.String("selectedWaypoint.Position", selectedWaypoint.Position.String()))
		return 0 // don't add duplicate waypoint
	}

	var targetPlanet *Planet
	if dest.MO.Type == MapObjectTypePlanet {
		targetPlanet = player.GetPlanetIntel(dest.MO.Num)
	}

	// if we are targeting a planet, record if we can colonize or remote mine it for later
	canColonize := f.CanColonize(targetPlanet)
	canRemoteMine := f.CanRemoteMine(player, targetPlanet)

	fuelAlreadyAllocated := f.GetFuelAllocated(player, index)
	var orbiting *Planet
	if selectedWaypoint.TargetType == MapObjectTypePlanet {
		orbiting = player.GetPlanetIntel(selectedWaypoint.TargetNum)
	}

	dist := math.Ceil(selectedWaypoint.Position.DistanceTo(position))

	// determine what warp we should set for this waypoint
	warpSpeed := f.GetWarpSpeed(player, dist, orbiting, targetPlanet, fuelAlreadyAllocated, fastestWaypoint)

	if dest.MO.Type != MapObjectTypeNone {
		wp := Waypoint{
			Position: dest.MO.Position,
			MapObjectTarget: MapObjectTarget{
				TargetName:      dest.MO.Name,
				TargetPlayerNum: dest.MO.PlayerNum,
				TargetNum:       dest.MO.Num,
				TargetType:      dest.MO.Type,
				TargetPosition:  dest.MO.Position,
			},
			WarpSpeed:      warpSpeed,
			Task:           selectedWaypoint.Task,
			TransportTasks: selectedWaypoint.TransportTasks,
		}
		if canColonize {
			wp.Task = WaypointTaskColonize
			wp.TransportTasks = WaypointTransportTasks{}
		} else if canRemoteMine {
			wp.Task = WaypointTaskRemoteMining
			wp.TransportTasks = WaypointTransportTasks{}
		}
		wp.EstFuelUsage = f.GetFuelCost(player, wp.WarpSpeed, dist)
		f.Waypoints = slices.Insert(f.Waypoints, index+1, wp)
	} else {
		wp := Waypoint{
			Position: position,
			MapObjectTarget: MapObjectTarget{
				TargetPosition: position,
			},
			WarpSpeed:      warpSpeed,
			Task:           selectedWaypoint.Task,
			TransportTasks: selectedWaypoint.TransportTasks,
		}
		wp.EstFuelUsage = f.GetFuelCost(player, wp.WarpSpeed, dist)
		f.Waypoints = slices.Insert(f.Waypoints, index+1, wp)
	}

	return index + 1
}

// UpdateWaypoint updates an existing waypoint from a drag/drop type operation
func (f *Fleet) UpdateWaypoint(
	player *Player,
	dest WaypointDest,
	currentSelectedWaypointIndex int,
	fastestWaypoint bool,
) UpdateWaypointResult {
	info := f.getSelectedWaypointInfo(currentSelectedWaypointIndex)

	selectedWaypoint := info.selectedWaypoint
	previousWaypoint := info.previousWaypoint
	nextWaypoint := info.nextWaypoint
	waypointIndex := info.waypointIndex

	if previousWaypoint == nil {
		// can't update wp0
		return UpdateWaypointResultNone
	}

	f.computeFuelUsage(player)

	// the position is either a position or the target's position
	position := dest.Position
	if dest.MO.Type != MapObjectTypeNone {
		position = dest.MO.Position
	}

	if position == (Vector{}) {
		return UpdateWaypointResultNone
	}

	if position == previousWaypoint.Position {
		// don't update a waypoint to be the same as a previous waypoint, this should just delete it
		return UpdateWaypointResultPreviousWaypoint
	}

	if nextWaypoint != nil && position == nextWaypoint.Position {
		// don't update a waypoint to be the same as a previous waypoint, this should just delete it
		return UpdateWaypointResultNextWaypoint
	}

	dist := math.Ceil(previousWaypoint.Position.DistanceTo(position))

	// get the fuel allocated up to but not including this waypoint since we're moving it around
	fuelAlreadyAllocated := f.GetFuelAllocated(player, waypointIndex-1)

	var orbiting *Planet
	if previousWaypoint.TargetType == MapObjectTypePlanet {
		orbiting = player.GetPlanetIntel(previousWaypoint.TargetNum)
	}

	var targetPlanet *Planet
	if dest.MO.Type == MapObjectTypePlanet {
		targetPlanet = player.GetPlanetIntel(dest.MO.Num)
	}

	// if we are targeting a planet, record if we can colonize or remote mine it for later
	canColonize := f.CanColonize(targetPlanet)
	canRemoteMine := f.CanRemoteMine(player, targetPlanet)

	// determine what warp we should set for this waypoint
	warpSpeed := f.GetWarpSpeed(player, dist, orbiting, targetPlanet, fuelAlreadyAllocated, fastestWaypoint)

	if dest.MO.Type != MapObjectTypeNone {
		selectedWaypoint.Position = dest.MO.Position
		selectedWaypoint.MapObjectTarget = dest.MO.ToTarget()
		selectedWaypoint.WarpSpeed = warpSpeed

		if canColonize {
			selectedWaypoint.Task = WaypointTaskColonize
			selectedWaypoint.TransportTasks = WaypointTransportTasks{}
		} else if canRemoteMine {
			selectedWaypoint.Task = WaypointTaskRemoteMining
			selectedWaypoint.TransportTasks = WaypointTransportTasks{}
		}
	} else {
		selectedWaypoint.Position = position
		selectedWaypoint.MapObjectTarget = MapObjectTarget{
			TargetPosition: position,
		}
		selectedWaypoint.WarpSpeed = warpSpeed
	}

	selectedWaypoint.EstFuelUsage = f.GetFuelCost(player, selectedWaypoint.WarpSpeed, dist)

	return UpdateWaypointResultUpdated
}

// GetFuelAllocated gets the fuel allocated up to the waypointIndex accounting for any refueling
func (f *Fleet) GetFuelAllocated(player *Player, waypointIndex int) int {
	fuelAllocated := 0
	for i := 0; i <= waypointIndex && i < len(f.Waypoints); i++ {
		wp := f.Waypoints[i]
		fuelAllocated += wp.EstFuelUsage

		var targetPlanet *Planet
		if wp.TargetType == MapObjectTypePlanet {
			targetPlanet = player.GetPlanetIntel(wp.TargetNum)
			if targetPlanet != nil && f.CanFuel(player, targetPlanet) {
				fuelAllocated = 0
			}
		}
	}
	return fuelAllocated
}

// GetWarpSpeed returns the warp speed this fleet should go given the distance, target, and orbiting planet
// this will go max speed if we can colonize the destination
func (f *Fleet) GetWarpSpeed(
	player *Player,
	dist float64,
	orbiting *Planet,
	targetPlanet *Planet,
	fuelAlreadyAllocated int,
	fastestWaypoint bool,
) int {
	canColonize := false
	canJump := false
	canFuel := false

	if targetPlanet != nil {
		canColonize = f.CanColonize(targetPlanet)
		canJump = f.CanJump(player, dist, orbiting, targetPlanet)
		canFuel = f.CanFuel(player, targetPlanet)
	}

	engine := f.Spec.Engine

	var warpSpeed int
	if canJump {
		warpSpeed = StargateWarpSpeed
	} else if canFuel || canColonize || fastestWaypoint {
		warpSpeed = f.GetMaxWarp(
			player,
			fuelAlreadyAllocated,
			dist,
			engine.FreeSpeed,
			engine.MaxSafeSpeed,
		)
	} else {
		warpSpeed = f.GetMinimalWarp(
			player,
			fuelAlreadyAllocated,
			dist,
			engine.IdealSpeed,
			engine.FreeSpeed,
			engine.MaxSafeSpeed,
		)
	}

	return warpSpeed
}

// GetMinimalWarp returns the highest useful speed less than or equal to a given warp speed
// to reach a given destination.
func (f *Fleet) GetMinimalWarp(
	player *Player,
	fuelAlreadyAllocated int,
	dist float64,
	startSpeed int,
	freeSpeed int,
	maxSafeSpeed int,
) int {
	idealYears := int(math.Ceil(float64(dist) / float64(startSpeed*startSpeed)))
	speed := startSpeed

	// Prefer slower speeds if they take the same time
	for i := startSpeed; i > freeSpeed; i-- {
		years := int(math.Ceil(float64(dist) / float64(i*i)))
		if years == idealYears {
			speed = i
		}
	}

	// Decrease speed until we can afford it
	for speed >= freeSpeed {
		fuelUsed := f.GetFuelCost(player, speed, dist)
		if fuelUsed+fuelAlreadyAllocated > f.Fuel {
			speed--
			continue
		}
		break
	}

	if speed > maxSafeSpeed {
		speed = maxSafeSpeed
	}

	return speed
}

// GetMaxWarp returns the max warp we have fuel for to make it to the destination
func (f *Fleet) GetMaxWarp(
	player *Player,
	fuelAlreadyAllocated int,
	dist float64,
	freeSpeed int,
	maxSafeSpeed int,
) int {
	speed := freeSpeed

	for i := speed + 1; i <= maxSafeSpeed; i++ {
		fuelUsed := f.GetFuelCost(player, i, dist)
		if fuelUsed+fuelAlreadyAllocated > f.Fuel {
			break
		}
		speed = i
	}

	idealSpeed := f.Spec.Engine.IdealSpeed
	idealFuelUsed := f.GetFuelCost(player, idealSpeed, dist)

	if freeSpeed > 1 && speed < idealSpeed && idealFuelUsed <= f.Fuel {
		speed = idealSpeed
	}

	// Don't go faster than needed
	return f.GetMinimalWarp(
		player,
		fuelAlreadyAllocated,
		dist,
		speed,
		freeSpeed,
		maxSafeSpeed,
	)
}
