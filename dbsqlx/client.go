package dbsqlx

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/game"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/mattn/go-sqlite3"
	sqldblogger "github.com/simukti/sqldb-logger"
)

type client struct {
	db        *sqlx.DB
	converter Converter
}

type Client interface {
	Connect(config *config.Config) error
	ExecSchema(schemaPath string)

	GetUsers() ([]game.User, error)
	GetUser(id int64) (*game.User, error)
	GetUserByUsername(username string) (*game.User, error)
	CreateUser(user *game.User) error
	UpdateUser(user *game.User) error
	DeleteUser(id int64) error

	GetRaces() ([]game.Race, error)
	GetRace(id int64) (*game.Race, error)
	CreateRace(race *game.Race) error
	UpdateRace(race *game.Race) error
	DeleteRace(id int64) error
}

func NewClient() Client {
	return &client{
		converter: &GameConverter{},
	}
}

// interface to support NamedExec as either a transaction or db
type NamedExecer interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

func (c *client) Connect(config *config.Config) error {
	if config.Database.Recreate && config.Database.Filename != ":memory:" {
		info, _ := os.Stat(config.Database.Filename)
		if info != nil {
			log.Debug().Msgf("Deleting existing database %s", config.Database.Filename)
			os.Remove(config.Database.Filename)
		}
	}

	log.Debug().Msgf("Connecting to database %s", config.Database.Filename)

	// create a new logger for logging database calls
	var zlogger zerolog.Logger
	if config.Database.DebugLogging {
		zlogger = zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	} else {
		zlogger = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.WarnLevel)
	}
	loggerAdapter := NewLoggerWithLogger(&zlogger)
	db := sqldblogger.OpenDriver(config.Database.Filename, &sqlite3.SQLiteDriver{}, loggerAdapter /*, ...options */)

	localdb := sqlx.NewDb(db, "sqlite3")

	if _, err := localdb.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return err
	}

	// localc.MapperFunc(MapperFunc()) // snake_case
	// localc.MapperFunc(func(s string) string { return s }) // identity

	// Create a new mapper which will use the struct field tag "json" instead of "db"
	localdb.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	c.db = localdb
	return nil
}

// execute a schema file to create/update tables
func (c *client) ExecSchema(schemaPath string) {
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		panic(fmt.Errorf("failed to load schema file %s, %w", schemaPath, err))
	}

	schema := string(schemaBytes)
	log.Info().Msg("Creating schema")
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	c.db.MustExec(schema)
}
