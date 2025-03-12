package cs

import (
	"os"
	"testing"
)

// Package-wide initialization function to clean out
// temp directory once before test start
func TestMain(m *testing.M) {
	m.Run()
	os.Rename("../tmp/got.jsonl", "../tmp/got_cs.jsonl")
	os.Rename("../tmp/want.jsonl", "../tmp/want_cs.jsonl")
	os.Rename("../tmp/diff.jsonl", "../tmp/diff_cs.jsonl")
}
