package cs

import (
	"math"
	"slices"
)

// The TechComparer interface compares techs and techHullComponents
// to determine the one most suitable for a particular ship's purpose.
// TODO: Migrate all old tech getters in TechStore and convert them into comparers
type TechComparer interface {
	compareStargates(hc, other *TechHullComponent) bool
	compareTorpedos(hc, other *TechHullComponent) (float64, float64)
	compareWeaponPowers(hc, other *TechHullComponent) bool
	compareFieldsByTag(hc, other *TechHullComponent, tag TechTag, light bool) bool
	GetBestComponentWithTag(design *ShipDesign, hullSlotType HullSlotType, tag TechTag) *TechHullComponent
	getMostNeededComponent(design *ShipDesign, hullSlotType HullSlotType, qty int) (*TechHullComponent, error)
}

func NewTechComparer(rules *Rules, player *Player) TechComparer {
	return &techCompare{rules, player}
}

type techCompare struct {
	rules  *Rules
	player *Player
}

// Compare 2 stargates and return true if the 2nd one is superior.
//
// 1st priority mass, 2nd priority distance; ranking used as tiebreaker
func (tc *techCompare) compareStargates(hc, other *TechHullComponent) bool {
	if hc != other {
		switch {
		case hc.SafeHullMass < other.SafeHullMass:
			return false
		case hc.SafeHullMass > other.SafeHullMass:
			return true
		case hc.SafeRange < other.SafeRange: // same safe mass; check safe distance
			return false
		case hc.SafeRange > other.SafeRange:
			return true
		case hc.Ranking > other.Ranking: // same distance & range; compare ranking
			return true
		}
	}
	return false
}

// Compare 2 torpedoes or missiles and return true if the 2nd one is superior.
// Breaks ties by item ranking if all else fails
func (tc *techCompare) compareTorpedos(hc, other *TechHullComponent) (float64, float64) {
	var hcPower float64
	var otherPower float64
	capMissileMulti := 1.5
	// TODO: Figure out a value multi for capital ship missiles that makes sense
	// missiles do 2x damage on shieldless foes, but enemies don't always have shields down

	hcPower = float64(hc.Power) * float64(hc.Accuracy+10) / 100 // add a bit of accuracy bonus to account for computing vs jamming
	if hc.CapitalShipMissile {
		hcPower *= capMissileMulti
	}
	otherPower = float64(other.Power) * float64(other.Accuracy+10) / 100
	if other.CapitalShipMissile {
		otherPower *= capMissileMulti
	}
	return hcPower, otherPower
}

// compare 2 weapons' estimated damage values
// and return true if the 2nd component is better
func (tc *techCompare) compareWeaponPowers(hc, other *TechHullComponent) bool {
	rules := tc.rules
	if hc == nil || hc.Category != other.Category {
		return false
	} else if other == nil {
		return true
	}

	var hcPower, otherPower float64
	if hc.Category == TechCategoryTorpedo {
		hcPower = float64(hc.Power) * float64(hc.Accuracy+10) / 100 // add a bit of accuracy bonus to account for computing vs jamming
		if hc.CapitalShipMissile {
			hcPower *= 1.5
		}
		otherPower := float64(other.Power) * float64(other.Accuracy+10) / 100
		if other.CapitalShipMissile {
			otherPower *= 1.5
		}
	} else {
		// TODO: Rework this once damageShieldsOnly gets refactored into a variable dmg multiplier
		hcPower = float64(hc.Power)
		if !hc.Gatling {
			hcPower *= 1 - rules.BeamRangeDropoff
		}
		otherPower = float64(other.Power)
		if !other.Gatling {
			otherPower *= 1 - rules.BeamRangeDropoff
		}
	}
	return otherPower >= hcPower
}

