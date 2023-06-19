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

	_ "github.com/mattn/go-sqlite3"
)

type client struct {
	sqlDB *gorm.DB
}

type Client interface {
	Connect(config *config.Config)
	EnableDebugLogging()

	GetUsers() ([]game.User, error)
	GetUser(id int64) (*game.User, error)
	GetUserByUsername(username string) (*game.User, error)
	CreateUser(user *game.User) error
	UpdateUser(user *game.User) error
	DeleteUser(id int64) error

	GetTechStores() ([]*game.TechStore, error)
	CreateTechStore(tech *game.TechStore) error
	FindTechStoreById(id int64) (*game.TechStore, error)

	GetGames() ([]*game.Game, error)
	GetGamesForHost(userID int64) ([]*game.Game, error)
	GetGamesByUser(userID int64) ([]*game.Game, error)
	GetOpenGames() ([]*game.Game, error)
	GetGame(id int64) (*game.Game, error)
	GetFullGame(id int64) (*game.FullGame, error)
	FindGameRulesByGameID(gameID int64) (*game.Rules, error)
	CreateGame(game *game.Game) error
	UpdateFullGame(game *game.FullGame) error
	DeleteGame(id int64) error

	GetRaces(userID int64) ([]*game.Race, error)
	GetRace(id int64) (*game.Race, error)
	CreateRace(race *game.Race) error
	UpdateRace(race *game.Race) error

	FindPlayerByGameId(gameID int64, userID int64) (*game.FullPlayer, error)
	FindPlayerByGameIdLight(gameID int64, userID int64) (*game.Player, error)
	SavePlayer(player *game.Player) error

	GetPlanet(id int64) (*game.Planet, error)
	UpdatePlanet(planet *game.Planet) error

	GetFleet(id int64) (*game.Fleet, error)
	UpdateFleet(fleet *game.Fleet) error
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func NewClient() Client {
	return &client{}
}

func (c *client) Connect(config *config.Config) {
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

	c.sqlDB = localdb

	if config.Database.Filename == ":memory:" || config.Database.Recreate {
		log.Debug().Msgf("Creating in memory database")
		c.migrateAll()
	}

}

func (c *client) EnableDebugLogging() {
	zlogger := zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	dblogger := NewWithLogger(&zlogger)
	c.sqlDB.Logger = dblogger
}
