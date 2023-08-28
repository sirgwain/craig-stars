package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

/*
Upgrade data in the database based on game updates

Version Info:

001 - Fix player discovers own starbase designs on planet discovery and adds them to intel


*/

type Version struct {
	ID        int64     `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Current   int       `json:"current,omitempty"`
}

const LATEST_VERSION = 1

func (dbClient *dbClient) mustUpgrade() {
	c, err := dbClient.BeginTransaction()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to begin upgrade transaction")
	}
	defer func() { dbClient.Rollback(c) }()

	if err := c.ensureUpgrade(); err != nil {
		log.Fatal().Err(err).Msg("failed to upgrade database data to latest version")
	}

	dbClient.Commit(c)
}

func (c *txClient) ensureUpgrade() error {
	version, err := c.getVersion()
	if err != nil {
		return err
	}

	if version.Current < LATEST_VERSION {
		for current := version.Current; current < LATEST_VERSION; current++ {
			// check each version and call the upgrade functionality
			switch current {
			case 0:
				log.Info().Msgf("upgrading database data from v0 to v1")
				err = c.upgrade1()
			}

			// check for any issues upgrading
			if err != nil {
				return fmt.Errorf("upgrade database %w", err)
			}
		}

		// update the version to the latest so our one time upgrade only runs once
		version.Current = LATEST_VERSION
		if err = c.updateVersion(version); err != nil {
			return fmt.Errorf("update version %w", err)
		}
	}

	return nil
}

// get the version of the database
func (c *txClient) getVersion() (Version, error) {
	item := Version{}
	if err := c.db.Get(&item, "SELECT * FROM versions"); err != nil {
		if err == sql.ErrNoRows {
			return Version{}, nil
		}
		return Version{}, err
	}

	return item, nil
}

func (c *txClient) updateVersion(version Version) error {
	if _, err := c.db.NamedExec(`
	UPDATE versions SET 
		updatedAt = CURRENT_TIMESTAMP,
		current = :current
	WHERE id = :id`, version); err != nil {
		return err
	}

	return nil
}

func (c *txClient) upgrade1() error {

	games, err := c.GetGames()
	if err != nil {
		return err
	}

	for _, game := range games {
		fg, err := c.GetFullGame(game.ID)
		if err != nil {
			return err
		}

		cleaner := cs.NewCleaner()
		cleaner.RemovePlayerDesignIntels(fg)

		// save changes to the DB
		if err := c.UpdateFullGame(fg); err != nil {
			return err
		}
	}

	return nil
}
