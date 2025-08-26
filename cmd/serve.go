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
	serveCmd.Flags().Bool("test-mode", false, "Use an in memory database")

	// bind this flag so config can pick it up for overrides
	viper.BindPFlag("test-mode", serveCmd.Flags().Lookup("test-mode"))

	return serveCmd
}

func init() {
	rootCmd.AddCommand(newServeCmd())
}