// Compare 2 TechHullComponents by a field determined by the specified TechTag
// (alongside cost efficiency in certain cases).
// Returns true if the 2nd component is superior;
// precedence is given to the higher rated component in case of a tie.
//
// light determines whether to check weight for shields/armors
func (tc *techCompare) compareFieldsByTag(hc, other *TechHullComponent, tag TechTag, light bool) bool {
	player := tc.player

	if other == nil {
		return false
	} else if hc == nil {
		return true
	}

	var score, otherScore float64
	// whether to care about cost eff calcs
	var costTypesToCheck = false

	switch tag {
	case TechTagArmor, TechTagShield:
		// grab shield and armor stats and see which one makes number beeeeger
		hcArmor, hcShield := getArmorShieldAmounts(float64(hc.Armor), float64(hc.Shield), 1, player.Race.Spec, hc.Category == TechCategoryArmor)
		otherArmor, otherShield := getArmorShieldAmounts(float64(other.Armor), float64(other.Shield), 1, player.Race.Spec, other.Category == TechCategoryArmor)
		score = hcArmor + hcShield
		otherScore = otherArmor + otherShield

		if light {
			if hc.Mass > 30 {
				score /= (1 + float64(hc.Mass-30)/10)
			}
			if other.Mass > 30 {
				otherScore /= 1 + float64(other.Mass-30)/10
			}
		}
	case TechTagBeamCapacitor:
		score = hc.BeamBonus
		otherScore = other.BeamBonus
	case TechTagBeamDeflector:
		score = hc.BeamDefense
		otherScore = other.BeamDefense
	case TechTagScanner:
		if hc.ScanRangePen > 0 {
			if other.ScanRangePen > 0 {
				score = float64(hc.ScanRangePen)
				otherScore = float64(other.ScanRangePen)
			} else {
				// 2nd tech doesn't pen scan; 1st wins by default
				return false
			}
		} else if other.ScanRangePen > 0 {
			// 1st tech doesn't pen scan; 2nd wins by default
			return true
		} else {
			// neither tech can pen scan; just use regular scan ranges
			score = float64(hc.ScanRange)
			otherScore = float64(other.ScanRange)
		}
	case TechTagInitiativeBonus:
		score = float64(hc.InitiativeBonus)
		otherScore = float64(other.InitiativeBonus)
	case TechTagTorpedoJammer:
		score = hc.TorpedoJamming
		otherScore = other.TorpedoJamming
	case TechTagBeamWeapon, TechTagShieldSapper, TechTagGatlingGun:
		score = float64(hc.Power) * math.Pow(float64(hc.Range), 2)
		otherScore = float64(other.Power) * math.Pow(float64(other.Range), 2)
	case TechTagTorpedo, TechTagCapitalShipMissile:
		score, otherScore = tc.compareTorpedos(hc, other)
	case TechTagColonyModule:
		score = 1
		otherScore = 1
		costTypesToCheck = true // literally ALL we care about is cost efficiency
	case TechTagCargoPod:
		score = float64(hc.CargoBonus)
		otherScore = float64(other.CargoBonus)
		costTypesToCheck = true
	case TechTagFuelTank:
		score = float64(hc.FuelBonus + 5*hc.FuelGeneration)
		otherScore = float64(other.FuelBonus + 5*other.FuelGeneration)
		costTypesToCheck = true
	case TechTagMineLayer, TechTagHeavyMineLayer, TechTagSpeedMineLayer:
		otherScore = float64(other.MineLayingRate)
		score = float64(hc.MineLayingRate)
		costTypesToCheck = true
	case TechTagBomb, TechTagSmartBomb:
		score = float64(hc.KillRate)
		otherScore = float64(other.KillRate)
		costTypesToCheck = true
	case TechTagStructureBomb:
		score = float64(hc.StructureDestroyRate)
		otherScore = float64(other.StructureDestroyRate)
		costTypesToCheck = true
	case TechTagCloak:
		score = float64(hc.CloakUnits)
		otherScore = float64(other.CloakUnits)
	case TechTagManeuveringJet:
		score = float64(hc.MovementBonus)
		otherScore = float64(other.MovementBonus)
	case TechTagMassDriver:
		score = float64(hc.PacketSpeed)
		otherScore = float64(other.PacketSpeed)
	case TechTagStargate:
		return tc.compareStargates(hc, other)
	case TechTagTerraformingRobot:
		score = float64(hc.TerraformRate)
		otherScore = float64(other.TerraformRate)
		costTypesToCheck = true
	case TechTagMiningRobot:
		score = float64(hc.MiningRate)
		otherScore = float64(other.MiningRate)
		costTypesToCheck = true
	}

	scoreRatio := otherScore / score
	costRatio := 1.0
	if costTypesToCheck {
		hcCost := getPlayerCostFloat64(hc.Tech, player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset)
		otherCost := getPlayerCostFloat64(other.Tech, player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset)
		costRatio = getCostEfficiencyRatio(player, otherCost, hcCost, Resources)
	}
	return scoreRatio > costRatio ||
		(scoreRatio == costRatio && other.Ranking >= hc.Ranking)
	// FOR THE RECORD, this works out to be equivalent to comparing unit prices
	// If no cost ratio is used, it also is equivalent to simply comparing the scores
}

