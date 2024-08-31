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
	// Num num BasicInt ignore: false
	obj.Num = GetInt[int](o, "num")
	// Name name BasicString ignore: false
	obj.Name = string(GetString(o, "name"))
	// PrimaryTarget primaryTarget Named ignore: false
	obj.PrimaryTarget = cs.BattleTarget(GetString(o, "primaryTarget"))
	// SecondaryTarget secondaryTarget Named ignore: false
	obj.SecondaryTarget = cs.BattleTarget(GetString(o, "secondaryTarget"))
	// Tactic tactic Named ignore: false
	obj.Tactic = cs.BattleTactic(GetString(o, "tactic"))
	// AttackWho attackWho Named ignore: false
	obj.AttackWho = cs.BattleAttackWho(GetString(o, "attackWho"))
	// DumpCargo dumpCargo BasicBool ignore: false
	obj.DumpCargo = bool(GetBool(o, "dumpCargo"))
	return obj
}

func SetBattlePlan(o js.Value, obj *cs.BattlePlan) {
	// Num num BasicInt ignore: false
	o.Set("num", obj.Num)
	// Name name BasicString ignore: false
	o.Set("name", obj.Name)
	// PrimaryTarget primaryTarget Named ignore: false
	o.Set("primaryTarget", string(obj.PrimaryTarget))
	// SecondaryTarget secondaryTarget Named ignore: false
	o.Set("secondaryTarget", string(obj.SecondaryTarget))
	// Tactic tactic Named ignore: false
	o.Set("tactic", string(obj.Tactic))
	// AttackWho attackWho Named ignore: false
	o.Set("attackWho", string(obj.AttackWho))
	// DumpCargo dumpCargo BasicBool ignore: false
	o.Set("dumpCargo", obj.DumpCargo)
}

func GetBattleRecord(o js.Value) cs.BattleRecord {
	obj := cs.BattleRecord{}
	// Num num BasicInt ignore: false
	obj.Num = GetInt[int](o, "num")
	// PlanetNum planetNum BasicInt ignore: false
	obj.PlanetNum = GetInt[int](o, "planetNum")
	// Position position Object ignore: false
	obj.Position = GetVector(o.Get("position"))
	// Tokens tokens Slice ignore: false
	obj.Tokens = GetSlice(o.Get("tokens"), GetBattleRecordToken)
	// ActionsPerRound actionsPerRound Slice ignore: false
	obj.ActionsPerRound = GetSliceSlice(o.Get("actionsPerRound"), GetBattleRecordTokenAction)
	// DestroyedTokens destroyedTokens Slice ignore: false
	obj.DestroyedTokens = GetSlice(o.Get("destroyedTokens"), GetBattleRecordDestroyedToken)
	// Stats stats Object ignore: false
	obj.Stats = GetBattleRecordStats(o.Get("stats"))
	return obj
}

func SetBattleRecord(o js.Value, obj *cs.BattleRecord) {
	// Num num BasicInt ignore: false
	o.Set("num", obj.Num)
	// PlanetNum planetNum BasicInt ignore: false
	o.Set("planetNum", obj.PlanetNum)
	// Position position Object ignore: false
	o.Set("position", map[string]any{})
	SetVector(o.Get("position"), &obj.Position)
	// Tokens tokens Slice ignore: false
	o.Set("tokens", []any{})
	SetSlice(o.Get("tokens"), obj.Tokens, SetBattleRecordToken)
	// ActionsPerRound actionsPerRound Slice ignore: false
	o.Set("actionsPerRound", []any{})
	SetSliceSlice(o.Get("actionsPerRound"), obj.ActionsPerRound, SetBattleRecordTokenAction)
	// DestroyedTokens destroyedTokens Slice ignore: false
	o.Set("destroyedTokens", []any{})
	SetSlice(o.Get("destroyedTokens"), obj.DestroyedTokens, SetBattleRecordDestroyedToken)
	// Stats stats Object ignore: false
	o.Set("stats", map[string]any{})
	SetBattleRecordStats(o.Get("stats"), &obj.Stats)
}

func GetBattleRecordDestroyedToken(o js.Value) cs.BattleRecordDestroyedToken {
	obj := cs.BattleRecordDestroyedToken{}
	// Num num BasicInt ignore: false
	obj.Num = GetInt[int](o, "num")
	// PlayerNum playerNum BasicInt ignore: false
	obj.PlayerNum = GetInt[int](o, "playerNum")
	// DesignNum designNum BasicInt ignore: false
	obj.DesignNum = GetInt[int](o, "designNum")
	// Quantity quantity BasicInt ignore: false
	obj.Quantity = GetInt[int](o, "quantity")
	// design   ignore: true
	return obj
}

func SetBattleRecordDestroyedToken(o js.Value, obj *cs.BattleRecordDestroyedToken) {
	// Num num BasicInt ignore: false
	o.Set("num", obj.Num)
	// PlayerNum playerNum BasicInt ignore: false
	o.Set("playerNum", obj.PlayerNum)
	// DesignNum designNum BasicInt ignore: false
	o.Set("designNum", obj.DesignNum)
	// Quantity quantity BasicInt ignore: false
	o.Set("quantity", obj.Quantity)
	// design   ignore: true
}

func GetBattleRecordStats(o js.Value) cs.BattleRecordStats {
	obj := cs.BattleRecordStats{}
	// NumPlayers numPlayers BasicInt ignore: false
	obj.NumPlayers = GetInt[int](o, "numPlayers")
	// NumShipsByPlayer numShipsByPlayer Map ignore: false

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
	// ShipsDestroyedByPlayer shipsDestroyedByPlayer Map ignore: false

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
	// DamageTakenByPlayer damageTakenByPlayer Map ignore: false

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
	// CargoLostByPlayer cargoLostByPlayer Map ignore: true
	return obj
}

func SetBattleRecordStats(o js.Value, obj *cs.BattleRecordStats) {
	// NumPlayers numPlayers BasicInt ignore: false
	o.Set("numPlayers", obj.NumPlayers)
	// NumShipsByPlayer numShipsByPlayer Map ignore: false
	numShipsByPlayerMap := js.ValueOf(map[string]any{})
	for key, value := range obj.NumShipsByPlayer {
		numShipsByPlayerMap.Set(fmt.Sprintf("%v", key), int(value))
	}
	o.Set("numShipsByPlayer", numShipsByPlayerMap)
	// ShipsDestroyedByPlayer shipsDestroyedByPlayer Map ignore: false
	shipsDestroyedByPlayerMap := js.ValueOf(map[string]any{})
	for key, value := range obj.ShipsDestroyedByPlayer {
		shipsDestroyedByPlayerMap.Set(fmt.Sprintf("%v", key), int(value))
	}
	o.Set("shipsDestroyedByPlayer", shipsDestroyedByPlayerMap)
	// DamageTakenByPlayer damageTakenByPlayer Map ignore: false
	damageTakenByPlayerMap := js.ValueOf(map[string]any{})
	for key, value := range obj.DamageTakenByPlayer {
		damageTakenByPlayerMap.Set(fmt.Sprintf("%v", key), int(value))
	}
	o.Set("damageTakenByPlayer", damageTakenByPlayerMap)
	// CargoLostByPlayer cargoLostByPlayer Map ignore: true
}

func GetBattleRecordToken(o js.Value) cs.BattleRecordToken {
	obj := cs.BattleRecordToken{}
	// Num num BasicInt ignore: false
	obj.Num = GetInt[int](o, "num")
	// PlayerNum playerNum BasicInt ignore: false
	obj.PlayerNum = GetInt[int](o, "playerNum")
	// DesignNum designNum BasicInt ignore: false
	obj.DesignNum = GetInt[int](o, "designNum")
	// Position position Object ignore: false
	obj.Position = GetBattleVector(o.Get("position"))
	// Initiative initiative BasicInt ignore: false
	obj.Initiative = GetInt[int](o, "initiative")
	// Mass mass BasicInt ignore: false
	obj.Mass = GetInt[int](o, "mass")
	// Armor armor BasicInt ignore: false
	obj.Armor = GetInt[int](o, "armor")
	// StackShields stackShields BasicInt ignore: false
	obj.StackShields = GetInt[int](o, "stackShields")
	// Movement movement BasicInt ignore: false
	obj.Movement = GetInt[int](o, "movement")
	// StartingQuantity startingQuantity BasicInt ignore: false
	obj.StartingQuantity = GetInt[int](o, "startingQuantity")
	// Tactic tactic Named ignore: false
	obj.Tactic = cs.BattleTactic(GetString(o, "tactic"))
	// PrimaryTarget primaryTarget Named ignore: false
	obj.PrimaryTarget = cs.BattleTarget(GetString(o, "primaryTarget"))
	// SecondaryTarget secondaryTarget Named ignore: false
	obj.SecondaryTarget = cs.BattleTarget(GetString(o, "secondaryTarget"))
	// AttackWho attackWho Named ignore: false
	obj.AttackWho = cs.BattleAttackWho(GetString(o, "attackWho"))
	return obj
}

func SetBattleRecordToken(o js.Value, obj *cs.BattleRecordToken) {
	// Num num BasicInt ignore: false
	o.Set("num", obj.Num)
	// PlayerNum playerNum BasicInt ignore: false
	o.Set("playerNum", obj.PlayerNum)
	// DesignNum designNum BasicInt ignore: false
	o.Set("designNum", obj.DesignNum)
	// Position position Object ignore: false
	o.Set("position", map[string]any{})
	SetBattleVector(o.Get("position"), &obj.Position)
	// Initiative initiative BasicInt ignore: false
	o.Set("initiative", obj.Initiative)
	// Mass mass BasicInt ignore: false
	o.Set("mass", obj.Mass)
	// Armor armor BasicInt ignore: false
	o.Set("armor", obj.Armor)
	// StackShields stackShields BasicInt ignore: false
	o.Set("stackShields", obj.StackShields)
	// Movement movement BasicInt ignore: false
	o.Set("movement", obj.Movement)
	// StartingQuantity startingQuantity BasicInt ignore: false
	o.Set("startingQuantity", obj.StartingQuantity)
	// Tactic tactic Named ignore: false
	o.Set("tactic", string(obj.Tactic))
	// PrimaryTarget primaryTarget Named ignore: false
	o.Set("primaryTarget", string(obj.PrimaryTarget))
	// SecondaryTarget secondaryTarget Named ignore: false
	o.Set("secondaryTarget", string(obj.SecondaryTarget))
	// AttackWho attackWho Named ignore: false
	o.Set("attackWho", string(obj.AttackWho))
}

func GetBattleRecordTokenAction(o js.Value) cs.BattleRecordTokenAction {
	obj := cs.BattleRecordTokenAction{}
	// Type type Named ignore: false
	obj.Type = cs.BattleRecordTokenActionType(GetInt[cs.BattleRecordTokenActionType](o, "type"))
	// TokenNum tokenNum BasicInt ignore: false
	obj.TokenNum = GetInt[int](o, "tokenNum")
	// Round round BasicInt ignore: false
	obj.Round = GetInt[int](o, "round")
	// From from Object ignore: false
	obj.From = GetBattleVector(o.Get("from"))
	// To to Object ignore: false
	obj.To = GetBattleVector(o.Get("to"))
	// Slot slot BasicInt ignore: false
	obj.Slot = GetInt[int](o, "slot")
	// TargetNum targetNum BasicInt ignore: false
	obj.TargetNum = GetInt[int](o, "targetNum")
	// Target target Object ignore: false
	targetVal := o.Get("target")
	if !targetVal.IsUndefined() {
		target := GetShipToken(targetVal)
		obj.Target = &target
	}
	// TokensDestroyed tokensDestroyed BasicInt ignore: false
	obj.TokensDestroyed = GetInt[int](o, "tokensDestroyed")
	// DamageDoneShields damageDoneShields BasicInt ignore: false
	obj.DamageDoneShields = GetInt[int](o, "damageDoneShields")
	// DamageDoneArmor damageDoneArmor BasicInt ignore: false
	obj.DamageDoneArmor = GetInt[int](o, "damageDoneArmor")
	// TorpedoHits torpedoHits BasicInt ignore: false
	obj.TorpedoHits = GetInt[int](o, "torpedoHits")
	// TorpedoMisses torpedoMisses BasicInt ignore: false
	obj.TorpedoMisses = GetInt[int](o, "torpedoMisses")
	return obj
}

func SetBattleRecordTokenAction(o js.Value, obj *cs.BattleRecordTokenAction) {
	// Type type Named ignore: false
	o.Set("type", int(obj.Type))
	// TokenNum tokenNum BasicInt ignore: false
	o.Set("tokenNum", obj.TokenNum)
	// Round round BasicInt ignore: false
	o.Set("round", obj.Round)
	// From from Object ignore: false
	o.Set("from", map[string]any{})
	SetBattleVector(o.Get("from"), &obj.From)
	// To to Object ignore: false
	o.Set("to", map[string]any{})
	SetBattleVector(o.Get("to"), &obj.To)
	// Slot slot BasicInt ignore: false
	o.Set("slot", obj.Slot)
	// TargetNum targetNum BasicInt ignore: false
	o.Set("targetNum", obj.TargetNum)
	// Target target Object ignore: false
	o.Set("target", map[string]any{})
	SetShipToken(o.Get("target"), obj.Target)
	// TokensDestroyed tokensDestroyed BasicInt ignore: false
	o.Set("tokensDestroyed", obj.TokensDestroyed)
	// DamageDoneShields damageDoneShields BasicInt ignore: false
	o.Set("damageDoneShields", obj.DamageDoneShields)
	// DamageDoneArmor damageDoneArmor BasicInt ignore: false
	o.Set("damageDoneArmor", obj.DamageDoneArmor)
	// TorpedoHits torpedoHits BasicInt ignore: false
	o.Set("torpedoHits", obj.TorpedoHits)
	// TorpedoMisses torpedoMisses BasicInt ignore: false
	o.Set("torpedoMisses", obj.TorpedoMisses)
}

func GetBattleVector(o js.Value) cs.BattleVector {
	obj := cs.BattleVector{}
	// X x BasicInt ignore: false
	obj.X = GetInt[int](o, "x")
	// Y y BasicInt ignore: false
	obj.Y = GetInt[int](o, "y")
	return obj
}

func SetBattleVector(o js.Value, obj *cs.BattleVector) {
	// X x BasicInt ignore: false
	o.Set("x", obj.X)
	// Y y BasicInt ignore: false
	o.Set("y", obj.Y)
}

func GetBomb(o js.Value) cs.Bomb {
	obj := cs.Bomb{}
	// Quantity quantity BasicInt ignore: false
	obj.Quantity = GetInt[int](o, "quantity")
	// KillRate killRate BasicFloat ignore: false
	obj.KillRate = GetFloat[float64](o, "killRate")
	// MinKillRate minKillRate BasicInt ignore: false
	obj.MinKillRate = GetInt[int](o, "minKillRate")
	// StructureDestroyRate structureDestroyRate BasicFloat ignore: false
	obj.StructureDestroyRate = GetFloat[float64](o, "structureDestroyRate")
	// UnterraformRate unterraformRate BasicInt ignore: false
	obj.UnterraformRate = GetInt[int](o, "unterraformRate")
	return obj
}

func SetBomb(o js.Value, obj *cs.Bomb) {
	// Quantity quantity BasicInt ignore: false
	o.Set("quantity", obj.Quantity)
	// KillRate killRate BasicFloat ignore: false
	o.Set("killRate", obj.KillRate)
	// MinKillRate minKillRate BasicInt ignore: false
	o.Set("minKillRate", obj.MinKillRate)
	// StructureDestroyRate structureDestroyRate BasicFloat ignore: false
	o.Set("structureDestroyRate", obj.StructureDestroyRate)
	// UnterraformRate unterraformRate BasicInt ignore: false
	o.Set("unterraformRate", obj.UnterraformRate)
}

func GetBombingResult(o js.Value) cs.BombingResult {
	obj := cs.BombingResult{}
	// BomberName bomberName BasicString ignore: false
	obj.BomberName = string(GetString(o, "bomberName"))
	// NumBombers numBombers BasicInt ignore: false
	obj.NumBombers = GetInt[int](o, "numBombers")
	// ColonistsKilled colonistsKilled BasicInt ignore: false
	obj.ColonistsKilled = GetInt[int](o, "colonistsKilled")
	// MinesDestroyed minesDestroyed BasicInt ignore: false
	obj.MinesDestroyed = GetInt[int](o, "minesDestroyed")
	// FactoriesDestroyed factoriesDestroyed BasicInt ignore: false
	obj.FactoriesDestroyed = GetInt[int](o, "factoriesDestroyed")
	// DefensesDestroyed defensesDestroyed BasicInt ignore: false
	obj.DefensesDestroyed = GetInt[int](o, "defensesDestroyed")
	// UnterraformAmount unterraformAmount Object ignore: false
	obj.UnterraformAmount = GetHab(o.Get("unterraformAmount"))
	// PlanetEmptied planetEmptied BasicBool ignore: false
	obj.PlanetEmptied = bool(GetBool(o, "planetEmptied"))
	// fleet   ignore: true
	return obj
}

func SetBombingResult(o js.Value, obj *cs.BombingResult) {
	// BomberName bomberName BasicString ignore: false
	o.Set("bomberName", obj.BomberName)
	// NumBombers numBombers BasicInt ignore: false
	o.Set("numBombers", obj.NumBombers)
	// ColonistsKilled colonistsKilled BasicInt ignore: false
	o.Set("colonistsKilled", obj.ColonistsKilled)
	// MinesDestroyed minesDestroyed BasicInt ignore: false
	o.Set("minesDestroyed", obj.MinesDestroyed)
	// FactoriesDestroyed factoriesDestroyed BasicInt ignore: false
	o.Set("factoriesDestroyed", obj.FactoriesDestroyed)
	// DefensesDestroyed defensesDestroyed BasicInt ignore: false
	o.Set("defensesDestroyed", obj.DefensesDestroyed)
	// UnterraformAmount unterraformAmount Object ignore: false
	o.Set("unterraformAmount", map[string]any{})
	SetHab(o.Get("unterraformAmount"), &obj.UnterraformAmount)
	// PlanetEmptied planetEmptied BasicBool ignore: false
	o.Set("planetEmptied", obj.PlanetEmptied)
	// fleet   ignore: true
}

func GetCargo(o js.Value) cs.Cargo {
	obj := cs.Cargo{}
	// Ironium ironium BasicInt ignore: false
	obj.Ironium = GetInt[int](o, "ironium")
	// Boranium boranium BasicInt ignore: false
	obj.Boranium = GetInt[int](o, "boranium")
	// Germanium germanium BasicInt ignore: false
	obj.Germanium = GetInt[int](o, "germanium")
	// Colonists colonists BasicInt ignore: false
	obj.Colonists = GetInt[int](o, "colonists")
	return obj
}

func SetCargo(o js.Value, obj *cs.Cargo) {
	// Ironium ironium BasicInt ignore: false
	o.Set("ironium", obj.Ironium)
	// Boranium boranium BasicInt ignore: false
	o.Set("boranium", obj.Boranium)
	// Germanium germanium BasicInt ignore: false
	o.Set("germanium", obj.Germanium)
	// Colonists colonists BasicInt ignore: false
	o.Set("colonists", obj.Colonists)
}

