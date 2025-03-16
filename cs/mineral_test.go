package cs

import (
	"reflect"
	"testing"
)

func TestMineral_HighestType(t *testing.T) {
	type args struct {
		mineral Mineral
		ranking int
	}
	tests := []struct {
		name     string
		args     args
		wantType MineralType
	}{
		{"Highest Amount", args{Mineral{1, 2, 3}, 1}, Germanium},
		{"3 way tie", args{Mineral{1, 1, 1}, 1}, Ironium},
		{"Middle amount", args{Mineral{100, 1, 99}, 2}, Germanium},
		{"Lowest amount", args{Mineral{100, 9, 888}, 3}, Boranium},
		{"negative index", args{Mineral{100, 9, 888}, -1}, Boranium},
		{"negative index tie", args{Mineral{100, 9, 100}, -2}, Ironium},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotAmount := tt.args.mineral.HighestType(tt.args.ranking)
			if gotType != tt.wantType {
				t.Errorf("Mineral.HighestType() returned MineralType %v, want %v", gotType, tt.wantType)
			}
			if wantAmount := tt.args.mineral.GetAmount(tt.wantType); gotAmount != wantAmount {
				t.Errorf("Mineral.HighestType() returned amount %v, want %v", gotAmount, wantAmount)
			}
		})
	}
}

func TestMineral_Equalize(t *testing.T) {
	tests := []struct {
		name     string
		mineral  Mineral
		amtToAdd int
		want     Mineral
	}{
		{"adds nothing", Mineral{0, 0, 333}, 0, Mineral{0, 0, 333}},
		{"all equal; spreads leftovers", Mineral{0, 0, 0}, 5, Mineral{2, 2, 1}},
		{"equalizes lowest 2, but not fully", Mineral{0, 4, 2}, 5, Mineral{4, 4, 3}},
		{"unequal; equalizes fully", Mineral{10, 50, 30}, 60, Mineral{50, 50, 50}},
		{"negative amtToAdd", Mineral{33, 34, 33}, -2, Mineral{32, 33, 33}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mineral.Equalize(tt.amtToAdd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mineral.Equalize() = %v, want %v", got, tt.want)
			}
		})
	}
}
