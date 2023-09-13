package ai

import (
	"math/rand"

	"github.com/sirgwain/craig-stars/cs"
)

var Races []cs.Race = []cs.Race{
	{
		Name:          "Berserker",
		PluralName:    "Berserkers",
		PRT:           cs.HE,
		LRTs:          cs.Bitmask(cs.IFE) | cs.Bitmask(cs.UR) | cs.Bitmask(cs.MA) | cs.Bitmask(cs.OBRM),
		ImmuneGrav:    true,
		ImmuneTemp:    true,
		ImmuneRad:     true,
		GrowthRate:    7,
		PopEfficiency: 8,
		FactoryOutput: 13,
		FactoryCost:   9,
		NumFactories:  16,
		MineOutput:    10,
		MineCost:      4,
		NumMines:      8,
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
		Name:       "Hooveron",
		PluralName: "Hooverons",
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
	},
	{
		Name:       "Humanoid",
		PluralName: "Humanoids",
		PRT:        cs.WM,
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
		Name:       "Humanoid",
		PluralName: "Humanoids",
		PRT:        cs.CA,
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
		Name:       "Humanoid",
		PluralName: "Humanoids",
		PRT:        cs.IS,
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
		Name:       "Humanoid",
		PluralName: "Humanoids",
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
		Name:       "Humanoid",
		PluralName: "Humanoids",
		PRT:        cs.PP,
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
		Name:       "Humanoid",
		PluralName: "Humanoids",
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
		Name:       "Humanoid",
		PluralName: "Humanoids",
		PRT:        cs.AR,
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

func GetRandomRace() cs.Race {
	return Races[rand.Intn(len(Races)-1)]
}