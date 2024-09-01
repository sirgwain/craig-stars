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
	obj.Num = getInt[int](o.Get("num"))
	obj.Name = string(getString(o.Get("name")))
	obj.PrimaryTarget = GetBattleTarget(o.Get("primaryTarget"))
	obj.SecondaryTarget = GetBattleTarget(o.Get("secondaryTarget"))
	obj.Tactic = GetBattleTactic(o.Get("tactic"))
	obj.AttackWho = GetBattleAttackWho(o.Get("attackWho"))
	obj.DumpCargo = getBool(o.Get("dumpCargo"))
	return obj
}

func SetBattlePlan(o js.Value, obj *cs.BattlePlan) {
}

func GetBattleRecord(o js.Value) cs.BattleRecord {
	var obj cs.BattleRecord
	if o.IsUndefined() {
		return obj
	}
	obj.Num = getInt[int](o.Get("num"))
	obj.PlanetNum = getInt[int](o.Get("planetNum"))
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
	obj.Num = getInt[int](o.Get("num"))
	obj.PlayerNum = getInt[int](o.Get("playerNum"))
	obj.DesignNum = getInt[int](o.Get("designNum"))
	obj.Quantity = getInt[int](o.Get("quantity"))
	return obj
}

func SetBattleRecordDestroyedToken(o js.Value, obj *cs.BattleRecordDestroyedToken) {
}

func GetBattleRecordStats(o js.Value) cs.BattleRecordStats {
	var obj cs.BattleRecordStats
	if o.IsUndefined() {
		return obj
	}
	obj.NumPlayers = getInt[int](o.Get("numPlayers"))
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
	obj.Num = getInt[int](o.Get("num"))
	obj.PlayerNum = getInt[int](o.Get("playerNum"))
	obj.DesignNum = getInt[int](o.Get("designNum"))
	obj.Position = GetBattleVector(o.Get("position"))
	obj.Initiative = getInt[int](o.Get("initiative"))
	obj.Mass = getInt[int](o.Get("mass"))
	obj.Armor = getInt[int](o.Get("armor"))
	obj.StackShields = getInt[int](o.Get("stackShields"))
	obj.Movement = getInt[int](o.Get("movement"))
	obj.StartingQuantity = getInt[int](o.Get("startingQuantity"))
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
	obj.TokenNum = getInt[int](o.Get("tokenNum"))
	obj.Round = getInt[int](o.Get("round"))
	obj.From = GetBattleVector(o.Get("from"))
	obj.To = GetBattleVector(o.Get("to"))
	obj.Slot = getInt[int](o.Get("slot"))
	obj.TargetNum = getInt[int](o.Get("targetNum"))
	obj.Target = getPointer(GetShipToken(o.Get("target")))
	obj.TokensDestroyed = getInt[int](o.Get("tokensDestroyed"))
	obj.DamageDoneShields = getInt[int](o.Get("damageDoneShields"))
	obj.DamageDoneArmor = getInt[int](o.Get("damageDoneArmor"))
	obj.TorpedoHits = getInt[int](o.Get("torpedoHits"))
	obj.TorpedoMisses = getInt[int](o.Get("torpedoMisses"))
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
	obj.X = getInt[int](o.Get("x"))
	obj.Y = getInt[int](o.Get("y"))
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
	obj.Quantity = getInt[int](o.Get("quantity"))
	obj.KillRate = getFloat[float64](o.Get("killRate"))
	obj.MinKillRate = getInt[int](o.Get("minKillRate"))
	obj.StructureDestroyRate = getFloat[float64](o.Get("structureDestroyRate"))
	obj.UnterraformRate = getInt[int](o.Get("unterraformRate"))
	return obj
}

func SetBomb(o js.Value, obj *cs.Bomb) {
}

func GetBombingResult(o js.Value) cs.BombingResult {
	var obj cs.BombingResult
	if o.IsUndefined() {
		return obj
	}
	obj.BomberName = string(getString(o.Get("bomberName")))
	obj.NumBombers = getInt[int](o.Get("numBombers"))
	obj.ColonistsKilled = getInt[int](o.Get("colonistsKilled"))
	obj.MinesDestroyed = getInt[int](o.Get("minesDestroyed"))
	obj.FactoriesDestroyed = getInt[int](o.Get("factoriesDestroyed"))
	obj.DefensesDestroyed = getInt[int](o.Get("defensesDestroyed"))
	obj.UnterraformAmount = GetHab(o.Get("unterraformAmount"))
	obj.PlanetEmptied = getBool(o.Get("planetEmptied"))
	return obj
}

func SetBombingResult(o js.Value, obj *cs.BombingResult) {
}

func GetCargo(o js.Value) cs.Cargo {
	var obj cs.Cargo
	if o.IsUndefined() {
		return obj
	}
	obj.Ironium = getInt[int](o.Get("ironium"))
	obj.Boranium = getInt[int](o.Get("boranium"))
	obj.Germanium = getInt[int](o.Get("germanium"))
	obj.Colonists = getInt[int](o.Get("colonists"))
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
	obj.AllMinerals = getInt[int](o.Get("minMinerals"))
	obj.AllRandomMinerals = getInt[int](o.Get("randomMinerals"))
	obj.BonusMinerals = getInt[int](o.Get("bonusMinerals"))
	obj.BonusRandomMinerals = getInt[int](o.Get("bonusRandomMinerals"))
	obj.BonusMinConcentration = getInt[int](o.Get("minConcentrationBonus"))
	obj.BonusRandomConcentration = getInt[int](o.Get("randomConcentrationBonus"))
	obj.BonusAffectsMinerals = getInt[int](o.Get("affectsMinerals"))
	obj.MinTerraform = getInt[int](o.Get("minTerraform"))
	obj.RandomTerraform = getInt[int](o.Get("randomTerraform"))
	obj.AffectsHabs = getInt[int](o.Get("affectsHabs"))
	obj.PopKilledPercent = getFloat[float64](o.Get("popKilledPercent"))
	return obj
}

func SetCometStats(o js.Value, obj *cs.CometStats) {
}

func GetCost(o js.Value) cs.Cost {
	var obj cs.Cost
	if o.IsUndefined() {
		return obj
	}
	obj.Ironium = getInt[int](o.Get("ironium"))
	obj.Boranium = getInt[int](o.Get("boranium"))
	obj.Germanium = getInt[int](o.Get("germanium"))
	obj.Resources = getInt[int](o.Get("resources"))
	return obj
}

func SetCost(o js.Value, obj *cs.Cost) {
}

