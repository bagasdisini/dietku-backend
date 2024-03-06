package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// Config is a struct to store configuration from .env file
type Config struct {
	AppHost      string `mapstructure:"APP_HOST"`
	AppPort      string `mapstructure:"PORT"`
	SwaggerHost  string `mapstructure:"SWAGGER_HOST"`
	DBUrl        string `mapstructure:"MONGODB_URI"`
	DBName       string `mapstructure:"MONGODB_NAME"`
	AllowOrigins string `mapstructure:"CORS_ALLOW_ORIGINS"`
}

// InitConfigApp loads configuration from .env file
func InitConfigApp(path string) (*Config, error) {
	var config Config

	// load .env file, if not found, use environment variables instead
	err := godotenv.Load(path)
	if err != nil {
		fmt.Println(".env file not found, using environment variables instead.")
	}

	config.AppHost = os.Getenv("APP_HOST")
	config.AppPort = os.Getenv("PORT")
	config.SwaggerHost = os.Getenv("SWAGGER_HOST")
	config.DBName = os.Getenv("MONGODB_NAME")
	config.AllowOrigins = os.Getenv("CORS_ALLOW_ORIGINS")
	config.DBUrl = os.Getenv("MONGODB_URI")

	if config.DBUrl == "" {
		return &Config{}, errors.New("please check your database setting")
	}

	return &config, err
}
