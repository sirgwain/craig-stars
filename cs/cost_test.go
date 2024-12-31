package cs

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCost_Divide(t *testing.T) {
	tests := []struct {
		name     string
		dividend Cost
		divisor  Cost
		want     float64
	}{
		{"0", Cost{0, 0, 0, 0}, Cost{0, 0, 0, 0}, math.Inf(1)},
		{"1/1 I", Cost{1, 0, 0, 0}, Cost{1, 0, 0, 0}, 1},
		{"1 B", Cost{0, 1, 0, 0}, Cost{0, 1, 0, 0}, 1},
		{"1 G", Cost{0, 0, 1, 0}, Cost{0, 0, 1, 0}, 1},
		{"1 R", Cost{0, 0, 0, 1}, Cost{0, 0, 0, 1}, 1},
		{"2 I", Cost{2, 0, 0, 0}, Cost{1, 0, 0, 0}, 2},
		{"2 B", Cost{0, 2, 0, 0}, Cost{0, 1, 0, 0}, 2},
		{"2 G", Cost{0, 0, 2, 0}, Cost{0, 0, 1, 0}, 2},
		{"2 R", Cost{0, 0, 0, 2}, Cost{0, 0, 0, 1}, 2},
		{"2 All", Cost{2, 2, 2, 2}, Cost{1, 1, 1, 1}, 2},
		{"1/2 I", Cost{1, 0, 0, 0}, Cost{2, 0, 0, 0}, .5},
		{"1/2 B", Cost{0, 1, 0, 0}, Cost{0, 2, 0, 0}, .5},
		{"1/2 G", Cost{0, 0, 1, 0}, Cost{0, 0, 2, 0}, .5},
		{"1/2 R", Cost{0, 0, 0, 1}, Cost{0, 0, 0, 2}, .5},
		{"1/2 All", Cost{1, 1, 1, 1}, Cost{2, 2, 2, 2}, .5},
		{"841 / 5887", Cost{199, 1555, 841, 92}, Cost{71, 5, 5887, 17}, .142857},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dividend.DivideCost(tt.divisor)
			assert.InDeltaf(t, got, tt.want, 0.01, fmt.Sprintf("Cost.Divide() = %v, want %v", got, tt.want))
		})
	}
}

func TestCost_Max(t *testing.T) {
	tests := []struct {
		name  string
		cost  Cost
		other Cost
		want  Cost
	}{
		{"0 case", Cost{}, Cost{}, Cost{}},
		{"cost greater", Cost{1, 2, 3, 4}, Cost{}, Cost{1, 2, 3, 4}},
		{"other greater", Cost{}, Cost{1, 2, 3, 4}, Cost{1, 2, 3, 4}},
		{"mix", Cost{1, 2, 3, 4}, Cost{2, 1, 4, 3}, Cost{2, 2, 4, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cost.Max(tt.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cost.Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCost_HighestType(t *testing.T) {
	type args struct {
		cost    Cost
		ranking int
	}
	tests := []struct {
		name string
		args args
		want CostType
	}{
		{"Highest Amount", args{Cost{1, 2, 3, 4}, 1}, Resources},
		{"4 way tie", args{Cost{1, 1, 1, 1}, 1}, Ironium},
		{"2nd highest amount", args{Cost{100, 1, 99, 88}, 2}, Germanium},
		{"lowest amount", args{Cost{100, 9, 100, 888}, 4}, Boranium},
		{"negative index", args{Cost{100, 9, 100, 888}, -1}, Boranium},
		{"negative index tie", args{Cost{100, 9, 100, 888}, -2}, Ironium},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.cost.HighestType(tt.args.ranking); got != tt.want {
				t.Errorf("Cost.HighestType() = %v, want %v", got, tt.want)
			}
		})
	}
}
