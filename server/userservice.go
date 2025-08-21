package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
)

func NewUserServiceHandler(db DBConnection, discordNotifier *discordNotifier) craig_starsv1connect.UserServiceHandler {
	return &userService{db, discordNotifier}
}

type userService struct {
	db              DBConnection
	discordNotifier *discordNotifier
}

// GetMe implements craig_starsv1connect.UserServiceHandler.
func (s *userService) GetMe(ctx context.Context, req *connect.Request[craig_starsv1.GetMeRequest]) (*connect.Response[craig_starsv1.GetMeResponse], error) {
	user := contextUserSession(ctx)

	return connect.NewResponse(&craig_starsv1.GetMeResponse{
		User: &craig_starsv1.User{
			Id:            user.ID,
			Username:      user.Username,
			Role:          converter.CSUserRoleToUserRole(cs.UserRole(user.Role)),
			DiscordId:     user.DiscordID,
			DiscordAvatar: user.DiscordAvatar,
		},
	}), nil
}

// UpdateUserSettings implements craig_starsv1connect.UserServiceHandler.
func (s *userService) UpdateUserSettings(ctx context.Context, req *connect.Request[craig_starsv1.UpdateUserSettingsRequest]) (*connect.Response[craig_starsv1.UpdateUserSettingsResponse], error) {
	me := contextUserSession(ctx)
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)

	user, err := dbClient.GetUser(ctx, me.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get user"))
	}

	user.UserSettings = converter.C.ConvertUserSettings(req.Msg.UserSettings)
	if err := dbWriteClient.UpdateUserSettings(ctx, user); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to save user"))
	}

	return connect.NewResponse(&craig_starsv1.UpdateUserSettingsResponse{
		User: converter.C.ConvertCSUser(*user),
	}), nil

}

// GetUser returns a single user by ID
func (s *userService) GetUser(ctx context.Context, req *connect.Request[craig_starsv1.GetUserRequest]) (*connect.Response[craig_starsv1.GetUserResponse], error) {
	dbClient := contextDb(ctx)
	me := contextUserSession(ctx)
	if req.Msg.UserId != me.ID && !me.isAdmin() {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("only admins can access other users"))
	}

	user, err := dbClient.GetUser(ctx, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get user"))
	}

	if user == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
	}

	return connect.NewResponse(&craig_starsv1.GetUserResponse{
		User: converter.C.ConvertCSUser(*user),
	}), nil
}

// TestDiscordWebhook implements craig_starsv1connect.UserServiceHandler.
func (s *userService) TestDiscordWebhook(ctx context.Context, req *connect.Request[craig_starsv1.TestDiscordWebhookRequest]) (*connect.Response[craig_starsv1.TestDiscordWebhookResponse], error) {
	dbClient := contextDb(ctx)
	me := contextUserSession(ctx)

	user, err := dbClient.GetUser(ctx, me.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get user"))
	}

	if user.DiscordWebhookURL == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("no discord webhook"))
	}

	if user.DiscordID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("not a discord user"))
	}

	if err := s.discordNotifier.SendTestWebhook(ctx, user.DiscordWebhookURL, user.DiscordID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to send test webhook: %w", err))
	}

	log.Info().Msgf("sending test discord message for %s", user.Username)

	return &connect.Response[craig_starsv1.TestDiscordWebhookResponse]{}, nil
}
