package cmd

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

func TestLogout(t *testing.T) {
	var revoked []string
	status := http.StatusOK
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/oauth/revoke" {
			http.NotFound(w, r)
			return
		}
		revoked = append(revoked, r.FormValue("token"))
		w.WriteHeader(status)
	}))
	t.Cleanup(srv.Close)

	path := filepath.Join(t.TempDir(), "cli.json")
	expires := time.Now().Add(time.Hour)
	if err := saveServerToken(path, srv.URL, tokenResponse{AccessToken: "cstars_api_mine"}, expires); err != nil {
		t.Fatal(err)
	}
	if err := saveServerToken(path, "https://other.example", tokenResponse{AccessToken: "cstars_api_other"}, expires); err != nil {
		t.Fatal(err)
	}
	cfg := cliClientConfig{Server: normalizeServerURL(srv.URL), ConfigPath: path}

	var out bytes.Buffer
	if err := logout(context.Background(), cfg, &out); err != nil {
		t.Fatalf("logout: %v", err)
	}
	if len(revoked) != 1 || revoked[0] != "cstars_api_mine" {
		t.Errorf("revoked %v, want the saved token", revoked)
	}
	fileCfg, err := readCLIConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if fileCfg.serverConfig(srv.URL).AccessToken != "" {
		t.Errorf("token for %s still saved", srv.URL)
	}
	if fileCfg.serverConfig("https://other.example").AccessToken != "cstars_api_other" {
		t.Errorf("logout removed another server's token")
	}

	// logging out again is a no-op
	if err := logout(context.Background(), cfg, &out); err != nil || len(revoked) != 1 {
		t.Errorf("second logout: %v, %d revoke calls", err, len(revoked))
	}

	// when the server fails, the token is still forgotten locally
	if err := saveServerToken(path, srv.URL, tokenResponse{AccessToken: "cstars_api_mine"}, expires); err != nil {
		t.Fatal(err)
	}
	status = http.StatusInternalServerError
	if err := logout(context.Background(), cfg, &out); err == nil {
		t.Errorf("logout with a failing server: want an error")
	}
	if fileCfg, _ = readCLIConfig(path); fileCfg.serverConfig(srv.URL).AccessToken != "" {
		t.Errorf("token kept after the server failed to revoke it")
	}
}
