package hold

import (
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/game"
	"github.com/timshannon/bolthold"
)

type DB struct {
	store *bolthold.Store
}

type Client interface {
	Connect(config *config.Config)
}

func (db *DB) Connect(config *config.Config) {
	log.Debug().Msgf("Connecting to database %s", config.Database.Filename)
	store, err := bolthold.Open(config.Database.Filename, 0666, nil)
	if err != nil {
		// handle error
		log.Fatal().Err(err).Msgf("Failed to open badger hold store %s", "data/bolt.db")
	}
	db.store = store
	defer store.Close()

	g := game.NewGame().WithSettings(*game.NewGameSettings())
	err = store.Insert("key", g)

	if err != nil {
		// handle error
		log.Fatal().Err(err).Msgf("Failed to insert %v", g)
	}

	var loaded game.Game
	err = store.Get("key", &loaded)
	if err != nil {
		// handle error
		log.Fatal().Err(err).Msgf("Failed to get key from store")
	}

	log.Info().Msgf("Loaded game %s", loaded.Name)
}
