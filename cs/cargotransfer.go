package cs

import (
	"fmt"
	"slices"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CargoTransfers map[string][]ByHandCargoTransfer

type ByHandCargoTransfer struct {
	MapObjectTarget `tstype:",extends"`
	SourceFleetNum  int   `json:"sourceFleetNum,omitempty"`
	Cargo           Cargo `json:"cargo"`
}

type cargoTransferer struct {
	log     zerolog.Logger
	invader invader
	game    *FullGame
}

type CargoTransferStatus int

const (
	CargoTransferStatusNone CargoTransferStatus = iota
	CargoTransferStatusOwned
	CargoTransferStatusCargo
	CargoTransferStatusCargoCapacity
	CargoTransferStatusDestCargo
	CargoTransferStatusDestCargoCapacity
	CargoTransferStatusDestStarbase
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
	default:
		return fmt.Sprintf("Unknown %d", r)
	}
}

type cargoTransferResult struct {
	status         CargoTransferStatus // if transfer fails, this is the reason
	fleet          *Fleet
	dest           CargoHolder
	cargoType      CargoType
	transferred    int
	wanted         int
	waitAtWaypoint bool
}

func newCargoTransferer(log zerolog.Logger, game *FullGame) cargoTransferer {
	return cargoTransferer{log: log, game: game, invader: newInvader()}
}

func (cargoTransfers CargoTransfers) getTransfers(position Vector) []ByHandCargoTransfer {
	return cargoTransfers[position.String()]
}

// getByHandTransfer sums all cargo load/unloads for this position to determine the total amount of by hand cargo here
func (cargoTransfers CargoTransfers) getByHandTransfer(target MapObjectTarget) Cargo {
	cargo := Cargo{}
	if cargoTransfers == nil {
		return cargo
	}
	key := target.TargetPosition.String()

	for _, transfer := range cargoTransfers[key] {
		if transfer.MapObjectTarget != target {
			continue
		}
		cargo = cargo.Add(transfer.Cargo)
	}
	return cargo
}

// transferByHand adds a byHand transfer to a target
func (cargoTransfers CargoTransfers) transferByHand(fleet *Fleet, target MapObjectTarget, cargo Cargo) {
	// add the new cargo transfer to the player
	key := fleet.Position.String()
	transfers := cargoTransfers[key]

	// if the last transfer is the same fleet/target, just update it instead of creating another one
	var lastTransfer *ByHandCargoTransfer
	if len(transfers) > 0 {
		lastTransfer = &transfers[len(transfers)-1]
	}
	if lastTransfer != nil && lastTransfer.SourceFleetNum == fleet.Num && lastTransfer.MapObjectTarget == target {
		lastTransfer.Cargo = lastTransfer.Cargo.Add(cargo)
		return
	}

	transfer := ByHandCargoTransfer{
		SourceFleetNum:  fleet.Num,
		Cargo:           cargo,
		MapObjectTarget: target,
	}

	// add a new immediate cargo transfer for this jettison
	cargoTransfers[key] = append(cargoTransfers[key], transfer)
}

// splitByHandTransfers splits the ByHandCargoTransfers for a source fleet into two
// ByHandCargoTransfers, based on capacity of each fleet
func (cargoTransfers CargoTransfers) splitByHandTransfers(source *Fleet, dest *Fleet) error {
	key := source.Position.String()
	transfers, ok := cargoTransfers[key]
	if !ok {
		// no transfers
		return nil
	}

	updatedTransfers := make([]ByHandCargoTransfer, 0, len(transfers))
	for _, transfer := range transfers {
		if transfer.SourceFleetNum != source.Num {
			updatedTransfers = append(updatedTransfers, transfer)
			continue
		}
		// split this transfer
		sourceCargo := transfer.Cargo.ToArray()
		cargo1, cargo2, err := splitValues(
			source.Spec.CargoCapacity+dest.Spec.CargoCapacity,
			source.Spec.CargoCapacity,
			dest.Spec.CargoCapacity,
			sourceCargo[:]...)
		if err != nil {
			return err
		}
		// copy this transfer into two transfers
		transfer1 := transfer
		transfer2 := transfer

		transfer1.Cargo = NewCargoFromArray([4]int(cargo1))
		if transfer1.Cargo.absSum() > 0 {
			updatedTransfers = append(updatedTransfers, transfer1)
		}

		transfer2.Cargo = NewCargoFromArray([4]int(cargo2))
		if transfer2.Cargo.absSum() > 0 {
			transfer2.SourceFleetNum = dest.Num
			updatedTransfers = append(updatedTransfers, transfer2)
		}
	}

	cargoTransfers[key] = updatedTransfers
	return nil
}

