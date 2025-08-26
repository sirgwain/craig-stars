package server

import (
	"crypto/rand"
	"crypto/sha1" //nolint
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/rest"
	"github.com/golang-jwt/jwt"

	"github.com/go-pkgz/auth/logger"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
)

const (
	// MaxHTTPBodySize defines max http body size
	MaxHTTPBodySize = 1024 * 1024
)

// GuestHandler implements non-oauth2 provider authorizing user in traditional way with storage
// with users and hashes
type GuestHandler struct {
	logger.L
	HashChecker  HashChecker
	ProviderName string
	TokenService provider.TokenService
	Issuer       string
}

// HashChecker defines interface to check credentials
type HashChecker interface {
	Check(hash string) (username string, attributes map[string]interface{}, err error)
}

// HashCheckerFunc type is an adapter to allow the use of ordinary functions as CredsChecker.
type HashCheckerFunc func(hash string) (username string, attributes map[string]interface{}, err error)

// Check calls f(user,passwd)
func (f HashCheckerFunc) Check(hash string) (username string, attributes map[string]interface{}, err error) {
	return f(hash)
}

// credentials holds user credentials
type credentials struct {
	Hash     string `json:"hash"`
	Audience string `json:"aud"`
}

// AddGuestProvider adds provider with guest invite link checks against a data store
// it doesn't do any handshake and uses provided hashChecker to verify a user exists with a hash from the request
func AddGuestProvider(s *auth.Service, issuer string, logger logger.L, name string, hashChecker HashChecker) {
	dh := GuestHandler{
		L:            logger,
		ProviderName: name,
		Issuer:       issuer,
		TokenService: s.TokenService(),
		HashChecker:  hashChecker,
	}
	s.AddCustomHandler(dh)
}

// Name of the handler
func (p GuestHandler) Name() string { return p.ProviderName }

// LoginHandler checks "hash" against data store and makes jwt if all passed.
//
// POST /something?sess[0|1]
// Accepts application/x-www-form-urlencoded or application/json encoded requests.
//
// application/x-www-form-urlencoded body example:
// user=name&passwd=xyz&aud=bar
//
// application/json body example:
//
//	{
//	  "hash": "123xyz",
//	  "aud": "bar",
//	}
func (p GuestHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	creds, err := p.getCredentials(w, r)
	if err != nil {
		rest.SendErrorJSON(w, r, p.L, http.StatusBadRequest, err, "failed to parse credentials")
		return
	}
	sessOnly := r.URL.Query().Get("sess") == "1"
	if p.HashChecker == nil {
		rest.SendErrorJSON(w, r, p.L, http.StatusInternalServerError,
			fmt.Errorf("no credential checker"), "no credential checker")
		return
	}
	username, attrs, err := p.HashChecker.Check(creds.Hash)
	if err != nil {
		rest.SendErrorJSON(w, r, p.L, http.StatusInternalServerError, err, "failed to check user credentials")
		return
	}
	if username == "" {
		rest.SendErrorJSON(w, r, p.L, http.StatusForbidden, nil, "incorrect user or password")
		return
	}

	userID := p.ProviderName + "_" + token.HashID(sha1.New(), username)

	u := token.User{
		Name:       username,
		ID:         userID,
		Role:       "guest",
		Attributes: attrs,
	}

	cid, err := randToken()
	if err != nil {
		rest.SendErrorJSON(w, r, p.L, http.StatusInternalServerError, err, "can't make token id")
		return
	}

	claims := token.Claims{
		User: &u,
		StandardClaims: jwt.StandardClaims{
			Id:       cid,
			Issuer:   p.Issuer,
			Audience: creds.Audience,
		},
		SessionOnly: sessOnly,
	}

	if _, err = p.TokenService.Set(w, claims); err != nil {
		rest.SendErrorJSON(w, r, p.L, http.StatusInternalServerError, err, "failed to set token")
		return
	}
	rest.RenderJSON(w, claims.User)
}

// getCredentials extracts user and password from request
func (p GuestHandler) getCredentials(w http.ResponseWriter, r *http.Request) (credentials, error) {

	if r.Method != "POST" {
		return credentials{}, fmt.Errorf("method %s not supported", r.Method)
	}

	if r.Body != nil {
		r.Body = http.MaxBytesReader(w, r.Body, MaxHTTPBodySize)
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "" {
		mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			return credentials{}, err
		}
		contentType = mt
	}

	// POST with json body
	if contentType == "application/json" {
		var creds credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			return credentials{}, fmt.Errorf("failed to parse request body: %w", err)
		}
		return creds, nil
	}

	// POST with form
	if err := r.ParseForm(); err != nil {
		return credentials{}, fmt.Errorf("failed to parse request: %w", err)
	}

	return credentials{
		Hash:     r.Form.Get("hash"),
		Audience: r.Form.Get("aud"),
	}, nil
}

// AuthHandler doesn't do anything for guest login as it has no callbacks
func (p GuestHandler) AuthHandler(http.ResponseWriter, *http.Request) {}

// LogoutHandler - GET /logout
func (p GuestHandler) LogoutHandler(w http.ResponseWriter, _ *http.Request) {
	p.TokenService.Reset(w)
}

// copied from go-pkgz/auth
func randToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("can't get random: %w", err)
	}
	s := sha1.New()
	if _, err := s.Write(b); err != nil {
		return "", fmt.Errorf("can't write randoms to sha1: %w", err)
	}
	return fmt.Sprintf("%x", s.Sum(nil)), nil
}
