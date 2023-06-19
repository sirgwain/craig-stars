package cs

import (
	"fmt"
	"math"
)

// warpfactor for using a stargate vs moving with warp drive
const StargateWarpFactor = 11

// time period to perform a task, like patrol
const Indefinite = -1

// use automatic warp factor for patrols
const PatrolWarpFactorAutomatic = -1

// target fleets in any range when patrolling
const PatrolRangeInfinite = -1

// fleet not orbiting a planet
const NotOrbitingPlanet = 0

// no target planet, player, etc
const None = 0

type Fleet struct {
	MapObject
	FleetOrders
	PlanetNum         int         `json:"planetNum"` // for starbase fleets that are owned by a planet
	BaseName          string      `json:"baseName"`
	Cargo             Cargo       `json:"cargo,omitempty"`
	Fuel              int         `json:"fuel"`
	Tokens            []ShipToken `json:"tokens"`
	Heading           Vector      `json:"heading,omitempty"`
	WarpSpeed         int         `json:"warpSpeed,omitempty"`
	PreviousPosition  *Vector     `json:"previousPosition,omitempty"`
	OrbitingPlanetNum int         `json:"orbitingPlanetNum,omitempty"`
	Starbase          bool        `json:"starbase,omitempty"`
	Spec              FleetSpec   `json:"spec,omitempty"`
	battlePlan        *BattlePlan
}

type FleetOrders struct {
	Waypoints      []Waypoint `json:"waypoints"`
	RepeatOrders   bool       `json:"repeatOrders,omitempty"`
	BattlePlanName string     `json:"battlePlanName"`
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
	Position          Vector                 `json:"position,omitempty"`
	WarpFactor        int                    `json:"warpFactor,omitempty"`
	EstFuelUsage      int                    `json:"estFuelUsage,omitempty"`
	Task              WaypointTask           `json:"task,omitempty"`
	TransportTasks    WaypointTransportTasks `json:"transportTasks,omitempty"`
	WaitAtWaypoint    bool                   `json:"waitAtWaypoint,omitempty"`
	TargetType        MapObjectType          `json:"targetType,omitempty"`
	TargetNum         int                    `json:"targetNum,omitempty"`
	TargetPlayerNum   int                    `json:"targetPlayerNum,omitempty"`
	TransferToPlayer  int                    `json:"transferToPlayer,omitempty"`
	TargetName        string                 `json:"targetName,omitempty"`
	PartiallyComplete bool                   `json:"partiallyComplete,omitempty"`
	processed         bool                   `json:"-"`
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
	Fuel      WaypointTransportTask `json:"fuel,omitempty"`
	Ironium   WaypointTransportTask `json:"ironium,omitempty"`
	Boranium  WaypointTransportTask `json:"boranium,omitempty"`
	Germanium WaypointTransportTask `json:"germanium,omitempty"`
	Colonists WaypointTransportTask `json:"colonists,omitempty"`
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

// create a new fleet
func newFleet(player *Player, design *ShipDesign, num int, name string, waypoints []Waypoint) Fleet {
	return Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
			Dirty:     true,
			Num:       num,
			Name:      fmt.Sprintf("%s #%d", name, num),
			Position:  waypoints[0].Position,
		},
		BaseName: name,
		Tokens: []ShipToken{
			{design: design, DesignNum: design.Num, Quantity: 1},
		},
		FleetOrders: FleetOrders{
			Waypoints:      waypoints,
			BattlePlanName: player.BattlePlans[0].Name,
		},
		OrbitingPlanetNum: NotOrbitingPlanet,
		battlePlan:        &player.BattlePlans[0],
	}
}

func newFleetForToken(player *Player, num int, token ShipToken, waypoints []Waypoint) Fleet {
	return Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
			Dirty:     true,
			Num:       num,
			Name:      fmt.Sprintf("%s #%d", token.design.Name, num),
			Position:  waypoints[0].Position,
		},
		BaseName: token.design.Name,
		Tokens:   []ShipToken{token},
		FleetOrders: FleetOrders{
			Waypoints: waypoints,
		},
		OrbitingPlanetNum: NotOrbitingPlanet,
	}
}

// create a new fleet that is a starbase
func newStarbase(player *Player, planet *Planet, design *ShipDesign, name string) Fleet {
	fleet := newFleet(player, design, 0, name, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, 1)})
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

