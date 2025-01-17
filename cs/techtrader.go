package cs

import "math"

// The techTrader interface handles checks for tech level increases from trading
type techTrader interface {
	// Perform an invasion tech trading check.
	//
	// Wrapper function for techLevelGained
	checkInvasionTechTrade(rules *Rules, player *Player, targetLevel TechLevel) (field TechField)

	// Check fleet-based tech trading for both parts & levels at once.
	//
	checkFleetTechTrade(rules *Rules, player *Player, tokens []ShipToken) (field TechField, acquiredPart *Tech)

	// Check for a tech level boost for a player's tech level and some target we scrapped, invaded, etc.
	//
	// Returns the field gained from this trade instance, if any
	techLevelGained(rules *Rules, current, target TechLevel) TechField

	// Checks for acquirable part gain from from this tech trade event
	//
	// Returns the part gained from this trade instance, if any
	acquirablePartGained(rules *Rules, player *Player, tokens []ShipToken) *Tech
}

type techTrade struct{}

func newTechTrader() techTrader {
	return &techTrade{}
}

// Perform an invasion tech trading check.
//
// Wrapper function for techLevelGained
func (t *techTrade) checkInvasionTechTrade(rules *Rules, player *Player, targetLevel TechLevel) (field TechField) {
	if player.techLevelGained {
		return TechFieldNone
	}

	field = t.techLevelGained(rules, player.TechLevels, targetLevel)
	return field
}

// Check fleet-based tech trading for both parts & levels at once.
//
// Returns the field and part gained from this trade instance, if any
func (t *techTrade) checkFleetTechTrade(rules *Rules, player *Player, tokens []ShipToken) (field TechField, acquiredPart *Tech) {
	if len(tokens) == 0 {
		return TechFieldNone, nil
	}

	field = TechFieldNone
	if !player.techLevelGained {
		for _, token := range tokens {
			tokenLevel := token.design.Spec.TechLevel
			field = t.techLevelGained(rules, player.TechLevels, tokenLevel)
			if field != TechFieldNone {
				break // we have the field; time to skedaddle
			}
		}
	}

	if !player.acquirablePartGained {
		acquiredPart = t.acquirablePartGained(rules, player, tokens)
	}

	return field, acquiredPart
}

// Check for a tech level boost for a player's tech level and some target we scrapped, invaded, etc.
//
// Returns the field gained from this trade instance, if any
func (t *techTrade) techLevelGained(rules *Rules, current, target TechLevel) TechField {
	diff := target.Subtract(current).MinZero()
	if diff.Sum() <= 0 { // targetLevel has no fields greater than playerLevels
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

// Checks for acquirable part gain from from this tech trade event
//
// Returns the part gained from this trade instance, if any
func (t *techTrade) acquirablePartGained(rules *Rules, player *Player, tokens []ShipToken) *Tech {
	if player == nil || tokens == nil || player.acquirablePartGained {
		return nil
	}

	qtyPerPart := map[*Tech]int{} // maps tech part to total quantity on fleet
	parts := []*Tech{}            // list of parts being checked for; allows for deterministic shuffling of part checks

	// tally up parts in our fleet
	for _, token := range tokens {
		for _, slot := range token.design.Slots {
			hc := rules.techs.GetHullComponent(slot.HullComponent)
			if hc == nil || // hull component doesn't exist
				!hc.Requirements.Acquirable || // component is not acquirable
				player.AcquiredTechs[hc.Name] { // we already have this part
				continue
			}

			tech := hc.Tech
			if _, ok := qtyPerPart[&tech]; !ok {
				parts = append(parts, &tech)
			}
			qtyPerPart[&tech] += slot.Quantity * token.Quantity
		}

		hull := rules.techs.GetHull(token.design.Hull)
		if hull != nil &&
			hull.Requirements.Acquirable &&
			!player.AcquiredTechs[hull.Name] {
			if _, ok := qtyPerPart[&hull.Tech]; !ok {
				parts = append(parts, &hull.Tech)
			}
			qtyPerPart[&hull.Tech] += token.Quantity
		}
	}

	if len(qtyPerPart) == 0 {
		return nil // no parts, no checks
	}

	// randomize the order of items to remove potential bias
	rules.random.Shuffle(len(parts), func(i, j int) { parts[i], parts[j] = parts[j], parts[i] })

	// loop through our shuffled part list and check parts one by one
	for _, part := range parts {
		qty := qtyPerPart[part]
		if checkAcquirablePartChance(rules, qty) {
			return part
		}
	}

	return nil
}

// get the chance of a tech level trade occurring based on base chance & levels above target
func techTradeChance(baseChance float64, level int) float64 {
	// If we are one level above this is:
	// .5 * (1 - .5) = .25
	// if we are two levels above this is:
	// .5 * (1 - .5*.5) = .375
	return baseChance * (1 - math.Pow(baseChance, float64(level)))
}

// Check if a successful part trade occurs based on trade chance and item quantity
//
// Allocates based on optimal distribution of items (ie as many in 1 check as you can fit)
func checkAcquirablePartChance(rules *Rules, qty int) bool {
	for check := 0; check < qty; check += rules.AcquirablePartTradeItemMax {
		// chance for 1 check = # of items (max 25) * 0.005
		tradeChance := rules.AcquirablePartTradeChanceBase * float64(Min(
			qty-check, rules.AcquirablePartTradeItemMax))
		if tradeChance >= rules.random.Float64() {
			// we traded the part!
			return true
		}
	}
	return false
}
