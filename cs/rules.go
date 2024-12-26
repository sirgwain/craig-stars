package cs

import (
	"fmt"
	"math/rand"
	"time"
)

// The Rules struct contains all the various constants and configuration values that determine
// how the game mechanics work. These are designed to be unique per game, if desired. Currently for testing, all
// games just use the default rule set.
type Rules struct {
	CostRules
	BattleRules
	UniverseGenerationRules
	ID                               int64                               `json:"id,omitempty"`
	CreatedAt                        time.Time                           `json:"createdAt,omitempty"`
	UpdatedAt                        time.Time                           `json:"updatedAt,omitempty"`
	GameID                           int64                               `json:"gameId,omitempty"`
	AcquirablePartTradeChanceBase    float64                             `json:"acquirablePartTradeChanceBase,omitempty"`
	AcquirablePartTradeItemMax       int                                 `json:"acquirablePartTradeItemMax,omitempty"`
	CometStatsBySize                 map[CometSize]CometStats            `json:"cometStatsBySize,omitempty"`
	FleetSafeSpeedExplosionChance    float64                             `json:"fleetSafeSpeedExplosionChance"`
	InvasionDefenseCoverageFactor    float64                             `json:"invasionDefenseCoverageFactor"`
	LRTSpecs                         map[LRT]LRTSpec                     `json:"lrtSpecs"`
	MaxPopulation                    int                                 `json:"maxPopulation"`
	MaxTechLevel                     int                                 `json:"maxTechLevel"`
	MineFieldCloak                   int                                 `json:"mineFieldCloak"`
	MineFieldStatsByType             map[MineFieldType]MineFieldStats    `json:"mineFieldStatsByType"`
	MineralDecayFactor               int                                 `json:"mineralDecayFactor"`
	MinMaxPopulationPercent          float64                             `json:"minMaxPopulationPercent"`
	MysteryTraderRules               MysteryTraderRules                  `json:"mysteryTraderRules"`
	PacketDecayRate                  map[int]float64                     `json:"packetDecayRate"`
	PacketMaxOverwarpSpeed           int                                 `json:"packetMaxOverwarpSpeed"`
	PacketMinDecay                   int                                 `json:"packetMinDecay"`
	PlanetMinDistance                int                                 `json:"planetMinDistance"`
	PopulationOvercrowdDieoffRate    float64                             `json:"populationOvercrowdDieoffRate"`
	PopulationOvercrowdDieoffRateMax float64                             `json:"populationOvercrowdDieoffRateMax"`
	PopulationScannerError           float64                             `json:"populationScannerError"`
	PRTSpecs                         map[PRT]PRTSpec                     `json:"prtSpecs"`
	RaceStartingPoints               int                                 `json:"raceStartingPoints"`
	RadiatingImmune                  int                                 `json:"radiatingImmune"`
	RandomArtifactResearchBonusRange [2]int                              `json:"randomArtifactResearchBonusRange"`
	RandomCometMinYear               int                                 `json:"randomCometMinYear,omitempty"`
	RandomCometMinYearPlayerWorld    int                                 `json:"randomCometMinYearPlayerWorld,omitempty"`
	RandomEventChances               map[RandomEvent]float64             `json:"randomEventChances,omitempty"`
	RandomMineralDepositBonusRange   [2]int                              `json:"randomMineralDepositBonusRange,omitempty"`
	RemoteMiningMineOutput           int                                 `json:"remoteMiningMineOutput,omitempty"`
	RepairRates                      map[RepairRate]float64              `json:"repairRates,omitempty"`
	SalvageDecayMin                  int                                 `json:"salvageDecayMin,omitempty"`
	SalvageDecayRate                 float64                             `json:"salvageDecayRate,omitempty"`
	SalvageFromBattleFactor          float64                             `json:"salvageFromBattleFactor,omitempty"`
	ScrapColonizeAmount              float64                             `json:"scrapColonizeAmount,omitempty"`
	ScrapMineralAmount               float64                             `json:"scrapMineralAmount,omitempty"`
	ScrapResourceAmount              float64                             `json:"scrapResourceAmount,omitempty"`
	ShowPublicScoresAfterYears       int                                 `json:"showPublicScoresAfterYears,omitempty"`
	SmartDefenseCoverageFactor       float64                             `json:"smartDefenseCoverageFactor,omitempty"`
	StargateMaxHullMassFactor        int                                 `json:"stargateMaxHullMassFactor,omitempty"`
	StargateMaxRangeFactor           int                                 `json:"stargateMaxRangeFactor,omitempty"`
	TachyonCloakReduction            float64                             `json:"tachyonCloakReduction,omitempty"`
	TachyonMaxCloakReduction         float64                             `json:"tachyonMaxCloakReduction,omitempty"`
	TechsID                          int64                               `json:"techsId,omitempty"`
	TechTradeChance                  float64                             `json:"techTradeChance,omitempty"`
	TorpedoSplashDamage              float64                             `json:"torpedoSplashDamage,omitempty"`
	WormholeCloak                    int                                 `json:"wormholeCloak,omitempty"`
	WormholePairsForSize             map[Size]int                        `json:"wormholePairsForSize,omitempty"`
	WormholeStatsByStability         map[WormholeStability]WormholeStats `json:"wormholeStatsByStability,omitempty"`
	random                           rng
	techs                            *TechStore
}

