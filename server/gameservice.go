package server

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
	"golang.org/x/sync/singleflight"
)

func NewGameServiceHandler(db DBConnection, config config.Config, discordNotifier *discordNotifier) craig_starsv1connect.GameServiceHandler {
	return &gameService{db, config, discordNotifier, singleflight.Group{}}
}

type gameService struct {
	db              DBConnection
	config          config.Config
	discordNotifier *discordNotifier
	sf              singleflight.Group
}

func (s *gameService) GetGame(ctx context.Context, req *connect.Request[craig_starsv1.GetGameRequest]) (*connect.Response[craig_starsv1.GetGameResponse], error) {
	game := contextGame(ctx)

	return connect.NewResponse(&craig_starsv1.GetGameResponse{
		Game: converter.C.ConvertCSGameWithPlayers(game),
	}), nil
}

// CreateGame implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) CreateGame(ctx context.Context, req *connect.Request[craig_starsv1.CreateGameRequest]) (*connect.Response[craig_starsv1.CreateGameResponse], error) {
	user := contextUserSession(ctx)
	dbClient := contextDb(ctx)

	// Convert proto settings to CS settings
	settings := converter.C.ConvertGameSettings(req.Msg.Settings)

	// make sure guests don't create multiplayer games
	if user.isGuest() && !settings.IsSinglePlayer() {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("guests cannot host multiplayer games"))
	}

	gr := NewGameRunner(s.db, s.config)
	game, err := gr.HostGame(user.ID, &settings)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to host game: %w", err))
	}

	gameWithPlayers, err := dbClient.GetGame(ctx, game.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load new game: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.CreateGameResponse{
		Game: converter.C.ConvertCSGameWithPlayers(gameWithPlayers),
	}), nil
}

// DeleteGame implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) DeleteGame(ctx context.Context, req *connect.Request[craig_starsv1.DeleteGameRequest]) (*connect.Response[craig_starsv1.DeleteGameResponse], error) {
	user := contextUserSession(ctx)
	game := contextGame(ctx)

	if game.HostID != user.ID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("only host can delete game"))
	}

	if err := s.db.WrapInTransaction(func(c db.Client) error {
		if err := c.DeleteGame(ctx, game.ID); err != nil {
			return fmt.Errorf("delete game from database: %w", err)
		}

		// delete any guest users
		if err := c.DeleteGameGuestUsers(ctx, game.ID); err != nil {
			return fmt.Errorf("delete game guest users from database: %w", err)
		}

		return nil
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to delete game: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.DeleteGameResponse{}), nil
}

// GetGameByInviteHash implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) GetGameByInviteHash(ctx context.Context, req *connect.Request[craig_starsv1.GetGameByInviteHashRequest]) (*connect.Response[craig_starsv1.GetGameByInviteHashResponse], error) {
	c := contextDb(ctx)

	game, err := c.GetGameByHash(ctx, req.Msg.Hash)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get game: %w", err))
	}
	return connect.NewResponse(&craig_starsv1.GetGameByInviteHashResponse{
		Game: converter.C.ConvertCSGameWithPlayers(game),
	}), nil
}

// GetGames implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) GetGames(ctx context.Context, req *connect.Request[craig_starsv1.GetGamesRequest]) (*connect.Response[craig_starsv1.GetGamesResponse], error) {
	c := contextDb(ctx)
	user := contextUserSession(ctx)

	var err error
	var games []cs.GameWithPlayers
	if req.Msg.Open {
		games, err = c.GetOpenGames(ctx)
	} else {
		games, err = c.GetGamesForUser(ctx, user.ID)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get games: %w", err))
	}
	return connect.NewResponse(&craig_starsv1.GetGamesResponse{
		Games: converter.C.ConvertCSGamesWithPlayers(games),
	}), nil
}

// UpdateGame implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) UpdateGame(ctx context.Context, req *connect.Request[craig_starsv1.UpdateGameRequest]) (*connect.Response[craig_starsv1.UpdateGameResponse], error) {
	dbClient := contextDb(ctx)
	dbWriteClient := contextDbWrite(ctx)
	user := contextUserSession(ctx)
	game := contextGame(ctx)

	if game.State != cs.GameStateSetup {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("game cannot be updated after setup"))
	}

	if game.HostID != user.ID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("user is not the host"))
	}

	// Update game settings from proto request
	game.Name = req.Msg.Settings.Name
	game.Public = req.Msg.Settings.Public
	game.Size = converter.SizeToCSSize(req.Msg.Settings.Size)
	game.Density = converter.DensityToCSDensity(req.Msg.Settings.Density)
	game.PlayerPositions = converter.PlayerPositionsToCSPlayerPositions(req.Msg.Settings.PlayerPositions)
	game.RandomEvents = req.Msg.Settings.RandomEvents
	game.ComputerPlayersFormAlliances = req.Msg.Settings.ComputerPlayersFormAlliances
	game.PublicPlayerScores = req.Msg.Settings.PublicPlayerScores
	game.MaxMinerals = req.Msg.Settings.MaxMinerals
	game.StartMode = converter.GameStartModeToCSGameStartMode(req.Msg.Settings.StartMode)
	game.QuickStartTurns = int(req.Msg.Settings.QuickStartTurns)
	game.VictoryConditions = converter.C.ConvertVictoryConditions(req.Msg.Settings.VictoryConditions)

	if err := dbWriteClient.SaveGame(ctx, &game.Game); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update game: %w", err))
	}

	// Reload the game to get the updated version with players
	updatedGame, err := dbClient.GetGame(ctx, game.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to reload game: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.UpdateGameResponse{
		Game: converter.C.ConvertCSGameWithPlayers(updatedGame),
	}), nil
}

