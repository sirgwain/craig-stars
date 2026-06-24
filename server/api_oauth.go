package server

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-pkgz/auth/v2/token"
	"github.com/go-pkgz/rest"
	"github.com/sirgwain/craig-stars/cs"
)

const authCodeTTL = 5 * time.Minute

// apiOAuthStore keeps short-lived authorization codes between the browser
// authorization redirect and an API client's token exchange.
type apiOAuthStore struct {
	mu    sync.Mutex
	codes map[string]apiOAuthCode
}

// apiOAuthCode records the server-side state for one OAuth authorization code.
type apiOAuthCode struct {
	CodeHash      string
	UserID        int64
	ClientID      string
	Registered    bool
	RedirectURI   string
	Resource      string
	CodeChallenge string
	ExpiresAt     time.Time
	Used          bool
}

// tokenResponse is the OAuth response containing a craig-stars API token.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

func newAPIOAuthStore() *apiOAuthStore {
	return &apiOAuthStore{codes: map[string]apiOAuthCode{}}
}

// apiOAuthAuthorizeHandler starts browser authorization and redirects the
// client with a short-lived code that can be exchanged for an API token.
func (s *server) apiOAuthAuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := token.GetUserInfo(r)
	if err != nil {
		loginURL := "/api/auth/discord/login?from=" + url.QueryEscape(s.absoluteRequestURL(r))
		http.Redirect(w, r, loginURL, http.StatusFound)
		return
	}

	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	resource := r.URL.Query().Get("resource")
	codeChallenge := r.URL.Query().Get("code_challenge")
	registered, err := s.validateMCPRedirectURI(r.Context(), clientID, redirectURI)
	if err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, err, "invalid redirect_uri")
		return
	}
	if resource != "" && resource != s.mcpResourceURI() {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, errors.New("invalid resource"), "invalid resource")
		return
	}
	if codeChallenge == "" {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, errors.New("missing code_challenge"), "missing code_challenge")
		return
	}
	if method := r.URL.Query().Get("code_challenge_method"); method != "" && method != "S256" {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, errors.New("unsupported code_challenge_method"), "unsupported code_challenge_method")
		return
	}

	code, err := randomURLToken(32)
	if err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusInternalServerError, err, "failed to create code")
		return
	}
	tokenUser := newTokenUser(&userInfo)
	userID := tokenUser.databaseID()
	if userID == 0 {
		rest.SendErrorJSON(w, r, nil, http.StatusUnauthorized, errors.New("missing user id"), "missing user id")
		return
	}
	s.apiOAuth.put(apiOAuthCode{
		CodeHash:      hashOAuthSecret(code),
		UserID:        userID,
		ClientID:      clientID,
		Registered:    registered,
		RedirectURI:   redirectURI,
		Resource:      resource,
		CodeChallenge: codeChallenge,
		ExpiresAt:     time.Now().Add(authCodeTTL),
	})

	u, _ := url.Parse(redirectURI)
	q := u.Query()
	q.Set("code", code)
	if state := r.URL.Query().Get("state"); state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

// apiOAuthTokenHandler exchanges an authorization code and PKCE verifier for
// an opaque craig-stars API bearer token.
func (s *server) apiOAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	req, err := parseTokenRequest(r)
	if err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, err, "invalid token request")
		return
	}
	if req["grant_type"] != "authorization_code" {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, errors.New("unsupported grant_type"), "unsupported grant_type")
		return
	}

	code := req["code"]
	clientID := req["client_id"]
	redirectURI := req["redirect_uri"]
	resource := req["resource"]
	verifier := req["code_verifier"]
	if code == "" || redirectURI == "" || verifier == "" {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, errors.New("missing required token field"), "missing required token field")
		return
	}

	authCode, ok := s.apiOAuth.consume(code)
	if !ok || authCode.RedirectURI != redirectURI || !verifyPKCE(verifier, authCode.CodeChallenge) {
		rest.SendErrorJSON(w, r, nil, http.StatusUnauthorized, errors.New("invalid authorization code"), "invalid authorization code")
		return
	}
	if authCode.Registered && authCode.ClientID != clientID {
		rest.SendErrorJSON(w, r, nil, http.StatusUnauthorized, errors.New("invalid client_id"), "invalid client_id")
		return
	}
	if authCode.Resource != "" && authCode.Resource != resource {
		rest.SendErrorJSON(w, r, nil, http.StatusUnauthorized, errors.New("invalid resource"), "invalid resource")
		return
	}

	user, err := s.db.NewReadClient().GetUser(r.Context(), authCode.UserID)
	if err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusInternalServerError, err, "failed to load user")
		return
	}
	if user == nil || user.Banned {
		rest.SendErrorJSON(w, r, nil, http.StatusUnauthorized, errors.New("user inactive"), "user inactive")
		return
	}

	raw, prefix, err := newRawAPIToken()
	if err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusInternalServerError, err, "failed to create token")
		return
	}
	expiresAt := time.Now().Add(apiTokenDefaultExpiry)
	_, err = s.db.NewReadWriteClient().CreateAPIToken(r.Context(), &cs.APIToken{
		UserID:      user.ID,
		Name:        apiTokenDefaultName,
		TokenPrefix: prefix,
		Scope:       apiTokenScopePlayer,
		ExpiresAt:   expiresAt,
	}, hashAPIToken(raw))
	if err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusInternalServerError, err, "failed to save token")
		return
	}

	rest.RenderJSON(w, tokenResponse{
		AccessToken: raw,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(expiresAt).Seconds()),
		Scope:       apiTokenScopePlayer,
	})
}

