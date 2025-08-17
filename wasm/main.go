//go:build wasi || wasm

package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/sirgwain/craig-stars/wasm/wasm"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	wasm.ExposeFunction("handleGrpc", wasm.HandleGRPCCall)
	wasm.Ready()

	// fmt.Println("wasm initialized")
	<-make(chan bool) // To use anything from Go WASM, the program may not exit.
}
