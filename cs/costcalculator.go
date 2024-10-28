package cs

import (
	"fmt"
	"maps"
	"math"
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

// A costFloat64 is otherwise identical to a regular Cost struct, but uses float64s instead of ints
// used for internal cost calculations before being cast back into a regular Cost
type costFloat64 struct {
	Ironium   float64 `json:"ironium,omitempty"`
	Boranium  float64 `json:"boranium,omitempty"`
	Germanium float64 `json:"germanium,omitempty"`
	Resources float64 `json:"resources,omitempty"`
}

func (c costFloat64) getAmount(costType CostType) float64 {
	switch costType {
	case Ironium:
		return c.Ironium
	case Boranium:
		return c.Boranium
	case Germanium:
		return c.Germanium
	case Resources:
		return c.Resources
	}
	panic(fmt.Sprintf("GetAmount called with invalid CostType %s", costType))
}

// convert a costFloat64 to an int using the specified rounding method
func (c costFloat64) toCost(roundFunc func(float64) float64) Cost {
	return Cost{
		Ironium:   int(roundFunc(c.Ironium)),
		Boranium:  int(roundFunc(c.Boranium)),
		Germanium: int(roundFunc(c.Germanium)),
		Resources: int(roundFunc(c.Resources)),
	}
}

func (c costFloat64) addFloat64(costType CostType, amount float64) costFloat64 {
	switch costType {
	case Ironium:
		c.Ironium += amount
	case Boranium:
		c.Boranium += amount
	case Germanium:
		c.Germanium += amount
	case Resources:
		c.Resources += amount
	default:
		panic(fmt.Sprintf("addFloat64 called with invalid CostType %s", costType))
	}
	return c
}

func (c costFloat64) add(other costFloat64) costFloat64 {
	return costFloat64{
		Ironium:   c.Ironium + other.Ironium,
		Boranium:  c.Boranium + other.Boranium,
		Germanium: c.Germanium + other.Germanium,
		Resources: c.Resources + other.Resources,
	}
}

func (c costFloat64) subtract(other costFloat64) costFloat64 {
	return costFloat64{
		Ironium:   c.Ironium - other.Ironium,
		Boranium:  c.Boranium - other.Boranium,
		Germanium: c.Germanium - other.Germanium,
		Resources: c.Resources - other.Resources,
	}
}

func (c costFloat64) multiply(factor float64) costFloat64 {
	return costFloat64{
		Ironium:   c.Ironium * factor,
		Boranium:  c.Boranium * factor,
		Germanium: c.Germanium * factor,
		Resources: c.Resources * factor,
	}
}

func (c costFloat64) divide(divisor float64) costFloat64 {
	return costFloat64{
		Ironium:   c.Ironium / divisor,
		Boranium:  c.Boranium / divisor,
		Germanium: c.Germanium / divisor,
		Resources: c.Resources / divisor,
	}
}

// Return greater of 2 cost structs for all ResourceTypes separately
func (c costFloat64) max(other costFloat64) costFloat64 {
	return costFloat64{
		Ironium:   math.Max(c.Ironium, other.Ironium),
		Boranium:  math.Max(c.Boranium, other.Boranium),
		Germanium: math.Max(c.Germanium, other.Germanium),
		Resources: math.Max(c.Resources, other.Resources),
	}
}

// round a cost struct's values with passed in function
func (c costFloat64) round(roundFunc func(float64) float64) costFloat64 {
	return costFloat64{
		Ironium:   roundFunc(c.Ironium),
		Boranium:  roundFunc(c.Boranium),
		Germanium: roundFunc(c.Germanium),
		Resources: roundFunc(c.Resources),
	}
}


// get the upgrade cost for replacing a starbase with another
//
// Takes into account part replacement costs and minimum costs
func (p *costCalculate) StarbaseUpgradeCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design, newDesign *ShipDesign) (Cost, error) {
	if design.SlotsEqual(newDesign.Slots) && design.Hull == newDesign.Hull {
		// Exact same base; no calcs needed
		return Cost{}, nil
	}

	credit := costFloat64{}
	cost := costFloat64{}
	minCost := costFloat64{}
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
	if design.Hull != newDesign.Hull {
		oldHullCost := oldHull.Tech.GetPlayerCostFloat(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset)
		newHullCost := newHull.Tech.GetPlayerCostFloat(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset)
		cost = cost.add(newHullCost).subtract((oldHullCost.multiply(rules.StarbaseHullRefundFactor)))
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
				// add extras to old base list
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
				cost = cost.add(item.GetPlayerCostFloat(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(qty)))
			} else {
				cost = cost.add(item.GetPlayerCostFloat(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(qty) * rules.StarbaseComponentCostReduction))
			}
		}
		return cost.multiply(raceSpec.StarbaseCostFactor).toCost(math.Ceil), nil
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
	for category := range categories {
		oldCost := costFloat64{}
		newCost := costFloat64{}

		for _, oldItem := range oldComponentsByCategory[category] {
			if category == TechCategoryOrbital {
				oldCost = oldCost.add(oldItem.GetPlayerCostFloat(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(oldComponents[oldItem])))
			} else {
				oldCost = oldCost.add(oldItem.GetPlayerCostFloat(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(oldComponents[oldItem]) * rules.StarbaseComponentCostReduction))
			}
		}
		for _, newItem := range newComponentsByCategory[category] {
			if category == TechCategoryOrbital {
				newCost = newCost.add(newItem.GetPlayerCostFloat(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(newComponents[newItem])))
			} else {
				newCost = newCost.add(newItem.GetPlayerCostFloat(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(newComponents[newItem]) * rules.StarbaseComponentCostReduction))
			}
		}

		// Apply lower (70%) rebate to credit tally (up to 70% of new item value)
		// Apply difference between 2 discounts (10%) to this item category only, up to 10% of the original item value
		// (for a total of 80% rebate for same category items)

		// Compute costs for each resource type separately (I/B/G/R)
		for _, costType := range CostTypes {
			// extract float values for items
			oldCostFloat := oldCost.getAmount(costType)
			newCostFloat := newCost.getAmount(costType)
			if oldCostFloat == 0 && newCostFloat == 0 {
				continue
			}

			differentCategoryRebate := 0.7 * oldCostFloat

			// add global rebate to credit tally
			credit = credit.addFloat64(costType, differentCategoryRebate)

			// Consume global credit tally to reduce new item price from 100% to 30%
			// If this turns credit negative, no problem!
			// We add it to Cost at the end anyways
			adjCost := 0.3 * newCostFloat
			credit = credit.addFloat64(costType, -(newCostFloat - adjCost))
			
			// add on category specific rebates and tack onto minimum cost 
			adjCost = math.Max(0.2*newCostFloat, adjCost-0.1*oldCostFloat)
			cost = cost.addFloat64(costType, adjCost)
			minCost = minCost.addFloat64(costType, adjCost)
		}
	}

	return cost.subtract(credit).max(minCost).toCost(math.Ceil).MinZero(), nil
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
	starbase := hull.Starbase


	cost := hull.Tech.GetPlayerCostFloat(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset)

	// iterate through slots and tally prices up
	for _, slot := range design.Slots {
		item := rules.techs.GetHullComponent(slot.HullComponent)
		if item == nil {
			return Cost{}, fmt.Errorf("component %s in design slots not found in tech store", slot.HullComponent)
		}
		hcCost := item.Tech.GetPlayerCostFloat(techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(slot.Quantity))
		if starbase && item.Category != TechCategoryOrbital {
			cost = cost.add(hcCost.multiply(rules.StarbaseComponentCostReduction))
		} else {
			cost = cost.add(hcCost)
		}
	}

	if starbase {
		cost = cost.multiply(raceSpec.StarbaseCostFactor)
	}
	return cost.toCost(math.Ceil), nil
}
