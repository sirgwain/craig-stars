package appcontext

import (
	"context"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/db"
)

type AppContext struct {
	context.Context
	// Config *config.Config
	DB     db.Service
	Config *config.Config
}

// Initialize a new context for this app. This should be done in cmd before starting
// a server or
func Initialize() *AppContext {
	db := &db.DB{}
	cfg := config.GetConfig()
	db.Connect(cfg)

	return &AppContext{
		DB:     db,
		Config: cfg,
	}
}
