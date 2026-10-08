package cs

import (
	"fmt"
	"log/slog"
)

type cargoTransferer struct {
	log     *slog.Logger
	invader invader
	game    *FullGame
}

// CargoTransferStatus will alert the user if a CargoTransfer didn't go through due to insufficient capacity or available cargo
type CargoTransferStatus int

const (
	CargoTransferStatusNone CargoTransferStatus = iota
	CargoTransferStatusOwned
	CargoTransferStatusCargo
	CargoTransferStatusCargoCapacity
	CargoTransferStatusDestCargo
	CargoTransferStatusDestCargoCapacity
	// if a starbase is present, you cannot drop invaders
	CargoTransferStatusDestStarbase
	CargoTransferStatusDestUnowned
	// colonists can't survive in deep space or salvage
	CargoTransferStatusDeepSpace
)

func (r CargoTransferStatus) String() string {
	switch r {
	case CargoTransferStatusNone:
		return "None"
	case CargoTransferStatusOwned:
		return "Owned"
	case CargoTransferStatusCargo:
		return "Insufficient Cargo"
	case CargoTransferStatusCargoCapacity:
		return "Insufficient Cargo Capacity"
	case CargoTransferStatusDestCargo:
		return "Insufficient Destination Cargo"
	case CargoTransferStatusDestCargoCapacity:
		return "Insufficient Destination Cargo Capacity"
	case CargoTransferStatusDestStarbase:
		return "Destination Has Starbase"
	case CargoTransferStatusDestUnowned:
		return "Destination Unowned"
	case CargoTransferStatusDeepSpace:
		return "Deep Space"
	default:
		return fmt.Sprintf("Unknown %d", r)
	}
}

// cargoTransferResult is the result of a single CargoType cargo transfer to a dest
type cargoTransferResult struct {
	status      CargoTransferStatus // if transfer fails, this is the reason
	fleet       *Fleet
	dest        CargoHolder
	cargoType   CargoType
	transferred int
	wanted      int
}

func newCargoTransferer(log *slog.Logger, game *FullGame) cargoTransferer {
	return cargoTransferer{log: log, game: game, invader: newInvader()}
}

// load does a fleet's load tasks at a waypoint. Like the original game, cargo loads in a fixed order
// (ironium, boranium, germanium, colonists, fuel) and dunnage loads last, once no other load is waiting.
// It returns whether the fleet should wait at the waypoint for more cargo
func (t *cargoTransferer) load(fleet *Fleet, dest CargoHolder, transportTasks WaypointTransportTasks) (results []cargoTransferResult, wait bool) {
	// a set amount task couldn't get enough from the dest
	waitForCargo := false
	// a wait for percent task isn't filled yet. The fleet waits only while it has room left
	waitForPercent := false

	dunnage := []transportTask{}
	for _, task := range transportTasks.ordered() {
		if task.Action == TransportActionLoadDunnage {
			dunnage = append(dunnage, task)
			continue
		}

		result, taskWait := t.loadTask(fleet, dest, task)
		results = append(results, result)
		if task.Action == TransportActionWaitForPercent && task.cargoType != Fuel {
			waitForPercent = waitForPercent || taskWait
		} else {
			waitForCargo = waitForCargo || taskWait
		}
	}

	wait = waitForCargo || (waitForPercent && fleet.availableCargoSpace() > 0)
	if !wait {
		for _, task := range dunnage {
			result, _ := t.loadTask(fleet, dest, task)
			results = append(results, result)
		}
	}

	// delete this salvage or packet if we emptied it
	if salvage, ok := dest.(*Salvage); ok && salvage.Cargo == (Cargo{}) {
		t.game.deleteSalvage(salvage)
	}
	if packet, ok := dest.(*MineralPacket); ok && packet.Cargo == (Cargo{}) {
		t.game.deletePacket(packet)
	}

	return results, wait
}

