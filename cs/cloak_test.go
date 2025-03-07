package cs

import (
	"testing"

	"github.com/sirgwain/craig-stars/test"
)

func Test_getCloakPercentForCloakUnits(t *testing.T) {
	tests := []struct {
		name       string
		cloakUnits int
		want       int
	}{
		{"no cloak", 0, 0},
		{"19 cloak units = 10% cloaking", 19, 10},
		{"70 cloak units = 35% cloaking", 70, 35},
		{"110 cloak units = 51% cloaking", 110, 51},
		{"200 cloak units = 62% cloaking", 200, 62},
		{"400 cloak units = 79% cloaking", 400, 79},
		{"1380 cloak units = 97% cloaking", 1380, 97},
		{"2000 cloak units = 98% cloaking", 2000, 98},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCloakPercentForCloakUnits(tt.cloakUnits); got != tt.want {
				t.Errorf("getCloakPercentForCloakUnits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCloakFactor(t *testing.T) {
	type args struct {
		cloakPercent         int
		cloakReductionFactor float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"none", args{cloakPercent: 0, cloakReductionFactor: 1}, 1},
		{"55% cloak", args{cloakPercent: 55, cloakReductionFactor: 1}, .45},
		{"75% cloak wormhole", args{cloakPercent: 75, cloakReductionFactor: 1}, .25},
		{"75% cloak wormhole, one tachyon", args{cloakPercent: 75, cloakReductionFactor: .95}, .2875},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCloakFactor(tt.args.cloakPercent, tt.args.cloakReductionFactor); !test.WithinTolerance(got, tt.want, .0001) {
				t.Errorf("getCloakFactor() = %v, want %v", got, tt.want)
			}
		})
	}
}
