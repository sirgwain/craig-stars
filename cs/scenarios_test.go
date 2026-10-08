//go:build !wasi && !wasm

package cs

import (
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBattleScenarios(t *testing.T) {
	tests := []struct {
		scenario TestScenario
		players  int
		fire     []BattleRecordTokenActionType
	}{
		{ScenarioBattleTorpedoes(), 2, []BattleRecordTokenActionType{TokenActionTorpedoFire}},
		{ScenarioBattleBeams(), 2, []BattleRecordTokenActionType{TokenActionBeamFire}},
		{ScenarioBattleStarbases(), 2, []BattleRecordTokenActionType{TokenActionBeamFire, TokenActionTorpedoFire}},
		{ScenarioBattleDumpCargo(), 2, []BattleRecordTokenActionType{TokenActionTorpedoFire}},
		{ScenarioBattleRunAway(), 2, []BattleRecordTokenActionType{TokenActionTorpedoFire, TokenActionRanAway}},
		{ScenarioBattleDisengageIfChallenged(), 2, []BattleRecordTokenActionType{TokenActionBeamFire, TokenActionTorpedoFire, TokenActionRanAway}},
		{ScenarioBattle3Players(), 3, []BattleRecordTokenActionType{TokenActionBeamFire, TokenActionTorpedoFire}},
		{ScenarioBattle16Players(), 16, []BattleRecordTokenActionType{TokenActionBeamFire}},
		{ScenarioBattle20Players(), 20, []BattleRecordTokenActionType{TokenActionBeamFire}},
		{ScenarioBattleMissilesAndChaff(), 2, []BattleRecordTokenActionType{TokenActionBeamFire, TokenActionTorpedoFire}},
		{ScenarioBattleShieldSappers(), 2, []BattleRecordTokenActionType{TokenActionBeamFire}},
		{ScenarioBattleDamagedStacks(), 2, []BattleRecordTokenActionType{TokenActionBeamFire, TokenActionTorpedoFire}},
		{ScenarioBattleAlliedSupport(), 3, []BattleRecordTokenActionType{TokenActionBeamFire, TokenActionTorpedoFire}},
		{ScenarioBattleMixedFleet(), 2, []BattleRecordTokenActionType{TokenActionBeamFire, TokenActionTorpedoFire}},
	}
	for _, tt := range tests {
		t.Run(tt.scenario.Name, func(t *testing.T) {
			u := newTestUniverse(t, tt.scenario)
			u.Run((*turnGenerator).fleetBattle)
			require.Len(t, u.Player(1).BattleRecords, 1)
			record := u.Player(1).BattleRecords[0]
			assert.Equal(t, tt.players, record.Stats.NumPlayers)
			assert.LessOrEqual(t, len(record.ActionsPerRound), u.Game.Rules.NumBattleRounds+1)
			actions := slices.Concat(record.ActionsPerRound...)
			for _, kind := range tt.fire {
				assert.True(t, slices.ContainsFunc(actions, func(a BattleRecordTokenAction) bool { return a.Type == kind }), "missing %s", kind)
			}
			for _, token := range record.Tokens {
				assert.GreaterOrEqual(t, token.Position.X, 0)
				assert.Less(t, token.Position.X, battleWidth)
				assert.GreaterOrEqual(t, token.Position.Y, 0)
				assert.Less(t, token.Position.Y, battleHeight)
			}
		})
	}
}

func TestBattleScenarioLargeFormations(t *testing.T) {
	for _, s := range []TestScenario{ScenarioBattle16Players(), ScenarioBattle20Players()} {
		t.Run(s.Name, func(t *testing.T) {
			u := newTestUniverse(t, s)
			colors := map[string]bool{}
			for _, player := range u.Game.Players {
				assert.False(t, colors[player.Color], "player %d reuses a color", player.Num)
				colors[player.Color] = true
			}
			u.Run((*turnGenerator).fleetBattle)
			record := u.Player(1).BattleRecords[0]
			positions := map[Vector]bool{}
			for _, token := range record.Tokens {
				assert.Equal(t, battleStartingPositions[16][(token.PlayerNum-1)%16], token.Position)
				positions[token.Position] = true
			}
			assert.Len(t, record.Tokens, len(s.Players))
			assert.Len(t, positions, 16)
			for _, player := range u.Game.Players {
				assert.Len(t, player.BattleRecords, 1)
			}
		})
	}
}

func TestBattleScenarioDumpCargo(t *testing.T) {
	s := ScenarioBattleDumpCargo()
	u := newTestUniverse(t, s)
	dump, keep := u.Fleet("Dump Cargo"), u.Fleet("Keep Cargo")
	// Inspect preparation directly so casualties cannot hide what was dumped.
	b := newBattler(testLogger, &u.Game.Rules, 1, map[int]*Player{1: u.Player(1), 2: u.Player(2)}, u.Game.Fleets, nil).(*battle)
	b.prepareBattle()
	assert.Equal(t, Cargo{Colonists: 150}, dump.Cargo)
	assert.Equal(t, s.Players[0].Fleets[1].Cargo, keep.Cargo)
	assert.Equal(t, Mineral{Ironium: 150, Boranium: 150, Germanium: 150}, b.record.dumpedMinerals)
	for _, token := range b.tokens {
		if token.fleet == dump {
			cargoPerShip := getCargoPerShip(dump.Cargo.Total(), dump.Spec.CargoCapacity, token.design.Spec.CargoCapacity)
			assert.Equal(t, token.design.Spec.Mass+cargoPerShip, token.Mass)
			assert.Equal(t, max(u.Game.Rules.MovementMin, token.design.getMovement(&u.Game.Rules, cargoPerShip)-1), token.Movement)
		}
	}
	// Overrides and fleet plan pointers must be independent across builds.
	other := BuildScenario(s)
	dump.battlePlan.DumpCargo = false
	assert.True(t, other.Fleets[0].battlePlan.DumpCargo)
	assert.True(t, s.Players[0].BattlePlans[0].DumpCargo)
}

func TestBattleScenarioAlliedSupport(t *testing.T) {
	u := newTestUniverse(t, ScenarioBattleAlliedSupport())
	observer := u.Game.Fleets[3]
	originalTokens := slices.Clone(observer.Tokens)
	u.Run((*turnGenerator).fleetBattle)
	record := u.Player(1).BattleRecords[0]
	playersByToken := map[int]int{}
	for _, token := range record.Tokens {
		assert.NotEqual(t, 4, token.PlayerNum)
		playersByToken[token.Num] = token.PlayerNum
	}
	assert.Equal(t, originalTokens, observer.Tokens)
	assert.False(t, observer.Delete)
	allyFired := false
	defenderFired := false
	for _, round := range record.ActionsPerRound {
		for _, action := range round {
			if action.Type != TokenActionBeamFire && action.Type != TokenActionTorpedoFire {
				continue
			}
			if playersByToken[action.TokenNum] == 3 {
				allyFired = true
				assert.Equal(t, 2, playersByToken[action.TargetNum])
			}
			if playersByToken[action.TokenNum] == 2 {
				defenderFired = true
				assert.Equal(t, 1, playersByToken[action.TargetNum])
			}
		}
	}
	assert.True(t, allyFired)
	assert.True(t, defenderFired)
}

func TestBattleScenarioWeaponBehavior(t *testing.T) {
	t.Run("torpedo volleys", func(t *testing.T) {
		u := newTestUniverse(t, ScenarioBattleTorpedoes())
		u.Run((*turnGenerator).fleetBattle)
		hits, misses, shields, armor := 0, 0, 0, 0
		partialKill := false
		for _, action := range slices.Concat(u.Player(1).BattleRecords[0].ActionsPerRound...) {
			if action.Type == TokenActionTorpedoFire {
				hits += action.TorpedoHits
				misses += action.TorpedoMisses
				shields += action.DamageDoneShields
				armor += action.DamageDoneArmor
				partialKill = partialKill || action.TokensDestroyed > 0 && action.Target.Quantity > 0
			}
		}
		assert.Positive(t, hits)
		assert.Positive(t, misses)
		assert.Positive(t, shields)
		assert.Positive(t, armor)
		assert.True(t, partialKill)
	})
	t.Run("sappers only damage shields", func(t *testing.T) {
		u := newTestUniverse(t, ScenarioBattleShieldSappers())
		u.Run((*turnGenerator).fleetBattle)
		record := u.Player(1).BattleRecords[0]
		attacker := slices.IndexFunc(record.Tokens, func(token BattleRecordToken) bool { return token.PlayerNum == 1 })
		require.NotEqual(t, -1, attacker)
		shieldDamage := 0
		for _, action := range slices.Concat(record.ActionsPerRound...) {
			if action.Type == TokenActionBeamFire && action.TokenNum == record.Tokens[attacker].Num && action.Slot == 2 {
				shieldDamage += action.DamageDoneShields
				assert.Zero(t, action.DamageDoneArmor)
			}
		}
		assert.Positive(t, shieldDamage)
	})
	t.Run("retreating ship fires before escaping", func(t *testing.T) {
		u := newTestUniverse(t, ScenarioBattleRunAway())
		u.Run((*turnGenerator).fleetBattle)
		record := u.Player(1).BattleRecords[0]
		retreater := slices.IndexFunc(record.Tokens, func(token BattleRecordToken) bool { return token.PlayerNum == 1 })
		require.NotEqual(t, -1, retreater)
		fired, escaped := false, false
		for _, action := range slices.Concat(record.ActionsPerRound...) {
			if action.TokenNum != record.Tokens[retreater].Num {
				continue
			}
			if action.Type == TokenActionTorpedoFire {
				assert.False(t, escaped)
				fired = true
			}
			if action.Type == TokenActionRanAway {
				assert.True(t, fired)
				escaped = true
			}
		}
		assert.True(t, escaped)
	})
}

func TestBattleScenarioDisengageIfChallenged(t *testing.T) {
	u := newTestUniverse(t, ScenarioBattleDisengageIfChallenged())
	u.Run((*turnGenerator).fleetBattle)
	record := u.Player(1).BattleRecords[0]
	challenged := slices.IndexFunc(record.Tokens, func(token BattleRecordToken) bool { return token.PlayerNum == 1 })
	require.NotEqual(t, -1, challenged)
	tokenNum := record.Tokens[challenged].Num

	// the token attacks until an enemy damages its armor, then retreats
	// for the full countdown before leaving the board
	hitRound, retreatMoves := 0, 0
	escaped := false
	for _, action := range slices.Concat(record.ActionsPerRound...) {
		if action.TargetNum == tokenNum && action.DamageDoneArmor > 0 && hitRound == 0 {
			hitRound = action.Round
		}
		if action.TokenNum != tokenNum {
			continue
		}
		switch action.Type {
		case TokenActionMove:
			if hitRound > 0 && action.Round > hitRound {
				retreatMoves++
			}
		case TokenActionRanAway:
			assert.Positive(t, hitRound, "ran away before its armor was damaged")
			escaped = true
		}
	}
	assert.True(t, escaped)
	assert.Equal(t, u.Game.Rules.MovesToRunAway, retreatMoves)
}

func TestBattleScenarioMixedFleet(t *testing.T) {
	u := newTestUniverse(t, ScenarioBattleMixedFleet())
	u.Run((*turnGenerator).fleetBattle)
	record := u.Player(1).BattleRecords[0]
	assert.Less(t, 5, len(record.ActionsPerRound))

	// damage is fractional during battle, but surviving ships keep whole points
	for _, fleet := range u.Game.Fleets {
		for _, token := range fleet.Tokens {
			assert.Equal(t, math.Floor(token.Damage), token.Damage, "%s %s", fleet.Name, token.design.Name)
		}
	}
}
