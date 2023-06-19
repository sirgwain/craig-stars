package server

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type playerOrdersRequest struct {
	*cs.PlayerOrders
}

func (req *playerOrdersRequest) Bind(r *http.Request) error {
	return nil
}

type playerPlansRequest struct {
	*cs.PlayerPlans
}

func (req *playerPlansRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id} calls that require a player
func (s *server) playerCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.contextUser(r)
		game := s.contextGame(r)

		player, err := s.db.GetLightPlayerForGame(game.ID, user.ID)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if player == nil {
			log.Error().Int64("GameID", game.ID).Int64("UserID", user.ID).Msg("player not found")
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), keyPlayer, player)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) contextPlayer(r *http.Request) *cs.Player {
	return r.Context().Value(keyPlayer).(*cs.Player)
}

func (s *server) player(w http.ResponseWriter, r *http.Request) {
	player := s.contextPlayer(r)
	rest.RenderJSON(w, player)
}

func (s *server) fullPlayer(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	game := s.contextGame(r)

	player, err := s.db.GetPlayerForGame(game.ID, user.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if player == nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	rest.RenderJSON(w, player)
}

// get mapObjects for a player
func (s *server) mapObjects(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)

	gameID, err := s.int64URLParam(r, "id")
	if gameID == nil || err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	mapObjects, err := s.db.GetPlayerMapObjects(*gameID, user.ID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", *gameID).Int64("UserID", user.ID).Msg("load player map objects database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if mapObjects == nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	rest.RenderJSON(w, mapObjects)
}

func (s *server) playerStatuses(w http.ResponseWriter, r *http.Request) {
	gameID, err := s.int64URLParam(r, "id")
	if gameID == nil || err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	players, err := s.db.GetPlayerStatusesForGame(*gameID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", *gameID).Msg("load players and game from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if len(players) == 0 {
		render.Render(w, r, ErrNotFound)
		return
	}

	rest.RenderJSON(w, rest.JSON{"players": players})
}

// submit a player turn and return the newly generated turn if there is one
func (s *server) submitTurn(w http.ResponseWriter, r *http.Request) {
	player := s.contextPlayer(r)

	// submit the turn
	player.SubmittedTurn = true
	if err := s.db.UpdatePlayerOrders(player); err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// TODO: this should probably be a goroutine or something
	_, err := s.gameRunner.CheckAndGenerateTurn(player.GameID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Msg("check and generate new turn")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	game, fullPlayer, err := s.gameRunner.LoadPlayerGame(player.GameID, player.UserID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Msg("load full game from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, rest.JSON{"game": game, "player": fullPlayer.Player, "mapObjects": fullPlayer.PlayerMapObjects})
}

// Submit a turn for the player
func (s *server) updatePlayerOrders(w http.ResponseWriter, r *http.Request) {
	game := s.contextGame(r)
	player := s.contextPlayer(r)

	orders := playerOrdersRequest{}
	if err := render.Bind(r, &orders); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if orders.ResearchAmount < 0 || orders.ResearchAmount > 100 {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("research ammount must be between 0 and 100")))
		return
	}

	planets, err := s.db.GetPlanetsForPlayer(player.GameID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Msg("loading player planets from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	orderer := cs.NewOrderer()
	orderer.UpdatePlayerOrders(player, planets, *orders.PlayerOrders, &game.Rules)

	// save the player to the database
	if err := s.db.UpdatePlayerOrders(player); err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	for _, planet := range planets {
		if planet.Dirty {
			// TODO: only update the planet spec? that's all that changes
			// TODO: do this all in one transaction?
			if err := s.db.UpdatePlanet(planet); err != nil {
				log.Error().Err(err).Int64("ID", player.ID).Msg("updating player planet in database")
				render.Render(w, r, ErrInternalServerError(err))
				return
			}
		}
	}

	log.Info().Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update orders")
	rest.RenderJSON(w, rest.JSON{"player": player, "planets": planets})
}

// Submit a turn for the player
func (s *server) updatePlayerPlans(w http.ResponseWriter, r *http.Request) {
	player := s.contextPlayer(r)

	plans := playerPlansRequest{}
	if err := render.Bind(r, &plans); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if len(plans.BattlePlans) == 0 {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("must have at least one battle plan")))
		return
	}

	if plans.BattlePlans[0].Num != 0 {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("must have a default battle plan")))
		return
	}

	// TODO: validate?
	// TODO: convert creates into a separate POST?
	// TODO: update fleets with deleted battle plans to use default battleplan
	nextNum := 0
	for i := range plans.BattlePlans {
		nextNum = int(math.Max(float64(plans.BattlePlans[i].Num+1), float64(nextNum)))
	}

	for i := range plans.BattlePlans {
		if plans.BattlePlans[i].Num == -1 {
			plans.BattlePlans[i].Num = nextNum
			nextNum++
		}
	}

	player.PlayerPlans = *plans.PlayerPlans

	// save the player to the database
	if err := s.db.UpdatePlayerOrders(player); err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	log.Info().Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update plans")
	rest.RenderJSON(w, player)
}
