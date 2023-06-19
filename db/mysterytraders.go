package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

type MysteryTrader struct {
	ID        int64              `json:"id,omitempty"`
	GameID    int64              `json:"gameId,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty"`
	X         float64            `json:"x,omitempty"`
	Y         float64            `json:"y,omitempty"`
	Name      string             `json:"name,omitempty"`
	Num       int                `json:"num,omitempty"`
	HeadingX  float64            `json:"headingX,omitempty"`
	HeadingY  float64            `json:"headingY,omitempty"`
	WarpSpeed int                `json:"warpSpeed,omitempty"`
	Spec      *MysteryTraderSpec `json:"spec,omitempty"`
}

// we json serialize these types with custom Scan/Value methods
type MysteryTraderSpec cs.MysteryTraderSpec

// db serializer to serialize this to JSON
func (item *MysteryTraderSpec) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *MysteryTraderSpec) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// get a mysteryTrader by id
func (c *client) GetMysteryTrader(id int64) (*cs.MysteryTrader, error) {
	item := MysteryTrader{}
	if err := c.db.Get(&item, "SELECT * FROM mysteryTraders WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	mysteryTrader := c.converter.ConvertMysteryTrader(&item)
	return mysteryTrader, nil
}

func (c *client) getMysteryTradersForGame(gameID int64) ([]*cs.MysteryTrader, error) {

	items := []MysteryTrader{}
	if err := c.db.Select(&items, `SELECT * FROM mysteryTraders WHERE gameId = ?`, gameID); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.MysteryTrader{}, nil
		}
		return nil, err
	}

	results := make([]*cs.MysteryTrader, len(items))
	for i := range items {
		results[i] = c.converter.ConvertMysteryTrader(&items[i])
	}

	return results, nil
}

// create a new game
func (c *client) createMysteryTrader(mysteryTrader *cs.MysteryTrader, tx SQLExecer) error {
	item := c.converter.ConvertGameMysteryTrader(mysteryTrader)
	result, err := tx.NamedExec(`
	INSERT INTO mysteryTraders (
		createdAt,
		updatedAt,
		gameId,
		x,
		y,
		name,
		num,
		headingX,
		headingY,
		warpSpeed,
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
		:headingX,
		:headingY,
		:warpSpeed,
		:spec
	)
	`, item)

	if err != nil {
		return err
	}

	// update the id of our passed in game
	mysteryTrader.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) UpdateMysteryTrader(mysteryTrader *cs.MysteryTrader) error {
	return c.updateMysteryTrader(mysteryTrader, c.db)
}

// update an existing mysteryTrader
func (c *client) updateMysteryTrader(mysteryTrader *cs.MysteryTrader, tx SQLExecer) error {

	item := c.converter.ConvertGameMysteryTrader(mysteryTrader)

	if _, err := tx.NamedExec(`
	UPDATE mysteryTraders SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		x = :x,
		y = :y,
		name = :name,
		num = :num,
		headingX = :headingX,
		headingY = :headingY,
		warpSpeed = :warpSpeed,
		spec = :spec
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

func (c *client) deleteMysteryTrader(mysteryTraderID int64, tx SQLExecer) error {
	if _, err := tx.Exec("DELETE FROM mysteryTraders where id = ?", mysteryTraderID); err != nil {
		return fmt.Errorf("delete mysteryTrader %d %w", mysteryTraderID, err)
	}
	return nil
}
