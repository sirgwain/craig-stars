package db

import (
	"database/sql"
	"time"

	"github.com/sirgwain/craig-stars/game"
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Username  string    `json:"username,omitempty" header:"Username"`
	Password  string    `json:"password,omitempty"`
	Role      string    `json:"role,omitempty"`
}

func (c *client) GetUsers() ([]game.User, error) {

	// don't include password in bulk select
	items := []User{}
	if err := c.db.Select(&items, `
	SELECT 
		createdAt,
		updatedAt,
		username,
		role
	FROM Users
	`); err != nil {
		if err == sql.ErrNoRows {
			return []game.User{}, nil
		}
		return nil, err
	}

	return c.converter.ConvertUsers(items), nil
}

// get a user by id
func (c *client) GetUser(id int64) (*game.User, error) {
	item := User{}
	if err := c.db.Get(&item, "SELECT * FROM users WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user := c.converter.ConvertUser(item)
	return &user, nil
}

// get a user by id
func (c *client) GetUserByUsername(username string) (*game.User, error) {
	item := User{}
	if err := c.db.Get(&item, "SELECT * FROM users WHERE username = ?", username); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user := c.converter.ConvertUser(item)
	return &user, nil
}

// create a new user
func (c *client) CreateUser(user *game.User) error {
	item := c.converter.ConvertGameUser(user)
	result, err := c.db.NamedExec(`
	INSERT INTO users (
		createdAt,
		updatedAt,
		username,
		password,
		role
	) 
	VALUES (
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP,
		:username,
		:password,
		:role
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
func (c *client) UpdateUser(user *game.User) error {

	item := c.converter.ConvertGameUser(user)

	if _, err := c.db.NamedExec(`
	UPDATE users SET 
		updatedAt = CURRENT_TIMESTAMP,
		username = :username,
		password = :password,
		role = :role
	WHERE id = :id
	`, item); err != nil {
		return err
	}

	return nil
}

// delete a user by id
func (c *client) DeleteUser(id int64) error {
	if _, err := c.db.Exec("DELETE FROM users WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}
