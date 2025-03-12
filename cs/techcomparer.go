package cs

import (
	"fmt"

	"github.com/rs/zerolog"
	"math"
)

// The TechComparer interface compares techs and techHullComponents
// to determine the one most suitable for a particular ship's purpose.
// TODO: Add cost benefit analysis and support
type TechComparer interface {
	GetBestComponentWithTag(design *ShipDesign, hullSlotType HullSlotType, qty int, tag TechTag) *TechHullComponent
	CompareWeaponPowers(hc, other *TechHullComponent) bool
	GetMostNeededComponent(design *ShipDesign, hullSlotType HullSlotType, qty int) (*TechHullComponent, error)
}

func NewTechComparer(rules *Rules, player *Player, log zerolog.Logger) TechComparer {
	comparerLogger := log.With().
		Int("PlayerNum", player.Num).
		Str("Player", player.Name).
		Logger()
	return &techCompare{rules, player, comparerLogger}
}

// TODO: Add more logging to functions (ideally not too much though)
type techCompare struct {
	rules  *Rules
	player *Player
	log    zerolog.Logger
}

// GetBestComponentWithTag finda ans returns the best usable TechHullComponent
// for the specified TechTag.
//
// Used for "apples to apples" comparisons of parts when we
// know what we are looking for ahead of time
func (tc *techCompare) GetBestComponentWithTag(design *ShipDesign, hullSlotType HullSlotType, qty int, tag TechTag) (bestTech *TechHullComponent) {
	rules := tc.rules
	player := tc.player
	store := rules.techs

	// get slice of all components we can use for this slot & hull type
	// TODO: Cache this maybe
	comps := store.GetHullComponentsByHullSlotType(player, hullSlotType, design.Hull)

	for _, hc := range comps {
		hasTag := false
		// manually cover cases for tags being subsets of other categories
		// so we don't end up with sapper only ships or smart bomb only bombers
		switch tag {
		case TechTagBomb:
			hasTag = hc.Tags.HasTag(TechTagBomb) && !hc.Tags.HasTag(TechTagStructureBomb) && !hc.Tags.HasTag(TechTagSmartBomb)
		case TechTagBeamWeapon, TechTagTorpedo, TechTagCapitalShipMissile: // backup for IF we get shield sapping torpedoes
			hasTag = hc.Tags.HasTag(tag) && !hc.Tags.HasTag(TechTagShieldSapper) // && !hc.Tags.HasTag(TechTagCapitalShipMissile)
		default:
			hasTag = hc.Tags.HasTag(tag)
		}
		if hasTag && (bestTech == nil || tc.compareFieldsByTag(design, bestTech, hc, qty, tag)) {
			// we have the tag and it's better than what we already have; tack it on
			bestTech = hc
		}
	}

	return bestTech
}