func GetDBObject(o js.Value) cs.DBObject {
	var obj cs.DBObject
	if o.IsUndefined() {
		return obj
	}
	obj.ID = getInt[int64](o.Get("id"))
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
	obj.DefenseCoverage = getFloat[float64](o.Get("defenseCoverage"))
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
	obj.IdealSpeed = getInt[int](o.Get("idealSpeed"))
	obj.FreeSpeed = getInt[int](o.Get("freeSpeed"))
	obj.MaxSafeSpeed = getInt[int](o.Get("maxSafeSpeed"))
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
	obj.PlanetNum = getInt[int](o.Get("planetNum"))
	obj.BaseName = string(getString(o.Get("baseName")))
	obj.Cargo = GetCargo(o.Get("cargo"))
	obj.Fuel = getInt[int](o.Get("fuel"))
	obj.Age = getInt[int](o.Get("age"))
	obj.Tokens = GetSlice(o.Get("tokens"), GetShipToken)
	obj.Heading = GetVector(o.Get("heading"))
	obj.WarpSpeed = getInt[int](o.Get("warpSpeed"))
	obj.PreviousPosition = getPointer(GetVector(o.Get("previousPosition")))
	obj.OrbitingPlanetNum = getInt[int](o.Get("orbitingPlanetNum"))
	obj.Starbase = getBool(o.Get("starbase"))
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
	obj.BaseName = string(getString(o.Get("baseName")))
	obj.Heading = GetVector(o.Get("heading"))
	obj.OrbitingPlanetNum = getInt[int](o.Get("orbitingPlanetNum"))
	obj.WarpSpeed = getInt[int](o.Get("warpSpeed"))
	obj.Mass = getInt[int](o.Get("mass"))
	obj.Cargo = GetCargo(o.Get("cargo"))
	obj.CargoDiscovered = getBool(o.Get("cargoDiscovered"))
	obj.Freighter = getBool(o.Get("freighter"))
	obj.ScanRange = getInt[int](o.Get("scanRange"))
	obj.ScanRangePen = getInt[int](o.Get("scanRangePen"))
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
	obj.RepeatOrders = getBool(o.Get("repeatOrders"))
	obj.BattlePlanNum = getInt[int](o.Get("battlePlanNum"))
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
	obj.BaseCloakedCargo = getInt[int](o.Get("baseCloakedCargo"))
	obj.BasePacketSpeed = getInt[int](o.Get("basePacketSpeed"))
	obj.HasMassDriver = getBool(o.Get("hasMassDriver"))
	obj.HasStargate = getBool(o.Get("hasStargate"))
	obj.MassDriver = string(getString(o.Get("massDriver")))
	obj.MassEmpty = getInt[int](o.Get("massEmpty"))
	obj.MaxHullMass = getInt[int](o.Get("maxHullMass"))
	obj.MaxRange = getInt[int](o.Get("maxRange"))
	obj.Purposes = GetMap[map[cs.ShipDesignPurpose]bool, cs.ShipDesignPurpose, bool](o.Get("purposes"), GetShipDesignPurpose, getBool)
	obj.SafeHullMass = getInt[int](o.Get("safeHullMass"))
	obj.SafeRange = getInt[int](o.Get("safeRange"))
	obj.Stargate = string(getString(o.Get("stargate")))
	obj.TotalShips = getInt[int](o.Get("totalShips"))
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
	obj.ID = getInt[int64](o.Get("id"))
	obj.GameID = getInt[int64](o.Get("gameId"))
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
	obj.Grav = getInt[int](o.Get("grav"))
	obj.Temp = getInt[int](o.Get("temp"))
	obj.Rad = getInt[int](o.Get("rad"))
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
	obj.Name = string(getString(o.Get("name")))
	obj.Num = getInt[int](o.Get("num"))
	obj.PlayerNum = getInt[int](o.Get("playerNum"))
	obj.ReportAge = getInt[int](o.Get("reportAge"))
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
	obj.PointCost = getInt[int](o.Get("pointCost"))
	obj.StartingTechLevels = GetTechLevel(o.Get("startingTechLevels"))
	obj.TechCostOffset = GetTechCostOffset(o.Get("techCostOffset"))
	obj.NewTechCostFactorOffset = getFloat[float64](o.Get("newTechCostFactorOffset"))
	obj.MiniaturizationMax = getFloat[float64](o.Get("miniaturizationMax"))
	obj.MiniaturizationPerLevel = getFloat[float64](o.Get("miniaturizationPerLevel"))
	obj.NoAdvancedScanners = getBool(o.Get("noAdvancedScanners"))
	obj.ScanRangeFactorOffset = getFloat[float64](o.Get("scanRangeFactorOffset"))
	obj.FuelEfficiencyOffset = getFloat[float64](o.Get("fuelEfficiencyOffset"))
	obj.MaxPopulationOffset = getFloat[float64](o.Get("maxPopulationOffset"))
	obj.TerraformCostOffset = GetCost(o.Get("terraformCostOffset"))
	obj.MineralAlchemyCostOffset = getInt[int](o.Get("mineralAlchemyCostOffset"))
	obj.ScrapMineralOffset = getFloat[float64](o.Get("scrapMineralOffset"))
	obj.ScrapMineralOffsetStarbase = getFloat[float64](o.Get("scrapMineralOffsetStarbase"))
	obj.ScrapResourcesOffset = getFloat[float64](o.Get("scrapResourcesOffset"))
	obj.ScrapResourcesOffsetStarbase = getFloat[float64](o.Get("scrapResourcesOffsetStarbase"))
	obj.StartingPopulationFactorDelta = getFloat[float64](o.Get("startingPopulationFactorDelta"))
	obj.StarbaseBuiltInCloakUnits = getInt[int](o.Get("starbaseBuiltInCloakUnits"))
	obj.StarbaseCostFactorOffset = getFloat[float64](o.Get("starbaseCostFactorOffset"))
	obj.ResearchFactorOffset = getFloat[float64](o.Get("researchFactorOffset"))
	obj.ResearchSplashDamage = getFloat[float64](o.Get("researchSplashDamage"))
	obj.ShieldStrengthFactorOffset = getFloat[float64](o.Get("shieldStrengthFactorOffset"))
	obj.ShieldRegenerationRateOffset = getFloat[float64](o.Get("shieldRegenerationRateOffset"))
	obj.ArmorStrengthFactorOffset = getFloat[float64](o.Get("armorStrengthFactorOffset"))
	obj.EngineFailureRateOffset = getFloat[float64](o.Get("engineFailureRateOffset"))
	obj.EngineReliableSpeed = getInt[int](o.Get("engineReliableSpeed"))
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
	obj.Num = getInt[int](o.Get("num"))
	obj.PlayerNum = getInt[int](o.Get("playerNum"))
	obj.Name = string(getString(o.Get("name")))
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
	obj.NumMines = getInt[int](o.Get("numMines"))
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
	obj.NumMines = getInt[int](o.Get("numMines"))
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
	obj.Detonate = getBool(o.Get("detonate"))
	return obj
}

func SetMineFieldOrders(o js.Value, obj *cs.MineFieldOrders) {
}

func GetMineFieldSpec(o js.Value) cs.MineFieldSpec {
	var obj cs.MineFieldSpec
	if o.IsUndefined() {
		return obj
	}
	obj.Radius = getFloat[float64](o.Get("radius"))
	obj.DecayRate = getInt[int](o.Get("decayRate"))
	return obj
}

func SetMineFieldSpec(o js.Value, obj *cs.MineFieldSpec) {
}

func GetMineFieldStats(o js.Value) cs.MineFieldStats {
	var obj cs.MineFieldStats
	if o.IsUndefined() {
		return obj
	}
	obj.MinDamagePerFleetRS = getInt[int](o.Get("minDamagePerFleetRS"))
	obj.DamagePerEngineRS = getInt[int](o.Get("damagePerEngineRS"))
	obj.MaxSpeed = getInt[int](o.Get("maxSpeed"))
	obj.ChanceOfHit = getFloat[float64](o.Get("chanceOfHit"))
	obj.MinDamagePerFleet = getInt[int](o.Get("minDamagePerFleet"))
	obj.DamagePerEngine = getInt[int](o.Get("damagePerEngine"))
	obj.SweepFactor = getFloat[float64](o.Get("sweepFactor"))
	obj.MinDecay = getInt[int](o.Get("minDecay"))
	obj.CanDetonate = getBool(o.Get("canDetonate"))
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
	obj.Ironium = getInt[int](o.Get("ironium"))
	obj.Boranium = getInt[int](o.Get("boranium"))
	obj.Germanium = getInt[int](o.Get("germanium"))
	return obj
}

func SetMineral(o js.Value, obj *cs.Mineral) {
}

func GetMineralPacketDamage(o js.Value) cs.MineralPacketDamage {
	var obj cs.MineralPacketDamage
	if o.IsUndefined() {
		return obj
	}
	obj.Killed = getInt[int](o.Get("killed"))
	obj.DefensesDestroyed = getInt[int](o.Get("defensesDestroyed"))
	obj.Uncaught = getInt[int](o.Get("uncaught"))
	return obj
}

func SetMineralPacketDamage(o js.Value, obj *cs.MineralPacketDamage) {
}

