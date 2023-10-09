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
002 - Ensure all AR planets have scanners
003 - Add random artifacts to undiscovered planets

*/

type Version struct {
	ID        int64     `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Current   int       `json:"current,omitempty"`
}

// game upgrader
type upgrade struct {
	tx *client
}

const LATEST_VERSION = 3

func (conn *dbConn) mustUpgrade() {

	if err := conn.WrapInTransaction(func(c Client) error {
		return c.ensureUpgrade()
	}); err != nil {
		panic(fmt.Sprintf("failed to upgrade database %v", err))
	}
}

func (tx *client) ensureUpgrade() error {
	version, err := tx.getVersion()
	if err != nil {
		return err
	}

	if version.Current < LATEST_VERSION {
		u := upgrade{tx: tx}
		for current := version.Current; current < LATEST_VERSION; current++ {
			log.Info().Msgf("upgrading database data from v%d to v%d", current, current+1)
			// check each version and call the upgrade functionality
			switch current {
			case 0:
				err = u.upgrade1()
			case 1:
				err = u.upgrade2()
			case 2:
				err = u.upgrade3()
			}

			// check for any issues upgrading
			if err != nil {
				return fmt.Errorf("upgrade database %w", err)
			}
		}

		// update the version to the latest so our one time upgrade only runs once
		version.Current = LATEST_VERSION
		if err = tx.updateVersion(version); err != nil {
			return fmt.Errorf("update version %w", err)
		}
	}

	return nil
}

// get the version of the database
func (c *client) getVersion() (Version, error) {
	item := Version{}
	if err := c.reader.Get(&item, "SELECT * FROM versions"); err != nil {
		if err == sql.ErrNoRows {
			return Version{}, nil
		}
		return Version{}, err
	}

	return item, nil
}

func (c *client) updateVersion(version Version) error {
	if _, err := c.writer.NamedExec(`
	UPDATE versions SET 
		updatedAt = CURRENT_TIMESTAMP,
		current = :current
	WHERE id = :id`, version); err != nil {
		return err
	}

	return nil
}

// helper function to get all games in the db and call an upgrade function on each game
// then save the game back to the db
func (u *upgrade) upgradeGames(upgradeGame func(fg *cs.FullGame) error) error {

	games, err := u.tx.GetGames()
	if err != nil {
		return err
	}

	for _, game := range games {
		fg, err := u.tx.GetFullGame(game.ID)
		if err != nil {
			return err
		}

		// call the passed in function
		if err := upgradeGame(fg); err != nil {
			return err
		}

		// save changes to the DB
		if err := u.tx.UpdateFullGame(fg); err != nil {
			return err
		}
	}
	return nil
}

func (u *upgrade) upgrade1() error {
	return u.upgradeGames(func(fg *cs.FullGame) error {
		cleaner := cs.NewCleaner()
		cleaner.RemovePlayerDesignIntels(fg)
		return nil
	})
}

func (u *upgrade) upgrade2() error {
	return u.upgradeGames(func(fg *cs.FullGame) error {
		cleaner := cs.NewCleaner()
		cleaner.AddScannerToInnateScannerPlanets(fg)
		return nil
	})
}

func (u *upgrade) upgrade3() error {
	return u.upgradeGames(func(fg *cs.FullGame) error {
		cleaner := cs.NewCleaner()
		cleaner.AddRandomArtifactsToPlanets(fg)
		return nil
	})
}
