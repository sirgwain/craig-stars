package server

import (
	"context"
	"fmt"

	"log/slog"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
)

func NewFleetServiceHandler(db DBConnection) craig_starsv1connect.FleetServiceHandler {
	return &fleetService{db}
}

type fleetService struct {
	db DBConnection
}

func (s *fleetService) GetFleet(ctx context.Context, req *connect.Request[craig_starsv1.GetFleetRequest]) (*connect.Response[craig_starsv1.GetFleetResponse], error) {
	dbClient := contextDb(ctx)
	gamePlayer := contextGamePlayer(ctx)

	fleet, err := dbClient.GetFleetByNum(ctx, req.Msg.GameId, gamePlayer.Num, int(req.Msg.FleetNum))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get fleet from db"))
	}

	if fleet.PlayerNum != gamePlayer.Num {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("player doesn't own fleet"))
	}

	return connect.NewResponse(&craig_starsv1.GetFleetResponse{
		Fleet: converter.C.ConvertCSFleet(fleet),
	}), nil
}

func (s *fleetService) UpdateFleetOrders(ctx context.Context, req *connect.Request[craig_starsv1.UpdateFleetOrdersRequest]) (*connect.Response[craig_starsv1.UpdateFleetOrdersResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	fleet, err := dbClient.GetFleetByNum(ctx, req.Msg.GameId, gamePlayer.Num, int(req.Msg.FleetNum))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get fleet from db"))
	}

	if fleet.PlayerNum != gamePlayer.Num {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("player doesn't own fleet"))
	}

	// load the full player to update fleet production estimates
	player, err := dbClient.GetLightPlayerForGameWithDesigns(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	fleet.InjectDesigns(player.Designs)

	orders := converter.C.ConvertFleetOrders(req.Msg.FleetOrders)
	orderer := cs.NewOrderer()
	orderer.UpdateFleetOrders(player, fleet, *orders)

	if err := dbWriteClient.SaveFleet(ctx, fleet); err != nil {
		slog.Error("update fleet in database", slog.Any("error", err), slog.Int64("ID", fleet.ID))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.UpdateFleetOrdersResponse{
		Fleet: converter.C.ConvertCSFleet(fleet),
	}), nil
}

// MergeFleets implements craig_starsv1connect.FleetServiceHandler.
func (s *fleetService) MergeFleets(ctx context.Context, req *connect.Request[craig_starsv1.MergeFleetsRequest]) (*connect.Response[craig_starsv1.MergeFleetsResponse], error) {
	dbClient := contextDb(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	// Get the target fleet (the one we're merging into)
	fleet, err := dbClient.GetFleetByNum(ctx, req.Msg.GameId, gamePlayer.Num, int(req.Msg.FleetNum))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get fleet from db"))
	}
	if fleet == nil || fleet.PlayerNum != gamePlayer.Num {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("fleet not found"))
	}

	// Validate that the source fleet is not included in the merge request
	for _, num := range req.Msg.FleetNums {
		if fleet.Num == int(num) {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("cannot include target fleet in merge request"))
		}
	}

	// Convert fleet nums to int slice
	fleetNums := make([]int, len(req.Msg.FleetNums))
	for i, num := range req.Msg.FleetNums {
		fleetNums[i] = int(num)
	}

	// Get the fleets to merge
	fleets, err := dbClient.GetFleetsByNums(ctx, game.ID, gamePlayer.Num, fleetNums)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get fleets for merge"))
	}

	// Load the full player to get designs
	player, err := dbClient.GetLightPlayerForGameWithDesigns(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Execute the merge
	orderer := cs.NewOrderer()
	fleets = append([]*cs.Fleet{fleet}, fleets...)

	updatedFleet, err := orderer.Merge(&game.Rules, player, fleets)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Save the fleets
	if _, err := s.saveSplitMergeFleets(ctx, player, fleets); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.MergeFleetsResponse{
		Fleet:          converter.C.ConvertCSFleet(updatedFleet),
		CargoTransfers: converter.C.ConvertCSCargoTransfers(player.CargoTransfers),
	}), nil
}

// RenameFleet implements craig_starsv1connect.FleetServiceHandler.
func (s *fleetService) RenameFleet(ctx context.Context, req *connect.Request[craig_starsv1.RenameFleetRequest]) (*connect.Response[craig_starsv1.RenameFleetResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	gamePlayer := contextGamePlayer(ctx)

	// Get the fleet
	fleet, err := dbClient.GetFleetByNum(ctx, req.Msg.GameId, gamePlayer.Num, int(req.Msg.FleetNum))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get fleet from db"))
	}
	if fleet == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("fleet not found"))
	}

	if fleet.PlayerNum != gamePlayer.Num {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("player doesn't own fleet"))
	}

	// Update fleet name
	fleet.BaseName = req.Msg.Name
	fleet.Name = fmt.Sprintf("%s #%d", req.Msg.Name, fleet.Num)

	// Save the fleet
	if err := dbWriteClient.SaveFleet(ctx, fleet); err != nil {
		slog.Error("update fleet in database", slog.Any("error", err), slog.Int64("ID", fleet.ID))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.RenameFleetResponse{
		Fleet: converter.C.ConvertCSFleet(fleet),
	}), nil
}

