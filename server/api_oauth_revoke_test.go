package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func postRevoke(s *server, form url.Values) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/oauth/revoke", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.apiOAuthRevokeHandler(rec, req)
	return rec
}

func TestAPIOAuthRevoke(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	cfg.Database.Filename = filepath.Join(t.TempDir(), "revoke.db")
	dbConn := db.NewConn()
	require.NoError(t, dbConn.Connect(ctx, cfg))
	t.Cleanup(func() { dbConn.Close() })
	s := &server{db: dbConn, apiOAuth: newAPIOAuthStore()}

	user, err := dbConn.NewReadWriteClient().CreateUser(ctx, &cs.User{Username: "revoker", Role: cs.RoleUser})
	require.NoError(t, err)
	issue := func() string {
		raw, prefix, err := newRawAPIToken()
		require.NoError(t, err)
		_, err = dbConn.NewReadWriteClient().CreateAPIToken(ctx, &cs.APIToken{
			UserID:      user.ID,
			Name:        apiTokenDefaultName,
			TokenPrefix: prefix,
			Scope:       apiTokenScopePlayer,
			ExpiresAt:   time.Now().Add(time.Hour),
		}, hashAPIToken(raw))
		require.NoError(t, err)
		return raw
	}
	revoked, kept := issue(), issue()

	_, err = s.userFromAPIToken(ctx, revoked)
	require.NoError(t, err, "token works before it is revoked")

	rec := postRevoke(s, url.Values{"token": {revoked}})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	_, err = s.userFromAPIToken(ctx, revoked)
	assert.Error(t, err, "a revoked token is rejected")
	_, err = s.userFromAPIToken(ctx, kept)
	assert.NoError(t, err, "revoking one token leaves the user's others alone")

	// revoking again, or a token that never existed, still answers 200 (RFC 7009)
	assert.Equal(t, http.StatusOK, postRevoke(s, url.Values{"token": {revoked}}).Code)
	assert.Equal(t, http.StatusOK, postRevoke(s, url.Values{"token": {apiTokenPrefix + "nope"}}).Code)
	assert.Equal(t, http.StatusOK, postRevoke(s, url.Values{"token": {"not-ours"}}).Code)
	assert.Equal(t, http.StatusBadRequest, postRevoke(s, url.Values{}).Code)
}
