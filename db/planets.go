package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) GetPlanet(id int64) (*game.Planet, error) {
	planet := game.Planet{}
	if err := c.sqlDB.Preload(clause.Associations).First(&planet, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &planet, nil
}

func (c *client) UpdatePlanet(planet *game.Planet) error {

	// save the planet and all production queue items
	if err := c.sqlDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(planet).Error; err != nil {
		return err
	}

	return nil
}
