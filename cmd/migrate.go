package cmd

import (
	"github.com/sirgwain/craig-stars/appcontext"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

// migrateCmd represents the serve command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database models",
	Long:  `Start a local gin-gonic webserver and serve requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := appcontext.Initialize()
		err := ctx.DB.MigrateAll()
		if err != nil {
			log.Err(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
