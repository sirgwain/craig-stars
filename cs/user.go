package cs

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          int64     `json:"id" header:"ID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Username    string    `json:"username" header:"Username"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	Role        Role      `json:"role"`
	Verified    bool      `json:"verified"`
	VerifyToken string    `json:"verifyToken"`
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

func NewUser(username string, password string, email string, role Role) *User {
	verifyToken := uuid.New().String()
	return &User{Username: username, Password: password, Email: email, Role: role, VerifyToken: verifyToken}
}
