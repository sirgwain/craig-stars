package cs

import "testing"

func TestClamp(t *testing.T) {
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

func TestAbsMin(t *testing.T) {
	tests := []struct {
		name string
		nums []float64
		want float64
	}{
		{"grabs closest to 0", []float64{1, 2, 3, 0}, 0},
		{"all negative", []float64{-1, -1.2, -0.31, -4}, -0.31},
		{"mix; lowest positive", []float64{1, -222, 3, -10.3333}, 1},
		{"mix; lowest negative", []float64{2025, -1997, 2001, -3}, -3},
		{"takes last of absolute equals", []float64{1, 1, 1, 1, -1}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AbsMin(tt.nums...); got != tt.want {
				t.Errorf("AbsMin() returned value %v, want %v", got, tt.want)
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

func TestPowInt(t *testing.T) {
	tests := []struct {
		name     string
		base     int
		exponent int
		want     int
	}{
		{"1^4", 1, 4, 1},
		{"2^3", 2, 3, 8},
		{"30^4", 30, 4, 810_000},
		{"30^2", 30, 2, 900},
		{"5^3", 5, 3, 125},
		{"2^20", 2, 20, 1_048_576},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PowInt(tt.base, tt.exponent); got != tt.want {
				t.Errorf("PowInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
