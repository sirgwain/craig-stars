package cs


type PRTSpec struct {
	PRT                              PRT              `json:"prt,omitempty"`
	PointCost                        int              `json:"pointCost,omitempty"`
	StartingTechLevels               TechLevel        `json:"startingTechLevels,omitempty"`
	StartingPlanets                  []StartingPlanet `json:"startingPlanets,omitempty"`
	TechCostOffset                   TechCostOffset   `json:"techCostOffset,omitempty"`
	MineralsPerSingleMineralPacket   int              `json:"mineralsPerSingleMineralPacket,omitempty"`
	MineralsPerMixedMineralPacket    int              `json:"mineralsPerMixedMineralPacket,omitempty"`
	PacketResourceCost               int              `json:"packetResourceCost,omitempty"`
	PacketMineralCostFactor          float64          `json:"packetMineralCostFactor,omitempty"`
	PacketReceiverFactor             float64          `json:"packetReceiverFactor,omitempty"`
	PacketDecayFactor                float64          `json:"packetDecayFactor,omitempty"`
	PacketOverSafeWarpPenalty        int              `json:"packetOverSafeWarpPenalty,omitempty"`
	PacketBuiltInScanner             bool             `json:"packetBuiltInScanner,omitempty"`
	DetectPacketDestinationStarbases bool             `json:"detectPacketDestinationStarbases,omitempty"`
	DetectAllPackets                 bool             `json:"detectAllPackets,omitempty"`
	PacketTerraformChance            float64          `json:"packetTerraformChance,omitempty"`
	PacketPermaformChance            float64          `json:"packetPermaformChance,omitempty"`
	PacketPermaTerraformSizeUnit     int              `json:"packetPermaTerraformSizeUnit,omitempty"`
	CanGateCargo                     bool             `json:"canGateCargo,omitempty"`
	CanDetectStargatePlanets         bool             `json:"canDetectStargatePlanets,omitempty"`
	ShipsVanishInVoid                bool             `json:"shipsVanishInVoid,omitempty"`
	BuiltInScannerMultiplier         int              `json:"builtInScannerMultiplier,omitempty"`
	TechsCostExtraLevel              int              `json:"techsCostExtraLevel,omitempty"`
	FreighterGrowthFactor            float64          `json:"freighterGrowthFactor,omitempty"`
	GrowthFactor                     float64          `json:"growthFactor,omitempty"`
	MaxPopulationOffset              float64          `json:"maxPopulationOffset,omitempty"`
	BuiltInCloakUnits                int              `json:"builtInCloakUnits,omitempty"`
	StealsResearch                   StealsResearch   `json:"stealsResearch,omitempty"`
	FreeCargoCloaking                bool             `json:"freeCargoCloaking,omitempty"`
	MineFieldsAreScanners            bool             `json:"mineFieldsAreScanners,omitempty"`
	MineFieldRateMoveFactor          float64          `json:"mineFieldRateMoveFactor,omitempty"`
	MineFieldSafeWarpBonus           int              `json:"mineFieldSafeWarpBonus,omitempty"`
	MineFieldMinDecayFactor          float64          `json:"mineFieldMinDecayFactor,omitempty"`
	MineFieldBaseDecayRate           float64          `json:"mineFieldBaseDecayRate,omitempty"`
	MineFieldPlanetDecayRate         float64          `json:"mineFieldPlanetDecayRate,omitempty"`
	MineFieldMaxDecayRate            float64          `json:"mineFieldMaxDecayRate,omitempty"`
	CanDetonateMineFields            bool             `json:"canDetonateMineFields,omitempty"`
	MineFieldDetonateDecayRate       float64          `json:"mineFieldDetonateDecayRate,omitempty"`
	DiscoverDesignOnScan             bool             `json:"discoverDesignOnScan,omitempty"`
	CanRemoteMineOwnPlanets          bool             `json:"canRemoteMineOwnPlanets,omitempty"`
	InvasionAttackBonus              float64          `json:"invasionAttackBonus,omitempty"`
	InvasionDefendBonus              float64          `json:"invasionDefendBonus,omitempty"`
	MovementBonus                    int              `json:"movementBonus,omitempty"`
	Instaforming                     bool             `json:"instaforming,omitempty"`
	PermaformChance                  float64          `json:"permaformChance,omitempty"`
	PermaformPopulation              int              `json:"permaformPopulation,omitempty"`
	RepairFactor                     float64          `json:"repairFactor,omitempty"`
	StarbaseRepairFactor             float64          `json:"starbaseRepairFactor,omitempty"`
	StarbaseCostFactor               float64          `json:"starbaseCostFactor,omitempty"`
	InnateMining                     bool             `json:"innateMining,omitempty"`
	InnateResources                  bool             `json:"innateResources,omitempty"`
	InnateScanner                    bool             `json:"innateScanner,omitempty"`
	InnatePopulationFactor           float64          `json:"innatePopulationFactor,omitempty"`
	CanBuildDefenses                 bool             `json:"canBuildDefenses,omitempty"`
	LivesOnStarbases                 bool             `json:"livesOnStarbases,omitempty"`
}