type UniverseGenerationRules struct {
	HighRadMineralConcentrationBonusThreshold int                           `json:"highRadMineralConcentrationBonusThreshold,omitempty"`
	LimitMineralConcentration                 int                           `json:"limitMineralConcentration,omitempty"`
	MaxExtraWorldDistance                     int                           `json:"maxExtraWorldDistance,omitempty"`
	MaxHab                                    int                           `json:"maxHab,omitempty"`
	MaxMineralConcentration                   int                           `json:"maxMineralConcentration,omitempty"`
	MaxStartingMineralConcentration           int                           `json:"maxStartingMineralConcentration,omitempty"`
	MaxStartingMineralSurface                 int                           `json:"maxStartingMineralSurface,omitempty"`
	MinExtraPlanetMineralConcentration        int                           `json:"minExtraPlanetMineralConcentration,omitempty"`
	MinExtraWorldDistance                     int                           `json:"minExtraWorldDistance,omitempty"`
	MinHab                                    int                           `json:"minHab,omitempty"`
	MinHomeworldMineralConcentration          int                           `json:"minHomeworldMineralConcentration,omitempty"`
	MinMineralConcentration                   int                           `json:"minMineralConcentration,omitempty"`
	MinStartingMineralConcentration           int                           `json:"minStartingMineralConcentration,omitempty"`
	MinStartingMineralSurface                 int                           `json:"minStartingMineralSurface,omitempty"`
	RaceLeftoverPointsPerItem                 map[SpendLeftoverPointsOn]int `json:"raceLeftoverPointsPerItem,omitempty"`
	StartingYear                              int                           `json:"startingYear,omitempty"`
	WormholeMinPlanetDistance                 int                           `json:"wormholeMinPlanetDistance,omitempty"`
}

type CostRules struct {
	DefenseCost                    Cost    `json:"defenseCost,omitempty"`
	FactoryCostGermanium           int     `json:"factoryCostGermanium,omitempty"`
	MineralAlchemyCost             int     `json:"mineralAlchemyCost,omitempty"`
	PlanetaryScannerCost           Cost    `json:"planetaryScannerCost,omitempty"`
	StarbaseComponentCostReduction float64 `json:"starbaseComponentCostReduction,omitempty"`
	StarbaseHullRefundFactor       float64 `json:"starbaseHullRefundFactor,omitempty"`
	TerraformCost                  Cost    `json:"terraformCost,omitempty"`
	TechBaseCost                   []int   `json:"techBaseCost,omitempty"`
}

