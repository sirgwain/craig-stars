package cs

import "math"

// interface for populating completion estimates in a planet's production queue
type CompletionEstimator interface {
	// Get an estimate to complete a single item based on the minerals on hand on the planet (includes bonus resources from scrapping)
	// and the amount available per year from mining and production
	GetCompletionEstimate(item ProductionQueueItem, mineralsOnHand Cost, yearlyAvailableToSpend Cost) QueueItemCompletionEstimate

	// get a ProductionQueue with estimates filled in
	GetProductionWithEstimates(items []ProductionQueueItem, mineralsOnHand Cost, yearlyAvailableToSpend Cost) []ProductionQueueItem
}

type completionEstimate struct {
}

func newCompletionEstimator() CompletionEstimator {
	return &completionEstimate{}
}

// get a completion estimate for a single item in the production queue
func (e *completionEstimate) GetCompletionEstimate(item ProductionQueueItem, mineralsOnHand Cost, yearlyAvailableToSpend Cost) QueueItemCompletionEstimate {

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

	costOfOneMinusAllocated := costOfOne.Minus(item.Allocated).Minus(mineralsOnHand).MinZero()
	_ = costOfOneMinusAllocated
	numYearsToBuildOne := yearlyAvailableToSpend.Divide(costOfOne.Minus(item.Allocated).Minus(mineralsOnHand).MinZero())
	if numYearsToBuildOne == 0 || math.IsInf(numYearsToBuildOne, 1) {
		yearsToBuildOne = math.MaxInt
	} else {
		yearsToBuildOne = int(math.Ceil(1 / numYearsToBuildOne))
	}

	numBuiltPerYear := yearlyAvailableToSpend.Divide(costOfAll.Minus(item.Allocated).Minus(mineralsOnHand).MinZero()) * float64(item.Quantity)
	if numBuiltPerYear == 0 || math.IsInf(numBuiltPerYear, 1) {
		yearsToBuildAll = math.MaxInt
	} else {
		yearsToBuildAll = int(math.Ceil(float64(item.Quantity) / numBuiltPerYear))
	}

	return QueueItemCompletionEstimate{
		YearsToBuildOne: yearsToBuildOne,
		YearsToBuildAll: yearsToBuildAll,
		PercentComplete: percentComplete,
	}
}

// For a set of ProductionQueueItems, calculate the YearsToBuild and YearsToBuildOne properties
func (e *completionEstimate) GetProductionWithEstimates(items []ProductionQueueItem, mineralsOnHand Cost, yearlyAvailableToSpend Cost) []ProductionQueueItem {

	// go through each item and update it's YearsToComplete field
	updatedItems := make([]ProductionQueueItem, len(items))
	for i := range items {
		item := items[i]

		// skip items we can't build
		if item.Type.IsAuto() && item.MaxBuildable == 0 {
			item.QueueItemCompletionEstimate = QueueItemCompletionEstimate{
				Skipped: true,
			}
			continue
		}

		// calculate the estimate for this item
		estimate := e.GetCompletionEstimate(item, mineralsOnHand, yearlyAvailableToSpend)
		item.QueueItemCompletionEstimate = estimate

		// reduce the minerals on hand for the next item the total cost of building this item
		// this may mean we have negative resources, which means it'll take more years
		// to build the next items
		costOfAll := item.CostOfOne.MultiplyInt(item.Quantity).Minus(item.Allocated)		
		if item.Type.IsAuto() {
			// for auto tasks, we skip if missing minerals, so we only account for 
			// resources not spent
			mineralsOnHand.Resources -= costOfAll.Resources
		} else {
			mineralsOnHand = mineralsOnHand.Minus(costOfAll)
		}
		updatedItems[i] = item
	}

	return updatedItems
}
