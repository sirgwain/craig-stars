package cs

import (
	"math"
	"sort"

	"github.com/rs/zerolog/log"
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

	// all the other tokens this weapon is targeting
	targets []*battleToken

	// The weapon slot
	slot ShipDesignSlot

	// how many of this weapon are in this slot
	slotQuantity int

	// The type of weapon this weapon slot is
	weaponType battleWeaponType

	// The range of this weapon
	weaponRange int

	// the power of the weapon
	power int

	// true if this weapon damages shields only (i.e. a sapper)
	damagesShieldsOnly bool

	// the accuracy of the weapon, if it's a torpedo
	accuracy float64

	// the accuracy bonus to the torpedo
	torpedoInaccuracyMulti float64

	// the initiative of the weapon
	initiative int

	// gattling guns hit all targets in range
	hitsAllTargets bool

	// capital ships missiles do double damage after shields are gone
	capitalShipMissile bool
}

type battleWeaponDamage struct {
	// the damage inflicted on shields
	shieldDamage int
	// the damage inflicted on armor
	armorDamage int
	// the new stack damage
	damage float64
	// the new stack quantity damaged
	quantityDamaged int
	// the number of tokens destroyed
	numDestroyed int
	// any leftover beam power or torpedos we have after destroying all ships in the stack
	leftover int
}

// newBattleWeaponSlot creates a new BattleWeaponSlot object
func newBattleWeaponSlot(token *battleToken, slot ShipDesignSlot, hc *TechHullComponent, rangeBonus int, torpedoInaccuracyMulti float64, beamBonus float64) *battleWeaponSlot {
	weaponSlot := &battleWeaponSlot{
		token:                  token,
		slot:                   slot,
		slotQuantity:           slot.Quantity,
		weaponRange:            hc.Range + rangeBonus,
		power:                  hc.Power,
		damagesShieldsOnly:     hc.DamageShieldsOnly,
		accuracy:               float64(hc.Accuracy) / 100.0, // accuracy as 0 to 1.0
		torpedoInaccuracyMulti: torpedoInaccuracyMulti,
		initiative:             token.Initiative + hc.Initiative,
		hitsAllTargets:         hc.HitsAllTargets,
		capitalShipMissile:     hc.CapitalShipMissile,
	}

	if hc.Category == TechCategoryBeamWeapon {
		weaponSlot.weaponType = battleWeaponTypeBeam
		weaponSlot.power = int(float64(weaponSlot.power) * (1 + beamBonus))
	} else if hc.Category == TechCategoryTorpedo {
		weaponSlot.weaponType = battleWeaponTypeTorpedo
	}

	return weaponSlot
}