func GetMineralPacketIntel(o js.Value) cs.MineralPacketIntel {
	var obj cs.MineralPacketIntel
	if o.IsUndefined() {
		return obj
	}
	obj.WarpSpeed = getInt[int](o.Get("warpSpeed"))
	obj.Heading = GetVector(o.Get("heading"))
	obj.Cargo = GetCargo(o.Get("cargo"))
	obj.TargetPlanetNum = getInt[int](o.Get("targetPlanetNum"))
	obj.ScanRange = getInt[int](o.Get("scanRange"))
	obj.ScanRangePen = getInt[int](o.Get("scanRangePen"))
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
	obj.NewTechCostFactor = getFloat[float64](o.Get("newTechCostFactor"))
	obj.MiniaturizationMax = getFloat[float64](o.Get("miniaturizationMax"))
	obj.MiniaturizationPerLevel = getFloat[float64](o.Get("miniaturizationPerLevel"))
	return obj
}

func SetMiniaturizationSpec(o js.Value, obj *cs.MiniaturizationSpec) {
}

func GetMysteryTrader(o js.Value) cs.MysteryTrader {
	var obj cs.MysteryTrader
	if o.IsUndefined() {
		return obj
	}
	obj.WarpSpeed = getInt[int](o.Get("warpSpeed"))
	obj.Destination = GetVector(o.Get("destination"))
	obj.RequestedBoon = getInt[int](o.Get("requestedBoon"))
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
	obj.WarpSpeed = getInt[int](o.Get("warpSpeed"))
	obj.Heading = GetVector(o.Get("heading"))
	obj.RequestedBoon = getInt[int](o.Get("requestedBoon"))
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
	obj.Tech = string(getString(o.Get("tech")))
	obj.Ship = GetShipDesign(o.Get("ship"))
	obj.ShipCount = getInt[int](o.Get("shipCount"))
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
	obj.ChanceMaxTechGetsPart = getInt[int](o.Get("chanceMaxTechGetsPart"))
	obj.ChanceCourseChange = getInt[int](o.Get("chanceCourseChange"))
	obj.ChanceSpeedUpOnly = getInt[int](o.Get("chanceSpeedUpOnly"))
	obj.ChanceAgain = getInt[int](o.Get("chanceAgain"))
	obj.MinYear = getInt[int](o.Get("minYear"))
	obj.EvenYearOnly = getBool(o.Get("evenYearOnly"))
	obj.MinWarp = getInt[int](o.Get("minWarp"))
	obj.MaxWarp = getInt[int](o.Get("maxWarp"))
	obj.MaxMysteryTraders = getInt[int](o.Get("maxMysteryTraders"))
	obj.RequestedBoon = getInt[int](o.Get("requestedBoon"))
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
	obj.MineralsGiven = getInt[int](o.Get("mineralsGiven"))
	obj.Reward = getInt[int](o.Get("reward"))
	return obj
}

func SetMysteryTraderTechBoonMineralsReward(o js.Value, obj *cs.MysteryTraderTechBoonMineralsReward) {
}

func GetMysteryTraderTechBoonRules(o js.Value) cs.MysteryTraderTechBoonRules {
	var obj cs.MysteryTraderTechBoonRules
	if o.IsUndefined() {
		return obj
	}
	obj.TechLevels = getInt[int](o.Get("techLevels"))
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
	obj.PointCost = getInt[int](o.Get("pointCost"))
	obj.StartingTechLevels = GetTechLevel(o.Get("startingTechLevels"))
	obj.StartingPlanets = GetSlice(o.Get("startingPlanets"), GetStartingPlanet)
	obj.TechCostOffset = GetTechCostOffset(o.Get("techCostOffset"))
	obj.MineralsPerSingleMineralPacket = getInt[int](o.Get("mineralsPerSingleMineralPacket"))
	obj.MineralsPerMixedMineralPacket = getInt[int](o.Get("mineralsPerMixedMineralPacket"))
	obj.PacketResourceCost = getInt[int](o.Get("packetResourceCost"))
	obj.PacketMineralCostFactor = getFloat[float64](o.Get("packetMineralCostFactor"))
	obj.PacketReceiverFactor = getFloat[float64](o.Get("packetReceiverFactor"))
	obj.PacketDecayFactor = getFloat[float64](o.Get("packetDecayFactor"))
	obj.PacketOverSafeWarpPenalty = getInt[int](o.Get("packetOverSafeWarpPenalty"))
	obj.PacketBuiltInScanner = getBool(o.Get("packetBuiltInScanner"))
	obj.DetectPacketDestinationStarbases = getBool(o.Get("detectPacketDestinationStarbases"))
	obj.DetectAllPackets = getBool(o.Get("detectAllPackets"))
	obj.PacketTerraformChance = getFloat[float64](o.Get("packetTerraformChance"))
	obj.PacketPermaformChance = getFloat[float64](o.Get("packetPermaformChance"))
	obj.PacketPermaTerraformSizeUnit = getInt[int](o.Get("packetPermaTerraformSizeUnit"))
	obj.CanGateCargo = getBool(o.Get("canGateCargo"))
	obj.CanDetectStargatePlanets = getBool(o.Get("canDetectStargatePlanets"))
	obj.ShipsVanishInVoid = getBool(o.Get("shipsVanishInVoid"))
	obj.BuiltInScannerMultiplier = getInt[int](o.Get("builtInScannerMultiplier"))
	obj.TechsCostExtraLevel = getInt[int](o.Get("techsCostExtraLevel"))
	obj.FreighterGrowthFactor = getFloat[float64](o.Get("freighterGrowthFactor"))
	obj.GrowthFactor = getFloat[float64](o.Get("growthFactor"))
	obj.MaxPopulationOffset = getFloat[float64](o.Get("maxPopulationOffset"))
	obj.BuiltInCloakUnits = getInt[int](o.Get("builtInCloakUnits"))
	obj.StealsResearch = GetStealsResearch(o.Get("stealsResearch"))
	obj.FreeCargoCloaking = getBool(o.Get("freeCargoCloaking"))
	obj.MineFieldsAreScanners = getBool(o.Get("mineFieldsAreScanners"))
	obj.MineFieldRateMoveFactor = getFloat[float64](o.Get("mineFieldRateMoveFactor"))
	obj.MineFieldSafeWarpBonus = getInt[int](o.Get("mineFieldSafeWarpBonus"))
	obj.MineFieldMinDecayFactor = getFloat[float64](o.Get("mineFieldMinDecayFactor"))
	obj.MineFieldBaseDecayRate = getFloat[float64](o.Get("mineFieldBaseDecayRate"))
	obj.MineFieldPlanetDecayRate = getFloat[float64](o.Get("mineFieldPlanetDecayRate"))
	obj.MineFieldMaxDecayRate = getFloat[float64](o.Get("mineFieldMaxDecayRate"))
	obj.CanDetonateMineFields = getBool(o.Get("canDetonateMineFields"))
	obj.MineFieldDetonateDecayRate = getFloat[float64](o.Get("mineFieldDetonateDecayRate"))
	obj.DiscoverDesignOnScan = getBool(o.Get("discoverDesignOnScan"))
	obj.CanRemoteMineOwnPlanets = getBool(o.Get("canRemoteMineOwnPlanets"))
	obj.InvasionAttackBonus = getFloat[float64](o.Get("invasionAttackBonus"))
	obj.InvasionDefendBonus = getFloat[float64](o.Get("invasionDefendBonus"))
	obj.MovementBonus = getInt[int](o.Get("movementBonus"))
	obj.Instaforming = getBool(o.Get("instaforming"))
	obj.PermaformChance = getFloat[float64](o.Get("permaformChance"))
	obj.PermaformPopulation = getInt[int](o.Get("permaformPopulation"))
	obj.RepairFactor = getFloat[float64](o.Get("repairFactor"))
	obj.StarbaseRepairFactor = getFloat[float64](o.Get("starbaseRepairFactor"))
	obj.StarbaseCostFactor = getFloat[float64](o.Get("starbaseCostFactor"))
	obj.InnateMining = getBool(o.Get("innateMining"))
	obj.InnateResources = getBool(o.Get("innateResources"))
	obj.InnateScanner = getBool(o.Get("innateScanner"))
	obj.InnatePopulationFactor = getFloat[float64](o.Get("innatePopulationFactor"))
	obj.CanBuildDefenses = getBool(o.Get("canBuildDefenses"))
	obj.LivesOnStarbases = getBool(o.Get("livesOnStarbases"))
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
	obj.Mines = getInt[int](o.Get("mines"))
	obj.Factories = getInt[int](o.Get("factories"))
	obj.Defenses = getInt[int](o.Get("defenses"))
	obj.Homeworld = getBool(o.Get("homeworld"))
	obj.Scanner = getBool(o.Get("scanner"))
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
	obj.CargoDiscovered = getBool(o.Get("cargoDiscovered"))
	obj.PlanetHabitability = getInt[int](o.Get("planetHabitability"))
	obj.PlanetHabitabilityTerraformed = getInt[int](o.Get("planetHabitabilityTerraformed"))
	obj.Homeworld = getBool(o.Get("homeworld"))
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
	obj.ContributesOnlyLeftoverToResearch = getBool(o.Get("contributesOnlyLeftoverToResearch"))
	obj.ProductionQueue = GetSlice(o.Get("productionQueue"), GetProductionQueueItem)
	obj.RouteTargetType = GetMapObjectType(o.Get("routeTargetType"))
	obj.RouteTargetNum = getInt[int](o.Get("routeTargetNum"))
	obj.RouteTargetPlayerNum = getInt[int](o.Get("routeTargetPlayerNum"))
	obj.PacketTargetNum = getInt[int](o.Get("packetTargetNum"))
	obj.PacketSpeed = getInt[int](o.Get("packetSpeed"))
	return obj
}

