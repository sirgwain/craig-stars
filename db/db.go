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

type dbConn struct {
	dbRead           *sqlx.DB
	dbWrite          *sqlx.DB
	databaseInMemory bool
	usersInMemory    bool
}

type client struct {
	reader    sqlReader
	writer    sqlWriter
	tx        *sqlx.Tx
	converter Converter
}

type sqlReader interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	Rebind(query string) string
}

type sqlWriter interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
}

type DBConn interface {
	Connect(config *config.Config) error
	Close() error

	// create a new read client
	NewReadClient() Client
	NewReadWriteClient() Client

	// for write clients we use transactions
	BeginTransaction() (Client, error)
	Rollback(c Client) error
	Commit(c Client) error

	// wrap a function call inside a transaction
	WrapInTransaction(wrap func(c Client) error) error
}

type Client interface {
	// private transaction management methods used by DBConn RollBack, Commit
	rollback() error
	commit() error

	// private method used during DBConn Connect to upgrade a client
	// this is
	ensureUpgrade() error

	GetUsers() ([]cs.User, error)
	GetUser(id int64) (*cs.User, error)
	GetUserByUsername(username string) (*cs.User, error)
	GetGuestUser(hash string) (*cs.User, error)
	GetGuestUserForGame(gameID int64, playerNum int) (*cs.User, error)
	GetGuestUsersForGame(gameID int64) ([]cs.User, error)
	CreateUser(user *cs.User) error
	UpdateUser(user *cs.User) error
	DeleteUser(id int64) error
	DeleteGameUsers(gameID int64) error
	GetUsersForGame(gameID int64) ([]cs.User, error)

	GetRaces() ([]cs.Race, error)
	GetRacesForUser(userID int64) ([]cs.Race, error)
	GetRace(id int64) (*cs.Race, error)
	CreateRace(race *cs.Race) error
	UpdateRace(race *cs.Race) error
	DeleteRace(id int64) error
	DeleteUserRaces(userID int64) error
	
	GetTechStores() ([]cs.TechStore, error)
	CreateTechStore(tech *cs.TechStore) error
	GetTechStore(id int64) (*cs.TechStore, error)

	GetRulesForGame(gameID int64) (*cs.Rules, error)

	GetGames() ([]cs.Game, error)
	GetGamesWithPlayers() ([]cs.GameWithPlayers, error)
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
	UpdateGameHost(gameID int64, hostId int64) error
	DeleteGame(id int64) error
	DeleteUserGames(hostID int64) error

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
	UpdatePlayerRelations(player *cs.Player) error
	UpdatePlayerSpec(player *cs.Player) error
	UpdatePlayerPlans(player *cs.Player) error
	UpdatePlayerSalvageIntels(player *cs.Player) error
	UpdateLightPlayer(player *cs.Player) error
	UpdatePlayerUserId(player *cs.Player) error
	DeletePlayer(id int64) error

	GetShipDesignsForPlayer(gameID int64, playerNum int) ([]*cs.ShipDesign, error)
	GetShipDesign(id int64) (*cs.ShipDesign, error)
	GetShipDesignByNum(gameID int64, playerNum, num int) (*cs.ShipDesign, error)
	CreateShipDesign(shipDesign *cs.ShipDesign) error
	UpdateShipDesign(shipDesign *cs.ShipDesign) error
	DeleteShipDesign(id int64) error

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
	GetFleetsOrbitingPlanet(gameID int64, planetNum int) ([]*cs.Fleet, error)
	
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

func NewConn() DBConn {
	return &dbConn{}
}

func (conn *dbConn) NewReadClient() Client {
	return &client{
		reader:    conn.dbRead,
		converter: &GameConverter{},
	}
}

func (conn *dbConn) NewReadWriteClient() Client {
	return &client{
		reader:    conn.dbRead,
		writer:    conn.dbWrite,
		converter: &GameConverter{},
	}
}

// create a new dbClient from a transaction
func newTransactionClient(tx *sqlx.Tx) *client {
	return &client{
		reader:    tx,
		writer:    tx,
		tx:        tx,
		converter: &GameConverter{},
	}
}

func (conn *dbConn) BeginTransaction() (Client, error) {

	// log.Debug().Msg("begin transaction")
	tx, err := conn.dbWrite.Beginx()
	if err != nil {
		return nil, err
	}
	return newTransactionClient(tx), nil
}

func (conn *dbConn) Rollback(c Client) error {
	// log.Debug().Msg("rollback transaction")
	return c.rollback()
}
func (conn *dbConn) Commit(c Client) error {
	// log.Debug().Msg("commit transaction")
	return c.commit()
}

// helper function to wrap a series of db calls in a transaction
func (conn *dbConn) WrapInTransaction(wrap func(c Client) error) error {
	c, err := conn.BeginTransaction()
	if err != nil {
		return err
	}
	defer func() { conn.Rollback(c) }()

	if err := wrap(c); err != nil {
		return err
	}

	return conn.Commit(c)
}

func (c *dbConn) Connect(cfg *config.Config) error {

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


	// create a new logger for logging database calls
	var zlogger zerolog.Logger
	if cfg.Database.DebugLogging {
		zlogger = zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	} else {
		zlogger = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.WarnLevel)
	}
	loggerAdapter := NewLoggerWithLogger(&zlogger)

	// dsn is like file::memory:?cache=shared, or file:data.db?_journal=WAL
	dsn := fmt.Sprintf("file:%s%s", cfg.Database.Filename, cfg.Database.ReadConnectionParams)
	log.Debug().Msgf("Connecting to database %s", dsn)
	connectHook := func(conn *sqlite3.SQLiteConn) error {
		log.Debug().Msgf("Attaching Users database %s", cfg.Database.UsersFilename)
		if _, err := conn.Exec(fmt.Sprintf("ATTACH DATABASE '%s' as users;", cfg.Database.UsersFilename), nil); err != nil {
			return err
		}
		if _, err := conn.Exec("PRAGMA foreign_keys = ON;", nil); err != nil {
			return err
		}
		return nil
	}

	dbRead := sqldblogger.OpenDriver(dsn, &sqlite3.SQLiteDriver{ConnectHook: connectHook}, loggerAdapter)
	dbWrite := sqldblogger.OpenDriver(dsn, &sqlite3.SQLiteDriver{ConnectHook: connectHook}, loggerAdapter)

	c.dbRead = sqlx.NewDb(dbRead, "sqlite3")
	if c.databaseInMemory {
		// no separate write connetion for in memory dbs
		c.dbWrite = c.dbRead
	} else {
		c.dbWrite = sqlx.NewDb(dbWrite, "sqlite3")
		c.dbWrite.SetMaxOpenConns(1)
	}

	// Create a new mapper which will use the struct field tag "json" instead of "db"
	c.dbRead.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	c.dbWrite.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	// do some special processing for in memory databases
	if c.databaseInMemory {
		c.setupInMemoryDatabase()
	}

	// make sure the data is updated
	if !cfg.Database.SkipUpgrade {
		c.mustUpgrade()
	}

	log.Info().Msg("connect() complete")

	return nil
}

func (c *dbConn) Close() error {
	if err := c.dbRead.Close(); err != nil {
		return err
	}
	if err := c.dbWrite.Close(); err != nil {
		return err
	}
	return nil
}

func (c *client) rollback() error {
	return c.tx.Rollback()
}

func (c *client) commit() error {
	return c.tx.Commit()
}
