package dbsqlx

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/sirgwain/craig-stars/game"
)

type Planet struct {
	ID                                int64                 `json:"id,omitempty"`
	GameID                            int64                 `json:"gameId,omitempty"`
	CreatedAt                         time.Time             `json:"createdAt,omitempty"`
	UpdatedAt                         time.Time             `json:"updatedAt,omitempty"`
	PlayerID                          int64                 `json:"playerId,omitempty"`
	X                                 float64               `json:"x,omitempty"`
	Y                                 float64               `json:"y,omitempty"`
	Name                              string                `json:"name,omitempty"`
	Num                               int                   `json:"num,omitempty"`
	PlayerNum                         int                   `json:"playerNum,omitempty"`
	Tags                              Tags                  `json:"tags,omitempty"`
	Grav                              int                   `json:"grav,omitempty"`
	Temp                              int                   `json:"temp,omitempty"`
	Rad                               int                   `json:"rad,omitempty"`
	BaseGrav                          int                   `json:"baseGrav,omitempty"`
	BaseTemp                          int                   `json:"baseTemp,omitempty"`
	BaseRad                           int                   `json:"baseRad,omitempty"`
	TerraformedAmountGrav             int                   `json:"terraformedAmountGrav,omitempty"`
	TerraformedAmountTemp             int                   `json:"terraformedAmountTemp,omitempty"`
	TerraformedAmountRad              int                   `json:"terraformedAmountRad,omitempty"`
	MineralConcIronium                int                   `json:"mineralConcIronium,omitempty"`
	MineralConcBoranium               int                   `json:"mineralConcBoranium,omitempty"`
	MineralConcGermanium              int                   `json:"mineralConcGermanium,omitempty"`
	MineYearsIronium                  int                   `json:"mineYearsIronium,omitempty"`
	MineYearsBoranium                 int                   `json:"mineYearsBoranium,omitempty"`
	MineYearsGermanium                int                   `json:"mineYearsGermanium,omitempty"`
	Ironium                           int                   `json:"ironium,omitempty"`
	Boranium                          int                   `json:"boranium,omitempty"`
	Germanium                         int                   `json:"germanium,omitempty"`
	Colonists                         int                   `json:"colonists,omitempty"`
	Mines                             int                   `json:"mines,omitempty"`
	Factories                         int                   `json:"factories,omitempty"`
	Defenses                          int                   `json:"defenses,omitempty"`
	Homeworld                         bool                  `json:"homeworld,omitempty"`
	ContributesOnlyLeftoverToResearch bool                  `json:"contributesOnlyLeftoverToResearch,omitempty"`
	Scanner                           bool                  `json:"scanner,omitempty"`
	PacketSpeed                       int                   `json:"packetSpeed,omitempty"`
	BonusResources                    int                   `json:"bonusResources,omitempty"`
	ProductionQueue                   *ProductionQueueItems `json:"productionQueue,omitempty"`
	Spec                              *PlanetSpec           `json:"spec,omitempty"`
}

// we json serialize these types with custom Scan/Value methods
type ProductionQueueItems []game.ProductionQueueItem
type PlanetSpec game.PlanetSpec
type Tags game.Tags

// db serializer to serialize this to JSON
func (item *ProductionQueueItems) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *ProductionQueueItems) Scan(src interface{}) error {
	return scanJSON(src, &item)

}

// db serializer to serialize this to JSON
func (item *PlanetSpec) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *PlanetSpec) Scan(src interface{}) error {
	return scanJSON(src, &item)
}