// Compare 2 TechHullComponents by a field determined by the specified TechTag
// (alongside cost efficiency in certain cases).
//
// Returns true if the 2nd component is superior;
// precedence is given to the higher rated component in case of a tie.
func (tc *techCompare) compareFieldsByTag(design *ShipDesign, hc, other *TechHullComponent, qty int, tag TechTag) bool {
	player := tc.player
	rules := tc.rules

	if other == nil {
		return false
	} else if hc == nil {
		return true
	}

	var score, otherScore float64
	// whether to care about cost eff calcs
	var checkCost = false

	switch tag {
	case TechTagArmor, TechTagShield:
		// grab shield and armor stats and see which one makes number beeeeger
		hcArmor, hcShield := getArmorShieldAmounts(float64(hc.Armor), float64(hc.Shield), 1, player.Race.Spec, hc.Category == TechCategoryArmor)
		otherArmor, otherShield := getArmorShieldAmounts(float64(other.Armor), float64(other.Shield), 1, player.Race.Spec, other.Category == TechCategoryArmor)
		score = hcArmor + hcShield
		otherScore = otherArmor + otherShield
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
	case TechTagTorpedoBonus:
		score = hc.TorpedoBonus
		otherScore = other.TorpedoBonus
	case TechTagTorpedoJammer:
		score = hc.TorpedoJamming
		otherScore = other.TorpedoJamming
	case TechTagTorpedo, TechTagCapitalShipMissile, TechTagBeamWeapon, TechTagGatlingGun, TechTagShieldSapper:
		return tc.CompareWeaponPowers(hc, other)
	case TechTagColonyModule:
		score = 1
		otherScore = 1
		checkCost = true // literally ALL we care about is cost efficiency
	case TechTagCargoPod:
		score = float64(hc.CargoBonus)
		otherScore = float64(other.CargoBonus)
		checkCost = true
	case TechTagFuelTank:
		score = float64(hc.FuelBonus + 5*hc.FuelGeneration)
		otherScore = float64(other.FuelBonus + 5*other.FuelGeneration)
		checkCost = true
	case TechTagMineLayer, TechTagHeavyMineLayer, TechTagSpeedMineLayer:
		otherScore = float64(other.MineLayingRate)
		score = float64(hc.MineLayingRate)
		checkCost = true
	case TechTagBomb, TechTagSmartBomb:
		score = float64(hc.KillRate)
		otherScore = float64(other.KillRate)
		checkCost = true
	case TechTagStructureBomb:
		score = float64(hc.StructureDestroyRate)
		otherScore = float64(other.StructureDestroyRate)
		checkCost = true
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
		checkCost = true
	case TechTagMiningRobot:
		score = float64(hc.MiningRate)
		otherScore = float64(other.MiningRate)
		checkCost = true
	}

	if design.Purpose.IsLightShip() {
		if hc.Mass > 30 && design.getMovement(rules, hc.Mass*qty+design.Spec.CargoCapacity) < rules.MovementMax { // don't penalize movement if we still have max movement with full cargo
			score /= (1 + float64(hc.Mass-30)/10)
		}
		if other.Mass > 30 && design.getMovement(rules, hc.Mass*qty+design.Spec.CargoCapacity) < rules.MovementMax {
			otherScore /= 1 + float64(other.Mass-30)/10
		}
	}

	if score == 0 {
		return true // prevents division by 0
	}
	scoreRatio := otherScore / score
	costRatio := 1.0
	if checkCost {
		hcCost := getPlayerCost(hc.Tech, player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset)
		otherCost := getPlayerCost(other.Tech, player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset)
		costRatio = GetCostEfficiencyRatio(otherCost, hcCost, Resources)
	}
	return scoreRatio > costRatio ||
		(scoreRatio == costRatio && other.Ranking > hc.Ranking)
	// This works out to be equivalent to comparing unit prices
	// If checkCost is false, it is also equivalent to simply comparing the scores numerically
	// a/b > 1 = a(b)/b > 1(b) = a > b
}

// Compare 2 stargates and return true if the 2nd one is superior.
//
// 1st priority mass, 2nd priority distance; cost used as final tiebreaker
func (tc *techCompare) compareStargates(hc, other *TechHullComponent) bool {
	if hc != other {
		switch {
		case other.SafeHullMass < hc.SafeHullMass:
			return false
		case other.SafeHullMass > hc.SafeHullMass:
			return true
		case other.SafeRange < hc.SafeRange: // same safe mass; check safe distance
			return false
		case other.SafeRange > hc.SafeRange:
			return true
		default: // same safe mass & distance; compare price point
			hcCost := getPlayerCost(hc.Tech, tc.player.TechLevels, tc.player.Race.Spec.MiniaturizationSpec, tc.player.Race.Spec.TechCostOffset)
			otherCost := getPlayerCost(other.Tech, tc.player.TechLevels, tc.player.Race.Spec.MiniaturizationSpec, tc.player.Race.Spec.TechCostOffset)
			return GetCostEfficiencyRatio(otherCost, hcCost, Resources) > 1
		}
	}
	return false
}

// compare 2 weapons' estimated damage values
// and return true if other is better than hc
func (tc *techCompare) CompareWeaponPowers(hc, other *TechHullComponent) bool {
	rules := tc.rules
	if hc == nil {
		return true
	} else if other == nil || hc.Category != other.Category {
		return false // can't compare apples to oranges!
	}

	var hcPower, otherPower float64
	if hc.Category == TechCategoryTorpedo {
		hcPower = float64(hc.Power*hc.Accuracy*hc.Range) / 100 / 4
		otherPower = float64(other.Power) * float64(other.Accuracy) / 100 * float64(other.Range) / 4
		// TODO: Rework this once damageShieldsOnly & CapitalShipMissile get refactored into an armor dmg multiplier
		if hc.CapitalShipMissile {
			hcPower *= 1.5
		}
		if other.CapitalShipMissile {
			otherPower *= 1.5
		}
	} else {
		hcPower = float64(hc.Power * hc.Range << 1) // squares it
		otherPower = float64(other.Power * other.Range << 1)
		if !hc.Gatling {
			hcPower *= 1 - rules.BeamRangeDropoff
		}
		if !other.Gatling {
			otherPower *= 1 - rules.BeamRangeDropoff
		}
	}

	return otherPower >= hcPower
}

