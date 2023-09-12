package ai

import (
	"fmt"

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

		for _, planet := range ai.Planets {
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

				if !ai.isShipInQueue(planet, ship.purpose, ship.quantity) {
					// put this ship at the top of the queue
					planet.ProductionQueue = append([]cs.ProductionQueueItem{{Type: cs.QueueItemTypeShipToken, Quantity: ship.quantity, DesignNum: design.Num}}, planet.ProductionQueue...)
					if err := ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders); err != nil {
						return err
					}
				}
			}
		}
	}

	for _, planet := range ai.Planets {
		yearlyAvailableToSpend := cs.FromMineralAndResources(planet.Spec.MiningOutput, planet.Spec.ResourcesPerYearAvailable)

		// try and build starbase if we don't have one built or queued
		// TODO: upgrade starbases
		if !(planet.Spec.HasStarbase || ai.isStarbaseInQueue(planet)) {
			// try and build a full starbase
			queued, err := ai.checkPlanetStarbaseBuild(planet, cs.ShipDesignPurposeStarbase, costCalculator, yearlyAvailableToSpend, completionEstimator)
			if err != nil {
				return err
			}
			if !queued {
				// try a fuel depot
				_, err := ai.checkPlanetStarbaseBuild(planet, cs.ShipDesignPurposeFuelDepot, costCalculator, yearlyAvailableToSpend, completionEstimator)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (ai *aiPlayer) checkPlanetStarbaseBuild(planet *cs.Planet, purpose cs.ShipDesignPurpose, costCalculator cs.CostCalculator, yearlyAvailableToSpend cs.Cost, completionEstimator cs.CompletionEstimator) (bool, error) {
	// design and upgrade this ship
	design, err := ai.designShip(ai.config.namesByPurpose[purpose], purpose, cs.FleetPurposeFromShipDesignPurpose(purpose))
	if err != nil {
		return false, fmt.Errorf("unable to design ship %v %w", purpose, err)
	}

	item := cs.ProductionQueueItem{Type: cs.QueueItemTypeStarbase, Quantity: 1, DesignNum: design.Num}
	item.SetDesign(design)
	cost, err := costCalculator.CostOfOne(ai.Player, item)
	if err != nil {
		return false, fmt.Errorf("calculate starbase cost %w", err)
	}
	item.CostOfOne = cost

	// if we can complete this soon, queue it
	yearsToBuild := completionEstimator.GetYearsToBuildOne(item, planet.Spec.MiningOutput, yearlyAvailableToSpend)
	if yearsToBuild != cs.Infinite && yearsToBuild <= ai.config.minYearsToQueueStarbase {
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
