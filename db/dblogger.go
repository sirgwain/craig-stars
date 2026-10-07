package db

import (
	"context"
	"log/slog"

	sqldblogger "github.com/simukti/sqldb-logger"
)

type slogAdapter struct {
	logger *slog.Logger
}

func newLoggerWithLogger(l *slog.Logger) sqldblogger.Logger {
	return &slogAdapter{
		logger: l,
	}
}

// Log implements sqldblogger.Logger.
func (zl *slogAdapter) Log(ctx context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	var lvl slog.Level

	switch level {
	case sqldblogger.LevelError:
		lvl = slog.LevelError
	case sqldblogger.LevelInfo:
		lvl = slog.LevelInfo
	case sqldblogger.LevelDebug:
		lvl = slog.LevelDebug
	case sqldblogger.LevelTrace:
		lvl = slog.LevelDebug // slog doesn't have trace, default to debug
	default:
		lvl = slog.LevelDebug
	}
	// sqldb-logger supplies error text rather than the original error value,
	// and uses context.Background for row iteration. Recognize its exact
	// cancellation error while leaving other DB failures at error level.
	if level == sqldblogger.LevelError && data["error"] == context.Canceled.Error() {
		lvl = slog.LevelDebug
	}

	attrs := make([]slog.Attr, len(data))
	i := 0
	for k, v := range data {
		attrs[i] = slog.Any(k, v)
		i++
	}

	zl.logger.LogAttrs(ctx, lvl, msg, attrs...)
}
