package cs

import (
	"fmt"
	"os"
	"testing"
)

// Package-wide initialization function to clean out
// temp directory once before test start
func TestMain(m *testing.M) {
	m.Run()
	fmt.Println("moving diff files after cs package run")
	os.Rename("../tmp/got.jsonl", "../tmp/got_cs.jsonl")
	os.Rename("../tmp/want.jsonl", "../tmp/want_cs.jsonl")
	os.Rename("../tmp/diff.jsonl", "../tmp/diff_cs.jsonl")
}
