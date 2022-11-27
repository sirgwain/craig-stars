package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/game"
)

type Game struct {
	ID                                        int64                `json:"id,omitempty"`
	CreatedAt                                 time.Time            `json:"createdAt,omitempty"`
	UpdatedAt                                 time.Time            `json:"updatedAt,omitempty"`
	Name                                      string               `json:"name,omitempty"`
	HostID                                    int64                `json:"hostId,omitempty"`
	QuickStartTurns                           int                  `json:"quickStartTurns,omitempty"`
	Size                                      game.Size            `json:"size,omitempty"`
	Density                                   game.Density         `json:"density,omitempty"`
	PlayerPositions                           game.PlayerPositions `json:"playerPositions,omitempty"`
	RandomEvents                              bool                 `json:"randomEvents,omitempty"`
	ComputerPlayersFormAlliances              bool                 `json:"computerPlayersFormAlliances,omitempty"`
	PublicPlayerScores                        bool                 `json:"publicPlayerScores,omitempty"`
	StartMode                                 game.GameStartMode   `json:"startMode,omitempty"`
	Year                                      int                  `json:"year,omitempty"`
	State                                     game.GameState       `json:"state,omitempty"`
	OpenPlayerSlots                           uint                 `json:"openPlayerSlots,omitempty"`
	NumPlayers                                int                  `json:"numPlayers,omitempty"`
	VictoryConditionsConditions               *VictoryConditions   `json:"victoryConditionsConditions,omitempty"`
	VictoryConditionsNumCriteriaRequired      int                  `json:"victoryConditionsNumCriteriaRequired,omitempty"`
	VictoryConditionsYearsPassed              int                  `json:"victoryConditionsYearsPassed,omitempty"`
	VictoryConditionsOwnPlanets               int                  `json:"victoryConditionsOwnPlanets,omitempty"`
	VictoryConditionsAttainTechLevel          int                  `json:"victoryConditionsAttainTechLevel,omitempty"`
	VictoryConditionsAttainTechLevelNumFields int                  `json:"victoryConditionsAttainTechLevelNumFields,omitempty"`
	VictoryConditionsExceedsScore             int                  `json:"victoryConditionsExceedsScore,omitempty"`
	VictoryConditionsExceedsSecondPlaceScore  int                  `json:"victoryConditionsExceedsSecondPlaceScore,omitempty"`
	VictoryConditionsProductionCapacity       int                  `json:"victoryConditionsProductionCapacity,omitempty"`
	VictoryConditionsOwnCapitalShips          int                  `json:"victoryConditionsOwnCapitalShips,omitempty"`
	VictoryConditionsHighestScoreAfterYears   int                  `json:"victoryConditionsHighestScoreAfterYears,omitempty"`
	VictorDeclared                            bool                 `json:"victorDeclared,omitempty"`
	Seed                                      int64                `json:"seed,omitempty"`
	Rules                                     *Rules               `json:"rules,omitempty"`
	AreaX                                     float64              `json:"areaX,omitempty"`
	AreaY                                     float64              `json:"areaY,omitempty"`
}

// we json serialize these types with custom Scan/Value methods
type VictoryConditions []game.VictoryCondition
type Rules game.Rules

