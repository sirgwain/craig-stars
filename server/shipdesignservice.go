package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
	"log/slog"
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

	design := converter.C.ConvertShipDesignP(req.Msg.Design)
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
		slog.Error("compute ship design spec", slog.Any("error", err), slog.Int64("GameID", game.ID), slog.Int("PlayerNum", gamePlayer.Num), slog.String("DesignName", design.Name))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := dbWriteClient.SaveShipDesign(ctx, design); err != nil {
		slog.Error("save new player design", slog.Any("error", err), slog.Int64("GameID", game.ID), slog.Int("PlayerNum", gamePlayer.Num), slog.String("DesignName", design.Name))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	slog.Info("created player design", slog.Int64("GameID", design.GameID), slog.Int("PlayerNum", player.Num), slog.String("DesignName", design.Name))

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

	design := converter.C.ConvertShipDesignP(req.Msg.Design)

	// Load existingDesign design
	existingDesign, err := dbClient.GetShipDesignByNum(ctx, game.ID, gamePlayer.Num, design.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if existingDesign == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("design not found"))
	}

	if existingDesign.Spec.NumInstances > 0 {
		slog.Error("design in use, cannot update",
			slog.Int64("GameID", existingDesign.GameID),
			slog.Int64("ID", existingDesign.ID),
			slog.Int("PlayerNum", existingDesign.PlayerNum),
			slog.String("DesignName", design.Name),
			slog.Int("instances", existingDesign.Spec.NumInstances))

		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("design is in use, cannot update"))
	}
	if existingDesign.OriginalPlayerNum != cs.None {
		slog.Error("design is transferred from another player, cannot update",
			slog.Int64("GameID", existingDesign.GameID),
			slog.Int64("ID", existingDesign.ID),
			slog.Int("PlayerNum", existingDesign.PlayerNum),
			slog.String("DesignName", design.Name))

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
		slog.Error("update player design", slog.Any("error", err), slog.Int64("GameID", game.ID), slog.Int("PlayerNum", player.Num), slog.String("DesignName", design.Name))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	slog.Info("updated player design", slog.Int64("GameID", design.GameID), slog.Int("PlayerNum", player.Num), slog.String("DesignName", design.Name))

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
			slog.Info("updated fleet after deleting design", slog.String("fleet", fleet.Name), slog.Int64("GameID", game.ID), slog.Int("PlayerNum", player.Num), slog.Int("Num", design.Num))
		}

		for _, fleet := range fleetsToDelete {
			if err := c.DeleteFleet(ctx, fleet.ID); err != nil {
				return fmt.Errorf("delete fleet from database: %w", err)
			}
			slog.Info("deleted fleet after deleting design", slog.String("fleet", fleet.Name), slog.Int64("GameID", game.ID), slog.Int("PlayerNum", player.Num), slog.Int("Num", design.Num))
		}

		for _, planet := range planetsToUpdate {
			if err := c.SavePlanet(ctx, planet); err != nil {
				return fmt.Errorf("update planet in database: %w", err)
			}
			slog.Info("updated planet after deleting design", slog.String("planet", planet.Name), slog.Int64("GameID", game.ID), slog.Int("PlayerNum", player.Num), slog.Int("Num", design.Num))

		}

		if err := c.DeleteShipDesign(ctx, design.ID); err != nil {
			return fmt.Errorf("delete design from database: %w", err)
		}
		slog.Info("deleted design", slog.String("design", design.Name), slog.Int64("GameID", game.ID), slog.Int("PlayerNum", player.Num), slog.Int("Num", design.Num))

		return nil
	}); err != nil {
		slog.Error("delete design from database", slog.Any("error", err), slog.Int64("GameID", game.ID), slog.Int("PlayerNum", player.Num))
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
