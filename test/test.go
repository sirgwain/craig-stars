// Package test contains some useful utility functions for testing things.
package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/nsf/jsondiff"
)

// Compare two objects as json outputs for testing.
//
// If the comparison fails, this marks the test as a failure
// and writes 3 json files to the tmp folder,
// containing both values being compared and a pretty-printed
// difference between them.
//
// These files are continuously appended to during a test run (sectioned off by test name),
// and should ideally be moved or removed after the package finishes testing.
// Invocation from parallel tests is untested and not recommended.
//
// The json difference is passed to t.Errorf, so no extra calls to t.Log or t.Error
// are needed after calling this.
func CompareAsJSON(t TestingT, got, want any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if got == nil && want == nil {
		return true
	} else if (got == nil) != (want == nil) { // one is nil and the other isn't
		t.Errorf("Unequal values (nilness): got = %v, want = %v", got, want)
		return false
	}

	gotJson, err := json.MarshalIndent(got, "", "\t")
	if err != nil {
		t.Errorf("compareAsJSON could not marshal got (%q) to json: \n%v", got, err)
		return false
	}
	wantJson, err := json.MarshalIndent(want, "", "\t")
	if err != nil {
		t.Errorf("compareAsJSON could not marshal want (%q) to json: \n%v", want, err)
		return false
	}

	if string(gotJson) == string(wantJson) {
		return true
	}

	diff := parseJSONDiff(gotJson, wantJson, t.Name())

	t.Errorf("JSONs not equal; diff between got & want: \n%s", diff)
	return false
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
	for i := range 3 {
		header := "// " + testName + "\n" // header containing test name & extra newlines
		var path, body string
		switch i {
		case 0:
			path = "../tmp/got.jsonl"
			body = string(gotJSON)
		case 1:
			path = "../tmp/want.jsonl"
			body = string(wantJSON)
		case 2:
			path = "../tmp/diff.jsonl"
			body = diff
		}
		if FileExists(path) {
			// add extra newline in header to properly delimit sections
			header = "\n" + header
		}
		_ = AppendFile(path, header+body+"\n")
	}

	return diff
}

// Appends a string or byte slice to the named file, creating it if necessary.
func AppendFile[S ~string | ~[]byte](path string, data S) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file %q; error: \n%w", path, err)
	}
	defer f.Close()

	if _, err := f.Write([]byte(data)); err != nil {
		return fmt.Errorf("could not append data to file %q; error: \n%w", path, err)
	}
	return nil
}

// FileExists reports whether a file at path exists or not.
// It does not actually open the file or modify it in any way.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// Check for the existence of an expected or unexpected error within a test,
// failing the test as appropriate.
func CheckUnexpectedError(t *testing.T, err error, wantErr bool) {
	t.Helper()
	if (err != nil) == wantErr {
		return
	}

	errMsg := fmt.Sprintf("%s() errored unexpectedly;\n", t.Name())
	if err != nil {
		errMsg += fmt.Sprintf("test produced error \"%v\" despite expecting none", err)
	} else {
		errMsg += "test failed to error when expected to"
	}
	t.Error(errMsg)
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
