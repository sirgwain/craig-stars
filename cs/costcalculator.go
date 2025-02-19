package cs

import (
	"fmt"
	"math"
)

// The CostCalculator interface is used to calculate costs of single items or starbase upgrades
// This is used by planetary production and estimating production queue completion
type CostCalculator interface {
	StarbaseUpgradeCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design, newDesign *ShipDesign) (Cost, error)
	CostOfOne(player *Player, item ProductionQueueItem) (Cost, error)
	GetDesignCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design *ShipDesign) (Cost, error)
	GetTechCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, tech Tech) Cost
}

func NewCostCalculator() CostCalculator {
	return &costCalculate{}
}

type costCalculate struct {
}

// Returns the cost efficiency ratio for 2 Cost structs as a float64
// by dividing their respective total costs
// (numeratorTotal / denominatorTotal)
//
// costTypes indicates the cost types to be considered in analysis (defaults to all);
// function will panic if too many are provided
func GetCostEfficiencyRatio[T number](numerator, denominator cost[T], costTypes ...CostType) (costRatio float64) {
	// TODO: Add weighting support by replacing CostTypes by a single CostFloat64 containing weight values
	if len(costTypes) > 4 {
		panic(fmt.Sprintf("GetCostEfficiencyRatio called with too many cost types: \n%v", costTypes))
	} else if len(costTypes) == 0 {
		costTypes = CostTypes[:] // no cost types provided means we include everything
	}
	var hcTally, otherTally T
	for _, ct := range costTypes {
		hcTally += numerator.GetAmount(ct)
		otherTally += denominator.GetAmount(ct)
	}
	return float64(hcTally) / float64(otherTally)
}

// GetTechCost is an exported method to get a player's cost for a tech, used by wasm
func (c *costCalculate) GetTechCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, tech Tech) Cost {
	return getPlayerCost(tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).ToCost()
}

// Get baseline cost for this technology given a player's tech levels, minaturization stats & racial cost modifiers
//
// Rounds value in accordance with base game's cost calcs,
// but returns it as a floating point cost
// to allow combination with other float multipliers down the line
func getPlayerCost(tech Tech, techLevels TechLevel, miniaturizationSpec MiniaturizationSpec, costOffset TechCostOffset) (techCost CostFloat64) {
	// figure out miniaturization discounts - base cost is reduced by
	// 4% per tech level we have above the tech's requirements.
	// We count the smallest difference among all fields, so if you have
	// level 10 energy & 12 bio and a tech costs 9 energy & 4 bio,
	// the smallest level difference you have is the 1 energy level (not 8 bio levels)
	numTechLevelsAboveRequired := techLevels.LevelsAbove(tech.Requirements.TechLevel)

	// for starter techs, they are all 0 requirements, so just use our lowest field
	if numTechLevelsAboveRequired == math.MaxInt {
		numTechLevelsAboveRequired = techLevels.LowestLevel()
	}

	var miniaturizationFactor float64
	if numTechLevelsAboveRequired > 0 {
		// Ex: 5 tech levels * 4% discount per level = 20% cheaper (0.8x price modifier)
		miniaturizationFactor = 1 - math.Min(miniaturizationSpec.MiniaturizationMax,
			miniaturizationSpec.MiniaturizationPerLevel*float64(numTechLevelsAboveRequired))
	} else {
		// New techs cost BET races 2x and will
		// have 0 for miniaturization
		miniaturizationFactor = miniaturizationSpec.NewTechCostFactor
	}

	techCost = MultiplyCost(tech.Cost.ToCostFloat64(), miniaturizationFactor).Round(func(f float64) float64 {
		if f > 0 && f < 1 {
			return 1 // prevents items costing <0.5 from rounding to 0
		}
		return roundHalfTowards0(f)
	})

	// apply any tech cost offsets multiplicatively,
	// using jank rounding to simulate OG Stars!' int calculations
	var costMulti float64 = 1
	for tag := range tech.Tags {
		costMulti *= 1 + costOffset[tag]
	}
	techCost = techCost.Add(MultiplyCost(techCost, costMulti-1).Round(roundHalfTowards0))

	return techCost.Round(func(f float64) float64 {
		if f > 0 && f < 1 {
			return 1 // prevents total item cost from going below 1
		}
		return f
	})
}