type BattleRules struct {
	BeamRangeDropoff    float64 `json:"beamRangeDropoff,omitempty"`
	BeamBonusCap        float64 `json:"beamBonusCap,omitempty"`
	JammerCap           BoolMap `json:"jammerCap,omitempty"`
	JammerMulti         BoolMap `json:"jammerMulti,omitempty"`
	MovementMin         int     `json:"movementMin,omitempty"`
	MovementMax         int     `json:"movementMax,omitempty"`
	MovesToRunAway      int     `json:"movesToRunAway,omitempty"`
	NumBattleRounds     int     `json:"numBattleRounds,omitempty"`
	TorpedoSplashDamage float64 `json:"torpedoSplashDamage,omitempty"`
}

type RandomEvent string

const (
	RandomEventComet           RandomEvent = "Comet"
	RandomEventMineralDeposit  RandomEvent = "MineralDeposit"
	RandomEventPlanetaryChange RandomEvent = "PlanetaryChange"
	RandomEventAncientArtifact RandomEvent = "AncientArtifact"
)

type CometSize string

const (
	CometSmall  CometSize = "Small"
	CometMedium CometSize = "Medium"
	CometLarge  CometSize = "Large"
	CometHuge   CometSize = "Huge"
)

var CometSizes = []CometSize{
	CometSmall,
	CometMedium,
	CometLarge,
	CometHuge,
}

// each type of comet has stats for minerals added to each mineral type
// as well as some additional mineral types that get bonuses
type CometStats struct {
	AllMinerals              int     `json:"allMinerals,omitempty"`
	AllRandomMinerals        int     `json:"allRandomMinerals,omitempty"`
	BonusMinerals            int     `json:"bonusMinerals,omitempty"`
	BonusRandomMinerals      int     `json:"bonusRandomMinerals,omitempty"`
	BonusMinConcentration    int     `json:"bonusMinConcentration,omitempty"`
	BonusRandomConcentration int     `json:"bonusRandomConcentration,omitempty"`
	BonusAffectsMinerals     int     `json:"bonusAffectsMinerals,omitempty"`
	MinTerraform             int     `json:"minTerraform,omitempty"`
	RandomTerraform          int     `json:"randomTerraform,omitempty"`
	AffectsHabs              int     `json:"affectsHabs,omitempty"`
	PopKilledPercent         float64 `json:"popKilledPercent,omitempty"`
}

type RepairRate string

const (
	RepairRateNone              RepairRate = "None"
	RepairRateMoving            RepairRate = "Moving"
	RepairRateStopped           RepairRate = "Stopped"
	RepairRateOrbiting          RepairRate = "Orbiting"
	RepairRateOrbitingOwnPlanet RepairRate = "OrbitingOwnPlanet"
	RepairRateStarbase          RepairRate = "Starbase"
)

type MysteryTraderRules struct {
	ChanceSpawn           []int                        `json:"chanceSpawn,omitempty"`
	ChanceMaxTechGetsPart int                          `json:"chanceMaxTechGetsPart,omitempty"`
	ChanceCourseChange    int                          `json:"chanceCourseChange,omitempty"`
	ChanceSpeedUpOnly     int                          `json:"chanceSpeedUpOnly,omitempty"`
	ChanceAgain           int                          `json:"chanceAgain,omitempty"`
	EvenYearOnly          bool                         `json:"evenYearOnly,omitempty"`
	GenesisDeviceCost     Cost                         `json:"genesisDeviceCost,omitempty"`
	MaxMysteryTraders     int                          `json:"maxMysteryTraders,omitempty"`
	MaxWarp               int                          `json:"maxWarp,omitempty"`
	MinWarp               int                          `json:"minWarp,omitempty"`
	MinYear               int                          `json:"minYear,omitempty"`
	RequestedBoon         int                          `json:"requestedBoon,omitempty"`
	TechBoon              []MysteryTraderTechBoonRules `json:"techBoon,omitempty"`
}

