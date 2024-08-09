package ai

import (
	"math/rand"

	"github.com/sirgwain/craig-stars/cs"
)

var Races []cs.Race = []cs.Race{
	{
		Name:              "HE🫥",
		PluralName:        "HE🫥s",
		PRT:               cs.HE,
		LRTs:              cs.Bitmask(cs.IFE) | cs.Bitmask(cs.ISB) | cs.Bitmask(cs.NRSE) | cs.Bitmask(cs.GR) | cs.Bitmask(cs.OBRM) | cs.Bitmask(cs.NAS),
		ImmuneGrav:        true,
		ImmuneTemp:        true,
		ImmuneRad:         true,
		GrowthRate:        4,
		PopEfficiency:     10,
		FactoryOutput:     15,
		FactoryCost:       6,
		NumFactories:      23,
		FactoriesCostLess: true,
		MineOutput:        15,
		MineCost:          2,
		NumMines:          21,
		ResearchCost: cs.ResearchCost{
			Energy:        cs.ResearchCostLess,
			Weapons:       cs.ResearchCostLess,
			Propulsion:    cs.ResearchCostLess,
			Construction:  cs.ResearchCostLess,
			Electronics:   cs.ResearchCostLess,
			Biotechnology: cs.ResearchCostLess,
		},
	},
	{
		Name:       "SS",
		PluralName: "SSs",
		PRT:        cs.SS,
		LRTs:       cs.Bitmask(cs.IFE) | cs.Bitmask(cs.ARM) | cs.Bitmask(cs.MA) | cs.Bitmask(cs.RS),
		HabLow: cs.Hab{
			Grav: 20,
			Temp: 5,
		},
		HabHigh: cs.Hab{
			Grav: 92,
			Temp: 53,
		},
		ImmuneRad:     true,
		GrowthRate:    15,
		PopEfficiency: 8,
		FactoryOutput: 15,
		FactoryCost:   10,
		NumFactories:  25,
		MineOutput:    10,
		MineCost:      5,
		NumMines:      9,
		ResearchCost: cs.ResearchCost{
			Energy:        cs.ResearchCostExtra,
			Weapons:       cs.ResearchCostExtra,
			Propulsion:    cs.ResearchCostExtra,
			Construction:  cs.ResearchCostExtra,
			Electronics:   cs.ResearchCostExtra,
			Biotechnology: cs.ResearchCostExtra,
		},
		TechsStartHigh: true,
	},
	{
		Name:       "WM",
		PluralName: "WMs",
		PRT:        cs.WM,
		LRTs:       cs.Bitmask(cs.OBRM) | cs.Bitmask(cs.NAS) | cs.Bitmask(cs.CE),
		HabLow: cs.Hab{
			Grav: 40,
			Temp: 40,
		},
		HabHigh: cs.Hab{
			Grav: 100,
			Temp: 100,
		},
		ImmuneRad:         true,
		GrowthRate:        13,
		PopEfficiency:     10,
		FactoryOutput:     11,
		FactoryCost:       9,
		NumFactories:      12,
		FactoriesCostLess: true,
		MineOutput:        10,
		MineCost:          3,
		NumMines:          12,
		ResearchCost: cs.ResearchCost{
			Energy:        cs.ResearchCostExtra,
			Weapons:       cs.ResearchCostLess,
			Propulsion:    cs.ResearchCostStandard,
			Construction:  cs.ResearchCostLess,
			Electronics:   cs.ResearchCostExtra,
			Biotechnology: cs.ResearchCostExtra,
		},
		TechsStartHigh: true,
	},
	{
		Name:       "CA",
		PluralName: "CAs",
		PRT:        cs.CA,
		LRTs:       cs.Bitmask(cs.TT) | cs.Bitmask(cs.GR) | cs.Bitmask(cs.OBRM) | cs.Bitmask(cs.LSP),
		HabLow: cs.Hab{
			Grav: 28,
			Temp: 28,
			Rad:  28,
		},
		HabHigh: cs.Hab{
			Grav: 72,
			Temp: 72,
			Rad:  72,
		},
		GrowthRate:        19,
		PopEfficiency:     9,
		FactoryOutput:     10,
		FactoryCost:       10,
		NumFactories:      13,
		FactoriesCostLess: true,
		MineOutput:        10,
		MineCost:          5,
		NumMines:          11,
		ResearchCost: cs.ResearchCost{
			Energy:        cs.ResearchCostExtra,
			Weapons:       cs.ResearchCostStandard,
			Propulsion:    cs.ResearchCostStandard,
			Construction:  cs.ResearchCostStandard,
			Electronics:   cs.ResearchCostExtra,
			Biotechnology: cs.ResearchCostLess,
		},
	},
	{
		Name:       "IS",
		PluralName: "ISs",
		PRT:        cs.IS,
		LRTs:       cs.Bitmask(cs.IFE) | cs.Bitmask(cs.NRSE) | cs.Bitmask(cs.GR) | cs.Bitmask(cs.NAS),
		HabLow: cs.Hab{
			Grav: 0,
			Rad:  0,
		},
		HabHigh: cs.Hab{
			Grav: 52,
			Rad:  52,
		},
		ImmuneTemp:    true,
		GrowthRate:    13,
		PopEfficiency: 10,
		FactoryOutput: 10,
		FactoryCost:   10,
		NumFactories:  12,
		MineOutput:    10,
		MineCost:      5,
		NumMines:      12,
		ResearchCost: cs.ResearchCost{
			Energy:        cs.ResearchCostStandard,
			Weapons:       cs.ResearchCostLess,
			Propulsion:    cs.ResearchCostStandard,
			Construction:  cs.ResearchCostLess,
			Electronics:   cs.ResearchCostStandard,
			Biotechnology: cs.ResearchCostStandard,
		},
	},
	{
		Name:       "SD",
		PluralName: "SDs",
		PRT:        cs.SD,
		LRTs:       cs.LRTNone,
		HabLow: cs.Hab{
			Grav: 15,
			Temp: 15,
			Rad:  15,
		},
		HabHigh: cs.Hab{
			Grav: 85,
			Temp: 85,
			Rad:  85,
		},
		GrowthRate:    15,
		PopEfficiency: 10,
		FactoryOutput: 10,
		FactoryCost:   10,
		NumFactories:  10,
		MineOutput:    10,
		MineCost:      5,
		NumMines:      10,
		ResearchCost: cs.ResearchCost{
			Energy:        cs.ResearchCostStandard,
			Weapons:       cs.ResearchCostStandard,
			Propulsion:    cs.ResearchCostStandard,
			Construction:  cs.ResearchCostStandard,
			Electronics:   cs.ResearchCostStandard,
			Biotechnology: cs.ResearchCostStandard,
		},
	},
	{
		Name:       "PP",
		PluralName: "PPs",
		PRT:        cs.PP,
		LRTs:       cs.Bitmask(cs.ISB) | cs.Bitmask(cs.GR) | cs.Bitmask(cs.ARM),
		HabLow: cs.Hab{
			Temp: 0,
			Rad:  46,
		},
		HabHigh: cs.Hab{
			Temp: 52,
			Rad:  100,
		},
		ImmuneGrav:    true,
		GrowthRate:    14,
		PopEfficiency: 11,
		FactoryOutput: 9,
		FactoryCost:   10,
		NumFactories:  10,
		MineOutput:    10,
		MineCost:      4,
		NumMines:      11,
		ResearchCost: cs.ResearchCost{
			Energy:        cs.ResearchCostLess,
			Weapons:       cs.ResearchCostExtra,
			Propulsion:    cs.ResearchCostExtra,
			Construction:  cs.ResearchCostLess,
			Electronics:   cs.ResearchCostExtra,
			Biotechnology: cs.ResearchCostExtra,
		},
	},
	{
		Name:       "IT",
		PluralName: "ITs",
		PRT:        cs.IT,
		LRTs:       cs.LRTNone,
		HabLow: cs.Hab{
			Grav: 15,
			Temp: 15,
			Rad:  15,
		},
		HabHigh: cs.Hab{
			Grav: 85,
			Temp: 85,
			Rad:  85,
		},
		GrowthRate:    15,
		PopEfficiency: 10,
		FactoryOutput: 10,
		FactoryCost:   10,
		NumFactories:  10,
		MineOutput:    10,
		MineCost:      5,
		NumMines:      10,
		ResearchCost: cs.ResearchCost{
			Energy:        cs.ResearchCostStandard,
			Weapons:       cs.ResearchCostStandard,
			Propulsion:    cs.ResearchCostStandard,
			Construction:  cs.ResearchCostStandard,
			Electronics:   cs.ResearchCostStandard,
			Biotechnology: cs.ResearchCostStandard,
		},
	},
	{
		Name:       "AR",
		PluralName: "ARs",
		PRT:        cs.AR,
		LRTs:       cs.Bitmask(cs.IFE) | cs.Bitmask(cs.ISB) | cs.Bitmask(cs.ARM) | cs.Bitmask(cs.NRSE),
		HabLow: cs.Hab{
			Temp: 62,
			Rad:  64,
		},
		HabHigh: cs.Hab{
			Temp: 98,
			Rad:  100,
		},
		ImmuneGrav:    true,
		GrowthRate:    18,
		PopEfficiency: 10,
		ResearchCost: cs.ResearchCost{
			Energy:        cs.ResearchCostLess,
			Weapons:       cs.ResearchCostStandard,
			Propulsion:    cs.ResearchCostExtra,
			Construction:  cs.ResearchCostStandard,
			Electronics:   cs.ResearchCostExtra,
			Biotechnology: cs.ResearchCostExtra,
		},
	},
	// humanoids
	{
		Name:       "Humanoid",
		PluralName: "Humanoids",
		PRT:        cs.JoaT,
		LRTs:       cs.LRTNone,
		HabLow: cs.Hab{
			Grav: 15,
			Temp: 15,
			Rad:  15,
		},
		HabHigh: cs.Hab{
			Grav: 85,
			Temp: 85,
			Rad:  85,
		},
		GrowthRate:    15,
		PopEfficiency: 10,
		FactoryOutput: 10,
		FactoryCost:   10,
		NumFactories:  10,
		MineOutput:    10,
		MineCost:      5,
		NumMines:      10,
		ResearchCost: cs.ResearchCost{
			Energy:        cs.ResearchCostStandard,
			Weapons:       cs.ResearchCostStandard,
			Propulsion:    cs.ResearchCostStandard,
			Construction:  cs.ResearchCostStandard,
			Electronics:   cs.ResearchCostStandard,
			Biotechnology: cs.ResearchCostStandard,
		},
	},
}

