//go:build wasi || wasm

package wasm

import (
	"syscall/js"
)

// NewError returns a JS Error with the provided Go error's error message.
func NewError(err error) js.Value {
	errFunc := js.Global().Get("Error")

	// create a new js Error object
	return errFunc.New(err.Error())
}
