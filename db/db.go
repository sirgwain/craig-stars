package db

import (
	"database/sql"
	"errors"
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

bad

type client struct {
	db        *sqlx.DB
	converter Converter
}

type Client interface {
	Connect(config *config.Config) error
	ExecSQL(schemaPath string)

	GetUsers() ([]game.User, error)
	GetUser(id int64) (*game.User, error)
	GetUserByUsername(username string) (*game.User, error)
	CreateUser(user *game.User) error
	UpdateUser(user *game.User) error
	DeleteUser(id int64) error

	GetRaces() ([]game.Race, error)
	GetRacesForUser(userID int64) ([]game.Race, error)
	GetRace(id int64) (*game.Race, error)
	CreateRace(race *game.Race) error
	UpdateRace(race *game.Race) error
	DeleteRace(id int64) error

	GetTechStores() ([]game.TechStore, error)
	CreateTechStore(tech *game.TechStore) error
	GetTechStore(id int64) (*game.TechStore, error)

	GetRulesForGame(gameID int64) (*game.Rules, error)

	GetGames() ([]game.Game, error)
	GetGamesForHost(userID int64) ([]game.Game, error)
	GetGamesForUser(userID int64) ([]game.Game, error)
	GetOpenGames() ([]game.Game, error)
	GetGame(id int64) (*game.Game, error)
	GetFullGame(id int64) (*game.FullGame, error)
	CreateGame(game *game.Game) error
	UpdateGame(game *game.Game) error
	UpdateFullGame(g *game.FullGame) error
	DeleteGame(id int64) error

	GetPlayers() ([]game.Player, error)
	GetPlayersForUser(userID int64) ([]game.Player, error)
	GetPlayer(id int64) (*game.Player, error)
	GetLightPlayerForGame(gameID, userID int64) (*game.Player, error)
	GetPlayerForGame(gameID, userID int64) (*game.Player, error)
	GetFullPlayerForGame(gameID, userID int64) (*game.FullPlayer, error)
	CreatePlayer(player *game.Player) error
	UpdatePlayer(player *game.Player) error
	UpdateLightPlayer(player *game.Player) error
	DeletePlayer(id int64) error

	GetPlanet(id int64) (*game.Planet, error)
	UpdatePlanet(planet *game.Planet) error

	GetFleet(id int64) (*game.Fleet, error)
	UpdateFleet(fleet *game.Fleet) error
}

func NewClient() Client {
	return &client{
		converter: &GameConverter{},
	}
}

// interface to support NamedExec as either a transaction or db
type SQLExecer interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
}

func (c *client) Connect(config *config.Config) error {

	// exec the create schema sql if we are recreating the DB or using an in memory db
	execSchemaSql := config.Database.Recreate || config.Database.Filename == ":memory:"

	// if we are using a file based db, we have to exec the schema sql when we first
	// set it up
	if config.Database.Filename != ":memory:" {
		// check if the db exists
		info, err := os.Stat(config.Database.Filename)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		// delete the db and recreate it if we are configured for that
		if info != nil && config.Database.Recreate {
			log.Debug().Msgf("Deleting existing database %s", config.Database.Filename)
			os.Remove(config.Database.Filename)
		}

		if info == nil {
			// first time creating the db, so exec schema for it
			log.Debug().Msgf("Executing create schema %s", config.Database.Schema)
			execSchemaSql = true
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

	c.db = sqlx.NewDb(db, "sqlite3")

	if _, err := c.db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return err
	}

	// Create a new mapper which will use the struct field tag "json" instead of "db"
	c.db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	// recreate the schema if we are in memory or recreated the db
	if execSchemaSql {
		c.ExecSQL(config.Database.Schema)
	}
	return nil
}

// execute a schema file to create/update tables
func (c *client) ExecSQL(schemaPath string) {
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		panic(fmt.Errorf("load schema file %s, %w", schemaPath, err))
	}

	schema := string(schemaBytes)
	log.Info().Str("schemaPath", schemaPath).Msg("Executing sql")

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	c.db.MustExec(schema)
}
