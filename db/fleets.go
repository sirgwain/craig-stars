package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) FindFleetByID(id uint64) (*game.Fleet, error) {
	fleet := game.Fleet{}
	if err := db.sqlDB.Preload(clause.Associations).First(&fleet, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &fleet, nil
}

func (db *DB) FindFleetByNum(gameID uint64, playerNum int, num int) (*game.Fleet, error) {
	fleet := game.Fleet{}
	if err := db.sqlDB.Preload(clause.Associations).Where("game_id = ? AND player_num = ? AND num = ?", gameID, playerNum, num).First(&fleet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &fleet, nil
}

func (db *DB) SaveFleet(gameID uint64, fleet *game.Fleet) error {

	// save the fleet and all production queue items
	if err := db.sqlDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(fleet).Error; err != nil {
		return err
	}

	err := db.sqlDB.Model(fleet).Association("Tokens").Replace(fleet.Tokens)

	return err
}