// mergeByHandTransfers merges cargo transfers from merging fleets into a source fleet
func (cargoTransfers CargoTransfers) mergeByHandTransfers(fleet *Fleet, mergingFleets []*Fleet) {
	key := fleet.Position.String()
	transfers, ok := cargoTransfers[key]
	if !ok {
		// no transfers
		return
	}

	updatedTransfers := make([]ByHandCargoTransfer, 0, len(transfers))
	for i, transfer := range transfers {
		if !slices.ContainsFunc(mergingFleets, func(f *Fleet) bool { return transfer.SourceFleetNum == f.Num }) {
			// this transfer isn't about a merging fleet, continue
			updatedTransfers = append(updatedTransfers, transfer)
			continue
		}

		var prevTransfer *ByHandCargoTransfer
		if i > 0 {
			prevTransfer = &updatedTransfers[i-1]
		}

		// if the previous transfer is the fleet we're merging into AND the target is the same, just merge the request
		// into the previous transfer
		if prevTransfer != nil && prevTransfer.SourceFleetNum == fleet.Num && prevTransfer.MapObjectTarget == transfer.MapObjectTarget {
			prevTransfer.Cargo = prevTransfer.Cargo.Add(transfer.Cargo)
			continue
		}

		// give this transfer to the merged fleet
		transfer.SourceFleetNum = fleet.Num
		updatedTransfers = append(updatedTransfers, transfer)
	}

	cargoTransfers[key] = updatedTransfers
}

// getUnloadTasks creates a WaypointTransportTask for each positive cargo value in the ByHandCargoTransfer
func (o ByHandCargoTransfer) getUnloadTasks() WaypointTransportTasks {
	tt := WaypointTransportTasks{}
	if o.Cargo.Ironium > 0 {
		tt.Ironium.Action = TransportActionUnloadAmount
		tt.Ironium.Amount = o.Cargo.Ironium
	}
	if o.Cargo.Boranium > 0 {
		tt.Boranium.Action = TransportActionUnloadAmount
		tt.Boranium.Amount = o.Cargo.Boranium
	}
	if o.Cargo.Germanium > 0 {
		tt.Germanium.Action = TransportActionUnloadAmount
		tt.Germanium.Amount = o.Cargo.Germanium
	}
	if o.Cargo.Colonists > 0 {
		tt.Colonists.Action = TransportActionUnloadAmount
		tt.Colonists.Amount = o.Cargo.Colonists
	}

	return tt
}

// getLoadTasks creates a WaypointTransportTask for each negative cargo value in the ByHandCargoTransfer
func (o ByHandCargoTransfer) getLoadTasks() WaypointTransportTasks {
	tt := WaypointTransportTasks{}
	if o.Cargo.Ironium < 0 {
		tt.Ironium.Action = TransportActionLoadAmount
		tt.Ironium.Amount = -o.Cargo.Ironium
	}
	if o.Cargo.Boranium < 0 {
		tt.Boranium.Action = TransportActionLoadAmount
		tt.Boranium.Amount = -o.Cargo.Boranium
	}
	if o.Cargo.Germanium < 0 {
		tt.Germanium.Action = TransportActionLoadAmount
		tt.Germanium.Amount = -o.Cargo.Germanium
	}
	if o.Cargo.Colonists < 0 {
		tt.Colonists.Action = TransportActionLoadAmount
		tt.Colonists.Amount = -o.Cargo.Colonists
	}

	return tt
}

