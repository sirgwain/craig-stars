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
		// identical hulls & components; no calcs needed
		return Cost{}
	}
	
	credit := Cost{}
	cost := Cost{}
	oldComponents := map[*TechHullComponent]int{} // Maps hull component to quantity
	newComponents := map[*TechHullComponent]int{}
	oldComponentsByCategory := map[TechCategory][]*TechHullComponent{}
	newComponentsByCategory := map[TechCategory][]*TechHullComponent{}


	// iterate through both designs' slots and tally up items in each 
	for i := 0; i <= MaxInt(len(design.Slots), len(newDesign.Slots)); i++ {
		// don't wanna index arrays out of bounds!
		if i <= len(design.Slots) {
			oldComponents[rules.techs.GetHullComponent(design.Slots[i].HullComponent)] += design.Slots[i].Quantity
		}
		if i <= len(newDesign.Slots) {
			newComponents[rules.techs.GetHullComponent(newDesign.Slots[i].HullComponent)] += newDesign.Slots[i].Quantity
		}
	}

	// Iterate through all new parts in list to see if they are present on the old base
	// to create a list of all unique components 
	for item, newQuantity := range newComponents {
		oldQuantity := oldComponents[item] // defaults to 0 if not on base
		if newQuantity == oldQuantity {
			// same amount of item in both bases; remove from both
			delete(oldComponents, item)
			delete(newComponents, item)
		} else if newQuantity > oldQuantity {
			// More copies of item in new design; add to new base list
			oldComponentsByCategory[item.Tech.Category] = append(oldComponentsByCategory[item.Tech.Category], item)
			newComponents[item] = (newQuantity - oldQuantity)
			delete(oldComponents, item)
		} else {
			// More copies of item in original design (or item doesn't exist on new base)
			newComponentsByCategory[item.Tech.Category] = append(newComponentsByCategory[item.Tech.Category], item)
			oldComponents[item] = (oldQuantity - newQuantity)
			delete(newComponents, item)
		}
	}

	// At this point, we should have 4 maps in total: 2 for each base design
	// ComponentsUnique contains all components unique to each one mapped to their quantity
	// ComponentsByCategory contains a list of all categories present in each base 
	// mapped to a slice of all components on the base for said category 
	// Now, all that's left is the cost calcs

	// Get categories present in either map so we don't have to iterate over them all
	categories := map[TechCategory][]*TechHullComponent{}
	maps.Copy(categories, oldComponentsByCategory)
	maps.Copy(categories, newComponentsByCategory)


	// Tally up costs per category 
	// We multiply everything by 2x here (and by 10x in the following step) to divide by 20 later
	for category := range categories {
		oldCost := Cost{}
		newCost := Cost{}
		for _, oldItem := range oldComponentsByCategory[category] {
			if category == "Orbital" {
				oldCost = oldCost.Add(oldItem.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(oldComponents[oldItem]))
			} else {
				oldCost = oldCost.Add(oldItem.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(oldComponents[oldItem]).MultiplyInt(rules.StarbaseComponentCostReduction))
			}
		}
		for _, newItem := range newComponentsByCategory[category] {
			if category == "Orbital" {
				newCost = newCost.Add(newItem.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(newComponents[newItem]))
			} else {
				newCost = newCost.Add(newItem.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(newComponents[newItem]).MultiplyInt(rules.StarbaseComponentCostReduction))
			}
		}
		
		if len(oldComponentsByCategory[category]) > 0 {
			if len(newComponentsByCategory[category]) > 0 {
				// category present in both bases
				// Apply lower discount to credit 
				// Apply difference between 2 discounts to this category only
				// Costs here are multiplied by factor of 10 (for total of 20x normal value)

				// Compute costs for each resource type separately (I/B/G/R)
				for costType := 0; costType <4; costType ++ {
					// extract float values for items
					oldCostInt := oldCost.GetAmount(costType)
					newCostInt := newCost.GetAmount(costType)
					sameCost := MaxInt(newCostInt*3, newCostInt*10-oldCostInt*8)
					differentCost := MaxInt(newCostInt*3, newCostInt*10-oldCostInt*7)
					
					// add them together
					credit.AddInt(costType, differentCost)
					cost.AddInt(costType, MaxInt(newCostInt-AbsInt(sameCost-differentCost), 0))
				}
			} else {
				// category present in old but not new design; add to credit
				credit.Add(oldCost.MultiplyInt(7))
			}
		} else {
			if len(newComponentsByCategory[category]) > 0 {
				// item category present in new but not old design; add to cost
				cost.Add(oldCost.MultiplyInt(10))
			} else {
				// item category not present in either design
				// Should never happen because we only iterate over 
				// categories present in either map, but shouldn't matter either way
				continue
			}
		}
	}
	
	// Finally, if the hulls are different, tack on hull price as well
	if design.Hull != newDesign.Hull {
		oldHullCost := rules.techs.GetHull(design.Hull).Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(5)
		newHullCost := rules.techs.GetHull(newDesign.Hull).Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(10)
		cost.Add(newHullCost.Minus(oldHullCost))
	}

	return cost.Minus(credit).DivideByInt(10*rules.StarbaseComponentCostReduction).MinZero()
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
			costOrbital = costOrbital.Add(item.Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(10))
		} else {
			cost = cost.Add(item.Tech.GetPlayerCost(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).MultiplyInt(10*rules.StarbaseComponentCostReduction))
		}
	}

	if design.Spec.Starbase {
		cost = cost.DivideByInt(10*rules.StarbaseComponentCostReduction)
	}
	return cost.Add(costOrbital)
}