// loadTask does a single load task
func (t *cargoTransferer) loadTask(fleet *Fleet, dest CargoHolder, task transportTask) (cargoTransferResult, bool) {
	transferAmount, wanted, wait := t.getCargoLoadAmount(fleet, dest, task.cargoType, task.WaypointTransportTask)
	transferred, status := t.transferCargo(fleet, -transferAmount, task.cargoType, dest)
	return cargoTransferResult{
		fleet:       fleet,
		dest:        dest,
		cargoType:   task.cargoType,
		wanted:      -wanted,
		transferred: transferred,
		status:      status,
	}, wait
}

// unload does a fleet's unload tasks at a waypoint, in the same order as loads
func (t *cargoTransferer) unload(fleet *Fleet, dest CargoHolder, transportTasks WaypointTransportTasks) (results []cargoTransferResult) {
	for _, task := range transportTasks.ordered() {
		transferAmount, wanted := t.getCargoUnloadAmount(fleet, dest, task.cargoType, task.WaypointTransportTask)
		transferred, status := t.transferCargo(fleet, transferAmount, task.cargoType, dest)
		results = append(results, cargoTransferResult{
			fleet:       fleet,
			dest:        dest,
			cargoType:   task.cargoType,
			wanted:      wanted,
			transferred: transferred,
			status:      status,
		})
	}

	return results
}

// transferCargo transfers a single cargo type to/from the fleet to/from the dest. A positive
// transferAmount unloads from the fleet, a negative one loads
func (t *cargoTransferer) transferCargo(fleet *Fleet, transferAmount int, cargoType CargoType, dest CargoHolder) (transferred int, invalid CargoTransferStatus) {
	if transferAmount == 0 {
		return 0, CargoTransferStatusNone
	}

	player := t.game.getPlayer(fleet.PlayerNum)
	if transferAmount > 0 {
		switch dest := dest.(type) {
		case *Salvage:
			if cargoType == Colonists {
				// colonists can't survive in space, they stay aboard
				return 0, CargoTransferStatusDeepSpace
			}
		case *Planet:
			if cargoType == Colonists && !dest.OwnedBy(fleet.PlayerNum) {
				invaders, status := t.invade(player, fleet, dest, transferAmount)
				fleet.Cargo.Colonists -= invaders
				return invaders, status
			}
		case *Fleet:
			if !dest.OwnedBy(fleet.PlayerNum) {
				// like the original game, colonists can't be given to another player's fleet, and a
				// player won't accept anything from an enemy
				if cargoType == Colonists || t.game.getPlayer(dest.PlayerNum).IsEnemy(fleet.PlayerNum) {
					t.log.Debug("fleet cannot unload to another player's fleet",
						slog.Int("Player", fleet.PlayerNum),
						slog.String("Fleet", fleet.Name),
						slog.String("Dest", dest.Name),
						slog.String("cargoType", cargoType.String()))
					return 0, CargoTransferStatusOwned
				}
			}
		}
	}

	if status := t.transferToDest(fleet, dest, cargoType, transferAmount); status != CargoTransferStatusNone {
		return 0, status
	}

	return transferAmount, CargoTransferStatusNone
}

// invade queues an invasion of another player's planet with colonists from a fleet. Like the original
// game, colonists can also take over a planet whose population died off this turn. It returns how many
// colonists (in kT) invade
func (t *cargoTransferer) invade(player *Player, fleet *Fleet, planet *Planet, colonists int) (int, CargoTransferStatus) {
	defender := t.game.getPlayer(planet.PlayerNum)
	if !planet.Owned() {
		defender = t.game.getPlayer(planet.ownerAtTurnStart)
	}

	status := CargoTransferStatusNone
	switch {
	case defender == nil:
		// the planet was empty when the turn started, don't beam colonists to their death
		status = CargoTransferStatusDestUnowned
	case planet.Spec.HasStarbase:
		status = CargoTransferStatusDestStarbase
	case player.Race.Spec.LivesOnStarbases:
		status = CargoTransferStatusOwned
	}
	if status != CargoTransferStatusNone {
		t.log.Debug("fleet cannot unload colonists on planet",
			slog.Int("Player", fleet.PlayerNum),
			slog.String("Fleet", fleet.Name),
			slog.String("Planet", planet.Name),
			slog.String("Status", status.String()))
		return 0, status
	}

	t.invader.addInvasion(invasion{
		planet:    planet,
		attacker:  player,
		defender:  defender,
		attackers: colonists * 100,
		fleets:    []*Fleet{fleet},
	})

	return colonists, CargoTransferStatusNone
}