func (f *Fleet) withCargo(cargo Cargo) *Fleet {
	f.Cargo = cargo
	return f
}

func (f *Fleet) withPosition(position Vector) *Fleet {
	f.Position = position
	// todo: should we set waypoints in a builder?
	f.Waypoints = []Waypoint{{Position: position}}
	return f
}

func (f *Fleet) withWaypoints(waypoints []Waypoint) *Fleet {
	f.Waypoints = waypoints
	return f
}

func (f *Fleet) withOrbitingPlanetNum(num int) *Fleet {
	f.OrbitingPlanetNum = num
	return f
}

func (f *Fleet) Orbiting() bool {
	return f.OrbitingPlanetNum != NotOrbitingPlanet
}

func NewPlanetWaypoint(position Vector, num int, name string, warpFactor int) Waypoint {
	return Waypoint{
		Position:        position,
		TargetType:      MapObjectTypePlanet,
		TargetNum:       num,
		TargetName:      name,
		TargetPlayerNum: None,
		WarpFactor:      warpFactor,
	}
}

func NewFleetWaypoint(position Vector, num int, playerNum int, name string, warpFactor int) Waypoint {
	return Waypoint{
		Position:        position,
		TargetType:      MapObjectTypeFleet,
		TargetNum:       num,
		TargetPlayerNum: playerNum,
		TargetName:      name,
		WarpFactor:      warpFactor,
	}
}

