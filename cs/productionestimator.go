package cs

import "math"

// The CompletionEstimator is used for populating completion estimates in a planet's production queue
type CompletionEstimator interface {
	// get the estimated years to build one item with minerals on hand and some yearly mineral/resource output
	GetYearsToBuildOne(item ProductionQueueItem, cost Cost, mineralsOnHand Mineral, yearlyAvailableToSpend Cost) int

	// get a ProductionQueue with estimates filled in and leftover resources
	GetProductionWithEstimates(rules *Rules, player *Player, planet Planet) ([]ProductionQueueItem, int)
}

type completionEstimate struct {
}

func NewCompletionEstimator() CompletionEstimator {
	return &completionEstimate{}
}

// get the estimated years to build one item
func (e *completionEstimate) GetYearsToBuildOne(item ProductionQueueItem, cost Cost, mineralsOnHand Mineral, yearlyAvailableToSpend Cost) int {
	numBuiltInAYear := yearlyAvailableToSpend.Divide(cost.Minus(item.Allocated).MinusMineral(mineralsOnHand).MinZero())
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
func (e *completionEstimate) GetProductionWithEstimates(rules *Rules, player *Player, planet Planet) (items []ProductionQueueItem, leftoverResourcesForResearch int) {

	// copy the queue so we can update it
	items = make([]ProductionQueueItem, len(planet.ProductionQueue))
	copy(items, planet.ProductionQueue)

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

	planet.Spec = computePlanetSpec(rules, player, &planet)

	// keep track of items built so we know how many auto items are completed
	numBuilt := make([]int, len(planet.ProductionQueue))
	producer := newProducer(&planet, player)
	for year := 1; year <= 100; year++ {
		// mine for minerals
		planet.mine(rules)
		// remote mine for AR
		//remoteMine()

		// build!
		result := producer.produce()

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
					if itemBuilt.numBuilt >= item.Quantity || (itemBuilt.numBuilt >= maxBuildable) {
						item.YearsToBuildAll = year
					}
				} else {
					if numBuiltSoFar >= item.Quantity || (itemBuilt.numBuilt >= maxBuildable) {
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

	return items, leftoverResourcesForResearch
}
