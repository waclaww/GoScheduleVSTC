package main

import (
	adapterhttp "gateway/internal/adapter/http"
	"gateway/internal/adapter/http/handlers"
	"gateway/internal/adapter/postgres"
	"gateway/internal/config"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := postgres.NewDB(cfg.DB_URL)
	if err != nil {
		return err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDb.Close()

	_ = postgres.NewScheduleRepo(db)
	app := adapterhttp.NewRouter(adapterhttp.Handlers{
		General:  handlers.NewGeneral(),
		Schedule: handlers.NewSchedule(),
	})

	return app.Listen(":" + cfg.HTTPPort)
}
