package cs

import "math"

// Check for tech level increases
type techTrader interface {
	// Return the tech field gained for this tech trade event, if any
	techLevelGained(rules *Rules, current, target TechLevel) TechField
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
			// get the chance of a tech trade. If we are one level above this is:
			// .5 * (1 - .5) = .25
			// if we are two levels above this is:
			// .5 * (1 - .5*.5) = .375
			chance := rules.TechTradeChance * (1 - math.Pow(rules.TechTradeChance, float64(level)))
			// check if our random number between 0 and 1 is under the above, i.e. < .375 for 2 levels above
			if rules.random.Float64() <= chance {
				return field
			}
		}
	}

	return TechFieldNone
}
