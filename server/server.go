package server

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"golang.org/x/oauth2"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/logger"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
)

type contextKey int

const (
	keyDb contextKey = iota
	keyUser
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
	db     DBConnection
	config config.Config
}

const userRejected = "rejected"

func Start(config config.Config) error {

	dbConn := db.NewConn()
	if err := dbConn.Connect(&config); err != nil {
		return fmt.Errorf("failed to connect to database %v", err)
	}

	// static resources
	sub, err := fs.Sub(assets, "frontend/build")
	if err != nil {
		panic(err)
	}

	// create a server
	server := &server{
		db:     dbConn,
		config: config,
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
	issuer := "craig-stars"
	options := auth.Opts{
		SecretReader: token.SecretFunc(func(_ string) (string, error) { // secret key for JWT, ignores aud
			return server.config.Auth.Secret, nil
		}),
		TokenDuration:     time.Minute,                           // short token, refreshed automatically
		CookieDuration:    cookieDuration,                        // cookie fine to keep for long time
		DisableXSRF:       config.Auth.DisableXSRF,               // don't disable XSRF in real-life applications!
		Issuer:            issuer,                                // part of token, just informational
		URL:               server.config.Auth.URL,                // base url of the protected service
		AvatarStore:       avatar.NewLocalFS("/tmp/craig-stars"), // stores avatars locally
		AvatarResizeLimit: 200,                                   // resizes avatars to 200x200
		ClaimsUpd: token.ClaimsUpdFunc(func(claims token.Claims) token.Claims { // modify issued token
			if claims.User != nil {
				tokenUser := newTokenUser(claims.User)

				client := server.db.NewReadClient()
				var user *cs.User
				var err error
				user, err = client.GetUserByUsername(claims.User.Name)
				if err != nil {
					log.Error().Err(err).Msgf("failed to load %s from database during claims update", claims.User.Name)
					claims.User.SetBoolAttr(userRejected, true)
					return claims
				}
				// create a new user for this oauth user if it's discord
				if user == nil {
					if tokenUser.discordID() != "" {
						if _, err = server.createNewDiscordUser(tokenUser); err != nil {
							log.Error().Err(err).Msgf("failed to load %s from database during claims update", claims.User.Name)
							claims.User.SetBoolAttr(userRejected, true)
						}
					} else {
						log.Error().Err(err).Msgf("failed to load %s from database during claims update", claims.User.Name)
						claims.User.SetBoolAttr(userRejected, true)
					}
					return claims
				} else {
					// make sure the claim knows about the database id
					tokenUser.setDatabaseID(user.ID)
					tokenUser.setGameID(user.GameID)
					tokenUser.setPlayerNum(user.PlayerNum)
					tokenUser.SetRole(string(user.Role))

					// if we're admin, set the admin claim
					if user.Role == cs.RoleAdmin { // set attributes for admin
						claims.User.SetAdmin(true)
					} else if user.Role == cs.RoleGuest {
						claims.User.SetRole("guest")
					} else if user.IsDiscordUser() {
						// update the discord user on auth
						if err := server.updateUser(tokenUser, user); err != nil {
							log.Error().Err(err).Msgf("failed to load %s from database during claims update", claims.User.Name)
							claims.User.SetBoolAttr("blocked", true)
						}
					}
					return claims
				}
			}
			return claims
		}),
		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool { // rejects some tokens
			if claims.User != nil {
				return !claims.User.BoolAttr(userRejected)
			}
			return false
		}),
		Logger:      authLogger, // optional logger for auth library
		UseGravatar: true,       // for verified provider use gravatar service
	}

	// create auth service
	service := auth.NewService(options)

	service.AddDirectProvider("local", provider.CredCheckerFunc(func(username, password string) (ok bool, err error) {

		client := server.db.NewReadClient()
		user, err := client.GetUserByUsername(username)
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

	AddGuestProvider(service, issuer, authLogger, "guest", HashCheckerFunc(func(hash string) (username string, attributes map[string]interface{}, err error) {
		client := server.db.NewReadClient()
		user, err := client.GetGuestUser(hash)
		if err != nil {
			log.Error().Err(err).Str("Hash", hash).Msg("get user from database")
			return "", nil, err
		}

		if user == nil {
			log.Error().Str("Hash", hash).Msg("user not found")
			return "", nil, nil
		}

		// Check for username and password match
		return user.Username, map[string]interface{}{attrGameID: user.GameID}, nil
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
				tokenUser := newTokenUser(&token.User{
					ID:   "discord_" + token.HashID(sha1.New(), username),
					Name: username,
				})
				tokenUser.setDiscordID(id)
				tokenUser.setDiscordAvatar(avatar)
				return *tokenUser.User
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

	// wrap requests in a transaction
	// do this before Recoverer so we can rollback panics
	r.Use(server.dbClientMiddleware)
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
			r.Post("/starbase-upgrade-cost", server.getStarbaseUpgradeCost)
		})

		// admin routes
		r.Route("/api/admin", func(r chi.Router) {
			r.Use(server.adminRequired)
			r.Get("/games", server.allGames)
			r.Get("/users", server.users)
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
				r.Get("/guest/{num:[0-9]+}", server.getGuestUser)
				r.Post("/join", server.joinGame)
				r.Post("/leave", server.leaveGame)
				r.Post("/add-ai", server.addOpenPlayerSlot)
				r.Post("/add-open-player-slot", server.addOpenPlayerSlot)
				r.Post("/add-guest-player", server.addGuestPlayer)
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
					r.Post("/research-cost", server.getResearchCost)

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

	// The HTTP Server
	httpServer := &http.Server{Addr: ":8080", Handler: r}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		log.Info().Msg("shutdown signal received")

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		log.Info().Msg("shutting down http server")
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal().Err(err).Msg("graceful shutdown failed")
		}
		// close the db
		log.Info().Msg("closing database")
		if err = dbConn.Close(); err != nil {
			log.Fatal().Err(err).Msg("close db failed")
		}
		serverStopCtx()
	}()

	// Run the httpServer
	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server closed")
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	log.Info().Msg("shutdown complete")

	return nil
}

// custom dbClient Middleware to begin a dbClientMiddleware and commit it if successful
func (s *server) dbClientMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// for GET, just use the read client, no need to create a transaction
		if strings.ToUpper(r.Method) == "GET" {
			ctx := context.WithValue(r.Context(), keyDb, s.db.NewReadClient())
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// for POST/PUT/DELETE, etc wrap this request in a transaction
		ctx := context.WithValue(r.Context(), keyDb, s.db.NewReadWriteClient())

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func (s *server) contextDb(r *http.Request) DBClient {
	return r.Context().Value(keyDb).(DBClient)
}

// create a new gameRunner for this request
func (s *server) newGameRunner() GameRunner {
	return NewGameRunner(s.db, s.config)
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
