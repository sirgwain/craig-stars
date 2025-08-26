//go:build !wasi && !wasm

package test

import "fmt"

// TestingT is an interface wrapper around *testing.T.
type TestingT interface {
	Errorf(msg string, args ...any)
	Fatalf(msg string, args ...any)
	Name() string
}

type tHelper interface {
	Helper()
}

type mockTestingT struct {
	name     string
	failed   bool
	canceled bool
}

func (m *mockTestingT) Errorf(msg string, args ...any) {
	m.failed = true
	fmt.Println("MOCK TEST: \n" + fmt.Sprintf(msg, args...))
}

func (m *mockTestingT) Fatalf(msg string, args ...any) {
	m.failed = true
	m.canceled = true
	fmt.Println("MOCK TEST: \n" + fmt.Sprintf(msg, args...))
}

func (m *mockTestingT) Name() string {
	return m.name
}
