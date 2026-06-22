package db

import (
	"context"
	"database/sql"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db/generated"
)

func (c *client) CreateAPIToken(ctx context.Context, token *cs.APIToken, tokenHash string) (*cs.APIToken, error) {
	params := c.converter.ConvertGameApiTokenToCreateParams(token)
	params.TokenHash = tokenHash
	item, err := c.writer.CreateAPIToken(ctx, params)
	if err != nil {
		return nil, err
	}
	t := c.converter.ConvertApiToken(item)
	return &t, nil
}

func (c *client) GetAPITokenByHash(ctx context.Context, tokenHash string) (*cs.APIToken, error) {
	item, err := c.reader.GetAPITokenByHash(ctx, tokenHash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t := c.converter.ConvertApiToken(item)
	return &t, nil
}

func (c *client) GetAPITokensForUser(ctx context.Context, userID int64) ([]cs.APIToken, error) {
	items, err := c.reader.GetAPITokensForUser(ctx, userID)
	if err == sql.ErrNoRows {
		return []cs.APIToken{}, nil
	}
	if err != nil {
		return nil, err
	}

	return c.converter.ConvertApiTokens(items), nil
}

func (c *client) TouchAPIToken(ctx context.Context, id int64) error {
	return c.writer.TouchAPIToken(ctx, id)
}

func (c *client) RevokeAPIToken(ctx context.Context, userID, id int64) error {
	return c.writer.RevokeAPIToken(ctx, generated.RevokeAPITokenParams{
		ID:     id,
		UserID: userID,
	})
}
