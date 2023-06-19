package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Filename string `yaml:"filename,omitempty"`
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

		// Set default value
		viper.SetDefault("Database.Filename", "data/data.db")

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
		log.Debug().Msgf("DataDir : %v", config)
		if config.GeneratedUserPassword != "" {
			log.Debug().Msgf("GeneratedUserPassword is set")
		}
	}
	return config
}

func init() {

}