// get the most needed component for a *warship* design based on relative bonuses of various parts
func (tc *techCompare) GetMostNeededComponent(design *ShipDesign, hst HullSlotType, qty int) (bestTech *TechHullComponent, err error) {
	// TODO: Add special "design modes" to change part grabbing behavior and allow for multiple "correct designs"
	rules := tc.rules
	player := tc.player
	store := rules.techs

	// these parts come out sorted & filtered, so we don't need to filter it on our end
	comps := store.GetHullComponentsByHullSlotType(player, hst, design.Hull)
	bestTech = comps[0]
	bestBonus := tc.getWarshipPartBonus(design, bestTech, qty)

	hull := rules.techs.GetHull(design.Hull)
	if hull == nil {
		return nil, fmt.Errorf("getMostNeededComponent cannot compare cost efficiency; design hull %s was not found in tech store", design.Hull)
	}

	for _, hc := range comps {
		// compare bonuses and part costs
		hcBonus := tc.getWarshipPartBonus(design, hc, qty)

		if hcBonus > bestBonus || (hcBonus == bestBonus &&
			GetCostEfficiencyRatio(getPlayerCost(bestTech.Tech, player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset),
				getPlayerCost(hc.Tech, player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset),
				// TODO: Change this once AI gets custom cost type check support
				CostTypes[:]...) > 1) {
			// all else being equal, use items that cost less
			bestTech = hc
			bestBonus = hcBonus
		}
	}

	if bestTech == nil {
		return nil, nil
	}

	// for defensive components, check if adding them is even efficient at all
	// compared to using the hull's baseline stats
	if ((bestTech.Armor > 0 && hull.Armor > 0) || (bestTech.Shield > 0 && hull.Shield > 0)) && // tech gives the same stat that our hull does (no apples to oranges)
		!design.Spec.Starbase && design.Purpose != ShipDesignPurposeStartingFighter && // we are neither a staircase nor a scripted starter ship
		// TODO: Remove 2nd part of conditional post starting fleet refactor
		bestTech.Tags.hasTags([]TechTag{TechTagArmor, TechTagShield}, CombatTechTags...) { // part gives no tangible benefit aside from shield/armor bonuses

		// extract armor/shield values and compute cost ratio
		hcArmor, hcShield := getArmorShieldAmounts(float64(bestTech.Armor), float64(bestTech.Shield), 1, player.Race.Spec, bestTech.Category == TechCategoryArmor)
		hullArmor, hullShield := getArmorShieldAmounts(float64(hull.Armor), float64(hull.Shield), 1, player.Race.Spec, false)
		costRatio := GetCostEfficiencyRatio(getPlayerCost(hull.Tech, player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset),
			getPlayerCost(bestTech.Tech, player.TechLevels, player.Race.Spec.MiniaturizationSpec, player.Race.Spec.TechCostOffset),
			// TODO: Change this once AI gets custom cost type check support
			CostTypes[:]...)
		var armorShieldRatio float64
		switch {
		case bestTech.Armor <= 0 || hull.Armor <= 0: // no armor; only consider sheld
			armorShieldRatio = hullShield / hcShield
		case bestTech.Shield <= 0 || hull.Shield <= 0: // no shield; only consider armor
			armorShieldRatio = hullArmor / hcArmor
		default: // we have both types and consider both
			armorShieldRatio = (hullArmor + hullShield) / (hcArmor + hcShield)
		}

		if armorShieldRatio >= costRatio {
			tc.log.Debug().
				Str("Best Component", bestTech.Name).
				Str("Hull", hull.Name).
				Float64("Armor ratio", armorShieldRatio).
				Float64("Cost ratio", costRatio).
				Msgf("Skipping defensive component; hull more efficient")

			// Our baseline hull is more efficient than the hull armor;
			// leave it off as we can always make more ships
			return nil, nil
		}
	}

	if unitRate := (bestBonus - 1) / float64(qty); unitRate <= 0.01 && design.Spec.CloakPercent < 98 {
		// our "best part" barely helps us; try and add a cloak for funsies
		cloak := tc.GetBestComponentWithTag(design, hst, qty, TechTagCloak)
		if cloak != nil {
			tc.log.Debug().
				Str("Best Component", bestTech.Name).
				Float64("Bonus Per Qty", unitRate).
				Str("Cloak", cloak.Name).
				Int("Cloak Percent", design.Spec.CloakPercent).
				Msg("Forgoing component for cloak")
			return cloak, nil
		}
	}

	tc.log.Debug().
		Str("Best Component", bestTech.Name).
		Float64("Bonus", bestBonus).
		Msg("Obtained best component")

	return bestTech, nil
}

