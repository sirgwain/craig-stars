package cs

import "math"

// interface for populating completion estimates in a planet's production queue
type CompletionEstimator interface {
	// Get an estimate to complete a single item based on the minerals on hand on the planet (includes bonus resources from scrapping)
	// and the amount available per year from mining and production
	GetCompletionEstimate(item ProductionQueueItem, mineralsOnHand Mineral, yearlyAvailableToSpend Cost) QueueItemCompletionEstimate

	// get the estimated years to build one item with minerals on hand and some yearly mineral/resource output
	GetYearsToBuildOne(item ProductionQueueItem, mineralsOnHand Mineral, yearlyAvailableToSpend Cost) int

	// get a ProductionQueue with estimates filled in
	GetProductionWithEstimates(rules *Rules, player *Player, planet Planet) []ProductionQueueItem
}

type completionEstimate struct {
}

func newCompletionEstimator() CompletionEstimator {
	return &completionEstimate{}
}

// get the estimated years to build one item
func (e *completionEstimate) GetYearsToBuildOne(item ProductionQueueItem, mineralsOnHand Mineral, yearlyAvailableToSpend Cost) int {
	numYearsToBuildOne := yearlyAvailableToSpend.Divide(item.CostOfOne.Minus(item.Allocated).MinusMineral(mineralsOnHand).MinZero())
	if numYearsToBuildOne == 0 || math.IsInf(numYearsToBuildOne, 1) {
		return Infinite
	}
	return int(math.Ceil(1 / numYearsToBuildOne))
}

// get a completion estimate for a single item in the production queue
func (e *completionEstimate) GetCompletionEstimate(item ProductionQueueItem, mineralsOnHand Mineral, yearlyAvailableToSpend Cost) QueueItemCompletionEstimate {

	// Get the total cost of this item minus how much we've already allocated
	costOfOne := item.CostOfOne
	costOfAll := costOfOne.MultiplyInt(item.Quantity)

	// update the percent complete based on how much we've allocated vs the total cost of all items
	var percentComplete float64
	if !(item.Allocated == Cost{}) {
		percentComplete = clampFloat64(item.Allocated.Divide(costOfAll), 0, 1)
	} else {
		percentComplete = 0
	}

	var yearsToBuildAll, yearsToBuildOne int

	numYearsToBuildOne := yearlyAvailableToSpend.Divide(costOfOne.Minus(item.Allocated).MinusMineral(mineralsOnHand).MinZero())
	if numYearsToBuildOne == 0 || math.IsInf(numYearsToBuildOne, 1) {
		yearsToBuildOne = Infinite
	} else {
		yearsToBuildOne = int(math.Ceil(1 / numYearsToBuildOne))
	}

	numBuiltPerYear := yearlyAvailableToSpend.Divide(costOfAll.Minus(item.Allocated).MinusMineral(mineralsOnHand).MinZero()) * float64(item.Quantity)
	if numBuiltPerYear == 0 || math.IsInf(numBuiltPerYear, 1) {
		yearsToBuildAll = Infinite
	} else {
		yearsToBuildAll = int(math.Ceil(float64(item.Quantity) / numBuiltPerYear))
	}

	return QueueItemCompletionEstimate{
		YearsToBuildOne: yearsToBuildOne,
		YearsToBuildAll: yearsToBuildAll,
		PercentComplete: percentComplete,
	}
}

func (e *completionEstimate) GetProductionWithEstimates(rules *Rules, player *Player, planet Planet) []ProductionQueueItem {

	// copy the queue so we can update it
	items := make([]ProductionQueueItem, len(planet.ProductionQueue))
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

		for _, itemBuilt := range result.itemsBuilt {
			if itemBuilt.skipped {
				items[itemBuilt.index].Skipped = true
				continue
			}
			numBuilt[itemBuilt.index] = numBuilt[itemBuilt.index] + itemBuilt.numBuilt

			// TODO: this index changes... dang
			// see if we already recorded when the first item was built
			first := items[itemBuilt.index].YearsToBuildOne
			if first == Infinite {
				// we built one, update the years to build one
				items[itemBuilt.index].YearsToBuildOne = year
			}

			// we built the last
			last := items[itemBuilt.index].YearsToBuildAll
			if last == Infinite && items[itemBuilt.index].Quantity == numBuilt[itemBuilt.index] {
				items[itemBuilt.index].YearsToBuildAll = year
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

	return items
}
