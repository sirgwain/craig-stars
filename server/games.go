package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type idBind struct {
	ID int64 `uri:"id"`
}

func (s *server) deleteGame(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate
	game, err := s.db.GetGame(id.ID)
	if err != nil {
		log.Error().Err(err).Int64("ID", id.ID).Msg("get game from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get game from database"})
		return
	}

	if game.HostID != user.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only the host can delete a game"})
		return
	}

	// delete it
	if err := s.db.DeleteGame(id.ID); err != nil {
		log.Error().Err(err).Int64("ID", id.ID).Msg("delete game from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to delete game from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
