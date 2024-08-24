//go:build !wasi && !wasm

package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/sirgwain/craig-stars/cmd"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	cmd.Execute()
}
