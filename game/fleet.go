package game

import (
	"fmt"
	"math"
	"time"
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
const NotOrbitingPlanet = -1

// no target planet, player, etc
const NoTarget = -1

type Fleet struct {
	MapObject
	FleetOrders
	PlanetID          uint64      `json:"-"` // for starbase fleets that are owned by a planet
	BaseName          string      `json:"baseName"`
	Cargo             Cargo       `json:"cargo,omitempty" gorm:"embedded;embeddedPrefix:cargo_"`
	Fuel              int         `json:"fuel"`
	Damage            int         `json:"damage"`
	BattlePlanID      uint64      `json:"battlePlan"`
	Tokens            []ShipToken `json:"tokens"`
	Heading           Vector      `json:"heading,omitempty" gorm:"embedded;embeddedPrefix:heading_"`
	WarpSpeed         int         `json:"warpSpeed,omitempty"`
	PreviousPosition  *Vector     `json:"previousPosition,omitempty" gorm:"embedded;embeddedPrefix:previous_position_"`
	OrbitingPlanetNum int         `json:"orbitingPlanetNum,omitempty"`
	Starbase          bool        `json:"starbase,omitempty"`
	Spec              *FleetSpec  `json:"spec" gorm:"serializer:json"`
}

type FleetOrders struct {
	Waypoints    []Waypoint `json:"waypoints" gorm:"serializer:json"`
	RepeatOrders bool       `json:"repeatOrders,omitempty"`
}

type FleetSpec struct {
	ShipDesignSpec
	Purposes         map[ShipDesignPurpose]bool `json:"purposes"`
	TotalShips       int                        `json:"totalShips"`
	MassEmpty        int                        `json:"massEmpty"`
	BasePacketSpeed  int                        `json:"basePacketSpeed"`
	SafePacketSpeed  int                        `json:"safePacketSpeed"`
	BaseCloakedCargo int                        `json:"baseCloakedCargo"`
	HasMassDriver    bool                       `json:"hasMassDriver,omitempty"`
	HasStargate      bool                       `json:"hasStargate,omitempty"`
	Stargate         string                     `json:"stargate,omitempty"`
}

type ShipToken struct {
	ID        uint64      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	FleetID   uint64      `json:"gameId"`
	DesignID  uint64      `json:"designId"`
	Quantity  int         `json:"quantity"`
	Design    *ShipDesign `json:"-" gorm:"foreignKey:DesignID"`
}

type Waypoint struct {
	Position          Vector                 `json:"position,omitempty" gorm:"embedded"`
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
func NewFleet(player *Player, design *ShipDesign, num int, name string, waypoints []Waypoint) Fleet {
	return Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			GameID:    player.GameID,
			PlayerID:  player.ID,
			PlayerNum: player.Num,
			Dirty:     true,
			Num:       num,
			Name:      fmt.Sprintf("%s #%d", name, num),
			Position:  waypoints[0].Position,
		},
		BaseName: name,
		Tokens: []ShipToken{
			{Design: design, Quantity: 1},
		},
		FleetOrders: FleetOrders{
			Waypoints: waypoints,
		},
		OrbitingPlanetNum: NotOrbitingPlanet,
	}
}

func NewFleetForToken(player *Player, num int, token ShipToken, waypoints []Waypoint) Fleet {
	return Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			GameID:    player.GameID,
			PlayerID:  player.ID,
			PlayerNum: player.Num,
			Dirty:     true,
			Num:       num,
			Name:      fmt.Sprintf("%s #%d", token.Design.Name, num),
			Position:  waypoints[0].Position,
		},
		BaseName: token.Design.Name,
		Tokens:   []ShipToken{token},
		FleetOrders: FleetOrders{
			Waypoints: waypoints,
		},
		OrbitingPlanetNum: NotOrbitingPlanet,
	}
}

// create a new fleet that is a starbase
func NewStarbase(player *Player, planet *Planet, design *ShipDesign, name string) Fleet {
	fleet := NewFleet(player, design, 0, name, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, 1)})
	fleet.PlanetID = planet.ID
	fleet.Starbase = true

	return fleet
}

func (f *Fleet) String() string {
	return fmt.Sprintf("Fleet %s #%d", f.BaseName, f.Num)
}

func (f *Fleet) WithPlayerNum(playerNum int) *Fleet {
	f.PlayerNum = playerNum
	return f
}

func (f *Fleet) WithCargo(cargo Cargo) *Fleet {
	f.Cargo = cargo
	return f
}

func (f *Fleet) WithPosition(position Vector) *Fleet {
	f.Position = position
	// todo: should we set waypoints in a builder?
	f.Waypoints = []Waypoint{{Position: position}}
	return f
}

