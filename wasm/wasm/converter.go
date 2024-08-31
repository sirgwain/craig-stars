//go:build wasi || wasm

package wasm

import (
	"fmt"
	"strconv"
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

func GetFloat[T ~float32 | ~float64](o js.Value, key string) T {
	var result T
	val := o.Get(key)
	if !val.IsUndefined() {
		result = T(val.Float())
	}
	return result
}

func GetIntArray[T ~uint | ~uint32 | ~uint64 | ~int | ~int32 | ~int64](o js.Value, key string) []T {
	val := o.Get(key)
	items := make([]T, val.Length())
	for i := 0; i < len(items); i++ {
		items[i] = T(val.Index(i).Int())
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

// GetSlice populates an array with a getter function
func GetSlice[T any](o js.Value, getter func(o js.Value) T) []T {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]T, o.Length())
	for i := 0; i < len(items); i++ {
		items[i] = getter(o.Index(i))
	}
	return items
}

// GetSliceSlice populates a 2d array
func GetSliceSlice[T any](o js.Value, getter func(o js.Value) T) [][]T {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([][]T, o.Length())
	for i := 0; i < len(items); i++ {
		items[i] = GetSlice(o.Index(i), getter)
	}
	return items
}

func GetBattlePlan(o js.Value) cs.BattlePlan {
	obj := cs.BattlePlan{}
    obj.Num = GetInt[int](o, "num")
    obj.Name = string(GetString(o, "name"))
    obj.DumpCargo = bool(GetBool(o, "dumpCargo"))
	return obj
}

func GetBattleRecord(o js.Value) cs.BattleRecord {
	obj := cs.BattleRecord{}
    obj.Num = GetInt[int](o, "num")
    obj.PlanetNum = GetInt[int](o, "planetNum")
    obj.Position = GetVector(o.Get("position"))
    obj.Tokens = GetSlice(o.Get("tokens"), GetBattleRecordToken)
    obj.ActionsPerRound = GetSliceSlice(o.Get("actionsPerRound"), GetBattleRecordTokenAction)
    obj.DestroyedTokens = GetSlice(o.Get("destroyedTokens"), GetBattleRecordDestroyedToken)
    obj.Stats = GetBattleRecordStats(o.Get("stats"))
	return obj
}

func GetBattleRecordDestroyedToken(o js.Value) cs.BattleRecordDestroyedToken {
	obj := cs.BattleRecordDestroyedToken{}
    obj.Num = GetInt[int](o, "num")
    obj.PlayerNum = GetInt[int](o, "playerNum")
    obj.DesignNum = GetInt[int](o, "designNum")
    obj.Quantity = GetInt[int](o, "quantity")
	return obj
}

func GetBattleRecordStats(o js.Value) cs.BattleRecordStats {
	obj := cs.BattleRecordStats{}
    obj.NumPlayers = GetInt[int](o, "numPlayers")
  
  	numShipsByPlayerObj := o.Get("numShipsByPlayer")
  	if !numShipsByPlayerObj.IsUndefined() {
  		numShipsByPlayer := make(map[int]int)
  		numShipsByPlayerKeys := js.Global().Get("Object").Call("keys", numShipsByPlayerObj)
  		for i := 0; i < numShipsByPlayerKeys.Length(); i++ {
  			key, _ := strconv.Atoi(numShipsByPlayerKeys.Index(i).String())
  			numShipsByPlayer[key] = GetInt[int](numShipsByPlayerObj, fmt.Sprintf("%v", key))
  		}
  		obj.NumShipsByPlayer = numShipsByPlayer
  	}
  
  	shipsDestroyedByPlayerObj := o.Get("shipsDestroyedByPlayer")
  	if !shipsDestroyedByPlayerObj.IsUndefined() {
  		shipsDestroyedByPlayer := make(map[int]int)
  		shipsDestroyedByPlayerKeys := js.Global().Get("Object").Call("keys", shipsDestroyedByPlayerObj)
  		for i := 0; i < shipsDestroyedByPlayerKeys.Length(); i++ {
  			key, _ := strconv.Atoi(shipsDestroyedByPlayerKeys.Index(i).String())
  			shipsDestroyedByPlayer[key] = GetInt[int](shipsDestroyedByPlayerObj, fmt.Sprintf("%v", key))
  		}
  		obj.ShipsDestroyedByPlayer = shipsDestroyedByPlayer
  	}
  
  	damageTakenByPlayerObj := o.Get("damageTakenByPlayer")
  	if !damageTakenByPlayerObj.IsUndefined() {
  		damageTakenByPlayer := make(map[int]int)
  		damageTakenByPlayerKeys := js.Global().Get("Object").Call("keys", damageTakenByPlayerObj)
  		for i := 0; i < damageTakenByPlayerKeys.Length(); i++ {
  			key, _ := strconv.Atoi(damageTakenByPlayerKeys.Index(i).String())
  			damageTakenByPlayer[key] = GetInt[int](damageTakenByPlayerObj, fmt.Sprintf("%v", key))
  		}
  		obj.DamageTakenByPlayer = damageTakenByPlayer
  	}
	return obj
}

func GetBattleRecordToken(o js.Value) cs.BattleRecordToken {
	obj := cs.BattleRecordToken{}
    obj.Num = GetInt[int](o, "num")
    obj.PlayerNum = GetInt[int](o, "playerNum")
    obj.DesignNum = GetInt[int](o, "designNum")
    obj.Position = GetBattleVector(o.Get("position"))
    obj.Initiative = GetInt[int](o, "initiative")
    obj.Mass = GetInt[int](o, "mass")
    obj.Armor = GetInt[int](o, "armor")
    obj.StackShields = GetInt[int](o, "stackShields")
    obj.Movement = GetInt[int](o, "movement")
    obj.StartingQuantity = GetInt[int](o, "startingQuantity")
	return obj
}

func GetBattleRecordTokenAction(o js.Value) cs.BattleRecordTokenAction {
	obj := cs.BattleRecordTokenAction{}
    obj.TokenNum = GetInt[int](o, "tokenNum")
    obj.Round = GetInt[int](o, "round")
    obj.From = GetBattleVector(o.Get("from"))
    obj.To = GetBattleVector(o.Get("to"))
    obj.Slot = GetInt[int](o, "slot")
    obj.TargetNum = GetInt[int](o, "targetNum")
    targetVal := o.Get("target")
    			if !targetVal.IsUndefined() {
    				target := GetShipToken(targetVal)
    				obj.Target = &target
    			}
    obj.TokensDestroyed = GetInt[int](o, "tokensDestroyed")
    obj.DamageDoneShields = GetInt[int](o, "damageDoneShields")
    obj.DamageDoneArmor = GetInt[int](o, "damageDoneArmor")
    obj.TorpedoHits = GetInt[int](o, "torpedoHits")
    obj.TorpedoMisses = GetInt[int](o, "torpedoMisses")
	return obj
}

func GetBattleVector(o js.Value) cs.BattleVector {
	obj := cs.BattleVector{}
    obj.X = GetInt[int](o, "x")
    obj.Y = GetInt[int](o, "y")
	return obj
}

func GetBomb(o js.Value) cs.Bomb {
	obj := cs.Bomb{}
    obj.Quantity = GetInt[int](o, "quantity")
    obj.KillRate = GetFloat[float64](o, "killRate")
    obj.MinKillRate = GetInt[int](o, "minKillRate")
    obj.StructureDestroyRate = GetFloat[float64](o, "structureDestroyRate")
    obj.UnterraformRate = GetInt[int](o, "unterraformRate")
	return obj
}

func GetBombingResult(o js.Value) cs.BombingResult {
	obj := cs.BombingResult{}
    obj.BomberName = string(GetString(o, "bomberName"))
    obj.NumBombers = GetInt[int](o, "numBombers")
    obj.ColonistsKilled = GetInt[int](o, "colonistsKilled")
    obj.MinesDestroyed = GetInt[int](o, "minesDestroyed")
    obj.FactoriesDestroyed = GetInt[int](o, "factoriesDestroyed")
    obj.DefensesDestroyed = GetInt[int](o, "defensesDestroyed")
    obj.UnterraformAmount = GetHab(o.Get("unterraformAmount"))
    obj.PlanetEmptied = bool(GetBool(o, "planetEmptied"))
	return obj
}

func GetCargo(o js.Value) cs.Cargo {
	obj := cs.Cargo{}
    obj.Ironium = GetInt[int](o, "ironium")
    obj.Boranium = GetInt[int](o, "boranium")
    obj.Germanium = GetInt[int](o, "germanium")
    obj.Colonists = GetInt[int](o, "colonists")
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

func GetDefense(o js.Value) cs.Defense {
	obj := cs.Defense{}
    obj.DefenseCoverage = GetFloat[float64](o, "defenseCoverage")
	return obj
}

func GetEngine(o js.Value) cs.Engine {
	obj := cs.Engine{}
    obj.IdealSpeed = GetInt[int](o, "idealSpeed")
    obj.FreeSpeed = GetInt[int](o, "freeSpeed")
    obj.MaxSafeSpeed = GetInt[int](o, "maxSafeSpeed")
    if !o.Get("fuelUsage").IsUndefined() && o.Get("fuelUsage").Length() != 0 {
    			obj.FuelUsage = [11]int(GetIntArray[int](o, "fuelUsage"))
    		}
	return obj
}

func GetFleet(o js.Value) cs.Fleet {
	obj := cs.Fleet{}
    obj.MapObject = GetMapObject(o)
    obj.FleetOrders = GetFleetOrders(o)
    obj.PlanetNum = GetInt[int](o, "planetNum")
    obj.BaseName = string(GetString(o, "baseName"))
    obj.Cargo = GetCargo(o.Get("cargo"))
    obj.Fuel = GetInt[int](o, "fuel")
    obj.Age = GetInt[int](o, "age")
    obj.Tokens = GetSlice(o.Get("tokens"), GetShipToken)
    obj.Heading = GetVector(o.Get("heading"))
    obj.WarpSpeed = GetInt[int](o, "warpSpeed")
    previousPositionVal := o.Get("previousPosition")
    			if !previousPositionVal.IsUndefined() {
    				previousPosition := GetVector(previousPositionVal)
    				obj.PreviousPosition = &previousPosition
    			}
    obj.OrbitingPlanetNum = GetInt[int](o, "orbitingPlanetNum")
    obj.Starbase = bool(GetBool(o, "starbase"))
    obj.Spec = GetFleetSpec(o.Get("spec"))
	return obj
}

func GetFleetIntel(o js.Value) cs.FleetIntel {
	obj := cs.FleetIntel{}
    obj.MapObjectIntel = GetMapObjectIntel(o)
    obj.PlanetIntelID = GetInt[int64](o, "")
    obj.BaseName = string(GetString(o, "baseName"))
    obj.Heading = GetVector(o.Get("heading"))
    obj.OrbitingPlanetNum = GetInt[int](o, "orbitingPlanetNum")
    obj.WarpSpeed = GetInt[int](o, "warpSpeed")
    obj.Mass = GetInt[int](o, "mass")
    obj.Cargo = GetCargo(o.Get("cargo"))
    obj.CargoDiscovered = bool(GetBool(o, "cargoDiscovered"))
    obj.Freighter = bool(GetBool(o, "freighter"))
    obj.ScanRange = GetInt[int](o, "scanRange")
    obj.ScanRangePen = GetInt[int](o, "scanRangePen")
    obj.Tokens = GetSlice(o.Get("tokens"), GetShipToken)
	return obj
}

func GetFleetOrders(o js.Value) cs.FleetOrders {
	obj := cs.FleetOrders{}
    obj.Waypoints = GetSlice(o.Get("waypoints"), GetWaypoint)
    obj.RepeatOrders = bool(GetBool(o, "repeatOrders"))
    obj.BattlePlanNum = GetInt[int](o, "battlePlanNum")
	return obj
}

func GetFleetSpec(o js.Value) cs.FleetSpec {
	obj := cs.FleetSpec{}
    obj.ShipDesignSpec = GetShipDesignSpec(o)
    obj.BaseCloakedCargo = GetInt[int](o, "baseCloakedCargo")
    obj.BasePacketSpeed = GetInt[int](o, "basePacketSpeed")
    obj.HasMassDriver = bool(GetBool(o, "hasMassDriver"))
    obj.HasStargate = bool(GetBool(o, "hasStargate"))
    obj.MassDriver = string(GetString(o, "massDriver"))
    obj.MassEmpty = GetInt[int](o, "massEmpty")
    obj.MaxHullMass = GetInt[int](o, "maxHullMass")
    obj.MaxRange = GetInt[int](o, "maxRange")
    obj.SafeHullMass = GetInt[int](o, "safeHullMass")
    obj.SafeRange = GetInt[int](o, "safeRange")
    obj.Stargate = string(GetString(o, "stargate"))
    obj.TotalShips = GetInt[int](o, "totalShips")
	return obj
}

func GetGameDBObject(o js.Value) cs.GameDBObject {
	obj := cs.GameDBObject{}
    obj.ID = GetInt[int64](o, "id")
    obj.GameID = GetInt[int64](o, "gameId")
    obj.CreatedAt, _ = GetTime(o, "createdAt")
    obj.UpdatedAt, _ = GetTime(o, "updatedAt")
	return obj
}

func GetHab(o js.Value) cs.Hab {
	obj := cs.Hab{}
    obj.Grav = GetInt[int](o, "grav")
    obj.Temp = GetInt[int](o, "temp")
    obj.Rad = GetInt[int](o, "rad")
	return obj
}

func GetIntel(o js.Value) cs.Intel {
	obj := cs.Intel{}
    obj.Name = string(GetString(o, "name"))
    obj.Num = GetInt[int](o, "num")
    obj.PlayerNum = GetInt[int](o, "playerNum")
    obj.ReportAge = GetInt[int](o, "reportAge")
	return obj
}

func GetMapObject(o js.Value) cs.MapObject {
	obj := cs.MapObject{}
    obj.GameDBObject = GetGameDBObject(o)
    obj.Position = GetVector(o.Get("position"))
    obj.Num = GetInt[int](o, "num")
    obj.PlayerNum = GetInt[int](o, "playerNum")
    obj.Name = string(GetString(o, "name"))
	return obj
}

func GetMapObjectIntel(o js.Value) cs.MapObjectIntel {
	obj := cs.MapObjectIntel{}
    obj.Intel = GetIntel(o)
    obj.Position = GetVector(o.Get("position"))
	return obj
}

func GetMineField(o js.Value) cs.MineField {
	obj := cs.MineField{}
    obj.MapObject = GetMapObject(o)
    obj.MineFieldOrders = GetMineFieldOrders(o)
    obj.NumMines = GetInt[int](o, "numMines")
    obj.Spec = GetMineFieldSpec(o.Get("spec"))
	return obj
}

func GetMineFieldIntel(o js.Value) cs.MineFieldIntel {
	obj := cs.MineFieldIntel{}
    obj.MapObjectIntel = GetMapObjectIntel(o)
    obj.NumMines = GetInt[int](o, "numMines")
    obj.Spec = GetMineFieldSpec(o.Get("spec"))
	return obj
}

func GetMineFieldOrders(o js.Value) cs.MineFieldOrders {
	obj := cs.MineFieldOrders{}
    obj.Detonate = bool(GetBool(o, "detonate"))
	return obj
}

func GetMineFieldSpec(o js.Value) cs.MineFieldSpec {
	obj := cs.MineFieldSpec{}
    obj.Radius = GetFloat[float64](o, "radius")
    obj.DecayRate = GetInt[int](o, "decayRate")
	return obj
}

func GetMineral(o js.Value) cs.Mineral {
	obj := cs.Mineral{}
    obj.Ironium = GetInt[int](o, "ironium")
    obj.Boranium = GetInt[int](o, "boranium")
    obj.Germanium = GetInt[int](o, "germanium")
	return obj
}

func GetMineralPacketDamage(o js.Value) cs.MineralPacketDamage {
	obj := cs.MineralPacketDamage{}
    obj.Killed = GetInt[int](o, "killed")
    obj.DefensesDestroyed = GetInt[int](o, "defensesDestroyed")
    obj.Uncaught = GetInt[int](o, "uncaught")
	return obj
}

func GetMineralPacketIntel(o js.Value) cs.MineralPacketIntel {
	obj := cs.MineralPacketIntel{}
    obj.MapObjectIntel = GetMapObjectIntel(o)
    obj.WarpSpeed = GetInt[int](o, "warpSpeed")
    obj.Heading = GetVector(o.Get("heading"))
    obj.Cargo = GetCargo(o.Get("cargo"))
    obj.TargetPlanetNum = GetInt[int](o, "targetPlanetNum")
    obj.ScanRange = GetInt[int](o, "scanRange")
    obj.ScanRangePen = GetInt[int](o, "scanRangePen")
	return obj
}

func GetMysteryTrader(o js.Value) cs.MysteryTrader {
	obj := cs.MysteryTrader{}
    obj.MapObject = GetMapObject(o)
    obj.WarpSpeed = GetInt[int](o, "warpSpeed")
    obj.Destination = GetVector(o.Get("destination"))
    obj.RequestedBoon = GetInt[int](o, "requestedBoon")
    obj.Heading = GetVector(o.Get("heading"))
  
  	playersRewardedObj := o.Get("playersRewarded")
  	if !playersRewardedObj.IsUndefined() {
  		playersRewarded := make(map[int]bool)
  		playersRewardedKeys := js.Global().Get("Object").Call("keys", playersRewardedObj)
  		for i := 0; i < playersRewardedKeys.Length(); i++ {
  			key, _ := strconv.Atoi(playersRewardedKeys.Index(i).String())
  			playersRewarded[key] = bool(GetBool(playersRewardedObj, fmt.Sprintf("%v", key)))
  		}
  		obj.PlayersRewarded = playersRewarded
  	}
    obj.Spec = GetMysteryTraderSpec(o.Get("spec"))
	return obj
}

func GetMysteryTraderIntel(o js.Value) cs.MysteryTraderIntel {
	obj := cs.MysteryTraderIntel{}
    obj.MapObjectIntel = GetMapObjectIntel(o)
    obj.WarpSpeed = GetInt[int](o, "warpSpeed")
    obj.Heading = GetVector(o.Get("heading"))
    obj.RequestedBoon = GetInt[int](o, "requestedBoon")
	return obj
}

func GetMysteryTraderReward(o js.Value) cs.MysteryTraderReward {
	obj := cs.MysteryTraderReward{}
    obj.TechLevels = GetTechLevel(o.Get("techLevels"))
    obj.Tech = string(GetString(o, "tech"))
    obj.Ship = GetShipDesign(o.Get("ship"))
    obj.ShipCount = GetInt[int](o, "shipCount")
	return obj
}

func GetMysteryTraderSpec(o js.Value) cs.MysteryTraderSpec {
	obj := cs.MysteryTraderSpec{}
	return obj
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
    obj.Homeworld = bool(GetBool(o, "homeworld"))
    obj.Scanner = bool(GetBool(o, "scanner"))
    obj.Spec = GetPlanetSpec(o.Get("spec"))
	return obj
}

func GetPlanetIntel(o js.Value) cs.PlanetIntel {
	obj := cs.PlanetIntel{}
    obj.MapObjectIntel = GetMapObjectIntel(o)
    obj.Hab = GetHab(o.Get("hab"))
    obj.BaseHab = GetHab(o.Get("baseHab"))
    obj.MineralConcentration = GetMineral(o.Get("mineralConcentration"))
    starbaseVal := o.Get("starbase")
    			if !starbaseVal.IsUndefined() {
    				starbase := GetFleetIntel(starbaseVal)
    				obj.Starbase = &starbase
    			}
    obj.Cargo = GetCargo(o.Get("cargo"))
    obj.CargoDiscovered = bool(GetBool(o, "cargoDiscovered"))
    obj.PlanetHabitability = GetInt[int](o, "planetHabitability")
    obj.PlanetHabitabilityTerraformed = GetInt[int](o, "planetHabitabilityTerraformed")
    obj.Homeworld = bool(GetBool(o, "homeworld"))
    obj.Spec = GetPlanetSpec(o.Get("spec"))
	return obj
}

func GetPlanetOrders(o js.Value) cs.PlanetOrders {
	obj := cs.PlanetOrders{}
    obj.ContributesOnlyLeftoverToResearch = bool(GetBool(o, "contributesOnlyLeftoverToResearch"))
    obj.ProductionQueue = GetSlice(o.Get("productionQueue"), GetProductionQueueItem)
    obj.RouteTargetNum = GetInt[int](o, "routeTargetNum")
    obj.RouteTargetPlayerNum = GetInt[int](o, "routeTargetPlayerNum")
    obj.PacketTargetNum = GetInt[int](o, "packetTargetNum")
    obj.PacketSpeed = GetInt[int](o, "packetSpeed")
	return obj
}

func GetPlanetSpec(o js.Value) cs.PlanetSpec {
	obj := cs.PlanetSpec{}
    obj.PlanetStarbaseSpec = GetPlanetStarbaseSpec(o)
    obj.CanTerraform = bool(GetBool(o, "canTerraform"))
    obj.Defense = string(GetString(o, "defense"))
    obj.DefenseCoverage = GetFloat[float64](o, "defenseCoverage")
    obj.DefenseCoverageSmart = GetFloat[float64](o, "defenseCoverageSmart")
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
    obj.PopulationDensity = GetFloat[float64](o, "populationDensity")
    obj.ResourcesPerYear = GetInt[int](o, "resourcesPerYear")
    obj.ResourcesPerYearAvailable = GetInt[int](o, "resourcesPerYearAvailable")
    obj.ResourcesPerYearResearch = GetInt[int](o, "resourcesPerYearResearch")
    obj.ResourcesPerYearResearchEstimatedLeftover = GetInt[int](o, "resourcesPerYearResearchEstimatedLeftover")
    obj.Scanner = string(GetString(o, "scanner"))
    obj.ScanRange = GetInt[int](o, "scanRange")
    obj.ScanRangePen = GetInt[int](o, "scanRangePen")
    obj.TerraformAmount = GetHab(o.Get("terraformAmount"))
    obj.MinTerraformAmount = GetHab(o.Get("minTerraformAmount"))
    obj.TerraformedHabitability = GetInt[int](o, "terraformedHabitability")
    obj.Contested = bool(GetBool(o, "contested"))
	return obj
}

func GetPlanetStarbaseSpec(o js.Value) cs.PlanetStarbaseSpec {
	obj := cs.PlanetStarbaseSpec{}
    obj.HasMassDriver = bool(GetBool(o, "hasMassDriver"))
    obj.HasStarbase = bool(GetBool(o, "hasStarbase"))
    obj.HasStargate = bool(GetBool(o, "hasStargate"))
    obj.StarbaseDesignName = string(GetString(o, "starbaseDesignName"))
    obj.StarbaseDesignNum = GetInt[int](o, "starbaseDesignNum")
    obj.DockCapacity = GetInt[int](o, "dockCapacity")
    obj.BasePacketSpeed = GetInt[int](o, "basePacketSpeed")
    obj.SafePacketSpeed = GetInt[int](o, "safePacketSpeed")
    obj.SafeHullMass = GetInt[int](o, "safeHullMass")
    obj.SafeRange = GetInt[int](o, "safeRange")
    obj.MaxRange = GetInt[int](o, "maxRange")
    obj.MaxHullMass = GetInt[int](o, "maxHullMass")
    obj.Stargate = string(GetString(o, "stargate"))
    obj.MassDriver = string(GetString(o, "massDriver"))
	return obj
}

func GetPlayer(o js.Value) cs.Player {
	obj := cs.Player{}
    obj.GameDBObject = GetGameDBObject(o)
    obj.PlayerOrders = GetPlayerOrders(o)
    obj.PlayerIntels = GetPlayerIntels(o)
    obj.PlayerPlans = GetPlayerPlans(o)
    obj.UserID = GetInt[int64](o, "userId")
    obj.Name = string(GetString(o, "name"))
    obj.Num = GetInt[int](o, "num")
    obj.Ready = bool(GetBool(o, "ready"))
    obj.AIControlled = bool(GetBool(o, "aiControlled"))
    obj.Guest = bool(GetBool(o, "guest"))
    obj.SubmittedTurn = bool(GetBool(o, "submittedTurn"))
    obj.Color = string(GetString(o, "color"))
    obj.DefaultHullSet = GetInt[int](o, "defaultHullSet")
    obj.Race = GetRace(o.Get("race"))
    obj.TechLevels = GetTechLevel(o.Get("techLevels"))
    obj.TechLevelsSpent = GetTechLevel(o.Get("techLevelsSpent"))
    obj.ResearchSpentLastYear = GetInt[int](o, "researchSpentLastYear")
    obj.Relations = GetSlice(o.Get("relations"), GetPlayerRelationship)
    obj.Messages = GetSlice(o.Get("messages"), GetPlayerMessage)
    designs := GetSlice(o.Get("designs"), GetShipDesign)
    			obj.Designs = make([]*cs.ShipDesign, len(designs))
    			for i := range designs {
    				obj.Designs[i] = &designs[i]
    			}
    obj.ScoreHistory = GetSlice(o.Get("scoreHistory"), GetPlayerScore)
  
  	acquiredTechsObj := o.Get("acquiredTechs")
  	if !acquiredTechsObj.IsUndefined() {
  		acquiredTechs := make(map[string]bool)
  		acquiredTechsKeys := js.Global().Get("Object").Call("keys", acquiredTechsObj)
  		for i := 0; i < acquiredTechsKeys.Length(); i++ {
  			key := string(acquiredTechsKeys.Index(i).String())
  			acquiredTechs[key] = bool(GetBool(acquiredTechsObj, fmt.Sprintf("%v", key)))
  		}
  		obj.AcquiredTechs = acquiredTechs
  	}
    obj.Victor = bool(GetBool(o, "victor"))
    statsVal := o.Get("stats")
    			if !statsVal.IsUndefined() {
    				stats := GetPlayerStats(statsVal)
    				obj.Stats = &stats
    			}
    obj.Spec = GetPlayerSpec(o.Get("spec"))
	return obj
}

func GetPlayerIntel(o js.Value) cs.PlayerIntel {
	obj := cs.PlayerIntel{}
    obj.Name = string(GetString(o, "name"))
    obj.Num = GetInt[int](o, "num")
    obj.Color = string(GetString(o, "color"))
    obj.Seen = bool(GetBool(o, "seen"))
    obj.RaceName = string(GetString(o, "raceName"))
    obj.RacePluralName = string(GetString(o, "racePluralName"))
	return obj
}

func GetPlayerIntels(o js.Value) cs.PlayerIntels {
	obj := cs.PlayerIntels{}
    obj.BattleRecords = GetSlice(o.Get("battleRecords"), GetBattleRecord)
    obj.PlayerIntels = GetSlice(o.Get("playerIntels"), GetPlayerIntel)
    obj.ScoreIntels = GetSlice(o.Get("scoreIntels"), GetScoreIntel)
    obj.PlanetIntels = GetSlice(o.Get("planetIntels"), GetPlanetIntel)
    obj.FleetIntels = GetSlice(o.Get("fleetIntels"), GetFleetIntel)
    obj.StarbaseIntels = GetSlice(o.Get("starbaseIntels"), GetFleetIntel)
    obj.ShipDesignIntels = GetSlice(o.Get("shipDesignIntels"), GetShipDesignIntel)
    obj.MineralPacketIntels = GetSlice(o.Get("mineralPacketIntels"), GetMineralPacketIntel)
    obj.MineFieldIntels = GetSlice(o.Get("mineFieldIntels"), GetMineFieldIntel)
    obj.WormholeIntels = GetSlice(o.Get("wormholeIntels"), GetWormholeIntel)
    obj.MysteryTraderIntels = GetSlice(o.Get("mysteryTraderIntels"), GetMysteryTraderIntel)
    obj.SalvageIntels = GetSlice(o.Get("salvageIntels"), GetSalvageIntel)
	return obj
}

func GetPlayerMessage(o js.Value) cs.PlayerMessage {
	obj := cs.PlayerMessage{}
    obj.Text = string(GetString(o, "text"))
    obj.BattleNum = GetInt[int](o, "battleNum")
    obj.Spec = GetPlayerMessageSpec(o.Get("spec"))
	return obj
}

func GetPlayerMessageSpec(o js.Value) cs.PlayerMessageSpec {
	obj := cs.PlayerMessageSpec{}
    obj.Amount = GetInt[int](o, "amount")
    obj.Amount2 = GetInt[int](o, "amount2")
    obj.PrevAmount = GetInt[int](o, "prevAmount")
    obj.SourcePlayerNum = GetInt[int](o, "sourcePlayerNum")
    obj.DestPlayerNum = GetInt[int](o, "destPlayerNum")
    obj.Name = string(GetString(o, "name"))
    costVal := o.Get("cost")
    			if !costVal.IsUndefined() {
    				cost := GetCost(costVal)
    				obj.Cost = &cost
    			}
    mineralVal := o.Get("mineral")
    			if !mineralVal.IsUndefined() {
    				mineral := GetMineral(mineralVal)
    				obj.Mineral = &mineral
    			}
    cargoVal := o.Get("cargo")
    			if !cargoVal.IsUndefined() {
    				cargo := GetCargo(cargoVal)
    				obj.Cargo = &cargo
    			}
    obj.TechGained = string(GetString(o, "techGained"))
    obj.Battle = GetBattleRecordStats(o.Get("battle"))
    cometVal := o.Get("comet")
    			if !cometVal.IsUndefined() {
    				comet := GetPlayerMessageSpecComet(cometVal)
    				obj.Comet = &comet
    			}
    bombingVal := o.Get("bombing")
    			if !bombingVal.IsUndefined() {
    				bombing := GetBombingResult(bombingVal)
    				obj.Bombing = &bombing
    			}
    mineralPacketDamageVal := o.Get("mineralPacketDamage")
    			if !mineralPacketDamageVal.IsUndefined() {
    				mineralPacketDamage := GetMineralPacketDamage(mineralPacketDamageVal)
    				obj.MineralPacketDamage = &mineralPacketDamage
    			}
    mysteryTraderVal := o.Get("mysteryTrader")
    			if !mysteryTraderVal.IsUndefined() {
    				mysteryTrader := GetPlayerMessageSpecMysteryTrader(mysteryTraderVal)
    				obj.MysteryTrader = &mysteryTrader
    			}
	return obj
}

func GetPlayerMessageSpecComet(o js.Value) cs.PlayerMessageSpecComet {
	obj := cs.PlayerMessageSpecComet{}
    obj.MineralsAdded = GetMineral(o.Get("mineralsAdded"))
    obj.MineralConcentrationIncreased = GetMineral(o.Get("mineralConcentrationIncreased"))
    obj.HabChanged = GetHab(o.Get("habChanged"))
    obj.ColonistsKilled = GetInt[int](o, "colonistsKilled")
	return obj
}

func GetPlayerMessageSpecMysteryTrader(o js.Value) cs.PlayerMessageSpecMysteryTrader {
	obj := cs.PlayerMessageSpecMysteryTrader{}
    obj.MysteryTraderReward = GetMysteryTraderReward(o)
    obj.FleetNum = GetInt[int](o, "fleetNum")
	return obj
}

func GetPlayerOrders(o js.Value) cs.PlayerOrders {
	obj := cs.PlayerOrders{}
    obj.ResearchAmount = GetInt[int](o, "researchAmount")
	return obj
}

func GetPlayerPlans(o js.Value) cs.PlayerPlans {
	obj := cs.PlayerPlans{}
    obj.ProductionPlans = GetSlice(o.Get("productionPlans"), GetProductionPlan)
    obj.BattlePlans = GetSlice(o.Get("battlePlans"), GetBattlePlan)
    obj.TransportPlans = GetSlice(o.Get("transportPlans"), GetTransportPlan)
	return obj
}

func GetPlayerRelationship(o js.Value) cs.PlayerRelationship {
	obj := cs.PlayerRelationship{}
    obj.ShareMap = bool(GetBool(o, "shareMap"))
	return obj
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
	return obj
}

func GetPlayerSpec(o js.Value) cs.PlayerSpec {
	obj := cs.PlayerSpec{}
    obj.PlanetaryScanner = GetTechPlanetaryScanner(o.Get("planetaryScanner"))
    obj.Defense = GetTechDefense(o.Get("defense"))
    obj.ResourcesPerYear = GetInt[int](o, "resourcesPerYear")
    obj.ResourcesPerYearResearch = GetInt[int](o, "resourcesPerYearResearch")
    obj.ResourcesPerYearResearchEstimated = GetInt[int](o, "resourcesPerYearResearchEstimated")
    obj.CurrentResearchCost = GetInt[int](o, "currentResearchCost")
	return obj
}

func GetPlayerStats(o js.Value) cs.PlayerStats {
	obj := cs.PlayerStats{}
    obj.FleetsBuilt = GetInt[int](o, "fleetsBuilt")
    obj.StarbasesBuilt = GetInt[int](o, "starbasesBuilt")
    obj.TokensBuilt = GetInt[int](o, "tokensBuilt")
    obj.PlanetsColonized = GetInt[int](o, "planetsColonized")
	return obj
}

func GetProductionPlan(o js.Value) cs.ProductionPlan {
	obj := cs.ProductionPlan{}
    obj.Num = GetInt[int](o, "num")
    obj.Name = string(GetString(o, "name"))
    obj.ContributesOnlyLeftoverToResearch = bool(GetBool(o, "contributesOnlyLeftoverToResearch"))
	return obj
}

func GetProductionQueueItem(o js.Value) cs.ProductionQueueItem {
	obj := cs.ProductionQueueItem{}
    obj.DesignNum = GetInt[int](o, "designNum")
    obj.Quantity = GetInt[int](o, "quantity")
    obj.Allocated = GetCost(o.Get("allocated"))
	return obj
}

func GetRace(o js.Value) cs.Race {
	obj := cs.Race{}
    obj.DBObject = GetDBObject(o)
    obj.UserID = GetInt[int64](o, "userId")
    obj.Name = string(GetString(o, "name"))
    obj.PluralName = string(GetString(o, "pluralName"))
    obj.HabLow = GetHab(o.Get("habLow"))
    obj.HabHigh = GetHab(o.Get("habHigh"))
    obj.GrowthRate = GetInt[int](o, "growthRate")
    obj.PopEfficiency = GetInt[int](o, "popEfficiency")
    obj.FactoryOutput = GetInt[int](o, "factoryOutput")
    obj.FactoryCost = GetInt[int](o, "factoryCost")
    obj.NumFactories = GetInt[int](o, "numFactories")
    obj.FactoriesCostLess = bool(GetBool(o, "factoriesCostLess"))
    obj.ImmuneGrav = bool(GetBool(o, "immuneGrav"))
    obj.ImmuneTemp = bool(GetBool(o, "immuneTemp"))
    obj.ImmuneRad = bool(GetBool(o, "immuneRad"))
    obj.MineOutput = GetInt[int](o, "mineOutput")
    obj.MineCost = GetInt[int](o, "mineCost")
    obj.NumMines = GetInt[int](o, "numMines")
    obj.ResearchCost = GetResearchCost(o.Get("researchCost"))
    obj.TechsStartHigh = bool(GetBool(o, "techsStartHigh"))
	return obj
}

func GetResearchCost(o js.Value) cs.ResearchCost {
	obj := cs.ResearchCost{}
	return obj
}

func GetSalvageIntel(o js.Value) cs.SalvageIntel {
	obj := cs.SalvageIntel{}
    obj.MapObjectIntel = GetMapObjectIntel(o)
    obj.Cargo = GetCargo(o.Get("cargo"))
	return obj
}

func GetScoreIntel(o js.Value) cs.ScoreIntel {
	obj := cs.ScoreIntel{}
    obj.ScoreHistory = GetSlice(o.Get("scoreHistory"), GetPlayerScore)
	return obj
}

func GetShipDesign(o js.Value) cs.ShipDesign {
	obj := cs.ShipDesign{}
    obj.GameDBObject = GetGameDBObject(o)
    obj.Num = GetInt[int](o, "num")
    obj.PlayerNum = GetInt[int](o, "playerNum")
    obj.OriginalPlayerNum = GetInt[int](o, "originalPlayerNum")
    obj.Name = string(GetString(o, "name"))
    obj.Version = GetInt[int](o, "version")
    obj.Hull = string(GetString(o, "hull"))
    obj.HullSetNumber = GetInt[int](o, "hullSetNumber")
    obj.CannotDelete = bool(GetBool(o, "cannotDelete"))
    obj.MysteryTrader = bool(GetBool(o, "mysteryTrader"))
    obj.Slots = GetSlice(o.Get("slots"), GetShipDesignSlot)
    obj.Spec = GetShipDesignSpec(o.Get("spec"))
    obj.Delete = bool(GetBool(o, ""))
	return obj
}

func GetShipDesignIntel(o js.Value) cs.ShipDesignIntel {
	obj := cs.ShipDesignIntel{}
    obj.Intel = GetIntel(o)
    obj.Hull = string(GetString(o, "hull"))
    obj.HullSetNumber = GetInt[int](o, "hullSetNumber")
    obj.Version = GetInt[int](o, "version")
    obj.Slots = GetSlice(o.Get("slots"), GetShipDesignSlot)
    obj.Spec = GetShipDesignSpec(o.Get("spec"))
	return obj
}

func GetShipDesignSlot(o js.Value) cs.ShipDesignSlot {
	obj := cs.ShipDesignSlot{}
    obj.HullComponent = string(GetString(o, "hullComponent"))
    obj.HullSlotIndex = GetInt[int](o, "hullSlotIndex")
    obj.Quantity = GetInt[int](o, "quantity")
	return obj
}

func GetShipDesignSpec(o js.Value) cs.ShipDesignSpec {
	obj := cs.ShipDesignSpec{}
    obj.AdditionalMassDrivers = GetInt[int](o, "additionalMassDrivers")
    obj.Armor = GetInt[int](o, "armor")
    obj.BasePacketSpeed = GetInt[int](o, "basePacketSpeed")
    obj.BeamBonus = GetFloat[float64](o, "beamBonus")
    obj.BeamDefense = GetFloat[float64](o, "beamDefense")
    obj.Bomber = bool(GetBool(o, "bomber"))
    obj.Bombs = GetSlice(o.Get("bombs"), GetBomb)
    obj.CanJump = bool(GetBool(o, "canJump"))
    obj.CanLayMines = bool(GetBool(o, "canLayMines"))
    obj.CanStealFleetCargo = bool(GetBool(o, "canStealFleetCargo"))
    obj.CanStealPlanetCargo = bool(GetBool(o, "canStealPlanetCargo"))
    obj.CargoCapacity = GetInt[int](o, "cargoCapacity")
    obj.CloakPercent = GetInt[int](o, "cloakPercent")
    obj.CloakPercentFullCargo = GetInt[int](o, "cloakPercentFullCargo")
    obj.CloakUnits = GetInt[int](o, "cloakUnits")
    obj.Colonizer = bool(GetBool(o, "colonizer"))
    obj.Cost = GetCost(o.Get("cost"))
    obj.Engine = GetEngine(o.Get("engine"))
    obj.EstimatedRange = GetInt[int](o, "estimatedRange")
    obj.EstimatedRangeFull = GetInt[int](o, "estimatedRangeFull")
    obj.FuelCapacity = GetInt[int](o, "fuelCapacity")
    obj.FuelGeneration = GetInt[int](o, "fuelGeneration")
    obj.HasWeapons = bool(GetBool(o, "hasWeapons"))
    obj.ImmuneToOwnDetonation = bool(GetBool(o, "immuneToOwnDetonation"))
    obj.Initiative = GetInt[int](o, "initiative")
    obj.InnateScanRangePenFactor = GetFloat[float64](o, "innateScanRangePenFactor")
    obj.Mass = GetInt[int](o, "mass")
    obj.MassDriver = string(GetString(o, "massDriver"))
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
    obj.OrbitalConstructionModule = bool(GetBool(o, "orbitalConstructionModule"))
    obj.PowerRating = GetInt[int](o, "powerRating")
    obj.Radiating = bool(GetBool(o, "radiating"))
    obj.ReduceCloaking = GetFloat[float64](o, "reduceCloaking")
    obj.ReduceMovement = GetInt[int](o, "reduceMovement")
    obj.RepairBonus = GetFloat[float64](o, "repairBonus")
    obj.RetroBombs = GetSlice(o.Get("retroBombs"), GetBomb)
    obj.SafeHullMass = GetInt[int](o, "safeHullMass")
    obj.SafePacketSpeed = GetInt[int](o, "safePacketSpeed")
    obj.SafeRange = GetInt[int](o, "safeRange")
    obj.Scanner = bool(GetBool(o, "scanner"))
    obj.ScanRange = GetInt[int](o, "scanRange")
    obj.ScanRangePen = GetInt[int](o, "scanRangePen")
    obj.Shields = GetInt[int](o, "shields")
    obj.SmartBombs = GetSlice(o.Get("smartBombs"), GetBomb)
    obj.SpaceDock = GetInt[int](o, "spaceDock")
    obj.Starbase = bool(GetBool(o, "starbase"))
    obj.Stargate = string(GetString(o, "stargate"))
    obj.TechLevel = GetTechLevel(o.Get("techLevel"))
    obj.TerraformRate = GetInt[int](o, "terraformRate")
    obj.TorpedoBonus = GetFloat[float64](o, "torpedoBonus")
    obj.TorpedoJamming = GetFloat[float64](o, "torpedoJamming")
    obj.WeaponSlots = GetSlice(o.Get("weaponSlots"), GetShipDesignSlot)
	return obj
}

func GetShipToken(o js.Value) cs.ShipToken {
	obj := cs.ShipToken{}
    obj.DesignNum = GetInt[int](o, "designNum")
    obj.Quantity = GetInt[int](o, "quantity")
    obj.Damage = GetFloat[float64](o, "damage")
    obj.QuantityDamaged = GetInt[int](o, "quantityDamaged")
	return obj
}

func GetTech(o js.Value) cs.Tech {
	obj := cs.Tech{}
    obj.Name = string(GetString(o, "name"))
    obj.Cost = GetCost(o.Get("cost"))
    obj.Ranking = GetInt[int](o, "ranking")
    obj.Origin = string(GetString(o, "origin"))
	return obj
}

func GetTechDefense(o js.Value) cs.TechDefense {
	obj := cs.TechDefense{}
    obj.TechPlanetary = GetTechPlanetary(o)
    obj.Defense = GetDefense(o)
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

func GetTechPlanetary(o js.Value) cs.TechPlanetary {
	obj := cs.TechPlanetary{}
    obj.Tech = GetTech(o)
    obj.ResetPlanet = bool(GetBool(o, "resetPlanet"))
	return obj
}

func GetTechPlanetaryScanner(o js.Value) cs.TechPlanetaryScanner {
	obj := cs.TechPlanetaryScanner{}
    obj.TechPlanetary = GetTechPlanetary(o)
    obj.ScanRange = GetInt[int](o, "scanRange")
    obj.ScanRangePen = GetInt[int](o, "scanRangePen")
	return obj
}

func GetTransportPlan(o js.Value) cs.TransportPlan {
	obj := cs.TransportPlan{}
    obj.Num = GetInt[int](o, "num")
    obj.Name = string(GetString(o, "name"))
    obj.Tasks = GetWaypointTransportTasks(o.Get("tasks"))
	return obj
}

func GetVector(o js.Value) cs.Vector {
	obj := cs.Vector{}
    obj.X = GetFloat[float64](o, "x")
    obj.Y = GetFloat[float64](o, "y")
	return obj
}

func GetWaypoint(o js.Value) cs.Waypoint {
	obj := cs.Waypoint{}
    obj.Position = GetVector(o.Get("position"))
    obj.WarpSpeed = GetInt[int](o, "warpSpeed")
    obj.EstFuelUsage = GetInt[int](o, "estFuelUsage")
    obj.TransportTasks = GetWaypointTransportTasks(o.Get("transportTasks"))
    obj.WaitAtWaypoint = bool(GetBool(o, "waitAtWaypoint"))
    obj.LayMineFieldDuration = GetInt[int](o, "layMineFieldDuration")
    obj.PatrolRange = GetInt[int](o, "patrolRange")
    obj.PatrolWarpSpeed = GetInt[int](o, "patrolWarpSpeed")
    obj.TargetNum = GetInt[int](o, "targetNum")
    obj.TargetPlayerNum = GetInt[int](o, "targetPlayerNum")
    obj.TargetName = string(GetString(o, "targetName"))
    obj.TransferToPlayer = GetInt[int](o, "transferToPlayer")
    obj.PartiallyComplete = bool(GetBool(o, "partiallyComplete"))
	return obj
}

func GetWaypointTransportTask(o js.Value) cs.WaypointTransportTask {
	obj := cs.WaypointTransportTask{}
    obj.Amount = GetInt[int](o, "amount")
	return obj
}

func GetWaypointTransportTasks(o js.Value) cs.WaypointTransportTasks {
	obj := cs.WaypointTransportTasks{}
    obj.Fuel = GetWaypointTransportTask(o.Get("fuel"))
    obj.Ironium = GetWaypointTransportTask(o.Get("ironium"))
    obj.Boranium = GetWaypointTransportTask(o.Get("boranium"))
    obj.Germanium = GetWaypointTransportTask(o.Get("germanium"))
    obj.Colonists = GetWaypointTransportTask(o.Get("colonists"))
	return obj
}

func GetWormholeIntel(o js.Value) cs.WormholeIntel {
	obj := cs.WormholeIntel{}
    obj.MapObjectIntel = GetMapObjectIntel(o)
    obj.DestinationNum = GetInt[int](o, "destinationNum")
	return obj
}

