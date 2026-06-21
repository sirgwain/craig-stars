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
	"github.com/sirgwain/craig-stars/test/testgames"
)

func NewTestServiceHandler(db DBConnection) craig_starsv1connect.TestServiceHandler {
	return &testService{db}
}

type testService struct {
	db DBConnection
}

// CreateTestGame creates a test game by name from the predefined test games
func (s *testService) CreateTestGame(ctx context.Context, req *connect.Request[craig_starsv1.CreateTestGameRequest]) (*connect.Response[craig_starsv1.CreateTestGameResponse], error) {
	user := contextUserSession(ctx)

	// Find the test game by test game name
	var testGame *cs.TestScenario
	for _, tg := range testgames.TestGames {
		if tg.Name == req.Msg.TestGameName {
			testGame = &tg
			break
		}
	}

	if testGame == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("test game '%s' not found", req.Msg.TestGameName))
	}

	var game *cs.GameWithPlayers
	if err := s.db.WrapInTransaction(func(c db.Client) error {
		fullGame := cs.BuildScenario(*testGame)

		// Set the custom game name and host user ID
		fullGame.Name = req.Msg.GameName
		fullGame.HostID = user.ID

		// Save the game
		err := c.SaveGame(ctx, fullGame.Game)
		if err != nil {
			return err
		}

		// Save players and update game IDs
		for _, player := range fullGame.Players {
			player.GameID = fullGame.ID

			if err := c.SavePlayer(ctx, player); err != nil {
				return err
			}
			for _, design := range player.Designs {
				design.GameID = fullGame.ID
			}
		}

		// Save to db
		if err := c.UpdateFullGame(ctx, fullGame); err != nil {
			return err
		}

		// // Convert to GameWithPlayers for response
		game = &cs.GameWithPlayers{
			Game:    *fullGame.Game,
			Players: make([]cs.GamePlayer, 0, len(fullGame.Players)),
		}
		for _, player := range fullGame.Players {
			game.Players = append(game.Players, cs.GamePlayer{
				ID:            player.ID,
				UpdatedAt:     &player.UpdatedAt,
				UserID:        player.UserID,
				Name:          player.Name,
				Num:           player.Num,
				Ready:         player.Ready,
				AIControlled:  player.AIControlled,
				AIDifficulty:  player.AIDifficulty,
				Guest:         player.Guest,
				SubmittedTurn: player.SubmittedTurn,
				Color:         player.Color,
				Victor:        player.Victor,
				Archived:      player.Archived,
			})
		}

		return nil
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create test game: %w", err))
	}

	return connect.NewResponse(&craig_starsv1.CreateTestGameResponse{
		Game: converter.C.ConvertCSGameWithPlayers(game),
	}), nil
}

// GetTestGameNames returns the list of available test game names
func (s *testService) GetTestGameNames(ctx context.Context, req *connect.Request[craig_starsv1.GetTestGameNamesRequest]) (*connect.Response[craig_starsv1.GetTestGameNamesResponse], error) {
	names := make([]string, len(testgames.TestGames))
	for i, testGame := range testgames.TestGames {
		names[i] = testGame.Name
	}

	return connect.NewResponse(&craig_starsv1.GetTestGameNamesResponse{
		Names: names,
	}), nil
}
