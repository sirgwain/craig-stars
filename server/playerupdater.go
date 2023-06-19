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
	updatePlayerOrders(gameID, userID int64, orders cs.PlayerOrders) (*cs.Player, []*cs.Planet, error)
	updatePlanetOrders(userID int64, planetID int64, orders cs.PlanetOrders) (*cs.Planet, error)
	updateFleetOrders(userID int64, fleetID int64, orders cs.FleetOrders) (*cs.Fleet, error)
	updateMineFieldOrders(userID int64, mineFieldID int64, orders cs.MineFieldOrders) (*cs.MineField, error)
	transferCargo(userID int64, fleetID int64, destID int64, mapObjectType cs.MapObjectType, transferAmount cs.Cargo) (*cs.Fleet, error)
	createShipDesign(gameID, userID int64, design *cs.ShipDesign) (*cs.ShipDesign, error)
	updateShipDesign(gameID, userID int64, design *cs.ShipDesign) (*cs.ShipDesign, error)
	deleteShipDesign(gameID, userID int64, num int) (fleets, starbases []*cs.Fleet, err error)
}

func newPlayerUpdater(db DBClient) PlayerUpdater {
	return &playerUpdater{orderer: cs.NewOrderer(), db: db}
}

// update a player's orders (i.e. research settings) and return the updated planets
func (pu *playerUpdater) updatePlayerOrders(gameID, userID int64, orders cs.PlayerOrders) (*cs.Player, []*cs.Planet, error) {
	player, err := pu.db.GetPlayerForGame(gameID, userID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", gameID).Int64("UserID", userID).Msg("load player from database")
		return nil, nil, fmt.Errorf("failed to load player from database")
	}

	if player.UserID != userID {
		return nil, nil, fmt.Errorf("user %d does not control player %d", userID, player.Num)
	}

	planets, err := pu.db.GetPlanetsForPlayer(player.GameID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Msg("loading player planets from database")
		return nil, nil, fmt.Errorf("failed to load player planets from database")
	}

	// update this player's orders
	rules, err := pu.db.GetRulesForGame(player.GameID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Msg("loading game rules from database")
		return nil, nil, fmt.Errorf("failed to load game rules from database")
	}

	pu.orderer.UpdatePlayerOrders(player, planets, orders, rules)

	if err := pu.db.UpdateLightPlayer(player); err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Msg("updating player orders in database")
		return nil, nil, fmt.Errorf("failed to update player orders in database")
	}

	for _, planet := range planets {
		if planet.Dirty {
			// TODO: only update the planet spec? that's all that changes
			if err := pu.db.UpdatePlanet(planet); err != nil {
				log.Error().Err(err).Int64("ID", player.ID).Msg("updating player planet in database")
				return nil, nil, fmt.Errorf("failed to update player planet in database")
			}
		}
	}

	return player, planets, nil
}

// update the orders for a planet, i.e. production queue and research
func (pu *playerUpdater) updatePlanetOrders(userID int64, planetID int64, orders cs.PlanetOrders) (*cs.Planet, error) {

	// find the planet planet by id
	planet, err := pu.db.GetPlanet(planetID)
	if err != nil {
		log.Error().Err(err).Int64("ID", planetID).Msg("load planet from database")
		return nil, fmt.Errorf("failed to load planet from database")
	}

	player, err := pu.db.GetLightPlayerForGame(planet.GameID, userID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", planet.GameID).Int64("UserID", userID).Msg("load player from database")
		return nil, fmt.Errorf("failed to load player from database")
	}

	// verify the user actually owns this planet
	if planet.PlayerNum != player.Num {
		return nil, fmt.Errorf("%s does not own planet %s", player, planet)
	}

	pu.orderer.UpdatePlanetOrders(player, planet, orders)

	if err := pu.db.UpdatePlanet(planet); err != nil {
		log.Error().Err(err).Int64("ID", planet.ID).Msg("update planet in database")
		return nil, fmt.Errorf("failed to save planet to database")
	}

	return planet, nil
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

func (pu *playerUpdater) createShipDesign(gameID, userID int64, design *cs.ShipDesign) (*cs.ShipDesign, error) {
	player, err := pu.db.GetPlayerWithDesignsForGame(gameID, userID)

	if err != nil || player == nil {
		log.Error().Err(err).Int64("GameID", gameID).Int64("UserID", userID).Msg("load player from database")
		return nil, fmt.Errorf("failed to load player from database")
	}

	rules, err := pu.db.GetRulesForGame(player.GameID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Msg("loading game rules from database")
		return nil, fmt.Errorf("failed to load game rules from database")
	}

	if err := design.Validate(rules, player); err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Str("DesignName", design.Name).Msg("validate new player design")
		return nil, fmt.Errorf("invalid design, %w", err)
	}

	design.PlayerNum = player.Num
	design.GameID = player.GameID
	design.Num = player.GetNextDesignNum()

	if err := pu.db.CreateShipDesign(design); err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Str("DesignName", design.Name).Msg("save new player design")
		return nil, fmt.Errorf("failed to save design")
	}

	return design, nil
}

