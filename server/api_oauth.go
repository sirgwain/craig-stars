package server

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
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

const (
	authCodeTTL = 5 * time.Minute
	consentTTL  = 10 * time.Minute
)

// apiOAuthStore keeps short-lived authorization codes between the browser
// authorization redirect and an API client's token exchange.
type apiOAuthStore struct {
	mu       sync.Mutex
	codes    map[string]apiOAuthCode
	consents map[string]apiOAuthConsent
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

// apiOAuthConsent records a validated authorization request while the user
// decides whether to approve it on the consent page.
type apiOAuthConsent struct {
	IDHash        string
	UserID        int64
	ClientID      string
	Registered    bool
	RedirectURI   string
	Resource      string
	CodeChallenge string
	State         string
	ExpiresAt     time.Time
}

// tokenResponse is the OAuth response containing a craig-stars API token.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

func newAPIOAuthStore() *apiOAuthStore {
	return &apiOAuthStore{codes: map[string]apiOAuthCode{}, consents: map[string]apiOAuthConsent{}}
}

// apiOAuthAuthorizeHandler validates a browser authorization request and asks
// the signed-in user to approve it. Codes are only issued by the consent POST,
// so a link alone can never hand an API token to a third party.
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
	client, err := s.validateMCPRedirectURI(r.Context(), clientID, redirectURI)
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

	tokenUser := newTokenUser(&userInfo)
	userID := tokenUser.databaseID()
	if userID == 0 {
		rest.SendErrorJSON(w, r, nil, http.StatusUnauthorized, errors.New("missing user id"), "missing user id")
		return
	}

	consentID, err := randomURLToken(32)
	if err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusInternalServerError, err, "failed to create consent")
		return
	}
	s.apiOAuth.putConsent(apiOAuthConsent{
		IDHash:        hashOAuthSecret(consentID),
		UserID:        userID,
		ClientID:      clientID,
		Registered:    client != nil,
		RedirectURI:   redirectURI,
		Resource:      resource,
		CodeChallenge: codeChallenge,
		State:         r.URL.Query().Get("state"),
		ExpiresAt:     time.Now().Add(consentTTL),
	})

	clientName := ""
	if client != nil {
		clientName = client.ClientName
	}
	if clientName == "" {
		clientName = "An unnamed application"
	}
	redirectHost := redirectURI
	if u, err := url.Parse(redirectURI); err == nil && u.Host != "" {
		redirectHost = u.Host
	}

	// the consent page must never be framed, or a clickjacking page could trick
	// the user into pressing Allow
	w.Header().Set("Content-Security-Policy", consentContentSecurityPolicy(redirectURI))
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := consentTemplate.Execute(w, consentPage{
		ClientName:   clientName,
		RedirectHost: redirectHost,
		Username:     userInfo.Name,
		ConsentID:    consentID,
	}); err != nil {
		slog.Error("failed to render oauth consent page", slog.Any("error", err))
	}
}

// apiOAuthConsentHandler handles the Allow/Deny form posted from the consent
// page. The single-use consent id is only readable by the page we rendered for
// this user, so a cross-site form post cannot approve a request.
func (s *server) apiOAuthConsentHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := token.GetUserInfo(r)
	if err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusUnauthorized, err, "not signed in")
		return
	}
	if err := r.ParseForm(); err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, err, "invalid consent request")
		return
	}

	tokenUser := newTokenUser(&userInfo)
	consent, ok := s.apiOAuth.consumeConsent(r.PostForm.Get("consent_id"))
	if !ok || consent.UserID != tokenUser.databaseID() {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, errors.New("invalid or expired consent"), "invalid or expired consent")
		return
	}

	u, err := url.Parse(consent.RedirectURI)
	if err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, err, "invalid redirect_uri")
		return
	}
	q := u.Query()
	if consent.State != "" {
		q.Set("state", consent.State)
	}

	if r.PostForm.Get("action") != "allow" {
		q.Set("error", "access_denied")
		u.RawQuery = q.Encode()
		http.Redirect(w, r, u.String(), http.StatusFound)
		return
	}

	code, err := randomURLToken(32)
	if err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusInternalServerError, err, "failed to create code")
		return
	}
	s.apiOAuth.put(apiOAuthCode{
		CodeHash:      hashOAuthSecret(code),
		UserID:        consent.UserID,
		ClientID:      consent.ClientID,
		Registered:    consent.Registered,
		RedirectURI:   consent.RedirectURI,
		Resource:      consent.Resource,
		CodeChallenge: consent.CodeChallenge,
		ExpiresAt:     time.Now().Add(authCodeTTL),
	})

	q.Set("code", code)
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

// consentPage is the data rendered into the OAuth consent page.
type consentPage struct {
	ClientName   string
	RedirectHost string
	Username     string
	ConsentID    string
}

var consentTemplate = template.Must(template.New("consent").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Authorize access - craig-stars</title>
<style>
body { font-family: system-ui, sans-serif; background: #1d232a; color: #a6adbb; display: flex; justify-content: center; padding: 3rem 1rem; margin: 0; }
main { max-width: 28rem; width: 100%; background: #2a323c; border-radius: 0.5rem; padding: 1.5rem; }
h1 { font-size: 1.25rem; margin-top: 0; color: #fff; }
strong { color: #fff; }
.buttons { display: flex; gap: 0.75rem; margin-top: 1.5rem; }
button { flex: 1; padding: 0.6rem; border: 0; border-radius: 0.375rem; font-size: 1rem; cursor: pointer; }
.allow { background: #661ae6; color: #fff; }
.deny { background: #3d4451; color: #fff; }
</style>
</head>
<body>
<main>
<h1>Authorize access to craig-stars</h1>
<p><strong>{{.ClientName}}</strong> wants to act as <strong>{{.Username}}</strong> in your craig-stars games.</p>
<p>If you allow this, it can view your games and submit turns for you. Access will be sent to <strong>{{.RedirectHost}}</strong>.</p>
<p>Only allow this if you started connecting this application yourself.</p>
<form method="post" action="/api/oauth/authorize">
<input type="hidden" name="consent_id" value="{{.ConsentID}}">
<div class="buttons">
<button class="deny" type="submit" name="action" value="deny">Deny</button>
<button class="allow" type="submit" name="action" value="allow">Allow</button>
</div>
</form>
</main>
</body>
</html>
`))

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

func (s *apiOAuthStore) putConsent(consent apiOAuthConsent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pruneLocked(time.Now())
	s.consents[consent.IDHash] = consent
}

// consumeConsent returns the pending consent for id and removes it so it can
// only be answered once.
func (s *apiOAuthStore) consumeConsent(id string) (apiOAuthConsent, bool) {
	if id == "" {
		return apiOAuthConsent{}, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.pruneLocked(now)
	hash := hashOAuthSecret(id)
	consent, ok := s.consents[hash]
	if !ok {
		return apiOAuthConsent{}, false
	}
	delete(s.consents, hash)
	return consent, true
}

func (s *apiOAuthStore) pruneLocked(now time.Time) {
	for hash, code := range s.codes {
		if code.Used || !code.ExpiresAt.After(now) {
			delete(s.codes, hash)
		}
	}
	for hash, consent := range s.consents {
		if !consent.ExpiresAt.After(now) {
			delete(s.consents, hash)
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
	if u.Fragment != "" {
		return errors.New("redirect_uri must not include a fragment")
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
