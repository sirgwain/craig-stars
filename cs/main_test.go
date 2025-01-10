package cs

import (
	"os"
	"testing"
)

// Package-wide i0nitialization function to clean out  
// temp directory once before test start
func TestMain(m *testing.M) {
	_ = os.RemoveAll("../tmp")
	_ = os.Mkdir("../tmp", 0644)
	m.Run()
}
