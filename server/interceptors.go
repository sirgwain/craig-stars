//go:build !wasi && !wasm

package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

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
				log.Error().
					Err(err).
					Str("Procedure", req.Spec().Procedure).
					Msg("grpc call failed")
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

				userIsPlayer := false
				for _, player := range game.Players {
					if player.UserID == user.ID {
						userIsPlayer = true
						ctx = context.WithValue(ctx, keyGamePlayer, &player)
						break
					}
				}

				// if we aren't in setup mode and this user doesn't have a player yet, error
				if game.State != cs.GameStateSetup && !userIsPlayer && game.HostID != user.ID {
					return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("access denied for game %d", gameID))
				}

				// update context with game
				ctx = context.WithValue(ctx, keyGame, game)
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