// getCargoLoadAmount gets the amount of cargo to transfer for loading a cargo type from a cargoholder.
// waitAtWaypoint is true when a set amount task can't get enough from the dest, or a wait for percent
// task didn't fill the hold to its percent
func (t *cargoTransferer) getCargoLoadAmount(fleet *Fleet, dest CargoHolder, cargoType CargoType, task WaypointTransportTask) (transferAmount int, wantToTransfer int, waitAtWaypoint bool) {
	availableCapacity := fleet.Spec.CargoCapacity - fleet.Cargo.Total()
	availableToLoad := dest.GetCargo().GetAmount(cargoType)
	currentAmount := fleet.Cargo.GetAmount(cargoType)
	totalCapacity := fleet.Spec.CargoCapacity

	// fuel transfers use different tanks
	if cargoType == Fuel {
		availableCapacity = fleet.Spec.FuelCapacity - fleet.Fuel
		// planets with starbases have Infinite fuel, but fuel only moves between fleets
		availableToLoad = max(0, dest.GetFuel())
		currentAmount = fleet.Fuel
		totalCapacity = fleet.Spec.FuelCapacity
	}

	switch task.Action {
	case TransportActionLoadOptimal:
		// fuel only
		// we set our fuel to whatever it takes to finish our waypoints and transfer the rest to the ICargoHolder target.
		// If the target is a planet or starbase (and has infinite fuel capacity), we skip this and don't give them our fuel
		if cargoType == Fuel && dest.GetFuelCapacity() != Infinite {
			fuelRequiredForWaypoints := 0
			for i := 1; i < len(fleet.Waypoints); i++ {
				fuelRequiredForWaypoints += fleet.Waypoints[i].EstFuelUsage
			}
			leftoverFuel := fleet.Fuel - fuelRequiredForWaypoints
			wantToTransfer = leftoverFuel
			fuelCapacityAvailable := dest.GetFuelCapacity() - dest.GetFuel()
			if leftoverFuel > 0 && fuelCapacityAvailable > 0 {
				// transfer the lowest of how much fuel capacity they have available or how much we can give
				// this is a bit weird because we are doing a "Load", but it's actually an unload of fuel
				// from us to a dest fleet, so make the transferAmount negative.
				transferAmount = max(-leftoverFuel, -(dest.GetFuelCapacity() - dest.GetFuel()))
			}
		}
	case TransportActionLoadAll:
		// load all available, based on our constraints
		wantToTransfer = availableToLoad
		transferAmount = min(availableToLoad, availableCapacity)
	case TransportActionLoadAmount:
		wantToTransfer = task.Amount
		transferAmount = min(min(availableToLoad, task.Amount), availableCapacity)
	case TransportActionWaitForPercent, TransportActionFillPercent:
		// we want a percent of our hold to be filled with some amount, figure out how
		// much that is in kT, i.e. 50% of 100kT would be 50kT of this mineral
		var taskAmountkT = int(float64(task.Amount) / 100 * float64(totalCapacity))
		wantToTransfer = taskAmountkT

		if currentAmount >= taskAmountkT {
			// no need to transfer any, move on
			return 0, wantToTransfer, false
		} else {

			// transfer up to our percent specified
			// wait here if we haven't loaded the amount we want. The caller moves on anyway if the hold is
			// full (in case the user suffers from innumeracy and said they wanted 50% 50% 50%)
			transferAmount = min(min(availableToLoad, taskAmountkT-currentAmount), availableCapacity)
			if (transferAmount+currentAmount) < taskAmountkT && task.Action == TransportActionWaitForPercent {
				waitAtWaypoint = true
			}
		}
	case TransportActionSetAmountTo:
		// only transfer the min of what we have, vs what we need, vs the capacity
		wantToTransfer = max(0, task.Amount-currentAmount)
		transferAmount = max(0, min(min(availableToLoad, task.Amount-currentAmount), availableCapacity))
		// like the original game, wait for the dest to have enough, not for room in our hold
		if availableToLoad < task.Amount-currentAmount {
			waitAtWaypoint = true
		}
	case TransportActionSetWaypointTo:
		// Check how much the destination has of what we want
		// if we SetWaypointTo 100kT germanium and they have 120kT, we load 20kT if we can fit it
		if availableToLoad <= task.Amount {
			// they are below the amount, we won't load (this TransportAction will possibly be used to unload later)
			break
		} else {
			wantToTransfer = availableToLoad - task.Amount
			// only transfer down to what we set
			transferAmount = min(min(availableToLoad, availableToLoad-task.Amount), availableCapacity)
		}

	case TransportActionLoadDunnage:
		// (minerals and colonists only) This command waits until all other loads and unloads are
		// complete, then loads as many colonists or amount of a mineral as will fit in the remaining
		// space. For example, setting Load All Germanium, Load Dunnage Ironium, will load all the
		// Germanium that is available, then as much Ironium as possible. If more than one dunnage cargo
		// is specified, they are loaded in the order of Ironium, Boranium, Germanium, and Colonists.
		wantToTransfer = availableToLoad
		transferAmount = min(availableToLoad, availableCapacity)
	}

	// let the caller know how much of this cargo we load
	return transferAmount, wantToTransfer, waitAtWaypoint
}

