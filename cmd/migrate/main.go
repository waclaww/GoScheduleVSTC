package main

import (
	"gateway/internal/adapter/postgres"
	"gateway/internal/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewDB(cfg.DB_URL)
	if err != nil {
		log.Fatal(err)
	}

	if err := postgres.AutoMigrate(db); err != nil {
		log.Fatal(err)
	}

	log.Println("Миграции применены!")
}
