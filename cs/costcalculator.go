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
}

func NewCostCalculator() CostCalculator {
	return &costCalculate{}
}

type costCalculate struct {
}

// Returns the cost efficiency ratio for 2 costFloat64 structs
// by dividing their respective total costs
// (numeratorTotal / denominatorTotal)
//
// costTypes indicates the cost types to be considered in analysis (defaults to all);
// function will panic if too many are provided
func GetCostEfficiencyRatio(numerator, denominator CostFloat64, costTypes ...CostType) (costRatio float64) {
	if len(costTypes) > 4 {
		panic(fmt.Sprintf("GetCostEfficiencyRatio called with too many cost types; %v", costTypes))
	} else if len(costTypes) == 0 {
		costTypes = CostTypes[:] // no cost types provided means we include everything
	}
	var hcTally, otherTally float64
	for _, ct := range costTypes {
		hcTally += numerator.GetAmount(ct)
		otherTally += denominator.GetAmount(ct)
	}
	return hcTally / otherTally
}

// Get baseline cost for this technology given a player's tech levels, minaturization stats & racial cost modifiers
//
// Returns floating point cost for extra precision
func getPlayerCost(tech Tech, techLevels TechLevel, miniaturizationSpec MiniaturizationSpec, costOffset TechCostOffset) CostFloat64 {
	// figure out miniaturization
	// this is 4% per level above the required tech we have.
	// We count the smallest diff, i.e. if you have
	// tech level 10 energy, 12 bio and the tech costs 9 energy, 4 bio
	// the smallest level difference you have is 1 energy level (not 8 bio levels)

	// From the diff between the player level and the requirements, find the lowest difference
	// i.e. 1 energy level in the example above
	numTechLevelsAboveRequired := techLevels.LevelsAbove(tech.Requirements.TechLevel)

	// for starter techs, they are all 0 requirements, so just use our lowest field
	if numTechLevelsAboveRequired == math.MaxInt {
		numTechLevelsAboveRequired = techLevels.LowestLevel()
	}

	// As players gain tech levels, lower leveled items get cheaper.
	// Players start off with full priced techs, but every additional level past the requirements
	// reduces the cost slightly, up to a certain fraction of the initial price.
	var miniaturizationFactor float64
	if numTechLevelsAboveRequired > 0 {
		// Ex: 5 tech levels * 4% discount per level = 20% cheaper (0.8x price modifier)
		miniaturizationFactor = 1 - math.Min(miniaturizationSpec.MiniaturizationMax,
			miniaturizationSpec.MiniaturizationPerLevel*float64(numTechLevelsAboveRequired))
	} else {
		// New techs cost BET races 2x and will
		// have 0 for miniaturization.
		miniaturizationFactor = miniaturizationSpec.NewTechCostFactor
	}

	techCost := MultiplyCost(tech.Cost.ToCostFloat64(), miniaturizationFactor).Round(roundHalfDown)

	// apply any tech cost offsets
	highestCostMulti := 1.0
	for tag := range tech.Tags {
		highestCostMulti = math.Min(1+costOffset[tag], highestCostMulti) // only take the lowest single bonus
	}

	return MultiplyCost(techCost, highestCostMulti)
}

