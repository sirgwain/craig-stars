package server

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"
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
		{name: "missing port rejected", uri: "http://localhost/callback", wantErr: true},
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

func TestAPITokenConstants(t *testing.T) {
	if apiTokenPrefix != "cstars_api_" {
		t.Fatalf("apiTokenPrefix = %q", apiTokenPrefix)
	}
	if apiTokenScopePlayer != "api:player" {
		t.Fatalf("apiTokenScopePlayer = %q", apiTokenScopePlayer)
	}
}
