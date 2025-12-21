package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db/generated"
	"github.com/sirgwain/craig-stars/hash"
)

/*
Upgrade data in the database based on game updates

Version Info:

001 - Fix player discovers own starbase designs on planet discovery and adds them to intel
002 - Ensure all AR planets have scanners
003 - Add random artifacts to undiscovered planets
004 - Set BaseHab of Homeworlds to Hab. They were accidentally 0
005 - Update MineralConcentrations on all worlds where they are low
006 - Cleanup duplicate discord_id users

*/

type Version = generated.Version

// game upgrader
type upgrade struct {
	tx *client
}

const LATEST_VERSION = int64(6)

func (conn *dbConn) mustUpgrade() {

	if err := conn.WrapInTransaction(func(c Client) error {
		return c.ensureUpgrade(context.Background())
	}); err != nil {
		panic(fmt.Sprintf("failed to upgrade database %v", err))
	}
}

func (tx *client) ensureUpgrade(ctx context.Context) error {
	version, err := tx.getVersion(ctx)
	if err != nil {
		return err
	}

	if version.Current >= LATEST_VERSION {
		// already at latest version; no need to upgrade further
		return nil
	}

	u := upgrade{tx: tx}
	for current := version.Current; current < LATEST_VERSION; current++ {
		slog.Info("upgrading database data", slog.Int64("from", current), slog.Int64("to", current+1))
		// check each version and call the upgrade functionality
		switch current {
		case 0:
			//? Maybe make the starter database version -1?
			// That would make the switch marginally cleaner
			if u.initStarterDB(ctx); err != nil {
				return fmt.Errorf("initializing starter database failed: %w", err)
			}
			err = u.upgrade1(ctx)
		case 1:
			err = u.upgrade2(ctx)
		case 2:
			err = u.upgrade3(ctx)
		case 3:
			err = u.upgrade4(ctx)
		case 4:
			err = u.upgrade5(ctx)
		case 5:
			err = u.upgrade6(ctx)
		}

		// check for any issues upgrading
		if err != nil {
			return fmt.Errorf("upgrading database from v%d to v%d failed: %w", current, current+1, err)
		}
	}

	// update the version to the latest so our one time upgrade only runs once
	version.Current = LATEST_VERSION
	if err = tx.updateVersion(ctx, version); err != nil {
		return fmt.Errorf("updating to latest version failed: %w", err)
	}

	return nil
}

// get the version of the database
func (c *client) getVersion(ctx context.Context) (Version, error) {
	item, err := c.reader.GetVersion(ctx)
	if err == sql.ErrNoRows {
		return Version{}, nil
	}
	if err != nil {
		return Version{}, nil
	}

	return item, nil
}

func (c *client) updateVersion(ctx context.Context, version Version) error {
	return c.writer.UpdateVersion(ctx, generated.UpdateVersionParams{
		ID:      version.ID,
		Current: int64(version.Current),
	})
}

// helper function to get all games in the db and call an upgrade function on each game
// then save the game back to the db
func (u *upgrade) upgradeGames(ctx context.Context, upgradeGame func(fg *cs.FullGame) error) error {

	games, err := u.tx.GetGames(ctx)
	if err != nil {
		return fmt.Errorf("error while getting all games: %w", err)
	}

	for _, game := range games {
		fg, err := u.tx.GetFullGame(ctx, game.ID)
		if err != nil {
			return fmt.Errorf("retrieving fullGame with ID %d failed: %w", game.ID, err)
		}

		// call the passed in function
		if err := upgradeGame(fg); err != nil {
			return fmt.Errorf("upgrading fullGame with ID %d failed: %w", game.ID, err)
		}

		// save changes to the DB
		if err := u.tx.UpdateFullGame(ctx, fg); err != nil {
			return fmt.Errorf("updating fullGame with ID %d failed: %w", game.ID, err)

		}
	}
	return nil
}

func (u *upgrade) initStarterDB(ctx context.Context) error {
	slog.Info("initializing starter database with admin user, 'admin' password")
	password, err := hash.HashPassword("admin")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := cs.NewUser("admin", password, "", cs.RoleAdmin)

	// create the admin user, 'admin' password
	newUser, err := u.tx.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	rules := cs.NewRules()
	if err := u.tx.SaveRace(ctx, cs.NewRace().WithUserID(newUser.ID).WithSpec(&rules)); err != nil {
		return err
	}

	return nil
}

func (u *upgrade) upgrade1(ctx context.Context) error {
	return u.upgradeGames(ctx, func(fg *cs.FullGame) error {
		cleaner := cs.NewCleaner()
		cleaner.RemovePlayerDesignIntels(fg)
		return nil
	})
}

func (u *upgrade) upgrade2(ctx context.Context) error {
	return u.upgradeGames(ctx, func(fg *cs.FullGame) error {
		cleaner := cs.NewCleaner()
		cleaner.AddScannerToInnateScannerPlanets(fg)
		return nil
	})
}

