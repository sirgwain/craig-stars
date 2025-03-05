package ai

import (
	"fmt"
	"os"
	"testing"
)

// Package-wide initialization function to move temp json files.
func TestMain(m *testing.M) {
	m.Run()
	fmt.Println("moving diff files after ai package run")
	os.Rename("../tmp/got.jsonl", "../tmp/got_ai.jsonl")
	os.Rename("../tmp/want.jsonl", "../tmp/want_ai.jsonl")
	os.Rename("../tmp/diff.jsonl", "../tmp/diff_ai.jsonl")
}
