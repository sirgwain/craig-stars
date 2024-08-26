//go:build wasi || wasm

package wasm

import (
	"syscall/js"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

func GetInt[T ~uint | ~uint32 | ~uint64 | ~int | ~int32 | ~int64](o js.Value, key string) T {
	var result T
	val := o.Get(key)
	if !val.IsUndefined() {
		result = T(val.Int())
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

func GetString(o js.Value, key string) string {
	var result string
	val := o.Get(key)
	if !val.IsUndefined() {
		result = val.String()
	}
	return result
}

func GetTime(o js.Value, key string) (time.Time, error) {
	var result time.Time
	val := o.Get(key)
	if !val.IsUndefined() {
		// time assumes json string has quotes
		if err := result.UnmarshalJSON([]byte("\"" + val.String() + "\"")); err != nil {
			return result, err
		}
	}
	return result, nil

}

func GetRace(o js.Value) cs.Race {
	obj := cs.Race{}
	obj.DBObject = GetDBObject(o)
	obj.UserID = GetInt[int64](o, "userId")
	obj.Name = GetString(o, "name")
	obj.PluralName = GetString(o, "pluralName")
	obj.SpendLeftoverPointsOn = cs.SpendLeftoverPointsOn(GetString(o, "spendLeftoverPointsOn"))
	obj.PRT = cs.PRT(GetString(o, "prt"))
	obj.LRTs = cs.Bitmask(GetInt[uint32](o, "lrts"))
	obj.HabLow = GetHab(o.Get("habLow"))
	obj.HabHigh = GetHab(o.Get("habHigh"))
	obj.GrowthRate = GetInt[int](o, "growthRate")
	obj.PopEfficiency = GetInt[int](o, "popEfficiency")
	obj.FactoryOutput = GetInt[int](o, "factoryOutput")
	obj.FactoryCost = GetInt[int](o, "factoryCost")
	obj.NumFactories = GetInt[int](o, "numFactories")
	obj.FactoriesCostLess = GetBool(o, "factoriesCostLess")
	obj.ImmuneGrav = GetBool(o, "immuneGrav")
	obj.ImmuneTemp = GetBool(o, "immuneTemp")
	obj.ImmuneRad = GetBool(o, "immuneRad")
	obj.MineOutput = GetInt[int](o, "mineOutput")
	obj.MineCost = GetInt[int](o, "mineCost")
	obj.NumMines = GetInt[int](o, "numMines")
	obj.ResearchCost = GetResearchCost(o.Get("researchCost"))
	obj.TechsStartHigh = GetBool(o, "techsStartHigh")
	return obj
}

func GetResearchCost(o js.Value) cs.ResearchCost {
	obj := cs.ResearchCost{}
	obj.Energy = cs.ResearchCostLevel(GetString(o, "energy"))
	obj.Weapons = cs.ResearchCostLevel(GetString(o, "weapons"))
	obj.Propulsion = cs.ResearchCostLevel(GetString(o, "propulsion"))
	obj.Construction = cs.ResearchCostLevel(GetString(o, "construction"))
	obj.Electronics = cs.ResearchCostLevel(GetString(o, "electronics"))
	obj.Biotechnology = cs.ResearchCostLevel(GetString(o, "biotechnology"))
	return obj
}

func GetHab(o js.Value) cs.Hab {
	obj := cs.Hab{}
	obj.Grav = GetInt[int](o, "grav")
	obj.Temp = GetInt[int](o, "temp")
	obj.Rad = GetInt[int](o, "rad")
	return obj
}

func GetTechLevel(o js.Value) cs.TechLevel {
	obj := cs.TechLevel{}
	obj.Energy = GetInt[int](o, "energy")
	obj.Weapons = GetInt[int](o, "weapons")
	obj.Propulsion = GetInt[int](o, "propulsion")
	obj.Construction = GetInt[int](o, "construction")
	obj.Electronics = GetInt[int](o, "electronics")
	obj.Biotechnology = GetInt[int](o, "biotechnology")
	return obj
}

func GetCost(o js.Value) cs.Cost {
	obj := cs.Cost{}
	obj.Ironium = GetInt[int](o, "ironium")
	obj.Boranium = GetInt[int](o, "boranium")
	obj.Germanium = GetInt[int](o, "germanium")
	obj.Resources = GetInt[int](o, "resources")
	return obj
}

func GetDBObject(o js.Value) cs.DBObject {
	obj := cs.DBObject{}
	obj.ID = GetInt[int64](o, "id")	
	obj.CreatedAt, _ = GetTime(o, "createdAt")	
	obj.UpdatedAt, _ = GetTime(o, "updatedAt")
	return obj
}
