//go:build !wasi && !wasm

package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
)

func NewshipDesignServiceHandler(db DBConnection) craig_starsv1connect.ShipDesignServiceHandler {
	return &shipDesignService{db}
}

type shipDesignService struct {
	db DBConnection
}

// GetShipDesign returns a single design by num for the current player.
func (s *shipDesignService) GetShipDesign(ctx context.Context, req *connect.Request[craig_starsv1.GetShipDesignRequest]) (*connect.Response[craig_starsv1.GetShipDesignResponse], error) {
	dbClient := contextDb(ctx)
	gamePlayer := contextGamePlayer(ctx)

	design, err := dbClient.GetShipDesignByNum(ctx, req.Msg.GameId, gamePlayer.Num, int(req.Msg.Num))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get design"))
	}
	if design == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("design not found"))
	}

	return connect.NewResponse(&craig_starsv1.GetShipDesignResponse{
		Design: converter.C.ConvertCSShipDesign(design),
	}), nil
}

// GetShipDesigns returns all designs for the current player.
func (s *shipDesignService) GetShipDesigns(ctx context.Context, req *connect.Request[craig_starsv1.GetShipDesignsRequest]) (*connect.Response[craig_starsv1.GetShipDesignsResponse], error) {
	dbClient := contextDb(ctx)
	gamePlayer := contextGamePlayer(ctx)

	designs, err := dbClient.GetShipDesignsForPlayer(ctx, req.Msg.GameId, gamePlayer.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get designs"))
	}

	return connect.NewResponse(&craig_starsv1.GetShipDesignsResponse{
		Design: converter.C.ConvertCSShipDesigns(designs),
	}), nil
}

// CreateShipDesign validates, computes spec, assigns num, and saves a new design.
func (s *shipDesignService) CreateShipDesign(ctx context.Context, req *connect.Request[craig_starsv1.CreateShipDesignRequest]) (*connect.Response[craig_starsv1.CreateShipDesignResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	player, err := dbClient.GetLightPlayerForGameWithDesigns(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load player: %w", err))
	}

	design := converter.C.ConvertShipDesign(req.Msg.Design)
	if err := design.Validate(&game.Rules, player); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	for _, existing := range player.Designs {
		if existing.Name == design.Name {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("design %s already exists", design.Name))
		}
	}

	design.PlayerNum = player.Num
	design.GameID = game.ID
	design.Num = player.GetNextDesignNum(player.Designs)

	// Compute Spec then Validate
	design.Spec, err = cs.ComputeShipDesignSpec(&game.Rules, player.TechLevels, player.Race.Spec, design)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", gamePlayer.Num).Str("DesignName", design.Name).Msg("compute ship design spec")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := dbWriteClient.SaveShipDesign(ctx, design); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", gamePlayer.Num).Str("DesignName", design.Name).Msg("save new player design")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	log.Info().Int64("GameID", design.GameID).Int("PlayerNum", player.Num).Str("DesignName", design.Name).Msg("created player design")

	return connect.NewResponse(&craig_starsv1.CreateShipDesignResponse{
		Design: converter.C.ConvertCSShipDesign(design),
	}), nil
}

// UpdateShipDesign updates an existing design if not in use or transferred.
func (s *shipDesignService) UpdateShipDesign(ctx context.Context, req *connect.Request[craig_starsv1.UpdateShipDesignRequest]) (*connect.Response[craig_starsv1.UpdateShipDesignResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	design := converter.C.ConvertShipDesign(req.Msg.Design)

	// Load existingDesign design
	existingDesign, err := dbClient.GetShipDesignByNum(ctx, game.ID, gamePlayer.Num, design.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if existingDesign == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("design not found"))
	}

	if existingDesign.Spec.NumInstances > 0 {
		log.Error().
			Int64("GameID", existingDesign.GameID).
			Int64("ID", existingDesign.ID).
			Int("PlayerNum", existingDesign.PlayerNum).
			Str("DesignName", design.Name).
			Msgf("design in use (%d instances), cannot update", existingDesign.Spec.NumInstances)

		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("design is in use, cannot update"))
	}
	if existingDesign.OriginalPlayerNum != cs.None {
		log.Error().
			Int64("GameID", existingDesign.GameID).
			Int64("ID", existingDesign.ID).
			Int("PlayerNum", existingDesign.PlayerNum).
			Str("DesignName", design.Name).
			Msg("design is transferred from another player, cannot update")

		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("design is transferred, cannot update"))
	}

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load player: %w", err))
	}
	player.Race.Spec = cs.ComputeRaceSpec(&player.Race, &game.Rules)

	if err := design.Validate(&game.Rules, player); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Recompute spec and save
	design.PlayerNum = gamePlayer.Num
	design.GameID = game.ID
	design.Spec, err = cs.ComputeShipDesignSpec(&game.Rules, player.TechLevels, player.Race.Spec, design)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := dbWriteClient.SaveShipDesign(ctx, design); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Str("DesignName", design.Name).Msg("update player design")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	log.Info().Int64("GameID", design.GameID).Int("PlayerNum", player.Num).Str("DesignName", design.Name).Msg("updated player design")

	return connect.NewResponse(&craig_starsv1.UpdateShipDesignResponse{
		Design: converter.C.ConvertCSShipDesign(design),
	}), nil
}

