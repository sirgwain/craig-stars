package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
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
			log.Error().Int64("GameID", *id).Msg("game not found")
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
				log.Error().Int64("GameID", *id).Str("User", user.Username).Msg("access denied for game")
				render.Render(w, r, ErrForbidden)
				return
			}
		}

		if (r.Method == "POST" || r.Method == "PUT") && (game.State == cs.GameStateGeneratingTurn || game.State == cs.GameStateGeneratingUniverse) {
			err := fmt.Errorf("game is generating universe or new turn, cannot update")
			log.Error().Err(err).Int64("GameID", *id).Msg("update game during turn generation")
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
		log.Error().Err(err).Msg("load full game")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	gamer := cs.NewGamer()
	if err := gamer.ComputeSpecs(fg); err != nil {
		log.Error().Err(err).Msg("compute specs")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if err := s.db.WrapInTransaction(func(c db.Client) error {
		return c.UpdateFullGame(r.Context(), fg)
	}); err != nil {
		log.Error().Err(err).Msg("update game in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}
	rest.RenderJSON(w, rest.JSON{"game": fg.Game})
}
