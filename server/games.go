package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

type hostGameRequest struct {
	*cs.GameSettings
}

func (req *hostGameRequest) Bind(r *http.Request) error {
	return nil
}

type updateGameRequest struct {
	*cs.GameSettings
}

func (req *updateGameRequest) Bind(r *http.Request) error {
	return nil
}

type joinGameRequest struct {
	Race cs.Race `json:"race"`
	Name string  `json:"name,omitempty"`
}

func (req *joinGameRequest) Bind(r *http.Request) error {
	return nil
}

type playerNumRequest struct {
	PlayerNum int `json:"playerNum"`
}

func (req *playerNumRequest) Bind(r *http.Request) error {
	return nil
}

type playerRequest struct {
	*cs.Player
}

func (req *playerRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id} calls
func (s *server) gameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.contextUser(r)
		db := s.contextDb(r)
		// load the game by id from the database
		id, err := s.int64URLParam(r, "id")
		if id == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		game, err := db.GetGame(*id)
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

func (s *server) games(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	db := s.contextDb(r)

	games, err := db.GetGamesForUser(user.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("get games from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, games)
}

func (s *server) hostedGames(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	db := s.contextDb(r)

	games, err := db.GetGamesForHost(user.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("get games from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, games)
}

func (s *server) openGames(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)

	games, err := db.GetOpenGames()
	if err != nil {
		log.Error().Err(err).Msg("get games from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, games)
}

func (s *server) openGamesByHash(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	// load open games by hash from the database
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("invalid invite hash in url")))
		return
	}

	games, err := db.GetOpenGamesByHash(hash)
	if err != nil {
		log.Error().Err(err).Str("Hash", hash).Msg("get open games by hash from database")
		rest.RenderJSON(w, games)
	}

	if len(games) == 0 {
		render.Render(w, r, ErrNotFound)
		return
	}

	// return games with this invite link
	rest.RenderJSON(w, games)
}

func (s *server) game(w http.ResponseWriter, r *http.Request) {
	game := s.contextGame(r)
	rest.RenderJSON(w, game)
}

func (s *server) getGuestUser(w http.ResponseWriter, r *http.Request) {
	game := s.contextGame(r)
	user := s.contextUser(r)
	db := s.contextDb(r)
	if user.ID != game.HostID {
		log.Error().Int64("GameID", game.ID).Str("User", user.Username).Msg("only host can load guests")
		render.Render(w, r, ErrForbidden)
		return
	}

	// load the game by id from the database
	num, err := s.intURLParam(r, "num")
	if num == nil || err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	guest, err := db.GetGuestUserForGame(game.ID, *num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", *num).Str("User", user.Username).Msgf("get guest for game")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, guest)
}

// Host a new game
func (s *server) createGame(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)

	settings := hostGameRequest{}
	if err := render.Bind(r, &settings); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// make sure guests don't create multiplayer games
	if user.isGuest() && !settings.IsSinglePlayer() {
		log.Error().Str("User", user.Username).Msg("cannot host multiplayer games")
		render.Render(w, r, ErrForbidden)
		return
	}

	gr := s.newGameRunner()
	game, err := gr.HostGame(user.ID, settings.GameSettings)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msgf("host game %v", settings.GameSettings)
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, game)
}

func (s *server) updateGame(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	game := s.contextGame(r)
	db := s.contextDb(r)

	update := updateGameRequest{}
	if err := render.Bind(r, &update); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if game.State != cs.GameStateSetup {
		log.Error().Int64("ID", game.ID).Msg("update game after setup")
		render.Render(w, r, ErrBadRequest(fmt.Errorf("game cannot be updated after setup")))
		return
	}

	if game.HostID != user.ID {
		log.Error().Int64("ID", game.ID).Int64("UserID", user.ID).Msgf("user %s is not the host", user.Username)
		render.Render(w, r, ErrForbidden)
		return
	}

	game.Name = update.Name
	game.Public = update.Public
	game.Size = update.Size
	game.Density = update.Density
	game.PlayerPositions = update.PlayerPositions
	game.RandomEvents = update.RandomEvents
	game.ComputerPlayersFormAlliances = update.ComputerPlayersFormAlliances
	game.PublicPlayerScores = update.PublicPlayerScores
	game.StartMode = update.StartMode
	game.QuickStartTurns = update.QuickStartTurns
	game.VictoryConditions = update.VictoryConditions

	if err := db.UpdateGame(&game.Game); err != nil {
		log.Error().Err(err).Int64("ID", game.ID).Msg("update game in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, game)
}

// Join an open game
func (s *server) joinGame(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	game := s.contextGame(r)
	gr := s.newGameRunner()

	join := joinGameRequest{}
	if err := render.Bind(r, &join); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if game.State != cs.GameStateSetup {
		err := fmt.Errorf("cannot join game after setup")
		log.Error().Err(err).Msg("join game")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if join.Name == "" {
		err := fmt.Errorf("name cannot be empty")
		log.Error().Err(err).Msg("join game")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// try and join this game
	if err := gr.JoinGame(game.ID, user.ID, join.Name, join.Race); err != nil {
		log.Error().Err(err).Msg("join game")
		render.Render(w, r, ErrBadRequest(err))
		return
	}
}

// Join an open game
func (s *server) leaveGame(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	game := s.contextGame(r)
	gr := s.newGameRunner()

	if game.State != cs.GameStateSetup {
		err := fmt.Errorf("cannot leave game after setup")
		log.Error().Err(err).Msg("leave game")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// try and join this game
	if err := gr.LeaveGame(game.ID, user.ID); err != nil {
		log.Error().Err(err).Msg("leave game")
		render.Render(w, r, ErrBadRequest(err))
		return
	}
}

// Join an open game
func (s *server) kickPlayer(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	user := s.contextUser(r)
	game := s.contextGame(r)
	gr := s.newGameRunner()

	if game.State != cs.GameStateSetup {
		err := fmt.Errorf("cannot leave game after setup")
		log.Error().Err(err).Msg("leave game")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if user.ID != game.HostID {
		log.Error().Int64("GameID", game.ID).Str("User", user.Username).Msg("user is not host")
		render.Render(w, r, ErrForbidden)
		return
	}

	kickPlayer := playerNumRequest{}
	if err := render.Bind(r, &kickPlayer); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// try and join this game
	if err := gr.KickPlayer(game.ID, kickPlayer.PlayerNum); err != nil {
		log.Error().Err(err).Msg("leave game")
		render.Render(w, r, ErrBadRequest(err))
	}

	// reload the game for the response
	game, err := db.GetGame(game.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, game)
}

func (s *server) addOpenPlayerSlot(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	user := s.contextUser(r)
	game := s.contextGame(r)
	gr := s.newGameRunner()

	if game.State != cs.GameStateSetup {
		err := fmt.Errorf("cannot leave game after setup")
		log.Error().Err(err).Msg("leave game")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if user.ID != game.HostID {
		log.Error().Int64("GameID", game.ID).Str("User", user.Username).Msg("user is not host")
		render.Render(w, r, ErrForbidden)
		return
	}

	if user.isGuest() {
		log.Error().Str("User", user.Username).Msg("cannot add open slots to games")
		render.Render(w, r, ErrForbidden)
		return
	}

	// add a new player slot to this game
	if _, err := gr.AddOpenPlayerSlot(game); err != nil {
		log.Error().Err(err).Msg("add player slot")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// reload the game for the response
	game, err := db.GetGame(game.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, game)
}

func (s *server) addGuestPlayer(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	user := s.contextUser(r)
	game := s.contextGame(r)
	gr := s.newGameRunner()

	if game.State != cs.GameStateSetup {
		err := fmt.Errorf("cannot leave game after setup")
		log.Error().Err(err).Msg("leave game")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if user.ID != game.HostID {
		log.Error().Int64("GameID", game.ID).Str("User", user.Username).Msg("user is not host")
		render.Render(w, r, ErrForbidden)
		return
	}

	if user.isGuest() {
		log.Error().Str("User", user.Username).Msg("cannot add guests to games")
		render.Render(w, r, ErrForbidden)
		return
	}

	// add a new player slot to this game
	if _, err := gr.AddGuestPlayer(game); err != nil {
		log.Error().Err(err).Msg("add player slot")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// reload the game for the response
	game, err := db.GetGame(game.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, game)
}

func (s *server) addAIPlayer(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	user := s.contextUser(r)
	game := s.contextGame(r)
	gr := s.newGameRunner()

	if game.State != cs.GameStateSetup {
		err := fmt.Errorf("cannot leave game after setup")
		log.Error().Err(err).Msg("leave game")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if user.ID != game.HostID {
		log.Error().Int64("GameID", game.ID).Str("User", user.Username).Msg("user is not host")
		render.Render(w, r, ErrForbidden)
		return
	}

	// add a new player slot to this game
	if _, err := gr.AddAIPlayer(game); err != nil {
		log.Error().Err(err).Msg("add player slot")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// reload the game for the response
	game, err := db.GetGame(game.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, game)
}

func (s *server) deletePlayerSlot(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	user := s.contextUser(r)
	game := s.contextGame(r)
	gr := s.newGameRunner()

	if game.State != cs.GameStateSetup {
		err := fmt.Errorf("cannot leave game after setup")
		log.Error().Err(err).Msg("leave game")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if user.ID != game.HostID {
		log.Error().Int64("GameID", game.ID).Str("User", user.Username).Msg("user is not host")
		render.Render(w, r, ErrForbidden)
		return
	}

	kickPlayer := playerNumRequest{}
	if err := render.Bind(r, &kickPlayer); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// add a new player slot to this game
	if err := gr.DeletePlayerSlot(game.ID, kickPlayer.PlayerNum); err != nil {
		log.Error().Err(err).Msg("delete player slot")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// reload the game for the response
	game, err := db.GetGame(game.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, game)
}

func (s *server) updatePlayerSlot(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	user := s.contextUser(r)
	game := s.contextGame(r)

	if game.State != cs.GameStateSetup {
		err := fmt.Errorf("cannot leave game after setup")
		log.Error().Err(err).Msg("leave game")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	player := playerRequest{}
	if err := render.Bind(r, &player); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// check the race
	if player.Race.ComputeRacePoints(game.Rules.RaceStartingPoints) < 0 {
		err := fmt.Errorf("race not valid, race points too high")
		log.Error().Err(err).Msg("update player")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	existing, err := db.GetPlayer(player.ID)
	if err != nil {
		log.Error().Int64("GameID", game.ID).Int64("PlayerID", player.ID).Msg("load player to update")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if existing == nil {
		log.Error().Int64("GameID", game.ID).Int64("PlayerID", player.ID).Msg("player not found")
		render.Render(w, r, ErrNotFound)
		return
	}

	if user.ID != game.HostID && existing.UserID != user.ID {
		log.Error().Int64("GameID", game.ID).Str("User", user.Username).Msg("user is not host or player owner")
		render.Render(w, r, ErrForbidden)
		return
	}

	// update all the fields allowed to be updating during game setup
	existing.UserID = player.UserID
	existing.Name = player.Name
	existing.Ready = player.Ready
	existing.AIControlled = player.AIControlled
	existing.Color = player.Color
	existing.DefaultHullSet = player.DefaultHullSet
	existing.Race = player.Race

	if err := db.UpdatePlayer(existing); err != nil {
		log.Error().Int64("GameID", game.ID).Int64("PlayerID", player.ID).Msg("updating player in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// reload the game for the response
	game, err = db.GetGame(game.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, game)
}

// Generate a universe for a host
func (s *server) startGame(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	game := s.contextGame(r)
	gr := s.newGameRunner()

	// validate
	if user.ID != game.HostID {
		render.Render(w, r, ErrForbidden)
		return
	}

	if user.isGuest() && !game.IsSinglePlayer() {
		log.Error().Str("User", user.Username).Msg("cannot start multiplayer game")
		render.Render(w, r, ErrForbidden)
		return
	}

	if err := gr.StartGame(&game.Game); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Msg("generating universe")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// send the full game to the host
	s.sendNewTurnNotification(r, game.ID)
	s.renderFullPlayerGame(w, r, game.ID, user.ID)
}

func (s *server) generateTurn(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	user := s.contextUser(r)
	game := s.contextGame(r)

	// validate
	if user.ID != game.HostID {
		render.Render(w, r, ErrForbidden)
		return
	}

	// only allow one GenerateTurn to run at a time for a game
	// TODO: handle this differently if you ever scale out beyond one instance. :)
	result, err, _ := s.sf.Do(strconv.FormatInt(game.ID, 10), func() (interface{}, error) {
		gr := s.newGameRunner()
		result, err := gr.GenerateTurn(game.ID)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		log.Error().Err(err).Msg("generate turn")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	player, err := db.GetPlayerForGame(game.ID, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("loading player after turn generation")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// return the game status
	game, err = db.GetGame(player.GameID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Msg("load game")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// if the host isn't a player, just return the game
	if player == nil {
		rest.RenderJSON(w, rest.JSON{"game": game})
		return
	}

	// return the new game
	if result == TurnGenerated {
		s.sendNewTurnNotification(r, game.ID)
		s.renderFullPlayerGame(w, r, player.GameID, player.UserID)
		return
	}

	rest.RenderJSON(w, rest.JSON{"game": game})
}

func (s *server) computeSpecs(w http.ResponseWriter, r *http.Request) {
	readWriteClient := s.contextDb(r)
	user := s.contextUser(r)
	game := s.contextGame(r)

	// validate
	if user.ID != game.HostID {
		render.Render(w, r, ErrForbidden)
		return
	}

	fg, err := readWriteClient.GetFullGame(game.ID)
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
		return c.UpdateFullGame(fg)
	}); err != nil {
		log.Error().Err(err).Msg("update game in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}
	rest.RenderJSON(w, rest.JSON{"game": fg.Game})
}

func (s *server) deleteGame(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	game := s.contextGame(r)

	if game.HostID != user.ID {
		log.Error().Int64("ID", game.ID).Int64("UserID", user.ID).Msg("only host can delete game")
		render.Render(w, r, ErrBadRequest(fmt.Errorf("only host can delete game")))
		return
	}

	if err := s.db.WrapInTransaction(func(c db.Client) error {
		if err := c.DeleteGame(game.ID); err != nil {
			log.Error().Err(err).Int64("ID", game.ID).Msg("delete game from database")
			return err
		}

		// delete any guest users
		if err := c.DeleteGameUsers(game.ID); err != nil {
			log.Error().Err(err).Int64("ID", game.ID).Msg("delete game guest users from database")
			return err
		}

		return nil

	}); err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}
}
