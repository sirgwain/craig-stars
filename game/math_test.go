package game

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
			if got := clamp(tt.args.value, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("clamp() = %v, want %v", got, tt.want)
			}
		})
	}
}
