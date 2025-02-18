//go:build !wasi && !wasm

package cmd

import (
	"bytes"
	"io"
	"testing"
)

func TestExecute(t *testing.T) {

	cmd := rootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"-h"})
	cmd.Execute()
	out, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) == "" {
		t.Fatalf("no usage")
	}
}
