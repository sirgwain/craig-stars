package server

import (
	"context"
	"encoding/gob"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type sessionUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// get the user from the context
func (s *server) contextUser(r *http.Request) sessionUser {
	return r.Context().Value(keyUser).(sessionUser)
}

func (s *server) mustGetUser(w http.ResponseWriter, r *http.Request) sessionUser {
	userInfo, err := token.GetUserInfo(r)
	if err != nil {
		panic("failed to load user")
	}

	userID, err := strconv.ParseInt(userInfo.StrAttr(databaseIDAttr), 10, 64)
	if err != nil {
		panic("failed to load user")
	}

	return sessionUser{
		ID:       userID,
		Username: userInfo.Name,
		Role:     userInfo.Role,
	}
}

func me_chi(w http.ResponseWriter, r *http.Request) {

	userInfo, err := token.GetUserInfo(r)
	if err != nil {
		log.Printf("failed to get user info, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, err := strconv.ParseInt(userInfo.StrAttr(databaseIDAttr), 10, 64)
	if err != nil {
		log.Printf("failed to get user info, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := sessionUser{
		ID:       userID,
		Username: userInfo.Name,
		Role:     userInfo.Role,
	}

	rest.RenderJSON(w, res)
}

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	// make sure we can serialize this to a cookie
	gob.Register(sessionUser{})
}

func (s *server) userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.mustGetUser(w, r)

		ctx := context.WithValue(r.Context(), keyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// create a new user from a token
func (s *server) createNewUser(tokenUser *token.User) (*cs.User, error) {
	user, err := cs.NewUser(tokenUser.Name, "", "", cs.RoleUser)
	if err != nil {
		log.Error().Err(err).Str("Username", user.Username).Msg("failed to create new user")
		return nil, err
	}
	err = s.db.CreateUser(user)
	if err != nil {
		log.Error().Err(err).Str("Username", user.Username).Msg("failed to create new user")
		return nil, err
	}
	log.Info().Str("Username", user.Username).Int64("ID", user.ID).Msg("created new user from token")

	// create a new test race
	race := cs.Humanoids()
	race.UserID = user.ID
	err = s.db.CreateRace(&race)
	log.Info().Str("Username", user.Username).Int64("ID", user.ID).Msg("created new race for user")
	if err != nil {
		return nil, err
	}
	return user, nil
}

// authRequired is a simple middleware to check the session
func (s *server) authRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(keyUser)
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
	match, err := user.ComparePassword(creds.Password)
	if err != nil {
		log.Error().Err(err).Str("Username", creds.Username).Msg("hash password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check credentials"})
		return
	}

	if user == nil || !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Save the username in the session
	sesionUser := &sessionUser{ID: user.ID, Username: user.Username, Role: user.Role}
	session.Set(keyUser, sesionUser)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": " save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

func (s *server) logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(keyUser)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(keyUser)
	if err := session.Save(); err != nil {
		log.Error().Err(err).Msg("save session")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (s *server) GetSessionUser(c *gin.Context) *sessionUser {
	session := sessions.Default(c)
	user := session.Get(keyUser).(sessionUser)

	return &user
}
func (s *server) me(c *gin.Context) {
	// session := sessions.Default(c)
	// user := session.Get(keyUser)
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
