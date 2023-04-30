package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

// context for /api/games/{id} calls that require a player
func (s *server) playerCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.contextUser(r)
		game := r.Context().Value(keyGame).(*cs.Game)

		player, err := s.db.GetLightPlayerForGame(game.ID, user.ID)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if player == nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), keyPlayer, player)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) player(w http.ResponseWriter, r *http.Request) {
	player := r.Context().Value(keyPlayer).(*cs.Player)
	rest.RenderJSON(w, player)
}

func (s *server) fullPlayer(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)
	game := r.Context().Value(keyGame).(*cs.Game)

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

// Submit a turn for the player
func (s *server) updatePlayerOrders(c *gin.Context) {
	user := s.GetSessionUser(c)

	var gameID idBind
	if err := c.ShouldBindUri(&gameID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orders := cs.PlayerOrders{}
	if err := c.ShouldBindJSON(&orders); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if orders.ResearchAmount < 0 || orders.ResearchAmount > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "research ammount must be between 0 and 100"})
		return
	}

	player, planets, err := s.playerUpdater.updatePlayerOrders(gameID.ID, user.ID, orders)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int64("GameID", gameID.ID).Int64("PlayerID", player.ID).Msg("update orders")
	c.JSON(http.StatusOK, gin.H{
		"player":  player,
		"planets": planets,
	})
}

// Submit a turn for the player
func (s *server) submitTurn(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.gameRunner.SubmitTurn(id.ID, user.ID); err != nil {
		log.Error().Err(err).Int64("GameID", id.ID).Int64("UserID", user.ID).Msg("submit turn")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to submit turn"})
		return
	}

	_, err := s.gameRunner.CheckAndGenerateTurn(id.ID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", id.ID).Msg("check and generate new turn")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, player, err := s.gameRunner.LoadPlayerGame(id.ID, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"game": game, "player": player})

}
