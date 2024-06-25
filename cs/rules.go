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
	ID                                        int64                               `json:"id"`
	CreatedAt                                 time.Time                           `json:"createdAt"`
	UpdatedAt                                 time.Time                           `json:"updatedAt"`
	GameID                                    int64                               `json:"gameId"`
	TachyonCloakReduction                     int                                 `json:"tachyonCloakReduction"`
	MaxPopulation                             int                                 `json:"maxPopulation"`
	MinMaxPopulationPercent                   float64                             `json:"minMaxPopulationPercent"`
	PopulationOvercrowdDieoffRate             float64                             `json:"populationOvercrowdDieoffRate"`
	PopulationOvercrowdDieoffRateMax          float64                             `json:"populationOvercrowdDieoffRateMax"`
	PopulationScannerError                    float64                             `json:"populationScannerError"`
	SmartDefenseCoverageFactor                float64                             `json:"smartDefenseCoverageFactor"`
	InvasionDefenseCoverageFactor             float64                             `json:"invasionDefenseCoverageFactor"`
	NumBattleRounds                           int                                 `json:"numBattleRounds"`
	MovesToRunAway                            int                                 `json:"movesToRunAway"`
	BeamRangeDropoff                          float64                             `json:"beamRangeDropoff"`
	TorpedoSplashDamage                       float64                             `json:"torpedoSplashDamage"`
	SalvageDecayRate                          float64                             `json:"salvageDecayRate"`
	SalvageDecayMin                           int                                 `json:"salvageDecayMin"`
	MineFieldCloak                            int                                 `json:"mineFieldCloak"`
	StargateMaxRangeFactor                    int                                 `json:"stargateMaxRangeFactor"`
	StargateMaxHullMassFactor                 int                                 `json:"stargateMaxHullMassFactor"`
	FleetSafeSpeedExplosionChance             float64                             `json:"fleetSafeSpeedExplosionChance"`
	RandomEventChances                        map[RandomEvent]float64             `json:"randomEventChances"`
	RandomMineralDepositBonusRange            [2]int                              `json:"randomMineralDepositBonusRange"`
	RandomArtifactResearchBonusRange          [2]int                              `json:"randomArtifactResearchBonusRange"`
	RandomCometMinYear                        int                                 `json:"randomCometMinYear,omitempty"`
	RandomCometMinYearPlayerWorld             int                                 `json:"randomCometMinYearPlayerWorld,omitempty"`
	CometStatsBySize                          map[CometSize]CometStats            `json:"cometStatsBySize,omitempty"`
	WormholeCloak                             int                                 `json:"wormholeCloak"`
	WormholeMinPlanetDistance                 int                                 `json:"wormholeMinDistance"`
	WormholeStatsByStability                  map[WormholeStability]WormholeStats `json:"wormholeStatsByStability"`
	WormholePairsForSize                      map[Size]int                        `json:"wormholePairsForSize"`
	MineFieldStatsByType                      map[MineFieldType]MineFieldStats    `json:"mineFieldStatsByType"`
	RepairRates                               map[RepairRate]float64              `json:"repairRates"`
	MaxPlayers                                int                                 `json:"maxPlayers"`
	StartingYear                              int                                 `json:"startingYear"`
	ShowPublicScoresAfterYears                int                                 `json:"showPublicScoresAfterYears"`
	PlanetMinDistance                         int                                 `json:"planetMinDistance"`
	MaxExtraWorldDistance                     int                                 `json:"maxExtraWorldDistance"`
	MinExtraWorldDistance                     int                                 `json:"minExtraWorldDistance"`
	MinHomeworldMineralConcentration          int                                 `json:"minHomeworldMineralConcentration"`
	MinExtraPlanetMineralConcentration        int                                 `json:"minExtraPlanetMineralConcentration"`
	MinHab                                    int                                 `json:"minHab"`
	MaxHab                                    int                                 `json:"maxHab"`
	MinMineralConcentration                   int                                 `json:"minMineralConcentration"`
	MaxMineralConcentration                   int                                 `json:"maxMineralConcentration"`
	MinStartingMineralConcentration           int                                 `json:"minStartingMineralConcentration"`
	MaxStartingMineralConcentration           int                                 `json:"maxStartingMineralConcentration"`
	HighRadMineralConcentrationBonusThreshold int                                 `json:"highRadGermaniumBonusThreshold"`
	RadiatingImmune                           int                                 `json:"radiatingImmune"`
	MaxStartingMineralSurface                 int                                 `json:"maxStartingMineralSurface"`
	MinStartingMineralSurface                 int                                 `json:"minStartingMineralSurface"`
	MineralDecayFactor                        int                                 `json:"mineralDecayFactor"`
	RemoteMiningMineOutput                    int                                 `json:"remoteMiningMineOutput"`
	StartingMines                             int                                 `json:"startingMines"`
	StartingFactories                         int                                 `json:"startingFactories"`
	StartingDefenses                          int                                 `json:"startingDefenses"`
	RaceStartingPoints                        int                                 `json:"raceStartingPoints"`
	ScrapMineralAmount                        float64                             `json:"scrapMineralAmount"`
	ScrapResourceAmount                       float64                             `json:"scrapResourceAmount"`
	FactoryCostGermanium                      int                                 `json:"factoryCostGermanium"`
	DefenseCost                               Cost                                `json:"defenseCost"`
	MineralAlchemyCost                        int                                 `json:"mineralAlchemyCost"`
	PlanetaryScannerCost                      Cost                                `json:"planetaryScannerCost"`
	TerraformCost                             Cost                                `json:"terraformCost"`
	StarbaseComponentCostFactor               float64                             `json:"starbaseComponentCostFactor"`
	SalvageFromBattleFactor                   float64                             `json:"salvageFromBattleFactor"`
	TechTradeChance                           float64                             `json:"techTradeChance"`
	PacketDecayRate                           map[int]float64                     `json:"packetDecayRate"`
	MaxTechLevel                              int                                 `json:"maxTechLevel"`
	TechBaseCost                              []int                               `json:"techBaseCost"`
	PRTSpecs                                  map[PRT]PRTSpec                     `json:"prtSpecs"`
	LRTSpecs                                  map[LRT]LRTSpec                     `json:"lrtSpecs"`
	TechsID                                   int64                               `json:"techsId"`
	random                                    *rand.Rand
	techs                                     *TechStore
}

