package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)


func (db *DB) FindPlanetByID(id uint64) (*game.Planet, error) {
	planet := game.Planet{}
	if err := db.sqlDB.Preload(clause.Associations).First(&planet, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &planet, nil
}

func (db *DB) FindPlanetByNum(gameID uint64, num int) (*game.Planet, error) {
	planet := game.Planet{}
	if err := db.sqlDB.Preload(clause.Associations).Where("game_id = ? AND num = ?", gameID, num).First(&planet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &planet, nil
}

func (db *DB) SavePlanet(gameID uint64, planet *game.Planet) error {

	// save the planet and all production queue items
	if err := db.sqlDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(planet).Error; err != nil {
		return err
	}

	return nil
}
