package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db/generated"
)

func (c *client) CreateAPIToken(ctx context.Context, token *cs.APIToken, tokenHash string) (*cs.APIToken, error) {
	item, err := c.writer.CreateAPIToken(ctx, generated.CreateAPITokenParams{
		UserID:      token.UserID,
		Name:        token.Name,
		TokenPrefix: token.TokenPrefix,
		TokenHash:   tokenHash,
		Scope:       token.Scope,
		ExpiresAt:   token.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}
	return apiTokenFromGenerated(item), nil
}

func (c *client) GetAPITokenByHash(ctx context.Context, tokenHash string) (*cs.APIToken, error) {
	item, err := c.reader.GetAPITokenByHash(ctx, tokenHash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return apiTokenFromGenerated(item), nil
}

func (c *client) GetAPITokensForUser(ctx context.Context, userID int64) ([]cs.APIToken, error) {
	items, err := c.reader.GetAPITokensForUser(ctx, userID)
	if err == sql.ErrNoRows {
		return []cs.APIToken{}, nil
	}
	if err != nil {
		return nil, err
	}

	tokens := make([]cs.APIToken, 0, len(items))
	for _, item := range items {
		tokens = append(tokens, *apiTokenFromGenerated(item))
	}
	return tokens, nil
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

func apiTokenFromGenerated(item generated.ApiToken) *cs.APIToken {
	return &cs.APIToken{
		DBObject: cs.DBObject{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		},
		UserID:      item.UserID,
		Name:        item.Name,
		TokenPrefix: item.TokenPrefix,
		Scope:       item.Scope,
		ExpiresAt:   item.ExpiresAt,
		LastUsedAt:  nullTimePtr(item.LastUsedAt),
		RevokedAt:   nullTimePtr(item.RevokedAt),
	}
}

func nullTimePtr(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
