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

func SetTime(o js.Value, key string, time time.Time) {
	json, _ := time.MarshalJSON()
	o.Set(key, string(json))
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

// SetSlice sets a jsarray with items using a setter function
func SetSlice[T any](o js.Value, items []T, setter func(o js.Value, item *T)) {

	for i := 0; i < len(items); i++ {
		oItem := js.ValueOf(map[string]any{})
		setter(oItem, &items[i])
		o.SetIndex(i, oItem)
	}
}

// SetPointerSlice sets a jsarray with pointer items using a setter function
func SetPointerSlice[T any](o js.Value, items []*T, setter func(o js.Value, item *T)) {

	for i := 0; i < len(items); i++ {
		oItem := js.ValueOf(map[string]any{})
		setter(oItem, items[i])
		o.SetIndex(i, oItem)
	}
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

// SetSliceSlice sets a 2d jsarray with items using a setter function
func SetSliceSlice[T any](o js.Value, items [][]T, setter func(o js.Value, item *T)) {

	for i := 0; i < len(items); i++ {
		oItem := js.ValueOf([]any{})
		SetSlice(oItem, items[i], setter)
		o.SetIndex(i, oItem)
	}
}

func GetBattlePlan(o js.Value) cs.BattlePlan {
	obj := cs.BattlePlan{}
	obj.Num = GetInt[int](o, "num")
	obj.Name = string(GetString(o, "name"))
	obj.PrimaryTarget = cs.BattleTarget(GetString(o, "primaryTarget"))
	obj.SecondaryTarget = cs.BattleTarget(GetString(o, "secondaryTarget"))
	obj.Tactic = cs.BattleTactic(GetString(o, "tactic"))
	obj.AttackWho = cs.BattleAttackWho(GetString(o, "attackWho"))
	obj.DumpCargo = bool(GetBool(o, "dumpCargo"))
	return obj
}

func SetBattlePlan(o js.Value, obj *cs.BattlePlan) {
	o.Set("num", obj.Num)
	o.Set("name", obj.Name)
	o.Set("primaryTarget", string(obj.PrimaryTarget))
	o.Set("secondaryTarget", string(obj.SecondaryTarget))
	o.Set("tactic", string(obj.Tactic))
	o.Set("attackWho", string(obj.AttackWho))
	o.Set("dumpCargo", obj.DumpCargo)
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

func SetBattleRecord(o js.Value, obj *cs.BattleRecord) {
	o.Set("num", obj.Num)
	o.Set("planetNum", obj.PlanetNum)
	o.Set("position", map[string]any{})
	SetVector(o.Get("position"), &obj.Position)
	o.Set("tokens", []any{})
	SetSlice(o.Get("tokens"), obj.Tokens, SetBattleRecordToken)
	o.Set("actionsPerRound", []any{})
	SetSliceSlice(o.Get("actionsPerRound"), obj.ActionsPerRound, SetBattleRecordTokenAction)
	o.Set("destroyedTokens", []any{})
	SetSlice(o.Get("destroyedTokens"), obj.DestroyedTokens, SetBattleRecordDestroyedToken)
	o.Set("stats", map[string]any{})
	SetBattleRecordStats(o.Get("stats"), &obj.Stats)
}

func GetBattleRecordDestroyedToken(o js.Value) cs.BattleRecordDestroyedToken {
	obj := cs.BattleRecordDestroyedToken{}
	obj.Num = GetInt[int](o, "num")
	obj.PlayerNum = GetInt[int](o, "playerNum")
	obj.DesignNum = GetInt[int](o, "designNum")
	obj.Quantity = GetInt[int](o, "quantity")
	return obj
}

func SetBattleRecordDestroyedToken(o js.Value, obj *cs.BattleRecordDestroyedToken) {
	o.Set("num", obj.Num)
	o.Set("playerNum", obj.PlayerNum)
	o.Set("designNum", obj.DesignNum)
	o.Set("quantity", obj.Quantity)
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
	// CargoLostByPlayer cargoLostByPlayer Map ignored
	return obj
}

func SetBattleRecordStats(o js.Value, obj *cs.BattleRecordStats) {
	o.Set("numPlayers", obj.NumPlayers)
	numShipsByPlayerMap := js.ValueOf(map[string]any{})
	for key, value := range obj.NumShipsByPlayer {
		numShipsByPlayerMap.Set(fmt.Sprintf("%v", key), int(value))
	}
	o.Set("numShipsByPlayer", numShipsByPlayerMap)
	shipsDestroyedByPlayerMap := js.ValueOf(map[string]any{})
	for key, value := range obj.ShipsDestroyedByPlayer {
		shipsDestroyedByPlayerMap.Set(fmt.Sprintf("%v", key), int(value))
	}
	o.Set("shipsDestroyedByPlayer", shipsDestroyedByPlayerMap)
	damageTakenByPlayerMap := js.ValueOf(map[string]any{})
	for key, value := range obj.DamageTakenByPlayer {
		damageTakenByPlayerMap.Set(fmt.Sprintf("%v", key), int(value))
	}
	o.Set("damageTakenByPlayer", damageTakenByPlayerMap)
	// CargoLostByPlayer cargoLostByPlayer Map ignored
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
	obj.Tactic = cs.BattleTactic(GetString(o, "tactic"))
	obj.PrimaryTarget = cs.BattleTarget(GetString(o, "primaryTarget"))
	obj.SecondaryTarget = cs.BattleTarget(GetString(o, "secondaryTarget"))
	obj.AttackWho = cs.BattleAttackWho(GetString(o, "attackWho"))
	return obj
}

func SetBattleRecordToken(o js.Value, obj *cs.BattleRecordToken) {
	o.Set("num", obj.Num)
	o.Set("playerNum", obj.PlayerNum)
	o.Set("designNum", obj.DesignNum)
	o.Set("position", map[string]any{})
	SetBattleVector(o.Get("position"), &obj.Position)
	o.Set("initiative", obj.Initiative)
	o.Set("mass", obj.Mass)
	o.Set("armor", obj.Armor)
	o.Set("stackShields", obj.StackShields)
	o.Set("movement", obj.Movement)
	o.Set("startingQuantity", obj.StartingQuantity)
	o.Set("tactic", string(obj.Tactic))
	o.Set("primaryTarget", string(obj.PrimaryTarget))
	o.Set("secondaryTarget", string(obj.SecondaryTarget))
	o.Set("attackWho", string(obj.AttackWho))
}

func GetBattleRecordTokenAction(o js.Value) cs.BattleRecordTokenAction {
	obj := cs.BattleRecordTokenAction{}
	obj.Type = cs.BattleRecordTokenActionType(GetInt[cs.BattleRecordTokenActionType](o, "type"))
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

func SetBattleRecordTokenAction(o js.Value, obj *cs.BattleRecordTokenAction) {
	o.Set("type", int(obj.Type))
	o.Set("tokenNum", obj.TokenNum)
	o.Set("round", obj.Round)
	o.Set("from", map[string]any{})
	SetBattleVector(o.Get("from"), &obj.From)
	o.Set("to", map[string]any{})
	SetBattleVector(o.Get("to"), &obj.To)
	o.Set("slot", obj.Slot)
	o.Set("targetNum", obj.TargetNum)
	o.Set("target", map[string]any{})
	SetShipToken(o.Get("target"), obj.Target)
	o.Set("tokensDestroyed", obj.TokensDestroyed)
	o.Set("damageDoneShields", obj.DamageDoneShields)
	o.Set("damageDoneArmor", obj.DamageDoneArmor)
	o.Set("torpedoHits", obj.TorpedoHits)
	o.Set("torpedoMisses", obj.TorpedoMisses)
}

func GetBattleVector(o js.Value) cs.BattleVector {
	obj := cs.BattleVector{}
	obj.X = GetInt[int](o, "x")
	obj.Y = GetInt[int](o, "y")
	return obj
}

func SetBattleVector(o js.Value, obj *cs.BattleVector) {
	o.Set("x", obj.X)
	o.Set("y", obj.Y)
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

func SetBomb(o js.Value, obj *cs.Bomb) {
	o.Set("quantity", obj.Quantity)
	o.Set("killRate", obj.KillRate)
	o.Set("minKillRate", obj.MinKillRate)
	o.Set("structureDestroyRate", obj.StructureDestroyRate)
	o.Set("unterraformRate", obj.UnterraformRate)
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

func SetBombingResult(o js.Value, obj *cs.BombingResult) {
	o.Set("bomberName", obj.BomberName)
	o.Set("numBombers", obj.NumBombers)
	o.Set("colonistsKilled", obj.ColonistsKilled)
	o.Set("minesDestroyed", obj.MinesDestroyed)
	o.Set("factoriesDestroyed", obj.FactoriesDestroyed)
	o.Set("defensesDestroyed", obj.DefensesDestroyed)
	o.Set("unterraformAmount", map[string]any{})
	SetHab(o.Get("unterraformAmount"), &obj.UnterraformAmount)
	o.Set("planetEmptied", obj.PlanetEmptied)
}

func GetCargo(o js.Value) cs.Cargo {
	obj := cs.Cargo{}
	obj.Ironium = GetInt[int](o, "ironium")
	obj.Boranium = GetInt[int](o, "boranium")
	obj.Germanium = GetInt[int](o, "germanium")
	obj.Colonists = GetInt[int](o, "colonists")
	return obj
}

func SetCargo(o js.Value, obj *cs.Cargo) {
	o.Set("ironium", obj.Ironium)
	o.Set("boranium", obj.Boranium)
	o.Set("germanium", obj.Germanium)
	o.Set("colonists", obj.Colonists)
}

func GetCost(o js.Value) cs.Cost {
	obj := cs.Cost{}
	obj.Ironium = GetInt[int](o, "ironium")
	obj.Boranium = GetInt[int](o, "boranium")
	obj.Germanium = GetInt[int](o, "germanium")
	obj.Resources = GetInt[int](o, "resources")
	return obj
}

func SetCost(o js.Value, obj *cs.Cost) {
	o.Set("ironium", obj.Ironium)
	o.Set("boranium", obj.Boranium)
	o.Set("germanium", obj.Germanium)
	o.Set("resources", obj.Resources)
}

func GetDBObject(o js.Value) cs.DBObject {
	obj := cs.DBObject{}
	obj.ID = GetInt[int64](o, "id")
	obj.CreatedAt, _ = GetTime(o, "createdAt")
	obj.UpdatedAt, _ = GetTime(o, "updatedAt")
	return obj
}

func SetDBObject(o js.Value, obj *cs.DBObject) {
	o.Set("id", obj.ID)
	SetTime(o, "createdAt", obj.CreatedAt)
	SetTime(o, "updatedAt", obj.UpdatedAt)
}

func GetDefense(o js.Value) cs.Defense {
	obj := cs.Defense{}
	obj.DefenseCoverage = GetFloat[float64](o, "defenseCoverage")
	return obj
}

func SetDefense(o js.Value, obj *cs.Defense) {
	o.Set("defenseCoverage", obj.DefenseCoverage)
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

func SetEngine(o js.Value, obj *cs.Engine) {
	o.Set("idealSpeed", obj.IdealSpeed)
	o.Set("freeSpeed", obj.FreeSpeed)
	o.Set("maxSafeSpeed", obj.MaxSafeSpeed)
	// FuelUsage fuelUsage Array ignored
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

func SetFleet(o js.Value, obj *cs.Fleet) {
	SetMapObject(o, &obj.MapObject)
	SetFleetOrders(o, &obj.FleetOrders)
	o.Set("planetNum", obj.PlanetNum)
	o.Set("baseName", obj.BaseName)
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), &obj.Cargo)
	o.Set("fuel", obj.Fuel)
	o.Set("age", obj.Age)
	o.Set("tokens", []any{})
	SetSlice(o.Get("tokens"), obj.Tokens, SetShipToken)
	o.Set("heading", map[string]any{})
	SetVector(o.Get("heading"), &obj.Heading)
	o.Set("warpSpeed", obj.WarpSpeed)
	o.Set("previousPosition", map[string]any{})
	SetVector(o.Get("previousPosition"), obj.PreviousPosition)
	o.Set("orbitingPlanetNum", obj.OrbitingPlanetNum)
	o.Set("starbase", obj.Starbase)
	o.Set("spec", map[string]any{})
	SetFleetSpec(o.Get("spec"), &obj.Spec)
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

func SetFleetIntel(o js.Value, obj *cs.FleetIntel) {
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	o.Set("", obj.PlanetIntelID)
	o.Set("baseName", obj.BaseName)
	o.Set("heading", map[string]any{})
	SetVector(o.Get("heading"), &obj.Heading)
	o.Set("orbitingPlanetNum", obj.OrbitingPlanetNum)
	o.Set("warpSpeed", obj.WarpSpeed)
	o.Set("mass", obj.Mass)
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), &obj.Cargo)
	o.Set("cargoDiscovered", obj.CargoDiscovered)
	o.Set("freighter", obj.Freighter)
	o.Set("scanRange", obj.ScanRange)
	o.Set("scanRangePen", obj.ScanRangePen)
	o.Set("tokens", []any{})
	SetSlice(o.Get("tokens"), obj.Tokens, SetShipToken)
}

func GetFleetOrders(o js.Value) cs.FleetOrders {
	obj := cs.FleetOrders{}
	obj.Waypoints = GetSlice(o.Get("waypoints"), GetWaypoint)
	obj.RepeatOrders = bool(GetBool(o, "repeatOrders"))
	obj.BattlePlanNum = GetInt[int](o, "battlePlanNum")
	obj.Purpose = cs.FleetPurpose(GetString(o, "purpose"))
	return obj
}

func SetFleetOrders(o js.Value, obj *cs.FleetOrders) {
	o.Set("waypoints", []any{})
	SetSlice(o.Get("waypoints"), obj.Waypoints, SetWaypoint)
	o.Set("repeatOrders", obj.RepeatOrders)
	o.Set("battlePlanNum", obj.BattlePlanNum)
	o.Set("purpose", string(obj.Purpose))
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
	// Purposes purposes Map ignored
	obj.SafeHullMass = GetInt[int](o, "safeHullMass")
	obj.SafeRange = GetInt[int](o, "safeRange")
	obj.Stargate = string(GetString(o, "stargate"))
	obj.TotalShips = GetInt[int](o, "totalShips")
	return obj
}

func SetFleetSpec(o js.Value, obj *cs.FleetSpec) {
	SetShipDesignSpec(o, &obj.ShipDesignSpec)
	o.Set("baseCloakedCargo", obj.BaseCloakedCargo)
	o.Set("basePacketSpeed", obj.BasePacketSpeed)
	o.Set("hasMassDriver", obj.HasMassDriver)
	o.Set("hasStargate", obj.HasStargate)
	o.Set("massDriver", obj.MassDriver)
	o.Set("massEmpty", obj.MassEmpty)
	o.Set("maxHullMass", obj.MaxHullMass)
	o.Set("maxRange", obj.MaxRange)
	// Purposes purposes Map ignored
	o.Set("safeHullMass", obj.SafeHullMass)
	o.Set("safeRange", obj.SafeRange)
	o.Set("stargate", obj.Stargate)
	o.Set("totalShips", obj.TotalShips)
}

func GetGameDBObject(o js.Value) cs.GameDBObject {
	obj := cs.GameDBObject{}
	obj.ID = GetInt[int64](o, "id")
	obj.GameID = GetInt[int64](o, "gameId")
	obj.CreatedAt, _ = GetTime(o, "createdAt")
	obj.UpdatedAt, _ = GetTime(o, "updatedAt")
	return obj
}

func SetGameDBObject(o js.Value, obj *cs.GameDBObject) {
	o.Set("id", obj.ID)
	o.Set("gameId", obj.GameID)
	SetTime(o, "createdAt", obj.CreatedAt)
	SetTime(o, "updatedAt", obj.UpdatedAt)
}

func GetHab(o js.Value) cs.Hab {
	obj := cs.Hab{}
	obj.Grav = GetInt[int](o, "grav")
	obj.Temp = GetInt[int](o, "temp")
	obj.Rad = GetInt[int](o, "rad")
	return obj
}

func SetHab(o js.Value, obj *cs.Hab) {
	o.Set("grav", obj.Grav)
	o.Set("temp", obj.Temp)
	o.Set("rad", obj.Rad)
}

func GetIntel(o js.Value) cs.Intel {
	obj := cs.Intel{}
	obj.Name = string(GetString(o, "name"))
	obj.Num = GetInt[int](o, "num")
	obj.PlayerNum = GetInt[int](o, "playerNum")
	obj.ReportAge = GetInt[int](o, "reportAge")
	return obj
}

func SetIntel(o js.Value, obj *cs.Intel) {
	o.Set("name", obj.Name)
	o.Set("num", obj.Num)
	o.Set("playerNum", obj.PlayerNum)
	o.Set("reportAge", obj.ReportAge)
}

func GetMapObject(o js.Value) cs.MapObject {
	obj := cs.MapObject{}
	obj.GameDBObject = GetGameDBObject(o)
	obj.Type = cs.MapObjectType(GetString(o, "type"))
	// Delete  BasicBool ignored
	obj.Position = GetVector(o.Get("position"))
	obj.Num = GetInt[int](o, "num")
	obj.PlayerNum = GetInt[int](o, "playerNum")
	obj.Name = string(GetString(o, "name"))
	// unknown type Tags cs.Tags map[string]string
	return obj
}

func SetMapObject(o js.Value, obj *cs.MapObject) {
	SetGameDBObject(o, &obj.GameDBObject)
	o.Set("type", string(obj.Type))
	// Delete  BasicBool ignored
	o.Set("position", map[string]any{})
	SetVector(o.Get("position"), &obj.Position)
	o.Set("num", obj.Num)
	o.Set("playerNum", obj.PlayerNum)
	o.Set("name", obj.Name)
	// unknown type Tags cs.Tags map[string]string
}

func GetMapObjectIntel(o js.Value) cs.MapObjectIntel {
	obj := cs.MapObjectIntel{}
	obj.Intel = GetIntel(o)
	obj.Type = cs.MapObjectType(GetString(o, "type"))
	obj.Position = GetVector(o.Get("position"))
	return obj
}

func SetMapObjectIntel(o js.Value, obj *cs.MapObjectIntel) {
	SetIntel(o, &obj.Intel)
	o.Set("type", string(obj.Type))
	o.Set("position", map[string]any{})
	SetVector(o.Get("position"), &obj.Position)
}

func GetMineField(o js.Value) cs.MineField {
	obj := cs.MineField{}
	obj.MapObject = GetMapObject(o)
	obj.MineFieldOrders = GetMineFieldOrders(o)
	obj.MineFieldType = cs.MineFieldType(GetString(o, "mineFieldType"))
	obj.NumMines = GetInt[int](o, "numMines")
	obj.Spec = GetMineFieldSpec(o.Get("spec"))
	return obj
}

func SetMineField(o js.Value, obj *cs.MineField) {
	SetMapObject(o, &obj.MapObject)
	SetMineFieldOrders(o, &obj.MineFieldOrders)
	o.Set("mineFieldType", string(obj.MineFieldType))
	o.Set("numMines", obj.NumMines)
	o.Set("spec", map[string]any{})
	SetMineFieldSpec(o.Get("spec"), &obj.Spec)
}

func GetMineFieldIntel(o js.Value) cs.MineFieldIntel {
	obj := cs.MineFieldIntel{}
	obj.MapObjectIntel = GetMapObjectIntel(o)
	obj.NumMines = GetInt[int](o, "numMines")
	obj.MineFieldType = cs.MineFieldType(GetString(o, "mineFieldType"))
	obj.Spec = GetMineFieldSpec(o.Get("spec"))
	return obj
}

func SetMineFieldIntel(o js.Value, obj *cs.MineFieldIntel) {
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	o.Set("numMines", obj.NumMines)
	o.Set("mineFieldType", string(obj.MineFieldType))
	o.Set("spec", map[string]any{})
	SetMineFieldSpec(o.Get("spec"), &obj.Spec)
}

func GetMineFieldOrders(o js.Value) cs.MineFieldOrders {
	obj := cs.MineFieldOrders{}
	obj.Detonate = bool(GetBool(o, "detonate"))
	return obj
}

func SetMineFieldOrders(o js.Value, obj *cs.MineFieldOrders) {
	o.Set("detonate", obj.Detonate)
}

func GetMineFieldSpec(o js.Value) cs.MineFieldSpec {
	obj := cs.MineFieldSpec{}
	obj.Radius = GetFloat[float64](o, "radius")
	obj.DecayRate = GetInt[int](o, "decayRate")
	return obj
}

func SetMineFieldSpec(o js.Value, obj *cs.MineFieldSpec) {
	o.Set("radius", obj.Radius)
	o.Set("decayRate", obj.DecayRate)
}

func GetMineral(o js.Value) cs.Mineral {
	obj := cs.Mineral{}
	obj.Ironium = GetInt[int](o, "ironium")
	obj.Boranium = GetInt[int](o, "boranium")
	obj.Germanium = GetInt[int](o, "germanium")
	return obj
}

func SetMineral(o js.Value, obj *cs.Mineral) {
	o.Set("ironium", obj.Ironium)
	o.Set("boranium", obj.Boranium)
	o.Set("germanium", obj.Germanium)
}

func GetMineralPacketDamage(o js.Value) cs.MineralPacketDamage {
	obj := cs.MineralPacketDamage{}
	obj.Killed = GetInt[int](o, "killed")
	obj.DefensesDestroyed = GetInt[int](o, "defensesDestroyed")
	obj.Uncaught = GetInt[int](o, "uncaught")
	return obj
}

func SetMineralPacketDamage(o js.Value, obj *cs.MineralPacketDamage) {
	o.Set("killed", obj.Killed)
	o.Set("defensesDestroyed", obj.DefensesDestroyed)
	o.Set("uncaught", obj.Uncaught)
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

func SetMineralPacketIntel(o js.Value, obj *cs.MineralPacketIntel) {
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	o.Set("warpSpeed", obj.WarpSpeed)
	o.Set("heading", map[string]any{})
	SetVector(o.Get("heading"), &obj.Heading)
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), &obj.Cargo)
	o.Set("targetPlanetNum", obj.TargetPlanetNum)
	o.Set("scanRange", obj.ScanRange)
	o.Set("scanRangePen", obj.ScanRangePen)
}

func GetMysteryTrader(o js.Value) cs.MysteryTrader {
	obj := cs.MysteryTrader{}
	obj.MapObject = GetMapObject(o)
	obj.WarpSpeed = GetInt[int](o, "warpSpeed")
	obj.Destination = GetVector(o.Get("destination"))
	obj.RequestedBoon = GetInt[int](o, "requestedBoon")
	obj.RewardType = cs.MysteryTraderRewardType(GetString(o, "rewardType"))
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

func SetMysteryTrader(o js.Value, obj *cs.MysteryTrader) {
	SetMapObject(o, &obj.MapObject)
	o.Set("warpSpeed", obj.WarpSpeed)
	o.Set("destination", map[string]any{})
	SetVector(o.Get("destination"), &obj.Destination)
	o.Set("requestedBoon", obj.RequestedBoon)
	o.Set("rewardType", string(obj.RewardType))
	o.Set("heading", map[string]any{})
	SetVector(o.Get("heading"), &obj.Heading)
	playersRewardedMap := js.ValueOf(map[string]any{})
	for key, value := range obj.PlayersRewarded {
		playersRewardedMap.Set(fmt.Sprintf("%v", key), bool(value))
	}
	o.Set("playersRewarded", playersRewardedMap)
	o.Set("spec", map[string]any{})
	SetMysteryTraderSpec(o.Get("spec"), &obj.Spec)
}

func GetMysteryTraderIntel(o js.Value) cs.MysteryTraderIntel {
	obj := cs.MysteryTraderIntel{}
	obj.MapObjectIntel = GetMapObjectIntel(o)
	obj.WarpSpeed = GetInt[int](o, "warpSpeed")
	obj.Heading = GetVector(o.Get("heading"))
	obj.RequestedBoon = GetInt[int](o, "requestedBoon")
	return obj
}

func SetMysteryTraderIntel(o js.Value, obj *cs.MysteryTraderIntel) {
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	o.Set("warpSpeed", obj.WarpSpeed)
	o.Set("heading", map[string]any{})
	SetVector(o.Get("heading"), &obj.Heading)
	o.Set("requestedBoon", obj.RequestedBoon)
}

func GetMysteryTraderReward(o js.Value) cs.MysteryTraderReward {
	obj := cs.MysteryTraderReward{}
	obj.Type = cs.MysteryTraderRewardType(GetString(o, "type"))
	obj.TechLevels = GetTechLevel(o.Get("techLevels"))
	obj.Tech = string(GetString(o, "tech"))
	obj.Ship = GetShipDesign(o.Get("ship"))
	obj.ShipCount = GetInt[int](o, "shipCount")
	return obj
}

func SetMysteryTraderReward(o js.Value, obj *cs.MysteryTraderReward) {
	o.Set("type", string(obj.Type))
	o.Set("techLevels", map[string]any{})
	SetTechLevel(o.Get("techLevels"), &obj.TechLevels)
	o.Set("tech", obj.Tech)
	o.Set("ship", map[string]any{})
	SetShipDesign(o.Get("ship"), &obj.Ship)
	o.Set("shipCount", obj.ShipCount)
}

func GetMysteryTraderSpec(o js.Value) cs.MysteryTraderSpec {
	obj := cs.MysteryTraderSpec{}
	return obj
}

func SetMysteryTraderSpec(o js.Value, obj *cs.MysteryTraderSpec) {
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
	// RandomArtifact  BasicBool ignored
	// Starbase  Object ignored
	// Dirty  BasicBool ignored
	return obj
}

func SetPlanet(o js.Value, obj *cs.Planet) {
	SetMapObject(o, &obj.MapObject)
	SetPlanetOrders(o, &obj.PlanetOrders)
	o.Set("hab", map[string]any{})
	SetHab(o.Get("hab"), &obj.Hab)
	o.Set("baseHab", map[string]any{})
	SetHab(o.Get("baseHab"), &obj.BaseHab)
	o.Set("terraformedAmount", map[string]any{})
	SetHab(o.Get("terraformedAmount"), &obj.TerraformedAmount)
	o.Set("mineralConcentration", map[string]any{})
	SetMineral(o.Get("mineralConcentration"), &obj.MineralConcentration)
	o.Set("mineYears", map[string]any{})
	SetMineral(o.Get("mineYears"), &obj.MineYears)
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), &obj.Cargo)
	o.Set("mines", obj.Mines)
	o.Set("factories", obj.Factories)
	o.Set("defenses", obj.Defenses)
	o.Set("homeworld", obj.Homeworld)
	o.Set("scanner", obj.Scanner)
	o.Set("spec", map[string]any{})
	SetPlanetSpec(o.Get("spec"), &obj.Spec)
	// RandomArtifact  BasicBool ignored
	// Starbase  Object ignored
	// Dirty  BasicBool ignored
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

func SetPlanetIntel(o js.Value, obj *cs.PlanetIntel) {
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	o.Set("hab", map[string]any{})
	SetHab(o.Get("hab"), &obj.Hab)
	o.Set("baseHab", map[string]any{})
	SetHab(o.Get("baseHab"), &obj.BaseHab)
	o.Set("mineralConcentration", map[string]any{})
	SetMineral(o.Get("mineralConcentration"), &obj.MineralConcentration)
	o.Set("starbase", map[string]any{})
	SetFleetIntel(o.Get("starbase"), obj.Starbase)
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), &obj.Cargo)
	o.Set("cargoDiscovered", obj.CargoDiscovered)
	o.Set("planetHabitability", obj.PlanetHabitability)
	o.Set("planetHabitabilityTerraformed", obj.PlanetHabitabilityTerraformed)
	o.Set("homeworld", obj.Homeworld)
	o.Set("spec", map[string]any{})
	SetPlanetSpec(o.Get("spec"), &obj.Spec)
}

func GetPlanetOrders(o js.Value) cs.PlanetOrders {
	obj := cs.PlanetOrders{}
	obj.ContributesOnlyLeftoverToResearch = bool(GetBool(o, "contributesOnlyLeftoverToResearch"))
	obj.ProductionQueue = GetSlice(o.Get("productionQueue"), GetProductionQueueItem)
	obj.RouteTargetType = cs.MapObjectType(GetString(o, "routeTargetType"))
	obj.RouteTargetNum = GetInt[int](o, "routeTargetNum")
	obj.RouteTargetPlayerNum = GetInt[int](o, "routeTargetPlayerNum")
	obj.PacketTargetNum = GetInt[int](o, "packetTargetNum")
	obj.PacketSpeed = GetInt[int](o, "packetSpeed")
	return obj
}

func SetPlanetOrders(o js.Value, obj *cs.PlanetOrders) {
	o.Set("contributesOnlyLeftoverToResearch", obj.ContributesOnlyLeftoverToResearch)
	o.Set("productionQueue", []any{})
	SetSlice(o.Get("productionQueue"), obj.ProductionQueue, SetProductionQueueItem)
	o.Set("routeTargetType", string(obj.RouteTargetType))
	o.Set("routeTargetNum", obj.RouteTargetNum)
	o.Set("routeTargetPlayerNum", obj.RouteTargetPlayerNum)
	o.Set("packetTargetNum", obj.PacketTargetNum)
	o.Set("packetSpeed", obj.PacketSpeed)
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

func SetPlanetSpec(o js.Value, obj *cs.PlanetSpec) {
	SetPlanetStarbaseSpec(o, &obj.PlanetStarbaseSpec)
	o.Set("canTerraform", obj.CanTerraform)
	o.Set("defense", obj.Defense)
	o.Set("defenseCoverage", obj.DefenseCoverage)
	o.Set("defenseCoverageSmart", obj.DefenseCoverageSmart)
	o.Set("growthAmount", obj.GrowthAmount)
	o.Set("habitability", obj.Habitability)
	o.Set("maxDefenses", obj.MaxDefenses)
	o.Set("maxFactories", obj.MaxFactories)
	o.Set("maxMines", obj.MaxMines)
	o.Set("maxPopulation", obj.MaxPopulation)
	o.Set("maxPossibleFactories", obj.MaxPossibleFactories)
	o.Set("maxPossibleMines", obj.MaxPossibleMines)
	o.Set("miningOutput", map[string]any{})
	SetMineral(o.Get("miningOutput"), &obj.MiningOutput)
	o.Set("population", obj.Population)
	o.Set("populationDensity", obj.PopulationDensity)
	o.Set("resourcesPerYear", obj.ResourcesPerYear)
	o.Set("resourcesPerYearAvailable", obj.ResourcesPerYearAvailable)
	o.Set("resourcesPerYearResearch", obj.ResourcesPerYearResearch)
	o.Set("resourcesPerYearResearchEstimatedLeftover", obj.ResourcesPerYearResearchEstimatedLeftover)
	o.Set("scanner", obj.Scanner)
	o.Set("scanRange", obj.ScanRange)
	o.Set("scanRangePen", obj.ScanRangePen)
	o.Set("terraformAmount", map[string]any{})
	SetHab(o.Get("terraformAmount"), &obj.TerraformAmount)
	o.Set("minTerraformAmount", map[string]any{})
	SetHab(o.Get("minTerraformAmount"), &obj.MinTerraformAmount)
	o.Set("terraformedHabitability", obj.TerraformedHabitability)
	o.Set("contested", obj.Contested)
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

func SetPlanetStarbaseSpec(o js.Value, obj *cs.PlanetStarbaseSpec) {
	o.Set("hasMassDriver", obj.HasMassDriver)
	o.Set("hasStarbase", obj.HasStarbase)
	o.Set("hasStargate", obj.HasStargate)
	o.Set("starbaseDesignName", obj.StarbaseDesignName)
	o.Set("starbaseDesignNum", obj.StarbaseDesignNum)
	o.Set("dockCapacity", obj.DockCapacity)
	o.Set("basePacketSpeed", obj.BasePacketSpeed)
	o.Set("safePacketSpeed", obj.SafePacketSpeed)
	o.Set("safeHullMass", obj.SafeHullMass)
	o.Set("safeRange", obj.SafeRange)
	o.Set("maxRange", obj.MaxRange)
	o.Set("maxHullMass", obj.MaxHullMass)
	o.Set("stargate", obj.Stargate)
	o.Set("massDriver", obj.MassDriver)
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
	obj.AIDifficulty = cs.AIDifficulty(GetString(o, "aiDifficulty"))
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
	obj.AchievedVictoryConditions = cs.Bitmask(GetInt[cs.Bitmask](o, "achievedVictoryConditions"))
	obj.Victor = bool(GetBool(o, "victor"))
	statsVal := o.Get("stats")
	if !statsVal.IsUndefined() {
		stats := GetPlayerStats(statsVal)
		obj.Stats = &stats
	}
	obj.Spec = GetPlayerSpec(o.Get("spec"))
	return obj
}

func SetPlayer(o js.Value, obj *cs.Player) {
	SetGameDBObject(o, &obj.GameDBObject)
	SetPlayerOrders(o, &obj.PlayerOrders)
	SetPlayerIntels(o, &obj.PlayerIntels)
	SetPlayerPlans(o, &obj.PlayerPlans)
	o.Set("userId", obj.UserID)
	o.Set("name", obj.Name)
	o.Set("num", obj.Num)
	o.Set("ready", obj.Ready)
	o.Set("aiControlled", obj.AIControlled)
	o.Set("aiDifficulty", string(obj.AIDifficulty))
	o.Set("guest", obj.Guest)
	o.Set("submittedTurn", obj.SubmittedTurn)
	o.Set("color", obj.Color)
	o.Set("defaultHullSet", obj.DefaultHullSet)
	o.Set("race", map[string]any{})
	SetRace(o.Get("race"), &obj.Race)
	o.Set("techLevels", map[string]any{})
	SetTechLevel(o.Get("techLevels"), &obj.TechLevels)
	o.Set("techLevelsSpent", map[string]any{})
	SetTechLevel(o.Get("techLevelsSpent"), &obj.TechLevelsSpent)
	o.Set("researchSpentLastYear", obj.ResearchSpentLastYear)
	o.Set("relations", []any{})
	SetSlice(o.Get("relations"), obj.Relations, SetPlayerRelationship)
	o.Set("messages", []any{})
	SetSlice(o.Get("messages"), obj.Messages, SetPlayerMessage)
	o.Set("designs", []any{})
	SetPointerSlice(o.Get("designs"), obj.Designs, SetShipDesign)
	o.Set("scoreHistory", []any{})
	SetSlice(o.Get("scoreHistory"), obj.ScoreHistory, SetPlayerScore)
	acquiredTechsMap := js.ValueOf(map[string]any{})
	for key, value := range obj.AcquiredTechs {
		acquiredTechsMap.Set(fmt.Sprintf("%v", key), bool(value))
	}
	o.Set("acquiredTechs", acquiredTechsMap)
	o.Set("achievedVictoryConditions", uint32(obj.AchievedVictoryConditions))
	o.Set("victor", obj.Victor)
	o.Set("stats", map[string]any{})
	SetPlayerStats(o.Get("stats"), obj.Stats)
	o.Set("spec", map[string]any{})
	SetPlayerSpec(o.Get("spec"), &obj.Spec)
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

func SetPlayerIntel(o js.Value, obj *cs.PlayerIntel) {
	o.Set("name", obj.Name)
	o.Set("num", obj.Num)
	o.Set("color", obj.Color)
	o.Set("seen", obj.Seen)
	o.Set("raceName", obj.RaceName)
	o.Set("racePluralName", obj.RacePluralName)
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

func SetPlayerIntels(o js.Value, obj *cs.PlayerIntels) {
	o.Set("battleRecords", []any{})
	SetSlice(o.Get("battleRecords"), obj.BattleRecords, SetBattleRecord)
	o.Set("playerIntels", []any{})
	SetSlice(o.Get("playerIntels"), obj.PlayerIntels, SetPlayerIntel)
	o.Set("scoreIntels", []any{})
	SetSlice(o.Get("scoreIntels"), obj.ScoreIntels, SetScoreIntel)
	o.Set("planetIntels", []any{})
	SetSlice(o.Get("planetIntels"), obj.PlanetIntels, SetPlanetIntel)
	o.Set("fleetIntels", []any{})
	SetSlice(o.Get("fleetIntels"), obj.FleetIntels, SetFleetIntel)
	o.Set("starbaseIntels", []any{})
	SetSlice(o.Get("starbaseIntels"), obj.StarbaseIntels, SetFleetIntel)
	o.Set("shipDesignIntels", []any{})
	SetSlice(o.Get("shipDesignIntels"), obj.ShipDesignIntels, SetShipDesignIntel)
	o.Set("mineralPacketIntels", []any{})
	SetSlice(o.Get("mineralPacketIntels"), obj.MineralPacketIntels, SetMineralPacketIntel)
	o.Set("mineFieldIntels", []any{})
	SetSlice(o.Get("mineFieldIntels"), obj.MineFieldIntels, SetMineFieldIntel)
	o.Set("wormholeIntels", []any{})
	SetSlice(o.Get("wormholeIntels"), obj.WormholeIntels, SetWormholeIntel)
	o.Set("mysteryTraderIntels", []any{})
	SetSlice(o.Get("mysteryTraderIntels"), obj.MysteryTraderIntels, SetMysteryTraderIntel)
	o.Set("salvageIntels", []any{})
	SetSlice(o.Get("salvageIntels"), obj.SalvageIntels, SetSalvageIntel)
}

func GetPlayerMessage(o js.Value) cs.PlayerMessage {
	obj := cs.PlayerMessage{}
	// Target  Object ignored
	obj.Type = cs.PlayerMessageType(GetInt[cs.PlayerMessageType](o, "type"))
	obj.Text = string(GetString(o, "text"))
	obj.BattleNum = GetInt[int](o, "battleNum")
	obj.Spec = GetPlayerMessageSpec(o.Get("spec"))
	return obj
}

func SetPlayerMessage(o js.Value, obj *cs.PlayerMessage) {
	// Target  Object ignored
	o.Set("type", int(obj.Type))
	o.Set("text", obj.Text)
	o.Set("battleNum", obj.BattleNum)
	o.Set("spec", map[string]any{})
	SetPlayerMessageSpec(o.Get("spec"), &obj.Spec)
}

func GetPlayerMessageSpec(o js.Value) cs.PlayerMessageSpec {
	obj := cs.PlayerMessageSpec{}
	// Target  Object ignored
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
	obj.QueueItemType = cs.QueueItemType(GetString(o, "queueItemType"))
	obj.Field = cs.TechField(GetString(o, "field"))
	obj.NextField = cs.TechField(GetString(o, "nextField"))
	obj.TechGained = string(GetString(o, "techGained"))
	obj.LostTargetType = cs.MapObjectType(GetString(o, "lostTargetType"))
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

func SetPlayerMessageSpec(o js.Value, obj *cs.PlayerMessageSpec) {
	// Target  Object ignored
	o.Set("amount", obj.Amount)
	o.Set("amount2", obj.Amount2)
	o.Set("prevAmount", obj.PrevAmount)
	o.Set("sourcePlayerNum", obj.SourcePlayerNum)
	o.Set("destPlayerNum", obj.DestPlayerNum)
	o.Set("name", obj.Name)
	o.Set("cost", map[string]any{})
	SetCost(o.Get("cost"), obj.Cost)
	o.Set("mineral", map[string]any{})
	SetMineral(o.Get("mineral"), obj.Mineral)
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), obj.Cargo)
	o.Set("queueItemType", string(obj.QueueItemType))
	o.Set("field", string(obj.Field))
	o.Set("nextField", string(obj.NextField))
	o.Set("techGained", obj.TechGained)
	o.Set("lostTargetType", string(obj.LostTargetType))
	o.Set("battle", map[string]any{})
	SetBattleRecordStats(o.Get("battle"), &obj.Battle)
	o.Set("comet", map[string]any{})
	SetPlayerMessageSpecComet(o.Get("comet"), obj.Comet)
	o.Set("bombing", map[string]any{})
	SetBombingResult(o.Get("bombing"), obj.Bombing)
	o.Set("mineralPacketDamage", map[string]any{})
	SetMineralPacketDamage(o.Get("mineralPacketDamage"), obj.MineralPacketDamage)
	o.Set("mysteryTrader", map[string]any{})
	SetPlayerMessageSpecMysteryTrader(o.Get("mysteryTrader"), obj.MysteryTrader)
}

func GetPlayerMessageSpecComet(o js.Value) cs.PlayerMessageSpecComet {
	obj := cs.PlayerMessageSpecComet{}
	obj.Size = cs.CometSize(GetString(o, "size"))
	obj.MineralsAdded = GetMineral(o.Get("mineralsAdded"))
	obj.MineralConcentrationIncreased = GetMineral(o.Get("mineralConcentrationIncreased"))
	obj.HabChanged = GetHab(o.Get("habChanged"))
	obj.ColonistsKilled = GetInt[int](o, "colonistsKilled")
	return obj
}

func SetPlayerMessageSpecComet(o js.Value, obj *cs.PlayerMessageSpecComet) {
	o.Set("size", string(obj.Size))
	o.Set("mineralsAdded", map[string]any{})
	SetMineral(o.Get("mineralsAdded"), &obj.MineralsAdded)
	o.Set("mineralConcentrationIncreased", map[string]any{})
	SetMineral(o.Get("mineralConcentrationIncreased"), &obj.MineralConcentrationIncreased)
	o.Set("habChanged", map[string]any{})
	SetHab(o.Get("habChanged"), &obj.HabChanged)
	o.Set("colonistsKilled", obj.ColonistsKilled)
}

func GetPlayerMessageSpecMysteryTrader(o js.Value) cs.PlayerMessageSpecMysteryTrader {
	obj := cs.PlayerMessageSpecMysteryTrader{}
	obj.MysteryTraderReward = GetMysteryTraderReward(o)
	obj.FleetNum = GetInt[int](o, "fleetNum")
	return obj
}

func SetPlayerMessageSpecMysteryTrader(o js.Value, obj *cs.PlayerMessageSpecMysteryTrader) {
	SetMysteryTraderReward(o, &obj.MysteryTraderReward)
	o.Set("fleetNum", obj.FleetNum)
}

func GetPlayerOrders(o js.Value) cs.PlayerOrders {
	obj := cs.PlayerOrders{}
	obj.Researching = cs.TechField(GetString(o, "researching"))
	obj.NextResearchField = cs.NextResearchField(GetString(o, "nextResearchField"))
	obj.ResearchAmount = GetInt[int](o, "researchAmount")
	return obj
}

func SetPlayerOrders(o js.Value, obj *cs.PlayerOrders) {
	o.Set("researching", string(obj.Researching))
	o.Set("nextResearchField", string(obj.NextResearchField))
	o.Set("researchAmount", obj.ResearchAmount)
}

func GetPlayerPlans(o js.Value) cs.PlayerPlans {
	obj := cs.PlayerPlans{}
	obj.ProductionPlans = GetSlice(o.Get("productionPlans"), GetProductionPlan)
	obj.BattlePlans = GetSlice(o.Get("battlePlans"), GetBattlePlan)
	obj.TransportPlans = GetSlice(o.Get("transportPlans"), GetTransportPlan)
	return obj
}

func SetPlayerPlans(o js.Value, obj *cs.PlayerPlans) {
	o.Set("productionPlans", []any{})
	SetSlice(o.Get("productionPlans"), obj.ProductionPlans, SetProductionPlan)
	o.Set("battlePlans", []any{})
	SetSlice(o.Get("battlePlans"), obj.BattlePlans, SetBattlePlan)
	o.Set("transportPlans", []any{})
	SetSlice(o.Get("transportPlans"), obj.TransportPlans, SetTransportPlan)
}

func GetPlayerRelationship(o js.Value) cs.PlayerRelationship {
	obj := cs.PlayerRelationship{}
	obj.Relation = cs.PlayerRelation(GetString(o, "relation"))
	obj.ShareMap = bool(GetBool(o, "shareMap"))
	return obj
}

func SetPlayerRelationship(o js.Value, obj *cs.PlayerRelationship) {
	o.Set("relation", string(obj.Relation))
	o.Set("shareMap", obj.ShareMap)
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
	obj.AchievedVictoryConditions = cs.Bitmask(GetInt[cs.Bitmask](o, "achievedVictoryConditions"))
	return obj
}

func SetPlayerScore(o js.Value, obj *cs.PlayerScore) {
	o.Set("planets", obj.Planets)
	o.Set("starbases", obj.Starbases)
	o.Set("unarmedShips", obj.UnarmedShips)
	o.Set("escortShips", obj.EscortShips)
	o.Set("capitalShips", obj.CapitalShips)
	o.Set("techLevels", obj.TechLevels)
	o.Set("resources", obj.Resources)
	o.Set("score", obj.Score)
	o.Set("rank", obj.Rank)
	o.Set("achievedVictoryConditions", uint32(obj.AchievedVictoryConditions))
}

func GetPlayerSpec(o js.Value) cs.PlayerSpec {
	obj := cs.PlayerSpec{}
	obj.PlanetaryScanner = GetTechPlanetaryScanner(o.Get("planetaryScanner"))
	obj.Defense = GetTechDefense(o.Get("defense"))
	// Terraform terraform Map ignored
	obj.ResourcesPerYear = GetInt[int](o, "resourcesPerYear")
	obj.ResourcesPerYearResearch = GetInt[int](o, "resourcesPerYearResearch")
	obj.ResourcesPerYearResearchEstimated = GetInt[int](o, "resourcesPerYearResearchEstimated")
	obj.CurrentResearchCost = GetInt[int](o, "currentResearchCost")
	return obj
}

func SetPlayerSpec(o js.Value, obj *cs.PlayerSpec) {
	o.Set("planetaryScanner", map[string]any{})
	SetTechPlanetaryScanner(o.Get("planetaryScanner"), &obj.PlanetaryScanner)
	o.Set("defense", map[string]any{})
	SetTechDefense(o.Get("defense"), &obj.Defense)
	// Terraform terraform Map ignored
	o.Set("resourcesPerYear", obj.ResourcesPerYear)
	o.Set("resourcesPerYearResearch", obj.ResourcesPerYearResearch)
	o.Set("resourcesPerYearResearchEstimated", obj.ResourcesPerYearResearchEstimated)
	o.Set("currentResearchCost", obj.CurrentResearchCost)
}

func GetPlayerStats(o js.Value) cs.PlayerStats {
	obj := cs.PlayerStats{}
	obj.FleetsBuilt = GetInt[int](o, "fleetsBuilt")
	obj.StarbasesBuilt = GetInt[int](o, "starbasesBuilt")
	obj.TokensBuilt = GetInt[int](o, "tokensBuilt")
	obj.PlanetsColonized = GetInt[int](o, "planetsColonized")
	return obj
}

func SetPlayerStats(o js.Value, obj *cs.PlayerStats) {
	o.Set("fleetsBuilt", obj.FleetsBuilt)
	o.Set("starbasesBuilt", obj.StarbasesBuilt)
	o.Set("tokensBuilt", obj.TokensBuilt)
	o.Set("planetsColonized", obj.PlanetsColonized)
}

func GetProductionPlan(o js.Value) cs.ProductionPlan {
	obj := cs.ProductionPlan{}
	obj.Num = GetInt[int](o, "num")
	obj.Name = string(GetString(o, "name"))
	obj.Items = GetSlice(o.Get("items"), GetProductionPlanItem)
	obj.ContributesOnlyLeftoverToResearch = bool(GetBool(o, "contributesOnlyLeftoverToResearch"))
	return obj
}

func SetProductionPlan(o js.Value, obj *cs.ProductionPlan) {
	o.Set("num", obj.Num)
	o.Set("name", obj.Name)
	o.Set("items", []any{})
	SetSlice(o.Get("items"), obj.Items, SetProductionPlanItem)
	o.Set("contributesOnlyLeftoverToResearch", obj.ContributesOnlyLeftoverToResearch)
}

func GetProductionPlanItem(o js.Value) cs.ProductionPlanItem {
	obj := cs.ProductionPlanItem{}
	obj.Type = cs.QueueItemType(GetString(o, "type"))
	obj.DesignNum = GetInt[int](o, "designNum")
	obj.Quantity = GetInt[int](o, "quantity")
	return obj
}

func SetProductionPlanItem(o js.Value, obj *cs.ProductionPlanItem) {
	o.Set("type", string(obj.Type))
	o.Set("designNum", obj.DesignNum)
	o.Set("quantity", obj.Quantity)
}

func GetProductionQueueItem(o js.Value) cs.ProductionQueueItem {
	obj := cs.ProductionQueueItem{}
	obj.QueueItemCompletionEstimate = GetQueueItemCompletionEstimate(o)
	obj.Type = cs.QueueItemType(GetString(o, "type"))
	obj.DesignNum = GetInt[int](o, "designNum")
	obj.Quantity = GetInt[int](o, "quantity")
	obj.Allocated = GetCost(o.Get("allocated"))
	// unknown type Tags cs.Tags map[string]string
	return obj
}

func SetProductionQueueItem(o js.Value, obj *cs.ProductionQueueItem) {
	SetQueueItemCompletionEstimate(o, &obj.QueueItemCompletionEstimate)
	o.Set("type", string(obj.Type))
	o.Set("designNum", obj.DesignNum)
	o.Set("quantity", obj.Quantity)
	o.Set("allocated", map[string]any{})
	SetCost(o.Get("allocated"), &obj.Allocated)
	// unknown type Tags cs.Tags map[string]string
}

func GetQueueItemCompletionEstimate(o js.Value) cs.QueueItemCompletionEstimate {
	obj := cs.QueueItemCompletionEstimate{}
	obj.Skipped = bool(GetBool(o, "skipped"))
	obj.YearsToBuildOne = GetInt[int](o, "yearsToBuildOne")
	obj.YearsToBuildAll = GetInt[int](o, "yearsToBuildAll")
	obj.YearsToSkipAuto = GetInt[int](o, "yearsToSkipAuto")
	return obj
}

func SetQueueItemCompletionEstimate(o js.Value, obj *cs.QueueItemCompletionEstimate) {
	o.Set("skipped", obj.Skipped)
	o.Set("yearsToBuildOne", obj.YearsToBuildOne)
	o.Set("yearsToBuildAll", obj.YearsToBuildAll)
	o.Set("yearsToSkipAuto", obj.YearsToSkipAuto)
}

func GetRace(o js.Value) cs.Race {
	obj := cs.Race{}
	obj.DBObject = GetDBObject(o)
	obj.UserID = GetInt[int64](o, "userId")
	obj.Name = string(GetString(o, "name"))
	obj.PluralName = string(GetString(o, "pluralName"))
	obj.SpendLeftoverPointsOn = cs.SpendLeftoverPointsOn(GetString(o, "spendLeftoverPointsOn"))
	obj.PRT = cs.PRT(GetString(o, "prt"))
	obj.LRTs = cs.Bitmask(GetInt[cs.Bitmask](o, "lrts"))
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
	// Spec spec Object ignored
	return obj
}

func SetRace(o js.Value, obj *cs.Race) {
	SetDBObject(o, &obj.DBObject)
	o.Set("userId", obj.UserID)
	o.Set("name", obj.Name)
	o.Set("pluralName", obj.PluralName)
	o.Set("spendLeftoverPointsOn", string(obj.SpendLeftoverPointsOn))
	o.Set("prt", string(obj.PRT))
	o.Set("lrts", uint32(obj.LRTs))
	o.Set("habLow", map[string]any{})
	SetHab(o.Get("habLow"), &obj.HabLow)
	o.Set("habHigh", map[string]any{})
	SetHab(o.Get("habHigh"), &obj.HabHigh)
	o.Set("growthRate", obj.GrowthRate)
	o.Set("popEfficiency", obj.PopEfficiency)
	o.Set("factoryOutput", obj.FactoryOutput)
	o.Set("factoryCost", obj.FactoryCost)
	o.Set("numFactories", obj.NumFactories)
	o.Set("factoriesCostLess", obj.FactoriesCostLess)
	o.Set("immuneGrav", obj.ImmuneGrav)
	o.Set("immuneTemp", obj.ImmuneTemp)
	o.Set("immuneRad", obj.ImmuneRad)
	o.Set("mineOutput", obj.MineOutput)
	o.Set("mineCost", obj.MineCost)
	o.Set("numMines", obj.NumMines)
	o.Set("researchCost", map[string]any{})
	SetResearchCost(o.Get("researchCost"), &obj.ResearchCost)
	o.Set("techsStartHigh", obj.TechsStartHigh)
	// Spec spec Object ignored
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

func SetResearchCost(o js.Value, obj *cs.ResearchCost) {
	o.Set("energy", string(obj.Energy))
	o.Set("weapons", string(obj.Weapons))
	o.Set("propulsion", string(obj.Propulsion))
	o.Set("construction", string(obj.Construction))
	o.Set("electronics", string(obj.Electronics))
	o.Set("biotechnology", string(obj.Biotechnology))
}

func GetSalvageIntel(o js.Value) cs.SalvageIntel {
	obj := cs.SalvageIntel{}
	obj.MapObjectIntel = GetMapObjectIntel(o)
	obj.Cargo = GetCargo(o.Get("cargo"))
	return obj
}

func SetSalvageIntel(o js.Value, obj *cs.SalvageIntel) {
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), &obj.Cargo)
}

func GetScoreIntel(o js.Value) cs.ScoreIntel {
	obj := cs.ScoreIntel{}
	obj.ScoreHistory = GetSlice(o.Get("scoreHistory"), GetPlayerScore)
	return obj
}

func SetScoreIntel(o js.Value, obj *cs.ScoreIntel) {
	o.Set("scoreHistory", []any{})
	SetSlice(o.Get("scoreHistory"), obj.ScoreHistory, SetPlayerScore)
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
	obj.Purpose = cs.ShipDesignPurpose(GetString(o, "purpose"))
	obj.Spec = GetShipDesignSpec(o.Get("spec"))
	obj.Delete = bool(GetBool(o, ""))
	return obj
}

func SetShipDesign(o js.Value, obj *cs.ShipDesign) {
	SetGameDBObject(o, &obj.GameDBObject)
	o.Set("num", obj.Num)
	o.Set("playerNum", obj.PlayerNum)
	o.Set("originalPlayerNum", obj.OriginalPlayerNum)
	o.Set("name", obj.Name)
	o.Set("version", obj.Version)
	o.Set("hull", obj.Hull)
	o.Set("hullSetNumber", obj.HullSetNumber)
	o.Set("cannotDelete", obj.CannotDelete)
	o.Set("mysteryTrader", obj.MysteryTrader)
	o.Set("slots", []any{})
	SetSlice(o.Get("slots"), obj.Slots, SetShipDesignSlot)
	o.Set("purpose", string(obj.Purpose))
	o.Set("spec", map[string]any{})
	SetShipDesignSpec(o.Get("spec"), &obj.Spec)
	o.Set("", obj.Delete)
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

func SetShipDesignIntel(o js.Value, obj *cs.ShipDesignIntel) {
	SetIntel(o, &obj.Intel)
	o.Set("hull", obj.Hull)
	o.Set("hullSetNumber", obj.HullSetNumber)
	o.Set("version", obj.Version)
	o.Set("slots", []any{})
	SetSlice(o.Get("slots"), obj.Slots, SetShipDesignSlot)
	o.Set("spec", map[string]any{})
	SetShipDesignSpec(o.Get("spec"), &obj.Spec)
}

func GetShipDesignSlot(o js.Value) cs.ShipDesignSlot {
	obj := cs.ShipDesignSlot{}
	obj.HullComponent = string(GetString(o, "hullComponent"))
	obj.HullSlotIndex = GetInt[int](o, "hullSlotIndex")
	obj.Quantity = GetInt[int](o, "quantity")
	return obj
}

func SetShipDesignSlot(o js.Value, obj *cs.ShipDesignSlot) {
	o.Set("hullComponent", obj.HullComponent)
	o.Set("hullSlotIndex", obj.HullSlotIndex)
	o.Set("quantity", obj.Quantity)
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
	obj.HullType = cs.TechHullType(GetString(o, "hullType"))
	obj.ImmuneToOwnDetonation = bool(GetBool(o, "immuneToOwnDetonation"))
	obj.Initiative = GetInt[int](o, "initiative")
	obj.InnateScanRangePenFactor = GetFloat[float64](o, "innateScanRangePenFactor")
	obj.Mass = GetInt[int](o, "mass")
	obj.MassDriver = string(GetString(o, "massDriver"))
	obj.MaxHullMass = GetInt[int](o, "maxHullMass")
	obj.MaxPopulation = GetInt[int](o, "maxPopulation")
	obj.MaxRange = GetInt[int](o, "maxRange")
	// MineLayingRateByMineType mineLayingRateByMineType Map ignored
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

func SetShipDesignSpec(o js.Value, obj *cs.ShipDesignSpec) {
	o.Set("additionalMassDrivers", obj.AdditionalMassDrivers)
	o.Set("armor", obj.Armor)
	o.Set("basePacketSpeed", obj.BasePacketSpeed)
	o.Set("beamBonus", obj.BeamBonus)
	o.Set("beamDefense", obj.BeamDefense)
	o.Set("bomber", obj.Bomber)
	o.Set("bombs", []any{})
	SetSlice(o.Get("bombs"), obj.Bombs, SetBomb)
	o.Set("canJump", obj.CanJump)
	o.Set("canLayMines", obj.CanLayMines)
	o.Set("canStealFleetCargo", obj.CanStealFleetCargo)
	o.Set("canStealPlanetCargo", obj.CanStealPlanetCargo)
	o.Set("cargoCapacity", obj.CargoCapacity)
	o.Set("cloakPercent", obj.CloakPercent)
	o.Set("cloakPercentFullCargo", obj.CloakPercentFullCargo)
	o.Set("cloakUnits", obj.CloakUnits)
	o.Set("colonizer", obj.Colonizer)
	o.Set("cost", map[string]any{})
	SetCost(o.Get("cost"), &obj.Cost)
	o.Set("engine", map[string]any{})
	SetEngine(o.Get("engine"), &obj.Engine)
	o.Set("estimatedRange", obj.EstimatedRange)
	o.Set("estimatedRangeFull", obj.EstimatedRangeFull)
	o.Set("fuelCapacity", obj.FuelCapacity)
	o.Set("fuelGeneration", obj.FuelGeneration)
	o.Set("hasWeapons", obj.HasWeapons)
	o.Set("hullType", string(obj.HullType))
	o.Set("immuneToOwnDetonation", obj.ImmuneToOwnDetonation)
	o.Set("initiative", obj.Initiative)
	o.Set("innateScanRangePenFactor", obj.InnateScanRangePenFactor)
	o.Set("mass", obj.Mass)
	o.Set("massDriver", obj.MassDriver)
	o.Set("maxHullMass", obj.MaxHullMass)
	o.Set("maxPopulation", obj.MaxPopulation)
	o.Set("maxRange", obj.MaxRange)
	// MineLayingRateByMineType mineLayingRateByMineType Map ignored
	o.Set("mineSweep", obj.MineSweep)
	o.Set("miningRate", obj.MiningRate)
	o.Set("movement", obj.Movement)
	o.Set("movementBonus", obj.MovementBonus)
	o.Set("movementFull", obj.MovementFull)
	o.Set("numBuilt", obj.NumBuilt)
	o.Set("numEngines", obj.NumEngines)
	o.Set("numInstances", obj.NumInstances)
	o.Set("orbitalConstructionModule", obj.OrbitalConstructionModule)
	o.Set("powerRating", obj.PowerRating)
	o.Set("radiating", obj.Radiating)
	o.Set("reduceCloaking", obj.ReduceCloaking)
	o.Set("reduceMovement", obj.ReduceMovement)
	o.Set("repairBonus", obj.RepairBonus)
	o.Set("retroBombs", []any{})
	SetSlice(o.Get("retroBombs"), obj.RetroBombs, SetBomb)
	o.Set("safeHullMass", obj.SafeHullMass)
	o.Set("safePacketSpeed", obj.SafePacketSpeed)
	o.Set("safeRange", obj.SafeRange)
	o.Set("scanner", obj.Scanner)
	o.Set("scanRange", obj.ScanRange)
	o.Set("scanRangePen", obj.ScanRangePen)
	o.Set("shields", obj.Shields)
	o.Set("smartBombs", []any{})
	SetSlice(o.Get("smartBombs"), obj.SmartBombs, SetBomb)
	o.Set("spaceDock", obj.SpaceDock)
	o.Set("starbase", obj.Starbase)
	o.Set("stargate", obj.Stargate)
	o.Set("techLevel", map[string]any{})
	SetTechLevel(o.Get("techLevel"), &obj.TechLevel)
	o.Set("terraformRate", obj.TerraformRate)
	o.Set("torpedoBonus", obj.TorpedoBonus)
	o.Set("torpedoJamming", obj.TorpedoJamming)
	o.Set("weaponSlots", []any{})
	SetSlice(o.Get("weaponSlots"), obj.WeaponSlots, SetShipDesignSlot)
}

func GetShipToken(o js.Value) cs.ShipToken {
	obj := cs.ShipToken{}
	obj.DesignNum = GetInt[int](o, "designNum")
	obj.Quantity = GetInt[int](o, "quantity")
	obj.Damage = GetFloat[float64](o, "damage")
	obj.QuantityDamaged = GetInt[int](o, "quantityDamaged")
	return obj
}

func SetShipToken(o js.Value, obj *cs.ShipToken) {
	o.Set("designNum", obj.DesignNum)
	o.Set("quantity", obj.Quantity)
	o.Set("damage", obj.Damage)
	o.Set("quantityDamaged", obj.QuantityDamaged)
}

func GetTech(o js.Value) cs.Tech {
	obj := cs.Tech{}
	obj.Name = string(GetString(o, "name"))
	obj.Cost = GetCost(o.Get("cost"))
	// Requirements requirements Object ignored
	obj.Ranking = GetInt[int](o, "ranking")
	obj.Category = cs.TechCategory(GetString(o, "category"))
	obj.Origin = string(GetString(o, "origin"))
	return obj
}

func SetTech(o js.Value, obj *cs.Tech) {
	o.Set("name", obj.Name)
	o.Set("cost", map[string]any{})
	SetCost(o.Get("cost"), &obj.Cost)
	// Requirements requirements Object ignored
	o.Set("ranking", obj.Ranking)
	o.Set("category", string(obj.Category))
	o.Set("origin", obj.Origin)
}

func GetTechDefense(o js.Value) cs.TechDefense {
	obj := cs.TechDefense{}
	obj.TechPlanetary = GetTechPlanetary(o)
	obj.Defense = GetDefense(o)
	return obj
}

func SetTechDefense(o js.Value, obj *cs.TechDefense) {
	SetTechPlanetary(o, &obj.TechPlanetary)
	SetDefense(o, &obj.Defense)
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

func SetTechLevel(o js.Value, obj *cs.TechLevel) {
	o.Set("energy", obj.Energy)
	o.Set("weapons", obj.Weapons)
	o.Set("propulsion", obj.Propulsion)
	o.Set("construction", obj.Construction)
	o.Set("electronics", obj.Electronics)
	o.Set("biotechnology", obj.Biotechnology)
}

func GetTechPlanetary(o js.Value) cs.TechPlanetary {
	obj := cs.TechPlanetary{}
	obj.Tech = GetTech(o)
	obj.ResetPlanet = bool(GetBool(o, "resetPlanet"))
	return obj
}

func SetTechPlanetary(o js.Value, obj *cs.TechPlanetary) {
	SetTech(o, &obj.Tech)
	o.Set("resetPlanet", obj.ResetPlanet)
}

func GetTechPlanetaryScanner(o js.Value) cs.TechPlanetaryScanner {
	obj := cs.TechPlanetaryScanner{}
	obj.TechPlanetary = GetTechPlanetary(o)
	obj.ScanRange = GetInt[int](o, "scanRange")
	obj.ScanRangePen = GetInt[int](o, "scanRangePen")
	return obj
}

func SetTechPlanetaryScanner(o js.Value, obj *cs.TechPlanetaryScanner) {
	SetTechPlanetary(o, &obj.TechPlanetary)
	o.Set("scanRange", obj.ScanRange)
	o.Set("scanRangePen", obj.ScanRangePen)
}

func GetTransportPlan(o js.Value) cs.TransportPlan {
	obj := cs.TransportPlan{}
	obj.Num = GetInt[int](o, "num")
	obj.Name = string(GetString(o, "name"))
	obj.Tasks = GetWaypointTransportTasks(o.Get("tasks"))
	return obj
}

func SetTransportPlan(o js.Value, obj *cs.TransportPlan) {
	o.Set("num", obj.Num)
	o.Set("name", obj.Name)
	o.Set("tasks", map[string]any{})
	SetWaypointTransportTasks(o.Get("tasks"), &obj.Tasks)
}

func GetVector(o js.Value) cs.Vector {
	obj := cs.Vector{}
	obj.X = GetFloat[float64](o, "x")
	obj.Y = GetFloat[float64](o, "y")
	return obj
}

func SetVector(o js.Value, obj *cs.Vector) {
	o.Set("x", obj.X)
	o.Set("y", obj.Y)
}

func GetWaypoint(o js.Value) cs.Waypoint {
	obj := cs.Waypoint{}
	obj.Position = GetVector(o.Get("position"))
	obj.WarpSpeed = GetInt[int](o, "warpSpeed")
	obj.EstFuelUsage = GetInt[int](o, "estFuelUsage")
	obj.Task = cs.WaypointTask(GetString(o, "task"))
	obj.TransportTasks = GetWaypointTransportTasks(o.Get("transportTasks"))
	obj.WaitAtWaypoint = bool(GetBool(o, "waitAtWaypoint"))
	obj.LayMineFieldDuration = GetInt[int](o, "layMineFieldDuration")
	obj.PatrolRange = GetInt[int](o, "patrolRange")
	obj.PatrolWarpSpeed = GetInt[int](o, "patrolWarpSpeed")
	obj.TargetType = cs.MapObjectType(GetString(o, "targetType"))
	obj.TargetNum = GetInt[int](o, "targetNum")
	obj.TargetPlayerNum = GetInt[int](o, "targetPlayerNum")
	obj.TargetName = string(GetString(o, "targetName"))
	obj.TransferToPlayer = GetInt[int](o, "transferToPlayer")
	obj.PartiallyComplete = bool(GetBool(o, "partiallyComplete"))
	return obj
}

func SetWaypoint(o js.Value, obj *cs.Waypoint) {
	o.Set("position", map[string]any{})
	SetVector(o.Get("position"), &obj.Position)
	o.Set("warpSpeed", obj.WarpSpeed)
	o.Set("estFuelUsage", obj.EstFuelUsage)
	o.Set("task", string(obj.Task))
	o.Set("transportTasks", map[string]any{})
	SetWaypointTransportTasks(o.Get("transportTasks"), &obj.TransportTasks)
	o.Set("waitAtWaypoint", obj.WaitAtWaypoint)
	o.Set("layMineFieldDuration", obj.LayMineFieldDuration)
	o.Set("patrolRange", obj.PatrolRange)
	o.Set("patrolWarpSpeed", obj.PatrolWarpSpeed)
	o.Set("targetType", string(obj.TargetType))
	o.Set("targetNum", obj.TargetNum)
	o.Set("targetPlayerNum", obj.TargetPlayerNum)
	o.Set("targetName", obj.TargetName)
	o.Set("transferToPlayer", obj.TransferToPlayer)
	o.Set("partiallyComplete", obj.PartiallyComplete)
}

func GetWaypointTransportTask(o js.Value) cs.WaypointTransportTask {
	obj := cs.WaypointTransportTask{}
	obj.Amount = GetInt[int](o, "amount")
	obj.Action = cs.WaypointTaskTransportAction(GetString(o, "action"))
	return obj
}

func SetWaypointTransportTask(o js.Value, obj *cs.WaypointTransportTask) {
	o.Set("amount", obj.Amount)
	o.Set("action", string(obj.Action))
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

func SetWaypointTransportTasks(o js.Value, obj *cs.WaypointTransportTasks) {
	o.Set("fuel", map[string]any{})
	SetWaypointTransportTask(o.Get("fuel"), &obj.Fuel)
	o.Set("ironium", map[string]any{})
	SetWaypointTransportTask(o.Get("ironium"), &obj.Ironium)
	o.Set("boranium", map[string]any{})
	SetWaypointTransportTask(o.Get("boranium"), &obj.Boranium)
	o.Set("germanium", map[string]any{})
	SetWaypointTransportTask(o.Get("germanium"), &obj.Germanium)
	o.Set("colonists", map[string]any{})
	SetWaypointTransportTask(o.Get("colonists"), &obj.Colonists)
}

func GetWormholeIntel(o js.Value) cs.WormholeIntel {
	obj := cs.WormholeIntel{}
	obj.MapObjectIntel = GetMapObjectIntel(o)
	obj.DestinationNum = GetInt[int](o, "destinationNum")
	obj.Stability = cs.WormholeStability(GetString(o, "stability"))
	return obj
}

func SetWormholeIntel(o js.Value, obj *cs.WormholeIntel) {
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	o.Set("destinationNum", obj.DestinationNum)
	o.Set("stability", string(obj.Stability))
}
