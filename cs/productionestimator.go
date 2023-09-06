package cs

import "math"

// interface for populating completion estimates in a planet's production queue
type CompletionEstimator interface {
	// get the estimated years to build one item with minerals on hand and some yearly mineral/resource output
	GetYearsToBuildOne(item ProductionQueueItem, mineralsOnHand Mineral, yearlyAvailableToSpend Cost) int

	// get a ProductionQueue with estimates filled in
	GetProductionWithEstimates(rules *Rules, player *Player, planet Planet) ([]ProductionQueueItem, int)
}

type completionEstimate struct {
}

func newCompletionEstimator() CompletionEstimator {
	return &completionEstimate{}
}

// get the estimated years to build one item
func (e *completionEstimate) GetYearsToBuildOne(item ProductionQueueItem, mineralsOnHand Mineral, yearlyAvailableToSpend Cost) int {
	numBuiltInAYear := yearlyAvailableToSpend.Divide(item.CostOfOne.Minus(item.Allocated).MinusMineral(mineralsOnHand).MinZero())
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
		}
		item.PercentComplete = item.percentComplete()
	}

	// keep track of items built so we know how many auto items are completed
	numBuilt := make([]int, len(planet.ProductionQueue))
	producer := newProducer(&planet, player)
	completedItems := 0
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
			item := &items[itemBuilt.index]
			maxBuildable := planet.maxBuildable(item.Type)

			// this will be skipped if we've hit the max allowed
			if itemBuilt.skipped {
				item.Skipped = true
				if item.YearsToBuildOne == Infinite {
					item.YearsToBuildOne = 0
				} else {
					item.YearsToBuildOne = year
				}
				if item.YearsToBuildAll == Infinite {
					item.YearsToBuildAll = 0
				} else {
					item.YearsToBuildAll = year
				}
				completedItems++
				continue
			}

			// this item will never complete. If it's auto, it'll be ignored
			if itemBuilt.never {
				continue
			}
			numBuiltSoFar := numBuilt[itemBuilt.index] + itemBuilt.numBuilt
			numBuilt[itemBuilt.index] = numBuiltSoFar

			// TODO: this index changes... dang
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
			if last == Infinite && (numBuiltSoFar >= item.Quantity || (numBuiltSoFar > 0 && maxBuildable == 0)) {
				item.YearsToBuildAll = year
				completedItems++
			}
		}

		if completedItems >= len(items) {
			// we have end build years for all items, no need to loop anymore
			break
		}

		// grow pop
		planet.grow(player)
		planet.Spec = computePlanetSpec(rules, player, &planet)
	}

	return items, leftoverResourcesForResearch
}
