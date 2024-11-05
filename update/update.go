package update

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/db"
)

// Update the host of a game.
func UpdateHost(gameID int64, userID int64) error {

	cfg := config.GetConfig()

	// create a new connection to the database
	dbConn := db.NewConn()
	if err := dbConn.Connect(cfg); err != nil {
		return err
	}
	defer func() { dbConn.Close() }()
	db := dbConn.NewReadWriteClient()

	user, err := db.GetUser(userID)
	if err != nil {
		return fmt.Errorf("failed to load user: %d, %v", userID, err)
	}

	if user == nil {
		return fmt.Errorf("user %d not found", userID)
	}

	game, err := db.GetGame(gameID)
	if err != nil {
		return fmt.Errorf("failed to load game: %d, %v", gameID, err)
	}

	if game == nil {
		return fmt.Errorf("game %d not found", gameID)
	}

	db.UpdateGameHost(game.ID, user.ID)
	log.Info().Msgf("updated game %d host to userID %d", game.ID, userID)

	return nil
}

// Update the host of a game.
func UpdatePlayer(gameID int64, playerNum int, userID int64) error {

	cfg := config.GetConfig()

	// create a new connection to the database
	dbConn := db.NewConn()
	if err := dbConn.Connect(cfg); err != nil {
		return err
	}
	defer func() { dbConn.Close() }()
	db := dbConn.NewReadWriteClient()

	user, err := db.GetUser(userID)
	if err != nil {
		return fmt.Errorf("failed to load user: %d, %v", userID, err)
	}

	if user == nil {
		return fmt.Errorf("user %d not found", userID)
	}

	player, err := db.GetPlayerByNum(gameID, playerNum)
	if err != nil {
		return fmt.Errorf("failed to load player %d from game %d, %v", playerNum, gameID, err)
	}

	if player == nil {
		return fmt.Errorf("player %d game %d not found", playerNum, gameID)
	}

	// update this player's userID
	player.UserID = userID
	db.UpdatePlayerUserId(player)
	log.Info().Msgf("updated game %d, player %d to userID %d", player.GameID, player.Num, userID)

	return nil
}
