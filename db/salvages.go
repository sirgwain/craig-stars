package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

type Salvage struct {
	ID        int64     `json:"id,omitempty"`
	GameID    int64     `json:"gameId,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	X         float64   `json:"x,omitempty"`
	Y         float64   `json:"y,omitempty"`
	Name      string    `json:"name,omitempty"`
	Num       int       `json:"num,omitempty"`
	PlayerNum int       `json:"playerNum,omitempty"`
	Tags      Tags      `json:"tags,omitempty"`
	Ironium   int       `json:"ironium,omitempty"`
	Boranium  int       `json:"boranium,omitempty"`
	Germanium int       `json:"germanium,omitempty"`
}

// get a salvage by id
func (c *client) GetSalvage(id int64) (*cs.Salvage, error) {
	item := Salvage{}
	if err := c.db.Get(&item, "SELECT * FROM salvages WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	salvage := c.converter.ConvertSalvage(&item)
	return salvage, nil
}

func (c *client) GetSalvageByNum(gameID int64, num int) (*cs.Salvage, error) {

	item := Salvage{}
	if err := c.db.Get(&item, `SELECT * FROM salvages WHERE gameId = ? AND num = ?`, gameID, num); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	salvage := c.converter.ConvertSalvage(&item)
	return salvage, nil

}

func (c *client) GetSalvagesForGame(gameID int64) ([]*cs.Salvage, error) {
	return c.getSalvagesForGame(c.db, gameID)
}

func (c *client) getSalvagesForGame(db SQLSelector, gameID int64) ([]*cs.Salvage, error) {

	items := []Salvage{}
	if err := db.Select(&items, `SELECT * FROM salvages WHERE gameId = ? ORDER BY num`, gameID); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.Salvage{}, nil
		}
		return nil, err
	}

	results := make([]*cs.Salvage, len(items))
	for i := range items {
		results[i] = c.converter.ConvertSalvage(&items[i])
	}

	return results, nil
}

func (c *client) GetSalvagesForPlayer(gameID int64, playerNum int) ([]*cs.Salvage, error) {

	items := []Salvage{}
	if err := c.db.Select(&items, `SELECT * FROM salvages WHERE gameId = ? AND playerNum = ?`, gameID, playerNum); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.Salvage{}, nil
		}
		return nil, err
	}

	results := make([]*cs.Salvage, len(items))
	for i := range items {
		results[i] = c.converter.ConvertSalvage(&items[i])
	}

	return results, nil
}

func (c *client) CreateSalvage(salvage *cs.Salvage) error {
	return c.createSalvage(salvage, c.db)
}

// create a new salvage
func (c *client) createSalvage(salvage *cs.Salvage, tx SQLExecer) error {
	item := c.converter.ConvertGameSalvage(salvage)
	result, err := tx.NamedExec(`
	INSERT INTO salvages (
		createdAt,
		updatedAt,
		gameId,
		x,
		y,
		name,
		num,
		playerNum,
		ironium,
		boranium,
		germanium
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
		:ironium,
		:boranium,
		:germanium
	)
	`, item)

	if err != nil {
		return err
	}

	// update the id of our passed in game
	salvage.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) UpdateSalvage(salvage *cs.Salvage) error {
	return c.updateSalvage(salvage, c.db)
}

// update an existing salvage
func (c *client) updateSalvage(salvage *cs.Salvage, tx SQLExecer) error {

	item := c.converter.ConvertGameSalvage(salvage)

	if _, err := tx.NamedExec(`
	UPDATE salvages SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		x = :x,
		y = :y,
		name = :name,
		num = :num,
		playerNum = :playerNum,
		ironium = :ironium,
		boranium = :boranium,
		germanium = :germanium
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

func (c *client) deleteSalvage(salvageID int64, tx SQLExecer) error {
	if _, err := tx.Exec("DELETE FROM salvages where id = ?", salvageID); err != nil {
		return fmt.Errorf("delete salvage %d %w", salvageID, err)
	}
	return nil
}
