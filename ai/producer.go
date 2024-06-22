package ai

import (
	"fmt"
	"math"
	"slices"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type fleetBuildRequest struct {
	fleet       fleet
	buildPlanet *cs.Planet
	target      *cs.MapObjectTarget
}

func (ai *aiPlayer) produce() error {

	// first queue planet specific build orders
	// this would be for things like colonizers that are requested per planet
	planetsDoneBuilding := make(map[int]bool, len(ai.planetsWithDocks))
	for _, request := range ai.fleetBuildRequests {
		if request.buildPlanet == nil {
			continue
		}
		if planetsDoneBuilding[request.buildPlanet.Num] {
			continue
		}

		queued, err := ai.produceFleetOnPlanet(request.buildPlanet, request)
		if err != nil {
			return err
		}

		// this planet can't build anymore
		if !queued {
			planetsDoneBuilding[request.buildPlanet.Num] = true
		}
	}

	// now produce any fleets that can be built by any world
	for _, request := range ai.fleetBuildRequests {
		if request.buildPlanet != nil {
			continue
		}
		for _, planet := range ai.Planets {
			queuedAllShips, err := ai.produceFleetOnPlanet(planet, request)
			if err != nil {
				return err
			}

			// this planet was able to queue all ships in this request, break
			// and move on to the next request
			if queuedAllShips {
				break
			}
		}
	}

	// build scanners and starbases on each planet, where applicable
	for _, planet := range ai.Planets {

		if !planet.Scanner && !ai.isItemInQueue(planet, cs.QueueItemTypePlanetaryScanner) {
			yearsToBuild, err := ai.getYearsToBuild(planet, cs.QueueItemTypePlanetaryScanner, 1)
			if err != nil {
				return err
			}
			if yearsToBuild <= ai.config.minYearsToBuildScanner {
				ai.addItemToTopOfQueue(planet, cs.QueueItemTypePlanetaryScanner, 1)
			}
		}

		// see if we should build a starbase or upgrade an existing one
		if err := ai.buildOrUpgradeStarbase(planet); err != nil {
			return err
		}

	}
	return nil
}

func (ai *aiPlayer) produceFleetOnPlanet(planet *cs.Planet, request fleetBuildRequest) (bool, error) {
	fleetMakeup := request.fleet.clone()

	// make sure this planet is ready to build this type of fleet
	if !ai.isPlanetReadyToBuildFleet(planet, fleetMakeup.purpose) {
		return false, nil
	}

	// remove any idle ships that match this makeup that are already in orbit around the planet
	for _, ship := range fleetMakeup.ships {

		// ensure we actually have a design for this ship
		// if not, log it and break
		if ship.design == nil {
			log.Debug().
				Int64("GameID", ai.GameID).
				Int("PlayerNum", ai.Num).
				Msgf("unable to design ship %v", ship.purpose)
			return false, nil
		}

		// ensure this planet can build this type of ship
		if !planet.CanBuild(ship.design.Spec.Mass) {
			return false, nil
		}

		if fleetMakeup.onlyQueueOnePerPlanet && ai.isShipInQueue(planet, fleetMakeup.purpose, ship.purpose, ship.quantity) {
			return false, nil
		}
	}

	yearsToBuild, err := ai.getYearsToBuildShips(planet, fleetMakeup)
	if err != nil {
		return false, err
	}

	// some fleets, like colonizers, we want to only build if they build quickly
	// other fleets, like bombers, we don't want to use all our production up by queuing them
	if fleetMakeup.mustBuildInYears != 0 && yearsToBuild > fleetMakeup.mustBuildInYears {
		return false, nil
	}

	// we've made it here, queue this fleet's ships
	for _, ship := range fleetMakeup.ships {
		log.Debug().
			Int64("GameID", ai.GameID).
			Int("PlayerNum", ai.Num).
			Str("FleetPurpose", string(fleetMakeup.purpose)).
			Str("Purpose", string(ship.purpose)).
			Int("PlayerNum", ai.Num).
			Msgf("adding %d %s to %s queue", ship.quantity, ship.design.Name, planet.Name)

		planet.ProductionQueue, _ = ai.addShipToQueue(planet.ProductionQueue, fleetMakeup.purpose, ship.design, ship.quantity, request.target)
		if err := ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders); err != nil {
			return false, err
		}
	}

	// return true if this planet built all needed ships in the request
	return true, nil
}

