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

func TestValidateRedirectURI(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{name: "native app callback", uri: "cstars://oauth/callback"},
		{name: "native app callback with query", uri: "cstars://oauth/callback?source=ios"},
		{name: "native app missing host rejected", uri: "cstars:///oauth/callback", wantErr: true},
		{name: "other custom scheme rejected", uri: "other://oauth/callback", wantErr: true},
		{name: "loopback callback", uri: "http://127.0.0.1:39123/callback"},
		{name: "remote http callback rejected", uri: "http://example.com:39123/callback", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRedirectURI(tt.uri)
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

func TestAPIOAuthMetadataUsesCanonicalEndpoints(t *testing.T) {
	metadata := apiOAuthMetadata("https://craig-stars.example")

	if got := metadata["authorization_endpoint"]; got != "https://craig-stars.example/api/oauth/authorize" {
		t.Fatalf("authorization_endpoint = %q", got)
	}
	if got := metadata["token_endpoint"]; got != "https://craig-stars.example/api/oauth/token" {
		t.Fatalf("token_endpoint = %q", got)
	}
	if got := metadata["registration_endpoint"]; got != "https://craig-stars.example/api/mcp/oauth/register" {
		t.Fatalf("registration_endpoint = %q", got)
	}
}
