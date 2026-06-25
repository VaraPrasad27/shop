package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DATABASE_URL string
	Port         string
}

func LoadConfig() (Config, error) {
	// godotenv.Load is a no-op if .env is absent; we treat .env as the
	// source of truth and let process env override it.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("loading .env: %w", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{DATABASE_URL: databaseURL, Port: port}, nil
}
