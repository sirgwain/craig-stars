package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

type MineField struct {
	ID            int64            `json:"id,omitempty"`
	GameID        int64            `json:"gameId,omitempty"`
	CreatedAt     time.Time        `json:"createdAt,omitempty"`
	UpdatedAt     time.Time        `json:"updatedAt,omitempty"`
	X             float64          `json:"x,omitempty"`
	Y             float64          `json:"y,omitempty"`
	Name          string           `json:"name,omitempty"`
	Num           int              `json:"num,omitempty"`
	PlayerNum     int              `json:"playerNum,omitempty"`
	Tags          Tags             `json:"tags,omitempty"`
	MineFieldType cs.MineFieldType `json:"mineFieldType,omitempty"`
	NumMines      int              `json:"numMines,omitempty"`
	Detonate      bool             `json:"detonate,omitempty"`
	Spec          *MineFieldSpec   `json:"spec,omitempty"`
}

type MineFieldSpec cs.MineFieldSpec

// db serializer to serialize this to JSON
func (item *MineFieldSpec) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// db deserializer to read this from JSON
func (item *MineFieldSpec) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// get a mineField by id
func (c *client) GetMineField(id int64) (*cs.MineField, error) {
	item := MineField{}
	if err := c.db.Get(&item, "SELECT * FROM mineFields WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	mineField := c.converter.ConvertMineField(&item)
	return mineField, nil
}

func (c *client) GetMineFieldsForPlayer(gameID int64, playerNum int) ([]*cs.MineField, error) {

	items := []MineField{}
	if err := c.db.Select(&items, `SELECT * FROM mineFields WHERE gameId = ? AND playerNum = ?`, gameID, playerNum); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.MineField{}, nil
		}
		return nil, err
	}

	results := make([]*cs.MineField, len(items))
	for i := range items {
		results[i] = c.converter.ConvertMineField(&items[i])
	}

	return results, nil
}

func (c *client) getMineFieldsForGame(gameID int64) ([]*cs.MineField, error) {

	items := []MineField{}
	if err := c.db.Select(&items, `SELECT * FROM mineFields WHERE gameId = ?`, gameID); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.MineField{}, nil
		}
		return nil, err
	}

	results := make([]*cs.MineField, len(items))
	for i := range items {
		results[i] = c.converter.ConvertMineField(&items[i])
	}

	return results, nil
}

// create a new game
func (c *client) createMineField(mineField *cs.MineField, tx SQLExecer) error {
	item := c.converter.ConvertGameMineField(mineField)
	result, err := tx.NamedExec(`
	INSERT INTO mineFields (
		createdAt,
		updatedAt,
		gameId,
		x,
		y,
		name,
		num,
		playerNum,
		mineFieldType,
		numMines,
		detonate,
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
		:mineFieldType,
		:numMines,
		:detonate,
		:spec
	)
	`, item)

	if err != nil {
		return err
	}

	// update the id of our passed in game
	mineField.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) UpdateMineField(mineField *cs.MineField) error {
	return c.updateMineField(mineField, c.db)
}

// update an existing mineField
func (c *client) updateMineField(mineField *cs.MineField, tx SQLExecer) error {

	item := c.converter.ConvertGameMineField(mineField)

	if _, err := tx.NamedExec(`
	UPDATE mineFields SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		x = :x,
		y = :y,
		name = :name,
		num = :num,
		playerNum = :playerNum,
		mineFieldType = :mineFieldType,
		numMines = :numMines,
		detonate = :detonate,
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

func (c *client) deleteMineField(mineFieldID int64, tx SQLExecer) error {
	if _, err := tx.Exec("DELETE FROM mineFields where id = ?", mineFieldID); err != nil {
		return fmt.Errorf("delete mineField %d %w", mineFieldID, err)
	}
	return nil
}
