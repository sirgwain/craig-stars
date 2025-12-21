// The `server` package configures webserver routes to access the database.
// It is the "glue" that ties the [cs] and [db] packages together.
package server

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"log/slog"

	"connectrpc.com/connect"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	configpkg "github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/hash"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
	"github.com/sirgwain/craig-stars/test/testgames"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/sync/singleflight"

	"github.com/go-pkgz/auth/v2"
	"github.com/go-pkgz/auth/v2/avatar"
	"github.com/go-pkgz/auth/v2/logger"
	"github.com/go-pkgz/auth/v2/provider"
	"github.com/go-pkgz/auth/v2/token"
)

type contextKey int

const (
	keyDb contextKey = iota
	keyDbRead
	keyDbWrite
	keyUserSession
	keyRace
	keyGame
	keyGamePlayer
	keyPlayer
	keyShipDesign
	keyBattlePlan
	keyProductionPlan
	keyTransportPlan
	keyPlanet
	keyFleet
	keyMinefield
	keyUser
)

type server struct {
	db              DBConnection
	config          configpkg.Config
	sf              singleflight.Group
	discordNotifier *discordNotifier
}

const userRejected = "rejected"

// Start the webserver, expose all the routes, inject all the middleware, etc.
func Start(config configpkg.Config) error {
	ctx := context.Background()
	dbConn := db.NewConn()
	if err := dbConn.Connect(ctx, &config); err != nil {
		return fmt.Errorf("failed to connect to database %v", err)
	}

	if viper.GetBool("test-mode") {
		if err := testgames.CreateTestGames(dbConn.NewReadWriteClient()); err != nil {
			return err
		}
	}

	discordNotifier := newDiscordNotifier(dbConn, config)

	// create a server
	server := &server{
		db:              dbConn,
		config:          config,
		discordNotifier: discordNotifier,
	}

	var authLogger = logger.Func(func(format string, args ...interface{}) { slog.Info(fmt.Sprintf(format, args...)) })

	cookieDuration := time.Hour * 24
	if config.Discord.CookieDuration != "" {
		duration, err := time.ParseDuration(config.Discord.CookieDuration)
		if err != nil {
			slog.Error("failed to load cookie duration from config", slog.Any("error", err), slog.String("cookieDuration", config.Discord.CookieDuration))
		} else {
			cookieDuration = duration

		}
	}
	issuer := "craig-stars"
	options := auth.Opts{
		SecretReader: token.SecretFunc(func(_ string) (string, error) { // secret key for JWT, ignores aud
			return server.config.Auth.Secret, nil
		}),
		TokenDuration:     time.Minute,                           // short token, refreshed automatically
		CookieDuration:    cookieDuration,                        // cookie fine to keep for long time
		SecureCookies:     server.config.Auth.SecureCookie,       // true for deployment, false for local dev
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
				// if DiscordID is there, query the user by that
				if tokenUser.discordID() != "" {
					user, err = client.GetUserByDiscordID(context.Background(), tokenUser.discordID())
				} else {
					user, err = client.GetUserByUsername(context.Background(), claims.User.Name)
				}
				if err != nil {
					slog.Error("failed to load user from database during claims update", slog.Any("error", err), slog.String("user", claims.User.Name))
					claims.User.SetBoolAttr(userRejected, true)
					return claims
				}
				// create a new user for this oauth user if it's discord
				if user == nil {
					if tokenUser.discordID() != "" {
						if _, err = server.createNewDiscordUser(context.Background(), tokenUser); err != nil {
							slog.Error("failed to load user from database during claims update", slog.Any("error", err), slog.String("user", claims.User.Name))
							claims.User.SetBoolAttr(userRejected, true)
						}
					} else {
						slog.Error("failed to load user from database during claims update", slog.Any("error", err), slog.String("user", claims.User.Name))
						claims.User.SetBoolAttr(userRejected, true)
					}
					return claims
				} else {
					// make sure the claim knows about the database id
					tokenUser.setDatabaseID(user.ID)
					tokenUser.SetRole(string(user.Role))

					// if we're admin, set the admin claim
					if user.Role == cs.RoleAdmin { // set attributes for admin
						claims.User.SetAdmin(true)
					} else if user.Role == cs.RoleGuest {
						claims.User.SetRole("guest")
					} else if user.IsDiscordUser() {
						// update the discord user on auth
						if err := server.updateUser(context.Background(), tokenUser, user); err != nil {
							slog.Error("failed to load user from database during claims update", slog.Any("error", err), slog.String("user", claims.User.Name))
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
		user, err := client.GetUserByUsername(context.Background(), username)
		if err != nil {
			slog.Error("get user from database", slog.Any("error", err), slog.String("Username", username))
			return false, err
		}

		if user == nil {
			slog.Error("user not found", slog.String("Username", username))
			return false, nil
		}

		// Check for username and password match
		return hash.ComparePassword(user.Password, password)
	}))

	AddGuestProvider(service, issuer, authLogger, "guest", HashCheckerFunc(func(hash string) (username string, attributes map[string]interface{}, err error) {
		client := server.db.NewReadClient()
		user, err := client.GetGuestUser(context.Background(), hash)
		if err != nil {
			slog.Error("get user from database", slog.Any("error", err), slog.String("Hash", hash))
			return "", nil, err
		}

		if user == nil {
			slog.Error("user not found", slog.String("Hash", hash))
			return "", nil, nil
		}

		// Check for username and password match
		return user.Username, nil, nil
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
	r.Use(requestLogger())

	// wrap requests in a transaction
	// do this before Recoverer so we can rollback panics
	r.Use(server.dbClientMiddleware)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(middleware.Heartbeat("/api/ping"))

	r.Group(func(r chi.Router) {
		r.Use(m.Auth)
		r.Use(server.userSessionCtx)
		r.Use(render.SetContentType(render.ContentTypeJSON))

		// route for a few leftover REST operations on games
		r.Route("/api/games", func(r chi.Router) {
			// game by id operations
			r.Route("/{id:[0-9]+}", func(r chi.Router) {
				r.Use(server.gameCtx)
				r.Get("/ping-discord", server.pingDiscordForGameUpdate)
				r.Get("/compute-specs", server.computeSpecs)
			})
		})

	})

	// Create a subrouter for /api/grpc for grpc calls
	grpc := http.NewServeMux()
	userInterceptors := []connect.Interceptor{newDbInterceptor(dbConn), newErrorLogInterceptor()}
	gameInterceptors := []connect.Interceptor{newDbInterceptor(dbConn), newGameInterceptor(), newErrorLogInterceptor()}

	grpc.Handle(craig_starsv1connect.NewAdminServiceHandler(NewAdminServiceHandler(dbConn), connect.WithInterceptors(userInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewTestServiceHandler(NewTestServiceHandler(dbConn), connect.WithInterceptors(userInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewTechServiceHandler(NewTechServiceHandler(), connect.WithInterceptors(newErrorLogInterceptor())))
	grpc.Handle(craig_starsv1connect.NewUserServiceHandler(NewUserServiceHandler(dbConn, discordNotifier), connect.WithInterceptors(userInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewRaceServiceHandler(NewRaceServiceHandler(dbConn), connect.WithInterceptors(userInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewGameServiceHandler(NewGameServiceHandler(dbConn, server.config, discordNotifier), connect.WithInterceptors(gameInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewPlayerServiceHandler(NewPlayerServiceHandler(dbConn, server.config, discordNotifier), connect.WithInterceptors(gameInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewShipDesignServiceHandler(NewshipDesignServiceHandler(dbConn), connect.WithInterceptors(gameInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewBattlePlanServiceHandler(NewBattlePlanServiceHandler(dbConn), connect.WithInterceptors(gameInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewProductionPlanServiceHandler(NewProductionPlanServiceHandler(), connect.WithInterceptors(gameInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewTransportPlanServiceHandler(NewTransportPlanServiceHandler(), connect.WithInterceptors(gameInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewPlanetServiceHandler(NewPlanetServiceHandler(dbConn), connect.WithInterceptors(gameInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewFleetServiceHandler(NewFleetServiceHandler(dbConn), connect.WithInterceptors(gameInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewMinefieldServiceHandler(NewMinefieldServiceHandler(dbConn), connect.WithInterceptors(gameInterceptors...)))
	grpc.Handle(craig_starsv1connect.NewBattleServiceHandler(NewBattleServiceHandler(), connect.WithInterceptors(userInterceptors...)))

	// Mount the grpc calls to /api/grpc
	r.Group(func(r chi.Router) {
		r.Use(m.Auth)
		r.Use(server.userSessionCtx)

		r.Mount("/api/grpc", http.StripPrefix("/api/grpc", grpc))
	})

	// setup auth routes
	authRoutes, avaRoutes := service.Handlers()
	r.Mount("/api/auth", authRoutes)  // add auth handlers
	r.Mount("/api/avatar", avaRoutes) // add avatar handler

	// The HTTP Server
	httpServer := &http.Server{Addr: config.Address, Handler: r}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		slog.Info("shutdown signal received")

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				slog.Error("graceful shutdown timed out - forcing exit.")
				os.Exit(1)
			}
		}()

		// Trigger graceful shutdown
		slog.Info("shutting down http server")
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			slog.Error("graceful shutdown failed", slog.Any("error", err))
			os.Exit(1)
		}
		// close the db
		slog.Info("closing database")
		if err = dbConn.Close(); err != nil {
			slog.Error("close db failed", slog.Any("error", err))
			os.Exit(1)
		}

		// cleanup temporary test database file if in test mode
		configpkg.CleanupTestDatabase()

		serverStopCtx()
	}()

	// Run the httpServer
	slog.Info("starting http server")
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server closed", slog.Any("error", err))
		os.Exit(1)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	slog.Info("shutdown complete")

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

// create a new request logger with slog
func requestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				t2 := time.Now()

				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					slog.Error("system error",
						slog.Any("info", rec),
						slog.String("stack", string(debug.Stack())))
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				// log end request
				attrs := []slog.Attr{
					slog.String("ip", r.RemoteAddr),
					slog.String("url", r.URL.Path),
					slog.String("method", r.Method),
					slog.Int("status", ww.Status()),
					slog.Float64("ms", float64(t2.Sub(t1).Nanoseconds())/1000000.0),
					slog.String("content-length", r.Header.Get("Content-Length")),
					slog.Int("resp_bytes", ww.BytesWritten()),
				}

				// always log the user_agent for now
				attrs = append(attrs, slog.String("user_agent", r.Header.Get("User-Agent")))

				if ww.Status() >= 400 {
					slog.LogAttrs(r.Context(), slog.LevelError, "", attrs...)
				} else {
					slog.LogAttrs(r.Context(), slog.LevelInfo, "", attrs...)
				}
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
