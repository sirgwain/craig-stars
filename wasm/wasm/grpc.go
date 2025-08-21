//go:build wasi || wasm

package wasm

import (
	"context"
	"syscall/js"

	"github.com/rs/zerolog/log"
)

// HandleGRPCCall handles grpc calls through the wasm transport
func HandleGRPCCall(args []js.Value) any {
	go func() {

		method := args[0].String() // e.g. "/craig_stars.v1.RaceService/CalculateRacePoints"
		input := args[1]           // Uint8Array of request body
		cb := args[2]              // JS callback for result (or a promise resolve)

		reqBytes := make([]byte, input.Get("byteLength").Int())
		js.CopyBytesToGo(reqBytes, input)

		if log.Debug().Enabled() {
			log.Debug().Msgf("calling %s", method)
		}
		resBytes, err := serviceHandler.Call(context.Background(), method, reqBytes)

		if err != nil {
			println("failed to process method", method, "size", len(reqBytes), err)
			jsErr := js.Global().Get("Error").New(err.Error())
			cb.Invoke(jsErr, js.Null()) // reject
			return
		}

		out := js.Global().Get("Uint8Array").New(len(resBytes))
		js.CopyBytesToJS(out, resBytes)
		cb.Invoke(js.Null(), out) // resolve
	}()
	return nil
}
