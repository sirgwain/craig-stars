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

func NewMinefieldServiceHandler(db DBConnection) craig_starsv1connect.MinefieldServiceHandler {
	return &minefieldService{db}
}

type minefieldService struct {
	db DBConnection
}

func (s *minefieldService) GetMinefield(ctx context.Context, req *connect.Request[craig_starsv1.GetMinefieldRequest]) (*connect.Response[craig_starsv1.GetMinefieldResponse], error) {
	dbClient := contextDb(ctx)
	gamePlayer := contextGamePlayer(ctx)

	minefield, err := dbClient.GetMinefieldByNum(ctx, req.Msg.GameId, gamePlayer.Num, int(req.Msg.MinefieldNum))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get minefield from db"))
	}

	if minefield.PlayerNum != gamePlayer.Num {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("player doesn't own minefield"))
	}

	return connect.NewResponse(&craig_starsv1.GetMinefieldResponse{
		Minefield: converter.C.ConvertCSMinefield(minefield),
	}), nil
}

func (s *minefieldService) UpdateMinefieldOrders(ctx context.Context, req *connect.Request[craig_starsv1.UpdateMinefieldOrdersRequest]) (*connect.Response[craig_starsv1.UpdateMinefieldOrdersResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	minefield, err := dbClient.GetMinefieldByNum(ctx, req.Msg.GameId, gamePlayer.Num, int(req.Msg.MinefieldNum))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get minefield from db"))
	}

	if minefield.PlayerNum != gamePlayer.Num {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("player doesn't own minefield"))
	}

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get player"))
	}
	player.Race.Spec = cs.ComputeRaceSpec(&player.Race, &game.Rules)

	orders := converter.C.ConvertMinefieldOrders(req.Msg.MinefieldOrders)
	orderer := cs.NewOrderer()
	if err := orderer.UpdateMinefieldOrders(player, minefield, *orders); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update minefield orders"))
	}

	if err := dbWriteClient.SaveMinefield(ctx, minefield); err != nil {
		log.Error().Err(err).Int64("ID", minefield.ID).Msg("update minefield in database")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.UpdateMinefieldOrdersResponse{
		Minefield: converter.C.ConvertCSMinefield(minefield),
	}), nil
}
