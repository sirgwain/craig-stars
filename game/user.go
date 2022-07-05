package game

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id" header:"ID"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Username  string         `json:"username" header:"Username"`
	Password  string         `json:"password"`
	Role      Role           `json:"role"`
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
