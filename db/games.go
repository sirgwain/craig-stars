package db

import (
	"errors"

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
	return db.sqlDB.Omit("Planets", "Fleets", "MineralPackets", "Players").Create(game).Error
}

func (db *DB) SaveGame(game *game.FullGame) error {
	// rules don't change after game is created
	err := db.sqlDB.Save(game.Game).Error
	if err != nil {
		return err
	}

	for i := range game.Universe.Planets {
		planet := &game.Universe.Planets[i]
		if planet.Dirty {
			planet.GameID = game.ID
			err = db.SavePlanet(planet)
			if err != nil {
				return err
			}
		}
	}
	for i := 0; i < len(game.Universe.Fleets); i++ {
		fleet := &game.Universe.Fleets[i]
		if fleet.Dirty {
			fleet.GameID = game.ID
			err = db.sqlDB.Save(fleet).Error
			if err != nil {
				return err
			}
		}
		if fleet.Delete {
			err = db.sqlDB.Delete(fleet).Error
			if err != nil {
				return err
			}
			game.Universe.Fleets = append(game.Universe.Fleets[:i], game.Universe.Fleets[i+1:]...)
			i--
		}
	}

	for i := range game.Universe.MineralPackets {
		packet := &game.Universe.MineralPackets[i]
		if packet.Dirty {
			err = db.sqlDB.Save(packet).Error
			if err != nil {
				return err
			}
		}
	}

	// save each player's data
	for _, player := range game.Players {
		player.GameID = game.ID
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

		for j := range player.Designs {
			design := player.Designs[j]
			if design.Dirty {
				design.GameID = game.ID
				design.PlayerID = player.ID
				err = db.sqlDB.Save(&player.Designs[j]).Error
				if err != nil {
					return err
				}
			}
		}

		for j := range player.PlanetIntels {
			intel := &player.PlanetIntels[j]
			if intel.Dirty {
				intel.GameID = game.ID
				intel.PlayerID = player.ID
				err = db.sqlDB.Save(&player.PlanetIntels[j]).Error
				if err != nil {
					return err
				}
			}
		}

		for j := range player.FleetIntels {
			intel := &player.FleetIntels[j]
			if intel.Dirty {
				intel.GameID = game.ID
				intel.PlayerID = player.ID
				err = db.sqlDB.Save(&player.FleetIntels[j]).Error
				if err != nil {
					return err
				}
			}
		}

		for j := range player.DesignIntels {
			intel := &player.DesignIntels[j]
			if intel.Dirty {
				intel.GameID = game.ID
				intel.PlayerID = player.ID
				err = db.sqlDB.Save(&player.DesignIntels[j]).Error
				if err != nil {
					return err
				}
			}
		}

		for j := range player.MineralPacketIntels {
			intel := &player.MineralPacketIntels[j]
			if intel.Dirty {
				intel.GameID = game.ID
				intel.PlayerID = player.ID
				err = db.sqlDB.Save(&player.MineralPacketIntels[j]).Error
				if err != nil {
					return err
				}
			}
		}

	}

	return nil
}

func (db *DB) FindGameById(id uint) (*game.FullGame, error) {
	g := game.FullGame{}
	if err := db.sqlDB.Table("games").
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
		g.Rules.WithTechStore(&game.StaticTechStore)
	} else {
		techs, err := db.FindTechStoreById(g.Rules.TechsID)
		if err != nil {
			return nil, err
		}
		g.Rules.WithTechStore(techs)
	}

	// init the random generator after load
	g.Rules.ResetSeed()

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
