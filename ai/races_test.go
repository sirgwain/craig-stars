package ai

import (
	"testing"
)

func TestGetRandomRaces(t *testing.T) {
	type args struct {
		numRaces int
	}
	tests := []struct {
		name string
		args args
	}{
		{"don't crash for lots of races", args{1000}},
		{"don't crash for one race", args{1}},
		{"don't crash for a few races", args{3}},
		{"don't crash for all ai races", args{len(Races)}},
		{"don't crash for 0 ai races", args{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRandomRaces(tt.args.numRaces); len(got) != tt.args.numRaces {
				t.Errorf("GetRandomRaces() = %d races, want %d races", len(got), tt.args.numRaces)
			}
		})
	}
}
