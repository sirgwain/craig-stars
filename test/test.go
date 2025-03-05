// Package test contains some useful utility functions for testing things.
package test

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/nsf/jsondiff"
)

// Compare two objects as json outputs for testing.
//
// If the comparison fails, this marks the test as a failure
// and writes a json file to the tmp folder containing a pretty-printed
// difference between the 2 values.
//
// The file is continuously appended to during a test run (sectioned off by test name),
// and should ideally be moved or removed after the package finishes testing.
// Invocation from parallel tests is untested and not recommended.
//
// The json difference is passed to t.Errorf, so no extra calls to t.Log or t.Error
// are needed after calling this.
func CompareAsJSON(t TestingT, got, want any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if got == nil && want == nil {
		return
	} else if (got == nil) != (want == nil) { // one is nil and the other isn't
		t.Fatalf("Unequal values (nilness): got = %v, want = %v", got, want)
	}

	gotJson, err := json.MarshalIndent(got, "", "\t")
	if err != nil {
		t.Fatalf("CompareAsJSON could not marshal got (%q) to json: \n%v", got, err)
	}
	wantJson, err := json.MarshalIndent(want, "", "\t")
	if err != nil {
		t.Fatalf("CompareAsJSON could not marshal want (%q) to json: \n%v", want, err)
	}

	if string(gotJson) == string(wantJson) {
		return
	}

	diff := parseJSONDiff(gotJson, wantJson, t.Name())

	t.Fatalf("JSONs not equal; diff between got & want: \n%s", diff)
}

// parsing options for jsondiff
var options = jsondiff.Options{
	Added:            jsondiff.Tag{Begin: "\"prop-added\": {", End: "}"},
	Removed:          jsondiff.Tag{Begin: "\"prop-removed\": {", End: "}"},
	Changed:          jsondiff.Tag{Begin: "{\"changed\": [", End: "]}"},
	ChangedSeparator: ", ",
	Indent:           "\t", // tab indentation
	SkipMatches:      true,
}

// Parse JSON diffs, creating files to log values as appropriate.
func parseJSONDiff(gotJSON, wantJSON []byte, testName string) string {
	_, diff := jsondiff.Compare(gotJSON, wantJSON, &options)

	os.MkdirAll("../tmp", 0755) // create temp folder
	// append files 1 by 1
	header := "// " + testName + "\n" // header containing test name & extra newlines
	path := "../tmp/diff.jsonl"
	if _, err := os.Stat(path); err == nil {
		// add extra newline in header to properly delimit sections
		header = "\n" + header
	}
	_ = AppendFile(path, header+diff+"\n")

	return diff
}

// Appends a string or byte slice to the named file, creating it if necessary.
func AppendFile[S ~string | ~[]byte](path string, data S) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %q: \n%w", path, err)
	}
	defer f.Close()

	if _, err := f.Write([]byte(data)); err != nil {
		return fmt.Errorf("error appending data to file %q: \n%w", path, err)
	}
	return nil
}

// compare two floats within a tolerance range
// source: Ricardo Gerardi - https://medium.com/pragmatic-programmers/testing-floating-point-numbers-in-go-9872fe6de17f
func WithinTolerance(a, b, tolerance float64) bool {
	if a == b {
		return true
	}
	diff := math.Abs(a - b)
	if b == 0 {
		return diff < tolerance
	}
	return (diff / math.Abs(b)) < tolerance
}
