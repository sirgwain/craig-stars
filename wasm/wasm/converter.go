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

func GetIntArray[T ~uint | ~uint32 | ~uint64 | ~int | ~int32 | ~int64](o js.Value, key string) []T {
	items := make([]T, o.Length())
	for i := 0; i < len(items); i++ {
		items[i] = T(o.Index(i).Int())
	}
	return items

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

func GetBattlePlan(o js.Value) cs.BattlePlan {
	obj := cs.BattlePlan{}
	obj.Num = GetInt[int](o, "num")
	obj.Name = GetString(o, "name")
	obj.PrimaryTarget = cs.BattleTarget(GetString(o, "primaryTarget"))
	obj.SecondaryTarget = cs.BattleTarget(GetString(o, "secondaryTarget"))
	obj.Tactic = cs.BattleTactic(GetString(o, "tactic"))
	obj.AttackWho = cs.BattleAttackWho(GetString(o, "attackWho"))
	obj.DumpCargo = GetBool(o, "dumpCargo")
	return obj
}

func GetBattlePlanArray(o js.Value) []cs.BattlePlan {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.BattlePlan, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetBattlePlan(o.Index(i))
	}
	return items
}

func GetBattlePlanPointerArray(o js.Value) []*cs.BattlePlan {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.BattlePlan, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetBattlePlan(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetBattleRecord(o js.Value) cs.BattleRecord {
	obj := cs.BattleRecord{}
	obj.Num = GetInt[int](o, "num")
	obj.PlanetNum = GetInt[int](o, "planetNum")
	return obj
}

func GetBattleRecordArray(o js.Value) []cs.BattleRecord {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.BattleRecord, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetBattleRecord(o.Index(i))
	}
	return items
}

func GetBattleRecordPointerArray(o js.Value) []*cs.BattleRecord {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.BattleRecord, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetBattleRecord(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetBomb(o js.Value) cs.Bomb {
	obj := cs.Bomb{}
	obj.Quantity = GetInt[int](o, "quantity")
	obj.MinKillRate = GetInt[int](o, "minKillRate")
	obj.UnterraformRate = GetInt[int](o, "unterraformRate")
	return obj
}

func GetBombArray(o js.Value) []cs.Bomb {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.Bomb, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetBomb(o.Index(i))
	}
	return items
}

func GetBombPointerArray(o js.Value) []*cs.Bomb {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.Bomb, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetBomb(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetCargo(o js.Value) cs.Cargo {
	obj := cs.Cargo{}
	obj.Ironium = GetInt[int](o, "ironium")
	obj.Boranium = GetInt[int](o, "boranium")
	obj.Germanium = GetInt[int](o, "germanium")
	obj.Colonists = GetInt[int](o, "colonists")
	return obj
}

func GetCargoArray(o js.Value) []cs.Cargo {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.Cargo, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetCargo(o.Index(i))
	}
	return items
}

func GetCargoPointerArray(o js.Value) []*cs.Cargo {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.Cargo, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetCargo(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetCost(o js.Value) cs.Cost {
	obj := cs.Cost{}
	obj.Ironium = GetInt[int](o, "ironium")
	obj.Boranium = GetInt[int](o, "boranium")
	obj.Germanium = GetInt[int](o, "germanium")
	obj.Resources = GetInt[int](o, "resources")
	return obj
}

func GetCostArray(o js.Value) []cs.Cost {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.Cost, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetCost(o.Index(i))
	}
	return items
}

func GetCostPointerArray(o js.Value) []*cs.Cost {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.Cost, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetCost(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetDBObject(o js.Value) cs.DBObject {
	obj := cs.DBObject{}
	obj.ID = GetInt[int64](o, "id")	
	obj.CreatedAt, _ = GetTime(o, "createdAt")	
	obj.UpdatedAt, _ = GetTime(o, "updatedAt")
	return obj
}

func GetDBObjectArray(o js.Value) []cs.DBObject {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.DBObject, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetDBObject(o.Index(i))
	}
	return items
}

func GetDBObjectPointerArray(o js.Value) []*cs.DBObject {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.DBObject, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetDBObject(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetEngine(o js.Value) cs.Engine {
	obj := cs.Engine{}
	obj.IdealSpeed = GetInt[int](o, "idealSpeed")
	obj.FreeSpeed = GetInt[int](o, "freeSpeed")
	obj.MaxSafeSpeed = GetInt[int](o, "maxSafeSpeed")
	return obj
}

func GetEngineArray(o js.Value) []cs.Engine {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.Engine, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetEngine(o.Index(i))
	}
	return items
}

func GetEnginePointerArray(o js.Value) []*cs.Engine {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.Engine, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetEngine(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetFleetIntel(o js.Value) cs.FleetIntel {
	obj := cs.FleetIntel{}
	obj.PlanetIntelID = GetInt[int64](o, "")
	obj.BaseName = GetString(o, "baseName")
	obj.OrbitingPlanetNum = GetInt[int](o, "orbitingPlanetNum")
	obj.WarpSpeed = GetInt[int](o, "warpSpeed")
	obj.Mass = GetInt[int](o, "mass")
	obj.Cargo = GetCargo(o.Get("cargo"))
	obj.CargoDiscovered = GetBool(o, "cargoDiscovered")
	obj.Freighter = GetBool(o, "freighter")
	obj.ScanRange = GetInt[int](o, "scanRange")
	obj.ScanRangePen = GetInt[int](o, "scanRangePen")
	return obj
}

func GetFleetIntelArray(o js.Value) []cs.FleetIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.FleetIntel, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetFleetIntel(o.Index(i))
	}
	return items
}

func GetFleetIntelPointerArray(o js.Value) []*cs.FleetIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.FleetIntel, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetFleetIntel(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetGameDBObject(o js.Value) cs.GameDBObject {
	obj := cs.GameDBObject{}
	obj.ID = GetInt[int64](o, "id")
	obj.GameID = GetInt[int64](o, "gameId")	
	obj.CreatedAt, _ = GetTime(o, "createdAt")	
	obj.UpdatedAt, _ = GetTime(o, "updatedAt")
	return obj
}

func GetGameDBObjectArray(o js.Value) []cs.GameDBObject {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.GameDBObject, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetGameDBObject(o.Index(i))
	}
	return items
}

func GetGameDBObjectPointerArray(o js.Value) []*cs.GameDBObject {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.GameDBObject, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetGameDBObject(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetHab(o js.Value) cs.Hab {
	obj := cs.Hab{}
	obj.Grav = GetInt[int](o, "grav")
	obj.Temp = GetInt[int](o, "temp")
	obj.Rad = GetInt[int](o, "rad")
	return obj
}

func GetHabArray(o js.Value) []cs.Hab {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.Hab, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetHab(o.Index(i))
	}
	return items
}

func GetHabPointerArray(o js.Value) []*cs.Hab {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.Hab, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetHab(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetMapObject(o js.Value) cs.MapObject {
	obj := cs.MapObject{}
	obj.GameDBObject = GetGameDBObject(o)
	obj.Type = cs.MapObjectType(GetString(o, "type"))
	obj.Num = GetInt[int](o, "num")
	obj.PlayerNum = GetInt[int](o, "playerNum")
	obj.Name = GetString(o, "name")
	return obj
}

func GetMapObjectArray(o js.Value) []cs.MapObject {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.MapObject, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetMapObject(o.Index(i))
	}
	return items
}

func GetMapObjectPointerArray(o js.Value) []*cs.MapObject {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.MapObject, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetMapObject(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetMineFieldIntel(o js.Value) cs.MineFieldIntel {
	obj := cs.MineFieldIntel{}
	obj.NumMines = GetInt[int](o, "numMines")
	obj.MineFieldType = cs.MineFieldType(GetString(o, "mineFieldType"))
	return obj
}

func GetMineFieldIntelArray(o js.Value) []cs.MineFieldIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.MineFieldIntel, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetMineFieldIntel(o.Index(i))
	}
	return items
}

func GetMineFieldIntelPointerArray(o js.Value) []*cs.MineFieldIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.MineFieldIntel, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetMineFieldIntel(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetMineral(o js.Value) cs.Mineral {
	obj := cs.Mineral{}
	obj.Ironium = GetInt[int](o, "ironium")
	obj.Boranium = GetInt[int](o, "boranium")
	obj.Germanium = GetInt[int](o, "germanium")
	return obj
}

func GetMineralArray(o js.Value) []cs.Mineral {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.Mineral, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetMineral(o.Index(i))
	}
	return items
}

func GetMineralPointerArray(o js.Value) []*cs.Mineral {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.Mineral, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetMineral(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetMineralPacketIntel(o js.Value) cs.MineralPacketIntel {
	obj := cs.MineralPacketIntel{}
	obj.WarpSpeed = GetInt[int](o, "warpSpeed")
	obj.Cargo = GetCargo(o.Get("cargo"))
	obj.TargetPlanetNum = GetInt[int](o, "targetPlanetNum")
	obj.ScanRange = GetInt[int](o, "scanRange")
	obj.ScanRangePen = GetInt[int](o, "scanRangePen")
	return obj
}

func GetMineralPacketIntelArray(o js.Value) []cs.MineralPacketIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.MineralPacketIntel, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetMineralPacketIntel(o.Index(i))
	}
	return items
}

func GetMineralPacketIntelPointerArray(o js.Value) []*cs.MineralPacketIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.MineralPacketIntel, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetMineralPacketIntel(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetMysteryTraderIntel(o js.Value) cs.MysteryTraderIntel {
	obj := cs.MysteryTraderIntel{}
	obj.WarpSpeed = GetInt[int](o, "warpSpeed")
	obj.RequestedBoon = GetInt[int](o, "requestedBoon")
	return obj
}

func GetMysteryTraderIntelArray(o js.Value) []cs.MysteryTraderIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.MysteryTraderIntel, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetMysteryTraderIntel(o.Index(i))
	}
	return items
}

func GetMysteryTraderIntelPointerArray(o js.Value) []*cs.MysteryTraderIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.MysteryTraderIntel, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetMysteryTraderIntel(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlanet(o js.Value) cs.Planet {
	obj := cs.Planet{}
	obj.MapObject = GetMapObject(o)
	obj.PlanetOrders = GetPlanetOrders(o)
	obj.Hab = GetHab(o.Get("hab"))
	obj.BaseHab = GetHab(o.Get("baseHab"))
	obj.TerraformedAmount = GetHab(o.Get("terraformedAmount"))
	obj.MineralConcentration = GetMineral(o.Get("mineralConcentration"))
	obj.MineYears = GetMineral(o.Get("mineYears"))
	obj.Cargo = GetCargo(o.Get("cargo"))
	obj.Mines = GetInt[int](o, "mines")
	obj.Factories = GetInt[int](o, "factories")
	obj.Defenses = GetInt[int](o, "defenses")
	obj.Homeworld = GetBool(o, "homeworld")
	obj.Scanner = GetBool(o, "scanner")
	obj.Spec = GetPlanetSpec(o.Get("spec"))
	return obj
}

func GetPlanetArray(o js.Value) []cs.Planet {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.Planet, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlanet(o.Index(i))
	}
	return items
}

func GetPlanetPointerArray(o js.Value) []*cs.Planet {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.Planet, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlanet(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlanetIntel(o js.Value) cs.PlanetIntel {
	obj := cs.PlanetIntel{}
	obj.Hab = GetHab(o.Get("hab"))
	obj.BaseHab = GetHab(o.Get("baseHab"))
	obj.MineralConcentration = GetMineral(o.Get("mineralConcentration"))
	starbase := GetFleetIntel(o.Get("starbase"))
	obj.Starbase = &starbase
	obj.Cargo = GetCargo(o.Get("cargo"))
	obj.CargoDiscovered = GetBool(o, "cargoDiscovered")
	obj.PlanetHabitability = GetInt[int](o, "planetHabitability")
	obj.PlanetHabitabilityTerraformed = GetInt[int](o, "planetHabitabilityTerraformed")
	obj.Homeworld = GetBool(o, "homeworld")
	obj.Spec = GetPlanetSpec(o.Get("spec"))
	return obj
}

func GetPlanetIntelArray(o js.Value) []cs.PlanetIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlanetIntel, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlanetIntel(o.Index(i))
	}
	return items
}

func GetPlanetIntelPointerArray(o js.Value) []*cs.PlanetIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlanetIntel, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlanetIntel(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlanetOrders(o js.Value) cs.PlanetOrders {
	obj := cs.PlanetOrders{}
	obj.ContributesOnlyLeftoverToResearch = GetBool(o, "contributesOnlyLeftoverToResearch")
	obj.ProductionQueue = GetProductionQueueItemArray(o.Get("productionQueue"))
	obj.RouteTargetType = cs.MapObjectType(GetString(o, "routeTargetType"))
	obj.RouteTargetNum = GetInt[int](o, "routeTargetNum")
	obj.RouteTargetPlayerNum = GetInt[int](o, "routeTargetPlayerNum")
	obj.PacketTargetNum = GetInt[int](o, "packetTargetNum")
	obj.PacketSpeed = GetInt[int](o, "packetSpeed")
	return obj
}

func GetPlanetOrdersArray(o js.Value) []cs.PlanetOrders {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlanetOrders, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlanetOrders(o.Index(i))
	}
	return items
}

func GetPlanetOrdersPointerArray(o js.Value) []*cs.PlanetOrders {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlanetOrders, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlanetOrders(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlanetSpec(o js.Value) cs.PlanetSpec {
	obj := cs.PlanetSpec{}
	obj.PlanetStarbaseSpec = GetPlanetStarbaseSpec(o)
	obj.CanTerraform = GetBool(o, "canTerraform")
	obj.Defense = GetString(o, "defense")
	obj.GrowthAmount = GetInt[int](o, "growthAmount")
	obj.Habitability = GetInt[int](o, "habitability")
	obj.MaxDefenses = GetInt[int](o, "maxDefenses")
	obj.MaxFactories = GetInt[int](o, "maxFactories")
	obj.MaxMines = GetInt[int](o, "maxMines")
	obj.MaxPopulation = GetInt[int](o, "maxPopulation")
	obj.MaxPossibleFactories = GetInt[int](o, "maxPossibleFactories")
	obj.MaxPossibleMines = GetInt[int](o, "maxPossibleMines")
	obj.MiningOutput = GetMineral(o.Get("miningOutput"))
	obj.Population = GetInt[int](o, "population")
	obj.ResourcesPerYear = GetInt[int](o, "resourcesPerYear")
	obj.ResourcesPerYearAvailable = GetInt[int](o, "resourcesPerYearAvailable")
	obj.ResourcesPerYearResearch = GetInt[int](o, "resourcesPerYearResearch")
	obj.ResourcesPerYearResearchEstimatedLeftover = GetInt[int](o, "resourcesPerYearResearchEstimatedLeftover")
	obj.Scanner = GetString(o, "scanner")
	obj.ScanRange = GetInt[int](o, "scanRange")
	obj.ScanRangePen = GetInt[int](o, "scanRangePen")
	obj.TerraformAmount = GetHab(o.Get("terraformAmount"))
	obj.MinTerraformAmount = GetHab(o.Get("minTerraformAmount"))
	obj.TerraformedHabitability = GetInt[int](o, "terraformedHabitability")
	obj.Contested = GetBool(o, "contested")
	return obj
}

func GetPlanetSpecArray(o js.Value) []cs.PlanetSpec {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlanetSpec, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlanetSpec(o.Index(i))
	}
	return items
}

func GetPlanetSpecPointerArray(o js.Value) []*cs.PlanetSpec {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlanetSpec, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlanetSpec(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlanetStarbaseSpec(o js.Value) cs.PlanetStarbaseSpec {
	obj := cs.PlanetStarbaseSpec{}
	obj.HasMassDriver = GetBool(o, "hasMassDriver")
	obj.HasStarbase = GetBool(o, "hasStarbase")
	obj.HasStargate = GetBool(o, "hasStargate")
	obj.StarbaseDesignName = GetString(o, "starbaseDesignName")
	obj.StarbaseDesignNum = GetInt[int](o, "starbaseDesignNum")
	obj.DockCapacity = GetInt[int](o, "dockCapacity")
	obj.BasePacketSpeed = GetInt[int](o, "basePacketSpeed")
	obj.SafePacketSpeed = GetInt[int](o, "safePacketSpeed")
	obj.SafeHullMass = GetInt[int](o, "safeHullMass")
	obj.SafeRange = GetInt[int](o, "safeRange")
	obj.MaxRange = GetInt[int](o, "maxRange")
	obj.MaxHullMass = GetInt[int](o, "maxHullMass")
	obj.Stargate = GetString(o, "stargate")
	obj.MassDriver = GetString(o, "massDriver")
	return obj
}

func GetPlanetStarbaseSpecArray(o js.Value) []cs.PlanetStarbaseSpec {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlanetStarbaseSpec, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlanetStarbaseSpec(o.Index(i))
	}
	return items
}

func GetPlanetStarbaseSpecPointerArray(o js.Value) []*cs.PlanetStarbaseSpec {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlanetStarbaseSpec, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlanetStarbaseSpec(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlayer(o js.Value) cs.Player {
	obj := cs.Player{}
	obj.GameDBObject = GetGameDBObject(o)
	obj.PlayerOrders = GetPlayerOrders(o)
	obj.PlayerIntels = GetPlayerIntels(o)
	obj.PlayerPlans = GetPlayerPlans(o)
	obj.UserID = GetInt[int64](o, "userId")
	obj.Name = GetString(o, "name")
	obj.Num = GetInt[int](o, "num")
	obj.Ready = GetBool(o, "ready")
	obj.AIControlled = GetBool(o, "aiControlled")
	obj.AIDifficulty = cs.AIDifficulty(GetString(o, "aiDifficulty"))
	obj.Guest = GetBool(o, "guest")
	obj.SubmittedTurn = GetBool(o, "submittedTurn")
	obj.Color = GetString(o, "color")
	obj.DefaultHullSet = GetInt[int](o, "defaultHullSet")
	obj.Race = GetRace(o.Get("race"))
	obj.TechLevels = GetTechLevel(o.Get("techLevels"))
	obj.TechLevelsSpent = GetTechLevel(o.Get("techLevelsSpent"))
	obj.ResearchSpentLastYear = GetInt[int](o, "researchSpentLastYear")
	obj.Relations = GetPlayerRelationshipArray(o.Get("relations"))
	obj.Messages = GetPlayerMessageArray(o.Get("messages"))
	obj.Designs = GetShipDesignPointerArray(o.Get("designs"))
	obj.ScoreHistory = GetPlayerScoreArray(o.Get("scoreHistory"))
	acquiredTechs := make(map[string]bool)
	acquiredTechsObj := o.Get("acquiredTechs")
	acquiredTechsKeys := js.Global().Get("Object").Call("keys", acquiredTechsObj)
	for i := 0; i < acquiredTechsKeys.Length(); i++ {
		key := acquiredTechsKeys.Index(i).String()
		acquiredTechs[key] = GetBool(acquiredTechsObj, key)
	}
	obj.AcquiredTechs = acquiredTechs
	obj.AchievedVictoryConditions = cs.Bitmask(GetInt[uint32](o, "achievedVictoryConditions"))
	obj.Victor = GetBool(o, "victor")
	stats := GetPlayerStats(o.Get("stats"))
	obj.Stats = &stats
	obj.Spec = GetPlayerSpec(o.Get("spec"))
	return obj
}

func GetPlayerArray(o js.Value) []cs.Player {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.Player, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlayer(o.Index(i))
	}
	return items
}

func GetPlayerPointerArray(o js.Value) []*cs.Player {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.Player, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlayer(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlayerIntel(o js.Value) cs.PlayerIntel {
	obj := cs.PlayerIntel{}
	obj.Name = GetString(o, "name")
	obj.Num = GetInt[int](o, "num")
	obj.Color = GetString(o, "color")
	obj.Seen = GetBool(o, "seen")
	obj.RaceName = GetString(o, "raceName")
	obj.RacePluralName = GetString(o, "racePluralName")
	return obj
}

func GetPlayerIntelArray(o js.Value) []cs.PlayerIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlayerIntel, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlayerIntel(o.Index(i))
	}
	return items
}

func GetPlayerIntelPointerArray(o js.Value) []*cs.PlayerIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlayerIntel, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlayerIntel(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlayerIntels(o js.Value) cs.PlayerIntels {
	obj := cs.PlayerIntels{}
	obj.BattleRecords = GetBattleRecordArray(o.Get("battleRecords"))
	obj.PlayerIntels = GetPlayerIntelArray(o.Get("playerIntels"))
	obj.ScoreIntels = GetScoreIntelArray(o.Get("scoreIntels"))
	obj.PlanetIntels = GetPlanetIntelArray(o.Get("planetIntels"))
	obj.FleetIntels = GetFleetIntelArray(o.Get("fleetIntels"))
	obj.StarbaseIntels = GetFleetIntelArray(o.Get("starbaseIntels"))
	obj.ShipDesignIntels = GetShipDesignIntelArray(o.Get("shipDesignIntels"))
	obj.MineralPacketIntels = GetMineralPacketIntelArray(o.Get("mineralPacketIntels"))
	obj.MineFieldIntels = GetMineFieldIntelArray(o.Get("mineFieldIntels"))
	obj.WormholeIntels = GetWormholeIntelArray(o.Get("wormholeIntels"))
	obj.MysteryTraderIntels = GetMysteryTraderIntelArray(o.Get("mysteryTraderIntels"))
	obj.SalvageIntels = GetSalvageIntelArray(o.Get("salvageIntels"))
	return obj
}

func GetPlayerIntelsArray(o js.Value) []cs.PlayerIntels {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlayerIntels, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlayerIntels(o.Index(i))
	}
	return items
}

func GetPlayerIntelsPointerArray(o js.Value) []*cs.PlayerIntels {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlayerIntels, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlayerIntels(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlayerMessage(o js.Value) cs.PlayerMessage {
	obj := cs.PlayerMessage{}
	obj.Type = cs.PlayerMessageType(GetInt[int](o, "type"))
	obj.Text = GetString(o, "text")
	obj.BattleNum = GetInt[int](o, "battleNum")
	return obj
}

func GetPlayerMessageArray(o js.Value) []cs.PlayerMessage {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlayerMessage, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlayerMessage(o.Index(i))
	}
	return items
}

func GetPlayerMessagePointerArray(o js.Value) []*cs.PlayerMessage {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlayerMessage, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlayerMessage(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlayerOrders(o js.Value) cs.PlayerOrders {
	obj := cs.PlayerOrders{}
	obj.Researching = cs.TechField(GetString(o, "researching"))
	obj.NextResearchField = cs.NextResearchField(GetString(o, "nextResearchField"))
	obj.ResearchAmount = GetInt[int](o, "researchAmount")
	return obj
}

func GetPlayerOrdersArray(o js.Value) []cs.PlayerOrders {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlayerOrders, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlayerOrders(o.Index(i))
	}
	return items
}

func GetPlayerOrdersPointerArray(o js.Value) []*cs.PlayerOrders {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlayerOrders, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlayerOrders(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlayerPlans(o js.Value) cs.PlayerPlans {
	obj := cs.PlayerPlans{}
	obj.ProductionPlans = GetProductionPlanArray(o.Get("productionPlans"))
	obj.BattlePlans = GetBattlePlanArray(o.Get("battlePlans"))
	obj.TransportPlans = GetTransportPlanArray(o.Get("transportPlans"))
	return obj
}

func GetPlayerPlansArray(o js.Value) []cs.PlayerPlans {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlayerPlans, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlayerPlans(o.Index(i))
	}
	return items
}

func GetPlayerPlansPointerArray(o js.Value) []*cs.PlayerPlans {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlayerPlans, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlayerPlans(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlayerRelationship(o js.Value) cs.PlayerRelationship {
	obj := cs.PlayerRelationship{}
	obj.Relation = cs.PlayerRelation(GetString(o, "relation"))
	obj.ShareMap = GetBool(o, "shareMap")
	return obj
}

func GetPlayerRelationshipArray(o js.Value) []cs.PlayerRelationship {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlayerRelationship, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlayerRelationship(o.Index(i))
	}
	return items
}

func GetPlayerRelationshipPointerArray(o js.Value) []*cs.PlayerRelationship {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlayerRelationship, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlayerRelationship(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlayerScore(o js.Value) cs.PlayerScore {
	obj := cs.PlayerScore{}
	obj.Planets = GetInt[int](o, "planets")
	obj.Starbases = GetInt[int](o, "starbases")
	obj.UnarmedShips = GetInt[int](o, "unarmedShips")
	obj.EscortShips = GetInt[int](o, "escortShips")
	obj.CapitalShips = GetInt[int](o, "capitalShips")
	obj.TechLevels = GetInt[int](o, "techLevels")
	obj.Resources = GetInt[int](o, "resources")
	obj.Score = GetInt[int](o, "score")
	obj.Rank = GetInt[int](o, "rank")
	obj.AchievedVictoryConditions = cs.Bitmask(GetInt[uint32](o, "achievedVictoryConditions"))
	return obj
}

func GetPlayerScoreArray(o js.Value) []cs.PlayerScore {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlayerScore, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlayerScore(o.Index(i))
	}
	return items
}

func GetPlayerScorePointerArray(o js.Value) []*cs.PlayerScore {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlayerScore, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlayerScore(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlayerSpec(o js.Value) cs.PlayerSpec {
	obj := cs.PlayerSpec{}
	obj.ResourcesPerYear = GetInt[int](o, "resourcesPerYear")
	obj.ResourcesPerYearResearch = GetInt[int](o, "resourcesPerYearResearch")
	obj.ResourcesPerYearResearchEstimated = GetInt[int](o, "resourcesPerYearResearchEstimated")
	obj.CurrentResearchCost = GetInt[int](o, "currentResearchCost")
	return obj
}

func GetPlayerSpecArray(o js.Value) []cs.PlayerSpec {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlayerSpec, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlayerSpec(o.Index(i))
	}
	return items
}

func GetPlayerSpecPointerArray(o js.Value) []*cs.PlayerSpec {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlayerSpec, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlayerSpec(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetPlayerStats(o js.Value) cs.PlayerStats {
	obj := cs.PlayerStats{}
	obj.FleetsBuilt = GetInt[int](o, "fleetsBuilt")
	obj.StarbasesBuilt = GetInt[int](o, "starbasesBuilt")
	obj.TokensBuilt = GetInt[int](o, "tokensBuilt")
	obj.PlanetsColonized = GetInt[int](o, "planetsColonized")
	return obj
}

func GetPlayerStatsArray(o js.Value) []cs.PlayerStats {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.PlayerStats, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetPlayerStats(o.Index(i))
	}
	return items
}

func GetPlayerStatsPointerArray(o js.Value) []*cs.PlayerStats {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.PlayerStats, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetPlayerStats(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetProductionPlan(o js.Value) cs.ProductionPlan {
	obj := cs.ProductionPlan{}
	obj.Num = GetInt[int](o, "num")
	obj.Name = GetString(o, "name")
	obj.ContributesOnlyLeftoverToResearch = GetBool(o, "contributesOnlyLeftoverToResearch")
	return obj
}

func GetProductionPlanArray(o js.Value) []cs.ProductionPlan {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.ProductionPlan, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetProductionPlan(o.Index(i))
	}
	return items
}

func GetProductionPlanPointerArray(o js.Value) []*cs.ProductionPlan {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.ProductionPlan, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetProductionPlan(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetProductionQueueItem(o js.Value) cs.ProductionQueueItem {
	obj := cs.ProductionQueueItem{}
	obj.Type = cs.QueueItemType(GetString(o, "type"))
	obj.DesignNum = GetInt[int](o, "designNum")
	obj.Quantity = GetInt[int](o, "quantity")
	obj.Allocated = GetCost(o.Get("allocated"))
	return obj
}

func GetProductionQueueItemArray(o js.Value) []cs.ProductionQueueItem {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.ProductionQueueItem, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetProductionQueueItem(o.Index(i))
	}
	return items
}

func GetProductionQueueItemPointerArray(o js.Value) []*cs.ProductionQueueItem {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.ProductionQueueItem, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetProductionQueueItem(o.Index(i))
		items[i] = &item
	}
	return items
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

func GetRaceArray(o js.Value) []cs.Race {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.Race, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetRace(o.Index(i))
	}
	return items
}

func GetRacePointerArray(o js.Value) []*cs.Race {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.Race, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetRace(o.Index(i))
		items[i] = &item
	}
	return items
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

func GetResearchCostArray(o js.Value) []cs.ResearchCost {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.ResearchCost, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetResearchCost(o.Index(i))
	}
	return items
}

func GetResearchCostPointerArray(o js.Value) []*cs.ResearchCost {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.ResearchCost, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetResearchCost(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetSalvageIntel(o js.Value) cs.SalvageIntel {
	obj := cs.SalvageIntel{}
	obj.Cargo = GetCargo(o.Get("cargo"))
	return obj
}

func GetSalvageIntelArray(o js.Value) []cs.SalvageIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.SalvageIntel, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetSalvageIntel(o.Index(i))
	}
	return items
}

func GetSalvageIntelPointerArray(o js.Value) []*cs.SalvageIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.SalvageIntel, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetSalvageIntel(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetScoreIntel(o js.Value) cs.ScoreIntel {
	obj := cs.ScoreIntel{}
	obj.ScoreHistory = GetPlayerScoreArray(o.Get("scoreHistory"))
	return obj
}

func GetScoreIntelArray(o js.Value) []cs.ScoreIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.ScoreIntel, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetScoreIntel(o.Index(i))
	}
	return items
}

func GetScoreIntelPointerArray(o js.Value) []*cs.ScoreIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.ScoreIntel, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetScoreIntel(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetShipDesign(o js.Value) cs.ShipDesign {
	obj := cs.ShipDesign{}
	obj.GameDBObject = GetGameDBObject(o)
	obj.Num = GetInt[int](o, "num")
	obj.PlayerNum = GetInt[int](o, "playerNum")
	obj.OriginalPlayerNum = GetInt[int](o, "originalPlayerNum")
	obj.Name = GetString(o, "name")
	obj.Version = GetInt[int](o, "version")
	obj.Hull = GetString(o, "hull")
	obj.HullSetNumber = GetInt[int](o, "hullSetNumber")
	obj.CannotDelete = GetBool(o, "cannotDelete")
	obj.MysteryTrader = GetBool(o, "mysteryTrader")
	obj.Slots = GetShipDesignSlotArray(o.Get("slots"))
	obj.Purpose = cs.ShipDesignPurpose(GetString(o, "purpose"))
	obj.Spec = GetShipDesignSpec(o.Get("spec"))
	obj.Delete = GetBool(o, "")
	return obj
}

func GetShipDesignArray(o js.Value) []cs.ShipDesign {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.ShipDesign, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetShipDesign(o.Index(i))
	}
	return items
}

func GetShipDesignPointerArray(o js.Value) []*cs.ShipDesign {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.ShipDesign, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetShipDesign(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetShipDesignIntel(o js.Value) cs.ShipDesignIntel {
	obj := cs.ShipDesignIntel{}
	obj.Hull = GetString(o, "hull")
	obj.HullSetNumber = GetInt[int](o, "hullSetNumber")
	obj.Version = GetInt[int](o, "version")
	obj.Slots = GetShipDesignSlotArray(o.Get("slots"))
	obj.Spec = GetShipDesignSpec(o.Get("spec"))
	return obj
}

func GetShipDesignIntelArray(o js.Value) []cs.ShipDesignIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.ShipDesignIntel, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetShipDesignIntel(o.Index(i))
	}
	return items
}

func GetShipDesignIntelPointerArray(o js.Value) []*cs.ShipDesignIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.ShipDesignIntel, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetShipDesignIntel(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetShipDesignSlot(o js.Value) cs.ShipDesignSlot {
	obj := cs.ShipDesignSlot{}
	obj.HullComponent = GetString(o, "hullComponent")
	obj.HullSlotIndex = GetInt[int](o, "hullSlotIndex")
	obj.Quantity = GetInt[int](o, "quantity")
	return obj
}

func GetShipDesignSlotArray(o js.Value) []cs.ShipDesignSlot {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.ShipDesignSlot, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetShipDesignSlot(o.Index(i))
	}
	return items
}

func GetShipDesignSlotPointerArray(o js.Value) []*cs.ShipDesignSlot {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.ShipDesignSlot, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetShipDesignSlot(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetShipDesignSpec(o js.Value) cs.ShipDesignSpec {
	obj := cs.ShipDesignSpec{}
	obj.AdditionalMassDrivers = GetInt[int](o, "additionalMassDrivers")
	obj.Armor = GetInt[int](o, "armor")
	obj.BasePacketSpeed = GetInt[int](o, "basePacketSpeed")
	obj.Bomber = GetBool(o, "bomber")
	obj.Bombs = GetBombArray(o.Get("bombs"))
	obj.CanJump = GetBool(o, "canJump")
	obj.CanLayMines = GetBool(o, "canLayMines")
	obj.CanStealFleetCargo = GetBool(o, "canStealFleetCargo")
	obj.CanStealPlanetCargo = GetBool(o, "canStealPlanetCargo")
	obj.CargoCapacity = GetInt[int](o, "cargoCapacity")
	obj.CloakPercent = GetInt[int](o, "cloakPercent")
	obj.CloakPercentFullCargo = GetInt[int](o, "cloakPercentFullCargo")
	obj.CloakUnits = GetInt[int](o, "cloakUnits")
	obj.Colonizer = GetBool(o, "colonizer")
	obj.Cost = GetCost(o.Get("cost"))
	obj.Engine = GetEngine(o.Get("engine"))
	obj.EstimatedRange = GetInt[int](o, "estimatedRange")
	obj.EstimatedRangeFull = GetInt[int](o, "estimatedRangeFull")
	obj.FuelCapacity = GetInt[int](o, "fuelCapacity")
	obj.FuelGeneration = GetInt[int](o, "fuelGeneration")
	obj.HasWeapons = GetBool(o, "hasWeapons")
	obj.HullType = cs.TechHullType(GetString(o, "hullType"))
	obj.ImmuneToOwnDetonation = GetBool(o, "immuneToOwnDetonation")
	obj.Initiative = GetInt[int](o, "initiative")
	obj.Mass = GetInt[int](o, "mass")
	obj.MassDriver = GetString(o, "massDriver")
	obj.MaxHullMass = GetInt[int](o, "maxHullMass")
	obj.MaxPopulation = GetInt[int](o, "maxPopulation")
	obj.MaxRange = GetInt[int](o, "maxRange")
	obj.MineSweep = GetInt[int](o, "mineSweep")
	obj.MiningRate = GetInt[int](o, "miningRate")
	obj.Movement = GetInt[int](o, "movement")
	obj.MovementBonus = GetInt[int](o, "movementBonus")
	obj.MovementFull = GetInt[int](o, "movementFull")
	obj.NumBuilt = GetInt[int](o, "numBuilt")
	obj.NumEngines = GetInt[int](o, "numEngines")
	obj.NumInstances = GetInt[int](o, "numInstances")
	obj.OrbitalConstructionModule = GetBool(o, "orbitalConstructionModule")
	obj.PowerRating = GetInt[int](o, "powerRating")
	obj.Radiating = GetBool(o, "radiating")
	obj.ReduceMovement = GetInt[int](o, "reduceMovement")
	obj.RetroBombs = GetBombArray(o.Get("retroBombs"))
	obj.SafeHullMass = GetInt[int](o, "safeHullMass")
	obj.SafePacketSpeed = GetInt[int](o, "safePacketSpeed")
	obj.SafeRange = GetInt[int](o, "safeRange")
	obj.Scanner = GetBool(o, "scanner")
	obj.ScanRange = GetInt[int](o, "scanRange")
	obj.ScanRangePen = GetInt[int](o, "scanRangePen")
	obj.Shields = GetInt[int](o, "shields")
	obj.SmartBombs = GetBombArray(o.Get("smartBombs"))
	obj.SpaceDock = GetInt[int](o, "spaceDock")
	obj.Starbase = GetBool(o, "starbase")
	obj.Stargate = GetString(o, "stargate")
	obj.TechLevel = GetTechLevel(o.Get("techLevel"))
	obj.TerraformRate = GetInt[int](o, "terraformRate")
	obj.WeaponSlots = GetShipDesignSlotArray(o.Get("weaponSlots"))
	return obj
}

func GetShipDesignSpecArray(o js.Value) []cs.ShipDesignSpec {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.ShipDesignSpec, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetShipDesignSpec(o.Index(i))
	}
	return items
}

func GetShipDesignSpecPointerArray(o js.Value) []*cs.ShipDesignSpec {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.ShipDesignSpec, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetShipDesignSpec(o.Index(i))
		items[i] = &item
	}
	return items
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

func GetTechLevelArray(o js.Value) []cs.TechLevel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.TechLevel, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetTechLevel(o.Index(i))
	}
	return items
}

func GetTechLevelPointerArray(o js.Value) []*cs.TechLevel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.TechLevel, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetTechLevel(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetTransportPlan(o js.Value) cs.TransportPlan {
	obj := cs.TransportPlan{}
	obj.Num = GetInt[int](o, "num")
	obj.Name = GetString(o, "name")
	return obj
}

func GetTransportPlanArray(o js.Value) []cs.TransportPlan {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.TransportPlan, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetTransportPlan(o.Index(i))
	}
	return items
}

func GetTransportPlanPointerArray(o js.Value) []*cs.TransportPlan {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.TransportPlan, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetTransportPlan(o.Index(i))
		items[i] = &item
	}
	return items
}

func GetWormholeIntel(o js.Value) cs.WormholeIntel {
	obj := cs.WormholeIntel{}
	obj.DestinationNum = GetInt[int](o, "destinationNum")
	obj.Stability = cs.WormholeStability(GetString(o, "stability"))
	return obj
}

func GetWormholeIntelArray(o js.Value) []cs.WormholeIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]cs.WormholeIntel, o.Length())	
	for i := 0; i < len(items); i++ {
		items[i] = GetWormholeIntel(o.Index(i))
	}
	return items
}

func GetWormholeIntelPointerArray(o js.Value) []*cs.WormholeIntel {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*cs.WormholeIntel, o.Length())
	for i := 0; i < len(items); i++ {
		item := GetWormholeIntel(o.Index(i))
		items[i] = &item
	}
	return items
}
