package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

type Game struct {
	ID                                        int64              `json:"id,omitempty"`
	CreatedAt                                 time.Time          `json:"createdAt,omitempty"`
	UpdatedAt                                 time.Time          `json:"updatedAt,omitempty"`
	HostID                                    int64              `json:"hostId,omitempty"`
	Name                                      string             `json:"name,omitempty"`
	State                                     cs.GameState       `json:"state,omitempty"`
	Public                                    bool               `json:"public,omitempty"`
	Hash                                      string             `json:"hash"`
	Size                                      cs.Size            `json:"size,omitempty"`
	Density                                   cs.Density         `json:"density,omitempty"`
	PlayerPositions                           cs.PlayerPositions `json:"playerPositions,omitempty"`
	RandomEvents                              bool               `json:"randomEvents,omitempty"`
	ComputerPlayersFormAlliances              bool               `json:"computerPlayersFormAlliances,omitempty"`
	PublicPlayerScores                        bool               `json:"publicPlayerScores,omitempty"`
	StartMode                                 cs.GameStartMode   `json:"startMode,omitempty"`
	QuickStartTurns                           int                `json:"quickStartTurns,omitempty"`
	OpenPlayerSlots                           uint               `json:"openPlayerSlots,omitempty"`
	NumPlayers                                int                `json:"numPlayers,omitempty"`
	VictoryConditionsConditions               cs.Bitmask         `json:"victoryConditionsConditions,omitempty"`
	VictoryConditionsNumCriteriaRequired      int                `json:"victoryConditionsNumCriteriaRequired,omitempty"`
	VictoryConditionsYearsPassed              int                `json:"victoryConditionsYearsPassed,omitempty"`
	VictoryConditionsOwnPlanets               int                `json:"victoryConditionsOwnPlanets,omitempty"`
	VictoryConditionsAttainTechLevel          int                `json:"victoryConditionsAttainTechLevel,omitempty"`
	VictoryConditionsAttainTechLevelNumFields int                `json:"victoryConditionsAttainTechLevelNumFields,omitempty"`
	VictoryConditionsExceedsScore             int                `json:"victoryConditionsExceedsScore,omitempty"`
	VictoryConditionsExceedsSecondPlaceScore  int                `json:"victoryConditionsExceedsSecondPlaceScore,omitempty"`
	VictoryConditionsProductionCapacity       int                `json:"victoryConditionsProductionCapacity,omitempty"`
	VictoryConditionsOwnCapitalShips          int                `json:"victoryConditionsOwnCapitalShips,omitempty"`
	VictoryConditionsHighestScoreAfterYears   int                `json:"victoryConditionsHighestScoreAfterYears,omitempty"`
	Seed                                      int64              `json:"seed,omitempty"`
	Rules                                     *Rules             `json:"rules,omitempty"`
	AreaX                                     float64            `json:"areaX,omitempty"`
	AreaY                                     float64            `json:"areaY,omitempty"`
	Year                                      int                `json:"year,omitempty"`
	VictorDeclared                            bool               `json:"victorDeclared,omitempty"`
}

// we json serialize these types with custom Scan/Value methods
type Rules cs.Rules

