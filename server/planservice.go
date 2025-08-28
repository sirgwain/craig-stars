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

// ------------------------
// BattlePlanService
// ------------------------

func NewBattlePlanServiceHandler(db DBConnection) craig_starsv1connect.BattlePlanServiceHandler {
	return &battlePlanService{db: db}
}

type battlePlanService struct {
	db DBConnection
}

// GetBattlePlan returns a single battle plan for the current player by num.
func (s *battlePlanService) GetBattlePlan(ctx context.Context, req *connect.Request[craig_starsv1.GetBattlePlanRequest]) (*connect.Response[craig_starsv1.GetBattlePlanResponse], error) {
	c := contextDb(ctx)
	user := contextUserSession(ctx)
	game := contextGame(ctx)

	player, err := c.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{UserID: user.ID})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var battlePlan *cs.BattlePlan
	for i := range player.BattlePlans {
		if player.BattlePlans[i].Num == int(req.Msg.Num) {
			battlePlan = &player.BattlePlans[i]
			break
		}
	}
	if battlePlan == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("plan not found"))
	}

	return connect.NewResponse(&craig_starsv1.GetBattlePlanResponse{
		Plan: converter.C.ConvertCSBattlePlan(*battlePlan),
	}), nil
}

// CreateBattlePlan validates and saves a new battle plan for the current player.
func (s *battlePlanService) CreateBattlePlan(ctx context.Context, req *connect.Request[craig_starsv1.CreateBattlePlanRequest]) (*connect.Response[craig_starsv1.CreateBattlePlanResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	if req.Msg.Plan == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing plan"))
	}

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if player == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player not found"))
	}

	plan := converter.C.ConvertBattlePlan(req.Msg.Plan)
	if err := plan.Validate(player); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	plan.Num = player.GetNextBattlePlanNum()
	player.BattlePlans = append(player.BattlePlans, plan)

	if err := dbWriteClient.UpdatePlayerPlans(ctx, player); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.CreateBattlePlanResponse{
		Plan: converter.C.ConvertCSBattlePlan(plan),
	}), nil
}

// UpdateBattlePlan updates an existing battle plan.
func (s *battlePlanService) UpdateBattlePlan(ctx context.Context, req *connect.Request[craig_starsv1.UpdateBattlePlanRequest]) (*connect.Response[craig_starsv1.UpdateBattlePlanResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	if req.Msg.Plan == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing plan"))
	}

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if player == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player not found"))
	}

	plan := converter.C.ConvertBattlePlan(req.Msg.Plan)
	if err := plan.Validate(player); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	for i := range player.BattlePlans {
		if player.BattlePlans[i].Num == plan.Num {
			player.BattlePlans[i] = plan
			break
		}
	}

	if err := dbWriteClient.UpdatePlayerPlans(ctx, player); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&craig_starsv1.UpdateBattlePlanResponse{
		Plan: converter.C.ConvertCSBattlePlan(plan),
	}), nil
}

