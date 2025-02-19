package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
)

func Test_CompareAsJSON(t *testing.T) {
	tests := []struct {
		name     string
		got      any
		want     any
		wantFail bool
	}{
		{
			name:     "2 different planets",
			got:      cs.NewPlanet().WithMines(40),
			want:     cs.NewPlanet().WithNum(20),
			wantFail: true,
		},
		{
			name:     "identical players",
			got:      cs.NewPlayer(22, cs.NewRace()),
			want:     cs.NewPlayer(22, cs.NewRace()),
			wantFail: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := CompareAsJSON(t, tt.got, tt.want)
			if res != tt.wantFail {
				t.Errorf("Test_CompareAsJSON() failed; function returned %v but expected %v", res, tt.wantFail)
			}
		})
	}
}

func Test_parseJSONDiff(t *testing.T) {
	tt := []struct {
		name     string
		got      any
		want     any
		wantDiff string
	}{
		{
			name: "2 different planets",
			got:  cs.NewPlanet().WithMines(40),
			want: cs.NewPlanet().WithNum(20),
			wantDiff: `// Test_CompareAsJSON/2_different_planets
{
	"prop-removed": {"mines": 40},
	"num": {"changed": [0, 20]}
}
`,
		},
		// no need for same item tests as this shouldn't even be called if the JSONs are equal
	}

	for _, tt := range tt {
		t.Run(tt.name, func(t *testing.T) {
			gotJSON, err := json.MarshalIndent(tt.got, "", "\t")
			if err != nil {
				t.Fatalf("could not marshal got to JSON: \n%v", err)
			}
			wantJSON, err := json.MarshalIndent(tt.want, "", "\t")
			if err != nil {
				t.Fatalf("could not marshal want to JSON: \n%v", err)
			}

			parseJSONDiff(gotJSON, wantJSON, t.Name())

			// check the diff file to make sure it outputted the correct text
			gotBytes, err := os.ReadFile("../tmp/diff_test.jsonl")
			if err != nil {
				t.Fatalf("error reading diff file: \n%v", err)
			}
			gotDiff := string(gotBytes)
			if gotDiff != tt.wantDiff {
				t.Errorf(`CompareAsJSON() outputted incorrect diff:
Got: %v
Want: %v`, gotDiff, tt.wantDiff)
			}
		})
	}
}
