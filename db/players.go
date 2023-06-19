package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) FindPlayerByGameId(gameID uint64, userID uint64) (*game.FullPlayer, error) {
	player := game.FullPlayer{}

	if err := db.sqlDB.
		Preload(clause.Associations).
		Preload("Designs").
		Preload("PlanetIntels").
		Preload("FleetIntels").
		Preload("DesignIntels").
		Preload("MineralPacketIntels").
		Preload("BattlePlans").
		Where("game_id = ? AND user_id = ?", gameID, userID).
		First(&player.Player).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if err := db.sqlDB.
		Preload("Tokens.Design").
		Where("player_id = ? and starbase != ?", player.ID, true).
		Order("num").
		Find(&player.Fleets).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			player.Fleets = []*game.Fleet{}
		} else {
			return nil, err
		}
	}

	if err := db.sqlDB.
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
func (db *DB) FindPlayerByGameIdLight(gameID uint64, userID uint64) (*game.Player, error) {
	player := game.Player{}

	if err := db.sqlDB.
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
func (db *DB) SavePlayer(player *game.Player) error {
	return db.sqlDB.Save(player).Error
}