// JoinGame implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) JoinGame(ctx context.Context, req *connect.Request[craig_starsv1.JoinGameRequest]) (*connect.Response[craig_starsv1.JoinGameResponse], error) {
	user := contextUserSession(ctx)
	game := contextGame(ctx)
	gr := NewGameRunner(s.db, s.config)

	if game.State != cs.GameStateSetup {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("cannot join game after setup"))
	}

	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name cannot be empty"))
	}

	race := converter.C.ConvertRaceP(req.Msg.Race)
	race.DBObject = cs.DBObject{}
	if err := gr.JoinGame(game.ID, user.ID, req.Msg.Name, *race); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to join game: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.JoinGameResponse{}), nil
}

// KickPlayer implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) KickPlayer(ctx context.Context, req *connect.Request[craig_starsv1.KickPlayerRequest]) (*connect.Response[craig_starsv1.KickPlayerResponse], error) {
	c := contextDb(ctx)
	user := contextUserSession(ctx)
	game := contextGame(ctx)
	gr := NewGameRunner(s.db, s.config)

	if game.State != cs.GameStateSetup {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("cannot kick player after setup"))
	}

	if user.ID != game.HostID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("user is not host"))
	}

	if err := gr.KickPlayer(game.ID, int(req.Msg.PlayerNum)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to kick player: %w", err))
	}

	// Reload the game for the response
	updatedGame, err := c.GetGame(ctx, game.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to reload game: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.KickPlayerResponse{
		Game: converter.C.ConvertCSGameWithPlayers(updatedGame),
	}), nil
}

// LeaveGame implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) LeaveGame(ctx context.Context, req *connect.Request[craig_starsv1.LeaveGameRequest]) (*connect.Response[craig_starsv1.LeaveGameResponse], error) {
	user := contextUserSession(ctx)
	game := contextGame(ctx)
	gr := NewGameRunner(s.db, s.config)

	if game.State != cs.GameStateSetup {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("cannot leave game after setup"))
	}

	if err := gr.LeaveGame(game.ID, user.ID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to leave game: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.LeaveGameResponse{}), nil
}

// StartGame implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) StartGame(ctx context.Context, req *connect.Request[craig_starsv1.StartGameRequest]) (*connect.Response[craig_starsv1.StartGameResponse], error) {
	user := contextUserSession(ctx)
	game := contextGame(ctx)
	gr := NewGameRunner(s.db, s.config)

	// validate
	if user.ID != game.HostID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("user is not host"))
	}

	if user.isGuest() && !game.IsSinglePlayer() {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("guests cannot start multiplayer games"))
	}

	if err := gr.StartGame(&game.Game); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to start game: %w", err))
	}

	s.discordNotifier.SendNewTurnNotification(game.ID)

	return connect.NewResponse(&craig_starsv1.StartGameResponse{}), nil

}

// AddPlayer implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) AddPlayer(ctx context.Context, req *connect.Request[craig_starsv1.AddPlayerRequest]) (*connect.Response[craig_starsv1.AddPlayerResponse], error) {
	c := contextDb(ctx)
	user := contextUserSession(ctx)
	game := contextGame(ctx)
	gr := NewGameRunner(s.db, s.config)

	if game.State != cs.GameStateSetup {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("cannot add player after setup"))
	}

	if user.ID != game.HostID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("user is not host"))
	}

	var err error
	switch req.Msg.PlayerType {
	case craig_starsv1.PlayerType_PLAYER_TYPE_AI:
		_, err = gr.AddAIPlayer(game)
	case craig_starsv1.PlayerType_PLAYER_TYPE_GUEST:
		if user.isGuest() {
			return nil, connect.NewError(connect.CodePermissionDenied, errors.New("guests cannot add guest players"))
		}
		_, err = gr.AddGuestPlayer(game)
	case craig_starsv1.PlayerType_PLAYER_TYPE_OPEN:
		if user.isGuest() {
			return nil, connect.NewError(connect.CodePermissionDenied, errors.New("guests cannot add open slots"))
		}
		_, err = gr.AddOpenPlayerSlot(game)
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid player type"))
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to add player: %w", err))
	}

	// Reload the game for the response
	updatedGame, err := c.GetGame(ctx, game.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to reload game: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.AddPlayerResponse{
		Game: converter.C.ConvertCSGameWithPlayers(updatedGame),
	}), nil
}

