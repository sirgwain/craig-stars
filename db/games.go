package db

import (
	"errors"
	"math/rand"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) GetGames() ([]game.Game, error) {

	games := []game.Game{}
	if err := db.sqlDB.Find(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (db *DB) GetGamesByUser(userID uint) ([]game.Game, error) {
	games := []game.Game{}
	if err := db.sqlDB.Raw("SELECT * from games g WHERE g.id in (SELECT game_id from players p WHERE p.user_id = ?)", userID).Scan(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (db *DB) GetGamesHostedByUser(userID uint) ([]game.Game, error) {
	games := []game.Game{}
	if err := db.sqlDB.Raw("SELECT * from games g WHERE g.host_id = ?", userID).Scan(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (db *DB) GetOpenGames() ([]game.Game, error) {
	games := []game.Game{}
	if err := db.sqlDB.Raw("SELECT * from games g WHERE g.state = ? AND g.open_player_slots > 0", game.GameStateSetup).Scan(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (db *DB) CreateGame(game *game.Game) error {
	err := db.sqlDB.Omit("Planets", "Fleets", "MineralPackets", "Players").Create(game).Error
	if err != nil {
		return err
	}

	err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20, FullSaveAssociations: true}).Model(game).Association("Players").Replace(game.Players)
	if err != nil {
		return err
	}
	err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20, FullSaveAssociations: true}).Model(game).Association("Planets").Replace(game.Planets)
	if err != nil {
		return err
	}
	err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20, FullSaveAssociations: true}).Model(game).Association("Fleets").Replace(game.Fleets)
	if err != nil {
		return err
	}
	err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20, FullSaveAssociations: true}).Model(game).Association("MineralPackets").Replace(game.MineralPackets)
	if err != nil {
		return err
	}

	// save each player's data
	for i := range game.Players {
		player := &game.Players[i]
		err = db.sqlDB.Omit("Messages", "Designs", "Planets", "Fleets", "MineralPackets", "PlanetIntels", "FleetIntels", "DesignIntels", "MineralPacketIntels").Save(player).Error
		if err != nil {
			return err
		}

		err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20, FullSaveAssociations: true}).Model(player).Association("Messages").Replace(player.Messages)
		if err != nil {
			return err
		}
		err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20, FullSaveAssociations: true}).Model(player).Association("PlanetIntels").Replace(player.PlanetIntels)
		if err != nil {
			return err
		}
		err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20, FullSaveAssociations: true}).Model(player).Association("FleetIntels").Replace(player.FleetIntels)
		if err != nil {
			return err
		}
		err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20, FullSaveAssociations: true}).Model(player).Association("DesignIntels").Replace(player.DesignIntels)
		if err != nil {
			return err
		}
		err = db.sqlDB.Session(&gorm.Session{CreateBatchSize: 20, FullSaveAssociations: true}).Model(player).Association("MineralPacketIntels").Replace(player.MineralPacketIntels)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) SaveGame(game *game.Game) error {

	// rules don't change after game is created
	err := db.sqlDB.Omit("Planets", "Fleets", "MineralPackets", "Players").Save(game).Error
	if err != nil {
		return err
	}

	for i := range game.Planets {
		planet := &game.Planets[i]
		if planet.Dirty {
			err = db.SavePlanet(planet)
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

	for i := range game.MineralPackets {
		packet := &game.MineralPackets[i]
		if packet.Dirty {
			err = db.sqlDB.Save(packet).Error
			if err != nil {
				return err
			}
		}
	}

	// save each player's data
	for i := range game.Players {
		player := &game.Players[i]
		// race doesn't change after game is created
		err = db.sqlDB.Omit("Race", "PlanetIntels", "FleetIntels", "DesignIntels", "MineralPacketIntels", "Messages", "Planets", "Fleets", "Designs").Save(player).Error
		if err != nil {
			return err
		}

		// replace messages
		err := db.sqlDB.Model(player).Association("Messages").Replace(player.Messages)
		if err != nil {
			return err
		}

		// if the AI player updated their planet/fleet orders, save them
		for j := range player.Planets {
			planet := player.Planets[j]
			if planet.Dirty {
				err = db.SavePlanet(planet)
				if err != nil {
					return err
				}
			}
		}

		for j := range player.Fleets {
			fleet := player.Fleets[j]
			if fleet.Dirty {
				err = db.SaveFleet(fleet)
				if err != nil {
					return err
				}
			}
		}

		for j := range player.Designs {
			design := player.Designs[j]
			if design.Dirty {
				err = db.sqlDB.Save(&player.Designs[j]).Error
				if err != nil {
					return err
				}
			}
		}

		for j := range player.PlanetIntels {
			intel := &player.PlanetIntels[j]
			if intel.Dirty {
				err = db.sqlDB.Save(&player.PlanetIntels[j]).Error
				if err != nil {
					return err
				}
			}
		}

		for j := range player.FleetIntels {
			intel := &player.FleetIntels[j]
			if intel.Dirty {
				err = db.sqlDB.Save(&player.FleetIntels[j]).Error
				if err != nil {
					return err
				}
			}
		}

		for j := range player.DesignIntels {
			intel := &player.DesignIntels[j]
			if intel.Dirty {
				err = db.sqlDB.Save(&player.DesignIntels[j]).Error
				if err != nil {
					return err
				}
			}
		}

		for j := range player.MineralPacketIntels {
			intel := &player.MineralPacketIntels[j]
			if intel.Dirty {
				err = db.sqlDB.Save(&player.MineralPacketIntels[j]).Error
				if err != nil {
					return err
				}
			}
		}

	}

	return nil
}

func (db *DB) FindGameById(id uint) (*game.Game, error) {
	g := game.Game{}
	if err := db.sqlDB.
		Preload(clause.Associations).Preload("Planets", func(db *gorm.DB) *gorm.DB {
		return db.Order("planets.num")
	}).
		Preload("Planets.ProductionQueue", func(db *gorm.DB) *gorm.DB {
			return db.Order("production_queue_items.sort_order")
		}).
		Preload(clause.Associations).Preload("Fleets", func(db *gorm.DB) *gorm.DB {
		return db.Order("fleets.num")
	}).
		Preload("Fleets.Tokens.Design").
		Preload("Planets.Starbase").
		Preload("Players.Fleets.Tokens").
		Preload("Players.Planets.Starbase").
		Preload("Players.Planets.ProductionQueue").
		Preload("Players.Race").
		Preload("Players.Designs").
		Preload("Players.MineralPackets").
		Preload("Players.PlanetIntels").
		Preload("Players.FleetIntels").
		Preload("Players.DesignIntels").
		Preload("Players.BattlePlans").
		Preload("Players.ProductionPlans").
		Preload("Players.Messages").
		First(&g, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if g.Rules.TechsID == 0 {
		g.Rules.Techs = &game.StaticTechStore
	} else {
		techs, err := db.FindTechStoreById(g.Rules.TechsID)
		if err != nil {
			return nil, err
		}
		g.Rules.Techs = techs
	}

	// init the random generator after load
	g.Rules.Random = rand.New(rand.NewSource(g.Rules.Seed))

	return &g, nil
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

func (db *DB) FindGameRulesByGameId(gameId uint) (*game.Rules, error) {
	rules := game.Rules{}
	if err := db.sqlDB.Where("game_id = ? ", gameId).First(&rules).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &rules, nil
}

func (db *DB) DeleteGameById(id uint) error {
	return db.sqlDB.Delete(&game.Game{ID: id}).Error
}
