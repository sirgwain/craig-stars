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

// A costFloat64 is otherwise identical to a regular Cost struct, but uses float64s instead of ints
// used for internal cost calculations before being cast back into a regular Cost
type costFloat64 struct {
	ironium   float64
	boranium  float64
	germanium float64
	resources float64
}

func (c costFloat64) getAmount(costType CostType) float64 {
	switch costType {
	case Ironium:
		return c.ironium
	case Boranium:
		return c.boranium
	case Germanium:
		return c.germanium
	case Resources:
		return c.resources
	}
	panic(fmt.Sprintf("getAmount called with invalid CostType %s", costType))
}

// convert a costFloat64 to an int using the specified rounding method
func (c costFloat64) toCost(roundFunc func(float64) float64) Cost {
	return Cost{
		Ironium:   int(roundFunc(c.ironium)),
		Boranium:  int(roundFunc(c.boranium)),
		Germanium: int(roundFunc(c.germanium)),
		Resources: int(roundFunc(c.resources)),
	}
}

// convert an int cost to a costfloat64 struct for internal calcs
func costFloat64fromCost(c Cost) costFloat64 {
	return costFloat64{
		ironium:   float64(c.Ironium),
		boranium:  float64(c.Boranium),
		germanium: float64(c.Germanium),
		resources: float64(c.Resources),
	}
}

func (c costFloat64) add(other costFloat64) costFloat64 {
	return costFloat64{
		ironium:   c.ironium + other.ironium,
		boranium:  c.boranium + other.boranium,
		germanium: c.germanium + other.germanium,
		resources: c.resources + other.resources,
	}
}

func (c costFloat64) subtract(other costFloat64) costFloat64 {
	return costFloat64{
		ironium:   c.ironium - other.ironium,
		boranium:  c.boranium - other.boranium,
		germanium: c.germanium - other.germanium,
		resources: c.resources - other.resources,
	}
}

func (c costFloat64) multiply(factor float64) costFloat64 {
	return costFloat64{
		ironium:   c.ironium * factor,
		boranium:  c.boranium * factor,
		germanium: c.germanium * factor,
		resources: c.resources * factor,
	}
}

// Return greater of 2 cost structs for all ResourceTypes separately
func (c costFloat64) max(other costFloat64) costFloat64 {
	return costFloat64{
		ironium:   math.Max(c.ironium, other.ironium),
		boranium:  math.Max(c.boranium, other.boranium),
		germanium: math.Max(c.germanium, other.germanium),
		resources: math.Max(c.resources, other.resources),
	}
}

// round a cost struct's values with passed in function
func (c costFloat64) round(roundFunc func(float64) float64) costFloat64 {
	return costFloat64{
		ironium:   roundFunc(c.ironium),
		boranium:  roundFunc(c.boranium),
		germanium: roundFunc(c.germanium),
		resources: roundFunc(c.resources),
	}
}

// Get baseline cost for this technology given a player's tech levels, minaturization stats & racial cost modifiers
//
// Returns floating point cost for extra precision
func getPlayerCostFloat64(tech Tech, techLevels TechLevel, spec MiniaturizationSpec, costOffset TechCostOffset) costFloat64 {
	// figure out miniaturization
	// this is 4% per level above the required tech we have.
	// We count the smallest diff, i.e. if you have
	// tech level 10 energy, 12 bio and the tech costs 9 energy, 4 bio
	// the smallest level difference you have is 1 energy level (not 8 bio levels)

	levelDiff := TechLevel{-1, -1, -1, -1, -1, -1}

	// From the diff between the player level and the requirements, find the lowest difference
	// i.e. 1 energy level in the example above
	numTechLevelsAboveRequired := math.MaxInt
	if tech.Requirements.Energy > 0 {
		levelDiff.Energy = techLevels.Energy - tech.Requirements.Energy
		numTechLevelsAboveRequired = MinInt(levelDiff.Energy, numTechLevelsAboveRequired)
	}
	if tech.Requirements.Weapons > 0 {
		levelDiff.Weapons = techLevels.Weapons - tech.Requirements.Weapons
		numTechLevelsAboveRequired = MinInt(levelDiff.Weapons, numTechLevelsAboveRequired)
	}
	if tech.Requirements.Propulsion > 0 {
		levelDiff.Propulsion = techLevels.Propulsion - tech.Requirements.Propulsion
		numTechLevelsAboveRequired = MinInt(levelDiff.Propulsion, numTechLevelsAboveRequired)
	}
	if tech.Requirements.Construction > 0 {
		levelDiff.Construction = techLevels.Construction - tech.Requirements.Construction
		numTechLevelsAboveRequired = MinInt(levelDiff.Construction, numTechLevelsAboveRequired)
	}
	if tech.Requirements.Electronics > 0 {
		levelDiff.Electronics = techLevels.Electronics - tech.Requirements.Electronics
		numTechLevelsAboveRequired = MinInt(levelDiff.Electronics, numTechLevelsAboveRequired)
	}
	if tech.Requirements.Biotechnology > 0 {
		levelDiff.Biotechnology = techLevels.Biotechnology - tech.Requirements.Biotechnology
		numTechLevelsAboveRequired = MinInt(levelDiff.Biotechnology, numTechLevelsAboveRequired)
	}

	// for starter techs, they are all 0 requirements, so just use our lowest field
	if numTechLevelsAboveRequired == math.MaxInt {
		numTechLevelsAboveRequired = techLevels.Min()
	}

	// As we learn techs, they get cheaper. We start off with full priced techs, but every additional level of research we learn makes
	// techs cost a little less, maxing out at some discount (i.e. 75% or 80% for races with BET)
	miniaturization := math.Min(spec.MiniaturizationMax, spec.MiniaturizationPerLevel*float64(numTechLevelsAboveRequired))
	// New techs cost BET races 2x
	// new techs will have 0 for miniaturization.
	miniaturizationFactor := spec.NewTechCostFactor
	if numTechLevelsAboveRequired > 0 {
		miniaturizationFactor = 1 - miniaturization
	}

	// apply any tech cost offsets
	cost := costFloat64fromCost(tech.Cost).multiply(miniaturizationFactor).round(roundHalfDown)
	var highestCostMulti float64
	for tag := range tech.Tags {
		highestCostMulti = math.Min(1+costOffset[tag], highestCostMulti)
	}

	return cost
}