func (t *cargoTransferer) loadByHands(player *Player, transfers []ByHandCargoTransfer) []cargoTransferResult {
	var results []cargoTransferResult

	// we iterate over transfers by their target so we can keep track of the running cargo amount from by hand unloads
	transfersByTarget := make(map[MapObjectTarget][]ByHandCargoTransfer)
	for _, transfer := range transfers {
		transfersByTarget[transfer.MapObjectTarget] = append(transfersByTarget[transfer.MapObjectTarget], transfer)
	}

	for _, transfers := range transfersByTarget {
		// as we do by hand loads, load from any by hand unloads first
		// do this by keeping track of a bucket of cargo for this transfer
		// the idea is to handle scenarios like this
		//
		// Fleet 1 unloads 50kT ironium
		// Fleet 2 loads 60kT ironium
		// in the above scenario, Fleet 2 loads 50kT from the bucket, and 10kT from the dest
		cargoInBucket := Cargo{}
		for _, transfer := range transfers {
			// add any by hand unloads to the bucket
			cargoInBucket = cargoInBucket.Add(transfer.Cargo.PositiveOnly())

			// find any loads. A fleet loading 5kT of ironium from a planet would have transfer of {ironium: -5}
			cargoToLoad := transfer.Cargo.NegativeOnly()

			// load from the bucket first
			for _, cargoType := range CargoTypes {
				amount := cargoToLoad.GetAmount(cargoType)
				amountInBucket := cargoInBucket.GetAmount(cargoType)
				// if we are loading cargo, take it from the bucket first
				if amount < 0 && amountInBucket > 0 {
					loadAmount := Min(0, amount+amountInBucket)
					cargoInBucket = cargoInBucket.SetAmount(cargoType, amountInBucket+loadAmount)
					cargoToLoad = cargoToLoad.SetAmount(cargoType, loadAmount)

					log.Debug().
						Int("Player", player.Num).
						Int("FleetNum", transfer.SourceFleetNum).
						Str("Target", transfer.MapObjectTarget.String()).
						Str("cargoInBucket", cargoInBucket.PrettyString()).
						Str("cargoToLoad", cargoToLoad.PrettyString()).
						Msgf("by hand load cargo")

				}
			}
			// update the transfer with the actual amount we're loading
			transfer.Cargo = cargoToLoad

			if cargoToLoad == (Cargo{}) {
				// skip any empty requests
				continue
			}

			fleet := t.game.Universe.getFleet(player.Num, transfer.SourceFleetNum)
			if fleet == nil {
				// don't kill turn processing for this because it's unclear how to fix it to unblock players
				t.log.Error().
					Int("Player", player.Num).
					Int("Fleet", transfer.SourceFleetNum).
					Msgf("fleet not found for ByHandCargoTransfer")
				continue
			}

			// convert all by hand load transfers into "Load Amount" style WaypointTransportTasks
			transportTasks := transfer.getLoadTasks()

			// for by hand transfers, the fleet already thinks it loaded this cargo, so take away the cargo and make the
			// fleet load it for real
			fleet.Cargo = fleet.Cargo.Add(cargoToLoad.NegativeOnly())

			dest, ok := t.game.getCargoHolder(transfer.TargetType, transfer.TargetNum, transfer.TargetPlayerNum)
			if !ok || dest.Deleted() {
				// can't load from space
				continue
			}

			mo := dest.GetMapObject()
			if mo.OwnedBy(fleet.PlayerNum) {
				// this transfer already happened so reverse it and add it again for real this time
				dest.SetCargo(dest.GetCargo().Subtract(cargoToLoad.NegativeOnly()))
			}
			log.Debug().
				Int("Player", fleet.PlayerNum).
				Str("Fleet", fleet.Name).
				Str("Dest", dest.GetMapObject().Name).
				Str("Cargo", fleet.Cargo.PrettyString()).
				Str("DestCargo", dest.GetCargo().PrettyString()).
				Str("CargoToLoad", cargoToLoad.PrettyString()).
				Msgf("by hand load cargo")

			results = append(results, t.load(fleet, dest, transportTasks)...)
		}
	}
	return results
}

