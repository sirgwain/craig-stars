package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"log/slog"
)

// header used to pick which player a game request acts as, for hot seat games or admin viewing
const asPlayerHeader = "X-As-Player"

// header sent by clients viewing a game in read-only mode, i.e. viewing a submitted turn
const readOnlyHeader = "X-Read-Only"

// header with the revision of the player's orders. Every change a player makes to a game increments it.
// Clients send back the revision they loaded so changes based on stale data, i.e. the player made
// more recent changes on another device, are rejected.
const playerRevisionHeader = "X-Player-Revision"

type GameIdRequest interface {
	GetGameId() int64
}

func contextUserSession(ctx context.Context) userSession {
	return ctx.Value(keyUserSession).(userSession)
}

func contextDb(ctx context.Context) db.ReadClient {
	return ctx.Value(keyDbRead).(db.ReadClient)
}

func contextDbWrite(ctx context.Context) db.Client {
	return ctx.Value(keyDbWrite).(db.Client)
}

func contextGame(ctx context.Context) *cs.GameWithPlayers {
	return ctx.Value(keyGame).(*cs.GameWithPlayers)
}

func contextGamePlayer(ctx context.Context) *cs.GamePlayer {
	return ctx.Value(keyGamePlayer).(*cs.GamePlayer)
}

func newErrorLogInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			res, err := next(ctx, req)
			if err != nil {
				// Navigating away can cancel an in-flight query. Preserve cancellation
				// through service errors that otherwise wrap every DB error as internal.
				if ctx.Err() == context.Canceled && errors.Is(err, context.Canceled) {
					slog.DebugContext(ctx, "grpc call canceled",
						slog.Any("error", err),
						slog.String("Procedure", req.Spec().Procedure))
					return res, connect.NewError(connect.CodeCanceled, err)
				}
				slog.Error("grpc call failed",
					slog.Any("error", err),
					slog.String("Procedure", req.Spec().Procedure))
			}
			return res, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func newDbInterceptor(db DBConnection) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			ctx = context.WithValue(ctx, keyDbRead, db.NewReadClient())
			ctx = context.WithValue(ctx, keyDbWrite, db.NewReadWriteClient())

			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func newGameInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			dbClient := contextDb(ctx)
			user := contextUserSession(ctx)

			if r, ok := req.Any().(GameIdRequest); ok {
				gameID := r.GetGameId()
				game, err := dbClient.GetGame(ctx, gameID)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load game %d from database %w", gameID, err))
				}

				if game == nil {
					return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("game not found %d", gameID))
				}

				asPlayerNum, err := parseAsPlayerHeader(req.Header())
				if err != nil {
					return nil, connect.NewError(connect.CodeInvalidArgument, err)
				}

				readOnlyRequest := req.Spec().IdempotencyLevel == connect.IdempotencyNoSideEffects
				if req.Header().Get(readOnlyHeader) == "true" && !readOnlyRequest {
					return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("game %d is being viewed read-only", gameID))
				}

				gamePlayer, err := resolveGamePlayer(game, user, asPlayerNum, readOnlyRequest)
				if err != nil {
					return nil, err
				}

				// update context with game
				ctx = context.WithValue(ctx, keyGame, game)

				if gamePlayer != nil {
					ctx = context.WithValue(ctx, keyGamePlayer, gamePlayer)
					return withPlayerRevision(ctx, req, next, gamePlayer, readOnlyRequest)
				}
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

// withPlayerRevision calls next, rejecting changes from clients that loaded an older revision of the
// player's orders than the one in the database. The current revision is returned in the response
// headers so clients know when they need to reload.
func withPlayerRevision(ctx context.Context, req connect.AnyRequest, next connect.UnaryFunc, gamePlayer *cs.GamePlayer, readOnlyRequest bool) (connect.AnyResponse, error) {
	// read the revision before loading any data so a client never thinks its data is newer than it is
	revision, err := contextDb(ctx).GetPlayerRevision(ctx, gamePlayer.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load player revision %w", err))
	}

	// clients that don't track revisions, i.e. the API, aren't checked
	if clientRevision := req.Header().Get(playerRevisionHeader); !readOnlyRequest && clientRevision != "" && clientRevision != strconv.FormatInt(revision, 10) {
		err := connect.NewError(connect.CodeAborted, fmt.Errorf("%s has newer changes from another device, reload to get the latest", gamePlayer.Name))
		err.Meta().Set(playerRevisionHeader, strconv.FormatInt(revision, 10))
		return nil, err
	}

	res, err := next(ctx, req)
	if err == nil && !readOnlyRequest {
		if revision, err = contextDbWrite(ctx).IncrementPlayerRevision(ctx, gamePlayer.ID); err != nil {
			// the player could have been removed by this request, i.e. leaving a game
			slog.Debug("increment player revision", slog.Any("error", err), slog.Int64("ID", gamePlayer.ID))
			return res, nil
		}
	}

	if err != nil {
		var connectErr *connect.Error
		if errors.As(err, &connectErr) {
			connectErr.Meta().Set(playerRevisionHeader, strconv.FormatInt(revision, 10))
		}
		return res, err
	}
	res.Header().Set(playerRevisionHeader, strconv.FormatInt(revision, 10))
	return res, nil
}

// parseAsPlayerHeader returns the player num requested in the X-As-Player header, or 0 if not set
func parseAsPlayerHeader(header http.Header) (int, error) {
	value := header.Get(asPlayerHeader)
	if value == "" {
		return 0, nil
	}
	num, err := strconv.Atoi(value)
	if err != nil || num < 1 {
		return 0, fmt.Errorf("invalid %s header %q", asPlayerHeader, value)
	}
	return num, nil
}

// resolveGamePlayer finds the player a request acts as. By default this is the user's player, but users
// controlling multiple players (hot seat) or admins can choose one with the X-As-Player header.
// Admins can view games they aren't part of, or view as players they don't control, but only read.
func resolveGamePlayer(game *cs.GameWithPlayers, user userSession, asPlayerNum int, readOnlyRequest bool) (*cs.GamePlayer, error) {
	var gamePlayer *cs.GamePlayer
	impersonating := false
	for i := range game.Players {
		player := &game.Players[i]
		if asPlayerNum == 0 {
			if player.UserID == user.ID {
				gamePlayer = player
				break
			}
			continue
		}
		if player.Num == asPlayerNum {
			if player.UserID == user.ID {
				gamePlayer = player
			} else if user.isAdmin() {
				gamePlayer = player
				impersonating = true
			} else {
				return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("access denied for player %d in game %d", asPlayerNum, game.ID))
			}
			break
		}
	}

	if asPlayerNum != 0 && gamePlayer == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player %d not found in game %d", asPlayerNum, game.ID))
	}

	readOnly := impersonating
	if game.State != cs.GameStateSetup && gamePlayer == nil && game.HostID != user.ID {
		if !user.isAdmin() {
			return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("access denied for game %d", game.ID))
		}
		readOnly = true
	}

	if readOnly && !readOnlyRequest {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("admin access to game %d is read-only", game.ID))
	}

	return gamePlayer, nil
}
