package cs

import "math"

// A Fleet contains multiple ShipTokens, each of which have a design and a quantity.
type ShipToken struct {
	DesignNum       int     `json:"designNum"`
	Quantity        int     `json:"quantity"`                  // number of ships in the token
	Damage          float64 `json:"damage,omitempty"`          // damage per damaged ship in the token
	QuantityDamaged int     `json:"quantityDamaged,omitempty"` // number of damaged ships in token
	design          *ShipDesign
}

type tokenDamage struct {
	damage         int
	shipsDestroyed int
}

// Apply mine damage to a token, updating quantity damaged and damage amount
func (st *ShipToken) applyMineDamage(damage int) tokenDamage {
	// mines do half damage to shields
	shields := st.design.Spec.Shields
	armor := st.design.Spec.Armor
	possibleDamageToShields := float64(damage) * 0.5
	actualDamageToShields := min(float64(shields), possibleDamageToShields)
	armorDamage := damage - int(actualDamageToShields)
	existingStackDamage := st.Damage * float64(st.QuantityDamaged) // get the total stack damage

	// get the new stackDamage spread across all ships int he stack
	stackDamage := math.Floor(float64(existingStackDamage) + float64(armorDamage))

	// from the new total stack damage, figure out how many ships were destroyed
	shipsDestroyed := int(min(float64(st.Quantity), math.Floor(float64(stackDamage)/float64(armor))))
	st.Quantity -= shipsDestroyed

	if st.Quantity > 0 {
		// Figure out how much damage we have leftover after destroying
		// ships. This will be applied to the rest of the ships
		// if we took 100 damage, and we have 40 armor, we lose 2 ships
		// and have 20 leftover damage to spread across ships
		leftoverStackDamage := stackDamage - float64(shipsDestroyed*armor)
		st.Damage = math.Floor(leftoverStackDamage / float64(st.Quantity))
		st.QuantityDamaged = st.Quantity
	}

	return tokenDamage{damage: armorDamage, shipsDestroyed: shipsDestroyed}
}

// Apply overgate damage (if any) to each token that overgated
func (st *ShipToken) applyOvergateDamage(dist float64, safeRange int, safeSourceMass int, safeDestMass int, maxMassFactor int) tokenDamage {
	if st.Quantity == 0 {
		// no ships means nothing to damage
		return tokenDamage{damage: 0, shipsDestroyed: 0}
	}

	// testing with overgating 12 scouts 479.5 ly with a 250ly gate
	// 1st run damaged all 12 by 20% (4 damage each);
	// subsequent runs went 40%, 60%, 80%, then 100% (destroying all ships).

	rangeDamageFactor := st.getStargateRangeDamageFactor(dist, safeRange)
	massDamageFactor := st.getStargateMassDamageFactor(safeSourceMass, safeDestMass, maxMassFactor)

	// damage capped at 98% for a single overgate
	totalDamageFactor := min(0.98, massDamageFactor+(1-massDamageFactor)*rangeDamageFactor)

	// apply damage as a percentage of armor to all tokens
	armor := st.design.Spec.Armor
	existingDamage := st.Damage

	var tokensDestroyed int
	damagePerShip := int(math.Round(totalDamageFactor * float64(armor)))

	// ships are never destroyed by overgating if they aren't already damaged
	if existingDamage == 0 && damagePerShip >= armor {
		damagePerShip = armor - 1
	}

	st.Damage += float64(damagePerShip)

	if st.Damage >= float64(armor) {
		// our damage exceeds our armor, destroy any previous damaged ships
		tokensDestroyed = st.QuantityDamaged
		st.Quantity -= st.QuantityDamaged
	}

	// apply overgate damage to any leftover tokens
	if damagePerShip > 0 {
		st.Damage = float64(damagePerShip)
		st.QuantityDamaged = st.Quantity

		if st.Quantity == 0 {
			// can't damage something that isn't there
			st.Damage = 0
		}
	}

	return tokenDamage{damagePerShip, tokensDestroyed}
}

