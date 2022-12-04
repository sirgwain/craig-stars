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

// Allow a user to update a fleet's orders
func (s *server) UpdateFleetOrders(c *gin.Context) {
	user := s.GetSessionUser(c)

	var fleetID idBind
	if err := c.ShouldBindUri(&fleetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find the existing fleet by id
	existing, err := s.db.GetFleet(fleetID.ID)
	if err != nil {
		log.Error().Err(err).Msg("get fleet from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get fleet from database"})
		return
	}

	// find the player for this user
	player, err := s.db.GetLightPlayerForGame(existing.GameID, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("load player from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to load player from database"})
		return
	}

	// verify the user actually owns this fleet
	if existing.PlayerNum != player.Num {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Errorf("%s does not own %s", player, existing)})
		return
	}

	orders := game.FleetOrders{}
	if err := c.ShouldBindJSON(&orders); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := game.NewClient()
	client.UpdateFleetOrders(player, existing, orders)
	if err := s.db.UpdateFleet(existing); err != nil {
		log.Error().Err(err).Msg("update fleet in database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update fleet in database"})
		return
	}

	c.JSON(http.StatusOK, existing)
}

// Transfer cargo to/from a player's fleet
func (s *server) TransferCargo(c *gin.Context) {
	user := s.GetSessionUser(c)

	var id idBind
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fleet, err := s.db.GetFleet(id.ID)
	if err != nil {
		log.Error().Err(err).Msg("get fleet from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get fleet from database"})
		return
	}

	// find the player for this user
	player, err := s.db.GetLightPlayerForGame(fleet.GameID, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("load player from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get player from database"})
		return
	}

	// verify the user actually owns this planet
	if fleet.PlayerNum != player.Num {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Errorf("%s does not own %s", player, fleet)})
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
	planet, err := s.db.GetPlanet(transfer.MO.ID)
	if err != nil {
		log.Error().Err(err).Msg("get planet from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get planet from database"})
		return
	}

	if planet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No planet for id %d found.", transfer.MO.ID)})
		return
	}

	if err := fleet.TransferPlanetCargo(planet, transfer.TransferAmount); err != nil {
		log.Error().Err(err).Msg("transfer cargo")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to transfer cargo"})
		return
	}

	if err := s.db.UpdatePlanet(planet); err != nil {
		log.Error().Err(err).Int64("ID", planet.ID).Msg("update planet in database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to save planet to database"})
		return
	}
	
	if err := s.db.UpdateFleet(fleet); err != nil {
		log.Error().Err(err).Msg("update fleet in database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update fleet in database"})
		return
	}


	log.Info().
		Int64("GameID", fleet.GameID).
		Int("Player", fleet.PlayerNum).
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
	dest, err := s.db.GetFleet(transfer.MO.ID)
	if err != nil {
		log.Error().Err(err).Msg("get dest fleet from database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get dest fleet from database"})
		return
	}

	if dest == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No planet for id %d found.", transfer.MO.ID)})
		return
	}

	if err := fleet.TransferFleetCargo(dest, transfer.TransferAmount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.db.UpdateFleet(dest); err != nil {
		log.Error().Err(err).Msg("update fleet in database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update dest fleet in database"})
		return
	}

	if err := s.db.UpdateFleet(fleet); err != nil {
		log.Error().Err(err).Msg("update fleet in database")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update fleet in database"})
		return
	}

	log.Info().
		Int64("GameID", fleet.GameID).
		Int("Player", fleet.PlayerNum).
		Str("Fleet", fleet.Name).
		Str("Planet", dest.Name).
		Str("TransferAmount", fmt.Sprintf("%v", transfer.TransferAmount)).
		Msgf("%s transfered %v to/from Planet %s", fleet.Name, transfer.TransferAmount, dest.Name)

	// success, return the updated fleet
	c.JSON(http.StatusOK, fleet)

}
