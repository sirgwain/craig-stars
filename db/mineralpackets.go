package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	generated "github.com/sirgwain/craig-stars/db/generated"
)

// get a mineralpacket by id
func (c *client) GetMineralPacket(ctx context.Context, id int64) (*cs.MineralPacket, error) {
	item, err := c.reader.GetMineralPacket(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMineralPacket(item), nil
}

func (c *client) GetMineralPacketByNum(ctx context.Context, gameID int64, playerNum int, num int) (*cs.MineralPacket, error) {

	item, err := c.reader.GetMineralPacketByNum(ctx, generated.GetMineralPacketByNumParams{
		GameID:    gameID,
		PlayerNum: int64(playerNum),
		Num:       int64(num),
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMineralPacket(item), nil

}

func (c *client) getMineralPacketsForGame(ctx context.Context, gameID int64) ([]*cs.MineralPacket, error) {
	items, err := c.reader.GetMineralPacketsForGame(ctx, gameID)

	if err == sql.ErrNoRows {
		return []*cs.MineralPacket{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMineralPackets(items), nil
}

func (c *client) GetMineralPacketsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.MineralPacket, error) {
	items, err := c.reader.GetMineralPacketsForPlayer(ctx, generated.GetMineralPacketsForPlayerParams{
		GameID:    gameID,
		PlayerNum: int64(playerNum),
	})

	if err == sql.ErrNoRows {
		return []*cs.MineralPacket{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertMineralPackets(items), nil
}

func (c *client) SaveMineralPacket(ctx context.Context, mineralPacket *cs.MineralPacket) error {
	if mineralPacket.ID == 0 {
		result, err := c.writer.CreateMineralPacket(ctx, c.converter.ConvertGameMineralPacketToCreateParams(mineralPacket))
		if err != nil {
			return err
		}
		mineralPacket.ID = result
	} else {
		_, err := c.writer.UpdateMineralPacket(ctx, c.converter.ConvertGameMineralPacketToUpdateParams(mineralPacket))
		if err != nil {
			return err
		}
	}

	return nil
}

// delete a mineralpacket by id
func (c *client) DeleteMineralPacket(ctx context.Context, id int64) error {
	_, err := c.writer.DeleteMineralPacket(ctx, id)
	return err
}
