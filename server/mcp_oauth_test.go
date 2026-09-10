package server

import "testing"

func TestValidateHostedRedirectURI(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{name: "chatgpt callback", uri: "https://chatgpt.com/connector/oauth/abcd0123-"},
		{name: "legacy chatgpt callback", uri: "https://chatgpt.com/connector_platform_oauth_redirect"},
		{name: "http rejected", uri: "http://chatgpt.com/connector/oauth/abcd0123-", wantErr: true},
		{name: "missing host rejected", uri: "https:///connector/oauth/abcd0123-", wantErr: true},
		{name: "fragment rejected", uri: "https://chatgpt.com/connector/oauth/abcd0123-#frag", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMCPRegistrationRedirectURI(tt.uri)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
