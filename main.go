package main

import (
	"log/slog"
	"os"

	"github.com/fabien-marty/slog-helpers/pkg/stacktrace"
	"github.com/phsym/console-slog"
	"github.com/sirgwain/craig-stars/cmd"
)

func main() {

	// pretty console handler
	consoleHandler := console.NewHandler(os.Stderr, &console.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// Wrap with stacktrace handler
	stackHandler := stacktrace.New(consoleHandler, &stacktrace.Options{
		HandlerOptions: slog.HandlerOptions{
			AddSource: true, // include caller file/line
		},
		Mode: stacktrace.ModePrintWithColors, // or ModeAddAttr for structured attr
	})

	// Set as default
	logger := slog.New(stackHandler)
	slog.SetDefault(logger)

	cmd.Execute()
}