type LRTSpec struct {
	LRT                           LRT             `json:"lrt,omitempty"`
	StartingFleets                []StartingFleet `json:"startingFleets,omitempty"`
	PointCost                     int             `json:"pointCost,omitempty"`
	StartingTechLevels            TechLevel       `json:"startingTechLevels,omitempty"`
	TechCostOffset                TechCostOffset  `json:"techCostOffset,omitempty"`
	NewTechCostFactorOffset       float64         `json:"newTechCostFactorOffset,omitempty"`
	MiniaturizationMax            float64         `json:"miniaturizationMax,omitempty"`
	MiniaturizationPerLevel       float64         `json:"miniaturizationPerLevel,omitempty"`
	NoAdvancedScanners            bool            `json:"noAdvancedScanners,omitempty"`
	ScanRangeFactorOffset         float64         `json:"scanRangeFactorOffset,omitempty"`
	FuelEfficiencyOffset          float64         `json:"fuelEfficiencyOffset,omitempty"`
	MaxPopulationOffset           float64         `json:"maxPopulationOffset,omitempty"`
	TerraformCostOffset           Cost            `json:"terraformCostOffset,omitempty"`
	MineralAlchemyCostOffset      int             `json:"mineralAlchemyCostOffset,omitempty"`
	ScrapMineralOffset            float64         `json:"scrapMineralOffset,omitempty"`
	ScrapMineralOffsetStarbase    float64         `json:"scrapMineralOffsetStarbase,omitempty"`
	ScrapResourcesOffset          float64         `json:"scrapResourcesOffset,omitempty"`
	ScrapResourcesOffsetStarbase  float64         `json:"scrapResourcesOffsetStarbase,omitempty"`
	StartingPopulationFactorDelta float64         `json:"startingPopulationFactorDelta,omitempty"`
	StarbaseBuiltInCloakUnits     int             `json:"starbaseBuiltInCloakUnits,omitempty"`
	StarbaseCostFactorOffset      float64         `json:"starbaseCostFactorOffset,omitempty"`
	ResearchFactorOffset          float64         `json:"researchFactorOffset,omitempty"`
	ResearchSplashDamage          float64         `json:"researchSplashDamage,omitempty"`
	ShieldStrengthFactorOffset    float64         `json:"shieldStrengthFactorOffset,omitempty"`
	ShieldRegenerationRateOffset  float64         `json:"shieldRegenerationRateOffset,omitempty"`
	ArmorStrengthFactorOffset     float64         `json:"armorStrengthFactorOffset,omitempty"`
	EngineFailureRateOffset       float64         `json:"engineFailureRateOffset,omitempty"`
	EngineReliableSpeed           int             `json:"engineReliableSpeed,omitempty"`
}

type TechCostOffset struct {
	Engine           float64 `json:"engine,omitempty"`
	BeamWeapon       float64 `json:"beamWeapon,omitempty"`
	Torpedo          float64 `json:"torpedo,omitempty"`
	Bomb             float64 `json:"bomb,omitempty"`
	PlanetaryDefense float64 `json:"planetaryDefense,omitempty"`
}

