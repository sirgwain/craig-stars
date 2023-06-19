package db

import (
	"errors"
	"time"

	"github.com/sirgwain/craig-stars/game"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) GetGames() ([]game.Game, error) {

	games := []game.Game{}
	if err := c.sqlDB.Find(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (c *client) GetGamesForUser(userID int64) ([]game.Game, error) {
	games := []game.Game{}
	if err := c.sqlDB.Raw("SELECT * from games g WHERE g.id in (SELECT game_id from players p WHERE p.user_id = ?)", userID).Scan(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (c *client) GetGamesForHost(userID int64) ([]game.Game, error) {
	games := []game.Game{}
	if err := c.sqlDB.Raw("SELECT * from games g WHERE g.host_id = ?", userID).Scan(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (c *client) GetOpenGames() ([]game.Game, error) {
	games := []game.Game{}
	if err := c.sqlDB.Raw("SELECT * from games g WHERE g.state = ? AND g.open_player_slots > 0", game.GameStateSetup).Scan(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (c *client) CreateGame(game *game.Game) error {
	return c.sqlDB.Create(game).Error
}

func (c *client) UpdateFullGame(g *game.FullGame) error {
	defer timeTrack(time.Now(), "UpdateFullGame")

	err := c.sqlDB.Save(g.Game).Error
	if err != nil {
		return err
	}

	for _, planet := range g.Universe.Planets {
		if planet.Dirty {
			planet.GameID = g.ID
			err = c.UpdatePlanet(planet)
			if err != nil {
				return err
			}
		}
	}
	for i := 0; i < len(g.Universe.Fleets); i++ {
		fleet := g.Universe.Fleets[i]
		if fleet.Dirty {
			fleet.GameID = g.ID
			err = c.sqlDB.Save(fleet).Error
			if err != nil {
				return err
			}
		}
		if fleet.Delete {
			err = c.sqlDB.Delete(fleet).Error
			if err != nil {
				return err
			}
			g.Universe.Fleets = append(g.Universe.Fleets[:i], g.Universe.Fleets[i+1:]...)
			i--
		}
	}

	for _, packet := range g.Universe.MineralPackets {
		if packet.Dirty {
			err = c.sqlDB.Save(packet).Error
			if err != nil {
				return err
			}
		}
	}

	// save each player's data
	for _, player := range g.Players {
		player.GameID = g.ID
		// race doesn't change after game is created
		err = c.sqlDB.Omit("Designs").Save(player).Error
		if err != nil {
			return err
		}

		for j := range player.Designs {
			design := player.Designs[j]
			if design.Dirty {
				design.PlayerID = player.ID
				err = c.sqlDB.Save(&player.Designs[j]).Error
				if err != nil {
					return err
				}
			}
		}

	}

	return nil
}

func (c *client) GetFullGame(id int64) (*game.FullGame, error) {
	g := game.FullGame{
		Game:     &game.Game{},
		Universe: &game.Universe{},
		Players:  []*game.Player{},
	}
	if err := c.sqlDB.First(&g.Game, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if err := c.sqlDB.
		Preload("Starbase").
		Where("game_id = ?", id).
		Order("planets.num").
		Find(&g.Universe.Planets).
		Error; err != nil {
		return nil, err
	}

	if err := c.sqlDB.
		Preload("Tokens").
		Where("game_id = ?", id).
		Order("fleets.num").
		Find(&g.Universe.Fleets).
		Error; err != nil {
		return nil, err
	}

	players := []game.Player{}
	if err := c.sqlDB.
		Preload(clause.Associations).
		Preload("Designs").
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
		techs, err := c.GetTechStore(g.Rules.TechsID)
		if err != nil {
			return nil, err
		}
		g.Rules.WithTechStore(techs)
	}

	// init the random generator after load
	(&g.Rules).ResetSeed(g.Seed)

	return &g, nil
}

func (c *client) GetGame(id int64) (*game.Game, error) {
	game := game.Game{}
	if err := c.sqlDB.First(&game, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &game, nil
}

func (c *client) GetRulesForGame(gameID int64) (*game.Rules, error) {
	rules := game.Rules{}
	if err := c.sqlDB.
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

func (c *client) DeleteGame(id int64) error {
	return c.sqlDB.Delete(&game.Game{ID: id}).Error
}
