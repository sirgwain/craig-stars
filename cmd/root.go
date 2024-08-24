package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

// semver is the semantic-release semver (added at compile time)
var (
	semver    string = "0.0.0-develop"
	commit    string = "local"
	buildTime string = time.Now().Format(time.RFC3339)
)

var debug bool

// prerun method for enabling debug logging
func debugPreRun(cmd *cobra.Command, args []string) {
	// enable debug mode if configured
	if debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled")
	}
	log.Info().Msgf(fmt.Sprintf("version: %s, build: %s (%s)", semver, commit, buildTime))
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: semver,
	Use:     "craig-stars",
	Short:   "A Stars! clone",
	Long: `
craig-stars will start a webserver for playing the game, or act as a
CLI for interacting with the server resources such as users.
`,
	PersistentPreRun: debugPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		// Show usage
		cmd.Help()
		os.Exit(1)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("version: {{ .Version }}\nbuild: %s (%s)", commit, buildTime))
	err := rootCmd.Execute()
	if err != nil {
		dbRollback()
		os.Exit(1)
	}
}

func init() {
	// all commands have debug mode
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug logging")
}
