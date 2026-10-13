package cmd

import (
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newServeCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the webserver",
		Long:  `Start a local webserver and serve requests.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()
			server.Start(*cfg)
			return nil
		},
	}
	serveCmd.Flags().Bool("test-mode", false, "Use a temporary file-based database for testing")

	serveCmd.Flags().Bool("ai-mode", false, "Use a temporary file-based database containing a finished game for each AI race, to view as an admin")

	// bind these flags so config can pick them up for overrides
	viper.BindPFlag("test-mode", serveCmd.Flags().Lookup("test-mode"))
	viper.BindPFlag("ai-mode", serveCmd.Flags().Lookup("ai-mode"))

	return serveCmd
}

func init() {
	rootCmd.AddCommand(newServeCmd())
}
