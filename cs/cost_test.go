package cs

import (
	"math"
	"testing"
)

func TestCost_NumBuildable(t *testing.T) {
	type fields struct {
		Ironium   int
		Boranium  int
		Germanium int
		Resources int
	}
	type args struct {
		cost Cost
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"Build 1", fields{1, 2, 3, 4}, args{Cost{1, 2, 3, 4}}, 1},
		{"Build 2", fields{2, 4, 6, 8}, args{Cost{1, 2, 3, 4}}, 2},
		{"Build 2 w/one limiting mineral", fields{2, 400, 600, 800}, args{Cost{1, 2, 3, 4}}, 2},
		{"Build none w/one missing mineral", fields{0, 400, 600, 800}, args{Cost{1, 2, 3, 4}}, 0},
		{"Build 1 w/only resource cost", fields{1, 2, 3, 4}, args{Cost{0, 0, 0, 4}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			available := Cost{
				Ironium:   tt.fields.Ironium,
				Boranium:  tt.fields.Boranium,
				Germanium: tt.fields.Germanium,
				Resources: tt.fields.Resources,
			}
			if got := available.NumBuildable(tt.args.cost); got != tt.want {
				t.Errorf("Cost.NumBuildable() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		{"0", fields{0, 0, 0, 0}, args{Cost{0, 0, 0, 0}}, math.MaxFloat64},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Cost{
				Ironium:   tt.fields.Ironium,
				Boranium:  tt.fields.Boranium,
				Germanium: tt.fields.Germanium,
				Resources: tt.fields.Resources,
			}
			if got := a.Divide(tt.args.b); got != tt.want {
				t.Errorf("Cost.Divide() = %v, want %v", got, tt.want)
			}
		})
	}
}