var CheaterRaces []cs.Race = []cs.Race{
	{
		Name:       "Waaagh!",
		PluralName: "Waaagh!",
		PRT:        cs.JoaT,
		LRTs:       cs.Bitmask(cs.IFE) | cs.Bitmask(cs.ISB) | cs.Bitmask(cs.OBRM),
		HabLow: cs.Hab{
			Temp: 15,
			Rad:  15,
		},
		HabHigh: cs.Hab{
			Temp: 85,
			Rad:  85,
		},
		ImmuneGrav:    true,
		GrowthRate:    20,
		PopEfficiency: 8,
		FactoryOutput: 13,
		FactoryCost:   5,
		NumFactories:  15,
		MineOutput:    13,
		MineCost:      2,
		NumMines:      15,
		ResearchCost: cs.ResearchCost{
			Energy:        cs.ResearchCostStandard,
			Weapons:       cs.ResearchCostLess,
			Propulsion:    cs.ResearchCostStandard,
			Construction:  cs.ResearchCostLess,
			Electronics:   cs.ResearchCostStandard,
			Biotechnology: cs.ResearchCostExtra,
		},
	},
	{
		Name:              "Space Ork",
		PluralName:        "Space Orks",
		PRT:               cs.HE,
		LRTs:              cs.Bitmask(cs.IFE) | cs.Bitmask(cs.ISB) | cs.Bitmask(cs.NRSE) | cs.Bitmask(cs.OBRM),
		ImmuneGrav:        true,
		ImmuneTemp:        true,
		ImmuneRad:         true,
		GrowthRate:        10,
		PopEfficiency:     10,
		FactoryOutput:     15,
		FactoryCost:       6,
		NumFactories:      15,
		FactoriesCostLess: true,
		MineOutput:        12,
		MineCost:          3,
		NumMines:          16,
		ResearchCost: cs.ResearchCost{
			Energy:        cs.ResearchCostStandard,
			Weapons:       cs.ResearchCostLess,
			Propulsion:    cs.ResearchCostExtra,
			Construction:  cs.ResearchCostLess,
			Electronics:   cs.ResearchCostExtra,
			Biotechnology: cs.ResearchCostExtra,
		},
	},
}

// Get a list of random races with minimal repeats from either the Cheater races or regular races
func GetRandomRaces(numRaces int, cheater bool) []cs.Race {
	if cheater {
		return getRandomRaces(numRaces, CheaterRaces)
	}
	return getRandomRaces(numRaces, Races)
}

func getRandomRaces(numRaces int, racesToPickFrom []cs.Race) []cs.Race {
	raceNums := make([]int, len(racesToPickFrom))
	races := make([]cs.Race, numRaces)

	for i := 0; i < len(racesToPickFrom); i++ {
		raceNums[i] = i
	}

	// get a shuffled list of races we can pull from
	rand.Shuffle(len(raceNums), func(i, j int) { raceNums[i], raceNums[j] = raceNums[j], raceNums[i] })

	// make sure we include each race only once, then repeat with a new randomized list of AI races if we reach the end
	count := 0
	for r := 0; r < numRaces; r++ {
		if count >= len(raceNums) {
			// we've reached the end of our race list
			// reset count and re-shuffle the list
			count = 0
			rand.Shuffle(len(raceNums), func(i, j int) { raceNums[i], raceNums[j] = raceNums[j], raceNums[i] })
		}
		races[r] = racesToPickFrom[raceNums[count]]
		count++
	}

	return races
}
