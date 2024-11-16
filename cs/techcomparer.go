package cs

import (
	"fmt"
	"math"
	"slices"
)

// The TechComparer interface compares techs and techHullComponents
// to determine the one most suitable for a particular purpose.
type TechComparer interface {
	CompareStargates(hc, other *TechHullComponent) bool
	CompareTorpedos(player *Player, hc, other *TechHullComponent) bool
	CompareEngine(player *Player, hc, other *TechEngine, purpose FleetPurpose) *TechEngine
	CompareFieldsByTag(player *Player,  hc, other *TechHullComponent, tag TechTag, light bool) bool
	GetBestComponentWithTags(rules *Rules, player *Player, design *ShipDesign, hullSlotType HullSlotType, tags ...TechTag) (*TechHullComponent, error)
}

func NewTechComparer() TechComparer {
	return &techCompare{}
}

type techCompare struct {
}

// compare 2 stargates and determine which one is better
// 
// 1st priority mass, 2nd priority distance, ranking as tiebreaker
func (tc *techCompare) CompareStargates(hc, other *TechHullComponent) bool {
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

// return the better of the 2 provided torpedo weapons.
// Defaults to 1st if both are equal or either one is null and breaks ties by item ranking
func (tc *techCompare) CompareTorpedos(player *Player, hc, other *TechHullComponent) bool {
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

	// if the 2nd torpedo is more cost efficient (in terms of avg damage/minerals spent) than the new weapon, use it
	return otherPower/hcPower > getCostEfficiencyRatio(player, hc, other, Ironium) ||
		otherPower/hcPower == getCostEfficiencyRatio(player, hc, other, Ironium) && other.Ranking >= hc.Ranking
}

// returns the better of the 2 engines
func (tc *techCompare) CompareEngine(player *Player, hc, other *TechEngine, purpose FleetPurpose) *TechEngine {
	tech := hc.TechHullComponent
	otherTech := other.TechHullComponent
	if player.HasTech(&tech.Tech) {
		// colony ships don't want radiating engines if we would lose colonists from it
		if ((purpose == FleetPurposeColonizer || purpose == FleetPurposeColonistFreighter) && tech.Radiating &&
			!(player.Race.ImmuneRad || player.Race.Spec.HabCenter.Rad >= 85)) ||
			otherTech.Ranking > tech.Ranking {
			return other
		}
	}
	return hc
}

// Compare 2 TechHullComponents by a field determined by the specified TechTag
// (alongside cost efficiency in certain cases).
// Returns true if the 2nd component is superior;
// precedence is given to the higher rated component in case of a tie.
//
// light determines whether to check weight for shields/armors
func (tc *techCompare) CompareFieldsByTag(player *Player, hc, other *TechHullComponent, tag TechTag, light bool) bool {
	if other == nil {
		return false
	} else if hc == nil {
		return true
	}

	var score, otherScore float64
	var costTypesToCheck []CostType
	// which cost types to care about for cost eff calcs, if any
	// usually only applies for items that make up the bulk of their respective ships' cost
	// and/or ones with a definitive quantifiable stat we can price

	switch tag {
	case TechTagArmor, TechTagShield:
		// grab shield and armor stats and see which one makes number beeeeger
		hcArmor, hcShield := getActualArmorAmount(float64(hc.Armor), float64(hc.Shield), 1, player.Race.Spec, hc.Category == TechCategoryArmor)
		otherArmor, otherShield := getActualArmorAmount(float64(other.Armor), float64(other.Shield), 1, player.Race.Spec, other.Category == TechCategoryArmor)
		score = hcArmor + hcShield
		otherScore = otherArmor + otherShield

		if light {
			if hc.Mass > 30 {
				score /= float64(hc.Mass-30) / 10
			}
			if other.Mass > 30 {
				otherScore /= float64(other.Mass-30) / 10
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
		costTypesToCheck = []CostType{Resources}
	case TechTagTorpedo, TechTagCapitalShipMissile:
		return tc.CompareTorpedos(player, hc, other)
	case TechTagColonyModule:
		score = 1
		otherScore = 1
		costTypesToCheck = []CostType{} // literally ALL we care about is cost efficiency
	case TechTagCargoPod:
		score = float64(hc.CargoBonus)
		otherScore = float64(other.CargoBonus)
		costTypesToCheck = []CostType{Resources}
	case TechTagFuelTank:
		score = float64(hc.FuelBonus + 5*hc.FuelGeneration)
		otherScore = float64(other.FuelBonus + 5*other.FuelGeneration)
		costTypesToCheck = []CostType{Resources}
	case TechTagMineLayer, TechTagHeavyMineLayer, TechTagSpeedMineLayer:
		otherScore = float64(other.MineLayingRate)
		score = float64(hc.MineLayingRate)
		costTypesToCheck = []CostType{Resources}
	case TechTagBomb, TechTagSmartBomb:
		score = float64(hc.KillRate)
		otherScore = float64(other.KillRate)
		costTypesToCheck = []CostType{}
	case TechTagStructureBomb:
		score = float64(hc.StructureDestroyRate)
		otherScore = float64(other.StructureDestroyRate)
		costTypesToCheck = []CostType{}
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
		return tc.CompareStargates(hc, other)
	case TechTagTerraformingRobot:
		score = float64(hc.TerraformRate)
		otherScore = float64(other.TerraformRate)
		costTypesToCheck = []CostType{Resources}
	case TechTagMiningRobot:
		score = float64(hc.MiningRate)
		otherScore = float64(other.MiningRate)
		costTypesToCheck = []CostType{Resources}
	}

	scoreRatio := otherScore / score
	costRatio := 1.0
	if costTypesToCheck != nil {
		costRatio = getCostEfficiencyRatio(player, other, hc, costTypesToCheck...)
	}
	return scoreRatio > costRatio ||
		(scoreRatio == costRatio && other.Ranking > hc.Ranking)
	// FOR THE RECORD, this works out to be equivalent to comparing unit prices
	// If you don't believe this yourself, do some algebra
}

// get the best TechHullComponent for the specified HullSlotType(s) that also contains the specified TechTag(s).
func (tc *techCompare) GetBestComponentWithTags(rules *Rules, player *Player, design *ShipDesign, hullSlotType HullSlotType, tags ...TechTag) (*TechHullComponent, error) {
	// PROGRAMMER'S NOTE: the reason I didn't make this a property of techStore is because
	// we already have to pass in the Rules struct to check for warship stat hardcaps anyways

	var bestTech *TechHullComponent

	store := rules.techs
	hull := store.GetHull(design.Hull)
	if hull == nil {
		return nil, fmt.Errorf("failed to get hull %v from tech store", hull)
	}

	// get list of components for the TechHullTypes we can use
	var comps []*TechHullComponent = store.GetHullComponentsByHullSlotType(player, hullSlotType, hull.Name)

	for _, hc := range comps {
		if !player.HasTech(&hc.Tech) ||
			(len(hc.Tech.Requirements.HullsAllowed) > 0 && !slices.Contains(hc.Tech.Requirements.HullsAllowed, hull.Name)) ||
			(len(hc.Tech.Requirements.HullsDenied) > 0 && slices.Contains(hc.Tech.Requirements.HullsDenied, hull.Name)) {
			// we cannot use this part; skip to the next item
			continue
		}

		// need only 1 tag to match
		// we set match to false and catalog the part as soon as a single tag matches our list
		// and is better than our current item
		hasTag := false
		for _, tag := range tags {
			// manually cover cases for tags being subsets of other categories so we don't end up with
			// sapper only ships
			switch tag {
			case TechTagBomb:
				hasTag = hc.Tags[TechTagBomb] && !hc.Tags[TechTagStructureBomb] && !hc.Tags[TechTagSmartBomb]
			case TechTagBeamWeapon:
				hasTag = hc.Tags[TechTagBeamWeapon] && !hc.Tags[TechTagShieldSapper]
			default:
				hasTag = hc.Tags[tag]
			}
			if hasTag && (bestTech == nil || tc.CompareFieldsByTag(player, bestTech, hc, tag, design.Purpose.IsBeamShip())) {
				// we have the tag and it's better than what we already have; tack it on
				bestTech = hc
				break
			}
		}
	}

	return bestTech, nil
}