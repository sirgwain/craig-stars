package ai

import (
	"fmt"
	"math"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

func (ai *aiPlayer) produce() error {

	// check for builds of whole fleets
	for fleetPurpose, quantity := range ai.requests.fleetBuilds {
		if quantity <= 0 {
			continue
		}

		for _, planet := range ai.Planets {
			fleetMakeup := ai.fleetsByPurpose[fleetPurpose].clone()

			// make sure this planet is ready to build this type of fleet
			if !ai.isPlanetReadyToBuildFleet(planet, fleetMakeup.purpose) {
				continue
			}

			for shipIndex, ship := range fleetMakeup.ships {

				idleShips := ai.getIdleShipCount(planet, fleetPurpose, ship.purpose)
				if idleShips > 0 {
					quantityNeeded := ship.quantity - idleShips
					if quantityNeeded > 0 {
						// queue it on this planet
						fleetMakeup.ships[shipIndex].quantity = quantityNeeded
					} else {
						fleetMakeup.ships[shipIndex].quantity = 0
					}
				}
			}

			for _, ship := range fleetMakeup.ships {
				if ship.quantity <= 0 {
					continue
				}

				// design and upgrade this ship
				design, err := ai.designShip(ai.config.namesByPurpose[ship.purpose], ship.purpose, fleetMakeup.purpose)
				if err != nil {
					return fmt.Errorf("unable to design ship %v %w", ship.purpose, err)
				}
				if design == nil {
					log.Debug().
						Int64("GameID", ai.GameID).
						Int("PlayerNum", ai.Num).
						Msgf("unable to design ship %v", ship.purpose)
					continue
				}

				if !planet.CanBuild(design.Spec.Mass) {
					continue
				}

				if !ai.isShipInQueue(planet, fleetMakeup.purpose, ship.purpose, ship.quantity) {
					log.Debug().
						Int64("GameID", ai.GameID).
						Int("PlayerNum", ai.Num).
						Str("FleetPurpose", string(fleetMakeup.purpose)).
						Str("Purpose", string(ship.purpose)).
						Int("PlayerNum", ai.Num).
						Msgf("adding %d %s to %s queue", ship.quantity, design.Name, planet.Name)

					ai.addShipToTopOfQueue(planet, fleetMakeup.purpose, design, ship.quantity)
					if err := ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders); err != nil {
						return err
					}
				}
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

// add a production queue item to the top of the planet queue
func (ai *aiPlayer) addItemToTopOfQueue(planet *cs.Planet, t cs.QueueItemType, quantity int) {
	item := cs.ProductionQueueItem{Type: cs.QueueItemTypePlanetaryScanner, Quantity: quantity}
	planet.ProductionQueue = append([]cs.ProductionQueueItem{item}, planet.ProductionQueue...)
}

// add a production queue item to the top of the planet queue
func (ai *aiPlayer) addShipToTopOfQueue(planet *cs.Planet, purpose cs.FleetPurpose, design *cs.ShipDesign, quantity int) {
	item := cs.ProductionQueueItem{Type: cs.QueueItemTypeShipToken, Quantity: quantity, DesignNum: design.Num}
	item.WithTag(cs.TagPurpose, string(purpose))
	planet.ProductionQueue = append([]cs.ProductionQueueItem{item}, planet.ProductionQueue...)
}

// add a production queue item to the top of the planet queue
func (ai *aiPlayer) addStarbaseToTopOfQueue(planet *cs.Planet, design *cs.ShipDesign) {
	item := cs.ProductionQueueItem{Type: cs.QueueItemTypeStarbase, Quantity: 1, DesignNum: design.Num}
	planet.ProductionQueue = append([]cs.ProductionQueueItem{item}, planet.ProductionQueue...)

	log.Debug().
		Int64("GameID", ai.GameID).
		Int("PlayerNum", ai.Num).
		Msgf("Planet %s added %s to production queue", planet.Name, design.Name)

}

// get the years to build a certain number of items
func (ai *aiPlayer) getYearsToBuild(planet *cs.Planet, t cs.QueueItemType, quantity int) (int, error) {
	yearlyAvailableToSpend := cs.FromMineralAndResources(planet.Spec.MiningOutput, planet.Spec.ResourcesPerYearAvailable)
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
	yearlyAvailableToSpend := cs.FromMineralAndResources(planet.Spec.MiningOutput, planet.Spec.ResourcesPerYearAvailable)
	costCalculator := cs.NewCostCalculator()
	completionEstimator := cs.NewCompletionEstimator()
	item := cs.ProductionQueueItem{Type: cs.QueueItemTypeStarbase, Quantity: 1, DesignNum: design.Num}
	item.SetDesign(design)

	var err error
	var cost cs.Cost
	if planet.Spec.HasStarbase {
		existingStarbase := ai.GetDesign(planet.Starbase.Tokens[0].DesignNum)
		cost = costCalculator.StarbaseUpgradeCost(&ai.game.Rules, ai.Player.TechLevels, ai.Player.Race.Spec, existingStarbase, design)
	} else {
		cost = design.Spec.Cost
		if err != nil {
			return math.MaxInt, fmt.Errorf("calculate starbase cost %w", err)
		} else {
			return int(math.Ceil(cost.Divide(yearlyAvailableToSpend))), nil
		}
	}

	// if we can complete this soon, queue it
	yearsToBuild := completionEstimator.GetYearsToBuildOne(item, cost, planet.Spec.MiningOutput, yearlyAvailableToSpend)
	// log.Debug().
	// 	Int64("GameID", ai.GameID).
	// 	Int("PlayerNum", ai.Num).
	// 	Msgf("Planet %s would take %d years to build %s", planet.Name, yearsToBuild, design.Name)

	// make our conditionals easier
	if yearsToBuild == cs.Infinite {
		yearsToBuild = math.MaxInt
	}

	return yearsToBuild, nil
}
