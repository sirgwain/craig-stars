package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirgwain/craig-stars/cs"
)

type Player struct {
	ID                           int64                `json:"id,omitempty"`
	CreatedAt                    time.Time            `json:"createdAt,omitempty"`
	UpdatedAt                    time.Time            `json:"updatedAt,omitempty"`
	GameID                       int64                `json:"gameId,omitempty"`
	UserID                       int64                `json:"userId,omitempty"`
	Name                         string               `json:"name,omitempty"`
	Num                          int                  `json:"num,omitempty"`
	Ready                        bool                 `json:"ready,omitempty"`
	AIControlled                 bool                 `json:"aiControlled,omitempty"`
	SubmittedTurn                bool                 `json:"submittedTurn,omitempty"`
	Color                        string               `json:"color,omitempty"`
	DefaultHullSet               int                  `json:"defaultHullSet,omitempty"`
	TechLevelsEnergy             int                  `json:"techLevelsEnergy,omitempty"`
	TechLevelsWeapons            int                  `json:"techLevelsWeapons,omitempty"`
	TechLevelsPropulsion         int                  `json:"techLevelsPropulsion,omitempty"`
	TechLevelsConstruction       int                  `json:"techLevelsConstruction,omitempty"`
	TechLevelsElectronics        int                  `json:"techLevelsElectronics,omitempty"`
	TechLevelsBiotechnology      int                  `json:"techLevelsBiotechnology,omitempty"`
	TechLevelsSpentEnergy        int                  `json:"techLevelsSpentEnergy,omitempty"`
	TechLevelsSpentWeapons       int                  `json:"techLevelsSpentWeapons,omitempty"`
	TechLevelsSpentPropulsion    int                  `json:"techLevelsSpentPropulsion,omitempty"`
	TechLevelsSpentConstruction  int                  `json:"techLevelsSpentConstruction,omitempty"`
	TechLevelsSpentElectronics   int                  `json:"techLevelsSpentElectronics,omitempty"`
	TechLevelsSpentBiotechnology int                  `json:"techLevelsSpentBiotechnology,omitempty"`
	ResearchAmount               int                  `json:"researchAmount,omitempty"`
	ResearchSpentLastYear        int                  `json:"researchSpentLastYear,omitempty"`
	NextResearchField            cs.NextResearchField `json:"nextResearchField,omitempty"`
	Researching                  cs.TechField         `json:"researching,omitempty"`
	BattlePlans                  *BattlePlans         `json:"battlePlans,omitempty"`
	ProductionPlans              *ProductionPlans     `json:"productionPlans,omitempty"`
	TransportPlans               *TransportPlans      `json:"transportPlans,omitempty"`
	Relations                    *PlayerRelationships `json:"relations,omitempty"`
	Messages                     *PlayerMessages      `json:"messages,omitempty"`
	Battles                      *BattleRecords       `json:"battles,omitempty"`
	PlayerIntels                 *PlayerIntels        `json:"playerIntels,omitempty"`
	PlanetIntels                 *PlanetIntels        `json:"planetIntels,omitempty"`
	FleetIntels                  *FleetIntels         `json:"fleetIntels,omitempty"`
	ShipDesignIntels             *ShipDesignIntels    `json:"shipDesignIntels,omitempty"`
	MineralPacketIntels          *MineralPacketIntels `json:"mineralPacketIntels,omitempty"`
	MineFieldIntels              *MineFieldIntels     `json:"mineFieldIntels,omitempty"`
	WormholeIntels               *WormholeIntels      `json:"wormholeIntels,omitempty"`
	Race                         *PlayerRace          `json:"race,omitempty"`
	Stats                        *PlayerStats         `json:"stats,omitempty"`
	Spec                         *PlayerSpec          `json:"spec,omitempty"`
}

// we json serialize these types with custom Scan/Value methods
type BattlePlans []cs.BattlePlan
type ProductionPlans []cs.ProductionPlan
type TransportPlans []cs.TransportPlan
type PlayerRelationships []cs.PlayerRelationship
type PlayerMessages []cs.PlayerMessage
type BattleRecords []cs.BattleRecord
type PlayerIntels []cs.PlayerIntel
type PlanetIntels []cs.PlanetIntel
type FleetIntels []cs.FleetIntel
type ShipDesignIntels []cs.ShipDesignIntel
type MineralPacketIntels []cs.MineralPacketIntel
type MineFieldIntels []cs.MineFieldIntel
type WormholeIntels []cs.WormholeIntel
type PlayerRace cs.Race
type PlayerSpec cs.PlayerSpec
type PlayerStats cs.PlayerStats

// db serializer to serialize this to JSON
func (item *BattlePlans) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *BattlePlans) Scan(src interface{}) error {
	return scanJSON(src, &item)
}

