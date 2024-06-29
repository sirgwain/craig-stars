package cs

import "testing"

func Test_techTrade_techLevelGained(t *testing.T) {
	rules := NewRules()
	// override the rules to always get a tech trade
	type args struct {
		current TechLevel
		target  TechLevel
	}
	tests := []struct {
		name string
		tr   *techTrade
		args args
		want TechField
	}{
		{name: "None", args: args{current: TechLevel{}, target: TechLevel{}}, want: TechFieldNone},
		// TODO: put these tests back when we have mocked random numbers
		// {name: "Energy", args: args{current: TechLevel{}, target: TechLevel{Energy: 1}}, want: Energy},
		// {name: "Weapons", args: args{
		// 	current: TechLevel{Energy: 1, Weapons: 1, Propulsion: 1, Construction: 1, Electronics: 1, Biotechnology: 1},
		// 	target:  TechLevel{Energy: 1, Weapons: 2, Propulsion: 3, Construction: 4, Electronics: 5, Biotechnology: 6}},
		// 	want: Weapons,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &techTrade{}
			if got := tr.techLevelGained(&rules, tt.args.current, tt.args.target); got != tt.want {
				t.Errorf("techTrade.techLevelGained() = %v, want %v", got, tt.want)
			}
		})
	}
}