// ArchiveGame implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) ArchiveGame(ctx context.Context, req *connect.Request[craig_starsv1.ArchiveGameRequest]) (*connect.Response[craig_starsv1.ArchiveGameResponse], error) {
	user := contextUserSession(ctx)
	game := contextGame(ctx)
	player := contextGamePlayer(ctx)

	if err := s.db.WrapInTransaction(func(c db.Client) error {
		// archive the whole game if the host requests it
		if user.ID == game.HostID {
			game.Archived = true
			return c.SaveGame(ctx, &game.Game)
		} else {
			return c.ArchivePlayer(ctx, game.ID, player.Num, true)
		}
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to archive game: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.ArchiveGameResponse{
		Game: converter.C.ConvertCSGameWithPlayers(game),
	}), nil
}

// DeletePlayerSlot implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) DeletePlayerSlot(ctx context.Context, req *connect.Request[craig_starsv1.DeletePlayerSlotRequest]) (*connect.Response[craig_starsv1.DeletePlayerSlotResponse], error) {
	c := contextDb(ctx)
	user := contextUserSession(ctx)
	game := contextGame(ctx)
	gr := NewGameRunner(s.db, s.config)

	if game.State != cs.GameStateSetup {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("cannot delete player slot after setup"))
	}

	if user.ID != game.HostID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("user is not host"))
	}

	if err := gr.DeletePlayerSlot(game.ID, int(req.Msg.PlayerNum)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to delete player slot: %w", err))
	}

	// Reload the game for the response
	updatedGame, err := c.GetGame(ctx, game.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to reload game: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.DeletePlayerSlotResponse{
		Game: converter.C.ConvertCSGameWithPlayers(updatedGame),
	}), nil
}

// GetGuestUser implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) GetGuestUser(ctx context.Context, req *connect.Request[craig_starsv1.GetGuestUserRequest]) (*connect.Response[craig_starsv1.GetGuestUserResponse], error) {
	c := contextDb(ctx)
	user := contextUserSession(ctx)
	game := contextGame(ctx)

	if user.ID != game.HostID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("only host can load guest users"))
	}

	guest, err := c.GetGuestUserForGame(ctx, game.ID, int(req.Msg.PlayerNum))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get guest user: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.GetGuestUserResponse{
		User: &craig_starsv1.GuestUser{
			Id:        guest.ID,
			Username:  guest.Username,
			Hash:      guest.Password,
			GameId:    guest.GameID,
			PlayerNum: int32(guest.PlayerNum),
		},
	}), nil
}

// UnarchiveGame implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) UnarchiveGame(ctx context.Context, req *connect.Request[craig_starsv1.UnarchiveGameRequest]) (*connect.Response[craig_starsv1.UnarchiveGameResponse], error) {
	user := contextUserSession(ctx)
	game := contextGame(ctx)
	player := contextGamePlayer(ctx)

	if err := s.db.WrapInTransaction(func(c db.Client) error {
		// unarchive the whole game if the host requests it
		if user.ID == game.HostID {
			game.Archived = false
			return c.SaveGame(ctx, &game.Game)
		} else {
			return c.ArchivePlayer(ctx, game.ID, player.Num, false)
		}
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to unarchive game: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.UnarchiveGameResponse{
		Game: converter.C.ConvertCSGameWithPlayers(game),
	}), nil
}

// ForceGenerateTurn implements craig_starsv1connect.GameServiceHandler.
func (s *gameService) ForceGenerateTurn(ctx context.Context, req *connect.Request[craig_starsv1.ForceGenerateTurnRequest]) (*connect.Response[craig_starsv1.ForceGenerateTurnResponse], error) {
	dbClient := contextDb(ctx)
	user := contextUserSession(ctx)
	game := contextGame(ctx)

	// validate
	if user.ID != game.HostID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("user is not host"))
	}

	// only allow one GenerateTurn to run at a time for a game
	// TODO: handle this differently if you ever scale out beyond one instance. :)
	result, err, _ := s.sf.Do(strconv.FormatInt(game.ID, 10), func() (interface{}, error) {
		gr := NewGameRunner(s.db, s.config)
		result, err := gr.GenerateTurn(game.ID)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate turn: %w", err))
	}

	// return the game status
	updatedGame, err := dbClient.GetGame(ctx, game.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load game: %w", err))
	}

	// if a new turn was generated, send discord notification
	if result == TurnGenerated {
		s.discordNotifier.SendNewTurnNotification(game.ID)
	}

	// Convert and return the updated game
	return connect.NewResponse(&craig_starsv1.ForceGenerateTurnResponse{
		Game: converter.C.ConvertCSGameWithPlayers(updatedGame),
	}), nil
}
