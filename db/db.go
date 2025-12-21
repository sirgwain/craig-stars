package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	gen "github.com/sirgwain/craig-stars/db/generated"

	"log/slog"

	"github.com/mattn/go-sqlite3"
	sqldblogger "github.com/simukti/sqldb-logger"
)

// DBConn represents a connection to the database
// Some database connections are read only, others are readwrite
// A readWrite connection can be wrapped in a transaction, but be warned, this locks the
// database until the transaction completes or fails.
type DBConn interface {
	Connect(ctx context.Context, config *config.Config) error
	Close() error

	// create a new read client
	NewReadClient() ReadClient
	NewReadWriteClient() Client

	// wrap a function call inside a transaction
	WrapInTransaction(wrap func(c Client) error) error
}

type ReadClient interface {
	GetGuestUser(ctx context.Context, hash string) (*cs.User, error)
	GetGuestUserForGame(ctx context.Context, gameID int64, playerNum int) (*cs.User, error)
	GetGuestUsersForGame(ctx context.Context, gameID int64) ([]cs.User, error)
	GetUser(ctx context.Context, id int64) (*cs.User, error)
	GetUserByUsername(ctx context.Context, username string) (*cs.User, error)
	GetUserByDiscordID(ctx context.Context, discord_id string) (*cs.User, error)
	GetUsers(ctx context.Context) ([]cs.User, error)
	GetUsersForGame(ctx context.Context, gameID int64) ([]cs.User, error)

	GetRace(ctx context.Context, id int64) (*cs.Race, error)
	GetRaces(ctx context.Context) ([]cs.Race, error)
	GetRacesForUser(ctx context.Context, userID int64) ([]cs.Race, error)

	GetTechStore(ctx context.Context, id int64) (*cs.TechStore, error)
	GetTechStores(ctx context.Context) ([]cs.TechStore, error)

	GetRulesForGame(ctx context.Context, gameID int64) (*cs.Rules, error)

	GetFullGame(ctx context.Context, id int64) (*cs.FullGame, error)
	GetGame(ctx context.Context, id int64) (*cs.GameWithPlayers, error)
	GetGameByHash(ctx context.Context, hash string) (*cs.GameWithPlayers, error)
	GetGames(ctx context.Context) ([]cs.Game, error)
	GetGamesForHost(ctx context.Context, userID int64) ([]cs.Game, error)
	GetGamesForUser(ctx context.Context, userID int64) ([]cs.GameWithPlayers, error)
	GetGamesWithPlayers(ctx context.Context) ([]cs.GameWithPlayers, error)
	GetOpenGames(ctx context.Context) ([]cs.GameWithPlayers, error)

	GetLightPlayerForGame(ctx context.Context, gameID int64, params GetPlayerParams) (*cs.Player, error)
	GetLightPlayerForGameWithDesigns(ctx context.Context, gameID int64, params GetPlayerParams) (*cs.Player, error)
	GetPlayer(ctx context.Context, id int64) (*cs.Player, error)
	GetPlayerForGame(ctx context.Context, gameID int64, playerNum int) (*cs.Player, error)
	GetPlayerIntel(ctx context.Context, gameID int64, playerNum int) (*cs.Intels, error)
	GetPlayerMapObjects(ctx context.Context, gameID int64, playerNum int) (*cs.PlayerMapObjects, error)
	GetPlayers(ctx context.Context) ([]*cs.Player, error)
	GetPlayersForUser(ctx context.Context, userID int64) ([]*cs.Player, error)
	GetPlayersStatusForGame(ctx context.Context, gameID int64) ([]*cs.Player, error)

	GetShipDesign(ctx context.Context, id int64) (*cs.ShipDesign, error)
	GetShipDesignByNum(ctx context.Context, gameID int64, playerNum, num int) (*cs.ShipDesign, error)
	GetShipDesignsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.ShipDesign, error)

	GetPlanet(ctx context.Context, id int64) (*cs.Planet, error)
	GetPlanetByNum(ctx context.Context, gameID int64, num int) (*cs.Planet, error)
	GetPlanetsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.Planet, error)

	GetFleet(ctx context.Context, id int64) (*cs.Fleet, error)
	GetFleetByNum(ctx context.Context, gameID int64, playerNum int, num int) (*cs.Fleet, error)
	GetFleetsByNums(ctx context.Context, gameID int64, playerNum int, nums []int) ([]*cs.Fleet, error)
	GetFleetsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.Fleet, error)

	GetMinefield(ctx context.Context, id int64) (*cs.Minefield, error)
	GetMinefieldByNum(ctx context.Context, gameID int64, playerNum int, num int) (*cs.Minefield, error)
	GetMinefieldsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.Minefield, error)

	GetMineralPacket(ctx context.Context, id int64) (*cs.MineralPacket, error)
	GetMineralPacketByNum(ctx context.Context, gameID int64, playerNum int, num int) (*cs.MineralPacket, error)
	GetMineralPacketsForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.MineralPacket, error)

	GetSalvageByNum(ctx context.Context, gameID int64, num int) (*cs.Salvage, error)
	GetSalvagesForGame(ctx context.Context, gameID int64) ([]*cs.Salvage, error)
	GetSalvagesForPlayer(ctx context.Context, gameID int64, playerNum int) ([]*cs.Salvage, error)
}

