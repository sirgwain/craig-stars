package db

import (
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/game"
)

func (c *client) migrateAll() error {

	types := []interface{}{
		&game.User{},
		&game.Game{},
		&game.Rules{},
		&game.Player{},
		&game.Race{},
		&game.PlayerMessage{},
		&game.BattlePlan{},
		&game.Fleet{},
		&game.ShipToken{},
		&game.ShipDesign{},
		&game.Planet{},
		&game.MineralPacket{},
		&game.Salvage{},
		&game.Wormohole{},
		&game.MineField{},
		&game.PlanetIntel{},
		&game.FleetIntel{},
		&game.ShipDesignIntel{},
		&game.MineralPacketIntel{},
		&game.MineFieldIntel{},
		&game.TechStore{},
		&game.TechEngine{},
		&game.TechPlanetaryScanner{},
		&game.TechDefense{},
		&game.TechHullComponent{},
		&game.TechHull{},
	}

	for _, t := range types {
		log.Info().Msgf("Migrating %v", reflect.TypeOf(t))
		err := c.sqlDB.AutoMigrate(t)
		if err != nil {
			return err
		}
	}

	return nil
}
