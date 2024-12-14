package cs

import (
	"math"

	"slices"
)

// Check for tech level increases
type techTrader interface {
	// Return the tech field gained for this tech trade event, if any
	techLevelGained(rules *Rules, current, target TechLevel) TechField

	// Return which tradable part was gained for this tech trade event (if any)
	//
	// Loops through the tokens' slots to count MT parts
	acquirablePartGained(rules *Rules, player *Player, tokens []ShipToken) *Tech
}

type techTrade struct {
}

func newTechTrader() techTrader {
	return &techTrade{}
}

// check for a tech level bonus from a player tech level and some target we scrapped, invaded, etc
// https://wiki.starsautohost.org/wiki/Guts_of_Tech_Trading
func (t *techTrade) techLevelGained(rules *Rules, current, target TechLevel) TechField {
	diff := target.Subtract(current).MinZero()
	if diff.Sum() <= 0 {
		return TechFieldNone
	}

	for _, field := range TechFields {
		level := diff.Get(field)
		if level > 0 {
			chance := techTradeChance(rules.TechTradeChance, level)
			// check if our chance is higher than a randomly rolled float;
			// ex: 0.375 >= R for 2 levels above
			if chance >= rules.random.Float64() {
				return field
			}
		}
	}

	return TechFieldNone
}

// return which tradable component(s) were obtained for this tech trading instance, if any
func (t *techTrade) acquirablePartGained(rules *Rules, player *Player, tokens []ShipToken) *Tech {
	if player == nil || tokens == nil || player.acquirablePartGained {
		return nil
	}

	qtyPerPart := map[*Tech]int{} // maps tech part to total quantity on fleet
	parts := []*Tech{}

	// tally up parts in our fleet
	for _, token := range tokens {
		// iterate through the token's slots and remove anything not explicitly tradeable
		slotsFiltered := slices.DeleteFunc(slices.Clone(token.design.Slots), func(slot ShipDesignSlot) bool {
			tc := rules.techs.GetHullComponent(slot.HullComponent)
			return tc == nil ||
				!tc.Requirements.Acquirable ||
				player.AcquiredTechs[tc.Name]
		})

		for _, slot := range slotsFiltered {
			tech := rules.techs.GetHullComponent(slot.HullComponent).Tech
			qtyPerPart[&tech] += slot.Quantity * token.Quantity
			parts = append(parts, &tech)
		}

		hull := rules.techs.GetHull(token.design.Hull)
		if hull != nil &&
			hull.Requirements.Acquirable &&
			!player.AcquiredTechs[hull.Name] {
			qtyPerPart[&hull.Tech] += token.Quantity
			parts = append(parts, &hull.Tech)
		}
	}

	if len(qtyPerPart) > 0 {
		// randomize the order of items to remove potential bias
		rules.random.Shuffle(len(parts), func(i, j int) { parts[i], parts[j] = parts[j], parts[i] })

		// loop through our part list
		for _, part := range parts {
			qty := qtyPerPart[part]
			if checkAcquirablePartChance(rules, qty) {
				return part
			}
		}
	}

	return nil
}

// get the chance of a tech trade. If we are one level above this is:
// .5 * (1 - .5) = .25
// if we are two levels above this is:
// .5 * (1 - .5*.5) = .375
func techTradeChance(baseChance float64, level int) float64 {
	return baseChance * (1 - math.Pow(baseChance, float64(level)))
}

// Check if a successful part trade occurs based on base trade chance and item quantity
//
// Allocates based on optimal distribution of items (ie as many in 1 check as you can fit)
func checkAcquirablePartChance(rules *Rules, qty int) bool {
	for check := 0; check < qty; check += rules.AcquirablePartTradeItemMax {
		// chance for 1 check = # of items (max 25) * 0.005
		tradeChance := rules.AcquirablePartTradeChanceBase * float64(MinInt(
			qty-check, rules.AcquirablePartTradeItemMax))
		if tradeChance >= rules.random.Float64() {
			// we traded the part!
			return true
		}
	}
	return false
}
