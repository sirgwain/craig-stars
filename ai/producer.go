package ai

import (
	"fmt"

	"github.com/sirgwain/craig-stars/cs"
	"golang.org/x/exp/slices"
)

type fleetPurpose uint

const (
	scout fleetPurpose = iota
	colonizer
	transport
	miner
	minelayer
	fighter
	bomber
)

type fleet struct {
	purpose fleetPurpose
	ships   []fleetShip
}

type fleetShip struct {
	purpose  cs.ShipDesignPurpose
	hull     *cs.TechHull
	quantity int
}

func (ai *aiPlayer) produce() error {
	buildablePlanets := ai.getPlanetsWithDocks()
	costCalculator := cs.NewCostCalculator()
	completionEstimator := cs.NewCompletionEstimator()

	for purpose, quantity := range ai.requests.shipBuilds {
		if quantity <= 0 {
			continue
		}

		// design and upgrade this ship
		design, err := ai.designShip(ai.config.namesByPurpose[purpose], purpose)
		if err != nil {
			return fmt.Errorf("unable to design ship %v %w", purpose, err)
		}

		for _, planet := range buildablePlanets {
			if !planet.CanBuild(design.Spec.Mass) {
				continue
			}

			// don't build colonizers or freighters on this planet if it doesn't have the pop to remove them
			if (purpose == cs.ShipDesignPurposeColonistFreighter || purpose == cs.ShipDesignPurposeColonizer) &&
				planet.Spec.PopulationDensity < ai.config.colonizerPopulationDensity {
				continue
			}

			existingQueueItemIndex := slices.IndexFunc(planet.ProductionQueue, func(item cs.ProductionQueueItem) bool {
				if item.Type == cs.QueueItemTypeShipToken {
					itemDesign := ai.GetDesign(item.DesignNum)
					if itemDesign != nil {
						// true if this design or one like it is already queued
						return item.DesignNum == design.Num || itemDesign.Purpose == purpose
					}
				}
				return false
			})
			if existingQueueItemIndex == -1 {
				// put a new scout at the front of the queue
				planet.ProductionQueue = append([]cs.ProductionQueueItem{{Type: cs.QueueItemTypeShipToken, Quantity: 1, DesignNum: design.Num}}, planet.ProductionQueue...)
				if err := ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders); err != nil {
					return err
				}
			}
		}
	}

	for _, planet := range ai.Planets {
		yearlyAvailableToSpend := cs.FromMineralAndResources(planet.Spec.MiningOutput, planet.Spec.ResourcesPerYearAvailable)

		// try and build starbases
		// TODO: upgrade starbases
		if !planet.Spec.HasStarbase && !ai.isStarbaseInQueue(planet) {
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

		// check for terraforming
		if planet.Spec.CanTerraform {
			item := cs.ProductionQueueItem{Type: cs.QueueItemTypeTerraformEnvironment, Quantity: 1}
			cost, err := costCalculator.CostOfOne(ai.Player, item)
			if err != nil {
				return fmt.Errorf("calculate terraform cost %w", err)
			}

			// if it takes only a year to build this, add it to the production queue
			item.CostOfOne = cost
			if completionEstimator.GetYearsToBuildOne(item, planet.Cargo.ToMineral(), yearlyAvailableToSpend) <= 1 {
				planet.ProductionQueue = append([]cs.ProductionQueueItem{item}, planet.ProductionQueue...)
			}
		}
	}
	return nil
}

func (ai *aiPlayer) checkPlanetStarbaseBuild(planet *cs.Planet, purpose cs.ShipDesignPurpose, costCalculator cs.CostCalculator, yearlyAvailableToSpend cs.Cost, completionEstimator cs.CompletionEstimator) (bool, error) {
	// design and upgrade this ship
	design, err := ai.designShip(ai.config.namesByPurpose[purpose], purpose)
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
func (ai *aiPlayer) addShipBuildRequest(purpose cs.ShipDesignPurpose, count int) {
	current := ai.requests.shipBuilds[purpose]
	ai.requests.shipBuilds[purpose] = current + count
}
