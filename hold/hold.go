package hold

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/game"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

type DB struct {
	store *bolthold.Store
}

type Client interface {
	Connect(config *config.Config)
	EnableDebugLogging()
	Close()

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

	FindPlanetByNum(gameID uint64, num int) (*game.Planet, error)
	SavePlanet(gameID uint64, planet *game.Planet) error

	FindFleetByNum(gameID uint64, playerNum int, num int) (*game.Fleet, error)
	SaveFleet(gameID uint64, fleet *game.Fleet) error
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// each game has a bucket with game objects in it
func (db *DB) getGameBucketName(gameID uint64) string {
	return fmt.Sprintf("Game-%d", gameID)
}

func (db *DB) getFleetKey(playerNum int, num int) string {
	return fmt.Sprintf("%d-%d", playerNum, num)
}

func (db *DB) Connect(config *config.Config) {
	if config.Database.Recreate {
		info, _ := os.Stat(config.Database.Filename)
		if info != nil {
			log.Debug().Msgf("Deleting existing database %s", config.Database.Filename)
			os.Remove(config.Database.Filename)
		}
	}
	log.Debug().Msgf("Connecting to database %s", config.Database.Filename)
	store, err := bolthold.Open(config.Database.Filename, 0666, nil)
	if err != nil {
		// handle error
		log.Fatal().Err(err).Str("Filename", config.Database.Filename).Msg("Failed to open bolthold store")
	}
	db.store = store
}

func (db *DB) MigrateAll() error {
	// no op
	return nil
}

func (db *DB) EnableDebugLogging() {
	// todo: figure it out
}

func (db *DB) Close() {
	db.store.Close()
}

func (db *DB) GetUsers() ([]*game.User, error) {
	users := []*game.User{}
	// load Players
	if err := db.store.Find(&users, nil); err != nil {
		return nil, err
	}

	return users, nil
}

func (db *DB) SaveUser(user *game.User) error {
	if user.ID != 0 {
		return db.store.Upsert(user.ID, user)
	} else {
		return db.store.Insert(bolthold.NextSequence(), user)
	}
}

func (db *DB) FindUserById(id uint64) (*game.User, error) {
	user := game.User{}
	if err := db.store.Get(id, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) FindUserByUsername(username string) (*game.User, error) {
	user := game.User{}
	if err := db.store.FindOne(&user, bolthold.Where("Username").Eq(username)); err != nil {
		if err == bolthold.ErrNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &user, nil
}

func (db *DB) DeleteUserById(id uint64) error {
	return db.store.Delete(id, &game.User{})
}

func (db *DB) GetTechStores() ([]*game.TechStore, error) {
	techstores := []*game.TechStore{}
	// load Players
	if err := db.store.Find(&techstores, nil); err != nil {
		return nil, err
	}

	return techstores, nil

}

func (db *DB) CreateTechStore(tech *game.TechStore) error {
	key := bolthold.NextSequence()
	if tech.ID != 0 {
		key = tech.ID
	}

	return db.store.Upsert(key, tech)
}

func (db *DB) FindTechStoreById(id uint64) (*game.TechStore, error) {
	techstore := game.TechStore{}
	if err := db.store.Get(id, &techstore); err != nil {
		return nil, err
	}
	return &techstore, nil
}

func (db *DB) GetGames() ([]*game.Game, error) {
	games := []*game.Game{}
	if err := db.store.Find(&games, nil); err != nil {
		return nil, err
	}

	return games, nil
}

func (db *DB) GetGamesHostedByUser(userID uint64) ([]*game.Game, error) {
	games := []*game.Game{}
	if err := db.store.Find(&games, bolthold.Where("HostID").Eq(userID)); err != nil {
		if err == bolthold.ErrNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return games, nil
}

func (db *DB) GetGamesByUser(userID uint64) ([]*game.Game, error) {
	games := []*game.Game{}
	userGames := []*game.Game{}
	if err := db.store.Find(&games, bolthold.Where("HostID").Eq(userID)); err != nil {
		if err == bolthold.ErrNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	for _, g := range games {
		if err := db.store.Bolt().View(func(tx *bbolt.Tx) error {
			gameBucketName := db.getGameBucketName(g.ID)
			gameBucket := tx.Bucket([]byte(gameBucketName))
			if gameBucket == nil {
				return fmt.Errorf("game bucket: %s does not exist", gameBucketName)
			}

			players := []*game.Player{}
			if err := db.store.FindInBucket(gameBucket, &players, bolthold.Where("UserID").Eq(userID)); err != nil {
				if err == bolthold.ErrNotFound {
					return nil
				}
				return err
			}
			if len(players) > 0 {
				userGames = append(userGames, g)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return userGames, nil
}

func (db *DB) GetOpenGames() ([]*game.Game, error) {
	games := []*game.Game{}
	if err := db.store.Find(&games, bolthold.Where("OpenPlayerSlots").Gt(0)); err != nil {
		if err == bolthold.ErrNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return games, nil
}

func (db *DB) FindGameById(id uint64) (*game.FullGame, error) {
	g := game.Game{}

	if err := db.store.Get(id, &g); err != nil {
		if err == bolthold.ErrNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	universe := game.Universe{}
	players := []*game.Player{}

	if err := db.store.Bolt().View(func(tx *bbolt.Tx) error {
		gameBucketName := db.getGameBucketName(g.ID)
		gameBucket := tx.Bucket([]byte(gameBucketName))
		if gameBucket == nil {
			return fmt.Errorf("game bucket: %s does not exist", gameBucketName)
		}

		// load Planets
		if err := db.store.FindInBucket(gameBucket, &universe.Planets, nil); err != nil {
			return err
		}
		// load Fleets
		if err := db.store.FindInBucket(gameBucket, &universe.Fleets, nil); err != nil {
			return err
		}
		// load Wormholes
		if err := db.store.FindInBucket(gameBucket, &universe.Wormholes, nil); err != nil {
			return err
		}
		// load MineralPackets
		if err := db.store.FindInBucket(gameBucket, &universe.MineralPackets, nil); err != nil {
			return err
		}
		// load MineFields
		if err := db.store.FindInBucket(gameBucket, &universe.MineFields, nil); err != nil {
			return err
		}
		// load Salvage
		if err := db.store.FindInBucket(gameBucket, &universe.Salvage, nil); err != nil {
			return err
		}

		// load Players
		if err := db.store.FindInBucket(gameBucket, &players, nil); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	g.Rules.ResetSeed()
	g.Rules.WithTechStore(&game.StaticTechStore)

	return &game.FullGame{
		Game:     &g,
		Universe: &universe,
		Players:  players,
	}, nil
}

func (db *DB) FindGameByIdLight(id uint64) (*game.Game, error) {
	g := game.Game{}
	if err := db.store.Get(id, &g); err != nil {
		if err == bolthold.ErrNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &g, nil
}

func (db *DB) FindGameRulesByGameID(gameID uint64) (*game.Rules, error) {
	g, err := db.FindGameByIdLight(gameID)
	if err != nil {
		return nil, err
	}

	if g != nil {
		return &g.Rules, nil
	}
	return nil, nil
}

func (db *DB) CreateGame(g *game.Game) error {
	err := db.store.Insert(bolthold.NextSequence(), g)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) SaveGame(g *game.FullGame) error {
	defer timeTrack(time.Now(), "SaveGame")

	if err := db.store.Upsert(g.ID, g.Game); err != nil {
		return err
	}

	return db.SaveGameBits(g)

	// return nil
}

// DefaultEncode is the default encoding func for bolthold (Gob)
func DefaultEncode(value interface{}) ([]byte, error) {
	var buff bytes.Buffer

	en := gob.NewEncoder(&buff)

	err := en.Encode(value)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// DefaultDecode is the default decoding func for bolthold (Gob)
func DefaultDecode(data []byte, value interface{}) error {
	var buff bytes.Buffer
	de := gob.NewDecoder(&buff)

	_, err := buff.Write(data)
	if err != nil {
		return err
	}

	return de.Decode(value)
}

func (db *DB) SaveGameBits(g *game.FullGame) error {
	defer timeTrack(time.Now(), "SaveGameBits")
	if err := db.store.Bolt().Batch(func(tx *bbolt.Tx) error {
		gameBucketName := db.getGameBucketName(g.ID)
		gameBucket, err := tx.CreateBucketIfNotExists([]byte(gameBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for _, planet := range g.Planets {
			if planet.Dirty {
				planet.Dirty = false
				if err := db.store.UpsertBucket(gameBucket, planet.Num, planet); err != nil {
					return err
				}
			}
		}
		for _, fleet := range g.Fleets {
			if fleet.Dirty {
				fleet.Dirty = false
				if err := db.store.UpsertBucket(gameBucket, db.getFleetKey(fleet.PlayerNum, fleet.Num), fleet); err != nil {
					return err
				}

			}
		}
		for _, wormhole := range g.Wormholes {
			if wormhole.Dirty {
				wormhole.Dirty = false
				if err := db.store.UpsertBucket(gameBucket, wormhole.Num, wormhole); err != nil {
					return err
				}

			}
		}
		for _, mineralPacket := range g.MineralPackets {
			if mineralPacket.Dirty {
				mineralPacket.Dirty = false
				if err := db.store.UpsertBucket(gameBucket, mineralPacket.Num, mineralPacket); err != nil {
					return err
				}

			}
		}
		for _, mineField := range g.MineFields {
			if mineField.Dirty {
				mineField.Dirty = false
				if err := db.store.UpsertBucket(gameBucket, mineField.Num, mineField); err != nil {
					return err
				}

			}
		}
		for _, salvage := range g.Salvage {
			if salvage.Dirty {
				salvage.Dirty = false
				if err := db.store.UpsertBucket(gameBucket, salvage.Num, salvage); err != nil {
					return err
				}

			}
		}

		// for _, player := range g.Players {
		// 	// player's need to know the game id in order to fetch/save
		// 	player.GameID = g.ID
		// 	if err := db.store.UpsertBucket(gameBucket, player.Num, player); err != nil {
		// 		return err
		// 	}
		// }
		return db.SavePlayers(g.ID, gameBucket, g.Players)

		// return nil
	}); err != nil {
		return err
	}
	return nil
}

func (db *DB) SavePlayers(gameID uint64, gameBucket *bbolt.Bucket, players []*game.Player) error {
	defer timeTrack(time.Now(), "SavePlayers")
	for _, player := range players {
		// player's need to know the game id in order to fetch/save
		player.GameID = gameID
		if err := db.store.UpsertBucket(gameBucket, player.Num, player); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) DeleteGameById(id uint64) error {
	// delete the game
	if err := db.store.Delete(id, &game.Game{}); err != nil {
		return err
	}

	if err := db.store.Bolt().Update(func(tx *bbolt.Tx) error {
		gameBucketName := db.getGameBucketName(id)
		return tx.DeleteBucket([]byte(gameBucketName))
	}); err != nil {
		return err
	}
	return nil
}

func (db *DB) GetRaces(userID uint64) ([]*game.Race, error) {
	races := []*game.Race{}
	if err := db.store.Find(&races, bolthold.Where("UserID").Eq(userID)); err != nil {
		if err == bolthold.ErrNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return races, nil

}

func (db *DB) FindRaceById(id uint64) (*game.Race, error) {
	race := game.Race{}
	if err := db.store.Get(id, &race); err != nil {
		return nil, err
	}
	return &race, nil
}

func (db *DB) SaveRace(race *game.Race) error {
	if race.ID != 0 {
		return db.store.Upsert(race.ID, race)
	} else {
		return db.store.Insert(bolthold.NextSequence(), race)
	}
}

func (db *DB) FindPlayerByGameId(gameID uint64, userID uint64) (*game.FullPlayer, error) {
	player := game.FullPlayer{}
	if err := db.store.Bolt().View(func(tx *bbolt.Tx) error {
		gameBucketName := db.getGameBucketName(gameID)
		gameBucket := tx.Bucket([]byte(gameBucketName))
		if gameBucket == nil {
			return fmt.Errorf("game bucket: %s does not exist", gameBucketName)
		}

		if err := db.store.FindOneInBucket(gameBucket, &player.Player, bolthold.Where("UserID").Eq(userID)); err != nil {
			return err
		}

		// load PlayerMapObjects
		if err := db.store.FindInBucket(gameBucket, &player.Planets, bolthold.Where("PlayerNum").Eq(player.Num)); err != nil {
			return err
		}
		if err := db.store.FindInBucket(gameBucket, &player.Fleets, bolthold.Where("PlayerNum").Eq(player.Num)); err != nil {
			return err
		}
		if err := db.store.FindInBucket(gameBucket, &player.MineFields, bolthold.Where("PlayerNum").Eq(player.Num)); err != nil {
			return err
		}
		if err := db.store.FindInBucket(gameBucket, &player.MineralPackets, bolthold.Where("PlayerNum").Eq(player.Num)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &player, nil
}

func (db *DB) FindPlayerByGameIdLight(gameID uint64, userID uint64) (*game.Player, error) {
	player := game.Player{}
	if err := db.store.Bolt().View(func(tx *bbolt.Tx) error {
		gameBucketName := db.getGameBucketName(gameID)
		gameBucket := tx.Bucket([]byte(gameBucketName))
		if gameBucket == nil {
			return fmt.Errorf("game bucket: %s does not exist", gameBucketName)
		}

		if err := db.store.FindOneInBucket(gameBucket, &player, bolthold.Where("UserID").Eq(userID)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &player, nil
}

// save a player into the game bucket
func (db *DB) SavePlayer(player *game.Player) error {
	return db.store.Bolt().Update(func(tx *bbolt.Tx) error {
		gameBucketName := db.getGameBucketName(player.GameID)
		gameBucket, err := tx.CreateBucketIfNotExists([]byte(gameBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		if err := db.store.UpdateBucket(gameBucket, player.Num, player); err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) FindPlanetByNum(gameID uint64, num int) (*game.Planet, error) {
	planet := game.Planet{}
	if err := db.store.Bolt().View(func(tx *bbolt.Tx) error {
		gameBucketName := db.getGameBucketName(gameID)
		gameBucket := tx.Bucket([]byte(gameBucketName))
		if gameBucket == nil {
			return fmt.Errorf("game bucket: %s does not exist", gameBucketName)
		}

		if err := db.store.GetFromBucket(gameBucket, num, &planet); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &planet, nil
}

func (db *DB) SavePlanet(gameID uint64, planet *game.Planet) error {
	return db.store.Bolt().Update(func(tx *bbolt.Tx) error {
		gameBucketName := db.getGameBucketName(gameID)
		gameBucket := tx.Bucket([]byte(gameBucketName))

		if err := db.store.UpdateBucket(gameBucket, planet.Num, planet); err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) FindFleetByNum(gameID uint64, playerNum int, num int) (*game.Fleet, error) {
	fleet := game.Fleet{}
	if err := db.store.Bolt().View(func(tx *bbolt.Tx) error {
		gameBucketName := db.getGameBucketName(gameID)
		gameBucket := tx.Bucket([]byte(gameBucketName))
		if gameBucket == nil {
			return fmt.Errorf("game bucket: %s does not exist", gameBucketName)
		}

		if err := db.store.GetFromBucket(gameBucket, db.getFleetKey(playerNum, num), &fleet); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &fleet, nil
}

func (db *DB) SaveFleet(gameID uint64, fleet *game.Fleet) error {
	return db.store.Bolt().Update(func(tx *bbolt.Tx) error {
		gameBucketName := db.getGameBucketName(gameID)
		gameBucket := tx.Bucket([]byte(gameBucketName))

		if err := db.store.UpdateBucket(gameBucket, db.getFleetKey(fleet.PlayerNum, fleet.Num), fleet); err != nil {
			return err
		}

		return nil
	})
}