// db serializer to serialize this to JSON
func (item *ProductionPlans) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *ProductionPlans) Scan(src interface{}) error {
	return scanJSON(src, &item)
}

// db serializer to serialize this to JSON
func (item *TransportPlans) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *TransportPlans) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// db serializer to serialize this to JSON
func (item *PlayerRace) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *PlayerRace) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// db serializer to serialize this to JSON
func (item *PlayerSpec) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *PlayerSpec) Scan(src interface{}) error {
	return scanJSON(src, item)

}

// db serializer to serialize this to JSON
func (item *PlayerStats) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *PlayerStats) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *PlayerRelationships) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *PlayerRelationships) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *PlayerMessages) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *PlayerMessages) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *BattleRecords) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *BattleRecords) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *PlayerIntels) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *PlayerIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *PlanetIntels) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *PlanetIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *FleetIntels) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *FleetIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *ShipDesignIntels) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *ShipDesignIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *MineralPacketIntels) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *MineralPacketIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *MineFieldIntels) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *MineFieldIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *WormholeIntels) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *WormholeIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (c *client) GetPlayers() ([]cs.Player, error) {

	items := []Player{}
	if err := c.db.Select(&items, `SELECT * FROM players`); err != nil {
		if err == sql.ErrNoRows {
			return []cs.Player{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertPlayers(items), nil
}

func (c *client) GetPlayersForUser(userID int64) ([]cs.Player, error) {

	items := []Player{}
	if err := c.db.Select(&items, `SELECT * FROM players WHERE userId = ?`, userID); err != nil {
		if err == sql.ErrNoRows {
			return []cs.Player{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertPlayers(items), nil
}

// get all the players for a game, with data loaded
func (c *client) getPlayersForGame(gameID int64) ([]*cs.Player, error) {

	items := []Player{}
	if err := c.db.Select(&items, `SELECT * FROM players WHERE gameId = ?`, gameID); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.Player{}, nil
		}
		return nil, err
	}

	players := make([]*cs.Player, len(items))
	for i := range items {
		player := c.converter.ConvertPlayer(items[i])
		players[i] = &player

		designs, err := c.GetShipDesignsForPlayer(gameID, player.Num)
		if err != nil {
			return nil, fmt.Errorf("get designs for player %w", err)
		}
		player.Designs = designs
	}

	return players, nil
}

// get all the players for a game, with data loaded
func (c *client) GetPlayerStatusesForGame(gameID int64) ([]*cs.Player, error) {

	items := []Player{}
	if err := c.db.Select(&items, `SELECT id, createdAt, updatedAt, gameId, userId, name, num, ready, aiControlled, submittedTurn, color FROM players WHERE gameId = ? ORDER BY num`, gameID); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.Player{}, nil
		}
		return nil, err
	}

	players := make([]*cs.Player, len(items))
	for i := range items {
		player := c.converter.ConvertPlayer(items[i])
		players[i] = &player
	}

	return players, nil
}

// get a player by id
func (c *client) GetPlayer(id int64) (*cs.Player, error) {
	item := Player{}
	if err := c.db.Get(&item, "SELECT * FROM players WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	player := c.converter.ConvertPlayer(item)
	return &player, nil
}

func (c *client) GetPlayerForGame(gameID, userID int64) (*cs.Player, error) {
	item := Player{}
	if err := c.db.Get(&item, "SELECT * FROM players WHERE gameId = ? AND userId = ?", gameID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	player := c.converter.ConvertPlayer(item)

	// get designs
	designs, err := c.GetShipDesignsForPlayer(gameID, player.Num)
	if err != nil {
		return nil, fmt.Errorf("get player designs %w", err)
	}
	player.Designs = designs

	return &player, nil
}

func (c *client) GetLightPlayerForGame(gameID, userID int64) (*cs.Player, error) {
	item := Player{}
	if err := c.db.Get(&item, `
	SELECT 
	id,
	createdAt,
	updatedAt,
	gameId,
	userId,
	name,
	num,
	ready,
	aiControlled,
	submittedTurn,
	color,
	defaultHullSet,
	race,
	techLevelsEnergy,
	techLevelsWeapons,
	techLevelsPropulsion,
	techLevelsConstruction,
	techLevelsElectronics,
	techLevelsBiotechnology	
	techLevelsSpentEnergy,
	techLevelsSpentWeapons,
	techLevelsSpentPropulsion,
	techLevelsSpentConstruction,
	techLevelsSpentElectronics,
	techLevelsSpentBiotechnology	
	researchSpentLastYear,
	researchAmount,
	nextResearchField,
	researching,
	battlePlans,
	productionPlans,
	transportPlans,
	spec,
	stats
	FROM players 
	WHERE gameId = ? AND userId = ?`, gameID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	player := c.converter.ConvertPlayer(item)
	return &player, nil
}

// get a full player by id with all dependencies loaded
func (c *client) GetFullPlayerForGame(gameID, userID int64) (*cs.FullPlayer, error) {
	player := cs.FullPlayer{}

	item := Player{}
	if err := c.db.Get(&item, "SELECT * FROM players WHERE gameId = ? AND userId = ?", gameID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// load the player component from the DB
	player.Player = c.converter.ConvertPlayer(item)

	designs, err := c.GetShipDesignsForPlayer(gameID, player.Num)
	if err != nil {
		return nil, fmt.Errorf("get player designs %w", err)
	}
	player.Designs = designs

	planets, err := c.GetPlanetsForPlayer(player.GameID, player.Num)
	if err != nil {
		return nil, fmt.Errorf("get player planets %w", err)
	}
	player.Planets = planets

	fleets, err := c.GetFleetsForPlayer(player.GameID, player.Num)
	if err != nil {
		return nil, fmt.Errorf("get player fleets %w", err)
	}
	// pre-instantiate the fleets/starbases arrays (make it a little bigger than necessary)
	player.Fleets = make([]*cs.Fleet, 0, len(fleets))
	player.Starbases = make([]*cs.Fleet, 0, len(planets))
	for i := range fleets {
		fleet := fleets[i]
		if fleet.Starbase {
			player.Starbases = append(player.Starbases, fleet)
		} else {
			player.Fleets = append(player.Fleets, fleet)
		}
	}

	return &player, nil
}

func (c *client) GetPlayerMapObjects(gameID, userID int64) (*cs.PlayerMapObjects, error) {
	mapObjects := cs.PlayerMapObjects{}
	var num int
	if err := c.db.Get(&num, "SELECT num FROM players WHERE gameId = ? AND userId = ?", gameID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	planets, err := c.GetPlanetsForPlayer(gameID, num)
	if err != nil {
		return nil, fmt.Errorf("get player planets %w", err)
	}
	mapObjects.Planets = planets

	fleets, err := c.GetFleetsForPlayer(gameID, num)
	if err != nil {
		return nil, fmt.Errorf("get player fleets %w", err)
	}
	// pre-instantiate the fleets/starbases arrays (make it a little bigger than necessary)
	mapObjects.Fleets = make([]*cs.Fleet, 0, len(fleets))
	mapObjects.Starbases = make([]*cs.Fleet, 0, len(planets))
	for i := range fleets {
		fleet := fleets[i]
		if fleet.Starbase {
			mapObjects.Starbases = append(mapObjects.Starbases, fleet)
		} else {
			mapObjects.Fleets = append(mapObjects.Fleets, fleet)
		}
	}

	return &mapObjects, nil
}

// get a player with designs loaded
func (c *client) GetPlayerWithDesignsForGame(gameID, userID int64) (*cs.Player, error) {
	player := cs.Player{}

	item := Player{}
	if err := c.db.Get(&item, "SELECT * FROM players WHERE gameId = ? AND userId = ?", gameID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// load the player component from the DB
	player = c.converter.ConvertPlayer(item)

	designs, err := c.GetShipDesignsForPlayer(gameID, player.Num)
	if err != nil {
		return nil, fmt.Errorf("get player designs %w", err)
	}
	player.Designs = designs

	return &player, nil
}

func (c *client) CreatePlayer(player *cs.Player) error {
	return c.createPlayer(player, c.db)
}

// create a new game
func (c *client) createPlayer(player *cs.Player, tx SQLExecer) error {
	item := c.converter.ConvertGamePlayer(player)
	result, err := tx.NamedExec(`
	INSERT INTO players (
		createdAt,
		updatedAt,
		gameId,
		userId,
		name,
		num,
		ready,
		aiControlled,
		submittedTurn,
		color,
		defaultHullSet,
		techLevelsEnergy,
		techLevelsWeapons,
		techLevelsPropulsion,
		techLevelsConstruction,
		techLevelsElectronics,
		techLevelsBiotechnology,
		techLevelsSpentEnergy,
		techLevelsSpentWeapons,
		techLevelsSpentPropulsion,
		techLevelsSpentConstruction,
		techLevelsSpentElectronics,
		techLevelsSpentBiotechnology,
		researchAmount,
		researchSpentLastYear,
		nextResearchField,
		researching,
		battlePlans,
		productionPlans,
		transportPlans,
		relations,
		messages,
		battles,
		playerIntels,
		planetIntels,
		fleetIntels,
		shipDesignIntels,
		mineralPacketIntels,
		mineFieldIntels,
		wormholeIntels,
		race,
		stats,
		spec
	)
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:gameId,
		:userId,
		:name,
		:num,
		:ready,
		:aiControlled,
		:submittedTurn,
		:color,
		:defaultHullSet,
		:techLevelsEnergy,
		:techLevelsWeapons,
		:techLevelsPropulsion,
		:techLevelsConstruction,
		:techLevelsElectronics,
		:techLevelsBiotechnology,
		:techLevelsSpentEnergy,
		:techLevelsSpentWeapons,
		:techLevelsSpentPropulsion,
		:techLevelsSpentConstruction,
		:techLevelsSpentElectronics,
		:techLevelsSpentBiotechnology,
		:researchAmount,
		:researchSpentLastYear,
		:nextResearchField,
		:researching,
		:battlePlans,
		:productionPlans,
		:transportPlans,
		:relations,
		:messages,
		:battles,
		:playerIntels,
		:planetIntels,
		:fleetIntels,
		:shipDesignIntels,
		:mineralPacketIntels,
		:mineFieldIntels,
		:wormholeIntels,
		:race,
		:stats,
		:spec
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

	player.ID = id

	return nil
}

// update an existing player
func (c *client) UpdatePlayer(player *cs.Player) error {
	return c.updatePlayerWithNamedExecer(player, c.db)
}

// update an existing player's lightweight fields
func (c *client) UpdateLightPlayer(player *cs.Player) error {
	item := c.converter.ConvertGamePlayer(player)

	if _, err := c.db.NamedExec(`
	UPDATE players SET
		updatedAt = CURRENT_TIMESTAMP,
		name = :name,
		num = :num,
		ready = :ready,
		aiControlled = :aiControlled,
		submittedTurn = :submittedTurn,
		color = :color,
		defaultHullSet = :defaultHullSet,
		researchAmount = :researchAmount,
		nextResearchField = :nextResearchField,
		researching = :researching,
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// helper to update a player using a transaction or DB
func (c *client) updatePlayerWithNamedExecer(player *cs.Player, tx SQLExecer) error {
	item := c.converter.ConvertGamePlayer(player)

	if _, err := tx.NamedExec(`
	UPDATE players SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		userId = :userId,
		name = :name,
		num = :num,
		ready = :ready,
		aiControlled = :aiControlled,
		submittedTurn = :submittedTurn,
		color = :color,
		defaultHullSet = :defaultHullSet,
		techLevelsEnergy = :techLevelsEnergy,
		techLevelsWeapons = :techLevelsWeapons,
		techLevelsPropulsion = :techLevelsPropulsion,
		techLevelsConstruction = :techLevelsConstruction,
		techLevelsElectronics = :techLevelsElectronics,
		techLevelsBiotechnology = :techLevelsBiotechnology,
		techLevelsSpentEnergy = :techLevelsSpentEnergy,
		techLevelsSpentWeapons = :techLevelsSpentWeapons,
		techLevelsSpentPropulsion = :techLevelsSpentPropulsion,
		techLevelsSpentConstruction = :techLevelsSpentConstruction,
		techLevelsSpentElectronics = :techLevelsSpentElectronics,
		techLevelsSpentBiotechnology = :techLevelsSpentBiotechnology,
		researchAmount = :researchAmount,
		researchSpentLastYear = :researchSpentLastYear,
		nextResearchField = :nextResearchField,
		researching = :researching,
		battlePlans = :battlePlans,
		productionPlans = :productionPlans,
		transportPlans = :transportPlans,
		relations = :relations,
		messages = :messages,
		battles = :battles,
		playerIntels = :playerIntels,
		planetIntels = :planetIntels,
		fleetIntels = :fleetIntels,
		shipDesignIntels = :shipDesignIntels,
		mineralPacketIntels = :mineralPacketIntels,
		mineFieldIntels = :mineFieldIntels,
		wormholeIntels = :wormholeIntels,
		race = :race,
		stats = :stats,
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

func (c *client) updateFullPlayer(player *cs.Player) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}

	if err := c.updateFullPlayerWithTransaction(player, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// update an existing player
func (c *client) updateFullPlayerWithTransaction(player *cs.Player, tx *sqlx.Tx) error {

	if err := c.updatePlayerWithNamedExecer(player, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("update player %w", err)
	}

	for i := range player.Designs {
		design := player.Designs[i]
		if design.ID == 0 {
			design.GameID = player.GameID
			if err := c.createShipDesign(design, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("create design %w", err)
			}
		} else if design.Dirty {
			if err := c.updateShipDesign(design, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("update design %w", err)
			}
		}
	}

	return nil
}

// delete a player by id
func (c *client) DeletePlayer(id int64) error {
	if _, err := c.db.Exec("DELETE FROM players WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}
