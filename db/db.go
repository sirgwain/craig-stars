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

type Service interface {
	Connect(config *config.Config)
	EnableDebugLogging()

	GetUsers() *[]game.User
	SaveUser(user *game.User) error
	FindUserById(id uint) (*game.User, error)
	FindUserByUsername(username string) *game.User
	DeleteUserById(id uint)

	Migrate(item interface{})
	MigrateAll() error

	GetTechStores() (*game.TechStore, error)
	CreateTechStore(tech *game.TechStore) error
	FindTechStoreById(id uint) (*game.TechStore, error)

	GetGames() []game.Game
	CreateGame(game *game.Game) error
	SaveGame(game *game.Game) error
	FindGameById(id uint) (*game.Game, error)
	FindGameByIdLight(id uint) (*game.Game, error)
	GetGamesByUser(userID uint) []game.Game
	DeleteGameById(id uint)

	FindPlayerByGameId(gameID uint, userID uint) (*game.Player, error)
	FindPlayerByGameIdLight(gameID uint, userID uint) (*game.Player, error)
	SavePlayer(player *game.Player) error

	FindPlanetById(id uint) (*game.Planet, error)
	SavePlanet(planet *game.Planet) error
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
}

func (db *DB) EnableDebugLogging() {
	zlogger := zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	dblogger := NewWithLogger(&zlogger)
	db.sqlDB.Logger = dblogger
}
