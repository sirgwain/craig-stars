package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
)

func TestMCPOAuthClientRoundTrip(t *testing.T) {
	c := connectTestDB()
	redirectURI := "https://chatgpt.com/connector/oauth/abcd0123-"
	client := &cs.MCPOAuthClient{
		ClientID:                "mcp_test",
		ClientName:              "ChatGPT",
		ClientURI:               "https://chatgpt.com",
		RedirectURIs:            []string{redirectURI},
		TokenEndpointAuthMethod: "none",
		Scope:                   "api:player",
		ClientIDIssuedAt:        12345,
	}

	created, err := c.CreateMCPOAuthClient(t.Context(), client)
	if err != nil {
		t.Fatalf("CreateMCPOAuthClient() error = %v", err)
	}
	if created.ClientID != client.ClientID {
		t.Fatalf("created client id = %q, want %q", created.ClientID, client.ClientID)
	}

	got, err := c.GetMCPOAuthClient(t.Context(), client.ClientID)
	if err != nil {
		t.Fatalf("GetMCPOAuthClient() error = %v", err)
	}
	if got == nil {
		t.Fatal("expected stored client")
	}
	if got.ClientName != client.ClientName || got.ClientURI != client.ClientURI {
		t.Fatalf("stored client = %+v, want name %q uri %q", got, client.ClientName, client.ClientURI)
	}

	ok, err := c.HasMCPOAuthRedirectURI(t.Context(), client.ClientID, redirectURI)
	if err != nil {
		t.Fatalf("HasMCPOAuthRedirectURI() error = %v", err)
	}
	if !ok {
		t.Fatal("expected registered redirect uri")
	}

	ok, err = c.HasMCPOAuthRedirectURI(t.Context(), client.ClientID, "https://chatgpt.com/connector/oauth/other")
	if err != nil {
		t.Fatalf("HasMCPOAuthRedirectURI() unexpected error = %v", err)
	}
	if ok {
		t.Fatal("unexpected unregistered redirect uri")
	}
}
