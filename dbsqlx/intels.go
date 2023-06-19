package dbsqlx

import (
	"database/sql"
	"time"

	"github.com/sirgwain/craig-stars/game"
)

type PlanetIntel struct {
	ID                   int64              `json:"id,omitempty"`
	CreatedAt            time.Time          `json:"createdAt,omitempty"`
	UpdatedAt            time.Time          `json:"updatedAt,omitempty"`
	Dirty                bool               `json:"dirty,omitempty"`
	PlayerID             int64              `json:"playerId,omitempty"`
	Name                 string             `json:"name,omitempty"`
	Num                  int                `json:"num,omitempty"`
	PlayerNum            int                `json:"playerNum,omitempty"`
	ReportAge            int                `json:"reportAge,omitempty"`
	Type                 game.MapObjectType `json:"type,omitempty"`
	X                    float64            `json:"x,omitempty"`
	Y                    float64            `json:"y,omitempty"`
	Grav                 int                `json:"grav,omitempty"`
	Temp                 int                `json:"temp,omitempty"`
	Rad                  int                `json:"rad,omitempty"`
	MineralConcIronium   int                `json:"mineralConcIronium,omitempty"`
	MineralConcBoranium  int                `json:"mineralConcBoranium,omitempty"`
	MineralConcGermanium int                `json:"mineralConcGermanium,omitempty"`
	Population           uint               `json:"population,omitempty"`
	Ironium              int                `json:"ironium,omitempty"`
	Boranium             int                `json:"boranium,omitempty"`
	Germanium            int                `json:"germanium,omitempty"`
	Colonists            int                `json:"colonists,omitempty"`
	CargoDiscovered      bool               `json:"cargoDiscovered,omitempty"`
}

func (c *client) GetPlanetIntelsForPlayer(playerID int64) ([]game.PlanetIntel, error) {

	items := []PlanetIntel{}
	if err := c.db.Select(&items, `SELECT * FROM planetIntels WHERE playerId = ?`, playerID); err != nil {
		if err == sql.ErrNoRows {
			return []game.PlanetIntel{}, nil
		}
		return nil, err
	}

	result := make([]game.PlanetIntel, len(items))

	for i := range items {
		result[i] = *c.converter.ConvertPlanetIntel(&items[i])
	}

	return result, nil
}

// get a planetIntel by id
func (c *client) GetPlanetIntel(id int64) (*game.PlanetIntel, error) {
	item := PlanetIntel{}
	if err := c.db.Get(&item, "SELECT * FROM planetIntels WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return c.converter.ConvertPlanetIntel(&item), nil
}

func (c *client) CreatePlanetIntel(planetIntel *game.PlanetIntel) error {
	return c.createPlanetIntel(planetIntel, c.db)
}

// create a new planetIntel given something that can execute NamedExec (either a DB or )
func (c *client) createPlanetIntel(planetIntel *game.PlanetIntel, tx SQLExecer) error {
	item := c.converter.ConvertGamePlanetIntel(planetIntel)

	result, err := tx.NamedExec(`
	INSERT INTO planetIntels (
		createdAt,
		updatedAt,
		playerId,
		name,
		num,
		playerNum,
		reportAge,
		x,
		y,
		grav,
		temp,
		rad,
		mineralConcIronium,
		mineralConcBoranium,
		mineralConcGermanium,
		population,
		ironium,
		boranium,
		germanium,
		colonists,
		cargoDiscovered
	)
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:playerId,
		:name,
		:num,
		:playerNum,
		:reportAge,
		:x,
		:y,
		:grav,
		:temp,
		:rad,
		:mineralConcIronium,
		:mineralConcBoranium,
		:mineralConcGermanium,
		:population,
		:ironium,
		:boranium,
		:germanium,
		:colonists,
		:cargoDiscovered		
	)
	`, item)

	if err != nil {
		return err
	}

	// update the id of our passed in planetIntel
	planetIntel.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) UpdatePlanetIntel(planetIntel *game.PlanetIntel) error {
	return c.updatePlanetIntel(planetIntel, c.db)
}

// update an existing planetIntel
func (c *client) updatePlanetIntel(planetIntel *game.PlanetIntel, tx SQLExecer) error {

	item := c.converter.ConvertGamePlanetIntel(planetIntel)

	if _, err := tx.NamedExec(`
	UPDATE planetIntels SET
		updatedAt = CURRENT_TIMESTAMP,
		playerId = :playerId,
		name = :name,
		num = :num,
		playerNum = :playerNum,
		reportAge = :reportAge,
		x = :x,
		y = :y,
		grav = :grav,
		temp = :temp,
		rad = :rad,
		mineralConcIronium = :mineralConcIronium,
		mineralConcBoranium = :mineralConcBoranium,
		mineralConcGermanium = :mineralConcGermanium,
		population = :population,
		ironium = :ironium,
		boranium = :boranium,
		germanium = :germanium,
		colonists = :colonists,
		cargoDiscovered = :cargoDiscovered
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// delete a planetIntel by id
func (c *client) DeletePlanetIntel(id int64) error {
	if _, err := c.db.Exec("DELETE FROM planetIntels WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}
