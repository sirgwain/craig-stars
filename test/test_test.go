package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
)

func Test_CompareAsJSON(t *testing.T) {
	tests := []struct {
		name       string
		got        any
		want       any
		wantFailed bool
	}{
		{
			name:       "2 different planets",
			got:        cs.NewPlanet().WithMines(40),
			want:       cs.NewPlanet().WithNum(20),
			wantFailed: true,
		},
		{
			name:       "identical players",
			got:        cs.NewPlayer(22, cs.NewRace()),
			want:       cs.NewPlayer(22, cs.NewRace()),
			wantFailed: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockTestingT)
			m.name = tt.name
			CompareAsJSON(m, tt.got, tt.want)
			// m.Failed is set to true when func would've normally failed a test
			if m.failed != tt.wantFailed {
				var s string
				if tt.wantFailed {
					s = "did not fail test when expected"
				} else {
					s = "failed test unexpectedly"
				}
				t.Errorf("Test_CompareAsJSON() %s; test failed flag returned %v instead of %v", s, m.failed, tt.wantFailed)
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
			wantDiff: `{
	"mines": {"changed": [40, 0]},
	"num": {"changed": [0, 20]}
}
`,
		},
		// no need for same item tests as this shouldn't even be called if the JSONs are equal
	}

	for _, tt := range tt {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantDiff = "// " + t.Name() + "\n" + tt.wantDiff
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
			gotBytes, err := os.ReadFile("../tmp/diff.jsonl")
			if err != nil {
				t.Fatalf("error reading diff file: \n%v", err)
			}
			gotDiff := string(gotBytes)
			if gotDiff != tt.wantDiff {
				t.Errorf("CompareAsJSON() outputted incorrect diff:\nGot: %v\nWant: %v", gotDiff, tt.wantDiff)
			}
		})
	}
}
