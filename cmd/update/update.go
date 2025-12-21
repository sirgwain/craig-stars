package update

import (
	"context"
	"fmt"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/db"
	"log/slog"
)

// Update the host of a game.
func UpdateHost(gameID int64, userID int64) error {

	ctx := context.Background()
	cfg := config.GetConfig()

	// create a new connection to the database
	dbConn := db.NewConn()
	if err := dbConn.Connect(ctx, cfg); err != nil {
		return err
	}
	defer func() { dbConn.Close() }()
	db := dbConn.NewReadWriteClient()

	user, err := db.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to load user: %d, %v", userID, err)
	}

	if user == nil {
		return fmt.Errorf("user %d not found", userID)
	}

	game, err := db.GetGame(ctx, gameID)
	if err != nil {
		return fmt.Errorf("failed to load game: %d, %v", gameID, err)
	}

	if game == nil {
		return fmt.Errorf("game %d not found", gameID)
	}

	db.UpdateGameHost(ctx, game.ID, user.ID)
	slog.Info("updated game host", slog.Int64("gameID", game.ID), slog.Int64("userID", userID))

	return nil
}

// Update the player of a game.
func UpdatePlayer(gameID int64, playerNum int, userID int64) error {

	ctx := context.Background()
	cfg := config.GetConfig()

	// create a new connection to the database
	dbConn := db.NewConn()
	if err := dbConn.Connect(ctx, cfg); err != nil {
		return err
	}
	defer func() { dbConn.Close() }()
	readWriteClient := dbConn.NewReadWriteClient()

	user, err := readWriteClient.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to load user: %d, %v", userID, err)
	}

	if user == nil {
		return fmt.Errorf("user %d not found", userID)
	}

	player, err := readWriteClient.GetPlayerForGame(ctx, gameID, playerNum)
	if err != nil {
		return fmt.Errorf("failed to load player %d from game %d, %v", playerNum, gameID, err)
	}

	if player == nil {
		return fmt.Errorf("player %d game %d not found", playerNum, gameID)
	}

	// update this player's userID
	player.UserID = userID
	readWriteClient.UpdatePlayerUserID(ctx, player)
	slog.Info("updated game player", slog.Int64("gameID", player.GameID), slog.Int("playerNum", player.Num), slog.Int64("userID", userID))

	return nil
}
