package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Recreate      bool   `yaml:"Recreate,omitempty"`
		Filename      string `yaml:"Filename,omitempty"`
		UsersFilename string `yaml:"UsersFilename,omitempty"`
		DebugLogging  bool   `yaml:"DebugLogging,omitempty"`
	}
	Auth struct {
		Secret      string `yaml:"Secret,omitempty"`
		URL         string `yaml:"URL,omitempty"`
		DisableXSRF bool   `yaml:"DisableXSRF,omitempty"`
	}
	Discord struct {
		Enabled               bool   `yaml:"Enabled,omitempty"`
		ClientID              string `yaml:"ClientID,omitempty"`
		ClientSecret          string `yaml:"ClientSecret,omitempty"`
		CookieDuration        string `yaml:"CookieDuration,omitempty"`
		WebhookNotify         bool   `yaml:"WebhookNotify,omitempty"`
		WebhookID             string `yaml:"WebhookID,omitempty"`
		WebhookToken          string `yaml:"WebhookToken,omitempty"`
		WebhookNotifyForAdmin bool   `yaml:"WebhookNotifyForAdmin,omitempty"`
	}
	Game struct {
		InviteLinkSalt string `yaml:"InviteLinkSalt,omitempty"`
	}
	GeneratedUserPassword string
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		path := "./data/config"
		viper.SetConfigName("config")        // config file name without extension
		viper.SetConfigType("yaml")          // yaml type
		viper.AddConfigPath("./data/config") // config file path
		viper.AutomaticEnv()                 // read value ENV variable

		// Set default values for local dev
		viper.SetDefault("Database.Filename", "data/data.db")
		viper.SetDefault("Database.UsersFilename", "data/users.db")
		viper.SetDefault("Auth.Secret", "secret")             // default for local dev
		viper.SetDefault("Auth.URL", "http://localhost:5173") // default for local dev
		viper.SetDefault("Auth.DisableXSRF", true)            // default for local dev
		viper.SetDefault("Discord.CookieDuration", "24h")     // default for local dev

		viper.SetDefault("Game.InviteLinkSalt", "salt") // default for local dev

		// write config if not present
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			panic(fmt.Sprintln("fatal error creating config file directory \n", err))
		}
		viper.SafeWriteConfig()

		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Sprintln("fatal error config file: default \n", err))
		}

		viper.Unmarshal(&config)

		// Config
		log.Debug().Msgf("Database.Filename : %v", config.Database.Filename)
		log.Debug().Msgf("Database.UsersFilename : %v", config.Database.UsersFilename)
		log.Debug().Msgf("DataDir : %v", config)
		if config.GeneratedUserPassword != "" {
			log.Debug().Msgf("GeneratedUserPassword is set")
		}
	}
	return config
}

func init() {

}
