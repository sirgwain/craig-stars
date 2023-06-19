package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	usage := "A server running in go with an embedded svelte front end"
	cmd := rootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"-h"})
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), usage) {
		t.Fatalf("expected \"%s\" got \"%s\"", usage, string(out))
	}
}
