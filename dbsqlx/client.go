package dbsqlx

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
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
		panic(fmt.Errorf("load schema file %s, %w", schemaPath, err))
	}

	schema := string(schemaBytes)
	log.Info().Msg("Creating schema")
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	c.db.MustExec(schema)
}

// helper to convert an item into JSON
func valueJSON(item interface{}) (driver.Value, error) {
	if isNil(item) {
		return nil, nil
	}

	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// helper to scan a text JSON column back into a struct
func scanJSON(src interface{}, dest interface{}) error {
	if src == nil {
		// leave empty
		return nil
	}

	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, dest)
	case string:
		return json.Unmarshal([]byte(v), dest)
	}
	return errors.New("type assertion failed")
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		//use of IsNil method
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
