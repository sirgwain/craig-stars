package db

import (
	"errors"
	"time"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) GetGames() ([]*game.Game, error) {

	games := []*game.Game{}
	if err := db.sqlDB.Find(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (db *DB) GetGamesByUser(userID uint64) ([]*game.Game, error) {
	games := []*game.Game{}
	if err := db.sqlDB.Raw("SELECT * from games g WHERE g.id in (SELECT game_id from players p WHERE p.user_id = ?)", userID).Scan(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (db *DB) GetGamesHostedByUser(userID uint64) ([]*game.Game, error) {
	games := []*game.Game{}
	if err := db.sqlDB.Raw("SELECT * from games g WHERE g.host_id = ?", userID).Scan(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (db *DB) GetOpenGames() ([]*game.Game, error) {
	games := []*game.Game{}
	if err := db.sqlDB.Raw("SELECT * from games g WHERE g.state = ? AND g.open_player_slots > 0", game.GameStateSetup).Scan(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (db *DB) CreateGame(game *game.Game) error {
	return db.sqlDB.Create(game).Error
}

func (db *DB) SaveGame(g *game.FullGame) error {
	defer timeTrack(time.Now(), "SaveGame")

	// rules don't change after game is created
	if (g.Game.Area == game.Vector{}) {
		g.Game.Area = g.Universe.Area
	}
	err := db.sqlDB.Save(g.Game).Error
	if err != nil {
		return err
	}

	for _, planet := range g.Universe.Planets {
		if planet.Dirty {
			planet.GameID = g.ID
			err = db.SavePlanet(g.ID, planet)
			if err != nil {
				return err
			}
		}
	}
	for i := 0; i < len(g.Universe.Fleets); i++ {
		fleet := g.Universe.Fleets[i]
		if fleet.Dirty {
			fleet.GameID = g.ID
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
			g.Universe.Fleets = append(g.Universe.Fleets[:i], g.Universe.Fleets[i+1:]...)
			i--
		}
	}

	for _, packet := range g.Universe.MineralPackets {
		if packet.Dirty {
			err = db.sqlDB.Save(packet).Error
			if err != nil {
				return err
			}
		}
	}

	// save each player's data
	for _, player := range g.Players {
		player.GameID = g.ID
		// race doesn't change after game is created
		err = db.sqlDB.Omit("PlanetIntels", "FleetIntels", "DesignIntels", "MineralPacketIntels", "Messages", "Designs").Save(player).Error
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
				design.GameID = g.ID
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
				intel.GameID = g.ID
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
				intel.GameID = g.ID
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
				intel.GameID = g.ID
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
				intel.GameID = g.ID
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

func (db *DB) FindGameById(id uint64) (*game.FullGame, error) {
	g := game.FullGame{
		Game:     &game.Game{},
		Universe: &game.Universe{},
		Players:  []*game.Player{},
	}
	if err := db.sqlDB.First(&g.Game, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if err := db.sqlDB.
		Preload("Starbase").
		Where("game_id = ?", id).
		Order("planets.num").
		Find(&g.Universe.Planets).
		Error; err != nil {
		return nil, err
	}

	if err := db.sqlDB.
		Preload("Tokens.Design").
		Where("game_id = ?", id).
		Order("fleets.num").
		Find(&g.Universe.Fleets).
		Error; err != nil {
		return nil, err
	}

	players := []game.Player{}
	if err := db.sqlDB.
		Preload(clause.Associations).
		Preload("Designs").
		Preload("PlanetIntels").
		Preload("FleetIntels").
		Preload("DesignIntels").
		Preload("BattlePlans").
		Preload("Messages").
		Order("num").
		Where("game_id = ?", id).
		Find(&players).
		Error; err != nil {
		return nil, err
	}
	for i := range players {
		g.Players = append(g.Players, &players[i])
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
	(&g.Rules).ResetSeed()

	return &g, nil
}

func (db *DB) FindGameByIdLight(id uint64) (*game.Game, error) {
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

func (db *DB) FindGameRulesByGameID(gameID uint64) (*game.Rules, error) {
	rules := game.Rules{}
	if err := db.sqlDB.
		Where("game_id = ? ", gameID).
		First(&rules).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &rules, nil
}

func (db *DB) DeleteGameById(id uint64) error {
	return db.sqlDB.Delete(&game.Game{ID: id}).Error
}