func NewPositionWaypoint(position Vector, warpFactor int) Waypoint {
	return Waypoint{
		Position:        position,
		WarpFactor:      warpFactor,
		TargetNum:       None,
		TargetPlayerNum: None,
	}
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
func (f *Fleet) InjectDesigns(designs []*ShipDesign) {

	designsByNum := make(map[int]*ShipDesign, len(designs))
	for i := range designs {
		design := designs[i]
		designsByNum[design.Num] = design
	}

	// inject the design into this
	for i := range f.Tokens {
		token := &f.Tokens[i]
		token.design = designsByNum[token.DesignNum]
	}

}

func ComputeFleetSpec(rules *Rules, player *Player, fleet *Fleet) FleetSpec {
	spec := FleetSpec{
		ShipDesignSpec: ShipDesignSpec{
			ScanRange:    NoScanner,
			ScanRangePen: NoScanner,
			SpaceDock:    UnlimitedSpaceDock,
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

		// use the lowest ideal speed for this fleet
		// if we have multiple engines
		if token.design.Spec.Engine != "" {
			if spec.IdealSpeed == 0 {
				spec.IdealSpeed = token.design.Spec.IdealSpeed
			} else {
				spec.IdealSpeed = minInt(spec.IdealSpeed, token.design.Spec.IdealSpeed)
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
		spec.Shield += token.design.Spec.Shield * token.Quantity

		// cargo
		spec.CargoCapacity += token.design.Spec.CargoCapacity * token.Quantity

		// fuel
		spec.FuelCapacity += token.design.Spec.FuelCapacity * token.Quantity

		// minesweep
		spec.MineSweep += token.design.Spec.MineSweep * token.Quantity

		// remote mining
		spec.MiningRate += token.design.Spec.MiningRate * token.Quantity

		// remote terraforming
		spec.TerraformRate += token.design.Spec.TerraformRate * token.Quantity

		// colonization
		spec.Colonizer = spec.Colonizer || token.design.Spec.Colonizer
		spec.OrbitalConstructionModule = spec.OrbitalConstructionModule || token.design.Spec.OrbitalConstructionModule

		// spec all mine layers in the fleet
		if token.design.Spec.CanLayMines {
			for key := range token.design.Spec.MineLayingRateByMineType {
				if _, ok := spec.MineLayingRateByMineType[key]; ok {
					spec.MineLayingRateByMineType[key] = 0
				}
				spec.MineLayingRateByMineType[key] += token.design.Spec.MineLayingRateByMineType[key] * token.Quantity
			}
		}

		// We should only have one ship stack with spacdock capabilities, but for this logic just go with the max
		spec.SpaceDock = maxInt(spec.SpaceDock, token.design.Spec.SpaceDock)

		// sadly, the fleet only gets the best repair bonus from one design
		spec.RepairBonus = math.Max(spec.RepairBonus, token.design.Spec.RepairBonus)

		spec.ScanRange = maxInt(spec.ScanRange, token.design.Spec.ScanRange)
		spec.ScanRangePen = maxInt(spec.ScanRangePen, token.design.Spec.ScanRangePen)
		if token.design.Spec.Scanner {
			spec.Scanner = true
		}

		// add bombs
		if token.design.Spec.Bomber {
			spec.Bomber = true
			spec.Bombs = append(spec.Bombs, token.design.Spec.Bombs...)
			spec.SmartBombs = append(spec.SmartBombs, token.design.Spec.SmartBombs...)
			spec.RetroBombs = append(spec.RetroBombs, token.design.Spec.RetroBombs...)
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
		spec.ReduceCloaking = math.Max(spec.ReduceCloaking, token.design.Spec.ReduceCloaking)

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
	spec.CloakPercent = fleet.computeFleetCloakPercent(&spec, player.Race.Spec.FreeCargoCloaking)

	return spec
}

// compute fuel usage for each waypoint
func (f *Fleet) computeFuelUsage(player *Player) {
	for i := range f.Waypoints {
		wp := &f.Waypoints[i]
		if i > 0 {
			wpPrevious := f.Waypoints[i-1]
			fuelUsage := f.GetFuelCost(player, wp.WarpFactor, wp.Position.DistanceTo(wpPrevious.Position))
			wp.EstFuelUsage = fuelUsage
		} else {
			wp.EstFuelUsage = 0
		}
	}
}

// compute a fleet's cloak percent based on its current cargo/mass/cloak units
func (f *Fleet) computeFleetCloakPercent(spec *FleetSpec, freeCargoCloaking bool) int {
	cloakUnits := spec.CloakUnits

	// starbases have no mass or cargo, but fleet cloaking is adjusted for it
	if spec.Mass > 0 {
		// figure out how much cargo we are cloaking
		cloakedCargo := 0
		if !freeCargoCloaking {
			cloakedCargo = f.Cargo.Total()
		}

		cloakUnits = int(math.Round(float64(cloakUnits) * float64(spec.MassEmpty) / float64(spec.MassEmpty+cloakedCargo)))
	}
	return getCloakPercentForCloakUnits(cloakUnits)
}

func (f *Fleet) availableCargoSpace() int {
	return clamp(f.Spec.CargoCapacity-f.Cargo.Total(), 0, f.Spec.CargoCapacity)
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

func (fleet *Fleet) moveFleet(mapObjectGetter mapObjectGetter, rules *Rules, player *Player) {
	wp0 := fleet.Waypoints[0]
	wp1 := fleet.Waypoints[1]
	totalDist := fleet.Position.DistanceTo(wp1.Position)

	fleet.PreviousPosition = &Vector{fleet.Position.X, fleet.Position.Y}
	dist := float64(wp1.WarpFactor * wp1.WarpFactor)
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
	if player.Race.Spec.EngineFailureRate > 0 && wp1.WarpFactor > player.Race.Spec.EngineReliableSpeed && player.Race.Spec.EngineFailureRate >= rules.random.Float64() {
		messager.fleetEngineFailure(player, fleet)
		return
	}

	// get the cost for the fleet
	fuelCost := fleet.GetFuelCost(player, wp1.WarpFactor, dist)
	var fuelGenerated int = 0
	if fuelCost > fleet.Fuel {
		// we will run out of fuel
		// if this distance would have cost us 10 fuel but we have 6 left, only travel 60% of the distance.
		distanceFactor := float64(fleet.Fuel) / float64(fuelCost)
		dist = dist * distanceFactor

		// collide with minefields on route, but don't hit a minefield if we run out of fuel beforehand
		// TODO: add back in minefield check
		// dist = CheckForMineFields(fleet, player, wp1, dist)

		fleet.Fuel = 0
		wp1.WarpFactor = fleet.getNoFuelWarpFactor(rules.techs, player)
		messager.fleetOutOfFuel(player, fleet, wp1.WarpFactor)

		// if we ran out of fuel 60% of the way to our normal distance, the remaining 40% of our time
		// was spent travelling at fuel generation speeds:
		remainingDistanceTravelled := (1 - distanceFactor) * float64(wp1.WarpFactor*wp1.WarpFactor)
		dist += remainingDistanceTravelled
		fuelGenerated = fleet.getFuelGeneration(rules.techs, player, wp1.WarpFactor, remainingDistanceTravelled)
	} else {
		// collide with minefields on route, but don't hit a minefield if we run out of fuel beforehand
		// TODO: add back in minefield check
		// actualDist := CheckForMineFields(fleet, player, wp1, dist)
		actualDist := dist
		if actualDist != dist {
			dist = actualDist
			fuelCost = fleet.GetFuelCost(player, wp1.WarpFactor, dist)
			// we hit a minefield, update fuel usage
		}

		fleet.Fuel -= fuelCost
		fuelGenerated = fleet.getFuelGeneration(rules.techs, player, wp1.WarpFactor, dist)
	}

	// message the player about fuel generation
	fuelGenerated = minInt(fuelGenerated, fleet.Spec.FuelCapacity-fleet.Fuel)
	if fuelGenerated > 0 {
		fleet.Fuel += fuelGenerated
		messager.fleetGeneratedFuel(player, fleet, fuelGenerated)
	}

	// assuming we move at all, make sure we are no longer orbiting any planets
	if dist > 0 && fleet.Orbiting() {
		fleet.OrbitingPlanetNum = NotOrbitingPlanet
	}

	// TODO: repeat orders, can we just append wp0 to waypoints when we repeat?
	// if wp0.OriginalTarget == nil || !wp0.OriginalPosition.HasValue {
	// 	wp0.OriginalTarget = wp0.Target
	// 	wp0.OriginalPosition = fleet.Position
	// }

	if totalDist == dist {
		fleet.completeMove(mapObjectGetter, wp0, wp1)
	} else {
		// move this fleet closer to the next waypoint
		fleet.WarpSpeed = wp1.WarpFactor
		fleet.Heading = (wp1.Position.Subtract(fleet.Position)).Normalized()
		wp0.TargetType = MapObjectTypeNone
		wp0.TargetNum = None
		wp0.TargetPlayerNum = None
		wp0.TargetName = ""
		wp0.PartiallyComplete = true

		fleet.Position = fleet.Position.Add(fleet.Heading.Scale(dist))
		fleet.Position = fleet.Position.Round()
		wp0.Position = fleet.Position
		fleet.Waypoints[0] = wp0
	}
}

// GateFleet moves the fleet the cool way, with stargates!
func (fleet *Fleet) gateFleet(mapObjectGetter mapObjectGetter, playerGetter playerGetter, rules *Rules, player *Player) {
	wp0 := fleet.Waypoints[0]
	wp1 := fleet.Waypoints[1]
	totalDist := fleet.Position.DistanceTo(wp1.Position)

	// if we got here, both source and dest have stargates
	sourcePlanet := mapObjectGetter.getPlanet(fleet.OrbitingPlanetNum)
	destPlanet := mapObjectGetter.getPlanet(wp1.TargetNum)

	if sourcePlanet == nil || !sourcePlanet.Spec.HasStargate {
		messager.fleetStargateInvalidSource(player, fleet, wp0)
		return
	}
	if destPlanet == nil || !destPlanet.Spec.HasStargate {
		messager.fleetStargateInvalidDest(player, fleet, wp0, wp1)
		return
	}

	sourcePlanetPlayer := playerGetter.getPlayer(sourcePlanet.PlayerNum)
	destPlanetPlayer := playerGetter.getPlayer(destPlanet.PlayerNum)

	if !sourcePlanetPlayer.IsFriend(player.Num) {
		messager.fleetStargateInvalidSourceOwner(player, fleet, wp0, wp1)
		return
	}
	if !destPlanetPlayer.IsFriend(player.Num) {
		messager.fleetStargateInvalidDestOwner(player, fleet, wp0, wp1)
		return
	}
	if fleet.Cargo.Colonists > 0 && !sourcePlanet.OwnedBy(player.Num) {
		messager.fleetStargateInvalidColonists(player, fleet, wp0, wp1)
		return
	}

	sourceStargate := sourcePlanet.Spec
	destStargate := destPlanet.Spec

	// only the source gate matters for range
	minSafeRange := sourceStargate.SafeRange
	minSafeHullMass := minInt(sourceStargate.SafeHullMass, destStargate.SafeHullMass)

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

	// dump cargo if we aren't IT
	if fleet.Cargo.Total() > 0 && !player.Race.Spec.CanGateCargo {
		messager.fleetStargateDumpedCargo(player, fleet, wp0, wp1, fleet.Cargo)
		sourcePlanet.Cargo = sourcePlanet.Cargo.Add(fleet.Cargo)
		fleet.Cargo = Cargo{}
	}

	// apply overgate damage and delete tokens (and possibly the fleet)
	// also vanish tokens for non IT races
	fleet.ApplyOvergatePenalty(player, rules, totalDist, wp0, wp1, sourceStargate, destStargate)

	if len(fleet.Tokens) > 0 {
		// if we survived, warp it!
		fleet.completeMove(mapObjectGetter, wp0, wp1)
	}
}

// ApplyOvergatePenalty applies damage (if any) to each token that overgated
func (fleet *Fleet) ApplyOvergatePenalty(player *Player, rules *Rules, distance float64, wp0, wp1 Waypoint, sourceStargate, destStargate PlanetSpec) {
	var totalDamage, shipsLostToDamage, shipsLostToTheVoid, startingShips int
	for _, token := range fleet.Tokens {
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
							// get rid of the damaged ships first and redistribute the damage
							// i.e. if we have 2 damaged ships with 20 total damage
							// we get rid of one of them and leave one with 10 damage
							token.Damage = math.Max(0, token.Damage/float64(token.QuantityDamaged))
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
	var newTokens []ShipToken
	for _, token := range fleet.Tokens {
		if token.Quantity > 0 {
			newTokens = append(newTokens, token)
		}
	}
	fleet.Tokens = newTokens

	if len(fleet.Tokens) == 0 {
		// gm.eventManager.PublishMapObjectDeletedEvent(fleet)
		messager.fleetStargateDestroyed(player, fleet, wp0, wp1)
	} else {
		if totalDamage > 0 || shipsLostToTheVoid > 0 {
			messager.fleetStargateDamaged(player, fleet, wp0, wp1, totalDamage, startingShips, shipsLostToDamage, shipsLostToTheVoid)
		}
	}
}

// Engine fuel usage calculation courtesy of m.a@stars
func (fleet *Fleet) getFuelCostForEngine(warpFactor int, mass int, dist float64, ifeFactor float64, fuelUsage [11]int) int {
	if warpFactor == 0 {
		return 0
	}
	// 1 mg of fuel will move 200kT of weight 1 LY at a Fuel Usage Number of 100.
	// Number of engines doesn't matter. Neither number of ships with the same engine.

	distanceCeiling := math.Ceil(dist) // rounding to next integer gives best graph fit
	// window.status = 'Actual distance used is ' + Distan + 'ly';

	// IFE is applied to drive specifications, just as the helpfile hints.
	// Stars! probably does it outside here once per turn per engine to save time.
	engineEfficiency := math.Ceil(ifeFactor * float64(fuelUsage[warpFactor]))

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
func (fleet *Fleet) GetFuelCost(player *Player, warpFactor int, distance float64) int {

	// figure out how much fuel we're going to use
	efficiencyFactor := 1 - player.Race.Spec.FuelEfficiencyOffset

	var fuelCost int = 0

	// compute each ship stack separately
	for _, token := range fleet.Tokens {
		// figure out this ship stack's mass as well as it's proportion of the cargo
		mass := token.design.Spec.Mass * token.Quantity
		fleetCargo := fleet.Cargo.Total()
		stackCapacity := token.design.Spec.CargoCapacity * token.Quantity
		fleetCapacity := fleet.Spec.CargoCapacity

		if fleetCapacity > 0 {
			mass += int(float64(fleetCargo) * (float64(stackCapacity) / float64(fleetCapacity)))
		}

		fuelCost += fleet.getFuelCostForEngine(warpFactor, mass, distance, efficiencyFactor, token.design.Spec.FuelUsage)
	}

	return fuelCost
}

// Get the warp factor for when we run out of fuel.
func (fleet *Fleet) getNoFuelWarpFactor(techStore *TechStore, player *Player) int {
	// find the lowest freeSpeed from all the fleet's engines
	var freeSpeed = math.MaxInt
	for _, token := range fleet.Tokens {
		engine := techStore.GetEngine(token.design.Spec.Engine)
		freeSpeed = minInt(freeSpeed, engine.FreeSpeed)
	}
	return freeSpeed
}

// Get the warp factor for when we run out of fuel.
func (fleet *Fleet) getFuelGeneration(techStore *TechStore, player *Player, warpFactor int, distance float64) int {
	// find the lowest freeSpeed from all the fleet's engines
	var freeSpeed = math.MaxInt
	for _, token := range fleet.Tokens {
		engine := techStore.GetEngine(token.design.Spec.Engine)
		freeSpeed = minInt(freeSpeed, engine.FreeSpeed)
	}
	return freeSpeed
}

// Complete a move from one waypoint to another
func (fleet *Fleet) completeMove(mapObjectGetter mapObjectGetter, wp0 Waypoint, wp1 Waypoint) {
	fleet.Position = wp1.Position

	// find out if we arrived at a planet, either by reaching our target fleet
	// or reaching a planet
	if wp1.TargetType == MapObjectTypeFleet && wp1.TargetPlayerNum != None && wp1.TargetNum != None {
		target := mapObjectGetter.getFleet(wp1.TargetPlayerNum, wp1.TargetNum)
		fleet.OrbitingPlanetNum = target.OrbitingPlanetNum

		// we are orbiting a friendly planet
		targetPlanet := mapObjectGetter.getPlanet(fleet.OrbitingPlanetNum)
		if fleet.PlayerNum == targetPlanet.PlayerNum && targetPlanet.Spec.HasStarbase {
			// refuel at starbases
			fleet.Fuel = fleet.Spec.FuelCapacity
		}
	} else if wp1.TargetType == MapObjectTypePlanet && wp1.TargetNum != None {
		target := mapObjectGetter.getPlanet(wp1.TargetNum)
		fleet.OrbitingPlanetNum = target.Num
		if fleet.PlayerNum == target.PlayerNum && target.Spec.HasStarbase {
			// refuel at starbases
			fleet.Fuel = fleet.Spec.FuelCapacity
		}
	} else if wp1.TargetType == MapObjectTypeWormhole && wp1.TargetNum != None {
		target := mapObjectGetter.getWormhole(wp1.TargetNum)
		dest := mapObjectGetter.getWormhole(target.DestinationNum)
		fleet.Position = dest.Position
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
		fleet.WarpSpeed = wp1.WarpFactor
		fleet.Heading = (wp1.Position.Subtract(fleet.Position)).Normalized()
	}
}

// colonize a planet
func (fleet *Fleet) colonizePlanet(rules *Rules, player *Player, planet *Planet) {
	planet.Dirty = true
	planet.PlayerNum = player.Num
	planet.ProductionQueue = []ProductionQueueItem{}
	planet.Cargo = planet.Cargo.Add(fleet.Cargo)

	if len(player.ProductionPlans) > 0 {
		// TODO: apply production plan
		plan := player.ProductionPlans[0]
		planet.ContributesOnlyLeftoverToResearch = plan.ContributesOnlyLeftoverToResearch
	}

	if player.Race.Spec.InnateMining {
		planet.Mines = planet.innateMines(player)
	}

	if fleet.Spec.OrbitalConstructionModule {
		design := player.GetLatestDesign(ShipDesignPurposeStarterColony)
		if design != nil {
			starbase := newStarbase(player, planet, design, design.Name)
			starbase.Spec = ComputeFleetSpec(rules, player, &starbase)
			planet.starbase = &starbase
		}
	}

	planet.Spec = computePlanetSpec(rules, player, planet)
}

// get the minerals and resources recovered from a scrapped fleet
func (fleet *Fleet) getScrapAmount(rules *Rules, player *Player, planet *Planet) Cost {

	// create a new cargo instance out of our fleet cost
	scrappedCost := fleet.Spec.Cost

	scrapMineralFactor := rules.ScrapMineralAmount
	scrapResourceFactor := rules.ScrapResourceAmount
	extraResources := 0
	planetResources := 0

	if planet != nil && planet.OwnedBy(player.Num) {
		planetResources = planet.Spec.ResourcesPerYear + planet.BonusResources
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
	}

	return scrappedCost
}

// getTransferAmount gets the amount of cargo to transfer for loading a cargo type from a cargoholder
func (fleet *Fleet) getCargoLoadAmount(dest cargoHolder, cargoType CargoType, task WaypointTransportTask) (transferAmount int, waitAtWaypoint bool) {
	capacity := fleet.Spec.CargoCapacity - fleet.Cargo.Total()
	availableToLoad := dest.getCargo().GetAmount(cargoType)
	currentAmount := fleet.Cargo.GetAmount(cargoType)
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
				transferAmount = maxInt(-leftoverFuel, -(dest.getFuelCapacity() - dest.getFuel()))
			}
		}
	case TransportActionLoadAll:
		// load all available, based on our constraints
		transferAmount = minInt(availableToLoad, capacity)
	case TransportActionLoadAmount:
		transferAmount = minInt(minInt(availableToLoad, task.Amount), capacity)
	case TransportActionWaitForPercent:
		fallthrough
	case TransportActionFillPercent:
		// we want a percent of our hold to be filled with some amount, figure out how
		// much that is in kT, i.e. 50% of 100kT would be 50kT of this mineral
		var taskAmountkT = int(float64(task.Amount) / 100 * float64(capacity))

		if currentAmount >= taskAmountkT {
			// no need to transfer any, move on
			return 0, false
		} else {

			transferAmount = minInt(minInt(availableToLoad, taskAmountkT-currentAmount), capacity)
			if transferAmount < taskAmountkT && task.Action == TransportActionWaitForPercent {
				waitAtWaypoint = true
			}
		}
	case TransportActionSetAmountTo:
		// only transfer the min of what we have, vs what we need, vs the capacity
		transferAmount = minInt(minInt(availableToLoad, task.Amount-currentAmount), capacity)
		if transferAmount < (task.Amount - currentAmount) {
			waitAtWaypoint = true
		}
	case TransportActionLoadDunnage:
		transferAmount = fleet.getCargoLoadDunnageAmount(dest, cargoType)
	}

	// let the caller know how much of this cargo we load
	return transferAmount, waitAtWaypoint
}

func (fleet *Fleet) getCargoLoadDunnageAmount(dest cargoHolder, cargoType CargoType) (transferAmount int) {
	capacity := fleet.Spec.CargoCapacity - fleet.Cargo.Total()
	availableToLoad := dest.getCargo().GetAmount(cargoType)

	// (minerals and colonists only) This command waits until all other loads and unloads are
	// complete, then loads as many colonists or amount of a mineral as will fit in the remaining
	// space. For example, setting Load All Germanium, Load Dunnage Ironium, will load all the
	// Germanium that is available, then as much Ironium as possible. If more than one dunnage cargo
	// is specified, they are loaded in the order of Ironium, Boranium, Germanium, and Colonists.
	return minInt(availableToLoad, capacity)
}

// getCargoUnloadAmount gets the amount of cargo to transfer for unloading a cargo type from a cargoholder
func (fleet *Fleet) getCargoUnloadAmount(dest cargoHolder, cargoType CargoType, task WaypointTransportTask) (transferAmount int, waitAtWaypoint bool) {

	var availableToUnload int
	if cargoType == Fuel {

		availableToUnload = fleet.Fuel
	} else {
		availableToUnload = fleet.Cargo.GetAmount(cargoType)
	}
	currentAmount := fleet.Cargo.GetAmount(cargoType)
	_ = currentAmount
	switch task.Action {
	case TransportActionUnloadAll:
		// unload all available, based on our constraints
		transferAmount = minInt(availableToUnload, dest.getCargoCapacity())
	case TransportActionUnloadAmount:
		// don't unload more than the task says
		transferAmount = minInt(minInt(availableToUnload, task.Amount), dest.getCargoCapacity())
	case TransportActionSetWaypointTo:
		// Make sure the waypoint has at least whatever we specified
		var currentAmount = dest.getCargo().GetAmount(cargoType)

		if currentAmount >= task.Amount {
			// no need to transfer any, move on
			break
		} else {
			// only transfer the min of what we have, vs what we need, vs the capacity
			transferAmount = minInt(minInt(availableToUnload, task.Amount-currentAmount), dest.getCargoCapacity())
			if transferAmount < task.Amount {
				waitAtWaypoint = true
			}
		}
	}
	return transferAmount, waitAtWaypoint
}