// db serializer to serialize this to JSON
func (item *Tags) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *Tags) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// get a planet by id
func (c *client) GetPlanet(id int64) (*game.Planet, error) {
	item := Planet{}
	if err := c.db.Get(&item, "SELECT * FROM planets WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	planet := c.converter.ConvertPlanet(&item)
	return planet, nil
}

func (c *client) getPlanetsForGame(gameID int64) ([]*game.Planet, error) {

	items := []Planet{}
	if err := c.db.Select(&items, `SELECT * FROM planets WHERE gameId = ?`, gameID); err != nil {
		if err == sql.ErrNoRows {
			return []*game.Planet{}, nil
		}
		return nil, err
	}

	results := make([]*game.Planet, len(items))
	for i := range items {
		results[i] = c.converter.ConvertPlanet(&items[i])
	}

	return results, nil
}

func (c *client) getPlanetsForPlayer(playerID int64) ([]*game.Planet, error) {

	items := []Planet{}
	if err := c.db.Select(&items, `SELECT * FROM planets WHERE playerId = ?`, playerID); err != nil {
		if err == sql.ErrNoRows {
			return []*game.Planet{}, nil
		}
		return nil, err
	}

	results := make([]*game.Planet, len(items))
	for i := range items {
		results[i] = c.converter.ConvertPlanet(&items[i])
	}

	return results, nil
}

// create a new game
func (c *client) createPlanet(planet *game.Planet, tx SQLExecer) error {
	item := c.converter.ConvertGamePlanet(planet)
	result, err := tx.NamedExec(`
	INSERT INTO planets (
		createdAt,
		updatedAt,
		gameId,
		playerId,
		x,
		y,
		name,
		num,
		playerNum,
		grav,
		temp,
		rad,
		baseGrav,
		baseTemp,
		baseRad,
		terraformedAmountGrav,
		terraformedAmountTemp,
		terraformedAmountRad,
		mineralConcIronium,
		mineralConcBoranium,
		mineralConcGermanium,
		mineYearsIronium,
		mineYearsBoranium,
		mineYearsGermanium,
		ironium,
		boranium,
		germanium,
		colonists,
		mines,
		factories,
		defenses,
		homeworld,
		contributesOnlyLeftoverToResearch,
		scanner,
		packetSpeed,
		productionQueue,
		spec
	)
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:gameId,
		:playerId,
		:x,
		:y,
		:name,
		:num,
		:playerNum,
		:grav,
		:temp,
		:rad,
		:baseGrav,
		:baseTemp,
		:baseRad,
		:terraformedAmountGrav,
		:terraformedAmountTemp,
		:terraformedAmountRad,
		:mineralConcIronium,
		:mineralConcBoranium,
		:mineralConcGermanium,
		:mineYearsIronium,
		:mineYearsBoranium,
		:mineYearsGermanium,
		:ironium,
		:boranium,
		:germanium,
		:colonists,
		:mines,
		:factories,
		:defenses,
		:homeworld,
		:contributesOnlyLeftoverToResearch,
		:scanner,
		:packetSpeed,
		:productionQueue,
		:spec
	)
	`, item)

	if err != nil {
		return err
	}

	// update the id of our passed in game
	planet.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) UpdatePlanet(planet *game.Planet) error {
	return c.updatePlanet(planet, c.db)
}

// update an existing planet
func (c *client) updatePlanet(planet *game.Planet, tx SQLExecer) error {

	item := c.converter.ConvertGamePlanet(planet)

	if _, err := tx.NamedExec(`
	UPDATE planets SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		playerId = :playerId,
		x = :x,
		y = :y,
		name = :name,
		num = :num,
		playerNum = :playerNum,
		grav = :grav,
		temp = :temp,
		rad = :rad,
		baseGrav = :baseGrav,
		baseTemp = :baseTemp,
		baseRad = :baseRad,
		terraformedAmountGrav = :terraformedAmountGrav,
		terraformedAmountTemp = :terraformedAmountTemp,
		terraformedAmountRad = :terraformedAmountRad,
		mineralConcIronium = :mineralConcIronium,
		mineralConcBoranium = :mineralConcBoranium,
		mineralConcGermanium = :mineralConcGermanium,
		mineYearsIronium = :mineYearsIronium,
		mineYearsBoranium = :mineYearsBoranium,
		mineYearsGermanium = :mineYearsGermanium,
		ironium = :ironium,
		boranium = :boranium,
		germanium = :germanium,
		colonists = :colonists,
		mines = :mines,
		factories = :factories,
		defenses = :defenses,
		homeworld = :homeworld,
		contributesOnlyLeftoverToResearch = :contributesOnlyLeftoverToResearch,
		scanner = :scanner,
		packetSpeed = :packetSpeed,
		productionQueue = :productionQueue,
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}
