package cs

import (
	"math"
	"reflect"
	"testing"
)

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

func Test_splitValues(t *testing.T) {
	type args struct {
		sourceCapacity int
		destCapacity1  int
		destCapacity2  int
		values         []int
	}
	tests := []struct {
		name    string
		args    args
		want1   []int
		want2   []int
		wantErr bool
	}{
		{
			name:    "invalid args",
			args:    args{sourceCapacity: 4, destCapacity1: 5, destCapacity2: 0},
			wantErr: true,
		},
		{
			name:  "no split",
			args:  args{sourceCapacity: 4, destCapacity1: 4, destCapacity2: 0, values: []int{1, 1, 1, 1}},
			want1: []int{1, 1, 1, 1},
			want2: []int{0, 0, 0, 0},
		},
		{
			name: "split all",

			args:  args{sourceCapacity: 4, destCapacity1: 0, destCapacity2: 4, values: []int{1, 1, 1, 1}},
			want1: []int{0, 0, 0, 0},
			want2: []int{1, 1, 1, 1},
		},
		{
			name: "split even",

			args:  args{sourceCapacity: 8, destCapacity1: 4, destCapacity2: 4, values: []int{2, 2, 2, 2}},
			want1: []int{1, 1, 1, 1},
			want2: []int{1, 1, 1, 1},
		},
		{
			name: "split 1",

			args:  args{sourceCapacity: 8, destCapacity1: 7, destCapacity2: 1, values: []int{2, 2, 2, 2}},
			want1: []int{2, 2, 2, 1},
			want2: []int{0, 0, 0, 1},
		},
		{
			name: "split 2",

			args:  args{sourceCapacity: 8, destCapacity1: 6, destCapacity2: 2, values: []int{2, 2, 2, 2}},
			want1: []int{2, 2, 2, 0},
			want2: []int{0, 0, 0, 2},
		},
		{
			name: "split 3",

			args:  args{sourceCapacity: 8, destCapacity1: 5, destCapacity2: 3, values: []int{2, 2, 2, 2}},
			want1: []int{1, 1, 1, 2},
			want2: []int{1, 1, 1, 0},
		},
		{
			name: "split 2/3",

			args:  args{sourceCapacity: 10, destCapacity1: 7, destCapacity2: 3, values: []int{3, 3, 3, 1}},
			want1: []int{2, 2, 2, 1},
			want2: []int{1, 1, 1, 0},
		},
		{
			name: "split many",

			args:  args{sourceCapacity: 2000, destCapacity1: 1500, destCapacity2: 500, values: []int{100, 200, 300, 400}},
			want1: []int{75, 150, 225, 300},
			want2: []int{25, 50, 75, 100},
		},
		{
			name: "split 33%",

			args:  args{sourceCapacity: 300, destCapacity1: 200, destCapacity2: 100, values: []int{100, 0, 0, 0}},
			want1: []int{67, 0, 0, 0},
			want2: []int{33, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, err := splitValues(tt.args.sourceCapacity, tt.args.destCapacity1, tt.args.destCapacity2, tt.args.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("splitValues() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("splitValues() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_divideRoundUp(t *testing.T) {
	tests := []struct {
		name     string
		dividend int
		divisor  int
		want     int
	}{
		{"10 ÷ 3", 10, 3, 4},
		{"9 ÷ 3", 9, 3, 3},
		{"1 ÷ 2", 1, 2, 1},
		{"0 ÷ 55", 0, 55, 0},
		{"100 ÷ 7", 100, 7, 15},
		{"negative dividend", -10, 3, -4},
		{"negative divisor", 10, -3, -4},
		{"both negative", -10, -3, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := divideRoundAway0(tt.dividend, tt.divisor); got != tt.want {
				t.Errorf("divideRoundUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogBase(t *testing.T) {
	tests := []struct {
		name string
		base float64
		x    float64
		want float64
	}{
		{"log2(8)", 2, 8, 3},
		{"log10(1000)", 10, 1000, 3},
		{"log3(27)", 3, 27, 3},
		{"log5(125)", 5, 125, 3},
		{"log10(0.1)", 10, 0.1, -1},
		{"invalid base <= 0", -2, 8, math.NaN()},
		{"invalid base == 1", 1, 8, math.NaN()},
		{"invalid x <= 0", 2, -8, math.NaN()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LogBase(tt.base, tt.x)
			if (math.IsNaN(got) && !math.IsNaN(tt.want)) || (!math.IsNaN(got) && math.Abs(got-tt.want) > 1e-9) {
				t.Errorf("LogBase() = %v, want %v", got, tt.want)
			}
		})
	}
}
