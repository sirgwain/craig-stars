package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type cargoTransferBind struct {
	MO             cs.MapObject `json:"mo,omitempty"`
	TransferAmount cs.Cargo     `json:"transferAmount,omitempty"`
}

// Allow a user to update a fleet's orders
func (s *server) UpdateFleetOrders(c *gin.Context) {
	user := s.GetSessionUser(c)

	var fleetID idBind
	if err := c.ShouldBindUri(&fleetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fleet := cs.Fleet{}
	if err := c.ShouldBindJSON(&fleet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := s.playerUpdater.updateFleetOrders(user.ID, fleetID.ID, fleet.FleetOrders)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)

}

// Transfer cargo from a player's fleet to/from a fleet or planet the player controls
func (s *server) transferCargo(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// figure out what type of object we have
	transfer := cargoTransferBind{}
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := s.playerUpdater.transferCargo(user.ID, id.ID, transfer.MO.ID, transfer.MO.Type, transfer.TransferAmount)
	if err != nil {
		if errors.Is(err, errNotFound) {
			log.Error().Err(err).Msg("planet not found")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Error().Err(err).Msg("transfer cargo")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updated)

}
