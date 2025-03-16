package cs

import (
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
