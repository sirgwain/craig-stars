package cs

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

type User struct {
	DBObject
	Username      string     `json:"username" header:"Username"`
	Password      string     `json:"password"`
	Email         string     `json:"email"`
	Role          UserRole   `json:"role"`
	Banned        bool       `json:"banned"`
	Verified      bool       `json:"verified"`
	GameID        int64      `json:"gameId,omitempty"`
	PlayerNum     int        `json:"playerNum,omitempty"`
	LastLogin     *time.Time `json:"lastLogin,omitempty"`
	DiscordID     *string    `json:"discordId,omitempty"`
	DiscordAvatar *string    `json:"discordAvatar,omitempty"`
}

type UserRole string

const (
	RoleNone  UserRole = ""
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
	RoleGuest UserRole = "guest"
)

func NewUser(username string, password string, email string, role UserRole) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{Username: username, Password: hashedPassword, Email: email, Role: role}, nil
}

func NewDiscordUser(username string, discordID string, discordAvatar string) (*User, error) {
	return &User{
		Username:      username,
		Role:          RoleUser,
		DiscordID:     &discordID,
		DiscordAvatar: &discordAvatar,
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
	return u.DiscordID != nil && len(*u.DiscordID) > 0
}

func (u *User) IsGuest() bool {
	return u.Role == RoleGuest
}

func (u *User) ComparePassword(password string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeHash(u.Password)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// password hashing from https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// Establish the parameters to use for Argon2.
var hashParams = &params{
	memory:      64 * 1024,
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func hashPassword(password string) (encodedHash string, err error) {
	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(hashParams.saltLength)
	if err != nil {
		return "", err
	}

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash := argon2.IDKey([]byte(password), salt, hashParams.iterations, hashParams.memory, hashParams.parallelism, hashParams.keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, hashParams.memory, hashParams.iterations, hashParams.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func decodeHash(encodedHash string) (p *params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
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
