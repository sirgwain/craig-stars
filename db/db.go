package db

import (
	"database/sql"
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

type dbClient struct {
	db               *sqlx.DB
	databaseInMemory bool
	usersInMemory    bool
}

type txClient struct {
	db               *sqlx.Tx
	converter        Converter
	databaseInMemory bool
	usersInMemory    bool
}

type DBClient interface {
	Connect(config *config.Config) error
	BeginTransaction() (Client, error)
	Rollback(c Client) error
	Commit(c Client) error

	// call within a single transaction
	GetUserByUsername(username string) (*cs.User, error)
	CreateUser(user *cs.User) error
	UpdateUser(user *cs.User) error
	CreateRace(race *cs.Race) error
}

type Client interface {
	// private transaction management methods
	rollback() error
	commit() error
	ensureUpgrade() error

	GetUsers() ([]cs.User, error)
	GetUser(id int64) (*cs.User, error)
	GetUserByUsername(username string) (*cs.User, error)
	CreateUser(user *cs.User) error
	UpdateUser(user *cs.User) error
	DeleteUser(id int64) error
	GetUsersForGame(gameID int64) ([]cs.User, error)

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
	GetGamesForHost(userID int64) ([]cs.GameWithPlayers, error)
	GetGamesForUser(userID int64) ([]cs.GameWithPlayers, error)
	GetOpenGames() ([]cs.GameWithPlayers, error)
	GetOpenGamesByHash(hash string) ([]cs.GameWithPlayers, error)
	GetGame(id int64) (*cs.GameWithPlayers, error)
	GetGameWithPlayersStatus(gameID int64) (*cs.GameWithPlayers, error)
	GetFullGame(id int64) (*cs.FullGame, error)
	CreateGame(game *cs.Game) error
	UpdateGame(game *cs.Game) error
	UpdateGameState(gameID int64, state cs.GameState) error
	UpdateFullGame(fullGame *cs.FullGame) error
	DeleteGame(id int64) error

	GetPlayers() ([]cs.Player, error)
	GetPlayersForUser(userID int64) ([]cs.Player, error)
	GetPlayer(id int64) (*cs.Player, error)
	GetLightPlayerForGame(gameID, userID int64) (*cs.Player, error)
	GetPlayersStatusForGame(gameID int64) ([]*cs.Player, error)
	GetPlayerForGame(gameID, userID int64) (*cs.Player, error)
	GetPlayerIntelsForGame(gameID, userID int64) (*cs.PlayerIntels, error)
	GetPlayerByNum(gameID int64, num int) (*cs.Player, error)
	GetFullPlayerForGame(gameID, userID int64) (*cs.FullPlayer, error)
	GetPlayerMapObjects(gameID, userID int64) (*cs.PlayerMapObjects, error)
	GetPlayerWithDesignsForGame(gameID int64, num int) (*cs.Player, error)
	CreatePlayer(player *cs.Player) error
	UpdatePlayer(player *cs.Player) error
	SubmitPlayerTurn(gameID int64, num int, submittedTurn bool) error
	UpdatePlayerOrders(player *cs.Player) error
	UpdatePlayerPlans(player *cs.Player) error
	UpdatePlayerSalvageIntels(player *cs.Player) error
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

	GetMineralPacket(id int64) (*cs.MineralPacket, error)
	GetMineralPacketsForPlayer(gameID int64, playerNum int) ([]*cs.MineralPacket, error)

	GetSalvagesForGame(gameID int64) ([]*cs.Salvage, error)
	GetSalvagesForPlayer(gameID int64, playerNum int) ([]*cs.Salvage, error)
	GetSalvageByNum(gameID int64, num int) (*cs.Salvage, error)
	CreateSalvage(salvage *cs.Salvage) error
	UpdateSalvage(salvage *cs.Salvage) error
}

func NewClient() DBClient {
	return &dbClient{}
}

func newTransaction(tx *sqlx.Tx) *txClient {
	return &txClient{
		db:        tx,
		converter: &GameConverter{},
	}
}

func (dbClient *dbClient) BeginTransaction() (Client, error) {

	tx, err := dbClient.db.Beginx()
	if err != nil {
		return nil, err
	}

	return newTransaction(tx), nil
}

func (dbClient *dbClient) Rollback(c Client) error {
	return c.rollback()
}
func (dbClient *dbClient) Commit(c Client) error {
	return c.commit()
}

// interface to support NamedExec as either a transaction or db
type SQLExecer interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
}

type SQLSelector interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
}

func (c *dbClient) Connect(cfg *config.Config) error {

	c.databaseInMemory = strings.Contains(cfg.Database.Filename, ":memory:")
	c.usersInMemory = strings.Contains(cfg.Database.UsersFilename, ":memory:")
	// if we are using a file based db, we have to exec the schema sql when we first
	// set it up
	if !c.databaseInMemory && cfg.Database.Recreate {
		// check if the db exists
		info, err := os.Stat(cfg.Database.Filename)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		// delete the db and recreate it if we are configured for that
		if info != nil {
			log.Debug().Msgf("Deleting existing database %s", cfg.Database.Filename)
			os.Remove(cfg.Database.Filename)
		}
	}

	// make sure the database is up to date
	c.mustMigrate(cfg)

	log.Debug().Msgf("Connecting to database %s", cfg.Database.Filename)

	// create a new logger for logging database calls
	var zlogger zerolog.Logger
	if cfg.Database.DebugLogging {
		zlogger = zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	} else {
		zlogger = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.WarnLevel)
	}
	loggerAdapter := NewLoggerWithLogger(&zlogger)

	db := sqldblogger.OpenDriver(cfg.Database.Filename, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {

			if _, err := conn.Exec(fmt.Sprintf("ATTACH DATABASE '%s' as users;", cfg.Database.UsersFilename), nil); err != nil {
				return err
			}
			if _, err := conn.Exec("PRAGMA foreign_keys = ON;", nil); err != nil {
				return err
			}
			return nil
		},
	}, loggerAdapter /*, ...options */)

	c.db = sqlx.NewDb(db, "sqlite3")

	// Create a new mapper which will use the struct field tag "json" instead of "db"
	c.db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	// do some special processing for in memory databases
	if c.databaseInMemory {
		c.setupInMemoryDatabase()
	}

	// make sure the data is updated
	if !cfg.Database.SkipUpgrade {
		c.mustUpgrade()
	}

	return nil
}

func (c *txClient) rollback() error {
	return c.db.Rollback()
}

func (c *txClient) commit() error {
	return c.db.Commit()
}

// helper function to wrap code in a transaction
func (dbClient *dbClient) wrapInTransaction(wrap func(c Client) error) error {
	c, err := dbClient.BeginTransaction()
	if err != nil {
		return err
	}
	defer func() { dbClient.Rollback(c) }()

	if err := wrap(c); err != nil {
		return err
	}

	return dbClient.Commit(c)
}

func (dbClient *dbClient) GetUserByUsername(username string) (*cs.User, error) {
	var user *cs.User

	err := dbClient.wrapInTransaction(func(c Client) error {
		var err error
		user, err = c.GetUserByUsername(username)
		return err
	})

	return user, err
}

func (dbClient *dbClient) CreateUser(user *cs.User) error {
	return dbClient.wrapInTransaction(func(c Client) error { return c.CreateUser(user) })
}

func (dbClient *dbClient) UpdateUser(user *cs.User) error {
	return dbClient.wrapInTransaction(func(c Client) error { return c.UpdateUser(user) })
}

func (dbClient *dbClient) CreateRace(race *cs.Race) error {
	return dbClient.wrapInTransaction(func(c Client) error { return c.CreateRace(race) })
}