// DeleteBattlePlan deletes a battle plan and updates fleets using it to default plan.
func (s *battlePlanService) DeleteBattlePlan(ctx context.Context, req *connect.Request[craig_starsv1.DeleteBattlePlanRequest]) (*connect.Response[craig_starsv1.DeleteBattlePlanResponse], error) {
	dbClient := contextDb(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if player == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player not found"))
	}

	num := int(req.Msg.Num)
	var plan *cs.BattlePlan
	for i := range player.BattlePlans {
		if player.BattlePlans[i].Num == num {
			plan = &player.BattlePlans[i]
			break
		}
	}
	if plan == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("plan not found"))
	}
	if plan.Num == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("cannot delete default battle plan"))
	}

	// set fleets with this plan to default, remove plan, save in tx
	playerFleets, err := dbClient.GetFleetsForPlayer(ctx, game.ID, gamePlayer.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	fleetsToUpdate := []*cs.Fleet{}
	for _, f := range playerFleets {
		if f.BattlePlanNum == plan.Num {
			f.BattlePlanNum = 0
			fleetsToUpdate = append(fleetsToUpdate, f)
		}
	}

	plans := make([]cs.BattlePlan, 0, len(player.BattlePlans)-1)
	for _, bp := range player.BattlePlans {
		if bp.Num != plan.Num {
			plans = append(plans, bp)
		}
	}
	player.BattlePlans = plans

	if err := s.db.WrapInTransaction(func(c db.Client) error {
		for _, f := range fleetsToUpdate {
			if err := c.SaveFleet(ctx, f); err != nil {
				slog.Error("update fleet in database", slog.Any("error", err))
				return err
			}
		}
		if err := c.UpdatePlayerPlans(ctx, player); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// reload fleets for response and split into fleets/starbases
	allFleets, err := dbClient.GetFleetsForPlayer(ctx, game.ID, gamePlayer.Num)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	fleets := make([]*cs.Fleet, 0, len(allFleets))
	starbases := make([]*cs.Fleet, 0)
	for _, f := range allFleets {
		if f.Starbase {
			starbases = append(starbases, f)
		} else {
			fleets = append(fleets, f)
		}
	}

	return connect.NewResponse(&craig_starsv1.DeleteBattlePlanResponse{
		Player:    converter.C.ConvertCSPlayer(player),
		Fleets:    converter.C.ConvertCSFleets(fleets),
		Starbases: converter.C.ConvertCSFleets(starbases),
	}), nil
}

// ------------------------
// ProductionPlanService
// ------------------------

func NewProductionPlanServiceHandler() craig_starsv1connect.ProductionPlanServiceHandler {
	return &productionPlanService{}
}

type productionPlanService struct{}

// GetProductionPlan returns a single production plan by num.
func (s *productionPlanService) GetProductionPlan(ctx context.Context, req *connect.Request[craig_starsv1.GetProductionPlanRequest]) (*connect.Response[craig_starsv1.GetProductionPlanResponse], error) {
	c := contextDb(ctx)
	user := contextUserSession(ctx)
	game := contextGame(ctx)

	player, err := c.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{UserID: user.ID})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	var plan *cs.ProductionPlan
	for i := range player.ProductionPlans {
		if player.ProductionPlans[i].Num == int(req.Msg.Num) {
			plan = &player.ProductionPlans[i]
			break
		}
	}
	if plan == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("plan not found"))
	}
	return connect.NewResponse(&craig_starsv1.GetProductionPlanResponse{
		Plan: converter.C.ConvertCSProductionPlan(*plan),
	}), nil
}

// CreateProductionPlan validates and saves a new production plan for the current player.
func (s *productionPlanService) CreateProductionPlan(ctx context.Context, req *connect.Request[craig_starsv1.CreateProductionPlanRequest]) (*connect.Response[craig_starsv1.CreateProductionPlanResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	if req.Msg.Plan == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing plan"))
	}

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if player == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player not found"))
	}

	newPlan := converter.C.ConvertProductionPlan(req.Msg.Plan)
	if err := newPlan.Validate(player); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	newPlan.Num = player.GetNextProductionPlanNum()
	player.ProductionPlans = append(player.ProductionPlans, newPlan)

	if err := dbWriteClient.UpdatePlayerPlans(ctx, player); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.CreateProductionPlanResponse{
		Plan: converter.C.ConvertCSProductionPlan(newPlan),
	}), nil
}

// UpdateProductionPlan updates an existing production plan.
func (s *productionPlanService) UpdateProductionPlan(ctx context.Context, req *connect.Request[craig_starsv1.UpdateProductionPlanRequest]) (*connect.Response[craig_starsv1.UpdateProductionPlanResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	if req.Msg.Plan == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing plan"))
	}

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if player == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player not found"))
	}

	plan := converter.C.ConvertProductionPlan(req.Msg.Plan)
	if err := plan.Validate(player); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	for i := range player.ProductionPlans {
		if player.ProductionPlans[i].Num == plan.Num {
			player.ProductionPlans[i] = plan
			break
		}
	}

	if err := dbWriteClient.UpdatePlayerPlans(ctx, player); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&craig_starsv1.UpdateProductionPlanResponse{
		Plan: converter.C.ConvertCSProductionPlan(plan),
	}), nil
}

// DeleteProductionPlan deletes a production plan.
func (s *productionPlanService) DeleteProductionPlan(ctx context.Context, req *connect.Request[craig_starsv1.DeleteProductionPlanRequest]) (*connect.Response[craig_starsv1.DeleteProductionPlanResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if player == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player not found"))
	}

	num := int(req.Msg.Num)
	var plan *cs.ProductionPlan
	for i := range player.ProductionPlans {
		if player.ProductionPlans[i].Num == num {
			plan = &player.ProductionPlans[i]
			break
		}
	}
	if plan == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("plan not found"))
	}
	if plan.Num == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("cannot delete default production plan"))
	}
	newPlans := make([]cs.ProductionPlan, 0, len(player.ProductionPlans)-1)
	for _, p := range player.ProductionPlans {
		if p.Num != plan.Num {
			newPlans = append(newPlans, p)
		}
	}
	player.ProductionPlans = newPlans
	if err := dbWriteClient.UpdatePlayerPlans(ctx, player); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	slog.Info("deleted ProductionPlan", slog.String("plan", plan.Name), slog.Int64("GameID", game.ID), slog.Int("PlayerNum", gamePlayer.Num), slog.Int("Num", plan.Num))
	return connect.NewResponse(&craig_starsv1.DeleteProductionPlanResponse{}), nil
}

