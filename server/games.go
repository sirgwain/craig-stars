package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type idBind struct {
	ID int64 `uri:"id"`
}

func (s *server) Games(c *gin.Context) {
	games, err := s.db.GetGames()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, games)
}

func (s *server) GameById(c *gin.Context) {
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

	if game == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, game)
}

func (s *server) DeleteGame(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate
	game, err := s.db.GetGame(id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if game.HostID != user.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only the host can delete a game"})
		return
	}

	// delete it
	if err := s.db.DeleteGame(id.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
