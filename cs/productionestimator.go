package cs

import (
	"fmt"
	"math"

	"github.com/rs/zerolog/log"
)

// The CompletionEstimator is used for populating completion estimates in a planet's production queue
type CompletionEstimator interface {
	// GetYearsToBuild returns the number of years required to build a given ProductionQueueItem,
	// given some amount of minerals on hand and yearly mineral/resource output.
	GetYearsToBuild(item ProductionQueueItem, cost Cost, mineralsOnHand Mineral, yearlyAvailableToSpend Cost) int

	// GetProductionWithEstimates populates a planet's production queue with estimates
	// for how long each item will take to build.
	// It simulates up to 100 years of growth, mining & production,
	// recording the first and last time each item was completed.
	//
	// It returns the new updated queue, the amount of resources left for research, and any error encountered.
	// An unsuccessful run will return nil, 0, err.
	GetProductionWithEstimates(rules *Rules, player *Player, planet Planet) ([]ProductionQueueItem, int, error)
}

type completionEstimate struct{}

func NewCompletionEstimator() CompletionEstimator {
	return &completionEstimate{}
}

// GetYearsToBuild returns the number of years required to build a given ProductionQueueItem,
// given some amount of minerals on hand and yearly mineral/resource output.
func (e *completionEstimate) GetYearsToBuild(item ProductionQueueItem, costPerItem Cost, mineralsOnHand Mineral, yearlyAvailableToSpend Cost) int {
	costPerItem = costPerItem.Subtract(item.Allocated).SubtractMineral(mineralsOnHand).MinZero()
	numBuiltPerYear := yearlyAvailableToSpend.DivideCost(costPerItem)
	if numBuiltPerYear <= 0 {
		return Infinite
	} else if math.IsInf(numBuiltPerYear, 0) {
		return 0
	}
	return int(math.Ceil(float64(item.Quantity) / numBuiltPerYear))
}

// GetProductionWithEstimates populates a planet's production queue with estimates
// for how long each item will take to build.
// It simulates up to 100 years of growth, mining & production,
// recording the first and last time each item was completed.
//
// It returns the new updated queue, the amount of resources left for research, and any error encountered.
// An unsuccessful run will return nil, 0, err.
func (e *completionEstimate) GetProductionWithEstimates(rules *Rules, player *Player, planet Planet) (items []ProductionQueueItem, leftoverResourcesForResearch int, err error) {
	if len(planet.ProductionQueue) == 0 {
		// no queue makes our job quite easy
		return planet.ProductionQueue, planet.Spec.ResourcesPerYear, nil
	}

	// reset any prior estimates and add indices to queue items
	for i := range planet.ProductionQueue {
		planet.ProductionQueue[i].index = i
		planet.ProductionQueue[i].QueueItemCompletionEstimate = QueueItemCompletionEstimate{
			YearsToBuildOne:     Infinite,
			YearsToBuildAll:     Infinite,
			YearsToSkipOrCancel: Infinite,
		}
	}

	// copy the queue so we can update it
	items = make([]ProductionQueueItem, len(planet.ProductionQueue))
	copy(items, planet.ProductionQueue)

	numBuilt := make([]int, len(planet.ProductionQueue)) // slice tracking items built for each production queue item
	producer := newProducer(log.Logger, rules, &planet, player)
	for year := 1; year <= 100; year++ {
		// mine for minerals
		planet.mine(rules, planet.Spec.MiningOutput, min(planet.Mines, planet.Spec.MaxPossibleMines))
		// TODO: Simulate AR remote mining (perhaps with a slice of mining rates passed down by the caller)

		// build!
		result, err := producer.produce()
		if err != nil {
			return nil, 0, fmt.Errorf("error while producing items: %w", err)
		}

		if year == 1 {
			leftoverResourcesForResearch = result.leftoverResources
		}

		// check everything built this turn, tacking on estimates if we haven't already done so
		for _, itemBuilt := range result.itemsBuilt {
			if itemBuilt.index == -1 {
				// item is a half-built concrete version of an auto item; skip
				continue
			}

			item := &items[itemBuilt.index]
			maxBuildable := planet.MaxBuildable(player, item.Type)

			// log auto items being skipped or concrete items being canceled
			if itemBuilt.skipped && item.YearsToSkipOrCancel == Infinite {
				item.YearsToSkipOrCancel = year
				continue
			}

			numBuilt[itemBuilt.index] += itemBuilt.numBuilt

			// record the year the first item was built (if not done beforehand)
			if item.YearsToBuildOne == Infinite {
				item.YearsToBuildOne = year
			}

			// check if we've built the last item in this group
			if item.YearsToBuildAll == Infinite {
				var num int
				if item.Type.IsAuto() {
					// auto items refresh each year, so we check how many
					// were built this current year
					num = itemBuilt.numBuilt
				} else {
					// concrete items retain quantities each year, so we check
					// the total items built across all years
					num = numBuilt[itemBuilt.index]
				}
				// if we've built up to the item's quantity or
				// maxBuildable, mark it as done
				if num >= item.Quantity || (maxBuildable != Infinite && num >= maxBuildable) {
					item.YearsToBuildAll = year
				}

			}
		}

		if result.completed {
			// we built (or skipped) the last item in the queue; no need to loop anymore
			break
		}

		// if we made a base, simulate adding it to the planet
		if result.starbase != nil {
			s := newStarbase(player, &planet, result.starbase, result.starbase.Name)
			s.Spec = ComputeFleetSpec(rules, player, s)
			planet.Starbase = s
		}

		// grow pop & compute spec
		planet.grow(player)
		planet.Spec = computePlanetSpec(rules, player, &planet)

		if planet.GetPopulation() <= 0 {
			// colonists died off completely, which should never happen under normal circumstances
			log.Warn().
				Int("Num", planet.Num).
				Str("Name", planet.Name).
				Int("PlayerNum", player.Num).
				Str("PlayerName", player.Race.PluralName).
				Int("Years Passed", year).
				Int("Population", planet.exactPopulation()).
				Msg("Planet population went negative during production estimates")

			break
		}
	}

	return items, leftoverResourcesForResearch, nil
}
