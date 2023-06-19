package db

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
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

//go:embed schema_users.sql
var schemaUsers string

type client struct {
	db        *sqlx.DB
	converter Converter
}

type Client interface {
	Connect(config *config.Config) error

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
	GetPlayerMapObjects(gameID, userID int64) (*cs.PlayerMapObjects, error)
	GetPlayerWithDesignsForGame(gameID, userID int64) (*cs.Player, error)
	CreatePlayer(player *cs.Player) error
	UpdatePlayer(player *cs.Player) error
	UpdatePlayerOrders(player *cs.Player) error
	UpdatePlayerPlans(player *cs.Player) error
	UpdateLightPlayer(player *cs.Player) error
	DeletePlayer(id int64) error

	GetShipDesignsForPlayer(gameID int64, playerNum int) ([]*cs.ShipDesign, error)
	GetShipDesign(id int64) (*cs.ShipDesign, error)
	GetShipDesignByNum(gameID int64, playerNum, num int) (*cs.ShipDesign, error)
	CreateShipDesign(shipDesign *cs.ShipDesign) error
	UpdateShipDesign(shipDesign *cs.ShipDesign) error
	DeleteShipDesign(id int64) error
	DeleteShipDesignWithFleets(id int64, fleetsToUpdate, fleetsToDelete []*cs.Fleet) error

	GetPlanet(id int64) (*cs.Planet, error)
	GetPlanetByNum(gameID int64, num int) (*cs.Planet, error)
	GetPlanetsForPlayer(gameID int64, playerNum int) ([]*cs.Planet, error)
	UpdatePlanet(planet *cs.Planet) error
	UpdatePlanetSpec(planet *cs.Planet) error

	GetFleet(id int64) (*cs.Fleet, error)
	GetFleetByNum(gameID int64, playerNum int, num int) (*cs.Fleet, error)
	GetFleetsByNums(gameID int64, playerNum int, nums []int) ([]*cs.Fleet, error)
	UpdateFleet(fleet *cs.Fleet) error
	CreateUpdateOrDeleteFleets(gameID int64, fleets []*cs.Fleet) error
	DeleteFleet(id int64) error
	GetFleetsForPlayer(gameID int64, playerNum int) ([]*cs.Fleet, error)

	GetMineField(id int64) (*cs.MineField, error)
	GetMineFieldsForPlayer(gameID int64, playerNum int) ([]*cs.MineField, error)
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
	execSchemaSql := config.Database.Recreate || strings.Contains(config.Database.Filename, ":memory:")
	execUsersSchema := strings.Contains(config.Database.UsersFilename, ":memory:")

	// if we are using a file based db, we have to exec the schema sql when we first
	// set it up
	if !strings.Contains(config.Database.Filename, ":memory:") {
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
			execSchemaSql = true
		}
	}

	if !strings.Contains(config.Database.UsersFilename, ":memory:") {
		info, err := os.Stat(config.Database.UsersFilename)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		if info == nil {
			// first time creating the db, so exec schema for it
			log.Debug().Msgf("Executing users create schema")
			execUsersSchema = true
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

	// attach the users database
	if _, err := c.db.Exec(fmt.Sprintf("ATTACH DATABASE '%s' as users;", config.Database.UsersFilename)); err != nil {
		return err
	}

	if _, err := c.db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return err
	}

	// Create a new mapper which will use the struct field tag "json" instead of "db"
	c.db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	// recreate the users schema if we are in memory or recreated the db
	if execUsersSchema {
		// exec the schema or fail
		log.Info().Msg("executing create user db schema")
		c.db.MustExec(schemaUsers)
	}

	// recreate the schema if we are in memory or recreated the db
	if execSchemaSql {
		log.Info().Msg("executing create db schema")
		c.db.MustExec(schema)
	}
	return nil
}
