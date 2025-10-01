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
	ID                                 int64                               `json:"id,omitempty"`
	CreatedAt                          time.Time                           `json:"createdAt,omitempty"`
	UpdatedAt                          time.Time                           `json:"updatedAt,omitempty"`
	GameID                             int64                               `json:"gameId,omitempty"`
	AcquirablePartTradeChanceBase      float64                             `json:"acquirablePartTradeChanceBase"`
	AcquirablePartTradeItemMax         int                                 `json:"acquirablePartTradeItemMax"`
	CometStatsBySize                   map[CometSize]CometStats            `json:"cometStatsBySize"`
	FleetSafeSpeedExplosionChance      float64                             `json:"fleetSafeSpeedExplosionChance"`
	InvasionDefenseCoverageFactor      float64                             `json:"invasionDefenseCoverageFactor"`
	LRTSpecs                           map[LRT]LRTSpec                     `json:"lrtSpecs"`
	MaxPopulation                      int                                 `json:"maxPopulation"`
	MinPopFloor                        int                                 `json:"minPopFloor"`
	MaxTechLevel                       int                                 `json:"maxTechLevel"`
	MinefieldCloak                     int                                 `json:"minefieldCloak"`
	MinefieldStatsByType               map[MinefieldType]MinefieldStats    `json:"minefieldStatsByType"`
	MineralDecayFactor                 int                                 `json:"mineralDecayFactor"`
	MinHabFloor                        int                                 `json:"minHabFloor"` //@sirgwain: Do we need this? It's only used as a default value for race generation
	MysteryTraderRules                 MysteryTraderRules                  `json:"mysteryTraderRules"`
	PacketDecayRate                    map[int]float64                     `json:"packetDecayRate"`
	PacketMaxOverwarpSpeed             int                                 `json:"packetMaxOverwarpSpeed"`
	PacketMinDecay                     int                                 `json:"packetMinDecay"`
	PlanetMinDistance                  int                                 `json:"planetMinDistance"`
	PopulationOvercrowdDieoffRate      float64                             `json:"populationOvercrowdDieoffRate"`
	PopulationOvercrowdDieoffRateMax   float64                             `json:"populationOvercrowdDieoffRateMax"`
	PopulationOvercrowdResourcePenalty float64                             `json:"populationOvercrowdResourcePenalty"`
	PopulationOvercrowdResourceMax     float64                             `json:"populationOvercrowdResourceMax"`
	PopulationScannerError             float64                             `json:"populationScannerError"`
	PRTSpecs                           map[PRT]PRTSpec                     `json:"prtSpecs"`
	RaceStartingPoints                 int                                 `json:"raceStartingPoints"` // TODO: Change this into a "handicap" system with bonuses/penalties per PRT/LRT
	RadiatingImmune                    int                                 `json:"radiatingImmune"`
	RandomArtifactResearchBonusRange   []int                               `json:"randomArtifactResearchBonusRange"`
	RandomCometMinYear                 int                                 `json:"randomCometMinYear"`
	RandomCometMinYearPlayerWorld      int                                 `json:"randomCometMinYearPlayerWorld"`
	RandomEventChances                 map[RandomEvent]float64             `json:"randomEventChances"`
	RandomMineralDepositBonusRange     []int                               `json:"randomMineralDepositBonusRange"`
	RemoteMiningMineOutput             int                                 `json:"remoteMiningMineOutput"`
	RepairRates                        map[RepairRate]float64              `json:"repairRates"`
	SalvageDecayMin                    int                                 `json:"salvageDecayMin"`
	SalvageDecayRate                   float64                             `json:"salvageDecayRate"`
	SalvageFromBattleFactor            float64                             `json:"salvageFromBattleFactor"`
	ScrapColonizeAmount                float64                             `json:"scrapColonizeAmount"`
	ScrapMineralAmount                 float64                             `json:"scrapMineralAmount"`
	ScrapResourceAmount                float64                             `json:"scrapResourceAmount"`
	ShowPublicScoresAfterYears         int                                 `json:"showPublicScoresAfterYears"`
	SmartDefenseCoverageFactor         float64                             `json:"smartDefenseCoverageFactor"`
	StargateMaxHullMassFactor          int                                 `json:"stargateMaxHullMassFactor"`
	StargateMaxRangeFactor             int                                 `json:"stargateMaxRangeFactor"`
	TachyonCloakReduction              float64                             `json:"tachyonCloakReduction"`
	TachyonMaxCloakReduction           float64                             `json:"tachyonMaxCloakReduction"`
	TechsID                            int64                               `json:"techsId,omitempty"`
	TechTradeChance                    float64                             `json:"techTradeChance"`
	TorpedoSplashDamage                float64                             `json:"torpedoSplashDamage"`
	WormholeCloak                      int                                 `json:"wormholeCloak"`
	WormholePairsForSize               map[Size]int                        `json:"wormholePairsForSize"`
	WormholeStatsByStability           map[WormholeStability]WormholeStats `json:"wormholeStatsByStability"`
	random                             rng
	techs                              *TechStore
}

