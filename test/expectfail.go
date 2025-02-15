package test

import (
	"strings"
	"sync"
	"testing"
)

// Run testFunc and expect it to fail the test,
// marking the entire test as having failed if this does not occur.
func ExpectFailure(t *testing.T, testFunc func(t *testing.T)) {
	t.Helper()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		testFunc(t)
	}()

	wg.Wait()

	if !t.Failed() {
		// split around the Test_XXX prefix
		prefix, tNameCut, found := strings.Cut(t.Name(), "_")
		if !found {
			// cover for really weird jank with test names that
			// *PHYSICALLY SHOULD NEVER HAPPEN*
			tNameCut = prefix
		}
		// extract subtest
		tMaintest, tSubtest, found := strings.Cut(tNameCut, "/")
		if found {
			tSubtest = "subtest (" + tSubtest + ")"
		}
		strings.Cut(tNameCut, "/")
		t.Errorf("%s() %s did not fail when expected to", tMaintest, tSubtest)
	}
}