// get the best TechHullComponent for the specified HullSlotType(s) based on the specified TechTag.
// Used for "apples to apples" comparisons of similar part types/roles
func (tc *techCompare) GetBestComponentWithTag(design *ShipDesign, hullSlotType HullSlotType, tag TechTag) *TechHullComponent {
	// PROGRAMMER'S NOTE: the reason I didn't make this a property of techStore is because
	// we already have to pass in the Rules struct to check for warship stat hardcaps anyways
	rules := tc.rules
	player := tc.player
	store := rules.techs
	var bestTech *TechHullComponent
	bestCost := costFloat64{1, 1, 1, 1}

	// get list of all components we can use for this slot & hull type
	comps := store.GetHullComponentsByHullSlotType(player, hullSlotType, design.Hull)

	for _, hc := range comps {
		// we set match to false and catalog the part as soon as a single tag matches our list
		// and is better than our current item
		hasTag := false
		// manually cover cases for tags being subsets of other categories so we don't end up with
		// sapper only ships
		switch tag {
		case TechTagBomb:
			hasTag = hc.Tags[TechTagBomb] && !hc.Tags[TechTagStructureBomb] && !hc.Tags[TechTagSmartBomb]
		case TechTagBeamWeapon, TechTagTorpedo, TechTagCapitalShipMissile: // backup for IF we get shield sapping torpedoes
			hasTag = hc.Tags[tag] && !hc.Tags[TechTagShieldSapper] // && !hc.Tags[TechTagCapitalShipMissile]
		default:
			hasTag = hc.Tags[tag]
		}
		if hasTag && tc.compareFieldsByTag(bestTech, hc, tag, design.Purpose.IsLightShip()) {
			// we have the tag and it's better than what we already have; tack it on
			bestTech = hc
			bestCost = getPlayerCostFloat64(bestTech.Tech, player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset)
		}
	}

	// for armor, check if adding armor is even efficient at all
	// TODO: Rework this into a full on cost benefit analysis later on
	if tag == TechTagArmor && bestTech != nil && bestTech.Armor > 0 &&
		design.Purpose != ShipDesignPurposeStartingFighter {
		hull := store.GetHull(design.Hull)
		bestArmor, bestShield := getArmorShieldAmounts(float64(bestTech.Armor), float64(bestTech.Shield), 1, player.Race.Spec, bestTech.Category == TechCategoryArmor)
		hullArmor, _ := getArmorShieldAmounts(float64(hull.Armor), float64(hull.Shield), 1, player.Race.Spec, false)
		hullCost := getPlayerCostFloat64(hull.Tech, player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset)
		if hullArmor/(bestArmor+bestShield) >
			getCostEfficiencyRatio(player, hullCost, bestCost, CostType(design.Spec.Cost.HighestType(1))) {
			// it's more efficient for us to use the base hull armor than our added item armor; leave it off
			return nil
		}
	}
	return bestTech
}

// get the most needed component for a *warship* design based on relative bonuses of various parts
func (tc *techCompare) getMostNeededComponent(design *ShipDesign, hullSlotType HullSlotType, qty int) (*TechHullComponent, error) {
	// TODO: Add special "design modes" to change part grabbing behavior and allow for multiple "correct designs"
	rules := tc.rules
	player := tc.player
	store := rules.techs

	comps := store.GetHullComponentsByHullSlotType(player, hullSlotType, design.Hull)
	var bestTech *TechHullComponent
	bestBonus := 1.0

	for _, hc := range comps {
		if !player.HasTech(&hc.Tech) ||
			(len(hc.Tech.Requirements.HullsAllowed) > 0 && !slices.Contains(hc.Tech.Requirements.HullsAllowed, design.Hull)) ||
			(len(hc.Tech.Requirements.HullsDenied) > 0 && slices.Contains(hc.Tech.Requirements.HullsDenied, design.Hull)) {
			// we do not have or cannot use this part; skip to the next item
			continue
		}
		// compare bonuses and part costs
		hcBonus := design.getWarshipPartBonus(rules, player, hc, qty)
		if roundFloat(hcBonus, 3) == 1 {
			continue // part does not benefit us in any way, shape or form
		}
		if hcBonus >= bestBonus {
			bestTech = hc
			bestBonus = hcBonus
		}
	}
	return bestTech, nil
}

