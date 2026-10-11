package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/sirgwain/craig-stars/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	mattnsqlite3 "github.com/mattn/go-sqlite3"
)

//go:embed schema/filesystem/*.sql
var gamesSchemaFiles embed.FS

//go:embed schema/memory/*.sql
var memorySchemaFiles embed.FS

func (c *dbConn) mustMigrate(cfg *config.Config) {
	if c.databaseInMemory {
		// no migration for in memory dbs
		return
	}

	c.mustMigrateDatabase(cfg.Database.Filename, gamesSchemaFiles, "schema/filesystem")
}

// in memory databases are different because the user and games database has to live in the same
// memory space so we have to "migrate" them together
func (c *dbConn) setupInMemoryDatabase() {
	schema, err := iofs.New(memorySchemaFiles, "schema/memory")
	if err != nil {
		slog.Error("loading embedded schema", slog.Any("error", err))
		os.Exit(1)
	}

	config := &sqlite3.Config{
		MigrationsTable: "my_migration_table",
	}

	driver, err := sqlite3.WithInstance(c.dbRead, config)
	if err != nil {
		slog.Error("creating database driver", slog.Any("error", err))
		os.Exit(1)
	}

	m, err := migrate.NewWithInstance("iofs", schema, "users", driver)
	if err != nil {
		slog.Error("creating migration", slog.Any("error", err))
		os.Exit(1)
	}

	err = m.Up()
	if err != nil {
		slog.Error("migrating users", slog.Any("error", err))
		os.Exit(1)
	}

}

func (c *dbConn) mustMigrateDatabase(datasource string, fs fs.FS, path string) {
	if err := c.migrateDatabase(datasource, fs, path); err != nil {
		slog.Error("migrating database", slog.String("path", path), slog.Any("error", err))
		os.Exit(1)
	}
}

// migrateDatabase brings a database up to the latest schema, backing it up first if there is
// anything to migrate.
func (c *dbConn) migrateDatabase(datasource string, fs fs.FS, path string) (err error) {
	d, err := iofs.New(fs, path)
	if err != nil {
		return fmt.Errorf("loading embedded schema: %w", err)
	}

	latest, err := latestMigrationVersion(d)
	if err != nil {
		return fmt.Errorf("finding latest schema version: %w", err)
	}

	config := &sqlite3.Config{
		MigrationsTable: "my_migration_table",
	}

	db, err := sql.Open("sqlite3", datasource)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer func() {
		// close this db connection, we'll open a new joined connection after migration
		if closeErr := db.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("closing database after migration: %w", closeErr)
		}
	}()

	driver, err := sqlite3.WithInstance(db, config)
	if err != nil {
		return fmt.Errorf("creating database driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, datasource, driver)
	if err != nil {
		return fmt.Errorf("creating migration: %w", err)
	}

	version, dirty, err := m.Version()
	newDatabase := err == migrate.ErrNilVersion
	if err != nil && !newDatabase {
		return fmt.Errorf("get database version: %w", err)
	}

	slog.Info("database version", slog.String("path", path), slog.Int("version", int(version)), slog.Int("latest", int(latest)))

	// A migration that failed leaves the database marked dirty, and it stays that way until it's
	// fixed by hand. The server is restarted when it fails to start, so don't back up the
	// database again on every attempt. The backup from the first attempt is still there.
	if dirty {
		return fmt.Errorf("a previous migration to version %d failed, restore the database from the backup made before it", version)
	}

	// backups are as big as the database, so only make one if we're about to change it
	if !newDatabase && version != latest {
		c.mustBackup(datasource, version)
	}

	err = m.Up()
	switch err {
	case migrate.ErrNoChange:
		slog.Info("database no migration required", slog.String("path", path))
		return nil
	case nil:
		slog.Info("database migrated", slog.String("path", path))
		db.Exec("VACUUM;")
		slog.Info("database vacuumed", slog.String("path", path))
		return nil
	}
	return err
}

// latestMigrationVersion returns the version of the newest migration, or 0 if there aren't any
func latestMigrationVersion(d source.Driver) (uint, error) {
	version, err := d.First()
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	for err == nil {
		var next uint
		if next, err = d.Next(version); err == nil {
			version = next
		}
	}
	if !errors.Is(err, os.ErrNotExist) {
		return 0, err
	}
	return version, nil
}

func (c *dbConn) mustBackup(filename string, version uint) string {

	if strings.Contains(filename, ":memory") {
		slog.Debug("not backing up in memory db")
		return ""
	}

	// timestamp code from https://gist.github.com/rustyeddy/77f17f4f0fb83cc87115eb72a23f18f7
	ts := time.Now().UTC().Format(time.RFC3339)
	backup := fmt.Sprintf("%s.%d.%s", filename, version, strings.Replace(ts, ":", "", -1))

	// register a connect hook so we can get the sqlite3 connection for backup
	// code from here: https://github.com/mattn/go-sqlite3/blob/master/_example/hook/hook.go
	register := "sqlite3_backup_" + backup
	sqlite3conn := []*mattnsqlite3.SQLiteConn{}
	sql.Register(register, &mattnsqlite3.SQLiteDriver{
		ConnectHook: func(conn *mattnsqlite3.SQLiteConn) error {
			sqlite3conn = append(sqlite3conn, conn)
			return nil
		},
	})

	// connect to the source
	srcDb, err := sql.Open(register, filename)
	if err != nil {
		slog.Error("failed to connect to backup", slog.String("filename", filename), slog.Any("error", err))
		os.Exit(1)
	}
	defer srcDb.Close()
	srcDb.Ping()

	// connect to the dest
	destDb, err := sql.Open(register, backup)
	if err != nil {
		slog.Error("failed to connect to backup", slog.String("backup", backup), slog.Any("error", err))
		os.Exit(1)
	}
	defer destDb.Close()
	destDb.Ping()

	// perform the backup
	bk, err := sqlite3conn[1].Backup("main", sqlite3conn[0], "main")
	if err != nil {
		slog.Error("failed to backup", slog.String("backup", backup), slog.Any("error", err))
		os.Exit(1)
	}

	_, err = bk.Step(-1)
	if err != nil {
		slog.Error("failed backup step", slog.String("backup", backup), slog.Any("error", err))
		os.Exit(1)
	}

	if err := bk.Finish(); err != nil {
		slog.Error("failed backup finish", slog.String("backup", backup), slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("backed up database", slog.String("from", filename), slog.String("to", backup))

	return backup
}
