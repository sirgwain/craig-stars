package cs

import (
	"fmt"
	"maps"
)

// The CostCalculator interface is used to calculate costs of single items or starbase upgrades
// This is used by planetary production and estimating production queue completion
type CostCalculator interface {
	StarbaseUpgradeCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design, newDesign *ShipDesign) Cost
	CostOfOne(player *Player, item ProductionQueueItem) (Cost, error)
	GetDesignCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design *ShipDesign) Cost
}

func NewCostCalculator() CostCalculator {
	return &costCalculate{}
}

type costCalculate struct {
}

// get the upgrade cost for replacing a starbase with another
//
// Takes into account part replacement costs
func (p *costCalculate) StarbaseUpgradeCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design, newDesign *ShipDesign) Cost {
	if design.SlotsEqual(newDesign.Slots) && design.Hull == newDesign.Hull {
		// Exact same base; no calcs needed
		return Cost{}
	}

	credit := Cost{}
	cost := Cost{}
	oldComponents := map[*TechHullComponent]int{} // Maps hull component to quantity
	newComponents := map[*TechHullComponent]int{}
	oldComponentsByCategory := map[TechCategory][]*TechHullComponent{}
	newComponentsByCategory := map[TechCategory][]*TechHullComponent{}

	// Firstly, if the hulls are different, add (newHullCost - 0.5*OldHullCost)
	// Multiplied by 10*componentCostReduction for rounding purposes
	if design.Hull != newDesign.Hull {
		oldHullCost := rules.techs.GetHull(design.Hull).Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(5 * rules.StarbaseComponentCostReduction)
		newHullCost := rules.techs.GetHull(newDesign.Hull).Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(10 * rules.StarbaseComponentCostReduction)
		cost = cost.Add(newHullCost).Minus(oldHullCost)
	}
	// iterate through both designs' slots and tally up items in each
	for i := 0; i < MaxInt(len(design.Slots), len(newDesign.Slots)); i++ {
		// don't wanna index arrays out of bounds!
		if i < len(design.Slots) {
			oldComponents[rules.techs.GetHullComponent(design.Slots[i].HullComponent)] += design.Slots[i].Quantity
		}
		if i < len(newDesign.Slots) {
			newComponents[rules.techs.GetHullComponent(newDesign.Slots[i].HullComponent)] += newDesign.Slots[i].Quantity
		}
	}

	// Iterate through all new parts in list to see if they are present on the old base
	// to create a list of all unique components
	if len(oldComponents) > 0 && len(newComponents) > 0 {
		for item, newQuantity := range newComponents {
			oldQuantity := oldComponents[item]
			if newQuantity == oldQuantity {
				// same amount of item in both bases; remove from both
				delete(oldComponents, item)
				delete(newComponents, item)
			} else if newQuantity > oldQuantity {
				// More copies of item in new design; add extras to new base list
				newComponentsByCategory[item.Tech.Category] = append(newComponentsByCategory[item.Tech.Category], item)
				newComponents[item] = (newQuantity - oldQuantity)
				delete(oldComponents, item)
			} else {
				// More copies of item in original design (or item doesn't exist on new base)
				// Add extras to old base list
				oldComponents[item] = (oldQuantity - newQuantity)
				delete(newComponents, item)
			}
		}
	}

	if len(oldComponents) == 0 {
		// no items in old base not also present in the new one
		// We can just tally up all our costs for the new stuff and be done for the day
		for item, qty := range newComponents {
			if item.Tech.Category == "Orbital" {
				cost = cost.Add(item.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(qty).MultiplyInt(10 * rules.StarbaseComponentCostReduction))
			} else {
				cost = cost.Add(item.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(qty).MultiplyInt(10))
			}
		}
		return cost.DivideByInt(10*rules.StarbaseComponentCostReduction, true)
	} else {
		// Loop through any remaining items from old base and add to category list
		for item := range oldComponents {
			oldComponentsByCategory[item.Tech.Category] = append(oldComponentsByCategory[item.Tech.Category], item)
		}
	}

	// At this point, we should have 4 maps in total: 2 for each base design
	// ComponentsUnique contains all components unique to each one mapped to their quantity
	// ComponentsByCategory contains a list of all categories present in each base,
	// mapped to a slice of all components on the base for said category
	// Now, all that's left is the cost calcs

	// Get categories present in either map type so we don't have to iterate over every single tachCategory

	categories := oldComponentsByCategory
	maps.Copy(categories, newComponentsByCategory)

	// Tally up costs per category
	// We multiply everything by 2x here (and by 10x in the following step) & divide by 20 later to minimize rounding errors
	for category := range categories {
		oldCost := Cost{}
		newCost := Cost{}

		for _, oldItem := range oldComponentsByCategory[category] {
			if category == "Orbital" {
				oldCost = oldCost.Add(oldItem.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(oldComponents[oldItem]).MultiplyInt(rules.StarbaseComponentCostReduction))
			} else {
				oldCost = oldCost.Add(oldItem.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(oldComponents[oldItem]))
			}
		}
		for _, newItem := range newComponentsByCategory[category] {
			if category == "Orbital" {
				newCost = newCost.Add(newItem.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(newComponents[newItem]).MultiplyInt(rules.StarbaseComponentCostReduction))
			} else {
				newCost = newCost.Add(newItem.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(newComponents[newItem]))
			}
		}

		if len(oldComponentsByCategory[category]) > 0 {
			if len(newComponentsByCategory[category]) > 0 {
				// Category present in both bases
				// Apply lower (70%) rebate to credit tally (up to 70% of the actual item value)
				// Apply difference between 2 discounts (10%) to this item category only, up to 10% of the original item value
				// Costs here are multiplied by factor of 10 to be divided later on

				// Compute costs for each resource type separately (I/B/G/R)
				for costType := 0; costType < 4; costType++ {
					// extract float values for items
					oldCostInt := oldCost.GetAmount(costType)
					newCostInt := newCost.GetAmount(costType)
					if oldCostInt == 0 && newCostInt == 0 {
						continue
					}

					sameCategoryRebate := oldCostInt * 8
					differentCategoryRebate := oldCostInt * 7

					if oldCostInt == 0 {
						cost = cost.AddInt(costType, newCostInt*10)
					} else if newCostInt == 0 {
						credit = credit.AddInt(costType, differentCategoryRebate)
					} else {
						// Add global rebate to credit tally
						credit = credit.AddInt(costType, differentCategoryRebate)

						// Consume global credit tally to reduce price from 100% to 30%
						// If this turns credit negative, no problem!
						// We add it to Cost at the end anyways
						adjCost := 3 * newCostInt
						credit = credit.AddInt(costType, -(10*newCostInt - adjCost))

						// Add on category specific rebate
						categoryRebate := sameCategoryRebate-differentCategoryRebate
						adjCost = MinInt(2*newCostInt, adjCost - categoryRebate)

						cost = cost.AddInt(costType, adjCost)
					}
				}
			} else {
				// category present in old but not new design; add to global credit tally
				credit = credit.Add(oldCost.MultiplyInt(7))
			}
		} else {
			if len(newComponentsByCategory[category]) > 0 {
				// item category present in new but not old design; add to global cost tally
				cost = cost.Add(newCost.MultiplyInt(10))
			} else {
				// item category not present in either design
				// Should never happen because we only iterate over
				// categories present in either map, but shouldn't matter either way
				continue
			}
		}
	}

	return cost.Minus(credit).DivideByInt(10*rules.StarbaseComponentCostReduction, true).MinZero()
}

