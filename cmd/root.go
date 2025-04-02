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

// Versioning information used by server.
var (
	semver    string = "0.0.0-develop"
	commit    string = "local"
	buildTime string = time.Now().Format(time.RFC3339)
)

var logFilePath string

// prerun method for enabling debug logging, using separate loggers
// for both logfiles and normal logging.
// TODO: Migrate to stdlib slog package
func logPreRun(cmd *cobra.Command, args []string) error {
	// create a console writer by default, piping log outputs to the logfile if not already doing so
	var writer io.Writer = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}

	if logFilePath != "" {
		if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
			return fmt.Errorf("failed to create log dir %s: \n%w", filepath.Base(logFilePath), err)
		}
		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			return fmt.Errorf("failed to create logfile at %q: \n%w", logFilePath, err)
		}

		writer = zerolog.MultiLevelWriter(writer, logFile)
	}
	log.Logger = log.Output(writer)
	/* 	// Override the default interface marshaler func to indent values
	   	zerolog.InterfaceMarshalFunc = func(v any) ([]byte, error) {
	   		var buf bytes.Buffer
	   		encoder := json.NewEncoder(&buf)
	   		encoder.SetEscapeHTML(false)
	   		encoder.SetIndent("", "\t")
	   		err := encoder.Encode(v)
	   		if err != nil {
	   			return nil, err
	   		}
	   		b := buf.Bytes()
	   		if len(b) > 0 {
	   			// Remove trailing \n which is added by Encode.
	   			return b[:len(b)-1], nil
	   		}
	   		return b, nil
	   	} */
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Debug().Msg("Debug logging enabled")
	log.Info().Msgf("version: %s, build: %s (%s)", semver, commit, buildTime)
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: semver,
	Use:     "craig-stars",
	Short:   "CLI for crajg-stars server",
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
	rootCmd.PersistentFlags().StringVarP(&logFilePath, "log", "", "", "log file to send structured logs to")
}
