package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database              databaseConfig
	Auth                  authConfig
	Discord               discordConfig
	Game                  gameConfig
	GeneratedUserPassword string
	Address               string
}

type databaseConfig struct {
	Filename              string `yaml:"Filename,omitempty"`
	ReadConnectionParams  string `yaml:"ReadConnectionParams"`
	WriteConnectionParams string `yaml:"WriteConnectionParams"`
	DebugLogging          bool   `yaml:"DebugLogging,omitempty"`
	SkipUpgrade           bool   `yaml:"SkipUpgrade,omitempty"`
}

type authConfig struct {
	Secret       string `yaml:"Secret,omitempty"`
	URL          string `yaml:"URL,omitempty"`
	DisableXSRF  bool   `yaml:"DisableXSRF,omitempty"`
	SecureCookie bool   `yaml:"SecureCookie,omitempty"`
}

type discordConfig struct {
	Enabled               bool   `yaml:"Enabled,omitempty"`
	ClientID              string `yaml:"ClientID,omitempty"`
	ClientSecret          string `yaml:"ClientSecret,omitempty"`
	CookieDuration        string `yaml:"CookieDuration,omitempty"`
	WebhookNotify         bool   `yaml:"WebhookNotify,omitempty"`
	WebhookID             string `yaml:"WebhookID,omitempty"`
	WebhookToken          string `yaml:"WebhookToken,omitempty"`
	WebhookNotifyForAdmin bool   `yaml:"WebhookNotifyForAdmin,omitempty"`
}

type gameConfig struct {
	InviteLinkSalt string `yaml:"InviteLinkSalt,omitempty"`
}

func getTestModeConfig() Config {
	// Create a unique temporary database file for this test instance
	tmpDir := os.TempDir()
	timestamp := time.Now().UnixNano()
	dbFile := filepath.Join(tmpDir, fmt.Sprintf("craigstars_test_%d.db", timestamp))

	return Config{
		Database: databaseConfig{
			Filename: dbFile,
		},
		Auth: authConfig{
			Secret:      "testSecret",
			DisableXSRF: true,
			URL:         "http://localhost:5173",
		},
		Game: gameConfig{
			InviteLinkSalt: "salt",
		},
		Address: "localhost:8080",
	}
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		if viper.GetBool("test-mode") {
			// test mode uses a temporary file-based db, no discord auth
			testConfig := getTestModeConfig()
			config = &testConfig
			slog.Debug("Config (test mode)", slog.Any("config", config))
			return config
		}

		// setup default config for running in local dev
		path := "./data/config"
		viper.SetConfigName("config")        // config file name without extension
		viper.SetConfigType("yaml")          // yaml type
		viper.AddConfigPath("./data/config") // config file path
		viper.AutomaticEnv()                 // read value ENV variable

		// Set default values for local dev
		viper.SetDefault("Database.Filename", "data/data.db")
		viper.SetDefault("Database.ReadConnectionParams", "?_txlock=deferred")
		viper.SetDefault("Database.WriteConnectionParams", "?_txlock=immediate&_busy_timeout=1200")
		viper.SetDefault("Auth.DisableXSRF", true)            // default for local dev
		viper.SetDefault("Auth.Secret", "secret")             // default for local dev
		viper.SetDefault("Auth.URL", "http://localhost:5173") // default for local dev
		viper.SetDefault("Discord.CookieDuration", "24h")     // default for local dev

		viper.SetDefault("Game.InviteLinkSalt", "salt") // default for local dev
		viper.SetDefault("Address", "localhost:8080")   // default for local dev, override to :8080 for deployment

		// write config if not present
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			panic(fmt.Sprintln("fatal error creating config file directory \n", err))
		}
		viper.SafeWriteConfig()

		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Sprintln("fatal error loading config file: \n", err))
		}

		viper.Unmarshal(&config)

		// Config
		slog.Debug(fmt.Sprintf("Database.Filename : %v", config.Database.Filename))
		slog.Debug(fmt.Sprintf("Config : %+v", config))
		if config.GeneratedUserPassword != "" {
			slog.Debug("GeneratedUserPassword is set")
		}
	}
	return config
}

// CleanupTestDatabase removes the temporary database file if in test mode
func CleanupTestDatabase() {
	if config != nil && viper.GetBool("test-mode") {
		if config.Database.Filename != "" && config.Database.Filename != ":memory:?cache=shared" {
			if err := os.Remove(config.Database.Filename); err != nil {
				slog.Warn("Failed to cleanup test database file", "file", config.Database.Filename, "error", err)
			} else {
				slog.Debug("Cleaned up test database file", "file", config.Database.Filename)
			}
		}
	}
}

func init() {

}
