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
func (st *ShipToken) applyOvergateDamage(dist float64, safeRange, safeSourceMass, safeDestMass, maxMassFactor, maxRangeFactor int) (damage tokenDamage) {
	if st.Quantity == 0 {
		// no ships means nothing to damage
		return tokenDamage{}
	}

	// testing with overgating 12 scouts 479.5 ly with a 250ly gate
	// 1st run damaged all 12 by 20% (4 damage each);
	// subsequent runs went 40%, 60%, 80%, then 100% (destroying all ships).

	rangeDamageFactor := st.getStargateRangeDamageFactor(dist, safeRange, maxRangeFactor)
	massDamageFactor := st.getStargateMassDamageFactor(safeSourceMass, safeDestMass, maxMassFactor)

	// range and mass damage stack multiplicatively, capping at 98% per instance
	totalDamageFactor := min(0.98, 1-(1-massDamageFactor)*(1-rangeDamageFactor))

	// apply damage as a percentage of armor to all tokens
	armor := st.design.Spec.Armor
	existingDamage := st.Damage

	damagePerShip := int(math.Round(totalDamageFactor * float64(armor)))

	if damagePerShip <= 0 {
		return tokenDamage{}
	}

	// ships are never destroyed by overgating if they aren't already damaged
	if existingDamage == 0 && damagePerShip >= armor {
		damagePerShip = armor - 1
	}

	st.Damage += float64(damagePerShip)

	var tokensDestroyed int
	if st.Damage >= float64(armor) {
		// our damage exceeds our armor, destroy all previously damaged tokens.
		// We don't need to update damagePerShip for the return since damage is ignored
		// in messages involving dead ships
		tokensDestroyed = st.QuantityDamaged
		st.Quantity -= st.QuantityDamaged
		st.Damage = float64(damagePerShip)
	}

	// update token damaged quantity, resetting damage to 0 if entirely wiped out
	st.QuantityDamaged = st.Quantity
	if st.Quantity == 0 {
		st.Damage = 0
	}

	return tokenDamage{damagePerShip * st.Quantity, tokensDestroyed}
}

func (t *ShipToken) getStargateRangeDamageFactor(dist float64, safeRange, maxSafeRange int) (rangeDamageFactor float64) {
	if safeRange == InfiniteGate || safeRange >= int(dist) {
		return 0
	}

	// Formula: (dist-safeRange)/(4*safeRange)
	return (dist - float64(safeRange)) / float64((maxSafeRange-1)*safeRange)
}

func (t *ShipToken) getStargateMassDamageFactor(safeSourceMass, safeDestMass, maxSafeMass int) float64 {
	mass := t.design.Spec.Mass
	sourceMassDamageFactor := 1.0
	destMassDamageFactor := 1.0
	// dmg% = 1 - (5*safeMass-shipMass) / (4*safeMass)
	// This occurs invididually for both source and dest gates
	// reaching 100% once either source or dest gates' capacities are exceeded by 5x.
	if safeSourceMass < mass {
		sourceMassDamageFactor = float64(maxSafeMass*safeSourceMass-mass) / float64((maxSafeMass-1)*safeSourceMass)
	}
	if safeDestMass < mass {
		destMassDamageFactor = float64(maxSafeMass*safeDestMass-mass) / float64((maxSafeMass-1)*safeDestMass)
	}

	return 1 - sourceMassDamageFactor*destMassDamageFactor
}

// applyOvergateVanishing vanishes overgating ship tokens exceeding safe limits,
// reducing token quanitity as appropriate.
// It returns the total number of tokens vanished (origQty - newQty).
func (token *ShipToken) applyOvergateVanishing(rules *Rules, distance float64, sourceRange, sourceMass int) (shipsLost int) {
	rangeVanishChance := token.getOvergateRangeVanishingChance(distance, sourceRange, rules.StargateMaxRangeFactor)
	massVanishChance := token.getOvergateMassVanishingChance(sourceMass, rules.StargateMaxHullMassFactor)
	if rangeVanishChance <= 0 && massVanishChance <= 0 {
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
func (t *ShipToken) getOvergateMassVanishingChance(safeSourceMass, maxMassFactor int) (massChance float64) {
	if safeSourceMass == InfiniteGate {
		return 0
	}
	// Mass Vanishing % = 100/3*[1-(5*maxMass-mass)^2/(4*maxMass)^2], rounded down to nearest 1%.
	// where maxMass is the maximum safe mass for the sending gate.
	// This caps out at 33% at 5x max mass
	vanishingChance := 100.0 / 3 * (1 -
		float64(PowInt(maxMassFactor*safeSourceMass-t.design.Spec.Mass, 2))/
			float64(PowInt(4*safeSourceMass, 2)))

	// return chance rounded down to nearest %
	return math.Floor(vanishingChance) / 100
}

// getOvergateRangeVanishingChance returns the range-based portion of this ShipToken's
// overgate vanishing chance.
func (t *ShipToken) getOvergateRangeVanishingChance(dist float64, safeRange, maxSafeRange int) (rangeChance float64) {
	// Range vanishing chance is roughly equal to 1/3 damage dealt -
	// 60% range damage factor = 20% loss chance.
	chance := 100 * t.getStargateRangeDamageFactor(dist, safeRange, maxSafeRange) / 3

	// return chance rounded down to nearest %
	return math.Floor(chance) / 100
}
