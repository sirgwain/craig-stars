package cs

import (
	"log/slog"
	"math"
	"slices"
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
	numBuiltInAYear := yearlyAvailableToSpend.ToCostFloat64().DivideCost(
		cost.Subtract(item.Allocated).SubtractMineral(mineralsOnHand).MinZero().ToCostFloat64())
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

	// the player spec is computed at turn generation time, but when updating planets we won't have it
	// so compute it on demand
	if player.Spec.Terraform == nil {
		player.Spec = ComputePlayerSpec(player, rules)
	}

	// the planet is a copy, but its queue is shared with the caller, so clone it before we
	// index it and produce from it
	planet.ProductionQueue = slices.Clone(planet.ProductionQueue)

	// in the browser the starbase only has its design, so compute its fleet spec for a copy of it.
	// Otherwise the planet spec we recompute each year would lose the starbase's dock and mass driver.
	if planet.Starbase != nil {
		starbase := *planet.Starbase
		starbase.Spec = ComputeFleetSpec(rules, player, &starbase)
		planet.Starbase = &starbase
	}

	// reset any estimates
	// we simulate until every item's estimate is settled
	settled := make([]bool, len(items))
	for i := range items {
		planet.ProductionQueue[i].index = i
		item := &items[i]
		item.QueueItemCompletionEstimate = QueueItemCompletionEstimate{
			YearsToBuildOne: Infinite,
			YearsToBuildAll: Infinite,
			YearsToSkipAuto: Infinite,
		}

		// auto alchemy in front of another item doesn't build anything itself
		settled[i] = item.Type == QueueItemTypeAutoMineralAlchemy && i < len(items)-1
	}

	producer := newProducer(slog.Default(), rules, &planet, player)
	producer.estimating = true
	for year := 1; year <= 100; year++ {
		// mine for minerals
		planet.mine(rules, planet.Spec.MiningOutput, min(planet.Mines, planet.Spec.MaxPossibleMines))
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
				// leftover resources put toward alchemy, not part of any item
				continue
			}
			item := &items[itemBuilt.index]

			if itemBuilt.queueItemType != item.Type {
				// a partial an auto item started, credit the auto item when it's finished
				if itemBuilt.numBuilt > 0 && item.YearsToBuildOne == Infinite {
					item.YearsToBuildOne = year
				}
				continue
			}

			switch {
			case itemBuilt.never:
				// this item will never complete
				settled[itemBuilt.index] = true
				continue
			case itemBuilt.status == productionStatusSkippedAuto:
				// an auto item at capacity is done
				if year == 1 {
					item.Skipped = true
				}
				if item.YearsToSkipAuto == Infinite {
					item.YearsToSkipAuto = year
				}
				settled[itemBuilt.index] = true
				continue
			case itemBuilt.skipped:
				// an auto item waiting on minerals may build later
				if item.YearsToSkipAuto == Infinite {
					item.YearsToSkipAuto = year
				}
				continue
			}

			// see if we already recorded when the first item was built
			if item.YearsToBuildOne == Infinite {
				item.YearsToBuildOne = year
			}

			switch {
			case item.Type == QueueItemTypeAutoMineralAlchemy:
				// trailing auto alchemy is continuous, its stored quantity is not a quota
				settled[itemBuilt.index] = true
			case itemBuilt.status == productionStatusComplete || itemBuilt.status == productionStatusCompleteAuto:
				// we built all of a concrete item, or all an auto item can build this year
				if item.YearsToBuildAll == Infinite {
					item.YearsToBuildAll = year
				}
				settled[itemBuilt.index] = true
			}
		}

		if !slices.Contains(settled, false) {
			// every item's estimate is settled, no need to loop anymore
			break
		}

		// grow pop
		planet.grow(player)
		planet.Spec = ComputePlanetSpec(rules, player, &planet)

		// colonists died off, no more production
		if planet.GetPopulation() < 0 {
			break
		}
	}

	return items, leftoverResourcesForResearch, nil
}
