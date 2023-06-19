package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type sessionUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// get the user from the context
func (s *server) contextUser(r *http.Request) sessionUser {
	return r.Context().Value(keyUser).(sessionUser)
}

func (s *server) mustGetUser(w http.ResponseWriter, r *http.Request) sessionUser {
	userInfo, err := token.GetUserInfo(r)
	if err != nil {
		panic("failed to load user")
	}

	userID, err := strconv.ParseInt(userInfo.StrAttr(databaseIDAttr), 10, 64)
	if err != nil {
		panic("failed to load user")
	}

	return sessionUser{
		ID:       userID,
		Username: userInfo.Name,
		Role:     userInfo.Role,
	}
}

func me(w http.ResponseWriter, r *http.Request) {

	userInfo, err := token.GetUserInfo(r)
	if err != nil {
		log.Printf("failed to get user info, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, err := strconv.ParseInt(userInfo.StrAttr(databaseIDAttr), 10, 64)
	if err != nil {
		log.Printf("failed to get user info, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := sessionUser{
		ID:       userID,
		Username: userInfo.Name,
		Role:     userInfo.Role,
	}

	rest.RenderJSON(w, res)
}

func (s *server) userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.mustGetUser(w, r)

		ctx := context.WithValue(r.Context(), keyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// create a new user from a token
func (s *server) createNewUser(tokenUser *token.User) (*cs.User, error) {
	user, err := cs.NewUser(tokenUser.Name, "", "", cs.RoleUser)
	if err != nil {
		log.Error().Err(err).Str("Username", user.Username).Msg("failed to create new user")
		return nil, err
	}
	err = s.db.CreateUser(user)
	if err != nil {
		log.Error().Err(err).Str("Username", user.Username).Msg("failed to create new user")
		return nil, err
	}
	log.Info().Str("Username", user.Username).Int64("ID", user.ID).Msg("created new user from token")

	// create a new test race
	race := cs.Humanoids()
	race.UserID = user.ID
	err = s.db.CreateRace(&race)
	log.Info().Str("Username", user.Username).Int64("ID", user.ID).Msg("created new race for user")
	if err != nil {
		return nil, err
	}
	return user, nil
}
