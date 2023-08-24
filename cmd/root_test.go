package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestExecute(t *testing.T) {

	cmd := rootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"-h"})
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) == "" {
		t.Fatalf("no usage")
	}
}
