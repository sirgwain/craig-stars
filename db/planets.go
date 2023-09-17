package db

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

type Planet struct {
	ID                                int64                 `json:"id,omitempty"`
	GameID                            int64                 `json:"gameId,omitempty"`
	CreatedAt                         time.Time             `json:"createdAt,omitempty"`
	UpdatedAt                         time.Time             `json:"updatedAt,omitempty"`
	X                                 float64               `json:"x,omitempty"`
	Y                                 float64               `json:"y,omitempty"`
	Name                              string                `json:"name,omitempty"`
	Num                               int                   `json:"num,omitempty"`
	PlayerNum                         int                   `json:"playerNum,omitempty"`
	Tags                              *Tags                 `json:"tags,omitempty"`
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
	RouteTargetType                   cs.MapObjectType      `json:"routeTargetType,omitempty"`
	RouteTargetNum                    int                   `json:"routeTargetNum,omitempty"`
	RouteTargetPlayerNum              int                   `json:"routeTargetPlayerNum,omitempty"`
	PacketTargetNum                   int                   `json:"packetTargetNum,omitempty"`
	PacketSpeed                       int                   `json:"packetSpeed,omitempty"`
	RandomArtifact                    bool                  `json:"randomArtifact,omitempty"`
	ProductionQueue                   *ProductionQueueItems `json:"productionQueue,omitempty"`
	Spec                              *PlanetSpec           `json:"spec,omitempty"`
}

// we json serialize these types with custom Scan/Value methods
type ProductionQueueItems []cs.ProductionQueueItem
type PlanetSpec cs.PlanetSpec
type Tags cs.Tags

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
func (c *client) GetPlanet(id int64) (*cs.Planet, error) {
	item := Planet{}
	if err := c.reader.Get(&item, "SELECT * FROM planets WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	planet := c.converter.ConvertPlanet(&item)
	return planet, nil
}

func (c *client) getPlanetsForGame(gameID int64) ([]*cs.Planet, error) {

	items := []Planet{}
	if err := c.reader.Select(&items, `SELECT * FROM planets WHERE gameId = ?`, gameID); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.Planet{}, nil
		}
		return nil, err
	}

	results := make([]*cs.Planet, len(items))
	for i := range items {
		results[i] = c.converter.ConvertPlanet(&items[i])
	}

	return results, nil
}

func (c *client) GetPlanetsForPlayer(gameID int64, playerNum int) ([]*cs.Planet, error) {

	items := []Planet{}
	if err := c.reader.Select(&items, `SELECT * FROM planets WHERE gameId = ? AND playerNum = ?`, gameID, playerNum); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.Planet{}, nil
		}
		return nil, err
	}

	results := make([]*cs.Planet, len(items))
	for i := range items {
		results[i] = c.converter.ConvertPlanet(&items[i])
	}

	return results, nil
}

