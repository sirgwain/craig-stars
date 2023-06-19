package cmd

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/server"

	"github.com/spf13/cobra"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debug().Msgf("%s took %s", name, elapsed)
}

func newServeCmd() *cobra.Command {
	var generateUniverse bool
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the webserver",
		Long:  `Start a local webserver and serve requests.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			db := db.NewClient()
			cfg := config.GetConfig()
			if err := db.Connect(cfg); err != nil {
				return err
			}

			if generateUniverse {
				if err := generateTestGame(db, *cfg); err != nil {
					return err
				}
			}
			server.Start(db, *cfg)
			return nil
		},
	}
	serveCmd.Flags().BoolVar(&generateUniverse, "generate-test-game", false, "Generate a test user and game")

	return serveCmd
}

func generateTestGame(db server.DBClient, config config.Config) error {
	defer timeTrack(time.Now(), "generateTestGame")

	admin, adminRace, err := createTestUser(db, "admin", config.GeneratedUserPassword, "admin@craig-stars.net", cs.RoleAdmin)
	if err != nil {
		return err
	}

	// create a game runner to host some games
	gameRunner := server.NewGameRunner(db)

	// admin user will host a game with an ai player
	if _, err := gameRunner.HostGame(admin.ID, cs.NewGameSettings().
		// WithSize(game.SizeTiny).
		// WithDensity(game.DensitySparse).
		WithHost(adminRace.ID).
		WithAIPlayer(cs.AIDifficultyNormal, 1)); err != nil {
		return err
	}

	// also create a medium size game with 25 turns generated
	mediumGame, err := gameRunner.HostGame(admin.ID, cs.NewGameSettings().
		WithName("Medium Game").
		WithSize(cs.SizeMedium).
		WithHost(adminRace.ID).
		WithAIPlayerRace(cs.HEs(), cs.AIDifficultyNormal, 0).
		WithAIPlayerRace(cs.SSs(), cs.AIDifficultyNormal, 1).
		WithAIPlayerRace(cs.WMs(), cs.AIDifficultyNormal, 2).
		WithAIPlayerRace(cs.CAs(), cs.AIDifficultyNormal, 3).
		WithAIPlayerRace(cs.ISs(), cs.AIDifficultyNormal, 0).
		WithAIPlayerRace(cs.SDs(), cs.AIDifficultyNormal, 1).
		WithAIPlayerRace(cs.PPs(), cs.AIDifficultyNormal, 2).
		WithAIPlayerRace(cs.ITs(), cs.AIDifficultyNormal, 3).
		WithAIPlayerRace(cs.ARs(), cs.AIDifficultyNormal, 0).
		WithAIPlayerRace(cs.JoaTs(), cs.AIDifficultyNormal, 1))

	if err != nil {
		return err
	}
	mediumGame.Players[0].AIControlled = true
	db.UpdateLightPlayer(mediumGame.Players[0])

	for i := 0; i < 40; i++ {
		gameRunner.SubmitTurn(mediumGame.ID, mediumGame.HostID)
		if _, err := gameRunner.CheckAndGenerateTurn(mediumGame.ID); err != nil {
			log.Error().Err(err).Msg("check and generate new turn")
		}
	}

	mediumGame.Players[0].AIControlled = false
	db.UpdateLightPlayer(mediumGame.Players[0])

	// also create a medium size game with 25 turns generated
	tinyGame, err := gameRunner.HostGame(admin.ID, cs.NewGameSettings().
		WithName("Tiny Game").
		WithSize(cs.SizeTiny).
		WithHost(adminRace.ID).
		WithAIPlayerRace(cs.HEs(), cs.AIDifficultyNormal, 0).
		WithAIPlayerRace(cs.SDs(), cs.AIDifficultyNormal, 1).
		WithAIPlayerRace(cs.PPs(), cs.AIDifficultyNormal, 2).
		WithAIPlayerRace(cs.ITs(), cs.AIDifficultyNormal, 3).
		WithAIPlayerRace(cs.JoaTs(), cs.AIDifficultyNormal, 1),
	)
	if err != nil {
		return err
	}

	tinyGame.Players[0].AIControlled = true
	db.UpdateLightPlayer(tinyGame.Players[0])

	for i := 0; i < 9; i++ {
		gameRunner.SubmitTurn(tinyGame.ID, tinyGame.HostID)
		if _, err := gameRunner.CheckAndGenerateTurn(tinyGame.ID); err != nil {
			log.Error().Err(err).Msg("check and generate new turn")
		}
	}

	tinyGame.Players[0].AIControlled = false
	db.UpdateLightPlayer(tinyGame.Players[0])

	// create a new tiny game for each race
	races, err := db.GetRacesForUser(admin.ID)
	if err != nil {
		return err
	}

	for _, race := range races {
		_, err := gameRunner.HostGame(admin.ID, cs.NewGameSettings().
			WithName(fmt.Sprintf("%s Test Game", race.PluralName)).
			WithSize(cs.SizeTiny).
			WithHost(race.ID).
			WithAIPlayer(cs.AIDifficultyNormal, 1).
			WithAIPlayer(cs.AIDifficultyNormal, 2))
		if err != nil {
			return err
		}
	}

	// user2, user2Race, err := createTestUser(db, "craig", config.GeneratedUserPassword, "craig@craig-stars.net", cs.RoleUser)
	// if err != nil {
	// 	return err
	// }


	// // user2 will also host a game so with an open player slot
	// _, err = gameRunner.HostGame(user2.ID, cs.NewGameSettings().
	// 	WithName("Joinable Game").
	// 	WithHost(user2Race.ID).
	// 	WithOpenPlayerSlot().
	// 	WithAIPlayer(cs.AIDifficultyNormal, 1))
	// if err != nil {
	// 	return err
	// }

	return nil
}

func createTestUser(db server.DBClient, username string, password string, email string, role string) (*cs.User, *cs.Race, error) {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		return nil, nil, err
	}

	// default the password to the username if it's empty
	if password == "" {
		password = username
	}

	if user == nil {
		user, err = cs.NewUser(username, password, email, role)
		if err != nil {
			return nil, nil, err
		}
		err := db.CreateUser(user)
		if err != nil {
			return nil, nil, err
		}

	}

	races, err := db.GetRacesForUser(user.ID)
	if err != nil {
		return nil, nil, err
	}

	var race cs.Race
	if len(races) == 0 {
		race = cs.Humanoids()
		race.UserID = user.ID

		if err := db.CreateRace(&race); err != nil {
			return nil, nil, err
		}
	} else {
		race = races[0]
	}

	// ensure we have other races
	if err := ensureTestRaces(db, user.ID); err != nil {
		return nil, nil, err
	}

	return user, &race, nil
}

// ensure this user has test races for all PRTs
func ensureTestRaces(db server.DBClient, userID int64) error {

	races, err := db.GetRacesForUser(userID)
	if err != nil {
		return nil
	}

	prts := make(map[cs.PRT]bool, len(cs.PRTs))
	for _, race := range races {
		prts[race.PRT] = true
	}

	for _, prt := range cs.PRTs {
		if found := prts[prt]; !found {
			race := cs.NewRace()
			race.UserID = userID
			race.PRT = prt
			race.Name = fmt.Sprintf("%v", prt)
			race.PluralName = fmt.Sprintf("%vs", prt)
			if err := db.CreateRace(race); err != nil {
				return err
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(newServeCmd())
}
