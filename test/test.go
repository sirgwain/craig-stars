// Package test contains some useful utility functions for testing things.
package test

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/nsf/jsondiff"
)

// Compare two objects as json outputs for testing.
//
// If the comparison fails, this marks the test as a failure and
// writes 3 JSONL files to ./tmp, containing serialized versions of got and want
// and a pretty-printed difference between them courtesy of [github.com/nsf/jsondiff].
//
// These files are continuously appended to during a test run (sectioned off by test name),
// and should be moved or removed after the package finishes testing.
// Invocation from parallel tests is untested and not recommended.
//
// The json difference is passed to [testing.T.Fatalf], so no extra function calls
// should be made after calling this.
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

	diff, err := parseJSONDiff(string(gotJson), string(wantJson), t.Name())
	if err != nil {
		t.Fatalf("error creating JSON diffs: \n%v", err)
	}

	// replace comments
	r := strings.NewReplacer("/* ", "", " */", ":")

	t.Fatalf("JSONs not equal; diff between got & want: \n%s", r.Replace(diff))
}

// Parsing options for jsondiff.
// Fun fact: this is guaranteed to produce valid JSONL syntax.
// Block comments are removed in the stdout version since we don't care about syntax there.
var options = jsondiff.Options{
	Added:            jsondiff.Tag{Begin: "/* Added */ ", End: ""},
	Removed:          jsondiff.Tag{Begin: "/* Removed */ ", End: ""},
	Changed:          jsondiff.Tag{Begin: "/* Changed */ [ ", End: " ]"},
	ChangedSeparator: ", ",
	Indent:           "\t", // tab indentation

	SkipMatches: true,
}

// Parse JSON diffs, creating files to log values as appropriate.
func parseJSONDiff(gotJSON, wantJSON, testName string) (diff string, err error) {
	// compare diff without whitespace because jsondiff REALLY hates whitespace
	_, diff = jsondiff.Compare([]byte(gotJSON), []byte(wantJSON), &options)

	os.MkdirAll("../tmp", 0755)
	for i := range 3 {
		var path string
		var body string
		switch i {
		case 0:
			path = "../tmp/got.jsonl"
			body = gotJSON
		case 1:
			path = "../tmp/want.jsonl"
			body = wantJSON
		case 2:
			path = "../tmp/diff.jsonl"
			body = diff
		}

		header := "// " + testName + "\n" // header containing test name & extra newlines
		if _, err := os.Stat(path); err == nil {
			// add extra newline in header to properly delimit sections on existing files
			header = "\n" + header
		}
		if err = AppendFile(path, header+body+"\n"); err != nil {
			return "", err
		}
	}

	return diff, nil
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
