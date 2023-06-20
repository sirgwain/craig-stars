package cs

import "math"

// interface for populating completion estimates in a planet's production queue
type CompletionEstimator interface {
	PopulateCompletionEstimates(items []ProductionQueueItem, yearlyAvailableToSpend Cost)
}

type completionEstimate struct {
}

func newCompletionEstimator() CompletionEstimator {
	return &completionEstimate{}
}

/// For a set of ProductionQueueItems, calculate the YearsToBuild and YearsToBuildOne properties
func (e *completionEstimate) PopulateCompletionEstimates(items []ProductionQueueItem, yearlyAvailableToSpend Cost) {

	// as we go, subtract the previous items
	previousItemsCost := yearlyAvailableToSpend.Negate()

	// go through each item and update it's YearsToComplete field
	for i := range items {
		item := &items[i]

		// skip items we can't build
		if item.Type.IsAuto() && item.MaxBuildable == 0 {
			item.QueueItemCompletionEstimate = QueueItemCompletionEstimate{
				Skipped: true,
			}
			continue
		}
		previousItemsCost = previousItemsCost.Minus(item.Allocated)

		var costOfOne = item.CostOfOne

		estimate := e.getCompletionEstimate(*item, yearlyAvailableToSpend, previousItemsCost)
		item.QueueItemCompletionEstimate = estimate

		// increase our previousItemsCost for the next item
		// reduce the available resources for next estimate
		previousItemsCost = previousItemsCost.Add(costOfOne.MultiplyInt(item.Quantity))
	}
}

// get a completion estimate for a single item in the production queue
func (e *completionEstimate) getCompletionEstimate(item ProductionQueueItem, yearlyAvailableToSpend Cost, previousItemsCost Cost) QueueItemCompletionEstimate {

	// Get the total cost of this item plus any previous items in the queue
	// and subtract what we have on hand (that will be applied this year)
	costOfOne := item.CostOfOne
	costOfAll := costOfOne.MultiplyInt(item.Quantity).Add(previousItemsCost)
	totalCostToBuildOne := costOfOne.Add(previousItemsCost)
	totalCostToBuildAll := costOfAll.Add(previousItemsCost)

	// If we have a bunch of leftover minerals because our planet is full, 0 those out
	totalCostToBuildOne = totalCostToBuildOne.MinZero()
	totalCostToBuildAll = totalCostToBuildAll.MinZero()

	var percentComplete float64
	if !(item.Allocated == Cost{}) {
		percentComplete = clampFloat64(item.Allocated.Divide(costOfAll), 0, 1)
	} else {
		percentComplete = 0
	}

	yearsToBuildAll := 1
	yearsToBuildOne := 1
	if !(totalCostToBuildAll == Cost{}) {
		numBuiltPerYear := yearlyAvailableToSpend.Divide(totalCostToBuildAll)
		if numBuiltPerYear == 0 || math.IsInf(numBuiltPerYear, 1) {
			yearsToBuildAll = math.MaxInt
		} else {
			yearsToBuildAll = int(math.Ceil(1.0 / numBuiltPerYear))
		}
	}
	if !(totalCostToBuildOne == Cost{}) {
		numBuiltPerYear := yearlyAvailableToSpend.Divide(totalCostToBuildOne)
		if numBuiltPerYear == 0 || math.IsInf(numBuiltPerYear, 1) {
			yearsToBuildOne = math.MaxInt
		} else {
			yearsToBuildOne = int(math.Ceil(1 / numBuiltPerYear))
		}
	}

	return QueueItemCompletionEstimate{
		YearsToBuildOne: yearsToBuildOne,
		YearsToBuildAll: yearsToBuildAll,
		PercentComplete: percentComplete,
	}
}
