package cs

import (
	"math"
	"sort"
)

type battleWeaponType int

const (
	battleWeaponTypeBeam battleWeaponType = iota
	battleWeaponTypeTorpedo
)

// A token firing weapons
type battleWeaponSlot struct {
	// The token with the weapon
	token *battleToken

	// The weapon slot
	slot ShipDesignSlot

	// how many of this weapon are in this slot
	slotQuantity int

	// The type of weapon this weapon slot is
	weaponType battleWeaponType

	// The range of this weapon
	weaponRange int
	// the hull's range bonus included in weaponRange
	rangeBonus int

	// the power of the weapon
	power     int
	beamBonus float64

	// true if this weapon damages shields only (i.e. a sapper)
	damagesShieldsOnly bool

	// the accuracy of the weapon, if it's a torpedo
	accuracy float64

	// the accuracy bonus to the torpedo
	torpedoBonus float64

	// the initiative of the weapon
	initiative int

	// gatling guns hit all targets in range
	hitsAllTargets bool

	// capital ships missiles do double damage after shields are gone
	capitalShipMissile bool
}

type battleWeaponDamage struct {
	// the damage inflicted on shields
	shieldDamage int
	// the damage inflicted on armor
	armorDamage int
	// the new stack damage, per ship in the stack
	damage float64
	// the new stack quantity damaged
	quantityDamaged int
	// the number of tokens destroyed
	numDestroyed int
	// any leftover beam power or torpedoes we have after destroying all ships in the stack
	leftover int
}

// newBattleWeaponSlot creates a new BattleWeaponSlot object
func newBattleWeaponSlot(token *battleToken, slot ShipDesignSlot, hc *TechHullComponent, rangeBonus int, torpedoBonus float64, beamBonus float64) *battleWeaponSlot {
	weaponSlot := &battleWeaponSlot{
		token:              token,
		slot:               slot,
		slotQuantity:       slot.Quantity,
		weaponRange:        hc.Range + rangeBonus,
		rangeBonus:         rangeBonus,
		power:              hc.Power,
		beamBonus:          beamBonus,
		damagesShieldsOnly: hc.DamageShieldsOnly,
		accuracy:           float64(hc.Accuracy) / 100.0, // accuracy as 0 to 1.0
		torpedoBonus:       torpedoBonus,
		initiative:         min(63, token.Initiative+hc.Initiative),
		hitsAllTargets:     hc.HitsAllTargets,
		capitalShipMissile: hc.CapitalShipMissile,
	}

	switch hc.Category {
	case TechCategoryBeamWeapon:
		weaponSlot.weaponType = battleWeaponTypeBeam
	case TechCategoryTorpedo:
		weaponSlot.weaponType = battleWeaponTypeTorpedo
	}

	return weaponSlot
}

// get beam damage with dropoff and defense included
func getBeamDamageAtDistance(damage, weaponRange, dist int, beamDefense float64, beamRangeDropoff float64) int {
	if weaponRange > 0 {
		return int(math.Round(float64(damage) * (1 - float64(dist)/float64(weaponRange)*beamRangeDropoff) * beamDefense))
	}
	return int(math.Round(float64(damage) * beamDefense))
}

// Return true if this weapon slot wiil damage this token
// if this is a sapper and the target is out of shields, return false
func (slot *battleWeaponSlot) willDamage(target *battleToken) bool {
	if target == nil {
		return false
	}
	return !slot.damagesShieldsOnly || (slot.damagesShieldsOnly && target.stackShields > 0)
}

// Return true if this weapon slot is in range of the token target
func (slot *battleWeaponSlot) isInRange(target *battleToken) bool {
	if target == nil {
		return false
	}
	return slot.isInRangePosition(target.Position)
}

func (slot *battleWeaponSlot) isInRangePosition(position Vector) bool {
	// diagonal shots count as one move, so we take the max distance on the x or y as our actual distance away
	// i.e. 4 over, 1 up is 4 range away, 3 over 2 up is 3 range away, etc.
	return slot.token.getDistanceAway(position) <= slot.weaponRange
}

