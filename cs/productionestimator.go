package cs

import (
	"fmt"
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

type completionEstimate struct{}

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

// Simulate up to 100 years of growth (including mining & production)
// on a planet to determine how long each production queue item will take to build.
//
// After each year of growth, it checks what was built and records the year of the first and
// last completion.
// Items not finished within 100 turns are labeled as "never completable".
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
		// mine for minerals &
		planet.mine(rules)
		// TODO: Simulate remote mining for AR (likely by including the mineral outputs of remote miners directly in the planet's spec)

		// build stuff
		result, err := producer.produce()
		if err != nil {
			return nil, 0, fmt.Errorf("could not simulate production queue status for %d years into the future; produce() returned error %w", year, err)
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

		// if all colonists died off, we can stop producing
		// should never happen as pop can never go below 100 from "natural" causes
		if planet.GetPopulation() <= 0 {
			log.Logger.Debug().
				Int("Num", planet.Num).
				Str("Name", planet.Name).
				Int("PlayerNum", player.Num).
				Str("PlayerName", player.Race.PluralName).
				Int("years ahead", year).
				Msgf("Planet ran out of population; breaking loop")
			break
		}
	}

	return items, leftoverResourcesForResearch, nil
}
