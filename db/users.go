package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db/generated"
)

func (c *client) GetUsers(ctx context.Context) ([]cs.User, error) {
	items, err := c.reader.GetUsers(ctx)
	if err == sql.ErrNoRows {
		return []cs.User{}, nil
	}
	if err != nil {
		return nil, err
	}

	// don't return the password
	users := c.converter.ConvertUsers(items)
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

// get a user by id
func (c *client) GetUser(ctx context.Context, id int64) (*cs.User, error) {
	item, err := c.reader.GetUser(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user := c.converter.ConvertUser(item)
	return &user, nil
}

// get a user by id
func (c *client) GetUserByUsername(ctx context.Context, username string) (*cs.User, error) {
	item, err := c.reader.GetUserByUsername(ctx, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user := c.converter.ConvertUser(item)
	return &user, nil
}

func (c *client) GetGuestUser(ctx context.Context, hash string) (*cs.User, error) {
	item, err := c.reader.GetGuestUser(ctx, sql.NullString{Valid: true, String: hash})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user := c.converter.ConvertUser(item)
	return &user, nil
}

func (c *client) GetGuestUsersForGame(ctx context.Context, gameID int64) ([]cs.User, error) {

	items, err := c.reader.GetGuestUsersForGame(ctx, gameID)
	if err == sql.ErrNoRows {
		return []cs.User{}, nil
	}
	if err != nil {
		return nil, err
	}

	// don't return the password
	users := c.converter.ConvertUsers(items)
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (c *client) GetGuestUserForGame(ctx context.Context, gameID int64, playerNum int) (*cs.User, error) {

	item, err := c.reader.GetGetGuestUserForGame(ctx, generated.GetGetGuestUserForGameParams{Gameid: gameID, Playernum: int64(playerNum)})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user := c.converter.ConvertUser(item)
	return &user, nil
}

func (c *client) GetUsersForGame(ctx context.Context, gameID int64) ([]cs.User, error) {
	items, err := c.reader.GetUsersForGame(ctx, gameID)
	if err == sql.ErrNoRows {
		return []cs.User{}, nil
	}
	if err != nil {
		return nil, err
	}

	// don't return the password
	users := c.converter.ConvertUsers(items)
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

// create a new user
func (c *client) CreateUser(ctx context.Context, user *cs.User) (*cs.User, error) {

	result, err := c.writer.CreateUser(ctx, c.converter.ConvertGameUserToCreateParams(user))
	if err != nil {
		return nil, err
	}

	created := c.converter.ConvertUser(result)
	return &created, nil
}

// update an existing user
func (c *client) UpdateUser(ctx context.Context, user *cs.User) error {
	if err := c.writer.UpdateUser(ctx, c.converter.ConvertGameUserToUpdateParams(user)); err != nil {
		return err
	}

	return nil
}

// UpdateUserSettings updates a user's webhook url
func (c *client) UpdateUserSettings(ctx context.Context, user *cs.User) error {
	return c.writer.UpdateUserSettings(ctx,
		generated.UpdateUserSettingsParams{
			ID:                user.ID,
			Discordwebhookurl: sql.NullString{Valid: true, String: user.DiscordWebhookURL},
		},
	)
}

// delete a user by id
func (c *client) DeleteUser(ctx context.Context, id int64) error {
	return c.writer.DeleteUser(ctx, id)
}

// delete a guest user for a game
func (c *client) DeleteGameUser(ctx context.Context, gameID int64, playerNum int) error {
	return c.writer.DeleteGameGuestUser(ctx, generated.DeleteGameGuestUserParams{Gameid: gameID, Playernum: int64(playerNum)})
}

// delete guest users for a game
func (c *client) DeleteGameUsers(ctx context.Context, gameID int64) error {
	return c.writer.DeleteGameGuestUsers(ctx, gameID)
}