func SetPlanetOrders(o js.Value, obj *cs.PlanetOrders) {
}

func GetPlanetSpec(o js.Value) cs.PlanetSpec {
	var obj cs.PlanetSpec
	if o.IsUndefined() {
		return obj
	}
	obj.CanTerraform = getBool(o.Get("canTerraform"))
	obj.Defense = string(getString(o.Get("defense")))
	obj.DefenseCoverage = getFloat[float64](o.Get("defenseCoverage"))
	obj.DefenseCoverageSmart = getFloat[float64](o.Get("defenseCoverageSmart"))
	obj.GrowthAmount = getInt[int](o.Get("growthAmount"))
	obj.Habitability = getInt[int](o.Get("habitability"))
	obj.MaxDefenses = getInt[int](o.Get("maxDefenses"))
	obj.MaxFactories = getInt[int](o.Get("maxFactories"))
	obj.MaxMines = getInt[int](o.Get("maxMines"))
	obj.MaxPopulation = getInt[int](o.Get("maxPopulation"))
	obj.MaxPossibleFactories = getInt[int](o.Get("maxPossibleFactories"))
	obj.MaxPossibleMines = getInt[int](o.Get("maxPossibleMines"))
	obj.MiningOutput = GetMineral(o.Get("miningOutput"))
	obj.Population = getInt[int](o.Get("population"))
	obj.PopulationDensity = getFloat[float64](o.Get("populationDensity"))
	obj.ResourcesPerYear = getInt[int](o.Get("resourcesPerYear"))
	obj.ResourcesPerYearAvailable = getInt[int](o.Get("resourcesPerYearAvailable"))
	obj.ResourcesPerYearResearch = getInt[int](o.Get("resourcesPerYearResearch"))
	obj.ResourcesPerYearResearchEstimatedLeftover = getInt[int](o.Get("resourcesPerYearResearchEstimatedLeftover"))
	obj.Scanner = string(getString(o.Get("scanner")))
	obj.ScanRange = getInt[int](o.Get("scanRange"))
	obj.ScanRangePen = getInt[int](o.Get("scanRangePen"))
	obj.TerraformAmount = GetHab(o.Get("terraformAmount"))
	obj.MinTerraformAmount = GetHab(o.Get("minTerraformAmount"))
	obj.TerraformedHabitability = getInt[int](o.Get("terraformedHabitability"))
	obj.Contested = getBool(o.Get("contested"))
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
	obj.HasMassDriver = getBool(o.Get("hasMassDriver"))
	obj.HasStarbase = getBool(o.Get("hasStarbase"))
	obj.HasStargate = getBool(o.Get("hasStargate"))
	obj.StarbaseDesignName = string(getString(o.Get("starbaseDesignName")))
	obj.StarbaseDesignNum = getInt[int](o.Get("starbaseDesignNum"))
	obj.DockCapacity = getInt[int](o.Get("dockCapacity"))
	obj.BasePacketSpeed = getInt[int](o.Get("basePacketSpeed"))
	obj.SafePacketSpeed = getInt[int](o.Get("safePacketSpeed"))
	obj.SafeHullMass = getInt[int](o.Get("safeHullMass"))
	obj.SafeRange = getInt[int](o.Get("safeRange"))
	obj.MaxRange = getInt[int](o.Get("maxRange"))
	obj.MaxHullMass = getInt[int](o.Get("maxHullMass"))
	obj.Stargate = string(getString(o.Get("stargate")))
	obj.MassDriver = string(getString(o.Get("massDriver")))
	return obj
}

func SetPlanetStarbaseSpec(o js.Value, obj *cs.PlanetStarbaseSpec) {
}

func GetPlayer(o js.Value) cs.Player {
	var obj cs.Player
	if o.IsUndefined() {
		return obj
	}
	obj.UserID = getInt[int64](o.Get("userId"))
	obj.Name = string(getString(o.Get("name")))
	obj.Num = getInt[int](o.Get("num"))
	obj.Ready = getBool(o.Get("ready"))
	obj.AIControlled = getBool(o.Get("aiControlled"))
	obj.AIDifficulty = GetAIDifficulty(o.Get("aiDifficulty"))
	obj.Guest = getBool(o.Get("guest"))
	obj.SubmittedTurn = getBool(o.Get("submittedTurn"))
	obj.Color = string(getString(o.Get("color")))
	obj.DefaultHullSet = getInt[int](o.Get("defaultHullSet"))
	obj.Race = GetRace(o.Get("race"))
	obj.TechLevels = GetTechLevel(o.Get("techLevels"))
	obj.TechLevelsSpent = GetTechLevel(o.Get("techLevelsSpent"))
	obj.ResearchSpentLastYear = getInt[int](o.Get("researchSpentLastYear"))
	obj.Relations = GetSlice(o.Get("relations"), GetPlayerRelationship)
	obj.Messages = GetSlice(o.Get("messages"), GetPlayerMessage)
	obj.Designs = GetPointerSlice(o.Get("designs"), GetShipDesign)
	obj.ScoreHistory = GetSlice(o.Get("scoreHistory"), GetPlayerScore)
	obj.AcquiredTechs = GetMap[map[string]bool, string, bool](o.Get("acquiredTechs"), getString, getBool)
	obj.AchievedVictoryConditions = GetBitmask(o.Get("achievedVictoryConditions"))
	obj.Victor = getBool(o.Get("victor"))
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
	obj.Name = string(getString(o.Get("name")))
	obj.Num = getInt[int](o.Get("num"))
	obj.Color = string(getString(o.Get("color")))
	obj.Seen = getBool(o.Get("seen"))
	obj.RaceName = string(getString(o.Get("raceName")))
	obj.RacePluralName = string(getString(o.Get("racePluralName")))
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
	obj.Text = string(getString(o.Get("text")))
	obj.BattleNum = getInt[int](o.Get("battleNum"))
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
	obj.Amount = getInt[int](o.Get("amount"))
	obj.Amount2 = getInt[int](o.Get("amount2"))
	obj.PrevAmount = getInt[int](o.Get("prevAmount"))
	obj.SourcePlayerNum = getInt[int](o.Get("sourcePlayerNum"))
	obj.DestPlayerNum = getInt[int](o.Get("destPlayerNum"))
	obj.Name = string(getString(o.Get("name")))
	obj.Cost = getPointer(GetCost(o.Get("cost")))
	obj.Mineral = getPointer(GetMineral(o.Get("mineral")))
	obj.Cargo = getPointer(GetCargo(o.Get("cargo")))
	obj.QueueItemType = GetQueueItemType(o.Get("queueItemType"))
	obj.Field = GetTechField(o.Get("field"))
	obj.NextField = GetTechField(o.Get("nextField"))
	obj.TechGained = string(getString(o.Get("techGained")))
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
	obj.ColonistsKilled = getInt[int](o.Get("colonistsKilled"))
	return obj
}

