//go:build !wasi && !wasm

package cs

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func simulationSides() []BattleSimulationSide {
	sides := []BattleSimulationSide{}
	for range 2 {
		sides = append(sides, BattleSimulationSide{Race: Humanoids(), TechLevels: TechLevel{26, 26, 26, 26, 26, 26}, Stacks: []BattleSimulationStack{{Design: DesignStalwartDefender, Quantity: 10, Plan: BattlePlan{PrimaryTarget: BattleTargetAny, SecondaryTarget: BattleTargetAny, Tactic: BattleTacticMaximizeDamage}}}})
	}
	return sides
}
func TestSimulateBattle(t *testing.T) {
	t.Run("reproducible and independent of game state", func(t *testing.T) {
		rules := NewRulesWithSeed(123)
		control := NewRulesWithSeed(123)
		sides := simulationSides()
		slots := append([]ShipDesignSlot(nil), sides[0].Stacks[0].Design.Slots...)
		first, err := SimulateBattle(rules, sides, 42)
		require.NoError(t, err)
		second, err := SimulateBattle(rules, sides, 42)
		require.NoError(t, err)
		require.Equal(t, first.Record.Tokens, second.Record.Tokens)
		require.Equal(t, first.Record.ActionsPerRound, second.Record.ActionsPerRound)
		require.Equal(t, first.Stacks, second.Stacks)
		require.Equal(t, control.random.Float64(), rules.random.Float64())
		require.Equal(t, 10, sides[0].Stacks[0].Quantity)
		require.Equal(t, slots, sides[0].Stacks[0].Design.Slots)
		require.Len(t, first.Stacks, 2)
		for _, s := range first.Stacks {
			require.Equal(t, 10, s.Survivors+s.Destroyed)
		}
	})
	t.Run("same design stays in separate stacks with separate plans", func(t *testing.T) {
		sides := simulationSides()
		split := sides[0].Stacks[0]
		split.Plan.Tactic = BattleTacticDisengage
		split.QuantityDamaged, split.Damage = 3, 10
		sides[0].Stacks = append(sides[0].Stacks, split)
		result, err := SimulateBattle(NewRulesWithSeed(1), sides, 1)
		require.NoError(t, err)
		require.Len(t, result.Record.Tokens, 3)
		found := false
		for _, token := range result.Record.Tokens {
			if token.PlayerNum == 1 && token.DesignNum == 2 {
				found = true
				require.Equal(t, BattleTacticDisengage, token.Tactic)
				require.Equal(t, 3, token.StartingQuantityDamaged)
				require.Equal(t, 10, token.StartingDamage)
			}
		}
		require.True(t, found)
	})
	t.Run("reject invalid input", func(t *testing.T) {
		cases := []struct {
			name   string
			change func([]BattleSimulationSide)
		}{
			{"empty side", func(s []BattleSimulationSide) { s[0].Stacks = nil }},
			{"quantity", func(s []BattleSimulationSide) { s[0].Stacks[0].Quantity = 0 }},
			{"damaged quantity", func(s []BattleSimulationSide) { s[0].Stacks[0].QuantityDamaged = 11 }},
			{"damage", func(s []BattleSimulationSide) { s[0].Stacks[0].QuantityDamaged = 1; s[0].Stacks[0].Damage = 100000 }},
			{"tech", func(s []BattleSimulationSide) { s[0].TechLevels.Energy = 27 }},
			{"hull", func(s []BattleSimulationSide) { s[0].Stacks[0].Design.Hull = "invalid" }},
			{"no hostile orders", func(s []BattleSimulationSide) {
				s[0].Stacks[0].Plan.PrimaryTarget = BattleTargetNone
				s[1].Stacks[0].Plan.PrimaryTarget = BattleTargetNone
			}},
		}
		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				sides := simulationSides()
				tt.change(sides)
				_, err := SimulateBattle(NewRules(), sides, 1)
				require.Error(t, err)
			})
		}
	})

	t.Run("existing ships can use technology their owner cannot build", func(t *testing.T) {
		sides := simulationSides()
		sides[0].TechLevels = TechLevel{}
		_, err := SimulateBattle(NewRules(), sides, 1)
		require.NoError(t, err)
	})
	t.Run("preserves rules and recomputes supplied specs", func(t *testing.T) {
		sides := simulationSides()
		sides[0].Stacks[0].Design.Spec.Armor = 999999
		result, err := SimulateBattle(NewRules(), sides, 7)
		require.NoError(t, err)
		require.NotEqual(t, 999999, result.Designs[0].Spec.Armor)
		require.Equal(t, 999999, sides[0].Stacks[0].Design.Spec.Armor)
	})
}
