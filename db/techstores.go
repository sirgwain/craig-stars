package db

import (
	"context"

	"github.com/sirgwain/craig-stars/cs"
)

func (c *client) GetTechStores(ctx context.Context) ([]cs.TechStore, error) {
	// TODO: implement
	return []cs.TechStore{cs.StaticTechStore}, nil
}

func (c *client) CreateTechStore(ctx context.Context, tech *cs.TechStore) (*cs.TechStore, error) {
	// TODO: implement
	return nil, nil
}

func (c *client) GetTechStore(ctx context.Context, id int64) (*cs.TechStore, error) {
	// TODO: implement
	return nil, nil
}
