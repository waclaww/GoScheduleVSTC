package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort string
	DB_URL   string
}

func Load() (Config, error) {
	_ = godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return Config{}, fmt.Errorf("Database URL is not set")
	}

	port := os.Getenv("HTTPPort")
	if port == "" {
		return Config{}, fmt.Errorf("Database PORT is not set")
	}

	return Config{HTTPPort: port, DB_URL: dbUrl}, nil

}
