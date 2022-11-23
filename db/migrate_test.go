package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/config"
)

func TestDB_migrateAll(t *testing.T) {
	c := &client{}
	cfg := &config.Config{}
	cfg.Database.Filename = ":memory:"
	c.Connect(cfg)
	if err := c.migrateAll(); err != nil {
		t.Errorf("c.MigrateAll() error = %v", err)
	}
}
