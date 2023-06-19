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

	GetUsers() ([]game.User, error)
	GetUser(id int64) (*game.User, error)
	GetUserByUsername(username string) (*game.User, error)
	CreateUser(user *game.User) error
	UpdateUser(user *game.User) error
	DeleteUser(id int64) error

	GetRacesForUser(userID int64) ([]game.Race, error)
	GetRace(id int64) (*game.Race, error)
	CreateRace(race *game.Race) error
	UpdateRace(race *game.Race) error

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
	UpdateFullGame(game *game.FullGame) error
	DeleteGame(id int64) error

	GetFullPlayerForGame(gameID int64, userID int64) (*game.FullPlayer, error)
	GetPlayerForGame(gameID int64, userID int64) (*game.Player, error)
	CreatePlayer(player *game.Player) error
	UpdatePlayer(player *game.Player) error

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

	// create a new logger for logging database calls
	var zlogger zerolog.Logger
	if config.Database.DebugLogging {
		zlogger = zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	} else {
		zlogger = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.WarnLevel)
	}

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