func (c *client) GetPlanetByNum(gameID int64, num int) (*cs.Planet, error) {

	type planetStarbaseJoin struct {
		Planet `json:"planet,omitempty"`
		Fleet  `json:"fleet,omitempty"`
	}

	item := planetStarbaseJoin{}
	if err := c.reader.Get(&item, `
	SELECT 
		p.id AS 'planet.id',
		p.createdAt AS 'planet.createdAt',
		p.updatedAt AS 'planet.updatedAt',
		p.gameId AS 'planet.gameId',
		p.x AS 'planet.x',
		p.y AS 'planet.y',
		p.name AS 'planet.name',
		p.num AS 'planet.num',
		p.playerNum AS 'planet.playerNum',
		p.tags AS 'planet.tags',
		p.grav AS 'planet.grav',
		p.temp AS 'planet.temp',
		p.rad AS 'planet.rad',
		p.baseGrav AS 'planet.baseGrav',
		p.baseTemp AS 'planet.baseTemp',
		p.baseRad AS 'planet.baseRad',
		p.terraformedAmountGrav AS 'planet.terraformedAmountGrav',
		p.terraformedAmountTemp AS 'planet.terraformedAmountTemp',
		p.terraformedAmountRad AS 'planet.terraformedAmountRad',
		p.mineralConcIronium AS 'planet.mineralConcIronium',
		p.mineralConcBoranium AS 'planet.mineralConcBoranium',
		p.mineralConcGermanium AS 'planet.mineralConcGermanium',
		p.mineYearsIronium AS 'planet.mineYearsIronium',
		p.mineYearsBoranium AS 'planet.mineYearsBoranium',
		p.mineYearsGermanium AS 'planet.mineYearsGermanium',
		p.ironium AS 'planet.ironium',
		p.boranium AS 'planet.boranium',
		p.germanium AS 'planet.germanium',
		p.colonists AS 'planet.colonists',
		p.mines AS 'planet.mines',
		p.factories AS 'planet.factories',
		p.defenses AS 'planet.defenses',
		p.homeworld AS 'planet.homeworld',
		p.contributesOnlyLeftoverToResearch AS 'planet.contributesOnlyLeftoverToResearch',
		p.scanner AS 'planet.scanner',
		p.routeTargetType AS 'planet.routeTargetType',
		p.routeTargetNum AS 'planet.routeTargetNum',
		p.routeTargetPlayerNum AS 'planet.routeTargetPlayerNum',
		p.packetTargetNum AS 'planet.packetTargetNum',
		p.packetSpeed AS 'planet.packetSpeed',
		p.randomArtifact AS 'planet.randomArtifact',
		p.productionQueue AS 'planet.productionQueue',
		p.spec AS 'planet.spec',

		f.createdAt AS 'fleet.createdAt',
		f.updatedAt AS 'fleet.updatedAt',
		COALESCE(f.id, 0) AS 'fleet.id',
		COALESCE(f.gameId, 0) AS 'fleet.gameId',
		COALESCE(f.battlePlanNum, 0) AS 'fleet.battlePlanNum',
		COALESCE(f.x, 0) AS 'fleet.x',
		COALESCE(f.y, 0) AS 'fleet.y',
		COALESCE(f.name, '') AS 'fleet.name',
		COALESCE(f.num, 0) AS 'fleet.num',
		COALESCE(f.playerNum, 0) AS 'fleet.playerNum',
		COALESCE(f.tags, '{}') AS 'fleet.tags',
		COALESCE(f.tokens, '{}') AS 'fleet.tokens',
		COALESCE(f.waypoints, '{}') AS 'fleet.waypoints',
		COALESCE(f.repeatOrders, 0) AS 'fleet.repeatOrders',
		COALESCE(f.planetNum, 0) AS 'fleet.planetNum',
		COALESCE(f.baseName, '') AS 'fleet.baseName',
		COALESCE(f.ironium, 0) AS 'fleet.ironium',
		COALESCE(f.boranium, 0) AS 'fleet.boranium',
		COALESCE(f.germanium, 0) AS 'fleet.germanium',
		COALESCE(f.colonists, 0) AS 'fleet.colonists',
		COALESCE(f.fuel, 0) AS 'fleet.fuel',
		COALESCE(f.age, 0) AS 'fleet.age',
		COALESCE(f.headingX, 0) AS 'fleet.headingX',
		COALESCE(f.headingY, 0) AS 'fleet.headingY',
		COALESCE(f.warpSpeed, 0) AS 'fleet.warpSpeed',
		COALESCE(f.orbitingPlanetNum, 0) AS 'fleet.orbitingPlanetNum',
		COALESCE(f.starbase, 0) AS 'fleet.starbase',
		COALESCE(f.spec, '{}') AS 'fleet.spec'	

	FROM planets p 
	LEFT JOIN fleets f 
	WHERE p.gameId = ? AND p.num = ?`, gameID, num); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	planet := c.converter.ConvertPlanet(&item.Planet)
	if item.Fleet.ID != 0 {
		planet.Starbase = c.converter.ConvertFleet(&item.Fleet)
	}
	return planet, nil

}

// create a new game
func (c *client) createPlanet(planet *cs.Planet) error {
	item := c.converter.ConvertGamePlanet(planet)
	result, err := c.writer.NamedExec(`
	INSERT INTO planets (
		createdAt,
		updatedAt,
		gameId,
		x,
		y,
		name,
		num,
		playerNum,
		tags,
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
		routeTargetType,
		routeTargetNum,
		routeTargetPlayerNum,
		packetTargetNum,
		packetSpeed,
		randomArtifact,
		productionQueue,
		spec
	)
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:gameId,
		:x,
		:y,
		:name,
		:num,
		:playerNum,
		:tags,
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
		:routeTargetType,
		:routeTargetNum,
		:routeTargetPlayerNum,
		:packetTargetNum,
		:packetSpeed,
		:randomArtifact,
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

// update an existing planet
func (c *client) UpdatePlanet(planet *cs.Planet) error {

	item := c.converter.ConvertGamePlanet(planet)

	if _, err := c.writer.NamedExec(`
	UPDATE planets SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		x = :x,
		y = :y,
		name = :name,
		num = :num,
		playerNum = :playerNum,
		tags = :tags,
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
		routeTargetType = :routeTargetType,
		routeTargetNum = :routeTargetNum,
		routeTargetPlayerNum = :routeTargetPlayerNum,
		packetTargetNum = :packetTargetNum,
		packetSpeed = :packetSpeed,
		randomArtifact = :randomArtifact,
		productionQueue = :productionQueue,
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// UpdatePlanetSpec updates only a planets spec field
func (c *client) UpdatePlanetSpec(planet *cs.Planet) error {
	item := c.converter.ConvertGamePlanet(planet)

	if _, err := c.writer.NamedExec(`
	UPDATE planets SET
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}
