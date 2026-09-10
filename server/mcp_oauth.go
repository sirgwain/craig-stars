package server

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-pkgz/rest"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

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

// validateMCPRedirectURI bridges MCP dynamic client registrations with the
// shared API OAuth flow. Unregistered native and CLI clients use safe local
// callback validation.
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
	return false, validateRedirectURI(redirectURI)
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
		"authorization_server_metadata_endpoint": base + "/api/oauth/metadata",
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
		if err := validateMCPRegistrationRedirectURI(redirectURI); err != nil {
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

// validateMCPRegistrationRedirectURI supports https:// schemes for chatgpt and localhost uris for claude
func validateMCPRegistrationRedirectURI(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}

	if u.Fragment != "" {
		return errors.New("redirect_uri must not include a fragment")
	}

	if u.Scheme == "https" {
		if u.Hostname() == "" {
			return errors.New("redirect_uri missing host")
		}
		return nil
	}

	return validateLoopbackRedirectURI(raw)
}

func stringSliceContains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
