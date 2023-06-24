package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type hostGameRequest struct {
	*cs.GameSettings
}

func (req *hostGameRequest) Bind(r *http.Request) error {
	return nil
}

type joinGameRequest struct {
	RaceID int64  `json:"raceId"`
	Color  string `json:"color"`
}

func (req *joinGameRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id} calls
func (s *server) gameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// load the game by id from the database
		id, err := s.int64URLParam(r, "id")
		if id == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		game, err := s.db.GetGame(*id)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if game == nil {
			log.Error().Int64("GameID", *id).Msg("game not found")
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), keyGame, game)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) contextGame(r *http.Request) *cs.Game {
	return r.Context().Value(keyGame).(*cs.Game)
}

func (s *server) games(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)

	games, err := s.db.GetGamesForUser(user.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("get games from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, games)
}

func (s *server) hostedGames(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)

	games, err := s.db.GetGamesForHost(user.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("get games from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, games)
}

func (s *server) openGames(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)

	games, err := s.db.GetOpenGames(user.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("get games from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, games)
}

func (s *server) openGamesByHash(w http.ResponseWriter, r *http.Request) {
	// load open games by hash from the database
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("invalid invite hash in url")))
		return
	}

	games, err := s.db.GetOpenGamesByHash(hash)
	if err != nil {
		log.Error().Err(err).Str("Hash", hash).Msg("get open games by hash from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, games)
}

func (s *server) game(w http.ResponseWriter, r *http.Request) {
	// TODO: any need to prevent a non-player from loading this game?
	game := s.contextGame(r)
	rest.RenderJSON(w, game)
}

// Host a new game
func (s *server) createGame(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)

	settings := hostGameRequest{}
	if err := render.Bind(r, &settings); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	game, err := s.gameRunner.HostGame(user.ID, settings.GameSettings)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msgf("host game %v", settings.GameSettings)
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, game)
}

// Join an open game
func (s *server) joinGame(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	game := s.contextGame(r)

	join := joinGameRequest{}
	if err := render.Bind(r, &join); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// try and join this game
	if err := s.gameRunner.JoinGame(game.ID, user.ID, join.RaceID, join.Color); err != nil {
		log.Error().Err(err).Msg("join game")
		render.Render(w, r, ErrBadRequest(err))
	}
}

// Generate a universe for a host
func (s *server) generateUniverse(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	game := s.contextGame(r)

	// validate
	if user.ID != game.HostID {
		render.Render(w, r, ErrForbidden)
		return
	}

	if err := s.gameRunner.GenerateUniverse(game); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Msg("generating universe")
		render.Render(w, r, ErrInternalServerError(err))
	}

	// send the full game to the host
	s.renderFullPlayerGame(w, r, game.ID, user.ID)
}

func (s *server) generateTurn(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	game := s.contextGame(r)

	// validate
	if user.ID != game.HostID {
		render.Render(w, r, ErrForbidden)
		return
	}

	if err := s.gameRunner.GenerateTurn(game.ID); err != nil {
		log.Error().Err(err).Msg("generate turn")
		render.Render(w, r, ErrBadRequest(err))
	}
}

func (s *server) deleteGame(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	game := s.contextGame(r)

	if game.HostID != user.ID {
		log.Error().Int64("ID", game.ID).Int64("UserID", user.ID).Msg("only host can delete game")
		render.Render(w, r, ErrBadRequest(fmt.Errorf("only host can delete game")))
		return
	}

	if err := s.db.DeleteGame(game.ID); err != nil {
		log.Error().Err(err).Int64("ID", game.ID).Msg("delete game from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}
}