func SetPlayerMessageSpecComet(o js.Value, obj *cs.PlayerMessageSpecComet) {
}

func GetPlayerMessageSpecMysteryTrader(o js.Value) cs.PlayerMessageSpecMysteryTrader {
	var obj cs.PlayerMessageSpecMysteryTrader
	if o.IsUndefined() {
		return obj
	}
	obj.FleetNum = getInt[int](o.Get("fleetNum"))
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
	obj.ResearchAmount = getInt[int](o.Get("researchAmount"))
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
	obj.ShareMap = getBool(o.Get("shareMap"))
	return obj
}

func SetPlayerRelationship(o js.Value, obj *cs.PlayerRelationship) {
}

func GetPlayerScore(o js.Value) cs.PlayerScore {
	var obj cs.PlayerScore
	if o.IsUndefined() {
		return obj
	}
	obj.Planets = getInt[int](o.Get("planets"))
	obj.Starbases = getInt[int](o.Get("starbases"))
	obj.UnarmedShips = getInt[int](o.Get("unarmedShips"))
	obj.EscortShips = getInt[int](o.Get("escortShips"))
	obj.CapitalShips = getInt[int](o.Get("capitalShips"))
	obj.TechLevels = getInt[int](o.Get("techLevels"))
	obj.Resources = getInt[int](o.Get("resources"))
	obj.Score = getInt[int](o.Get("score"))
	obj.Rank = getInt[int](o.Get("rank"))
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
	obj.ResourcesPerYear = getInt[int](o.Get("resourcesPerYear"))
	obj.ResourcesPerYearResearch = getInt[int](o.Get("resourcesPerYearResearch"))
	obj.ResourcesPerYearResearchEstimated = getInt[int](o.Get("resourcesPerYearResearchEstimated"))
	obj.CurrentResearchCost = getInt[int](o.Get("currentResearchCost"))
	return obj
}

func SetPlayerSpec(o js.Value, obj *cs.PlayerSpec) {
}

func GetPlayerStats(o js.Value) cs.PlayerStats {
	var obj cs.PlayerStats
	if o.IsUndefined() {
		return obj
	}
	obj.FleetsBuilt = getInt[int](o.Get("fleetsBuilt"))
	obj.StarbasesBuilt = getInt[int](o.Get("starbasesBuilt"))
	obj.TokensBuilt = getInt[int](o.Get("tokensBuilt"))
	obj.PlanetsColonized = getInt[int](o.Get("planetsColonized"))
	return obj
}

func SetPlayerStats(o js.Value, obj *cs.PlayerStats) {
}