// unloadByHands processes all by hand unloads for a player for a location
// these transfers should be recorded in the order they are performed
func (t *cargoTransferer) unloadByHands(player *Player, transfers []ByHandCargoTransfer) []cargoTransferResult {
	var results []cargoTransferResult

	// we iterate over transfers by their target so we can keep track of the running cargo amount from by hand loads
	transfersByTarget := make(map[MapObjectTarget][]ByHandCargoTransfer)
	for _, transfer := range transfers {
		transfersByTarget[transfer.MapObjectTarget] = append(transfersByTarget[transfer.MapObjectTarget], transfer)
	}

	for _, transfers := range transfersByTarget {
		// as we do by hand unloads, don't unload any cargo we by hand loaded
		// do this by keeping track of a bucket of cargo for this transfer
		// the idea is to handle scenarios like this
		//
		// Fleet 1 unloads 50kT ironium
		// Fleet 2 loads 10kT ironium
		// in the above scenario, Fleet 2 loads 10kT from the bucket, and Fleet 1 only unloads 40kT

		// loads happen first, so for unloads, handle it in reverse
		cargoInBucket := Cargo{}
		for i := len(transfers) - 1; i >= 0; i-- {
			transfer := transfers[i]

			// add any by hand loads to the bucket (these would be negative)
			cargoInBucket = cargoInBucket.Add(transfer.Cargo.NegativeOnly())

			cargoToUnload := transfer.Cargo.PositiveOnly()

			// account for any by hand loads from the bucket
			for _, cargoType := range CargoTypes {
				amount := cargoToUnload.GetAmount(cargoType)
				amountInBucket := cargoInBucket.GetAmount(cargoType)
				// if we are loading cargo, take it from the bucket first
				if amount > 0 && amountInBucket < 0 {
					unloadAmount := Max(0, amount+amountInBucket)
					cargoInBucket = cargoInBucket.SetAmount(cargoType, amountInBucket+unloadAmount)
					cargoToUnload = cargoToUnload.SetAmount(cargoType, unloadAmount)

					log.Debug().
						Int("Player", player.Num).
						Int("FleetNum", transfer.SourceFleetNum).
						Str("Target", transfer.MapObjectTarget.String()).
						Str("cargoInBucket", cargoInBucket.PrettyString()).
						Str("cargoToUnload", cargoToUnload.PrettyString()).
						Msgf("by hand unload cargo")

				}
			}
			// update the transfer with the actual amount we're loading
			transfer.Cargo = cargoToUnload

			if cargoToUnload == (Cargo{}) {
				// skip any empty requests
				continue
			}

			fleet := t.game.Universe.getFleet(player.Num, transfer.SourceFleetNum)
			if fleet == nil {
				// don't kill turn processing for this because it's unclear how to fix it to unblock players
				t.log.Error().
					Int("Player", player.Num).
					Int("Fleet", transfer.SourceFleetNum).
					Msgf("fleet not found for ByHandCargoTransfer")
				continue
			}

			// convert all by hand unload transfers into "Unload Amount" style WaypointTransportTasks
			transportTasks := transfer.getUnloadTasks()

			// for by hand transfers, the fleet already thinks it unloaded this cargo, so add back the cargo and make the
			// fleet unload it for real
			fleet.Cargo = fleet.Cargo.Add(cargoToUnload.PositiveOnly())

			dest, ok := t.game.getCargoHolder(transfer.TargetType, transfer.TargetNum, transfer.TargetPlayerNum)
			newlyCreatedSalvage := false
			if !ok && transfer.TargetType == MapObjectTypeNone {
				// create a salvage
				dest = t.game.getOrCreateSalvage(fleet.Position, fleet.PlayerNum, Cargo{})
				newlyCreatedSalvage = true
			}

			mo := dest.GetMapObject()
			if mo.OwnedBy(fleet.PlayerNum) && !newlyCreatedSalvage {
				// this transfer already happened so reverse it and transfer it again for real this time
				dest.SetCargo(dest.GetCargo().Subtract(cargoToUnload.PositiveOnly()))
			}

			log.Debug().
				Int("Player", fleet.PlayerNum).
				Str("Fleet", fleet.Name).
				Str("Dest", dest.GetMapObject().Name).
				Str("Cargo", fleet.Cargo.PrettyString()).
				Str("DestCargo", dest.GetCargo().PrettyString()).
				Str("CargoToUnload", cargoToUnload.PrettyString()).
				Msgf("by hand unload cargo")

			results = append(results, t.unload(fleet, dest, transportTasks)...)
		}
	}
	return results
}

