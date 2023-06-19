package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type designNumBind struct {
	DesignNum int `uri:"designnum"`
}

// CRUD for ship designs
func (s *server) createShipDesign(c *gin.Context) {
	user := s.GetSessionUser(c)

	var gameID idBind
	if err := c.ShouldBindUri(&gameID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	design := cs.ShipDesign{}
	if err := c.ShouldBindJSON(&design); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := s.playerUpdater.createShipDesign(gameID.ID, user.ID, &design)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, created)
}

func (s *server) getShipDesign(c *gin.Context) {
	user := s.GetSessionUser(c)

	var gameID idBind
	if err := c.ShouldBindUri(&gameID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var designNum designNumBind
	if err := c.ShouldBindUri(&designNum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := s.db.GetPlayerForGame(gameID.ID, user.ID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", gameID.ID).Int64("UserID", user.ID).Msg("load player from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to load player from database"})
		return
	}

	design, err := s.db.GetShipDesignByNum(gameID.ID, player.Num, designNum.DesignNum)
	if err != nil {
		log.Error().Err(err).Int("DesignNum", designNum.DesignNum).Int64("GameID", gameID.ID).Int64("UserID", user.ID).Msg("load design from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to load design from database"})
		return
	}

	if design == nil {
		log.Error().Err(err).Int("DesignNum", designNum.DesignNum).Int64("GameID", gameID.ID).Int64("UserID", user.ID).Msg("not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "design not found"})
		return
	}

	c.JSON(http.StatusOK, design)
}

func (s *server) getShipDesigns(c *gin.Context) {
	user := s.GetSessionUser(c)

	var gameID idBind
	if err := c.ShouldBindUri(&gameID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := s.db.GetPlayerWithDesignsForGame(gameID.ID, user.ID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", gameID.ID).Int64("UserID", user.ID).Msg("load player from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to load player from database"})
		return
	}

	c.JSON(http.StatusOK, player.Designs)

}

// update a shipdesign
func (s *server) updateShipDesign(c *gin.Context) {
	user := s.GetSessionUser(c)

	var gameID idBind
	if err := c.ShouldBindUri(&gameID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	design := cs.ShipDesign{}
	if err := c.ShouldBindJSON(&design); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := s.playerUpdater.updateShipDesign(gameID.ID, user.ID, &design)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)

}

func (s *server) deleteShipDesign(c *gin.Context) {
	user := s.GetSessionUser(c)

	var gameID idBind
	if err := c.ShouldBindUri(&gameID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var designNum designNumBind
	if err := c.ShouldBindUri(&designNum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// delete the ship design
	fleets, starbases, err := s.playerUpdater.deleteShipDesign(gameID.ID, user.ID, designNum.DesignNum)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fleets": fleets, "starbases": starbases})
}

func (s *server) computeShipDesignSpec(c *gin.Context) {
	user := s.GetSessionUser(c)

	var gameID idBind
	if err := c.ShouldBindUri(&gameID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	design := cs.ShipDesign{}
	if err := c.ShouldBindJSON(&design); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := s.db.GetPlayerForGame(gameID.ID, user.ID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Msg("loading player from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rules, err := s.db.GetRulesForGame(player.GameID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Msg("loading game rules from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	design.Spec = cs.ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, &design)

	c.JSON(http.StatusOK, design.Spec)
}
