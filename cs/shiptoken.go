package cs

import "math"

type ShipToken struct {
	DesignNum       int     `json:"designNum,omitempty"`
	Quantity        int     `json:"quantity"`
	Damage          float64 `json:"damage"`
	QuantityDamaged int     `json:"quantityDamaged"`
	design          *ShipDesign
}

type tokenDamage struct {
	damage         int
	shipsDestroyed int
}

// Apply damage to a token, updating quantity damaged and damage amount
func (st *ShipToken) applyMineDamage(damage int) tokenDamage {
	// mines do half damage to shields
	shields := st.design.Spec.Shield
	armor := st.design.Spec.Armor
	possibleDamageToShields := float64(damage) * 0.5
	actualDamageToShields := math.Min(float64(shields), possibleDamageToShields)
	remainingDamage := damage - int(actualDamageToShields)
	existingDamage := st.Damage * float64(st.QuantityDamaged)

	st.Damage = float64(existingDamage) + float64(remainingDamage)

	tokensDestroyed := int(math.Min(float64(st.Quantity), math.Floor(float64(st.Damage)/float64(armor))))
	st.Quantity -= tokensDestroyed

	if st.Quantity > 0 {
		// Figure out how much damage we have leftover after destroying
		// tokens. This will be applied to the rest of the tokens
		// if we took 100 damage, and we have 40 armor, we lose 2 tokens
		// and have 20 leftover damage to spread across tokens
		leftoverDamage := st.Damage - float64(tokensDestroyed*armor)
		st.Damage = leftoverDamage
		st.QuantityDamaged = st.Quantity
	}

	return tokenDamage{damage: remainingDamage, shipsDestroyed: tokensDestroyed}
}

// Apply damage (if any) to each token that overgated
func (st *ShipToken) applyOvergateDamage(dist float64, safeRange int, safeSourceMass int, safeDestMass int, maxMassFactor int) tokenDamage {
	rangeDamageFactor := st.getStargateRangeDamageFactor(dist, safeRange)
	massDamageFactor := st.getStargateMassDamageFactor(dist, safeSourceMass, safeDestMass, maxMassFactor)

	totalDamageFactor := math.Min(0.98, massDamageFactor+(1.0-massDamageFactor)*rangeDamageFactor)

	// apply damage as a percentage of armor to all tokens
	armor := st.design.Spec.Armor
	newDamage := int(math.Round(totalDamageFactor * float64(st.Quantity*armor)))
	st.QuantityDamaged = st.Quantity
	st.Damage += float64(newDamage)

	tokensDestroyed := int(math.Min(float64(st.Quantity), math.Floor(float64(st.Damage)/float64(armor))))

	st.Quantity = st.Quantity - tokensDestroyed
	if st.Quantity > 0 {
		// Figure out how much damage we have leftover after destroying
		// tokens. This will be applied to the rest of the tokens
		// if we took 100 damage, and we have 40 armor, we lose 2 tokens
		// and have 20 leftover damage to spread across tokens
		leftoverDamage := st.Damage - float64(tokensDestroyed*armor)
		st.Damage = leftoverDamage
		st.QuantityDamaged = st.Quantity
	}

	return tokenDamage{newDamage, tokensDestroyed}
}

func (t *ShipToken) getStargateRangeDamageFactor(dist float64, safeRange int) float64 {
	rangeDamageFactor := 0.0
	if safeRange == InfinteGate || safeRange >= int(dist) {
		rangeDamageFactor = 0
	} else {
		rangeDamageFactor = (dist - float64(safeRange)) / (4.0 * float64(safeRange))
	}

	return rangeDamageFactor
}

func (t *ShipToken) getStargateMassDamageFactor(dist float64, safeSourceMass int, safeDestMass int, maxMassFactor int) float64 {
	mass := t.design.Spec.Mass
	sourceMassDamageFactor := 1.0
	destMassDamageFactor := 1.0
	if safeSourceMass != InfinteGate && safeSourceMass < mass {
		sourceMassDamageFactor = (float64(maxMassFactor)*float64(safeSourceMass) - float64(mass)) / (4.0 * float64(safeSourceMass))
	}
	if safeDestMass != InfinteGate && safeDestMass < mass {
		destMassDamageFactor *= (float64(maxMassFactor)*float64(safeDestMass) - float64(mass)) / (4.0 * float64(safeDestMass))
	}

	return 1 - (sourceMassDamageFactor * destMassDamageFactor)
}

// Vanishing% = 100/3*[1-(5*maxMass-mass)^2/(4*maxMass)^2], rounded down to nearest 1%.
// where maxMass is the maximum safe mass for the sending gate.
func (t *ShipToken) getStargateMassVanishingChance(safeSourceMass int, maxMassFactor int) float64 {
	mass := t.design.Spec.Mass
	vanishingChance := 100.0 / 3 * (1 -
		math.Pow((float64)(maxMassFactor*safeSourceMass-mass), 2)/
			(math.Pow((float64)(4*safeSourceMass), 2)))

	// chance as percent
	return math.Floor(vanishingChance) / 100.0
}

// For distance overgating the probability of ships being lost to the void is roughly equal to the damage divided by 3.
// For example, if the overgating causes 60% damage then there will be a 20% chance of losing the ship.
func (t *ShipToken) getStargateRangeVanishingChance(dist float64, safeRange int) float64 {
	// chance as percent, rounded down to 1%
	return math.Floor(100.0/3*t.getStargateRangeDamageFactor(dist, safeRange)) / 100.0
}