type StartingPlanet struct {
	Population         int             `json:"population,omitempty"`
	HabPenaltyFactor   float64         `json:"habPenaltyFactor,omitempty"`
	HasStargate        bool            `json:"hasStargate,omitempty"`
	HasMassDriver      bool            `json:"hasMassDriver,omitempty"`
	StarbaseDesignName string          `json:"starbaseDesignName,omitempty"`
	StarbaseHull       string          `json:"starbaseHull,omitempty"`
	StartingFleets     []StartingFleet `json:"startingFleets,omitempty"`
}

type StartingFleet struct {
	Name          string            `json:"name,omitempty"`
	HullName      StartingFleetHull `json:"hullName,omitempty"`
	HullSetNumber uint              `json:"hullSetNumber,omitempty"`
	Purpose       ShipDesignPurpose `json:"purpose,omitempty"`
}

type StealsResearch struct {
	Energy        float64 `json:"energy,omitempty"`
	Weapons       float64 `json:"weapons,omitempty"`
	Propulsion    float64 `json:"propulsion,omitempty"`
	Construction  float64 `json:"construction,omitempty"`
	Electronics   float64 `json:"electronics,omitempty"`
	Biotechnology float64 `json:"biotechnology,omitempty"`
}

type StartingFleetHull string

const (
	StartingFleetHullColonyShip      StartingFleetHull = "Colony Ship"
	StartingFleetHullDestroyer       StartingFleetHull = "Destroyer"
	StartingFleetHullMediumFreighter StartingFleetHull = "Medium Freighter"
	StartingFleetHullMiniColonyShip  StartingFleetHull = "Mini-Colony Ship"
	StartingFleetHullMiniMineLayer   StartingFleetHull = "Mini Mine Layer"
	StartingFleetHullMiniMiner       StartingFleetHull = "Mini-Miner"
	StartingFleetHullMidgetMiner     StartingFleetHull = "Midget-Miner"
	StartingFleetHullPrivateer       StartingFleetHull = "Privateer"
	StartingFleetHullScout           StartingFleetHull = "Scout"
)

func defaultPRTSpec() PRTSpec {
	return PRTSpec{
		StartingPlanets: []StartingPlanet{{Population: 25000, StarbaseHull: SpaceStation.Name, StarbaseDesignName: "Starbase"}},

		PointCost:                        66,
		MineralsPerSingleMineralPacket:   100,
		MineralsPerMixedMineralPacket:    40,
		PacketResourceCost:               10,
		PacketMineralCostFactor:          1,
		PacketReceiverFactor:             1,
		PacketDecayFactor:                1,
		PacketOverSafeWarpPenalty:        0,
		PacketBuiltInScanner:             false,
		DetectPacketDestinationStarbases: false,
		DetectAllPackets:                 false,
		PacketTerraformChance:            0,
		PacketPermaformChance:            0,
		PacketPermaTerraformSizeUnit:     100,
		CanGateCargo:                     false,
		CanDetectStargatePlanets:         false,
		ShipsVanishInVoid:                true,
		BuiltInScannerMultiplier:         0,
		TechsCostExtraLevel:              3,
		FreighterGrowthFactor:            0,
		GrowthFactor:                     1,
		MaxPopulationOffset:              0,
		BuiltInCloakUnits:                0,
		FreeCargoCloaking:                false,
		MineFieldsAreScanners:            false,
		MineFieldRateMoveFactor:          0,
		MineFieldSafeWarpBonus:           0,
		MineFieldMinDecayFactor:          1,
		MineFieldBaseDecayRate:           .02,
		MineFieldPlanetDecayRate:         .04,
		MineFieldMaxDecayRate:            .5,
		CanDetonateMineFields:            false,
		MineFieldDetonateDecayRate:       .25,
		DiscoverDesignOnScan:             false,
		CanRemoteMineOwnPlanets:          false,
		InvasionAttackBonus:              1.1,
		InvasionDefendBonus:              1,
		MovementBonus:                    0,
		Instaforming:                     false,
		PermaformChance:                  0,
		PermaformPopulation:              0,
		RepairFactor:                     1,
		StarbaseRepairFactor:             1,
		StarbaseCostFactor:               1,
		InnateMining:                     false,
		InnateResources:                  false,
		InnateScanner:                    false,
		InnatePopulationFactor:           1,
		CanBuildDefenses:                 true,
		LivesOnStarbases:                 false,
	}
}