// Calculate the upgrade cost for replacing one starbase design with another
//
// Takes into account part replacement refunds and minimum costs
func (c *costCalculate) StarbaseUpgradeCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design, newDesign *ShipDesign) (Cost, error) {
	if design.SlotsEqual(newDesign.Slots) && design.Hull == newDesign.Hull {
		// Exact same base; no calcs needed
		return Cost{}, nil
	}

	cost := CostFloat64{}
	minCost := CostFloat64{}
	oldComponents := map[*TechHullComponent]int{} // Maps hull component to quantity
	newComponents := map[*TechHullComponent]int{}
	oldComponentsByCategory := map[TechCategory][]*TechHullComponent{} // Maps component category to hull components
	newComponentsByCategory := map[TechCategory][]*TechHullComponent{}
	categories := map[TechCategory]bool{}
	// wrapper function so I don't have to write everything out all the time
	getItemCost := func(t Tech) CostFloat64 {
		return getPlayerCost(t, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset)
	}

	// First of all, check to see if the hulls even EXIST in the first place
	// and return an error if they don't
	oldHull := rules.techs.GetHull(design.Hull)
	newHull := rules.techs.GetHull(newDesign.Hull)
	if oldHull == nil {
		return Cost{}, fmt.Errorf("starbase hull %q of old design was not found in tech store", design.Hull)
	} else if newHull == nil {
		return Cost{}, fmt.Errorf("starbase hull %q of new design was not found in tech store", newDesign.Hull)
	}

	// If the hulls are different, add (newHullCost - 0.5*OldHullCost) to our conversion cost
	if design.Hull != newDesign.Hull {
		oldHullCost := getItemCost(oldHull.Tech)
		newHullCost := getItemCost(newHull.Tech)
		cost = cost.Add(newHullCost).Subtract(MultiplyCost(oldHullCost, rules.StarbaseHullRefundFactor))
	}

	// Next, iterate through both designs' slots and tally up items in each
	// Also check if they even exist (and return error if so)
	for i := range Max(len(design.Slots), len(newDesign.Slots)) {
		// don't wanna index arrays out of bounds!
		if i < len(design.Slots) {
			hc := rules.techs.GetHullComponent(design.Slots[i].HullComponent)
			if hc == nil && design.Slots[i].HullComponent != "" { // if for whatever reason the slot has no component, just assume it's empty
				return Cost{}, fmt.Errorf("component %q of old design was not found in tech store", design.Slots[i].HullComponent)
			}
			oldComponents[hc] += design.Slots[i].Quantity
		}
		if i < len(newDesign.Slots) {
			hc := rules.techs.GetHullComponent(newDesign.Slots[i].HullComponent)
			if hc == nil && design.Slots[i].HullComponent != "" {
				return Cost{}, fmt.Errorf("component %q of new design was not found in tech store", newDesign.Slots[i].HullComponent)
			}
			newComponents[hc] += newDesign.Slots[i].Quantity
		}
	}

	// Iterate through all new parts in our new base list
	// to see if they are present on the old base
	// and remove any duplicates
	if len(oldComponents) > 0 && len(newComponents) > 0 {
		for item, newQuantity := range newComponents {
			oldQuantity := oldComponents[item]
			switch {
			case newQuantity == oldQuantity:
				// same amount of item in both bases; remove from both
				delete(oldComponents, item)
				delete(newComponents, item)
			case newQuantity > oldQuantity:
				// More copies of item in new design; remove duplicates from new base list
				newComponentsByCategory[item.Tech.Category] = append(newComponentsByCategory[item.Tech.Category], item)
				categories[item.Tech.Category] = true
				newComponents[item] -= oldQuantity
				delete(oldComponents, item)
			default:
				// More copies of item in original design (or item doesn't exist on new base)
				// remove duplicates from old base list
				oldComponents[item] -= newQuantity
				delete(newComponents, item)
			}
		}
	}

	if len(oldComponents) == 0 {
		// every item in our old base is also present in the new one;
		// we can just tally up all our costs for the new stuff and be done for the day
		for item, leftoverQty := range newComponents {
			if item.Tech.Category == TechCategoryOrbital {
				cost = cost.Add(MultiplyCost(getItemCost(item.Tech), leftoverQty))
			} else {
				cost = cost.Add(MultiplyCost(getItemCost(item.Tech), float64(leftoverQty)*rules.StarbaseComponentCostReduction))
			}
		}
		return MultiplyCost(cost, raceSpec.StarbaseCostFactor).Round(math.Ceil).ToCost().MinZero(), nil
	} else {
		// We have some leftovers to take care of...
		// Loop through any remaining items from the old base and add them to our category list
		// everything from the new base is already present, so this just ensures everything gets from both
		for item := range oldComponents {
			oldComponentsByCategory[item.Tech.Category] = append(oldComponentsByCategory[item.Tech.Category], item)
			categories[item.Tech.Category] = true
		}
	}

	// At this point, we should have 4 maps in total: 2 for each base design
	// XXXComponents contains all components unique to each base mapped to their respective quantities
	// XXXComponentsByCategory contains all categories present in each base
	// mapped to a slice of all components of that category on said base
	// Now, all that's left to do are the cost calcs!

	// Tally up costs per category present on either base
	for category := range categories {
		oldCost := CostFloat64{}
		newCost := CostFloat64{}

		for _, oldItem := range oldComponentsByCategory[category] {
			if category == TechCategoryOrbital {
				oldCost = oldCost.Add(MultiplyCost(getItemCost(oldItem.Tech), oldComponents[oldItem]))
			} else {
				oldCost = oldCost.Add(MultiplyCost(getItemCost(oldItem.Tech), float64(oldComponents[oldItem])*rules.StarbaseComponentCostReduction))
			}
		}
		for _, newItem := range newComponentsByCategory[category] {
			if category == TechCategoryOrbital {
				newCost = newCost.Add(MultiplyCost(getItemCost(newItem.Tech), newComponents[newItem]))
			} else {
				newCost = newCost.Add(MultiplyCost(getItemCost(newItem.Tech), float64(newComponents[newItem])*rules.StarbaseComponentCostReduction))
			}
		}

		// apply first part of costs to tally (70% of new item cost - 70% of old item cost);
		// this is the portion that can be reduced by normal rebates
		cost = cost.Add(MultiplyCost(newCost.Subtract(oldCost), 0.7))

		// add on rest of the cost after category specific rebates
		// higher of (20% new item cost, 30% new item cost - 10% old item cost)
		// if no old item exists, you pay 100% of this chunk
		adjCost := MultiplyCost(newCost, 0.2).Max(
			MultiplyCost(newCost, 0.3).Subtract(MultiplyCost(oldCost, 0.1)))
		cost = cost.Add(adjCost)
		minCost = minCost.Add(adjCost)
	}

	return MultiplyCost(cost.Max(minCost), raceSpec.StarbaseCostFactor).Round(math.Ceil).ToCost().MinZero(), nil
}

