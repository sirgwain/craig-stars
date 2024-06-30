package cs

import (
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
		tr   *techTrade
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
