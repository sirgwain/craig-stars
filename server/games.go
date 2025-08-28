package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"log/slog"
)

// context for /api/games/{id} calls
func (s *server) gameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.contextUserSession(r)
		db := s.contextDb(r)
		// load the game by id from the database
		id, err := s.int64URLParam(r, "id")
		if id == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		game, err := db.GetGame(r.Context(), *id)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if game == nil {
			slog.Error("game not found", slog.Int64("GameID", *id))
			render.Render(w, r, ErrNotFound)
			return
		}

		if game.State != cs.GameStateSetup && game.HostID != user.ID {
			userIsPlayer := false
			for _, player := range game.Players {
				if player.UserID == user.ID {
					userIsPlayer = true
				}
			}

			if !userIsPlayer {
				slog.Error("access denied for game", slog.Int64("GameID", *id), slog.String("User", user.Username))
				render.Render(w, r, ErrForbidden)
				return
			}
		}

		if (r.Method == "POST" || r.Method == "PUT") && (game.State == cs.GameStateGeneratingTurn || game.State == cs.GameStateGeneratingUniverse) {
			err := fmt.Errorf("game is generating universe or new turn, cannot update")
			slog.Error("update game during turn generation", slog.Any("error", err), slog.Int64("GameID", *id))
			render.Render(w, r, ErrConflict(err))
			return
		}

		ctx := context.WithValue(r.Context(), keyGame, game)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) contextGame(r *http.Request) *cs.GameWithPlayers {
	return r.Context().Value(keyGame).(*cs.GameWithPlayers)
}

func (s *server) computeSpecs(w http.ResponseWriter, r *http.Request) {
	readWriteClient := s.contextDb(r)
	user := s.contextUserSession(r)
	game := s.contextGame(r)

	// validate
	if user.ID != game.HostID {
		render.Render(w, r, ErrForbidden)
		return
	}

	fg, err := readWriteClient.GetFullGame(r.Context(), game.ID)
	if err != nil {
		slog.Error("load full game", slog.Any("error", err))
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	gamer := cs.NewGamer()
	if err := gamer.ComputeSpecs(fg); err != nil {
		slog.Error("compute specs", slog.Any("error", err))
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if err := s.db.WrapInTransaction(func(c db.Client) error {
		return c.UpdateFullGame(r.Context(), fg)
	}); err != nil {
		slog.Error("update game in database", slog.Any("error", err))
		render.Render(w, r, ErrInternalServerError(err))
		return
	}
	rest.RenderJSON(w, rest.JSON{"game": fg.Game})
}
