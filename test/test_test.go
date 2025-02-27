package test

import (
	"os"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
)

func TestCompareAsJSON(t *testing.T) {
	tests := []struct {
		name       string
		got        any
		want       any
		wantFailed bool
		wantDiff   string
	}{
		{
			name:       "2 different planets",
			got:        cs.NewPlanet().WithMines(40),
			want:       cs.NewPlanet().WithNum(20),
			wantFailed: true,
			wantDiff: `// TestCompareAsJSON/2_different_planets
{
	"mines": {"changed": [40, 0]},
	"num": {"changed": [0, 20]}
}
`,
		},
		{
			name:       "identical players",
			got:        cs.NewPlayer(22, cs.NewRace()),
			want:       cs.NewPlayer(22, cs.NewRace()),
			wantFailed: false,
			wantDiff:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fixup by adding doc comment
			m := new(mockTestingT)
			m.name = t.Name()
			CompareAsJSON(m, tt.got, tt.want)
			// Check if function failed;
			// m.failed is set to true when func would've normally failed a test
			if m.failed != tt.wantFailed {
				var s string
				if tt.wantFailed {
					s = "did not fail test when expected"
				} else {
					s = "failed test unexpectedly"
				}
				t.Errorf("CompareAsJSON() %s; test failed flag returned %v instead of %v", s, m.failed, tt.wantFailed)
			}

			// check the diff file to make sure it outputted the correct text
			gotBytes, err := os.ReadFile("../tmp/diff.jsonl")
			if err != nil {
				t.Fatalf("error reading diff file: \n%v", err)
			}
			gotDiff := string(gotBytes)
			if gotDiff != tt.wantDiff {
				t.Errorf("CompareAsJSON() outputted incorrect diff:\nGot: \n%v\nWant: \n%v", gotDiff, tt.wantDiff)
			}
			os.RemoveAll("../tmp/diff.jsonl") // delete temp diff afterwards
		})

	}
}
