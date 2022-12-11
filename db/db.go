package db

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/mattn/go-sqlite3"
	sqldblogger "github.com/simukti/sqldb-logger"
)

//go:embed schema.sql
var schema string

type client struct {
	db        *sqlx.DB
	converter Converter
}

type Client interface {
	Connect(config *config.Config) error
	ExecSQL(schemaPath string)

	GetUsers() ([]cs.User, error)
	GetUser(id int64) (*cs.User, error)
	GetUserByUsername(username string) (*cs.User, error)
	CreateUser(user *cs.User) error
	UpdateUser(user *cs.User) error
	DeleteUser(id int64) error

	GetRaces() ([]cs.Race, error)
	GetRacesForUser(userID int64) ([]cs.Race, error)
	GetRace(id int64) (*cs.Race, error)
	CreateRace(race *cs.Race) error
	UpdateRace(race *cs.Race) error
	DeleteRace(id int64) error

	GetTechStores() ([]cs.TechStore, error)
	CreateTechStore(tech *cs.TechStore) error
	GetTechStore(id int64) (*cs.TechStore, error)

	GetRulesForGame(gameID int64) (*cs.Rules, error)

	GetGames() ([]cs.Game, error)
	GetGamesForHost(userID int64) ([]cs.Game, error)
	GetGamesForUser(userID int64) ([]cs.Game, error)
	GetOpenGames(userID int64) ([]cs.Game, error)
	GetGame(id int64) (*cs.Game, error)
	GetFullGame(id int64) (*cs.FullGame, error)
	CreateGame(game *cs.Game) error
	UpdateGame(game *cs.Game) error
	UpdateFullGame(fullGame *cs.FullGame) error
	DeleteGame(id int64) error

	GetPlayers() ([]cs.Player, error)
	GetPlayersForUser(userID int64) ([]cs.Player, error)
	GetPlayer(id int64) (*cs.Player, error)
	GetLightPlayerForGame(gameID, userID int64) (*cs.Player, error)
	GetPlayerStatusesForGame(gameID int64) ([]*cs.Player, error)
	GetPlayerForGame(gameID, userID int64) (*cs.Player, error)
	GetFullPlayerForGame(gameID, userID int64) (*cs.FullPlayer, error)
	CreatePlayer(player *cs.Player) error
	UpdatePlayer(player *cs.Player) error
	UpdateLightPlayer(player *cs.Player) error
	DeletePlayer(id int64) error

	GetPlanet(id int64) (*cs.Planet, error)
	GetPlanetsForPlayer(playerID int64) ([]*cs.Planet, error)
	UpdatePlanet(planet *cs.Planet) error

	GetFleet(id int64) (*cs.Fleet, error)
	UpdateFleet(fleet *cs.Fleet) error

	GetMineField(id int64) (*cs.MineField, error)
	UpdateMineField(fleet *cs.MineField) error
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

	sql := schema
	if schemaPath != "" {
		schemaBytes, err := ioutil.ReadFile(schemaPath)
		if err != nil {
			panic(fmt.Errorf("load schema file %s, %w", schemaPath, err))
		}
		sql = string(schemaBytes)
		log.Info().Str("schemaPath", schemaPath).Msg("Executing sql")
	} else {
		log.Info().Msg("Executing embedded sql")
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	c.db.MustExec(sql)
}
