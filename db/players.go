package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) GetFullPlayerForGame(gameID int64, userID int64) (*game.FullPlayer, error) {
	player := game.FullPlayer{}

	if err := c.sqlDB.
		Preload(clause.Associations).
		Preload("Designs").
		Where("game_id = ? AND user_id = ?", gameID, userID).
		First(&player.Player).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if err := c.sqlDB.
		Preload("Tokens").
		Where("player_id = ? and starbase != ?", player.ID, true).
		Order("num").
		Find(&player.Fleets).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			player.Fleets = []*game.Fleet{}
		} else {
			return nil, err
		}
	}

	if err := c.sqlDB.
		Preload("Starbase").
		Preload("Starbase.Tokens").
		Where("player_id = ?", player.ID).
		Order("num").
		Find(&player.Planets).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			player.Planets = []*game.Planet{}
		} else {
			return nil, err
		}
	}

	return &player, nil
}

// find a plyer for a game without loading all data
func (c *client) GetPlayerForGame(gameID int64, userID int64) (*game.Player, error) {
	player := game.Player{}

	if err := c.sqlDB.
		Where("game_id = ? AND user_id = ?", gameID, userID).
		First(&player).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &player, nil
}

// Save a player to the db
func (c *client) CreatePlayer(player *game.Player) error {
	return c.sqlDB.Save(player).Error
}

// Save a player to the db
func (c *client) UpdatePlayer(player *game.Player) error {
	return c.sqlDB.Save(player).Error
}