// load executes cargo load operations for a single fleet, returning a result for each cargoType transfered
func (t *cargoTransferer) load(fleet *Fleet, dest CargoHolder, transportTasks WaypointTransportTasks) []cargoTransferResult {
	results := []cargoTransferResult{}

	// dunnage tasks are done after regular tasks
	type dunnageTask struct {
		cargoType CargoType
		task      WaypointTransportTask
	}
	dunnageTasks := []dunnageTask{}

	// process regular load tasks
	for cargoType, task := range transportTasks.getTransportTasks() {
		if task.Action == TransportActionLoadDunnage {
			dunnageTasks = append(dunnageTasks, dunnageTask{cargoType, task})
			continue
		}

		// get the load amount
		transferAmount, wanted, wait := t.getCargoLoadAmount(fleet, dest, cargoType, task)

		// do the transfer and record the result
		transferred, status := t.transferCargo(fleet, -transferAmount, cargoType, dest)
		results = append(results, cargoTransferResult{
			fleet:          fleet,
			dest:           dest,
			cargoType:      cargoType,
			wanted:         -wanted,
			waitAtWaypoint: wait,
			transferred:    transferred,
			status:         status,
		})
	}

	// process dunnage tasks after all other loads
	for _, dunnageTask := range dunnageTasks {
		cargoType, task := dunnageTask.cargoType, dunnageTask.task

		// get any dunnage load amount
		transferAmount, wanted, wait := t.getCargoLoadAmount(fleet, dest, cargoType, task)

		// do the transfer and record the result
		transferred, status := t.transferCargo(fleet, -transferAmount, cargoType, dest)
		results = append(results, cargoTransferResult{
			fleet:          fleet,
			dest:           dest,
			cargoType:      cargoType,
			wanted:         wanted,
			waitAtWaypoint: wait,
			transferred:    transferred,
			status:         status,
		})
	}

	// delete this salvage if we emptied it
	if salvage, ok := dest.(*Salvage); ok && salvage.Cargo == (Cargo{}) {
		t.game.deleteSalvage(salvage)

		t.log.Debug().
			Int("Player", salvage.PlayerNum).
			Str("Salvage", salvage.Name).
			Msgf("deleted salvage")

	}
	// delete this packet if we emptied it
	if packet, ok := dest.(*MineralPacket); ok && packet.Cargo == (Cargo{}) {
		t.game.deletePacket(packet)

		t.log.Debug().
			Int("Player", packet.PlayerNum).
			Str("Packet", packet.Name).
			Msgf("deleted salvage")

	}

	return results
}

func (t *cargoTransferer) unload(fleet *Fleet, dest CargoHolder, transportTasks WaypointTransportTasks) []cargoTransferResult {
	results := []cargoTransferResult{}

	for cargoType, task := range transportTasks.getTransportTasks() {
		// get how much this order wants to unload
		transferAmount, wanted, wait := t.getCargoUnloadAmount(fleet, dest, cargoType, task)

		// perform the transfer and record the result
		transferred, status := t.transferCargo(fleet, transferAmount, cargoType, dest)
		results = append(results, cargoTransferResult{
			fleet:          fleet,
			dest:           dest,
			cargoType:      cargoType,
			wanted:         wanted,
			waitAtWaypoint: wait,
			transferred:    transferred,
			status:         status,
		})
	}

	return results
}

// transferCargo transfers a single cargo type to/from the fleet to/from the dest
func (t *cargoTransferer) transferCargo(fleet *Fleet, transferAmount int, cargoType CargoType, dest CargoHolder) (transferred int, invalid CargoTransferStatus) {
	if transferAmount == 0 {
		return 0, CargoTransferStatusNone
	}

	// check for invasion
	player := t.game.Players[fleet.PlayerNum-1]
	planet, ok := dest.(*Planet)
	if transferAmount > 0 && cargoType == Colonists && ok && planet.Owned() && !planet.OwnedBy(fleet.PlayerNum) {
		if planet.Spec.HasStarbase {
			// can't invade a planet with a starbase
			t.log.Debug().
				Int("Player", fleet.PlayerNum).
				Str("Fleet", fleet.Name).
				Str("Dest", dest.GetMapObject().Name).
				Str("cargoType", cargoType.String()).
				Msgf("fleet %s cannot unload %d00 colonists to %s, starbase is in orbit", fleet.Name, transferAmount, dest.GetMapObject().Name)

			return 0, CargoTransferStatusDestStarbase
		}

		// invasion!
		attacker := player
		defender := t.game.getPlayer(planet.PlayerNum)

		t.invader.addInvasion(invasion{
			planet:    planet,
			attacker:  attacker,
			defender:  defender,
			attackers: transferAmount * 100,
			fleets:    []*Fleet{fleet},
		})

		fleet.Cargo.Colonists -= transferAmount
		return transferAmount, CargoTransferStatusNone
	}

	if status := t.transferToDest(fleet, dest, cargoType, transferAmount); status != CargoTransferStatusNone {
		return 0, status
	}
	return transferAmount, CargoTransferStatusNone
}

