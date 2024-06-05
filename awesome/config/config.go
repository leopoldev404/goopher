package config

import (
	"log"

	"github.com/spf13/viper"
)

type ApiConfig struct {
	Address string
}

type DatabaseConfig struct {
	ConnectionString string
}

type Configuration struct {
	Api      *ApiConfig
	Database *DatabaseConfig
}

var config *Configuration

func LoadConfiguration() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	log.Println("Init Configuration 🧭!")
}

func GetDatabaseConfiguration() *DatabaseConfig {
	return config.Database
}

func GetApiConfiguration() *DatabaseConfig {
	return config.Database
}