// Get the cost of one item in a production queue, for a player
func (c *costCalculate) CostOfOne(player *Player, item ProductionQueueItem) (Cost, error) {
	cost := player.Race.Spec.Costs[item.Type]
	if item.Type == QueueItemTypeStarbase || item.Type == QueueItemTypeShipToken {
		if item.design != nil {
			cost = item.design.Spec.Cost // should never happen since it isn't called for designs
		} else {
			return Cost{}, fmt.Errorf("ship design #%d not populated in production queue during CostOfOne", item.DesignNum)
		}
	}
	return cost, nil
}

// Get cost of a given ship or new starbase design
func (c *costCalculate) GetDesignCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design *ShipDesign) (Cost, error) {

	hull := rules.techs.GetHull(design.Hull)
	if hull == nil {
		return Cost{}, fmt.Errorf("hull design %q was not found in tech store", design.Hull)
	}
	starbase := hull.Starbase

	costTally := getPlayerCost(hull.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset)

	// iterate through slots and tally prices up one by one
	for _, slot := range design.Slots {
		if slot.HullComponent == "" {
			// slot is empty; move on
			continue
		}

		item := rules.techs.GetHullComponent(slot.HullComponent)
		if item == nil {
			return Cost{}, fmt.Errorf("component %q in design slots was not found in tech store", slot.HullComponent)
		}
		hcCost := MultiplyCost(getPlayerCost(item.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset), slot.Quantity)
		if starbase && item.Category != TechCategoryOrbital {
			hcCost = MultiplyCost(hcCost, rules.StarbaseComponentCostReduction)
		}
		costTally = costTally.Add(hcCost)
	}

	if starbase {
		costTally = MultiplyCost(costTally, raceSpec.StarbaseCostFactor)
	}
	return costTally.Round(math.Ceil).ToCost(), nil
}
