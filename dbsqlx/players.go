package dbsqlx

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirgwain/craig-stars/game"
)

type Player struct {
	ID                           int64                  `json:"id,omitempty"`
	CreatedAt                    time.Time              `json:"createdAt,omitempty"`
	UpdatedAt                    time.Time              `json:"updatedAt,omitempty"`
	GameID                       int64                  `json:"gameId,omitempty"`
	UserID                       int64                  `json:"userId,omitempty"`
	Name                         string                 `json:"name,omitempty"`
	Num                          int                    `json:"num,omitempty"`
	Ready                        bool                   `json:"ready,omitempty"`
	AIControlled                 bool                   `json:"aiControlled,omitempty"`
	SubmittedTurn                bool                   `json:"submittedTurn,omitempty"`
	Color                        string                 `json:"color,omitempty"`
	DefaultHullSet               int                    `json:"defaultHullSet,omitempty"`
	TechLevelsEnergy             int                    `json:"techLevelsEnergy,omitempty"`
	TechLevelsWeapons            int                    `json:"techLevelsWeapons,omitempty"`
	TechLevelsPropulsion         int                    `json:"techLevelsPropulsion,omitempty"`
	TechLevelsConstruction       int                    `json:"techLevelsConstruction,omitempty"`
	TechLevelsElectronics        int                    `json:"techLevelsElectronics,omitempty"`
	TechLevelsBiotechnology      int                    `json:"techLevelsBiotechnology,omitempty"`
	TechLevelsSpentEnergy        int                    `json:"techLevelsSpentEnergy,omitempty"`
	TechLevelsSpentWeapons       int                    `json:"techLevelsSpentWeapons,omitempty"`
	TechLevelsSpentPropulsion    int                    `json:"techLevelsSpentPropulsion,omitempty"`
	TechLevelsSpentConstruction  int                    `json:"techLevelsSpentConstruction,omitempty"`
	TechLevelsSpentElectronics   int                    `json:"techLevelsSpentElectronics,omitempty"`
	TechLevelsSpentBiotechnology int                    `json:"techLevelsSpentBiotechnology,omitempty"`
	ResearchAmount               int                    `json:"researchAmount,omitempty"`
	ResearchSpentLastYear        int                    `json:"researchSpentLastYear,omitempty"`
	NextResearchField            game.NextResearchField `json:"nextResearchField,omitempty"`
	Researching                  game.TechField         `json:"researching,omitempty"`
	BattlePlans                  *BattlePlans           `json:"battlePlans,omitempty"`
	ProductionPlans              *ProductionPlans       `json:"productionPlans,omitempty"`
	TransportPlans               *TransportPlans        `json:"transportPlans,omitempty"`
	Messages                     *PlayerMessages        `json:"messages,omitempty"`
	PlanetIntels                 *PlanetIntels          `json:"planetIntels,omitempty"`
	FleetIntels                  *FleetIntels           `json:"fleetIntels,omitempty"`
	ShipDesignIntels             *ShipDesignIntels      `json:"shipDesignIntels,omitempty"`
	MineralPacketIntels          *MineralPacketIntels   `json:"mineralPacketIntels,omitempty"`
	MineFieldIntels              *MineFieldIntels       `json:"mineFieldIntels,omitempty"`
	Race                         *PlayerRace            `json:"race,omitempty"`
	Stats                        *PlayerStats           `json:"stats,omitempty"`
	Spec                         *PlayerSpec            `json:"spec,omitempty"`
}

// we json serialize these types with custom Scan/Value methods
type BattlePlans []game.BattlePlan
type ProductionPlans []game.ProductionPlan
type TransportPlans []game.TransportPlan
type PlayerMessages []game.PlayerMessage
type PlanetIntels []game.PlanetIntel
type FleetIntels []game.FleetIntel
type ShipDesignIntels []game.ShipDesignIntel
type MineralPacketIntels []game.MineralPacketIntel
type MineFieldIntels []game.MineFieldIntel
type PlayerRace game.Race
type PlayerSpec game.PlayerSpec
type PlayerStats game.PlayerStats

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

