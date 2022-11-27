package db

import "github.com/sirgwain/craig-stars/game"

func (c *client) GetTechStores() ([]game.TechStore, error) {
	// TODO: implement
	return []game.TechStore{game.StaticTechStore}, nil
}

func (c *client) CreateTechStore(tech *game.TechStore) error {
	// TODO: implement
	return nil
}

func (c *client) GetTechStore(id int64) (*game.TechStore, error) {
	// TODO: implement
	return nil, nil
}
