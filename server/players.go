package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/game"
)

type HostGameBind struct {
	Settings game.GameSettings `json:"settings"`
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
	games, err := s.db.GetOpenGames()
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

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body := JoinGameBind{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := s.db.GetFullGame(id.ID)
	if err != nil {
		log.Error().Err(err).Int64("ID", id.ID).Msg("get game from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get game from database"})
		return
	}

	race, err := s.db.GetRace(body.RaceID)
	if err != nil {
		log.Error().Err(err).Int64("ID", body.RaceID).Msg("get race from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get race from database"})
		return
	}

	// validate
	if race == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("a race with id %d not found", body.RaceID)})
		return
	}

	if race.UserID != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("you do not own the race with id %d not found", body.RaceID)})
		return
	}

	if game == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("a game with id %d not found", id.ID)})
		return
	}

	if game.OpenPlayerSlots == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No slots left in this game"})
		return
	}

	// all good, add the player
	if err := s.gameRunner.AddPlayer(game.ID, user.ID, race); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int64("UserID", user.ID).Msgf("addPlayer, race: %v", race)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to join game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Submit a turn for the player
func (s *server) updatePlayerOrders(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orders := game.PlayerOrders{}
	if err := c.ShouldBindJSON(&orders); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, planets, err := s.playerUpdater.UpdatePlayerOrders(user.ID, id.ID, orders)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	status, err := s.gameRunner.CheckAndGenerateTurn(id.ID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", id.ID).Msg("check and generate new turn")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if status == TurnGenerated {
		game, player, err := s.gameRunner.LoadPlayerGame(id.ID, user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"game": game, "player": player})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}

}
