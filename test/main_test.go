package test

import (
	"fmt"
	"os"
	"testing"
)

// Package-wide initialization function to clean out
// temp directory once before test start
func TestMain(m *testing.M) {
	m.Run()
	fmt.Println("moving diff files after test package run")
	os.Rename("../tmp/got.jsonl", "../tmp/got_test.jsonl")
	os.Rename("../tmp/want.jsonl", "../tmp/want_test.jsonl")
	os.Rename("../tmp/diff.jsonl", "../tmp/diff_test.jsonl")
}
