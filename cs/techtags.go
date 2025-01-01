package cs

// TechTags are functional labels used to
// categorize tech items based on their function.
type TechTag string

const (
	TechTagNone               TechTag = "None"
	TechTagArmor              TechTag = "Armor"
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
	TechTagTorpedoBonus       TechTag = "TorpedoBonus"
)

// a list of all TechTags that benefit ships in combat
var CombatTechTags = []TechTag{
	TechTagArmor,
	TechTagBeamCapacitor,
	TechTagBeamDeflector,
	TechTagBeamWeapon,
	TechTagCapitalShipMissile,
	TechTagEngine,
	TechTagGatlingGun,
	TechTagInitiativeBonus,
	TechTagManeuveringJet,
	TechTagShieldSapper,
	TechTagShield,
	TechTagTorpedo,
	TechTagTorpedoJammer,
	TechTagTorpedoBonus,
}

// A collection of an object's TechTags (like on a tech part)
type TechTags map[TechTag]bool

// Create a new TechTags map from a list of TechTag items, or an empty map if none are specified
func newTechTags(tags ...TechTag) TechTags {
	newTechTags := TechTags{}
	for _, t := range tags {
		newTechTags[t] = true
	}
	return newTechTags
}

// returns true if tt has at least 1 of the specified TechTags
// and none of the tags in tagsToExclude
//
// Tag(s) contained in both lists will not be banned
func (tt TechTags) hasTags(tagsToInclude []TechTag, tagsToExclude ...TechTag) bool {
	blacklist := newTechTags(tagsToExclude...)
	whitelist := newTechTags(tagsToInclude...)
	hasTag := false

	for _, tag := range tt.GetTags() {
		if _, ok := whitelist[tag]; ok {
			hasTag = true
		} else if _, ok := blacklist[tag]; ok {
			// our TechTags has a blacklisted tag not
			// also in our whitelist; automatic fail
			return false
		}
	}
	return hasTag
}

// return true if tt has this tag
func (tt TechTags) HasTag(tag TechTag) bool {
	_, ok := tt[tag]
	return ok
}

// return unsorted list of all tags in tt
func (tt TechTags) GetTags() []TechTag {
	list := make([]TechTag, 0, len(tt))
	for k := range tt {
		list = append(list, k)
	}
	return list
}

// return number of unique tags in tt
func (tt TechTags) Count() int {
	return len(tt)
}
