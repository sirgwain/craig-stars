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

			existingQueueItemIndex := slices.IndexFunc(planet.ProductionQueue, func(item cs.ProductionQueueItem) bool { return item.DesignNum == design.Num })
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

// add a new ship build request
func (ai *aiPlayer) addShipBuildRequest(purpose cs.ShipDesignPurpose, count int) {
	current := ai.requests.shipBuilds[purpose]
	ai.requests.shipBuilds[purpose] = current + count
}
