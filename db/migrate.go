package db

import (
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/game"
)

func (db *DB) MigrateAll() error {

	types := []interface{}{
		&game.User{},
		&game.Race{},
		&game.Player{},
		&game.TechStore{},
		&game.TechEngine{},
		&game.TechPlanetaryScanner{},
		&game.TechDefense{},
		&game.TechHullComponent{},
		&game.Rules{},
		&game.Game{},
		&game.Planet{},
		&game.PlanetIntel{},
		&game.FleetIntel{},
		&game.ProductionQueueItem{},
		&game.Fleet{},
		&game.ShipDesign{},
	}

	for _, t := range types {
		log.Info().Msgf("Migrating %v", reflect.TypeOf(t))
		err := db.sqlDB.AutoMigrate(t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) Migrate(item interface{}) {
	db.sqlDB.AutoMigrate(&item)
}