func GetCost(o js.Value) cs.Cost {
	obj := cs.Cost{}
	// Ironium ironium BasicInt ignore: false
	obj.Ironium = GetInt[int](o, "ironium")
	// Boranium boranium BasicInt ignore: false
	obj.Boranium = GetInt[int](o, "boranium")
	// Germanium germanium BasicInt ignore: false
	obj.Germanium = GetInt[int](o, "germanium")
	// Resources resources BasicInt ignore: false
	obj.Resources = GetInt[int](o, "resources")
	return obj
}

func SetCost(o js.Value, obj *cs.Cost) {
	// Ironium ironium BasicInt ignore: false
	o.Set("ironium", obj.Ironium)
	// Boranium boranium BasicInt ignore: false
	o.Set("boranium", obj.Boranium)
	// Germanium germanium BasicInt ignore: false
	o.Set("germanium", obj.Germanium)
	// Resources resources BasicInt ignore: false
	o.Set("resources", obj.Resources)
}

func GetDBObject(o js.Value) cs.DBObject {
	obj := cs.DBObject{}
	// ID id BasicInt ignore: false
	obj.ID = GetInt[int64](o, "id")
	// CreatedAt createdAt  ignore: false
	obj.CreatedAt, _ = GetTime(o, "createdAt")
	// UpdatedAt updatedAt  ignore: false
	obj.UpdatedAt, _ = GetTime(o, "updatedAt")
	return obj
}

func SetDBObject(o js.Value, obj *cs.DBObject) {
	// ID id BasicInt ignore: false
	o.Set("id", obj.ID)
	// CreatedAt createdAt  ignore: false
	SetTime(o, "createdAt", obj.CreatedAt)
	// UpdatedAt updatedAt  ignore: false
	SetTime(o, "updatedAt", obj.UpdatedAt)
}

func GetDefense(o js.Value) cs.Defense {
	obj := cs.Defense{}
	// DefenseCoverage defenseCoverage BasicFloat ignore: false
	obj.DefenseCoverage = GetFloat[float64](o, "defenseCoverage")
	return obj
}

func SetDefense(o js.Value, obj *cs.Defense) {
	// DefenseCoverage defenseCoverage BasicFloat ignore: false
	o.Set("defenseCoverage", obj.DefenseCoverage)
}

func GetEngine(o js.Value) cs.Engine {
	obj := cs.Engine{}
	// IdealSpeed idealSpeed BasicInt ignore: false
	obj.IdealSpeed = GetInt[int](o, "idealSpeed")
	// FreeSpeed freeSpeed BasicInt ignore: false
	obj.FreeSpeed = GetInt[int](o, "freeSpeed")
	// MaxSafeSpeed maxSafeSpeed BasicInt ignore: false
	obj.MaxSafeSpeed = GetInt[int](o, "maxSafeSpeed")
	// FuelUsage fuelUsage Array ignore: false
	if !o.Get("fuelUsage").IsUndefined() && o.Get("fuelUsage").Length() != 0 {
		obj.FuelUsage = [11]int(GetIntArray[int](o, "fuelUsage"))
	}
	return obj
}

func SetEngine(o js.Value, obj *cs.Engine) {
	// IdealSpeed idealSpeed BasicInt ignore: false
	o.Set("idealSpeed", obj.IdealSpeed)
	// FreeSpeed freeSpeed BasicInt ignore: false
	o.Set("freeSpeed", obj.FreeSpeed)
	// MaxSafeSpeed maxSafeSpeed BasicInt ignore: false
	o.Set("maxSafeSpeed", obj.MaxSafeSpeed)
	// FuelUsage fuelUsage Array ignore: false
}

func GetFleet(o js.Value) cs.Fleet {
	obj := cs.Fleet{}
	// MapObject  Object ignore: false
	obj.MapObject = GetMapObject(o)
	// FleetOrders  Object ignore: false
	obj.FleetOrders = GetFleetOrders(o)
	// PlanetNum planetNum BasicInt ignore: false
	obj.PlanetNum = GetInt[int](o, "planetNum")
	// BaseName baseName BasicString ignore: false
	obj.BaseName = string(GetString(o, "baseName"))
	// Cargo cargo Object ignore: false
	obj.Cargo = GetCargo(o.Get("cargo"))
	// Fuel fuel BasicInt ignore: false
	obj.Fuel = GetInt[int](o, "fuel")
	// Age age BasicInt ignore: false
	obj.Age = GetInt[int](o, "age")
	// Tokens tokens Slice ignore: false
	obj.Tokens = GetSlice(o.Get("tokens"), GetShipToken)
	// Heading heading Object ignore: false
	obj.Heading = GetVector(o.Get("heading"))
	// WarpSpeed warpSpeed BasicInt ignore: false
	obj.WarpSpeed = GetInt[int](o, "warpSpeed")
	// PreviousPosition previousPosition Object ignore: false
	previousPositionVal := o.Get("previousPosition")
	if !previousPositionVal.IsUndefined() {
		previousPosition := GetVector(previousPositionVal)
		obj.PreviousPosition = &previousPosition
	}
	// OrbitingPlanetNum orbitingPlanetNum BasicInt ignore: false
	obj.OrbitingPlanetNum = GetInt[int](o, "orbitingPlanetNum")
	// Starbase starbase BasicBool ignore: false
	obj.Starbase = bool(GetBool(o, "starbase"))
	// Spec spec Object ignore: false
	obj.Spec = GetFleetSpec(o.Get("spec"))
	// battlePlan   ignore: true
	// struckMineField   ignore: true
	// remoteMined   ignore: true
	return obj
}

func SetFleet(o js.Value, obj *cs.Fleet) {
	// MapObject  Object ignore: false
	SetMapObject(o, &obj.MapObject)
	// FleetOrders  Object ignore: false
	SetFleetOrders(o, &obj.FleetOrders)
	// PlanetNum planetNum BasicInt ignore: false
	o.Set("planetNum", obj.PlanetNum)
	// BaseName baseName BasicString ignore: false
	o.Set("baseName", obj.BaseName)
	// Cargo cargo Object ignore: false
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), &obj.Cargo)
	// Fuel fuel BasicInt ignore: false
	o.Set("fuel", obj.Fuel)
	// Age age BasicInt ignore: false
	o.Set("age", obj.Age)
	// Tokens tokens Slice ignore: false
	o.Set("tokens", []any{})
	SetSlice(o.Get("tokens"), obj.Tokens, SetShipToken)
	// Heading heading Object ignore: false
	o.Set("heading", map[string]any{})
	SetVector(o.Get("heading"), &obj.Heading)
	// WarpSpeed warpSpeed BasicInt ignore: false
	o.Set("warpSpeed", obj.WarpSpeed)
	// PreviousPosition previousPosition Object ignore: false
	o.Set("previousPosition", map[string]any{})
	SetVector(o.Get("previousPosition"), obj.PreviousPosition)
	// OrbitingPlanetNum orbitingPlanetNum BasicInt ignore: false
	o.Set("orbitingPlanetNum", obj.OrbitingPlanetNum)
	// Starbase starbase BasicBool ignore: false
	o.Set("starbase", obj.Starbase)
	// Spec spec Object ignore: false
	o.Set("spec", map[string]any{})
	SetFleetSpec(o.Get("spec"), &obj.Spec)
	// battlePlan   ignore: true
	// struckMineField   ignore: true
	// remoteMined   ignore: true
}

func GetFleetIntel(o js.Value) cs.FleetIntel {
	obj := cs.FleetIntel{}
	// MapObjectIntel  Object ignore: false
	obj.MapObjectIntel = GetMapObjectIntel(o)
	// PlanetIntelID  BasicInt ignore: false
	obj.PlanetIntelID = GetInt[int64](o, "")
	// BaseName baseName BasicString ignore: false
	obj.BaseName = string(GetString(o, "baseName"))
	// Heading heading Object ignore: false
	obj.Heading = GetVector(o.Get("heading"))
	// OrbitingPlanetNum orbitingPlanetNum BasicInt ignore: false
	obj.OrbitingPlanetNum = GetInt[int](o, "orbitingPlanetNum")
	// WarpSpeed warpSpeed BasicInt ignore: false
	obj.WarpSpeed = GetInt[int](o, "warpSpeed")
	// Mass mass BasicInt ignore: false
	obj.Mass = GetInt[int](o, "mass")
	// Cargo cargo Object ignore: false
	obj.Cargo = GetCargo(o.Get("cargo"))
	// CargoDiscovered cargoDiscovered BasicBool ignore: false
	obj.CargoDiscovered = bool(GetBool(o, "cargoDiscovered"))
	// Freighter freighter BasicBool ignore: false
	obj.Freighter = bool(GetBool(o, "freighter"))
	// ScanRange scanRange BasicInt ignore: false
	obj.ScanRange = GetInt[int](o, "scanRange")
	// ScanRangePen scanRangePen BasicInt ignore: false
	obj.ScanRangePen = GetInt[int](o, "scanRangePen")
	// Tokens tokens Slice ignore: false
	obj.Tokens = GetSlice(o.Get("tokens"), GetShipToken)
	return obj
}

func SetFleetIntel(o js.Value, obj *cs.FleetIntel) {
	// MapObjectIntel  Object ignore: false
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	// PlanetIntelID  BasicInt ignore: false
	o.Set("", obj.PlanetIntelID)
	// BaseName baseName BasicString ignore: false
	o.Set("baseName", obj.BaseName)
	// Heading heading Object ignore: false
	o.Set("heading", map[string]any{})
	SetVector(o.Get("heading"), &obj.Heading)
	// OrbitingPlanetNum orbitingPlanetNum BasicInt ignore: false
	o.Set("orbitingPlanetNum", obj.OrbitingPlanetNum)
	// WarpSpeed warpSpeed BasicInt ignore: false
	o.Set("warpSpeed", obj.WarpSpeed)
	// Mass mass BasicInt ignore: false
	o.Set("mass", obj.Mass)
	// Cargo cargo Object ignore: false
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), &obj.Cargo)
	// CargoDiscovered cargoDiscovered BasicBool ignore: false
	o.Set("cargoDiscovered", obj.CargoDiscovered)
	// Freighter freighter BasicBool ignore: false
	o.Set("freighter", obj.Freighter)
	// ScanRange scanRange BasicInt ignore: false
	o.Set("scanRange", obj.ScanRange)
	// ScanRangePen scanRangePen BasicInt ignore: false
	o.Set("scanRangePen", obj.ScanRangePen)
	// Tokens tokens Slice ignore: false
	o.Set("tokens", []any{})
	SetSlice(o.Get("tokens"), obj.Tokens, SetShipToken)
}

func GetFleetOrders(o js.Value) cs.FleetOrders {
	obj := cs.FleetOrders{}
	// Waypoints waypoints Slice ignore: false
	obj.Waypoints = GetSlice(o.Get("waypoints"), GetWaypoint)
	// RepeatOrders repeatOrders BasicBool ignore: false
	obj.RepeatOrders = bool(GetBool(o, "repeatOrders"))
	// BattlePlanNum battlePlanNum BasicInt ignore: false
	obj.BattlePlanNum = GetInt[int](o, "battlePlanNum")
	// Purpose purpose Named ignore: false
	obj.Purpose = cs.FleetPurpose(GetString(o, "purpose"))
	return obj
}

func SetFleetOrders(o js.Value, obj *cs.FleetOrders) {
	// Waypoints waypoints Slice ignore: false
	o.Set("waypoints", []any{})
	SetSlice(o.Get("waypoints"), obj.Waypoints, SetWaypoint)
	// RepeatOrders repeatOrders BasicBool ignore: false
	o.Set("repeatOrders", obj.RepeatOrders)
	// BattlePlanNum battlePlanNum BasicInt ignore: false
	o.Set("battlePlanNum", obj.BattlePlanNum)
	// Purpose purpose Named ignore: false
	o.Set("purpose", string(obj.Purpose))
}

func GetFleetSpec(o js.Value) cs.FleetSpec {
	obj := cs.FleetSpec{}
	// ShipDesignSpec  Object ignore: false
	obj.ShipDesignSpec = GetShipDesignSpec(o)
	// BaseCloakedCargo baseCloakedCargo BasicInt ignore: false
	obj.BaseCloakedCargo = GetInt[int](o, "baseCloakedCargo")
	// BasePacketSpeed basePacketSpeed BasicInt ignore: false
	obj.BasePacketSpeed = GetInt[int](o, "basePacketSpeed")
	// HasMassDriver hasMassDriver BasicBool ignore: false
	obj.HasMassDriver = bool(GetBool(o, "hasMassDriver"))
	// HasStargate hasStargate BasicBool ignore: false
	obj.HasStargate = bool(GetBool(o, "hasStargate"))
	// MassDriver massDriver BasicString ignore: false
	obj.MassDriver = string(GetString(o, "massDriver"))
	// MassEmpty massEmpty BasicInt ignore: false
	obj.MassEmpty = GetInt[int](o, "massEmpty")
	// MaxHullMass maxHullMass BasicInt ignore: false
	obj.MaxHullMass = GetInt[int](o, "maxHullMass")
	// MaxRange maxRange BasicInt ignore: false
	obj.MaxRange = GetInt[int](o, "maxRange")
	// Purposes purposes Map ignore: true
	// SafeHullMass safeHullMass BasicInt ignore: false
	obj.SafeHullMass = GetInt[int](o, "safeHullMass")
	// SafeRange safeRange BasicInt ignore: false
	obj.SafeRange = GetInt[int](o, "safeRange")
	// Stargate stargate BasicString ignore: false
	obj.Stargate = string(GetString(o, "stargate"))
	// TotalShips totalShips BasicInt ignore: false
	obj.TotalShips = GetInt[int](o, "totalShips")
	return obj
}

func SetFleetSpec(o js.Value, obj *cs.FleetSpec) {
	// ShipDesignSpec  Object ignore: false
	SetShipDesignSpec(o, &obj.ShipDesignSpec)
	// BaseCloakedCargo baseCloakedCargo BasicInt ignore: false
	o.Set("baseCloakedCargo", obj.BaseCloakedCargo)
	// BasePacketSpeed basePacketSpeed BasicInt ignore: false
	o.Set("basePacketSpeed", obj.BasePacketSpeed)
	// HasMassDriver hasMassDriver BasicBool ignore: false
	o.Set("hasMassDriver", obj.HasMassDriver)
	// HasStargate hasStargate BasicBool ignore: false
	o.Set("hasStargate", obj.HasStargate)
	// MassDriver massDriver BasicString ignore: false
	o.Set("massDriver", obj.MassDriver)
	// MassEmpty massEmpty BasicInt ignore: false
	o.Set("massEmpty", obj.MassEmpty)
	// MaxHullMass maxHullMass BasicInt ignore: false
	o.Set("maxHullMass", obj.MaxHullMass)
	// MaxRange maxRange BasicInt ignore: false
	o.Set("maxRange", obj.MaxRange)
	// Purposes purposes Map ignore: true
	// SafeHullMass safeHullMass BasicInt ignore: false
	o.Set("safeHullMass", obj.SafeHullMass)
	// SafeRange safeRange BasicInt ignore: false
	o.Set("safeRange", obj.SafeRange)
	// Stargate stargate BasicString ignore: false
	o.Set("stargate", obj.Stargate)
	// TotalShips totalShips BasicInt ignore: false
	o.Set("totalShips", obj.TotalShips)
}

func GetGameDBObject(o js.Value) cs.GameDBObject {
	obj := cs.GameDBObject{}
	// ID id BasicInt ignore: false
	obj.ID = GetInt[int64](o, "id")
	// GameID gameId BasicInt ignore: false
	obj.GameID = GetInt[int64](o, "gameId")
	// CreatedAt createdAt  ignore: false
	obj.CreatedAt, _ = GetTime(o, "createdAt")
	// UpdatedAt updatedAt  ignore: false
	obj.UpdatedAt, _ = GetTime(o, "updatedAt")
	return obj
}

func SetGameDBObject(o js.Value, obj *cs.GameDBObject) {
	// ID id BasicInt ignore: false
	o.Set("id", obj.ID)
	// GameID gameId BasicInt ignore: false
	o.Set("gameId", obj.GameID)
	// CreatedAt createdAt  ignore: false
	SetTime(o, "createdAt", obj.CreatedAt)
	// UpdatedAt updatedAt  ignore: false
	SetTime(o, "updatedAt", obj.UpdatedAt)
}

func GetHab(o js.Value) cs.Hab {
	obj := cs.Hab{}
	// Grav grav BasicInt ignore: false
	obj.Grav = GetInt[int](o, "grav")
	// Temp temp BasicInt ignore: false
	obj.Temp = GetInt[int](o, "temp")
	// Rad rad BasicInt ignore: false
	obj.Rad = GetInt[int](o, "rad")
	return obj
}

func SetHab(o js.Value, obj *cs.Hab) {
	// Grav grav BasicInt ignore: false
	o.Set("grav", obj.Grav)
	// Temp temp BasicInt ignore: false
	o.Set("temp", obj.Temp)
	// Rad rad BasicInt ignore: false
	o.Set("rad", obj.Rad)
}

func GetIntel(o js.Value) cs.Intel {
	obj := cs.Intel{}
	// Name name BasicString ignore: false
	obj.Name = string(GetString(o, "name"))
	// Num num BasicInt ignore: false
	obj.Num = GetInt[int](o, "num")
	// PlayerNum playerNum BasicInt ignore: false
	obj.PlayerNum = GetInt[int](o, "playerNum")
	// ReportAge reportAge BasicInt ignore: false
	obj.ReportAge = GetInt[int](o, "reportAge")
	return obj
}

func SetIntel(o js.Value, obj *cs.Intel) {
	// Name name BasicString ignore: false
	o.Set("name", obj.Name)
	// Num num BasicInt ignore: false
	o.Set("num", obj.Num)
	// PlayerNum playerNum BasicInt ignore: false
	o.Set("playerNum", obj.PlayerNum)
	// ReportAge reportAge BasicInt ignore: false
	o.Set("reportAge", obj.ReportAge)
}

func GetMapObject(o js.Value) cs.MapObject {
	obj := cs.MapObject{}
	// GameDBObject  Object ignore: false
	obj.GameDBObject = GetGameDBObject(o)
	// Type type Named ignore: false
	obj.Type = cs.MapObjectType(GetString(o, "type"))
	// Delete  BasicBool ignore: true
	// Position position Object ignore: false
	obj.Position = GetVector(o.Get("position"))
	// Num num BasicInt ignore: false
	obj.Num = GetInt[int](o, "num")
	// PlayerNum playerNum BasicInt ignore: false
	obj.PlayerNum = GetInt[int](o, "playerNum")
	// Name name BasicString ignore: false
	obj.Name = string(GetString(o, "name"))
	// Tags tags Named ignore: false
	// unknown type Tags cs.Tags map[string]string
	return obj
}

func SetMapObject(o js.Value, obj *cs.MapObject) {
	// GameDBObject  Object ignore: false
	SetGameDBObject(o, &obj.GameDBObject)
	// Type type Named ignore: false
	o.Set("type", string(obj.Type))
	// Delete  BasicBool ignore: true
	// Position position Object ignore: false
	o.Set("position", map[string]any{})
	SetVector(o.Get("position"), &obj.Position)
	// Num num BasicInt ignore: false
	o.Set("num", obj.Num)
	// PlayerNum playerNum BasicInt ignore: false
	o.Set("playerNum", obj.PlayerNum)
	// Name name BasicString ignore: false
	o.Set("name", obj.Name)
	// Tags tags Named ignore: false
	// unknown type Tags cs.Tags map[string]string
}

