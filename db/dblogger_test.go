package db

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/stretchr/testify/assert"
)

func TestDBLoggerCancellation(t *testing.T) {
	tests := []struct {
		name     string
		canceled bool
		err      string
		level    string
	}{
		{"canceled query", true, context.Canceled.Error(), "DEBUG"},
		{"real DB error despite canceled request", true, "database is locked", "ERROR"},
		{"row cancellation without request context", false, context.Canceled.Error(), "DEBUG"},
		{"deadline exceeded remains an error", true, context.DeadlineExceeded.Error(), "ERROR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			logger := newLoggerWithLogger(slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})))
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if tt.canceled {
				cancel()
			}
			logger.Log(ctx, sqldblogger.LevelError, "RowsNext", map[string]interface{}{"error": tt.err})
			assert.Contains(t, logs.String(), `"level":"`+tt.level+`"`)
			assert.Contains(t, logs.String(), tt.err)
		})
	}
}