type UniverseGenerationRules struct {
	BorderInset                               int                           `json:"borderInset"`
	HabDropoffRange                           Hab                           `json:"habDropoffRange"` // Controls up to how many clicks (inclusive) away from MinHab & MaxHab planet habs become linearly less likely
	HighRadMineralConcentrationBonusThreshold int                           `json:"highRadMineralConcentrationBonusThreshold"`
	LimitMineralConcentration                 int                           `json:"limitMineralConcentration"`
	MaxHab                                    int                           `json:"maxHab"`
	MaxMineralConcentration                   int                           `json:"maxMineralConcentration"`
	MaxStartingMineralConcentration           int                           `json:"maxStartingMineralConcentration"`
	MaxStartingMineralSurface                 int                           `json:"maxStartingMineralSurface"`
	MinExtraPlanetMineralConcentration        int                           `json:"minExtraPlanetMineralConcentration"`
	MinHab                                    int                           `json:"minHab"`
	MinHomeworldMineralConcentration          int                           `json:"minHomeworldMineralConcentration"`
	MinMineralConcentration                   int                           `json:"minMineralConcentration"`
	MinPlanetSpacing                          int                           `json:"minPlanetSpacing"`
	MinStartingMineralConcentration           int                           `json:"minStartingMineralConcentration"`
	MinStartingMineralSurface                 int                           `json:"minStartingMineralSurface"`
	SqLyPerPlanet                             int                           `json:"planetsPerSqLy"`
	RaceLeftoverPointsPerItem                 map[SpendLeftoverPointsOn]int `json:"raceLeftoverPointsPerItem"` // amount of points required for 1 starting point increase; for surface minerals this is instead the unit rate in kT/point
	StartingYear                              int                           `json:"startingYear"`
	WormholeMinPlanetDistance                 int                           `json:"wormholeMinPlanetDistance"`
}

type CostRules struct {
	DefenseCost                    Cost    `json:"defenseCost"`
	FactoryCostGermanium           int     `json:"factoryCostGermanium"`
	MineralAlchemyCost             int     `json:"mineralAlchemyCost"`
	PlanetaryScannerCost           Cost    `json:"planetaryScannerCost"`
	StarbaseComponentCostReduction float64 `json:"starbaseComponentCostReduction"`
	StarbaseHullRefundFactor       float64 `json:"starbaseHullRefundFactor"`
	TerraformCost                  Cost    `json:"terraformCost"`
	TechBaseCost                   []int   `json:"techBaseCost"`
}

// A slightly fancier map[bool]float64 that can be serialized to JSON
type JammerCap struct {
	Ship     float64 `json:"ship,omitempty"`
	Starbase float64 `json:"starbase,omitempty"`
}

func (jc JammerCap) Get(starbase bool) float64 {
	if starbase {
		return jc.Starbase
	}
	return jc.Ship
}

type BattleRules struct {
	BeamRangeDropoff    float64   `json:"beamRangeDropoff"`
	BeamBonusCap        float64   `json:"beamBonusCap"`
	JammerCap           JammerCap `json:"jammerCap"`
	JammerMulti         JammerCap `json:"jammerMulti"`
	MovementMin         int       `json:"movementMin"`
	MovementMax         int       `json:"movementMax"`
	MovesToRunAway      int       `json:"movesToRunAway"`
	NumBattleRounds     int       `json:"numBattleRounds"`
	TorpedoSplashDamage float64   `json:"torpedoSplashDamage"`
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
	CometUnspecified CometSize = ""
	CometSmall       CometSize = "Small"
	CometMedium      CometSize = "Medium"
	CometLarge       CometSize = "Large"
	CometHuge        CometSize = "Huge"
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
	AllMinerals              int     `json:"allMinerals"`
	AllRandomMinerals        int     `json:"allRandomMinerals"`
	BonusMinerals            int     `json:"bonusMinerals"`
	BonusRandomMinerals      int     `json:"bonusRandomMinerals"`
	BonusMinConcentration    int     `json:"bonusMinConcentration"`
	BonusRandomConcentration int     `json:"bonusRandomConcentration"`
	BonusAffectsMinerals     int     `json:"bonusAffectsMinerals"`
	MinTerraform             int     `json:"minTerraform"`
	RandomTerraform          int     `json:"randomTerraform"`
	AffectsHabs              int     `json:"affectsHabs"`
	PopKilledPercent         float64 `json:"popKilledPercent"`
}