func heSpec() PRTSpec {
	spec := defaultPRTSpec()

	spec.StartingPlanets[0].StartingFleets = []StartingFleet{
		{"Deep Space Probe", StartingFleetHullScout, 0, ShipDesignPurposeScout},
		{"Spore Cloud", StartingFleetHullMiniColonyShip, 0, ShipDesignPurposeColonizer},
		{"Spore Cloud", StartingFleetHullMiniColonyShip, 0, ShipDesignPurposeColonizer},
		{"Spore Cloud", StartingFleetHullMiniColonyShip, 0, ShipDesignPurposeColonizer},
	}

	spec.GrowthFactor = 2
	spec.MaxPopulationOffset = -.5
	return spec
}

func ssSpec() PRTSpec {
	spec := defaultPRTSpec()

	spec.StartingTechLevels = TechLevel{
		Electronics: 5,
	}

	spec.StartingPlanets[0].StartingFleets = []StartingFleet{
		{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
		{"Santa Maria", StartingFleetHullColonyShip, 0, ShipDesignPurposeColonizer},
	}

	spec.BuiltInCloakUnits = 300
	spec.FreeCargoCloaking = true
	spec.MineFieldSafeWarpBonus = 1
	spec.StealsResearch = StealsResearch{
		Energy:        .5,
		Weapons:       .5,
		Propulsion:    .5,
		Construction:  .5,
		Electronics:   .5,
		Biotechnology: .5,
	}

	return spec
}

func wmSpec() PRTSpec {
	spec := defaultPRTSpec()

	spec.StartingTechLevels = TechLevel{
		Energy:     1,
		Weapons:    6,
		Propulsion: 1,
	}

	spec.StartingPlanets[0].StartingFleets = []StartingFleet{
		{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
		{"Santa Maria", StartingFleetHullColonyShip, 0, ShipDesignPurposeColonizer},
		{"Armored Probe", StartingFleetHullScout, 1, ShipDesignPurposeFighterScout},
	}

	spec.TechCostOffset = TechCostOffset{
		BeamWeapon: -.25,
		Torpedo:    -.25,
		Bomb:       -.25,
	}
	spec.DiscoverDesignOnScan = true
	spec.InvasionAttackBonus = 1.65
	spec.MovementBonus = 2

	return spec
}

func caSpec() PRTSpec {
	spec := defaultPRTSpec()

	spec.StartingTechLevels = TechLevel{
		Energy:        1,
		Weapons:       1,
		Propulsion:    1,
		Construction:  2,
		Biotechnology: 6,
	}

	spec.StartingPlanets[0].StartingFleets = []StartingFleet{
		{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
		{"Santa Maria", StartingFleetHullColonyShip, 0, ShipDesignPurposeColonizer},
		{"Change of Heart", StartingFleetHullMiniMiner, 1, ShipDesignPurposeTerraformer},
	}

	spec.Instaforming = true
	spec.PermaformChance = .1 // chance is 10% if pop is over 100k
	spec.PermaformPopulation = 100_000

	return spec
}

func isSpec() PRTSpec {
	spec := defaultPRTSpec()

	spec.StartingPlanets[0].StartingFleets = []StartingFleet{
		{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
		{"Santa Maria", StartingFleetHullColonyShip, 0, ShipDesignPurposeColonizer},
	}

	spec.TechCostOffset = TechCostOffset{
		PlanetaryDefense: -.4, // defenses cost 40% less
		BeamWeapon:       .25, // weapons cost 25% less
		Torpedo:          .25, // weapons cost 25% less
		Bomb:             .25, // weapons cost 25% less
	}

	spec.FreighterGrowthFactor = .5
	spec.InvasionDefendBonus = 2
	spec.RepairFactor = 2 // double repairs!
	spec.StarbaseRepairFactor = 1.5

	return spec
}

func sdSpec() PRTSpec {
	spec := defaultPRTSpec()

	spec.StartingTechLevels = TechLevel{
		Propulsion:    2,
		Biotechnology: 2,
	}

	spec.StartingPlanets[0].StartingFleets = []StartingFleet{
		{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
		{"Santa Maria", StartingFleetHullColonyShip, 0, ShipDesignPurposeColonizer},
		{"Little Hen", StartingFleetHullMiniMineLayer, 0, ShipDesignPurposeDamageMineLayer},
		{"Speed Turtle", StartingFleetHullMiniMineLayer, 0, ShipDesignPurposeSpeedMineLayer},
	}

	spec.MineFieldsAreScanners = true
	spec.CanDetonateMineFields = true
	spec.MineFieldRateMoveFactor = .5
	spec.MineFieldMinDecayFactor = .25
	spec.MineFieldSafeWarpBonus = 2

	return spec
}

func ppSpec() PRTSpec {
	spec := defaultPRTSpec()

	spec.StartingTechLevels = TechLevel{Energy: 4}

	spec.StartingPlanets = []StartingPlanet{
		// one homeworld, 20k people, no hab penalty
		{Population: 25000, HabPenaltyFactor: 0, HasMassDriver: true, StarbaseHull: SpaceStation.Name, StarbaseDesignName: "Starbase",
			StartingFleets: []StartingFleet{
				{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
				{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
				{"Santa Maria", StartingFleetHullColonyShip, 0, ShipDesignPurposeColonizer},
			},
		},
		// on extra world where hab varies by 1/2 of the range
		{
			Population: 10000, HabPenaltyFactor: 1, HasMassDriver: true, StarbaseHull: OrbitalFort.Name, StarbaseDesignName: "Accelerator Platform",
			StartingFleets: []StartingFleet{
				{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
			},
		},
	}
	spec.MineralsPerSingleMineralPacket = 70
	spec.MineralsPerMixedMineralPacket = 25
	spec.PacketResourceCost = 5
	spec.PacketMineralCostFactor = 1
	spec.PacketDecayFactor = .5
	spec.PacketBuiltInScanner = true
	spec.DetectPacketDestinationStarbases = true
	spec.DetectAllPackets = true
	spec.PacketTerraformChance = .5   // 50% per 100kT uncaught
	spec.PacketPermaformChance = .001 // .1% per 100kT uncaught

	return spec
}

func itSpec() PRTSpec {
	spec := defaultPRTSpec()

	spec.StartingTechLevels = TechLevel{
		Propulsion:   5,
		Construction: 5,
	}

	spec.StartingPlanets = []StartingPlanet{
		// one homeworld, 20k people, no hab penalty
		{Population: 25000, HabPenaltyFactor: 0, HasStargate: true, StarbaseHull: SpaceStation.Name, StarbaseDesignName: "Starbase",
			StartingFleets: []StartingFleet{
				{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
				{"Santa Maria", StartingFleetHullColonyShip, 0, ShipDesignPurposeColonizer},
				{"Swashbuckler", StartingFleetHullPrivateer, 0, ShipDesignPurposeArmedFreighter},
				{"Stalwart Defender", StartingFleetHullDestroyer, 0, ShipDesignPurposeFighter},
			},
		},
		// on extra world where hab varies by 1/2 of the range
		{
			Population: 10000, HabPenaltyFactor: 1, HasStargate: true, StarbaseHull: OrbitalFort.Name, StarbaseDesignName: "Accelerator Platform",
			StartingFleets: []StartingFleet{
				{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
			},
		},
	}

	spec.CanGateCargo = true
	spec.CanDetectStargatePlanets = true
	spec.ShipsVanishInVoid = false
	spec.PacketMineralCostFactor = 1.2
	spec.PacketReceiverFactor = .5
	spec.PacketOverSafeWarpPenalty = 1

	return spec
}

func arSpec() PRTSpec {
	spec := defaultPRTSpec()

	spec.StartingTechLevels = TechLevel{
		Energy: 1,
	}

	spec.StartingPlanets[0].StartingFleets = []StartingFleet{
		{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
		{"Santa Maria", StartingFleetHullColonyShip, 0, ShipDesignPurposeColonizer},
	}

	spec.CanRemoteMineOwnPlanets = true
	spec.StarbaseCostFactor = .8
	spec.InnateMining = true
	spec.InnateResources = true
	spec.InnateScanner = true
	spec.InnatePopulationFactor = .1
	spec.CanBuildDefenses = false
	spec.LivesOnStarbases = true

	return spec
}

func joatSpec() PRTSpec {
	spec := defaultPRTSpec()

	spec.StartingTechLevels = TechLevel{
		Energy:        3,
		Weapons:       3,
		Propulsion:    3,
		Construction:  3,
		Electronics:   3,
		Biotechnology: 3,
	}
	spec.StartingPlanets[0].StartingFleets = []StartingFleet{
		{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
		{"Santa Maria", StartingFleetHullColonyShip, 0, ShipDesignPurposeColonizer},
		{"Teamster", StartingFleetHullMediumFreighter, 0, ShipDesignPurposeFreighter},
		{"Cotton Picker", StartingFleetHullMiniMiner, 0, ShipDesignPurposeMiner},
		{"Armored Probe", StartingFleetHullScout, 1, ShipDesignPurposeFighterScout},
		{"Stalwart Defender", StartingFleetHullDestroyer, 0, ShipDesignPurposeFighter},
	}

	spec.MaxPopulationOffset = .2
	spec.BuiltInScannerMultiplier = 20
	spec.TechsCostExtraLevel = 4
	return spec
}

func ifeSpec() LRTSpec {
	return LRTSpec{
		StartingTechLevels:   TechLevel{Propulsion: 1},
		FuelEfficiencyOffset: -.15,
	}
}

func ttSpec() LRTSpec {
	return LRTSpec{
		TerraformCostOffset: Cost{Resources: -30},
	}
}

func armSpec() LRTSpec {
	spec := LRTSpec{}
	spec.StartingFleets = []StartingFleet{
		{"Potato Bug", StartingFleetHullMidgetMiner, 0, ShipDesignPurposeMiner},
		{"Potato Bug", StartingFleetHullMidgetMiner, 0, ShipDesignPurposeMiner},
	}

	return spec
}

func isbSpec() LRTSpec {
	return LRTSpec{
		StarbaseBuiltInCloakUnits: 40,  // 20% built in cloaking
		StarbaseCostFactorOffset:  -.2, // starbases cost 20% less
	}
}

func grSpec() LRTSpec {
	return LRTSpec{
		ResearchFactorOffset: -.5, // research is 50% less effective
		ResearchSplashDamage: .15, // our research applies 15% to all other fields
	}
}

func urSpec() LRTSpec {
	return LRTSpec{
		// UR gives us 45% of scrapped minerals and resources, versus 1/3 for races without UR
		ScrapMineralOffset:           .45 - (1.0 / 3),
		ScrapMineralOffsetStarbase:   .9 - (1.0 / 3),
		ScrapResourcesOffset:         .35,
		ScrapResourcesOffsetStarbase: .7,
	}
}

func nrseSpec() LRTSpec {
	return LRTSpec{}
}

func obrmSpec() LRTSpec {
	return LRTSpec{
		MaxPopulationOffset: .1,
	}
}

func nasSpec() LRTSpec {
	return LRTSpec{
		NoAdvancedScanners:    true,
		ScanRangeFactorOffset: 1, // add one to scan range factor scan range is doubled
	}
}

func lspSpec() LRTSpec {
	return LRTSpec{
		StartingPopulationFactorDelta: -.3, // start with 30% fewer colonists
	}
}

func betSpec() LRTSpec {
	return LRTSpec{
		NewTechCostFactorOffset: 1,   // new techs cost twice as much
		MiniaturizationMax:      .05, // go from 75% to 80% miniaturization max
		MiniaturizationPerLevel: .01, // go from 4% per level to 5% per level
	}
}

func rsSpec() LRTSpec {
	return LRTSpec{
		ShieldStrengthFactorOffset:   .4,
		ShieldRegenerationRateOffset: .1,
		ArmorStrengthFactorOffset:    -.5,
	}
}

func maSpec() LRTSpec {
	return LRTSpec{
		MineralAlchemyCostOffset: -75,
	}
}

func ceSpec() LRTSpec {
	return LRTSpec{
		StartingTechLevels: TechLevel{Propulsion: 1},

		TechCostOffset: TechCostOffset{
			Engine: -.5, // engines cost 50% less
		},

		EngineFailureRateOffset: .1,
		EngineReliableSpeed:     6,
	}
}
