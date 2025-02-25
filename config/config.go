package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
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
	UsersFilename         string `yaml:"UsersFilename,omitempty"`
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

var testModeConfig = Config{
	Database: databaseConfig{
		Filename: ":memory:?cache=shared",
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

var config *Config

func GetConfig() *Config {
	if config == nil {
		if viper.GetBool("test-mode") {
			// test mode uses an in memory db, no discord auth
			config = &testModeConfig
			log.Debug().Msgf("Config (test mode) : %+v", config)
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
		viper.SetDefault("Database.UsersFilename", "data/users.db")
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
		log.Debug().Msgf("Database.Filename : %v", config.Database.Filename)
		log.Debug().Msgf("Database.UsersFilename : %v", config.Database.UsersFilename)
		log.Debug().Msgf("Config : %+v", config)
		if config.GeneratedUserPassword != "" {
			log.Debug().Msgf("GeneratedUserPassword is set")
		}
	}
	return config
}

func init() {

}