func GetProductionPlan(o js.Value) cs.ProductionPlan {
	var obj cs.ProductionPlan
	if o.IsUndefined() {
		return obj
	}
	obj.Num = getInt[int](o.Get("num"))
	obj.Name = string(getString(o.Get("name")))
	obj.Items = GetSlice(o.Get("items"), GetProductionPlanItem)
	obj.ContributesOnlyLeftoverToResearch = getBool(o.Get("contributesOnlyLeftoverToResearch"))
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
	obj.DesignNum = getInt[int](o.Get("designNum"))
	obj.Quantity = getInt[int](o.Get("quantity"))
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
	obj.DesignNum = getInt[int](o.Get("designNum"))
	obj.Quantity = getInt[int](o.Get("quantity"))
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
	obj.Skipped = getBool(o.Get("skipped"))
	obj.YearsToBuildOne = getInt[int](o.Get("yearsToBuildOne"))
	obj.YearsToBuildAll = getInt[int](o.Get("yearsToBuildAll"))
	obj.YearsToSkipAuto = getInt[int](o.Get("yearsToSkipAuto"))
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
	obj.UserID = getInt[int64](o.Get("userId"))
	obj.Name = string(getString(o.Get("name")))
	obj.PluralName = string(getString(o.Get("pluralName")))
	obj.SpendLeftoverPointsOn = GetSpendLeftoverPointsOn(o.Get("spendLeftoverPointsOn"))
	obj.PRT = GetPRT(o.Get("prt"))
	obj.LRTs = GetBitmask(o.Get("lrts"))
	obj.HabLow = GetHab(o.Get("habLow"))
	obj.HabHigh = GetHab(o.Get("habHigh"))
	obj.GrowthRate = getInt[int](o.Get("growthRate"))
	obj.PopEfficiency = getInt[int](o.Get("popEfficiency"))
	obj.FactoryOutput = getInt[int](o.Get("factoryOutput"))
	obj.FactoryCost = getInt[int](o.Get("factoryCost"))
	obj.NumFactories = getInt[int](o.Get("numFactories"))
	obj.FactoriesCostLess = getBool(o.Get("factoriesCostLess"))
	obj.ImmuneGrav = getBool(o.Get("immuneGrav"))
	obj.ImmuneTemp = getBool(o.Get("immuneTemp"))
	obj.ImmuneRad = getBool(o.Get("immuneRad"))
	obj.MineOutput = getInt[int](o.Get("mineOutput"))
	obj.MineCost = getInt[int](o.Get("mineCost"))
	obj.NumMines = getInt[int](o.Get("numMines"))
	obj.ResearchCost = GetResearchCost(o.Get("researchCost"))
	obj.TechsStartHigh = getBool(o.Get("techsStartHigh"))
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
	obj.MineralsPerSingleMineralPacket = getInt[int](o.Get("mineralsPerSingleMineralPacket"))
	obj.MineralsPerMixedMineralPacket = getInt[int](o.Get("mineralsPerMixedMineralPacket"))
	obj.PacketResourceCost = getInt[int](o.Get("packetResourceCost"))
	obj.PacketMineralCostFactor = getFloat[float64](o.Get("packetMineralCostFactor"))
	obj.PacketReceiverFactor = getFloat[float64](o.Get("packetReceiverFactor"))
	obj.PacketDecayFactor = getFloat[float64](o.Get("packetDecayFactor"))
	obj.PacketOverSafeWarpPenalty = getInt[int](o.Get("packetOverSafeWarpPenalty"))
	obj.PacketBuiltInScanner = getBool(o.Get("packetBuiltInScanner"))
	obj.DetectPacketDestinationStarbases = getBool(o.Get("detectPacketDestinationStarbases"))
	obj.DetectAllPackets = getBool(o.Get("detectAllPackets"))
	obj.PacketTerraformChance = getFloat[float64](o.Get("packetTerraformChance"))
	obj.PacketPermaformChance = getFloat[float64](o.Get("packetPermaformChance"))
	obj.PacketPermaTerraformSizeUnit = getInt[int](o.Get("packetPermaTerraformSizeUnit"))
	obj.CanGateCargo = getBool(o.Get("canGateCargo"))
	obj.CanDetectStargatePlanets = getBool(o.Get("canDetectStargatePlanets"))
	obj.ShipsVanishInVoid = getBool(o.Get("shipsVanishInVoid"))
	obj.TechsCostExtraLevel = getInt[int](o.Get("techsCostExtraLevel"))
	obj.FreighterGrowthFactor = getFloat[float64](o.Get("freighterGrowthFactor"))
	obj.GrowthFactor = getFloat[float64](o.Get("growthFactor"))
	obj.MaxPopulationOffset = getFloat[float64](o.Get("maxPopulationOffset"))
	obj.BuiltInCloakUnits = getInt[int](o.Get("builtInCloakUnits"))
	obj.StealsResearch = GetStealsResearch(o.Get("stealsResearch"))
	obj.FreeCargoCloaking = getBool(o.Get("freeCargoCloaking"))
	obj.MineFieldsAreScanners = getBool(o.Get("mineFieldsAreScanners"))
	obj.MineFieldRateMoveFactor = getFloat[float64](o.Get("mineFieldRateMoveFactor"))
	obj.MineFieldSafeWarpBonus = getInt[int](o.Get("mineFieldSafeWarpBonus"))
	obj.MineFieldMinDecayFactor = getFloat[float64](o.Get("mineFieldMinDecayFactor"))
	obj.MineFieldBaseDecayRate = getFloat[float64](o.Get("mineFieldBaseDecayRate"))
	obj.MineFieldPlanetDecayRate = getFloat[float64](o.Get("mineFieldPlanetDecayRate"))
	obj.MineFieldMaxDecayRate = getFloat[float64](o.Get("mineFieldMaxDecayRate"))
	obj.CanDetonateMineFields = getBool(o.Get("canDetonateMineFields"))
	obj.MineFieldDetonateDecayRate = getFloat[float64](o.Get("mineFieldDetonateDecayRate"))
	obj.DiscoverDesignOnScan = getBool(o.Get("discoverDesignOnScan"))
	obj.CanRemoteMineOwnPlanets = getBool(o.Get("canRemoteMineOwnPlanets"))
	obj.InvasionAttackBonus = getFloat[float64](o.Get("invasionAttackBonus"))
	obj.InvasionDefendBonus = getFloat[float64](o.Get("invasionDefendBonus"))
	obj.MovementBonus = getInt[int](o.Get("movementBonus"))
	obj.Instaforming = getBool(o.Get("instaforming"))
	obj.PermaformChance = getFloat[float64](o.Get("permaformChance"))
	obj.PermaformPopulation = getInt[int](o.Get("permaformPopulation"))
	obj.RepairFactor = getFloat[float64](o.Get("repairFactor"))
	obj.StarbaseRepairFactor = getFloat[float64](o.Get("starbaseRepairFactor"))
	obj.InnateMining = getBool(o.Get("innateMining"))
	obj.InnateResources = getBool(o.Get("innateResources"))
	obj.InnateScanner = getBool(o.Get("innateScanner"))
	obj.InnatePopulationFactor = getFloat[float64](o.Get("innatePopulationFactor"))
	obj.CanBuildDefenses = getBool(o.Get("canBuildDefenses"))
	obj.LivesOnStarbases = getBool(o.Get("livesOnStarbases"))
	obj.FuelEfficiencyOffset = getFloat[float64](o.Get("fuelEfficiencyOffset"))
	obj.TerraformCostOffset = GetCost(o.Get("terraformCostOffset"))
	obj.MineralAlchemyCostOffset = getInt[int](o.Get("mineralAlchemyCostOffset"))
	obj.ScrapMineralOffset = getFloat[float64](o.Get("scrapMineralOffset"))
	obj.ScrapMineralOffsetStarbase = getFloat[float64](o.Get("scrapMineralOffsetStarbase"))
	obj.ScrapResourcesOffset = getFloat[float64](o.Get("scrapResourcesOffset"))
	obj.ScrapResourcesOffsetStarbase = getFloat[float64](o.Get("scrapResourcesOffsetStarbase"))
	obj.StartingPopulationFactor = getFloat[float64](o.Get("startingPopulationFactor"))
	obj.StarbaseBuiltInCloakUnits = getInt[int](o.Get("starbaseBuiltInCloakUnits"))
	obj.StarbaseCostFactor = getFloat[float64](o.Get("starbaseCostFactor"))
	obj.ResearchFactor = getFloat[float64](o.Get("researchFactor"))
	obj.ResearchSplashDamage = getFloat[float64](o.Get("researchSplashDamage"))
	obj.ArmorStrengthFactor = getFloat[float64](o.Get("armorStrengthFactor"))
	obj.ShieldStrengthFactor = getFloat[float64](o.Get("shieldStrengthFactor"))
	obj.ShieldRegenerationRate = getFloat[float64](o.Get("shieldRegenerationRate"))
	obj.EngineFailureRate = getFloat[float64](o.Get("engineFailureRate"))
	obj.EngineReliableSpeed = getInt[int](o.Get("engineReliableSpeed"))
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
	obj.ID = getInt[int64](o.Get("id"))
	obj.CreatedAt = getTime(o.Get("createdAt"))
	obj.UpdatedAt = getTime(o.Get("updatedAt"))
	obj.GameID = getInt[int64](o.Get("gameId"))
	obj.TachyonCloakReduction = getInt[int](o.Get("tachyonCloakReduction"))
	obj.MaxPopulation = getInt[int](o.Get("maxPopulation"))
	obj.MinMaxPopulationPercent = getFloat[float64](o.Get("minMaxPopulationPercent"))
	obj.PopulationOvercrowdDieoffRate = getFloat[float64](o.Get("populationOvercrowdDieoffRate"))
	obj.PopulationOvercrowdDieoffRateMax = getFloat[float64](o.Get("populationOvercrowdDieoffRateMax"))
	obj.PopulationScannerError = getFloat[float64](o.Get("populationScannerError"))
	obj.SmartDefenseCoverageFactor = getFloat[float64](o.Get("smartDefenseCoverageFactor"))
	obj.InvasionDefenseCoverageFactor = getFloat[float64](o.Get("invasionDefenseCoverageFactor"))
	obj.NumBattleRounds = getInt[int](o.Get("numBattleRounds"))
	obj.MovesToRunAway = getInt[int](o.Get("movesToRunAway"))
	obj.BeamRangeDropoff = getFloat[float64](o.Get("beamRangeDropoff"))
	obj.TorpedoSplashDamage = getFloat[float64](o.Get("torpedoSplashDamage"))
	obj.SalvageDecayRate = getFloat[float64](o.Get("salvageDecayRate"))
	obj.SalvageDecayMin = getInt[int](o.Get("salvageDecayMin"))
	obj.MineFieldCloak = getInt[int](o.Get("mineFieldCloak"))
	obj.StargateMaxRangeFactor = getInt[int](o.Get("stargateMaxRangeFactor"))
	obj.StargateMaxHullMassFactor = getInt[int](o.Get("stargateMaxHullMassFactor"))
	obj.FleetSafeSpeedExplosionChance = getFloat[float64](o.Get("fleetSafeSpeedExplosionChance"))
	obj.RandomEventChances = GetMap[map[cs.RandomEvent]float64, cs.RandomEvent, float64](o.Get("randomEventChances"), GetRandomEvent, getFloat)
	obj.RandomMineralDepositBonusRange = [2]int(GetSlice[int](o.Get("randomMineralDepositBonusRange"), getInt))
	obj.RandomArtifactResearchBonusRange = [2]int(GetSlice[int](o.Get("randomArtifactResearchBonusRange"), getInt))
	obj.RandomCometMinYear = getInt[int](o.Get("randomCometMinYear"))
	obj.RandomCometMinYearPlayerWorld = getInt[int](o.Get("randomCometMinYearPlayerWorld"))
	obj.MysteryTraderRules = GetMysteryTraderRules(o.Get("mysteryTraderRules"))
	obj.CometStatsBySize = GetMap[map[cs.CometSize]cs.CometStats, cs.CometSize, cs.CometStats](o.Get("cometStatsBySize"), GetCometSize, GetCometStats)
	obj.WormholeCloak = getInt[int](o.Get("wormholeCloak"))
	obj.WormholeMinPlanetDistance = getInt[int](o.Get("wormholeMinDistance"))
	obj.WormholeStatsByStability = GetMap[map[cs.WormholeStability]cs.WormholeStats, cs.WormholeStability, cs.WormholeStats](o.Get("wormholeStatsByStability"), GetWormholeStability, GetWormholeStats)
	obj.WormholePairsForSize = GetMap[map[cs.Size]int, cs.Size, int](o.Get("wormholePairsForSize"), GetSize, getInt)
	obj.MineFieldStatsByType = GetMap[map[cs.MineFieldType]cs.MineFieldStats, cs.MineFieldType, cs.MineFieldStats](o.Get("mineFieldStatsByType"), GetMineFieldType, GetMineFieldStats)
	obj.RepairRates = GetMap[map[cs.RepairRate]float64, cs.RepairRate, float64](o.Get("repairRates"), GetRepairRate, getFloat)
	obj.MaxPlayers = getInt[int](o.Get("maxPlayers"))
	obj.StartingYear = getInt[int](o.Get("startingYear"))
	obj.ShowPublicScoresAfterYears = getInt[int](o.Get("showPublicScoresAfterYears"))
	obj.PlanetMinDistance = getInt[int](o.Get("planetMinDistance"))
	obj.MaxExtraWorldDistance = getInt[int](o.Get("maxExtraWorldDistance"))
	obj.MinExtraWorldDistance = getInt[int](o.Get("minExtraWorldDistance"))
	obj.MinHomeworldMineralConcentration = getInt[int](o.Get("minHomeworldMineralConcentration"))
	obj.MinExtraPlanetMineralConcentration = getInt[int](o.Get("minExtraPlanetMineralConcentration"))
	obj.MinHab = getInt[int](o.Get("minHab"))
	obj.MaxHab = getInt[int](o.Get("maxHab"))
	obj.MinMineralConcentration = getInt[int](o.Get("minMineralConcentration"))
	obj.MaxMineralConcentration = getInt[int](o.Get("maxMineralConcentration"))
	obj.MinStartingMineralConcentration = getInt[int](o.Get("minStartingMineralConcentration"))
	obj.MaxStartingMineralConcentration = getInt[int](o.Get("maxStartingMineralConcentration"))
	obj.HighRadMineralConcentrationBonusThreshold = getInt[int](o.Get("highRadGermaniumBonusThreshold"))
	obj.RadiatingImmune = getInt[int](o.Get("radiatingImmune"))
	obj.MaxStartingMineralSurface = getInt[int](o.Get("maxStartingMineralSurface"))
	obj.MinStartingMineralSurface = getInt[int](o.Get("minStartingMineralSurface"))
	obj.MineralDecayFactor = getInt[int](o.Get("mineralDecayFactor"))
	obj.RemoteMiningMineOutput = getInt[int](o.Get("remoteMiningMineOutput"))
	obj.StartingMines = getInt[int](o.Get("startingMines"))
	obj.StartingFactories = getInt[int](o.Get("startingFactories"))
	obj.StartingDefenses = getInt[int](o.Get("startingDefenses"))
	obj.RaceStartingPoints = getInt[int](o.Get("raceStartingPoints"))
	obj.ScrapMineralAmount = getFloat[float64](o.Get("scrapMineralAmount"))
	obj.ScrapResourceAmount = getFloat[float64](o.Get("scrapResourceAmount"))
	obj.FactoryCostGermanium = getInt[int](o.Get("factoryCostGermanium"))
	obj.DefenseCost = GetCost(o.Get("defenseCost"))
	obj.MineralAlchemyCost = getInt[int](o.Get("mineralAlchemyCost"))
	obj.PlanetaryScannerCost = GetCost(o.Get("planetaryScannerCost"))
	obj.TerraformCost = GetCost(o.Get("terraformCost"))
	obj.StarbaseComponentCostFactor = getFloat[float64](o.Get("starbaseComponentCostFactor"))
	obj.SalvageFromBattleFactor = getFloat[float64](o.Get("salvageFromBattleFactor"))
	obj.TechTradeChance = getFloat[float64](o.Get("techTradeChance"))
	obj.PacketDecayRate = GetMap[map[int]float64, int, float64](o.Get("packetDecayRate"), getInt, getFloat)
	obj.PacketMinDecay = getInt[int](o.Get("packetMinDecay"))
	obj.MaxTechLevel = getInt[int](o.Get("maxTechLevel"))
	obj.TechBaseCost = GetSlice[int](o.Get("techBaseCost"), getInt)
	obj.PRTSpecs = GetMap[map[cs.PRT]cs.PRTSpec, cs.PRT, cs.PRTSpec](o.Get("prtSpecs"), GetPRT, GetPRTSpec)
	obj.LRTSpecs = GetMap[map[cs.LRT]cs.LRTSpec, cs.LRT, cs.LRTSpec](o.Get("lrtSpecs"), GetLRT, GetLRTSpec)
	obj.TechsID = getInt[int64](o.Get("techsId"))
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
	obj.BuiltInScannerMultiplier = getInt[int](o.Get("builtInScannerMultiplier"))
	obj.NoAdvancedScanners = getBool(o.Get("noAdvancedScanners"))
	obj.ScanRangeFactor = getFloat[float64](o.Get("scanRangeFactor"))
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
	obj.Num = getInt[int](o.Get("num"))
	obj.PlayerNum = getInt[int](o.Get("playerNum"))
	obj.OriginalPlayerNum = getInt[int](o.Get("originalPlayerNum"))
	obj.Name = string(getString(o.Get("name")))
	obj.Version = getInt[int](o.Get("version"))
	obj.Hull = string(getString(o.Get("hull")))
	obj.HullSetNumber = getInt[int](o.Get("hullSetNumber"))
	obj.CannotDelete = getBool(o.Get("cannotDelete"))
	obj.MysteryTrader = getBool(o.Get("mysteryTrader"))
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
	obj.Hull = string(getString(o.Get("hull")))
	obj.HullSetNumber = getInt[int](o.Get("hullSetNumber"))
	obj.Version = getInt[int](o.Get("version"))
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
	obj.HullComponent = string(getString(o.Get("hullComponent")))
	obj.HullSlotIndex = getInt[int](o.Get("hullSlotIndex"))
	obj.Quantity = getInt[int](o.Get("quantity"))
	return obj
}

