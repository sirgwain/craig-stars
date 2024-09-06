//go:build wasi || wasm
//go:generate go run github.com/sirgwain/craig-stars/wasm/generator ../../cs ./converter.go
// code heavily inspired by golang-wasm
package wasm

import (
	"syscall/js"
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

// Expose exposes a copy of the provided value in JS.
func ExposeFunction(property string, x js.Func) {
	bridge.Set(property, x)
}
