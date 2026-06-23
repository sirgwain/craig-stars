package cs

// MCPOAuthClient is a dynamically registered OAuth client for hosted MCP
// clients such as ChatGPT. It is not a craig-stars user and does not grant
// access until a user completes OAuth and receives an APIToken.
type MCPOAuthClient struct {
	DBObject
	ClientID                string   `json:"clientId"`
	ClientName              string   `json:"clientName,omitempty"`
	ClientURI               string   `json:"clientUri,omitempty"`
	RedirectURIs            []string `json:"redirectUris"`
	TokenEndpointAuthMethod string   `json:"tokenEndpointAuthMethod"`
	Scope                   string   `json:"scope,omitempty"`
	ClientIDIssuedAt        int64    `json:"clientIdIssuedAt"`
}
