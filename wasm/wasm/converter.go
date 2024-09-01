//go:build wasi || wasm

package wasm

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

func getPointer[T any](val T) *T {
	return &val
}

func getInt[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int | ~int8 | ~int16 | ~int32 | ~int64](o js.Value) T {
	if o.IsUndefined() {
		return 0
	}

	return T(o.Int())
}

func getFloat[T ~float32 | ~float64](o js.Value) T {
	if o.IsUndefined() {
		return 0
	}

	return T(o.Float())
}

func getBool(o js.Value) bool {
	if o.IsUndefined() {
		return false
	}
	return o.Bool()
}

func getString(o js.Value) string {
	if o.IsUndefined() {
		return ""
	}
	return o.String()
}

func getTime(o js.Value) time.Time {
	var result time.Time
	if o.IsUndefined() || o.IsNull() {
		return result
	}
	// time assumes json string has quotes
	result.UnmarshalJSON([]byte("\"" + o.String() + "\""))
	return result
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

// GetSlice populates an array with a getter function
func GetPointerSlice[T any](o js.Value, getter func(o js.Value) T) []*T {
	if o.IsUndefined() || o.IsNull() {
		return nil
	}

	items := make([]*T, o.Length())
	for i := 0; i < len(items); i++ {
		items[i] = getPointer(getter(o.Index(i)))
	}
	return items
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

func GetMap[M ~map[K]V, K comparable, V any](o js.Value, keyGetter func(o js.Value) K, valueGetter func(o js.Value) V) M {
	result := make(M)
	if !o.IsUndefined() {
		resultKeys := js.Global().Get("Object").Call("keys", o)
		for i := 0; i < resultKeys.Length(); i++ {
			key := keyGetter(o)
			result[key] = valueGetter(o.Get(fmt.Sprintf("%v", key)))
		}
	}
	return result
}

func GetAIDifficulty(o js.Value) cs.AIDifficulty {
	var obj cs.AIDifficulty
	if o.IsUndefined() {
		return obj
	}
	obj = cs.AIDifficulty(getString(o))
	return obj
}

func SetAIDifficulty(o js.Value, obj *cs.AIDifficulty) {
}

func GetBattleAttackWho(o js.Value) cs.BattleAttackWho {
	var obj cs.BattleAttackWho
	if o.IsUndefined() {
		return obj
	}
	obj = cs.BattleAttackWho(getString(o))
	return obj
}

func SetBattleAttackWho(o js.Value, obj *cs.BattleAttackWho) {
}

func GetBattlePlan(o js.Value) cs.BattlePlan {
	var obj cs.BattlePlan
	if o.IsUndefined() {
		return obj
	}
	obj.Num = getInt[int](o)
	obj.Name = string(getString(o))
	obj.PrimaryTarget = GetBattleTarget(o.Get("primaryTarget"))
	obj.SecondaryTarget = GetBattleTarget(o.Get("secondaryTarget"))
	obj.Tactic = GetBattleTactic(o.Get("tactic"))
	obj.AttackWho = GetBattleAttackWho(o.Get("attackWho"))
	obj.DumpCargo = getBool(o)
	return obj
}

func SetBattlePlan(o js.Value, obj *cs.BattlePlan) {
}

func GetBattleRecord(o js.Value) cs.BattleRecord {
	var obj cs.BattleRecord
	if o.IsUndefined() {
		return obj
	}
	obj.Num = getInt[int](o)
	obj.PlanetNum = getInt[int](o)
	obj.Position = GetVector(o.Get("position"))
	obj.Tokens = GetSlice(o.Get("tokens"), GetBattleRecordToken)
	obj.ActionsPerRound = GetSliceSlice(o.Get("actionsPerRound"), GetBattleRecordTokenAction)
	obj.DestroyedTokens = GetSlice(o.Get("destroyedTokens"), GetBattleRecordDestroyedToken)
	obj.Stats = GetBattleRecordStats(o.Get("stats"))
	return obj
}

func SetBattleRecord(o js.Value, obj *cs.BattleRecord) {
}

func GetBattleRecordDestroyedToken(o js.Value) cs.BattleRecordDestroyedToken {
	var obj cs.BattleRecordDestroyedToken
	if o.IsUndefined() {
		return obj
	}
	obj.Num = getInt[int](o)
	obj.PlayerNum = getInt[int](o)
	obj.DesignNum = getInt[int](o)
	obj.Quantity = getInt[int](o)
	return obj
}

func SetBattleRecordDestroyedToken(o js.Value, obj *cs.BattleRecordDestroyedToken) {
}

func GetBattleRecordStats(o js.Value) cs.BattleRecordStats {
	var obj cs.BattleRecordStats
	if o.IsUndefined() {
		return obj
	}
	obj.NumPlayers = getInt[int](o)
	obj.NumShipsByPlayer = GetMap[map[int]int, int, int](o.Get("numShipsByPlayer"), getInt, getInt)
	obj.ShipsDestroyedByPlayer = GetMap[map[int]int, int, int](o.Get("shipsDestroyedByPlayer"), getInt, getInt)
	obj.DamageTakenByPlayer = GetMap[map[int]int, int, int](o.Get("damageTakenByPlayer"), getInt, getInt)
	obj.CargoLostByPlayer = GetMap[map[int]cs.Cargo, int, cs.Cargo](o.Get("cargoLostByPlayer"), getInt, GetCargo)
	return obj
}

func SetBattleRecordStats(o js.Value, obj *cs.BattleRecordStats) {
}

func GetBattleRecordToken(o js.Value) cs.BattleRecordToken {
	var obj cs.BattleRecordToken
	if o.IsUndefined() {
		return obj
	}
	obj.Num = getInt[int](o)
	obj.PlayerNum = getInt[int](o)
	obj.DesignNum = getInt[int](o)
	obj.Position = GetBattleVector(o.Get("position"))
	obj.Initiative = getInt[int](o)
	obj.Mass = getInt[int](o)
	obj.Armor = getInt[int](o)
	obj.StackShields = getInt[int](o)
	obj.Movement = getInt[int](o)
	obj.StartingQuantity = getInt[int](o)
	obj.Tactic = GetBattleTactic(o.Get("tactic"))
	obj.PrimaryTarget = GetBattleTarget(o.Get("primaryTarget"))
	obj.SecondaryTarget = GetBattleTarget(o.Get("secondaryTarget"))
	obj.AttackWho = GetBattleAttackWho(o.Get("attackWho"))
	return obj
}

func SetBattleRecordToken(o js.Value, obj *cs.BattleRecordToken) {
}

func GetBattleRecordTokenAction(o js.Value) cs.BattleRecordTokenAction {
	var obj cs.BattleRecordTokenAction
	if o.IsUndefined() {
		return obj
	}
	obj.Type = GetBattleRecordTokenActionType(o.Get("type"))
	obj.TokenNum = getInt[int](o)
	obj.Round = getInt[int](o)
	obj.From = GetBattleVector(o.Get("from"))
	obj.To = GetBattleVector(o.Get("to"))
	obj.Slot = getInt[int](o)
	obj.TargetNum = getInt[int](o)
	obj.Target = getPointer(GetShipToken(o.Get("target")))
	obj.TokensDestroyed = getInt[int](o)
	obj.DamageDoneShields = getInt[int](o)
	obj.DamageDoneArmor = getInt[int](o)
	obj.TorpedoHits = getInt[int](o)
	obj.TorpedoMisses = getInt[int](o)
	return obj
}

func SetBattleRecordTokenAction(o js.Value, obj *cs.BattleRecordTokenAction) {
}

func GetBattleRecordTokenActionType(o js.Value) cs.BattleRecordTokenActionType {
	var obj cs.BattleRecordTokenActionType
	if o.IsUndefined() {
		return obj
	}
	obj = getInt[cs.BattleRecordTokenActionType](o)
	return obj
}

func SetBattleRecordTokenActionType(o js.Value, obj *cs.BattleRecordTokenActionType) {
}

func GetBattleTactic(o js.Value) cs.BattleTactic {
	var obj cs.BattleTactic
	if o.IsUndefined() {
		return obj
	}
	obj = cs.BattleTactic(getString(o))
	return obj
}

func SetBattleTactic(o js.Value, obj *cs.BattleTactic) {
}

func GetBattleTarget(o js.Value) cs.BattleTarget {
	var obj cs.BattleTarget
	if o.IsUndefined() {
		return obj
	}
	obj = cs.BattleTarget(getString(o))
	return obj
}

func SetBattleTarget(o js.Value, obj *cs.BattleTarget) {
}

func GetBattleVector(o js.Value) cs.BattleVector {
	var obj cs.BattleVector
	if o.IsUndefined() {
		return obj
	}
	obj.X = getInt[int](o)
	obj.Y = getInt[int](o)
	return obj
}

func SetBattleVector(o js.Value, obj *cs.BattleVector) {
}

func GetBitmask(o js.Value) cs.Bitmask {
	var obj cs.Bitmask
	if o.IsUndefined() {
		return obj
	}
	obj = getInt[cs.Bitmask](o)
	return obj
}

func SetBitmask(o js.Value, obj *cs.Bitmask) {
}

func GetBomb(o js.Value) cs.Bomb {
	var obj cs.Bomb
	if o.IsUndefined() {
		return obj
	}
	obj.Quantity = getInt[int](o)
	obj.KillRate = getFloat[float64](o)
	obj.MinKillRate = getInt[int](o)
	obj.StructureDestroyRate = getFloat[float64](o)
	obj.UnterraformRate = getInt[int](o)
	return obj
}

func SetBomb(o js.Value, obj *cs.Bomb) {
}

func GetBombingResult(o js.Value) cs.BombingResult {
	var obj cs.BombingResult
	if o.IsUndefined() {
		return obj
	}
	obj.BomberName = string(getString(o))
	obj.NumBombers = getInt[int](o)
	obj.ColonistsKilled = getInt[int](o)
	obj.MinesDestroyed = getInt[int](o)
	obj.FactoriesDestroyed = getInt[int](o)
	obj.DefensesDestroyed = getInt[int](o)
	obj.UnterraformAmount = GetHab(o.Get("unterraformAmount"))
	obj.PlanetEmptied = getBool(o)
	return obj
}

func SetBombingResult(o js.Value, obj *cs.BombingResult) {
}

func GetCargo(o js.Value) cs.Cargo {
	var obj cs.Cargo
	if o.IsUndefined() {
		return obj
	}
	obj.Ironium = getInt[int](o)
	obj.Boranium = getInt[int](o)
	obj.Germanium = getInt[int](o)
	obj.Colonists = getInt[int](o)
	return obj
}

func SetCargo(o js.Value, obj *cs.Cargo) {
}

func GetCargoType(o js.Value) cs.CargoType {
	var obj cs.CargoType
	if o.IsUndefined() {
		return obj
	}
	obj = getInt[cs.CargoType](o)
	return obj
}

func SetCargoType(o js.Value, obj *cs.CargoType) {
}

func GetCometSize(o js.Value) cs.CometSize {
	var obj cs.CometSize
	if o.IsUndefined() {
		return obj
	}
	obj = cs.CometSize(getString(o))
	return obj
}

func SetCometSize(o js.Value, obj *cs.CometSize) {
}

func GetCometStats(o js.Value) cs.CometStats {
	var obj cs.CometStats
	if o.IsUndefined() {
		return obj
	}
	obj.AllMinerals = getInt[int](o)
	obj.AllRandomMinerals = getInt[int](o)
	obj.BonusMinerals = getInt[int](o)
	obj.BonusRandomMinerals = getInt[int](o)
	obj.BonusMinConcentration = getInt[int](o)
	obj.BonusRandomConcentration = getInt[int](o)
	obj.BonusAffectsMinerals = getInt[int](o)
	obj.MinTerraform = getInt[int](o)
	obj.RandomTerraform = getInt[int](o)
	obj.AffectsHabs = getInt[int](o)
	obj.PopKilledPercent = getFloat[float64](o)
	return obj
}

func SetCometStats(o js.Value, obj *cs.CometStats) {
}

func GetCost(o js.Value) cs.Cost {
	var obj cs.Cost
	if o.IsUndefined() {
		return obj
	}
	obj.Ironium = getInt[int](o)
	obj.Boranium = getInt[int](o)
	obj.Germanium = getInt[int](o)
	obj.Resources = getInt[int](o)
	return obj
}

func SetCost(o js.Value, obj *cs.Cost) {
}

func GetDBObject(o js.Value) cs.DBObject {
	var obj cs.DBObject
	if o.IsUndefined() {
		return obj
	}
	obj.ID = getInt[int64](o)
	obj.CreatedAt = getTime(o.Get("createdAt"))
	obj.UpdatedAt = getTime(o.Get("updatedAt"))
	return obj
}

func SetDBObject(o js.Value, obj *cs.DBObject) {
}

func GetDefense(o js.Value) cs.Defense {
	var obj cs.Defense
	if o.IsUndefined() {
		return obj
	}
	obj.DefenseCoverage = getFloat[float64](o)
	return obj
}

func SetDefense(o js.Value, obj *cs.Defense) {
}

func GetDensity(o js.Value) cs.Density {
	var obj cs.Density
	if o.IsUndefined() {
		return obj
	}
	obj = cs.Density(getString(o))
	return obj
}

func SetDensity(o js.Value, obj *cs.Density) {
}

func GetEngine(o js.Value) cs.Engine {
	var obj cs.Engine
	if o.IsUndefined() {
		return obj
	}
	obj.IdealSpeed = getInt[int](o)
	obj.FreeSpeed = getInt[int](o)
	obj.MaxSafeSpeed = getInt[int](o)
	obj.FuelUsage = [11]int(GetSlice[int](o.Get("fuelUsage"), getInt))
	return obj
}

func SetEngine(o js.Value, obj *cs.Engine) {
}

func GetFleet(o js.Value) cs.Fleet {
	var obj cs.Fleet
	if o.IsUndefined() {
		return obj
	}
	obj.PlanetNum = getInt[int](o)
	obj.BaseName = string(getString(o))
	obj.Cargo = GetCargo(o.Get("cargo"))
	obj.Fuel = getInt[int](o)
	obj.Age = getInt[int](o)
	obj.Tokens = GetSlice(o.Get("tokens"), GetShipToken)
	obj.Heading = GetVector(o.Get("heading"))
	obj.WarpSpeed = getInt[int](o)
	obj.PreviousPosition = getPointer(GetVector(o.Get("previousPosition")))
	obj.OrbitingPlanetNum = getInt[int](o)
	obj.Starbase = getBool(o)
	obj.Spec = GetFleetSpec(o.Get("spec"))
	return obj
}

func SetFleet(o js.Value, obj *cs.Fleet) {
	// MapObject  Object ignored
	// FleetOrders  Object ignored
}

func GetFleetIntel(o js.Value) cs.FleetIntel {
	var obj cs.FleetIntel
	if o.IsUndefined() {
		return obj
	}
	obj.BaseName = string(getString(o))
	obj.Heading = GetVector(o.Get("heading"))
	obj.OrbitingPlanetNum = getInt[int](o)
	obj.WarpSpeed = getInt[int](o)
	obj.Mass = getInt[int](o)
	obj.Cargo = GetCargo(o.Get("cargo"))
	obj.CargoDiscovered = getBool(o)
	obj.Freighter = getBool(o)
	obj.ScanRange = getInt[int](o)
	obj.ScanRangePen = getInt[int](o)
	obj.Tokens = GetSlice(o.Get("tokens"), GetShipToken)
	return obj
}

func SetFleetIntel(o js.Value, obj *cs.FleetIntel) {
	// MapObjectIntel  Object ignored
}

func GetFleetOrders(o js.Value) cs.FleetOrders {
	var obj cs.FleetOrders
	if o.IsUndefined() {
		return obj
	}
	obj.Waypoints = GetSlice(o.Get("waypoints"), GetWaypoint)
	obj.RepeatOrders = getBool(o)
	obj.BattlePlanNum = getInt[int](o)
	obj.Purpose = GetFleetPurpose(o.Get("purpose"))
	return obj
}

func SetFleetOrders(o js.Value, obj *cs.FleetOrders) {
}

func GetFleetPurpose(o js.Value) cs.FleetPurpose {
	var obj cs.FleetPurpose
	if o.IsUndefined() {
		return obj
	}
	obj = cs.FleetPurpose(getString(o))
	return obj
}

func SetFleetPurpose(o js.Value, obj *cs.FleetPurpose) {
}

func GetFleetSpec(o js.Value) cs.FleetSpec {
	var obj cs.FleetSpec
	if o.IsUndefined() {
		return obj
	}
	obj.BaseCloakedCargo = getInt[int](o)
	obj.BasePacketSpeed = getInt[int](o)
	obj.HasMassDriver = getBool(o)
	obj.HasStargate = getBool(o)
	obj.MassDriver = string(getString(o))
	obj.MassEmpty = getInt[int](o)
	obj.MaxHullMass = getInt[int](o)
	obj.MaxRange = getInt[int](o)
	obj.Purposes = GetMap[map[cs.ShipDesignPurpose]bool, cs.ShipDesignPurpose, bool](o.Get("purposes"), GetShipDesignPurpose, getBool)
	obj.SafeHullMass = getInt[int](o)
	obj.SafeRange = getInt[int](o)
	obj.Stargate = string(getString(o))
	obj.TotalShips = getInt[int](o)
	return obj
}

func SetFleetSpec(o js.Value, obj *cs.FleetSpec) {
	// ShipDesignSpec  Object ignored
}

func GetGameDBObject(o js.Value) cs.GameDBObject {
	var obj cs.GameDBObject
	if o.IsUndefined() {
		return obj
	}
	obj.ID = getInt[int64](o)
	obj.GameID = getInt[int64](o)
	obj.CreatedAt = getTime(o.Get("createdAt"))
	obj.UpdatedAt = getTime(o.Get("updatedAt"))
	return obj
}

func SetGameDBObject(o js.Value, obj *cs.GameDBObject) {
}

func GetGameStartMode(o js.Value) cs.GameStartMode {
	var obj cs.GameStartMode
	if o.IsUndefined() {
		return obj
	}
	obj = cs.GameStartMode(getString(o))
	return obj
}

func SetGameStartMode(o js.Value, obj *cs.GameStartMode) {
}

func GetGameState(o js.Value) cs.GameState {
	var obj cs.GameState
	if o.IsUndefined() {
		return obj
	}
	obj = cs.GameState(getString(o))
	return obj
}

func SetGameState(o js.Value, obj *cs.GameState) {
}

func GetHab(o js.Value) cs.Hab {
	var obj cs.Hab
	if o.IsUndefined() {
		return obj
	}
	obj.Grav = getInt[int](o)
	obj.Temp = getInt[int](o)
	obj.Rad = getInt[int](o)
	return obj
}

func SetHab(o js.Value, obj *cs.Hab) {
}

func GetHabType(o js.Value) cs.HabType {
	var obj cs.HabType
	if o.IsUndefined() {
		return obj
	}
	obj = getInt[cs.HabType](o)
	return obj
}

func SetHabType(o js.Value, obj *cs.HabType) {
}

func GetHullSlotType(o js.Value) cs.HullSlotType {
	var obj cs.HullSlotType
	if o.IsUndefined() {
		return obj
	}
	obj = getInt[cs.HullSlotType](o)
	return obj
}

func SetHullSlotType(o js.Value, obj *cs.HullSlotType) {
}

func GetIntel(o js.Value) cs.Intel {
	var obj cs.Intel
	if o.IsUndefined() {
		return obj
	}
	obj.Name = string(getString(o))
	obj.Num = getInt[int](o)
	obj.PlayerNum = getInt[int](o)
	obj.ReportAge = getInt[int](o)
	return obj
}

func SetIntel(o js.Value, obj *cs.Intel) {
}

func GetLRT(o js.Value) cs.LRT {
	var obj cs.LRT
	if o.IsUndefined() {
		return obj
	}
	obj = getInt[cs.LRT](o)
	return obj
}

func SetLRT(o js.Value, obj *cs.LRT) {
}

func GetLRTSpec(o js.Value) cs.LRTSpec {
	var obj cs.LRTSpec
	if o.IsUndefined() {
		return obj
	}
	obj.LRT = GetLRT(o.Get("lrt"))
	obj.StartingFleets = GetSlice(o.Get("startingFleets"), GetStartingFleet)
	obj.PointCost = getInt[int](o)
	obj.StartingTechLevels = GetTechLevel(o.Get("startingTechLevels"))
	obj.TechCostOffset = GetTechCostOffset(o.Get("techCostOffset"))
	obj.NewTechCostFactorOffset = getFloat[float64](o)
	obj.MiniaturizationMax = getFloat[float64](o)
	obj.MiniaturizationPerLevel = getFloat[float64](o)
	obj.NoAdvancedScanners = getBool(o)
	obj.ScanRangeFactorOffset = getFloat[float64](o)
	obj.FuelEfficiencyOffset = getFloat[float64](o)
	obj.MaxPopulationOffset = getFloat[float64](o)
	obj.TerraformCostOffset = GetCost(o.Get("terraformCostOffset"))
	obj.MineralAlchemyCostOffset = getInt[int](o)
	obj.ScrapMineralOffset = getFloat[float64](o)
	obj.ScrapMineralOffsetStarbase = getFloat[float64](o)
	obj.ScrapResourcesOffset = getFloat[float64](o)
	obj.ScrapResourcesOffsetStarbase = getFloat[float64](o)
	obj.StartingPopulationFactorDelta = getFloat[float64](o)
	obj.StarbaseBuiltInCloakUnits = getInt[int](o)
	obj.StarbaseCostFactorOffset = getFloat[float64](o)
	obj.ResearchFactorOffset = getFloat[float64](o)
	obj.ResearchSplashDamage = getFloat[float64](o)
	obj.ShieldStrengthFactorOffset = getFloat[float64](o)
	obj.ShieldRegenerationRateOffset = getFloat[float64](o)
	obj.ArmorStrengthFactorOffset = getFloat[float64](o)
	obj.EngineFailureRateOffset = getFloat[float64](o)
	obj.EngineReliableSpeed = getInt[int](o)
	return obj
}

func SetLRTSpec(o js.Value, obj *cs.LRTSpec) {
}

func GetMapObject(o js.Value) cs.MapObject {
	var obj cs.MapObject
	if o.IsUndefined() {
		return obj
	}
	obj.Type = GetMapObjectType(o.Get("type"))
	obj.Position = GetVector(o.Get("position"))
	obj.Num = getInt[int](o)
	obj.PlayerNum = getInt[int](o)
	obj.Name = string(getString(o))
	obj.Tags = GetTags(o.Get("tags"))
	return obj
}

func SetMapObject(o js.Value, obj *cs.MapObject) {
	// GameDBObject  Object ignored
	// Delete  BasicBool ignored
}

func GetMapObjectIntel(o js.Value) cs.MapObjectIntel {
	var obj cs.MapObjectIntel
	if o.IsUndefined() {
		return obj
	}
	obj.Type = GetMapObjectType(o.Get("type"))
	obj.Position = GetVector(o.Get("position"))
	return obj
}

func SetMapObjectIntel(o js.Value, obj *cs.MapObjectIntel) {
	// Intel  Object ignored
}

func GetMapObjectType(o js.Value) cs.MapObjectType {
	var obj cs.MapObjectType
	if o.IsUndefined() {
		return obj
	}
	obj = cs.MapObjectType(getString(o))
	return obj
}

func SetMapObjectType(o js.Value, obj *cs.MapObjectType) {
}

func GetMineField(o js.Value) cs.MineField {
	var obj cs.MineField
	if o.IsUndefined() {
		return obj
	}
	obj.MineFieldType = GetMineFieldType(o.Get("mineFieldType"))
	obj.NumMines = getInt[int](o)
	obj.Spec = GetMineFieldSpec(o.Get("spec"))
	return obj
}

func SetMineField(o js.Value, obj *cs.MineField) {
	// MapObject  Object ignored
	// MineFieldOrders  Object ignored
}

func GetMineFieldIntel(o js.Value) cs.MineFieldIntel {
	var obj cs.MineFieldIntel
	if o.IsUndefined() {
		return obj
	}
	obj.NumMines = getInt[int](o)
	obj.MineFieldType = GetMineFieldType(o.Get("mineFieldType"))
	obj.Spec = GetMineFieldSpec(o.Get("spec"))
	return obj
}

func SetMineFieldIntel(o js.Value, obj *cs.MineFieldIntel) {
	// MapObjectIntel  Object ignored
}

func GetMineFieldOrders(o js.Value) cs.MineFieldOrders {
	var obj cs.MineFieldOrders
	if o.IsUndefined() {
		return obj
	}
	obj.Detonate = getBool(o)
	return obj
}

func SetMineFieldOrders(o js.Value, obj *cs.MineFieldOrders) {
}

func GetMineFieldSpec(o js.Value) cs.MineFieldSpec {
	var obj cs.MineFieldSpec
	if o.IsUndefined() {
		return obj
	}
	obj.Radius = getFloat[float64](o)
	obj.DecayRate = getInt[int](o)
	return obj
}

func SetMineFieldSpec(o js.Value, obj *cs.MineFieldSpec) {
}

func GetMineFieldStats(o js.Value) cs.MineFieldStats {
	var obj cs.MineFieldStats
	if o.IsUndefined() {
		return obj
	}
	obj.MinDamagePerFleetRS = getInt[int](o)
	obj.DamagePerEngineRS = getInt[int](o)
	obj.MaxSpeed = getInt[int](o)
	obj.ChanceOfHit = getFloat[float64](o)
	obj.MinDamagePerFleet = getInt[int](o)
	obj.DamagePerEngine = getInt[int](o)
	obj.SweepFactor = getFloat[float64](o)
	obj.MinDecay = getInt[int](o)
	obj.CanDetonate = getBool(o)
	return obj
}

func SetMineFieldStats(o js.Value, obj *cs.MineFieldStats) {
}

func GetMineFieldType(o js.Value) cs.MineFieldType {
	var obj cs.MineFieldType
	if o.IsUndefined() {
		return obj
	}
	obj = cs.MineFieldType(getString(o))
	return obj
}

func SetMineFieldType(o js.Value, obj *cs.MineFieldType) {
}

func GetMineral(o js.Value) cs.Mineral {
	var obj cs.Mineral
	if o.IsUndefined() {
		return obj
	}
	obj.Ironium = getInt[int](o)
	obj.Boranium = getInt[int](o)
	obj.Germanium = getInt[int](o)
	return obj
}

func SetMineral(o js.Value, obj *cs.Mineral) {
}

func GetMineralPacketDamage(o js.Value) cs.MineralPacketDamage {
	var obj cs.MineralPacketDamage
	if o.IsUndefined() {
		return obj
	}
	obj.Killed = getInt[int](o)
	obj.DefensesDestroyed = getInt[int](o)
	obj.Uncaught = getInt[int](o)
	return obj
}

func SetMineralPacketDamage(o js.Value, obj *cs.MineralPacketDamage) {
}

func GetMineralPacketIntel(o js.Value) cs.MineralPacketIntel {
	var obj cs.MineralPacketIntel
	if o.IsUndefined() {
		return obj
	}
	obj.WarpSpeed = getInt[int](o)
	obj.Heading = GetVector(o.Get("heading"))
	obj.Cargo = GetCargo(o.Get("cargo"))
	obj.TargetPlanetNum = getInt[int](o)
	obj.ScanRange = getInt[int](o)
	obj.ScanRangePen = getInt[int](o)
	return obj
}

func SetMineralPacketIntel(o js.Value, obj *cs.MineralPacketIntel) {
	// MapObjectIntel  Object ignored
}

func GetMineralType(o js.Value) cs.MineralType {
	var obj cs.MineralType
	if o.IsUndefined() {
		return obj
	}
	obj = getInt[cs.MineralType](o)
	return obj
}

func SetMineralType(o js.Value, obj *cs.MineralType) {
}

func GetMiniaturizationSpec(o js.Value) cs.MiniaturizationSpec {
	var obj cs.MiniaturizationSpec
	if o.IsUndefined() {
		return obj
	}
	obj.NewTechCostFactor = getFloat[float64](o)
	obj.MiniaturizationMax = getFloat[float64](o)
	obj.MiniaturizationPerLevel = getFloat[float64](o)
	return obj
}

func SetMiniaturizationSpec(o js.Value, obj *cs.MiniaturizationSpec) {
}

func GetMysteryTrader(o js.Value) cs.MysteryTrader {
	var obj cs.MysteryTrader
	if o.IsUndefined() {
		return obj
	}
	obj.WarpSpeed = getInt[int](o)
	obj.Destination = GetVector(o.Get("destination"))
	obj.RequestedBoon = getInt[int](o)
	obj.RewardType = GetMysteryTraderRewardType(o.Get("rewardType"))
	obj.Heading = GetVector(o.Get("heading"))
	obj.PlayersRewarded = GetMap[map[int]bool, int, bool](o.Get("playersRewarded"), getInt, getBool)
	obj.Spec = GetMysteryTraderSpec(o.Get("spec"))
	return obj
}

func SetMysteryTrader(o js.Value, obj *cs.MysteryTrader) {
	// MapObject  Object ignored
}

func GetMysteryTraderIntel(o js.Value) cs.MysteryTraderIntel {
	var obj cs.MysteryTraderIntel
	if o.IsUndefined() {
		return obj
	}
	obj.WarpSpeed = getInt[int](o)
	obj.Heading = GetVector(o.Get("heading"))
	obj.RequestedBoon = getInt[int](o)
	return obj
}

func SetMysteryTraderIntel(o js.Value, obj *cs.MysteryTraderIntel) {
	// MapObjectIntel  Object ignored
}

func GetMysteryTraderReward(o js.Value) cs.MysteryTraderReward {
	var obj cs.MysteryTraderReward
	if o.IsUndefined() {
		return obj
	}
	obj.Type = GetMysteryTraderRewardType(o.Get("type"))
	obj.TechLevels = GetTechLevel(o.Get("techLevels"))
	obj.Tech = string(getString(o))
	obj.Ship = GetShipDesign(o.Get("ship"))
	obj.ShipCount = getInt[int](o)
	return obj
}

func SetMysteryTraderReward(o js.Value, obj *cs.MysteryTraderReward) {
}

func GetMysteryTraderRewardType(o js.Value) cs.MysteryTraderRewardType {
	var obj cs.MysteryTraderRewardType
	if o.IsUndefined() {
		return obj
	}
	obj = cs.MysteryTraderRewardType(getString(o))
	return obj
}

func SetMysteryTraderRewardType(o js.Value, obj *cs.MysteryTraderRewardType) {
}

func GetMysteryTraderRules(o js.Value) cs.MysteryTraderRules {
	var obj cs.MysteryTraderRules
	if o.IsUndefined() {
		return obj
	}
	obj.ChanceSpawn = GetSlice[int](o.Get("chanceSpawn"), getInt)
	obj.ChanceMaxTechGetsPart = getInt[int](o)
	obj.ChanceCourseChange = getInt[int](o)
	obj.ChanceSpeedUpOnly = getInt[int](o)
	obj.ChanceAgain = getInt[int](o)
	obj.MinYear = getInt[int](o)
	obj.EvenYearOnly = getBool(o)
	obj.MinWarp = getInt[int](o)
	obj.MaxWarp = getInt[int](o)
	obj.MaxMysteryTraders = getInt[int](o)
	obj.RequestedBoon = getInt[int](o)
	obj.GenesisDeviceCost = GetCost(o.Get("genesisDeviceCost"))
	obj.TechBoon = GetSlice(o.Get("techBoon"), GetMysteryTraderTechBoonRules)
	return obj
}

func SetMysteryTraderRules(o js.Value, obj *cs.MysteryTraderRules) {
}

func GetMysteryTraderSpec(o js.Value) cs.MysteryTraderSpec {
	var obj cs.MysteryTraderSpec
	if o.IsUndefined() {
		return obj
	}
	return obj
}

func SetMysteryTraderSpec(o js.Value, obj *cs.MysteryTraderSpec) {
}

func GetMysteryTraderTechBoonMineralsReward(o js.Value) cs.MysteryTraderTechBoonMineralsReward {
	var obj cs.MysteryTraderTechBoonMineralsReward
	if o.IsUndefined() {
		return obj
	}
	obj.MineralsGiven = getInt[int](o)
	obj.Reward = getInt[int](o)
	return obj
}

func SetMysteryTraderTechBoonMineralsReward(o js.Value, obj *cs.MysteryTraderTechBoonMineralsReward) {
}

func GetMysteryTraderTechBoonRules(o js.Value) cs.MysteryTraderTechBoonRules {
	var obj cs.MysteryTraderTechBoonRules
	if o.IsUndefined() {
		return obj
	}
	obj.TechLevels = getInt[int](o)
	obj.Rewards = GetSlice(o.Get("rewards"), GetMysteryTraderTechBoonMineralsReward)
	return obj
}

func SetMysteryTraderTechBoonRules(o js.Value, obj *cs.MysteryTraderTechBoonRules) {
}

func GetNewGamePlayerType(o js.Value) cs.NewGamePlayerType {
	var obj cs.NewGamePlayerType
	if o.IsUndefined() {
		return obj
	}
	obj = cs.NewGamePlayerType(getString(o))
	return obj
}

func SetNewGamePlayerType(o js.Value, obj *cs.NewGamePlayerType) {
}

func GetNextResearchField(o js.Value) cs.NextResearchField {
	var obj cs.NextResearchField
	if o.IsUndefined() {
		return obj
	}
	obj = cs.NextResearchField(getString(o))
	return obj
}

func SetNextResearchField(o js.Value, obj *cs.NextResearchField) {
}

func GetPRT(o js.Value) cs.PRT {
	var obj cs.PRT
	if o.IsUndefined() {
		return obj
	}
	obj = cs.PRT(getString(o))
	return obj
}

func SetPRT(o js.Value, obj *cs.PRT) {
}

func GetPRTSpec(o js.Value) cs.PRTSpec {
	var obj cs.PRTSpec
	if o.IsUndefined() {
		return obj
	}
	obj.PRT = GetPRT(o.Get("prt"))
	obj.PointCost = getInt[int](o)
	obj.StartingTechLevels = GetTechLevel(o.Get("startingTechLevels"))
	obj.StartingPlanets = GetSlice(o.Get("startingPlanets"), GetStartingPlanet)
	obj.TechCostOffset = GetTechCostOffset(o.Get("techCostOffset"))
	obj.MineralsPerSingleMineralPacket = getInt[int](o)
	obj.MineralsPerMixedMineralPacket = getInt[int](o)
	obj.PacketResourceCost = getInt[int](o)
	obj.PacketMineralCostFactor = getFloat[float64](o)
	obj.PacketReceiverFactor = getFloat[float64](o)
	obj.PacketDecayFactor = getFloat[float64](o)
	obj.PacketOverSafeWarpPenalty = getInt[int](o)
	obj.PacketBuiltInScanner = getBool(o)
	obj.DetectPacketDestinationStarbases = getBool(o)
	obj.DetectAllPackets = getBool(o)
	obj.PacketTerraformChance = getFloat[float64](o)
	obj.PacketPermaformChance = getFloat[float64](o)
	obj.PacketPermaTerraformSizeUnit = getInt[int](o)
	obj.CanGateCargo = getBool(o)
	obj.CanDetectStargatePlanets = getBool(o)
	obj.ShipsVanishInVoid = getBool(o)
	obj.BuiltInScannerMultiplier = getInt[int](o)
	obj.TechsCostExtraLevel = getInt[int](o)
	obj.FreighterGrowthFactor = getFloat[float64](o)
	obj.GrowthFactor = getFloat[float64](o)
	obj.MaxPopulationOffset = getFloat[float64](o)
	obj.BuiltInCloakUnits = getInt[int](o)
	obj.StealsResearch = GetStealsResearch(o.Get("stealsResearch"))
	obj.FreeCargoCloaking = getBool(o)
	obj.MineFieldsAreScanners = getBool(o)
	obj.MineFieldRateMoveFactor = getFloat[float64](o)
	obj.MineFieldSafeWarpBonus = getInt[int](o)
	obj.MineFieldMinDecayFactor = getFloat[float64](o)
	obj.MineFieldBaseDecayRate = getFloat[float64](o)
	obj.MineFieldPlanetDecayRate = getFloat[float64](o)
	obj.MineFieldMaxDecayRate = getFloat[float64](o)
	obj.CanDetonateMineFields = getBool(o)
	obj.MineFieldDetonateDecayRate = getFloat[float64](o)
	obj.DiscoverDesignOnScan = getBool(o)
	obj.CanRemoteMineOwnPlanets = getBool(o)
	obj.InvasionAttackBonus = getFloat[float64](o)
	obj.InvasionDefendBonus = getFloat[float64](o)
	obj.MovementBonus = getInt[int](o)
	obj.Instaforming = getBool(o)
	obj.PermaformChance = getFloat[float64](o)
	obj.PermaformPopulation = getInt[int](o)
	obj.RepairFactor = getFloat[float64](o)
	obj.StarbaseRepairFactor = getFloat[float64](o)
	obj.StarbaseCostFactor = getFloat[float64](o)
	obj.InnateMining = getBool(o)
	obj.InnateResources = getBool(o)
	obj.InnateScanner = getBool(o)
	obj.InnatePopulationFactor = getFloat[float64](o)
	obj.CanBuildDefenses = getBool(o)
	obj.LivesOnStarbases = getBool(o)
	return obj
}

func SetPRTSpec(o js.Value, obj *cs.PRTSpec) {
}

func GetPlanet(o js.Value) cs.Planet {
	var obj cs.Planet
	if o.IsUndefined() {
		return obj
	}
	obj.Hab = GetHab(o.Get("hab"))
	obj.BaseHab = GetHab(o.Get("baseHab"))
	obj.TerraformedAmount = GetHab(o.Get("terraformedAmount"))
	obj.MineralConcentration = GetMineral(o.Get("mineralConcentration"))
	obj.MineYears = GetMineral(o.Get("mineYears"))
	obj.Cargo = GetCargo(o.Get("cargo"))
	obj.Mines = getInt[int](o)
	obj.Factories = getInt[int](o)
	obj.Defenses = getInt[int](o)
	obj.Homeworld = getBool(o)
	obj.Scanner = getBool(o)
	obj.Spec = GetPlanetSpec(o.Get("spec"))
	return obj
}

func SetPlanet(o js.Value, obj *cs.Planet) {
	// MapObject  Object ignored
	// PlanetOrders  Object ignored
	// RandomArtifact  BasicBool ignored
	// Starbase  Object ignored
	// Dirty  BasicBool ignored
}

func GetPlanetIntel(o js.Value) cs.PlanetIntel {
	var obj cs.PlanetIntel
	if o.IsUndefined() {
		return obj
	}
	obj.Hab = GetHab(o.Get("hab"))
	obj.BaseHab = GetHab(o.Get("baseHab"))
	obj.MineralConcentration = GetMineral(o.Get("mineralConcentration"))
	obj.Starbase = getPointer(GetFleetIntel(o.Get("starbase")))
	obj.Cargo = GetCargo(o.Get("cargo"))
	obj.CargoDiscovered = getBool(o)
	obj.PlanetHabitability = getInt[int](o)
	obj.PlanetHabitabilityTerraformed = getInt[int](o)
	obj.Homeworld = getBool(o)
	obj.Spec = GetPlanetSpec(o.Get("spec"))
	return obj
}

func SetPlanetIntel(o js.Value, obj *cs.PlanetIntel) {
	// MapObjectIntel  Object ignored
}

func GetPlanetOrders(o js.Value) cs.PlanetOrders {
	var obj cs.PlanetOrders
	if o.IsUndefined() {
		return obj
	}
	obj.ContributesOnlyLeftoverToResearch = getBool(o)
	obj.ProductionQueue = GetSlice(o.Get("productionQueue"), GetProductionQueueItem)
	obj.RouteTargetType = GetMapObjectType(o.Get("routeTargetType"))
	obj.RouteTargetNum = getInt[int](o)
	obj.RouteTargetPlayerNum = getInt[int](o)
	obj.PacketTargetNum = getInt[int](o)
	obj.PacketSpeed = getInt[int](o)
	return obj
}

func SetPlanetOrders(o js.Value, obj *cs.PlanetOrders) {
}

func GetPlanetSpec(o js.Value) cs.PlanetSpec {
	var obj cs.PlanetSpec
	if o.IsUndefined() {
		return obj
	}
	obj.CanTerraform = getBool(o)
	obj.Defense = string(getString(o))
	obj.DefenseCoverage = getFloat[float64](o)
	obj.DefenseCoverageSmart = getFloat[float64](o)
	obj.GrowthAmount = getInt[int](o)
	obj.Habitability = getInt[int](o)
	obj.MaxDefenses = getInt[int](o)
	obj.MaxFactories = getInt[int](o)
	obj.MaxMines = getInt[int](o)
	obj.MaxPopulation = getInt[int](o)
	obj.MaxPossibleFactories = getInt[int](o)
	obj.MaxPossibleMines = getInt[int](o)
	obj.MiningOutput = GetMineral(o.Get("miningOutput"))
	obj.Population = getInt[int](o)
	obj.PopulationDensity = getFloat[float64](o)
	obj.ResourcesPerYear = getInt[int](o)
	obj.ResourcesPerYearAvailable = getInt[int](o)
	obj.ResourcesPerYearResearch = getInt[int](o)
	obj.ResourcesPerYearResearchEstimatedLeftover = getInt[int](o)
	obj.Scanner = string(getString(o))
	obj.ScanRange = getInt[int](o)
	obj.ScanRangePen = getInt[int](o)
	obj.TerraformAmount = GetHab(o.Get("terraformAmount"))
	obj.MinTerraformAmount = GetHab(o.Get("minTerraformAmount"))
	obj.TerraformedHabitability = getInt[int](o)
	obj.Contested = getBool(o)
	return obj
}

func SetPlanetSpec(o js.Value, obj *cs.PlanetSpec) {
	// PlanetStarbaseSpec  Object ignored
}

func GetPlanetStarbaseSpec(o js.Value) cs.PlanetStarbaseSpec {
	var obj cs.PlanetStarbaseSpec
	if o.IsUndefined() {
		return obj
	}
	obj.HasMassDriver = getBool(o)
	obj.HasStarbase = getBool(o)
	obj.HasStargate = getBool(o)
	obj.StarbaseDesignName = string(getString(o))
	obj.StarbaseDesignNum = getInt[int](o)
	obj.DockCapacity = getInt[int](o)
	obj.BasePacketSpeed = getInt[int](o)
	obj.SafePacketSpeed = getInt[int](o)
	obj.SafeHullMass = getInt[int](o)
	obj.SafeRange = getInt[int](o)
	obj.MaxRange = getInt[int](o)
	obj.MaxHullMass = getInt[int](o)
	obj.Stargate = string(getString(o))
	obj.MassDriver = string(getString(o))
	return obj
}

func SetPlanetStarbaseSpec(o js.Value, obj *cs.PlanetStarbaseSpec) {
}

func GetPlayer(o js.Value) cs.Player {
	var obj cs.Player
	if o.IsUndefined() {
		return obj
	}
	obj.UserID = getInt[int64](o)
	obj.Name = string(getString(o))
	obj.Num = getInt[int](o)
	obj.Ready = getBool(o)
	obj.AIControlled = getBool(o)
	obj.AIDifficulty = GetAIDifficulty(o.Get("aiDifficulty"))
	obj.Guest = getBool(o)
	obj.SubmittedTurn = getBool(o)
	obj.Color = string(getString(o))
	obj.DefaultHullSet = getInt[int](o)
	obj.Race = GetRace(o.Get("race"))
	obj.TechLevels = GetTechLevel(o.Get("techLevels"))
	obj.TechLevelsSpent = GetTechLevel(o.Get("techLevelsSpent"))
	obj.ResearchSpentLastYear = getInt[int](o)
	obj.Relations = GetSlice(o.Get("relations"), GetPlayerRelationship)
	obj.Messages = GetSlice(o.Get("messages"), GetPlayerMessage)
	obj.Designs = GetPointerSlice(o.Get("designs"), GetShipDesign)
	obj.ScoreHistory = GetSlice(o.Get("scoreHistory"), GetPlayerScore)
	obj.AcquiredTechs = GetMap[map[string]bool, string, bool](o.Get("acquiredTechs"), getString, getBool)
	obj.AchievedVictoryConditions = GetBitmask(o.Get("achievedVictoryConditions"))
	obj.Victor = getBool(o)
	obj.Stats = getPointer(GetPlayerStats(o.Get("stats")))
	obj.Spec = GetPlayerSpec(o.Get("spec"))
	return obj
}

func SetPlayer(o js.Value, obj *cs.Player) {
	// GameDBObject  Object ignored
	// PlayerOrders  Object ignored
	// PlayerIntels  Object ignored
	// PlayerPlans  Object ignored
}

func GetPlayerIntel(o js.Value) cs.PlayerIntel {
	var obj cs.PlayerIntel
	if o.IsUndefined() {
		return obj
	}
	obj.Name = string(getString(o))
	obj.Num = getInt[int](o)
	obj.Color = string(getString(o))
	obj.Seen = getBool(o)
	obj.RaceName = string(getString(o))
	obj.RacePluralName = string(getString(o))
	return obj
}

func SetPlayerIntel(o js.Value, obj *cs.PlayerIntel) {
}

func GetPlayerIntels(o js.Value) cs.PlayerIntels {
	var obj cs.PlayerIntels
	if o.IsUndefined() {
		return obj
	}
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
}

func GetPlayerMessage(o js.Value) cs.PlayerMessage {
	var obj cs.PlayerMessage
	if o.IsUndefined() {
		return obj
	}
	obj.Type = GetPlayerMessageType(o.Get("type"))
	obj.Text = string(getString(o))
	obj.BattleNum = getInt[int](o)
	obj.Spec = GetPlayerMessageSpec(o.Get("spec"))
	return obj
}

func SetPlayerMessage(o js.Value, obj *cs.PlayerMessage) {
	// Target  Object ignored
}

func GetPlayerMessageSpec(o js.Value) cs.PlayerMessageSpec {
	var obj cs.PlayerMessageSpec
	if o.IsUndefined() {
		return obj
	}
	obj.Amount = getInt[int](o)
	obj.Amount2 = getInt[int](o)
	obj.PrevAmount = getInt[int](o)
	obj.SourcePlayerNum = getInt[int](o)
	obj.DestPlayerNum = getInt[int](o)
	obj.Name = string(getString(o))
	obj.Cost = getPointer(GetCost(o.Get("cost")))
	obj.Mineral = getPointer(GetMineral(o.Get("mineral")))
	obj.Cargo = getPointer(GetCargo(o.Get("cargo")))
	obj.QueueItemType = GetQueueItemType(o.Get("queueItemType"))
	obj.Field = GetTechField(o.Get("field"))
	obj.NextField = GetTechField(o.Get("nextField"))
	obj.TechGained = string(getString(o))
	obj.LostTargetType = GetMapObjectType(o.Get("lostTargetType"))
	obj.Battle = GetBattleRecordStats(o.Get("battle"))
	obj.Comet = getPointer(GetPlayerMessageSpecComet(o.Get("comet")))
	obj.Bombing = getPointer(GetBombingResult(o.Get("bombing")))
	obj.MineralPacketDamage = getPointer(GetMineralPacketDamage(o.Get("mineralPacketDamage")))
	obj.MysteryTrader = getPointer(GetPlayerMessageSpecMysteryTrader(o.Get("mysteryTrader")))
	return obj
}

func SetPlayerMessageSpec(o js.Value, obj *cs.PlayerMessageSpec) {
	// Target  Object ignored
}

func GetPlayerMessageSpecComet(o js.Value) cs.PlayerMessageSpecComet {
	var obj cs.PlayerMessageSpecComet
	if o.IsUndefined() {
		return obj
	}
	obj.Size = GetCometSize(o.Get("size"))
	obj.MineralsAdded = GetMineral(o.Get("mineralsAdded"))
	obj.MineralConcentrationIncreased = GetMineral(o.Get("mineralConcentrationIncreased"))
	obj.HabChanged = GetHab(o.Get("habChanged"))
	obj.ColonistsKilled = getInt[int](o)
	return obj
}

func SetPlayerMessageSpecComet(o js.Value, obj *cs.PlayerMessageSpecComet) {
}

func GetPlayerMessageSpecMysteryTrader(o js.Value) cs.PlayerMessageSpecMysteryTrader {
	var obj cs.PlayerMessageSpecMysteryTrader
	if o.IsUndefined() {
		return obj
	}
	obj.FleetNum = getInt[int](o)
	return obj
}

func SetPlayerMessageSpecMysteryTrader(o js.Value, obj *cs.PlayerMessageSpecMysteryTrader) {
	// MysteryTraderReward  Object ignored
}

func GetPlayerMessageTargetType(o js.Value) cs.PlayerMessageTargetType {
	var obj cs.PlayerMessageTargetType
	if o.IsUndefined() {
		return obj
	}
	obj = cs.PlayerMessageTargetType(getString(o))
	return obj
}

func SetPlayerMessageTargetType(o js.Value, obj *cs.PlayerMessageTargetType) {
}

func GetPlayerMessageType(o js.Value) cs.PlayerMessageType {
	var obj cs.PlayerMessageType
	if o.IsUndefined() {
		return obj
	}
	obj = getInt[cs.PlayerMessageType](o)
	return obj
}

func SetPlayerMessageType(o js.Value, obj *cs.PlayerMessageType) {
}

func GetPlayerOrders(o js.Value) cs.PlayerOrders {
	var obj cs.PlayerOrders
	if o.IsUndefined() {
		return obj
	}
	obj.Researching = GetTechField(o.Get("researching"))
	obj.NextResearchField = GetNextResearchField(o.Get("nextResearchField"))
	obj.ResearchAmount = getInt[int](o)
	return obj
}

func SetPlayerOrders(o js.Value, obj *cs.PlayerOrders) {
}

func GetPlayerPlans(o js.Value) cs.PlayerPlans {
	var obj cs.PlayerPlans
	if o.IsUndefined() {
		return obj
	}
	obj.ProductionPlans = GetSlice(o.Get("productionPlans"), GetProductionPlan)
	obj.BattlePlans = GetSlice(o.Get("battlePlans"), GetBattlePlan)
	obj.TransportPlans = GetSlice(o.Get("transportPlans"), GetTransportPlan)
	return obj
}

func SetPlayerPlans(o js.Value, obj *cs.PlayerPlans) {
}

func GetPlayerPositions(o js.Value) cs.PlayerPositions {
	var obj cs.PlayerPositions
	if o.IsUndefined() {
		return obj
	}
	obj = cs.PlayerPositions(getString(o))
	return obj
}

func SetPlayerPositions(o js.Value, obj *cs.PlayerPositions) {
}

func GetPlayerRelation(o js.Value) cs.PlayerRelation {
	var obj cs.PlayerRelation
	if o.IsUndefined() {
		return obj
	}
	obj = cs.PlayerRelation(getString(o))
	return obj
}

func SetPlayerRelation(o js.Value, obj *cs.PlayerRelation) {
}

func GetPlayerRelationship(o js.Value) cs.PlayerRelationship {
	var obj cs.PlayerRelationship
	if o.IsUndefined() {
		return obj
	}
	obj.Relation = GetPlayerRelation(o.Get("relation"))
	obj.ShareMap = getBool(o)
	return obj
}

func SetPlayerRelationship(o js.Value, obj *cs.PlayerRelationship) {
}

func GetPlayerScore(o js.Value) cs.PlayerScore {
	var obj cs.PlayerScore
	if o.IsUndefined() {
		return obj
	}
	obj.Planets = getInt[int](o)
	obj.Starbases = getInt[int](o)
	obj.UnarmedShips = getInt[int](o)
	obj.EscortShips = getInt[int](o)
	obj.CapitalShips = getInt[int](o)
	obj.TechLevels = getInt[int](o)
	obj.Resources = getInt[int](o)
	obj.Score = getInt[int](o)
	obj.Rank = getInt[int](o)
	obj.AchievedVictoryConditions = GetBitmask(o.Get("achievedVictoryConditions"))
	return obj
}

func SetPlayerScore(o js.Value, obj *cs.PlayerScore) {
}

func GetPlayerSpec(o js.Value) cs.PlayerSpec {
	var obj cs.PlayerSpec
	if o.IsUndefined() {
		return obj
	}
	obj.PlanetaryScanner = GetTechPlanetaryScanner(o.Get("planetaryScanner"))
	obj.Defense = GetTechDefense(o.Get("defense"))
	obj.Terraform = GetMap[map[cs.TerraformHabType]*cs.TechTerraform, cs.TerraformHabType, *cs.TechTerraform](o.Get("terraform"), GetTerraformHabType, func(o js.Value) *cs.TechTerraform { return getPointer(GetTechTerraform(o)) })
	obj.ResourcesPerYear = getInt[int](o)
	obj.ResourcesPerYearResearch = getInt[int](o)
	obj.ResourcesPerYearResearchEstimated = getInt[int](o)
	obj.CurrentResearchCost = getInt[int](o)
	return obj
}

func SetPlayerSpec(o js.Value, obj *cs.PlayerSpec) {
}

func GetPlayerStats(o js.Value) cs.PlayerStats {
	var obj cs.PlayerStats
	if o.IsUndefined() {
		return obj
	}
	obj.FleetsBuilt = getInt[int](o)
	obj.StarbasesBuilt = getInt[int](o)
	obj.TokensBuilt = getInt[int](o)
	obj.PlanetsColonized = getInt[int](o)
	return obj
}

func SetPlayerStats(o js.Value, obj *cs.PlayerStats) {
}

func GetProductionPlan(o js.Value) cs.ProductionPlan {
	var obj cs.ProductionPlan
	if o.IsUndefined() {
		return obj
	}
	obj.Num = getInt[int](o)
	obj.Name = string(getString(o))
	obj.Items = GetSlice(o.Get("items"), GetProductionPlanItem)
	obj.ContributesOnlyLeftoverToResearch = getBool(o)
	return obj
}

func SetProductionPlan(o js.Value, obj *cs.ProductionPlan) {
}

func GetProductionPlanItem(o js.Value) cs.ProductionPlanItem {
	var obj cs.ProductionPlanItem
	if o.IsUndefined() {
		return obj
	}
	obj.Type = GetQueueItemType(o.Get("type"))
	obj.DesignNum = getInt[int](o)
	obj.Quantity = getInt[int](o)
	return obj
}

func SetProductionPlanItem(o js.Value, obj *cs.ProductionPlanItem) {
}

func GetProductionQueueItem(o js.Value) cs.ProductionQueueItem {
	var obj cs.ProductionQueueItem
	if o.IsUndefined() {
		return obj
	}
	obj.Type = GetQueueItemType(o.Get("type"))
	obj.DesignNum = getInt[int](o)
	obj.Quantity = getInt[int](o)
	obj.Allocated = GetCost(o.Get("allocated"))
	obj.Tags = GetTags(o.Get("tags"))
	return obj
}

func SetProductionQueueItem(o js.Value, obj *cs.ProductionQueueItem) {
	// QueueItemCompletionEstimate  Object ignored
}

func GetQueueItemCompletionEstimate(o js.Value) cs.QueueItemCompletionEstimate {
	var obj cs.QueueItemCompletionEstimate
	if o.IsUndefined() {
		return obj
	}
	obj.Skipped = getBool(o)
	obj.YearsToBuildOne = getInt[int](o)
	obj.YearsToBuildAll = getInt[int](o)
	obj.YearsToSkipAuto = getInt[int](o)
	return obj
}

func SetQueueItemCompletionEstimate(o js.Value, obj *cs.QueueItemCompletionEstimate) {
}

func GetQueueItemType(o js.Value) cs.QueueItemType {
	var obj cs.QueueItemType
	if o.IsUndefined() {
		return obj
	}
	obj = cs.QueueItemType(getString(o))
	return obj
}

func SetQueueItemType(o js.Value, obj *cs.QueueItemType) {
}

func GetRace(o js.Value) cs.Race {
	var obj cs.Race
	if o.IsUndefined() {
		return obj
	}
	obj.UserID = getInt[int64](o)
	obj.Name = string(getString(o))
	obj.PluralName = string(getString(o))
	obj.SpendLeftoverPointsOn = GetSpendLeftoverPointsOn(o.Get("spendLeftoverPointsOn"))
	obj.PRT = GetPRT(o.Get("prt"))
	obj.LRTs = GetBitmask(o.Get("lrts"))
	obj.HabLow = GetHab(o.Get("habLow"))
	obj.HabHigh = GetHab(o.Get("habHigh"))
	obj.GrowthRate = getInt[int](o)
	obj.PopEfficiency = getInt[int](o)
	obj.FactoryOutput = getInt[int](o)
	obj.FactoryCost = getInt[int](o)
	obj.NumFactories = getInt[int](o)
	obj.FactoriesCostLess = getBool(o)
	obj.ImmuneGrav = getBool(o)
	obj.ImmuneTemp = getBool(o)
	obj.ImmuneRad = getBool(o)
	obj.MineOutput = getInt[int](o)
	obj.MineCost = getInt[int](o)
	obj.NumMines = getInt[int](o)
	obj.ResearchCost = GetResearchCost(o.Get("researchCost"))
	obj.TechsStartHigh = getBool(o)
	obj.Spec = GetRaceSpec(o.Get("spec"))
	return obj
}

func SetRace(o js.Value, obj *cs.Race) {
	// DBObject  Object ignored
}

func GetRaceSpec(o js.Value) cs.RaceSpec {
	var obj cs.RaceSpec
	if o.IsUndefined() {
		return obj
	}
	obj.HabCenter = GetHab(o.Get("habCenter"))
	obj.Costs = GetMap[map[cs.QueueItemType]cs.Cost, cs.QueueItemType, cs.Cost](o.Get("costs"), GetQueueItemType, GetCost)
	obj.StartingTechLevels = GetTechLevel(o.Get("startingTechLevels"))
	obj.StartingPlanets = GetSlice(o.Get("startingPlanets"), GetStartingPlanet)
	obj.TechCostOffset = GetTechCostOffset(o.Get("techCostOffset"))
	obj.MineralsPerSingleMineralPacket = getInt[int](o)
	obj.MineralsPerMixedMineralPacket = getInt[int](o)
	obj.PacketResourceCost = getInt[int](o)
	obj.PacketMineralCostFactor = getFloat[float64](o)
	obj.PacketReceiverFactor = getFloat[float64](o)
	obj.PacketDecayFactor = getFloat[float64](o)
	obj.PacketOverSafeWarpPenalty = getInt[int](o)
	obj.PacketBuiltInScanner = getBool(o)
	obj.DetectPacketDestinationStarbases = getBool(o)
	obj.DetectAllPackets = getBool(o)
	obj.PacketTerraformChance = getFloat[float64](o)
	obj.PacketPermaformChance = getFloat[float64](o)
	obj.PacketPermaTerraformSizeUnit = getInt[int](o)
	obj.CanGateCargo = getBool(o)
	obj.CanDetectStargatePlanets = getBool(o)
	obj.ShipsVanishInVoid = getBool(o)
	obj.TechsCostExtraLevel = getInt[int](o)
	obj.FreighterGrowthFactor = getFloat[float64](o)
	obj.GrowthFactor = getFloat[float64](o)
	obj.MaxPopulationOffset = getFloat[float64](o)
	obj.BuiltInCloakUnits = getInt[int](o)
	obj.StealsResearch = GetStealsResearch(o.Get("stealsResearch"))
	obj.FreeCargoCloaking = getBool(o)
	obj.MineFieldsAreScanners = getBool(o)
	obj.MineFieldRateMoveFactor = getFloat[float64](o)
	obj.MineFieldSafeWarpBonus = getInt[int](o)
	obj.MineFieldMinDecayFactor = getFloat[float64](o)
	obj.MineFieldBaseDecayRate = getFloat[float64](o)
	obj.MineFieldPlanetDecayRate = getFloat[float64](o)
	obj.MineFieldMaxDecayRate = getFloat[float64](o)
	obj.CanDetonateMineFields = getBool(o)
	obj.MineFieldDetonateDecayRate = getFloat[float64](o)
	obj.DiscoverDesignOnScan = getBool(o)
	obj.CanRemoteMineOwnPlanets = getBool(o)
	obj.InvasionAttackBonus = getFloat[float64](o)
	obj.InvasionDefendBonus = getFloat[float64](o)
	obj.MovementBonus = getInt[int](o)
	obj.Instaforming = getBool(o)
	obj.PermaformChance = getFloat[float64](o)
	obj.PermaformPopulation = getInt[int](o)
	obj.RepairFactor = getFloat[float64](o)
	obj.StarbaseRepairFactor = getFloat[float64](o)
	obj.InnateMining = getBool(o)
	obj.InnateResources = getBool(o)
	obj.InnateScanner = getBool(o)
	obj.InnatePopulationFactor = getFloat[float64](o)
	obj.CanBuildDefenses = getBool(o)
	obj.LivesOnStarbases = getBool(o)
	obj.FuelEfficiencyOffset = getFloat[float64](o)
	obj.TerraformCostOffset = GetCost(o.Get("terraformCostOffset"))
	obj.MineralAlchemyCostOffset = getInt[int](o)
	obj.ScrapMineralOffset = getFloat[float64](o)
	obj.ScrapMineralOffsetStarbase = getFloat[float64](o)
	obj.ScrapResourcesOffset = getFloat[float64](o)
	obj.ScrapResourcesOffsetStarbase = getFloat[float64](o)
	obj.StartingPopulationFactor = getFloat[float64](o)
	obj.StarbaseBuiltInCloakUnits = getInt[int](o)
	obj.StarbaseCostFactor = getFloat[float64](o)
	obj.ResearchFactor = getFloat[float64](o)
	obj.ResearchSplashDamage = getFloat[float64](o)
	obj.ArmorStrengthFactor = getFloat[float64](o)
	obj.ShieldStrengthFactor = getFloat[float64](o)
	obj.ShieldRegenerationRate = getFloat[float64](o)
	obj.EngineFailureRate = getFloat[float64](o)
	obj.EngineReliableSpeed = getInt[int](o)
	return obj
}

func SetRaceSpec(o js.Value, obj *cs.RaceSpec) {
	// MiniaturizationSpec  Object ignored
	// ScannerSpec  Object ignored
}

func GetRandomCometSize(o js.Value) cs.RandomCometSize {
	var obj cs.RandomCometSize
	if o.IsUndefined() {
		return obj
	}
	obj = getInt[cs.RandomCometSize](o)
	return obj
}

func SetRandomCometSize(o js.Value, obj *cs.RandomCometSize) {
}

func GetRandomEvent(o js.Value) cs.RandomEvent {
	var obj cs.RandomEvent
	if o.IsUndefined() {
		return obj
	}
	obj = cs.RandomEvent(getString(o))
	return obj
}

func SetRandomEvent(o js.Value, obj *cs.RandomEvent) {
}

func GetRandomEventType(o js.Value) cs.RandomEventType {
	var obj cs.RandomEventType
	if o.IsUndefined() {
		return obj
	}
	obj = getInt[cs.RandomEventType](o)
	return obj
}

func SetRandomEventType(o js.Value, obj *cs.RandomEventType) {
}

func GetRepairRate(o js.Value) cs.RepairRate {
	var obj cs.RepairRate
	if o.IsUndefined() {
		return obj
	}
	obj = cs.RepairRate(getString(o))
	return obj
}

func SetRepairRate(o js.Value, obj *cs.RepairRate) {
}

func GetResearchCost(o js.Value) cs.ResearchCost {
	var obj cs.ResearchCost
	if o.IsUndefined() {
		return obj
	}
	obj.Energy = GetResearchCostLevel(o.Get("energy"))
	obj.Weapons = GetResearchCostLevel(o.Get("weapons"))
	obj.Propulsion = GetResearchCostLevel(o.Get("propulsion"))
	obj.Construction = GetResearchCostLevel(o.Get("construction"))
	obj.Electronics = GetResearchCostLevel(o.Get("electronics"))
	obj.Biotechnology = GetResearchCostLevel(o.Get("biotechnology"))
	return obj
}

func SetResearchCost(o js.Value, obj *cs.ResearchCost) {
}

func GetResearchCostLevel(o js.Value) cs.ResearchCostLevel {
	var obj cs.ResearchCostLevel
	if o.IsUndefined() {
		return obj
	}
	obj = cs.ResearchCostLevel(getString(o))
	return obj
}

func SetResearchCostLevel(o js.Value, obj *cs.ResearchCostLevel) {
}

func GetRules(o js.Value) cs.Rules {
	var obj cs.Rules
	if o.IsUndefined() {
		return obj
	}
	obj.ID = getInt[int64](o)
	obj.CreatedAt = getTime(o.Get("createdAt"))
	obj.UpdatedAt = getTime(o.Get("updatedAt"))
	obj.GameID = getInt[int64](o)
	obj.TachyonCloakReduction = getInt[int](o)
	obj.MaxPopulation = getInt[int](o)
	obj.MinMaxPopulationPercent = getFloat[float64](o)
	obj.PopulationOvercrowdDieoffRate = getFloat[float64](o)
	obj.PopulationOvercrowdDieoffRateMax = getFloat[float64](o)
	obj.PopulationScannerError = getFloat[float64](o)
	obj.SmartDefenseCoverageFactor = getFloat[float64](o)
	obj.InvasionDefenseCoverageFactor = getFloat[float64](o)
	obj.NumBattleRounds = getInt[int](o)
	obj.MovesToRunAway = getInt[int](o)
	obj.BeamRangeDropoff = getFloat[float64](o)
	obj.TorpedoSplashDamage = getFloat[float64](o)
	obj.SalvageDecayRate = getFloat[float64](o)
	obj.SalvageDecayMin = getInt[int](o)
	obj.MineFieldCloak = getInt[int](o)
	obj.StargateMaxRangeFactor = getInt[int](o)
	obj.StargateMaxHullMassFactor = getInt[int](o)
	obj.FleetSafeSpeedExplosionChance = getFloat[float64](o)
	obj.RandomEventChances = GetMap[map[cs.RandomEvent]float64, cs.RandomEvent, float64](o.Get("randomEventChances"), GetRandomEvent, getFloat)
	obj.RandomMineralDepositBonusRange = [2]int(GetSlice[int](o.Get("randomMineralDepositBonusRange"), getInt))
	obj.RandomArtifactResearchBonusRange = [2]int(GetSlice[int](o.Get("randomArtifactResearchBonusRange"), getInt))
	obj.RandomCometMinYear = getInt[int](o)
	obj.RandomCometMinYearPlayerWorld = getInt[int](o)
	obj.MysteryTraderRules = GetMysteryTraderRules(o.Get("mysteryTraderRules"))
	obj.CometStatsBySize = GetMap[map[cs.CometSize]cs.CometStats, cs.CometSize, cs.CometStats](o.Get("cometStatsBySize"), GetCometSize, GetCometStats)
	obj.WormholeCloak = getInt[int](o)
	obj.WormholeMinPlanetDistance = getInt[int](o)
	obj.WormholeStatsByStability = GetMap[map[cs.WormholeStability]cs.WormholeStats, cs.WormholeStability, cs.WormholeStats](o.Get("wormholeStatsByStability"), GetWormholeStability, GetWormholeStats)
	obj.WormholePairsForSize = GetMap[map[cs.Size]int, cs.Size, int](o.Get("wormholePairsForSize"), GetSize, getInt)
	obj.MineFieldStatsByType = GetMap[map[cs.MineFieldType]cs.MineFieldStats, cs.MineFieldType, cs.MineFieldStats](o.Get("mineFieldStatsByType"), GetMineFieldType, GetMineFieldStats)
	obj.RepairRates = GetMap[map[cs.RepairRate]float64, cs.RepairRate, float64](o.Get("repairRates"), GetRepairRate, getFloat)
	obj.MaxPlayers = getInt[int](o)
	obj.StartingYear = getInt[int](o)
	obj.ShowPublicScoresAfterYears = getInt[int](o)
	obj.PlanetMinDistance = getInt[int](o)
	obj.MaxExtraWorldDistance = getInt[int](o)
	obj.MinExtraWorldDistance = getInt[int](o)
	obj.MinHomeworldMineralConcentration = getInt[int](o)
	obj.MinExtraPlanetMineralConcentration = getInt[int](o)
	obj.MinHab = getInt[int](o)
	obj.MaxHab = getInt[int](o)
	obj.MinMineralConcentration = getInt[int](o)
	obj.MaxMineralConcentration = getInt[int](o)
	obj.MinStartingMineralConcentration = getInt[int](o)
	obj.MaxStartingMineralConcentration = getInt[int](o)
	obj.HighRadMineralConcentrationBonusThreshold = getInt[int](o)
	obj.RadiatingImmune = getInt[int](o)
	obj.MaxStartingMineralSurface = getInt[int](o)
	obj.MinStartingMineralSurface = getInt[int](o)
	obj.MineralDecayFactor = getInt[int](o)
	obj.RemoteMiningMineOutput = getInt[int](o)
	obj.StartingMines = getInt[int](o)
	obj.StartingFactories = getInt[int](o)
	obj.StartingDefenses = getInt[int](o)
	obj.RaceStartingPoints = getInt[int](o)
	obj.ScrapMineralAmount = getFloat[float64](o)
	obj.ScrapResourceAmount = getFloat[float64](o)
	obj.FactoryCostGermanium = getInt[int](o)
	obj.DefenseCost = GetCost(o.Get("defenseCost"))
	obj.MineralAlchemyCost = getInt[int](o)
	obj.PlanetaryScannerCost = GetCost(o.Get("planetaryScannerCost"))
	obj.TerraformCost = GetCost(o.Get("terraformCost"))
	obj.StarbaseComponentCostFactor = getFloat[float64](o)
	obj.SalvageFromBattleFactor = getFloat[float64](o)
	obj.TechTradeChance = getFloat[float64](o)
	obj.PacketDecayRate = GetMap[map[int]float64, int, float64](o.Get("packetDecayRate"), getInt, getFloat)
	obj.PacketMinDecay = getInt[int](o)
	obj.MaxTechLevel = getInt[int](o)
	obj.TechBaseCost = GetSlice[int](o.Get("techBaseCost"), getInt)
	obj.PRTSpecs = GetMap[map[cs.PRT]cs.PRTSpec, cs.PRT, cs.PRTSpec](o.Get("prtSpecs"), GetPRT, GetPRTSpec)
	obj.LRTSpecs = GetMap[map[cs.LRT]cs.LRTSpec, cs.LRT, cs.LRTSpec](o.Get("lrtSpecs"), GetLRT, GetLRTSpec)
	obj.TechsID = getInt[int64](o)
	return obj
}

func SetRules(o js.Value, obj *cs.Rules) {
}

func GetSalvageIntel(o js.Value) cs.SalvageIntel {
	var obj cs.SalvageIntel
	if o.IsUndefined() {
		return obj
	}
	obj.Cargo = GetCargo(o.Get("cargo"))
	return obj
}

func SetSalvageIntel(o js.Value, obj *cs.SalvageIntel) {
	// MapObjectIntel  Object ignored
}

func GetScannerSpec(o js.Value) cs.ScannerSpec {
	var obj cs.ScannerSpec
	if o.IsUndefined() {
		return obj
	}
	obj.BuiltInScannerMultiplier = getInt[int](o)
	obj.NoAdvancedScanners = getBool(o)
	obj.ScanRangeFactor = getFloat[float64](o)
	return obj
}

func SetScannerSpec(o js.Value, obj *cs.ScannerSpec) {
}

func GetScoreIntel(o js.Value) cs.ScoreIntel {
	var obj cs.ScoreIntel
	if o.IsUndefined() {
		return obj
	}
	obj.ScoreHistory = GetSlice(o.Get("scoreHistory"), GetPlayerScore)
	return obj
}

func SetScoreIntel(o js.Value, obj *cs.ScoreIntel) {
}

func GetShipDesign(o js.Value) cs.ShipDesign {
	var obj cs.ShipDesign
	if o.IsUndefined() {
		return obj
	}
	obj.Num = getInt[int](o)
	obj.PlayerNum = getInt[int](o)
	obj.OriginalPlayerNum = getInt[int](o)
	obj.Name = string(getString(o))
	obj.Version = getInt[int](o)
	obj.Hull = string(getString(o))
	obj.HullSetNumber = getInt[int](o)
	obj.CannotDelete = getBool(o)
	obj.MysteryTrader = getBool(o)
	obj.Slots = GetSlice(o.Get("slots"), GetShipDesignSlot)
	obj.Purpose = GetShipDesignPurpose(o.Get("purpose"))
	obj.Spec = GetShipDesignSpec(o.Get("spec"))
	return obj
}

func SetShipDesign(o js.Value, obj *cs.ShipDesign) {
	// GameDBObject  Object ignored
	// Delete  BasicBool ignored
}

func GetShipDesignIntel(o js.Value) cs.ShipDesignIntel {
	var obj cs.ShipDesignIntel
	if o.IsUndefined() {
		return obj
	}
	obj.Hull = string(getString(o))
	obj.HullSetNumber = getInt[int](o)
	obj.Version = getInt[int](o)
	obj.Slots = GetSlice(o.Get("slots"), GetShipDesignSlot)
	obj.Spec = GetShipDesignSpec(o.Get("spec"))
	return obj
}

func SetShipDesignIntel(o js.Value, obj *cs.ShipDesignIntel) {
	// Intel  Object ignored
}

func GetShipDesignPurpose(o js.Value) cs.ShipDesignPurpose {
	var obj cs.ShipDesignPurpose
	if o.IsUndefined() {
		return obj
	}
	obj = cs.ShipDesignPurpose(getString(o))
	return obj
}

func SetShipDesignPurpose(o js.Value, obj *cs.ShipDesignPurpose) {
}

func GetShipDesignSlot(o js.Value) cs.ShipDesignSlot {
	var obj cs.ShipDesignSlot
	if o.IsUndefined() {
		return obj
	}
	obj.HullComponent = string(getString(o))
	obj.HullSlotIndex = getInt[int](o)
	obj.Quantity = getInt[int](o)
	return obj
}

func SetShipDesignSlot(o js.Value, obj *cs.ShipDesignSlot) {
}

func GetShipDesignSpec(o js.Value) cs.ShipDesignSpec {
	var obj cs.ShipDesignSpec
	if o.IsUndefined() {
		return obj
	}
	obj.AdditionalMassDrivers = getInt[int](o)
	obj.Armor = getInt[int](o)
	obj.BasePacketSpeed = getInt[int](o)
	obj.BeamBonus = getFloat[float64](o)
	obj.BeamDefense = getFloat[float64](o)
	obj.Bomber = getBool(o)
	obj.Bombs = GetSlice(o.Get("bombs"), GetBomb)
	obj.CanJump = getBool(o)
	obj.CanLayMines = getBool(o)
	obj.CanStealFleetCargo = getBool(o)
	obj.CanStealPlanetCargo = getBool(o)
	obj.CargoCapacity = getInt[int](o)
	obj.CloakPercent = getInt[int](o)
	obj.CloakPercentFullCargo = getInt[int](o)
	obj.CloakUnits = getInt[int](o)
	obj.Colonizer = getBool(o)
	obj.Cost = GetCost(o.Get("cost"))
	obj.Engine = GetEngine(o.Get("engine"))
	obj.EstimatedRange = getInt[int](o)
	obj.EstimatedRangeFull = getInt[int](o)
	obj.FuelCapacity = getInt[int](o)
	obj.FuelGeneration = getInt[int](o)
	obj.HasWeapons = getBool(o)
	obj.HullType = GetTechHullType(o.Get("hullType"))
	obj.ImmuneToOwnDetonation = getBool(o)
	obj.Initiative = getInt[int](o)
	obj.InnateScanRangePenFactor = getFloat[float64](o)
	obj.Mass = getInt[int](o)
	obj.MassDriver = string(getString(o))
	obj.MaxHullMass = getInt[int](o)
	obj.MaxPopulation = getInt[int](o)
	obj.MaxRange = getInt[int](o)
	obj.MineLayingRateByMineType = GetMap[map[cs.MineFieldType]int, cs.MineFieldType, int](o.Get("mineLayingRateByMineType"), GetMineFieldType, getInt)
	obj.MineSweep = getInt[int](o)
	obj.MiningRate = getInt[int](o)
	obj.Movement = getInt[int](o)
	obj.MovementBonus = getInt[int](o)
	obj.MovementFull = getInt[int](o)
	obj.NumBuilt = getInt[int](o)
	obj.NumEngines = getInt[int](o)
	obj.NumInstances = getInt[int](o)
	obj.OrbitalConstructionModule = getBool(o)
	obj.PowerRating = getInt[int](o)
	obj.Radiating = getBool(o)
	obj.ReduceCloaking = getFloat[float64](o)
	obj.ReduceMovement = getInt[int](o)
	obj.RepairBonus = getFloat[float64](o)
	obj.RetroBombs = GetSlice(o.Get("retroBombs"), GetBomb)
	obj.SafeHullMass = getInt[int](o)
	obj.SafePacketSpeed = getInt[int](o)
	obj.SafeRange = getInt[int](o)
	obj.Scanner = getBool(o)
	obj.ScanRange = getInt[int](o)
	obj.ScanRangePen = getInt[int](o)
	obj.Shields = getInt[int](o)
	obj.SmartBombs = GetSlice(o.Get("smartBombs"), GetBomb)
	obj.SpaceDock = getInt[int](o)
	obj.Starbase = getBool(o)
	obj.Stargate = string(getString(o))
	obj.TechLevel = GetTechLevel(o.Get("techLevel"))
	obj.TerraformRate = getInt[int](o)
	obj.TorpedoBonus = getFloat[float64](o)
	obj.TorpedoJamming = getFloat[float64](o)
	obj.WeaponSlots = GetSlice(o.Get("weaponSlots"), GetShipDesignSlot)
	return obj
}

func SetShipDesignSpec(o js.Value, obj *cs.ShipDesignSpec) {
}

func GetShipToken(o js.Value) cs.ShipToken {
	var obj cs.ShipToken
	if o.IsUndefined() {
		return obj
	}
	obj.DesignNum = getInt[int](o)
	obj.Quantity = getInt[int](o)
	obj.Damage = getFloat[float64](o)
	obj.QuantityDamaged = getInt[int](o)
	return obj
}

func SetShipToken(o js.Value, obj *cs.ShipToken) {
}

func GetSize(o js.Value) cs.Size {
	var obj cs.Size
	if o.IsUndefined() {
		return obj
	}
	obj = cs.Size(getString(o))
	return obj
}

func SetSize(o js.Value, obj *cs.Size) {
}

func GetSpendLeftoverPointsOn(o js.Value) cs.SpendLeftoverPointsOn {
	var obj cs.SpendLeftoverPointsOn
	if o.IsUndefined() {
		return obj
	}
	obj = cs.SpendLeftoverPointsOn(getString(o))
	return obj
}

func SetSpendLeftoverPointsOn(o js.Value, obj *cs.SpendLeftoverPointsOn) {
}

func GetStartingFleet(o js.Value) cs.StartingFleet {
	var obj cs.StartingFleet
	if o.IsUndefined() {
		return obj
	}
	obj.Name = string(getString(o))
	obj.HullName = GetStartingFleetHull(o.Get("hullName"))
	obj.HullSetNumber = getInt[uint](o)
	obj.Purpose = GetShipDesignPurpose(o.Get("purpose"))
	return obj
}

func SetStartingFleet(o js.Value, obj *cs.StartingFleet) {
}

func GetStartingFleetHull(o js.Value) cs.StartingFleetHull {
	var obj cs.StartingFleetHull
	if o.IsUndefined() {
		return obj
	}
	obj = cs.StartingFleetHull(getString(o))
	return obj
}

func SetStartingFleetHull(o js.Value, obj *cs.StartingFleetHull) {
}

func GetStartingPlanet(o js.Value) cs.StartingPlanet {
	var obj cs.StartingPlanet
	if o.IsUndefined() {
		return obj
	}
	obj.Population = getInt[int](o)
	obj.HabPenaltyFactor = getFloat[float64](o)
	obj.HasStargate = getBool(o)
	obj.HasMassDriver = getBool(o)
	obj.StarbaseDesignName = string(getString(o))
	obj.StarbaseHull = string(getString(o))
	obj.StartingFleets = GetSlice(o.Get("startingFleets"), GetStartingFleet)
	return obj
}

func SetStartingPlanet(o js.Value, obj *cs.StartingPlanet) {
}

func GetStealsResearch(o js.Value) cs.StealsResearch {
	var obj cs.StealsResearch
	if o.IsUndefined() {
		return obj
	}
	obj.Energy = getFloat[float64](o)
	obj.Weapons = getFloat[float64](o)
	obj.Propulsion = getFloat[float64](o)
	obj.Construction = getFloat[float64](o)
	obj.Electronics = getFloat[float64](o)
	obj.Biotechnology = getFloat[float64](o)
	return obj
}

func SetStealsResearch(o js.Value, obj *cs.StealsResearch) {
}

func GetTags(o js.Value) cs.Tags {
	var obj cs.Tags
	if o.IsUndefined() {
		return obj
	}
	obj = GetMap[map[string]string, string, string](o, getString, getString)
	return obj
}

func SetTags(o js.Value, obj *cs.Tags) {
}

func GetTech(o js.Value) cs.Tech {
	var obj cs.Tech
	if o.IsUndefined() {
		return obj
	}
	obj.Name = string(getString(o))
	obj.Cost = GetCost(o.Get("cost"))
	obj.Requirements = GetTechRequirements(o.Get("requirements"))
	obj.Ranking = getInt[int](o)
	obj.Category = GetTechCategory(o.Get("category"))
	obj.Origin = string(getString(o))
	return obj
}

func SetTech(o js.Value, obj *cs.Tech) {
}

func GetTechCategory(o js.Value) cs.TechCategory {
	var obj cs.TechCategory
	if o.IsUndefined() {
		return obj
	}
	obj = cs.TechCategory(getString(o))
	return obj
}

func SetTechCategory(o js.Value, obj *cs.TechCategory) {
}

func GetTechCostOffset(o js.Value) cs.TechCostOffset {
	var obj cs.TechCostOffset
	if o.IsUndefined() {
		return obj
	}
	obj.Engine = getFloat[float64](o)
	obj.BeamWeapon = getFloat[float64](o)
	obj.Torpedo = getFloat[float64](o)
	obj.Bomb = getFloat[float64](o)
	obj.PlanetaryDefense = getFloat[float64](o)
	return obj
}

func SetTechCostOffset(o js.Value, obj *cs.TechCostOffset) {
}

func GetTechDefense(o js.Value) cs.TechDefense {
	var obj cs.TechDefense
	if o.IsUndefined() {
		return obj
	}
	return obj
}

func SetTechDefense(o js.Value, obj *cs.TechDefense) {
	// TechPlanetary  Object ignored
	// Defense  Object ignored
}

func GetTechField(o js.Value) cs.TechField {
	var obj cs.TechField
	if o.IsUndefined() {
		return obj
	}
	obj = cs.TechField(getString(o))
	return obj
}

func SetTechField(o js.Value, obj *cs.TechField) {
}

func GetTechHullType(o js.Value) cs.TechHullType {
	var obj cs.TechHullType
	if o.IsUndefined() {
		return obj
	}
	obj = cs.TechHullType(getString(o))
	return obj
}

func SetTechHullType(o js.Value, obj *cs.TechHullType) {
}

func GetTechLevel(o js.Value) cs.TechLevel {
	var obj cs.TechLevel
	if o.IsUndefined() {
		return obj
	}
	obj.Energy = getInt[int](o)
	obj.Weapons = getInt[int](o)
	obj.Propulsion = getInt[int](o)
	obj.Construction = getInt[int](o)
	obj.Electronics = getInt[int](o)
	obj.Biotechnology = getInt[int](o)
	return obj
}

func SetTechLevel(o js.Value, obj *cs.TechLevel) {
}

func GetTechPlanetary(o js.Value) cs.TechPlanetary {
	var obj cs.TechPlanetary
	if o.IsUndefined() {
		return obj
	}
	obj.ResetPlanet = getBool(o)
	return obj
}

func SetTechPlanetary(o js.Value, obj *cs.TechPlanetary) {
	// Tech  Object ignored
}

func GetTechPlanetaryScanner(o js.Value) cs.TechPlanetaryScanner {
	var obj cs.TechPlanetaryScanner
	if o.IsUndefined() {
		return obj
	}
	obj.ScanRange = getInt[int](o)
	obj.ScanRangePen = getInt[int](o)
	return obj
}

func SetTechPlanetaryScanner(o js.Value, obj *cs.TechPlanetaryScanner) {
	// TechPlanetary  Object ignored
}

func GetTechRequirements(o js.Value) cs.TechRequirements {
	var obj cs.TechRequirements
	if o.IsUndefined() {
		return obj
	}
	obj.PRTsDenied = GetSlice[cs.PRT](o.Get("prtsDenied"), GetPRT)
	obj.LRTsRequired = GetLRT(o.Get("lrtsRequired"))
	obj.LRTsDenied = GetLRT(o.Get("lrtsDenied"))
	obj.PRTsRequired = GetSlice[cs.PRT](o.Get("prtsRequired"), GetPRT)
	obj.HullsAllowed = GetSlice[string](o.Get("hullsAllowed"), getString)
	obj.HullsDenied = GetSlice[string](o.Get("hullsDenied"), getString)
	obj.Acquirable = getBool(o)
	return obj
}

func SetTechRequirements(o js.Value, obj *cs.TechRequirements) {
	// TechLevel  Object ignored
}

func GetTechTerraform(o js.Value) cs.TechTerraform {
	var obj cs.TechTerraform
	if o.IsUndefined() {
		return obj
	}
	obj.Ability = getInt[int](o)
	obj.HabType = GetTerraformHabType(o.Get("habType"))
	return obj
}

func SetTechTerraform(o js.Value, obj *cs.TechTerraform) {
	// Tech  Object ignored
}

func GetTerraformHabType(o js.Value) cs.TerraformHabType {
	var obj cs.TerraformHabType
	if o.IsUndefined() {
		return obj
	}
	obj = cs.TerraformHabType(getString(o))
	return obj
}

func SetTerraformHabType(o js.Value, obj *cs.TerraformHabType) {
}

func GetTransportPlan(o js.Value) cs.TransportPlan {
	var obj cs.TransportPlan
	if o.IsUndefined() {
		return obj
	}
	obj.Num = getInt[int](o)
	obj.Name = string(getString(o))
	obj.Tasks = GetWaypointTransportTasks(o.Get("tasks"))
	return obj
}

func SetTransportPlan(o js.Value, obj *cs.TransportPlan) {
}

func GetUserRole(o js.Value) cs.UserRole {
	var obj cs.UserRole
	if o.IsUndefined() {
		return obj
	}
	obj = cs.UserRole(getString(o))
	return obj
}

func SetUserRole(o js.Value, obj *cs.UserRole) {
}

func GetVector(o js.Value) cs.Vector {
	var obj cs.Vector
	if o.IsUndefined() {
		return obj
	}
	obj.X = getFloat[float64](o)
	obj.Y = getFloat[float64](o)
	return obj
}

func SetVector(o js.Value, obj *cs.Vector) {
}

func GetVictoryCondition(o js.Value) cs.VictoryCondition {
	var obj cs.VictoryCondition
	if o.IsUndefined() {
		return obj
	}
	obj = getInt[cs.VictoryCondition](o)
	return obj
}

func SetVictoryCondition(o js.Value, obj *cs.VictoryCondition) {
}

func GetWaypoint(o js.Value) cs.Waypoint {
	var obj cs.Waypoint
	if o.IsUndefined() {
		return obj
	}
	obj.Position = GetVector(o.Get("position"))
	obj.WarpSpeed = getInt[int](o)
	obj.EstFuelUsage = getInt[int](o)
	obj.Task = GetWaypointTask(o.Get("task"))
	obj.TransportTasks = GetWaypointTransportTasks(o.Get("transportTasks"))
	obj.WaitAtWaypoint = getBool(o)
	obj.LayMineFieldDuration = getInt[int](o)
	obj.PatrolRange = getInt[int](o)
	obj.PatrolWarpSpeed = getInt[int](o)
	obj.TargetType = GetMapObjectType(o.Get("targetType"))
	obj.TargetNum = getInt[int](o)
	obj.TargetPlayerNum = getInt[int](o)
	obj.TargetName = string(getString(o))
	obj.TransferToPlayer = getInt[int](o)
	obj.PartiallyComplete = getBool(o)
	return obj
}

func SetWaypoint(o js.Value, obj *cs.Waypoint) {
}

func GetWaypointTask(o js.Value) cs.WaypointTask {
	var obj cs.WaypointTask
	if o.IsUndefined() {
		return obj
	}
	obj = cs.WaypointTask(getString(o))
	return obj
}

func SetWaypointTask(o js.Value, obj *cs.WaypointTask) {
}

func GetWaypointTaskTransportAction(o js.Value) cs.WaypointTaskTransportAction {
	var obj cs.WaypointTaskTransportAction
	if o.IsUndefined() {
		return obj
	}
	obj = cs.WaypointTaskTransportAction(getString(o))
	return obj
}

func SetWaypointTaskTransportAction(o js.Value, obj *cs.WaypointTaskTransportAction) {
}

func GetWaypointTransportTask(o js.Value) cs.WaypointTransportTask {
	var obj cs.WaypointTransportTask
	if o.IsUndefined() {
		return obj
	}
	obj.Amount = getInt[int](o)
	obj.Action = GetWaypointTaskTransportAction(o.Get("action"))
	return obj
}

func SetWaypointTransportTask(o js.Value, obj *cs.WaypointTransportTask) {
}

func GetWaypointTransportTasks(o js.Value) cs.WaypointTransportTasks {
	var obj cs.WaypointTransportTasks
	if o.IsUndefined() {
		return obj
	}
	obj.Fuel = GetWaypointTransportTask(o.Get("fuel"))
	obj.Ironium = GetWaypointTransportTask(o.Get("ironium"))
	obj.Boranium = GetWaypointTransportTask(o.Get("boranium"))
	obj.Germanium = GetWaypointTransportTask(o.Get("germanium"))
	obj.Colonists = GetWaypointTransportTask(o.Get("colonists"))
	return obj
}

func SetWaypointTransportTasks(o js.Value, obj *cs.WaypointTransportTasks) {
}

func GetWormholeIntel(o js.Value) cs.WormholeIntel {
	var obj cs.WormholeIntel
	if o.IsUndefined() {
		return obj
	}
	obj.DestinationNum = getInt[int](o)
	obj.Stability = GetWormholeStability(o.Get("stability"))
	return obj
}

func SetWormholeIntel(o js.Value, obj *cs.WormholeIntel) {
	// MapObjectIntel  Object ignored
}

func GetWormholeStability(o js.Value) cs.WormholeStability {
	var obj cs.WormholeStability
	if o.IsUndefined() {
		return obj
	}
	obj = cs.WormholeStability(getString(o))
	return obj
}

func SetWormholeStability(o js.Value, obj *cs.WormholeStability) {
}

func GetWormholeStats(o js.Value) cs.WormholeStats {
	var obj cs.WormholeStats
	if o.IsUndefined() {
		return obj
	}
	obj.YearsToDegrade = getInt[int](o)
	obj.ChanceToJump = getFloat[float64](o)
	obj.JiggleDistance = getInt[int](o)
	return obj
}

func SetWormholeStats(o js.Value, obj *cs.WormholeStats) {
}
