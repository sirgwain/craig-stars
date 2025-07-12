//go:build !wasi && !wasm

package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

type userSession struct {
	ID            int64  `json:"id"`
	Role          string `json:"role"`
	Username      string `json:"username,omitempty"`
	DiscordID     string `json:"discordId,omitempty"`
	DiscordAvatar string `json:"discordAvatar,omitempty"`
}

func (u *userSession) isGuest() bool {
	return u.Role == string(cs.RoleGuest)
}

func (u *userSession) isAdmin() bool {
	return u.Role == string(cs.RoleAdmin)
}

const (
	attrDatabaseID    string = "database_id"
	attrDiscordID     string = "discord_id"
	attrDiscordAvatar string = "discord_avatar"
)

type tokenUser struct {
	*token.User
}

func newTokenUser(u *token.User) tokenUser {
	return tokenUser{User: u}
}

func (u *tokenUser) setDatabaseID(val int64) {
	u.SetStrAttr(attrDatabaseID, strconv.FormatInt(val, 10))
}

func (u *tokenUser) databaseID() int64 {
	if val := u.StrAttr(attrDatabaseID); val != "" {
		i, _ := strconv.ParseInt(val, 10, 64)
		return i
	}

	return 0
}

func (u *tokenUser) setDiscordID(val string) {
	u.SetStrAttr(attrDiscordID, val)
}

func (u *tokenUser) discordID() string {
	return u.StrAttr(attrDiscordID)
}

func (u *tokenUser) setDiscordAvatar(val string) {
	u.SetStrAttr(attrDiscordAvatar, val)
}

func (u *tokenUser) discordAvatar() string {
	return u.StrAttr(attrDiscordAvatar)
}

// get the user from the context
func (s *server) contextUserSession(r *http.Request) userSession {
	return r.Context().Value(keyUserSession).(userSession)
}

func (s *server) userSessionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.mustGetUserSession(w, r)

		ctx := context.WithValue(r.Context(), keyUserSession, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) mustGetUserSession(_ http.ResponseWriter, r *http.Request) userSession {
	userInfo, err := token.GetUserInfo(r)
	if err != nil {
		panic("failed to load user")
	}

	tokenUser := newTokenUser(&userInfo)

	userID := tokenUser.databaseID()
	discordID := tokenUser.discordID()
	discordAvatar := tokenUser.discordAvatar()

	return userSession{
		ID:            userID,
		Username:      userInfo.Name,
		Role:          userInfo.Role,
		DiscordID:     discordID,
		DiscordAvatar: discordAvatar,
	}
}

func me(w http.ResponseWriter, r *http.Request) {
	userInfo, err := token.GetUserInfo(r)
	if err != nil {
		log.Printf("failed to get user info, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, err := strconv.ParseInt(userInfo.StrAttr(attrDatabaseID), 10, 64)
	if err != nil {
		log.Printf("failed to get user info, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var discordID string
	var discordAvatar string

	if val, ok := userInfo.Attributes["discord_id"]; ok {
		discordID = val.(string)
	}
	if val, ok := userInfo.Attributes["discord_avatar"]; ok {
		discordAvatar = val.(string)
	}

	res := userSession{
		ID:            userID,
		Username:      userInfo.Name,
		Role:          userInfo.Role,
		DiscordID:     discordID,
		DiscordAvatar: discordAvatar,
	}

	rest.RenderJSON(w, res)
}

// create a new user from a token
func (s *server) createNewDiscordUser(ctx context.Context, tokenUser tokenUser) (*cs.User, error) {

	user, err := cs.NewDiscordUser(tokenUser.Name, tokenUser.discordID(), tokenUser.discordAvatar())
	if err != nil {
		log.Error().Err(err).Str("Username", user.Username).Msg("failed to create new user")
		return nil, err
	}

	if err := s.db.WrapInTransaction(func(c db.Client) error {
		user, err := c.CreateUser(ctx, user)
		if err != nil {
			log.Error().Err(err).Str("Username", user.Username).Msg("failed to create new user")
			return err
		}
		log.Info().Str("Username", user.Username).Int64("ID", user.ID).Msg("created new user from token")

		// create a new test race
		race := cs.Humanoids()
		race.UserID = user.ID
		_, err = c.CreateRace(ctx, &race)
		if err != nil {
			return err
		}
		log.Info().Str("Username", user.Username).Int64("ID", user.ID).Msg("created new race for user")
		return nil
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *server) updateUser(ctx context.Context, tokenUser tokenUser, user *cs.User) error {

	idStr := tokenUser.discordID()
	avatarStr := tokenUser.discordAvatar()
	user.DiscordID = &idStr
	user.DiscordAvatar = &avatarStr
	now := time.Now()
	user.LastLogin = &now

	readWriteClient := s.db.NewReadWriteClient()
	if err := readWriteClient.UpdateUser(ctx, user); err != nil {
		log.Error().Err(err).Str("Username", user.Username).Msg("failed to update user")
		return err
	}
	log.Info().Str("Username", user.Username).Int64("ID", user.ID).Str("DiscordID", *user.DiscordID).Str("DiscordAvatar", *user.DiscordAvatar).Msg("updated")

	return nil
}