// get the attractiveness of a token versus a weapon
func (weapon *battleWeaponSlot) getAttractiveness(target *battleToken) float64 {
	quantity := max(1, target.Quantity)
	armor := max(1, float64(target.armor*quantity)-target.Damage*float64(target.QuantityDamaged))
	shields := float64(target.stackShields)
	value := min(100000, float64((target.cost.Boranium+target.cost.Resources)*quantity))
	if weapon.weaponType == battleWeaponTypeBeam {
		if weapon.damagesShieldsOnly {
			if shields == 0 {
				return 0
			}
			return value * target.beamDamageMultiplier() / shields
		}
		return value * target.beamDamageMultiplier() / (armor + shields + 1)
	}
	accuracy := weapon.getAccuracy(target.torpedoJamming)
	if accuracy == 0 {
		return 0
	}
	armorShots := armor * 2 / accuracy
	shieldShots := shields / (accuracy/2 + (1-accuracy)/8)
	multiplier := 1.0
	if weapon.capitalShipMissile {
		multiplier = 2
	}
	remainingShots := (armor - shieldShots*accuracy/2) / (accuracy * multiplier)
	defense := min(armorShots, shieldShots+remainingShots)
	if defense <= 0 {
		return 0
	}
	return value / defense
}

// findTargets returns the hostile tokens in range that this weapon will fire at,
// primary targets first, each group sorted by attractiveness.
func (weapon *battleWeaponSlot) findTargets(tokens []*battleToken) (targets []*battleToken) {
	attacker := weapon.token
	primaryTarget := attacker.PrimaryTarget
	secondaryTarget := attacker.SecondaryTarget

	var primaryTargets []*battleToken
	var secondaryTargets []*battleToken

	// Find all enemy tokens
	for _, token := range tokens {
		if !token.isStillInBattle() || !attacker.willAttack(token.PlayerNum) || !weapon.isInRange(token) {
			continue
		}

		// if we will target this
		if token.isTargetOf(primaryTarget) && weapon.willDamage(token) {
			primaryTargets = append(primaryTargets, token)
		} else if token.isTargetOf(secondaryTarget) && weapon.willDamage(token) {
			secondaryTargets = append(secondaryTargets, token)
		}
	}

	sort.SliceStable(primaryTargets, func(i, j int) bool {
		return weapon.getAttractiveness(primaryTargets[i]) > weapon.getAttractiveness(primaryTargets[j])
	})
	sort.SliceStable(secondaryTargets, func(i, j int) bool {
		return weapon.getAttractiveness(secondaryTargets[i]) > weapon.getAttractiveness(secondaryTargets[j])
	})

	targets = make([]*battleToken, 0, len(primaryTargets)+len(secondaryTargets))
	targets = append(targets, primaryTargets...)
	targets = append(targets, secondaryTargets...)
	return targets
}

// getBeamDamageToTarget returns how much of a beam volley, after range dropoff and
// deflectors, goes to the target's shields and armor, and how much is left over.
func (weapon *battleWeaponSlot) getBeamDamageToTarget(damage int, target *battleToken, beamRangeDropoff float64) battleWeaponDamage {
	dist := weapon.token.getDistanceAway(target.Position)
	if weapon.hitsAllTargets {
		dist = 0
	}
	// attenuation uses the weapon's base range without the hull bonus
	damage = getBeamDamageAtDistance(damage, weapon.weaponRange-weapon.rangeBonus, dist, target.beamDamageMultiplier(), beamRangeDropoff)
	shieldDamage := min(target.stackShields, damage)
	if weapon.damagesShieldsOnly || damage <= target.stackShields {
		leftover := 0
		if weapon.damagesShieldsOnly {
			leftover = damage - shieldDamage
		}
		return battleWeaponDamage{shieldDamage: shieldDamage, damage: target.Damage, quantityDamaged: target.QuantityDamaged, leftover: leftover}
	}
	result := getBattleArmorDamage(target, float64(damage-shieldDamage), target.Quantity)
	result.shieldDamage = shieldDamage
	if target.attributes&battleTokenAttributeStarbase != 0 {
		result.leftover = 0
	}
	return result
}

