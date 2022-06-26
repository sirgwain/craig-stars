package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
)

func (db *DB) GetRaces(userID uint) ([]game.Race, error) {
	races := []game.Race{}
	if err := db.sqlDB.Where("user_id = ? AND player_id = 0", userID).Find(&races).Error; err != nil {
		return nil, err
	}

	return races, nil
}

func (db *DB) FindRaceById(id uint) (*game.Race, error) {
	race := game.Race{}
	if err := db.sqlDB.First(&race, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &race, nil
}

func (db *DB) CreateRace(race *game.Race) error {
	err := db.sqlDB.Create(race).Error
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) SaveRace(race *game.Race) error {
	err := db.sqlDB.Save(race).Error
	if err != nil {
		return err
	}

	return nil
}
