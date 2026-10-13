package ai

import (
	"math"

	"github.com/sirgwain/craig-stars/cs"
)

// How much the AI commits to at once grows with what its empire can support
// and with the size of the universe. Difficulty sets the pace.

// planets with a starbase that can build ships
func (ai *aiPlayer) shipyards() int {
	count := 0
	for _, planet := range ai.Planets {
		if planet.Spec.HasStarbase && planet.Spec.DockCapacity != 0 {
			count++
		}
	}
	return count
}

// spareShipyards counts developed shipyards whose ship and starbase orders
// will take less than a year of their resources: worlds with capacity going
// unused.
func (ai *aiPlayer) spareShipyards() int {
	count := 0
	for _, planet := range ai.Planets {
		if !planet.Spec.HasStarbase || planet.Spec.DockCapacity == 0 || ai.infrastructureDevelopment(planet) < ai.config.fleetProductionCutoff {
			continue
		}
		ordered := 0
		for _, item := range planet.ProductionQueue {
			if design := ai.GetDesign(item.DesignNum); design != nil && (item.Type == cs.QueueItemTypeShipToken || item.Type == cs.QueueItemTypeStarbase) {
				ordered += item.Quantity*design.Spec.Cost.Resources - item.Allocated.Resources
			}
		}
		if ordered < planet.Spec.ResourcesPerYearAvailable {
			count++
		}
	}
	return count
}

// scale a count by the difficulty's pace, never below one
func (ai *aiPlayer) paced(count float64) int {
	return max(1, int(math.Round(ai.config.pace*count)))
}

// Colony ships in flight or waiting at once. A shipyard with capacity going
// unused is better spent on more colony ships than on research.
func (ai *aiPlayer) colonyLimit() int {
	limit := ai.config.colonyPipeline
	if ai.doctrine.colonyShips > 0 {
		limit = int(math.Round(float64(limit) * ai.doctrine.colonyShips))
	}
	if ai.enabled("spare-colony") {
		if spare := ai.spareShipyards(); spare > 0 {
			extra := float64(spareColonyShips * spare)
			if ai.doctrine.colonyShips > 0 {
				extra *= ai.doctrine.colonyShips
			}
			limit += ai.paced(extra)
		}
	}
	return limit
}

// scouts owned at once: enough to cover what is left to explore, however large the universe
func (ai *aiPlayer) scoutLimit() int {
	unexplored := 0
	for _, intel := range ai.PlanetIntels {
		if intel.Unexplored() {
			unexplored++
		}
	}
	return max(2, ai.paced(float64(unexplored)/unexploredPerScout))
}

// miner or mine layer fleets requested at once
func (ai *aiPlayer) supportLimit() int {
	return ai.paced(float64(ai.shipyards()))
}

// War fleets to have under construction at once: one for each shipyard
// developed enough to build them, and none while the ships already on order
// would take more than a few years of the empire's war budget.
func (ai *aiPlayer) warLimit() int {
	yards, income, ordered := 0, 0, 0
	for _, planet := range ai.Planets {
		income += planet.Spec.ResourcesPerYear
		if ai.isPlanetReadyToBuildFleet(planet, cs.FleetPurposeBomber) {
			yards++
		}
		for _, item := range planet.ProductionQueue {
			if item.Type == cs.QueueItemTypeShipToken && item.GetTag(cs.TagPurpose) == string(cs.FleetPurposeBomber) {
				if design := ai.GetDesign(item.DesignNum); design != nil {
					ordered += item.Quantity*design.Spec.Cost.Resources - item.Allocated.Resources
				}
			}
		}
	}
	if float64(ordered) > warBudgetYears*warBudgetShare*ai.config.pace*float64(income) {
		return 0
	}
	fleets := float64(yards)
	if ai.doctrine.warFleets > 0 {
		fleets *= ai.doctrine.warFleets
	}
	return ai.paced(fleets)
}

const (
	unexploredPerScout = 10  // unexplored planets for each scout
	spareColonyShips   = 8   // extra colony ships at once for each shipyard with unused capacity
	warBudgetShare     = .3  // share of empire income for war fleets at normal pace
	warBudgetYears     = 3.0 // years of that budget that may be on order at once
)
