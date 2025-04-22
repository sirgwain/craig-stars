//go:build !wasi && !wasm

package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type userSettingsRequest struct {
	UserSettings *cs.UserSettings `json:"userSettings,omitempty"`
}

func (req *userSettingsRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/users/{id} calls
func (s *server) userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := s.contextDb(r)
		userSession := s.contextUserSession(r)

		// load the user by id from the database
		id, err := s.int64URLParam(r, "id")
		if id == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		user, err := db.GetUser(*id)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if user == nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		if user.ID != userSession.ID && !userSession.isAdmin() {
			render.Render(w, r, ErrForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), keyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// contextUser returns the user from the request context
func (s *server) contextUser(r *http.Request) *cs.User {
	return r.Context().Value(keyUser).(*cs.User)
}

// users returns all users in the database (admin only)
func (s *server) users(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	user := s.contextUserSession(r)

	if !user.isAdmin() {
		log.Error().Str("User", user.Username).Msg("only admins can view all games")
		render.Render(w, r, ErrForbidden)
		return
	}

	users, err := db.GetUsers()
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("get users from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// don't ever send a password over the wire
	for i := range users {
		users[i].Password = ""
	}

	rest.RenderJSON(w, users)
}

func (s *server) user(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	user.Password = ""
	rest.RenderJSON(w, user)
}

// updateUserSettings will update a user editable settings.
// currently just supports updating the webhook, but could have other defaults in the future.
func (s *server) updateUserSettings(w http.ResponseWriter, r *http.Request) {
	readWriteClient := s.contextDb(r)
	user := s.contextUser(r)

	request := userSettingsRequest{}
	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if request.UserSettings == nil {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("invalid request body")))
		return
	}

	dbUser, err := readWriteClient.GetUser(user.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if dbUser == nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("user not found")
		render.Render(w, r, ErrNotFound)
		return
	}

	dbUser.UserSettings = *request.UserSettings

	if err := readWriteClient.UpdateUserSettings(dbUser); err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("failed to update webhook")
		render.Render(w, r, ErrNotFound)
		return
	}

	log.Info().
		Int64("UserID", user.ID).
		Msgf("updated user settings %v", request)

	rest.RenderJSON(w, rest.JSON{})
}
