package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	mattnsqlite3 "github.com/mattn/go-sqlite3"
)

//go:embed schema/users/*.sql
var usersSchemaFiles embed.FS

//go:embed schema/games/*.sql
var gamesSchemaFiles embed.FS

//go:embed schema/memory/*.sql
var memorySchemaFiles embed.FS

func (c *dbConn) mustMigrate(cfg *config.Config) {
	if !c.usersInMemory {
		c.mustMigrateDatabase(cfg.Database.UsersFilename, usersSchemaFiles, "schema/users", !cfg.Database.Recreate)
	}
	if !c.databaseInMemory {
		c.mustMigrateDatabase(cfg.Database.Filename, gamesSchemaFiles, "schema/games", !cfg.Database.Recreate)
	}
}

// in memory databases are different because the user and games database has to live in the same
// memory space so we have to "migrate" them together
func (c *dbConn) setupInMemoryDatabase() {
	schema, err := iofs.New(memorySchemaFiles, "schema/memory")
	if err != nil {
		log.Fatal().Err(err).Msg("loading embedded schema")
	}

	config := &sqlite3.Config{
		MigrationsTable: "my_migration_table",
	}

	driver, err := sqlite3.WithInstance(c.dbRead.DB, config)
	if err != nil {
		log.Fatal().Err(err).Msg("creating database driver")
	}

	m, err := migrate.NewWithInstance("iofs", schema, "users", driver)
	if err != nil {
		log.Fatal().Err(err).Msg("creating migration")
	}

	err = m.Up()
	if err != nil {
		log.Fatal().Err(err).Msg("migrating users")
	}

}

func (c *dbConn) mustMigrateDatabase(datasource string, fs embed.FS, path string, backup bool) {
	d, err := iofs.New(fs, path)
	if err != nil {
		log.Fatal().Err(err).Msg("loading embedded schema")
	}

	config := &sqlite3.Config{
		MigrationsTable: "my_migration_table",
	}

	db, err := sql.Open("sqlite3", datasource)
	defer func() {
		// close this db connection, we'll open a new joined connection after migration
		if err := db.Close(); err != nil {
			log.Fatal().Err(err).Msg("failed to close database after migration")
		}
	}()

	if err != nil {
		log.Fatal().Err(err).Msg("opening database")
	}

	driver, err := sqlite3.WithInstance(db, config)
	if err != nil {
		log.Fatal().Err(err).Msg("creating database driver")
	}

	m, err := migrate.NewWithInstance("iofs", d, datasource, driver)
	if err != nil {
		log.Fatal().Err(err).Msg("creating migration")
	}

	version, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatal().Err(err).Msg("get database version")
	}

	log.Info().Msgf("database %s is version %d", path, version)
	var backupFile string
	if backup {
		backupFile = c.mustBackup(datasource, version)
	}
	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Info().Msgf("database %s, no migration required", path)
		if backup {
			// remove the backup, we don't need it
			os.Remove(backupFile)
		}
	} else if err == nil {
		log.Info().Msgf("database %s migrated", path)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("migrating database")
	}

}

func (c *dbConn) mustBackup(filename string, version uint) string {

	if strings.Contains(filename, ":memory") {
		log.Debug().Msg("not backing up in memory db")
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
		log.Fatal().Err(err).Msgf("failed to connect to %s to backup", filename)
	}
	defer srcDb.Close()
	srcDb.Ping()

	// connect to the dest
	destDb, err := sql.Open(register, backup)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to connect to %s to backup", backup)
	}
	defer destDb.Close()
	destDb.Ping()

	// perform the backup
	bk, err := sqlite3conn[1].Backup("main", sqlite3conn[0], "main")
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to backup %s", backup)
	}

	_, err = bk.Step(-1)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed backup step %s", backup)
	}

	if err := bk.Finish(); err != nil {
		log.Fatal().Err(err).Msgf("failed backup finish %s", backup)
	}

	log.Info().Msgf("backed up database %s -> %s", filename, backup)

	return backup
}