// return the total % amount this TechHullComponent boosts our ship's performance,
// using its TechTags to evaluate individual bonuses
//
// Used to assess parts from different categories to determine
// which one is the best for us to use
func (tc *techCompare) getWarshipPartBonus(design *ShipDesign, hc *TechHullComponent, qty int) (relativeBoost float64) {
	rules := tc.rules
	player := tc.player

	relativeBoost = 1.0
	checkedShield := false
	techTagsToCheck := newTechTags(
		TechTagArmor,
		TechTagShield,
		TechTagTorpedoBonus,
		TechTagTorpedoJammer,
		TechTagBeamCapacitor,
		TechTagBeamDeflector,
		TechTagManeuveringJet)

	// check tags individually and tally up the numbers
tagLoop:
	for tag := range hc.Tags {
		if !techTagsToCheck.HasTag(tag) {
			continue
		}

		switch {
		case (tag == TechTagArmor || tag == TechTagShield) && !checkedShield:
			// grab shield and armor stats
			hcArmor, hcShield := getArmorShieldAmounts(float64(hc.Armor), float64(hc.Shield), qty, player.Race.Spec, hc.Category == TechCategoryArmor)
			oldArmor := float64(design.Spec.Armor)
			oldShield := float64(Max(design.Spec.Shields, 1)) // prevents divide by 0 errors
			if oldArmor/oldShield < 1.2 && oldArmor/oldShield > 1/1.2 && !design.Spec.Starbase {
				// our armor ratio is good enough that we don't really need
				// more armor/shields; can just build more ships for more overall chung
				// TODO: Remove once I figure out how to make the AI realize cost exists
				continue tagLoop
			}
			newArmor := math.Max(hcArmor+oldArmor, 1) // prevents divide by 0 errors
			newShield := math.Max(hcShield+oldShield, 1)

			// apply scaling score penalty for adding more armor/shield when we already have lots
			// margin of error before penalty kicks in is 30%
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
				continue // beam bonus meaningless on missile boats
			}
			// boost *= new beam bonus / old beam bonus
			relativeBoost *= (math.Min(getNewBeamBonus(design.Spec.BeamBonus, hc.BeamBonus, qty), rules.BeamBonusCap) / design.Spec.BeamBonus)
		case tag == TechTagTorpedoBonus:
			if design.Purpose.IsBeamShip() {
				continue // torpedo bonus meaningless on beam ships
			}
			fallthrough
		case tag == TechTagTorpedoJammer, tag == TechTagBeamDeflector:
			relativeBoost *= design.Spec.getJamOrComputerBonus(rules, hc, qty, tag)
		case tag == TechTagManeuveringJet && design.Spec.Movement < rules.MovementMax && !design.Spec.Starbase:
			// add a small, staple boost to jets
			oldMove := float64(design.Spec.Movement)
			moveBoost := float64(getBattleMovement(rules.MovementMin, rules.MovementMax, design.Spec.Engine.IdealSpeed, design.Spec.MovementBonus+hc.MovementBonus*float64(qty), design.Spec.Mass+hc.Mass*qty, design.Spec.NumEngines)) - oldMove
			multi := 0.6
			if design.Purpose.IsTorpedoShip() {
				multi = 0.35 // reduce bonus multi for torp ships; they don't need it as much
			}
			relativeBoost *= (1 + multi*moveBoost/oldMove)
		}
		// all other cases are either painful to quantify or don't get checked in the actual design function

		// reduce bonus if component is too heavy and we need to go fast
		if design.Purpose.IsBeamShip() && hc.Mass > 30 {
			relativeBoost = 1 + (relativeBoost-1)/(1+float64(hc.Mass-30)/10)
		}
		roundFloat(relativeBoost, 8)
	}
	return relativeBoost
}
