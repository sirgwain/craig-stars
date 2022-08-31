package db

import (
	"errors"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) FindPlayerByGameId(gameID uint, userID uint) (*game.FullPlayer, error) {
	player := game.FullPlayer{}

	if err := db.sqlDB.
		Preload(clause.Associations).
		Preload("Designs").
		Preload("PlanetIntels").
		Preload("FleetIntels").
		Preload("DesignIntels").
		Preload("MineralPacketIntels").
		Preload("ProductionPlans").
		Preload("BattlePlans").
		Preload("Fleets", func(db *gorm.DB) *gorm.DB {
			return db.Where("fleets.starbase != ?", true).Order("fleets.num")
		}).
		Preload("Fleets.Tokens.Design").
		Preload("Planets", func(db *gorm.DB) *gorm.DB {
			return db.Order("planets.num")
		}).
		Preload("Planets.Starbase").
		Preload("Planets.Starbase.Tokens").
		Preload("Planets.ProductionQueue", func(db *gorm.DB) *gorm.DB {
			return db.Order("production_queue_items.sort_order")
		}).Where("game_id = ? AND user_id = ?", gameID, userID).First(&player).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	// build the non-serialized player map objects
	player.BuildMaps()

	return &player, nil
}

// find a plyer for a game without loading all data
func (db *DB) FindPlayerByGameIdLight(gameID uint, userID uint) (*game.Player, error) {
	player := game.Player{}

	if err := db.sqlDB.Where("game_id = ? AND user_id = ?", gameID, userID).Preload("Race").First(&player).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &player, nil
}

// Save a player to the db
func (db *DB) SavePlayer(player *game.Player) error {
	return db.sqlDB.Save(player).Error
}
