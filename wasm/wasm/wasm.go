//go:build wasi || wasm

// code heavily inspired by golang-wasm
package wasm

import (
	"fmt"
	"runtime/debug"
	"syscall/js"

	"github.com/rs/zerolog/log"
)

// Magic values to communicate with the JS library.
const (
	globalIdent = "__go_wasm__"
	readyHint   = "__ready__"
)

var (
	bridge js.Value
)

func init() {
	bridge = js.Global().Get(globalIdent)
	if bridge.IsUndefined() || bridge.IsNull() {
		panic("JS wrapper " + globalIdent + " not found")
	}

}

// Ready notifies the JS bridge that the WASM is ready.
// It should be called when every value and function is exposed.
func Ready() {
	bridge.Set(readyHint, js.ValueOf(true))
}

// jsFunctionWrapper wraps a wasm function to capture panics and dump them to the console, but not crash the wasm
func jsFunctionWrapper(wrappedFunc func(args []js.Value) interface{}) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		defer func() {
			// if we panic, update the game state to fail
			if r := recover(); r != nil {
				log.Error().
					Interface("info", r).
					Msgf("%s", debug.Stack())

				js.Global().Set("wasmError", NewError(fmt.Errorf("wasm module failed to run, check console")))
			}
		}()

		return wrappedFunc(args)
	})
}

// Expose exposes a copy of the provided value in JS.
func ExposeFunction(property string, wrappedFunc func(args []js.Value) interface{}) {
	bridge.Set(property, jsFunctionWrapper(wrappedFunc))
}
