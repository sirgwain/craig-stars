//go:build wasi || wasm

package main

import (
	wasm "github.com/sirgwain/craig-stars/wasm/wasm"
)

func main() {

	wasm.ExposeFunction("handleGrpc", wasm.HandleGRPCCall)
	wasm.Ready()

	// fmt.Println("wasm initialized")
	<-make(chan bool) // To use anything from Go WASM, the program may not exit.
}
