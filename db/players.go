package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

type PlayerStatus struct {
	UpdatedAt     time.Time `json:"updatedAt,omitempty"`
	UserID        int64     `json:"userId,omitempty"`
	Name          string    `json:"name,omitempty"`
	Num           int       `json:"num,omitempty"`
	Ready         bool      `json:"ready,omitempty"`
	AIControlled  bool      `json:"aiControlled,omitempty"`
	SubmittedTurn bool      `json:"submittedTurn,omitempty"`
	Color         string    `json:"color,omitempty"`
	Victor        bool      `json:"victor,omitempty"`
}

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
	BattleRecords                *BattleRecords       `json:"battleRecords,omitempty"`
	PlayerIntels                 *PlayerIntels        `json:"playerIntels,omitempty"`
	ScoreIntels                  *ScoreIntels         `json:"scoreIntels,omitempty"`
	PlanetIntels                 *PlanetIntels        `json:"planetIntels,omitempty"`
	FleetIntels                  *FleetIntels         `json:"fleetIntels,omitempty"`
	StarbaseIntels               *FleetIntels         `json:"starbaseIntels,omitempty"`
	ShipDesignIntels             *ShipDesignIntels    `json:"shipDesignIntels,omitempty"`
	MineralPacketIntels          *MineralPacketIntels `json:"mineralPacketIntels,omitempty"`
	MineFieldIntels              *MineFieldIntels     `json:"mineFieldIntels,omitempty"`
	WormholeIntels               *WormholeIntels      `json:"wormholeIntels,omitempty"`
	MysteryTraderIntels          *MysteryTraderIntels `json:"mysteryTraderIntels,omitempty"`
	SalvageIntels                *SalvageIntels       `json:"salvageIntels,omitempty"`
	Race                         *PlayerRace          `json:"race,omitempty"`
	Stats                        *PlayerStats         `json:"stats,omitempty"`
	ScoreHistory                 *PlayerScores        `json:"scoreHistory,omitempty"`
	AchievedVictoryConditions    cs.Bitmask           `json:"achievedVictoryConditions,omitempty"`
	Victor                       bool                 `json:"victor,omitempty"`
	Spec                         *PlayerSpec          `json:"spec,omitempty"`
}

// we json serialize these types with custom Scan/Value methods
type BattlePlans []cs.BattlePlan
type ProductionPlans []cs.ProductionPlan
type TransportPlans []cs.TransportPlan
type PlayerRelationships []cs.PlayerRelationship
type PlayerMessages []cs.PlayerMessage
type PlayerScores []cs.PlayerScore
type BattleRecords []cs.BattleRecord
type PlayerIntels []cs.PlayerIntel
type ScoreIntels []cs.ScoreIntel
type PlanetIntels []cs.PlanetIntel
type FleetIntels []cs.FleetIntel
type ShipDesignIntels []cs.ShipDesignIntel
type MineralPacketIntels []cs.MineralPacketIntel
type SalvageIntels []cs.SalvageIntel
type MineFieldIntels []cs.MineFieldIntel
type MysteryTraderIntels []cs.MysteryTraderIntel
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

func (item *PlayerScores) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *PlayerScores) Scan(src interface{}) error {
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

func (item *ScoreIntels) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *ScoreIntels) Scan(src interface{}) error {
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

func (item *SalvageIntels) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *SalvageIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *MineFieldIntels) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *MineFieldIntels) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (item *MysteryTraderIntels) Value() (driver.Value, error) {
	return valueJSON(item)
}

