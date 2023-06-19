package db

import (
	"testing"

	"github.com/sirgwain/craig-stars/config"
)

func TestDB_MigrateAll(t *testing.T) {
	db := &DB{}
	cfg := &config.Config{}
	cfg.Database.Filename = ":memory:"
	db.Connect(cfg)
	if err := db.MigrateAll(); err != nil {
		t.Errorf("DB.MigrateAll() error = %v", err)
	}
}
