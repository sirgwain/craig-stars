package ai

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/sirgwain/craig-stars/cs"
)

func (ai *aiPlayer) produce() error {
	type shipBuild struct {
		design   *cs.ShipDesign
		quantity int
	}

	// Allocate ship production for each requested fleet purpose.
	for fleetPurpose, quantity := range ai.requests.fleetBuilds {
		if quantity <= 0 {
			continue
		}

		fleetMakeup := ai.fleetsByPurpose[fleetPurpose]
		if len(fleetMakeup.requiredShips()) == 0 {
			continue
		}

		// Managers already account for complete idle fleets. Count only the
		// additional fleets supplied by pending production against this demand.
		for _, planet := range ai.Planets {
			idle, available := ai.fleetBuildShipCounts(planet, fleetMakeup)
			quantity -= fleetMakeup.countCompleteFleets(available) - fleetMakeup.countCompleteFleets(idle)
		}

		// Use shipyards with existing orders first, then consider other shipyards.
		for _, hasQueuedFleet := range []bool{true, false} {
			for _, planet := range ai.Planets {
				if quantity <= 0 {
					break
				}
				if ai.isFleetInQueue(planet, fleetPurpose) != hasQueuedFleet || !ai.isPlanetReadyToBuildFleet(planet, fleetPurpose) {
					continue
				}

				// Find the ships needed for one more fleet using local idle ships and orders.
				_, available := ai.fleetBuildShipCounts(planet, fleetMakeup)
				fleetsAtPlanet := fleetMakeup.countCompleteFleets(available)
				required := fleetMakeup.requiredShips()
				builds := make([]shipBuild, 0, len(fleetMakeup.ships))
				canBuild := true
				for _, ship := range fleetMakeup.ships {
					shipQuantity, found := required[ship.purpose]
					if !found {
						continue
					}
					// Handle each ship purpose once, including duplicate recipe entries.
					delete(required, ship.purpose)
					needed := (fleetsAtPlanet+1)*shipQuantity - available[ship.purpose]
					if needed <= 0 {
						continue
					}

					// Check every missing design before adding any orders at this shipyard.
					design, err := ai.designShip(ai.config.namesByPurpose[ship.purpose], ship.purpose, fleetPurpose)
					if err != nil {
						return fmt.Errorf("unable to design ship %v: %w", ship.purpose, err)
					}
					if design == nil || !planet.CanBuild(design.Spec.Mass) {
						canBuild = false
						break
					}
					builds = append(builds, shipBuild{design, needed})
				}
				if !canBuild {
					continue
				}

				// Queue the missing ships together so they can assemble at this planet.
				for _, build := range builds {
					ai.log.Debug("adding ship to queue",
						slog.String("FleetPurpose", string(fleetPurpose)),
						slog.String("Purpose", string(build.design.Purpose)),
						slog.Int("PlayerNum", ai.Num),
						slog.Int("quantity", build.quantity),
						slog.String("design", build.design.Name),
						slog.String("planet", planet.Name))
					ai.addShipToTopOfQueue(planet, fleetPurpose, build.design, build.quantity)
				}
				if err := ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders); err != nil {
					return err
				}
				// This shipyard now supplies one of the requested fleets.
				quantity--
			}
		}
	}

	// Add planetary infrastructure independently of fleet requests.
	for _, planet := range ai.Planets {

		// Add a missing scanner if it can be built within the configured time.
		if !planet.Scanner && !ai.isItemInQueue(planet, cs.QueueItemTypePlanetaryScanner) {
			yearsToBuild, err := ai.getYearsToBuild(planet)
			if err != nil {
				return err
			}
			if yearsToBuild <= ai.config.minYearsToBuildScanner {
				ai.addItemToTopOfQueue(planet, 1)
			}
		}

		// Build or upgrade the starbase according to development and threats.
		if err := ai.buildOrUpgradeStarbase(planet); err != nil {
			return err
		}

	}
	return nil
}

// fleetBuildShipCounts includes all design versions and queue rows assigned to
// this fleet purpose. Ships at different planets cannot form the same fleet.
func (ai *aiPlayer) fleetBuildShipCounts(planet *cs.Planet, makeup fleet) (idle, available map[cs.ShipDesignPurpose]int) {
	idle = make(map[cs.ShipDesignPurpose]int, len(makeup.ships))
	available = make(map[cs.ShipDesignPurpose]int, len(makeup.ships))
	for purpose := range makeup.requiredShips() {
		idle[purpose] = ai.getIdleShipCount(planet, makeup.purpose, purpose)
		available[purpose] = idle[purpose]
	}
	for _, item := range planet.ProductionQueue {
		if item.Type != cs.QueueItemTypeShipToken || item.GetTag(cs.TagPurpose) != string(makeup.purpose) {
			continue
		}
		if design := ai.GetDesign(item.DesignNum); design != nil {
			available[design.Purpose] += item.Quantity
		}
	}
	return idle, available
}

