package db

import (
	"context"

	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
)

type zerologAdapter struct {
	logger *zerolog.Logger
}

func newLoggerWithLogger(l *zerolog.Logger) sqldblogger.Logger {
	return &zerologAdapter{
		logger: l,
	}
}

// Log implement sqldblogger.Logger and log it as is.
// To use context.Context values, please copy this file and adjust to your needs.
func (zl *zerologAdapter) Log(_ context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	var lvl zerolog.Level

	switch level {
	case sqldblogger.LevelError:
		lvl = zerolog.ErrorLevel
	case sqldblogger.LevelInfo:
		lvl = zerolog.InfoLevel
	case sqldblogger.LevelDebug:
		lvl = zerolog.DebugLevel
	case sqldblogger.LevelTrace:
		lvl = zerolog.TraceLevel
	default:
		lvl = zerolog.DebugLevel
	}

	zl.logger.WithLevel(lvl).Fields(data).Msg(msg)
}