func GetMapObjectIntel(o js.Value) cs.MapObjectIntel {
	obj := cs.MapObjectIntel{}
	// Intel  Object ignore: false
	obj.Intel = GetIntel(o)
	// Type type Named ignore: false
	obj.Type = cs.MapObjectType(GetString(o, "type"))
	// Position position Object ignore: false
	obj.Position = GetVector(o.Get("position"))
	return obj
}

func SetMapObjectIntel(o js.Value, obj *cs.MapObjectIntel) {
	// Intel  Object ignore: false
	SetIntel(o, &obj.Intel)
	// Type type Named ignore: false
	o.Set("type", string(obj.Type))
	// Position position Object ignore: false
	o.Set("position", map[string]any{})
	SetVector(o.Get("position"), &obj.Position)
}

func GetMineField(o js.Value) cs.MineField {
	obj := cs.MineField{}
	// MapObject  Object ignore: false
	obj.MapObject = GetMapObject(o)
	// MineFieldOrders  Object ignore: false
	obj.MineFieldOrders = GetMineFieldOrders(o)
	// MineFieldType mineFieldType Named ignore: false
	obj.MineFieldType = cs.MineFieldType(GetString(o, "mineFieldType"))
	// NumMines numMines BasicInt ignore: false
	obj.NumMines = GetInt[int](o, "numMines")
	// Spec spec Object ignore: false
	obj.Spec = GetMineFieldSpec(o.Get("spec"))
	return obj
}

func SetMineField(o js.Value, obj *cs.MineField) {
	// MapObject  Object ignore: false
	SetMapObject(o, &obj.MapObject)
	// MineFieldOrders  Object ignore: false
	SetMineFieldOrders(o, &obj.MineFieldOrders)
	// MineFieldType mineFieldType Named ignore: false
	o.Set("mineFieldType", string(obj.MineFieldType))
	// NumMines numMines BasicInt ignore: false
	o.Set("numMines", obj.NumMines)
	// Spec spec Object ignore: false
	o.Set("spec", map[string]any{})
	SetMineFieldSpec(o.Get("spec"), &obj.Spec)
}

func GetMineFieldIntel(o js.Value) cs.MineFieldIntel {
	obj := cs.MineFieldIntel{}
	// MapObjectIntel  Object ignore: false
	obj.MapObjectIntel = GetMapObjectIntel(o)
	// NumMines numMines BasicInt ignore: false
	obj.NumMines = GetInt[int](o, "numMines")
	// MineFieldType mineFieldType Named ignore: false
	obj.MineFieldType = cs.MineFieldType(GetString(o, "mineFieldType"))
	// Spec spec Object ignore: false
	obj.Spec = GetMineFieldSpec(o.Get("spec"))
	return obj
}

func SetMineFieldIntel(o js.Value, obj *cs.MineFieldIntel) {
	// MapObjectIntel  Object ignore: false
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	// NumMines numMines BasicInt ignore: false
	o.Set("numMines", obj.NumMines)
	// MineFieldType mineFieldType Named ignore: false
	o.Set("mineFieldType", string(obj.MineFieldType))
	// Spec spec Object ignore: false
	o.Set("spec", map[string]any{})
	SetMineFieldSpec(o.Get("spec"), &obj.Spec)
}

func GetMineFieldOrders(o js.Value) cs.MineFieldOrders {
	obj := cs.MineFieldOrders{}
	// Detonate detonate BasicBool ignore: false
	obj.Detonate = bool(GetBool(o, "detonate"))
	return obj
}

func SetMineFieldOrders(o js.Value, obj *cs.MineFieldOrders) {
	// Detonate detonate BasicBool ignore: false
	o.Set("detonate", obj.Detonate)
}

func GetMineFieldSpec(o js.Value) cs.MineFieldSpec {
	obj := cs.MineFieldSpec{}
	// Radius radius BasicFloat ignore: false
	obj.Radius = GetFloat[float64](o, "radius")
	// DecayRate decayRate BasicInt ignore: false
	obj.DecayRate = GetInt[int](o, "decayRate")
	return obj
}

func SetMineFieldSpec(o js.Value, obj *cs.MineFieldSpec) {
	// Radius radius BasicFloat ignore: false
	o.Set("radius", obj.Radius)
	// DecayRate decayRate BasicInt ignore: false
	o.Set("decayRate", obj.DecayRate)
}

func GetMineral(o js.Value) cs.Mineral {
	obj := cs.Mineral{}
	// Ironium ironium BasicInt ignore: false
	obj.Ironium = GetInt[int](o, "ironium")
	// Boranium boranium BasicInt ignore: false
	obj.Boranium = GetInt[int](o, "boranium")
	// Germanium germanium BasicInt ignore: false
	obj.Germanium = GetInt[int](o, "germanium")
	return obj
}

func SetMineral(o js.Value, obj *cs.Mineral) {
	// Ironium ironium BasicInt ignore: false
	o.Set("ironium", obj.Ironium)
	// Boranium boranium BasicInt ignore: false
	o.Set("boranium", obj.Boranium)
	// Germanium germanium BasicInt ignore: false
	o.Set("germanium", obj.Germanium)
}

func GetMineralPacketDamage(o js.Value) cs.MineralPacketDamage {
	obj := cs.MineralPacketDamage{}
	// Killed killed BasicInt ignore: false
	obj.Killed = GetInt[int](o, "killed")
	// DefensesDestroyed defensesDestroyed BasicInt ignore: false
	obj.DefensesDestroyed = GetInt[int](o, "defensesDestroyed")
	// Uncaught uncaught BasicInt ignore: false
	obj.Uncaught = GetInt[int](o, "uncaught")
	return obj
}

func SetMineralPacketDamage(o js.Value, obj *cs.MineralPacketDamage) {
	// Killed killed BasicInt ignore: false
	o.Set("killed", obj.Killed)
	// DefensesDestroyed defensesDestroyed BasicInt ignore: false
	o.Set("defensesDestroyed", obj.DefensesDestroyed)
	// Uncaught uncaught BasicInt ignore: false
	o.Set("uncaught", obj.Uncaught)
}

func GetMineralPacketIntel(o js.Value) cs.MineralPacketIntel {
	obj := cs.MineralPacketIntel{}
	// MapObjectIntel  Object ignore: false
	obj.MapObjectIntel = GetMapObjectIntel(o)
	// WarpSpeed warpSpeed BasicInt ignore: false
	obj.WarpSpeed = GetInt[int](o, "warpSpeed")
	// Heading heading Object ignore: false
	obj.Heading = GetVector(o.Get("heading"))
	// Cargo cargo Object ignore: false
	obj.Cargo = GetCargo(o.Get("cargo"))
	// TargetPlanetNum targetPlanetNum BasicInt ignore: false
	obj.TargetPlanetNum = GetInt[int](o, "targetPlanetNum")
	// ScanRange scanRange BasicInt ignore: false
	obj.ScanRange = GetInt[int](o, "scanRange")
	// ScanRangePen scanRangePen BasicInt ignore: false
	obj.ScanRangePen = GetInt[int](o, "scanRangePen")
	return obj
}

func SetMineralPacketIntel(o js.Value, obj *cs.MineralPacketIntel) {
	// MapObjectIntel  Object ignore: false
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	// WarpSpeed warpSpeed BasicInt ignore: false
	o.Set("warpSpeed", obj.WarpSpeed)
	// Heading heading Object ignore: false
	o.Set("heading", map[string]any{})
	SetVector(o.Get("heading"), &obj.Heading)
	// Cargo cargo Object ignore: false
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), &obj.Cargo)
	// TargetPlanetNum targetPlanetNum BasicInt ignore: false
	o.Set("targetPlanetNum", obj.TargetPlanetNum)
	// ScanRange scanRange BasicInt ignore: false
	o.Set("scanRange", obj.ScanRange)
	// ScanRangePen scanRangePen BasicInt ignore: false
	o.Set("scanRangePen", obj.ScanRangePen)
}

func GetMysteryTrader(o js.Value) cs.MysteryTrader {
	obj := cs.MysteryTrader{}
	// MapObject  Object ignore: false
	obj.MapObject = GetMapObject(o)
	// WarpSpeed warpSpeed BasicInt ignore: false
	obj.WarpSpeed = GetInt[int](o, "warpSpeed")
	// Destination destination Object ignore: false
	obj.Destination = GetVector(o.Get("destination"))
	// RequestedBoon requestedBoon BasicInt ignore: false
	obj.RequestedBoon = GetInt[int](o, "requestedBoon")
	// RewardType rewardType Named ignore: false
	obj.RewardType = cs.MysteryTraderRewardType(GetString(o, "rewardType"))
	// Heading heading Object ignore: false
	obj.Heading = GetVector(o.Get("heading"))
	// PlayersRewarded playersRewarded Map ignore: false

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
	// Spec spec Object ignore: false
	obj.Spec = GetMysteryTraderSpec(o.Get("spec"))
	return obj
}

func SetMysteryTrader(o js.Value, obj *cs.MysteryTrader) {
	// MapObject  Object ignore: false
	SetMapObject(o, &obj.MapObject)
	// WarpSpeed warpSpeed BasicInt ignore: false
	o.Set("warpSpeed", obj.WarpSpeed)
	// Destination destination Object ignore: false
	o.Set("destination", map[string]any{})
	SetVector(o.Get("destination"), &obj.Destination)
	// RequestedBoon requestedBoon BasicInt ignore: false
	o.Set("requestedBoon", obj.RequestedBoon)
	// RewardType rewardType Named ignore: false
	o.Set("rewardType", string(obj.RewardType))
	// Heading heading Object ignore: false
	o.Set("heading", map[string]any{})
	SetVector(o.Get("heading"), &obj.Heading)
	// PlayersRewarded playersRewarded Map ignore: false
	playersRewardedMap := js.ValueOf(map[string]any{})
	for key, value := range obj.PlayersRewarded {
		playersRewardedMap.Set(fmt.Sprintf("%v", key), bool(value))
	}
	o.Set("playersRewarded", playersRewardedMap)
	// Spec spec Object ignore: false
	o.Set("spec", map[string]any{})
	SetMysteryTraderSpec(o.Get("spec"), &obj.Spec)
}

func GetMysteryTraderIntel(o js.Value) cs.MysteryTraderIntel {
	obj := cs.MysteryTraderIntel{}
	// MapObjectIntel  Object ignore: false
	obj.MapObjectIntel = GetMapObjectIntel(o)
	// WarpSpeed warpSpeed BasicInt ignore: false
	obj.WarpSpeed = GetInt[int](o, "warpSpeed")
	// Heading heading Object ignore: false
	obj.Heading = GetVector(o.Get("heading"))
	// RequestedBoon requestedBoon BasicInt ignore: false
	obj.RequestedBoon = GetInt[int](o, "requestedBoon")
	return obj
}

func SetMysteryTraderIntel(o js.Value, obj *cs.MysteryTraderIntel) {
	// MapObjectIntel  Object ignore: false
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	// WarpSpeed warpSpeed BasicInt ignore: false
	o.Set("warpSpeed", obj.WarpSpeed)
	// Heading heading Object ignore: false
	o.Set("heading", map[string]any{})
	SetVector(o.Get("heading"), &obj.Heading)
	// RequestedBoon requestedBoon BasicInt ignore: false
	o.Set("requestedBoon", obj.RequestedBoon)
}

func GetMysteryTraderReward(o js.Value) cs.MysteryTraderReward {
	obj := cs.MysteryTraderReward{}
	// Type type Named ignore: false
	obj.Type = cs.MysteryTraderRewardType(GetString(o, "type"))
	// TechLevels techLevels Object ignore: false
	obj.TechLevels = GetTechLevel(o.Get("techLevels"))
	// Tech tech BasicString ignore: false
	obj.Tech = string(GetString(o, "tech"))
	// Ship ship Object ignore: false
	obj.Ship = GetShipDesign(o.Get("ship"))
	// ShipCount shipCount BasicInt ignore: false
	obj.ShipCount = GetInt[int](o, "shipCount")
	return obj
}

func SetMysteryTraderReward(o js.Value, obj *cs.MysteryTraderReward) {
	// Type type Named ignore: false
	o.Set("type", string(obj.Type))
	// TechLevels techLevels Object ignore: false
	o.Set("techLevels", map[string]any{})
	SetTechLevel(o.Get("techLevels"), &obj.TechLevels)
	// Tech tech BasicString ignore: false
	o.Set("tech", obj.Tech)
	// Ship ship Object ignore: false
	o.Set("ship", map[string]any{})
	SetShipDesign(o.Get("ship"), &obj.Ship)
	// ShipCount shipCount BasicInt ignore: false
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
	// MapObject  Object ignore: false
	obj.MapObject = GetMapObject(o)
	// PlanetOrders  Object ignore: false
	obj.PlanetOrders = GetPlanetOrders(o)
	// Hab hab Object ignore: false
	obj.Hab = GetHab(o.Get("hab"))
	// BaseHab baseHab Object ignore: false
	obj.BaseHab = GetHab(o.Get("baseHab"))
	// TerraformedAmount terraformedAmount Object ignore: false
	obj.TerraformedAmount = GetHab(o.Get("terraformedAmount"))
	// MineralConcentration mineralConcentration Object ignore: false
	obj.MineralConcentration = GetMineral(o.Get("mineralConcentration"))
	// MineYears mineYears Object ignore: false
	obj.MineYears = GetMineral(o.Get("mineYears"))
	// Cargo cargo Object ignore: false
	obj.Cargo = GetCargo(o.Get("cargo"))
	// Mines mines BasicInt ignore: false
	obj.Mines = GetInt[int](o, "mines")
	// Factories factories BasicInt ignore: false
	obj.Factories = GetInt[int](o, "factories")
	// Defenses defenses BasicInt ignore: false
	obj.Defenses = GetInt[int](o, "defenses")
	// Homeworld homeworld BasicBool ignore: false
	obj.Homeworld = bool(GetBool(o, "homeworld"))
	// Scanner scanner BasicBool ignore: false
	obj.Scanner = bool(GetBool(o, "scanner"))
	// Spec spec Object ignore: false
	obj.Spec = GetPlanetSpec(o.Get("spec"))
	// RandomArtifact  BasicBool ignore: true
	// Starbase  Object ignore: true
	// Dirty  BasicBool ignore: true
	// bonusResources   ignore: true
	return obj
}

func SetPlanet(o js.Value, obj *cs.Planet) {
	// MapObject  Object ignore: false
	SetMapObject(o, &obj.MapObject)
	// PlanetOrders  Object ignore: false
	SetPlanetOrders(o, &obj.PlanetOrders)
	// Hab hab Object ignore: false
	o.Set("hab", map[string]any{})
	SetHab(o.Get("hab"), &obj.Hab)
	// BaseHab baseHab Object ignore: false
	o.Set("baseHab", map[string]any{})
	SetHab(o.Get("baseHab"), &obj.BaseHab)
	// TerraformedAmount terraformedAmount Object ignore: false
	o.Set("terraformedAmount", map[string]any{})
	SetHab(o.Get("terraformedAmount"), &obj.TerraformedAmount)
	// MineralConcentration mineralConcentration Object ignore: false
	o.Set("mineralConcentration", map[string]any{})
	SetMineral(o.Get("mineralConcentration"), &obj.MineralConcentration)
	// MineYears mineYears Object ignore: false
	o.Set("mineYears", map[string]any{})
	SetMineral(o.Get("mineYears"), &obj.MineYears)
	// Cargo cargo Object ignore: false
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), &obj.Cargo)
	// Mines mines BasicInt ignore: false
	o.Set("mines", obj.Mines)
	// Factories factories BasicInt ignore: false
	o.Set("factories", obj.Factories)
	// Defenses defenses BasicInt ignore: false
	o.Set("defenses", obj.Defenses)
	// Homeworld homeworld BasicBool ignore: false
	o.Set("homeworld", obj.Homeworld)
	// Scanner scanner BasicBool ignore: false
	o.Set("scanner", obj.Scanner)
	// Spec spec Object ignore: false
	o.Set("spec", map[string]any{})
	SetPlanetSpec(o.Get("spec"), &obj.Spec)
	// RandomArtifact  BasicBool ignore: true
	// Starbase  Object ignore: true
	// Dirty  BasicBool ignore: true
	// bonusResources   ignore: true
}

func GetPlanetIntel(o js.Value) cs.PlanetIntel {
	obj := cs.PlanetIntel{}
	// MapObjectIntel  Object ignore: false
	obj.MapObjectIntel = GetMapObjectIntel(o)
	// Hab hab Object ignore: false
	obj.Hab = GetHab(o.Get("hab"))
	// BaseHab baseHab Object ignore: false
	obj.BaseHab = GetHab(o.Get("baseHab"))
	// MineralConcentration mineralConcentration Object ignore: false
	obj.MineralConcentration = GetMineral(o.Get("mineralConcentration"))
	// Starbase starbase Object ignore: false
	starbaseVal := o.Get("starbase")
	if !starbaseVal.IsUndefined() {
		starbase := GetFleetIntel(starbaseVal)
		obj.Starbase = &starbase
	}
	// Cargo cargo Object ignore: false
	obj.Cargo = GetCargo(o.Get("cargo"))
	// CargoDiscovered cargoDiscovered BasicBool ignore: false
	obj.CargoDiscovered = bool(GetBool(o, "cargoDiscovered"))
	// PlanetHabitability planetHabitability BasicInt ignore: false
	obj.PlanetHabitability = GetInt[int](o, "planetHabitability")
	// PlanetHabitabilityTerraformed planetHabitabilityTerraformed BasicInt ignore: false
	obj.PlanetHabitabilityTerraformed = GetInt[int](o, "planetHabitabilityTerraformed")
	// Homeworld homeworld BasicBool ignore: false
	obj.Homeworld = bool(GetBool(o, "homeworld"))
	// Spec spec Object ignore: false
	obj.Spec = GetPlanetSpec(o.Get("spec"))
	return obj
}

func SetPlanetIntel(o js.Value, obj *cs.PlanetIntel) {
	// MapObjectIntel  Object ignore: false
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	// Hab hab Object ignore: false
	o.Set("hab", map[string]any{})
	SetHab(o.Get("hab"), &obj.Hab)
	// BaseHab baseHab Object ignore: false
	o.Set("baseHab", map[string]any{})
	SetHab(o.Get("baseHab"), &obj.BaseHab)
	// MineralConcentration mineralConcentration Object ignore: false
	o.Set("mineralConcentration", map[string]any{})
	SetMineral(o.Get("mineralConcentration"), &obj.MineralConcentration)
	// Starbase starbase Object ignore: false
	o.Set("starbase", map[string]any{})
	SetFleetIntel(o.Get("starbase"), obj.Starbase)
	// Cargo cargo Object ignore: false
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), &obj.Cargo)
	// CargoDiscovered cargoDiscovered BasicBool ignore: false
	o.Set("cargoDiscovered", obj.CargoDiscovered)
	// PlanetHabitability planetHabitability BasicInt ignore: false
	o.Set("planetHabitability", obj.PlanetHabitability)
	// PlanetHabitabilityTerraformed planetHabitabilityTerraformed BasicInt ignore: false
	o.Set("planetHabitabilityTerraformed", obj.PlanetHabitabilityTerraformed)
	// Homeworld homeworld BasicBool ignore: false
	o.Set("homeworld", obj.Homeworld)
	// Spec spec Object ignore: false
	o.Set("spec", map[string]any{})
	SetPlanetSpec(o.Get("spec"), &obj.Spec)
}

