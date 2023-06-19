package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) FindPlanetById(id uint) (*game.Planet, error) {
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

func (db *DB) SavePlanet(planet *game.Planet) error {
	return db.sqlDB.Save(planet).Error
}
