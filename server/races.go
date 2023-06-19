package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

func (s *server) races(c *gin.Context) {
	user := s.GetSessionUser(c)

	races, err := s.db.GetRacesForUser(user.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("get races from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get races from database"})
		return
	}

	c.JSON(http.StatusOK, races)
}

func (s *server) race(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	race, err := s.db.GetRace(id.ID)
	if err != nil {
		log.Error().Err(err).Int64("ID", id.ID).Msg("get race from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get race from database"})
		return
	}

	if race.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("user %d does not own race %d", user.ID, id.ID)})
		return
	}

	c.JSON(http.StatusOK, race)
}

// Submit a turn for the player
func (s *server) createRace(c *gin.Context) {
	user := s.GetSessionUser(c)

	race := cs.Race{}
	if err := c.ShouldBindJSON(&race); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	race.UserID = user.ID
	if err := s.db.CreateRace(&race); err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("create race")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create race"})
		return
	}

	c.JSON(http.StatusOK, race)
}

func (s *server) updateRace(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	race := cs.Race{}
	if err := c.ShouldBindJSON(&race); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// load in the existing race from the database
	existingRace, err := s.db.GetRace(id.ID)

	if err != nil {
		log.Error().Err(err).Int64("ID", id.ID).Msg("get race from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get race from database"})
		return
	}

	// validate

	if existingRace == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("race with id %d not found", id.ID)})
		return
	}

	if id.ID != race.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id in body does not match id in url"})
		return
	}

	if user.ID != existingRace.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("user %d does not own race %d", user.ID, existingRace.ID)})
		return
	}

	if err := s.db.UpdateRace(&race); err != nil {
		log.Error().Err(err).Int64("ID", id.ID).Msg("update race in database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update race in database"})
		return
	}

	c.JSON(http.StatusOK, race)
}
