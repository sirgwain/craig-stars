package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/sirgwain/craig-stars/appcontext"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type server struct {
	ctx        *appcontext.AppContext
	gameRunner *GameRunner
}

func Start(ctx *appcontext.AppContext) {

	server := &server{
		ctx:        ctx,
		gameRunner: NewGameRunner(ctx.DB),
	}

	// create the data dir
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("data/session.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	store := gormsessions.NewStore(db, true, []byte("secret"))

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

	ar.GET("/rules", server.Rules)

	ar.GET("/races", server.Races)
	ar.GET("/races/:id", server.Race)
	ar.POST("/races", server.CreateRace)
	ar.PUT("/races/:id", server.UpdateRace)

	ar.GET("/games", server.PlayerGames)
	ar.GET("/games/hosted", server.HostedGames)
	ar.GET("/games/open", server.OpenGames)
	ar.GET("/games/open/:id", server.OpenGame)
	ar.POST("/games", server.HostGame)
	ar.POST("/games/open/:id", server.JoinGame)
	ar.GET("/games/:id", server.PlayerGame)
	ar.DELETE("/games/:id", server.DeleteGameById)
	ar.POST("/games/:id/submit-turn", server.SubmitTurn)
	ar.PUT("/planets/:id", server.UpdatePlanetOrders)

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
