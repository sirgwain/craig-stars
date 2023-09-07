package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
)

// only allow admin requests through
func (s *server) adminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.contextUser(r)

		if !user.isAdmin() {
			log.Error().Str("User", user.Username).Msg("only admins can view all games")
			render.Render(w, r, ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) allGames(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)

	games, err := db.GetGamesWithPlayers()
	if err != nil {
		log.Error().Err(err).Msg("get games from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, games)
}

func (s *server) users(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	db := s.contextDb(r)

	users, err := db.GetUsers()
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("get users from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, users)
}
