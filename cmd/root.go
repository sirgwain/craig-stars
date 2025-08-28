package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/phsym/console-slog"
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/cobra"
)

// semver is the semantic-release semver (added at compile time)
var (
	semver    string = "0.0.0-develop"
	commit    string = "local"
	buildTime string = time.Now().Format(time.RFC3339)
)

var logFile string
var debugEnabled = true

// prerun method for enabling slog logging
func logPreRun(cmd *cobra.Command, args []string) error {
	// pick level based on your flag/env/config
	level := slog.LevelInfo
	if debugEnabled { // e.g. from a flag
		level = slog.LevelDebug
	}

	// base console handler (pretty printing)
	consoleHandler := console.NewHandler(os.Stderr, &console.HandlerOptions{
		Level: level,
	})

	var handler slog.Handler = consoleHandler

	if logFile != "" {
		if err := os.MkdirAll(filepath.Dir(logFile), 0o755); err != nil {
			return fmt.Errorf("failed to create log dir %s: \n%w", filepath.Base(logFile), err)
		}
		logFileWriter, err := os.OpenFile(
			logFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0o664,
		)
		if err != nil {
			return fmt.Errorf("failed to create log file %s: \n%w", logFile, err)
		}

		// add a JSON handler for the log file
		fileHandler := slog.NewJSONHandler(logFileWriter, &slog.HandlerOptions{
			Level: level,
		})

		// wrap both handlers with MultiHandler
		handler = slogmulti.Fanout(consoleHandler, fileHandler)
	}

	// build the logger and set global
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// emit startup logs
	slog.Debug("Debug logging enabled")
	slog.Info("starting up",
		slog.String("version", semver),
		slog.String("commit", commit),
		slog.String("build", buildTime),
	)

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
