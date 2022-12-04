package game

type PRTSpec struct {
	PRT                              PRT              `json:"prt"`
	PointCost                        int              `json:"pointCost"`
	StartingTechLevels               TechLevel        `json:"startingTechLevels"`
	StartingPlanets                  []StartingPlanet `json:"startingPlanets"`
	TechCostOffset                   TechCostOffset   `json:"techCostOffset"`
	MineralsPerSingleMineralPacket   int              `json:"mineralsPerSingleMineralPacket"`
	MineralsPerMixedMineralPacket    int              `json:"mineralsPerMixedMineralPacket"`
	PacketResourceCost               int              `json:"packetResourceCost"`
	PacketMineralCostFactor          float64          `json:"packetMineralCostFactor"`
	PacketReceiverFactor             float64          `json:"packetReceiverFactor"`
	PacketDecayFactor                float64          `json:"packetDecayFactor"`
	PacketOverSafeWarpPenalty        int              `json:"packetOverSafeWarpPenalty"`
	PacketBuiltInScanner             bool             `json:"packetBuiltInScanner"`
	DetectPacketDestinationStarbases bool             `json:"detectPacketDestinationStarbases"`
	DetectAllPackets                 bool             `json:"detectAllPackets"`
	PacketTerraformChance            float64          `json:"packetTerraformChance"`
	PacketPermaformChance            float64          `json:"packetPermaformChance"`
	PacketPermaTerraformSizeUnit     int              `json:"packetPermaTerraformSizeUnit"`
	CanGateCargo                     bool             `json:"canGateCargo"`
	CanDetectStargatePlanets         bool             `json:"canDetectStargatePlanets"`
	ShipsVanishInVoid                bool             `json:"shipsVanishInVoid"`
	BuiltInScannerMultiplier         int              `json:"builtInScannerMultiplier"`
	TechsCostExtraLevel              int              `json:"techsCostExtraLevel"`
	FreighterGrowthFactor            float64          `json:"freighterGrowthFactor"`
	GrowthFactor                     float64          `json:"growthFactor"`
	MaxPopulationOffset              float64          `json:"maxPopulationOffset"`
	BuiltInCloakUnits                int              `json:"builtInCloakUnits"`
	StealsResearch                   TechLevel        `json:"stealsResearch"`
	FreeCargoCloaking                bool             `json:"freeCargoCloaking"`
	MineFieldsAreScanners            bool             `json:"mineFieldsAreScanners"`
	MineFieldRateMoveFactor          float64          `json:"mineFieldRateMoveFactor"`
	MineFieldSafeWarpBonus           int              `json:"mineFieldSafeWarpBonus"`
	MineFieldMinDecayFactor          float64          `json:"mineFieldMinDecayFactor"`
	MineFieldBaseDecayRate           float64          `json:"mineFieldBaseDecayRate"`
	MineFieldPlanetDecayRate         float64          `json:"mineFieldPlanetDecayRate"`
	MineFieldMaxDecayRate            float64          `json:"mineFieldMaxDecayRate"`
	CanDetonateMineFields            bool             `json:"canDetonateMineFields"`
	MineFieldDetonateDecayRate       float64          `json:"mineFieldDetonateDecayRate"`
	DiscoverDesignOnScan             bool             `json:"discoverDesignOnScan"`
	CanRemoteMineOwnPlanets          bool             `json:"canRemoteMineOwnPlanets"`
	InvasionAttackBonus              float64          `json:"invasionAttackBonus"`
	InvasionDefendBonus              float64          `json:"invasionDefendBonus"`
	MovementBonus                    int              `json:"movementBonus"`
	Instaforming                     bool             `json:"instaforming"`
	PermaformChance                  float64          `json:"permaformChance"`
	PermaformPopulation              int              `json:"permaformPopulation"`
	RepairFactor                     float64          `json:"repairFactor"`
	StarbaseRepairFactor             float64          `json:"starbaseRepairFactor"`
	StarbaseCostFactor               float64          `json:"starbaseCostFactor"`
	InnateMining                     bool             `json:"innateMining"`
	InnateResources                  bool             `json:"innateResources"`
	InnateScanner                    bool             `json:"innateScanner"`
	InnatePopulationFactor           float64          `json:"innatePopulationFactor"`
	CanBuildDefenses                 bool             `json:"canBuildDefenses"`
	LivesOnStarbases                 bool             `json:"livesOnStarbases"`
}

