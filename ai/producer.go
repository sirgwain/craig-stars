package ai

import (
	"fmt"
	"math"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

func (ai *aiPlayer) produce() error {
	costCalculator := cs.NewCostCalculator()
	completionEstimator := cs.NewCompletionEstimator()

	// check for builds of whole fleets
	for fleetPurpose, quantity := range ai.requests.fleetBuilds {
		if quantity <= 0 {
			continue
		}

		for _, planet := range ai.getPlanetsWithDocks() {
			fleetMakeup := ai.fleetsByPurpose[fleetPurpose].clone()
			for shipIndex, ship := range fleetMakeup.ships {

				idleShips := ai.getIdleShipCount(planet, ship.purpose)
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

				// don't build colonizers or freighters on this planet if it doesn't have the pop to remove them
				if (ship.purpose == cs.ShipDesignPurposeColonistFreighter || ship.purpose == cs.ShipDesignPurposeColonizer) &&
					planet.Spec.PopulationDensity < ai.config.colonizerPopulationDensity {
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

					// put this ship at the top of the queue
					planet.ProductionQueue = append([]cs.ProductionQueueItem{
						*cs.NewProductionQueueItemShip(ship.quantity, design).WithTag("purpose", string(fleetMakeup.purpose)),
					}, planet.ProductionQueue...)
					if err := ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders); err != nil {
						return err
					}
				}
			}
		}
	}

	for _, planet := range ai.Planets {
		// don't build starbases unless this planet has half its economy built
		planetaryStructuresBuilt := math.Min(float64(planet.Mines) / float64(planet.Spec.MaxMines), float64(planet.Factories) / float64(planet.Spec.MaxFactories))
		if planetaryStructuresBuilt < ai.config.fleetProductionCutoff {
			continue
		}

		yearlyAvailableToSpend := cs.FromMineralAndResources(planet.Spec.MiningOutput, planet.Spec.ResourcesPerYearAvailable)

		// try and build starbase if we don't have one built or queued
		if !ai.isStarbaseInQueue(planet) {

			enemyOrbitingFleets := ai.enemyShipsAbovePlanet(planet)
			attackShipsInOrbit := ai.hasAttackShips(enemyOrbitingFleets)

			// if we don't have a starbase yet, try and build one
			if !(planet.Spec.HasStarbase) {

				if attackShipsInOrbit {
					// try and build a half, quarter, or fort starbase to repel these attacking ships
					queued, err := ai.checkPlanetStarbaseBuild(planet, cs.ShipDesignPurposeStarbaseHalf, ai.config.minYearsToQueueStarbaseWarTime, costCalculator, yearlyAvailableToSpend, completionEstimator)
					if err != nil {
						return err
					}
					if !queued {
						// try an orbital fort
						_, err := ai.checkPlanetStarbaseBuild(planet, cs.ShipDesignPurposeStarbaseQuarter, ai.config.minYearsToQueueStarbaseWarTime, costCalculator, yearlyAvailableToSpend, completionEstimator)
						if err != nil {
							return err
						}
					}
					if !queued {
						// try an orbital fort
						_, err := ai.checkPlanetStarbaseBuild(planet, cs.ShipDesignPurposeFort, ai.config.minYearsToQueueStarbaseWarTime, costCalculator, yearlyAvailableToSpend, completionEstimator)
						if err != nil {
							return err
						}
					}
				} else {
					// try a fuel depot
					_, err := ai.checkPlanetStarbaseBuild(planet, cs.ShipDesignPurposeFuelDepot, ai.config.minYearsToQueueStarbasePeaceTime, costCalculator, yearlyAvailableToSpend, completionEstimator)
					if err != nil {
						return err
					}
				}
			} else {
				design := ai.GetDesign(planet.Spec.StarbaseDesignNum)
				if design.Purpose == cs.ShipDesignPurposeFort || design.Purpose == cs.ShipDesignPurposeFuelDepot {
					// we have a fort or fuel depot, try and upgrade
					timeToWait := ai.config.minYearsToQueueStarbasePeaceTime
					if attackShipsInOrbit {
						timeToWait = ai.config.minYearsToQueueStarbaseWarTime
					}

					// try and build a full
					queued, err := ai.checkPlanetStarbaseBuild(planet, cs.ShipDesignPurposeStarbase, timeToWait, costCalculator, yearlyAvailableToSpend, completionEstimator)
					if err != nil {
						return err
					}

					// try and build a half starbase
					if !queued {
						// try an orbital fort
						_, err := ai.checkPlanetStarbaseBuild(planet, cs.ShipDesignPurposeStarbaseHalf, ai.config.minYearsToQueueStarbaseWarTime, costCalculator, yearlyAvailableToSpend, completionEstimator)
						if err != nil {
							return err
						}
					}

					// try and build a quarter starbase
					if !queued {
						// try an orbital fort
						_, err := ai.checkPlanetStarbaseBuild(planet, cs.ShipDesignPurposeStarbaseQuarter, ai.config.minYearsToQueueStarbaseWarTime, costCalculator, yearlyAvailableToSpend, completionEstimator)
						if err != nil {
							return err
						}
					}
				}
			}
		}

	}
	return nil
}

func (ai *aiPlayer) checkPlanetStarbaseBuild(planet *cs.Planet, purpose cs.ShipDesignPurpose, maxYearsToWait int, costCalculator cs.CostCalculator, yearlyAvailableToSpend cs.Cost, completionEstimator cs.CompletionEstimator) (bool, error) {
	// design and upgrade this ship
	design, err := ai.designShip(ai.config.namesByPurpose[purpose], purpose, cs.FleetPurposeFromShipDesignPurpose(purpose))
	if err != nil {
		return false, fmt.Errorf("unable to design ship %v %w", purpose, err)
	}

	item := cs.ProductionQueueItem{Type: cs.QueueItemTypeStarbase, Quantity: 1, DesignNum: design.Num}
	item.SetDesign(design)

	var cost cs.Cost

	if planet.Spec.HasStarbase {
		existingStarbase := ai.GetDesign(planet.Starbase.Tokens[0].DesignNum)
		cost = costCalculator.StarbaseUpgradeCost(design, existingStarbase)
	} else {
		cost, err = costCalculator.CostOfOne(ai.Player, item)
		if err != nil {
			return false, fmt.Errorf("calculate starbase cost %w", err)
		}
	}

	item.CostOfOne = cost

	// if we can complete this soon, queue it
	yearsToBuild := completionEstimator.GetYearsToBuildOne(item, planet.Spec.MiningOutput, yearlyAvailableToSpend)
	log.Debug().
		Int64("GameID", ai.GameID).
		Int("PlayerNum", ai.Num).
		Msgf("Planet %s would take %d years to build %s", planet.Name, yearsToBuild, design.Name)

	if yearsToBuild != cs.Infinite && yearsToBuild <= maxYearsToWait {
		log.Debug().
			Int64("GameID", ai.GameID).
			Int("PlayerNum", ai.Num).
			Msgf("Planet %s queuing starbase %s build, will take %d years to build", planet.Name, design.Name, yearsToBuild)
		planet.ProductionQueue = append([]cs.ProductionQueueItem{item}, planet.ProductionQueue...)
		return true, nil
	}

	return false, nil
}

// add a new ship build request
func (ai *aiPlayer) addFleetBuildRequest(purpose cs.FleetPurpose, count int) {
	current := ai.requests.fleetBuilds[purpose]
	ai.requests.fleetBuilds[purpose] = current + count
}