// Get the cost of one item in a production queue, for a player
func (p *costCalculate) CostOfOne(player *Player, item ProductionQueueItem) (Cost, error) {
	cost := player.Race.Spec.Costs[item.Type]
	if item.Type == QueueItemTypeStarbase || item.Type == QueueItemTypeShipToken {
		if item.design != nil {
			cost = item.design.Spec.Cost
		} else {
			return Cost{}, fmt.Errorf("design %d not populated in queue item", item.DesignNum)
		}
	}
	return cost, nil
}

// Get cost of a given ship or new starbase design
func (p *costCalculate) GetDesignCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design *ShipDesign) Cost {
	cost := Cost{}
	costOrbital := rules.techs.GetHullComponent(design.Hull).Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(10)
	// iterate through slots and tally prices up
	for _, slot := range design.Slots {
		item := rules.techs.GetHullComponent(slot.HullComponent)
		if design.Spec.Starbase && item.Category == "Orbital" {
			costOrbital = costOrbital.Add(item.Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(10 * rules.StarbaseComponentCostReduction))
		} else {
			cost = cost.Add(item.Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(10 * rules.StarbaseComponentCostReduction))
		}
	}

	if design.Spec.Starbase {
		cost = cost.DivideByInt(10*rules.StarbaseComponentCostReduction, true)
	}
	return cost.Add(costOrbital)
}