type RepairRate string

const (
	RepairRateNone              RepairRate = ""
	RepairRateMoving            RepairRate = "Moving"
	RepairRateStopped           RepairRate = "Stopped"
	RepairRateOrbiting          RepairRate = "Orbiting"
	RepairRateOrbitingOwnPlanet RepairRate = "OrbitingOwnPlanet"
	RepairRateStarbase          RepairRate = "Starbase"
)

type MysteryTraderRules struct {
	ChanceSpawn           []int                        `json:"chanceSpawn"`
	ChanceMaxTechGetsPart int                          `json:"chanceMaxTechGetsPart"`
	ChanceCourseChange    int                          `json:"chanceCourseChange"`
	ChanceSpeedUpOnly     int                          `json:"chanceSpeedUpOnly"`
	ChanceAgain           int                          `json:"chanceAgain"`
	EvenYearOnly          bool                         `json:"evenYearOnly"`
	GenesisDeviceCost     Cost                         `json:"genesisDeviceCost"`
	MaxMysteryTraders     int                          `json:"maxMysteryTraders"`
	MaxWarp               int                          `json:"maxWarp"`
	MinWarp               int                          `json:"minWarp"`
	MinYear               int                          `json:"minYear"`
	RequestedBoon         int                          `json:"requestedBoon"`
	TechBoon              []MysteryTraderTechBoonRules `json:"techBoon"`
}

type MysteryTraderTechBoonRules struct {
	TechLevels int                                   `json:"techLevels"`
	Rewards    []MysteryTraderTechBoonMineralsReward `json:"rewards"`
}

