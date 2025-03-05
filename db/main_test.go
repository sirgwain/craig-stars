package db

import (
	"fmt"
	"os"
	"testing"
)

// Package-wide initialization function to clean out
// temp directory once before test start
func TestMain(m *testing.M) {
	m.Run()
	fmt.Println("moving diff files after db package run")
	os.Rename("../tmp/got.jsonl", "../tmp/got_db.jsonl")
	os.Rename("../tmp/want.jsonl", "../tmp/want_db.jsonl")
	os.Rename("../tmp/diff.jsonl", "../tmp/diff_db.jsonl")

}
