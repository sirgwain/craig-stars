package db

import (
	"database/sql"
	"time"

	"github.com/sirgwain/craig-stars/game"
)

type MineField struct {
	ID        int64              `json:"id,omitempty"`
	GameID    int64              `json:"gameId,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty"`
	PlayerID  int64              `json:"playerId,omitempty"`
	X         float64            `json:"x,omitempty"`
	Y         float64            `json:"y,omitempty"`
	Name      string             `json:"name,omitempty"`
	Num       int                `json:"num,omitempty"`
	PlayerNum int                `json:"playerNum,omitempty"`
	Tags      Tags               `json:"tags,omitempty"`
	Type      game.MineFieldType `json:"type,omitempty"`
	NumMines  int                `json:"numMines,omitempty"`
	Detonate  bool               `json:"detonate,omitempty"`
}

// get a mineField by id
func (c *client) GetMineField(id int64) (*game.MineField, error) {
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

func (c *client) getMineFieldsForGame(gameID int64) ([]*game.MineField, error) {

	items := []MineField{}
	if err := c.db.Select(&items, `SELECT * FROM mineFields WHERE gameId = ?`, gameID); err != nil {
		if err == sql.ErrNoRows {
			return []*game.MineField{}, nil
		}
		return nil, err
	}

	results := make([]*game.MineField, len(items))
	for i := range items {
		results[i] = c.converter.ConvertMineField(&items[i])
	}

	return results, nil
}

// create a new game
func (c *client) createMineField(mineField *game.MineField, tx SQLExecer) error {
	item := c.converter.ConvertGameMineField(mineField)
	result, err := tx.NamedExec(`
	INSERT INTO mineFields (
		createdAt,
		updatedAt,
		gameId,
		playerId,
		x,
		y,
		name,
		num,
		playerNum,
		type,
		numMines,
		detonate
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
		:type,
		:numMines,
		:detonate
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

func (c *client) UpdateMineField(mineField *game.MineField) error {
	return c.updateMineField(mineField, c.db)
}

// update an existing mineField
func (c *client) updateMineField(mineField *game.MineField, tx SQLExecer) error {

	item := c.converter.ConvertGameMineField(mineField)

	if _, err := tx.NamedExec(`
	UPDATE mineFields SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		playerId = :playerId,
		x = :x,
		y = :y,
		name = :name,
		num = :num,
		playerNum = :playerNum,
		type = :type,
		numMines = :numMines,
		detonate = :detonate
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}
