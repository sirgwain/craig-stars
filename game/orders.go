package game

import "fmt"

// handle player orders
type orderer interface {
	updatePlayerOrders(player *Player, playerPlanets []*Planet, order PlayerOrders)
	updatePlanetOrders(player *Player, planet *Planet, orders PlanetOrders)
	updateFleetOrders(player *Player, fleet *Fleet, orders FleetOrders)
	updateMineFieldOrders(player *Player, minefield *MineField, orders MineFieldOrders)
	transferFleetCargo(source *Fleet, dest *Fleet, transferAmount Cargo) error
	transferPlanetCargo(source *Fleet, dest *Planet, transferAmount Cargo) error
}

type orders struct {
}

// update a player's orders
func (o *orders) updatePlayerOrders(player *Player, playerPlanets []*Planet, orders PlayerOrders) {
	researchAmountUpdated := orders.ResearchAmount != player.ResearchAmount

	// save the new orders
	player.PlayerOrders = orders

	// if we updated the research amount, update the planet specs
	if researchAmountUpdated {
		for _, planet := range playerPlanets {
			if planet.PlayerNum == player.Num {
				spec := &planet.Spec
				spec.computeResourcesPerYearAvailable(player, planet)
				planet.Dirty = true
			}
		}
	}

}

// update a planet orders
func (o *orders) updatePlanetOrders(player *Player, planet *Planet, orders PlanetOrders) {
	planet.PlanetOrders = orders

	spec := &planet.Spec
	spec.computeResourcesPerYearAvailable(player, planet)
}

// update the orders to a fleet
func (o *orders) updateFleetOrders(player *Player, fleet *Fleet, orders FleetOrders) {
	// copy user modifiable things to the fleet fleet
	fleet.RepeatOrders = orders.RepeatOrders
	wp0 := &fleet.Waypoints[0]
	newWP0 := orders.Waypoints[0]

	// TODO: do we want to lookup the target?
	wp0.WarpFactor = newWP0.WarpFactor
	wp0.Task = newWP0.Task
	wp0.TransportTasks = newWP0.TransportTasks
	wp0.WaitAtWaypoint = newWP0.WaitAtWaypoint
	wp0.TargetName = newWP0.TargetName
	wp0.TargetType = newWP0.TargetType
	wp0.TargetNum = newWP0.TargetNum
	wp0.TargetPlayerNum = newWP0.TargetPlayerNum
	wp0.TransferToPlayer = newWP0.TransferToPlayer

	fleet.Waypoints = append(fleet.Waypoints[:1], orders.Waypoints[1:]...)
	fleet.computeFuelUsage(player)
}

func (o *orders) updateMineFieldOrders(player *Player, minefield *MineField, orders MineFieldOrders) {
	minefield.MineFieldOrders = orders
}

// transfer cargo from a fleet to/from a fleet
func (o *orders) transferFleetCargo(source *Fleet, dest *Fleet, transferAmount Cargo) error {

	if source.availableCargoSpace() < transferAmount.Total() {
		return fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", source.Name, source.availableCargoSpace(), transferAmount.Total(), dest.Name)
	}

	if !dest.Cargo.CanTransfer(transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %v from %s, there is not enough to transfer", source.Name, transferAmount, dest.Name)
	}

	// transfer the cargo
	source.Cargo = source.Cargo.Add(transferAmount)
	dest.Cargo = dest.Cargo.Subtract(transferAmount)

	return nil
}

// transfer cargo from a planet to/from a fleet
func (o *orders) transferPlanetCargo(source *Fleet, dest *Planet, transferAmount Cargo) error {
	if source.availableCargoSpace() < transferAmount.Total() {
		return fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", source.Name, source.availableCargoSpace(), transferAmount.Total(), dest.Name)
	}

	if !dest.Cargo.CanTransfer(transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %v from %s, there is not enough to transfer", source.Name, transferAmount, dest.Name)
	}

	// transfer the cargo
	source.Cargo = source.Cargo.Add(transferAmount)
	dest.Cargo = dest.Cargo.Subtract(transferAmount)

	return nil
}