// add a new ship build request per count
func (ai *aiPlayer) addFleetBuildRequest(purpose cs.FleetPurpose, count int) {
	for i := 0; i < count; i++ {
		ai.fleetBuildRequests = append(ai.fleetBuildRequests, fleetBuildRequest{fleet: ai.fleetsByPurpose[purpose]})
	}
}

// check if an item type is in the queue
func (ai *aiPlayer) isItemInQueue(planet *cs.Planet, t cs.QueueItemType) bool {
	for _, item := range planet.ProductionQueue {
		if item.Type == t {
			return true
		}
	}
	return false
}

// check if a planet is in a state where it will build a fleet
func (ai *aiPlayer) isPlanetReadyToBuildFleet(planet *cs.Planet, purpose cs.FleetPurpose) bool {
	if !planet.Spec.HasStarbase || planet.Spec.DockCapacity == 0 {
		return false
	}

	// don't worry about production, let this planet build scouts
	if purpose == cs.FleetPurposeScout {
		return true
	}

	// don't build colonizers or freighters on this planet if it doesn't have the pop to move them
	if purpose == cs.FleetPurposeColonistFreighter &&
		planet.Spec.PopulationDensity >= ai.config.colonistTransportDensity {
		return true
	}
	if purpose == cs.FleetPurposeColonizer &&
		planet.Spec.PopulationDensity >= ai.config.colonizerPopulationDensity {
		return true
	}

	// don't build certain things unless we meet some requirements
	planetaryStructuresBuilt := math.Min(float64(planet.Mines)/float64(planet.Spec.MaxMines), float64(planet.Factories)/float64(planet.Spec.MaxFactories))

	// bombers require the planet to be very mature
	if purpose == cs.FleetPurposeBomber {
		if planetaryStructuresBuilt < ai.config.bomberProductionCutoff {
			return false
		}
	}

	// make sure our planetary productivity is still in line to build fleets
	if planetaryStructuresBuilt < ai.config.fleetProductionCutoff {
		return false
	}

	return true
}

// check the existing starbase and build or upgrade it
func (ai *aiPlayer) buildOrUpgradeStarbase(planet *cs.Planet) error {
	// if we're already building a starbase, don't do anything
	if ai.isStarbaseInQueue(planet) {
		return nil
	}

	// check if we are targetted or being bombed by bombers
	enemyOrbitingFleets := ai.enemyShipsAbovePlanet(planet)
	attackShipsInOrbit := ai.hasAttackShips(enemyOrbitingFleets)
	_, targetted := ai.targetedPlanets[planet.Num]

	// don't build starbases unless this planet has not moved forward enough economically
	// if we are being targetted for bombing though, we want to try and build a starbase
	planetaryStructuresBuilt := math.Min(float64(planet.Mines)/float64(planet.Spec.MaxMines), float64(planet.Factories)/float64(planet.Spec.MaxFactories))
	if !(targetted || attackShipsInOrbit) && planetaryStructuresBuilt < ai.config.fleetProductionCutoff {
		return nil
	}

	timeToWait := ai.config.minYearsToQueueStarbasePeaceTime
	if attackShipsInOrbit {
		timeToWait = ai.config.minYearsToQueueStarbaseWarTime
	}

	if targetted || attackShipsInOrbit {
		// this planet is being threatened
		if planet.Spec.HasStarbase {
			ai.upgradeStarbase(planet, timeToWait)
		} else {
			yearsToBuild, err := ai.getYearsToBuildStarbase(planet, ai.fortDesign)
			if err != nil {
				return err
			}

			if yearsToBuild < ai.config.minYearsToBuildFort {
				planet.ProductionQueue, _ = ai.addStarbaseToQueue(planet.ProductionQueue, ai.fortDesign)
			}
		}
	} else {
		if planet.Spec.HasStarbase {
			ai.upgradeStarbase(planet, timeToWait)
		} else {
			yearsToBuild, err := ai.getYearsToBuildStarbase(planet, ai.fuelDepotDesign)
			if err != nil {
				return err
			}
			if yearsToBuild <= timeToWait {
				planet.ProductionQueue, _ = ai.addStarbaseToQueue(planet.ProductionQueue, ai.fuelDepotDesign)
			}
		}
	}

	return nil
}