// apiOAuthMetadataHandler publishes metadata for clients that discover the
// shared API-token authorization flow.
func (s *server) apiOAuthMetadataHandler(w http.ResponseWriter, r *http.Request) {
	base := strings.TrimRight(s.config.Auth.URL, "/")
	rest.RenderJSON(w, apiOAuthMetadata(base))
}

func apiOAuthMetadata(base string) map[string]any {
	return map[string]any{
		"issuer":                                base,
		"authorization_endpoint":                base + "/api/oauth/authorize",
		"token_endpoint":                        base + "/api/oauth/token",
		"registration_endpoint":                 base + "/api/mcp/oauth/register",
		"response_types_supported":              []string{"code"},
		"grant_types_supported":                 []string{"authorization_code"},
		"code_challenge_methods_supported":      []string{"S256"},
		"token_endpoint_auth_methods_supported": []string{"none"},
		"scopes_supported":                      []string{apiTokenScopePlayer},
	}
}

// absoluteRequestURL reconstructs the current request URL for the post-login
// return path, honoring proxy protocol headers.
func (s *server) absoluteRequestURL(r *http.Request) string {
	u := *r.URL
	u.Scheme = "https"
	if r.TLS == nil {
		u.Scheme = "http"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		u.Scheme = proto
	}
	u.Host = r.Host
	return u.String()
}

func (s *apiOAuthStore) put(code apiOAuthCode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pruneLocked(time.Now())
	s.codes[code.CodeHash] = code
}

func (s *apiOAuthStore) consume(code string) (apiOAuthCode, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.pruneLocked(now)
	hash := hashOAuthSecret(code)
	authCode, ok := s.codes[hash]
	if !ok || authCode.Used || !authCode.ExpiresAt.After(now) {
		return apiOAuthCode{}, false
	}
	authCode.Used = true
	s.codes[hash] = authCode
	return authCode, true
}

func (s *apiOAuthStore) pruneLocked(now time.Time) {
	for hash, code := range s.codes {
		if code.Used || !code.ExpiresAt.After(now) {
			delete(s.codes, hash)
		}
	}
}

// parseTokenRequest accepts either JSON or form-encoded OAuth token requests.
func parseTokenRequest(r *http.Request) (map[string]string, error) {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return req, nil
	}
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	return map[string]string{
		"grant_type":    r.Form.Get("grant_type"),
		"code":          r.Form.Get("code"),
		"client_id":     r.Form.Get("client_id"),
		"redirect_uri":  r.Form.Get("redirect_uri"),
		"resource":      r.Form.Get("resource"),
		"code_verifier": r.Form.Get("code_verifier"),
	}, nil
}

// validateRedirectURI accepts native craig-stars app callbacks in addition to
// local HTTP listeners used by desktop CLI clients.
func validateRedirectURI(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.Scheme == "cstars" {
		if u.Host == "" {
			return errors.New("cstars redirect_uri missing host")
		}
		return nil
	}
	return validateLoopbackRedirectURI(raw)
}

// validateLoopbackRedirectURI restricts callbacks to local HTTP listeners.
func validateLoopbackRedirectURI(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.Scheme != "http" {
		return errors.New("redirect_uri must use http for loopback")
	}
	host := u.Hostname()
	if host == "" {
		return errors.New("redirect_uri missing host")
	}
	ip := net.ParseIP(host)
	if host != "localhost" && (ip == nil || !ip.IsLoopback()) {
		return errors.New("redirect_uri must be loopback")
	}
	if u.Port() == "" {
		return errors.New("redirect_uri must include a port")
	}
	return nil
}

func verifyPKCE(verifier, expectedChallenge string) bool {
	sum := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])
	return subtle.ConstantTimeCompare([]byte(challenge), []byte(expectedChallenge)) == 1
}

func randomURLToken(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("make random token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashOAuthSecret(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hexString(sum[:])
}

func hexString(b []byte) string {
	const alphabet = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = alphabet[v>>4]
		out[i*2+1] = alphabet[v&0x0f]
	}
	return string(out)
}
