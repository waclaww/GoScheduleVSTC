package config

import (
	"fmt"
	"net"
	"net/url"
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
		dbUrl = dbUrlFromParts()
	}
	if dbUrl == "" {
		return Config{}, fmt.Errorf("Database URL is not set")
	}

	port := os.Getenv("HTTPPort")
	if port == "" {
		port = "8000"
	}

	return Config{HTTPPort: port, DB_URL: dbUrl}, nil
}

func dbUrlFromParts() string {
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	if host == "" || name == "" {
		return ""
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")),
		Host:     net.JoinHostPort(host, port),
		Path:     name,
		RawQuery: "sslmode=disable",
	}

	return u.String()
}
