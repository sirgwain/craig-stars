//go:build !wasi && !wasm

package cs

import (
	"reflect"
	"testing"
)

func Test_techTrade_techLevelGained(t *testing.T) {
	// override the rules to always get a tech trade
	type args struct {
		current TechLevel
		target  TechLevel
	}
	tests := []struct {
		name string
		args args
		rng  rng
		want TechField
	}{
		{name: "None", args: args{current: TechLevel{}, target: TechLevel{}}, rng: newFloat64Random(0), want: TechFieldNone},
		{name: "Energy", args: args{current: TechLevel{}, target: TechLevel{Energy: 1}}, rng: newFloat64Random(.25), want: Energy},
		{name: "Weapons", args: args{
			current: TechLevel{Energy: 1, Weapons: 1, Propulsion: 1, Construction: 1, Electronics: 1, Biotechnology: 1},
			target:  TechLevel{Energy: 1, Weapons: 2, Propulsion: 3, Construction: 4, Electronics: 5, Biotechnology: 6}},
			rng:  &testRandom{floatsToReturn: []float64{.25}, intsToReturn: []int{1}},
			want: Weapons,
		},
		{name: "None, random didn't work out", args: args{current: TechLevel{}, target: TechLevel{Energy: 1}}, rng: newFloat64Random(.5), want: TechFieldNone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &techTrade{}
			rules := NewRulesWithSeed(0)
			rules.random = tt.rng
			if got := tr.techLevelGained(&rules, tt.args.current, tt.args.target); got != tt.want {
				t.Errorf("techTrade.techLevelGained() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_techTrade_acquirablePartGained(t *testing.T) {
	type token struct {
		hull  TechHull
		slots []ShipDesignSlot
		qty   int
	}
	tests := []struct {
		name            string
		acquiredTechs   []string
		tokens          []token
		partChanceRolls []float64
		partDraws       []int
		acquiredTech    bool
		want            *Tech
	}{
		{
			name:          "No parts to give",
			acquiredTechs: []string{},
			tokens: []token{
				{
					hull: Scout,
					slots: []ShipDesignSlot{
						{
							HullComponent: QuickJump5.Name,
							HullSlotIndex: 1,
							Quantity:      1,
						},
					}, qty: 1,
				},
			},
			partChanceRolls: []float64{0},
			acquiredTech:    false,
			want:            nil,
		},
		{
			name:          "Random roll failed",
			acquiredTechs: []string{},
			tokens: []token{
				{
					hull: Scout, slots: []ShipDesignSlot{
						{
							HullComponent: QuickJump5.Name,
							HullSlotIndex: 1,
							Quantity:      1,
						},
					}, qty: 1,
				},
			},
			partChanceRolls: []float64{1},
			acquiredTech:    false,
			want:            nil,
		},
		{
			name:          "Part already acquired",
			acquiredTechs: []string{EnigmaPulsar.Name},
			tokens: []token{
				{
					hull: Scout,
					slots: []ShipDesignSlot{
						{
							HullComponent: EnigmaPulsar.Name,
							HullSlotIndex: 1,
							Quantity:      1,
						},
					}, qty: 1,
				},
			},
			partChanceRolls: []float64{0},
			acquiredTech:    false,
			want:            nil,
		},
		{
			name:          "Already aquired a part this turn",
			acquiredTechs: []string{AlienMiner.Name},
			tokens: []token{
				{
					hull: Scout,
					slots: []ShipDesignSlot{
						{
							HullComponent: EnigmaPulsar.Name,
							HullSlotIndex: 1,
							Quantity:      1,
						},
					}, qty: 1,
				},
			},
			partChanceRolls: []float64{0},
			acquiredTech:    true,
			want:            nil,
		},
		{
			name:          "Trader hull cannot be learned from scrapping",
			acquiredTechs: []string{},
			tokens: []token{
				{
					hull: MiniMorph,
					slots: []ShipDesignSlot{
						{
							HullComponent: LongHump6.Name,
							HullSlotIndex: 1,
							Quantity:      1,
						},
					}, qty: 12,
				},
				{
					hull: MiniMorph,
					slots: []ShipDesignSlot{
						{
							HullComponent: QuickJump5.Name,
							HullSlotIndex: 1,
							Quantity:      1,
						},
					}, qty: 13,
				},
			},
			partChanceRolls: []float64{0.125},
			acquiredTech:    false,
			want:            nil,
		},
		{
			name:          "Component counts ignore ship count; randomized draws can retry",
			acquiredTechs: []string{MiniMorph.Name},
			tokens: []token{
				{
					hull: MidgetMiner,
					slots: []ShipDesignSlot{
						{
							HullComponent: EnigmaPulsar.Name,
							HullSlotIndex: 1,
							Quantity:      1,
						},
						{
							HullComponent: AlienMiner.Name,
							HullSlotIndex: 2,
							Quantity:      2,
						},
					}, qty: 5,
				},
				{
					hull: Miner,
					slots: []ShipDesignSlot{
						{
							HullComponent: EnigmaPulsar.Name,
							HullSlotIndex: 1,
							Quantity:      1,
						},
						{
							HullComponent: AlienMiner.Name,
							HullSlotIndex: 2,
							Quantity:      3,
						},
					}, qty: 10,
				},
			},
			partChanceRolls: []float64{1, 1, 0.04},
			partDraws:       []int{9, 4, 4},
			acquiredTech:    false,
			want:            &AlienMiner.Tech,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &techTrade{}
			rules := NewRulesWithSeed(0)
			player := NewPlayer(1, NewRace())
			player.acquirablePartGained = tt.acquiredTech
			for _, t := range tt.acquiredTechs {
				player.AcquiredTechs[t] = true
			}
			rng := newFloat64Random(tt.partChanceRolls...)
			rng.intsToReturn = tt.partDraws

			tokens := []ShipToken{}
			for n, token := range tt.tokens {
				design := NewShipDesign(player.Num, n+1).WithSlots(token.slots).WithHull(token.hull.Name)
				tokens = append(tokens, ShipToken{
					DesignNum: n + 1,
					Quantity:  token.qty,
					design:    design,
				})
			}

			rules.random = rng

			if got := tr.acquirablePartGained(&rules, player, tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("techTrade.acquirablePartGained() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkAcquirablePartChance(t *testing.T) {
	tests := []struct {
		name string
		qty  int
		rng  rng
		want bool
	}{
		{name: "1 part; failed roll", qty: 1, rng: newFloat64Random(1), want: false},
		{name: "10 parts; 5% roll", qty: 10, rng: newFloat64Random(0.05), want: true},
		{name: "50 parts capped at 25; no second independent batch", qty: 50, rng: newFloat64Random(1, 0.125), want: false},
	}
	for _, tt := range tests {
		rules := NewRulesWithSeed(0)
		rules.random = tt.rng
		t.Run(tt.name, func(t *testing.T) {
			if got := checkAcquirablePartChance(&rules, tt.qty); got != tt.want {
				t.Errorf("checkAcquirablePartChance() = %v, want %v", got, tt.want)
			}
		})
	}
}
