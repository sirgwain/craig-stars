package cs

import (
	"math"
	"reflect"
	"testing"
)

func TestTechLevel_Lowest(t *testing.T) {
	tests := []struct {
		name string
		tl   TechLevel
		want TechField
	}{
		{"energy lowest by default", TechLevel{}, Energy},
		{"biotech lowest", TechLevel{
			Energy:        6,
			Weapons:       5,
			Propulsion:    4,
			Construction:  3,
			Electronics:   2,
			Biotechnology: 1,
		}, Biotechnology},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tl.Lowest(); got != tt.want {
				t.Errorf("TechLevel.Lowest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTechLevel_LevelsAbove(t *testing.T) {
	tests := []struct {
		name  string
		tl    TechLevel
		other TechLevel
		want  int
	}{
		{"starter tech", TechLevel{}, TechLevel{}, math.MaxInt},
		{"prop 5 humanoid vs DL7", TechLevel{Propulsion: 5}, TechLevel{3, 3, 5, 3, 3, 3}, 0},
		{"prop 6 humanoid vs DL7", TechLevel{Propulsion: 5}, TechLevel{3, 3, 6, 3, 3, 3}, 1},
		{"prop 6 humanoid vs RHRS", TechLevel{Propulsion: 6, Energy: 2}, TechLevel{3, 3, 6, 3, 3, 3}, 0},
		{"starter humanoid vs RHRS", TechLevel{Propulsion: 6, Energy: 2}, TechLevel{3, 3, 3, 3, 3, 3}, -3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tl.LevelsAbove(tt.other); got != tt.want {
				t.Errorf("TechLevel.LevelsAbove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTechLevel_LevelsAboveField(t *testing.T) {
	type args struct {
		other TechLevel
		field TechField
	}
	tests := []struct {
		name string
		tl   TechLevel
		args args
		want int
	}{
		{"none", TechLevel{}, args{TechLevel{}, Energy}, 0},
		{"prop 5 humanoid vs DL7", TechLevel{Propulsion: 5}, args{TechLevel{3, 3, 5, 3, 3, 3}, Propulsion}, 0},
		{"prop 6 humanoid vs DL7", TechLevel{Propulsion: 5}, args{TechLevel{3, 3, 6, 3, 3, 3}, Propulsion}, 1},
		{"prop 6 humanoid vs RHRS", TechLevel{Propulsion: 6, Energy: 2}, args{TechLevel{3, 3, 6, 3, 3, 3}, Propulsion}, 0},
		{"check other field", TechLevel{Energy: 2, Propulsion: 6}, args{TechLevel{2, 3, 6, 3, 3, 3}, Energy}, 0},
		{"check other field", TechLevel{Energy: 2, Propulsion: 6}, args{TechLevel{2, 3, 6, 3, 3, 10}, Biotechnology}, 10},
		{"starter humanoid vs RHRS", TechLevel{Propulsion: 6, Energy: 2}, args{TechLevel{3, 3, 3, 3, 3, 3}, Propulsion}, -3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tl.LevelsAboveField(tt.args.other, tt.args.field); got != tt.want {
				t.Errorf("TechLevel.LevelsAboveField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTechLevel_GetLearnableTechFields(t *testing.T) {

	tests := []struct {
		name      string
		techLevel TechLevel
		want      []TechField
	}{
		{"All techs", TechLevel{}, TechFields},
		{"Some techs", TechLevel{Energy: 26, Weapons: 25, Propulsion: 26, Construction: 10, Electronics: 5, Biotechnology: 0},
			[]TechField{
				Weapons,
				Construction,
				Electronics,
				Biotechnology,
			}},
		{"No techs", TechLevel{Energy: 26, Weapons: 26, Propulsion: 26, Construction: 26, Electronics: 26, Biotechnology: 26}, []TechField{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.techLevel.LearnableTechFields(&rules); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Player.GetLearnableTechFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