// SplitAllFleets implements craig_starsv1connect.FleetServiceHandler.
func (s *fleetService) SplitAllFleets(ctx context.Context, req *connect.Request[craig_starsv1.SplitAllFleetsRequest]) (*connect.Response[craig_starsv1.SplitAllFleetsResponse], error) {
	dbClient := contextDb(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	// Get the fleet
	fleet, err := dbClient.GetFleetByNum(ctx, req.Msg.GameId, gamePlayer.Num, int(req.Msg.FleetNum))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get fleet from db"))
	}
	if fleet == nil || fleet.PlayerNum != gamePlayer.Num {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("fleet not found"))
	}

	// Load the full player to get designs and fleets
	player, err := dbClient.GetLightPlayerForGameWithDesigns(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	fleets, err := dbClient.GetFleetsForPlayer(ctx, game.ID, gamePlayer.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get fleets for player"))
	}

	// Execute the split all
	orderer := cs.NewOrderer()
	newFleets, err := orderer.SplitAll(&game.Rules, player, fleets, fleet)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Save all the fleets
	remainingFleets, err := s.saveSplitMergeFleets(ctx, player, append(newFleets, fleet))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.SplitAllFleetsResponse{
		Fleets:         converter.C.ConvertCSFleets(remainingFleets),
		CargoTransfers: converter.C.ConvertCSCargoTransfers(player.CargoTransfers),
	}), nil
}

