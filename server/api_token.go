package server

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	authmiddleware "github.com/go-pkgz/auth/v2/middleware"
	"github.com/go-pkgz/auth/v2/token"
	"github.com/sirgwain/craig-stars/cs"
)

// apiBearerContextKey stores the raw bearer token on a request context so MCP
// tool calls can forward it to the internal Connect clients.
type apiBearerContextKey struct{}

const (
	apiTokenPrefix        = "cstars_api_"
	apiTokenScopePlayer   = "api:player"
	apiTokenDefaultName   = "API client"
	apiTokenDefaultExpiry = 90 * 24 * time.Hour
)

// authAPIOrSession accepts either a craig-stars API bearer token or the normal
// browser session authentication used by the web app.
func (s *server) authAPIOrSession(m authmiddleware.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		sessionAuth := m.Auth(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := bearerToken(r)
			if raw == "" {
				sessionAuth.ServeHTTP(w, r)
				return
			}

			userInfo, err := s.userFromAPIToken(r.Context(), raw)
			if err != nil {
				http.Error(w, "invalid bearer token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), apiBearerContextKey{}, raw)
			next.ServeHTTP(w, token.SetUserInfo(r.WithContext(ctx), userInfo))
		})
	}
}

// authAPITokenOnly requires a valid API bearer token and advertises the MCP
// resource metadata endpoint when a client needs to authenticate.
func (s *server) authAPITokenOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resourceMetadata := s.mcpResourceMetadataURI()
		raw := bearerToken(r)
		if raw == "" {
			w.Header().Set("WWW-Authenticate", `Bearer resource_metadata="`+resourceMetadata+`"`)
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}
		userInfo, err := s.userFromAPIToken(r.Context(), raw)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Bearer resource_metadata="`+resourceMetadata+`"`)
			http.Error(w, "invalid bearer token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), apiBearerContextKey{}, raw)
		next.ServeHTTP(w, token.SetUserInfo(r.WithContext(ctx), userInfo))
	})
}

// bearerToken extracts the token value from an Authorization: Bearer header.
func bearerToken(r *http.Request) string {
	authz := r.Header.Get("Authorization")
	if authz == "" {
		return ""
	}
	scheme, value, ok := strings.Cut(authz, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return ""
	}
	return strings.TrimSpace(value)
}

// userFromAPIToken validates an opaque API token against the database and
// returns the authenticated auth token user.
func (s *server) userFromAPIToken(ctx context.Context, raw string) (token.User, error) {
	if !strings.HasPrefix(raw, apiTokenPrefix) {
		return token.User{}, errors.New("unsupported token prefix")
	}

	client := s.db.NewReadWriteClient()
	apiToken, err := client.GetAPITokenByHash(ctx, hashAPIToken(raw))
	if err != nil {
		return token.User{}, err
	}
	if apiToken == nil {
		return token.User{}, errors.New("token not found")
	}
	now := time.Now()
	if apiToken.RevokedAt != nil || !apiToken.ExpiresAt.After(now) || apiToken.Scope != apiTokenScopePlayer {
		return token.User{}, errors.New("token inactive")
	}

	user, err := client.GetUser(ctx, apiToken.UserID)
	if err != nil {
		return token.User{}, err
	}
	if user == nil || user.Banned {
		return token.User{}, errors.New("user inactive")
	}
	if err := client.TouchAPIToken(ctx, apiToken.ID); err != nil {
		return token.User{}, err
	}
	return tokenUserFromCSUser(user), nil
}

// tokenUserFromCSUser converts a domain user into the token user shape expected
// by existing request context helpers.
func tokenUserFromCSUser(user *cs.User) token.User {
	u := token.User{
		ID:   fmt.Sprintf("db_%d", user.ID),
		Name: user.Username,
		Role: string(user.Role),
	}
	tu := newTokenUser(&u)
	tu.setDatabaseID(user.ID)
	if user.DiscordID != "" {
		tu.setDiscordID(user.DiscordID)
	}
	if user.DiscordAvatar != "" {
		tu.setDiscordAvatar(user.DiscordAvatar)
	}
	if user.Role == cs.RoleAdmin {
		tu.SetAdmin(true)
	}
	return *tu.User
}

// newRawAPIToken creates a new opaque API token and the short prefix stored for
// user-facing token identification.
func newRawAPIToken() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("make api token: %w", err)
	}
	secret := base64.RawURLEncoding.EncodeToString(b)
	raw := apiTokenPrefix + secret
	return raw, raw[:min(len(raw), len(apiTokenPrefix)+8)], nil
}

// hashAPIToken returns the SHA-256 hex digest stored in the database instead of
// the raw bearer token.
func hashAPIToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