// getTransferAmount gets the amount of cargo to transfer for loading a cargo type from a cargoholder
func (t *cargoTransferer) getCargoLoadAmount(fleet *Fleet, dest CargoHolder, cargoType CargoType, task WaypointTransportTask) (transferAmount int, wantToTransfer int, waitAtWaypoint bool) {
	availableCapacity := fleet.Spec.CargoCapacity - fleet.Cargo.Total()
	availableToLoad := dest.GetCargo().GetAmount(cargoType)
	currentAmount := fleet.Cargo.GetAmount(cargoType)
	totalCapacity := fleet.Spec.CargoCapacity

	// fuel transfers use different tanks
	if cargoType == Fuel {
		availableCapacity = fleet.Spec.FuelCapacity - fleet.Fuel
		availableToLoad = dest.GetFuel()
		currentAmount = fleet.Fuel
		totalCapacity = fleet.Spec.FuelCapacity
	}

	switch task.Action {
	case TransportActionLoadOptimal:
		// fuel only
		// we set our fuel to whatever it takes to finish our waypoints and transfer the rest to the ICargoHolder target.
		// If the target is a planet or starbase (and has infinite fuel capacity), we skip this and don't give them our fuel
		if cargoType == Fuel && dest.GetFuelCapacity() != Unlimited {
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
				transferAmount = Max(-leftoverFuel, -(dest.GetFuelCapacity() - dest.GetFuel()))
			}
		}
	case TransportActionLoadAll:
		// load all available, based on our constraints
		wantToTransfer = availableToLoad
		transferAmount = Min(availableToLoad, availableCapacity)
	case TransportActionLoadAmount:
		wantToTransfer = task.Amount
		transferAmount = Min(Min(availableToLoad, task.Amount), availableCapacity)
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
			// wait here if we haven't loaded the amount we want
			// but move on if we are out of cargo space (in case the user suffers from innumeracy and said they wanted 50% 50% 50%)
			transferAmount = Min(Min(availableToLoad, taskAmountkT-currentAmount), availableCapacity)
			if (transferAmount+currentAmount) < taskAmountkT && task.Action == TransportActionWaitForPercent && (availableCapacity-transferAmount) > 0 {
				waitAtWaypoint = true
			}
		}
	case TransportActionSetAmountTo:
		// only transfer the min of what we have, vs what we need, vs the capacity
		wantToTransfer = Max(0, task.Amount-currentAmount)
		transferAmount = Max(0, Min(Min(availableToLoad, task.Amount-currentAmount), availableCapacity))
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
			wantToTransfer = availableToLoad - task.Amount
			// only transfer down to what we set
			transferAmount = Min(Min(availableToLoad, availableToLoad-task.Amount), availableCapacity)
		}

	case TransportActionLoadDunnage:
		// (minerals and colonists only) This command waits until all other loads and unloads are
		// complete, then loads as many colonists or amount of a mineral as will fit in the remaining
		// space. For example, setting Load All Germanium, Load Dunnage Ironium, will load all the
		// Germanium that is available, then as much Ironium as possible. If more than one dunnage cargo
		// is specified, they are loaded in the order of Ironium, Boranium, Germanium, and Colonists.
		wantToTransfer = availableToLoad
		transferAmount = Min(availableToLoad, availableCapacity)
	}

	// let the caller know how much of this cargo we load
	return transferAmount, wantToTransfer, waitAtWaypoint
}