func GetPlanetOrders(o js.Value) cs.PlanetOrders {
	obj := cs.PlanetOrders{}
	// ContributesOnlyLeftoverToResearch contributesOnlyLeftoverToResearch BasicBool ignore: false
	obj.ContributesOnlyLeftoverToResearch = bool(GetBool(o, "contributesOnlyLeftoverToResearch"))
	// ProductionQueue productionQueue Slice ignore: false
	obj.ProductionQueue = GetSlice(o.Get("productionQueue"), GetProductionQueueItem)
	// RouteTargetType routeTargetType Named ignore: false
	obj.RouteTargetType = cs.MapObjectType(GetString(o, "routeTargetType"))
	// RouteTargetNum routeTargetNum BasicInt ignore: false
	obj.RouteTargetNum = GetInt[int](o, "routeTargetNum")
	// RouteTargetPlayerNum routeTargetPlayerNum BasicInt ignore: false
	obj.RouteTargetPlayerNum = GetInt[int](o, "routeTargetPlayerNum")
	// PacketTargetNum packetTargetNum BasicInt ignore: false
	obj.PacketTargetNum = GetInt[int](o, "packetTargetNum")
	// PacketSpeed packetSpeed BasicInt ignore: false
	obj.PacketSpeed = GetInt[int](o, "packetSpeed")
	return obj
}

func SetPlanetOrders(o js.Value, obj *cs.PlanetOrders) {
	// ContributesOnlyLeftoverToResearch contributesOnlyLeftoverToResearch BasicBool ignore: false
	o.Set("contributesOnlyLeftoverToResearch", obj.ContributesOnlyLeftoverToResearch)
	// ProductionQueue productionQueue Slice ignore: false
	o.Set("productionQueue", []any{})
	SetSlice(o.Get("productionQueue"), obj.ProductionQueue, SetProductionQueueItem)
	// RouteTargetType routeTargetType Named ignore: false
	o.Set("routeTargetType", string(obj.RouteTargetType))
	// RouteTargetNum routeTargetNum BasicInt ignore: false
	o.Set("routeTargetNum", obj.RouteTargetNum)
	// RouteTargetPlayerNum routeTargetPlayerNum BasicInt ignore: false
	o.Set("routeTargetPlayerNum", obj.RouteTargetPlayerNum)
	// PacketTargetNum packetTargetNum BasicInt ignore: false
	o.Set("packetTargetNum", obj.PacketTargetNum)
	// PacketSpeed packetSpeed BasicInt ignore: false
	o.Set("packetSpeed", obj.PacketSpeed)
}

func GetPlanetSpec(o js.Value) cs.PlanetSpec {
	obj := cs.PlanetSpec{}
	// PlanetStarbaseSpec  Object ignore: false
	obj.PlanetStarbaseSpec = GetPlanetStarbaseSpec(o)
	// CanTerraform canTerraform BasicBool ignore: false
	obj.CanTerraform = bool(GetBool(o, "canTerraform"))
	// Defense defense BasicString ignore: false
	obj.Defense = string(GetString(o, "defense"))
	// DefenseCoverage defenseCoverage BasicFloat ignore: false
	obj.DefenseCoverage = GetFloat[float64](o, "defenseCoverage")
	// DefenseCoverageSmart defenseCoverageSmart BasicFloat ignore: false
	obj.DefenseCoverageSmart = GetFloat[float64](o, "defenseCoverageSmart")
	// GrowthAmount growthAmount BasicInt ignore: false
	obj.GrowthAmount = GetInt[int](o, "growthAmount")
	// Habitability habitability BasicInt ignore: false
	obj.Habitability = GetInt[int](o, "habitability")
	// MaxDefenses maxDefenses BasicInt ignore: false
	obj.MaxDefenses = GetInt[int](o, "maxDefenses")
	// MaxFactories maxFactories BasicInt ignore: false
	obj.MaxFactories = GetInt[int](o, "maxFactories")
	// MaxMines maxMines BasicInt ignore: false
	obj.MaxMines = GetInt[int](o, "maxMines")
	// MaxPopulation maxPopulation BasicInt ignore: false
	obj.MaxPopulation = GetInt[int](o, "maxPopulation")
	// MaxPossibleFactories maxPossibleFactories BasicInt ignore: false
	obj.MaxPossibleFactories = GetInt[int](o, "maxPossibleFactories")
	// MaxPossibleMines maxPossibleMines BasicInt ignore: false
	obj.MaxPossibleMines = GetInt[int](o, "maxPossibleMines")
	// MiningOutput miningOutput Object ignore: false
	obj.MiningOutput = GetMineral(o.Get("miningOutput"))
	// Population population BasicInt ignore: false
	obj.Population = GetInt[int](o, "population")
	// PopulationDensity populationDensity BasicFloat ignore: false
	obj.PopulationDensity = GetFloat[float64](o, "populationDensity")
	// ResourcesPerYear resourcesPerYear BasicInt ignore: false
	obj.ResourcesPerYear = GetInt[int](o, "resourcesPerYear")
	// ResourcesPerYearAvailable resourcesPerYearAvailable BasicInt ignore: false
	obj.ResourcesPerYearAvailable = GetInt[int](o, "resourcesPerYearAvailable")
	// ResourcesPerYearResearch resourcesPerYearResearch BasicInt ignore: false
	obj.ResourcesPerYearResearch = GetInt[int](o, "resourcesPerYearResearch")
	// ResourcesPerYearResearchEstimatedLeftover resourcesPerYearResearchEstimatedLeftover BasicInt ignore: false
	obj.ResourcesPerYearResearchEstimatedLeftover = GetInt[int](o, "resourcesPerYearResearchEstimatedLeftover")
	// Scanner scanner BasicString ignore: false
	obj.Scanner = string(GetString(o, "scanner"))
	// ScanRange scanRange BasicInt ignore: false
	obj.ScanRange = GetInt[int](o, "scanRange")
	// ScanRangePen scanRangePen BasicInt ignore: false
	obj.ScanRangePen = GetInt[int](o, "scanRangePen")
	// TerraformAmount terraformAmount Object ignore: false
	obj.TerraformAmount = GetHab(o.Get("terraformAmount"))
	// MinTerraformAmount minTerraformAmount Object ignore: false
	obj.MinTerraformAmount = GetHab(o.Get("minTerraformAmount"))
	// TerraformedHabitability terraformedHabitability BasicInt ignore: false
	obj.TerraformedHabitability = GetInt[int](o, "terraformedHabitability")
	// Contested contested BasicBool ignore: false
	obj.Contested = bool(GetBool(o, "contested"))
	return obj
}

func SetPlanetSpec(o js.Value, obj *cs.PlanetSpec) {
	// PlanetStarbaseSpec  Object ignore: false
	SetPlanetStarbaseSpec(o, &obj.PlanetStarbaseSpec)
	// CanTerraform canTerraform BasicBool ignore: false
	o.Set("canTerraform", obj.CanTerraform)
	// Defense defense BasicString ignore: false
	o.Set("defense", obj.Defense)
	// DefenseCoverage defenseCoverage BasicFloat ignore: false
	o.Set("defenseCoverage", obj.DefenseCoverage)
	// DefenseCoverageSmart defenseCoverageSmart BasicFloat ignore: false
	o.Set("defenseCoverageSmart", obj.DefenseCoverageSmart)
	// GrowthAmount growthAmount BasicInt ignore: false
	o.Set("growthAmount", obj.GrowthAmount)
	// Habitability habitability BasicInt ignore: false
	o.Set("habitability", obj.Habitability)
	// MaxDefenses maxDefenses BasicInt ignore: false
	o.Set("maxDefenses", obj.MaxDefenses)
	// MaxFactories maxFactories BasicInt ignore: false
	o.Set("maxFactories", obj.MaxFactories)
	// MaxMines maxMines BasicInt ignore: false
	o.Set("maxMines", obj.MaxMines)
	// MaxPopulation maxPopulation BasicInt ignore: false
	o.Set("maxPopulation", obj.MaxPopulation)
	// MaxPossibleFactories maxPossibleFactories BasicInt ignore: false
	o.Set("maxPossibleFactories", obj.MaxPossibleFactories)
	// MaxPossibleMines maxPossibleMines BasicInt ignore: false
	o.Set("maxPossibleMines", obj.MaxPossibleMines)
	// MiningOutput miningOutput Object ignore: false
	o.Set("miningOutput", map[string]any{})
	SetMineral(o.Get("miningOutput"), &obj.MiningOutput)
	// Population population BasicInt ignore: false
	o.Set("population", obj.Population)
	// PopulationDensity populationDensity BasicFloat ignore: false
	o.Set("populationDensity", obj.PopulationDensity)
	// ResourcesPerYear resourcesPerYear BasicInt ignore: false
	o.Set("resourcesPerYear", obj.ResourcesPerYear)
	// ResourcesPerYearAvailable resourcesPerYearAvailable BasicInt ignore: false
	o.Set("resourcesPerYearAvailable", obj.ResourcesPerYearAvailable)
	// ResourcesPerYearResearch resourcesPerYearResearch BasicInt ignore: false
	o.Set("resourcesPerYearResearch", obj.ResourcesPerYearResearch)
	// ResourcesPerYearResearchEstimatedLeftover resourcesPerYearResearchEstimatedLeftover BasicInt ignore: false
	o.Set("resourcesPerYearResearchEstimatedLeftover", obj.ResourcesPerYearResearchEstimatedLeftover)
	// Scanner scanner BasicString ignore: false
	o.Set("scanner", obj.Scanner)
	// ScanRange scanRange BasicInt ignore: false
	o.Set("scanRange", obj.ScanRange)
	// ScanRangePen scanRangePen BasicInt ignore: false
	o.Set("scanRangePen", obj.ScanRangePen)
	// TerraformAmount terraformAmount Object ignore: false
	o.Set("terraformAmount", map[string]any{})
	SetHab(o.Get("terraformAmount"), &obj.TerraformAmount)
	// MinTerraformAmount minTerraformAmount Object ignore: false
	o.Set("minTerraformAmount", map[string]any{})
	SetHab(o.Get("minTerraformAmount"), &obj.MinTerraformAmount)
	// TerraformedHabitability terraformedHabitability BasicInt ignore: false
	o.Set("terraformedHabitability", obj.TerraformedHabitability)
	// Contested contested BasicBool ignore: false
	o.Set("contested", obj.Contested)
}

func GetPlanetStarbaseSpec(o js.Value) cs.PlanetStarbaseSpec {
	obj := cs.PlanetStarbaseSpec{}
	// HasMassDriver hasMassDriver BasicBool ignore: false
	obj.HasMassDriver = bool(GetBool(o, "hasMassDriver"))
	// HasStarbase hasStarbase BasicBool ignore: false
	obj.HasStarbase = bool(GetBool(o, "hasStarbase"))
	// HasStargate hasStargate BasicBool ignore: false
	obj.HasStargate = bool(GetBool(o, "hasStargate"))
	// StarbaseDesignName starbaseDesignName BasicString ignore: false
	obj.StarbaseDesignName = string(GetString(o, "starbaseDesignName"))
	// StarbaseDesignNum starbaseDesignNum BasicInt ignore: false
	obj.StarbaseDesignNum = GetInt[int](o, "starbaseDesignNum")
	// DockCapacity dockCapacity BasicInt ignore: false
	obj.DockCapacity = GetInt[int](o, "dockCapacity")
	// BasePacketSpeed basePacketSpeed BasicInt ignore: false
	obj.BasePacketSpeed = GetInt[int](o, "basePacketSpeed")
	// SafePacketSpeed safePacketSpeed BasicInt ignore: false
	obj.SafePacketSpeed = GetInt[int](o, "safePacketSpeed")
	// SafeHullMass safeHullMass BasicInt ignore: false
	obj.SafeHullMass = GetInt[int](o, "safeHullMass")
	// SafeRange safeRange BasicInt ignore: false
	obj.SafeRange = GetInt[int](o, "safeRange")
	// MaxRange maxRange BasicInt ignore: false
	obj.MaxRange = GetInt[int](o, "maxRange")
	// MaxHullMass maxHullMass BasicInt ignore: false
	obj.MaxHullMass = GetInt[int](o, "maxHullMass")
	// Stargate stargate BasicString ignore: false
	obj.Stargate = string(GetString(o, "stargate"))
	// MassDriver massDriver BasicString ignore: false
	obj.MassDriver = string(GetString(o, "massDriver"))
	return obj
}

func SetPlanetStarbaseSpec(o js.Value, obj *cs.PlanetStarbaseSpec) {
	// HasMassDriver hasMassDriver BasicBool ignore: false
	o.Set("hasMassDriver", obj.HasMassDriver)
	// HasStarbase hasStarbase BasicBool ignore: false
	o.Set("hasStarbase", obj.HasStarbase)
	// HasStargate hasStargate BasicBool ignore: false
	o.Set("hasStargate", obj.HasStargate)
	// StarbaseDesignName starbaseDesignName BasicString ignore: false
	o.Set("starbaseDesignName", obj.StarbaseDesignName)
	// StarbaseDesignNum starbaseDesignNum BasicInt ignore: false
	o.Set("starbaseDesignNum", obj.StarbaseDesignNum)
	// DockCapacity dockCapacity BasicInt ignore: false
	o.Set("dockCapacity", obj.DockCapacity)
	// BasePacketSpeed basePacketSpeed BasicInt ignore: false
	o.Set("basePacketSpeed", obj.BasePacketSpeed)
	// SafePacketSpeed safePacketSpeed BasicInt ignore: false
	o.Set("safePacketSpeed", obj.SafePacketSpeed)
	// SafeHullMass safeHullMass BasicInt ignore: false
	o.Set("safeHullMass", obj.SafeHullMass)
	// SafeRange safeRange BasicInt ignore: false
	o.Set("safeRange", obj.SafeRange)
	// MaxRange maxRange BasicInt ignore: false
	o.Set("maxRange", obj.MaxRange)
	// MaxHullMass maxHullMass BasicInt ignore: false
	o.Set("maxHullMass", obj.MaxHullMass)
	// Stargate stargate BasicString ignore: false
	o.Set("stargate", obj.Stargate)
	// MassDriver massDriver BasicString ignore: false
	o.Set("massDriver", obj.MassDriver)
}

func GetPlayer(o js.Value) cs.Player {
	obj := cs.Player{}
	// GameDBObject  Object ignore: false
	obj.GameDBObject = GetGameDBObject(o)
	// PlayerOrders  Object ignore: false
	obj.PlayerOrders = GetPlayerOrders(o)
	// PlayerIntels  Object ignore: false
	obj.PlayerIntels = GetPlayerIntels(o)
	// PlayerPlans  Object ignore: false
	obj.PlayerPlans = GetPlayerPlans(o)
	// UserID userId BasicInt ignore: false
	obj.UserID = GetInt[int64](o, "userId")
	// Name name BasicString ignore: false
	obj.Name = string(GetString(o, "name"))
	// Num num BasicInt ignore: false
	obj.Num = GetInt[int](o, "num")
	// Ready ready BasicBool ignore: false
	obj.Ready = bool(GetBool(o, "ready"))
	// AIControlled aiControlled BasicBool ignore: false
	obj.AIControlled = bool(GetBool(o, "aiControlled"))
	// AIDifficulty aiDifficulty Named ignore: false
	obj.AIDifficulty = cs.AIDifficulty(GetString(o, "aiDifficulty"))
	// Guest guest BasicBool ignore: false
	obj.Guest = bool(GetBool(o, "guest"))
	// SubmittedTurn submittedTurn BasicBool ignore: false
	obj.SubmittedTurn = bool(GetBool(o, "submittedTurn"))
	// Color color BasicString ignore: false
	obj.Color = string(GetString(o, "color"))
	// DefaultHullSet defaultHullSet BasicInt ignore: false
	obj.DefaultHullSet = GetInt[int](o, "defaultHullSet")
	// Race race Object ignore: false
	obj.Race = GetRace(o.Get("race"))
	// TechLevels techLevels Object ignore: false
	obj.TechLevels = GetTechLevel(o.Get("techLevels"))
	// TechLevelsSpent techLevelsSpent Object ignore: false
	obj.TechLevelsSpent = GetTechLevel(o.Get("techLevelsSpent"))
	// ResearchSpentLastYear researchSpentLastYear BasicInt ignore: false
	obj.ResearchSpentLastYear = GetInt[int](o, "researchSpentLastYear")
	// Relations relations Slice ignore: false
	obj.Relations = GetSlice(o.Get("relations"), GetPlayerRelationship)
	// Messages messages Slice ignore: false
	obj.Messages = GetSlice(o.Get("messages"), GetPlayerMessage)
	// Designs designs Slice ignore: false
	designs := GetSlice(o.Get("designs"), GetShipDesign)
	obj.Designs = make([]*cs.ShipDesign, len(designs))
	for i := range designs {
		obj.Designs[i] = &designs[i]
	}
	// ScoreHistory scoreHistory Slice ignore: false
	obj.ScoreHistory = GetSlice(o.Get("scoreHistory"), GetPlayerScore)
	// AcquiredTechs acquiredTechs Map ignore: false

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
	// AchievedVictoryConditions achievedVictoryConditions Named ignore: false
	obj.AchievedVictoryConditions = cs.Bitmask(GetInt[cs.Bitmask](o, "achievedVictoryConditions"))
	// Victor victor BasicBool ignore: false
	obj.Victor = bool(GetBool(o, "victor"))
	// Stats stats Object ignore: false
	statsVal := o.Get("stats")
	if !statsVal.IsUndefined() {
		stats := GetPlayerStats(statsVal)
		obj.Stats = &stats
	}
	// Spec spec Object ignore: false
	obj.Spec = GetPlayerSpec(o.Get("spec"))
	// leftoverResources   ignore: true
	// techLevelGained   ignore: true
	// discoverer   ignore: true
	return obj
}

