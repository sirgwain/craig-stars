package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type idBind struct {
	ID uint `uri:"id"`
}

func (s *server) Games(c *gin.Context) {
	games, err := s.ctx.DB.GetGames()
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

	game, err := s.ctx.DB.FindGameById(id.ID)
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
