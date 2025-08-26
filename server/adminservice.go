package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
)

// NewAdminServiceHandler constructs the AdminService handler.
func NewAdminServiceHandler(db DBConnection) craig_starsv1connect.AdminServiceHandler {
	return &adminService{db: db}
}

type adminService struct {
	db DBConnection
}

// GetAllGames implements AdminService.GetAllGames:
// - requires admin role
// - returns all games with players
func (s *adminService) GetAllGames(ctx context.Context, req *connect.Request[craig_starsv1.GetAllGamesRequest]) (*connect.Response[craig_starsv1.GetAllGamesResponse], error) {
	user := contextUserSession(ctx)
	if !user.isAdmin() {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("only admins can view all games"))
	}

	c := contextDb(ctx)
	games, err := c.GetGamesWithPlayers(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("get games from database: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.GetAllGamesResponse{
		Games: converter.C.ConvertCSGamesWithPlayers(games),
	}), nil
}

// GetUserGames implements AdminService.GetUserGames:
// - requires admin role
// - returns all games for a specified user
func (s *adminService) GetUserGames(ctx context.Context, req *connect.Request[craig_starsv1.GetUserGamesRequest]) (*connect.Response[craig_starsv1.GetUserGamesResponse], error) {
	user := contextUserSession(ctx)
	if !user.isAdmin() {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("only admins can view user games"))
	}

	c := contextDb(ctx)
	games, err := c.GetGamesForUser(ctx, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("get games from database: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.GetUserGamesResponse{
		Games: converter.C.ConvertCSGamesWithPlayers(games),
	}), nil
}

// GetUsers returns all users for the authenticated user
func (s *adminService) GetUsers(ctx context.Context, req *connect.Request[craig_starsv1.GetUsersRequest]) (*connect.Response[craig_starsv1.GetUsersResponse], error) {
	me := contextUserSession(ctx)
	if !me.isAdmin() {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("access denied"))
	}
	dbClient := contextDb(ctx)

	users, err := dbClient.GetUsers(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get users"))
	}

	// we don't actually include a password in the proto user, but just as a precaution...
	for i := range users {
		users[i].Password = ""
	}

	return connect.NewResponse(&craig_starsv1.GetUsersResponse{
		Users: converter.C.ConvertCSUsers(users),
	}), nil
}

// ConvertGuestUser implements AdminService.ConvertGuestUser:
// - requires admin role
// - converts a guest user's data to a real user, mirroring REST logic in server/admin.go
func (s *adminService) ConvertGuestUser(ctx context.Context, req *connect.Request[craig_starsv1.ConvertGuestUserRequest]) (*connect.Response[craig_starsv1.ConvertGuestUserResponse], error) {
	dbClient := contextDb(ctx)
	admin := contextUserSession(ctx)
	if !admin.isAdmin() {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("only admins can convert guest users"))
	}

	guestID := req.Msg.GuestUserId
	targetUserID := req.Msg.UserId

	guestUser, err := dbClient.GetUser(ctx, guestID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", guestID).Msg("load guest user")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("load guest user"))
	}
	if guestUser == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("guest user not found"))
	}

	// must be a guest user
	if !guestUser.IsGuest() {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%s is not a guest user", guestUser.Username))
	}

	targetUser, err := dbClient.GetUser(ctx, targetUserID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", targetUserID).Msg("load user")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("load user"))
	}
	if targetUser == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
	}

	// target user cannot be a guest
	if targetUser.IsGuest() {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%s is a guest user", targetUser.Username))
	}

	// load players belonging to guest
	players, err := dbClient.GetPlayersForUser(ctx, guestUser.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", targetUserID).Msg("load players")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("load players"))
	}

	// load guest user games as player
	playerGames, err := dbClient.GetGamesForUser(ctx, guestUser.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", targetUserID).Msg("load player games")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("load player games"))
	}

	// ensure target user is not already playing in the same games as guest
	for _, game := range playerGames {
		for _, p := range game.Players {
			if p.UserID == targetUser.ID {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%s is already playing in game %d (%s)", targetUser.Username, game.ID, game.Name))
			}
		}
	}

	// load races belonging to guest
	races, err := dbClient.GetRacesForUser(ctx, guestUser.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", targetUserID).Msg("load races")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("load races"))
	}

	// load games hosted by guest
	hostedGames, err := dbClient.GetGamesForHost(ctx, guestUser.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", targetUserID).Msg("load games for host")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("load games for host"))
	}

	// perform updates in a single transaction
	if err := s.db.WrapInTransaction(func(c db.Client) error {
		// move players
		for _, player := range players {
			player.UserID = targetUser.ID
			if err := c.UpdatePlayerUserID(ctx, player); err != nil {
				return fmt.Errorf("update Player UserID: %w", err)
			}
		}

		// move races
		for _, race := range races {
			race.UserID = targetUser.ID
			if err := c.SaveRace(ctx, &race); err != nil {
				return fmt.Errorf("update Race UserID: %w", err)
			}
		}

		// move hosted games
		for _, game := range hostedGames {
			game.HostID = targetUser.ID
			if err := c.UpdateGameHost(ctx, game.ID, game.HostID); err != nil {
				return fmt.Errorf("update Game HostID: %w", err)
			}
		}

		// delete guest user
		if err := c.DeleteUser(ctx, guestUser.ID); err != nil {
			return fmt.Errorf("delete guest user: %w", err)
		}

		log.Info().Msgf("converted guest %s to %s (%d players, %d races, %d games)", guestUser.Username, targetUser.Username, len(players), len(races), len(hostedGames))

		return nil
	}); err != nil {
		log.Error().Err(err).Msg("update user objects in database")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	log.Info().
		Int64("GuestUserID", guestUser.ID).
		Int64("UserID", targetUser.ID).
		Msgf("moved guest %s games and races to %s, deleted %s", guestUser.Username, targetUser.Username, guestUser.Username)

	return connect.NewResponse(&craig_starsv1.ConvertGuestUserResponse{}), nil
}
