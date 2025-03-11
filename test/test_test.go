package test

import (
	"errors"
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
			name:       "identical players",
			got:        cs.NewPlayer(22, cs.NewRace()),
			want:       cs.NewPlayer(22, cs.NewRace()),
			wantFailed: false,
			wantDiff:   "",
		},
		{
			name:       "2 different planets",
			got:        cs.NewPlanet().WithMines(40).WithHomeworld(true),
			want:       cs.NewPlanet().WithNum(20).WithContributesOnlyLeftoverToResearch(true),
			wantFailed: true,
			wantDiff: `// TestCompareAsJSON/2_different_planets
{
	/* Added */ "contributesOnlyLeftoverToResearch": true,
	/* Removed */ "homeworld": true,
	"mines": /* Changed */ [ 40, 0 ],
	"num": /* Changed */ [ 0, 20 ]
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				// remove json afterwards to make sure successive runs still pass
				if err := os.RemoveAll("tmp"); err != nil {
					t.Fatalf("error cleaning temp files: \n%v", err)
				}
				if err := os.Mkdir("tmp", 0755); err != nil {
					t.Fatalf("error cleaning temp files: \n%v", err)
				}
			})

			m := new(mockTestingT)
			m.name = t.Name() // fixup to add doc comment
			CompareAsJSON(m, tt.got, tt.want)
			// Check if function failed;
			// m.failed is set to true when function would've normally failed a test
			if m.failed != tt.wantFailed {
				var s string
				if tt.wantFailed {
					s = "did not fail test when expected"
				} else {
					s = "failed test unexpectedly"
				}
				t.Errorf("CompareAsJSON() %s; test failed flag returned %v instead of %v", s, m.failed, tt.wantFailed)
			}

			if tt.wantFailed {
				// if we wanted test to fail, check the diff file to make sure it contains the correct text
				gotBytes, err := os.ReadFile("../tmp/diff.jsonl")
				if errors.Is(err, os.ErrNotExist) {
					t.Fatal("Test failed to create diff file when expected")
				} else if err != nil {
					t.Fatalf("error reading diff file: \n%v", err)
				}
				gotDiff := string(gotBytes)
				if gotDiff != tt.wantDiff {
					t.Fatalf("CompareAsJSON() outputted incorrect diff:\nGot: \n%v\nWant: \n%v", gotDiff, tt.wantDiff)
				}
				return
			}

			// if we want the test to succeed, check to ensure json files don't exist
			for _, path := range []string{"../tmp/got.jsonl", "../tmp/want.jsonl", "../tmp/diff.jsonl"} {
				if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
					// no diff = happy test
					continue
				} else if err != nil {
					t.Fatalf("error checking JSON file existence at %s: \n%v", path, err)
				}

				// uh oh, we found a file when we weren't supposed to...
				gotBytes, err := os.ReadFile(path)
				if err != nil {
					t.Errorf("error reading JSON file at %s: \n%v", path, err)
				}
				t.Fatalf("CompareAsJSON() created file at %s when it wasn't supposed to; file contents: \n%s", path, gotBytes)
			}

		})

	}
}