type RandomEvent string

const (
	RandomEventComet           RandomEvent = "Comet"
	RandomEventMineralDeposit  RandomEvent = "MineralDeposit"
	RandomEventPlanetaryChange RandomEvent = "PlanetaryChange"
	RandomEventAncientArtifact RandomEvent = "AncientArtifact"
	RandomEventMysteryTrader   RandomEvent = "MysteryTrader"
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
	AllMinerals              int     `json:"minMinerals,omitempty"`
	AllRandomMinerals        int     `json:"randomMinerals,omitempty"`
	BonusMinerals            int     `json:"bonusMinerals,omitempty"`
	BonusRandomMinerals      int     `json:"bonusRandomMinerals,omitempty"`
	BonusMinConcentration    int     `json:"minConcentrationBonus,omitempty"`
	BonusRandomConcentration int     `json:"randomConcentrationBonus,omitempty"`
	BonusAffectsMinerals     int     `json:"affectsMinerals,omitempty"`
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
		random:                           random,
		TachyonCloakReduction:            5,
		MaxPopulation:                    1000000,
		MinMaxPopulationPercent:          .05,
		PopulationOvercrowdDieoffRate:    .04, // overcrowded pops die off at 4% per doubling
		PopulationOvercrowdDieoffRateMax: .12, // overcrowded pops will not die off more than 12% (3x pop) in a year
		PopulationScannerError:           0.2,
		SmartDefenseCoverageFactor:       0.5,
		InvasionDefenseCoverageFactor:    0.75,
		NumBattleRounds:                  16,
		MovesToRunAway:                   7,
		BeamRangeDropoff:                 0.1,
		TorpedoSplashDamage:              0.125,
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
			RandomEventMysteryTrader:   .05,
		},
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
		WormholeCloak:                    75,
		WormholeMinPlanetDistance:        30,
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
		MaxPlayers:                                16,
		StartingYear:                              2400,
		ShowPublicScoresAfterYears:                20,
		PlanetMinDistance:                         15,
		MaxExtraWorldDistance:                     180,
		MinExtraWorldDistance:                     130,
		MinHomeworldMineralConcentration:          30,
		MinExtraPlanetMineralConcentration:        30,
		MinMineralConcentration:                   1,
		MaxMineralConcentration:                   200,
		MinHab:                                    1,
		MaxHab:                                    99,
		MinStartingMineralConcentration:           31,
		MaxStartingMineralConcentration:           121,
		HighRadMineralConcentrationBonusThreshold: 90,
		MaxStartingMineralSurface:                 1000,
		MinStartingMineralSurface:                 300,
		MineralDecayFactor:                        1_500_000,
		RemoteMiningMineOutput:                    10,
		StartingMines:                             10,
		StartingFactories:                         10,
		StartingDefenses:                          10,
		RaceStartingPoints:                        1650,
		ScrapMineralAmount:                        0.333333343,
		ScrapResourceAmount:                       0.0,
		FactoryCostGermanium:                      4,
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
		TerraformCost: Cost{
			Ironium:   0,
			Boranium:  0,
			Germanium: 0,
			Resources: 100,
		},
		StarbaseComponentCostFactor: 0.5,
		SalvageFromBattleFactor:     .3,
		PacketDecayRate: map[int]float64{
			1: 0.1,
			2: 0.25,
			3: 0.5,
		},
		MaxTechLevel: 26,
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
