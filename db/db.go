package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	gen "github.com/sirgwain/craig-stars/db/generated"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	sqldblogger "github.com/simukti/sqldb-logger"
)

// DBConn represents a connection to the database
// Some database connections are read only, others are readwrite
// A readWrite connection can be wrapped in a transaction, but be warned, this locks the
// database until the transaction completes or fails.
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

// A database Client interface is used to make all calls that modify the database
type Client interface {
	// private transaction management methods used by DBConn RollBack, Commit
	rollback() error
	commit() error

	// private method used during DBConn Connect to upgrade a client
	// this is
	ensureUpgrade(context.Context) error

	CreateUser(ctx context.Context, user *cs.User) (*cs.User, error)
	DeleteGameUsers(ctx context.Context, gameID int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetGuestUser(ctx context.Context, hash string) (*cs.User, error)
	GetGuestUserForGame(ctx context.Context, gameID int64, playerNum int) (*cs.User, error)
	GetGuestUsersForGame(ctx context.Context, gameID int64) ([]cs.User, error)
	GetUser(ctx context.Context, id int64) (*cs.User, error)
	GetUserByUsername(ctx context.Context, username string) (*cs.User, error)
	GetUsers(ctx context.Context) ([]cs.User, error)
	GetUsersForGame(ctx context.Context, gameID int64) ([]cs.User, error)
	UpdateUser(ctx context.Context, user *cs.User) error
	UpdateUserSettings(ctx context.Context, user *cs.User) error

	CreateRace(ctx context.Context, race *cs.Race) (*cs.Race, error)
	DeleteRace(ctx context.Context, id int64) error
	DeleteUserRaces(ctx context.Context, userID int64) error
	GetRace(ctx context.Context, id int64) (*cs.Race, error)
	GetRaces(ctx context.Context) ([]cs.Race, error)
	GetRacesForUser(ctx context.Context, userID int64) ([]cs.Race, error)
	UpdateRace(ctx context.Context, race *cs.Race) error

	CreateTechStore(ctx context.Context, tech *cs.TechStore) (*cs.TechStore, error)
	GetTechStore(ctx context.Context, id int64) (*cs.TechStore, error)
	GetTechStores(ctx context.Context) ([]cs.TechStore, error)

	GetRulesForGame(ctx context.Context, gameID int64) (*cs.Rules, error)

	CreateGame(ctx context.Context, game *cs.Game) (*cs.Game, error)
	DeleteGame(ctx context.Context, id int64) error
	DeleteUserGames(ctx context.Context, hostID int64) error
	GetFullGame(ctx context.Context, id int64) (*cs.FullGame, error)
	GetGame(ctx context.Context, id int64) (*cs.GameWithPlayers, error)
	GetGames(ctx context.Context) ([]cs.Game, error)
	GetGamesForHost(ctx context.Context, userID int64) ([]cs.GameWithPlayers, error)
	GetGamesForUser(ctx context.Context, userID int64) ([]cs.GameWithPlayers, error)
	GetGamesWithPlayers(ctx context.Context) ([]cs.GameWithPlayers, error)
	GetOpenGames(ctx context.Context) ([]cs.GameWithPlayers, error)
	GetOpenGamesByHash(ctx context.Context, hash string) ([]cs.GameWithPlayers, error)
	UpdateFullGame(ctx context.Context, fullGame *cs.FullGame) error
	UpdateGame(ctx context.Context, game *cs.Game) error
	UpdateGameHost(ctx context.Context, gameID int64, hostId int64) error
	UpdateGameState(ctx context.Context, gameID int64, state cs.GameState) error

	ArchivePlayer(ctx context.Context, gameID int64, num int, archived bool) error
	CreatePlayer(ctx context.Context, player *cs.Player) (*cs.Player, error)
	DeletePlayer(ctx context.Context, id int64) error
	GetFullPlayerForGame(ctx context.Context, gameID int64, params GetPlayerParams) (*cs.FullPlayer, error)
	GetLightPlayerForGame(ctx context.Context, gameID int64, params GetPlayerParams) (*cs.Player, error)
	GetPlayer(ctx context.Context, id int64) (*cs.Player, error)
	GetPlayerForGame(ctx context.Context, gameID int64, params GetPlayerParams) (*cs.Player, error)
	GetPlayerMapObjects(ctx context.Context, gameID, userID int64) (*cs.PlayerMapObjects, error)
	GetPlayers(ctx context.Context) ([]*cs.Player, error)
	GetPlayersForUser(ctx context.Context, userID int64) ([]*cs.Player, error)
	GetPlayersStatusForGame(ctx context.Context, gameID int64) ([]*cs.Player, error)
	SubmitPlayerTurn(ctx context.Context, gameID int64, num int, submittedTurn bool) error
	UpdateLightPlayer(ctx context.Context, player *cs.Player) error
	UpdatePlayer(ctx context.Context, player *cs.Player) error
	UpdatePlayerCargoTransfers(ctx context.Context, player *cs.Player) error
	UpdatePlayerFleetIntels(ctx context.Context, player *cs.Player) error
	UpdatePlayerMineralPacketIntels(ctx context.Context, player *cs.Player) error
	UpdatePlayerOrders(ctx context.Context, player *cs.Player) error
	UpdatePlayerPlanetIntels(ctx context.Context, player *cs.Player) error
	UpdatePlayerPlans(ctx context.Context, player *cs.Player) error
	UpdatePlayerRelations(ctx context.Context, player *cs.Player) error
	UpdatePlayerSalvageIntels(ctx context.Context, player *cs.Player) error
	UpdatePlayerSpec(ctx context.Context, player *cs.Player) error
	UpdatePlayerUserId(ctx context.Context, player *cs.Player) error

	CreateShipDesign(ctx context.Context, shipDesign *cs.ShipDesign) (*cs.ShipDesign, error)
	DeleteShipDesign(ctx context.Context, id int64) error
	GetShipDesign(ctx context.Context, id int64) (*cs.ShipDesign, error)
	GetShipDesignByNum(ctx context.Context, gameID int64, playerNum, num int) (*cs.ShipDesign, error)
	GetShipDesignsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.ShipDesign, error)
	UpdateShipDesign(ctx context.Context, shipDesign *cs.ShipDesign) error

	GetPlanet(ctx context.Context, id int64) (*cs.Planet, error)
	GetPlanetByNum(ctx context.Context, gameID int64, num int) (*cs.Planet, error)
	GetPlanetsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.Planet, error)
	UpdatePlanet(ctx context.Context, planet *cs.Planet) error
	UpdatePlanetSpec(ctx context.Context, planet *cs.Planet) error

	CreateFleet(ctx context.Context, fleet *cs.Fleet) (*cs.Fleet, error)
	CreateUpdateOrDeleteFleets(ctx context.Context, gameID int64, fleets []*cs.Fleet) error
	DeleteFleet(ctx context.Context, id int64) error
	GetFleet(ctx context.Context, id int64) (*cs.Fleet, error)
	GetFleetByNum(ctx context.Context, gameID int64, playerNum int, num int) (*cs.Fleet, error)
	GetFleetsByNums(ctx context.Context, gameID int64, playerNum int, nums []int) ([]*cs.Fleet, error)
	GetFleetsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.Fleet, error)
	UpdateFleet(ctx context.Context, fleet *cs.Fleet) error

	GetMineField(ctx context.Context, id int64) (*cs.MineField, error)
	GetMineFieldByNum(ctx context.Context, gameID int64, playerNum int, num int) (*cs.MineField, error)
	GetMineFieldsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.MineField, error)
	UpdateMineField(ctx context.Context, fleet *cs.MineField) error

	GetMineralPacket(ctx context.Context, id int64) (*cs.MineralPacket, error)
	GetMineralPacketByNum(ctx context.Context, gameID int64, playerNum int, num int) (*cs.MineralPacket, error)
	GetMineralPacketsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.MineralPacket, error)
	UpdateMineralPacket(ctx context.Context, mineralPacket *cs.MineralPacket) error

	CreateSalvage(ctx context.Context, salvage *cs.Salvage) (*cs.Salvage, error)
	GetSalvageByNum(ctx context.Context, gameID int64, num int) (*cs.Salvage, error)
	GetSalvagesForGame(ctx context.Context, gameID int64) ([]*cs.Salvage, error)
	GetSalvagesForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.Salvage, error)
	UpdateSalvage(ctx context.Context, salvage *cs.Salvage) error
}

