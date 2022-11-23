package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) GetFleet(id int64) (*game.Fleet, error) {
	fleet := game.Fleet{}
	if err := c.sqlDB.Preload(clause.Associations).Preload("Tokens.Design").First(&fleet, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &fleet, nil
}

func (c *client) UpdateFleet(fleet *game.Fleet) error {

	// save the fleet and all production queue items
	if err := c.sqlDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(fleet).Error; err != nil {
		return err
	}

	err := c.sqlDB.Model(fleet).Association("Tokens").Replace(fleet.Tokens)

	return err
}
