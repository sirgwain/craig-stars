package server

import (
	"encoding/gob"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

const userkey = "user"

type sessionUser struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	Role     cs.Role `json:"role"`
}

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	// make sure we can serialize this to a cookie
	gob.Register(sessionUser{})
}

// authRequired is a simple middleware to check the session
func (s *server) authRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}

// login is a handler that parses a form and checks for specific data
func (s *server) login(c *gin.Context) {
	session := sessions.Default(c)
	var creds creds

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&creds); err != nil {
		return
	}

	// Validate form input
	if strings.Trim(creds.Username, " ") == "" || strings.Trim(creds.Password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields can't be empty"})
		return
	}

	user, err := s.db.GetUserByUsername(creds.Username)
	if err != nil {
		log.Error().Err(err).Str("Username", creds.Username).Msg("get user from database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user from database"})
		return
	}

	// Check for username and password match, usually from a database
	if user == nil || user.Password != creds.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Save the username in the session
	sesionUser := &sessionUser{ID: user.ID, Username: user.Username, Role: user.Role}
	session.Set(userkey, sesionUser)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": " save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

func (s *server) logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		log.Error().Err(err).Msg("save session")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (s *server) GetSessionUser(c *gin.Context) *sessionUser {
	session := sessions.Default(c)
	user := session.Get(userkey).(sessionUser)

	return &user
}
func (s *server) me(c *gin.Context) {
	// session := sessions.Default(c)
	// user := session.Get(userkey)
	user := s.GetSessionUser(c)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, s.GetSessionUser(c))
}

// get users
func (s *server) users(c *gin.Context) {
	user := s.GetSessionUser(c)

	if user == nil {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	if user.Role != cs.RoleAdmin {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	users, err := s.db.GetUsers()
	if err != nil {
		log.Error().Err(err).Msg("load users")
		c.JSON(http.StatusBadRequest, gin.H{"error": " load users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