// add a new ship build request
func (ai *aiPlayer) addFleetBuildRequest(purpose cs.FleetPurpose, count int) {
	current := ai.requests.fleetBuilds[purpose]
	ai.requests.fleetBuilds[purpose] = current + count
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
	planetaryStructuresBuilt := min(float64(planet.Mines)/float64(planet.Spec.MaxMines), float64(planet.Factories)/float64(planet.Spec.MaxFactories))

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

	// check if we are targeted or being bombed by bombers
	enemyOrbitingFleets := ai.enemyShipsAbovePlanet(planet)
	attackShipsInOrbit := ai.hasAttackShips(enemyOrbitingFleets)
	_, targeted := ai.targetedPlanets[planet.Num]

	// don't build starbases if this planet has not moved forward enough economically
	// if we are being targeted for bombing though, we want to try and build a starbase regardless
	// TODO: Add ability to build fuel depots and infrastructure based on a (lower) cutoff
	// This will be useful for IT/PP and desperately necessary for AR
	if !(targeted || attackShipsInOrbit) {
		if ai.Player.Race.Spec.InnateResources {
			if planet.Spec.CanTerraform {
				// Terraforming is economic development for AR planets
				return nil
			}
			existingDesign := ai.GetDesign(planet.Spec.StarbaseDesignNum)
			if len(existingDesign.Slots) == 0 && existingDesign != ai.fuelDepotDesign {
				// Bigger starbase allows bigger population
				ai.addStarbaseToTopOfQueue(planet, ai.fuelDepotDesign)
			} else {
				ai.upgradeStarbase(planet, ai.config.minYearsToQueueStarbasePeaceTime)
			}
			return nil
		} else {
			planetaryStructuresBuilt := min(float64(planet.Mines)/float64(planet.Spec.MaxMines), float64(planet.Factories)/float64(planet.Spec.MaxFactories))
			if planetaryStructuresBuilt < ai.config.fleetProductionCutoff {
				return nil
			}
		}
	}

	timeToWait := ai.config.minYearsToQueueStarbasePeaceTime
	if attackShipsInOrbit {
		timeToWait = ai.config.minYearsToQueueStarbaseWarTime
	}

	if targeted || attackShipsInOrbit {
		// this planet is being threatened
		if planet.Spec.HasStarbase {
			ai.upgradeStarbase(planet, timeToWait)
		} else {
			yearsToBuild, err := ai.getYearsToBuildStarbase(planet, ai.fortDesign)
			if err != nil {
				return err
			}

			if yearsToBuild <= ai.config.minYearsToBuildFort {
				ai.addStarbaseToTopOfQueue(planet, ai.fortDesign)
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
				ai.addStarbaseToTopOfQueue(planet, ai.fuelDepotDesign)
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
		slog.Error("design not found",
			slog.Any("error", err),
			slog.Int64("GameID", ai.GameID),
			slog.Int("PlayerNum", ai.Num),
			slog.Int("PlanetNum", planet.Num),
			slog.String("PlanetName", planet.Name),
			slog.Int("DesignNum", planet.Spec.StarbaseDesignNum),
			slog.String("DesignName", planet.Spec.StarbaseDesignName))

		return err
	}
	if existingDesign.Purpose == cs.ShipDesignPurposeFort || existingDesign.Purpose == cs.ShipDesignPurposeFuelDepot || existingDesign.Purpose == cs.ShipDesignPurposeStarbaseUnarmed {
		// try and upgrade our fort/fuel depot to a quarter filled out starbase
		yearsToBuild, err := ai.getYearsToBuildStarbase(planet, ai.starbaseQuarterDesign)
		if err != nil {
			return err
		}
		if yearsToBuild <= timeToWait {
			ai.addStarbaseToTopOfQueue(planet, ai.starbaseQuarterDesign)
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
				ai.addStarbaseToTopOfQueue(planet, ai.starbaseHalfDesign)
			}
		case cs.ShipDesignPurposeStarbaseHalf:
			// upgrade 1/2 -> full
			yearsToBuild, err := ai.getYearsToBuildStarbase(planet, ai.starbaseDesign)
			if err != nil {
				return err
			}
			if yearsToBuild <= timeToWait {
				ai.addStarbaseToTopOfQueue(planet, ai.starbaseDesign)
			}

		case cs.ShipDesignPurposeStarbase:
			if existingDesign.Num != ai.starbaseDesign.Num {
				// we have a new full design, check for upgrade
				yearsToBuild, err := ai.getYearsToBuildStarbase(planet, ai.starbaseDesign)
				if err != nil {
					return err
				}
				if yearsToBuild <= timeToWait {
					ai.addStarbaseToTopOfQueue(planet, ai.starbaseDesign)
				}
			}
		}
	}

	return nil
}

// add a normal production queue item to the top of the planet queue
func (ai *aiPlayer) addItemToTopOfQueue(planet *cs.Planet, quantity int) {
	item := cs.ProductionQueueItem{Type: cs.QueueItemTypePlanetaryScanner, Quantity: quantity}
	planet.ProductionQueue = append([]cs.ProductionQueueItem{item}, planet.ProductionQueue...)
}

// add one or more ships to the top of a planet's production queue
func (ai *aiPlayer) addShipToTopOfQueue(planet *cs.Planet, purpose cs.FleetPurpose, design *cs.ShipDesign, quantity int) {
	item := cs.ProductionQueueItem{Type: cs.QueueItemTypeShipToken, Quantity: quantity, DesignNum: design.Num}
	item.WithTag(cs.TagPurpose, string(purpose))
	planet.ProductionQueue = append([]cs.ProductionQueueItem{item}, planet.ProductionQueue...)
}

// add a starbase to the top of a planet's production queue
func (ai *aiPlayer) addStarbaseToTopOfQueue(planet *cs.Planet, design *cs.ShipDesign) {
	item := cs.ProductionQueueItem{Type: cs.QueueItemTypeStarbase, Quantity: 1, DesignNum: design.Num}
	planet.ProductionQueue = append([]cs.ProductionQueueItem{item}, planet.ProductionQueue...)

	ai.log.Debug("added starbase to production queue",
		slog.Int64("GameID", ai.GameID),
		slog.Int("PlayerNum", ai.Num),
		slog.String("planet", planet.Name),
		slog.String("design", design.Name))

}

// get the years to build a certain number of items
func (ai *aiPlayer) getYearsToBuild(planet *cs.Planet) (int, error) {
	yearlyAvailableToSpend := cs.NewCostFromMineralAndResources(planet.Spec.MiningOutput, planet.Spec.ResourcesPerYearAvailable)
	costCalculator := cs.NewCostCalculator()
	completionEstimator := cs.NewCompletionEstimator()

	item := cs.ProductionQueueItem{Type: cs.QueueItemTypePlanetaryScanner, Quantity: 1}
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
	yearlyAvailableToSpend := cs.NewCostFromMineralAndResources(planet.Spec.MiningOutput, planet.Spec.ResourcesPerYearAvailable)
	costCalculator := cs.NewCostCalculator()
	completionEstimator := cs.NewCompletionEstimator()
	item := cs.ProductionQueueItem{Type: cs.QueueItemTypeStarbase, Quantity: 1, DesignNum: design.Num}
	item.SetDesign(design)

	var err error
	var cost cs.Cost
	if planet.Spec.HasStarbase {
		existingStarbase := ai.GetDesign(planet.Spec.StarbaseDesignNum)
		cost, err = costCalculator.StarbaseUpgradeCost(&ai.game.Rules, ai.Player.TechLevels, ai.Player.Race.Spec, existingStarbase, design)
	} else {
		cost, err = costCalculator.GetDesignCost(&ai.game.Rules, ai.Player.TechLevels, ai.Player.Race.Spec, design)
	}
	if err != nil {
		return math.MaxInt, fmt.Errorf("calculate starbase cost: %w", err)
	}

	// calculate how long it take to build
	yearsToBuild := completionEstimator.GetYearsToBuildOne(item, cost, planet.Spec.MiningOutput, yearlyAvailableToSpend)
	// ai.log.Debug().
	// 	Int64("GameID", ai.GameID).
	// 	Int("PlayerNum", ai.Num).
	// 	Msgf("Planet %s would take %d years to build %s", planet.Name, yearsToBuild, design.Name)

	// make our conditionals easier
	if yearsToBuild == cs.Infinite {
		yearsToBuild = math.MaxInt
	}

	return yearsToBuild, nil
}