func (item *MysteryTraderIntels) Scan(src interface{}) error {
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
	if err := c.reader.Select(&items, `SELECT * FROM players`); err != nil {
		if err == sql.ErrNoRows {
			return []cs.Player{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertPlayers(items), nil
}

func (c *client) GetPlayersForUser(userID int64) ([]cs.Player, error) {

	items := []Player{}
	if err := c.reader.Select(&items, `SELECT * FROM players WHERE userId = ?`, userID); err != nil {
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
	if err := c.reader.Select(&items, `SELECT * FROM players WHERE gameId = ?`, gameID); err != nil {
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
func (c *client) GetPlayersStatusForGame(gameID int64) ([]*cs.Player, error) {

	items := []Player{}
	if err := c.reader.Select(&items, `
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
	color 
	FROM players WHERE gameId = ? ORDER BY num`, gameID); err != nil {
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

func (c *client) getPlayerWithDesigns(where string, args ...interface{}) ([]cs.Player, error) {
	type playerDesignsJoin struct {
		Player     `json:"player,omitempty"`
		ShipDesign `json:"shipDesign,omitempty"`
	}

	rows := []playerDesignsJoin{}

	err := c.reader.Select(&rows, fmt.Sprintf(`
	SELECT 
		p.id AS 'player.id',
		p.createdAt AS 'player.createdAt',
		p.updatedAt AS 'player.updatedAt',
		p.gameId AS 'player.gameId',
		p.userId AS 'player.userId',
		p.name AS 'player.name',
		p.num AS 'player.num',
		p.ready AS 'player.ready',
		p.aiControlled AS 'player.aiControlled',
		p.submittedTurn AS 'player.submittedTurn',
		p.color AS 'player.color',
		p.defaultHullSet AS 'player.defaultHullSet',
		p.techLevelsEnergy AS 'player.techLevelsEnergy',
		p.techLevelsWeapons AS 'player.techLevelsWeapons',
		p.techLevelsPropulsion AS 'player.techLevelsPropulsion',
		p.techLevelsConstruction AS 'player.techLevelsConstruction',
		p.techLevelsElectronics AS 'player.techLevelsElectronics',
		p.techLevelsBiotechnology AS 'player.techLevelsBiotechnology',
		p.techLevelsSpentEnergy AS 'player.techLevelsSpentEnergy',
		p.techLevelsSpentWeapons AS 'player.techLevelsSpentWeapons',
		p.techLevelsSpentPropulsion AS 'player.techLevelsSpentPropulsion',
		p.techLevelsSpentConstruction AS 'player.techLevelsSpentConstruction',
		p.techLevelsSpentElectronics AS 'player.techLevelsSpentElectronics',
		p.techLevelsSpentBiotechnology AS 'player.techLevelsSpentBiotechnology',
		p.researchAmount AS 'player.researchAmount',
		p.researchSpentLastYear AS 'player.researchSpentLastYear',
		p.nextResearchField AS 'player.nextResearchField',
		p.researching AS 'player.researching',
		p.battlePlans AS 'player.battlePlans',
		p.productionPlans AS 'player.productionPlans',
		p.transportPlans AS 'player.transportPlans',
		p.relations AS 'player.relations',
		p.messages AS 'player.messages',
		p.battleRecords AS 'player.battleRecords',
		p.playerIntels AS 'player.playerIntels',
		p.scoreIntels AS 'player.scoreIntels',
		p.planetIntels AS 'player.planetIntels',
		p.fleetIntels AS 'player.fleetIntels',
		p.starbaseIntels AS 'player.starbaseIntels',
		p.shipDesignIntels AS 'player.shipDesignIntels',
		p.mineralPacketIntels AS 'player.mineralPacketIntels',
		p.mineFieldIntels AS 'player.mineFieldIntels',
		p.wormholeIntels AS 'player.wormholeIntels',
		p.mysteryTraderIntels AS 'player.mysteryTraderIntels',
		p.salvageIntels AS 'player.salvageIntels',
		p.race AS 'player.race',
		p.stats AS 'player.stats',
		p.scoreHistory AS 'player.scoreHistory',
		p.achievedVictoryConditions AS 'player.achievedVictoryConditions',
		p.victor AS 'player.victor',
		p.spec AS 'player.spec',

		
		COALESCE(d.id, 0) AS 'shipDesign.id',
		d.createdAt AS 'shipDesign.createdAt',
		d.updatedAt AS 'shipDesign.updatedAt',
		COALESCE(d.gameId, 0) AS 'shipDesign.gameId',
		COALESCE(d.num, 0) AS 'shipDesign.num',
		COALESCE(d.playerNum, 0) AS 'shipDesign.playerNum',
		COALESCE(d.name, '') AS 'shipDesign.name',
		COALESCE(d.version, 0) AS 'shipDesign.version',
		COALESCE(d.hull, '') AS 'shipDesign.hull',
		COALESCE(d.hullSetNumber, 0) AS 'shipDesign.hullSetNumber',
		COALESCE(d.canDelete, 0) AS 'shipDesign.canDelete',
		COALESCE(d.slots, '[]') AS 'shipDesign.slots',
		COALESCE(d.purpose, '') AS 'shipDesign.purpose',
		COALESCE(d.spec, '{}') AS 'shipDesign.spec'

	FROM players p
	LEFT JOIN shipDesigns d
		ON p.gameId = d.gameId AND p.num = d.playerNum
	WHERE %s
`, where), args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []cs.Player{}, nil
		}
		return nil, err
	}

	// check if we have a game
	if len(rows) == 0 {
		return []cs.Player{}, nil
	}

	// join results give a row per item, so if we have 2 players
	// one with 2 designs, one with 3, we'll end up with 5 rows
	// row 0 - player1, design 1
	// row 1 - player1, design 2
	// row 2 - player2, design 1
	// row 3 - player2, design 2
	// row 4 - player2, design 3
	players := []cs.Player{}
	var item Player
	var player *cs.Player
	for _, row := range rows {

		if row.Player.ID != item.ID {
			// convert this row into a game
			item = row.Player
			p := c.converter.ConvertPlayer(item)
			p.Designs = []*cs.ShipDesign{}
			players = append(players, p)
			player = &players[len(players)-1]
		}

		if row.ShipDesign.ID != 0 {
			design := c.converter.ConvertShipDesign(&row.ShipDesign)
			player.Designs = append(player.Designs, design)
		}
	}

	return players, nil
}

// get a player by id
func (c *client) GetPlayer(id int64) (*cs.Player, error) {
	item := Player{}
	if err := c.reader.Get(&item, "SELECT * FROM players WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	player := c.converter.ConvertPlayer(item)
	return &player, nil
}

// Get all player data except universe intel
func (c *client) GetPlayerForGame(gameID, userID int64) (*cs.Player, error) {
	item := Player{}
	if err := c.reader.Get(&item, `
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
	race,
	stats,
	scoreHistory,
	achievedVictoryConditions,
	victor,
	spec
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

// Get player intel
func (c *client) GetPlayerIntelsForGame(gameID, userID int64) (*cs.PlayerIntels, error) {
	item := Player{}
	if err := c.reader.Get(&item, `
	SELECT
	battleRecords,
	playerIntels,
	scoreIntels,
	planetIntels,
	fleetIntels,
	starbaseIntels,
	shipDesignIntels,
	mineralPacketIntels,
	mineFieldIntels,
	wormholeIntels,
	mysteryTraderIntels,
	salvageIntels
	FROM players 
	WHERE gameId = ? AND userId = ?`, gameID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	player := c.converter.ConvertPlayer(item)

	return &player.PlayerIntels, nil
}

func (c *client) GetPlayerByNum(gameID int64, num int) (*cs.Player, error) {
	item := Player{}
	if err := c.reader.Get(&item, "SELECT * FROM players WHERE gameId = ? AND num = ?", gameID, num); err != nil {
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
	if err := c.reader.Get(&item, `
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
	techLevelsBiotechnology,
	techLevelsSpentEnergy,
	techLevelsSpentWeapons,
	techLevelsSpentPropulsion,
	techLevelsSpentConstruction,
	techLevelsSpentElectronics,
	techLevelsSpentBiotechnology,
	researchSpentLastYear,
	researchAmount,
	nextResearchField,
	researching,
	battlePlans,
	productionPlans,
	transportPlans,
	relations,
	stats,
	scoreHistory,
	achievedVictoryConditions,
	victor,
	spec
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
	if err := c.reader.Get(&item, "SELECT * FROM players WHERE gameId = ? AND userId = ?", gameID, userID); err != nil {
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

	mineFields, err := c.GetMineFieldsForPlayer(player.GameID, player.Num)
	if err != nil {
		return nil, fmt.Errorf("get player mineFields %w", err)
	}
	player.MineFields = mineFields

	mineralPackets, err := c.GetMineralPacketsForPlayer(player.GameID, player.Num)
	if err != nil {
		return nil, fmt.Errorf("get player mineralPackets %w", err)
	}
	player.MineralPackets = mineralPackets

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
	if err := c.reader.Get(&num, "SELECT num FROM players WHERE gameId = ? AND userId = ?", gameID, userID); err != nil {
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

	mineFields, err := c.GetMineFieldsForPlayer(gameID, num)
	if err != nil {
		return nil, fmt.Errorf("get player mineFields %w", err)
	}
	mapObjects.MineFields = mineFields

	mineralPackets, err := c.GetMineralPacketsForPlayer(gameID, num)
	if err != nil {
		return nil, fmt.Errorf("get player mineralPackets %w", err)
	}
	mapObjects.MineralPackets = mineralPackets

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
func (c *client) GetPlayerWithDesignsForGame(gameID int64, num int) (*cs.Player, error) {
	player := cs.Player{}

	item := Player{}
	if err := c.reader.Get(&item, "SELECT * FROM players WHERE gameId = ? AND num = ?", gameID, num); err != nil {
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
	item := c.converter.ConvertGamePlayer(player)
	result, err := c.writer.NamedExec(`
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
		battleRecords,
		playerIntels,
		scoreIntels,
		planetIntels,
		fleetIntels,
		starbaseIntels,
		shipDesignIntels,
		mineralPacketIntels,
		mineFieldIntels,
		wormholeIntels,
		mysteryTraderIntels,
		salvageIntels,
		race,
		stats,
		scoreHistory,
		achievedVictoryConditions,
		victor,
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
		:battleRecords,
		:playerIntels,
		:scoreIntels,
		:planetIntels,
		:fleetIntels,
		:starbaseIntels,
		:shipDesignIntels,
		:mineralPacketIntels,
		:mineFieldIntels,
		:wormholeIntels,
		:mysteryTraderIntels,
		:salvageIntels,
		:race,
		:stats,
		:scoreHistory,
		:achievedVictoryConditions,
		:victor,
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

// update an existing player's lightweight fields
func (c *client) UpdateLightPlayer(player *cs.Player) error {
	item := c.converter.ConvertGamePlayer(player)

	if _, err := c.writer.NamedExec(`
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

// update an existing player's lightweight fields
func (c *client) UpdatePlayerOrders(player *cs.Player) error {
	item := c.converter.ConvertGamePlayer(player)

	if _, err := c.writer.NamedExec(`
	UPDATE players SET
		updatedAt = CURRENT_TIMESTAMP,
		submittedTurn = :submittedTurn,
		defaultHullSet = :defaultHullSet,
		researchAmount = :researchAmount,
		nextResearchField = :nextResearchField,
		researching = :researching,
		battlePlans = :battlePlans,
		productionPlans = :productionPlans,
		transportPlans = :transportPlans,
		relations = :relations,
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// update an existing player's lightweight fields
func (c *client) SubmitPlayerTurn(gameID int64, num int, submittedTurn bool) error {
	type submitData struct {
		GameID        int64 `json:"gameID"`
		Num           int   `json:"num"`
		SubmittedTurn bool  `json:"submittedTurn"`
	}

	if _, err := c.writer.NamedExec(`
	UPDATE players SET
		updatedAt = CURRENT_TIMESTAMP,
		submittedTurn = :submittedTurn
	WHERE gameId = :gameID AND num = :num
	`, submitData{GameID: gameID, Num: num, SubmittedTurn: submittedTurn}); err != nil {
		return err
	}

	return nil
}

// update an existing player's lightweight fields
func (c *client) UpdatePlayerPlans(player *cs.Player) error {
	item := c.converter.ConvertGamePlayer(player)

	if _, err := c.writer.NamedExec(`
	UPDATE players SET
		updatedAt = CURRENT_TIMESTAMP,
		battlePlans = :battlePlans,
		productionPlans = :productionPlans,
		transportPlans = :transportPlans
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// update a players salvage intels (used after creating a new salvage)
func (c *client) UpdatePlayerSalvageIntels(player *cs.Player) error {
	item := c.converter.ConvertGamePlayer(player)

	if _, err := c.writer.NamedExec(`
	UPDATE players SET
		updatedAt = CURRENT_TIMESTAMP,
		salvageIntels = :salvageIntels
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// helper to update a player using a transaction or DB
// update an existing player
func (c *client) UpdatePlayer(player *cs.Player) error {
	item := c.converter.ConvertGamePlayer(player)

	if _, err := c.writer.NamedExec(`
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
		battleRecords = :battleRecords,
		playerIntels = :playerIntels,
		scoreIntels = :scoreIntels,
		planetIntels = :planetIntels,
		fleetIntels = :fleetIntels,
		starbaseIntels = :starbaseIntels,
		shipDesignIntels = :shipDesignIntels,
		mineralPacketIntels = :mineralPacketIntels,
		mineFieldIntels = :mineFieldIntels,
		wormholeIntels = :wormholeIntels,
		mysteryTraderIntels = :mysteryTraderIntels,
		salvageIntels = :salvageIntels,
		messages = :messages,
		race = :race,
		stats = :stats,
		scoreHistory = :scoreHistory,
		achievedVictoryConditions = :achievedVictoryConditions,
		victor = :victor,
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}


// delete a player by id
func (c *client) DeletePlayer(id int64) error {
	if _, err := c.writer.Exec("DELETE FROM players WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}
