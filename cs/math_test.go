package cs

import (
	"testing"
)

func Test_clamp(t *testing.T) {
	type args struct {
		value int
		min   int
		max   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"clamp 1 -> 0 to 100", args{1, 0, 100}, 1},
		{"clamp -1 -> 0 to 100", args{-1, 0, 100}, 0},
		{"clamp 101 -> 0 to 100", args{101, 0, 100}, 100},
		{"clamp -25 -> -50 to 0", args{-25, -50, 0}, -25},
		{"clamp -51 -> -50 to 0", args{-51, -50, 0}, -50},
		{"clamp 1 -> -50 to 0", args{1, -50, 0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clamp(tt.args.value, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("clamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNewJamming(t *testing.T) {
	type args struct {
		prevBonus      float64
		componentBonus float64
		multi          float64
		qty            int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"2x 20%", args{0, 0.2, 1, 2}, 0.36},
		{"4x 30%; 0.75x multi", args{0, 0.3, 0.75, 4}, 0.5699},
		{"2x 20%; prev 10%", args{0.1, 0.2, 1, 2}, 0.424},
		{"10x 30%; prev 10%", args{0.1, 0.3, 1, 10}, 0.9746},
		{"2x 20%; prev 10.8%, 0.75x multi", args{0.10875, 0.2, 0.75, 2}, 0.3396},
		{"4x 30%; prev 57%, 0.75x multi", args{0.569925, 0.3, 0.75, 4}, 0.7068},
		{"4x 20%; prev 57%, 0.75x multi", args{0.569925, 0.2, 0.75, 4}, 0.6762},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := roundFloat(getNewJamming(tt.args.prevBonus, tt.args.componentBonus, tt.args.multi, tt.args.qty), 4); got != tt.want {
				t.Errorf("getNewJamming() = %v, want %v", got, tt.want)
			}
		})
	}
}
