package config

//config.go defines a struct that loads environmental variables from .env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                  string `env:"ENV,required"`
	DBName               string `env:"DATABASE_NAME,required"`
	DBUser               string `env:"DATABASE_USER,required"`
	DBPassword           string `env:"DATABASE_PASSWORD,required"`
	DBHost               string `env:"DATABASE_HOST,required"`
	DBPort               string `env:"DATABASE_PORT,required"`
	HTTPDomain           string `env:"HTTP_DOMAIN,required"`
	HTTPPort             string `env:"HTTP_PORT,required"`
	HTTPShutdownDuration int
}

func NewConfig() (Config, error) {
	godotenv.Load()

	newConfig := Config{
		Env:                  os.Getenv("ENV"),
		DBName:               os.Getenv("DATABASE_NAME"),
		DBUser:               os.Getenv("DATABASE_USER"),
		DBPassword:           os.Getenv("DATABASE_PASSWORD"),
		DBHost:               os.Getenv("DATABASE_HOST"),
		DBPort:               os.Getenv("DATABASE_PORT"),
		HTTPDomain:           os.Getenv("HTTP_DOMAIN"),
		HTTPPort:             os.Getenv("HTTP_PORT"),
		HTTPShutdownDuration: 10,
	}
	if newConfig.Env == "" || newConfig.DBName == "" || newConfig.DBUser == "" ||
		newConfig.DBPassword == "" || newConfig.DBHost == "" ||
		newConfig.DBPort == "" || newConfig.HTTPDomain == "" || newConfig.HTTPPort == "" {
		return Config{}, fmt.Errorf("missing required field")
	}

	return newConfig, nil

}