type WriteClient interface {
	// private method used during DBConn Connect to upgrade a client
	ensureUpgrade(context.Context) error

	CreateUser(ctx context.Context, user *cs.User) (*cs.User, error)
	DeleteGameGuestUsers(ctx context.Context, gameID int64) error
	DeleteUser(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, user *cs.User) error
	UpdateUserSettings(ctx context.Context, user *cs.User) error

	DeleteRace(ctx context.Context, id int64) error
	DeleteUserRaces(ctx context.Context, userID int64) error
	SaveRace(ctx context.Context, race *cs.Race) error

	CreateTechStore(ctx context.Context, tech *cs.TechStore) (*cs.TechStore, error)

	DeleteGame(ctx context.Context, id int64) error
	DeleteUserGames(ctx context.Context, hostID int64) error
	UpdateFullGame(ctx context.Context, fullGame *cs.FullGame) error
	SaveGame(ctx context.Context, game *cs.Game) error
	UpdateGameHost(ctx context.Context, gameID int64, hostID int64) error
	UpdateGameState(ctx context.Context, gameID int64, state cs.GameState) error

	ArchivePlayer(ctx context.Context, gameID int64, num int, archived bool) error
	DeletePlayer(ctx context.Context, id int64) error
	SavePlayer(ctx context.Context, player *cs.Player) error
	SubmitPlayerTurn(ctx context.Context, gameID int64, num int, submittedTurn bool) error
	UpdateLightPlayer(ctx context.Context, player *cs.Player) error
	UpdatePlayerCargoTransfers(ctx context.Context, player *cs.Player) error
	UpdatePlayerFleetIntels(ctx context.Context, player *cs.Player) error
	UpdatePlayerMineralPacketIntels(ctx context.Context, player *cs.Player) error
	UpdatePlayerOrders(ctx context.Context, player *cs.Player) error
	UpdatePlayerPlanetIntels(ctx context.Context, player *cs.Player) error
	UpdatePlayerPlans(ctx context.Context, player *cs.Player) error
	UpdatePlayerRelations(ctx context.Context, player *cs.Player) error
	UpdatePlayerSalvageIntels(ctx context.Context, player *cs.Player) error
	UpdatePlayerUserID(ctx context.Context, player *cs.Player) error

	DeleteShipDesign(ctx context.Context, id int64) error
	SaveShipDesign(ctx context.Context, shipDesign *cs.ShipDesign) error

	SavePlanet(ctx context.Context, planet *cs.Planet) error
	UpdatePlanetSpec(ctx context.Context, planet *cs.Planet) error

	DeleteFleet(ctx context.Context, id int64) error
	SaveFleet(ctx context.Context, fleet *cs.Fleet) error

	SaveMinefield(ctx context.Context, minefield *cs.Minefield) error

	SaveMineralPacket(ctx context.Context, mineralPacket *cs.MineralPacket) error

	SaveSalvage(ctx context.Context, salvage *cs.Salvage) error
}