type dbConn struct {
	dbRead           *sqlx.DB
	dbWrite          *sqlx.DB
	databaseInMemory bool
}

type sqlReader interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	Rebind(query string) string
}

type client struct {
	tx        *sqlx.Tx
	readConn  sqlReader
	reader    *gen.Queries
	writer    *gen.Queries
	converter Converter
}

func NewConn() DBConn {
	return &dbConn{}
}

func (conn *dbConn) NewReadClient() Client {
	return &client{
		readConn:  conn.dbRead,
		reader:    gen.New(conn.dbRead),
		converter: c,
	}
}

func (conn *dbConn) NewReadWriteClient() Client {
	return &client{
		readConn:  conn.dbRead,
		reader:    gen.New(conn.dbRead),
		writer:    gen.New(conn.dbWrite),
		converter: c,
	}
}

// create a new dbClient from a transaction
func newTransactionClient(tx *sqlx.Tx) *client {
	return &client{
		tx:        tx,
		readConn:  tx,
		reader:    gen.New(tx),
		writer:    gen.New(tx),
		converter: c,
	}
}

func (conn *dbConn) BeginTransaction() (Client, error) {
	tx, err := conn.dbWrite.Beginx()
	if err != nil {
		return nil, err
	}
	return newTransactionClient(tx), nil
}

func (conn *dbConn) Rollback(c Client) error {
	return c.rollback()
}
func (conn *dbConn) Commit(c Client) error {
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
	// make sure the database is up to date
	c.mustMigrate(cfg)

	// create a new logger for logging database calls
	var zlogger zerolog.Logger
	if cfg.Database.DebugLogging {
		zlogger = zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	} else {
		zlogger = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.WarnLevel)
	}
	loggerAdapter := newLoggerWithLogger(&zlogger)

	// dsn is like file::memory:?cache=shared, or file:data.db?_journal=WAL
	dsn := fmt.Sprintf("file:%s%s", cfg.Database.Filename, cfg.Database.ReadConnectionParams)
	log.Debug().Msgf("Connecting to database %s", dsn)
	connectHook := func(conn *sqlite3.SQLiteConn) error {
		if c.databaseInMemory {
			// no need to attach
			return nil
		}
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
