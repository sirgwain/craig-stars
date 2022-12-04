package game

// handle player orders
type orderer interface {
	updatePlanetOrders(player *Player, planet *Planet, productionQueue []ProductionQueueItem, contributesOnlyLeftoverToResearch bool)
	updateFleetOrders(player *Player, fleet *Fleet, orders FleetOrders)
}

type orders struct {

}

// update a planet orders
func (o *orders) updatePlanetOrders(player *Player, planet *Planet, productionQueue []ProductionQueueItem, contributesOnlyLeftoverToResearch bool) {
	planet.ContributesOnlyLeftoverToResearch = contributesOnlyLeftoverToResearch
	planet.ProductionQueue = productionQueue

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