func SetPlayer(o js.Value, obj *cs.Player) {
	// GameDBObject  Object ignore: false
	SetGameDBObject(o, &obj.GameDBObject)
	// PlayerOrders  Object ignore: false
	SetPlayerOrders(o, &obj.PlayerOrders)
	// PlayerIntels  Object ignore: false
	SetPlayerIntels(o, &obj.PlayerIntels)
	// PlayerPlans  Object ignore: false
	SetPlayerPlans(o, &obj.PlayerPlans)
	// UserID userId BasicInt ignore: false
	o.Set("userId", obj.UserID)
	// Name name BasicString ignore: false
	o.Set("name", obj.Name)
	// Num num BasicInt ignore: false
	o.Set("num", obj.Num)
	// Ready ready BasicBool ignore: false
	o.Set("ready", obj.Ready)
	// AIControlled aiControlled BasicBool ignore: false
	o.Set("aiControlled", obj.AIControlled)
	// AIDifficulty aiDifficulty Named ignore: false
	o.Set("aiDifficulty", string(obj.AIDifficulty))
	// Guest guest BasicBool ignore: false
	o.Set("guest", obj.Guest)
	// SubmittedTurn submittedTurn BasicBool ignore: false
	o.Set("submittedTurn", obj.SubmittedTurn)
	// Color color BasicString ignore: false
	o.Set("color", obj.Color)
	// DefaultHullSet defaultHullSet BasicInt ignore: false
	o.Set("defaultHullSet", obj.DefaultHullSet)
	// Race race Object ignore: false
	o.Set("race", map[string]any{})
	SetRace(o.Get("race"), &obj.Race)
	// TechLevels techLevels Object ignore: false
	o.Set("techLevels", map[string]any{})
	SetTechLevel(o.Get("techLevels"), &obj.TechLevels)
	// TechLevelsSpent techLevelsSpent Object ignore: false
	o.Set("techLevelsSpent", map[string]any{})
	SetTechLevel(o.Get("techLevelsSpent"), &obj.TechLevelsSpent)
	// ResearchSpentLastYear researchSpentLastYear BasicInt ignore: false
	o.Set("researchSpentLastYear", obj.ResearchSpentLastYear)
	// Relations relations Slice ignore: false
	o.Set("relations", []any{})
	SetSlice(o.Get("relations"), obj.Relations, SetPlayerRelationship)
	// Messages messages Slice ignore: false
	o.Set("messages", []any{})
	SetSlice(o.Get("messages"), obj.Messages, SetPlayerMessage)
	// Designs designs Slice ignore: false
	o.Set("designs", []any{})
	SetPointerSlice(o.Get("designs"), obj.Designs, SetShipDesign)
	// ScoreHistory scoreHistory Slice ignore: false
	o.Set("scoreHistory", []any{})
	SetSlice(o.Get("scoreHistory"), obj.ScoreHistory, SetPlayerScore)
	// AcquiredTechs acquiredTechs Map ignore: false
	acquiredTechsMap := js.ValueOf(map[string]any{})
	for key, value := range obj.AcquiredTechs {
		acquiredTechsMap.Set(fmt.Sprintf("%v", key), bool(value))
	}
	o.Set("acquiredTechs", acquiredTechsMap)
	// AchievedVictoryConditions achievedVictoryConditions Named ignore: false
	o.Set("achievedVictoryConditions", uint32(obj.AchievedVictoryConditions))
	// Victor victor BasicBool ignore: false
	o.Set("victor", obj.Victor)
	// Stats stats Object ignore: false
	o.Set("stats", map[string]any{})
	SetPlayerStats(o.Get("stats"), obj.Stats)
	// Spec spec Object ignore: false
	o.Set("spec", map[string]any{})
	SetPlayerSpec(o.Get("spec"), &obj.Spec)
	// leftoverResources   ignore: true
	// techLevelGained   ignore: true
	// discoverer   ignore: true
}

func GetPlayerIntel(o js.Value) cs.PlayerIntel {
	obj := cs.PlayerIntel{}
	// Name name BasicString ignore: false
	obj.Name = string(GetString(o, "name"))
	// Num num BasicInt ignore: false
	obj.Num = GetInt[int](o, "num")
	// Color color BasicString ignore: false
	obj.Color = string(GetString(o, "color"))
	// Seen seen BasicBool ignore: false
	obj.Seen = bool(GetBool(o, "seen"))
	// RaceName raceName BasicString ignore: false
	obj.RaceName = string(GetString(o, "raceName"))
	// RacePluralName racePluralName BasicString ignore: false
	obj.RacePluralName = string(GetString(o, "racePluralName"))
	return obj
}

func SetPlayerIntel(o js.Value, obj *cs.PlayerIntel) {
	// Name name BasicString ignore: false
	o.Set("name", obj.Name)
	// Num num BasicInt ignore: false
	o.Set("num", obj.Num)
	// Color color BasicString ignore: false
	o.Set("color", obj.Color)
	// Seen seen BasicBool ignore: false
	o.Set("seen", obj.Seen)
	// RaceName raceName BasicString ignore: false
	o.Set("raceName", obj.RaceName)
	// RacePluralName racePluralName BasicString ignore: false
	o.Set("racePluralName", obj.RacePluralName)
}

func GetPlayerIntels(o js.Value) cs.PlayerIntels {
	obj := cs.PlayerIntels{}
	// BattleRecords battleRecords Slice ignore: false
	obj.BattleRecords = GetSlice(o.Get("battleRecords"), GetBattleRecord)
	// PlayerIntels playerIntels Slice ignore: false
	obj.PlayerIntels = GetSlice(o.Get("playerIntels"), GetPlayerIntel)
	// ScoreIntels scoreIntels Slice ignore: false
	obj.ScoreIntels = GetSlice(o.Get("scoreIntels"), GetScoreIntel)
	// PlanetIntels planetIntels Slice ignore: false
	obj.PlanetIntels = GetSlice(o.Get("planetIntels"), GetPlanetIntel)
	// FleetIntels fleetIntels Slice ignore: false
	obj.FleetIntels = GetSlice(o.Get("fleetIntels"), GetFleetIntel)
	// StarbaseIntels starbaseIntels Slice ignore: false
	obj.StarbaseIntels = GetSlice(o.Get("starbaseIntels"), GetFleetIntel)
	// ShipDesignIntels shipDesignIntels Slice ignore: false
	obj.ShipDesignIntels = GetSlice(o.Get("shipDesignIntels"), GetShipDesignIntel)
	// MineralPacketIntels mineralPacketIntels Slice ignore: false
	obj.MineralPacketIntels = GetSlice(o.Get("mineralPacketIntels"), GetMineralPacketIntel)
	// MineFieldIntels mineFieldIntels Slice ignore: false
	obj.MineFieldIntels = GetSlice(o.Get("mineFieldIntels"), GetMineFieldIntel)
	// WormholeIntels wormholeIntels Slice ignore: false
	obj.WormholeIntels = GetSlice(o.Get("wormholeIntels"), GetWormholeIntel)
	// MysteryTraderIntels mysteryTraderIntels Slice ignore: false
	obj.MysteryTraderIntels = GetSlice(o.Get("mysteryTraderIntels"), GetMysteryTraderIntel)
	// SalvageIntels salvageIntels Slice ignore: false
	obj.SalvageIntels = GetSlice(o.Get("salvageIntels"), GetSalvageIntel)
	return obj
}

func SetPlayerIntels(o js.Value, obj *cs.PlayerIntels) {
	// BattleRecords battleRecords Slice ignore: false
	o.Set("battleRecords", []any{})
	SetSlice(o.Get("battleRecords"), obj.BattleRecords, SetBattleRecord)
	// PlayerIntels playerIntels Slice ignore: false
	o.Set("playerIntels", []any{})
	SetSlice(o.Get("playerIntels"), obj.PlayerIntels, SetPlayerIntel)
	// ScoreIntels scoreIntels Slice ignore: false
	o.Set("scoreIntels", []any{})
	SetSlice(o.Get("scoreIntels"), obj.ScoreIntels, SetScoreIntel)
	// PlanetIntels planetIntels Slice ignore: false
	o.Set("planetIntels", []any{})
	SetSlice(o.Get("planetIntels"), obj.PlanetIntels, SetPlanetIntel)
	// FleetIntels fleetIntels Slice ignore: false
	o.Set("fleetIntels", []any{})
	SetSlice(o.Get("fleetIntels"), obj.FleetIntels, SetFleetIntel)
	// StarbaseIntels starbaseIntels Slice ignore: false
	o.Set("starbaseIntels", []any{})
	SetSlice(o.Get("starbaseIntels"), obj.StarbaseIntels, SetFleetIntel)
	// ShipDesignIntels shipDesignIntels Slice ignore: false
	o.Set("shipDesignIntels", []any{})
	SetSlice(o.Get("shipDesignIntels"), obj.ShipDesignIntels, SetShipDesignIntel)
	// MineralPacketIntels mineralPacketIntels Slice ignore: false
	o.Set("mineralPacketIntels", []any{})
	SetSlice(o.Get("mineralPacketIntels"), obj.MineralPacketIntels, SetMineralPacketIntel)
	// MineFieldIntels mineFieldIntels Slice ignore: false
	o.Set("mineFieldIntels", []any{})
	SetSlice(o.Get("mineFieldIntels"), obj.MineFieldIntels, SetMineFieldIntel)
	// WormholeIntels wormholeIntels Slice ignore: false
	o.Set("wormholeIntels", []any{})
	SetSlice(o.Get("wormholeIntels"), obj.WormholeIntels, SetWormholeIntel)
	// MysteryTraderIntels mysteryTraderIntels Slice ignore: false
	o.Set("mysteryTraderIntels", []any{})
	SetSlice(o.Get("mysteryTraderIntels"), obj.MysteryTraderIntels, SetMysteryTraderIntel)
	// SalvageIntels salvageIntels Slice ignore: false
	o.Set("salvageIntels", []any{})
	SetSlice(o.Get("salvageIntels"), obj.SalvageIntels, SetSalvageIntel)
}

func GetPlayerMessage(o js.Value) cs.PlayerMessage {
	obj := cs.PlayerMessage{}
	// Target  Object ignore: true
	// Type type Named ignore: false
	obj.Type = cs.PlayerMessageType(GetInt[cs.PlayerMessageType](o, "type"))
	// Text text BasicString ignore: false
	obj.Text = string(GetString(o, "text"))
	// BattleNum battleNum BasicInt ignore: false
	obj.BattleNum = GetInt[int](o, "battleNum")
	// Spec spec Object ignore: false
	obj.Spec = GetPlayerMessageSpec(o.Get("spec"))
	return obj
}

func SetPlayerMessage(o js.Value, obj *cs.PlayerMessage) {
	// Target  Object ignore: true
	// Type type Named ignore: false
	o.Set("type", int(obj.Type))
	// Text text BasicString ignore: false
	o.Set("text", obj.Text)
	// BattleNum battleNum BasicInt ignore: false
	o.Set("battleNum", obj.BattleNum)
	// Spec spec Object ignore: false
	o.Set("spec", map[string]any{})
	SetPlayerMessageSpec(o.Get("spec"), &obj.Spec)
}

func GetPlayerMessageSpec(o js.Value) cs.PlayerMessageSpec {
	obj := cs.PlayerMessageSpec{}
	// Target  Object ignore: true
	// Amount amount BasicInt ignore: false
	obj.Amount = GetInt[int](o, "amount")
	// Amount2 amount2 BasicInt ignore: false
	obj.Amount2 = GetInt[int](o, "amount2")
	// PrevAmount prevAmount BasicInt ignore: false
	obj.PrevAmount = GetInt[int](o, "prevAmount")
	// SourcePlayerNum sourcePlayerNum BasicInt ignore: false
	obj.SourcePlayerNum = GetInt[int](o, "sourcePlayerNum")
	// DestPlayerNum destPlayerNum BasicInt ignore: false
	obj.DestPlayerNum = GetInt[int](o, "destPlayerNum")
	// Name name BasicString ignore: false
	obj.Name = string(GetString(o, "name"))
	// Cost cost Object ignore: false
	costVal := o.Get("cost")
	if !costVal.IsUndefined() {
		cost := GetCost(costVal)
		obj.Cost = &cost
	}
	// Mineral mineral Object ignore: false
	mineralVal := o.Get("mineral")
	if !mineralVal.IsUndefined() {
		mineral := GetMineral(mineralVal)
		obj.Mineral = &mineral
	}
	// Cargo cargo Object ignore: false
	cargoVal := o.Get("cargo")
	if !cargoVal.IsUndefined() {
		cargo := GetCargo(cargoVal)
		obj.Cargo = &cargo
	}
	// QueueItemType queueItemType Named ignore: false
	obj.QueueItemType = cs.QueueItemType(GetString(o, "queueItemType"))
	// Field field Named ignore: false
	obj.Field = cs.TechField(GetString(o, "field"))
	// NextField nextField Named ignore: false
	obj.NextField = cs.TechField(GetString(o, "nextField"))
	// TechGained techGained BasicString ignore: false
	obj.TechGained = string(GetString(o, "techGained"))
	// LostTargetType lostTargetType Named ignore: false
	obj.LostTargetType = cs.MapObjectType(GetString(o, "lostTargetType"))
	// Battle battle Object ignore: false
	obj.Battle = GetBattleRecordStats(o.Get("battle"))
	// Comet comet Object ignore: false
	cometVal := o.Get("comet")
	if !cometVal.IsUndefined() {
		comet := GetPlayerMessageSpecComet(cometVal)
		obj.Comet = &comet
	}
	// Bombing bombing Object ignore: false
	bombingVal := o.Get("bombing")
	if !bombingVal.IsUndefined() {
		bombing := GetBombingResult(bombingVal)
		obj.Bombing = &bombing
	}
	// MineralPacketDamage mineralPacketDamage Object ignore: false
	mineralPacketDamageVal := o.Get("mineralPacketDamage")
	if !mineralPacketDamageVal.IsUndefined() {
		mineralPacketDamage := GetMineralPacketDamage(mineralPacketDamageVal)
		obj.MineralPacketDamage = &mineralPacketDamage
	}
	// MysteryTrader mysteryTrader Object ignore: false
	mysteryTraderVal := o.Get("mysteryTrader")
	if !mysteryTraderVal.IsUndefined() {
		mysteryTrader := GetPlayerMessageSpecMysteryTrader(mysteryTraderVal)
		obj.MysteryTrader = &mysteryTrader
	}
	return obj
}

func SetPlayerMessageSpec(o js.Value, obj *cs.PlayerMessageSpec) {
	// Target  Object ignore: true
	// Amount amount BasicInt ignore: false
	o.Set("amount", obj.Amount)
	// Amount2 amount2 BasicInt ignore: false
	o.Set("amount2", obj.Amount2)
	// PrevAmount prevAmount BasicInt ignore: false
	o.Set("prevAmount", obj.PrevAmount)
	// SourcePlayerNum sourcePlayerNum BasicInt ignore: false
	o.Set("sourcePlayerNum", obj.SourcePlayerNum)
	// DestPlayerNum destPlayerNum BasicInt ignore: false
	o.Set("destPlayerNum", obj.DestPlayerNum)
	// Name name BasicString ignore: false
	o.Set("name", obj.Name)
	// Cost cost Object ignore: false
	o.Set("cost", map[string]any{})
	SetCost(o.Get("cost"), obj.Cost)
	// Mineral mineral Object ignore: false
	o.Set("mineral", map[string]any{})
	SetMineral(o.Get("mineral"), obj.Mineral)
	// Cargo cargo Object ignore: false
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), obj.Cargo)
	// QueueItemType queueItemType Named ignore: false
	o.Set("queueItemType", string(obj.QueueItemType))
	// Field field Named ignore: false
	o.Set("field", string(obj.Field))
	// NextField nextField Named ignore: false
	o.Set("nextField", string(obj.NextField))
	// TechGained techGained BasicString ignore: false
	o.Set("techGained", obj.TechGained)
	// LostTargetType lostTargetType Named ignore: false
	o.Set("lostTargetType", string(obj.LostTargetType))
	// Battle battle Object ignore: false
	o.Set("battle", map[string]any{})
	SetBattleRecordStats(o.Get("battle"), &obj.Battle)
	// Comet comet Object ignore: false
	o.Set("comet", map[string]any{})
	SetPlayerMessageSpecComet(o.Get("comet"), obj.Comet)
	// Bombing bombing Object ignore: false
	o.Set("bombing", map[string]any{})
	SetBombingResult(o.Get("bombing"), obj.Bombing)
	// MineralPacketDamage mineralPacketDamage Object ignore: false
	o.Set("mineralPacketDamage", map[string]any{})
	SetMineralPacketDamage(o.Get("mineralPacketDamage"), obj.MineralPacketDamage)
	// MysteryTrader mysteryTrader Object ignore: false
	o.Set("mysteryTrader", map[string]any{})
	SetPlayerMessageSpecMysteryTrader(o.Get("mysteryTrader"), obj.MysteryTrader)
}

func GetPlayerMessageSpecComet(o js.Value) cs.PlayerMessageSpecComet {
	obj := cs.PlayerMessageSpecComet{}
	// Size size Named ignore: false
	obj.Size = cs.CometSize(GetString(o, "size"))
	// MineralsAdded mineralsAdded Object ignore: false
	obj.MineralsAdded = GetMineral(o.Get("mineralsAdded"))
	// MineralConcentrationIncreased mineralConcentrationIncreased Object ignore: false
	obj.MineralConcentrationIncreased = GetMineral(o.Get("mineralConcentrationIncreased"))
	// HabChanged habChanged Object ignore: false
	obj.HabChanged = GetHab(o.Get("habChanged"))
	// ColonistsKilled colonistsKilled BasicInt ignore: false
	obj.ColonistsKilled = GetInt[int](o, "colonistsKilled")
	return obj
}

func SetPlayerMessageSpecComet(o js.Value, obj *cs.PlayerMessageSpecComet) {
	// Size size Named ignore: false
	o.Set("size", string(obj.Size))
	// MineralsAdded mineralsAdded Object ignore: false
	o.Set("mineralsAdded", map[string]any{})
	SetMineral(o.Get("mineralsAdded"), &obj.MineralsAdded)
	// MineralConcentrationIncreased mineralConcentrationIncreased Object ignore: false
	o.Set("mineralConcentrationIncreased", map[string]any{})
	SetMineral(o.Get("mineralConcentrationIncreased"), &obj.MineralConcentrationIncreased)
	// HabChanged habChanged Object ignore: false
	o.Set("habChanged", map[string]any{})
	SetHab(o.Get("habChanged"), &obj.HabChanged)
	// ColonistsKilled colonistsKilled BasicInt ignore: false
	o.Set("colonistsKilled", obj.ColonistsKilled)
}

func GetPlayerMessageSpecMysteryTrader(o js.Value) cs.PlayerMessageSpecMysteryTrader {
	obj := cs.PlayerMessageSpecMysteryTrader{}
	// MysteryTraderReward  Object ignore: false
	obj.MysteryTraderReward = GetMysteryTraderReward(o)
	// FleetNum fleetNum BasicInt ignore: false
	obj.FleetNum = GetInt[int](o, "fleetNum")
	return obj
}

func SetPlayerMessageSpecMysteryTrader(o js.Value, obj *cs.PlayerMessageSpecMysteryTrader) {
	// MysteryTraderReward  Object ignore: false
	SetMysteryTraderReward(o, &obj.MysteryTraderReward)
	// FleetNum fleetNum BasicInt ignore: false
	o.Set("fleetNum", obj.FleetNum)
}

func GetPlayerOrders(o js.Value) cs.PlayerOrders {
	obj := cs.PlayerOrders{}
	// Researching researching Named ignore: false
	obj.Researching = cs.TechField(GetString(o, "researching"))
	// NextResearchField nextResearchField Named ignore: false
	obj.NextResearchField = cs.NextResearchField(GetString(o, "nextResearchField"))
	// ResearchAmount researchAmount BasicInt ignore: false
	obj.ResearchAmount = GetInt[int](o, "researchAmount")
	return obj
}

func SetPlayerOrders(o js.Value, obj *cs.PlayerOrders) {
	// Researching researching Named ignore: false
	o.Set("researching", string(obj.Researching))
	// NextResearchField nextResearchField Named ignore: false
	o.Set("nextResearchField", string(obj.NextResearchField))
	// ResearchAmount researchAmount BasicInt ignore: false
	o.Set("researchAmount", obj.ResearchAmount)
}