// DeleteShipDesign removes a design and updates fleets/planets accordingly.
func (s *shipDesignService) DeleteShipDesign(ctx context.Context, req *connect.Request[craig_starsv1.DeleteShipDesignRequest]) (*connect.Response[craig_starsv1.DeleteShipDesignResponse], error) {
	dbClient := contextDb(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	design, err := dbClient.GetShipDesignByNum(ctx, game.ID, gamePlayer.Num, int(req.Msg.Num))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if design == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("design not found"))
	}

	if design.CannotDelete {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("shipDesign cannot be deleted"))
	}

	playerFleets, err := dbClient.GetFleetsForPlayer(ctx, game.ID, gamePlayer.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	playerDesigns, err := dbClient.GetShipDesignsForPlayer(ctx, game.ID, gamePlayer.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	playerPlanets, err := dbClient.GetPlanetsForPlayer(ctx, game.ID, gamePlayer.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load player: %w", err))
	}
	player.Race.Spec = cs.ComputeRaceSpec(&player.Race, &game.Rules)

	fleetsToDelete := []*cs.Fleet{}
	fleetsToUpdate := []*cs.Fleet{}
	leftoverPlayerFleets := []*cs.Fleet{}
	for _, fleet := range playerFleets {
		// find any tokens using this design
		updatedTokens := make([]cs.ShipToken, 0, len(fleet.Tokens))
		for _, token := range fleet.Tokens {
			if token.DesignNum != design.Num {
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
				fleet.InjectDesigns(playerDesigns)
				fleet.Spec = cs.ComputeFleetSpec(&game.Rules, player, fleet)
				fleetsToUpdate = append(fleetsToUpdate, fleet)
			}
			leftoverPlayerFleets = append(leftoverPlayerFleets, fleet)
		}
	}

	// remove this design from any production queues
	planetsToUpdate := []*cs.Planet{}
	for _, planet := range playerPlanets {
		newQueue := make([]cs.ProductionQueueItem, 0, len(planet.ProductionQueue))
		for _, item := range planet.ProductionQueue {
			if item.DesignNum == design.Num {
				// remove this design from the queue and mark this planet for update
				planetsToUpdate = append(planetsToUpdate, planet)
				continue
			}
			newQueue = append(newQueue, item)
		}
		planet.ProductionQueue = newQueue
	}

	// update fleets to remove tokens, delete fleets without tokens, update planets with production queue changes
	if err := s.db.WrapInTransaction(func(c db.Client) error {

		for _, fleet := range fleetsToUpdate {
			if err := c.SaveFleet(ctx, fleet); err != nil {
				return fmt.Errorf("update fleet in database: %w", err)
			}
			log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", design.Num).Msgf("updated fleet %s after deleting design", fleet.Name)
		}

		for _, fleet := range fleetsToDelete {
			if err := c.DeleteFleet(ctx, fleet.ID); err != nil {
				return fmt.Errorf("delete fleet from database: %w", err)
			}
			log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", design.Num).Msgf("deleted fleet %s after deleting design", fleet.Name)
		}

		for _, planet := range planetsToUpdate {
			if err := c.SavePlanet(ctx, planet); err != nil {
				return fmt.Errorf("update planet in database: %w", err)
			}
			log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", design.Num).Msgf("updated planet %s after deleting design", planet.Name)

		}

		if err := c.DeleteShipDesign(ctx, design.ID); err != nil {
			return fmt.Errorf("delete design from database: %w", err)
		}
		log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", design.Num).Msgf("deleted design %s", design.Name)

		return nil
	}); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("delete design from database")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// split the player fleets into fleets and starbases
	fleets := make([]*cs.Fleet, 0, len(leftoverPlayerFleets))
	starbases := make([]*cs.Fleet, 0)
	for i := range leftoverPlayerFleets {
		fleet := leftoverPlayerFleets[i]
		if fleet.Starbase {
			starbases = append(starbases, fleet)
		} else {
			fleets = append(fleets, fleet)
		}
	}

	return connect.NewResponse(&craig_starsv1.DeleteShipDesignResponse{
		Fleets:    converter.C.ConvertCSFleets(fleets),
		Starbases: converter.C.ConvertCSFleets(starbases),
		Planets:   converter.C.ConvertCSPlanets(playerPlanets),
	}), nil
}