type LRTSpec struct {
	LRT                          LRT            `json:"lrt"`
	PointCost                    int            `json:"pointCost"`
	StartingTechLevels           TechLevel      `json:"startingTechLevels"`
	TechCostOffset               TechCostOffset `json:"techCostOffset"`
	NewTechCostFactor            float64        `json:"newTechCostFactor"`
	MiniaturizationMax           float64        `json:"miniaturizationMax"`
	MiniaturizationPerLevel      float64        `json:"miniaturizationPerLevel"`
	NoAdvancedScanners           bool           `json:"noAdvancedScanners"`
	ScanRangeFactor              float64        `json:"scanRangeFactor"`
	FuelEfficiencyOffset         float64        `json:"fuelEfficiencyOffset"`
	MaxPopulationOffset          float64        `json:"maxPopulationOffset"`
	TerraformCostOffset          Cost           `json:"terraformCostOffset"`
	MineralAlchemyCostOffset     int            `json:"mineralAlchemyCostOffset"`
	ScrapMineralOffset           float64        `json:"scrapMineralOffset"`
	ScrapMineralOffsetStarbase   float64        `json:"scrapMineralOffsetStarbase"`
	ScrapResourcesOffset         float64        `json:"scrapResourcesOffset"`
	ScrapResourcesOffsetStarbase float64        `json:"scrapResourcesOffsetStarbase"`
	StartingPopulationFactor     float64        `json:"startingPopulationFactor"`
	StarbaseBuiltInCloakUnits    int            `json:"starbaseBuiltInCloakUnits"`
	StarbaseCostFactor           float64        `json:"starbaseCostFactor"`
	ResearchFactor               float64        `json:"researchFactor"`
	ResearchSplashDamage         float64        `json:"researchSplashDamage"`
	ShieldStrengthFactor         float64        `json:"shieldStrengthFactor"`
	ShieldRegenerationRate       float64        `json:"shieldRegenerationRate"`
	EngineFailureRate            float64        `json:"engineFailureRate"`
	EngineReliableSpeed          int            `json:"engineReliableSpeed"`
}

type TechCostOffset struct {
	Engine           float64 `json:"engine,omitempty"`
	BeamWeapon       float64 `json:"beamWeapon,omitempty"`
	Torpedo          float64 `json:"torpedo,omitempty"`
	Bomb             float64 `json:"bomb,omitempty"`
	PlanetaryDefense float64 `json:"planetaryDefense,omitempty"`
}

type StartingPlanet struct {
	Population         int             `json:"population"`
	HabPenaltyFactor   float64         `json:"habPenaltyFactor"`
	HasStargate        bool            `json:"hasStargate"`
	HasMassDriver      bool            `json:"hasMassDriver"`
	StarbaseDesignName string          `json:"starbaseDesignName"`
	StarbaseHull       string          `json:"starbaseHull"`
	StartingFleets     []StartingFleet `json:"startingFleets"`
}

type StartingFleet struct {
	Name          string            `json:"name"`
	HullName      StartingFleetHull `json:"hullName"`
	HullSetNumber uint              `json:"hullSetNumber"`
	Purpose       ShipDesignPurpose `json:"purpose"`
}

type StartingFleetHull string

const (
	StartingFleetHullColonyShip      StartingFleetHull = "Colony Ship"
	StartingFleetHullDestroyer       StartingFleetHull = "Destroyer"
	StartingFleetHullMediumFreighter StartingFleetHull = "Medium Freighter"
	StartingFleetHullMiniColonyShip  StartingFleetHull = "Mini-Colony Ship"
	StartingFleetHullMiniMineLayer   StartingFleetHull = "Mini Mine Layer"
	StartingFleetHullMiniMiner       StartingFleetHull = "Mini-Miner"
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
		InvasionAttackBonus:              1,
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

	return spec
}

func ssSpec() PRTSpec {
	spec := defaultPRTSpec()

	return spec
}

func wmSpec() PRTSpec {
	spec := defaultPRTSpec()

	return spec
}

func caSpec() PRTSpec {
	spec := defaultPRTSpec()

	return spec
}

func isSpec() PRTSpec {
	spec := defaultPRTSpec()

	return spec
}

func sdSpec() PRTSpec {
	spec := defaultPRTSpec()

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

	return spec
}

func arSpec() PRTSpec {
	spec := defaultPRTSpec()

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
		{"Stalwart Defender", StartingFleetHullDestroyer, 0, ShipDesignPurposeFighterScout},
	}

	spec.MaxPopulationOffset = .2
	spec.BuiltInScannerMultiplier = 20
	spec.TechsCostExtraLevel = 4
	return spec
}

func ifeSpec() LRTSpec {
	return LRTSpec{}
}

func ttSpec() LRTSpec {
	return LRTSpec{}
}

func armSpec() LRTSpec {
	return LRTSpec{}
}

func isbSpec() LRTSpec {
	return LRTSpec{}
}

func grSpec() LRTSpec {
	return LRTSpec{}
}

func urSpec() LRTSpec {
	return LRTSpec{}
}

func nrseSpec() LRTSpec {
	return LRTSpec{}
}

func obrmSpec() LRTSpec {
	return LRTSpec{}
}

func nasSpec() LRTSpec {
	return LRTSpec{}
}

func lspSpec() LRTSpec {
	return LRTSpec{}
}

func betSpec() LRTSpec {
	return LRTSpec{}
}

func rsSpec() LRTSpec {
	return LRTSpec{}
}

func maSpec() LRTSpec {
	return LRTSpec{}
}

func ceSpec() LRTSpec {
	return LRTSpec{}
}
