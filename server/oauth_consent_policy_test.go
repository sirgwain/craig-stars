package server

import "testing"

func TestConsentContentSecurityPolicy(t *testing.T) {
	tests := []struct {
		name, redirectURI, formAction string
	}{
		{"hosted callback", "https://chatgpt.com/connector/oauth/callback?state=abc", "'self' https://chatgpt.com"},
		{"callback query is not a policy", "https://chatgpt.com/callback?next=;script-src%20*", "'self' https://chatgpt.com"},
		{"loopback callback", "http://127.0.0.1:39123/callback", "'self' http://127.0.0.1:39123"},
		{"IPv6 loopback callback", "http://[::1]:39123/callback", "'self' http://[::1]:39123"},
		{"native callback", "cstars://oauth/callback", "'self' cstars:"},
		{"malformed callback", "://invalid", "'self'"},
		{"unsupported scheme", "javascript:alert(1)", "'self'"},
		{"directive punctuation", "https://example.com;script-src/callback", "'self'"},
		{"policy punctuation", "https://example.com,other/callback", "'self'"},
		{"wildcard host", "https://*.example.com/callback", "'self'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := "default-src 'none'; style-src 'unsafe-inline'; form-action " + tt.formAction + "; frame-ancestors 'none'"
			if got := consentContentSecurityPolicy(tt.redirectURI); got != want {
				t.Fatalf("policy = %q, want %q", got, want)
			}
		})
	}
}
