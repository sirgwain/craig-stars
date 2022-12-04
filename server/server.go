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
	db         DBClient
	config     config.Config
	gameRunner GameRunner
}

func Start(db DBClient, config config.Config) {

	server := &server{
		db:         db,
		config:     config,
		gameRunner: NewGameRunner(db),
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

	r.POST("/api/login", server.Login)
	r.GET("/api/logout", server.Logout)

	// techs are public
	r.GET("/api/techs", server.Techs)

	// authorized routes
	ar := r.Group("/api")
	ar.Use(server.AuthRequired)
	ar.GET("/me", server.Me)
	ar.GET("/users", server.Users)

	ar.GET("/rules", server.Rules)

	ar.GET("/races", server.Races)
	ar.GET("/races/:id", server.Race)
	ar.POST("/races", server.CreateRace)
	ar.PUT("/races/:id", server.UpdateRace)

	// get various lists of games
	ar.GET("/games", server.PlayerGames)
	ar.GET("/games/hosted", server.HostedGames)
	ar.GET("/games/open", server.OpenGames)
	ar.GET("/games/open/:id", server.OpenGame)

	// host, join, delete games
	ar.POST("/games", server.HostGame)
	ar.POST("/games/open/:id", server.JoinGame)
	ar.DELETE("/games/:id", server.DeleteGame)

	// player load/submit turn
	ar.GET("/games/:id", server.PlayerGame)
	ar.POST("/games/:id/submit-turn", server.SubmitTurn)
	// update player reserarch, settings
	ar.PUT("/games/:id", noop)

	// update planet production
	ar.PUT("/planets/:id", server.UpdatePlanetOrders)

	// transfer cargo, update fleet orders
	ar.POST("/fleets/:id/transfer-cargo", server.TransferCargo)
	ar.PUT("/fleets/:id", server.UpdateFleetOrders)
	ar.PUT("/fleets/:id/split", noop)
	ar.PUT("/fleets/:id/merge", noop)
	ar.PUT("/fleets/:id/rename", noop)

	// CRUD for ship designs
	ar.GET("/games/:id/designs", noop)
	ar.GET("/games/:id/designs/:designid", noop)
	ar.POST("/games/:id/designs", noop)
	ar.PUT("/games/:id/designs/:designid", noop)

	// CRUD for battle plans
	ar.GET("/games/:id/battle-plans", noop)
	ar.GET("/games/:id/battle-plans/:name", noop)
	ar.POST("/games/:id/battle-plans", noop)
	ar.PUT("/games/:id/battle-plans/:name", noop)

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
