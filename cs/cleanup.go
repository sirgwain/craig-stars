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
	AddScannerToInnateScannerPlanets(game *FullGame)
	AddRandomArtifactsToPlanets(game *FullGame)
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

// ensure AR planets have scanners
func (c *cleanup) AddScannerToInnateScannerPlanets(game *FullGame) {
	for _, planet := range game.Planets {
		if !planet.Owned() {
			continue
		}

		player := game.getPlayer(planet.PlayerNum)
		if !player.Race.Spec.InnateScanner {
			continue
		}

		planet.Scanner = true
		planet.MarkDirty()
		log.Info().
			Int64("GameID", game.ID).
			Str("Name", game.Name).
			Msgf("cleanup: setting scanner to true for player %d's %s planet", player.Num, planet.Name)
	}
}

// we support random artifacts now, yay! upgrade our existing games
func (c *cleanup) AddRandomArtifactsToPlanets(game *FullGame) {
	if !game.RandomEvents {
		return
	}

	for _, planet := range game.Planets {
		if planet.Owned() {
			continue
		}

		// check if this planet should have a random artifact
		if game.Rules.RandomEventChances[RandomEventAncientArtifact] > game.Rules.random.Float64() {
			planet.RandomArtifact = true
			planet.MarkDirty()

			log.Info().
				Int64("GameID", game.ID).
				Str("Name", game.Name).
				Msgf("cleanup: added random artifact to planet %s", planet.Name)
		}

	}
}
