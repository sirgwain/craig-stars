package game

import (
	"errors"
	"time"
)

type User struct {
	ID        int64     `gorm:"primaryKey" json:"id" header:"ID" boltholdKey:"ID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Username  string    `json:"username" header:"Username"`
	Password  string    `json:"password"`
	Role      Role      `json:"role"`
}

type Role string

const (
	RoleAdmin Role = "Admin"
	RoleUser  Role = "User"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *Role) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *Role) Set(v string) error {
	switch v {
	case "Admin", "User":
		*e = Role(v)
		return nil
	default:
		return errors.New(`must be one of "Admin", or "User"`)
	}
}

// Type is only used in help text
func (e *Role) Type() string {
	return "Role"
}

func NewUser(username string, password string, role Role) *User {
	return &User{Username: username, Password: password, Role: role}
}
