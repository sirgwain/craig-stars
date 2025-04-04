package ai

import (
	"encoding/json"
	"testing"
)

func TestGetRandomRaces(t *testing.T) {
	tests := []struct {
		name     string
		numRaces int
		cheater  bool
	}{
		{"0 races", 0, false},
		{"1 race", 1, false},
		{"3 races", 3, false},
		{"all normal races", len(Races), false},
		{"1000 races", 1000, false},
		{"all cheater races", len(CheaterRaces), true},
		{"1000 cheater races", 1000, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRandomRaces(tt.numRaces, tt.cheater); len(got) != tt.numRaces {
				t.Errorf("GetRandomRaces() returned total of %d races, want %d races", len(got), tt.numRaces)
				races, _ := json.MarshalIndent(got, "", "\t")
				t.Logf("Races: \n%s", string(races))
			}
		})
	}
}