func GetPlayerPlans(o js.Value) cs.PlayerPlans {
	obj := cs.PlayerPlans{}
	// ProductionPlans productionPlans Slice ignore: false
	obj.ProductionPlans = GetSlice(o.Get("productionPlans"), GetProductionPlan)
	// BattlePlans battlePlans Slice ignore: false
	obj.BattlePlans = GetSlice(o.Get("battlePlans"), GetBattlePlan)
	// TransportPlans transportPlans Slice ignore: false
	obj.TransportPlans = GetSlice(o.Get("transportPlans"), GetTransportPlan)
	return obj
}

func SetPlayerPlans(o js.Value, obj *cs.PlayerPlans) {
	// ProductionPlans productionPlans Slice ignore: false
	o.Set("productionPlans", []any{})
	SetSlice(o.Get("productionPlans"), obj.ProductionPlans, SetProductionPlan)
	// BattlePlans battlePlans Slice ignore: false
	o.Set("battlePlans", []any{})
	SetSlice(o.Get("battlePlans"), obj.BattlePlans, SetBattlePlan)
	// TransportPlans transportPlans Slice ignore: false
	o.Set("transportPlans", []any{})
	SetSlice(o.Get("transportPlans"), obj.TransportPlans, SetTransportPlan)
}

func GetPlayerRelationship(o js.Value) cs.PlayerRelationship {
	obj := cs.PlayerRelationship{}
	// Relation relation Named ignore: false
	obj.Relation = cs.PlayerRelation(GetString(o, "relation"))
	// ShareMap shareMap BasicBool ignore: false
	obj.ShareMap = bool(GetBool(o, "shareMap"))
	return obj
}

func SetPlayerRelationship(o js.Value, obj *cs.PlayerRelationship) {
	// Relation relation Named ignore: false
	o.Set("relation", string(obj.Relation))
	// ShareMap shareMap BasicBool ignore: false
	o.Set("shareMap", obj.ShareMap)
}

func GetPlayerScore(o js.Value) cs.PlayerScore {
	obj := cs.PlayerScore{}
	// Planets planets BasicInt ignore: false
	obj.Planets = GetInt[int](o, "planets")
	// Starbases starbases BasicInt ignore: false
	obj.Starbases = GetInt[int](o, "starbases")
	// UnarmedShips unarmedShips BasicInt ignore: false
	obj.UnarmedShips = GetInt[int](o, "unarmedShips")
	// EscortShips escortShips BasicInt ignore: false
	obj.EscortShips = GetInt[int](o, "escortShips")
	// CapitalShips capitalShips BasicInt ignore: false
	obj.CapitalShips = GetInt[int](o, "capitalShips")
	// TechLevels techLevels BasicInt ignore: false
	obj.TechLevels = GetInt[int](o, "techLevels")
	// Resources resources BasicInt ignore: false
	obj.Resources = GetInt[int](o, "resources")
	// Score score BasicInt ignore: false
	obj.Score = GetInt[int](o, "score")
	// Rank rank BasicInt ignore: false
	obj.Rank = GetInt[int](o, "rank")
	// AchievedVictoryConditions achievedVictoryConditions Named ignore: false
	obj.AchievedVictoryConditions = cs.Bitmask(GetInt[cs.Bitmask](o, "achievedVictoryConditions"))
	return obj
}

func SetPlayerScore(o js.Value, obj *cs.PlayerScore) {
	// Planets planets BasicInt ignore: false
	o.Set("planets", obj.Planets)
	// Starbases starbases BasicInt ignore: false
	o.Set("starbases", obj.Starbases)
	// UnarmedShips unarmedShips BasicInt ignore: false
	o.Set("unarmedShips", obj.UnarmedShips)
	// EscortShips escortShips BasicInt ignore: false
	o.Set("escortShips", obj.EscortShips)
	// CapitalShips capitalShips BasicInt ignore: false
	o.Set("capitalShips", obj.CapitalShips)
	// TechLevels techLevels BasicInt ignore: false
	o.Set("techLevels", obj.TechLevels)
	// Resources resources BasicInt ignore: false
	o.Set("resources", obj.Resources)
	// Score score BasicInt ignore: false
	o.Set("score", obj.Score)
	// Rank rank BasicInt ignore: false
	o.Set("rank", obj.Rank)
	// AchievedVictoryConditions achievedVictoryConditions Named ignore: false
	o.Set("achievedVictoryConditions", uint32(obj.AchievedVictoryConditions))
}

func GetPlayerSpec(o js.Value) cs.PlayerSpec {
	obj := cs.PlayerSpec{}
	// PlanetaryScanner planetaryScanner Object ignore: false
	obj.PlanetaryScanner = GetTechPlanetaryScanner(o.Get("planetaryScanner"))
	// Defense defense Object ignore: false
	obj.Defense = GetTechDefense(o.Get("defense"))
	// Terraform terraform Map ignore: true
	// ResourcesPerYear resourcesPerYear BasicInt ignore: false
	obj.ResourcesPerYear = GetInt[int](o, "resourcesPerYear")
	// ResourcesPerYearResearch resourcesPerYearResearch BasicInt ignore: false
	obj.ResourcesPerYearResearch = GetInt[int](o, "resourcesPerYearResearch")
	// ResourcesPerYearResearchEstimated resourcesPerYearResearchEstimated BasicInt ignore: false
	obj.ResourcesPerYearResearchEstimated = GetInt[int](o, "resourcesPerYearResearchEstimated")
	// CurrentResearchCost currentResearchCost BasicInt ignore: false
	obj.CurrentResearchCost = GetInt[int](o, "currentResearchCost")
	return obj
}

func SetPlayerSpec(o js.Value, obj *cs.PlayerSpec) {
	// PlanetaryScanner planetaryScanner Object ignore: false
	o.Set("planetaryScanner", map[string]any{})
	SetTechPlanetaryScanner(o.Get("planetaryScanner"), &obj.PlanetaryScanner)
	// Defense defense Object ignore: false
	o.Set("defense", map[string]any{})
	SetTechDefense(o.Get("defense"), &obj.Defense)
	// Terraform terraform Map ignore: true
	// ResourcesPerYear resourcesPerYear BasicInt ignore: false
	o.Set("resourcesPerYear", obj.ResourcesPerYear)
	// ResourcesPerYearResearch resourcesPerYearResearch BasicInt ignore: false
	o.Set("resourcesPerYearResearch", obj.ResourcesPerYearResearch)
	// ResourcesPerYearResearchEstimated resourcesPerYearResearchEstimated BasicInt ignore: false
	o.Set("resourcesPerYearResearchEstimated", obj.ResourcesPerYearResearchEstimated)
	// CurrentResearchCost currentResearchCost BasicInt ignore: false
	o.Set("currentResearchCost", obj.CurrentResearchCost)
}

func GetPlayerStats(o js.Value) cs.PlayerStats {
	obj := cs.PlayerStats{}
	// FleetsBuilt fleetsBuilt BasicInt ignore: false
	obj.FleetsBuilt = GetInt[int](o, "fleetsBuilt")
	// StarbasesBuilt starbasesBuilt BasicInt ignore: false
	obj.StarbasesBuilt = GetInt[int](o, "starbasesBuilt")
	// TokensBuilt tokensBuilt BasicInt ignore: false
	obj.TokensBuilt = GetInt[int](o, "tokensBuilt")
	// PlanetsColonized planetsColonized BasicInt ignore: false
	obj.PlanetsColonized = GetInt[int](o, "planetsColonized")
	return obj
}

func SetPlayerStats(o js.Value, obj *cs.PlayerStats) {
	// FleetsBuilt fleetsBuilt BasicInt ignore: false
	o.Set("fleetsBuilt", obj.FleetsBuilt)
	// StarbasesBuilt starbasesBuilt BasicInt ignore: false
	o.Set("starbasesBuilt", obj.StarbasesBuilt)
	// TokensBuilt tokensBuilt BasicInt ignore: false
	o.Set("tokensBuilt", obj.TokensBuilt)
	// PlanetsColonized planetsColonized BasicInt ignore: false
	o.Set("planetsColonized", obj.PlanetsColonized)
}

func GetProductionPlan(o js.Value) cs.ProductionPlan {
	obj := cs.ProductionPlan{}
	// Num num BasicInt ignore: false
	obj.Num = GetInt[int](o, "num")
	// Name name BasicString ignore: false
	obj.Name = string(GetString(o, "name"))
	// Items items Slice ignore: true
	// ContributesOnlyLeftoverToResearch contributesOnlyLeftoverToResearch BasicBool ignore: false
	obj.ContributesOnlyLeftoverToResearch = bool(GetBool(o, "contributesOnlyLeftoverToResearch"))
	return obj
}

func SetProductionPlan(o js.Value, obj *cs.ProductionPlan) {
	// Num num BasicInt ignore: false
	o.Set("num", obj.Num)
	// Name name BasicString ignore: false
	o.Set("name", obj.Name)
	// Items items Slice ignore: true
	// ContributesOnlyLeftoverToResearch contributesOnlyLeftoverToResearch BasicBool ignore: false
	o.Set("contributesOnlyLeftoverToResearch", obj.ContributesOnlyLeftoverToResearch)
}

func GetProductionQueueItem(o js.Value) cs.ProductionQueueItem {
	obj := cs.ProductionQueueItem{}
	// QueueItemCompletionEstimate  Object ignore: false
	obj.QueueItemCompletionEstimate = GetQueueItemCompletionEstimate(o)
	// Type type Named ignore: false
	obj.Type = cs.QueueItemType(GetString(o, "type"))
	// DesignNum designNum BasicInt ignore: false
	obj.DesignNum = GetInt[int](o, "designNum")
	// Quantity quantity BasicInt ignore: false
	obj.Quantity = GetInt[int](o, "quantity")
	// Allocated allocated Object ignore: false
	obj.Allocated = GetCost(o.Get("allocated"))
	// Tags tags Named ignore: false
	// unknown type Tags cs.Tags map[string]string
	// index   ignore: true
	// design   ignore: true
	return obj
}

func SetProductionQueueItem(o js.Value, obj *cs.ProductionQueueItem) {
	// QueueItemCompletionEstimate  Object ignore: false
	SetQueueItemCompletionEstimate(o, &obj.QueueItemCompletionEstimate)
	// Type type Named ignore: false
	o.Set("type", string(obj.Type))
	// DesignNum designNum BasicInt ignore: false
	o.Set("designNum", obj.DesignNum)
	// Quantity quantity BasicInt ignore: false
	o.Set("quantity", obj.Quantity)
	// Allocated allocated Object ignore: false
	o.Set("allocated", map[string]any{})
	SetCost(o.Get("allocated"), &obj.Allocated)
	// Tags tags Named ignore: false
	// unknown type Tags cs.Tags map[string]string
	// index   ignore: true
	// design   ignore: true
}

func GetQueueItemCompletionEstimate(o js.Value) cs.QueueItemCompletionEstimate {
	obj := cs.QueueItemCompletionEstimate{}
	// Skipped skipped BasicBool ignore: false
	obj.Skipped = bool(GetBool(o, "skipped"))
	// YearsToBuildOne yearsToBuildOne BasicInt ignore: false
	obj.YearsToBuildOne = GetInt[int](o, "yearsToBuildOne")
	// YearsToBuildAll yearsToBuildAll BasicInt ignore: false
	obj.YearsToBuildAll = GetInt[int](o, "yearsToBuildAll")
	// YearsToSkipAuto yearsToSkipAuto BasicInt ignore: false
	obj.YearsToSkipAuto = GetInt[int](o, "yearsToSkipAuto")
	return obj
}

func SetQueueItemCompletionEstimate(o js.Value, obj *cs.QueueItemCompletionEstimate) {
	// Skipped skipped BasicBool ignore: false
	o.Set("skipped", obj.Skipped)
	// YearsToBuildOne yearsToBuildOne BasicInt ignore: false
	o.Set("yearsToBuildOne", obj.YearsToBuildOne)
	// YearsToBuildAll yearsToBuildAll BasicInt ignore: false
	o.Set("yearsToBuildAll", obj.YearsToBuildAll)
	// YearsToSkipAuto yearsToSkipAuto BasicInt ignore: false
	o.Set("yearsToSkipAuto", obj.YearsToSkipAuto)
}

func GetRace(o js.Value) cs.Race {
	obj := cs.Race{}
	// DBObject  Object ignore: false
	obj.DBObject = GetDBObject(o)
	// UserID userId BasicInt ignore: false
	obj.UserID = GetInt[int64](o, "userId")
	// Name name BasicString ignore: false
	obj.Name = string(GetString(o, "name"))
	// PluralName pluralName BasicString ignore: false
	obj.PluralName = string(GetString(o, "pluralName"))
	// SpendLeftoverPointsOn spendLeftoverPointsOn Named ignore: false
	obj.SpendLeftoverPointsOn = cs.SpendLeftoverPointsOn(GetString(o, "spendLeftoverPointsOn"))
	// PRT prt Named ignore: false
	obj.PRT = cs.PRT(GetString(o, "prt"))
	// LRTs lrts Named ignore: false
	obj.LRTs = cs.Bitmask(GetInt[cs.Bitmask](o, "lrts"))
	// HabLow habLow Object ignore: false
	obj.HabLow = GetHab(o.Get("habLow"))
	// HabHigh habHigh Object ignore: false
	obj.HabHigh = GetHab(o.Get("habHigh"))
	// GrowthRate growthRate BasicInt ignore: false
	obj.GrowthRate = GetInt[int](o, "growthRate")
	// PopEfficiency popEfficiency BasicInt ignore: false
	obj.PopEfficiency = GetInt[int](o, "popEfficiency")
	// FactoryOutput factoryOutput BasicInt ignore: false
	obj.FactoryOutput = GetInt[int](o, "factoryOutput")
	// FactoryCost factoryCost BasicInt ignore: false
	obj.FactoryCost = GetInt[int](o, "factoryCost")
	// NumFactories numFactories BasicInt ignore: false
	obj.NumFactories = GetInt[int](o, "numFactories")
	// FactoriesCostLess factoriesCostLess BasicBool ignore: false
	obj.FactoriesCostLess = bool(GetBool(o, "factoriesCostLess"))
	// ImmuneGrav immuneGrav BasicBool ignore: false
	obj.ImmuneGrav = bool(GetBool(o, "immuneGrav"))
	// ImmuneTemp immuneTemp BasicBool ignore: false
	obj.ImmuneTemp = bool(GetBool(o, "immuneTemp"))
	// ImmuneRad immuneRad BasicBool ignore: false
	obj.ImmuneRad = bool(GetBool(o, "immuneRad"))
	// MineOutput mineOutput BasicInt ignore: false
	obj.MineOutput = GetInt[int](o, "mineOutput")
	// MineCost mineCost BasicInt ignore: false
	obj.MineCost = GetInt[int](o, "mineCost")
	// NumMines numMines BasicInt ignore: false
	obj.NumMines = GetInt[int](o, "numMines")
	// ResearchCost researchCost Object ignore: false
	obj.ResearchCost = GetResearchCost(o.Get("researchCost"))
	// TechsStartHigh techsStartHigh BasicBool ignore: false
	obj.TechsStartHigh = bool(GetBool(o, "techsStartHigh"))
	// Spec spec Object ignore: true
	return obj
}

func SetRace(o js.Value, obj *cs.Race) {
	// DBObject  Object ignore: false
	SetDBObject(o, &obj.DBObject)
	// UserID userId BasicInt ignore: false
	o.Set("userId", obj.UserID)
	// Name name BasicString ignore: false
	o.Set("name", obj.Name)
	// PluralName pluralName BasicString ignore: false
	o.Set("pluralName", obj.PluralName)
	// SpendLeftoverPointsOn spendLeftoverPointsOn Named ignore: false
	o.Set("spendLeftoverPointsOn", string(obj.SpendLeftoverPointsOn))
	// PRT prt Named ignore: false
	o.Set("prt", string(obj.PRT))
	// LRTs lrts Named ignore: false
	o.Set("lrts", uint32(obj.LRTs))
	// HabLow habLow Object ignore: false
	o.Set("habLow", map[string]any{})
	SetHab(o.Get("habLow"), &obj.HabLow)
	// HabHigh habHigh Object ignore: false
	o.Set("habHigh", map[string]any{})
	SetHab(o.Get("habHigh"), &obj.HabHigh)
	// GrowthRate growthRate BasicInt ignore: false
	o.Set("growthRate", obj.GrowthRate)
	// PopEfficiency popEfficiency BasicInt ignore: false
	o.Set("popEfficiency", obj.PopEfficiency)
	// FactoryOutput factoryOutput BasicInt ignore: false
	o.Set("factoryOutput", obj.FactoryOutput)
	// FactoryCost factoryCost BasicInt ignore: false
	o.Set("factoryCost", obj.FactoryCost)
	// NumFactories numFactories BasicInt ignore: false
	o.Set("numFactories", obj.NumFactories)
	// FactoriesCostLess factoriesCostLess BasicBool ignore: false
	o.Set("factoriesCostLess", obj.FactoriesCostLess)
	// ImmuneGrav immuneGrav BasicBool ignore: false
	o.Set("immuneGrav", obj.ImmuneGrav)
	// ImmuneTemp immuneTemp BasicBool ignore: false
	o.Set("immuneTemp", obj.ImmuneTemp)
	// ImmuneRad immuneRad BasicBool ignore: false
	o.Set("immuneRad", obj.ImmuneRad)
	// MineOutput mineOutput BasicInt ignore: false
	o.Set("mineOutput", obj.MineOutput)
	// MineCost mineCost BasicInt ignore: false
	o.Set("mineCost", obj.MineCost)
	// NumMines numMines BasicInt ignore: false
	o.Set("numMines", obj.NumMines)
	// ResearchCost researchCost Object ignore: false
	o.Set("researchCost", map[string]any{})
	SetResearchCost(o.Get("researchCost"), &obj.ResearchCost)
	// TechsStartHigh techsStartHigh BasicBool ignore: false
	o.Set("techsStartHigh", obj.TechsStartHigh)
	// Spec spec Object ignore: true
}

func GetResearchCost(o js.Value) cs.ResearchCost {
	obj := cs.ResearchCost{}
	// Energy energy Named ignore: false
	obj.Energy = cs.ResearchCostLevel(GetString(o, "energy"))
	// Weapons weapons Named ignore: false
	obj.Weapons = cs.ResearchCostLevel(GetString(o, "weapons"))
	// Propulsion propulsion Named ignore: false
	obj.Propulsion = cs.ResearchCostLevel(GetString(o, "propulsion"))
	// Construction construction Named ignore: false
	obj.Construction = cs.ResearchCostLevel(GetString(o, "construction"))
	// Electronics electronics Named ignore: false
	obj.Electronics = cs.ResearchCostLevel(GetString(o, "electronics"))
	// Biotechnology biotechnology Named ignore: false
	obj.Biotechnology = cs.ResearchCostLevel(GetString(o, "biotechnology"))
	return obj
}

func SetResearchCost(o js.Value, obj *cs.ResearchCost) {
	// Energy energy Named ignore: false
	o.Set("energy", string(obj.Energy))
	// Weapons weapons Named ignore: false
	o.Set("weapons", string(obj.Weapons))
	// Propulsion propulsion Named ignore: false
	o.Set("propulsion", string(obj.Propulsion))
	// Construction construction Named ignore: false
	o.Set("construction", string(obj.Construction))
	// Electronics electronics Named ignore: false
	o.Set("electronics", string(obj.Electronics))
	// Biotechnology biotechnology Named ignore: false
	o.Set("biotechnology", string(obj.Biotechnology))
}

