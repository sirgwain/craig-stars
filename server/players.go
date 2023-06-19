package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirgwain/craig-stars/game"
)

func (s *server) PlayerGames(c *gin.Context) {
	user := s.GetSessionUser(c)

	games := s.ctx.DB.GetGamesByUser(user.ID)
	c.JSON(http.StatusOK, games)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if game == nil || player == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"game": game, "player": player})

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

// Allow a user to update a planet's orders
func (s *server) UpdatePlanetOrders(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find the player for this user
	player, err := s.ctx.DB.FindPlayerByGameIdLight(id.ID, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	planet := game.Planet{}
	if err := c.ShouldBindJSON(&planet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find the existing planet by id
	existing, err := s.ctx.DB.FindPlanetById(planet.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verify the user actually owns this planet
	if *existing.PlayerNum != player.Num {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("%s does not own %s", player, existing)})
		return
	}

	// copy user modifiable things to the existing planet
	existing.ContributesOnlyLeftoverToResearch = planet.ContributesOnlyLeftoverToResearch
	existing.ProductionQueue = planet.ProductionQueue
	s.ctx.DB.SavePlanet(existing)

	c.JSON(http.StatusOK, existing)
}
