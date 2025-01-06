package cs

import "testing"

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

func Test_roundHalfTowards0(t *testing.T) {
	tests := []struct {
		name string
		num  float64
		want float64
	}{
		{"positive, <0.5", 0.2, 0},
		{"positive, =0.5", 1.5, 1},
		{"positive, >0.5", 3.6, 4},
		{"negative, <0.5", -1.4, -1},
		{"negative, =0.5", -71.5, -71},
		{"negative, >0.5", -1.6, -2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := roundHalfTowards0(tt.num); got != tt.want {
				t.Errorf("roundHalfTowards0() = %v, want %v", got, tt.want)
			}
		})
	}
}
