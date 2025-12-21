package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// mergeUsersIntoData merges the users.db into data.db
// for historical reasons, they were separated, but they no longer
// need to be
func mergeUsersIntoData(dataDSN, usersDSN string) error {
	if usersDSN == "" {
		return nil
	}
	if _, err := os.Stat(usersDSN); err != nil {
		return nil
	}

	db, err := sql.Open("sqlite3", dataDSN)
	if err != nil {
		return fmt.Errorf("open data.db: %w", err)
	}
	defer db.Close()

	// Attach OUTSIDE the transaction.
	absUsers, _ := filepath.Abs(usersDSN)
	if _, err := db.Exec(fmt.Sprintf(`ATTACH DATABASE %q AS users;`, absUsers)); err != nil {
		return fmt.Errorf("attach users db: %w", err)
	}
	// Ensure we always detach (also outside tx). If detach fails, log it and continue closing.
	defer func() {
		if _, derr := db.Exec(`DETACH DATABASE users;`); derr != nil {
			slog.Warn("detach users db failed", slog.Any("error", derr))
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`PRAGMA foreign_keys = OFF;`); err != nil {
		return fmt.Errorf("pragma foreign_keys off: %w", err)
	}

	// Copy users rows (id-preserving)
	if _, err := tx.Exec(`
		INSERT OR IGNORE INTO main.users (
			id, created_at, updated_at, username, password, email, verified, banned, role,
			last_login, discord_id, discord_avatar, discord_webhook_url, game_id, player_num
		)
		SELECT
			id, created_at, updated_at, username, password, email, verified, banned, role,
			last_login, discord_id, discord_avatar, discord_webhook_url, game_id, player_num
		FROM users.users;
	`); err != nil {
		return fmt.Errorf("copy users: %w", err)
	}

	// Copy races rows (spec is intentionally not copied)
	if _, err := tx.Exec(`
		INSERT OR IGNORE INTO main.races (
			id, created_at, updated_at, user_id, name, plural_name, spend_leftover_points_on,
			prt, lrts,
			hab_low_grav, hab_low_temp, hab_low_rad,
			hab_high_grav, hab_high_temp, hab_high_rad,
			growth_rate, pop_efficiency,
			factory_output, factory_cost, num_factories, factories_cost_less,
			immune_grav, immune_temp, immune_rad,
			mine_output, mine_cost, num_mines,
			research_cost_energy, research_cost_weapons, research_cost_propulsion,
			research_cost_construction, research_cost_electronics, research_cost_biotechnology,
			techs_start_high
		)
		SELECT
			id, created_at, updated_at, user_id, name, plural_name, spend_leftover_points_on,
			prt, lrts,
			hab_low_grav, hab_low_temp, hab_low_rad,
			hab_high_grav, hab_high_temp, hab_high_rad,
			growth_rate, pop_efficiency,
			factory_output, factory_cost, num_factories, factories_cost_less,
			immune_grav, immune_temp, immune_rad,
			mine_output, mine_cost, num_mines,
			research_cost_energy, research_cost_weapons, research_cost_propulsion,
			research_cost_construction, research_cost_electronics, research_cost_biotechnology,
			techs_start_high
		FROM users.races;
	`); err != nil {
		return fmt.Errorf("copy races: %w", err)
	}

	// Keep AUTOINCREMENT sequence aligned (optional but recommended)
	_, _ = tx.Exec(`
		INSERT OR REPLACE INTO main.sqlite_sequence(name, seq)
		SELECT name, seq FROM users.sqlite_sequence
		WHERE name IN ('users', 'races');
	`)

	if _, err := tx.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("pragma foreign_keys on: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	slog.Info("merged users.db into main db", slog.String("main", dataDSN), slog.String("users", usersDSN))
	return nil
}
