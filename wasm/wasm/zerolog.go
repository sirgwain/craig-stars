//go:build !tinygo && (wasi || wasm)

package wasm

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func EnableDebug() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime, NoColor: true})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Debug().Msg("enabled debug mode")
}