type MysteryTraderTechBoonMineralsReward struct {
	MineralsGiven int `json:"mineralsGiven"`
	Reward        int `json:"reward"`
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

	// @sirgwain: We should consider moving these comments to the corresponding
	// struct field definitions for editor syntax highlighting
	// (also just more comments never hurts)
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
			StarbaseHullRefundFactor:       0.5, // 50% of the old base's hull cost goes towards the new base
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
			BeamRangeDropoff: 0.1,  // 10% pro-rated damage penalty
			BeamBonusCap:     2.55, // 2.55x damage max from beam capacitors
			JammerCap: JammerCap{
				Starbase: 1,    // starbases have no explicit jamming hardcap, but an innate 0.75x jamming multi
				Ship:     0.95, // ships hardcap at 95% jamming
			},
			JammerMulti: JammerCap{
				Starbase: 0.75, // starbases have innate 0.75x jamming multipler by default
				Ship:     1,    // ships have no innate jamming multipler
			},
			MovementMin:         2,  // minimum of 2 battle board movement (1, 0, 1, 0...)
			MovementMax:         10, // minimum of 10 battle board movement (3, 2, 3, 2...)
			MovesToRunAway:      7,
			NumBattleRounds:     16,
			TorpedoSplashDamage: 0.125,
		},
		UniverseGenerationRules: UniverseGenerationRules{
			BorderInset: 20,
			// The first 9 Grav/Temp hab values from either edge (1-9 & 91-99) are linearly less likely to generate.
			// More specifically, a hab value N clicks away from MinHab/MaxHab with dropoff range of H
			// becomes (N+1/H+1)x as likely as a normal mid-value hab
			// Ex: 6 temp is 5 clicks away from min (1) and is thus 6/10x as likely to generate;
			// 99 temp is 1 click away from max (100) and is thus 1/10x as likely.
			HabDropoffRange: Hab{
				Grav: 9,
				Temp: 9,
				Rad:  0,
			},
			HighRadMineralConcentrationBonusThreshold: 90,
			MinPlanetSpacing:                   12,
			MinHomeworldMineralConcentration:   30,
			MinExtraPlanetMineralConcentration: 30,
			MinMineralConcentration:            1,
			MaxMineralConcentration:            200,
			MinHab:                             1,
			MaxHab:                             99,
			MinStartingMineralConcentration:    1,
			MaxStartingMineralConcentration:    121,
			LimitMineralConcentration:          30,
			MaxStartingMineralSurface:          1000,
			MinStartingMineralSurface:          300,
			SqLyPerPlanet:                      5000, // base numplanets is area / 5000
			RaceLeftoverPointsPerItem: map[SpendLeftoverPointsOn]int{
				SpendLeftoverPointsOnMines:                 2,
				SpendLeftoverPointsOnFactories:             5,
				SpendLeftoverPointsOnDefenses:              10,
				SpendLeftoverPointsOnMineralConcentrations: 2,  // Due to some high level source code chicanery
				SpendLeftoverPointsOnSurfaceMinerals:       10, // special case; denotes kT/point
			},
			StartingYear:              2400,
			WormholeMinPlanetDistance: 30,
		},
		// TODO: Change tachyon cloak reduction to a property of the technology itself
		TachyonCloakReduction:              .05, // 5% diminishing cloak reduction per detector
		TachyonMaxCloakReduction:           .81, // tachyon detectors cap at 81% cloaking reduction
		MaxPopulation:                      1_000_000,
		MinPopFloor:                        100,  // low value planets cannot fall below 100 pop from natural growth
		MinHabFloor:                        5,    // minimum of 5% effective habitability for inhabited planet productivity/maxpop
		PopulationOvercrowdDieoffRate:      .04,  // overcrowded pops die off at 4% per 100% over cap
		PopulationOvercrowdDieoffRateMax:   .12,  // overcrowded pops will not die off more than 12% (400% capacity) per year
		PopulationOvercrowdResourcePenalty: 0.5,  // overcrowded pop produce resources at 50% efficiency
		PopulationOvercrowdResourceMax:     1,    // maximum 100% extra resources from overcrowded pop
		PopulationScannerError:             0.2,  // opponents' scanners have ±20% error on pop readings
		SmartDefenseCoverageFactor:         0.5,  // smart bombs penetrate 50% enemy defenses
		InvasionDefenseCoverageFactor:      0.75, // invasions penetrate 25% enemy defenses
		SalvageDecayRate:                   0.1,
		SalvageDecayMin:                    10,
		MinefieldCloak:                     75,
		StargateMaxRangeFactor:             5,  // ships can only gate up to 5x gate safe range
		StargateMaxHullMassFactor:          5,  // ships can only gate up to 5x gate safe mass
		TechTradeChance:                    .5, // 50% chance of tech trading per level
		FleetSafeSpeedExplosionChance:      .1, // 10% chance of losing a ship
		// TODO: Make this a property of the TechHullComponent
		RadiatingImmune: 85, // hab center of >= 85 are immune to radating damage
		RandomEventChances: map[RandomEvent]float64{
			RandomEventComet:           .05, // 5% chance of a planet being struck by a comet in a given turn
			RandomEventMineralDeposit:  .05,
			RandomEventPlanetaryChange: .05,
			RandomEventAncientArtifact: 1.0 / 3, // 1 in 3 planets have random artifacts
		},
		AcquirablePartTradeChanceBase: 0.005, // 0.5% chance per item in fleet
		AcquirablePartTradeItemMax:    25,    // 25 items max per trade instance (12.5% chance)
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
		RandomMineralDepositBonusRange:   []int{20, 50},
		RandomArtifactResearchBonusRange: []int{120, 400},
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
		MinefieldStatsByType: map[MinefieldType]MinefieldStats{
			MinefieldTypeStandard: {
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
			MinefieldTypeHeavy: {
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
			MinefieldTypeSpeedBump: {
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
	dim, err := rules.GetArea(size)
	if err != nil {
		return 0, err
	}

	// start with something like 1200/5000 for medium
	base := int(dim.X*dim.Y) / rules.SqLyPerPlanet

	// add 25% less increments of more planets based on density
	switch density {
	case DensitySparse:
		return base - base/4, nil
	case DensityNormal:
		return base, nil
	case DensityDense:
		return base + base/4, nil
	case DensityPacked:
		return base + base*3/4, nil
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
