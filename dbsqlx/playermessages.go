package dbsqlx

import (
	"database/sql"

	"github.com/sirgwain/craig-stars/game"
)

func (c *client) GetPlayerMessages() ([]game.PlayerMessage, error) {

	items := []game.PlayerMessage{}
	if err := c.db.Select(&items, `SELECT * FROM playerMessages`); err != nil {
		if err == sql.ErrNoRows {
			return []game.PlayerMessage{}, nil
		}
		return nil, err
	}

	return items, nil
}

func (c *client) GetPlayerMessagesForPlayer(playerID int64) ([]game.PlayerMessage, error) {

	items := []game.PlayerMessage{}
	if err := c.db.Select(&items, `SELECT * FROM playerMessages WHERE playerId = ?`, playerID); err != nil {
		if err == sql.ErrNoRows {
			return []game.PlayerMessage{}, nil
		}
		return nil, err
	}

	return items, nil
}

// get a playerMessage by id
func (c *client) GetPlayerMessage(id int64) (*game.PlayerMessage, error) {
	item := game.PlayerMessage{}
	if err := c.db.Get(&item, "SELECT * FROM playerMessages WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

// create a new playerMessage given something that can execute NamedExec (either a DB or )
func (c *client) CreatePlayerMessage(playerMessage *game.PlayerMessage, tx SQLExecer) error {

	result, err := tx.NamedExec(`
	INSERT INTO playerMessages (
		createdAt,
		updatedAt,		
		playerId,
		type,
		text,
		targetMapObjectNum,
		targetPlayerNum,
		targetType
	)
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:playerId,
		:type,
		:text,
		:targetMapObjectNum,
		:targetPlayerNum,
		:targetType
	)
	`, playerMessage)

	if err != nil {
		return err
	}

	// update the id of our passed in playerMessage
	playerMessage.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// delete a playerMessage by id
func (c *client) DeletePlayerMessage(id int64) error {
	if _, err := c.db.Exec("DELETE FROM playerMessages WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}

func (c *client) DeletePlayerMessagesForPlayer(playerID int64) error {
	if _, err := c.db.Exec("DELETE FROM playerMessages WHERE playerID = ?", playerID); err != nil {
		return err
	}

	return nil
}
