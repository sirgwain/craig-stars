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
	// Does not take a part as an argument as the only way MT parts can be transferred
	// is if a ship with one dies in battle or is scrapped
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
	diff := target.Minus(current).MinZero()
	if diff.Sum() <= 0 {
		return TechFieldNone
	}

	for _, field := range TechFields {
		level := diff.Get(field)
		if level > 0 {
			chance := techTradeChance(rules.TechTradeChance, level)
			// check if our random number between 0 and 1 is under the above, i.e. < .375 for 2 levels above
			if rules.random.Float64() <= chance {
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

	var qtyPerPart map[*Tech]int
	var parts []*Tech

	// tally up parts in our fleet
	for _, token := range tokens {
		// make copy of token design's slots for part checking
		slots := slices.Clone(token.design.Slots)
		// iterate through the token's slots and remove anything not explicitly tradable 
		slots = slices.DeleteFunc(slots, func(slot ShipDesignSlot) bool {
			return rules.techs.GetHullComponent(slot.HullComponent).Tech.Requirements.Acquirable
		})

		for _, slot := range slots {
			if slot.Quantity == 0 { // indicates field has been zeroed out by slices.DeleteFunc
				continue
			}
			tech := rules.techs.GetHullComponent(slot.HullComponent).Tech
			qtyPerPart[&tech] += slot.Quantity * token.Quantity
			parts = append(parts, &tech)
		}

		hull := rules.techs.GetHull(token.design.Hull).Tech
		if hull.Origin == OriginMysteryTrader {
			qtyPerPart[&hull] += token.Quantity
			parts = append(parts, &hull)
		}
	}

	if len(qtyPerPart) > 0 {
		// randmize the order of items to remove potential bias
		rules.random.Shuffle(len(parts), func(i, j int) { parts[i], parts[j] = parts[j], parts[i] })

		// loop through our part list  
		for _, part := range parts {
			qty := qtyPerPart[part]
			if checkAcquirablePartChance(rules, qty) &&
				!player.AcquiredTechs[part.Name] {
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

// get the chance of a part trade, based on the base chance and the TOTAL number of items in the fleet
// 
// allocates based on optimal distribution of items (ie as many in 1 check as you can)
func checkAcquirablePartChance(rules *Rules, qty int) bool {
	for check := 0; check < qty; check += rules.AcquirablePartTradeItemMax {
		// chance for 1 check = # of items (max 25) * 0.005
		tradeChance := rules.TechTradeChance * rules.AcquirablePartTradeChanceBase * math.Min(
		float64((qty-check)/rules.AcquirablePartTradeItemMax), 1)
		if tradeChance > rules.random.Float64() {
			// we traded the part! 
			return true
		}
	}
	return false
}