func (f *Fleet) WithWaypoints(waypoints []Waypoint) *Fleet {
	f.Waypoints = waypoints
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
		TargetPlayerNum: NoTarget,
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
		TargetNum:       NoTarget,
		TargetPlayerNum: NoTarget,
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

func ComputeFleetSpec(rules *Rules, player *Player, fleet *Fleet) *FleetSpec {
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

		if token.Design.Purpose != ShipDesignPurposeNone {
			spec.Purposes[token.Design.Purpose] = true
		}

		// use the lowest ideal speed for this fleet
		// if we have multiple engines
		if token.Design.Spec.Engine != "" {
			if spec.IdealSpeed == 0 {
				spec.IdealSpeed = token.Design.Spec.IdealSpeed
			} else {
				spec.IdealSpeed = MinInt(spec.IdealSpeed, token.Design.Spec.IdealSpeed)
			}
		}
		// cost
		spec.Cost = token.Design.Spec.Cost.MultiplyInt(token.Quantity)

		// mass
		spec.Mass += token.Design.Spec.Mass * token.Quantity
		spec.MassEmpty += token.Design.Spec.Mass * token.Quantity

		// armor
		spec.Armor += token.Design.Spec.Armor * token.Quantity

		// shield
		spec.Shield += token.Design.Spec.Shield * token.Quantity

		// cargo
		spec.CargoCapacity += token.Design.Spec.CargoCapacity * token.Quantity

		// fuel
		spec.FuelCapacity += token.Design.Spec.FuelCapacity * token.Quantity

		// minesweep
		spec.MineSweep += token.Design.Spec.MineSweep * token.Quantity

		// remote mining
		spec.MiningRate += token.Design.Spec.MiningRate * token.Quantity

		// remote terraforming
		spec.TerraformRate += token.Design.Spec.TerraformRate * token.Quantity

		// colonization
		spec.Colonizer = spec.Colonizer || token.Design.Spec.Colonizer
		spec.OrbitalConstructionModule = spec.OrbitalConstructionModule || token.Design.Spec.OrbitalConstructionModule

		// spec all mine layers in the fleet
		if token.Design.Spec.CanLayMines {
			for key := range token.Design.Spec.MineLayingRateByMineType {
				if _, ok := spec.MineLayingRateByMineType[key]; ok {
					spec.MineLayingRateByMineType[key] = 0
				}
				spec.MineLayingRateByMineType[key] += token.Design.Spec.MineLayingRateByMineType[key] * token.Quantity
			}
		}

		// We should only have one ship stack with spacdock capabilities, but for this logic just go with the max
		spec.SpaceDock = MaxInt(spec.SpaceDock, token.Design.Spec.SpaceDock)

		// sadly, the fleet only gets the best repair bonus from one design
		spec.RepairBonus = math.Max(spec.RepairBonus, token.Design.Spec.RepairBonus)

		spec.ScanRange = MaxInt(spec.ScanRange, token.Design.Spec.ScanRange)
		spec.ScanRangePen = MaxInt(spec.ScanRangePen, token.Design.Spec.ScanRangePen)
		if token.Design.Spec.Scanner {
			spec.Scanner = true
		}

		// add bombs
		if token.Design.Spec.Bomber {
			spec.Bomber = true
			spec.Bombs = append(spec.Bombs, token.Design.Spec.Bombs...)
			spec.SmartBombs = append(spec.SmartBombs, token.Design.Spec.SmartBombs...)
			spec.RetroBombs = append(spec.RetroBombs, token.Design.Spec.RetroBombs...)
		}

		// check if any tokens have weapons
		// we process weapon slots per stack, so we don't need to spec all
		// weapons in a fleet
		if token.Design.Spec.HasWeapons {
			spec.HasWeapons = true
		}

		if token.Design.Spec.CloakUnits > 0 {
			// calculate the cloak units for this token based on the design's cloak units (i.e. 70 cloak units / kT for a stealh cloak)
			spec.CloakUnits += token.Design.Spec.CloakUnits
		} else {
			// if this ship doesn't have cloaking, it counts as cargo (except for races with free cargo cloaking)
			if !player.Race.Spec.FreeCargoCloaking {
				spec.BaseCloakedCargo += token.Design.Spec.Mass * token.Quantity
			}
		}

		// choose the best tachyon detector ship
		spec.ReduceCloaking = math.Max(spec.ReduceCloaking, token.Design.Spec.ReduceCloaking)

		spec.CanStealFleetCargo = spec.CanStealFleetCargo || token.Design.Spec.CanStealFleetCargo
		spec.CanStealPlanetCargo = spec.CanStealPlanetCargo || token.Design.Spec.CanStealPlanetCargo
	}

	// compute the cloaking based on the cloak units and cargo
	spec.CloakPercent = fleet.computeFleetCloakPercent(&spec, player.Race.Spec.FreeCargoCloaking)

	return &spec
}

