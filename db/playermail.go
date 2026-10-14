package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	generated "github.com/sirgwain/craig-stars/db/generated"
)

// get a player mail by id
func (c *client) GetPlayerMail(ctx context.Context, id int64) (*cs.PlayerMail, error) {
	item, err := c.reader.GetPlayerMail(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertPlayerMail(item), nil
}

// get all the mail a player wrote, and the mail other players wrote them that was delivered before year
func (c *client) GetPlayerMailForPlayer(ctx context.Context, gameID int64, playerNum int, year int) ([]*cs.PlayerMail, error) {
	items, err := c.reader.GetPlayerMailForPlayer(ctx, generated.GetPlayerMailForPlayerParams{
		GameID:    gameID,
		PlayerNum: int64(playerNum),
		Year:      int64(year),
	})
	if err == sql.ErrNoRows {
		return []*cs.PlayerMail{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertPlayerMails(items), nil
}

// get the mail players wrote in a year, which is delivered when that year's turn is generated
func (c *client) GetPendingPlayerMailForGame(ctx context.Context, gameID int64, year int) ([]*cs.PlayerMail, error) {
	items, err := c.reader.GetPlayerMailForGameYear(ctx, generated.GetPlayerMailForGameYearParams{
		GameID: gameID,
		Year:   int64(year),
	})
	if err == sql.ErrNoRows {
		return []*cs.PlayerMail{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertPlayerMails(items), nil
}

func (c *client) SavePlayerMail(ctx context.Context, mail *cs.PlayerMail) error {
	if mail.ID == 0 {
		result, err := c.writer.CreatePlayerMail(ctx, c.converter.ConvertGamePlayerMailToCreateParams(mail))
		if err != nil {
			return err
		}
		mail.ID = result
	} else {
		_, err := c.writer.UpdatePlayerMail(ctx, c.converter.ConvertGamePlayerMailToUpdateParams(mail))
		if err != nil {
			return err
		}
	}

	return nil
}

// delete a player mail by id
func (c *client) DeletePlayerMail(ctx context.Context, id int64) error {
	_, err := c.writer.DeletePlayerMail(ctx, id)
	return err
}