func SetShipDesignSlot(o js.Value, obj *cs.ShipDesignSlot) {
}

func GetShipDesignSpec(o js.Value) cs.ShipDesignSpec {
	var obj cs.ShipDesignSpec
	if o.IsUndefined() {
		return obj
	}
	obj.AdditionalMassDrivers = getInt[int](o.Get("additionalMassDrivers"))
	obj.Armor = getInt[int](o.Get("armor"))
	obj.BasePacketSpeed = getInt[int](o.Get("basePacketSpeed"))
	obj.BeamBonus = getFloat[float64](o.Get("beamBonus"))
	obj.BeamDefense = getFloat[float64](o.Get("beamDefense"))
	obj.Bomber = getBool(o.Get("bomber"))
	obj.Bombs = GetSlice(o.Get("bombs"), GetBomb)
	obj.CanJump = getBool(o.Get("canJump"))
	obj.CanLayMines = getBool(o.Get("canLayMines"))
	obj.CanStealFleetCargo = getBool(o.Get("canStealFleetCargo"))
	obj.CanStealPlanetCargo = getBool(o.Get("canStealPlanetCargo"))
	obj.CargoCapacity = getInt[int](o.Get("cargoCapacity"))
	obj.CloakPercent = getInt[int](o.Get("cloakPercent"))
	obj.CloakPercentFullCargo = getInt[int](o.Get("cloakPercentFullCargo"))
	obj.CloakUnits = getInt[int](o.Get("cloakUnits"))
	obj.Colonizer = getBool(o.Get("colonizer"))
	obj.Cost = GetCost(o.Get("cost"))
	obj.Engine = GetEngine(o.Get("engine"))
	obj.EstimatedRange = getInt[int](o.Get("estimatedRange"))
	obj.EstimatedRangeFull = getInt[int](o.Get("estimatedRangeFull"))
	obj.FuelCapacity = getInt[int](o.Get("fuelCapacity"))
	obj.FuelGeneration = getInt[int](o.Get("fuelGeneration"))
	obj.HasWeapons = getBool(o.Get("hasWeapons"))
	obj.HullType = GetTechHullType(o.Get("hullType"))
	obj.ImmuneToOwnDetonation = getBool(o.Get("immuneToOwnDetonation"))
	obj.Initiative = getInt[int](o.Get("initiative"))
	obj.InnateScanRangePenFactor = getFloat[float64](o.Get("innateScanRangePenFactor"))
	obj.Mass = getInt[int](o.Get("mass"))
	obj.MassDriver = string(getString(o.Get("massDriver")))
	obj.MaxHullMass = getInt[int](o.Get("maxHullMass"))
	obj.MaxPopulation = getInt[int](o.Get("maxPopulation"))
	obj.MaxRange = getInt[int](o.Get("maxRange"))
	obj.MineLayingRateByMineType = GetMap[map[cs.MineFieldType]int, cs.MineFieldType, int](o.Get("mineLayingRateByMineType"), GetMineFieldType, getInt)
	obj.MineSweep = getInt[int](o.Get("mineSweep"))
	obj.MiningRate = getInt[int](o.Get("miningRate"))
	obj.Movement = getInt[int](o.Get("movement"))
	obj.MovementBonus = getInt[int](o.Get("movementBonus"))
	obj.MovementFull = getInt[int](o.Get("movementFull"))
	obj.NumBuilt = getInt[int](o.Get("numBuilt"))
	obj.NumEngines = getInt[int](o.Get("numEngines"))
	obj.NumInstances = getInt[int](o.Get("numInstances"))
	obj.OrbitalConstructionModule = getBool(o.Get("orbitalConstructionModule"))
	obj.PowerRating = getInt[int](o.Get("powerRating"))
	obj.Radiating = getBool(o.Get("radiating"))
	obj.ReduceCloaking = getFloat[float64](o.Get("reduceCloaking"))
	obj.ReduceMovement = getInt[int](o.Get("reduceMovement"))
	obj.RepairBonus = getFloat[float64](o.Get("repairBonus"))
	obj.RetroBombs = GetSlice(o.Get("retroBombs"), GetBomb)
	obj.SafeHullMass = getInt[int](o.Get("safeHullMass"))
	obj.SafePacketSpeed = getInt[int](o.Get("safePacketSpeed"))
	obj.SafeRange = getInt[int](o.Get("safeRange"))
	obj.Scanner = getBool(o.Get("scanner"))
	obj.ScanRange = getInt[int](o.Get("scanRange"))
	obj.ScanRangePen = getInt[int](o.Get("scanRangePen"))
	obj.Shields = getInt[int](o.Get("shields"))
	obj.SmartBombs = GetSlice(o.Get("smartBombs"), GetBomb)
	obj.SpaceDock = getInt[int](o.Get("spaceDock"))
	obj.Starbase = getBool(o.Get("starbase"))
	obj.Stargate = string(getString(o.Get("stargate")))
	obj.TechLevel = GetTechLevel(o.Get("techLevel"))
	obj.TerraformRate = getInt[int](o.Get("terraformRate"))
	obj.TorpedoBonus = getFloat[float64](o.Get("torpedoBonus"))
	obj.TorpedoJamming = getFloat[float64](o.Get("torpedoJamming"))
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
	obj.DesignNum = getInt[int](o.Get("designNum"))
	obj.Quantity = getInt[int](o.Get("quantity"))
	obj.Damage = getFloat[float64](o.Get("damage"))
	obj.QuantityDamaged = getInt[int](o.Get("quantityDamaged"))
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
	obj.Name = string(getString(o.Get("name")))
	obj.HullName = GetStartingFleetHull(o.Get("hullName"))
	obj.HullSetNumber = getInt[uint](o.Get("hullSetNumber"))
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
	obj.Population = getInt[int](o.Get("population"))
	obj.HabPenaltyFactor = getFloat[float64](o.Get("habPenaltyFactor"))
	obj.HasStargate = getBool(o.Get("hasStargate"))
	obj.HasMassDriver = getBool(o.Get("hasMassDriver"))
	obj.StarbaseDesignName = string(getString(o.Get("starbaseDesignName")))
	obj.StarbaseHull = string(getString(o.Get("starbaseHull")))
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
	obj.Energy = getFloat[float64](o.Get("energy"))
	obj.Weapons = getFloat[float64](o.Get("weapons"))
	obj.Propulsion = getFloat[float64](o.Get("propulsion"))
	obj.Construction = getFloat[float64](o.Get("construction"))
	obj.Electronics = getFloat[float64](o.Get("electronics"))
	obj.Biotechnology = getFloat[float64](o.Get("biotechnology"))
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
	obj.Name = string(getString(o.Get("name")))
	obj.Cost = GetCost(o.Get("cost"))
	obj.Requirements = GetTechRequirements(o.Get("requirements"))
	obj.Ranking = getInt[int](o.Get("ranking"))
	obj.Category = GetTechCategory(o.Get("category"))
	obj.Origin = string(getString(o.Get("origin")))
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
	obj.Engine = getFloat[float64](o.Get("engine"))
	obj.BeamWeapon = getFloat[float64](o.Get("beamWeapon"))
	obj.Torpedo = getFloat[float64](o.Get("torpedo"))
	obj.Bomb = getFloat[float64](o.Get("bomb"))
	obj.PlanetaryDefense = getFloat[float64](o.Get("planetaryDefense"))
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
	obj.Energy = getInt[int](o.Get("energy"))
	obj.Weapons = getInt[int](o.Get("weapons"))
	obj.Propulsion = getInt[int](o.Get("propulsion"))
	obj.Construction = getInt[int](o.Get("construction"))
	obj.Electronics = getInt[int](o.Get("electronics"))
	obj.Biotechnology = getInt[int](o.Get("biotechnology"))
	return obj
}