// Returns the cost efficiency ratio for 2 TechHullComponents 
// by dividing the techs' total costs
// (numeratorTotal / denominatorTotal)
//
// costTypes indicate the cost types to be considered (defaults to all); 
// function will panic if too many are provided
func getCostEfficiencyRatio(player *Player, numerator, denominator *TechHullComponent, costTypes ...CostType) float64 {
	if len(costTypes) > 4 {
		panic(fmt.Sprintf("getCostEfficiencyRatio called with incorrect amount of cost types; %v", costTypes))
	} else if len(costTypes) == 0 {
		costTypes = CostTypes[:] // no cost types provided means we include everything
	}
	hcCost := getPlayerCostFloat64(numerator.Tech, player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset)
	otherCost := getPlayerCostFloat64(denominator.Tech, player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset)
	hcTally, otherTally := 0., 0.
	for _, ct := range costTypes {
		hcTally += hcCost.getAmount(ct)
		otherTally += otherCost.getAmount(ct)
	}
	return hcTally / otherTally
}

// get the upgrade cost for replacing a starbase with another
//
// Takes into account part replacement costs and minimum costs
func (c *costCalculate) StarbaseUpgradeCost(rules *Rules, techLevels TechLevel, raceSpec RaceSpec, design, newDesign *ShipDesign) (Cost, error) {
	if design.SlotsEqual(newDesign.Slots) && design.Hull == newDesign.Hull {
		// Exact same base; no calcs needed
		return Cost{}, nil
	}

	cost := costFloat64{}
	minCost := costFloat64{}
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

	// If the hulls are different, add (newHullCost - 0.5*OldHullCost)
	if design.Hull != newDesign.Hull {
		oldHullCost := getPlayerCostFloat64(oldHull.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset)
		newHullCost := getPlayerCostFloat64(newHull.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset)
		cost = cost.add(newHullCost).subtract((oldHullCost.multiply(rules.StarbaseHullRefundFactor)))
	}

	// Next, iterate through both designs' slots and tally up items in each
	// Also check if they even exist (and return error if so)
	for i := 0; i < MaxInt(len(design.Slots), len(newDesign.Slots)); i++ {
		// don't wanna index arrays out of bounds!
		if i < len(design.Slots) {
			hc := rules.techs.GetHullComponent(design.Slots[i].HullComponent)
			if hc != nil { // todo: reverse conditional to have break go first
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
				newComponents[item] = (newQuantity - oldQuantity)
				delete(oldComponents, item)
			default:
				// More copies of item in original design (or item doesn't exist on new base)
				// remove duplicates from old base list
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
				cost = cost.add(getPlayerCostFloat64(item.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(qty)))
			} else {
				cost = cost.add(getPlayerCostFloat64(item.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(qty) * rules.StarbaseComponentCostReduction))
			}
		}
		return cost.multiply(raceSpec.StarbaseCostFactor).toCost(math.Ceil).MinZero(), nil
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
		oldCost := costFloat64{}
		newCost := costFloat64{}

		for _, oldItem := range oldComponentsByCategory[category] {
			if category == TechCategoryOrbital {
				oldCost = oldCost.add(getPlayerCostFloat64(oldItem.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(oldComponents[oldItem])))
			} else {
				oldCost = oldCost.add(getPlayerCostFloat64(oldItem.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(oldComponents[oldItem]) * rules.StarbaseComponentCostReduction))
			}
		}
		for _, newItem := range newComponentsByCategory[category] {
			if category == TechCategoryOrbital {
				newCost = newCost.add(getPlayerCostFloat64(newItem.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(newComponents[newItem])))
			} else {
				newCost = newCost.add(getPlayerCostFloat64(newItem.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(newComponents[newItem]) * rules.StarbaseComponentCostReduction))
			}
		}

		// apply first part of costs to tally (70% of new item cost - 70% of old item cost)
		// this is the part that can be reduced by normal rebates
		cost = cost.add(newCost.subtract(oldCost).multiply(0.7))

		// add on rest of the cost after category specific rebates
		// higher of (20% new item cost, 30% new item cost - 10% old item cost)
		// if no old item exists, you pay 100%
		adjCost := newCost.multiply(0.2).max(
			newCost.multiply(0.3).subtract(oldCost.multiply(0.1)))
		cost = cost.add(adjCost)
		minCost = minCost.add(adjCost)
	}

	return cost.max(minCost).multiply(raceSpec.StarbaseCostFactor).toCost(math.Ceil).MinZero(), nil
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

	cost := getPlayerCostFloat64(hull.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset)

	// iterate through slots and tally prices up
	for _, slot := range design.Slots {
		item := rules.techs.GetHullComponent(slot.HullComponent)
		if item == nil {
			return Cost{}, fmt.Errorf("component %s in design slots not found in tech store", slot.HullComponent)
		}
		hcCost := getPlayerCostFloat64(item.Tech, techLevels, raceSpec.MiniaturizationSpec, raceSpec.TechCostOffset).multiply(float64(slot.Quantity))
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
