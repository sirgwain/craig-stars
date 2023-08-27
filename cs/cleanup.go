package cs

import "github.com/rs/zerolog/log"

/*
Because the game is actively running, bugs are introduced and we need to fix them as
we go. For example, a bug was introduced to "discover" player starbase designs leading
to duplciate designs on the UI.

This file will hold cleanup functions for cleaning up this sort of data on the game
*/

type Cleaner interface {
	RemovePlayerDesignIntels(game *FullGame)
}

type cleanup struct {
}

func NewCleaner() Cleaner {
	return &cleanup{}
}

// cleanup design intels owned by the player
func (c *cleanup) RemovePlayerDesignIntels(game *FullGame) {
	for _, player := range game.Players {
		designIntels := make([]ShipDesignIntel, 0, len(player.ShipDesignIntels))
		for _, design := range player.ShipDesignIntels {
			if design.PlayerNum != player.Num {
				designIntels = append(designIntels, design)
			} else {
				log.Info().
					Int64("GameID", game.ID).
					Str("Name", game.Name).
					Msgf("cleanup: removing design intel for player %d's %s design", player.Num, design.Name)
			}
		}
		player.ShipDesignIntels = designIntels
	}
}
