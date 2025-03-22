// The `cs` package houses all internal game logic and data types for craig-stars.
// It serves as the "backbone" upon which all the other packages build and rely on.
package cs

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

// Package-wide initialization function to move temp json files after test runs.
func TestMain(m *testing.M) {
	m.Run()
	// Don't move files on CI runs if the target files already exist
	// (since that likely indicates gotestsum rerunning failed cases)
	if _, err := os.Stat("../tmp/got_cs.jsonl"); os.Getenv("CI") == "" || errors.Is(err, os.ErrNotExist) {
		fmt.Println("moving diff files after cs package run")
		os.Rename("../tmp/got.jsonl", "../tmp/got_cs.jsonl")
		os.Rename("../tmp/want.jsonl", "../tmp/want_cs.jsonl")
		os.Rename("../tmp/diff.jsonl", "../tmp/diff_cs.jsonl")
	}
}