func (u *upgrade) upgrade3(ctx context.Context) error {
	return u.upgradeGames(ctx, func(fg *cs.FullGame) error {
		cleaner := cs.NewCleaner()
		cleaner.AddRandomArtifactsToPlanets(fg)
		return nil
	})
}

func (u *upgrade) upgrade4(ctx context.Context) error {
	return u.upgradeGames(ctx, func(fg *cs.FullGame) error {
		cleaner := cs.NewCleaner()
		cleaner.ResetHomeworldBaseHab(fg)
		return nil
	})
}

func (u *upgrade) upgrade5(ctx context.Context) error {
	return u.upgradeGames(ctx, func(fg *cs.FullGame) error {
		cleaner := cs.NewCleaner()
		cleaner.FixMineralConc(fg)
		return nil
	})
}

func (u *upgrade) upgrade6(ctx context.Context) error {
	// Find duplicate discord_id groups with dupe_id=min(id), keep_id=max(id), cnt users
	tx := u.tx.tx
	rows, err := tx.QueryContext(ctx, `
	WITH dupes AS (
		SELECT
			discord_id,
			MIN(id) AS dupe_id,
			MAX(id) AS keep_id,
			COUNT(*) AS cnt
		FROM users
		WHERE discord_id <> ''
		GROUP BY discord_id
		HAVING COUNT(*) > 1
	)
	SELECT
		d.discord_id,
		d.dupe_id,
		du.username AS dupe_username,
		d.keep_id,
		ku.username AS keep_username,
		d.cnt
	FROM dupes d
	JOIN users du ON du.id = d.dupe_id
	JOIN users ku ON ku.id = d.keep_id;
`)
	if err != nil {
		return err
	}
	defer rows.Close()

	type pair struct {
		DiscordID    string
		DupeID       int64
		DupeUsername string
		KeepID       int64
		KeepUsername string
		Count        int64
	}

	var pairs []pair
	for rows.Next() {
		var p pair
		if err := rows.Scan(
			&p.DiscordID,
			&p.DupeID,
			&p.DupeUsername,
			&p.KeepID,
			&p.KeepUsername,
			&p.Count,
		); err != nil {
			return err
		}
		pairs = append(pairs, p)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// If you really only expect pairs, enforce it.
	for _, p := range pairs {
		if p.Count != 2 {
			return fmt.Errorf("discord_id %s has %d users; expected 2", p.DiscordID, p.Count)
		}
		// Rewire references
		if _, err := tx.ExecContext(ctx, `UPDATE races SET user_id=? WHERE user_id=?`, p.KeepID, p.DupeID); err != nil {
			return fmt.Errorf("update races %d->%d: %w", p.DupeID, p.KeepID, err)
		}
		if _, err := tx.ExecContext(ctx, `UPDATE players SET user_id=? WHERE user_id=?`, p.KeepID, p.DupeID); err != nil {
			return fmt.Errorf("update players %d->%d: %w", p.DupeID, p.KeepID, err)
		}
		if _, err := tx.ExecContext(ctx, `UPDATE games SET host_id=? WHERE host_id=?`, p.KeepID, p.DupeID); err != nil {
			return fmt.Errorf("update games.host_id %d->%d: %w", p.DupeID, p.KeepID, err)
		}

		// Optionally merge some fields from the old into the new (only if empty on keep)
		// Example: discord_webhook_url
		_, _ = tx.ExecContext(ctx, `
			UPDATE users
			SET discord_webhook_url = (
				SELECT discord_webhook_url FROM users WHERE id = ?
			)
			WHERE id = ?
			  AND (discord_webhook_url IS NULL OR discord_webhook_url = '')
			  AND (SELECT discord_webhook_url FROM users WHERE id = ?) <> '';
		`, p.DupeID, p.KeepID, p.DupeID)

		// Delete the dupe user
		if _, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id=?`, p.DupeID); err != nil {
			return fmt.Errorf("delete dupe user %d: %w", p.DupeID, err)
		}
		slog.InfoContext(ctx, "De-duped user",
			"Username", p.KeepUsername,
			"ID", p.KeepID,
			"DupeUsername", p.DupeUsername,
			"DupeID", p.DupeID,
		)
	}

	// Ensure nothing remains duplicated
	var remaining int
	if err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM (
			SELECT 1
			FROM users
			WHERE discord_id <> ''
			GROUP BY discord_id
			HAVING COUNT(*) > 1
		)
	`).Scan(&remaining); err != nil {
		return err
	}
	if remaining != 0 {
		return fmt.Errorf("discord_id duplicates remain after merge: %d groups", remaining)
	}

	// Add partial unique index
	if _, err := tx.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_discord_id_unique
		ON users(discord_id)
		WHERE discord_id <> '';
	`); err != nil {
		return fmt.Errorf("create unique index: %w", err)
	}

	return nil
}