// return the total % amount this TechHullComponent boosts our ship's performance,
// using its TechTags to evaluate individual bonuses
//
// Used to assess parts from different categories to determine
// which one is the best for us to use
func (design *ShipDesign) getWarshipPartBonus(rules *Rules, player *Player, hc *TechHullComponent, qty int) float64 {
	relativeBoost := 1.0
	checkedShield := false

	// check tags individually and tally up the numbers
tagLoop:
	for tag := range hc.Tags {
		switch {
		case (tag == TechTagArmor || tag == TechTagShield) && !checkedShield:
			// grab shield and armor stats
			hcArmor, hcShield := getArmorShieldAmounts(float64(hc.Armor), float64(hc.Shield), qty, player.Race.Spec, hc.Category == TechCategoryArmor)
			oldArmor := float64(design.Spec.Armor)
			oldShield := float64(design.Spec.Shields)
			newArmor := math.Max(hcArmor+oldArmor, 1) // prevents divide by 0 errors
			newShield := math.Max(hcShield+oldShield, 1)

			if oldArmor/oldShield < 1.2 && oldArmor/oldShield > 1/1.2 {
				// our armor ratio is close enough that we don't really need
				// more armor/shield; why reward messing with perfection?
				continue tagLoop
			}

			// apply scaling score penalty for adding more armor/shield when we already have lots
			// margin of error before penalty kicks in is 30%
			// TODO: Make this look at costs later once I have energy ig
			if newArmor/newShield > 1.3 {
				// reduce our effective armor bonus for adding too much armor
				hcArmor /= 1 + (math.Min(newArmor/newShield, 4.3) - 1.3)
				newArmor = hcArmor + oldArmor
			} else if newShield/newArmor > 1.3 {
				// reduce our effective shield bonus for adding too much shield
				hcShield /= 1 + (math.Min(newShield/newArmor, 4.3) - 1.3)
				newShield = hcShield + oldShield
			}

			// wrap up by adding relative armor & shield boosts to get overall relative bonus
			relativeBoost *= (newArmor + newShield) / (oldArmor + oldShield)
			checkedShield = true // needed to prevent shield/armor parts from double counting themselves
		case tag == TechTagBeamCapacitor:
			if design.Purpose.IsTorpedoShip() {
				break // beam bonus meaningless on missile boats
			}
			// new beam bonus / old beam bonus
			relativeBoost *= (math.Min(getNewBeamBonus(design.Spec.BeamBonus, hc.BeamBonus, qty), rules.BeamBonusCap) / design.Spec.BeamBonus)
		case tag == TechTagBeamDeflector:
			beamDefense := math.Pow(1-hc.BeamDefense, float64(qty))
			if getNewBeamDefense(design.Spec.BeamDefense, hc.BeamDefense, qty) <= 0.3 {
				// penalize having too much beam defense in favor of adding other things
				// up to 3x penalty
				// TODO: Rework after BeamDefense cleanup
				relativeBoost *= 1 + ((1/beamDefense)-1)/
					(3-(getNewBeamDefense(design.Spec.BeamDefense, hc.BeamDefense, qty)/2))
			} else {
				relativeBoost *= 1 / beamDefense
			}
		case tag == TechTagInitiativeBonus:
			// do nothing lol
		case tag == TechTagTorpedoBonus:
			if design.Purpose.IsBeamShip() {
				break // torpedo bonus meaningless on beamships
			}
			fallthrough
		case tag == TechTagTorpedoJammer:
			relativeBoost *= design.Spec.getJamOrComputerIncrease(rules, hc, qty, tag == TechTagTorpedoJammer)

			// all other cases are either stupidly painful to quantify or don't get checked in the actual design function
		}

		if design.Purpose.IsBeamShip() && hc.Mass > 30 {
			// reduce bonus if component is too heavy and we need to go fast
			relativeBoost = 1 + (relativeBoost-1)/(1+float64(hc.Mass-30)/10)
		}
	}
	return roundFloat(relativeBoost, 4)
}

// return relative factor by which a jammer or computer boosts our relative torpedo defense/offense
//
// Formula: 1+oldJamming / 1+newJamming (https://www.desmos.com/calculator/cpcyiloqeg)
//
// fieldToCheck determines which stat is being calculated for; true checks jamming whereas false checks computing
func (spec *ShipDesignSpec) getJamOrComputerIncrease(rules *Rules, hc *TechHullComponent, qty int, fieldToCheck bool) float64 {
	var oldBonus, hcBonus, cap, jamMulti float64
	if fieldToCheck {
		jamMulti = rules.JammerMulti[spec.Starbase]
		cap = rules.JammerCap[spec.Starbase] * jamMulti
		oldBonus = spec.TorpedoJamming
		hcBonus = hc.TorpedoJamming
	} else {
		jamMulti = 1
		cap = 1
		oldBonus = spec.TorpedoBonus
		hcBonus = hc.TorpedoBonus
	}

	if oldBonus == cap {
		return 1
	}

	// *I HATE FLOATING POINT ROUNDING ERRORS*
	newBonus := roundFloat(math.Min(getNewJamming(oldBonus, hcBonus, jamMulti, qty), cap), 7)

	// return final score, rounded to 4 decimal places for ease of comparison
	return roundFloat((1+newBonus)/(1+oldBonus), 4)
}
