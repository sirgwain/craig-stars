package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

func newLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with a craig-stars server",
		Long:  `Authenticate with a craig-stars server and save an API bearer token for CLI use.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := resolveCLIClientConfig(false)
			if err != nil {
				return err
			}

			verifier, challenge, err := makePKCEPair()
			if err != nil {
				return err
			}
			state, err := randomURLToken(16)
			if err != nil {
				return err
			}

			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				return fmt.Errorf("start callback listener: %w", err)
			}
			defer listener.Close()

			callback := make(chan url.Values, 1)
			callbackErr := make(chan error, 1)
			mux := http.NewServeMux()
			mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
				if got := r.URL.Query().Get("state"); got != state {
					http.Error(w, "invalid state", http.StatusBadRequest)
					callbackErr <- fmt.Errorf("invalid OAuth state")
					return
				}
				if errMsg := r.URL.Query().Get("error"); errMsg != "" {
					http.Error(w, errMsg, http.StatusBadRequest)
					callbackErr <- fmt.Errorf("authorization failed: %s", errMsg)
					return
				}
				fmt.Fprintln(w, "craig-stars CLI login complete. You can close this window.")
				callback <- r.URL.Query()
			})
			server := &http.Server{Handler: mux}
			defer server.Close()
			go func() {
				if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
					callbackErr <- err
				}
			}()

			redirectURI := "http://" + listener.Addr().String() + "/callback"
			authURL, err := url.Parse(endpointURL(cfg.Server, "/api/mcp/oauth/authorize"))
			if err != nil {
				return err
			}
			q := authURL.Query()
			q.Set("response_type", "code")
			q.Set("redirect_uri", redirectURI)
			q.Set("code_challenge", challenge)
			q.Set("code_challenge_method", "S256")
			q.Set("state", state)
			authURL.RawQuery = q.Encode()

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Opening browser for craig-stars login:\n%s\n", authURL.String())
			if err := openBrowser(authURL.String()); err != nil {
				fmt.Fprintf(out, "Could not open browser automatically: %v\n", err)
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Minute)
			defer cancel()
			var values url.Values
			select {
			case values = <-callback:
			case err := <-callbackErr:
				return err
			case <-ctx.Done():
				return fmt.Errorf("timed out waiting for OAuth callback")
			}

			code := values.Get("code")
			if code == "" {
				return fmt.Errorf("OAuth callback did not include a code")
			}
			token, err := exchangeOAuthCode(ctx, cfg.Server, code, redirectURI, verifier)
			if err != nil {
				return err
			}

			expiresAt := time.Time{}
			if token.ExpiresIn > 0 {
				expiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
			}
			if err := saveServerToken(cfg.ConfigPath, cfg.Server, token, expiresAt); err != nil {
				return err
			}

			fmt.Fprintf(out, "Saved craig-stars CLI token for %s to %s\n", cfg.Server, cfg.ConfigPath)
			return nil
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newLoginCmd())
}
