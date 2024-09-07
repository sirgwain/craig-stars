//go:build !wasi && !wasm

package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

var logFile string

// prerun method for enabling debug logging
func logPreRun(cmd *cobra.Command, args []string) error {
	// enable debug mode if configured
	var writer io.Writer
	writer = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}
	if logFile != "" {
		if err := os.MkdirAll(filepath.Dir(logFile), 0755); err != nil {
			return fmt.Errorf("failed to create log dir %s %w", filepath.Base(logFile), err)
		}
		logFileWriter, err := os.OpenFile(
			logFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			return fmt.Errorf("failed to create log file %s %w", logFile, err)
		}

		// make a file and console writer
		writer = io.MultiWriter(writer, logFileWriter)
	}
	log.Logger = log.Output(writer)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Debug().Msg("Debug logging enabled")
	log.Info().Msgf(fmt.Sprintf("version: %s, build: %s (%s)", semver, commit, buildTime))
	return nil
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
	PersistentPreRunE: logPreRun,
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
		os.Exit(1)
	}
}

func init() {
	// all commands have debug mode
	rootCmd.PersistentFlags().StringVarP(&logFile, "log", "", "", "log file to send structured logs to")
}
