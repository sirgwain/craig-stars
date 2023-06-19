package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) GetGames() []game.Game {

	games := []game.Game{}
	db.sqlDB.Find(&games)

	return games
}

func (db *DB) CreateGame(game *game.Game) error {
	err := db.sqlDB.Omit("Planets", "Fleets", "Players").Create(game).Error
	if err != nil {
		return err
	}

	err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20}).Model(game).Association("Players").Replace(game.Players)
	if err != nil {
		return err
	}
	err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20}).Model(game).Association("Planets").Replace(game.Planets)
	if err != nil {
		return err
	}
	err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20}).Model(game).Association("Fleets").Replace(game.Fleets)
	if err != nil {
		return err
	}

	// save each player's data
	for i := range game.Players {
		player := &game.Players[i]
		err = db.sqlDB.Omit("PlanetIntels", "Fleets").Save(player).Error
		if err != nil {
			return err
		}
		err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20}).Model(player).Association("PlanetIntels").Replace(player.PlanetIntels)
		if err != nil {
			return err
		}
		err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20}).Model(player).Association("Fleets").Replace(player.Fleets)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) SaveGame(game *game.Game) error {

	// rules don't change after game is created
	err := db.sqlDB.Omit("Planets", "Fleets", "Players").Save(game).Error
	if err != nil {
		return err
	}

	for i := range game.Planets {
		planet := &game.Planets[i]
		if planet.Dirty {
			err = db.sqlDB.Save(planet).Error
			if err != nil {
				return err
			}
		}
	}
	for i := range game.Fleets {
		fleet := &game.Fleets[i]
		if fleet.Dirty {
			err = db.sqlDB.Save(fleet).Error
			if err != nil {
				return err
			}
		}
	}

	// save each player's data
	for i := range game.Players {
		player := &game.Players[i]
		// race doesn't change after game is created
		err = db.sqlDB.Omit("Race", "PlanetIntels", "Planets", "Fleets").Save(player).Error
		if err != nil {
			return err
		}

		for j := range player.PlanetIntels {
			planet := &player.PlanetIntels[j]
			if planet.Dirty {
				err = db.sqlDB.Save(&player.PlanetIntels[j]).Error
				if err != nil {
					return err
				}
			}
		}
		for j := range player.Planets {
			planet := &player.Planets[j]
			if planet.Dirty {
				err = db.sqlDB.Save(&player.Planets[j]).Error
				if err != nil {
					return err
				}
			}
		}
		for j := range player.Fleets {
			fleet := &player.Fleets[j]
			if fleet.Dirty {
				err = db.sqlDB.Save(&player.Fleets[j]).Error
				if err != nil {
					return err
				}
			}
		}

	}

	return nil
}

func (db *DB) FindGameById(id uint) (*game.Game, error) {
	game := game.Game{}
	if err := db.sqlDB.Preload(clause.Associations).Preload("Planets.ProductionQueue").Preload("Players.Race").Preload("Players.PlanetIntels").Preload("Players.Fleets").First(&game, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	

	return &game, nil
}

func (db *DB) FindGameByIdLight(id uint) (*game.Game, error) {
	game := game.Game{}
	if err := db.sqlDB.First(&game, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &game, nil
}

func (db *DB) GetGamesByUser(userID uint) []game.Game {
	games := []game.Game{}
	db.sqlDB.Raw("SELECT * from games g WHERE g.id in (SELECT game_id from players p WHERE p.user_id = ?)", userID).Scan(&games)

	return games
}

func (db *DB) DeleteGameById(id uint) {
	db.sqlDB.Delete(&game.Game{ID: id})
}
