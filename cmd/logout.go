package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

func newLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of a craig-stars server",
		Long: `Revoke the API token saved for a craig-stars server and remove it from the
CLI config. Other servers' tokens are kept. Your browser stays logged in to
the web app; log out there too if you want the next login to use a
different account.`,
		Example: `  craig-stars logout
  craig-stars logout --server http://localhost:8080`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := resolveCLIClientConfig(false)
			if err != nil {
				return err
			}
			return logout(cmd.Context(), cfg, cmd.OutOrStdout())
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newLogoutCmd())
}

// logout revokes the token saved for cfg.Server and removes it from the config
// file. The saved token is removed even when the server can't be reached, so
// the next command logs in again; the returned error then says the token may
// still work until it expires. Tokens from flags or the environment are not
// touched, since they aren't ours to forget.
func logout(ctx context.Context, cfg cliClientConfig, out io.Writer) error {
	fileCfg, err := readCLIConfig(cfg.ConfigPath)
	if err != nil {
		return err
	}
	saved := fileCfg.serverConfig(cfg.Server)
	if saved.AccessToken == "" {
		fmt.Fprintf(out, "Not logged in to %s.\n", cfg.Server)
		return nil
	}

	revokeErr := revokeToken(ctx, cfg.Server, saved.AccessToken)
	if err := removeServerToken(cfg.ConfigPath, cfg.Server); err != nil {
		return err
	}
	if revokeErr != nil {
		return fmt.Errorf("removed the saved token for %s, but the server didn't revoke it, so it works until it expires: %w", cfg.Server, revokeErr)
	}
	fmt.Fprintf(out, "Logged out of %s.\n", cfg.Server)
	return nil
}

// revokeToken asks the server to revoke an API token (RFC 7009).
func revokeToken(ctx context.Context, server, token string) error {
	form := url.Values{}
	form.Set("token", token)
	form.Set("token_type_hint", "access_token")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointURL(server, "/api/oauth/revoke"), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("revoke failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

// removeServerToken deletes the saved credentials for one server, leaving the
// others alone. Removing a server that has none is not an error.
func removeServerToken(path, server string) error {
	cfg, err := readCLIConfig(path)
	if err != nil {
		return err
	}
	name := serverName(server)
	if _, ok := cfg.Servers[name]; !ok {
		return nil
	}
	delete(cfg.Servers, name)
	return writeCLIConfig(path, cfg)
}