func SetTechLevel(o js.Value, obj *cs.TechLevel) {
}

func GetTechPlanetary(o js.Value) cs.TechPlanetary {
	var obj cs.TechPlanetary
	if o.IsUndefined() {
		return obj
	}
	obj.ResetPlanet = getBool(o.Get("resetPlanet"))
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
	obj.ScanRange = getInt[int](o.Get("scanRange"))
	obj.ScanRangePen = getInt[int](o.Get("scanRangePen"))
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
	obj.Acquirable = getBool(o.Get("acquirable"))
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
	obj.Ability = getInt[int](o.Get("ability"))
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
	obj.Num = getInt[int](o.Get("num"))
	obj.Name = string(getString(o.Get("name")))
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
	obj.X = getFloat[float64](o.Get("x"))
	obj.Y = getFloat[float64](o.Get("y"))
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
	obj.WarpSpeed = getInt[int](o.Get("warpSpeed"))
	obj.EstFuelUsage = getInt[int](o.Get("estFuelUsage"))
	obj.Task = GetWaypointTask(o.Get("task"))
	obj.TransportTasks = GetWaypointTransportTasks(o.Get("transportTasks"))
	obj.WaitAtWaypoint = getBool(o.Get("waitAtWaypoint"))
	obj.LayMineFieldDuration = getInt[int](o.Get("layMineFieldDuration"))
	obj.PatrolRange = getInt[int](o.Get("patrolRange"))
	obj.PatrolWarpSpeed = getInt[int](o.Get("patrolWarpSpeed"))
	obj.TargetType = GetMapObjectType(o.Get("targetType"))
	obj.TargetNum = getInt[int](o.Get("targetNum"))
	obj.TargetPlayerNum = getInt[int](o.Get("targetPlayerNum"))
	obj.TargetName = string(getString(o.Get("targetName")))
	obj.TransferToPlayer = getInt[int](o.Get("transferToPlayer"))
	obj.PartiallyComplete = getBool(o.Get("partiallyComplete"))
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
	obj.Amount = getInt[int](o.Get("amount"))
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
	obj.DestinationNum = getInt[int](o.Get("destinationNum"))
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
	obj.YearsToDegrade = getInt[int](o.Get("yearsToDegrade"))
	obj.ChanceToJump = getFloat[float64](o.Get("chanceToJump"))
	obj.JiggleDistance = getInt[int](o.Get("jiggleDistance"))
	return obj
}

func SetWormholeStats(o js.Value, obj *cs.WormholeStats) {
}
