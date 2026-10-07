package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
)

func TestErrorLogInterceptorCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	tests := []struct {
		name  string
		ctx   context.Context
		err   error
		code  connect.Code
		level string
	}{
		{"canceled query", ctx, context.Canceled, connect.CodeCanceled, "DEBUG"},
		{"wrapped canceled query", ctx, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get games: %w", context.Canceled)), connect.CodeCanceled, "DEBUG"},
		{"real DB error despite canceled request", ctx, connect.NewError(connect.CodeInternal, errors.New("database is locked")), connect.CodeInternal, "ERROR"},
		{"cancellation without canceled request", t.Context(), connect.NewError(connect.CodeInternal, context.Canceled), connect.CodeInternal, "ERROR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			previous := slog.Default()
			slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})))
			t.Cleanup(func() { slog.SetDefault(previous) })
			handler := newErrorLogInterceptor().WrapUnary(func(context.Context, connect.AnyRequest) (connect.AnyResponse, error) {
				return nil, tt.err
			})
			_, err := handler(tt.ctx, connect.NewRequest(&struct{}{}))
			assert.Equal(t, tt.code, connect.CodeOf(err))
			assert.ErrorIs(t, err, tt.err)
			assert.Contains(t, logs.String(), `"level":"`+tt.level+`"`)
		})
	}
}

func TestRequestLoggerCancellation(t *testing.T) {
	tests := []struct {
		name     string
		canceled bool
		status   int
		level    string
	}{
		{"client cancellation", true, 499, "DEBUG"},
		{"real failure after cancellation", true, http.StatusInternalServerError, "ERROR"},
		{"499 without cancellation", false, 499, "ERROR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			previous := slog.Default()
			slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})))
			t.Cleanup(func() { slog.SetDefault(previous) })
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if tt.canceled {
				cancel()
			}
			req := httptest.NewRequest(http.MethodPost, "/api/grpc/test", nil).WithContext(ctx)
			requestLogger()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			})).ServeHTTP(httptest.NewRecorder(), req)
			assert.Contains(t, logs.String(), `"level":"`+tt.level+`"`)
		})
	}
}