func (t *ShipToken) getStargateRangeDamageFactor(dist float64, safeRange int) (rangeDamageFactor float64) {
	if safeRange == InfiniteGate || safeRange >= int(dist) {
		return 0
	}

	// Formula: (dist-safeRange)/(4*safeRange)
	return (dist - float64(safeRange)) / (4.0 * float64(safeRange))
}

func (t *ShipToken) getStargateMassDamageFactor(safeSourceMass int, safeDestMass int, maxMassFactor int) float64 {
	mass := t.design.Spec.Mass
	sourceMassDamageFactor := 1.0
	destMassDamageFactor := 1.0
	if safeSourceMass != InfiniteGate && safeSourceMass < mass {
		sourceMassDamageFactor = (float64(maxMassFactor)*float64(safeSourceMass) - float64(mass)) / (4.0 * float64(safeSourceMass))
	}
	if safeDestMass != InfiniteGate && safeDestMass < mass {
		destMassDamageFactor *= (float64(maxMassFactor)*float64(safeDestMass) - float64(mass)) / (4.0 * float64(safeDestMass))
	}

	return 1 - (sourceMassDamageFactor * destMassDamageFactor)
}

// applyOvergateVanishing vanishes overgating ship tokens exceeding safe limits,
// reducing token quanitity as appropriate.
// It returns the total number of tokens vanished (origQty - newQty).
func (token *ShipToken) applyOvergateVanishing(rules *Rules, distance float64, sourceRange, sourceMass int) (shipsLost int) {
	rangeVanishChance := max(0, token.getOvergateRangeVanishingChance(distance, sourceRange))
	massVanishChance := max(0, token.getOvergateMassVanishingChance(sourceMass, rules.StargateMaxHullMassFactor))
	if rangeVanishChance == 0 && massVanishChance == 0 {
		// neither range nor mass can harm us; return
		return
	}

	// Combined vanishing chance formula courtesy of ekolis
	// Both checks fire independently, so the chance of both passing is
	// 1-(rangeFailChance*massFailChance)
	vanishingChance := 1 - (1-rangeVanishChance)*(1-massVanishChance)

	// check each token one by one to see if it kersplodes
	for range token.Quantity {
		if vanishingChance >= rules.random.Float64() {
			shipsLost++
		}
	}

	// reduce token quantity by however many ships died,
	// prioritizing damaged ones if possible.
	token.Quantity -= shipsLost
	token.QuantityDamaged -= shipsLost
	if token.QuantityDamaged <= 0 {
		// reset token damage to 0 if none remain
		token.Damage = 0
		token.QuantityDamaged = 0
	}

	return shipsLost
}

// getOvergateMassVanishingChance returns the mass-based portion of this ShipToken's
// overgate vanishing chance.
// Graph: https://www.desmos.com/calculator/ftqvsbkmj5
func (t *ShipToken) getOvergateMassVanishingChance(safeSourceMass int, maxMassFactor int) (massChance float64) {
	if safeSourceMass == InfiniteGate {
		return 0
	}
	// Mass Vanishing % = 100/3*[1-(5*maxMass-mass)^2/(4*maxMass)^2], rounded down to nearest 1%.
	// where maxMass is the maximum safe mass for the sending gate.
	vanishingChance := 100.0 / 3 * (1 -
		float64(PowInt(maxMassFactor*safeSourceMass-t.design.Spec.Mass, 2))/
			float64(PowInt(4*safeSourceMass, 2)))

	// return chance rounded down to nearest %
	return math.Floor(vanishingChance) / 100
}

// getOvergateRangeVanishingChance returns the range-based portion of this ShipToken's
// overgate vanishing chance.
func (t *ShipToken) getOvergateRangeVanishingChance(dist float64, safeRange int) (rangeChance float64) {
	// Range vanishing chance is roughly equal to 1/3 damage dealt -
	// 60% range damage factor = 20% loss chance.
	chance := 100 * t.getStargateRangeDamageFactor(dist, safeRange) / 3

	// return chance rounded down to nearest %
	return math.Floor(chance) / 100
}