// get the accuracy of a torpedo against a target
func (weapon *battleWeaponSlot) getAccuracy(torpedoJamming float64) float64 {
	if weapon.accuracy == 0 {
		return 0
	}
	accuracy := weapon.accuracy
	if torpedoJamming >= weapon.torpedoBonus {
		accuracy *= 1 - (torpedoJamming - weapon.torpedoBonus)
	} else {
		accuracy += (1 - accuracy) * (weapon.torpedoBonus - torpedoJamming)
	}
	return Clamp(accuracy, 0.01, 1.0)
}

// battleArmorDamageSteps is the precision of a ship's stored armor damage.
const battleArmorDamageSteps = 500

// roundBattleArmorDamage rounds a ship's damage up to the next 1/500th of its armor,
// doing at least one point of damage. Every hit that spreads armor damage across a
// stack therefore damages each surviving ship by at least 0.2%.
func roundBattleArmorDamage(damage, armor float64) float64 {
	step := armor / battleArmorDamageSteps
	steps := math.Ceil(max(1, damage)/step - 1e-9)
	return Clamp(steps, 1, battleArmorDamageSteps-1) * armor / battleArmorDamageSteps
}

// getBattleArmorDamage spends new damage against damaged ships first, then spreads
// the remainder across survivors. A volley can destroy at most killLimit ships.
func getBattleArmorDamage(target *battleToken, damage float64, killLimit int) battleWeaponDamage {
	result := battleWeaponDamage{damage: target.Damage, quantityDamaged: target.QuantityDamaged}
	if damage <= 0 || target.Quantity == 0 {
		return result
	}
	armor := float64(target.armor)
	damaged := min(target.Quantity, target.QuantityDamaged)
	remaining := damage
	damagedArmor := max(math.SmallestNonzeroFloat64, armor-target.Damage)
	damagedKills := min(damaged, killLimit, int(remaining/damagedArmor))
	remaining -= float64(damagedKills) * damagedArmor
	damaged -= damagedKills
	result.numDestroyed = damagedKills
	healthyKills := min(max(0, target.Quantity-target.QuantityDamaged), killLimit-damagedKills, int(remaining/max(1, armor)))
	remaining -= float64(healthyKills) * armor
	result.numDestroyed += healthyKills
	quantity := target.Quantity - result.numDestroyed
	result.armorDamage = int(math.Round(min(damage, float64(target.armor*target.Quantity)-target.Damage*float64(target.QuantityDamaged))))
	if quantity == 0 {
		result.damage = 0
		result.quantityDamaged = 0
		result.leftover = int(math.Round(remaining))
	} else if result.numDestroyed >= killLimit {
		result.armorDamage = int(math.Round(damage - remaining))
		result.damage = target.Damage
		result.quantityDamaged = damaged
		if damaged == 0 {
			result.damage = 0
		}
	} else if remaining > 0 {
		result.damage = roundBattleArmorDamage((remaining+target.Damage*float64(damaged))/float64(quantity), armor)
		result.quantityDamaged = quantity
	} else {
		result.quantityDamaged = damaged
		if damaged == 0 {
			result.damage = 0
		}
	}
	return result
}

// getTorpedoVolleyDamage calculates pooled splash and hit damage without changing
// the target, limiting ship kills to the number of torpedoes fired.
func (weapon *battleWeaponSlot) getTorpedoVolleyDamage(target *battleToken, hits, misses float64, count int, splash float64) battleWeaponDamage {
	power := float64(weapon.power)
	if weapon.capitalShipMissile && target.stackShields == 0 {
		power *= 2
	}
	splashDamage := min(float64(target.stackShields), misses*power*splash)
	shieldsLeft := float64(target.stackShields) - splashDamage
	hitDamage := hits * power
	hitShieldDamage := min(shieldsLeft, hitDamage/2)
	result := getBattleArmorDamage(target, hitDamage-hitShieldDamage, count)
	result.shieldDamage = int(math.Round(splashDamage + hitShieldDamage))
	result.leftover = 0
	return result
}

// beamPower returns the slot's total volley power after applying the beam bonus.
func (weapon *battleWeaponSlot) beamPower() int {
	return int(math.Round(float64(weapon.power*weapon.slotQuantity*weapon.token.Quantity) * weapon.beamBonus))
}
