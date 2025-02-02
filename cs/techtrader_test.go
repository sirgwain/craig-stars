package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
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
			rng:  newFloat64Random(.25),
			want: Weapons,
		},
		{name: "None, random didn't work out", args: args{current: TechLevel{}, target: TechLevel{Energy: 1}}, rng: newFloat64Random(.3), want: TechFieldNone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &techTrade{}
			rules := NewRules()
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
			name:          "25 parts split in 2 fleets, 12.5% roll",
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
			want:            &MiniMorph.Tech,
		},
		{
			name:          "15 Enigmas/40 Alien Miners, fail x2 -> success",
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
			partChanceRolls: []float64{1, 1, 0.075},
			acquiredTech:    false,
			want:            &AlienMiner.Tech,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &techTrade{}
			rules := NewRules()
			player := NewPlayer(1, NewRace())
			player.acquirablePartGained = tt.acquiredTech
			for _, t := range tt.acquiredTechs {
				player.AcquiredTechs[t] = true
			}
			rng := newFloat64Random(tt.partChanceRolls...)

			tokens := []ShipToken{}
			for n, token := range tt.tokens {
				design := NewShipDesign(player.Num, n+1).WithSlots(token.slots).WithHull(token.hull.Name)
				tokens = append(tokens, ShipToken{
					DesignNum: n + 1,
					Quantity:  token.qty,
					design:    design,
				})
				for i := len(token.slots); i > 0; i-- {
					rng.addInts(i) // keeps list in order as each element is swapped with itself in list sorting
				}
			}

			rules.random = rng

			if got := tr.acquirablePartGained(&rules, player, tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("techTrade.acquirablePartGained() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_techTradeChance(t *testing.T) {
	type args struct {
		baseChance float64
		level      int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"no levels", args{baseChance: .5, level: 0}, 0},
		{"one level", args{baseChance: .5, level: 1}, .25},
		{"two levels", args{baseChance: .5, level: 2}, .375},
		{"six levels", args{baseChance: .5, level: 6}, .492},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := techTradeChance(tt.args.baseChance, tt.args.level); !test.WithinTolerance(got, tt.want, .001) {
				t.Errorf("techTradeChance() = %v, want %v", got, tt.want)
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
		{name: "50 parts; fail -> success", qty: 50, rng: newFloat64Random(1, 0.125), want: true},
	}
	for _, tt := range tests {
		rules := NewRules()
		rules.random = tt.rng
		t.Run(tt.name, func(t *testing.T) {
			if got := checkAcquirablePartChance(&rules, tt.qty); got != tt.want {
				t.Errorf("checkAcquirablePartChance() = %v, want %v", got, tt.want)
			}
		})
	}
}