// upgrade an existing starbase to a better or newer model
func (ai *aiPlayer) upgradeStarbase(planet *cs.Planet, timeToWait int) error {
	existingDesign := ai.GetDesign(planet.Spec.StarbaseDesignNum)
	if existingDesign == nil {
		err := fmt.Errorf("failed to find existing starbase design")
		log.Err(err).
			Int64("GameID", ai.GameID).
			Int("PlayerNum", ai.Num).
			Int("PlanetNum", planet.Num).
			Str("PlanetName", planet.Name).
			Int("DesignNum", planet.Spec.StarbaseDesignNum).
			Str("DesignName", planet.Spec.StarbaseDesignName).
			Msgf("design not found")

		return err
	}
	if existingDesign.Purpose == cs.ShipDesignPurposeFort || existingDesign.Purpose == cs.ShipDesignPurposeFuelDepot {
		// try and upgrade our fort/fueldepot to a quarter filled out starbase
		yearsToBuild, err := ai.getYearsToBuildStarbase(planet, ai.starbaseQuarterDesign)
		if err != nil {
			return err
		}
		if yearsToBuild <= timeToWait {
			planet.ProductionQueue, _ = ai.addStarbaseToQueue(planet.ProductionQueue, ai.starbaseQuarterDesign)
		}
	} else {
		switch existingDesign.Purpose {
		case cs.ShipDesignPurposeStarbaseQuarter:
			// upgrade 1/4 -> 1/2
			yearsToBuild, err := ai.getYearsToBuildStarbase(planet, ai.starbaseHalfDesign)
			if err != nil {
				return err
			}
			if yearsToBuild <= timeToWait {
				planet.ProductionQueue, _ = ai.addStarbaseToQueue(planet.ProductionQueue, ai.starbaseHalfDesign)
			}
		case cs.ShipDesignPurposeStarbaseHalf:
			// upgrade 1/2 -> full
			yearsToBuild, err := ai.getYearsToBuildStarbase(planet, ai.starbaseDesign)
			if err != nil {
				return err
			}
			if yearsToBuild <= timeToWait {
				planet.ProductionQueue, _ = ai.addStarbaseToQueue(planet.ProductionQueue, ai.starbaseDesign)
			}

		case cs.ShipDesignPurposeStarbase:
			if existingDesign.Num != ai.starbaseDesign.Num {
				// we have a new full design, check for upgrade
				yearsToBuild, err := ai.getYearsToBuildStarbase(planet, ai.starbaseDesign)
				if err != nil {
					return err
				}
				if yearsToBuild <= timeToWait {
					planet.ProductionQueue, _ = ai.addStarbaseToQueue(planet.ProductionQueue, ai.starbaseDesign)
				}
			}
		}
	}

	return nil
}

// add a production queue item to the top of the planet queue
func (ai *aiPlayer) addItemToTopOfQueue(planet *cs.Planet, t cs.QueueItemType, quantity int) {
	item := cs.ProductionQueueItem{Type: t, Quantity: quantity}
	planet.ProductionQueue = append([]cs.ProductionQueueItem{item}, planet.ProductionQueue...)
}

// add a production queue item to the bottom of the queue, but before auto items
func (ai *aiPlayer) addShipToQueue(queue []cs.ProductionQueueItem, purpose cs.FleetPurpose, design *cs.ShipDesign, quantity int, target *cs.MapObjectTarget) ([]cs.ProductionQueueItem, int) {
	item := cs.ProductionQueueItem{Type: cs.QueueItemTypeShipToken, Quantity: quantity, DesignNum: design.Num}
	item = item.WithTag(cs.TagPurpose, string(purpose))
	if target != nil {
		item = item.WithTag(cs.TagTarget, target.String())
	}

	// add the item before the first auto
	index := slices.IndexFunc(queue, func(item cs.ProductionQueueItem) bool { return item.Type.IsAuto() })
	if index == -1 {
		index = len(queue)
	}
	return slices.Insert(queue, index, item), index
}

// add a production queue item to the bottom of the queue, but before auto items
func (ai *aiPlayer) addStarbaseToQueue(queue []cs.ProductionQueueItem, design *cs.ShipDesign) ([]cs.ProductionQueueItem, int) {
	item := cs.ProductionQueueItem{Type: cs.QueueItemTypeStarbase, Quantity: 1, DesignNum: design.Num}

	// add the item before the first auto
	index := slices.IndexFunc(queue, func(item cs.ProductionQueueItem) bool { return item.Type.IsAuto() })
	if index == -1 {
		index = len(queue)
	}
	return slices.Insert(queue, index, item), index
}

