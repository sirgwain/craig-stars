package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

type MysteryTrader struct {
	ID            int64                    `json:"id,omitempty"`
	GameID        int64                    `json:"gameId,omitempty"`
	CreatedAt     time.Time                `json:"createdAt,omitempty"`
	UpdatedAt     time.Time                `json:"updatedAt,omitempty"`
	X             float64                  `json:"x,omitempty"`
	Y             float64                  `json:"y,omitempty"`
	Name          string                   `json:"name,omitempty"`
	Num           int                      `json:"num,omitempty"`
	Tags          *Tags                    `json:"tags,omitempty"`
	HeadingX      float64                  `json:"headingX,omitempty"`
	HeadingY      float64                  `json:"headingY,omitempty"`
	WarpSpeed     int                      `json:"warpSpeed,omitempty"`
	RequestedBoon int                      `json:"requestedBoon,omitempty"`
	DestinationX  float64                  `json:"destinationX,omitempty"`
	DestinationY  float64                  `json:"destinationY,omitempty"`
	RewardType    *MysteryTraderRewardType `json:"rewardType,omitempty"`
	Spec          *MysteryTraderSpec       `json:"spec,omitempty"`
}

// we json serialize these types with custom Scan/Value methods
type MysteryTraderSpec cs.MysteryTraderSpec
type MysteryTraderRewardType cs.MysteryTraderRewardType

// db serializer to serialize this to JSON
func (item *MysteryTraderSpec) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *MysteryTraderSpec) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// db serializer to serialize this to JSON
func (item *MysteryTraderRewardType) Value() (driver.Value, error) {
	return valueJSON(item)
}

// db deserializer to read this from JSON
func (item *MysteryTraderRewardType) Scan(src interface{}) error {
	return scanJSON(src, item)
}

// get a mysteryTrader by id
func (c *client) GetMysteryTrader(id int64) (*cs.MysteryTrader, error) {
	item := MysteryTrader{}
	if err := c.reader.Get(&item, "SELECT * FROM mysteryTraders WHERE id = ?", id); err != nil {
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
	if err := c.reader.Select(&items, `SELECT * FROM mysteryTraders WHERE gameId = ?`, gameID); err != nil {
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
func (c *client) createMysteryTrader(mysteryTrader *cs.MysteryTrader) error {
	item := c.converter.ConvertGameMysteryTrader(mysteryTrader)
	result, err := c.writer.NamedExec(`
	INSERT INTO mysteryTraders (
		createdAt,
		updatedAt,
		gameId,
		x,
		y,
		name,
		num,
		tags,
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
		:tags,
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

// update an existing mysteryTrader
func (c *client) updateMysteryTrader(mysteryTrader *cs.MysteryTrader) error {

	item := c.converter.ConvertGameMysteryTrader(mysteryTrader)

	if _, err := c.writer.NamedExec(`
	UPDATE mysteryTraders SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		x = :x,
		y = :y,
		name = :name,
		num = :num,
		tags = :tags,
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

func (c *client) deleteMysteryTrader(mysteryTraderID int64) error {
	if _, err := c.writer.Exec("DELETE FROM mysteryTraders where id = ?", mysteryTraderID); err != nil {
		return fmt.Errorf("delete mysteryTrader %d %w", mysteryTraderID, err)
	}
	return nil
}
