//go:build wasi || wasm

package wasm

import (
	"syscall/js"

	"github.com/sirgwain/craig-stars/cs"
)



func GetRace(o js.Value) cs.Race {
	race := cs.Race{
		UserID:                int64(o.Get("userId").Int()),
		Name:                  o.Get("name").String(),
		PluralName:            o.Get("pluralName").String(),
		SpendLeftoverPointsOn: cs.SpendLeftoverPointsOn(o.Get("spendLeftoverPointsOn").String()),
		PRT:                   cs.PRT(o.Get("prt").String()),
		LRTs:                  cs.Bitmask(GetInt(o, "lrts")),
		HabLow:                GetHab(o.Get("habLow")),
		HabHigh:               GetHab(o.Get("habHigh")),
		GrowthRate:            GetInt(o, "growthRate"),
		PopEfficiency:         GetInt(o, "popEfficiency"),
		FactoryOutput:         GetInt(o, "factoryOutput"),
		FactoryCost:           GetInt(o, "factoryCost"),
		NumFactories:          GetInt(o, "numFactories"),
		FactoriesCostLess:     GetBool(o, "factoriesCostLess"),
		ImmuneGrav:            GetBool(o, "immuneGrav"),
		ImmuneTemp:            GetBool(o, "immuneTemp"),
		ImmuneRad:             GetBool(o, "immuneRad"),
		MineOutput:            GetInt(o, "mineOutput"),
		MineCost:              GetInt(o, "mineCost"),
		NumMines:              GetInt(o, "numMines"),
		ResearchCost:          GetResearchCost(o.Get("researchCost")),
		TechsStartHigh:        GetBool(o, "techsStartHigh"),
	}

	return race
}

func GetHab(o js.Value) cs.Hab {
	return cs.Hab{
		Grav: GetInt(o, "grav"),
		Temp: GetInt(o, "temp"),
		Rad:  GetInt(o, "rad"),
	}
}

func GetResearchCost(o js.Value) cs.ResearchCost {
	return cs.ResearchCost{
		Energy:        cs.ResearchCostLevel(o.Get("Energy").String()),
		Weapons:       cs.ResearchCostLevel(o.Get("Weapons").String()),
		Propulsion:    cs.ResearchCostLevel(o.Get("Propulsion").String()),
		Construction:  cs.ResearchCostLevel(o.Get("Construction").String()),
		Electronics:   cs.ResearchCostLevel(o.Get("Electronics").String()),
		Biotechnology: cs.ResearchCostLevel(o.Get("Biotechnology").String()),
	}
}

func GetInt(o js.Value, key string) int {
	var result int
	val := o.Get(key)
	if !val.IsUndefined() {
		result = val.Int()
	}
	return result
}

func GetBool(o js.Value, key string) bool {
	var result bool
	val := o.Get(key)
	if !val.IsUndefined() {
		result = val.Bool()
	}
	return result
}