type MysteryTraderTechBoonRules struct {
	TechLevels int                                   `json:"techLevels,omitempty"`
	Rewards    []MysteryTraderTechBoonMineralsReward `json:"rewards,omitempty"`
}

type MysteryTraderTechBoonMineralsReward struct {
	MineralsGiven int `json:"mineralsGiven,omitempty"`
	Reward        int `json:"reward,omitempty"`
}

var StandardRules = NewRules()

// Seed the random number generator with the rules Seed value
// This should be called after deserializing
// This can be used to generate the same world repeatedly (hopefully)
func (r *Rules) ResetSeed(seed int64) {
	r.random = rand.New(rand.NewSource(seed))
}

func (r *Rules) SetTechStore(techStore *TechStore) *Rules {
	r.techs = techStore
	return r
}

func NewRules() Rules {
	// create the random number generator for these rules
	seed := time.Now().UnixNano()
	return NewRulesWithSeed(seed)
}

func NewRulesWithSeed(seed int64) Rules {
	random := rand.New(rand.NewSource(seed))

	return Rules{
		random: random,
		CostRules: CostRules{
			FactoryCostGermanium: 4,
			DefenseCost: Cost{
				Ironium:   5,
				Boranium:  5,
				Germanium: 5,
				Resources: 15,
			},
			MineralAlchemyCost: 100,
			PlanetaryScannerCost: Cost{
				Ironium:   10,
				Boranium:  10,
				Germanium: 70,
				Resources: 100,
			},
			StarbaseComponentCostReduction: 0.5, // 50% discount on non-orbital components
			StarbaseHullRefundFactor:       0.5, // 50% of the old base's cost goes towards the new base
			TechBaseCost: []int{
				0,
				50,
				80,
				130,
				210,
				340,
				550,
				890,
				1440,
				2330,
				3770,
				6100,
				9870,
				13850,
				18040,
				22440,
				27050,
				31870,
				36900,
				42140,
				47590,
				53250,
				59120,
				65200,
				71490,
				77990,
				84700,
			},
			TerraformCost: Cost{
				Ironium:   0,
				Boranium:  0,
				Germanium: 0,
				Resources: 100,
			},
		},
		BattleRules: BattleRules{
			BeamRangeDropoff: 0.1,
			BeamBonusCap:     2.55, // 2.55x damage max from caps
			JammerCap: BoolMap{
				valueIfTrue:  1,    // starbases have 100 jamming max, but an innate 0.75x jam penalty
				valueIfFalse: 0.95, // non-starbases (ie fleets) have 95% jamming max
			},
			JammerMulti: BoolMap{
				valueIfTrue:  0.75, // starbases have innate 0.75x jam penalty by default
				valueIfFalse: 1,    // non-starbases (ie fleets) have no penalty
			},
			MovementMin:         2,
			MovementMax:         10,
			MovesToRunAway:      7,
			NumBattleRounds:     16,
			TorpedoSplashDamage: 0.125,
		},
		UniverseGenerationRules: UniverseGenerationRules{
			MaxExtraWorldDistance:                     180,
			MinExtraWorldDistance:                     130,
			MinHomeworldMineralConcentration:          30,
			MinExtraPlanetMineralConcentration:        30,
			MinMineralConcentration:                   1,
			MaxMineralConcentration:                   200,
			MinHab:                                    1,
			MaxHab:                                    99,
			MinStartingMineralConcentration:           1,
			MaxStartingMineralConcentration:           121,
			LimitMineralConcentration:                 30,
			HighRadMineralConcentrationBonusThreshold: 90,
			MaxStartingMineralSurface:                 1000,
			MinStartingMineralSurface:                 300,
			RaceLeftoverPointsPerItem: map[SpendLeftoverPointsOn]int{
				SpendLeftoverPointsOnMines:                 2,
				SpendLeftoverPointsOnFactories:             5,
				SpendLeftoverPointsOnDefenses:              10,
				SpendLeftoverPointsOnMineralConcentrations: 3,
				SpendLeftoverPointsOnSurfaceMinerals:       10, // special case - indicates kT per point leftover
			},
			StartingYear:              2400,
			WormholeMinPlanetDistance: 30,
		},
		// TODO: Change tachyon cloak reduction to a property of the technology itself
		TachyonCloakReduction:            .05, // 5% diminishing cloak reduction per detector
		TachyonMaxCloakReduction:         .81, // tachyon detectors cap at 81% cloaking reduction
		MaxPopulation:                    1000000,
		MinMaxPopulationPercent:          .05, // red worlds have 5% max pop
		PopulationOvercrowdDieoffRate:    .04, // overcrowded pops die off at 4% per doubling
		PopulationOvercrowdDieoffRateMax: .12, // overcrowded pops will not die off more than 12% (3x pop) in a year
		PopulationScannerError:           0.2,
		SmartDefenseCoverageFactor:       0.5,
		InvasionDefenseCoverageFactor:    0.75,
		SalvageDecayRate:                 0.1,
		SalvageDecayMin:                  10,
		MineFieldCloak:                   75,
		StargateMaxRangeFactor:           5,
		StargateMaxHullMassFactor:        5,
		TechTradeChance:                  .5, // 50% chance of tech trading per level
		FleetSafeSpeedExplosionChance:    .1, // 10% chance of losing a ship
		RadiatingImmune:                  85, // hab center of > 85 are immune to radating damage
		RandomEventChances: map[RandomEvent]float64{
			RandomEventComet:           .05, // 1 in 20 chance of a planet being struck by a comet in a given turn
			RandomEventMineralDeposit:  .05,
			RandomEventPlanetaryChange: .05,
			RandomEventAncientArtifact: .33, // 1 in 3 planets have random artifacts
		},
		AcquirablePartTradeChanceBase: 0.005, // 0.5% chance per item in fleet
		AcquirablePartTradeItemMax:    25,    // 25 items max per trade instance
		RandomCometMinYear:            10,
		RandomCometMinYearPlayerWorld: 20,
		CometStatsBySize: map[CometSize]CometStats{
			CometSmall: {
				AllMinerals:              50, // adds 50 minerals to 300 minerals (>> 4) to all types
				AllRandomMinerals:        250,
				BonusMinerals:            3000, // adds (3000 to 20000) >> 4 bonus minerals
				BonusRandomMinerals:      17000,
				BonusMinConcentration:    50, // adds 50 to 100 mineral concentration
				BonusRandomConcentration: 50,
				BonusAffectsMinerals:     1,   // only one mineral gets a bonus + concentration
				MinTerraform:             3,   // terraforms by +/- 3 points
				RandomTerraform:          3,   // randomly terraforms by an additional +/- 3 points
				AffectsHabs:              1,   // terraforming affects one hab
				PopKilledPercent:         .25, // 25% pop killed
			},
			CometMedium: {
				AllMinerals:              50, // adds 50 minerals to 300 minerals (>> 4) to all types
				AllRandomMinerals:        250,
				BonusMinerals:            3000, // adds (3000 to 20000) >> 4 bonus minerals
				BonusRandomMinerals:      17000,
				BonusMinConcentration:    50, // adds 50 to 100 mineral concentration
				BonusRandomConcentration: 50,
				BonusAffectsMinerals:     2, // two minerals gets a bonus + concentration
				MinTerraform:             3,
				RandomTerraform:          3,
				AffectsHabs:              2, // terraforming affects two habs
				PopKilledPercent:         .45,
			},
			CometLarge: {
				AllMinerals:              50, // adds 50 minerals to 300 minerals (>> 4) to all types
				AllRandomMinerals:        250,
				BonusMinerals:            3000, // adds (3000 to 20000) >> 4 bonus minerals
				BonusRandomMinerals:      17000,
				BonusMinConcentration:    50, // adds 50 to 100 mineral concentration
				BonusRandomConcentration: 50,
				BonusAffectsMinerals:     3, // three minerals gets a bonus + concentration
				MinTerraform:             3,
				RandomTerraform:          3,
				AffectsHabs:              3, // terraforming affects three habs
				PopKilledPercent:         .65,
			},
			CometHuge: {
				AllMinerals:              50, // adds 50 minerals to 300 minerals (>> 4) to all types
				AllRandomMinerals:        250,
				BonusMinerals:            3000, // adds (3000 to 20000) >> 4 bonus minerals
				BonusRandomMinerals:      17000,
				BonusMinConcentration:    65, // adds 65 to 130 mineral concentration
				BonusRandomConcentration: 65,
				BonusAffectsMinerals:     3, // three minerals gets a bonus + concentration
				MinTerraform:             6, // terraforms 6 to 12 in a random direction
				RandomTerraform:          6,
				AffectsHabs:              3, // terraforming affects three habs
				PopKilledPercent:         .85,
			},
		},
		RandomMineralDepositBonusRange:   [2]int{20, 50},
		RandomArtifactResearchBonusRange: [2]int{120, 400},
		MysteryTraderRules: MysteryTraderRules{
			// ChanceSpawn:      []int{1}, // force it
			ChanceSpawn:           []int{7, 7, 7, 7, 7, 7, 7, 4, 4, 3, 2}, // randomly pick a random chance to spawn an MT. It's not the same every turn
			ChanceMaxTechGetsPart: 5,                                      // 1 in 5 chance a player with max tech gets a part if they get a research trader
			ChanceCourseChange:    20,                                     // 1 in 20 chance the MT speeds up/changes course
			ChanceSpeedUpOnly:     3,                                      // if change course, 1 in 3 chance it's speed up only
			ChanceAgain:           2,                                      // 1 in 2 chance an MT makes another trip through the universe
			EvenYearOnly:          true,                                   // true for only spawning mystery traders during even years
			GenesisDeviceCost:     Cost{0, 0, 0, 5000},                    // no miniaturization, always costs this much
			MaxMysteryTraders:     5,                                      // the maximum number of mystery traders spawned in a universe at one time
			MaxWarp:               13,                                     // the fastest warp a mystery trader will go
			MinWarp:               7,                                      // the slowest warp a mystery trader will go
			MinYear:               40,                                     // the earliest year a mystery trader will spawn
			RequestedBoon:         5000,                                   // how many minerals a player must give the MT to get a reward
			TechBoon: []MysteryTraderTechBoonRules{
				{
					TechLevels: 59,
					Rewards: []MysteryTraderTechBoonMineralsReward{
						{MineralsGiven: 5000, Reward: 6},
						{MineralsGiven: 6200, Reward: 7},
						{MineralsGiven: 7400, Reward: 8},
						{MineralsGiven: 8600, Reward: 9},
						{MineralsGiven: 9800, Reward: 10},
					},
				},
				{
					TechLevels: 71,
					Rewards: []MysteryTraderTechBoonMineralsReward{
						{MineralsGiven: 5000, Reward: 5},
						{MineralsGiven: 6200, Reward: 6},
						{MineralsGiven: 7400, Reward: 7},
						{MineralsGiven: 8600, Reward: 8},
						{MineralsGiven: 9800, Reward: 9},
					},
				},
				{
					TechLevels: 83,
					Rewards: []MysteryTraderTechBoonMineralsReward{
						{MineralsGiven: 5000, Reward: 4},
						{MineralsGiven: 6200, Reward: 5},
						{MineralsGiven: 7400, Reward: 6},
						{MineralsGiven: 8600, Reward: 7},
						{MineralsGiven: 9800, Reward: 8},
					},
				},
				{
					TechLevels: 95,
					Rewards: []MysteryTraderTechBoonMineralsReward{
						{MineralsGiven: 5000, Reward: 3},
						{MineralsGiven: 6200, Reward: 4},
						{MineralsGiven: 7400, Reward: 5},
						{MineralsGiven: 8600, Reward: 6},
						{MineralsGiven: 9800, Reward: 7},
					},
				},
				{
					TechLevels: 107,
					Rewards: []MysteryTraderTechBoonMineralsReward{
						{MineralsGiven: 5000, Reward: 2},
						{MineralsGiven: 6200, Reward: 2},
						{MineralsGiven: 7400, Reward: 2},
						{MineralsGiven: 8600, Reward: 2},
						{MineralsGiven: 9800, Reward: 2},
					},
				},
				{
					TechLevels: 108,
					Rewards: []MysteryTraderTechBoonMineralsReward{
						{MineralsGiven: 5000, Reward: 1},
						{MineralsGiven: 6200, Reward: 1},
						{MineralsGiven: 7400, Reward: 1},
						{MineralsGiven: 8600, Reward: 1},
						{MineralsGiven: 9800, Reward: 1},
					},
				},
			},
		},
		WormholeCloak: 75,
		WormholeStatsByStability: map[WormholeStability]WormholeStats{
			WormholeStabilityRockSolid: {
				YearsToDegrade: 10,
				ChanceToJump:   0,
				JiggleDistance: 10,
			},
			WormholeStabilityStable: {
				YearsToDegrade: 5,
				ChanceToJump:   0.005,
				JiggleDistance: 10,
			},
			WormholeStabilityMostlyStable: {
				YearsToDegrade: 5,
				ChanceToJump:   0.02,
				JiggleDistance: 10,
			},
			WormholeStabilityAverage: {
				YearsToDegrade: 5,
				ChanceToJump:   0.04,
				JiggleDistance: 10,
			},
			WormholeStabilitySlightlyVolatile: {
				YearsToDegrade: 5,
				ChanceToJump:   0.03,
				JiggleDistance: 10,
			},
			WormholeStabilityVolatile: {
				YearsToDegrade: 5,
				ChanceToJump:   0.06,
				JiggleDistance: 10,
			},
			WormholeStabilityExtremelyVolatile: {
				YearsToDegrade: Infinite,
				ChanceToJump:   0.04,
				JiggleDistance: 10,
			},
		},
		WormholePairsForSize: map[Size]int{
			SizeTiny:       1,
			SizeTinyWide:   1,
			SizeSmall:      3,
			SizeSmallWide:  3,
			SizeMedium:     4,
			SizeMediumWide: 4,
			SizeLarge:      5,
			SizeLargeWide:  5,
			SizeHuge:       6,
			SizeHugeWide:   6,
		},
		MineFieldStatsByType: map[MineFieldType]MineFieldStats{
			MineFieldTypeStandard: {
				MinDamagePerFleetRS: 600,
				DamagePerEngineRS:   125,
				MaxSpeed:            4,
				ChanceOfHit:         0.003,
				MinDamagePerFleet:   500,
				DamagePerEngine:     100,
				SweepFactor:         1.0,
				MinDecay:            10,
				CanDetonate:         true,
			},
			MineFieldTypeHeavy: {
				MinDamagePerFleetRS: 2500,
				DamagePerEngineRS:   600,
				MaxSpeed:            6,
				ChanceOfHit:         0.01,
				MinDamagePerFleet:   2000,
				DamagePerEngine:     500,
				SweepFactor:         1.0,
				MinDecay:            10,
				CanDetonate:         false,
			},
			MineFieldTypeSpeedBump: {
				MinDamagePerFleetRS: 0,
				DamagePerEngineRS:   0,
				MaxSpeed:            5,
				ChanceOfHit:         0.035,
				MinDamagePerFleet:   0,
				DamagePerEngine:     0,
				SweepFactor:         0.333333343,
				MinDecay:            0,
				CanDetonate:         false,
			},
		},
		RepairRates: map[RepairRate]float64{
			RepairRateNone:              0.0,
			RepairRateMoving:            0.01,
			RepairRateStopped:           0.02,
			RepairRateOrbiting:          0.03,
			RepairRateOrbitingOwnPlanet: 0.05,
			RepairRateStarbase:          0.1,
		},
		ShowPublicScoresAfterYears: 20,
		PlanetMinDistance:          15,
		MineralDecayFactor:         1_500_000,
		RemoteMiningMineOutput:     10,
		RaceStartingPoints:         1650,
		ScrapMineralAmount:         0.333333343,
		ScrapResourceAmount:        0.0,
		ScrapColonizeAmount:        0.75,
		SalvageFromBattleFactor:    .3,
		PacketDecayRate: map[int]float64{
			1: 0.1,
			2: 0.25,
			3: 0.5,
		},
		PacketMaxOverwarpSpeed: 3,
		PacketMinDecay:         10,
		MaxTechLevel:           26,
		PRTSpecs: map[PRT]PRTSpec{
			HE:   heSpec(),
			SS:   ssSpec(),
			WM:   wmSpec(),
			CA:   caSpec(),
			IS:   isSpec(),
			SD:   sdSpec(),
			PP:   ppSpec(),
			IT:   itSpec(),
			AR:   arSpec(),
			JoaT: joatSpec(),
		},
		LRTSpecs: map[LRT]LRTSpec{
			IFE:  ifeSpec(),
			TT:   ttSpec(),
			ARM:  armSpec(),
			ISB:  isbSpec(),
			GR:   grSpec(),
			UR:   urSpec(),
			NRSE: nrseSpec(),
			OBRM: obrmSpec(),
			NAS:  nasSpec(),
			LSP:  lspSpec(),
			BET:  betSpec(),
			RS:   rsSpec(),
			MA:   maSpec(),
			CE:   ceSpec(),
		},
		techs: &StaticTechStore,
	}
}

