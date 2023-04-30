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

type idBind struct {
	ID int64 `uri:"id"`
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

	//
	// authorized routes
	//
	ar := r.Group("/api")
	ar.Use(server.authRequired)

	// user stuff
	ar.GET("/me", server.me)
	ar.GET("/users", server.users)

	// transfer cargo, update fleet orders
	ar.PUT("/fleets/:id", server.UpdateFleetOrders)
	ar.POST("/fleets/:id/transfer-cargo", server.transferCargo)

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
