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

// Log implement sqldblogger.Logger and log it as is.
// To use context.Context values, please copy this file and adjust to your needs.
func (zl *slogAdapter) Log(_ context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
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

	attrs := make([]slog.Attr, len(data))
	i := 0
	for k, v := range data {
		attrs[i] = slog.Any(k, v)
		i++
	}

	zl.logger.LogAttrs(context.Background(), lvl, msg, attrs...)
}
