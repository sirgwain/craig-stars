package cs

import (
	"fmt"
	"io"
	"log/slog"
	"math"
	"slices"
)

// BattleSimulationSide describes a participant without any live game objects.
type BattleSimulationSide struct {
	Race          Race
	TechLevels    TechLevel
	AcquiredTechs map[string]bool
	Stacks        []BattleSimulationStack
}

type BattleSimulationStack struct {
	Design          ShipDesign
	Quantity        int
	QuantityDamaged int
	Damage          float64
	Plan            BattlePlan
}

type BattleSimulationStackResult struct {
	PlayerNum int
	DesignNum int
	Survivors int
	Retreated int
	Destroyed int
	Loss      Cost
}

type BattleSimulationResult struct {
	Record  *BattleRecord
	Designs []*ShipDesign
	Stacks  []BattleSimulationStackResult
}

// SimulateBattle uses the normal battle engine on private copies of every input.
// Each stack is wrapped in a private fleet solely to supply its battle plan.
// The copied rules have an independent RNG, leaving the game's random stream intact.
func SimulateBattle(rules Rules, sides []BattleSimulationSide, seed int64) (*BattleSimulationResult, error) {
	if len(sides) != 2 {
		return nil, fmt.Errorf("a simulation needs exactly two sides")
	}
	rules.ResetSeed(seed)
	players := make(map[int]*Player, 2)
	fleets := []*Fleet{}
	result := &BattleSimulationResult{}
	for sideIndex, side := range sides {
		if len(side.Stacks) == 0 || len(side.Stacks) > 64 {
			return nil, fmt.Errorf("side %d needs 1 to 64 stacks", sideIndex+1)
		}
		for _, level := range []int{side.TechLevels.Energy, side.TechLevels.Weapons, side.TechLevels.Propulsion, side.TechLevels.Construction, side.TechLevels.Electronics, side.TechLevels.Biotechnology} {
			if level < 0 || level > 26 {
				return nil, fmt.Errorf("tech levels must be between 0 and 26")
			}
		}
		player := &Player{Num: sideIndex + 1, Race: side.Race, TechLevels: side.TechLevels, AcquiredTechs: side.AcquiredTechs}
		player.Race.Spec = ComputeRaceSpec(&player.Race, &rules)
		players[player.Num] = player
		for stackIndex, stack := range side.Stacks {
			if stack.Quantity < 1 || stack.Quantity > 10000 || stack.QuantityDamaged < 0 || stack.QuantityDamaged > stack.Quantity || math.IsNaN(stack.Damage) || math.IsInf(stack.Damage, 0) || stack.Damage < 0 {
				return nil, fmt.Errorf("side %d stack %d has invalid quantity or damage", player.Num, stackIndex+1)
			}
			design := stack.Design
			design.Slots = slices.Clone(design.Slots)
			design.PlayerNum, design.Num = player.Num, stackIndex+1
			if err := design.validate(&rules, player, false); err != nil {
				return nil, fmt.Errorf("side %d stack %d: %w", player.Num, stackIndex+1, err)
			}
			var err error
			design.Spec, err = ComputeShipDesignSpec(&rules, player.TechLevels, player.Race.Spec, &design)
			if err != nil {
				return nil, err
			}
			if stack.Damage >= float64(design.Spec.Armor) && stack.QuantityDamaged > 0 {
				return nil, fmt.Errorf("starting damage must be less than ship armor")
			}
			if design.Spec.Starbase && stack.Quantity != 1 {
				return nil, fmt.Errorf("a starbase stack must contain exactly one starbase")
			}
			plan := stack.Plan
			if plan.Tactic == BattleTacticUnspecified {
				plan.Tactic = BattleTacticMaximizeDamage
			}
			plan.AttackWho = BattleAttackWhoEveryone
			fleet := &Fleet{MapObject: MapObject{Num: len(fleets) + 1, PlayerNum: player.Num, Type: MapObjectTypeFleet}, Starbase: design.Spec.Starbase, battlePlan: &plan, Tokens: []ShipToken{{DesignNum: design.Num, Quantity: stack.Quantity, QuantityDamaged: stack.QuantityDamaged, Damage: stack.Damage, design: &design}}}
			fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
			fleets = append(fleets, fleet)
			result.Designs = append(result.Designs, &design)
		}
	}
	b := newBattler(slog.New(slog.NewTextHandler(io.Discard, nil)), &rules, 1, players, fleets, nil).(*battle)
	if !b.hasHostility() {
		return nil, fmt.Errorf("no battle: at least one armed ship stack must have orders to attack the other side")
	}
	result.Record = b.runBattle()
	for _, token := range b.tokens {
		loss := token.design.Spec.Cost
		loss.Ironium *= token.quantityDestroyed
		loss.Boranium *= token.quantityDestroyed
		loss.Germanium *= token.quantityDestroyed
		loss.Resources *= token.quantityDestroyed
		retreated := 0
		if token.ranAway {
			retreated = token.Quantity
		}
		result.Stacks = append(result.Stacks, BattleSimulationStackResult{PlayerNum: token.PlayerNum, DesignNum: token.ShipToken.DesignNum, Survivors: token.Quantity, Retreated: retreated, Destroyed: token.quantityDestroyed, Loss: loss})
	}
	slices.SortFunc(result.Stacks, func(a, b BattleSimulationStackResult) int {
		if a.PlayerNum != b.PlayerNum {
			return a.PlayerNum - b.PlayerNum
		}
		return a.DesignNum - b.DesignNum
	})
	return result, nil
}