// get the upgrade cost for replacing a starbase with another
//
// Takes into account part replacement costs and minimum costs
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

	// First of all, check to see if the hulls even EXIST in the first place
	// and return an error if they don't
	oldHull := rules.techs.GetHull(design.Hull)
	newHull := rules.techs.GetHull(newDesign.Hull)
	if oldHull == nil {
		return Cost{}, fmt.Errorf("starbase hull %s of old design not found in tech store", design.Hull)
	} else if newHull == nil {
		return Cost{}, fmt.Errorf("starbase hull %s of new design not found in tech store", newDesign.Hull)
	}

	// If the hulls are different, add (newHullCost - 0.5*OldHullCost) to our conversion cost
	if design.Hull != newDesign.Hull {
		oldHullCost := getPlayerCost(oldHull.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset)
		newHullCost := getPlayerCost(newHull.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset)
		cost = cost.Add(newHullCost).Subtract(MultiplyCost(oldHullCost, rules.StarbaseHullRefundFactor))
	}

	// Next, iterate through both designs' slots and tally up items in each
	// Also check if they even exist (and return error if so)
	for i := range Max(len(design.Slots), len(newDesign.Slots)) {
		// don't wanna index arrays out of bounds!
		if i < len(design.Slots) {
			hc := rules.techs.GetHullComponent(design.Slots[i].HullComponent)
			if hc == nil {
				return Cost{}, fmt.Errorf("component %s of old design not found in tech store", design.Slots[i].HullComponent)
			}
			oldComponents[hc] += design.Slots[i].Quantity
		}
		if i < len(newDesign.Slots) {
			hc := rules.techs.GetHullComponent(newDesign.Slots[i].HullComponent)
			if hc == nil {
				return Cost{}, fmt.Errorf("component %s of new design not found in tech store", newDesign.Slots[i].HullComponent)
			}
			newComponents[hc] += newDesign.Slots[i].Quantity
		}
	}

	// Iterate through all new parts in list to see if they are present on the old base
	// and remove any duplicates we find
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
				newComponents[item] = newQuantity - oldQuantity
				delete(oldComponents, item)
			default:
				// More copies of item in original design (or item doesn't exist on new base)
				// remove duplicates from old base list
				oldComponents[item] = oldQuantity - newQuantity
				delete(newComponents, item)
			}
		}
	}

	if len(oldComponents) == 0 {
		// no items in old base not also present in the new one
		// We can just tally up all our costs for the new stuff and be done for the day
		for item, qty := range newComponents {
			if item.Tech.Category == TechCategoryOrbital {
				cost = cost.Add(MultiplyCost(getPlayerCost(item.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset), qty))
			} else {
				cost = cost.Add(MultiplyCost(getPlayerCost(item.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset), float64(qty)*rules.StarbaseComponentCostReduction))
			}
		}
		return MultiplyCost(cost, raceSpec.StarbaseCostFactor).Round(math.Ceil).ToCost().MinZero(), nil
	} else {
		// Loop through any remaining items from old base and add to category list
		// everything from the new base is already on, so this ensures everything gets checked
		for item := range oldComponents {
			oldComponentsByCategory[item.Tech.Category] = append(oldComponentsByCategory[item.Tech.Category], item)
			categories[item.Tech.Category] = true
		}
	}

	// At this point, we should have 4 maps in total: 2 for each base design
	// Components contains all components unique to each base mapped to their respective quantities
	// ComponentsByCategory contains all categories present in each base
	// mapped to a list of all components of that category on said base
	// Now, all that's left is the cost calcs

	// Tally up costs per category
	for category := range categories {
		oldCost := CostFloat64{}
		newCost := CostFloat64{}

		for _, oldItem := range oldComponentsByCategory[category] {
			if category == TechCategoryOrbital {
				oldCost = oldCost.Add(MultiplyCost(getPlayerCost(oldItem.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset), oldComponents[oldItem]))
			} else {
				oldCost = oldCost.Add(MultiplyCost(getPlayerCost(oldItem.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset), float64(oldComponents[oldItem])*rules.StarbaseComponentCostReduction))
			}
		}
		for _, newItem := range newComponentsByCategory[category] {
			if category == TechCategoryOrbital {
				newCost = newCost.Add(MultiplyCost(getPlayerCost(newItem.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset), newComponents[newItem]))
			} else {
				newCost = newCost.Add(MultiplyCost(getPlayerCost(newItem.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset), float64(newComponents[newItem])*rules.StarbaseComponentCostReduction))
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
func (p *costCalculate) CostOfOne(player *Player, item ProductionQueueItem) (Cost, error) {
	cost := player.Race.Spec.Costs[item.Type]
	if item.Type == QueueItemTypeStarbase || item.Type == QueueItemTypeShipToken {
		if item.design != nil {
			cost = item.design.Spec.Cost // should never happen since it isn't called for designs
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
		return Cost{}, fmt.Errorf("hull design \"%s\" not found in tech store", design.Hull)
	}
	starbase := hull.Starbase

	cost := getPlayerCost(hull.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset)

	// iterate through slots and tally prices up
	for _, slot := range design.Slots {
		item := rules.techs.GetHullComponent(slot.HullComponent)
		if slot.HullComponent == "" {
			// slot is empty; move on
			continue
		}
		if item == nil {
			return Cost{}, fmt.Errorf("component \"%s\" in design slots not found in tech store", slot.HullComponent)
		}
		hcCost := MultiplyCost(getPlayerCost(item.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset), slot.Quantity)
		if starbase && item.Category != TechCategoryOrbital {
			cost = cost.Add(MultiplyCost(hcCost, rules.StarbaseComponentCostReduction))
		} else {
			cost = cost.Add(hcCost)
		}
	}

	if starbase {
		cost = MultiplyCost(cost, raceSpec.StarbaseCostFactor)
	}
	return cost.Round(math.Ceil).ToCost(), nil
}
