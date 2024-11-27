package cs

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCost_Divide(t *testing.T) {
	type fields struct {
		Ironium   int
		Boranium  int
		Germanium int
		Resources int
	}
	type args struct {
		b Cost
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{"0", fields{0, 0, 0, 0}, args{Cost{0, 0, 0, 0}}, math.Inf(1)},
		{"1 I", fields{1, 0, 0, 0}, args{Cost{1, 0, 0, 0}}, 1},
		{"1 B", fields{0, 1, 0, 0}, args{Cost{0, 1, 0, 0}}, 1},
		{"1 G", fields{0, 0, 1, 0}, args{Cost{0, 0, 1, 0}}, 1},
		{"1 R", fields{0, 0, 0, 1}, args{Cost{0, 0, 0, 1}}, 1},
		{"2 I", fields{2, 0, 0, 0}, args{Cost{1, 0, 0, 0}}, 2},
		{"2 B", fields{0, 2, 0, 0}, args{Cost{0, 1, 0, 0}}, 2},
		{"2 G", fields{0, 0, 2, 0}, args{Cost{0, 0, 1, 0}}, 2},
		{"2 R", fields{0, 0, 0, 2}, args{Cost{0, 0, 0, 1}}, 2},
		{"2 All", fields{2, 2, 2, 2}, args{Cost{1, 1, 1, 1}}, 2},
		{"1/2 I", fields{1, 0, 0, 0}, args{Cost{2, 0, 0, 0}}, .5},
		{"1/2 B", fields{0, 1, 0, 0}, args{Cost{0, 2, 0, 0}}, .5},
		{"1/2 G", fields{0, 0, 1, 0}, args{Cost{0, 0, 2, 0}}, .5},
		{"1/2 R", fields{0, 0, 0, 1}, args{Cost{0, 0, 0, 2}}, .5},
		{"1/2 All", fields{1, 1, 1, 1}, args{Cost{2, 2, 2, 2}}, .5},
		{"5887 / 841", fields{199, 1555, 841, 92}, args{Cost{71, 5, 5887, 17}}, .142857},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Cost{
				Ironium:   tt.fields.Ironium,
				Boranium:  tt.fields.Boranium,
				Germanium: tt.fields.Germanium,
				Resources: tt.fields.Resources,
			}
			got := a.Divide(tt.args.b)
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