// SplitFleet implements craig_starsv1connect.FleetServiceHandler.
func (s *fleetService) SplitFleet(ctx context.Context, req *connect.Request[craig_starsv1.SplitFleetRequest]) (*connect.Response[craig_starsv1.SplitFleetResponse], error) {
	dbClient := contextDb(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	// Get the source fleet
	sourceFleet, err := dbClient.GetFleetByNum(ctx, req.Msg.GameId, gamePlayer.Num, int(req.Msg.SourceFleetNum))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get source fleet from db"))
	}
	if sourceFleet == nil || sourceFleet.PlayerNum != gamePlayer.Num {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("source fleet not found"))
	}

	// Load the dest fleet if specified
	var destFleet *cs.Fleet
	if req.Msg.DestFleetNum != 0 {
		destFleet, err = dbClient.GetFleetByNum(ctx, game.ID, gamePlayer.Num, int(req.Msg.DestFleetNum))
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get dest fleet from db"))
		}
	}

	// Load the full player to get designs and fleets
	player, err := dbClient.GetLightPlayerForGameWithDesigns(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	playerFleets, err := dbClient.GetFleetsForPlayer(ctx, game.ID, gamePlayer.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get fleets for player"))
	}

	// Convert proto tokens to CS tokens
	sourceTokens := converter.C.ConvertShipTokens(req.Msg.SourceTokens)
	destTokens := converter.C.ConvertShipTokens(req.Msg.DestTokens)

	// Convert transfer amount manually
	transferAmount := cs.CargoTransferRequest{
		Cargo: converter.C.ConvertCargo(req.Msg.TransferAmount),
		Fuel:  int(req.Msg.FuelTransferAmount),
	}

	// Execute the split
	orderer := cs.NewOrderer()
	splitFleetRequest := cs.SplitFleetRequest{
		Source:         sourceFleet,
		Dest:           destFleet,
		SourceTokens:   sourceTokens,
		DestTokens:     destTokens,
		DestBaseName:   req.Msg.DestBaseName,
		TransferAmount: transferAmount,
	}

	source, dest, err := orderer.SplitFleet(&game.Rules, player, playerFleets, splitFleetRequest)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Save the updated fleets back to the database
	if err := s.db.WrapInTransaction(func(c db.Client) error {
		if source.Delete {
			if err := c.DeleteFleet(ctx, source.ID); err != nil {
				slog.Error("delete fleet in database", slog.Any("error", err))
				return err
			}
		} else {
			if err := c.SaveFleet(ctx, source); err != nil {
				slog.Error("update fleet in database", slog.Any("error", err))
				return err
			}
		}

		// It's possible the user sent a "split" request with no dest but didn't split any
		// in this case dest will be marked for deletion but won't ever have been saved
		// to the database, so just ignore it
		if dest.Delete && dest.ID != 0 {
			if err := c.DeleteFleet(ctx, dest.ID); err != nil {
				slog.Error("delete fleet in database", slog.Any("error", err))
				return err
			}
		} else {
			dest.GameID = game.ID
			if err := c.SaveFleet(ctx, dest); err != nil {
				slog.Error("update fleet in database", slog.Any("error", err))
				return err
			}
		}

		if err := c.UpdatePlayerCargoTransfers(ctx, player); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update database after split"))
	}

	if source.Delete {
		// we are deleting the source, so the new source to return to the player
		// is the dest (it should not be deleted)
		source = dest
		dest = nil
	} else if dest.Delete {
		dest = nil
	}

	var responseSource, responseDest *craig_starsv1.Fleet
	if source != nil {
		responseSource = converter.C.ConvertCSFleet(source)
	}
	if dest != nil {
		responseDest = converter.C.ConvertCSFleet(dest)
	}

	return connect.NewResponse(&craig_starsv1.SplitFleetResponse{
		Source:         responseSource,
		Dest:           responseDest,
		CargoTransfers: converter.C.ConvertCSCargoTransfers(player.CargoTransfers),
	}), nil
}

// saveSplitMergeFleets saves or deletes fleets from a list and updates a player's byhand cargo transfers in a single transaction
func (s *fleetService) saveSplitMergeFleets(ctx context.Context, player *cs.Player, fleets []*cs.Fleet) ([]*cs.Fleet, error) {
	remainingFleets := make([]*cs.Fleet, 0, len(fleets))
	// save/delete the fleets
	err := s.db.WrapInTransaction(func(c db.Client) error {
		// first delete fleets. This way if we end up creating a new fleet
		// with an in use unique index, we'll delete the old one first
		for _, fleet := range fleets {
			if !fleet.Delete {
				remainingFleets = append(remainingFleets, fleet)
				continue
			}
			// possible a fleet was created and destroyed in one turn
			if fleet.ID != 0 {
				if err := c.DeleteFleet(ctx, fleet.ID); err != nil {
					return fmt.Errorf("delete fleet: %w", err)
				}
			}
		}

		for _, fleet := range remainingFleets {
			fleet.GameID = player.GameID
			if err := c.SaveFleet(ctx, fleet); err != nil {
				return fmt.Errorf("save fleet: %w", err)
			}
		}

		// update by hand transfers
		if err := c.UpdatePlayerCargoTransfers(ctx, player); err != nil {
			return err
		}

		return nil
	})

	return remainingFleets, err
}

// TransferCargo implements craig_starsv1connect.FleetServiceHandler.
func (s *fleetService) TransferCargo(ctx context.Context, req *connect.Request[craig_starsv1.TransferCargoRequest]) (*connect.Response[craig_starsv1.TransferCargoResponse], error) {
	dbClient := contextDb(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	// Get the fleet
	fleet, err := dbClient.GetFleetByNum(ctx, req.Msg.GameId, gamePlayer.Num, int(req.Msg.FleetNum))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get fleet from db"))
	}
	if fleet == nil || fleet.PlayerNum != gamePlayer.Num {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("fleet not found"))
	}

	// Load the full player with designs
	player, err := dbClient.GetPlayerForGame(ctx, game.ID, gamePlayer.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	player.Race.Spec = cs.ComputeRaceSpec(&player.Race, &game.Rules)
	fleet.InjectDesigns(player.Designs)

	// Convert transfer amount manually
	transferAmount := cs.CargoTransferRequest{
		Cargo: converter.C.ConvertCargo(req.Msg.TransferAmount),
		Fuel:  int(req.Msg.FuelTransferAmount),
	}

	// Convert MapObject
	mo := converter.C.ConvertMapObject(req.Msg.Mo)
	var dest cs.CargoHolder
	var destIsIntel bool

	// Handle different destination types - simplified implementation
	switch mo.Type {
	case cs.MapObjectTypeNone:
		// Jettison - dest is nil
		dest = nil
	case cs.MapObjectTypePlanet:
		// Get planet
		planet, err := dbClient.GetPlanetByNum(ctx, game.ID, mo.Num)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		if planet == nil {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("planet not found"))
		}
		if !planet.OwnedBy(gamePlayer.Num) {
			destPlanet := player.GetPlanetIntel(mo.Num)
			if destPlanet.Spec.HasStarbase {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("cannot transfer cargo to/from an enemy planet with a starbase"))
			}
			dest = destPlanet
			destIsIntel = true
		} else {
			dest = planet
		}
	case cs.MapObjectTypeFleet:
		// Get fleet
		if mo.PlayerNum == gamePlayer.Num {
			destFleet, err := dbClient.GetFleetByNum(ctx, game.ID, mo.PlayerNum, mo.Num)
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			if destFleet == nil {
				return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("destination fleet not found"))
			}
			destFleet.InjectDesigns(player.Designs)
			dest = destFleet
		} else {
			dest = player.GetFleetIntel(mo.PlayerNum, mo.Num)
			destIsIntel = true
		}
	case cs.MapObjectTypeMineralPacket:
		// Get mineralPacket
		if mo.PlayerNum == gamePlayer.Num {
			destPacket, err := dbClient.GetMineralPacketByNum(ctx, game.ID, mo.PlayerNum, mo.Num)
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			if destPacket == nil {
				return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("destination fleet not found"))
			}
			dest = destPacket
		} else {
			dest = player.GetMineralPacketIntel(mo.PlayerNum, mo.Num)
			destIsIntel = true
		}
	case cs.MapObjectTypeSalvage:
		// Get salvage
		salvage := player.GetSalvageIntel(mo.Num)
		if salvage == nil {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("salvage not found"))
		}
		dest = salvage
		destIsIntel = true
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("transfer type %s not yet implemented", mo.Type))
	}

	// Execute the transfer
	orderer := cs.NewOrderer()
	if err := orderer.TransferByHand(&game.Rules, player, fleet, dest, transferAmount); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// update this planet and the player's intel
	if planet, ok := dest.(*cs.Planet); ok && planet != nil && !planet.OwnedBy(player.Num) {
		player.PlanetIntels[planet.Num-1] = planet
	}

	// Save changes in a transaction
	err = s.db.WrapInTransaction(func(c db.Client) error {
		if err := c.SaveFleet(ctx, fleet); err != nil {
			return err
		}

		if destFleet, ok := dest.(*cs.Fleet); ok && destFleet != nil && !destIsIntel {
			if err := c.SaveFleet(ctx, destFleet); err != nil {
				return err
			}
		}

		if destPlanet, ok := dest.(*cs.Planet); ok && destPlanet != nil && !destIsIntel {
			if err := c.SavePlanet(ctx, destPlanet); err != nil {
				return err
			}
		}

		if destFleetIntel, ok := dest.(*cs.Fleet); ok && destFleetIntel != nil {
			if err := c.UpdatePlayerFleetIntels(ctx, player); err != nil {
				return err
			}
		}

		if destPlanetIntel, ok := dest.(*cs.Planet); ok && destPlanetIntel != nil {
			if err := c.UpdatePlayerPlanetIntels(ctx, player); err != nil {
				return err
			}
		}

		if destMineralPacket, ok := dest.(*cs.MineralPacket); ok && destMineralPacket != nil && !destIsIntel {
			if err := c.SaveMineralPacket(ctx, destMineralPacket); err != nil {
				return err
			}
		}

		if destMineralPacketIntel, ok := dest.(*cs.MineralPacket); ok && destMineralPacketIntel != nil {
			if err := c.UpdatePlayerMineralPacketIntels(ctx, player); err != nil {
				return err
			}
		}

		if destSalvage, ok := dest.(*cs.Salvage); ok && destSalvage != nil {
			if err := c.UpdatePlayerSalvageIntels(ctx, player); err != nil {
				return err
			}
		}

		if err := c.UpdatePlayerCargoTransfers(ctx, player); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Prepare response
	response := &craig_starsv1.TransferCargoResponse{
		Player: converter.C.ConvertCSPlayer(player),
		Fleet:  converter.C.ConvertCSFleet(fleet),
	}

	// Set the appropriate destination in the oneof field
	if destFleet, ok := dest.(*cs.Fleet); ok && destFleet != nil {
		response.Dest = &craig_starsv1.TransferCargoResponse_DestFleet{
			DestFleet: converter.C.ConvertCSFleet(destFleet),
		}
	} else if destPlanet, ok := dest.(*cs.Planet); ok && destPlanet != nil {
		response.Dest = &craig_starsv1.TransferCargoResponse_DestPlanet{
			DestPlanet: converter.C.ConvertCSPlanet(destPlanet),
		}
	} else if destMineralPacket, ok := dest.(*cs.MineralPacket); ok && destMineralPacket != nil {
		response.Dest = &craig_starsv1.TransferCargoResponse_DestMineralPacket{
			DestMineralPacket: converter.C.ConvertCSMineralPacket(destMineralPacket),
		}
	} else if destSalvage, ok := dest.(*cs.Salvage); ok && destSalvage != nil {
		response.Dest = &craig_starsv1.TransferCargoResponse_DestSalvage{
			DestSalvage: converter.C.ConvertCSSalvage(destSalvage),
		}
	}

	return connect.NewResponse(response), nil
}