// getCargoUnloadAmount gets the amount of cargo to transfer for unloading a cargo type from a cargoholder
func (t *cargoTransferer) getCargoUnloadAmount(fleet *Fleet, dest CargoHolder, cargoType CargoType, task WaypointTransportTask) (transferAmount int, wantToTransfer int) {

	capacity := dest.GetCargoCapacity()
	if capacity != Infinite {
		capacity = capacity - dest.GetCargo().Total()
	}
	currentAmount := fleet.Cargo.GetAmount(cargoType)

	var availableToUnload int
	if cargoType == Fuel {
		availableToUnload = fleet.Fuel
		capacity = max(0, dest.GetFuelCapacity()-dest.GetFuel())
		currentAmount = fleet.Fuel
	} else {
		availableToUnload = fleet.Cargo.GetAmount(cargoType)
	}
	switch task.Action {
	case TransportActionUnloadAll:
		// unload all available, based on our constraints
		wantToTransfer = availableToUnload
		if capacity == Infinite {
			transferAmount = availableToUnload
		} else {
			transferAmount = min(availableToUnload, capacity)
		}
	case TransportActionUnloadAmount:
		// don't unload more than the task says
		wantToTransfer = task.Amount
		if capacity == Infinite {
			transferAmount = min(availableToUnload, task.Amount)
		} else {
			transferAmount = min(min(availableToUnload, task.Amount), capacity)
		}
	case TransportActionSetAmountTo:
		// set the amount in our hold to amount, or do nothing if we have under that amount
		// unload what the dest has room for
		wantToTransfer = max(0, currentAmount-task.Amount)
		transferAmount = max(0, min(availableToUnload, currentAmount-task.Amount))
		if capacity != Infinite {
			transferAmount = min(transferAmount, capacity)
		}
	case TransportActionSetWaypointTo:
		// Make sure the waypoint has at least whatever we specified
		var currentAmount = dest.GetCargo().GetAmount(cargoType)

		if currentAmount >= task.Amount {
			// no need to transfer any, move on
			break
		} else {
			// only transfer the min of what we have, vs what we need, vs the capacity
			wantToTransfer = task.Amount - currentAmount
			if capacity == Infinite {
				transferAmount = min(availableToUnload, task.Amount-currentAmount)
			} else {
				transferAmount = min(min(availableToUnload, task.Amount-currentAmount), capacity)
			}
		}
	}
	return transferAmount, wantToTransfer
}

