package server

import (
	"crypto/sha1"
	"fmt"
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
	"github.com/go-pkgz/rest"
)

const databaseIDAttr = "DatabaseID"

type contextKey int

const (
	keyUser contextKey = iota
	keyRace
)

func StartChi(db DBClient, config config.Config) {

	// static resources
	sub, err := fs.Sub(assets, "frontend/build")
	if err != nil {
		panic(err)
	}

	// create a server
	server := &server{
		db:            db,
		config:        config,
		gameRunner:    NewGameRunner(db),
		playerUpdater: newPlayerUpdater(db),
	}
	_ = server

	var authLogger = logger.Func(func(format string, args ...interface{}) { log.Info().Msgf(format, args...) })

	options := auth.Opts{
		SecretReader: token.SecretFunc(func(_ string) (string, error) { // secret key for JWT, ignores aud
			return server.config.Auth.Secret, nil
		}),
		TokenDuration:     time.Minute,                           // short token, refreshed automatically
		CookieDuration:    time.Hour * 24,                        // cookie fine to keep for long time
		DisableXSRF:       true,                                  // don't disable XSRF in real-life applications!
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
				}

				if claims.User != nil && claims.User.Name == "admin" { // set attributes for admin
					claims.User.SetAdmin(true)
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
				discriminator := data.Value("discriminator")
				fullUsername := fmt.Sprintf("%s#%s", username, discriminator)
				id := data.Value("id")
				avatar := data.Value("avatar")
				userInfo := token.User{
					ID: "discord_" + token.HashID(sha1.New(),
						fullUsername),
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

	r.Group(func(r chi.Router) {
		r.Use(m.Auth)
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Use(server.userCtx)
		r.Get("/api/me", me_chi)
		r.Get("/api/games", server.games_chi)
		r.Get("/api/games/open", server.gamesOpen_chi)

		// race CRUD
		r.Route("/api/races", func(r chi.Router) {
			r.Post("/", server.createRace)
			r.Get("/", server.races)
			r.Get("/points", server.getRacePoints)

			// Subrouters:
			r.Route("/{id:[0-9]+}", func(r chi.Router) {
				r.Use(server.raceCtx)
				r.Get("/", server.race)
				r.Put("/", server.updateRace)
				r.Delete("/", server.deleteRace)
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
						Str("type", "error").
						Timestamp().
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Msg("log system error")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				// log end request
				var event *zerolog.Event
				if ww.Status() >= 400 {
					event = log.Error()
				} else {
					event = log.Info()
				}
				event.
					Timestamp().
					Fields(map[string]interface{}{
						"remote_ip":  r.RemoteAddr,
						"url":        r.URL.Path,
						"method":     r.Method,
						"user_agent": r.Header.Get("User-Agent"),
						"status":     ww.Status(),
						"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
						"bytes_in":   r.Header.Get("Content-Length"),
						"bytes_out":  ww.BytesWritten(),
					}).Msg("")
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

func (s *server) games_chi(w http.ResponseWriter, r *http.Request) {
	user := s.mustGetUser(w, r)

	games, err := s.db.GetGamesForUser(user.ID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, games)
}

func (s *server) gamesOpen_chi(w http.ResponseWriter, r *http.Request) {
	user := s.mustGetUser(w, r)

	games, err := s.db.GetOpenGames(user.ID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, games)
}

// get the user from the context
func (s *server) contextUser(r *http.Request) sessionUser {
	return r.Context().Value(keyUser).(sessionUser)
}

func (s *server) mustGetUser(w http.ResponseWriter, r *http.Request) sessionUser {
	userInfo, err := token.GetUserInfo(r)
	if err != nil {
		panic("failed to load user")
	}

	userID, err := strconv.ParseInt(userInfo.StrAttr(databaseIDAttr), 10, 64)
	if err != nil {
		panic("failed to load user")
	}

	return sessionUser{
		ID:       userID,
		Username: userInfo.Name,
		Role:     userInfo.Role,
	}
}

func me_chi(w http.ResponseWriter, r *http.Request) {

	userInfo, err := token.GetUserInfo(r)
	if err != nil {
		log.Printf("failed to get user info, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, err := strconv.ParseInt(userInfo.StrAttr(databaseIDAttr), 10, 64)
	if err != nil {
		log.Printf("failed to get user info, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := sessionUser{
		ID:       userID,
		Username: userInfo.Name,
		Role:     userInfo.Role,
	}

	rest.RenderJSON(w, res)
}
