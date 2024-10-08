package cs

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
	TechTagDefense            TechTag = "Defense"
	TechTagEngine             TechTag = "Engine"
	TechTagFuelTank           TechTag = "FuelTank"
	TechTagGatlingGun         TechTag = "GatlingGun"
	TechTagHeavyMineLayer     TechTag = "HeavyMineLayer"
	TechTagInitiativeBonus    TechTag = "InitiativeBonus"
	TechTagMassDriver         TechTag = "MassDriver"
	TechTagManeuveringJet     TechTag = "ManeuveringJet"
	TechTagMineLayer          TechTag = "MineLayer"
	TechTagMiningRobot        TechTag = "MiningRobot"
	TechTagPlanetaryScanner   TechTag = "PlanetaryScanner"
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

/* returns true if tt has ALL of the specified tags
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
}*/

// return unsorted list of all tags in tt
func (tt TechTags) GetTags() []TechTag {
	var list []TechTag
	for k, v := range tt {
		if v {
			list = append(list, k)
		}
	}
	return list
}

// return number of unique tags in tt
func (tt TechTags) CountTags() int {
	count := 0
	checkedTags := TechTags{}
	for k, v := range tt {
		if v && !checkedTags[k] {
			count += 1
			checkedTags[k] = true
		}
	}
	return count
}
