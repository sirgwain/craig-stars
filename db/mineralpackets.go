package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

type MineralPacket struct {
	ID                int64     `json:"id,omitempty"`
	GameID            int64     `json:"gameId,omitempty"`
	CreatedAt         time.Time `json:"createdAt,omitempty"`
	UpdatedAt         time.Time `json:"updatedAt,omitempty"`
	X                 float64   `json:"x,omitempty"`
	Y                 float64   `json:"y,omitempty"`
	Name              string    `json:"name,omitempty"`
	Num               int       `json:"num,omitempty"`
	PlayerNum         int       `json:"playerNum,omitempty"`
	Tags              Tags      `json:"tags,omitempty"`
	TargetPlanetNum   int       `json:"targetPlanetNum,omitempty"`
	Ironium           int       `json:"ironium,omitempty"`
	Boranium          int       `json:"boranium,omitempty"`
	Germanium         int       `json:"germanium,omitempty"`
	SafeWarpSpeed     int       `json:"safeWarpSpeed,omitempty"`
	WarpFactor        int       `json:"warpFactor,omitempty"`
	HeadingX          float64   `json:"headingX,omitempty"`
	HeadingY          float64   `json:"headingY,omitempty"`
}

// get a mineralPacket by id
func (c *client) GetMineralPacket(id int64) (*cs.MineralPacket, error) {
	item := MineralPacket{}
	if err := c.db.Get(&item, "SELECT * FROM mineralPackets WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	mineralPacket := c.converter.ConvertMineralPacket(&item)
	return mineralPacket, nil
}

func (c *client) getMineralPacketsForGame(gameID int64) ([]*cs.MineralPacket, error) {

	items := []MineralPacket{}
	if err := c.db.Select(&items, `SELECT * FROM mineralPackets WHERE gameId = ?`, gameID); err != nil {
		if err == sql.ErrNoRows {
			return []*cs.MineralPacket{}, nil
		}
		return nil, err
	}

	results := make([]*cs.MineralPacket, len(items))
	for i := range items {
		results[i] = c.converter.ConvertMineralPacket(&items[i])
	}

	return results, nil
}

// create a new game
func (c *client) createMineralPacket(mineralPacket *cs.MineralPacket, tx SQLExecer) error {
	item := c.converter.ConvertGameMineralPacket(mineralPacket)
	result, err := tx.NamedExec(`
	INSERT INTO mineralPackets (
		createdAt,
		updatedAt,
		gameId,
		x,
		y,
		name,
		num,
		playerNum,
		targetPlanetNum,
		ironium,
		boranium,
		germanium,
		safeWarpSpeed,
		warpFactor,
		headingX,
		headingY
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
		:targetPlanetNum,
		:ironium,
		:boranium,
		:germanium,
		:safeWarpSpeed,
		:warpFactor,
		:headingX,
		:headingY
	)
	`, item)

	if err != nil {
		return err
	}

	// update the id of our passed in game
	mineralPacket.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) UpdateMineralPacket(mineralPacket *cs.MineralPacket) error {
	return c.updateMineralPacket(mineralPacket, c.db)
}

// update an existing mineralPacket
func (c *client) updateMineralPacket(mineralPacket *cs.MineralPacket, tx SQLExecer) error {

	item := c.converter.ConvertGameMineralPacket(mineralPacket)

	if _, err := tx.NamedExec(`
	UPDATE mineralPackets SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		x = :x,
		y = :y,
		name = :name,
		num = :num,
		playerNum = :playerNum,
		targetPlanetNum = :targetPlanetNum,
		ironium = :ironium,
		boranium = :boranium,
		germanium = :germanium,
		safeWarpSpeed = :safeWarpSpeed,
		warpFactor = :warpFactor,
		headingX = :headingX,
		headingY = :headingY
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

func (c *client) deleteMineralPacket(mineralPacketID int64, tx SQLExecer) error {
	if _, err := tx.Exec("DELETE FROM mineralPackets where id = ?", mineralPacketID); err != nil {
		return fmt.Errorf("delete mineralPacket %d %w", mineralPacketID, err)
	}
	return nil
}
