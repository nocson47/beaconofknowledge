package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Configuration holds the application configuration values.
type Configuration struct {
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        int    `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	MongoHost     string `mapstructure:"MONGO_HOST"`
	MongoPort     int    `mapstructure:"MONGO_PORT"`
	MongoUser     string `mapstructure:"MONGO_USER"`
	MongoPassword string `mapstructure:"MONGO_PASSWORD"`
	MongoDBName   string `mapstructure:"MONGO_DBNAME"`
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     int    `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
}

// LoadConfig loads the configuration from a .env file or environment variables.
func LoadConfig() (Configuration, error) {
	v := viper.New()

	// Set .env file configuration
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables if set")
	}

	// Unmarshal configuration into the struct
	var cfg Configuration
	if err := v.Unmarshal(&cfg); err != nil {
		return Configuration{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return cfg, nil
}
