package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type Version struct {
	ID        int64     `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Current   int       `json:"current,omitempty"`
}

const LATEST_VERSION = 1

func (c *client) mustUpgrade() {
	if err := c.ensureUpgrade(); err != nil {
		log.Fatal().Err(err).Msg("failed to upgrade database data to latest version")
	}
}

func (c *client) ensureUpgrade() error {
	version, err := c.getVersion()
	if err != nil {
		return err
	}

	if version.Current < LATEST_VERSION {
		var err error
		// do all updates in a transaction in case we need to rollback
		tx, err := c.db.Beginx()
		if err != nil {
			return err
		}

		for current := version.Current; current < LATEST_VERSION; current++ {
			// check each version and call the upgrade functionality
			switch current {
			case 0:
				log.Info().Msgf("upgrading database data from v0 to v1")
				err = c.upgrade1(tx)
			}

			// check for any issues upgrading
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("upgrade database %w", err)
			}
		}

		// update the version to the latest so our one time upgrade only runs once
		version.Current = LATEST_VERSION
		if err = c.updateVersion(version, tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("update version %w", err)
		}

		if err := tx.Commit(); err != nil {
			return err
		}

	}

	return nil
}

// get the version of the database
func (c *client) getVersion() (Version, error) {
	item := Version{}
	if err := c.db.Get(&item, "SELECT * FROM versions"); err != nil {
		if err == sql.ErrNoRows {
			return Version{}, nil
		}
		return Version{}, err
	}

	return item, nil
}

func (c *client) updateVersion(version Version, tx SQLExecer) error {
	if _, err := tx.NamedExec(`
	UPDATE versions SET 
		updatedAt = CURRENT_TIMESTAMP,
		current = :current
	WHERE id = :id`, version); err != nil {
		return err
	}

	return nil
}

func (c *client) upgrade1(tx *sqlx.Tx) error {

	games, err := c.getGames(tx)
	if err != nil {
		return err
	}

	for _, game := range games {
		fg, err := c.getFullGame(tx, game.ID)
		if err != nil {
			return err
		}

		cleaner := cs.NewCleaner()
		cleaner.RemovePlayerDesignIntels(fg)

		// save changes to the DB
		if err := c.updateFullGame(fg, tx); err != nil {
			return err
		}
	}

	return nil
}
