package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/config"
)

func TestDB_MigrateAll(t *testing.T) {
	c := &client{}
	cfg := &config.Config{}
	cfg.Database.Filename = ":memory:"
	c.Connect(cfg)
	if err := c.MigrateAll(); err != nil {
		t.Errorf("c.MigrateAll() error = %v", err)
	}
}
