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

func (s *server) PlayerGames(c *gin.Context) {
	user := s.GetSessionUser(c)

	games, err := s.db.GetGamesForUser(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, games)
}

func (s *server) HostedGames(c *gin.Context) {
	user := s.GetSessionUser(c)

	games, err := s.db.GetGamesForHost(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, games)
}

func (s *server) OpenGames(c *gin.Context) {
	games, err := s.db.GetOpenGames()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, games)
}

func (s *server) OpenGame(c *gin.Context) {

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := s.db.GetGame(id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, game)
}

func (s *server) PlayerGame(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, player, err := s.gameRunner.LoadPlayerGame(id.ID, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to load player and game from database")
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
func (s *server) HostGame(c *gin.Context) {
	user := s.GetSessionUser(c)

	body := HostGameBind{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := s.gameRunner.HostGame(user.ID, &body.Settings)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = game

	c.JSON(http.StatusOK, gin.H{})

}

// Join an open game
func (s *server) JoinGame(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	race, err := s.db.GetRace(body.RaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// Submit a turn for the player
func (s *server) SubmitTurn(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.gameRunner.SubmitTurn(id.ID, user.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := s.gameRunner.CheckAndGenerateTurn(id.ID)
	if err != nil {
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