// get the years to build a certain number of items
func (ai *aiPlayer) getYearsToBuild(planet *cs.Planet, t cs.QueueItemType, quantity int) (int, error) {
	yearlyAvailableToSpend := cs.FromMineralAndResources(planet.Spec.MiningOutput, planet.Spec.ResourcesPerYearAvailable)
	costCalculator := cs.NewCostCalculator()
	completionEstimator := cs.NewCompletionEstimator()

	item := cs.ProductionQueueItem{Type: t, Quantity: 1}
	cost, err := costCalculator.CostOfOne(ai.Player, item)
	if err != nil {
		return 0, err
	}

	// get the years to build one of these
	yearsToBuild := completionEstimator.GetYearsToBuildOne(item, cost, planet.Spec.MiningOutput, yearlyAvailableToSpend)

	// make our conditionals easier
	if yearsToBuild == cs.Infinite {
		yearsToBuild = math.MaxInt
	}
	return yearsToBuild, nil
}

// get the years it will take to build or upgrade to this starbase
func (ai *aiPlayer) getYearsToBuildStarbase(planet *cs.Planet, design *cs.ShipDesign) (int, error) {
	queue := make([]cs.ProductionQueueItem, len(planet.ProductionQueue), len(planet.ProductionQueue)+1)
	copy(queue, planet.ProductionQueue)

	queue, index := ai.addStarbaseToQueue(queue, design)

	// estsimate how long it will take to build our last item added to the queue
	completionEstimator := cs.NewCompletionEstimator()
	estimaterPlanet := *planet
	estimaterPlanet.ProductionQueue = queue
	if err := estimaterPlanet.PopulateProductionQueueDesigns(ai.Player); err != nil {
		return 0, err
	}

	// see how many years it'll take to build this ship
	queue, _ = completionEstimator.GetProductionWithEstimates(&ai.game.Rules, ai.Player, estimaterPlanet)
	// log.Debug().
	// 	Int64("GameID", ai.GameID).
	// 	Int("PlayerNum", ai.Num).
	// 	Msgf("Planet %s would take %d years to build %s", planet.Name, yearsToBuild, design.Name)

	if index >= 0 && index < len(queue) {
		yearsToBuild := queue[index].YearsToBuildAll
		if yearsToBuild == cs.Infinite {
			return math.MaxInt, nil
		}
		return yearsToBuild, nil
	}
	return 0, fmt.Errorf("failed to estimate years to build starbase, item index not found")

}

// get the years it will take to build ships
func (ai *aiPlayer) getYearsToBuildShips(planet *cs.Planet, fleetMakup fleet) (int, error) {

	queue := make([]cs.ProductionQueueItem, len(planet.ProductionQueue), len(planet.ProductionQueue)+len(fleetMakup.ships))
	copy(queue, planet.ProductionQueue)

	var index = -1
	for _, ship := range fleetMakup.ships {
		if ship.quantity > 0 {
			queue, index = ai.addShipToQueue(queue, cs.FleetPurposeNone, ship.design, ship.quantity, nil)
		}
	}

	// estsimate how long it will take to build our last item added to the queue
	completionEstimator := cs.NewCompletionEstimator()
	estimaterPlanet := *planet
	estimaterPlanet.ProductionQueue = queue
	if err := estimaterPlanet.PopulateProductionQueueDesigns(ai.Player); err != nil {
		return 0, err
	}

	// see how many years it'll take to build this ship
	queue, _ = completionEstimator.GetProductionWithEstimates(&ai.game.Rules, ai.Player, estimaterPlanet)
	// log.Debug().
	// 	Int64("GameID", ai.GameID).
	// 	Int("PlayerNum", ai.Num).
	// 	Msgf("Planet %s would take %d years to build %s", planet.Name, yearsToBuild, design.Name)

	if index >= 0 && index < len(queue) {
		yearsToBuild := queue[index].YearsToBuildAll
		if yearsToBuild == cs.Infinite {
			return math.MaxInt, nil
		}
		return yearsToBuild, nil
	}
	return 0, fmt.Errorf("failed to estimate years to build ship, item index not found")
}
