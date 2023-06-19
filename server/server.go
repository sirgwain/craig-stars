package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/sirgwain/craig-stars/config"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type server struct {
	db            DBClient
	config        config.Config
	gameRunner    GameRunner
	playerUpdater PlayerUpdater
}

func Start(db DBClient, config config.Config) {

	server := &server{
		db:            db,
		config:        config,
		gameRunner:    NewGameRunner(db),
		playerUpdater: newPlayerUpdater(db),
	}

	// create the data dir
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		panic(err)
	}

	sessionDB, err := gorm.Open(sqlite.Open("data/session.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	store := gormsessions.NewStore(sessionDB, true, []byte("secret"))

	r := gin.Default()
	r.Use(sessions.Sessions("session", store))

	// create a static folder for embedded assets
	embeddedFolder := EmbedFolder(GetAssets(), "frontend/build")
	staticServer := static.Serve("/", embeddedFolder)
	r.Use(staticServer)

	// weird workaround for embedded static assets (from gin-contrib/static github issue)
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet &&
			!strings.ContainsRune(c.Request.URL.Path, '.') &&
			!strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Request.URL.Path = "/"
			staticServer(c)
		}
	})

	r.POST("/api/login", server.login)
	r.GET("/api/logout", server.logout)

	// techs are public
	r.GET("/api/techs", server.techs)
	r.GET("/api/techs/:name", server.tech)

	//
	// authorized routes
	//
	ar := r.Group("/api")
	ar.Use(server.authRequired)

	// user stuff
	ar.GET("/me", server.me)
	ar.GET("/users", server.users)

	ar.GET("/rules", server.rules)

	// race CRUD
	ar.GET("/races", server.races)
	ar.GET("/races/:id", server.race)
	ar.PUT("/races/:id", server.updateRace)
	ar.POST("/races", server.createRace)
	ar.POST("/races/points", server.getRacePoints)

	// get various lists of games
	ar.GET("/games", server.playerGames)
	ar.GET("/games/hosted", server.hostedGames)
	ar.GET("/games/open", server.openGames)
	ar.GET("/games/open/:id", server.openGame)

	// host, join, generate, delete games
	ar.POST("/games", server.hostGame)
	ar.POST("/games/open/:id", server.joinGame)
	ar.POST("/games/:id/generate", server.generateUniverse)
	ar.DELETE("/games/:id", server.deleteGame)

	// player load/submit turn
	ar.GET("/games/:id", server.game)
	ar.GET("/games/:id/player", server.lightPlayer)
	ar.GET("/games/:id/full-player", server.fullPlayer)
	ar.GET("/games/:id/mapobjects", server.mapObjects)
	ar.GET("/games/:id/player-statuses", server.playerStatuses)
	ar.POST("/games/:id/submit-turn", server.submitTurn)

	// update planet production
	ar.PUT("/planets/:id", server.UpdatePlanetOrders)

	// transfer cargo, update fleet orders
	ar.PUT("/fleets/:id", server.UpdateFleetOrders)
	ar.POST("/fleets/:id/transfer-cargo", server.transferCargo)
	ar.POST("/fleets/:id/split", noop)
	ar.POST("/fleets/:id/merge", noop)
	ar.POST("/fleets/:id/rename", noop)

	// CRUD for ship designs
	ar.GET("/games/:id/designs", server.getShipDesigns)
	ar.GET("/games/:id/designs/:designnum", server.getShipDesign)
	ar.DELETE("/games/:id/designs/:designnum", server.deleteShipDesign)
	ar.POST("/games/:id/designs", server.createShipDesign)
	ar.POST("/games/:id/designs/spec", server.computeShipDesignSpec)
	ar.PUT("/games/:id/designs", server.updateShipDesign)

	// Update player plans
	ar.PUT("/games/:id/battle-plans", noop)
	ar.PUT("/games/:id/transport-plans", noop)
	ar.PUT("/games/:id/production-plans", noop)

	// update player reserarch, settings
	ar.PUT("/games/:id", server.updatePlayerOrders)

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

// noop function to create test api handlers
func noop(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "noop",
	})
}
