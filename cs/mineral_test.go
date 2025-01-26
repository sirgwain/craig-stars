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
		name string
		args args
		want MineralType
	}{
		{"Highest Amount", args{Mineral{1, 2, 3}, 1}, Germanium},
		{"4 way tie", args{Mineral{1, 1, 1}, 1}, Ironium},
		{"2nd highest amount", args{Mineral{100, 1, 99}, 2}, Germanium},
		{"lowest amount", args{Mineral{100, 9, 888}, 3}, Boranium},
		{"negative index", args{Mineral{100, 9, 888}, -1}, Boranium},
		{"negative index tie", args{Mineral{100, 9, 100}, -2}, Ironium},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.mineral.HighestType(tt.args.ranking); got != tt.want {
				t.Errorf("Mineral.HighestType() = %v, want %v", got, tt.want)
			}
		})
	}
}
