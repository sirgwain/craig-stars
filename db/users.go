package db

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
)

func (db *DB) GetUsers() ([]*game.User, error) {

	users := []*game.User{}
	if err := db.sqlDB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (db *DB) SaveUser(user *game.User) error {
	log.Debug().Msgf("Creating user %s", user.Username)
	err := db.sqlDB.Save(&user).Error
	return err
}

func (db *DB) FindUserById(id uint64) (*game.User, error) {
	user := game.User{}

	if err := db.sqlDB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (db *DB) FindUserByUsername(username string) (*game.User, error) {
	user := game.User{}
	if err := db.sqlDB.
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

func (db *DB) DeleteUserById(id uint64) error {
	return db.sqlDB.Delete(&game.User{ID: id}).Error
}
