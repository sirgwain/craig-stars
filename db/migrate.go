package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/sirgwain/craig-stars/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
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

func (c *dbConn) mustMigrateDatabase(datasource string, fs embed.FS, path string) {
	d, err := iofs.New(fs, path)
	if err != nil {
		slog.Error("loading embedded schema", slog.Any("error", err))
		os.Exit(1)
	}

	config := &sqlite3.Config{
		MigrationsTable: "my_migration_table",
	}

	db, err := sql.Open("sqlite3", datasource)
	defer func() {
		// close this db connection, we'll open a new joined connection after migration
		if err := db.Close(); err != nil {
			slog.Error("failed to close database after migration", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	if err != nil {
		slog.Error("opening database", slog.Any("error", err))
		os.Exit(1)
	}

	driver, err := sqlite3.WithInstance(db, config)
	if err != nil {
		slog.Error("creating database driver", slog.Any("error", err))
		os.Exit(1)
	}

	m, err := migrate.NewWithInstance("iofs", d, datasource, driver)
	if err != nil {
		slog.Error("creating migration", slog.Any("error", err))
		os.Exit(1)
	}

	version, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		slog.Error("get database version", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("database version", slog.String("path", path), slog.Int("version", int(version)))
	backupFile := c.mustBackup(datasource, version)
	err = m.Up()
	switch err {
	case migrate.ErrNoChange:
		slog.Info("database no migration required", slog.String("path", path))
		// remove the backup, we don't need it
		os.Remove(backupFile)
	case nil:
		slog.Info("database migrated", slog.String("path", path))
		db.Exec("VACUUM;")
		slog.Info("database vacuumed", slog.String("path", path))
	}

	if err != nil && err != migrate.ErrNoChange {
		slog.Error("migrating database", slog.Any("error", err))
		os.Exit(1)
	}
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
	register := "sqlite3_backup_" + filename
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
