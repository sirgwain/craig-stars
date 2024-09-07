//go:build wasi || wasm

package wasm

import (
	"syscall/js"

	"github.com/rs/zerolog/log"
)

// NewError returns a JS Error with the provided Go error's error message.
func NewError(err error) js.Value {
	log.Error().Err(err).Msg("")
	errFunc := js.Global().Get("Error")

	// create a new js Error object
	return errFunc.New(err.Error())
}
