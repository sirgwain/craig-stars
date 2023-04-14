package server

import (
	"fmt"
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

	shipdesign := cs.ShipDesign{}
	if err := c.ShouldBindJSON(&shipdesign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// load in the existing shipdesign from the database
	existingShipDesign, err := s.db.GetShipDesign(shipdesign.ID)

	if err != nil {
		log.Error().Err(err).Int64("ID", shipdesign.ID).Msg("get shipdesign from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get shipdesign from database"})
		return
	}

	player, err := s.db.GetLightPlayerForGame(gameID.ID, user.ID)
	if err != nil || player == nil {
		log.Error().Err(err).Int64("GameID", gameID.ID).Int64("UserID", user.ID).Msg("load player from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to load player from database"})
		return
	}

	// validate
	if existingShipDesign == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("shipdesign with id %d not found", shipdesign.ID)})
		return
	}

	if player.Num != existingShipDesign.PlayerNum {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("user %d does not own shipdesign %d", user.ID, existingShipDesign.ID)})
		return
	}

	if err := s.db.UpdateShipDesign(&shipdesign); err != nil {
		log.Error().Err(err).Int64("ID", shipdesign.ID).Msg("update shipdesign in database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update shipdesign in database"})
		return
	}

	c.JSON(http.StatusOK, shipdesign)
}