func (pu *playerUpdater) updateShipDesign(gameID, userID int64, design *cs.ShipDesign) (*cs.ShipDesign, error) {
	player, err := pu.db.GetLightPlayerForGame(gameID, userID)

	if err != nil || player == nil {
		log.Error().Err(err).Int64("GameID", gameID).Int64("UserID", userID).Msg("load player from database")
		return nil, fmt.Errorf("failed to load player from database")
	}

	existingShipDesign, err := pu.db.GetShipDesign(design.ID)
	if err != nil {
		log.Error().Err(err).Int64("ID", design.ID).Msg("get shipdesign from database")
		return nil, fmt.Errorf("failed to load design from database")
	}

	if existingShipDesign == nil {
		return nil, fmt.Errorf("design doesn't exist")
	}

	if player.Num != existingShipDesign.PlayerNum {
		return nil, fmt.Errorf("user %d does not own shipdesign %d", userID, existingShipDesign.ID)
	}

	rules, err := pu.db.GetRulesForGame(player.GameID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Msg("loading game rules from database")
		return nil, fmt.Errorf("failed to load game rules from database")
	}

	if err := design.Validate(rules, player); err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Str("DesignName", design.Name).Msg("validate player design")
		return nil, fmt.Errorf("invalid design, %w", err)
	}

	design.PlayerNum = player.Num
	design.GameID = player.GameID
	design.Spec = cs.ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)

	if err := pu.db.UpdateShipDesign(design); err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Str("DesignName", design.Name).Msg("save design")
		return nil, fmt.Errorf("failed to save design")
	}

	return design, nil
}

func (pu *playerUpdater) deleteShipDesign(gameID, userID int64, num int) (fleets, starbases []*cs.Fleet, err error) {
	player, err := pu.db.GetPlayerForGame(gameID, userID)

	if err != nil || player == nil {
		log.Error().Err(err).Int64("GameID", gameID).Int64("UserID", userID).Msg("load player from database")
		return nil, nil, fmt.Errorf("failed to load player from database")
	}

	design, err := pu.db.GetShipDesignByNum(gameID, player.Num, num)
	if err != nil || design == nil {
		log.Error().Err(err).Int64("GameID", gameID).Int64("UserID", userID).Int("Num", num).Msg("load design from database")
		return nil, nil, fmt.Errorf("failed to load design from database")
	}

	fleets, err = pu.db.GetFleetsForPlayer(gameID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", gameID).Int64("UserID", userID).Int("Num", num).Msg("load fleets from database")
		return nil, nil, fmt.Errorf("failed to load fleets from database")
	}

	fleetsToDelete := []*cs.Fleet{}
	fleetsToUpdate := []*cs.Fleet{}
	for _, fleet := range fleets {
		// find any tokens using this design
		updatedTokens := make([]cs.ShipToken, 0, len(fleet.Tokens))
		for _, token := range fleet.Tokens {
			if token.DesignNum != num {
				updatedTokens = append(updatedTokens, token)
			}
		}
		// if we have no tokens left, delete the fleet
		if len(updatedTokens) == 0 {
			fleetsToDelete = append(fleetsToDelete, fleet)
		} else {
			// if we have a different number of tokens than we
			// had before, update this fleet
			if len(updatedTokens) != len(fleet.Tokens) {
				fleet.Tokens = updatedTokens
				fleetsToUpdate = append(fleetsToUpdate, fleet)
			}
		}
	}

	if err := pu.db.DeleteShipDesignWithFleets(design.ID, fleetsToUpdate, fleetsToDelete); err != nil {
		log.Error().Err(err).Int64("GameID", gameID).Int64("UserID", userID).Int("Num", num).Msg("delete design from database")
		return nil, nil, fmt.Errorf("failed to delete design from database")
	}

	// log what we did
	log.Info().Int64("GameID", gameID).Int64("UserID", userID).Int("Num", num).Msgf("deleted design %s", design.Name)

	for _, fleet := range fleetsToUpdate {
		log.Info().Int64("GameID", gameID).Int64("UserID", userID).Int("Num", num).Msgf("updated fleet %s after deleting design", fleet.Name)
	}

	for _, fleet := range fleetsToDelete {
		log.Info().Int64("GameID", gameID).Int64("UserID", userID).Int("Num", num).Msgf("deleted fleet %s after deleting design", fleet.Name)
	}

	allFleets, err := pu.db.GetFleetsForPlayer(gameID, player.Num)

	if err != nil {
		log.Error().Err(err).Int64("GameID", gameID).Int64("UserID", userID).Int("Num", num).Msg("load fleets from database")
	}

	// split the player fleets into fleets and starbases
	fleets = make([]*cs.Fleet, 0, len(allFleets))
	starbases = make([]*cs.Fleet, 0)
	for i := range allFleets {
		fleet := allFleets[i]
		if fleet.Starbase {
			starbases = append(starbases, fleet)
		} else {
			fleets = append(fleets, fleet)
		}
	}
	return fleets, starbases, nil
}
