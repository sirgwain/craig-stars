package db

import (
	"database/sql"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

type User struct {
	ID            int64      `json:"id" header:"ID"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	Username      string     `json:"username" header:"Username"`
	Password      string     `json:"password"`
	Email         string     `json:"email"`
	Role          string     `json:"role"`
	Banned        bool       `json:"banned"`
	Verified      bool       `json:"verified"`
	LastLogin     *time.Time `json:"lastLogin,omitempty"`
	DiscordID     *string    `json:"discordId,omitempty"`
	DiscordAvatar *string    `json:"discordAvatar,omitempty"`
}

func (c *client) GetUsers() ([]cs.User, error) {

	// don't include password in bulk select
	items := []User{}
	if err := c.reader.Select(&items, `
	SELECT 
		createdAt,
		updatedAt,
		username,
		email,
		role,
		banned,
		verified,
		lastLogin,
		discordId,
		discordAvatar
	FROM Users
	`); err != nil {
		if err == sql.ErrNoRows {
			return []cs.User{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertUsers(items), nil
}

// get a user by id
func (c *client) GetUser(id int64) (*cs.User, error) {
	item := User{}
	if err := c.reader.Get(&item, "SELECT * FROM users WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user := c.converter.ConvertUser(item)
	return &user, nil
}

// get a user by id
func (c *client) GetUserByUsername(username string) (*cs.User, error) {
	item := User{}
	if err := c.reader.Get(&item, "SELECT * FROM users WHERE username = ?", username); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user := c.converter.ConvertUser(item)
	return &user, nil
}

func (c *client) GetUsersForGame(gameID int64) ([]cs.User, error) {

	// don't include password in bulk select
	items := []User{}
	if err := c.reader.Select(&items, `
	SELECT 
		createdAt,
		updatedAt,
		username,
		email,
		role,
		banned,
		verified,
		lastLogin,
		discordId,
		discordAvatar
	FROM users WHERE id IN (SELECT userId FROM players p WHERE p.gameId = ?)
	`, gameID); err != nil {
		if err == sql.ErrNoRows {
			return []cs.User{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertUsers(items), nil
}

// create a new user
func (c *client) CreateUser(user *cs.User) error {
	item := c.converter.ConvertGameUser(user)

	result, err := c.writer.NamedExec(`
	INSERT INTO users (
		createdAt,
		updatedAt,
		username,
		password,
		email,
		role,
		banned,
		verified,
		lastLogin,
		discordId,
		discordAvatar
	) 
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:username,
		:password,
		:email,
		:role,
		:banned,
		:verified,
		:lastLogin,
		:discordId,
		:discordAvatar
	)
	`, item)

	if err != nil {
		return err
	}

	// update the id of our passed in user
	user.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// update an existing user
func (c *client) UpdateUser(user *cs.User) error {

	item := c.converter.ConvertGameUser(user)

	if _, err := c.writer.NamedExec(`
	UPDATE users SET 
		updatedAt = CURRENT_TIMESTAMP,
		username = :username,
		password = :password,
		username = :username,
		password = :password,
		email = :email,
		role = :role,
		banned = :banned,
		verified = :verified,
		lastLogin = :lastLogin,
		discordId = :discordId,
		discordAvatar = :discordAvatar
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// delete a user by id
func (c *client) DeleteUser(id int64) error {
	if _, err := c.writer.Exec("DELETE FROM users WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}
