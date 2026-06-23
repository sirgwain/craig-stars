package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db/generated"
)

func (c *client) CreateMCPOAuthClient(ctx context.Context, client *cs.MCPOAuthClient) (*cs.MCPOAuthClient, error) {
	item, err := c.writer.CreateMCPOAuthClient(ctx, generated.CreateMCPOAuthClientParams{
		ClientID:                client.ClientID,
		ClientName:              client.ClientName,
		ClientUri:               client.ClientURI,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
		Scope:                   client.Scope,
		ClientIDIssuedAt:        client.ClientIDIssuedAt,
	})
	if err != nil {
		return nil, err
	}
	for _, redirectURI := range client.RedirectURIs {
		err = c.writer.CreateMCPOAuthRedirectURI(ctx, generated.CreateMCPOAuthRedirectURIParams{
			ClientID:    client.ClientID,
			RedirectUri: redirectURI,
		})
		if err != nil {
			return nil, err
		}
	}

	result := c.converter.ConvertMCPOAuthClient(item)
	result.RedirectURIs = append([]string(nil), client.RedirectURIs...)
	return &result, nil
}

func (c *client) GetMCPOAuthClient(ctx context.Context, clientID string) (*cs.MCPOAuthClient, error) {
	item, err := c.reader.GetMCPOAuthClient(ctx, clientID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	result := c.converter.ConvertMCPOAuthClient(item)
	return &result, nil
}

func (c *client) HasMCPOAuthRedirectURI(ctx context.Context, clientID, redirectURI string) (bool, error) {
	_, err := c.reader.GetMCPOAuthRedirectURI(ctx, generated.GetMCPOAuthRedirectURIParams{
		ClientID:    clientID,
		RedirectUri: redirectURI,
	})
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
