package dbsqlx

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/sirgwain/craig-stars/game"
)

type Race struct {
	ID                        int64                  `json:"id,omitempty"`
	CreatedAt                 time.Time              `json:"createdAt,omitempty"`
	UpdatedAt                 time.Time              `json:"updatedAt,omitempty"`
	UserID                    int64                  `json:"userId,omitempty"`
	Name                      string                 `json:"name,omitempty"`
	PluralName                string                 `json:"pluralName,omitempty"`
	PRT                       game.PRT               `json:"prt,omitempty"`
	LRTs                      game.Bitmask           `json:"lrts,omitempty"`
	HabLowGrav                int                    `json:"habLowGrav,omitempty"`
	HabLowTemp                int                    `json:"habLowTemp,omitempty"`
	HabLowRad                 int                    `json:"habLowRad,omitempty"`
	HabHighGrav               int                    `json:"habHighGrav,omitempty"`
	HabHighTemp               int                    `json:"habHighTemp,omitempty"`
	HabHighRad                int                    `json:"habHighRad,omitempty"`
	GrowthRate                int                    `json:"growthRate,omitempty"`
	PopEfficiency             int                    `json:"popEfficiency,omitempty"`
	FactoryOutput             int                    `json:"factoryOutput,omitempty"`
	FactoryCost               int                    `json:"factoryCost,omitempty"`
	NumFactories              int                    `json:"numFactories,omitempty"`
	FactoriesCostLess         bool                   `json:"factoriesCostLess,omitempty"`
	ImmuneGrav                bool                   `json:"immuneGrav,omitempty"`
	ImmuneTemp                bool                   `json:"immuneTemp,omitempty"`
	ImmuneRad                 bool                   `json:"immuneRad,omitempty"`
	MineOutput                int                    `json:"mineOutput,omitempty"`
	MineCost                  int                    `json:"mineCost,omitempty"`
	NumMines                  int                    `json:"numMines,omitempty"`
	ResearchCostEnergy        game.ResearchCostLevel `json:"researchCostEnergy,omitempty"`
	ResearchCostWeapons       game.ResearchCostLevel `json:"researchCostWeapons,omitempty"`
	ResearchCostPropulsion    game.ResearchCostLevel `json:"researchCostPropulsion,omitempty"`
	ResearchCostConstruction  game.ResearchCostLevel `json:"researchCostConstruction,omitempty"`
	ResearchCostElectronics   game.ResearchCostLevel `json:"researchCostElectronics,omitempty"`
	ResearchCostBiotechnology game.ResearchCostLevel `json:"researchCostBiotechnology,omitempty"`
	TechsStartHigh            bool                   `json:"techsStartHigh,omitempty"`
	Spec                      *RaceSpec              `json:"spec,omitempty"`
}

// we json serialize these types with custom Scan/Value methods
type RaceSpec game.RaceSpec

// db serializer to serialize this to JSON
func (item *RaceSpec) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *RaceSpec) Scan(src interface{}) error {
	return scanJSON(src, item)
}

func (c *client) GetRaces() ([]game.Race, error) {

	items := []Race{}
	if err := c.db.Select(&items, `SELECT * FROM races`); err != nil {
		if err == sql.ErrNoRows {
			return []game.Race{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertRaces(items), nil
}

func (c *client) GetRacesForUser(userID int64) ([]game.Race, error) {

	items := []Race{}
	if err := c.db.Select(&items, `SELECT * FROM races WHERE userId = ?`, userID); err != nil {
		if err == sql.ErrNoRows {
			return []game.Race{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertRaces(items), nil
}

// get a race by id
func (c *client) GetRace(id int64) (*game.Race, error) {
	item := Race{}
	if err := c.db.Get(&item, "SELECT * FROM races WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	race := c.converter.ConvertRace(item)
	return &race, nil
}

func (c *client) CreateRace(race *game.Race) error {
	return c.createRace(race, c.db)
}

// create a new race
func (c *client) createRace(race *game.Race, tx SQLExecer) error {

	item := c.converter.ConvertGameRace(race)
	result, err := tx.NamedExec(`
	INSERT INTO races (
		createdAt,
		updatedAt,
		userId,
		name,
		pluralName,
		prt,
		lrts,
		habLowGrav,
		habLowTemp,
		habLowRad,
		habHighGrav,
		habHighTemp,
		habHighRad,
		growthRate,
		popEfficiency,
		factoryOutput,
		factoryCost,
		numFactories,
		factoriesCostLess,
		immuneGrav,
		immuneTemp,
		immuneRad,
		mineOutput,
		mineCost,
		numMines,
		researchCostEnergy,
		researchCostWeapons,
		researchCostPropulsion,
		researchCostConstruction,
		researchCostElectronics,
		researchCostBiotechnology,
		techsStartHigh,
		spec
	)
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:userId,
		:name,
		:pluralName,
		:prt,
		:lrts,
		:habLowGrav,
		:habLowTemp,
		:habLowRad,
		:habHighGrav,
		:habHighTemp,
		:habHighRad,
		:growthRate,
		:popEfficiency,
		:factoryOutput,
		:factoryCost,
		:numFactories,
		:factoriesCostLess,
		:immuneGrav,
		:immuneTemp,
		:immuneRad,
		:mineOutput,
		:mineCost,
		:numMines,
		:researchCostEnergy,
		:researchCostWeapons,
		:researchCostPropulsion,
		:researchCostConstruction,
		:researchCostElectronics,
		:researchCostBiotechnology,
		:techsStartHigh,
		:spec
	)
	`, item)

	if err != nil {
		return err
	}

	// update the id of our passed in race
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	race.ID = id

	return nil
}

// update an existing race
func (c *client) UpdateRace(race *game.Race) error {

	item := c.converter.ConvertGameRace(race)

	if _, err := c.db.NamedExec(`
	UPDATE races SET
		updatedAt = CURRENT_TIMESTAMP,
		userId = :userId,
		name = :name,
		pluralName = :pluralName,
		prt = :prt,
		lrts = :lrts,
		habLowGrav = :habLowGrav,
		habLowTemp = :habLowTemp,
		habLowRad = :habLowRad,
		habHighGrav = :habHighGrav,
		habHighTemp = :habHighTemp,
		habHighRad = :habHighRad,
		growthRate = :growthRate,
		popEfficiency = :popEfficiency,
		factoryOutput = :factoryOutput,
		factoryCost = :factoryCost,
		numFactories = :numFactories,
		factoriesCostLess = :factoriesCostLess,
		immuneGrav = :immuneGrav,
		immuneTemp = :immuneTemp,
		immuneRad = :immuneRad,
		mineOutput = :mineOutput,
		mineCost = :mineCost,
		numMines = :numMines,
		researchCostEnergy = :researchCostEnergy,
		researchCostWeapons = :researchCostWeapons,
		researchCostPropulsion = :researchCostPropulsion,
		researchCostConstruction = :researchCostConstruction,
		researchCostElectronics = :researchCostElectronics,
		researchCostBiotechnology = :researchCostBiotechnology,
		techsStartHigh = :techsStartHigh,
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// delete a race by id
func (c *client) DeleteRace(id int64) error {
	if _, err := c.db.Exec("DELETE FROM races WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}
