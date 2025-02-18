package test

import (
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
		wantDiff string
	}{
		{
			name:     "2 different planets",
			got:      cs.NewPlanet().WithMines(40),
			want:     cs.NewPlanet().WithNum(20),
			wantFail: true,
			wantDiff: `// Test_CompareAsJSON/2_different_planets
{
	"prop-removed": {"mines": 40},
	"num": {"changed": [0, 20]}
}
`,
		},
		{
			name:     "identical players",
			got:      cs.NewPlayer(22, cs.NewRace()),
			want:     cs.NewPlayer(22, cs.NewRace()),
			wantFail: false,
			wantDiff: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if we don't want the test to fail, call the function as normal
			if !tt.wantFail {
				CompareAsJSON(t, tt.got, tt.want)
				return
			}

			ExpectFailure(t, func(t *testing.T) {
				CompareAsJSON(t, tt.got, tt.want)
			})
			// check the diff to make sure it outputted the correct stuff
			gotBytes, err := os.ReadFile("../tmp/diff.jsonl")
			if err != nil {
				t.Fatalf("error reading diff file: \n%v", err)
			}
			gotDiff := string(gotBytes)
			if gotDiff != tt.wantDiff {
				t.Errorf("CompareAsJSON() outputted incorrect diff:\n"+
					"Got: %v\n"+
					"Want: %v", gotDiff, tt.wantDiff)
			}
		})
	}
}

func Test_F(t *testing.T) {
	t.Run("a", func(t *testing.T) {
		for n := 0; n < 3; n++ {
			t.Run(fmt.Sprintf("test_%d", n), func(t *testing.T) {

				t.Logf("Executing test %d", n)
			})
		}
	})
	t.Run("b", func(t *testing.T) {
		for n := 0; n < 3; n++ {
			t.Run(fmt.Sprintf("test_%d", n), func(t *testing.T) {

				t.Logf("Executing test %d", n)
			})
		}
	})
}
