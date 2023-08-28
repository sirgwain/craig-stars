package db

import "github.com/sirgwain/craig-stars/cs"

func (c *txClient) GetTechStores() ([]cs.TechStore, error) {
	// TODO: implement
	return []cs.TechStore{cs.StaticTechStore}, nil
}

func (c *txClient) CreateTechStore(tech *cs.TechStore) error {
	// TODO: implement
	return nil
}

func (c *txClient) GetTechStore(id int64) (*cs.TechStore, error) {
	// TODO: implement
	return nil, nil
}