// Get the number of planets for a universe based on size and density
func (rules *Rules) GetNumPlanets(size Size, density Density) (int, error) {
	switch size {
	case SizeTiny, SizeTinyWide:
		switch density {
		case DensitySparse:
			return 24, nil
		case DensityNormal:
			return 32, nil
		case DensityDense:
			return 40, nil
		case DensityPacked:
			return 60, nil
		}
	case SizeSmall, SizeSmallWide:
		switch density {
		case DensitySparse:
			return 96, nil
		case DensityNormal:
			return 128, nil
		case DensityDense:
			return 160, nil
		case DensityPacked:
			return 240, nil
		}
	case SizeMedium, SizeMediumWide:
		switch density {
		case DensitySparse:
			return 216, nil
		case DensityNormal:
			return 288, nil
		case DensityDense:
			return 360, nil
		case DensityPacked:
			return 540, nil
		}
	case SizeLarge, SizeLargeWide:
		switch density {
		case DensitySparse:
			return 384, nil
		case DensityNormal:
			return 512, nil
		case DensityDense:
			return 640, nil
		case DensityPacked:
			return 910, nil
		}
	case SizeHuge, SizeHugeWide:
		switch density {
		case DensitySparse:
			return 600, nil
		case DensityNormal:
			return 800, nil
		case DensityDense:
			return 940, nil
		case DensityPacked:
			return 945, nil
		}

	}

	return 0, fmt.Errorf("unable to GetNumPlanets for Size: %v, Density: %v", size, density)
}

// Get the area of a universe based on size
func (rules *Rules) GetArea(size Size) (Vector, error) {

	switch size {
	case SizeTiny:
		return Vector{400, 400}, nil
	case SizeTinyWide:
		return Vector{500, 300}, nil
	case SizeSmall:
		return Vector{800, 800}, nil
	case SizeSmallWide:
		return Vector{1000, 600}, nil
	case SizeMedium:
		return Vector{1200, 1200}, nil
	case SizeMediumWide:
		return Vector{1500, 900}, nil
	case SizeLarge:
		return Vector{1600, 1600}, nil
	case SizeLargeWide:
		return Vector{2000, 1200}, nil
	case SizeHuge:
		return Vector{2000, 2000}, nil
	case SizeHugeWide:
		return Vector{2500, 1500}, nil
	}

	return Vector{}, fmt.Errorf("unable to GetArea for Size: %v", size)

}
