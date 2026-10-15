package server

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/go-pkgz/auth/v2/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateLoopbackRedirectURI(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{name: "localhost with port", uri: "http://localhost:39123/callback"},
		{name: "ipv4 loopback with port", uri: "http://127.0.0.1:39123/callback"},
		{name: "ipv6 loopback with port", uri: "http://[::1]:39123/callback"},
		{name: "https rejected", uri: "https://localhost:39123/callback", wantErr: true},
		{name: "remote host rejected", uri: "http://example.com:39123/callback", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLoopbackRedirectURI(tt.uri)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateRedirectURI(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{name: "native app callback", uri: "cstars://oauth/callback"},
		{name: "native app callback with query", uri: "cstars://oauth/callback?source=ios"},
		{name: "native app missing host rejected", uri: "cstars:///oauth/callback", wantErr: true},
		{name: "other custom scheme rejected", uri: "other://oauth/callback", wantErr: true},
		{name: "loopback callback", uri: "http://127.0.0.1:39123/callback"},
		{name: "remote http callback rejected", uri: "http://example.com:39123/callback", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRedirectURI(tt.uri)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestVerifyPKCE(t *testing.T) {
	verifier := "test-verifier"
	sum := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])

	if !verifyPKCE(verifier, challenge) {
		t.Fatal("expected verifier to match challenge")
	}
	if verifyPKCE("wrong-verifier", challenge) {
		t.Fatal("expected wrong verifier to fail")
	}
}

func TestAPIOAuthMetadataUsesCanonicalEndpoints(t *testing.T) {
	metadata := apiOAuthMetadata("https://craig-stars.example")

	if got := metadata["authorization_endpoint"]; got != "https://craig-stars.example/api/oauth/authorize" {
		t.Fatalf("authorization_endpoint = %q", got)
	}
	if got := metadata["token_endpoint"]; got != "https://craig-stars.example/api/oauth/token" {
		t.Fatalf("token_endpoint = %q", got)
	}
	if got := metadata["registration_endpoint"]; got != "https://craig-stars.example/api/mcp/oauth/register" {
		t.Fatalf("registration_endpoint = %q", got)
	}
}

var consentIDPattern = regexp.MustCompile(`name="consent_id" value="([^"]+)"`)

func oauthTestUser(id int64) token.User {
	u := newTokenUser(&token.User{Name: "player"})
	u.setDatabaseID(id)
	return *u.User
}

// startOAuthConsent requests the authorize endpoint as user and returns the
// rendered consent id.
func startOAuthConsent(t *testing.T, s *server, user token.User) string {
	t.Helper()
	q := url.Values{
		"redirect_uri":          {"http://127.0.0.1:39123/callback"},
		"code_challenge":        {"challenge"},
		"code_challenge_method": {"S256"},
		"state":                 {"xyz"},
	}
	rec := httptest.NewRecorder()
	req := token.SetUserInfo(httptest.NewRequest(http.MethodGet, "/api/oauth/authorize?"+q.Encode(), nil), user)
	s.apiOAuthAuthorizeHandler(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	assert.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	assert.Contains(t, rec.Header().Get("Content-Security-Policy"), "frame-ancestors 'none'")
	assert.Contains(t, rec.Header().Get("Content-Security-Policy"), "form-action 'self' http://127.0.0.1:39123;")
	assert.Empty(t, s.apiOAuth.codes, "authorize GET must not issue a code")
	m := consentIDPattern.FindStringSubmatch(rec.Body.String())
	require.Len(t, m, 2)
	return m[1]
}

func postOAuthConsent(s *server, user token.User, consentID, action string) *httptest.ResponseRecorder {
	form := url.Values{"consent_id": {consentID}, "action": {action}}
	req := httptest.NewRequest(http.MethodPost, "/api/oauth/authorize", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	s.apiOAuthConsentHandler(rec, token.SetUserInfo(req, user))
	return rec
}

func TestAPIOAuthAuthorizeRequiresConsent(t *testing.T) {
	s := &server{apiOAuth: newAPIOAuthStore()}
	user := oauthTestUser(12)
	consentID := startOAuthConsent(t, s, user)

	rec := postOAuthConsent(s, user, consentID, "allow")
	require.Equal(t, http.StatusFound, rec.Code, rec.Body.String())
	loc, err := url.Parse(rec.Header().Get("Location"))
	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1:39123", loc.Host)
	assert.Equal(t, "xyz", loc.Query().Get("state"))
	code := loc.Query().Get("code")
	require.NotEmpty(t, code)

	authCode, ok := s.apiOAuth.consume(code)
	require.True(t, ok)
	assert.Equal(t, int64(12), authCode.UserID)
	assert.Equal(t, "challenge", authCode.CodeChallenge)

	// a consent can only be answered once
	rec = postOAuthConsent(s, user, consentID, "allow")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAPIOAuthConsentDeny(t *testing.T) {
	s := &server{apiOAuth: newAPIOAuthStore()}
	user := oauthTestUser(12)
	consentID := startOAuthConsent(t, s, user)

	rec := postOAuthConsent(s, user, consentID, "deny")
	require.Equal(t, http.StatusFound, rec.Code)
	loc, err := url.Parse(rec.Header().Get("Location"))
	require.NoError(t, err)
	assert.Equal(t, "access_denied", loc.Query().Get("error"))
	assert.Equal(t, "xyz", loc.Query().Get("state"))
	assert.Empty(t, loc.Query().Get("code"))
	assert.Empty(t, s.apiOAuth.codes)
}

func TestAPIOAuthConsentRejectsForgedApproval(t *testing.T) {
	s := &server{apiOAuth: newAPIOAuthStore()}
	consentID := startOAuthConsent(t, s, oauthTestUser(12))

	// a guessed consent id is rejected
	rec := postOAuthConsent(s, oauthTestUser(12), "not-a-real-consent", "allow")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// another user can't approve someone else's consent
	rec = postOAuthConsent(s, oauthTestUser(99), consentID, "allow")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Empty(t, s.apiOAuth.codes)
}