func GetSalvageIntel(o js.Value) cs.SalvageIntel {
	obj := cs.SalvageIntel{}
	// MapObjectIntel  Object ignore: false
	obj.MapObjectIntel = GetMapObjectIntel(o)
	// Cargo cargo Object ignore: false
	obj.Cargo = GetCargo(o.Get("cargo"))
	return obj
}

func SetSalvageIntel(o js.Value, obj *cs.SalvageIntel) {
	// MapObjectIntel  Object ignore: false
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	// Cargo cargo Object ignore: false
	o.Set("cargo", map[string]any{})
	SetCargo(o.Get("cargo"), &obj.Cargo)
}

func GetScoreIntel(o js.Value) cs.ScoreIntel {
	obj := cs.ScoreIntel{}
	// ScoreHistory scoreHistory Slice ignore: false
	obj.ScoreHistory = GetSlice(o.Get("scoreHistory"), GetPlayerScore)
	return obj
}

func SetScoreIntel(o js.Value, obj *cs.ScoreIntel) {
	// ScoreHistory scoreHistory Slice ignore: false
	o.Set("scoreHistory", []any{})
	SetSlice(o.Get("scoreHistory"), obj.ScoreHistory, SetPlayerScore)
}

func GetShipDesign(o js.Value) cs.ShipDesign {
	obj := cs.ShipDesign{}
	// GameDBObject  Object ignore: false
	obj.GameDBObject = GetGameDBObject(o)
	// Num num BasicInt ignore: false
	obj.Num = GetInt[int](o, "num")
	// PlayerNum playerNum BasicInt ignore: false
	obj.PlayerNum = GetInt[int](o, "playerNum")
	// OriginalPlayerNum originalPlayerNum BasicInt ignore: false
	obj.OriginalPlayerNum = GetInt[int](o, "originalPlayerNum")
	// Name name BasicString ignore: false
	obj.Name = string(GetString(o, "name"))
	// Version version BasicInt ignore: false
	obj.Version = GetInt[int](o, "version")
	// Hull hull BasicString ignore: false
	obj.Hull = string(GetString(o, "hull"))
	// HullSetNumber hullSetNumber BasicInt ignore: false
	obj.HullSetNumber = GetInt[int](o, "hullSetNumber")
	// CannotDelete cannotDelete BasicBool ignore: false
	obj.CannotDelete = bool(GetBool(o, "cannotDelete"))
	// MysteryTrader mysteryTrader BasicBool ignore: false
	obj.MysteryTrader = bool(GetBool(o, "mysteryTrader"))
	// Slots slots Slice ignore: false
	obj.Slots = GetSlice(o.Get("slots"), GetShipDesignSlot)
	// Purpose purpose Named ignore: false
	obj.Purpose = cs.ShipDesignPurpose(GetString(o, "purpose"))
	// Spec spec Object ignore: false
	obj.Spec = GetShipDesignSpec(o.Get("spec"))
	// Delete  BasicBool ignore: false
	obj.Delete = bool(GetBool(o, ""))
	return obj
}

func SetShipDesign(o js.Value, obj *cs.ShipDesign) {
	// GameDBObject  Object ignore: false
	SetGameDBObject(o, &obj.GameDBObject)
	// Num num BasicInt ignore: false
	o.Set("num", obj.Num)
	// PlayerNum playerNum BasicInt ignore: false
	o.Set("playerNum", obj.PlayerNum)
	// OriginalPlayerNum originalPlayerNum BasicInt ignore: false
	o.Set("originalPlayerNum", obj.OriginalPlayerNum)
	// Name name BasicString ignore: false
	o.Set("name", obj.Name)
	// Version version BasicInt ignore: false
	o.Set("version", obj.Version)
	// Hull hull BasicString ignore: false
	o.Set("hull", obj.Hull)
	// HullSetNumber hullSetNumber BasicInt ignore: false
	o.Set("hullSetNumber", obj.HullSetNumber)
	// CannotDelete cannotDelete BasicBool ignore: false
	o.Set("cannotDelete", obj.CannotDelete)
	// MysteryTrader mysteryTrader BasicBool ignore: false
	o.Set("mysteryTrader", obj.MysteryTrader)
	// Slots slots Slice ignore: false
	o.Set("slots", []any{})
	SetSlice(o.Get("slots"), obj.Slots, SetShipDesignSlot)
	// Purpose purpose Named ignore: false
	o.Set("purpose", string(obj.Purpose))
	// Spec spec Object ignore: false
	o.Set("spec", map[string]any{})
	SetShipDesignSpec(o.Get("spec"), &obj.Spec)
	// Delete  BasicBool ignore: false
	o.Set("", obj.Delete)
}

func GetShipDesignIntel(o js.Value) cs.ShipDesignIntel {
	obj := cs.ShipDesignIntel{}
	// Intel  Object ignore: false
	obj.Intel = GetIntel(o)
	// Hull hull BasicString ignore: false
	obj.Hull = string(GetString(o, "hull"))
	// HullSetNumber hullSetNumber BasicInt ignore: false
	obj.HullSetNumber = GetInt[int](o, "hullSetNumber")
	// Version version BasicInt ignore: false
	obj.Version = GetInt[int](o, "version")
	// Slots slots Slice ignore: false
	obj.Slots = GetSlice(o.Get("slots"), GetShipDesignSlot)
	// Spec spec Object ignore: false
	obj.Spec = GetShipDesignSpec(o.Get("spec"))
	return obj
}

func SetShipDesignIntel(o js.Value, obj *cs.ShipDesignIntel) {
	// Intel  Object ignore: false
	SetIntel(o, &obj.Intel)
	// Hull hull BasicString ignore: false
	o.Set("hull", obj.Hull)
	// HullSetNumber hullSetNumber BasicInt ignore: false
	o.Set("hullSetNumber", obj.HullSetNumber)
	// Version version BasicInt ignore: false
	o.Set("version", obj.Version)
	// Slots slots Slice ignore: false
	o.Set("slots", []any{})
	SetSlice(o.Get("slots"), obj.Slots, SetShipDesignSlot)
	// Spec spec Object ignore: false
	o.Set("spec", map[string]any{})
	SetShipDesignSpec(o.Get("spec"), &obj.Spec)
}

func GetShipDesignSlot(o js.Value) cs.ShipDesignSlot {
	obj := cs.ShipDesignSlot{}
	// HullComponent hullComponent BasicString ignore: false
	obj.HullComponent = string(GetString(o, "hullComponent"))
	// HullSlotIndex hullSlotIndex BasicInt ignore: false
	obj.HullSlotIndex = GetInt[int](o, "hullSlotIndex")
	// Quantity quantity BasicInt ignore: false
	obj.Quantity = GetInt[int](o, "quantity")
	return obj
}

func SetShipDesignSlot(o js.Value, obj *cs.ShipDesignSlot) {
	// HullComponent hullComponent BasicString ignore: false
	o.Set("hullComponent", obj.HullComponent)
	// HullSlotIndex hullSlotIndex BasicInt ignore: false
	o.Set("hullSlotIndex", obj.HullSlotIndex)
	// Quantity quantity BasicInt ignore: false
	o.Set("quantity", obj.Quantity)
}

func GetShipDesignSpec(o js.Value) cs.ShipDesignSpec {
	obj := cs.ShipDesignSpec{}
	// AdditionalMassDrivers additionalMassDrivers BasicInt ignore: false
	obj.AdditionalMassDrivers = GetInt[int](o, "additionalMassDrivers")
	// Armor armor BasicInt ignore: false
	obj.Armor = GetInt[int](o, "armor")
	// BasePacketSpeed basePacketSpeed BasicInt ignore: false
	obj.BasePacketSpeed = GetInt[int](o, "basePacketSpeed")
	// BeamBonus beamBonus BasicFloat ignore: false
	obj.BeamBonus = GetFloat[float64](o, "beamBonus")
	// BeamDefense beamDefense BasicFloat ignore: false
	obj.BeamDefense = GetFloat[float64](o, "beamDefense")
	// Bomber bomber BasicBool ignore: false
	obj.Bomber = bool(GetBool(o, "bomber"))
	// Bombs bombs Slice ignore: false
	obj.Bombs = GetSlice(o.Get("bombs"), GetBomb)
	// CanJump canJump BasicBool ignore: false
	obj.CanJump = bool(GetBool(o, "canJump"))
	// CanLayMines canLayMines BasicBool ignore: false
	obj.CanLayMines = bool(GetBool(o, "canLayMines"))
	// CanStealFleetCargo canStealFleetCargo BasicBool ignore: false
	obj.CanStealFleetCargo = bool(GetBool(o, "canStealFleetCargo"))
	// CanStealPlanetCargo canStealPlanetCargo BasicBool ignore: false
	obj.CanStealPlanetCargo = bool(GetBool(o, "canStealPlanetCargo"))
	// CargoCapacity cargoCapacity BasicInt ignore: false
	obj.CargoCapacity = GetInt[int](o, "cargoCapacity")
	// CloakPercent cloakPercent BasicInt ignore: false
	obj.CloakPercent = GetInt[int](o, "cloakPercent")
	// CloakPercentFullCargo cloakPercentFullCargo BasicInt ignore: false
	obj.CloakPercentFullCargo = GetInt[int](o, "cloakPercentFullCargo")
	// CloakUnits cloakUnits BasicInt ignore: false
	obj.CloakUnits = GetInt[int](o, "cloakUnits")
	// Colonizer colonizer BasicBool ignore: false
	obj.Colonizer = bool(GetBool(o, "colonizer"))
	// Cost cost Object ignore: false
	obj.Cost = GetCost(o.Get("cost"))
	// Engine engine Object ignore: false
	obj.Engine = GetEngine(o.Get("engine"))
	// EstimatedRange estimatedRange BasicInt ignore: false
	obj.EstimatedRange = GetInt[int](o, "estimatedRange")
	// EstimatedRangeFull estimatedRangeFull BasicInt ignore: false
	obj.EstimatedRangeFull = GetInt[int](o, "estimatedRangeFull")
	// FuelCapacity fuelCapacity BasicInt ignore: false
	obj.FuelCapacity = GetInt[int](o, "fuelCapacity")
	// FuelGeneration fuelGeneration BasicInt ignore: false
	obj.FuelGeneration = GetInt[int](o, "fuelGeneration")
	// HasWeapons hasWeapons BasicBool ignore: false
	obj.HasWeapons = bool(GetBool(o, "hasWeapons"))
	// HullType hullType Named ignore: false
	obj.HullType = cs.TechHullType(GetString(o, "hullType"))
	// ImmuneToOwnDetonation immuneToOwnDetonation BasicBool ignore: false
	obj.ImmuneToOwnDetonation = bool(GetBool(o, "immuneToOwnDetonation"))
	// Initiative initiative BasicInt ignore: false
	obj.Initiative = GetInt[int](o, "initiative")
	// InnateScanRangePenFactor innateScanRangePenFactor BasicFloat ignore: false
	obj.InnateScanRangePenFactor = GetFloat[float64](o, "innateScanRangePenFactor")
	// Mass mass BasicInt ignore: false
	obj.Mass = GetInt[int](o, "mass")
	// MassDriver massDriver BasicString ignore: false
	obj.MassDriver = string(GetString(o, "massDriver"))
	// MaxHullMass maxHullMass BasicInt ignore: false
	obj.MaxHullMass = GetInt[int](o, "maxHullMass")
	// MaxPopulation maxPopulation BasicInt ignore: false
	obj.MaxPopulation = GetInt[int](o, "maxPopulation")
	// MaxRange maxRange BasicInt ignore: false
	obj.MaxRange = GetInt[int](o, "maxRange")
	// MineLayingRateByMineType mineLayingRateByMineType Map ignore: true
	// MineSweep mineSweep BasicInt ignore: false
	obj.MineSweep = GetInt[int](o, "mineSweep")
	// MiningRate miningRate BasicInt ignore: false
	obj.MiningRate = GetInt[int](o, "miningRate")
	// Movement movement BasicInt ignore: false
	obj.Movement = GetInt[int](o, "movement")
	// MovementBonus movementBonus BasicInt ignore: false
	obj.MovementBonus = GetInt[int](o, "movementBonus")
	// MovementFull movementFull BasicInt ignore: false
	obj.MovementFull = GetInt[int](o, "movementFull")
	// NumBuilt numBuilt BasicInt ignore: false
	obj.NumBuilt = GetInt[int](o, "numBuilt")
	// NumEngines numEngines BasicInt ignore: false
	obj.NumEngines = GetInt[int](o, "numEngines")
	// NumInstances numInstances BasicInt ignore: false
	obj.NumInstances = GetInt[int](o, "numInstances")
	// OrbitalConstructionModule orbitalConstructionModule BasicBool ignore: false
	obj.OrbitalConstructionModule = bool(GetBool(o, "orbitalConstructionModule"))
	// PowerRating powerRating BasicInt ignore: false
	obj.PowerRating = GetInt[int](o, "powerRating")
	// Radiating radiating BasicBool ignore: false
	obj.Radiating = bool(GetBool(o, "radiating"))
	// ReduceCloaking reduceCloaking BasicFloat ignore: false
	obj.ReduceCloaking = GetFloat[float64](o, "reduceCloaking")
	// ReduceMovement reduceMovement BasicInt ignore: false
	obj.ReduceMovement = GetInt[int](o, "reduceMovement")
	// RepairBonus repairBonus BasicFloat ignore: false
	obj.RepairBonus = GetFloat[float64](o, "repairBonus")
	// RetroBombs retroBombs Slice ignore: false
	obj.RetroBombs = GetSlice(o.Get("retroBombs"), GetBomb)
	// SafeHullMass safeHullMass BasicInt ignore: false
	obj.SafeHullMass = GetInt[int](o, "safeHullMass")
	// SafePacketSpeed safePacketSpeed BasicInt ignore: false
	obj.SafePacketSpeed = GetInt[int](o, "safePacketSpeed")
	// SafeRange safeRange BasicInt ignore: false
	obj.SafeRange = GetInt[int](o, "safeRange")
	// Scanner scanner BasicBool ignore: false
	obj.Scanner = bool(GetBool(o, "scanner"))
	// ScanRange scanRange BasicInt ignore: false
	obj.ScanRange = GetInt[int](o, "scanRange")
	// ScanRangePen scanRangePen BasicInt ignore: false
	obj.ScanRangePen = GetInt[int](o, "scanRangePen")
	// Shields shields BasicInt ignore: false
	obj.Shields = GetInt[int](o, "shields")
	// SmartBombs smartBombs Slice ignore: false
	obj.SmartBombs = GetSlice(o.Get("smartBombs"), GetBomb)
	// SpaceDock spaceDock BasicInt ignore: false
	obj.SpaceDock = GetInt[int](o, "spaceDock")
	// Starbase starbase BasicBool ignore: false
	obj.Starbase = bool(GetBool(o, "starbase"))
	// Stargate stargate BasicString ignore: false
	obj.Stargate = string(GetString(o, "stargate"))
	// TechLevel techLevel Object ignore: false
	obj.TechLevel = GetTechLevel(o.Get("techLevel"))
	// TerraformRate terraformRate BasicInt ignore: false
	obj.TerraformRate = GetInt[int](o, "terraformRate")
	// TorpedoBonus torpedoBonus BasicFloat ignore: false
	obj.TorpedoBonus = GetFloat[float64](o, "torpedoBonus")
	// TorpedoJamming torpedoJamming BasicFloat ignore: false
	obj.TorpedoJamming = GetFloat[float64](o, "torpedoJamming")
	// WeaponSlots weaponSlots Slice ignore: false
	obj.WeaponSlots = GetSlice(o.Get("weaponSlots"), GetShipDesignSlot)
	return obj
}