func (item *PlayerMessages) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *PlayerMessages) Scan(src interface{}) error {
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

func (c *client) GetPlayers() ([]game.Player, error) {

	items := []Player{}
	if err := c.db.Select(&items, `SELECT * FROM players`); err != nil {
		if err == sql.ErrNoRows {
			return []game.Player{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertPlayers(items), nil
}

func (c *client) GetPlayersForUser(userID int64) ([]game.Player, error) {

	items := []Player{}
	if err := c.db.Select(&items, `SELECT * FROM players WHERE userId = ?`, userID); err != nil {
		if err == sql.ErrNoRows {
			return []game.Player{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertPlayers(items), nil
}

// get all the players for a game, with data loaded
func (c *client) getPlayersForGame(gameID int64) ([]*game.Player, error) {

	items := []Player{}
	if err := c.db.Select(&items, `SELECT * FROM players WHERE gameId = ?`, gameID); err != nil {
		if err == sql.ErrNoRows {
			return []*game.Player{}, nil
		}
		return nil, err
	}

	players := make([]*game.Player, len(items))
	for i := range items {
		player := c.converter.ConvertPlayer(items[i])
		players[i] = &player

		designs, err := c.GetShipDesignsForPlayer(player.ID)
		if err != nil {
			return nil, fmt.Errorf("get designs for player %w", err)
		}
		player.Designs = designs
	}

	return players, nil
}

// get a player by id
func (c *client) GetPlayer(id int64) (*game.Player, error) {
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

func (c *client) GetPlayerForGame(gameID, userID int64) (*game.Player, error) {
	item := Player{}
	if err := c.db.Get(&item, "SELECT * FROM players WHERE gameId = ? AND userId = ?", gameID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	player := c.converter.ConvertPlayer(item)
	return &player, nil
}

// get a full player by id with all dependencies loaded
func (c *client) GetFullPlayerForGame(gameID, userID int64) (*game.FullPlayer, error) {
	player := game.FullPlayer{}

	item := Player{}
	if err := c.db.Get(&item, "SELECT * FROM players WHERE gameId = ? AND userId = ?", gameID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// load the player component from the DB
	player.Player = c.converter.ConvertPlayer(item)

	designs, err := c.GetShipDesignsForPlayer(player.ID)
	if err != nil {
		return nil, fmt.Errorf("get player designs %w", err)
	}
	player.Designs = designs

	planets, err := c.getPlanetsForPlayer(player.ID)
	if err != nil {
		return nil, fmt.Errorf("get player planets %w", err)
	}
	player.Planets = planets

	fleets, err := c.getFleetsForPlayer(player.ID)
	if err != nil {
		return nil, fmt.Errorf("get player fleets %w", err)
	}
	player.Fleets = fleets

	return &player, nil
}

func (c *client) CreatePlayer(player *game.Player) error {
	return c.createPlayer(player, c.db)
}

// create a new game
func (c *client) createPlayer(player *game.Player, tx SQLExecer) error {
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
		messages,
		planetIntels,
		fleetIntels,
		shipDesignIntels,
		mineralPacketIntels,
		mineFieldIntels,
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
		:messages,
		:planetIntels,
		:fleetIntels,
		:shipDesignIntels,
		:mineralPacketIntels,
		:mineFieldIntels,
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
func (c *client) UpdatePlayer(player *game.Player) error {
	return c.updatePlayerWithNamedExecer(player, c.db)
}

// helper to update a player using a transaction or DB
func (c *client) updatePlayerWithNamedExecer(player *game.Player, tx SQLExecer) error {
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
		messages = :messages,
		planetIntels = :planetIntels,
		fleetIntels = :fleetIntels,
		shipDesignIntels = :shipDesignIntels,
		mineralPacketIntels = :mineralPacketIntels,
		mineFieldIntels = :mineFieldIntels,
		race = :race,
		stats = :stats,
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

func (c *client) updateFullPlayer(player *game.Player) error {
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
func (c *client) updateFullPlayerWithTransaction(player *game.Player, tx *sqlx.Tx) error {

	if err := c.updatePlayerWithNamedExecer(player, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("update player %w", err)
	}

	for i := range player.Designs {
		design := &player.Designs[i]
		if design.ID == 0 {
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
