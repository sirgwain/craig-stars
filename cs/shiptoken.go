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
	if st.Quantity == 0 {
		return tokenDamage{}
	}

	// shields absorb up to half the damage
	shields := st.design.Spec.Shields * st.Quantity
	armorDamage := damage - min(shields, damage/2)

	// mines spread their damage evenly over the stack, on top of any damage it already has.
	// If that's more than a ship's armor, the whole stack is destroyed
	stackDamage := float64(armorDamage) + st.Damage*float64(st.QuantityDamaged)
	damagePerShip := stackDamage / float64(st.Quantity)
	if damagePerShip > float64(st.design.Spec.Armor) {
		shipsDestroyed := st.Quantity
		st.Quantity = 0
		st.Damage = 0
		st.QuantityDamaged = 0
		return tokenDamage{damage: armorDamage, shipsDestroyed: shipsDestroyed}
	}

	st.Damage = math.Floor(damagePerShip)
	st.QuantityDamaged = st.Quantity
	return tokenDamage{damage: armorDamage}
}

// overgateDamagePercent returns the percent of its armor a ship loses gating dist ly from a source to a dest
// gate. Past a gate's safe range or mass, the chance of arriving unharmed drops from 100% at the safe value to
// 0% at the gate's max (rules.StargateMaxRangeFactor or StargateMaxHullMassFactor times the safe value). Range
// only counts for the source gate, mass counts for both.
func (st *ShipToken) overgateDamagePercent(rules *Rules, dist float64, source, dest PlanetStarbaseSpec) int {
	survival := 1.0
	if source.SafeRange != InfiniteGate && dist > float64(source.SafeRange) {
		survival *= overgateSurvival(dist, source.SafeRange, rules.StargateMaxRangeFactor)
	}
	mass := float64(st.design.Spec.Mass)
	for _, safeMass := range []int{source.SafeHullMass, dest.SafeHullMass} {
		if safeMass != InfiniteGate && safeMass > 0 && mass > float64(safeMass) {
			survival *= overgateSurvival(mass, safeMass, rules.StargateMaxHullMassFactor)
		}
	}
	return int(math.Floor(100 * (1 - max(0, survival))))
}

// overgateSurvival is the chance of surviving going value past a gate's safe value, from 1 at the safe value
// to 0 at maxFactor times it
func overgateSurvival(value float64, safe, maxFactor int) float64 {
	return (float64(maxFactor*safe) - value) / float64((maxFactor-1)*safe)
}

// applyOvergateDamage damages every ship in the token by damagePercent of its armor, on top of any damage it
// already has. Damaged ships that can't take it are destroyed, but a single overgate never destroys an
// undamaged ship unless it's 100%.
func (st *ShipToken) applyOvergateDamage(damagePercent int) tokenDamage {
	if st.Quantity == 0 || damagePercent <= 0 {
		return tokenDamage{}
	}

	armor := st.design.Spec.Armor
	if damagePercent >= 100 {
		destroyed := st.Quantity
		st.Quantity, st.QuantityDamaged, st.Damage = 0, 0, 0
		return tokenDamage{damage: armor, shipsDestroyed: destroyed}
	}

	damagePerShip := min(armor-1, max(1, armor*damagePercent/100))

	// damaged ships that can't take more damage are destroyed
	destroyed := 0
	if st.QuantityDamaged > 0 && int(st.Damage)+damagePerShip >= armor {
		destroyed = st.QuantityDamaged
		st.Quantity -= destroyed
		st.QuantityDamaged, st.Damage = 0, 0
	}

	// spread the new damage and any old damage over the ships left
	if st.Quantity > 0 {
		stackDamage := float64(damagePerShip*st.Quantity) + st.Damage*float64(st.QuantityDamaged)
		st.Damage = math.Floor(stackDamage / float64(st.Quantity))
		st.QuantityDamaged = st.Quantity
	}

	return tokenDamage{damage: damagePerShip, shipsDestroyed: destroyed}
}

// applyOvergateVanishing loses ships to the void when overgating. Each ship has a chance of a third of the
// overgate damage percent of vanishing. Returns how many ships vanished.
func (token *ShipToken) applyOvergateVanishing(rules *Rules, damagePercent int) (shipsLost int) {
	vanishingChance := float64(damagePercent/3) / 100
	if vanishingChance <= 0 {
		return 0
	}

	// check each ship one by one to see if it kersplodes
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
