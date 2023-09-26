package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/db"
)

type convertGuestUserRequest struct {
	UserID int64 `json:"userID,omitempty"`
}

func (req *convertGuestUserRequest) Bind(r *http.Request) error {
	return nil
}

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

func (s *server) userGames(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)

	// load the games for users with id from the database
	id, err := s.int64URLParam(r, "id")
	if id == nil || err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	games, err := db.GetGamesForUser(*id)
	if err != nil {
		log.Error().Err(err).Msg("get games from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, games)
}

// convert a guest user into a full user
// this takes a guest user's games and races and assigns them to a new player
func (s *server) convertGuestUser(w http.ResponseWriter, r *http.Request) {
	readWriteClient := s.contextDb(r)

	request := convertGuestUserRequest{}
	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// load the games for users with id from the database
	id, err := s.int64URLParam(r, "id")
	if id == nil || err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	guestUser, err := readWriteClient.GetUser(*id)
	if err != nil {
		log.Error().Err(err).Int64("UserID", *id).Msg("load guest user")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if guestUser == nil {
		log.Error().Err(err).Int64("UserID", *id).Msg("guest user not found")
		render.Render(w, r, ErrNotFound)
		return
	}

	// make sure this is a guest user
	if !guestUser.IsGuest() {
		err := fmt.Errorf("%s is not a guest user", guestUser.Username)
		log.Error().Err(err).Int64("ID", guestUser.ID).Str("Username", guestUser.Username).Msg("user is not a guest")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	user, err := readWriteClient.GetUser(request.UserID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", request.UserID).Msg("load user")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if user == nil {
		log.Error().Err(err).Int64("UserID", request.UserID).Msg("user not found")
		render.Render(w, r, ErrNotFound)
		return
	}

	// make sure the user to convert into is NOT a guest user
	if user.IsGuest() {
		err := fmt.Errorf("%s is a guest user", user.Username)
		log.Error().Err(err).Int64("ID", user.ID).Str("Username", user.Username).Msg("user is a guest")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// load players
	players, err := readWriteClient.GetPlayersForUser(guestUser.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", request.UserID).Msg("load players")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}
	// load single player games
	playerGames, err := readWriteClient.GetGamesForUser(guestUser.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", request.UserID).Msg("load player games")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// make sure this user isn't playing in any of the same games as the guest
	for _, game := range playerGames {
		for _, player := range game.Players {
			if player.UserID == user.ID {
				err := fmt.Errorf("%s is already playing in game %d (%s)", user.Username, game.ID, game.Name)
				log.Error().Err(err).Msg("user cannot play in multiple games")
				render.Render(w, r, ErrBadRequest(err))
				return
			}
		}
	}

	// load races
	races, err := readWriteClient.GetRacesForUser(guestUser.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", request.UserID).Msg("load races")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// load single player games
	games, err := readWriteClient.GetGamesForHost(guestUser.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", request.UserID).Msg("load games")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if err := s.db.WrapInTransaction(func(c db.Client) error {

		// update each player UserID
		for _, player := range players {
			player.UserID = user.ID
			if err := c.UpdatePlayerUserId(&player); err != nil {
				return fmt.Errorf("update Player UserID %w", err)
			}
		}

		// update each race UserID
		for _, race := range races {
			race.UserID = user.ID
			if err := c.UpdateRace(&race); err != nil {
				return fmt.Errorf("update Race UserID %w", err)
			}
		}

		// update each game UserID
		for _, game := range games {
			game.HostID = user.ID
			if err := c.UpdateGameHost(game.ID, game.HostID); err != nil {
				return fmt.Errorf("update Game HostID %w", err)
			}
		}

		if err := c.DeleteUser(guestUser.ID); err != nil {
			return fmt.Errorf("delete guest user %w", err)
		}

		return nil
	}); err != nil {
		log.Error().Err(err).Msg("update user objects in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	log.Info().
		Int64("GuestUserID", guestUser.ID).
		Int64("UserID", user.ID).
		Msgf("moved guest %s games and races to %s, deleted %s", guestUser.Username, user.Username, guestUser.Username)
	rest.RenderJSON(w, rest.JSON{})
}
