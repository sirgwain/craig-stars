package server

import (
	"crypto/sha1"
	"io/fs"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/config"
	"golang.org/x/oauth2"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/logger"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
)

const databaseIDAttr = "DatabaseID"

type contextKey int

const (
	keyUser contextKey = iota
	keyRace
	keyGame
	keyPlayer
	keyShipDesign
	keyBattlePlan
	keyProductionPlan
	keyTransportPlan
	keyPlanet
	keyFleet
	keyMineField
)

type server struct {
	db         DBClient
	config     config.Config
	gameRunner GameRunner
}

func Start(db DBClient, config config.Config) {

	// static resources
	sub, err := fs.Sub(assets, "frontend/build")
	if err != nil {
		panic(err)
	}

	// create a server
	server := &server{
		db:         db,
		config:     config,
		gameRunner: NewGameRunner(db, config),
	}
	_ = server

	var authLogger = logger.Func(func(format string, args ...interface{}) { log.Info().Msgf(format, args...) })

	cookieDuration := time.Hour * 24
	duration, err := time.ParseDuration(config.Discord.CookieDuration)
	if err != nil {
		log.Error().Err(err).Msgf("failed to load cookie duration from config %s", config.Discord.CookieDuration)
	} else {
		cookieDuration = duration
	}
	options := auth.Opts{
		SecretReader: token.SecretFunc(func(_ string) (string, error) { // secret key for JWT, ignores aud
			return server.config.Auth.Secret, nil
		}),
		TokenDuration:     time.Minute,                           // short token, refreshed automatically
		CookieDuration:    cookieDuration,                        // cookie fine to keep for long time
		DisableXSRF:       config.Auth.DisableXSRF,               // don't disable XSRF in real-life applications!
		Issuer:            "craig-stars",                         // part of token, just informational
		URL:               server.config.Auth.URL,                // base url of the protected service
		AvatarStore:       avatar.NewLocalFS("/tmp/craig-stars"), // stores avatars locally
		AvatarResizeLimit: 200,                                   // resizes avatars to 200x200
		ClaimsUpd: token.ClaimsUpdFunc(func(claims token.Claims) token.Claims { // modify issued token
			if claims.User != nil {
				user, err := server.db.GetUserByUsername(claims.User.Name)
				if err != nil {
					log.Error().Err(err).Msgf("failed to load %s from database during claims update", claims.User.Name)
					// TODO: add a rejection or something?
					return claims
				}
				// create a new user for this oauth user
				if user == nil {
					user, err = server.createNewUser(claims.User)
					if err != nil {
						return claims
					}
				} else {
					// if we're admin, set the admin claim
					if claims.User != nil && claims.User.Name == "admin" { // set attributes for admin
						claims.User.SetAdmin(true)
					} else if err := server.updateUser(claims.User, user); err != nil {
						return claims
					}
				}

				claims.User.SetStrAttr(databaseIDAttr, strconv.FormatInt(user.ID, 10))
			}

			return claims
		}),
		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool { // rejects some tokens

			if claims.User != nil {
				// TODO: if the user is blocked, reject them here
				return true
			}
			return false
		}),
		Logger:      authLogger, // optional logger for auth library
		UseGravatar: true,       // for verified provider use gravatar service
	}

	// create auth service
	service := auth.NewService(options)

	service.AddDirectProvider("local", provider.CredCheckerFunc(func(username, password string) (ok bool, err error) {

		user, err := server.db.GetUserByUsername(username)
		if err != nil {
			log.Error().Err(err).Str("Username", username).Msg("get user from database")
			return false, err
		}

		if user == nil {
			log.Error().Str("Username", username).Msg("user not found")
			return false, nil
		}

		// Check for username and password match
		return user.ComparePassword(password)
	}))

	if server.config.Discord.Enabled {
		c := auth.Client{
			Cid:     server.config.Discord.ClientID,
			Csecret: server.config.Discord.ClientSecret,
		}

		service.AddCustomProvider("discord", c, provider.CustomHandlerOpt{
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/api/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
			InfoURL: "https://discord.com/api/users/@me",
			MapUserFn: func(data provider.UserData, _ []byte) token.User {
				username := data.Value("username")
				id := data.Value("id")
				avatar := data.Value("avatar")
				userInfo := token.User{
					ID:   "discord_" + token.HashID(sha1.New(), username),
					Name: username,
				}
				userInfo.SetStrAttr("discord_id", id)
				userInfo.SetStrAttr("discord_avatar", avatar)
				return userInfo
			},
			Scopes: []string{"identify"},
		})
	}

	// retrieve auth middleware
	m := service.Middleware()

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(requestLogger(&log.Logger))
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// techs are public
	r.Route("/api/techs", func(r chi.Router) {
		r.Get("/", server.techs)
		r.Get("/{name:[a-zA-Z0-9-\\s]+}", server.tech)
	})
	r.Route("/api/battles", func(r chi.Router) {
		r.Get("/test", server.testBattle)
	})

	r.Group(func(r chi.Router) {
		r.Use(m.Auth)
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Use(server.userCtx)
		r.Get("/api/me", me)

		// race CRUD
		r.Route("/api/races", func(r chi.Router) {
			r.Post("/", server.createRace)
			r.Get("/", server.races)
			r.Post("/points", server.getRacePoints)

			// race by id operations
			r.Route("/{id:[0-9]+}", func(r chi.Router) {
				r.Use(server.raceCtx)
				r.Get("/", server.race)
				r.Put("/", server.updateRace)
				r.Delete("/", server.deleteRace)
			})
		})

		r.Route("/api/calculators", func(r chi.Router) {
			r.Post("/planet-production-estimate", server.getPlanetProductionEstimate)
		})

		// route for all operations that act on a game
		r.Route("/api/games", func(r chi.Router) {
			r.Post("/", server.createGame)
			r.Get("/", server.games)
			r.Get("/hosted", server.hostedGames)
			r.Get("/open", server.openGames)
			r.Get("/invite/{hash:[a-zA-Z0-9]+}", server.openGamesByHash)

			// game by id operations
			r.Route("/{id:[0-9]+}", func(r chi.Router) {
				r.Use(server.gameCtx)
				r.Get("/ping-discord", server.pingDiscordForGameUpdate)
				r.Get("/", server.game)
				r.Put("/", server.updateGame)
				r.Post("/join", server.joinGame)
				r.Post("/leave", server.leaveGame)
				r.Post("/add-ai", server.addOpenPlayerSlot)
				r.Post("/add-player", server.addOpenPlayerSlot)
				r.Post("/add-ai-player", server.addAIPlayer)
				r.Post("/kick-player", server.kickPlayer)
				r.Post("/delete-player", server.deletePlayerSlot)
				r.Post("/update-player", server.updatePlayerSlot)
				r.Post("/start-game", server.startGame)
				r.Post("/generate-turn", server.generateTurn)
				r.Get("/compute-specs", server.computeSpecs)
				r.Delete("/", server.deleteGame)

				// routes requiring a player and game
				r.Group(func(r chi.Router) {
					r.Use(server.playerCtx)
					r.Get("/player", server.player)
					r.Get("/player/intels", server.playerIntels)
					r.Put("/player", server.updatePlayerOrders)
					r.Put("/player/plans", server.updatePlayerPlans)
					r.Post("/submit-turn", server.submitTurn)
					r.Post("/unsubmit-turn", server.unSubmitTurn)

					// ship designs
					r.Route("/designs", func(r chi.Router) {
						r.Get("/", server.shipDesigns)
						r.Post("/", server.createShipDesign)
						r.Post("/spec", server.computeShipDesignSpec)

						// shipdesign by num operations
						r.Route("/{num:[0-9]+}", func(r chi.Router) {
							r.Use(server.shipdDesignCtx)
							r.Get("/", server.shipDesign)
							r.Put("/", server.updateShipDesign)
							r.Delete("/", server.deleteShipDesign)
						})
					})

					// battle plans
					r.Route("/battle-plans", func(r chi.Router) {
						r.Post("/", server.createBattlePlan)

						// shipdesign by num operations
						r.Route("/{num:[0-9]+}", func(r chi.Router) {
							r.Use(server.battlePlanCtx)
							r.Put("/", server.updateBattlePlan)
							r.Delete("/", server.deleteBattlePlan)
						})
					})

					// production plans
					r.Route("/production-plans", func(r chi.Router) {
						r.Post("/", server.createProductionPlan)

						// shipdesign by num operations
						r.Route("/{num:[0-9]+}", func(r chi.Router) {
							r.Use(server.productionPlanCtx)
							r.Put("/", server.updateProductionPlan)
							r.Delete("/", server.deleteProductionPlan)
						})
					})

					// transport plans
					r.Route("/transport-plans", func(r chi.Router) {
						r.Post("/", server.createTransportPlan)

						// shipdesign by num operations
						r.Route("/{num:[0-9]+}", func(r chi.Router) {
							r.Use(server.transportPlanCtx)
							r.Put("/", server.updateTransportPlan)
							r.Delete("/", server.deleteTransportPlan)
						})
					})

					// planet order updates
					r.Route("/planets", func(r chi.Router) {
						r.Route("/{num:[0-9]+}", func(r chi.Router) {
							r.Use(server.planetCtx)
							r.Get("/", server.planet)
							r.Put("/", server.updatePlanetOrders)
						})
					})

					// planet order updates
					r.Route("/fleets", func(r chi.Router) {
						r.Route("/{num:[0-9]+}", func(r chi.Router) {
							r.Use(server.fleetCtx)
							r.Get("/", server.fleet)
							r.Put("/", server.updateFleetOrders)
							r.Post("/split-all", server.splitAll)
							r.Post("/merge", server.merge)
							r.Post("/transfer-cargo", server.transferCargo)
							r.Post("/rename", server.renameFleet)
						})
					})

				})

				r.Get("/full-player", server.fullPlayer)
				r.Get("/mapobjects", server.mapObjects)
				r.Get("/universe", server.universe)

			})
		})

	})

	// setup auth routes
	authRoutes, avaRoutes := service.Handlers()
	r.Mount("/api/auth", authRoutes)  // add auth handlers
	r.Mount("/api/avatar", avaRoutes) // add avatar handler

	r.Handle("/*", http.FileServer(http.FS(sub)))

	http.ListenAndServe(":8080", r)
}

