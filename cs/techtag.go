package cs

// TechTags are functional labels to categorize tech items
// based on their function.
// They are used by the game when determining cost discounts,
// as well as for categorizing parts during ship designing.
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

// (This is in a separate file for tygo reasons)
