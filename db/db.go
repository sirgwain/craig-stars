package db

import (
	"os"
	"time"

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
	MigrateAll() error

	Connect(config *config.Config)
	EnableDebugLogging()

	GetUsers() ([]*game.User, error)
	SaveUser(user *game.User) error
	FindUserById(id uint64) (*game.User, error)
	FindUserByUsername(username string) (*game.User, error)
	DeleteUserById(id uint64) error

	GetTechStores() ([]*game.TechStore, error)
	CreateTechStore(tech *game.TechStore) error
	FindTechStoreById(id uint64) (*game.TechStore, error)

	GetGames() ([]*game.Game, error)
	GetGamesHostedByUser(userID uint64) ([]*game.Game, error)
	GetGamesByUser(userID uint64) ([]*game.Game, error)
	GetOpenGames() ([]*game.Game, error)
	FindGameById(id uint64) (*game.FullGame, error)
	FindGameByIdLight(id uint64) (*game.Game, error)
	FindGameRulesByGameID(gameID uint64) (*game.Rules, error)
	CreateGame(game *game.Game) error
	SaveGame(game *game.FullGame) error
	DeleteGameById(id uint64) error

	GetRaces(userID uint64) ([]*game.Race, error)
	FindRaceById(id uint64) (*game.Race, error)
	SaveRace(race *game.Race) error

	FindPlayerByGameId(gameID uint64, userID uint64) (*game.FullPlayer, error)
	FindPlayerByGameIdLight(gameID uint64, userID uint64) (*game.Player, error)
	SavePlayer(player *game.Player) error

	FindPlanetByID(id uint64) (*game.Planet, error)
	FindPlanetByNum(gameID uint64, num int) (*game.Planet, error)
	SavePlanet(gameID uint64, planet *game.Planet) error

	FindFleetByID(id uint64) (*game.Fleet, error)
	FindFleetByNum(gameID uint64, playerNum int, num int) (*game.Fleet, error)
	SaveFleet(gameID uint64, fleet *game.Fleet) error
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func (db *DB) Connect(config *config.Config) {
	if config.Database.Recreate && config.Database.Filename != ":memory:" {
		info, _ := os.Stat(config.Database.Filename)
		if info != nil {
			log.Debug().Msgf("Deleting existing database %s", config.Database.Filename)
			os.Remove(config.Database.Filename)
		}
	}

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