// get beam damage with dropoff and defense included
func getBeamDamageAtDistance(damage, weaponRange, dist int, beamDefense float64, beamRangeDropoff float64) int {
	if weaponRange > 0 {
		return int(math.Round(float64(damage) * (1 - float64(dist)/float64(weaponRange)*beamRangeDropoff) * (1 - beamDefense)))
	}
	return int(math.Round(float64(damage) * (1 - beamDefense)))
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

func (slot *battleWeaponSlot) isInRangePosition(position BattleVector) bool {
	// diagonal shots count as one move, so we take the max distance on the x or y as our actual distance away
	// i.e. 4 over, 1 up is 4 range away, 3 over 2 up is 3 range away, etc.
	return slot.token.getDistanceAway(position) <= slot.weaponRange
}

func (slot *battleWeaponSlot) isInRangeValue(rangeValue int) bool {
	return rangeValue <= slot.weaponRange
}

// get the attractiveness of a token versus a weapon
func (weapon *battleWeaponSlot) getAttractiveness(target *battleToken) float64 {

	var defense float64
	// increase the defense for jammers and beam deflectors
	switch weapon.weaponType {
	case battleWeaponTypeBeam:
		defense = float64((target.armor + target.shields)) * (1 + target.beamDefense)
	case battleWeaponTypeTorpedo:
		accuracy := weapon.getAccuracy(target.torpedoJamming)
		if target.shields >= target.armor {
			defense = float64(target.armor*2) / (accuracy)
		} else {
			capitalShipMissileFactor := 1.0
			if weapon.capitalShipMissile {
				capitalShipMissileFactor = 2
			}
			defense = float64(target.shields*2)/(accuracy) + (float64(target.armor)-float64(target.shields))/(accuracy*capitalShipMissileFactor)
		}
	}

	cost := target.cost
	attractiveNess := float64(cost.Boranium+cost.Resources) / float64(defense)
	log.Debug().Msgf("weapon %s attractiveness to %s = %f", weapon.slot.HullComponent, target.designName, attractiveNess)
	return attractiveNess
}

// Find all the targets for this weapon
func (weapon *battleWeaponSlot) findTargets(tokens []*battleToken) (targets []*battleToken) {
	attacker := weapon.token
	primaryTarget := attacker.PrimaryTarget
	secondaryTarget := attacker.SecondaryTarget

	var primaryTargets []*battleToken
	var secondaryTargets []*battleToken

	// Find all enemy tokens
	for _, token := range tokens {
		if !token.isStillInBattle() {
			continue
		}
		if !attacker.willAttack(token.PlayerNum) {
			continue
		}

		// if we will target this
		if token.isTargetOf(primaryTarget) && weapon.willDamage(token) {
			primaryTargets = append(primaryTargets, token)
		} else if token.isTargetOf(secondaryTarget) && weapon.willDamage(token) {
			secondaryTargets = append(secondaryTargets, token)
		}
	}

	// our list of available targets is all primary and all secondary targets in range
	sort.Slice(primaryTargets, func(i, j int) bool {
		return weapon.getAttractiveness(primaryTargets[i]) > weapon.getAttractiveness(primaryTargets[j])
	})
	sort.Slice(secondaryTargets, func(i, j int) bool {
		return weapon.getAttractiveness(secondaryTargets[i]) > weapon.getAttractiveness(secondaryTargets[j])
	})

	targets = make([]*battleToken, 0, len(primaryTargets)+len(secondaryTargets))
	targets = append(targets, primaryTargets...)
	targets = append(targets, secondaryTargets...)
	return targets
}

// get all targets in range of this token.
func (weapon *battleWeaponSlot) getTargetsInRange() []*battleToken {
	tokensInRange := make([]*battleToken, 0, len(weapon.targets))

	for _, token := range weapon.targets {
		if !weapon.isInRange(token) {
			continue
		}
		tokensInRange = append(tokensInRange, token)
	}
	return tokensInRange
}

// get the accuracy of a torpedo against a target
func (weapon *battleWeaponSlot) getDamage(dist int, beamDefense, beamDropoff float64) int {
	if weapon.weaponType == battleWeaponTypeTorpedo {
		return weapon.power
	}
	// we're a beam
	if weapon.weaponRange > 0 {
		return int(math.Ceil(float64(weapon.power) * (1 - float64(dist)/float64(weapon.weaponRange)*beamDropoff) * (1 - beamDefense)))
	}
	return int(math.Ceil(float64(weapon.power) * (1 - beamDefense)))
}

func (weapon *battleWeaponSlot) getEstimatedTorpedoDamageToTarget(target *battleToken) battleWeaponDamage {
	numTorpedos := weapon.slotQuantity * weapon.token.Quantity
	accuracy := weapon.getAccuracy(target.torpedoJamming)
	hits := int(float64(numTorpedos) * accuracy)
	misses := numTorpedos - hits

	// estimate how much damage we'll actually do
	damage := weapon.power * hits

	// figure out how much total armor this stack has
	totalArmor := target.armor*target.Quantity - int(float64(target.QuantityDamaged)*target.Damage)

	var bwd battleWeaponDamage
	bwd.shieldDamage = MinInt(target.stackShields, int(float64(damage)/2))
	bwd.armorDamage = MinInt(totalArmor, damage-bwd.shieldDamage)

	// for any missed torpedos, they damage shields at 1/8th, so add that
	// to shield damage if still there
	missShieldDamage := int(math.Round(float64(weapon.power*misses) / 8))
	bwd.shieldDamage = MinInt(target.stackShields, bwd.shieldDamage+missShieldDamage)

	return bwd
}

// get the damage of a single torpedo to a target. Not currently being used... I'm not sure it
// makes sense to have a separate single torpedo damage calc since the torpedos really need to be fired
// in order accumulating damage as they go, destroying ships, etc
func (weapon *battleWeaponSlot) getTorpedoDamageToTarget(target *battleToken) battleWeaponDamage {

	bwd := battleWeaponDamage{}
	damage := weapon.power
	armor := target.armor
	shields := float64(target.stackShields)
	shipDamage := target.Damage

	// torpedos do half damage to shields, half to armor (until shields are gone, when they do full armor damage)
	var shieldDamage float64
	armorDamage := float64(damage)
	if target.stackShields > 0 {
		shieldDamage = float64(0.5) * float64(damage)
		armorDamage = float64(0.5) * float64(damage)
	}

	// apply up to half our damage to shields
	// anything leftover goes to armor
	afterShieldsDamaged := shields - shieldDamage
	var actualShieldDamage float64
	if afterShieldsDamaged < 0 {
		// We did more damage to shields than they had remaining
		// apply the difference to armor
		actualShieldDamage = shieldDamage + afterShieldsDamaged
		armorDamage += float64(-afterShieldsDamaged)
	} else {
		actualShieldDamage = shieldDamage
	}

	// we have our final shield damage, record it and apply it
	bwd.shieldDamage = int(math.Round(actualShieldDamage))
	shields -= actualShieldDamage

	if shields <= 0 && weapon.capitalShipMissile {
		// capital ship missiles double damage after shields are gone
		armorDamage *= 2
	}

	// we have our final shield damage, record it and apply it
	bwd.armorDamage = int(math.Round(armorDamage))
	shipDamage += armorDamage

	// this torpedo blew up a ship, hooray!
	if shipDamage >= float64(armor) {
		// remove a ship from this stack
		bwd.numDestroyed = 1
		bwd.armorDamage = armor // reset armor damage to the armor of this ship, everything else is lost
	} else if shipDamage > 0 {
		// we werne't able to destroy a ship, but we damaged one
		bwd.quantityDamaged = 1
		bwd.damage = shipDamage
	}

	return bwd
}

// given a volly of beam damage, get how much will be applied to this target's shields and armor
// and how much is leftover after
func (weapon *battleWeaponSlot) getBeamDamageToTarget(damage int, target *battleToken, beamRangeDropoff float64) battleWeaponDamage {
	dist := weapon.token.getDistanceAway(target.Position)

	return weapon.getBeamDamageToTargetAtDistance(damage, target, dist, beamRangeDropoff)
}

// get beam damage to a target adjusted for distance
func (weapon *battleWeaponSlot) getBeamDamageToTargetAtDistance(damage int, target *battleToken, dist int, beamRangeDropoff float64) battleWeaponDamage {
	// drain any range/defelctor penalty from beam damage
	damage = getBeamDamageAtDistance(damage, weapon.weaponRange, dist, target.beamDefense, beamRangeDropoff)

	// sappers only damage shields, can't damage more shields than we have
	if weapon.damagesShieldsOnly {
		return battleWeaponDamage{shieldDamage: MinInt(target.stackShields, damage)}
	}

	armor := target.armor
	shields := target.stackShields
	if damage >= shields {
		bwd := battleWeaponDamage{}
		bwd.shieldDamage = shields
		bwd.armorDamage = damage - shields

		// if this stack is damaged froma previous hit, account for that
		existingDamage := target.Damage * float64(target.QuantityDamaged)
		newDamage := float64(bwd.armorDamage) + existingDamage

		// see how many ships were destroyed
		bwd.numDestroyed = int(newDamage / float64(armor))
		newDamage -= float64(bwd.numDestroyed) * float64(armor)
		if newDamage > 0 {
			bwd.quantityDamaged = target.Quantity - bwd.numDestroyed
			if bwd.quantityDamaged > 0 {
				bwd.damage = newDamage / float64(bwd.quantityDamaged)
			}
		}

		if bwd.numDestroyed >= target.Quantity {
			// we killed the whole stack, make sure our damage numbers reflect that
			bwd.numDestroyed = target.Quantity

			// we destroyed all armor remaining and have some possible leftover damage
			// at this point our armor damage and damage are the same
			bwd.armorDamage = (armor * bwd.numDestroyed) - int(existingDamage)
			bwd.leftover = damage - bwd.armorDamage - bwd.shieldDamage
			bwd.quantityDamaged = 0
			bwd.damage = 0
		}

		return bwd
	}

	// we didn't get through the shields
	return battleWeaponDamage{shieldDamage: damage}
}

// get the accuracy of a torpedo against a target
func (weapon *battleWeaponSlot) getAccuracy(torpedoJamming float64) float64 {
	if torpedoJamming >= weapon.torpedoInaccuracyMulti {
		return weapon.accuracy * (1 - (torpedoJamming - weapon.torpedoInaccuracyMulti))
	} else {
		return weapon.accuracy + (1-(weapon.accuracy))*(weapon.torpedoInaccuracyMulti-torpedoJamming)
	}
}