// db serializer to serialize this to JSON
func (item *VictoryConditions) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *VictoryConditions) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// db serializer to serialize this to JSON
func (item *Rules) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *Rules) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (c *client) GetGames() ([]game.Game, error) {

	items := []Game{}
	if err := c.db.Select(&items, `SELECT * FROM Games`); err != nil {
		if err == sql.ErrNoRows {
			return []game.Game{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertGames(items), nil
}

func (c *client) GetGamesForHost(userID int64) ([]game.Game, error) {

	items := []Game{}
	if err := c.db.Select(&items, `SELECT * FROM Games WHERE hostId = ?`, userID); err != nil {
		if err == sql.ErrNoRows {
			return []game.Game{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertGames(items), nil
}

func (c *client) GetGamesForUser(userID int64) ([]game.Game, error) {

	items := []Game{}
	if err := c.db.Select(&items, `SELECT * from games g WHERE g.id in (SELECT gameId from players p WHERE p.userId = ?)`, userID); err != nil {
		if err == sql.ErrNoRows {
			return []game.Game{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertGames(items), nil
}

func (c *client) GetOpenGames() ([]game.Game, error) {
	items := []Game{}
	if err := c.db.Select(&items, `SELECT * from games g WHERE g.state = ? AND g.openPlayerSlots > 0`, game.GameStateSetup); err != nil {
		if err == sql.ErrNoRows {
			return []game.Game{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertGames(items), nil
}

// get a game by id
func (c *client) GetGame(id int64) (*game.Game, error) {
	item := Game{}
	if err := c.db.Get(&item, "SELECT * FROM games WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	game := c.converter.ConvertGame(item)
	return &game, nil
}

// get a game by id
func (c *client) GetFullGame(id int64) (*game.FullGame, error) {
	item := Game{}
	if err := c.db.Get(&item, "SELECT * FROM games WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	g := c.converter.ConvertGame(item)

	players, err := c.getPlayersForGame(g.ID)
	if err != nil {
		return nil, fmt.Errorf("load players for game %w", err)
	}

	universe := game.Universe{}

	planets, err := c.getPlanetsForGame(g.ID)
	if err != nil {
		return nil, fmt.Errorf("load planets for game %w", err)
	}
	universe.Planets = planets

	// load fleets and starbases
	fleets, err := c.getFleetsForGame(g.ID)
	if err != nil {
		return nil, fmt.Errorf("load fleets for game %w", err)
	}
	// pre-instantiate the fleets/starbases arrays (make it a little bigger than necessary)
	universe.Fleets = make([]*game.Fleet, 0, len(fleets))
	universe.Starbases = make([]*game.Fleet, 0, len(planets))
	for i := range fleets {
		fleet := fleets[i]
		if fleet.Starbase {
			universe.Starbases = append(universe.Starbases, fleet)
		} else {
			universe.Fleets = append(universe.Fleets, fleet)
		}
	}

	wormholes, err := c.getWormholesForGame(g.ID)
	if err != nil {
		return nil, fmt.Errorf("load wormholes for game %w", err)
	}
	universe.Wormholes = wormholes

	salvages, err := c.getSalvagesForGame(g.ID)
	if err != nil {
		return nil, fmt.Errorf("load salvages for game %w", err)
	}
	universe.Salvages = salvages

	mineFields, err := c.getMineFieldsForGame(g.ID)
	if err != nil {
		return nil, fmt.Errorf("load mineFields for game %w", err)
	}
	universe.MineFields = mineFields

	mineralPackets, err := c.getMineralPacketsForGame(g.ID)
	if err != nil {
		return nil, fmt.Errorf("load mineralPackets for game %w", err)
	}
	universe.MineralPackets = mineralPackets

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

	fg := game.FullGame{
		Game:     &g,
		Players:  players,
		Universe: &universe,
	}

	return &fg, nil
}

// create a new game
func (c *client) CreateGame(game *game.Game) error {

	item := c.converter.ConvertGameGame(game)
	result, err := c.db.NamedExec(`
	INSERT INTO games (
		createdAt,
		updatedAt,
		name,
		hostId,
		quickStartTurns,
		size,
		density,
		playerPositions,
		randomEvents,
		computerPlayersFormAlliances,
		publicPlayerScores,
		startMode,
		year,
		state,
		openPlayerSlots,
		numPlayers,
		victoryConditionsConditions,
		victoryConditionsNumCriteriaRequired,
		victoryConditionsYearsPassed,
		victoryConditionsOwnPlanets,
		victoryConditionsAttainTechLevel,
		victoryConditionsAttainTechLevelNumFields,
		victoryConditionsExceedsScore,
		victoryConditionsExceedsSecondPlaceScore,
		victoryConditionsProductionCapacity,
		victoryConditionsOwnCapitalShips,
		victoryConditionsHighestScoreAfterYears,
		victorDeclared,
		seed,
		rules,
		areaX,
		areaY
	)
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:name,
		:hostId,
		:quickStartTurns,
		:size,
		:density,
		:playerPositions,
		:randomEvents,
		:computerPlayersFormAlliances,
		:publicPlayerScores,
		:startMode,
		:year,
		:state,
		:openPlayerSlots,
		:numPlayers,
		:victoryConditionsConditions,
		:victoryConditionsNumCriteriaRequired,
		:victoryConditionsYearsPassed,
		:victoryConditionsOwnPlanets,
		:victoryConditionsAttainTechLevel,
		:victoryConditionsAttainTechLevelNumFields,
		:victoryConditionsExceedsScore,
		:victoryConditionsExceedsSecondPlaceScore,
		:victoryConditionsProductionCapacity,
		:victoryConditionsOwnCapitalShips,
		:victoryConditionsHighestScoreAfterYears,
		:victorDeclared,
		:seed,
		:rules,
		:areaX,
		:areaY
	)
	`, item)

	if err != nil {
		return err
	}

	// update the id of our passed in game
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	game.ID = id

	return nil
}

// update an existing game
func (c *client) UpdateGame(game *game.Game) error {
	return c.updateGameWithNamedExecer(game, c.db)
}

// update a game inside a transaction
func (c *client) updateGameWithNamedExecer(game *game.Game, tx SQLExecer) error {

	item := c.converter.ConvertGameGame(game)

	if _, err := tx.NamedExec(`
	UPDATE games SET
		updatedAt = CURRENT_TIMESTAMP,
		name = :name,
		hostId = :hostId,
		quickStartTurns = :quickStartTurns,
		size = :size,
		density = :density,
		playerPositions = :playerPositions,
		randomEvents = :randomEvents,
		computerPlayersFormAlliances = :computerPlayersFormAlliances,
		publicPlayerScores = :publicPlayerScores,
		startMode = :startMode,
		year = :year,
		state = :state,
		openPlayerSlots = :openPlayerSlots,
		numPlayers = :numPlayers,
		victoryConditionsConditions = :victoryConditionsConditions,
		victoryConditionsNumCriteriaRequired = :victoryConditionsNumCriteriaRequired,
		victoryConditionsYearsPassed = :victoryConditionsYearsPassed,
		victoryConditionsOwnPlanets = :victoryConditionsOwnPlanets,
		victoryConditionsAttainTechLevel = :victoryConditionsAttainTechLevel,
		victoryConditionsAttainTechLevelNumFields = :victoryConditionsAttainTechLevelNumFields,
		victoryConditionsExceedsScore = :victoryConditionsExceedsScore,
		victoryConditionsExceedsSecondPlaceScore = :victoryConditionsExceedsSecondPlaceScore,
		victoryConditionsProductionCapacity = :victoryConditionsProductionCapacity,
		victoryConditionsOwnCapitalShips = :victoryConditionsOwnCapitalShips,
		victoryConditionsHighestScoreAfterYears = :victoryConditionsHighestScoreAfterYears,
		victorDeclared = :victorDeclared,
		rules = :rules,
		seed = :seed,
		areaX = :areaX,
		areaY = :areaY
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// update a full game
func (c *client) UpdateFullGame(g *game.FullGame) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}

	if err := c.updateGameWithNamedExecer(g.Game, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("update game %w", err)
	}

	for _, player := range g.Players {
		if player.ID == 0 {
			player.GameID = g.ID
			if err := c.createPlayer(player, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create player %w", err)
			}
		}
		if err := c.updateFullPlayerWithTransaction(player, tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("update player %w", err)
		}
	}

	for _, planet := range g.Planets {
		if planet.ID == 0 {
			planet.GameID = g.ID
			if err := c.createPlanet(planet, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create planet %w", err)
			}
			log.Debug().Int64("GameID", planet.GameID).Int64("ID", planet.ID).Msgf("Created planet %s", planet.Name)
		} else if planet.Dirty {
			if err := c.updatePlanet(planet, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update planet %w", err)
			}
			log.Debug().Int64("GameID", planet.GameID).Int64("ID", planet.ID).Msgf("Updated planet %s", planet.Name)
		}
	}

	// save fleets and starbases
	for _, fleet := range append(g.Fleets, g.Starbases...) {
		if fleet.ID == 0 {
			fleet.GameID = g.ID
			if err := c.createFleet(fleet, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create fleet %w", err)
			}
			log.Debug().Int64("GameID", fleet.GameID).Int64("ID", fleet.ID).Msgf("Created fleet %s", fleet.Name)
		} else if fleet.Dirty {
			if err := c.updateFleet(fleet, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update fleet %w", err)
			}
			log.Debug().Int64("GameID", fleet.GameID).Int64("ID", fleet.ID).Msgf("Updated fleet %s", fleet.Name)
		}
	}

	// save wormholes
	for _, wormhole := range g.Wormholes {
		if wormhole.ID == 0 {
			wormhole.GameID = g.ID
			if err := c.createWormhole(wormhole, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create wormhole %w", err)
			}
			log.Debug().Int64("GameID", wormhole.GameID).Int64("ID", wormhole.ID).Msgf("Created wormhole %s", wormhole.Name)
		} else if wormhole.Dirty {
			if err := c.updateWormhole(wormhole, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update wormhole %w", err)
			}
			log.Debug().Int64("GameID", wormhole.GameID).Int64("ID", wormhole.ID).Msgf("Updated wormhole %s", wormhole.Name)
		}
	}

	// save salvages
	for _, salvage := range g.Salvages {
		if salvage.ID == 0 {
			salvage.GameID = g.ID
			if err := c.createSalvage(salvage, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create salvage %w", err)
			}
			log.Debug().Int64("GameID", salvage.GameID).Int64("ID", salvage.ID).Msgf("Created salvage %s", salvage.Name)
		} else if salvage.Dirty {
			if err := c.updateSalvage(salvage, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update salvage %w", err)
			}
			log.Debug().Int64("GameID", salvage.GameID).Int64("ID", salvage.ID).Msgf("Updated salvage %s", salvage.Name)
		}
	}

	// save mineFields
	for _, mineField := range g.MineFields {
		if mineField.ID == 0 {
			mineField.GameID = g.ID
			if err := c.createMineField(mineField, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create mineField %w", err)
			}
			log.Debug().Int64("GameID", mineField.GameID).Int64("ID", mineField.ID).Msgf("Created mineField %s", mineField.Name)
		} else if mineField.Dirty {
			if err := c.updateMineField(mineField, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update mineField %w", err)
			}
			log.Debug().Int64("GameID", mineField.GameID).Int64("ID", mineField.ID).Msgf("Updated mineField %s", mineField.Name)
		}
	}

	// save mineralPackets
	for _, mineralPacket := range g.MineralPackets {
		if mineralPacket.ID == 0 {
			mineralPacket.GameID = g.ID
			if err := c.createMineralPacket(mineralPacket, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create mineralPacket %w", err)
			}
			log.Debug().Int64("GameID", mineralPacket.GameID).Int64("ID", mineralPacket.ID).Msgf("Created mineralPacket %s", mineralPacket.Name)
		} else if mineralPacket.Dirty {
			if err := c.updateMineralPacket(mineralPacket, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update mineralPacket %w", err)
			}
			log.Debug().Int64("GameID", mineralPacket.GameID).Int64("ID", mineralPacket.ID).Msgf("Updated mineralPacket %s", mineralPacket.Name)
		}
	}
	tx.Commit()
	return nil

}

// delete a game by id
func (c *client) DeleteGame(id int64) error {
	if _, err := c.db.Exec("DELETE FROM games WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}