package cs

import "testing"

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