// getCargoUnloadAmount gets the amount of cargo to transfer for unloading a cargo type from a cargoholder
func (t *cargoTransferer) getCargoUnloadAmount(fleet *Fleet, dest CargoHolder, cargoType CargoType, task WaypointTransportTask) (transferAmount int, wantToTransfer int, waitAtWaypoint bool) {

	capacity := dest.GetCargoCapacity()
	currentAmount := fleet.Cargo.GetAmount(cargoType)

	var availableToUnload int
	if cargoType == Fuel {
		availableToUnload = fleet.Fuel
		capacity = Max(0, dest.GetFuelCapacity()-dest.GetFuel())
		currentAmount = fleet.Fuel
	} else {
		availableToUnload = fleet.Cargo.GetAmount(cargoType)
	}
	switch task.Action {
	case TransportActionUnloadAll:
		// unload all available, based on our constraints
		wantToTransfer = availableToUnload
		if capacity == Unlimited {
			transferAmount = availableToUnload
		} else {
			transferAmount = Min(availableToUnload, capacity)
		}
	case TransportActionUnloadAmount:
		// don't unload more than the task says
		wantToTransfer = task.Amount
		if capacity == Unlimited {
			transferAmount = Min(availableToUnload, task.Amount)
		} else {
			transferAmount = Min(Min(availableToUnload, task.Amount), capacity)
		}
	case TransportActionSetAmountTo:
		// set the amount in our hold to amount, or do nothing if we have under that amount
		wantToTransfer = Max(0, currentAmount-task.Amount)
		transferAmount = Max(0, Min(availableToUnload, currentAmount-task.Amount))
	case TransportActionSetWaypointTo:
		// Make sure the waypoint has at least whatever we specified
		var currentAmount = dest.GetCargo().GetAmount(cargoType)

		if currentAmount >= task.Amount {
			// no need to transfer any, move on
			break
		} else {
			// only transfer the min of what we have, vs what we need, vs the capacity
			wantToTransfer = task.Amount - currentAmount
			if capacity == Unlimited {
				transferAmount = Min(availableToUnload, task.Amount-currentAmount)
			} else {
				transferAmount = Min(Min(availableToUnload, task.Amount-currentAmount), capacity)
			}
		}
	}
	return transferAmount, wantToTransfer, waitAtWaypoint
}

// transferToDest performs a transfer of a single cargo type to/from a destination
// returns a status other than None if the transfer fails for some reason
func (t *cargoTransferer) transferToDest(fleet *Fleet, dest CargoHolder, cargoType CargoType, transferAmount int) CargoTransferStatus {
	destCargo := dest.GetCargo()

	// check for load from owned dest
	if transferAmount < 0 && !dest.CanLoad(fleet) {
		// can't load from things we don't own
		t.log.Debug().
			Int("Player", fleet.PlayerNum).
			Str("Fleet", fleet.Name).
			Str("Dest", dest.GetMapObject().Name).
			Str("cargoType", cargoType.String()).
			Msgf("fleet %s cannot load %d to %s, does not own dest", fleet.Name, transferAmount, dest.GetMapObject().Name)

		return CargoTransferStatusOwned
	}

	if transferAmount > 0 && !fleet.Cargo.CanTransferAmount(cargoType, transferAmount) {
		t.log.Debug().
			Int("Player", fleet.PlayerNum).
			Str("Fleet", fleet.Name).
			Str("Dest", dest.GetMapObject().Name).
			Str("cargoType", cargoType.String()).
			Msgf("fleet %s cannot transfer %d to %s, there is not enough in the fleet to transfer", fleet.Name, transferAmount, dest.GetMapObject().Name)
		return CargoTransferStatusCargo
	}

	if transferAmount < 0 && fleet.availableCargoSpace() < -transferAmount {
		t.log.Debug().
			Int("Player", fleet.PlayerNum).
			Str("Fleet", fleet.Name).
			Str("Dest", dest.GetMapObject().Name).
			Str("cargoType", cargoType.String()).
			Msgf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", fleet.Name, fleet.availableCargoSpace(), transferAmount, dest.GetMapObject().Name)
		return CargoTransferStatusCargoCapacity
	}

	if transferAmount < 0 && !destCargo.CanTransferAmount(cargoType, -transferAmount) {
		t.log.Debug().
			Int("Player", fleet.PlayerNum).
			Str("Fleet", fleet.Name).
			Str("Dest", dest.GetMapObject().Name).
			Str("cargoType", cargoType.String()).
			Msgf("fleet %s cannot transfer %d from %s, there is not enough to transfer", fleet.Name, transferAmount, dest.GetMapObject().Name)
		return CargoTransferStatusDestCargo
	}

	if transferAmount > 0 && dest.GetCargoCapacity() != Unlimited && (dest.GetCargoCapacity()-destCargo.Total()) < transferAmount {
		t.log.Debug().
			Int("Player", fleet.PlayerNum).
			Str("Fleet", fleet.Name).
			Str("Dest", dest.GetMapObject().Name).
			Str("cargoType", cargoType.String()).
			Msgf("fleet %s cannot transfer %d to %s, there is not enough to space to hold the cargo", fleet.Name, transferAmount, dest.GetMapObject().Name)
		return CargoTransferStatusDestCargoCapacity

	}

	// transfer the cargo
	fleet.Cargo = fleet.Cargo.SubtractAmount(cargoType, transferAmount)
	dest.SetCargo(destCargo.AddAmount(cargoType, transferAmount))

	return CargoTransferStatusNone
}
