package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

type sessionUser struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Role          string `json:"role"`
	DiscordID     string `json:"discordId"`
	DiscordAvatar string `json:"discordAvatar"`
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

	var discordID string
	var discordAvatar string

	if val, ok := userInfo.Attributes["discord_id"]; ok {
		discordID = val.(string)
	}
	if val, ok := userInfo.Attributes["discord_avatar"]; ok {
		discordAvatar = val.(string)
	}

	return sessionUser{
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

	userID, err := strconv.ParseInt(userInfo.StrAttr(databaseIDAttr), 10, 64)
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

	res := sessionUser{
		ID:            userID,
		Username:      userInfo.Name,
		Role:          userInfo.Role,
		DiscordID:     discordID,
		DiscordAvatar: discordAvatar,
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

	discordID, foundDiscordID := tokenUser.Attributes["discord_id"]
	discordAvatar, foundDiscordAvatar := tokenUser.Attributes["discord_avatar"]
	if !foundDiscordID || !foundDiscordAvatar {
		return nil, fmt.Errorf("trying to create new user that isn't a discord user")
	}

	user, err := cs.NewDiscordUser(tokenUser.Name, discordID.(string), discordAvatar.(string))
	if err != nil {
		log.Error().Err(err).Str("Username", user.Username).Msg("failed to create new user")
		return nil, err
	}

	if err := s.db.WrapInTransaction(func(c db.Client) error {
		if err := c.CreateUser(user); err != nil {
			log.Error().Err(err).Str("Username", user.Username).Msg("failed to create new user")
			return err
		}
		log.Info().Str("Username", user.Username).Int64("ID", user.ID).Msg("created new user from token")

		// create a new test race
		race := cs.Humanoids()
		race.UserID = user.ID
		if err = c.CreateRace(&race); err != nil {
			return err
		}
		log.Info().Str("Username", user.Username).Int64("ID", user.ID).Msg("created new race for user")
		return nil
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *server) updateUser(tokenUser *token.User, user *cs.User) error {

	discordID, foundDiscordID := tokenUser.Attributes["discord_id"]
	discordAvatar, foundDiscordAvatar := tokenUser.Attributes["discord_avatar"]
	if !foundDiscordID || !foundDiscordAvatar {
		return fmt.Errorf("trying to create new user that isn't a discord user")
	}

	idStr := discordID.(string)
	avatarStr := discordAvatar.(string)
	user.DiscordID = &idStr
	user.DiscordAvatar = &avatarStr
	now := time.Now()
	user.LastLogin = &now

	readWriteClient := s.db.NewReadWriteClient()
	if err := readWriteClient.UpdateUser(user); err != nil {
		log.Error().Err(err).Str("Username", user.Username).Msg("failed to update user")
		return err
	}
	log.Info().Str("Username", user.Username).Int64("ID", user.ID).Str("DiscordID", *user.DiscordID).Str("DiscordAvatar", *user.DiscordAvatar).Msg("updated")

	return nil
}
