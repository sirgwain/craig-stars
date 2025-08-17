//go:build !wasi && !wasm

package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
)

func NewRaceServiceHandler(db DBConnection) craig_starsv1connect.RaceServiceHandler {
	return &raceService{db}
}

type raceService struct {
	db DBConnection
}

// GetRaces returns all races for the authenticated user
func (s *raceService) GetRaces(ctx context.Context, req *connect.Request[craig_starsv1.GetRacesRequest]) (*connect.Response[craig_starsv1.GetRacesResponse], error) {
	dbClient := contextDb(ctx)
	user := contextUserSession(ctx)

	races, err := dbClient.GetRacesForUser(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("get races from database")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get races from database"))
	}

	return connect.NewResponse(&craig_starsv1.GetRacesResponse{
		Races: converter.C.ConvertCSRaces(races),
	}), nil
}

// GetRace returns a single race by ID
func (s *raceService) GetRace(ctx context.Context, req *connect.Request[craig_starsv1.GetRaceRequest]) (*connect.Response[craig_starsv1.GetRaceResponse], error) {
	dbClient := contextDb(ctx)
	user := contextUserSession(ctx)

	race, err := dbClient.GetRace(ctx, req.Msg.RaceId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get race from database"))
	}

	if race == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("race not found"))
	}

	if race.UserID != user.ID {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("user does not own race"))
	}

	return connect.NewResponse(&craig_starsv1.GetRaceResponse{
		Race: converter.C.ConvertCSRace(*race),
	}), nil
}

// CreateRace creates a new race for the authenticated user
func (s *raceService) CreateRace(ctx context.Context, req *connect.Request[craig_starsv1.CreateRaceRequest]) (*connect.Response[craig_starsv1.CreateRaceResponse], error) {
	dbWriteClient := contextDbWrite(ctx)
	user := contextUserSession(ctx)

	// Convert proto race to CS race
	race := converter.C.ConvertRaceP(req.Msg.Race)
	race.UserID = user.ID

	if err := dbWriteClient.SaveRace(ctx, race); err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("create race")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewResponse(&craig_starsv1.CreateRaceResponse{
		Race: converter.C.ConvertCSRace(*race),
	}), nil
}

// UpdateRace updates an existing race
func (s *raceService) UpdateRace(ctx context.Context, req *connect.Request[craig_starsv1.UpdateRaceRequest]) (*connect.Response[craig_starsv1.UpdateRaceResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	user := contextUserSession(ctx)

	// Convert proto race to CS race
	race := converter.C.ConvertRaceP(req.Msg.Race)

	// Load existing race for validation
	existingRace, err := dbClient.GetRace(ctx, race.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get existing race from database"))
	}

	if existingRace == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("race not found"))
	}

	if existingRace.UserID != user.ID {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("user does not own race"))
	}

	if err := dbWriteClient.SaveRace(ctx, race); err != nil {
		log.Error().Err(err).Int64("ID", race.ID).Msg("update race in database")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.UpdateRaceResponse{
		Race: converter.C.ConvertCSRace(*race),
	}), nil
}

// DeleteRace deletes a race by ID
func (s *raceService) DeleteRace(ctx context.Context, req *connect.Request[craig_starsv1.DeleteRaceRequest]) (*connect.Response[craig_starsv1.DeleteRaceResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	user := contextUserSession(ctx)

	race, err := dbClient.GetRace(ctx, req.Msg.RaceId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get race from database"))
	}

	if race == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("race not found"))
	}

	if race.UserID != user.ID {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("user does not own race"))
	}

	if err := dbWriteClient.DeleteRace(ctx, race.ID); err != nil {
		log.Error().Err(err).Int64("ID", race.ID).Msg("delete race from database")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&craig_starsv1.DeleteRaceResponse{}), nil
}
