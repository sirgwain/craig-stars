package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

type MineralPacket struct {
	ID              int64     `json:"id,omitempty"`
	GameID          int64     `json:"gameId,omitempty"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
	UpdatedAt       time.Time `json:"updatedAt,omitempty"`
	X               float64   `json:"x,omitempty"`
	Y               float64   `json:"y,omitempty"`
	Name            string    `json:"name,omitempty"`
	Num             int       `json:"num,omitempty"`
	PlayerNum       int       `json:"playerNum,omitempty"`
	Tags            *Tags      `json:"tags,omitempty"`
	TargetPlanetNum int       `json:"targetPlanetNum,omitempty"`
	Ironium         int       `json:"ironium,omitempty"`
	Boranium        int       `json:"boranium,omitempty"`
	Germanium       int       `json:"germanium,omitempty"`
	SafeWarpSpeed   int       `json:"safeWarpSpeed,omitempty"`
	WarpSpeed       int       `json:"warpSpeed,omitempty"`
	ScanRange       int       `json:"scanRange,omitempty"`
	ScanRangePen    int       `json:"scanRangePen,omitempty"`
	HeadingX        float64   `json:"headingX,omitempty"`
	HeadingY        float64   `json:"headingY,omitempty"`
}

// get a mineralPacket by id
func (c *client) GetMineralPacket(id int64) (*cs.MineralPacket, error) {
	item := MineralPacket{}
	if err := c.reader.Get(&item, "SELECT * FROM mineralPackets WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	mineralPacket := c.converter.ConvertMineralPacket(&item)
	return mineralPacket, nil
}

func (c *client) GetMineralPacketByNum(gameID int64, playerNum int, num int) (*cs.MineralPacket, error) {

	item := MineralPacket{}
	if err := c.reader.Get(&item, `SELECT * FROM mineralPackets WHERE gameId = ? AND playerNum = ? AND num = ?`, gameID, playerNum, num); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	mineralPacket := c.converter.ConvertMineralPacket(&item)
	return mineralPacket, nil

}

func (c *client) GetMineralPacketsForPlayer(gameID int64, playerNum int) ([]*cs.MineralPacket, error) {

	items := []MineralPacket{}
	if err := c.reader.Select(&items, `SELECT * FROM mineralPackets WHERE gameId = ? AND playerNum = ?`, gameID, playerNum); err != nil {
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

func (c *client) getMineralPacketsForGame(gameID int64) ([]*cs.MineralPacket, error) {

	items := []MineralPacket{}
	if err := c.reader.Select(&items, `SELECT * FROM mineralPackets WHERE gameId = ?`, gameID); err != nil {
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
func (c *client) createMineralPacket(mineralPacket *cs.MineralPacket) error {
	item := c.converter.ConvertGameMineralPacket(mineralPacket)
	result, err := c.writer.NamedExec(`
	INSERT INTO mineralPackets (
		createdAt,
		updatedAt,
		gameId,
		x,
		y,
		name,
		num,
		playerNum,
		tags,
		targetPlanetNum,
		ironium,
		boranium,
		germanium,
		safeWarpSpeed,
		warpSpeed,
		scanRange,
		scanRangePen,
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
		:tags,
		:targetPlanetNum,
		:ironium,
		:boranium,
		:germanium,
		:safeWarpSpeed,
		:warpSpeed,
		:scanRange,
		:scanRangePen,
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

// update an existing mineralPacket
func (c *client) UpdateMineralPacket(mineralPacket *cs.MineralPacket) error {

	item := c.converter.ConvertGameMineralPacket(mineralPacket)

	if _, err := c.writer.NamedExec(`
	UPDATE mineralPackets SET
		updatedAt = CURRENT_TIMESTAMP,
		gameId = :gameId,
		x = :x,
		y = :y,
		name = :name,
		num = :num,
		playerNum = :playerNum,
		tags = :tags,
		targetPlanetNum = :targetPlanetNum,
		ironium = :ironium,
		boranium = :boranium,
		germanium = :germanium,
		safeWarpSpeed = :safeWarpSpeed,
		warpSpeed = :warpSpeed,
		scanRange = :scanRange,
		scanRangePen = :scanRangePen,
		headingX = :headingX,
		headingY = :headingY
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

func (c *client) deleteMineralPacket(mineralPacketID int64) error {
	if _, err := c.writer.Exec("DELETE FROM mineralPackets where id = ?", mineralPacketID); err != nil {
		return fmt.Errorf("delete mineralPacket %d %w", mineralPacketID, err)
	}
	return nil
}