// db serializer to serialize this to JSON
func (item *Rules) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *Rules) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (c *client) GetGames() ([]cs.Game, error) {

	items := []Game{}
	if err := c.db.Select(&items, `SELECT * FROM Games`); err != nil {
		if err == sql.ErrNoRows {
			return []cs.Game{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertGames(items), nil
}

func (c *client) GetGamesForHost(userID int64) ([]cs.GameWithPlayers, error) {
	return c.getGameWithPlayersStatus(`g.hostId = ?`, userID)
}

func (c *client) GetGamesForUser(userID int64) ([]cs.GameWithPlayers, error) {
	return c.getGameWithPlayersStatus(`g.hostId = ? OR g.id in (SELECT gameId from players p WHERE p.userId = ?)`, userID, userID)
}

func (c *client) GetOpenGames() ([]cs.GameWithPlayers, error) {
	return c.getGameWithPlayersStatus(`g.state = ? AND g.openPlayerSlots > 0`, cs.GameStateSetup)

}

func (c *client) GetOpenGamesByHash(hash string) ([]cs.GameWithPlayers, error) {
	return c.getGameWithPlayersStatus(`g.state = ? AND g.openPlayerSlots > 0 AND g.hash = ?`, cs.GameStateSetup, hash)
}

// get a game by id
func (c *client) GetGame(id int64) (*cs.GameWithPlayers, error) {
	games, err := c.getGameWithPlayersStatus("g.id = ?", id)
	if err != nil {
		return nil, err
	}
	if len(games) == 0 {
		return nil, nil
	}

	return &games[0], nil
}

func (c *client) GetGameWithPlayersStatus(gameID int64) (*cs.GameWithPlayers, error) {
	games, err := c.getGameWithPlayersStatus("g.id = ?", gameID)
	if err != nil {
		return nil, err
	}
	if len(games) == 0 {
		return nil, nil
	}

	return &games[0], nil
}

func (c *client) getGameWithPlayersStatus(where string, args ...interface{}) ([]cs.GameWithPlayers, error) {
	type gamePlayersJoin struct {
		Game            `json:"game,omitempty"`
		cs.PlayerStatus `json:"player,omitempty"`
	}

	rows := []gamePlayersJoin{}

	err := c.db.Select(&rows, fmt.Sprintf(`
	SELECT 
		g.id AS 'game.id',
		g.createdAt AS 'game.createdAt',
		g.updatedAt AS 'game.updatedAt',
		g.hostId AS 'game.hostId',
		g.name AS 'game.name',
		g.state AS 'game.state',
		g.public AS 'game.public',
		g.hash AS 'game.hash',
		g.size AS 'game.size',
		g.density AS 'game.density',
		g.playerPositions AS 'game.playerPositions',
		g.randomEvents AS 'game.randomEvents',
		g.computerPlayersFormAlliances AS 'game.computerPlayersFormAlliances',
		g.publicPlayerScores AS 'game.publicPlayerScores',
		g.startMode AS 'game.startMode',
		g.quickStartTurns AS 'game.quickStartTurns',
		g.openPlayerSlots AS 'game.openPlayerSlots',
		g.numPlayers AS 'game.numPlayers',
		g.victoryConditionsConditions AS 'game.victoryConditionsConditions',
		g.victoryConditionsNumCriteriaRequired AS 'game.victoryConditionsNumCriteriaRequired',
		g.victoryConditionsYearsPassed AS 'game.victoryConditionsYearsPassed',
		g.victoryConditionsOwnPlanets AS 'game.victoryConditionsOwnPlanets',
		g.victoryConditionsAttainTechLevel AS 'game.victoryConditionsAttainTechLevel',
		g.victoryConditionsAttainTechLevelNumFields AS 'game.victoryConditionsAttainTechLevelNumFields',
		g.victoryConditionsExceedsScore AS 'game.victoryConditionsExceedsScore',
		g.victoryConditionsExceedsSecondPlaceScore AS 'game.victoryConditionsExceedsSecondPlaceScore',
		g.victoryConditionsProductionCapacity AS 'game.victoryConditionsProductionCapacity',
		g.victoryConditionsOwnCapitalShips AS 'game.victoryConditionsOwnCapitalShips',
		g.victoryConditionsHighestScoreAfterYears AS 'game.victoryConditionsHighestScoreAfterYears',
		g.seed AS 'game.seed',
		g.rules AS 'game.rules',
		g.areaX AS 'game.areaX',
		g.areaY AS 'game.areaY',
		g.year AS 'game.year',
		g.victorDeclared AS 'game.victorDeclared',
		
		p.updatedAt AS 'player.updatedAt',
		COALESCE(p.userId, 0) AS 'player.userId',
		COALESCE(p.name, '') AS 'player.name',
		COALESCE(p.num, 0) AS 'player.num',
		COALESCE(p.ready, 0) AS 'player.ready',
		COALESCE(p.aiControlled, 0) AS 'player.aiControlled',
		COALESCE(p.submittedTurn, 0) AS 'player.submittedTurn',
		COALESCE(p.color, '') AS 'player.color',
		COALESCE(p.victor, 0) AS 'player.victor'

	FROM games g
	LEFT JOIN players p
		ON g.id = p.gameId
	WHERE %s
`, where), args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []cs.GameWithPlayers{}, nil
		}
		return nil, err
	}

	// check if we have a game
	if len(rows) == 0 {
		return []cs.GameWithPlayers{}, nil
	}

	// join results give a row per item, so if we have 2 games
	// one with 2 players, one with 3, we'll end up with 5 rows
	// row 0 - game1, player 1
	// row 1 - game1, player 2
	// row 2 - game2, player 1
	// row 3 - game2, player 2
	// row 4 - game2, player 3
	games := []cs.GameWithPlayers{}
	var item Game
	var game *cs.GameWithPlayers
	for _, row := range rows {

		if row.ID != item.ID {
			// convert this row into a game
			item = row.Game
			g := c.converter.ConvertGame(item)
			games = append(games, cs.GameWithPlayers{Game: g, Players: []cs.PlayerStatus{}})
			game = &games[len(games)-1]
		}

		if row.PlayerStatus.Num != 0 {
			game.Players = append(game.Players, row.PlayerStatus)
		}
	}

	return games, nil
}

// get a game by id
func (c *client) GetFullGame(id int64) (*cs.FullGame, error) {
	item := Game{}
	if err := c.db.Get(&item, "SELECT * FROM games WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	game := c.converter.ConvertGame(item)

	players, err := c.getPlayersForGame(game.ID)
	if err != nil {
		return nil, fmt.Errorf("load players for game %w", err)
	}

	universe := cs.NewUniverse(&game.Rules)

	planets, err := c.getPlanetsForGame(game.ID)
	if err != nil {
		return nil, fmt.Errorf("load planets for game %w", err)
	}
	universe.Planets = planets

	// load fleets and starbases
	fleets, err := c.getFleetsForGame(game.ID)
	if err != nil {
		return nil, fmt.Errorf("load fleets for game %w", err)
	}
	// pre-instantiate the fleets/starbases arrays (make it a little bigger than necessary)
	universe.Fleets = make([]*cs.Fleet, 0, len(fleets))
	universe.Starbases = make([]*cs.Fleet, 0, len(planets))
	for i := range fleets {
		fleet := fleets[i]
		if fleet.Starbase {
			universe.Starbases = append(universe.Starbases, fleet)
		} else {
			universe.Fleets = append(universe.Fleets, fleet)
		}
	}

	wormholes, err := c.getWormholesForGame(game.ID)
	if err != nil {
		return nil, fmt.Errorf("load wormholes for game %w", err)
	}
	universe.Wormholes = wormholes

	salvages, err := c.getSalvagesForGame(game.ID)
	if err != nil {
		return nil, fmt.Errorf("load salvages for game %w", err)
	}
	universe.Salvages = salvages

	mineFields, err := c.getMineFieldsForGame(game.ID)
	if err != nil {
		return nil, fmt.Errorf("load mineFields for game %w", err)
	}
	universe.MineFields = mineFields

	mineralPackets, err := c.getMineralPacketsForGame(game.ID)
	if err != nil {
		return nil, fmt.Errorf("load mineralPackets for game %w", err)
	}
	universe.MineralPackets = mineralPackets

	// load a tech store if this game has a separate one
	techStore := &cs.StaticTechStore
	if game.Rules.TechsID != 0 {
		techStore, err = c.GetTechStore(game.Rules.TechsID)
		if err != nil {
			return nil, err
		}
	}

	game.Rules.WithTechStore(techStore)

	// init the random generator after load
	(&game.Rules).ResetSeed(game.Seed)

	fg := cs.FullGame{
		Game:      &game,
		Universe:  &universe,
		TechStore: techStore,
		Players:   players,
	}

	return &fg, nil
}

// create a new game
func (c *client) CreateGame(game *cs.Game) error {

	item := c.converter.ConvertGameGame(game)
	result, err := c.db.NamedExec(`
	INSERT INTO games (
		createdAt,
		updatedAt,
		hostId,
		name,
		state,
		public,
		hash,
		size,
		density,
		playerPositions,
		randomEvents,
		computerPlayersFormAlliances,
		publicPlayerScores,
		startMode,
		quickStartTurns,
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
		seed,
		rules,
		areaX,
		areaY,
		year,
		victorDeclared
	)
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:hostId,
		:name,
		:state,
		:public,
		:hash,
		:size,
		:density,
		:playerPositions,
		:randomEvents,
		:computerPlayersFormAlliances,
		:publicPlayerScores,
		:startMode,
		:quickStartTurns,
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
		:seed,
		:rules,
		:areaX,
		:areaY,
		:year,
		:victorDeclared
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
func (c *client) UpdateGame(game *cs.Game) error {
	return c.updateGameWithNamedExecer(game, c.db)
}

// update a game inside a transaction
func (c *client) updateGameWithNamedExecer(game *cs.Game, tx SQLExecer) error {

	item := c.converter.ConvertGameGame(game)

	if _, err := tx.NamedExec(`
	UPDATE games SET
		updatedAt = CURRENT_TIMESTAMP,
		hostId = :hostId,
		name = :name,
		state = :state,
		public = :public,
		hash = :hash,
		size = :size,
		density = :density,
		playerPositions = :playerPositions,
		randomEvents = :randomEvents,
		computerPlayersFormAlliances = :computerPlayersFormAlliances,
		publicPlayerScores = :publicPlayerScores,
		startMode = :startMode,
		quickStartTurns = :quickStartTurns,
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
		seed = :seed,
		rules = :rules,
		areaX = :areaX,
		areaY = :areaY,
		year = :year,
		victorDeclared = :victorDeclared

	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// update a full game
func (c *client) UpdateFullGame(fullGame *cs.FullGame) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}

	if err := c.updateGameWithNamedExecer(fullGame.Game, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("update game %w", err)
	}

	for _, player := range fullGame.Players {
		if player.ID == 0 {
			player.GameID = fullGame.ID
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

	for _, planet := range fullGame.Planets {
		if planet.ID == 0 {
			planet.GameID = fullGame.ID
			if err := c.createPlanet(planet, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create planet %w", err)
			}
			// log.Debug().Int64("GameID", planet.GameID).Int64("ID", planet.ID).Msgf("Created planet %s", planet.Name)
		} else if planet.Dirty {
			if err := c.updatePlanet(planet, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update planet %w", err)
			}
			// log.Debug().Int64("GameID", planet.GameID).Int64("ID", planet.ID).Msgf("Updated planet %s", planet.Name)
		}
	}

	// save fleets and starbases
	remainingFleets := make([]*cs.Fleet, 0, len(fullGame.Fleets))
	for _, fleet := range append(fullGame.Fleets, fullGame.Starbases...) {
		if fleet.ID == 0 {
			fleet.GameID = fullGame.ID
			if err := c.createFleet(fleet, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create fleet %w", err)
			}
			remainingFleets = append(remainingFleets, fleet)
			// log.Debug().Int64("GameID", fleet.GameID).Int64("ID", fleet.ID).Msgf("Created fleet %s", fleet.Name)
		} else if fleet.Delete {
			if err := c.deleteFleet(fleet.ID, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("delete fleet %w", err)
			}
			// log.Debug().Int64("GameID", fleet.GameID).Int64("ID", fleet.ID).Msgf("Deleted fleet %s", fleet.Name)
		} else if fleet.Dirty {
			if err := c.updateFleet(fleet, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update fleet %w", err)
			}
			remainingFleets = append(remainingFleets, fleet)
			// log.Debug().Int64("GameID", fleet.GameID).Int64("ID", fleet.ID).Msgf("Updated fleet %s", fleet.Name)
		}
	}
	fullGame.Fleets = remainingFleets

	// save wormholes
	for _, wormhole := range fullGame.Wormholes {
		if wormhole.ID == 0 {
			wormhole.GameID = fullGame.ID
			if err := c.createWormhole(wormhole, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create wormhole %w", err)
			}
			// log.Debug().Int64("GameID", wormhole.GameID).Int64("ID", wormhole.ID).Msgf("Created wormhole %v", wormhole)
		} else if wormhole.Delete {
			if err := c.deleteWormhole(wormhole.ID, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("delete wormhole %w", err)
			}
			// log.Debug().Int64("GameID", wormhole.GameID).Int64("ID", wormhole.ID).Msgf("Deleted wormhole %s", wormhole.Name)
		} else if wormhole.Dirty {
			if err := c.updateWormhole(wormhole, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update wormhole %w", err)
			}
			// log.Debug().Int64("GameID", wormhole.GameID).Int64("ID", wormhole.ID).Msgf("Updated wormhole %v", wormhole)
		}
	}

	// save salvages
	for _, salvage := range fullGame.Salvages {
		if salvage.ID == 0 {
			salvage.GameID = fullGame.ID
			if err := c.createSalvage(salvage, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create salvage %w", err)
			}
			// log.Debug().Int64("GameID", salvage.GameID).Int64("ID", salvage.ID).Msgf("Created salvage %s", salvage.Name)
		} else if salvage.Delete {
			if err := c.deleteSalvage(salvage.ID, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("delete salvage %w", err)
			}
			// log.Debug().Int64("GameID", salvage.GameID).Int64("ID", salvage.ID).Msgf("Deleted salvage %s", salvage.Name)
		} else if salvage.Dirty {
			if err := c.updateSalvage(salvage, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update salvage %w", err)
			}
			// log.Debug().Int64("GameID", salvage.GameID).Int64("ID", salvage.ID).Msgf("Updated salvage %s", salvage.Name)
		}
	}

	// save mineFields
	for _, mineField := range fullGame.MineFields {
		if mineField.ID == 0 {
			mineField.GameID = fullGame.ID
			if err := c.createMineField(mineField, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create mineField %w", err)
			}
			// log.Debug().Int64("GameID", mineField.GameID).Int64("ID", mineField.ID).Msgf("Created mineField %s", mineField.Name)
		} else if mineField.Delete {
			if err := c.deleteMineField(mineField.ID, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("delete mineField %w", err)
			}
			// log.Debug().Int64("GameID", mineField.GameID).Int64("ID", mineField.ID).Msgf("Deleted mineField %s", mineField.Name)
		} else if mineField.Dirty {
			if err := c.updateMineField(mineField, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update mineField %w", err)
			}
			// log.Debug().Int64("GameID", mineField.GameID).Int64("ID", mineField.ID).Msgf("Updated mineField %s", mineField.Name)
		}
	}

	// save mineralPackets
	for _, mineralPacket := range fullGame.MineralPackets {
		if mineralPacket.ID == 0 {
			mineralPacket.GameID = fullGame.ID
			if err := c.createMineralPacket(mineralPacket, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create mineralPacket %w", err)
			}
			// log.Debug().Int64("GameID", mineralPacket.GameID).Int64("ID", mineralPacket.ID).Msgf("Created mineralPacket %s", mineralPacket.Name)
		} else if mineralPacket.Delete {
			if err := c.deleteMineralPacket(mineralPacket.ID, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("delete mineralPacket %w", err)
			}
			// log.Debug().Int64("GameID", mineralPacket.GameID).Int64("ID", mineralPacket.ID).Msgf("Deleted mineralPacket %s", mineralPacket.Name)
		} else if mineralPacket.Dirty {
			if err := c.updateMineralPacket(mineralPacket, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update mineralPacket %w", err)
			}
			// log.Debug().Int64("GameID", mineralPacket.GameID).Int64("ID", mineralPacket.ID).Msgf("Updated mineralPacket %s", mineralPacket.Name)
		}
	}

	// save mysteryTraders
	for _, mysteryTrader := range fullGame.MysteryTraders {
		if mysteryTrader.ID == 0 {
			mysteryTrader.GameID = fullGame.ID
			if err := c.createMysteryTrader(mysteryTrader, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create mysteryTrader %w", err)
			}
			// log.Debug().Int64("GameID", mysteryTrader.GameID).Int64("ID", mysteryTrader.ID).Msgf("Created mysteryTrader %s", mysteryTrader.Name)
		} else if mysteryTrader.Delete {
			if err := c.deleteMysteryTrader(mysteryTrader.ID, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("delete mysteryTrader %w", err)
			}
			// log.Debug().Int64("GameID", mysteryTrader.GameID).Int64("ID", mysteryTrader.ID).Msgf("Deleted mysteryTrader %s", mysteryTrader.Name)
		} else if mysteryTrader.Dirty {
			if err := c.updateMysteryTrader(mysteryTrader, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update mysteryTrader %w", err)
			}
			// log.Debug().Int64("GameID", mysteryTrader.GameID).Int64("ID", mysteryTrader.ID).Msgf("Updated mysteryTrader %s", mysteryTrader.Name)
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

// delete a game by id
func (c *client) DeleteGame(id int64) error {
	if _, err := c.db.Exec("DELETE FROM games WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}