func (f *Fleet) ComputeFuelUsage(player *Player) {
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

func (f *Fleet) AvailableCargoSpace() int {
	return clamp(f.Spec.CargoCapacity-f.Cargo.Total(), 0, f.Spec.CargoCapacity)
}

// transfer cargo from a fleet to a cargo holder
func (f *Fleet) TransferCargoItem(dest CargoHolder, cargoType CargoType, transferAmount int) error {
	destCargo := dest.GetCargo()
	if f.AvailableCargoSpace() < transferAmount {
		return fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", f.Name, f.AvailableCargoSpace(), transferAmount, dest.GetMapObject().Name)
	}

	if !destCargo.CanTransferAmount(cargoType, transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %v from %s, there is not enough to transfer", f.Name, transferAmount, dest.GetMapObject().Name)
	}

	// transfer the cargo
	f.Cargo = f.Cargo.AddAmount(cargoType, transferAmount)

	switch cargoType {
	case Ironium:
		destCargo.Ironium -= transferAmount
	case Boranium:
		destCargo.Boranium -= transferAmount
	case Germanium:
		destCargo.Germanium -= transferAmount
	case Colonists:
		destCargo.Colonists -= transferAmount
	}

	return nil
}

// transfer cargo from a fleet to a cargo holder
func (f *Fleet) TransferCargo(dest CargoHolder, transferAmount Cargo) error {
	destCargo := dest.GetCargo()
	if f.AvailableCargoSpace() < transferAmount.Total() {
		return fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", f.Name, f.AvailableCargoSpace(), transferAmount.Total(), dest.GetMapObject().Name)
	}

	if !destCargo.CanTransfer(transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %v from %s, there is not enough to transfer", f.Name, transferAmount, dest.GetMapObject().Name)
	}

	// transfer the cargo
	f.Cargo = f.Cargo.Add(transferAmount)

	// update the dest cargo. It's a pointer to a cargo, so
	// we need to update each value
	updatedDestCargo := destCargo.Subtract(transferAmount)
	destCargo.Ironium = updatedDestCargo.Ironium
	destCargo.Boranium = updatedDestCargo.Boranium
	destCargo.Germanium = updatedDestCargo.Germanium
	destCargo.Colonists = updatedDestCargo.Colonists

	return nil
}

// transfer cargo from a planet to/from a fleet
func (f *Fleet) TransferPlanetCargo(planet *Planet, transferAmount Cargo) error {

	if f.AvailableCargoSpace() < transferAmount.Total() {
		return fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", f.Name, f.AvailableCargoSpace(), transferAmount.Total(), planet.Name)
	}

	if !planet.Cargo.CanTransfer(transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %v from %s, there is not enough to transfer", f.Name, transferAmount, planet.Name)
	}

	// transfer the cargo
	f.Cargo = f.Cargo.Add(transferAmount)
	planet.Cargo = planet.Cargo.Subtract(transferAmount)

	return nil
}

// transfer cargo from a fleet to/from a fleet
func (f *Fleet) TransferFleetCargo(fleet *Fleet, transferAmount Cargo) error {

	if f.AvailableCargoSpace() < transferAmount.Total() {
		return fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", f.Name, f.AvailableCargoSpace(), transferAmount.Total(), fleet.Name)
	}

	if !fleet.Cargo.CanTransfer(transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %v from %s, there is not enough to transfer", f.Name, transferAmount, fleet.Name)
	}

	// transfer the cargo
	f.Cargo = f.Cargo.Add(transferAmount)
	fleet.Cargo = fleet.Cargo.Subtract(transferAmount)

	return nil
}

func (fleet *Fleet) moveFleet(mapObjectGetter MapObjectGetter, rules *Rules, player *Player) {
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
	fuelGenerated = MinInt(fuelGenerated, fleet.Spec.FuelCapacity-fleet.Fuel)
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
		wp0.TargetNum = NoTarget
		wp0.TargetPlayerNum = NoTarget
		wp0.TargetName = ""
		wp0.PartiallyComplete = true

		fleet.Position = fleet.Position.Add(fleet.Heading.Scale(dist))
		fleet.Position = fleet.Position.Round()
		wp0.Position = fleet.Position
		fleet.Waypoints[0] = wp0
	}
}

