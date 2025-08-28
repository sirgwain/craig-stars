//go:build wasi || wasm

package wasm

import (
	"log/slog"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
}

func EnableDebug() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Debug("enabled debug mode")
}
