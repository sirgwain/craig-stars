package cmd

import (
	"github.com/sirgwain/craig-stars/appcontext"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/game"
	"github.com/sirgwain/craig-stars/server"

	"github.com/spf13/cobra"
)

func newServeCmd() *cobra.Command {
	var generateUniverse bool
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the webserver",
		Long:  `Start a local gin-gonic webserver and serve requests.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := appcontext.Initialize()
			// ctx.DB.EnableDebugLogging()

			if generateUniverse {
				if err := generateTestGame(ctx); err != nil {
					return err
				}
			}
			server.Start(ctx)
			return nil
		},
	}
	serveCmd.Flags().BoolVar(&generateUniverse, "generate-test-game", false, "Generate a test user and game")

	return serveCmd
}

func generateTestGame(ctx *appcontext.AppContext) error {
	ctx.DB.MigrateAll()

	admin, adminRace, err := createTestUser(ctx.DB, "admin", "admin", game.RoleAdmin)
	if err != nil {
		return err
	}

	user2, user2Race, err := createTestUser(ctx.DB, "craig", "craig", game.RoleUser)
	if err != nil {
		return err
	}

	// create a game runner to host some games
	gameRunner := server.NewGameRunner(ctx.DB)

	// admin user will host a game with an ai player
	game1, err := gameRunner.HostGame(admin.ID, game.NewGameSettings().
		WithHost(adminRace.ID).
		WithAIPlayer(game.AIDifficultyNormal))
	if err != nil {
		return err
	}

	// generate the universe for game 1
	if err := gameRunner.GenerateUniverse(game1.ID); err != nil {
		return err
	}

	// user2 will also host a game so with an open player slot
	_, err = gameRunner.HostGame(user2.ID, game.NewGameSettings().
		WithName("Joinable Game").
		WithHost(user2Race.ID).
		WithOpenPlayerSlot())
	if err != nil {
		return err
	}

	return nil
}

func createTestUser(db db.Service, username string, password string, role game.Role) (*game.User, *game.Race, error) {
	user, err := db.FindUserByUsername(username)
	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		user = game.NewUser(username, password, role)
		err := db.SaveUser(user)
		if err != nil {
			return nil, nil, err
		}

	}

	races, err := db.GetRaces(user.ID)
	if err != nil {
		return nil, nil, err
	}

	var race *game.Race
	if len(races) == 0 {
		race = game.NewRace()
		race.UserID = user.ID

		if err := db.CreateRace(race); err != nil {
			return nil, nil, err
		}
	} else {
		race = &races[0]
	}

	return user, race, nil
}

func init() {
	rootCmd.AddCommand(newServeCmd())
}
