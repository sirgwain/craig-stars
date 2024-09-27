package cs

import (
	"slices"
)

type TechTag string

const (
	TechTagArmor              TechTag = "Armor"
	TechTagTorpedoBonus       TechTag = "TorpedoBonus"
	TechTagBeamCapacitor      TechTag = "BeamCapacitor"
	TechTagBeamDeflector      TechTag = "BeamDeflector"
	TechTagBeamWeapon         TechTag = "BeamWeapon"
	TechTagBomb               TechTag = "Bomb"
	TechTagCapitalShipMissile TechTag = "CapitalShipMissile"
	TechTagCargoPod           TechTag = "CargoPod"
	TechTagCloak              TechTag = "Cloak"
	TechTagColonyModule       TechTag = "ColonyModule"
	TechTagEngine             TechTag = "Engine"
	TechTagFuelTank           TechTag = "FuelTank"
	TechTagGatlingGun         TechTag = "GatlingGun"
	TechTagHeavyMineLayer     TechTag = "HeavyMineLayer"
	TechTagInitiativeBonus    TechTag = "InitiativeBonus"
	TechTagMassDriver         TechTag = "MassDriver"
	TechTagManeuveringJet     TechTag = "ManeuveringJet"
	TechTagMineLayer          TechTag = "MineLayer"
	TechTagMiningRobot        TechTag = "MiningRobot"
	TechTagRamscoop           TechTag = "Ramscoop"
	TechTagTerraformingRobot  TechTag = "TerraformingRobot"
	TechTagScanner            TechTag = "Scanner"
	TechTagShield             TechTag = "Shield"
	TechTagShieldSapper       TechTag = "ShieldSapper"
	TechTagSmartBomb          TechTag = "SmartBomb"
	TechTagSpeedMineLayer     TechTag = "SpeedMineLayer"
	TechTagStargate           TechTag = "Stargate"
	TechTagStructureBomb      TechTag = "StructureBomb"
	TechTagTerraforming       TechTag = "Terraforming"
	TechTagTorpedo            TechTag = "Torpedo"
	TechTagTorpedoJammer      TechTag = "TorpedoJammer"
)

type TechTags map[TechTag]bool

// Create a new TechTags map from a list of TechTag items, or an empty map if none are specified
func newTechTags(tags ...TechTag) TechTags {
	newMap := TechTags{}
	for _, t := range tags {
		newMap[t] = true
	}
	return newMap
}

// returns true if tt has ALL of the specified tags
func (tt TechTags) hasAllTags(tags ...TechTag) bool {
	for _, tag := range tags {
		if !tt[tag] {
			return false
		}
	}
	return true
}

// returns true if tt has AT LEAST 1 of the specified tags
func (tt TechTags) hasTag(tags ...TechTag) bool {
	for _, tag := range tags {
		if tt[tag] {
			return true
		}
	}
	return false
}

// returns list of all tags in tt, sorted alphabetically
func (tt TechTags) GetTags() []string {
	var list []string
	for k, v := range tt {
		if v {
			list = append(list, string(k))
		}
	}
	slices.Sort(list)
	return list
}

// return number of unique tags in tt
func (tt TechTags) CountTags() int {
	count := 0
	for _, v := range tt {
		if v {
			count += 1
		}
	}
	return count
}

// Compare 2 techHullComponents based on a field determined by the specified TechTag
// Precedence is given to 2nd component in case of tie
//
// NOTE: ONLY WORKS ON NON-COMBAT-FOCUSED PARTS. USE spec.getbestWarshipPart() instead
func CompareFieldsByTag(player *Player, hc, other *TechHullComponent, tag TechTag) bool {
	if other == nil {
		return false
	} else if hc == nil {
		return true
	}
	switch tag {
	case TechTagScanner:
		if hc.ScanRangePen > 0 {
			if other.ScanRangePen > 0 {
				return other.ScanRangePen >= hc.ScanRangePen
			} else {
				// 2nd tech doesn't pen scan; 1st wins by default
				return false
			}
		} else if other.ScanRangePen > 0 {
			// 1st tech doesn't pen scan; 2nd wins by default
			return true
		}
		return other.ScanRange >= hc.ScanRange
	case TechTagColonyModule:
		return getMineralEfficiencyRatio(player, hc, other, true) >= 1
	case TechTagCargoPod:
		return float64(other.CargoBonus)/float64(hc.CargoBonus) >= getMineralEfficiencyRatio(player, hc, other, true)
		// FOR THE RECORD, this works out to be equivalent to comparing unit prices
	case TechTagFuelTank:
		return float64(other.FuelBonus+5*other.FuelGeneration)/float64(hc.FuelBonus+5*hc.FuelGeneration) >=
			getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagMineLayer, TechTagHeavyMineLayer, TechTagSpeedMineLayer:
		return float64(other.MineLayingRate)/float64(hc.MineLayingRate) >= getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagBomb, TechTagSmartBomb:
		return float64(other.KillRate)/float64(hc.KillRate) >= getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagStructureBomb:
		return float64(other.StructureDestroyRate)/float64(hc.StructureDestroyRate) > getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagCloak:
		return float64(other.CloakUnits) >= float64(hc.CloakUnits) 
	case TechTagMassDriver:
		return other.PacketSpeed >= hc.PacketSpeed
	case TechTagStargate:
		return hc.GetBestStargate(other)
	case TechTagTerraformingRobot:
		return float64(other.TerraformRate)/float64(hc.TerraformRate) >= getMineralEfficiencyRatio(player, hc, other, true)
	case TechTagMiningRobot:
		return float64(other.MiningRate)/float64(hc.MiningRate) >= getMineralEfficiencyRatio(player, hc, other, false)
	}
	return false
}