// ------------------------
// TransportPlanService
// ------------------------

func NewTransportPlanServiceHandler() craig_starsv1connect.TransportPlanServiceHandler {
	return &transportPlanService{}
}

type transportPlanService struct{}

// GetTransportPlan returns a single transport plan by num.
func (s *transportPlanService) GetTransportPlan(ctx context.Context, req *connect.Request[craig_starsv1.GetTransportPlanRequest]) (*connect.Response[craig_starsv1.GetTransportPlanResponse], error) {
	c := contextDb(ctx)
	user := contextUserSession(ctx)
	game := contextGame(ctx)

	player, err := c.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{UserID: user.ID})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	var plan *cs.TransportPlan
	for i := range player.TransportPlans {
		if player.TransportPlans[i].Num == int(req.Msg.Num) {
			plan = &player.TransportPlans[i]
			break
		}
	}
	if plan == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("plan not found"))
	}
	return connect.NewResponse(&craig_starsv1.GetTransportPlanResponse{
		Plan: converter.C.ConvertCSTransportPlan(*plan),
	}), nil
}

// CreateTransportPlan validates and saves a new transport plan for the current player.
func (s *transportPlanService) CreateTransportPlan(ctx context.Context, req *connect.Request[craig_starsv1.CreateTransportPlanRequest]) (*connect.Response[craig_starsv1.CreateTransportPlanResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	if req.Msg.Plan == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing plan"))
	}

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if player == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player not found"))
	}

	newPlan := converter.C.ConvertTransportPlan(req.Msg.Plan)
	if err := newPlan.Validate(player); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	newPlan.Num = player.GetNextTransportPlanNum()
	player.TransportPlans = append(player.TransportPlans, newPlan)

	if err := dbWriteClient.UpdatePlayerPlans(ctx, player); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.CreateTransportPlanResponse{
		Plan: converter.C.ConvertCSTransportPlan(newPlan),
	}), nil
}

// UpdateTransportPlan updates an existing transport plan.
func (s *transportPlanService) UpdateTransportPlan(ctx context.Context, req *connect.Request[craig_starsv1.UpdateTransportPlanRequest]) (*connect.Response[craig_starsv1.UpdateTransportPlanResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	if req.Msg.Plan == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing plan"))
	}

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if player == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player not found"))
	}

	plan := converter.C.ConvertTransportPlan(req.Msg.Plan)
	if err := plan.Validate(player); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	for i := range player.TransportPlans {
		if player.TransportPlans[i].Num == plan.Num {
			player.TransportPlans[i] = plan
			break
		}
	}

	if err := dbWriteClient.UpdatePlayerPlans(ctx, player); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&craig_starsv1.UpdateTransportPlanResponse{
		Plan: converter.C.ConvertCSTransportPlan(plan),
	}), nil
}

// DeleteTransportPlan deletes a transport plan.
func (s *transportPlanService) DeleteTransportPlan(ctx context.Context, req *connect.Request[craig_starsv1.DeleteTransportPlanRequest]) (*connect.Response[craig_starsv1.DeleteTransportPlanResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	game := contextGame(ctx)
	gamePlayer := contextGamePlayer(ctx)

	player, err := dbClient.GetLightPlayerForGame(ctx, game.ID, db.GetPlayerParams{PlayerNum: gamePlayer.Num})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if player == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player not found"))
	}

	num := int(req.Msg.Num)
	var plan *cs.TransportPlan
	for i := range player.TransportPlans {
		if player.TransportPlans[i].Num == num {
			plan = &player.TransportPlans[i]
			break
		}
	}
	if plan == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("plan not found"))
	}
	if plan.Num == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("cannot delete default transport plan"))
	}
	newPlans := make([]cs.TransportPlan, 0, len(player.TransportPlans)-1)
	for _, p := range player.TransportPlans {
		if p.Num != plan.Num {
			newPlans = append(newPlans, p)
		}
	}
	player.TransportPlans = newPlans
	if err := dbWriteClient.UpdatePlayerPlans(ctx, player); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	slog.Info("deleted TransportPlan", slog.String("plan", plan.Name), slog.Int64("GameID", game.ID), slog.Int("PlayerNum", gamePlayer.Num), slog.Int("Num", plan.Num))
	return connect.NewResponse(&craig_starsv1.DeleteTransportPlanResponse{}), nil
}
