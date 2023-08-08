package cs

import "testing"

func Test_getCloakPercentForCloakUnits(t *testing.T) {
	type args struct {
		cloakUnits int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"no cloak", args{0}, 0},
		{"70 cloak units = 35% cloaking", args{70}, 35},
		{"110 cloak units = 51% cloaking", args{110}, 51},
		{"200 cloak units = 62% cloaking", args{200}, 62},
		{"400 cloak units = 50% cloaking", args{400}, 79},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCloakPercentForCloakUnits(tt.args.cloakUnits); got != tt.want {
				t.Errorf("getCloakPercentForCloakUnits() = %v, want %v", got, tt.want)
			}
		})
	}
}
