package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) FindPlayerByGameId(gameID uint, userID uint) (*game.Player, error) {
	player := game.Player{}

	if err := db.sqlDB.Preload(clause.Associations).Preload("Planets.ProductionQueue").Where("game_id = ? AND user_id = ?", gameID, userID).First(&player).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &player, nil
}

// find a plyer for a game without loading all data
func (db *DB) FindPlayerByGameIdLight(gameID uint, userID uint) (*game.Player, error) {
	player := game.Player{}

	if err := db.sqlDB.Where("game_id = ? AND user_id = ?", gameID, userID).First(&player).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &player, nil
}

// Save a player to the db
func (db *DB) SavePlayer(player *game.Player) error {
	return db.sqlDB.Save(player).Error
}
