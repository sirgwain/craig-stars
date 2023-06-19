package db

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
)

func (c *client) GetUsers() ([]game.User, error) {

	users := []game.User{}
	if err := c.sqlDB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (c *client) CreateUser(user *game.User) error {
	log.Debug().Msgf("Creating user %s", user.Username)
	err := c.sqlDB.Save(&user).Error
	return err
}

func (c *client) UpdateUser(user *game.User) error {
	log.Debug().Msgf("Updating user %s", user.Username)
	err := c.sqlDB.Save(&user).Error
	return err
}

func (c *client) GetUser(id int64) (*game.User, error) {
	user := game.User{}

	if err := c.sqlDB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (c *client) GetUserByUsername(username string) (*game.User, error) {
	user := game.User{}
	if err := c.sqlDB.
		Where("username = ?", username).
		First(&user).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (c *client) DeleteUser(id int64) error {
	return c.sqlDB.Delete(&game.User{ID: id}).Error
}
