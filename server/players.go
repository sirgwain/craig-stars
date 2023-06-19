package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type HostGameBind struct {
	Settings cs.GameSettings `json:"settings"`
}

type JoinGameBind struct {
	RaceID int64 `json:"raceId"`
}

func (s *server) playerGames(c *gin.Context) {
	user := s.GetSessionUser(c)

	games, err := s.db.GetGamesForUser(user.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("get games from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get games from database"})
		return
	}

	c.JSON(http.StatusOK, games)
}

func (s *server) hostedGames(c *gin.Context) {
	user := s.GetSessionUser(c)

	games, err := s.db.GetGamesForHost(user.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("get games from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get games from database"})
		return
	}

	c.JSON(http.StatusOK, games)
}

func (s *server) openGames(c *gin.Context) {
	user := s.GetSessionUser(c)

	games, err := s.db.GetOpenGames(user.ID)
	if err != nil {
		log.Error().Err(err).Msg("get open games from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get open games from database"})
		return
	}

	c.JSON(http.StatusOK, games)
}

func (s *server) openGame(c *gin.Context) {

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := s.db.GetGame(id.ID)
	if err != nil {
		log.Error().Err(err).Int64("ID", id.ID).Msg("get game from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get game from database"})
		return
	}

	if game == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("game with id %d not found", id.ID)})
	}

	c.JSON(http.StatusOK, game)
}

func (s *server) playerGame(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, player, err := s.gameRunner.LoadPlayerGame(id.ID, user.ID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", id.ID).Int64("UserID", user.ID).Msg("load player and game from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to load game from database"})
		return
	}

	if game == nil || player == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"game": game, "player": player})

}

func (s *server) playerStatuses(c *gin.Context) {

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	players, err := s.db.GetPlayerStatusesForGame(id.ID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", id.ID).Msg("load players and game from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to load players from database"})
		return
	}

	if len(players) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"players": players})

}

// Host a new game
func (s *server) hostGame(c *gin.Context) {
	user := s.GetSessionUser(c)

	body := HostGameBind{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := s.gameRunner.HostGame(user.ID, &body.Settings)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msgf("host game %v", body.Settings)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to host game"})
		return
	}

	_ = game

	c.JSON(http.StatusOK, gin.H{})

}

// Join an open game
func (s *server) joinGame(c *gin.Context) {
	user := s.GetSessionUser(c)

	var game idBind
	if err := c.ShouldBindUri(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body := JoinGameBind{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// try and join this game
	if err := s.gameRunner.JoinGame(game.ID, user.ID, body.RaceID); err != nil {
		if errors.Is(err, errNotFound) {
			log.Error().Err(err).Msg("not found")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Error().Err(err).Msg("join game")
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to join game"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Generate a universe for a host
func (s *server) generateUniverse(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.gameRunner.GenerateUniverse(id.ID, user.ID); err != nil {
		log.Error().Err(err).Int64("GameID", id.ID).Int64("UserID", user.ID).Msg("submit turn")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to submit turn"})
		return
	}

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

	player, planets, err := s.playerUpdater.updatePlayerOrders(user.ID, gameID.ID, orders)

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
