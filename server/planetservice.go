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

func NewPlanetServiceHandler(db DBConnection) craig_starsv1connect.PlanetServiceHandler {
	return &planetService{db}
}

type planetService struct {
	db DBConnection
}

func (s *planetService) GetPlanet(ctx context.Context, req *connect.Request[craig_starsv1.GetPlanetRequest]) (*connect.Response[craig_starsv1.GetPlanetResponse], error) {
	dbClient := contextDb(ctx)
	player := contextGamePlayer(ctx)

	planet, err := dbClient.GetPlanetByNum(ctx, req.Msg.GameId, int(req.Msg.PlanetNum))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get planet from db"))
	}

	if planet.PlayerNum != player.Num {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("player doesn't own planet"))
	}

	return connect.NewResponse(&craig_starsv1.GetPlanetResponse{
		Planet: converter.C.ConvertCSPlanet(planet),
	}), nil
}

func (s *planetService) UpdatePlanetOrders(ctx context.Context, req *connect.Request[craig_starsv1.UpdatePlanetOrdersRequest]) (*connect.Response[craig_starsv1.UpdatePlanetOrdersResponse], error) {
	dbClient := contextDb(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	planet, err := dbClient.GetPlanetByNum(ctx, req.Msg.GameId, int(req.Msg.PlanetNum))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get planet from db"))
	}

	if planet.PlayerNum != gamePlayer.Num {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("player doesn't own planet"))
	}

	// load the full player to update planet production estimates
	player, err := dbClient.GetPlayerForGame(ctx, game.ID, gamePlayer.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	player.Race.Spec = cs.ComputeRaceSpec(&player.Race, &game.Rules)

	// load all a player's planets so we can recompute research estimates
	planets, err := dbClient.GetPlanetsForPlayer(ctx, game.ID, player.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	orders := converter.C.ConvertPlanetOrders(req.Msg.PlanetOrders)
	orderer := cs.NewOrderer()
	if err := orderer.UpdatePlanetOrders(&game.Rules, player, planet, *orders, planets); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Str("Planet", planet.Name).Msg("update planet orders")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// update this planet and the player's spec in the database
	if err := s.db.WrapInTransaction(func(c db.Client) error {
		if err := c.SavePlanet(ctx, planet); err != nil {
			log.Error().Err(err).Int64("ID", planet.ID).Msg("update planet in database")
			return err
		}

		// update the player spec as well because changes in planet orders impact resources
		// available for research
		if err := c.UpdatePlayerSpec(ctx, player); err != nil {
			log.Error().Err(err).Int64("ID", planet.ID).Msg("update player spec in database")
			return err
		}
		return nil
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.UpdatePlanetOrdersResponse{
		Planet: converter.C.ConvertCSPlanet(planet),
		Player: converter.C.ConvertCSPlayer(player),
	}), nil
}
