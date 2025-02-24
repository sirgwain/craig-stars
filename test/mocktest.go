package test

import "fmt"

// TestingT is an interface wrapper around *testing.T.
type TestingT interface {
	Errorf(msg string, args ...any)
	Name() string
}

type tHelper interface {
	Helper()
}

type mockTestingT struct {
	name   string
	failed bool
}

func (m *mockTestingT) Errorf(msg string, args ...any) {
	m.failed = true
	fmt.Printf(msg, args...)
}

func (m *mockTestingT) Name() string {
	return m.name
}
