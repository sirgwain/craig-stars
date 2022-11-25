package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
)

func (c *client) GetRacesForUser(userID int64) ([]game.Race, error) {
	races := []game.Race{}
	if err := c.sqlDB.
		Where("user_id = ? AND player_id = 0", userID).
		Find(&races).
		Error; err != nil {
		return nil, err
	}

	return races, nil
}

func (c *client) GetRace(id int64) (*game.Race, error) {
	race := game.Race{}
	if err := c.sqlDB.First(&race, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &race, nil
}

func (c *client) CreateRace(race *game.Race) error {
	err := c.sqlDB.Save(race).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *client) UpdateRace(race *game.Race) error {
	err := c.sqlDB.Save(race).Error
	if err != nil {
		return err
	}

	return nil
}
