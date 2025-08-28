package cs

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
)

// A User corresponds to a human logged into craig-stars. This might not actually belong
// in the cs package, but I didn't feel like breaking it out into a new package.
type User struct {
	DBObject
	UserSettings
	Username      string     `json:"username" header:"Username"`
	Password      string     `json:"password"`
	Email         string     `json:"email"`
	Role          UserRole   `json:"role"`
	Banned        bool       `json:"banned"`
	Verified      bool       `json:"verified"`
	GameID        int64      `json:"gameId,omitempty"`
	PlayerNum     int        `json:"playerNum,omitempty"`
	LastLogin     *time.Time `json:"lastLogin,omitempty"`
	DiscordID     string     `json:"discordId,omitempty"`
	DiscordAvatar string     `json:"discordAvatar,omitempty"`
}

type UserSettings struct {
	DiscordWebhookURL string `json:"discordWebhookUrl,omitempty"`
}

type UserRole string

const (
	RoleNone  UserRole = ""
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
	RoleGuest UserRole = "guest"
)

func NewUser(username string, password string, email string, role UserRole) *User {
	return &User{Username: username, Password: password, Email: email, Role: role}
}

func NewDiscordUser(username string, discordID string, discordAvatar string) (*User, error) {
	return &User{
		Username:      username,
		Role:          RoleUser,
		DiscordID:     discordID,
		DiscordAvatar: discordAvatar,
	}, nil
}

// create a new guest user for a game. This user will have an invite link generated
// for their username
func NewGuestUser(username string, gameID int64, playerNum int) *User {
	return &User{Username: username, GameID: gameID, PlayerNum: playerNum, Role: RoleGuest}
}

// generate an invite hash for this user based on id and store it in the user password
func (u *User) GenerateHash(salt string) {
	hasher := sha1.New()
	hasher.Write([]byte(fmt.Sprintf("%d-%s-%s", u.ID, u.Username, salt)))
	sha := hex.EncodeToString(hasher.Sum(nil))
	hash := sha[10:]

	u.Password = hash
}

func (u *User) IsDiscordUser() bool {
	return u.DiscordID != ""
}

func (u *User) IsGuest() bool {
	return u.Role == RoleGuest
}

func UserRoleFromString(role string) UserRole {
	switch role {
	case "user":
		return RoleUser
	case "admin":
		return RoleAdmin
	case "guest":
		return RoleGuest
	}

	return RoleNone
}
