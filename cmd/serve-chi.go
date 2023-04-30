package cmd

import (
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/server"

	"github.com/spf13/cobra"
)

func newServeChiCmd() *cobra.Command {
	var generateUniverse bool
	serveCmd := &cobra.Command{
		Use:   "serve-chi",
		Short: "Start the webserver",
		Long:  `Start a webserver and serve requests.`,
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

			server.StartChi(db, *cfg)
			return nil
		},
	}
	serveCmd.Flags().BoolVar(&generateUniverse, "generate-test-game", false, "Generate a test user and game")

	return serveCmd
}

func init() {
	rootCmd.AddCommand(newServeChiCmd())
}
