//go:build !wasi && !wasm

package server

import (
	"context"
	"fmt"
	"strconv"

	"connectrpc.com/connect"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
	"golang.org/x/sync/singleflight"
)

func NewPlayerServiceHandler(db DBConnection, config config.Config, discordNotifier *discordNotifier) craig_starsv1connect.PlayerServiceHandler {
	return &playerService{db, config, discordNotifier, singleflight.Group{}}
}

type playerService struct {
	db              DBConnection
	config          config.Config
	discordNotifier *discordNotifier
	sf              singleflight.Group
}

func (s *playerService) GetPlayer(ctx context.Context, req *connect.Request[craig_starsv1.GetPlayerRequest]) (*connect.Response[craig_starsv1.GetPlayerResponse], error) {
	c := contextDb(ctx)
	gamePlayer := contextGamePlayer(ctx)

	player, err := c.GetPlayerForGame(ctx, req.Msg.GameId, gamePlayer.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get player from database: %w", err))
	}
	return connect.NewResponse(&craig_starsv1.GetPlayerResponse{
		Player:  converter.C.ConvertCSPlayer(player),
		Designs: converter.C.ConvertCSShipDesigns(player.Designs),
		Intels:  converter.C.ConvertCSIntels(player.Intels),
	}), nil
}

func (s *playerService) GetUniverse(ctx context.Context, req *connect.Request[craig_starsv1.GetUniverseRequest]) (*connect.Response[craig_starsv1.GetUniverseResponse], error) {
	c := contextDb(ctx)
	user := contextUserSession(ctx)

	pmos, err := c.GetPlayerMapObjects(ctx, req.Msg.GameId, user.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get player from database: %w", err))
	}
	return connect.NewResponse(&craig_starsv1.GetUniverseResponse{
		Universe: &craig_starsv1.PlayerUniverse{
			Planets:        converter.C.ConvertCSPlanets(pmos.Planets),
			Fleets:         converter.C.ConvertCSFleets(pmos.Fleets),
			Starbases:      converter.C.ConvertCSFleets(pmos.Starbases),
			Minefields:     converter.C.ConvertCSMinefields(pmos.Minefields),
			MineralPackets: converter.C.ConvertCSMineralPackets(pmos.MineralPackets),
		},
	}), nil
}

// SubmitTurn mirrors REST submitTurn: submit current player's turn and possibly generate a new turn.
func (s *playerService) SubmitTurn(ctx context.Context, req *connect.Request[craig_starsv1.SubmitTurnRequest]) (*connect.Response[craig_starsv1.SubmitTurnResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	// update the player in the db
	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("failed to get player")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get player: %w", err))
	}

	// submit the turn
	player.SubmittedTurn = true
	if err := dbWriteClient.SubmitPlayerTurn(ctx, game.ID, player.Num, true); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("update player submit")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("submit turn"))
	}

	// only allow one CheckAndGenerate to run at a time
	// TODO: handle this differently if you ever scale out beyond one instance. :)
	result, err, _ := s.sf.Do(strconv.FormatInt(game.ID, 10), func() (interface{}, error) {
		gr := NewGameRunner(s.db, s.config)
		result, err := gr.CheckAndGenerateTurn(game.ID)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Msg("check and generate new turn")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("check and generate new turn: %w", err))
	}

	// return the game status
	updatedGame, err := dbClient.GetGame(ctx, game.ID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Msg("load game")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("load game: %w", err))
	}

	if result == TurnGenerated {
		s.discordNotifier.SendNewTurnNotification(ctx, game.ID)
	}

	return connect.NewResponse(&craig_starsv1.SubmitTurnResponse{
		Game:   converter.C.ConvertCSGameWithPlayers(updatedGame),
		Player: converter.C.ConvertCSPlayer(player),
	}), nil
}

// UnsubmitTurn mirrors REST unSubmitTurn: mark the player's turn as not submitted and return player.
func (s *playerService) UnsubmitTurn(ctx context.Context, req *connect.Request[craig_starsv1.UnsubmitTurnRequest]) (*connect.Response[craig_starsv1.UnsubmitTurnResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("failed to get player")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get player: %w", err))
	}

	player.SubmittedTurn = false
	if err := dbWriteClient.SubmitPlayerTurn(ctx, player.GameID, player.Num, false); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unsubmit turn"))
	}

	return connect.NewResponse(&craig_starsv1.UnsubmitTurnResponse{
		Player: converter.C.ConvertCSPlayer(player),
	}), nil
}

// UpdatePlayerOrders mirrors REST updatePlayerOrders.
func (s *playerService) UpdatePlayerOrders(ctx context.Context, req *connect.Request[craig_starsv1.UpdatePlayerOrdersRequest]) (*connect.Response[craig_starsv1.UpdatePlayerOrdersResponse], error) {
	dbClient := contextDb(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	player, err := dbClient.GetPlayerForGame(ctx, game.ID, gamePlayer.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("failed to get player")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get player: %w", err))
	}
	player.Race.Spec = cs.ComputeRaceSpec(&player.Race, &game.Rules)

	// validation
	if req.Msg.Orders == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing orders"))
	}
	if req.Msg.Orders.ResearchAmount < 0 || req.Msg.Orders.ResearchAmount > 100 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("research amount must be between 0 and 100"))
	}

	planets, err := dbClient.GetPlanetsForPlayer(ctx, player.GameID, player.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("load player planets"))
	}

	orderer := cs.NewOrderer()
	// Build cs.PlayerOrders from proto without new converter helpers
	orders := converter.C.ConvertPlayerOrders(req.Msg.Orders)
	orders.CargoTransfers = player.CargoTransfers // don't let the player update cargoTransfers outside of a fleet transfer function
	orderer.UpdatePlayerOrders(player, planets, orders, &game.Rules)

	// save the updated fleets back to the database
	if err := s.db.WrapInTransaction(func(c db.Client) error {
		// save the player to the database
		if err := c.UpdatePlayerOrders(ctx, player); err != nil {
			log.Error().Err(err).Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update player")
			return err
		}

		for _, planet := range planets {
			if planet.Dirty {
				// TODO: only update the planet spec? that's all that changes
				if err := c.UpdatePlanetSpec(ctx, planet); err != nil {
					log.Error().Err(err).Int64("ID", player.ID).Msg("updating player planet in database")
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("save player orders"))
	}

	return connect.NewResponse(&craig_starsv1.UpdatePlayerOrdersResponse{
		Player:  converter.C.ConvertCSPlayer(player),
		Planets: converter.C.ConvertCSPlanets(planets),
	}), nil
}

// UpdatePlayerRelations mirrors REST updatePlayerRelations.
func (s *playerService) UpdatePlayerRelations(ctx context.Context, req *connect.Request[craig_starsv1.UpdatePlayerRelationsRequest]) (*connect.Response[craig_starsv1.UpdatePlayerRelationsResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("failed to get player")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get player: %w", err))
	}

	relations := converter.C.ConvertPlayerRelations(req.Msg.Relations)
	if len(relations) != len(player.Relations) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("must include all player relations"))
	}

	player.Relations = relations
	if err := dbWriteClient.UpdatePlayerRelations(ctx, player); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("update player relations"))
	}

	// return only relations, like REST
	return connect.NewResponse(&craig_starsv1.UpdatePlayerRelationsResponse{
		Relations: req.Msg.Relations,
	}), nil
}
