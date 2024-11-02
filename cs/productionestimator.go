package cs

import (
	"math"

	"github.com/rs/zerolog/log"
)

// The CompletionEstimator is used for populating completion estimates in a planet's production queue
type CompletionEstimator interface {
	// get the estimated years to build one item with minerals on hand and some yearly mineral/resource output
	GetYearsToBuildOne(item ProductionQueueItem, cost Cost, mineralsOnHand Mineral, yearlyAvailableToSpend Cost) int

	// get a ProductionQueue with estimates filled in
	GetProductionWithEstimates(rules *Rules, player *Player, planet Planet) ([]ProductionQueueItem, int, error)
}

type completionEstimate struct {
}

func NewCompletionEstimator() CompletionEstimator {
	return &completionEstimate{}
}

// get the estimated years to build one item
func (e *completionEstimate) GetYearsToBuildOne(item ProductionQueueItem, cost Cost, mineralsOnHand Mineral, yearlyAvailableToSpend Cost) int {
	numBuiltInAYear := yearlyAvailableToSpend.Divide(cost.Subtract(item.Allocated).SubtractMineral(mineralsOnHand).MinZero())
	if numBuiltInAYear == 0 || math.IsInf(numBuiltInAYear, 1) {
		return Infinite
	}
	return int(math.Ceil(1 / numBuiltInAYear))
}

// simulate up to 100 years of production to determine the time each item will take to build
// this function will take a copy of the planet and do the following:
// * clone the production queue
// * add an index to each production queue item so we can track it in the produce() result
// * default each item to never being completed
// * simulate 100 years of growth
//   - mine for resources
//   - run production (including terraforming the planet, building mines and factories, etc)
//   - grow pop on the planet
//
// For each year of growth, it checks what was built. If an item was built for the first time
// it records the year. If the item completed building, it records the last year
// when all items are complete or 100 years have passed, iit returns
func (e *completionEstimate) GetProductionWithEstimates(rules *Rules, player *Player, planet Planet) (items []ProductionQueueItem, leftoverResourcesForResearch int, err error) {

	// copy the queue so we can update it
	items = make([]ProductionQueueItem, len(planet.ProductionQueue))
	copy(items, planet.ProductionQueue)

	if len(items) == 0 {
		return items, planet.Spec.ResourcesPerYear, nil
	}

	// reset any estimates
	for i := range items {
		planet.ProductionQueue[i].index = i
		item := &items[i]
		item.QueueItemCompletionEstimate = QueueItemCompletionEstimate{
			YearsToBuildOne: Infinite,
			YearsToBuildAll: Infinite,
			YearsToSkipAuto: Infinite,
		}
	}

	// keep track of items built so we know how many auto items are completed
	numBuilt := make([]int, len(planet.ProductionQueue))
	producer := newProducer(log.Logger, rules, &planet, player)
	for year := 1; year <= 100; year++ {
		// mine for minerals
		planet.mine(rules)
		// remote mine for AR
		//remoteMine()

		// build!
		result, err := producer.produce()
		if err != nil {
			return nil, 0, err
		}

		if year == 1 {
			leftoverResourcesForResearch = result.leftoverResources
		}

		for _, itemBuilt := range result.itemsBuilt {
			if itemBuilt.index == -1 {
				// skip partial auto builds
				continue
			}
			item := &items[itemBuilt.index]
			maxBuildable := planet.maxBuildable(player, item.Type)

			// this will be skipped if we've hit the max allowed
			if itemBuilt.skipped {
				if year == 1 && maxBuildable == 0 {
					item.Skipped = true
					item.YearsToSkipAuto = 1
				} else {
					if item.YearsToSkipAuto == Infinite {
						item.YearsToSkipAuto = year
					}
				}
				continue
			}

			// this item will never complete
			if itemBuilt.never {
				continue
			}
			numBuiltSoFar := numBuilt[itemBuilt.index] + itemBuilt.numBuilt
			numBuilt[itemBuilt.index] = numBuiltSoFar

			// see if we already recorded when the first item was built
			first := item.YearsToBuildOne
			if first == Infinite {
				// we built one, update the years to build one
				item.YearsToBuildOne = year
			}

			// check if we built the last one of this group
			// if we've built the item's original quantity, or we've built some and the maxBuildable remaining is 0
			// we're done
			last := item.YearsToBuildAll
			if last == Infinite {
				if item.Type.IsAuto() {
					if itemBuilt.numBuilt >= item.Quantity || (maxBuildable != Infinite && itemBuilt.numBuilt >= maxBuildable) {
						item.YearsToBuildAll = year
					}
				} else {
					if numBuiltSoFar >= item.Quantity || (maxBuildable != Infinite && itemBuilt.numBuilt >= maxBuildable) {
						item.YearsToBuildAll = year
					}
				}
			}
		}

		if result.completed {
			// we built everything in the queue, no need to loop anymore
			break
		}

		// grow pop
		planet.grow(player)
		planet.Spec = computePlanetSpec(rules, player, &planet)

		// colonists died off, no more production
		if planet.population() < 0 {
			break
		}
	}

	return items, leftoverResourcesForResearch, nil
}
