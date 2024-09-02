package cs

import (
	"fmt"
	"maps"
)

// The CostCalculator interface is used to calculate costs of single items or starbase upgrades
// This is used by planetary production and estimating production queue completion
type CostCalculator interface {
	StarbaseUpgradeCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design, newDesign *ShipDesign) (Cost, error)
	CostOfOne(player *Player, item ProductionQueueItem) (Cost, error)
	GetDesignCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design *ShipDesign) (Cost, error)
}

func NewCostCalculator() CostCalculator {
	return &costCalculate{}
}

type costCalculate struct {
}

// get the upgrade cost for replacing a starbase with another
//
// Takes into account part replacement costs and minimum costs
func (p *costCalculate) StarbaseUpgradeCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design, newDesign *ShipDesign) (Cost, error) {
	if design.SlotsEqual(newDesign.Slots) && design.Hull == newDesign.Hull {
		// Exact same base; no calcs needed
		return Cost{}, nil
	}

	credit := Cost{}
	cost := Cost{}
	minCost := Cost{}
	oldComponents := map[*TechHullComponent]int{} // Maps hull component to quantity
	newComponents := map[*TechHullComponent]int{}
	oldComponentsByCategory := map[TechCategory][]*TechHullComponent{} // Maps component category to hull components
	newComponentsByCategory := map[TechCategory][]*TechHullComponent{}

	// First of all, check to see if the hulls even EXIST in the first place
	// and return an error if they don't
	oldHull := rules.techs.GetHull(design.Hull)
	newHull := rules.techs.GetHull(newDesign.Hull)
	if oldHull == nil {
		return Cost{}, fmt.Errorf("starbase hull %s of old design not found in tech store", design.Hull)
	} else if newHull == nil {
		return Cost{}, fmt.Errorf("starbase hull %s of new design not found in tech store", newDesign.Hull)
	}

	// If the hulls are different, add (newHullCost - 0.5*OldHullCost)
	// Multiplied by 1000 * componentCostReduction for rounding purposes
	if design.Hull != newDesign.Hull {
		oldHullCost := oldHull.Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(500 * rules.StarbaseComponentCostReduction)
		newHullCost := newHull.Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(1000 * rules.StarbaseComponentCostReduction)
		cost = cost.Add(newHullCost).Minus(oldHullCost)
	}

	// Next, iterate through both designs' slots and tally up items in each
	// Also check if they even exist (and return error if so)
	for i := 0; i < MaxInt(len(design.Slots), len(newDesign.Slots)); i++ {
		// don't wanna index arrays out of bounds!
		if i < len(design.Slots) {
			hc := rules.techs.GetHullComponent(design.Slots[i].HullComponent)
			if hc != nil {
				oldComponents[hc] += design.Slots[i].Quantity
			} else {
				return Cost{}, fmt.Errorf("component %s of old design not found in tech store", design.Slots[i].HullComponent)
			}
		}
		if i < len(newDesign.Slots) {
			hc := rules.techs.GetHullComponent(newDesign.Slots[i].HullComponent)
			if hc != nil {
				newComponents[hc] += newDesign.Slots[i].Quantity
			} else {
				return Cost{}, fmt.Errorf("component %s of new design not found in tech store", newDesign.Slots[i].HullComponent)
			}
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
			if item.Tech.Category == TechCategoryOrbital {
				cost = cost.Add(item.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(qty).MultiplyInt(1000 * rules.StarbaseComponentCostReduction))
			} else {
				cost = cost.Add(item.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(qty).MultiplyInt(1000))
			}
		}
		return cost.DivideByInt(int(roundHalfDown(1000*float64(rules.StarbaseComponentCostReduction)/raceSpec.StarbaseCostFactor)), true), nil
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
	categories := map[TechCategory][]*TechHullComponent{}
	maps.Copy(categories, oldComponentsByCategory)
	maps.Copy(categories, newComponentsByCategory)

	// Tally up costs per category
	// We multiply everything by 200x here (and by 10x in the following step) & divide by 2000 later to minimize rounding errors
	for category := range categories {
		oldCost := Cost{}
		newCost := Cost{}

		for _, oldItem := range oldComponentsByCategory[category] {
			if category == TechCategoryOrbital {
				oldCost = oldCost.Add(oldItem.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(oldComponents[oldItem] * 100 * rules.StarbaseComponentCostReduction))
			} else {
				oldCost = oldCost.Add(oldItem.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(oldComponents[oldItem] * 100))
			}
		}
		for _, newItem := range newComponentsByCategory[category] {
			if category == TechCategoryOrbital {
				newCost = newCost.Add(newItem.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(newComponents[newItem] * 100 * rules.StarbaseComponentCostReduction))
			} else {
				newCost = newCost.Add(newItem.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(newComponents[newItem] * 100))
			}
		}

		// Apply lower (70%) rebate to credit tally (up to 70% of the actual item value)
		// Apply difference between 2 discounts (10%) to this item category only, up to 10% of the original item value
		// Costs here are multiplied by factor of 10 to be divided later on

		// Compute costs for each resource type separately (I/B/G/R)
		for _, costType := range CostTypes {
			// extract float values for items
			oldCostInt := oldCost.GetAmount(costType)
			newCostInt := newCost.GetAmount(costType)
			if oldCostInt == 0 && newCostInt == 0 {
				continue
			}

			sameCategoryRebate := oldCostInt * 8
			differentCategoryRebate := oldCostInt * 7

			// Add global rebate to credit tally
			credit = credit.AddInt(costType, differentCategoryRebate)

			// Consume global credit tally to reduce new item price from 100% to 30%
			// If this turns credit negative, no problem!
			// We add it to Cost at the end anyways
			adjCost := 3 * newCostInt
			credit = credit.AddInt(costType, -(10*newCostInt - adjCost))

			// Add on category specific rebate
			categoryRebate := sameCategoryRebate - differentCategoryRebate
			adjCost = MaxInt(2*newCostInt, adjCost-categoryRebate)
			cost = cost.AddInt(costType, adjCost)
			minCost = minCost.AddInt(costType, adjCost)
		}
	}
	cost = cost.Minus(credit).MinZero()
	cost = cost.Max(minCost).DivideByInt(int(roundHalfDown(1000*float64(rules.StarbaseComponentCostReduction)/raceSpec.StarbaseCostFactor)), true)
	return cost, nil
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
func (p *costCalculate) GetDesignCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design *ShipDesign) (Cost, error) {

	hull := rules.techs.GetHull(design.Hull)
	if hull == nil {
		return Cost{}, fmt.Errorf("hull design %s not found in tech store", design.Hull)
	}

	cost := hull.Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(1000)
	if design.Spec.Starbase {
		// multiply by cost factor because we divide by it later
		cost = cost.MultiplyInt(rules.StarbaseComponentCostReduction)
	}

	// iterate through slots and tally prices up
	for _, slot := range design.Slots {
		item := rules.techs.GetHullComponent(slot.HullComponent)
		if item == nil {
			return Cost{}, fmt.Errorf("component %s in design slots not found in tech store", slot.HullComponent)
		}
		hcCost := item.Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(slot.Quantity)
		if design.Spec.Starbase && item.Category == TechCategoryOrbital {
			cost = cost.Add(hcCost.MultiplyInt(1000 * rules.StarbaseComponentCostReduction))
		} else {
			cost = cost.Add(hcCost.MultiplyInt(1000))
		}
	}

	if design.Spec.Starbase {
		return cost.DivideByInt(int(roundHalfDown(1000*float64(rules.StarbaseComponentCostReduction)/raceSpec.StarbaseCostFactor)), true), nil
	} else {
		return cost.DivideByInt(1000, true), nil
	}
}
