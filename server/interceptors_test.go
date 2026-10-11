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
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestWithPlayerRevision(t *testing.T) {
	ctx := t.Context()
	dbConn := db.NewConn()
	cfg := config.Config{}
	cfg.Database.Filename = ":memory:"
	require.NoError(t, dbConn.Connect(ctx, &cfg))
	t.Cleanup(func() { dbConn.Close() })
	client := dbConn.NewReadWriteClient()
	user, err := client.CreateUser(ctx, cs.NewUser("host", "", "", cs.RoleUser))
	require.NoError(t, err)
	fullGame, err := NewGameRunner(dbConn, cfg).HostGame(user.ID, cs.NewGameSettings().WithHost(cs.Humanoids()))
	require.NoError(t, err)
	game, err := client.GetGame(ctx, fullGame.ID)
	require.NoError(t, err)
	gamePlayer := &game.Players[0]
	ctx = context.WithValue(ctx, keyDbRead, dbConn.NewReadClient())
	ctx = context.WithValue(ctx, keyDbWrite, client)

	called := 0
	var handlerErr error
	call := func(clientRevision string, readOnly bool) (string, error) {
		req := connect.NewRequest(&struct{}{})
		if clientRevision != "" {
			req.Header().Set(playerRevisionHeader, clientRevision)
		}
		res, err := withPlayerRevision(ctx, req, func(context.Context, connect.AnyRequest) (connect.AnyResponse, error) {
			called++
			if handlerErr != nil {
				return nil, handlerErr
			}
			return connect.NewResponse(&struct{}{}), nil
		}, gamePlayer, readOnly)
		if err != nil {
			var connectErr *connect.Error
			require.ErrorAs(t, err, &connectErr)
			return connectErr.Meta().Get(playerRevisionHeader), err
		}
		return res.Header().Get(playerRevisionHeader), nil
	}

	// reads report the revision without changing it
	revision, err := call("", true)
	require.NoError(t, err)
	assert.Equal(t, "0", revision)

	// a change from a client with the current revision increments it
	revision, err = call("0", false)
	require.NoError(t, err)
	assert.Equal(t, "1", revision)

	// a change from another device that still has the old revision is rejected before the handler runs
	called = 0
	revision, err = call("0", false)
	assert.Equal(t, connect.CodeAborted, connect.CodeOf(err))
	assert.Equal(t, "1", revision)
	assert.Zero(t, called)

	// but that device can still read, and learns about the newer revision
	revision, err = call("0", true)
	require.NoError(t, err)
	assert.Equal(t, "1", revision)

	// failed changes don't increment the revision
	handlerErr = connect.NewError(connect.CodeInvalidArgument, errors.New("invalid orders"))
	revision, err = call("1", false)
	assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	assert.Equal(t, "1", revision)
	handlerErr = nil

	// clients that don't send a revision aren't checked, but their changes still increment it
	revision, err = call("", false)
	require.NoError(t, err)
	assert.Equal(t, "2", revision)
}
