package config

import (
	"log"

	"github.com/spf13/viper"
)

// A struct which holds the app configuration.
type Config struct {
	PostgresConfig PostgresConfig `mapstructure:"postgresConfig"`
	Port           string         `mapstructure:"port"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBname   string `mapstructure:"dbName"`
}

// Function which reads configuration from config.json.
func LoadAppConfig() Config {
	log.Println("loading server configuration...")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
