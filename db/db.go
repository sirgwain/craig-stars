package db

import (
	"os"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/game"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	sqlDB *gorm.DB
}

type Client interface {
	Connect(config *config.Config)
	EnableDebugLogging()

	GetUsers() *[]game.User
	SaveUser(user *game.User) error
	FindUserById(id uint) (*game.User, error)
	FindUserByUsername(username string) (*game.User, error)
	DeleteUserById(id uint)

	Migrate(item interface{})
	MigrateAll() error

	GetTechStores() (*game.TechStore, error)
	CreateTechStore(tech *game.TechStore) error
	FindTechStoreById(id uint) (*game.TechStore, error)

	GetGames() ([]game.Game, error)
	GetGamesHostedByUser(userID uint) ([]game.Game, error)
	GetGamesByUser(userID uint) ([]game.Game, error)
	GetOpenGames() ([]game.Game, error)
	FindGameById(id uint) (*game.FullGame, error)
	FindGameByIdLight(id uint) (*game.Game, error)
	FindGameRulesByGameId(gameId uint) (*game.Rules, error)
	CreateGame(game *game.Game) error
	SaveGame(game *game.FullGame) error
	DeleteGameById(id uint) error

	GetRaces(userID uint) ([]game.Race, error)
	FindRaceById(id uint) (*game.Race, error)
	CreateRace(race *game.Race) error
	SaveRace(race *game.Race) error

	FindPlayerByGameId(gameID uint, userID uint) (*game.FullPlayer, error)
	FindPlayerByGameIdLight(gameID uint, userID uint) (*game.Player, error)
	SavePlayer(player *game.Player) error

	FindPlanetById(id uint) (*game.Planet, error)
	SavePlanet(planet *game.Planet) error

	FindFleetById(id uint) (*game.Fleet, error)
	SaveFleet(fleet *game.Fleet) error
}

func (db *DB) Connect(config *config.Config) {

	log.Debug().Msgf("Connecting to database %s", config.Database.Filename)
	zlogger := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.WarnLevel)
	dblogger := NewWithLogger(&zlogger)

	localdb, err := gorm.Open(sqlite.Open(config.Database.Filename), &gorm.Config{
		Logger: dblogger,
	})
	if err != nil {
		panic(err)
	}

	db.sqlDB = localdb

	if config.Database.Filename == ":memory:" {
		log.Debug().Msgf("Creating in memory database")
		db.MigrateAll()
	}
}

func (db *DB) EnableDebugLogging() {
	zlogger := zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	dblogger := NewWithLogger(&zlogger)
	db.sqlDB.Logger = dblogger
}
