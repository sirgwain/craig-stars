package server

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/sirgwain/craig-stars/db"
)

const authCodeTTL = 5 * time.Minute

// mcpAuthStore keeps short-lived MCP authorization codes between the browser
// authorization redirect and the local client's token exchange.
type mcpAuthStore struct {
	mu    sync.Mutex
	codes map[string]mcpAuthCode
}

// mcpAuthCode records the server-side state for one OAuth authorization code.
type mcpAuthCode struct {
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

type mcpRegisteredClient struct {
	ClientID                string   `json:"client_id"`
	ClientName              string   `json:"client_name,omitempty"`
	ClientURI               string   `json:"client_uri,omitempty"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	Scope                   string   `json:"scope,omitempty"`
	ClientIDIssuedAt        int64    `json:"client_id_issued_at"`
}

// tokenResponse is the OAuth token response returned to MCP clients.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

// newMCPAuthStore creates an empty in-memory authorization code store.
func newMCPAuthStore() *mcpAuthStore {
	return &mcpAuthStore{codes: map[string]mcpAuthCode{}}
}

// mcpAuthorizeHandler starts the MCP browser authorization flow, using the
// normal craig-stars login and then redirecting back to the local MCP client.
func (s *server) mcpAuthorizeHandler(w http.ResponseWriter, r *http.Request) {
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
	s.mcpAuth.put(mcpAuthCode{
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

// mcpTokenHandler exchanges a valid MCP authorization code and PKCE verifier
// for an opaque craig-stars API bearer token.
func (s *server) mcpTokenHandler(w http.ResponseWriter, r *http.Request) {
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

	authCode, ok := s.mcpAuth.consume(code)
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

// mcpOAuthMetadataHandler publishes OAuth server metadata for MCP clients that
// discover the hosted authorization flow.
func (s *server) mcpOAuthMetadataHandler(w http.ResponseWriter, r *http.Request) {
	base := strings.TrimRight(s.config.Auth.URL, "/")
	rest.RenderJSON(w, map[string]any{
		"issuer":                                base,
		"authorization_endpoint":                base + "/api/mcp/oauth/authorize",
		"token_endpoint":                        base + "/api/mcp/oauth/token",
		"registration_endpoint":                 base + "/api/mcp/oauth/register",
		"response_types_supported":              []string{"code"},
		"grant_types_supported":                 []string{"authorization_code"},
		"code_challenge_methods_supported":      []string{"S256"},
		"token_endpoint_auth_methods_supported": []string{"none"},
		"scopes_supported":                      []string{apiTokenScopePlayer},
	})
}

// mcpResourceMetadataHandler publishes protected resource metadata for the
// hosted Streamable HTTP MCP endpoint.
func (s *server) mcpResourceMetadataHandler(w http.ResponseWriter, r *http.Request) {
	base := strings.TrimRight(s.config.Auth.URL, "/")
	rest.RenderJSON(w, map[string]any{
		"resource":                               s.mcpResourceURI(),
		"authorization_servers":                  []string{base},
		"bearer_methods_supported":               []string{"header"},
		"resource_documentation":                 base + "/docs/mcp",
		"resource_name":                          "craig-stars MCP",
		"scopes_supported":                       []string{apiTokenScopePlayer},
		"mcp_streamable_http_endpoint":           base + "/api/mcp",
		"mcp_protocol_transport":                 "streamable-http",
		"mcp_protocol_version":                   "2025-06-18",
		"authorization_server_metadata_endpoint": base + "/api/mcp/oauth/metadata",
	})
}

// mcpRegisterHandler implements OAuth Dynamic Client Registration for ChatGPT
// and other hosted MCP clients. Registered clients are public PKCE clients.
func (s *server) mcpRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req mcpRegisteredClient
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, err, "invalid registration request")
		return
	}
	if len(req.RedirectURIs) == 0 {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, errors.New("missing redirect_uris"), "missing redirect_uris")
		return
	}
	for _, redirectURI := range req.RedirectURIs {
		if err := validateHostedRedirectURI(redirectURI); err != nil {
			rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, err, "invalid redirect_uri")
			return
		}
	}
	if len(req.GrantTypes) == 0 {
		req.GrantTypes = []string{"authorization_code"}
	}
	if len(req.ResponseTypes) == 0 {
		req.ResponseTypes = []string{"code"}
	}
	if req.TokenEndpointAuthMethod == "" {
		req.TokenEndpointAuthMethod = "none"
	}
	if req.TokenEndpointAuthMethod != "none" {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, errors.New("unsupported token_endpoint_auth_method"), "unsupported token_endpoint_auth_method")
		return
	}
	if !stringSliceContains(req.GrantTypes, "authorization_code") || !stringSliceContains(req.ResponseTypes, "code") {
		rest.SendErrorJSON(w, r, nil, http.StatusBadRequest, errors.New("unsupported OAuth flow"), "unsupported OAuth flow")
		return
	}

	clientID, err := randomURLToken(24)
	if err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusInternalServerError, err, "failed to create client")
		return
	}
	req.ClientID = "mcp_" + clientID
	req.ClientIDIssuedAt = time.Now().Unix()
	var created *cs.MCPOAuthClient
	err = s.db.WrapInTransaction(func(c db.Client) error {
		var err error
		created, err = c.CreateMCPOAuthClient(r.Context(), &cs.MCPOAuthClient{
			ClientID:                req.ClientID,
			ClientName:              req.ClientName,
			ClientURI:               req.ClientURI,
			RedirectURIs:            req.RedirectURIs,
			TokenEndpointAuthMethod: req.TokenEndpointAuthMethod,
			Scope:                   req.Scope,
			ClientIDIssuedAt:        req.ClientIDIssuedAt,
		})
		return err
	})
	if err != nil {
		rest.SendErrorJSON(w, r, nil, http.StatusInternalServerError, err, "failed to save client")
		return
	}
	req.ClientID = created.ClientID
	req.ClientIDIssuedAt = created.ClientIDIssuedAt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(req); err != nil {
		slog.Error("failed to write MCP OAuth registration response", "error", err)
	}
}

func (s *server) mcpResourceURI() string {
	return strings.TrimRight(s.config.Auth.URL, "/") + "/api/mcp"
}

func (s *server) mcpResourceMetadataURI() string {
	return strings.TrimRight(s.config.Auth.URL, "/") + "/.well-known/oauth-protected-resource/api/mcp"
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

// put stores an authorization code and prunes stale entries.
func (s *mcpAuthStore) put(code mcpAuthCode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pruneLocked(time.Now())
	s.codes[code.CodeHash] = code
}

// consume marks an authorization code as used and returns its stored state.
func (s *mcpAuthStore) consume(code string) (mcpAuthCode, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.pruneLocked(now)
	hash := hashOAuthSecret(code)
	authCode, ok := s.codes[hash]
	if !ok || authCode.Used || !authCode.ExpiresAt.After(now) {
		return mcpAuthCode{}, false
	}
	authCode.Used = true
	s.codes[hash] = authCode
	return authCode, true
}

// pruneLocked removes expired or already-used authorization codes.
func (s *mcpAuthStore) pruneLocked(now time.Time) {
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

func (s *server) validateMCPRedirectURI(ctx context.Context, clientID, redirectURI string) (bool, error) {
	if clientID != "" {
		client, err := s.db.NewReadClient().GetMCPOAuthClient(ctx, clientID)
		if err != nil {
			return false, err
		}
		if client != nil {
			ok, err := s.db.NewReadClient().HasMCPOAuthRedirectURI(ctx, clientID, redirectURI)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
			return true, errors.New("redirect_uri not registered for client")
		}
	}
	return false, validateLoopbackRedirectURI(redirectURI)
}

// validateLoopbackRedirectURI restricts MCP authorization callbacks to local
// HTTP listeners owned by desktop clients.
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

func validateHostedRedirectURI(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.Scheme != "https" {
		return errors.New("redirect_uri must use https")
	}
	if u.Hostname() == "" {
		return errors.New("redirect_uri missing host")
	}
	if u.Fragment != "" {
		return errors.New("redirect_uri must not include a fragment")
	}
	return nil
}

func stringSliceContains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

// verifyPKCE checks an OAuth S256 PKCE verifier against the stored challenge.
func verifyPKCE(verifier, expectedChallenge string) bool {
	sum := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])
	return subtle.ConstantTimeCompare([]byte(challenge), []byte(expectedChallenge)) == 1
}

// randomURLToken returns URL-safe random bytes encoded without padding.
func randomURLToken(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("make random token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// hashOAuthSecret returns the SHA-256 hex digest used to store short-lived
// OAuth authorization codes.
func hashOAuthSecret(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hexString(sum[:])
}

// hexString encodes bytes as lowercase hexadecimal.
func hexString(b []byte) string {
	const alphabet = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = alphabet[v>>4]
		out[i*2+1] = alphabet[v&0x0f]
	}
	return string(out)
}