func (fleet *Fleet) gateFleet(rules *Rules, player *Player) {
	panic("unimplemented")
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
		mass := token.Design.Spec.Mass * token.Quantity
		fleetCargo := fleet.Cargo.Total()
		stackCapacity := token.Design.Spec.CargoCapacity * token.Quantity
		fleetCapacity := fleet.Spec.CargoCapacity

		if fleetCapacity > 0 {
			mass += int(float64(fleetCargo) * (float64(stackCapacity) / float64(fleetCapacity)))
		}

		fuelCost += fleet.getFuelCostForEngine(warpFactor, mass, distance, efficiencyFactor, token.Design.Spec.FuelUsage)
	}

	return fuelCost
}

// Get the warp factor for when we run out of fuel.
func (fleet *Fleet) getNoFuelWarpFactor(techStore *TechStore, player *Player) int {
	// find the lowest freeSpeed from all the fleet's engines
	var freeSpeed = math.MaxInt
	for _, token := range fleet.Tokens {
		engine := techStore.GetEngine(token.Design.Spec.Engine)
		freeSpeed = MinInt(freeSpeed, engine.FreeSpeed)
	}
	return freeSpeed
}

// Get the warp factor for when we run out of fuel.
func (fleet *Fleet) getFuelGeneration(techStore *TechStore, player *Player, warpFactor int, distance float64) int {
	// find the lowest freeSpeed from all the fleet's engines
	var freeSpeed = math.MaxInt
	for _, token := range fleet.Tokens {
		engine := techStore.GetEngine(token.Design.Spec.Engine)
		freeSpeed = MinInt(freeSpeed, engine.FreeSpeed)
	}
	return freeSpeed
}

// Complete a move from one waypoint to another
func (fleet *Fleet) completeMove(mapObjectGetter MapObjectGetter, wp0 Waypoint, wp1 Waypoint) {
	fleet.Position = wp1.Position

	// find out if we arrived at a planet, either by reaching our target fleet
	// or reaching a planet
	if wp1.TargetType == MapObjectTypeFleet && wp1.TargetPlayerNum != NoTarget && wp1.TargetNum != NoTarget {
		target := mapObjectGetter.GetFleet(wp1.TargetPlayerNum, wp1.TargetNum)
		fleet.OrbitingPlanetNum = target.OrbitingPlanetNum

		// we are orbiting a friendly planet
		targetPlanet := mapObjectGetter.GetPlanet(fleet.OrbitingPlanetNum)
		if fleet.PlayerNum == targetPlanet.PlayerNum && targetPlanet.Spec.HasStarbase {
			// refuel at starbases
			fleet.Fuel = fleet.Spec.FuelCapacity
		}
	} else if wp1.TargetType == MapObjectTypePlanet && wp1.TargetNum != NoTarget {
		target := mapObjectGetter.GetPlanet(wp1.TargetNum)
		fleet.OrbitingPlanetNum = target.Num
		if fleet.PlayerNum == target.PlayerNum && target.Spec.HasStarbase {
			// refuel at starbases
			fleet.Fuel = fleet.Spec.FuelCapacity
		}
	} else if wp1.TargetType == MapObjectTypeWormhole && wp1.TargetNum != NoTarget {
		target := mapObjectGetter.GetWormhole(wp1.TargetNum)
		dest := mapObjectGetter.GetWormhole(target.DestinationNum)
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
	planet.PlayerID = player.ID
	planet.ProductionQueue = []ProductionQueueItem{}
	planet.Cargo = planet.Cargo.Add(fleet.Cargo)

	if len(player.ProductionPlans) > 0 {
		// TODO: apply production plan
		plan := player.ProductionPlans[0]
		planet.ContributesOnlyLeftoverToResearch = plan.ContributesOnlyLeftoverToResearch
	}

	if player.Race.Spec.InnateMining {
		planet.Mines = planet.GetInnateMines(player)
	}

	if fleet.Spec.OrbitalConstructionModule {
		design := player.GetLatestDesign(ShipDesignPurposeStarterColony)
		if design != nil {
			starbase := NewStarbase(player, planet, design, design.Name)
			starbase.Spec = ComputeFleetSpec(rules, player, &starbase)
			planet.Starbase = &starbase
		}
	}

	planet.Spec = ComputePlanetSpec(rules, planet, player)
}

// scrap a fleet giving a planet resources or returning salvage
func (fleet *Fleet) scrap(rules *Rules, player *Player, planet *Planet) *Salvage {
	cost := fleet.getScrapAmount(rules, player, planet)

	if planet != nil {
		// scrap over a planet
		planet.Cargo = planet.Cargo.AddCostMinerals(cost)
		if planet.OwnedBy(player.Num) {
			planet.BonusResources += cost.Resources
		}
	} else {
		// create salvage
		salvage := NewSalvage(player.Num, fleet.Position, cost.ToCargo())
		return &salvage
	}

	fleet.Delete = true
	return nil
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