// create a new request logger with zerolog. Inspired by https://github.com/ironstar-io/chizerolog
func requestLogger(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.With().Logger()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				t2 := time.Now()

				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					log.Error().
						Timestamp().
						Interface("info", rec).
						Bytes("stack", debug.Stack()).
						Msg("system error")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				// log end request
				var event *zerolog.Event
				if ww.Status() >= 400 {
					event = log.Error()
				} else {
					event = log.Info()
				}

				fields := map[string]interface{}{
					"ip":             r.RemoteAddr,
					"url":            r.URL.Path,
					"method":         r.Method,
					"status":         ww.Status(),
					"ms":             float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
					"content-length": r.Header.Get("Content-Length"),
					"resp_bytes":     ww.BytesWritten(),
				}

				// don't log the user_agent while we're debugging, we should know what it is
				if zerolog.GlobalLevel() != zerolog.DebugLevel {
					fields["user_agent"] = r.Header.Get("User-Agent")
				}

				event.
					Timestamp().
					Fields(fields).Msg("")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func (s *server) int64URLParam(r *http.Request, key string) (*int64, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return nil, nil
	}
	var num int64
	num, err := strconv.ParseInt(param, 10, 64)
	if err != nil {

		return nil, err
	}

	return &num, nil
}

func (s *server) intURLParam(r *http.Request, key string) (*int, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return nil, nil
	}
	var num int
	num, err := strconv.Atoi(param)
	if err != nil {

		return nil, err
	}

	return &num, nil
}
