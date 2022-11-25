package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirgwain/craig-stars/game"
)

func (s *server) Races(c *gin.Context) {
	user := s.GetSessionUser(c)

	races, err := s.db.GetRacesForUser(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, races)
}

func (s *server) Race(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	race, err := s.db.GetRace(id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if race.UserID != user.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("user %d does not own race %d", user.ID, id.ID)})
		return
	}

	c.JSON(http.StatusOK, race)
}

// Submit a turn for the player
func (s *server) CreateRace(c *gin.Context) {
	user := s.GetSessionUser(c)

	race := game.Race{}
	if err := c.ShouldBindJSON(&race); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	race.UserID = user.ID
	err := s.db.CreateRace(&race)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, race)
}

func (s *server) UpdateRace(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	race := game.Race{}
	if err := c.ShouldBindJSON(&race); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// load in the existing race from the database
	existingRace, err := s.db.GetRace(id.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("user %d does not own race %d", user.ID, existingRace.ID)})
		return
	}

	if err := s.db.UpdateRace(&race); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, race)
}
