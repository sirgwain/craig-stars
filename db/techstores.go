package db

import "github.com/sirgwain/craig-stars/cs"

func (c *client) GetTechStores() ([]cs.TechStore, error) {
	// TODO: implement
	return []cs.TechStore{cs.StaticTechStore}, nil
}

func (c *client) CreateTechStore(tech *cs.TechStore) error {
	// TODO: implement
	return nil
}


func (c *client) getTechStore(db SQLSelector, id int64) (*cs.TechStore, error) {
	return nil, nil
}

func (c *client) GetTechStore(id int64) (*cs.TechStore, error) {
	// TODO: implement
	return nil, nil
}
