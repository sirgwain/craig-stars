package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) GetTechStores() ([]*game.TechStore, error) {

	techs := []*game.TechStore{}
	if err := c.sqlDB.Find(&techs).Error; err != nil {
		return nil, err
	}

	return techs, nil
}

func (c *client) CreateTechStore(tech *game.TechStore) error {
	err := c.sqlDB.Omit("Engines", "PlanetaryScanners", "Terraforms", "Defenses", "HullComponents", "Hulls").Create(tech).Error
	if err != nil {
		return err
	}

	err = c.sqlDB.Session(&gorm.Session{CreateBatchSize: 10}).Model(tech).Association("Engines").Replace(tech.Engines)
	if err != nil {
		return err
	}
	err = c.sqlDB.Session(&gorm.Session{CreateBatchSize: 10}).Model(tech).Association("PlanetaryScanners").Replace(tech.PlanetaryScanners)
	if err != nil {
		return err
	}
	err = c.sqlDB.Session(&gorm.Session{CreateBatchSize: 10}).Model(tech).Association("Terraforms").Replace(tech.Terraforms)
	if err != nil {
		return err
	}
	err = c.sqlDB.Session(&gorm.Session{CreateBatchSize: 10}).Model(tech).Association("Defenses").Replace(tech.Defenses)
	if err != nil {
		return err
	}
	err = c.sqlDB.Session(&gorm.Session{CreateBatchSize: 10}).Model(tech).Association("HullComponents").Replace(tech.HullComponents)
	if err != nil {
		return err
	}
	err = c.sqlDB.Session(&gorm.Session{CreateBatchSize: 10}).Model(tech).Association("Hulls").Replace(tech.Hulls)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) FindTechStoreById(id uint64) (*game.TechStore, error) {
	techs := game.TechStore{}
	if err := c.sqlDB.Preload(clause.Associations).First(&techs, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	techs.Init()
	return &techs, nil
}