func SetShipDesignSpec(o js.Value, obj *cs.ShipDesignSpec) {
	// AdditionalMassDrivers additionalMassDrivers BasicInt ignore: false
	o.Set("additionalMassDrivers", obj.AdditionalMassDrivers)
	// Armor armor BasicInt ignore: false
	o.Set("armor", obj.Armor)
	// BasePacketSpeed basePacketSpeed BasicInt ignore: false
	o.Set("basePacketSpeed", obj.BasePacketSpeed)
	// BeamBonus beamBonus BasicFloat ignore: false
	o.Set("beamBonus", obj.BeamBonus)
	// BeamDefense beamDefense BasicFloat ignore: false
	o.Set("beamDefense", obj.BeamDefense)
	// Bomber bomber BasicBool ignore: false
	o.Set("bomber", obj.Bomber)
	// Bombs bombs Slice ignore: false
	o.Set("bombs", []any{})
	SetSlice(o.Get("bombs"), obj.Bombs, SetBomb)
	// CanJump canJump BasicBool ignore: false
	o.Set("canJump", obj.CanJump)
	// CanLayMines canLayMines BasicBool ignore: false
	o.Set("canLayMines", obj.CanLayMines)
	// CanStealFleetCargo canStealFleetCargo BasicBool ignore: false
	o.Set("canStealFleetCargo", obj.CanStealFleetCargo)
	// CanStealPlanetCargo canStealPlanetCargo BasicBool ignore: false
	o.Set("canStealPlanetCargo", obj.CanStealPlanetCargo)
	// CargoCapacity cargoCapacity BasicInt ignore: false
	o.Set("cargoCapacity", obj.CargoCapacity)
	// CloakPercent cloakPercent BasicInt ignore: false
	o.Set("cloakPercent", obj.CloakPercent)
	// CloakPercentFullCargo cloakPercentFullCargo BasicInt ignore: false
	o.Set("cloakPercentFullCargo", obj.CloakPercentFullCargo)
	// CloakUnits cloakUnits BasicInt ignore: false
	o.Set("cloakUnits", obj.CloakUnits)
	// Colonizer colonizer BasicBool ignore: false
	o.Set("colonizer", obj.Colonizer)
	// Cost cost Object ignore: false
	o.Set("cost", map[string]any{})
	SetCost(o.Get("cost"), &obj.Cost)
	// Engine engine Object ignore: false
	o.Set("engine", map[string]any{})
	SetEngine(o.Get("engine"), &obj.Engine)
	// EstimatedRange estimatedRange BasicInt ignore: false
	o.Set("estimatedRange", obj.EstimatedRange)
	// EstimatedRangeFull estimatedRangeFull BasicInt ignore: false
	o.Set("estimatedRangeFull", obj.EstimatedRangeFull)
	// FuelCapacity fuelCapacity BasicInt ignore: false
	o.Set("fuelCapacity", obj.FuelCapacity)
	// FuelGeneration fuelGeneration BasicInt ignore: false
	o.Set("fuelGeneration", obj.FuelGeneration)
	// HasWeapons hasWeapons BasicBool ignore: false
	o.Set("hasWeapons", obj.HasWeapons)
	// HullType hullType Named ignore: false
	o.Set("hullType", string(obj.HullType))
	// ImmuneToOwnDetonation immuneToOwnDetonation BasicBool ignore: false
	o.Set("immuneToOwnDetonation", obj.ImmuneToOwnDetonation)
	// Initiative initiative BasicInt ignore: false
	o.Set("initiative", obj.Initiative)
	// InnateScanRangePenFactor innateScanRangePenFactor BasicFloat ignore: false
	o.Set("innateScanRangePenFactor", obj.InnateScanRangePenFactor)
	// Mass mass BasicInt ignore: false
	o.Set("mass", obj.Mass)
	// MassDriver massDriver BasicString ignore: false
	o.Set("massDriver", obj.MassDriver)
	// MaxHullMass maxHullMass BasicInt ignore: false
	o.Set("maxHullMass", obj.MaxHullMass)
	// MaxPopulation maxPopulation BasicInt ignore: false
	o.Set("maxPopulation", obj.MaxPopulation)
	// MaxRange maxRange BasicInt ignore: false
	o.Set("maxRange", obj.MaxRange)
	// MineLayingRateByMineType mineLayingRateByMineType Map ignore: true
	// MineSweep mineSweep BasicInt ignore: false
	o.Set("mineSweep", obj.MineSweep)
	// MiningRate miningRate BasicInt ignore: false
	o.Set("miningRate", obj.MiningRate)
	// Movement movement BasicInt ignore: false
	o.Set("movement", obj.Movement)
	// MovementBonus movementBonus BasicInt ignore: false
	o.Set("movementBonus", obj.MovementBonus)
	// MovementFull movementFull BasicInt ignore: false
	o.Set("movementFull", obj.MovementFull)
	// NumBuilt numBuilt BasicInt ignore: false
	o.Set("numBuilt", obj.NumBuilt)
	// NumEngines numEngines BasicInt ignore: false
	o.Set("numEngines", obj.NumEngines)
	// NumInstances numInstances BasicInt ignore: false
	o.Set("numInstances", obj.NumInstances)
	// OrbitalConstructionModule orbitalConstructionModule BasicBool ignore: false
	o.Set("orbitalConstructionModule", obj.OrbitalConstructionModule)
	// PowerRating powerRating BasicInt ignore: false
	o.Set("powerRating", obj.PowerRating)
	// Radiating radiating BasicBool ignore: false
	o.Set("radiating", obj.Radiating)
	// ReduceCloaking reduceCloaking BasicFloat ignore: false
	o.Set("reduceCloaking", obj.ReduceCloaking)
	// ReduceMovement reduceMovement BasicInt ignore: false
	o.Set("reduceMovement", obj.ReduceMovement)
	// RepairBonus repairBonus BasicFloat ignore: false
	o.Set("repairBonus", obj.RepairBonus)
	// RetroBombs retroBombs Slice ignore: false
	o.Set("retroBombs", []any{})
	SetSlice(o.Get("retroBombs"), obj.RetroBombs, SetBomb)
	// SafeHullMass safeHullMass BasicInt ignore: false
	o.Set("safeHullMass", obj.SafeHullMass)
	// SafePacketSpeed safePacketSpeed BasicInt ignore: false
	o.Set("safePacketSpeed", obj.SafePacketSpeed)
	// SafeRange safeRange BasicInt ignore: false
	o.Set("safeRange", obj.SafeRange)
	// Scanner scanner BasicBool ignore: false
	o.Set("scanner", obj.Scanner)
	// ScanRange scanRange BasicInt ignore: false
	o.Set("scanRange", obj.ScanRange)
	// ScanRangePen scanRangePen BasicInt ignore: false
	o.Set("scanRangePen", obj.ScanRangePen)
	// Shields shields BasicInt ignore: false
	o.Set("shields", obj.Shields)
	// SmartBombs smartBombs Slice ignore: false
	o.Set("smartBombs", []any{})
	SetSlice(o.Get("smartBombs"), obj.SmartBombs, SetBomb)
	// SpaceDock spaceDock BasicInt ignore: false
	o.Set("spaceDock", obj.SpaceDock)
	// Starbase starbase BasicBool ignore: false
	o.Set("starbase", obj.Starbase)
	// Stargate stargate BasicString ignore: false
	o.Set("stargate", obj.Stargate)
	// TechLevel techLevel Object ignore: false
	o.Set("techLevel", map[string]any{})
	SetTechLevel(o.Get("techLevel"), &obj.TechLevel)
	// TerraformRate terraformRate BasicInt ignore: false
	o.Set("terraformRate", obj.TerraformRate)
	// TorpedoBonus torpedoBonus BasicFloat ignore: false
	o.Set("torpedoBonus", obj.TorpedoBonus)
	// TorpedoJamming torpedoJamming BasicFloat ignore: false
	o.Set("torpedoJamming", obj.TorpedoJamming)
	// WeaponSlots weaponSlots Slice ignore: false
	o.Set("weaponSlots", []any{})
	SetSlice(o.Get("weaponSlots"), obj.WeaponSlots, SetShipDesignSlot)
}

func GetShipToken(o js.Value) cs.ShipToken {
	obj := cs.ShipToken{}
	// DesignNum designNum BasicInt ignore: false
	obj.DesignNum = GetInt[int](o, "designNum")
	// Quantity quantity BasicInt ignore: false
	obj.Quantity = GetInt[int](o, "quantity")
	// Damage damage BasicFloat ignore: false
	obj.Damage = GetFloat[float64](o, "damage")
	// QuantityDamaged quantityDamaged BasicInt ignore: false
	obj.QuantityDamaged = GetInt[int](o, "quantityDamaged")
	// design   ignore: true
	return obj
}

func SetShipToken(o js.Value, obj *cs.ShipToken) {
	// DesignNum designNum BasicInt ignore: false
	o.Set("designNum", obj.DesignNum)
	// Quantity quantity BasicInt ignore: false
	o.Set("quantity", obj.Quantity)
	// Damage damage BasicFloat ignore: false
	o.Set("damage", obj.Damage)
	// QuantityDamaged quantityDamaged BasicInt ignore: false
	o.Set("quantityDamaged", obj.QuantityDamaged)
	// design   ignore: true
}

func GetTech(o js.Value) cs.Tech {
	obj := cs.Tech{}
	// Name name BasicString ignore: false
	obj.Name = string(GetString(o, "name"))
	// Cost cost Object ignore: false
	obj.Cost = GetCost(o.Get("cost"))
	// Requirements requirements Object ignore: true
	// Ranking ranking BasicInt ignore: false
	obj.Ranking = GetInt[int](o, "ranking")
	// Category category Named ignore: false
	obj.Category = cs.TechCategory(GetString(o, "category"))
	// Origin origin BasicString ignore: false
	obj.Origin = string(GetString(o, "origin"))
	return obj
}

func SetTech(o js.Value, obj *cs.Tech) {
	// Name name BasicString ignore: false
	o.Set("name", obj.Name)
	// Cost cost Object ignore: false
	o.Set("cost", map[string]any{})
	SetCost(o.Get("cost"), &obj.Cost)
	// Requirements requirements Object ignore: true
	// Ranking ranking BasicInt ignore: false
	o.Set("ranking", obj.Ranking)
	// Category category Named ignore: false
	o.Set("category", string(obj.Category))
	// Origin origin BasicString ignore: false
	o.Set("origin", obj.Origin)
}

func GetTechDefense(o js.Value) cs.TechDefense {
	obj := cs.TechDefense{}
	// TechPlanetary  Object ignore: false
	obj.TechPlanetary = GetTechPlanetary(o)
	// Defense  Object ignore: false
	obj.Defense = GetDefense(o)
	return obj
}

func SetTechDefense(o js.Value, obj *cs.TechDefense) {
	// TechPlanetary  Object ignore: false
	SetTechPlanetary(o, &obj.TechPlanetary)
	// Defense  Object ignore: false
	SetDefense(o, &obj.Defense)
}

func GetTechLevel(o js.Value) cs.TechLevel {
	obj := cs.TechLevel{}
	// Energy energy BasicInt ignore: false
	obj.Energy = GetInt[int](o, "energy")
	// Weapons weapons BasicInt ignore: false
	obj.Weapons = GetInt[int](o, "weapons")
	// Propulsion propulsion BasicInt ignore: false
	obj.Propulsion = GetInt[int](o, "propulsion")
	// Construction construction BasicInt ignore: false
	obj.Construction = GetInt[int](o, "construction")
	// Electronics electronics BasicInt ignore: false
	obj.Electronics = GetInt[int](o, "electronics")
	// Biotechnology biotechnology BasicInt ignore: false
	obj.Biotechnology = GetInt[int](o, "biotechnology")
	return obj
}

func SetTechLevel(o js.Value, obj *cs.TechLevel) {
	// Energy energy BasicInt ignore: false
	o.Set("energy", obj.Energy)
	// Weapons weapons BasicInt ignore: false
	o.Set("weapons", obj.Weapons)
	// Propulsion propulsion BasicInt ignore: false
	o.Set("propulsion", obj.Propulsion)
	// Construction construction BasicInt ignore: false
	o.Set("construction", obj.Construction)
	// Electronics electronics BasicInt ignore: false
	o.Set("electronics", obj.Electronics)
	// Biotechnology biotechnology BasicInt ignore: false
	o.Set("biotechnology", obj.Biotechnology)
}

func GetTechPlanetary(o js.Value) cs.TechPlanetary {
	obj := cs.TechPlanetary{}
	// Tech  Object ignore: false
	obj.Tech = GetTech(o)
	// ResetPlanet resetPlanet BasicBool ignore: false
	obj.ResetPlanet = bool(GetBool(o, "resetPlanet"))
	return obj
}

func SetTechPlanetary(o js.Value, obj *cs.TechPlanetary) {
	// Tech  Object ignore: false
	SetTech(o, &obj.Tech)
	// ResetPlanet resetPlanet BasicBool ignore: false
	o.Set("resetPlanet", obj.ResetPlanet)
}

func GetTechPlanetaryScanner(o js.Value) cs.TechPlanetaryScanner {
	obj := cs.TechPlanetaryScanner{}
	// TechPlanetary  Object ignore: false
	obj.TechPlanetary = GetTechPlanetary(o)
	// ScanRange scanRange BasicInt ignore: false
	obj.ScanRange = GetInt[int](o, "scanRange")
	// ScanRangePen scanRangePen BasicInt ignore: false
	obj.ScanRangePen = GetInt[int](o, "scanRangePen")
	return obj
}

func SetTechPlanetaryScanner(o js.Value, obj *cs.TechPlanetaryScanner) {
	// TechPlanetary  Object ignore: false
	SetTechPlanetary(o, &obj.TechPlanetary)
	// ScanRange scanRange BasicInt ignore: false
	o.Set("scanRange", obj.ScanRange)
	// ScanRangePen scanRangePen BasicInt ignore: false
	o.Set("scanRangePen", obj.ScanRangePen)
}

func GetTransportPlan(o js.Value) cs.TransportPlan {
	obj := cs.TransportPlan{}
	// Num num BasicInt ignore: false
	obj.Num = GetInt[int](o, "num")
	// Name name BasicString ignore: false
	obj.Name = string(GetString(o, "name"))
	// Tasks tasks Object ignore: false
	obj.Tasks = GetWaypointTransportTasks(o.Get("tasks"))
	return obj
}

func SetTransportPlan(o js.Value, obj *cs.TransportPlan) {
	// Num num BasicInt ignore: false
	o.Set("num", obj.Num)
	// Name name BasicString ignore: false
	o.Set("name", obj.Name)
	// Tasks tasks Object ignore: false
	o.Set("tasks", map[string]any{})
	SetWaypointTransportTasks(o.Get("tasks"), &obj.Tasks)
}

func GetVector(o js.Value) cs.Vector {
	obj := cs.Vector{}
	// X x BasicFloat ignore: false
	obj.X = GetFloat[float64](o, "x")
	// Y y BasicFloat ignore: false
	obj.Y = GetFloat[float64](o, "y")
	return obj
}

func SetVector(o js.Value, obj *cs.Vector) {
	// X x BasicFloat ignore: false
	o.Set("x", obj.X)
	// Y y BasicFloat ignore: false
	o.Set("y", obj.Y)
}

func GetWaypoint(o js.Value) cs.Waypoint {
	obj := cs.Waypoint{}
	// Position position Object ignore: false
	obj.Position = GetVector(o.Get("position"))
	// WarpSpeed warpSpeed BasicInt ignore: false
	obj.WarpSpeed = GetInt[int](o, "warpSpeed")
	// EstFuelUsage estFuelUsage BasicInt ignore: false
	obj.EstFuelUsage = GetInt[int](o, "estFuelUsage")
	// Task task Named ignore: false
	obj.Task = cs.WaypointTask(GetString(o, "task"))
	// TransportTasks transportTasks Object ignore: false
	obj.TransportTasks = GetWaypointTransportTasks(o.Get("transportTasks"))
	// WaitAtWaypoint waitAtWaypoint BasicBool ignore: false
	obj.WaitAtWaypoint = bool(GetBool(o, "waitAtWaypoint"))
	// LayMineFieldDuration layMineFieldDuration BasicInt ignore: false
	obj.LayMineFieldDuration = GetInt[int](o, "layMineFieldDuration")
	// PatrolRange patrolRange BasicInt ignore: false
	obj.PatrolRange = GetInt[int](o, "patrolRange")
	// PatrolWarpSpeed patrolWarpSpeed BasicInt ignore: false
	obj.PatrolWarpSpeed = GetInt[int](o, "patrolWarpSpeed")
	// TargetType targetType Named ignore: false
	obj.TargetType = cs.MapObjectType(GetString(o, "targetType"))
	// TargetNum targetNum BasicInt ignore: false
	obj.TargetNum = GetInt[int](o, "targetNum")
	// TargetPlayerNum targetPlayerNum BasicInt ignore: false
	obj.TargetPlayerNum = GetInt[int](o, "targetPlayerNum")
	// TargetName targetName BasicString ignore: false
	obj.TargetName = string(GetString(o, "targetName"))
	// TransferToPlayer transferToPlayer BasicInt ignore: false
	obj.TransferToPlayer = GetInt[int](o, "transferToPlayer")
	// PartiallyComplete partiallyComplete BasicBool ignore: false
	obj.PartiallyComplete = bool(GetBool(o, "partiallyComplete"))
	// processed   ignore: true
	return obj
}

func SetWaypoint(o js.Value, obj *cs.Waypoint) {
	// Position position Object ignore: false
	o.Set("position", map[string]any{})
	SetVector(o.Get("position"), &obj.Position)
	// WarpSpeed warpSpeed BasicInt ignore: false
	o.Set("warpSpeed", obj.WarpSpeed)
	// EstFuelUsage estFuelUsage BasicInt ignore: false
	o.Set("estFuelUsage", obj.EstFuelUsage)
	// Task task Named ignore: false
	o.Set("task", string(obj.Task))
	// TransportTasks transportTasks Object ignore: false
	o.Set("transportTasks", map[string]any{})
	SetWaypointTransportTasks(o.Get("transportTasks"), &obj.TransportTasks)
	// WaitAtWaypoint waitAtWaypoint BasicBool ignore: false
	o.Set("waitAtWaypoint", obj.WaitAtWaypoint)
	// LayMineFieldDuration layMineFieldDuration BasicInt ignore: false
	o.Set("layMineFieldDuration", obj.LayMineFieldDuration)
	// PatrolRange patrolRange BasicInt ignore: false
	o.Set("patrolRange", obj.PatrolRange)
	// PatrolWarpSpeed patrolWarpSpeed BasicInt ignore: false
	o.Set("patrolWarpSpeed", obj.PatrolWarpSpeed)
	// TargetType targetType Named ignore: false
	o.Set("targetType", string(obj.TargetType))
	// TargetNum targetNum BasicInt ignore: false
	o.Set("targetNum", obj.TargetNum)
	// TargetPlayerNum targetPlayerNum BasicInt ignore: false
	o.Set("targetPlayerNum", obj.TargetPlayerNum)
	// TargetName targetName BasicString ignore: false
	o.Set("targetName", obj.TargetName)
	// TransferToPlayer transferToPlayer BasicInt ignore: false
	o.Set("transferToPlayer", obj.TransferToPlayer)
	// PartiallyComplete partiallyComplete BasicBool ignore: false
	o.Set("partiallyComplete", obj.PartiallyComplete)
	// processed   ignore: true
}

func GetWaypointTransportTask(o js.Value) cs.WaypointTransportTask {
	obj := cs.WaypointTransportTask{}
	// Amount amount BasicInt ignore: false
	obj.Amount = GetInt[int](o, "amount")
	// Action action Named ignore: false
	obj.Action = cs.WaypointTaskTransportAction(GetString(o, "action"))
	return obj
}

func SetWaypointTransportTask(o js.Value, obj *cs.WaypointTransportTask) {
	// Amount amount BasicInt ignore: false
	o.Set("amount", obj.Amount)
	// Action action Named ignore: false
	o.Set("action", string(obj.Action))
}

func GetWaypointTransportTasks(o js.Value) cs.WaypointTransportTasks {
	obj := cs.WaypointTransportTasks{}
	// Fuel fuel Object ignore: false
	obj.Fuel = GetWaypointTransportTask(o.Get("fuel"))
	// Ironium ironium Object ignore: false
	obj.Ironium = GetWaypointTransportTask(o.Get("ironium"))
	// Boranium boranium Object ignore: false
	obj.Boranium = GetWaypointTransportTask(o.Get("boranium"))
	// Germanium germanium Object ignore: false
	obj.Germanium = GetWaypointTransportTask(o.Get("germanium"))
	// Colonists colonists Object ignore: false
	obj.Colonists = GetWaypointTransportTask(o.Get("colonists"))
	return obj
}

func SetWaypointTransportTasks(o js.Value, obj *cs.WaypointTransportTasks) {
	// Fuel fuel Object ignore: false
	o.Set("fuel", map[string]any{})
	SetWaypointTransportTask(o.Get("fuel"), &obj.Fuel)
	// Ironium ironium Object ignore: false
	o.Set("ironium", map[string]any{})
	SetWaypointTransportTask(o.Get("ironium"), &obj.Ironium)
	// Boranium boranium Object ignore: false
	o.Set("boranium", map[string]any{})
	SetWaypointTransportTask(o.Get("boranium"), &obj.Boranium)
	// Germanium germanium Object ignore: false
	o.Set("germanium", map[string]any{})
	SetWaypointTransportTask(o.Get("germanium"), &obj.Germanium)
	// Colonists colonists Object ignore: false
	o.Set("colonists", map[string]any{})
	SetWaypointTransportTask(o.Get("colonists"), &obj.Colonists)
}

func GetWormholeIntel(o js.Value) cs.WormholeIntel {
	obj := cs.WormholeIntel{}
	// MapObjectIntel  Object ignore: false
	obj.MapObjectIntel = GetMapObjectIntel(o)
	// DestinationNum destinationNum BasicInt ignore: false
	obj.DestinationNum = GetInt[int](o, "destinationNum")
	// Stability stability Named ignore: false
	obj.Stability = cs.WormholeStability(GetString(o, "stability"))
	return obj
}

func SetWormholeIntel(o js.Value, obj *cs.WormholeIntel) {
	// MapObjectIntel  Object ignore: false
	SetMapObjectIntel(o, &obj.MapObjectIntel)
	// DestinationNum destinationNum BasicInt ignore: false
	o.Set("destinationNum", obj.DestinationNum)
	// Stability stability Named ignore: false
	o.Set("stability", string(obj.Stability))
}
