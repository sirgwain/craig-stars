package test

import (
	"encoding/json"
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
// These files are continuously appended to during a test run (sectioned off by test name), 
// and are cleaned out before a new test starts.
//
// The json difference is printed to stderr as well, so no extra calls to t.Log or t.Error
// are needed after calling this.
func CompareAsJSON(t *testing.T, got, want any) {
	if got == nil && want == nil {
		return
	} else if (got == nil) != (want == nil) { // one is nil and the other isn't
		t.Errorf("Unequal values (nilness): got = %v, want = %v", got, want)
	}

	gotJson, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		t.Errorf("compareAsJSON could not marshal got (%q) to json, error = %v", got, err)
	}
	wantJson, err := json.MarshalIndent(want, "", "  ")
	if err != nil {
		t.Errorf("compareAsJSON could not marshal want (%q) to json, error = %v", want, err)
	}

	if string(gotJson) == string(wantJson) {
		return
	}

	options := jsondiff.Options{
		Added:            jsondiff.Tag{Begin: "\"prop-added\": {", End: "}"},
		Removed:          jsondiff.Tag{Begin: "\"prop-removed\": {", End: "}"},
		Changed:          jsondiff.Tag{Begin: "{\"changed\": [", End: "]}"},
		ChangedSeparator: ", ",
		Indent:           "    ",
		SkipMatches:      true,
	}

	_, diff := jsondiff.Compare(gotJson, wantJson, &options)

	header := []byte("// " + t.Name() + "\n") // header containing test name & extra newlines
	_ = AppendFile("../tmp/got.jsonl", append(append(header, gotJson...), "\n\n"...))
	_ = AppendFile("../tmp/want.jsonl", append(append(header, wantJson...), "\n\n"...))
	_ = AppendFile("../tmp/diff.jsonl", append(append(header, diff...), "\n\n"...))

	t.Errorf("JSONs not equal; diff between got & want: \n%s", diff)
}

// Appends data to the named file, creating it if necessary.
func AppendFile[S ~string | ~[]byte](name string, data S) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file %s; error: %w", name, err)
	}
	defer f.Close()

	if _, err := f.Write([]byte(data)); err != nil {
		return fmt.Errorf("could not append bytes %q to file %s; error: %w", string(data), name, err)
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
