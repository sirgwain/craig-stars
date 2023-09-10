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

func newServeCmd() *cobra.Command {
	var generate bool
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the webserver",
		Long:  `Start a local webserver and serve requests.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()

			// generate test games if asked
			if generate {
				if err := generateTestGame(*cfg); err != nil {
					return err
				}

				// if we have recreate configured, turn it off after we generate a test game
				// so we don't recreate the db conn again when the server starts
				cfg.Database.Recreate = false
			}

			server.Start(*cfg)
			return nil
		},
	}
	serveCmd.Flags().BoolVar(&generate, "generate-test-game", false, "Generate a test user and game")

	return serveCmd
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debug().Msgf("%s took %s", name, elapsed)
}

func generateTestGame(config config.Config) error {
	defer timeTrack(time.Now(), "generateTestGame")
	// create a new connection for the server to use
	dbConn := db.NewConn()
	if err := dbConn.Connect(&config); err != nil {
		return err
	}
	defer func() { dbConn.Close() }()

	db := dbConn.NewReadWriteClient()

	admin, adminRace, err := createTestUser(db, "admin", config.GeneratedUserPassword, "admin@craig-stars.net", cs.RoleAdmin)
	if err != nil {
		return err
	}

	// create a game runner to host some games
	gameRunner := server.NewGameRunner(dbConn, config)

	// admin user will host a game with an ai player
	if _, err := gameRunner.HostGame(admin.ID, cs.NewGameSettings().
		WithHost(*adminRace).
		WithPublicPlayerScores(true).
		WithAIPlayer(cs.AIDifficultyNormal, 1)); err != nil {
		return err
	}

	// also create a medium size game with 25 turns generated
	mediumGame, err := gameRunner.HostGame(admin.ID, cs.NewGameSettings().
		WithName("Medium Game").
		WithSize(cs.SizeSmall).
		WithDensity(cs.DensityPacked).
		WithPublicPlayerScores(true).
		WithHost(*adminRace).
		// WithAIPlayerRace(cs.HEs().WithName(cs.AINames[0]), cs.AIDifficultyNormal, 0).
		// WithAIPlayerRace(cs.SSs().WithName(cs.AINames[1]), cs.AIDifficultyNormal, 1).
		// WithAIPlayerRace(cs.WMs().WithName(cs.AINames[2]), cs.AIDifficultyNormal, 2).
		// WithAIPlayerRace(cs.CAs().WithName(cs.AINames[3]), cs.AIDifficultyNormal, 3).
		// WithAIPlayerRace(cs.ISs().WithName(cs.AINames[4]), cs.AIDifficultyNormal, 0).
		WithAIPlayerRace(cs.SDs().WithName(cs.AINames[5]), cs.AIDifficultyNormal, 1).
		WithAIPlayerRace(cs.PPs().WithName(cs.AINames[6]), cs.AIDifficultyNormal, 2),
	// 	WithAIPlayerRace(cs.ITs().WithName(cs.AINames[7]), cs.AIDifficultyNormal, 3).
	// 	WithAIPlayerRace(cs.ARs().WithName(cs.AINames[8]), cs.AIDifficultyNormal, 0).
	// 	WithAIPlayerRace(cs.JoaTs().WithName(cs.AINames[9]), cs.AIDifficultyNormal, 1),
	)

	if err != nil {
		return err
	}
	mediumGame.Players[0].AIControlled = true
	db.UpdateLightPlayer(mediumGame.Players[0])

	for i := 0; i < 75; i++ {
		gameRunner.SubmitTurn(mediumGame.ID, mediumGame.HostID)
		if _, err := gameRunner.CheckAndGenerateTurn(mediumGame.ID); err != nil {
			log.Error().Err(err).Msg("check and generate new turn")
		}
	}

	mediumGame.Players[0].AIControlled = false
	db.UpdateLightPlayer(mediumGame.Players[0])

	// // also create a medium size game with 25 turns generated
	// tinyGame, err := gameRunner.HostGame(admin.ID, cs.NewGameSettings().
	// 	WithName("Tiny Game").
	// 	WithSize(cs.SizeTiny).
	// 	WithPublicPlayerScores(true).
	// 	WithHost(*adminRace).
	// 	WithAIPlayerRace(cs.HEs(), cs.AIDifficultyNormal, 0).
	// 	WithAIPlayerRace(cs.SDs(), cs.AIDifficultyNormal, 1).
	// 	WithAIPlayerRace(cs.PPs(), cs.AIDifficultyNormal, 2).
	// 	WithAIPlayerRace(cs.ITs(), cs.AIDifficultyNormal, 3).
	// 	WithAIPlayerRace(cs.JoaTs(), cs.AIDifficultyNormal, 1),
	// )
	// if err != nil {
	// 	return err
	// }

	// tinyGame.Players[0].AIControlled = true
	// db.UpdateLightPlayer(tinyGame.Players[0])

	// for i := 0; i < 9; i++ {
	// 	gameRunner.SubmitTurn(tinyGame.ID, tinyGame.HostID)
	// 	if _, err := gameRunner.CheckAndGenerateTurn(tinyGame.ID); err != nil {
	// 		log.Error().Err(err).Msg("check and generate new turn")
	// 	}
	// }

	// tinyGame.Players[0].AIControlled = false
	// db.UpdateLightPlayer(tinyGame.Players[0])

	// // create a new tiny game for each race
	// races, err := db.GetRacesForUser(admin.ID)
	// if err != nil {
	// 	return err
	// }

	// for _, race := range races {
	// 	_, err := gameRunner.HostGame(admin.ID, cs.NewGameSettings().
	// 		WithName(fmt.Sprintf("%s Test Game", race.PluralName)).
	// 		WithSize(cs.SizeTiny).
	// 		WithHost(race).
	// 		WithAIPlayer(cs.AIDifficultyNormal, 1).
	// 		WithAIPlayer(cs.AIDifficultyNormal, 2))
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// // also create a medium size game with 25 turns generated
	// multiplayerGame, err := gameRunner.HostGame(admin.ID, cs.NewGameSettings().
	// 	WithName("Multiplayer Game").
	// 	WithSize(cs.SizeTiny).
	// 	WithPublicPlayerScores(true).
	// 	WithHost(*adminRace).
	// 	WithOpenPlayerSlot().
	// 	WithAIPlayerRace(cs.SDs(), cs.AIDifficultyNormal, 1).
	// 	WithAIPlayerRace(cs.PPs(), cs.AIDifficultyNormal, 2),
	// )
	// if err != nil {
	// 	return err
	// }
	// _ = multiplayerGame

	return nil
}

func createTestUser(db server.DBClient, username string, password string, email string, role cs.UserRole) (*cs.User, *cs.Race, error) {
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
