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
	"github.com/sirgwain/craig-stars/cs"
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

func TestResolveGamePlayer(t *testing.T) {
	// user 1 is the host and plays players 1 and 2 (hot seat), user 2 plays player 3, player 4 is an AI
	game := &cs.GameWithPlayers{
		Game: cs.Game{DBObject: cs.DBObject{ID: 1}, HostID: 1, State: cs.GameStateWaitingForPlayers},
		Players: []cs.GamePlayer{
			{UserID: 1, Num: 1},
			{UserID: 1, Num: 2},
			{UserID: 2, Num: 3},
			{AIControlled: true, Num: 4},
		},
	}
	host := userSession{ID: 1, Role: string(cs.RoleUser)}
	player := userSession{ID: 2, Role: string(cs.RoleUser)}
	outsider := userSession{ID: 3, Role: string(cs.RoleUser)}
	admin := userSession{ID: 4, Role: string(cs.RoleAdmin)}

	tests := []struct {
		name        string
		user        userSession
		asPlayerNum int
		readOnly    bool
		wantNum     int
		wantCode    connect.Code
	}{
		{"default player", player, 0, false, 3, 0},
		{"default hot seat player", host, 0, false, 1, 0},
		{"hot seat player", host, 2, false, 2, 0},
		{"other user's player", host, 3, true, 0, connect.CodePermissionDenied},
		{"missing player", host, 5, true, 0, connect.CodeNotFound},
		{"outsider", outsider, 0, true, 0, connect.CodePermissionDenied},
		{"admin view as ai", admin, 4, true, 4, 0},
		{"admin write as ai", admin, 4, false, 0, connect.CodePermissionDenied},
		{"admin view without player", admin, 0, true, 0, 0},
		{"admin write without player", admin, 0, false, 0, connect.CodePermissionDenied},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gamePlayer, err := resolveGamePlayer(game, tt.user, tt.asPlayerNum, tt.readOnly)
			if tt.wantCode != 0 {
				assert.Equal(t, tt.wantCode, connect.CodeOf(err))
				return
			}
			assert.NoError(t, err)
			if tt.wantNum == 0 {
				assert.Nil(t, gamePlayer)
			} else if assert.NotNil(t, gamePlayer) {
				assert.Equal(t, tt.wantNum, gamePlayer.Num)
			}
		})
	}
}
