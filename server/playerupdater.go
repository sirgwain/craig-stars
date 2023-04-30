package server

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

var errNotFound = errors.New("resource was not found")

type playerUpdater struct {
	orderer cs.Orderer
	db      DBClient
}

type PlayerUpdater interface {
	updateFleetOrders(userID int64, fleetID int64, orders cs.FleetOrders) (*cs.Fleet, error)
	updateMineFieldOrders(userID int64, mineFieldID int64, orders cs.MineFieldOrders) (*cs.MineField, error)
	transferCargo(userID int64, fleetID int64, destID int64, mapObjectType cs.MapObjectType, transferAmount cs.Cargo) (*cs.Fleet, error)
}

func newPlayerUpdater(db DBClient) PlayerUpdater {
	return &playerUpdater{orderer: cs.NewOrderer(), db: db}
}

// update the orders for a fleet, i.e. waypoints and battle plan
func (pu *playerUpdater) updateFleetOrders(userID int64, fleetID int64, orders cs.FleetOrders) (*cs.Fleet, error) {
	// find the fleet fleet by id
	fleet, err := pu.db.GetFleet(fleetID)
	if err != nil {
		log.Error().Err(err).Int64("ID", fleetID).Msg("load fleet from database")
		return nil, fmt.Errorf("failed to load fleet from database")
	}

	player, err := pu.db.GetLightPlayerForGame(fleet.GameID, userID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", fleet.GameID).Int64("UserID", userID).Msg("load player from database")
		return nil, fmt.Errorf("failed to load player from database")
	}

	// verify the user actually owns this fleet
	if fleet.PlayerNum != player.Num {
		return nil, fmt.Errorf("%s does not own fleet %s", player, fleet)
	}

	pu.orderer.UpdateFleetOrders(player, fleet, orders)

	if err := pu.db.UpdateFleet(fleet); err != nil {
		log.Error().Err(err).Int64("ID", fleet.ID).Msg("update fleet in database")
		return nil, fmt.Errorf("failed to save fleet to database")
	}

	return fleet, nil
}

// update the orders for a minefield, i.e. whether it detonates
func (pu *playerUpdater) updateMineFieldOrders(userID int64, mineFieldID int64, orders cs.MineFieldOrders) (*cs.MineField, error) {
	// find the mineField mineField by id
	mineField, err := pu.db.GetMineField(mineFieldID)
	if err != nil {
		log.Error().Err(err).Int64("ID", mineFieldID).Msg("load mineField from database")
		return nil, fmt.Errorf("failed to load mineField from database")
	}

	player, err := pu.db.GetLightPlayerForGame(mineField.GameID, userID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", mineField.GameID).Int64("UserID", userID).Msg("load player from database")
		return nil, fmt.Errorf("failed to load player from database")
	}

	// verify the user actually owns this mineField
	if mineField.PlayerNum != player.Num {
		return nil, fmt.Errorf("%s does not own minefield %s", player, mineField)
	}

	pu.orderer.UpdateMineFieldOrders(player, mineField, orders)

	if err := pu.db.UpdateMineField(mineField); err != nil {
		log.Error().Err(err).Int64("ID", mineField.ID).Msg("update mineField in database")
		return nil, fmt.Errorf("failed to save mineField to database")
	}

	return mineField, nil
}

func (pu *playerUpdater) transferCargo(userID int64, fleetID int64, destID int64, mapObjectType cs.MapObjectType, transferAmount cs.Cargo) (*cs.Fleet, error) {
	// find the fleet fleet by id
	fleet, err := pu.db.GetFleet(fleetID)
	if err != nil {
		log.Error().Err(err).Int64("ID", fleetID).Msg("load fleet from database")
		return nil, fmt.Errorf("failed to load fleet from database")
	}

	player, err := pu.db.GetLightPlayerForGame(fleet.GameID, userID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", fleet.GameID).Int64("UserID", userID).Msg("load player from database")
		return nil, fmt.Errorf("failed to load player from database")
	}

	// verify the user actually owns this fleet
	if fleet.PlayerNum != player.Num {
		return nil, fmt.Errorf("%s does not own fleet %s", player, fleet)
	}

	switch mapObjectType {
	case cs.MapObjectTypePlanet:
		err = pu.transferCargoFleetPlanet(fleet, destID, transferAmount)
	case cs.MapObjectTypeFleet:
		err = pu.transferCargoFleetFleet(fleet, destID, transferAmount)
	}

	return fleet, err
}

// transfer cargo from a fleet to/from a planet
func (pu *playerUpdater) transferCargoFleetPlanet(fleet *cs.Fleet, destID int64, transferAmount cs.Cargo) error {
	// find the planet planet by id so we can perform the transfer
	planet, err := pu.db.GetPlanet(destID)
	if err != nil {
		log.Error().Err(err).Msg("get planet from database")
		return fmt.Errorf("failed to get planet from database")
	}

	if planet == nil {
		return fmt.Errorf("no planet for id %d found. %w", destID, errNotFound)
	}

	if err := pu.orderer.TransferPlanetCargo(fleet, planet, transferAmount); err != nil {
		log.Error().Err(err).Msg("transfer cargo")
		return fmt.Errorf("failed to transfer cargo")
	}

	if err := pu.db.UpdatePlanet(planet); err != nil {
		log.Error().Err(err).Int64("ID", planet.ID).Msg("update planet in database")
		return fmt.Errorf("failed to save planet to database")
	}

	if err := pu.db.UpdateFleet(fleet); err != nil {
		log.Error().Err(err).Msg("update fleet in database")
		return fmt.Errorf("failed to update fleet in database")
	}

	log.Info().
		Int64("GameID", fleet.GameID).
		Int("Player", fleet.PlayerNum).
		Str("Fleet", fleet.Name).
		Str("Planet", planet.Name).
		Str("TransferAmount", fmt.Sprintf("%v", transferAmount)).
		Msgf("%s transfered %v to/from Planet %s", fleet.Name, transferAmount, planet.Name)

	// success
	return nil
}

// transfer cargo from a fleet to/from a fleet
func (pu *playerUpdater) transferCargoFleetFleet(fleet *cs.Fleet, destID int64, transferAmount cs.Cargo) error {
	// find the dest dest by id so we can perform the transfer
	dest, err := pu.db.GetFleet(destID)
	if err != nil {
		log.Error().Err(err).Msg("get dest fleet from database")
		return fmt.Errorf("failed to get dest fleet from database")
	}

	if dest == nil {
		return fmt.Errorf("no fleet for id %d found. %w", destID, errNotFound)
	}

	if err := pu.orderer.TransferFleetCargo(fleet, dest, transferAmount); err != nil {
		return fmt.Errorf(err.Error())
	}

	if err := pu.db.UpdateFleet(dest); err != nil {
		log.Error().Err(err).Msg("update fleet in database")
		return fmt.Errorf("failed to update dest fleet in database")
	}

	if err := pu.db.UpdateFleet(fleet); err != nil {
		log.Error().Err(err).Msg("update fleet in database")
		return fmt.Errorf("failed to update fleet in database")
	}

	log.Info().
		Int64("GameID", fleet.GameID).
		Int("Player", fleet.PlayerNum).
		Str("Fleet", fleet.Name).
		Str("Planet", dest.Name).
		Str("TransferAmount", fmt.Sprintf("%v", transferAmount)).
		Msgf("%s transfered %v to/from Planet %s", fleet.Name, transferAmount, dest.Name)

	// success
	return nil
}