// transferToDest performs a transfer of a single cargo type to/from a destination. A positive
// transferAmount unloads from the fleet, a negative one loads.
// returns a status other than None if the transfer fails for some reason
func (t *cargoTransferer) transferToDest(fleet *Fleet, dest CargoHolder, cargoType CargoType, transferAmount int) CargoTransferStatus {
	status := t.checkTransferToDest(fleet, dest, cargoType, transferAmount)
	if status != CargoTransferStatusNone {
		t.log.Debug("fleet cannot transfer cargo",
			slog.Int("Player", fleet.PlayerNum),
			slog.String("Fleet", fleet.Name),
			slog.String("Dest", dest.GetMapObject().Name),
			slog.String("cargoType", cargoType.String()),
			slog.Int("TransferAmount", transferAmount),
			slog.String("Status", status.String()))
		return status
	}

	if cargoType == Fuel {
		fleet.Fuel -= transferAmount
		dest.(*Fleet).Fuel += transferAmount
		return CargoTransferStatusNone
	}

	fleet.Cargo = fleet.Cargo.SubtractAmount(cargoType, transferAmount)
	dest.SetCargo(dest.GetCargo().AddAmount(cargoType, transferAmount))
	return CargoTransferStatusNone
}

// checkTransferToDest checks whether a fleet can transfer a single cargo type to/from a destination
func (t *cargoTransferer) checkTransferToDest(fleet *Fleet, dest CargoHolder, cargoType CargoType, transferAmount int) CargoTransferStatus {
	if transferAmount < 0 {
		if !dest.CanLoad(fleet) {
			return CargoTransferStatusOwned
		}
		// like the original game (with its bug fixed), a thief can't take colonists or fuel
		if (cargoType == Colonists || cargoType == Fuel) && !dest.GetMapObject().OwnedBy(fleet.PlayerNum) {
			return CargoTransferStatusOwned
		}
	}

	if cargoType == Fuel {
		_, destIsFleet := dest.(*Fleet)
		switch {
		case transferAmount > 0 && fleet.Fuel < transferAmount:
			return CargoTransferStatusCargo
		case !destIsFleet:
			// like the original game, fuel only moves between fleets
			if transferAmount > 0 {
				return CargoTransferStatusDestCargoCapacity
			}
			return CargoTransferStatusDestCargo
		case transferAmount > 0 && dest.GetFuelCapacity()-dest.GetFuel() < transferAmount:
			return CargoTransferStatusDestCargoCapacity
		case transferAmount < 0 && fleet.availableFuelSpace() < -transferAmount:
			return CargoTransferStatusCargoCapacity
		case transferAmount < 0 && dest.GetFuel() < -transferAmount:
			return CargoTransferStatusDestCargo
		}
		return CargoTransferStatusNone
	}

	destCargo := dest.GetCargo()
	switch {
	case transferAmount > 0 && !fleet.Cargo.CanTransferAmount(cargoType, transferAmount):
		return CargoTransferStatusCargo
	case transferAmount > 0 && dest.GetCargoCapacity() != Infinite && (dest.GetCargoCapacity()-destCargo.Total()) < transferAmount:
		return CargoTransferStatusDestCargoCapacity
	case transferAmount < 0 && fleet.availableCargoSpace() < -transferAmount:
		return CargoTransferStatusCargoCapacity
	case transferAmount < 0 && !destCargo.CanTransferAmount(cargoType, -transferAmount):
		return CargoTransferStatusDestCargo
	}
	return CargoTransferStatusNone
}