// A database Client interface is used to make all calls that modify the database
type Client interface {
	ReadClient
	WriteClient
}

type dbConn struct {
	dbRead           *sql.DB
	dbWrite          *sql.DB
	databaseInMemory bool
}

type sqlReader interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type client struct {
	tx        *sql.Tx
	readConn  sqlReader
	reader    *gen.Queries
	writer    *gen.Queries
	converter Converter
}

func NewConn() DBConn {
	return &dbConn{}
}

func (conn *dbConn) NewReadClient() ReadClient {
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
func newTransactionClient(tx *sql.Tx) *client {
	return &client{
		tx:        tx,
		readConn:  tx,
		reader:    gen.New(tx),
		writer:    gen.New(tx),
		converter: c,
	}
}

// helper function to wrap a series of db calls in a transaction
func (conn *dbConn) WrapInTransaction(wrap func(c Client) error) error {
	tx, err := conn.dbWrite.Begin()
	if err != nil {
		return err
	}
	defer func() { tx.Rollback() }()

	if err := wrap(newTransactionClient(tx)); err != nil {
		return err
	}

	return tx.Commit()
}

func (c *dbConn) Connect(ctx context.Context, cfg *config.Config) error {

	c.databaseInMemory = strings.Contains(cfg.Database.Filename, ":memory:")
	// make sure the database is up to date
	c.mustMigrate(cfg)

	// create a new logger for logging database calls
	var logger *slog.Logger
	if cfg.Database.DebugLogging {
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
	}
	loggerAdapter := newLoggerWithLogger(logger)

	// dsnRead is like file::memory:?cache=shared, or file:data.db?_journal=WAL
	dsnRead := fmt.Sprintf("file:%s%s", cfg.Database.Filename, cfg.Database.ReadConnectionParams)
	dsnWrite := fmt.Sprintf("file:%s%s", cfg.Database.Filename, cfg.Database.WriteConnectionParams)
	slog.DebugContext(ctx, "Connecting to database", slog.String("dsnRead", dsnRead), slog.String("dsnWrite", dsnWrite))
	connectHook := func(conn *sqlite3.SQLiteConn) error {
		if c.databaseInMemory {
			// no need to attach
			return nil
		}

		// enforce foreign keys and WAL mode checkpointing
		if _, err := conn.Exec("PRAGMA foreign_keys = ON;", nil); err != nil {
			return err
		}
		if _, err := conn.Exec("PRAGMA journal_mode = WAL;", nil); err != nil {
			return err
		}
		if _, err := conn.Exec("PRAGMA synchronous = NORMAL;", nil); err != nil {
			return err
		}

		// cheap, non-blocking checkpoint attempt
		// (won't truncate; just nudges checkpointing)
		// ignore an error, we just want to attempt to checkpoint
		_, _ = conn.Exec("PRAGMA wal_checkpoint(PASSIVE);", nil)
		return nil
	}

	dbRead := sqldblogger.OpenDriver(dsnRead, &sqlite3.SQLiteDriver{ConnectHook: connectHook}, loggerAdapter)
	dbWrite := sqldblogger.OpenDriver(dsnWrite, &sqlite3.SQLiteDriver{ConnectHook: connectHook}, loggerAdapter)

	c.dbRead = dbRead
	if c.databaseInMemory {
		// no separate write connetion for in memory dbs
		c.dbWrite = c.dbRead
	} else {
		c.dbWrite = dbWrite
		c.dbWrite.SetMaxOpenConns(1)
	}

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

func (c *dbConn) Close() error {
	if err := c.dbRead.Close(); err != nil {
		return err
	}
	if err := c.dbWrite.Close(); err != nil {
		return err
	}
	return nil
}
