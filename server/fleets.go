package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/game"
)

type cargoTransferBind struct {
	MO             game.MapObject `json:"mo,omitempty"`
	TransferAmount game.Cargo     `json:"transferAmount,omitempty"`
}

// Transfer cargo to/from a player's fleet
func (s *server) TransferCargo(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fleet, err := s.ctx.DB.FindFleetById(id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find the player for this user
	player, err := s.ctx.DB.FindPlayerByGameIdLight(fleet.GameID, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verify the user actually owns this planet
	if fleet.PlayerNum != player.Num {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("%s does not own %s", player, fleet)})
		return
	}

	// figure out what type of object we have
	transfer := cargoTransferBind{}
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch transfer.MO.Type {
	case game.MapObjectTypePlanet:
		s.transferCargoFleetPlanet(c, fleet, transfer)
	case game.MapObjectTypeFleet:
		s.transferCargoFleetFleet(c, fleet, transfer)
	}

}

// transfer cargo from a fleet to/from a planet
func (s *server) transferCargoFleetPlanet(c *gin.Context, fleet *game.Fleet, transfer cargoTransferBind) {
	// find the planet planet by id so we can perform the transfer
	planet, err := s.ctx.DB.FindPlanetById(transfer.MO.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if planet == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No planet for id %d found.", transfer.MO.ID)})
		return
	}

	if err := fleet.TransferPlanetCargo(planet, transfer.TransferAmount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.ctx.DB.SavePlanet(planet)
	s.ctx.DB.SaveFleet(fleet)

	log.Info().
		Uint("GameID", fleet.GameID).
		Uint("Player", uint(fleet.PlayerNum)).
		Str("Fleet", fleet.Name).
		Str("Planet", planet.Name).
		Str("TransferAmount", fmt.Sprintf("%v", transfer.TransferAmount)).
		Msgf("%s transfered %v to/from Planet %s", fleet.Name, transfer.TransferAmount, planet.Name)

	// success, return the updated fleet
	c.JSON(http.StatusOK, fleet)
}

// transfer cargo from a fleet to/from a fleet
func (s *server) transferCargoFleetFleet(c *gin.Context, fleet *game.Fleet, transfer cargoTransferBind) {
	// find the dest dest by id so we can perform the transfer
	dest, err := s.ctx.DB.FindFleetById(transfer.MO.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if dest == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No fleet for id %d found.", transfer.MO.ID)})
		return
	}

	if err := fleet.TransferFleetCargo(dest, transfer.TransferAmount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.ctx.DB.SaveFleet(dest)
	s.ctx.DB.SaveFleet(fleet)

	log.Info().
		Uint("GameID", fleet.GameID).
		Uint("Player", uint(fleet.PlayerNum)).
		Str("Fleet", fleet.Name).
		Str("Planet", dest.Name).
		Str("TransferAmount", fmt.Sprintf("%v", transfer.TransferAmount)).
		Msgf("%s transfered %v to/from Planet %s", fleet.Name, transfer.TransferAmount, dest.Name)

	// success, return the updated fleet
	c.JSON(http.StatusOK, fleet)

}
