package db

import (
	"os"
	"testing"
)

// Package-wide initialization function to clean out
// temp directory once before test start
func TestMain(m *testing.M) {
	m.Run()
	os.Rename("../tmp/diff.jsonl", "../tmp/diff_db.jsonl")
}
