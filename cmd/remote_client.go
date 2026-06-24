package cmd

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
)

const (
	// defaultCLIServer is the production craig-stars server used when no
	// command flag, environment variable, or config value is set.
	defaultCLIServer = "https://craig-stars.net"

	// cliEnvServer names the environment variable that overrides the server URL.
	cliEnvServer = "CRAIG_STARS_SERVER"

	// cliEnvToken names the environment variable that supplies a bearer token.
	cliEnvToken = "CRAIG_STARS_TOKEN"
)

// cliConfig is the on-disk config format for the remote craig-stars CLI.
type cliConfig struct {
	DefaultServer string                     `json:"default_server,omitempty"`
	Servers       map[string]cliServerConfig `json:"servers,omitempty"`
}

// cliServerConfig stores credentials for one named server.
type cliServerConfig struct {
	Server      string    `json:"server"`
	AccessToken string    `json:"access_token"`
	Scope       string    `json:"scope,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

// tokenResponse is the OAuth token response returned by the craig-stars server.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

// cliClientConfig is the resolved runtime configuration for remote commands.
type cliClientConfig struct {
	Server      string
	ConfigPath  string
	AccessToken string
}

// remoteClients groups the generated Connect clients used by remote commands.
type remoteClients struct {
	Games craig_starsv1connect.GameServiceClient
	Races craig_starsv1connect.RaceServiceClient
}

// newRemoteClients constructs generated Connect clients with bearer auth.
func newRemoteClients(cfg cliClientConfig) remoteClients {
	opts := []connect.ClientOption{connect.WithInterceptors(newAuthInterceptor(cfg.AccessToken))}
	baseURL := connectBaseURL(cfg.Server)
	return remoteClients{
		Games: craig_starsv1connect.NewGameServiceClient(http.DefaultClient, baseURL, opts...),
		Races: craig_starsv1connect.NewRaceServiceClient(http.DefaultClient, baseURL, opts...),
	}
}

// resolveCLIClientConfig resolves flags, environment variables, and config file
// values into the settings required by a remote CLI command.
func resolveCLIClientConfig(requireToken bool) (cliClientConfig, error) {
	path, err := cliConfigFilePath(cliConfigPath)
	if err != nil {
		return cliClientConfig{}, err
	}
	fileCfg, err := readCLIConfig(path)
	if err != nil {
		return cliClientConfig{}, err
	}

	server := firstNonEmpty(cliServer, os.Getenv(cliEnvServer), fileCfg.DefaultServer, defaultCLIServer)
	server = normalizeServerURL(server)
	serverCfg := fileCfg.serverConfig(server)
	token := firstNonEmpty(cliToken, os.Getenv(cliEnvToken), serverCfg.AccessToken)
	cfg := cliClientConfig{
		Server:      server,
		ConfigPath:  path,
		AccessToken: token,
	}
	if requireToken && cfg.AccessToken == "" {
		return cliClientConfig{}, fmt.Errorf("missing API token; run `craig-stars login --server %s` or set %s", cfg.Server, cliEnvToken)
	}
	return cfg, nil
}

// serverConfig returns credentials for a normalized server URL.
func (c cliConfig) serverConfig(server string) cliServerConfig {
	name := serverName(server)
	if c.Servers != nil {
		if cfg, ok := c.Servers[name]; ok {
			return cfg
		}
	}
	return cliServerConfig{}
}

// newAuthInterceptor adds an Authorization header to outgoing Connect calls.
func newAuthInterceptor(token string) connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if token != "" {
				req.Header().Set("Authorization", "Bearer "+token)
			}
			return next(ctx, req)
		}
	})
}

// cliConfigFilePath returns the configured CLI config path, or the default path
// under the user's config directory.
func cliConfigFilePath(path string) (string, error) {
	if path != "" {
		return path, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "craig-stars", "cli.json"), nil
}

// readCLIConfig loads CLI config from disk; a missing file is treated as empty
// config so first-run commands can still use flags or environment variables.
func readCLIConfig(path string) (cliConfig, error) {
	var cfg cliConfig
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if len(strings.TrimSpace(string(b))) == 0 {
		return cfg, nil
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, fmt.Errorf("read CLI config %s: %w", path, err)
	}
	return cfg, nil
}

// writeCLIConfig stores CLI config with private file permissions.
func writeCLIConfig(path string, cfg cliConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o600)
}

// saveServerToken saves or replaces the token for a single server without
// removing credentials for other servers.
func saveServerToken(path, server string, token tokenResponse, expiresAt time.Time) error {
	cfg, err := readCLIConfig(path)
	if err != nil {
		return err
	}
	server = normalizeServerURL(server)
	if cfg.Servers == nil {
		cfg.Servers = map[string]cliServerConfig{}
	}
	if cfg.DefaultServer == "" {
		cfg.DefaultServer = defaultCLIServer
	}
	cfg.Servers[serverName(server)] = cliServerConfig{
		Server:      server,
		AccessToken: token.AccessToken,
		Scope:       token.Scope,
		ExpiresAt:   expiresAt,
	}
	return writeCLIConfig(path, cfg)
}

// normalizeServerURL normalizes user-supplied server values into absolute URLs.
func normalizeServerURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		raw = defaultCLIServer
	}
	raw = strings.TrimRight(raw, "/")
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	return raw
}

// serverName returns the stable config key for a normalized server URL.
func serverName(server string) string {
	u, err := url.Parse(normalizeServerURL(server))
	if err != nil || u.Host == "" {
		return normalizeServerURL(server)
	}
	return strings.ToLower(u.Host)
}

// connectBaseURL returns the base URL expected by generated Connect clients.
func connectBaseURL(server string) string {
	return strings.TrimRight(normalizeServerURL(server), "/") + "/api/grpc"
}

// endpointURL joins a normalized server URL with an absolute API path.
func endpointURL(server, path string) string {
	return strings.TrimRight(normalizeServerURL(server), "/") + path
}

// firstNonEmpty returns the first string with non-whitespace content.
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// makePKCEPair creates a verifier and S256 challenge for OAuth PKCE.
func makePKCEPair() (verifier string, challenge string, err error) {
	verifier, err = randomURLToken(32)
	if err != nil {
		return "", "", err
	}
	sum := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(sum[:])
	return verifier, challenge, nil
}

// randomURLToken returns cryptographically random URL-safe token text.
func randomURLToken(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("make random token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// openBrowser opens the user's default browser to the provided URL.
func openBrowser(rawURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", rawURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL)
	default:
		cmd = exec.Command("xdg-open", rawURL)
	}
	return cmd.Start()
}

// exchangeOAuthCode exchanges an OAuth authorization code for an API token.
func exchangeOAuthCode(ctx context.Context, server, code, redirectURI, verifier string) (tokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)
	form.Set("code_verifier", verifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointURL(server, "/api/oauth/token"), strings.NewReader(form.Encode()))
	if err != nil {
		return tokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return tokenResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tokenResponse{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return tokenResponse{}, fmt.Errorf("token exchange failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var token tokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return tokenResponse{}, err
	}
	if !strings.EqualFold(token.TokenType, "Bearer") || token.AccessToken == "" {
		return tokenResponse{}, errors.New("token response did not include a bearer access token")
	}
	return token, nil
}
