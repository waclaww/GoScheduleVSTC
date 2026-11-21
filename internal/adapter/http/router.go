package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func NewRouter(h Handlers, apiToken string) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Use(recover.New())
	app.Use(requestid.New())

	app.Get("/", h.General.Index)
	app.Get("/health", h.General.Health)

	schedule := app.Group("/schedule")
	schedule.Get("/vstechnical_get", h.Schedule.List)
	schedule.Post("/vstechical_insert", h.Schedule.Create)
	schedule.Post("/vstechical_update", h.Schedule.Update)
	schedule.Post("/vstechnical_group", h.Schedule.ByGroup)
	schedule.Get("/vstechnical_archive_get", h.Archive.List)
	schedule.Post("/vstechnical_archive_insert", h.Archive.Create)
	schedule.Post("/vstechnical_archive_get_date", h.Archive.ByDate)

	v1 := app.Group("/api/v1")
	v1.Get("/groups", h.Day.Groups)
	v1.Get("/groups/:name/schedule", h.Day.GroupSchedule)
	v1.Get("/days/:date", h.Day.Get)
	v1.Put("/days/:date", bearerAuth(apiToken), h.Day.Upsert)
	v1.Get("/pages/:name", h.Page.Get)
	v1.Put("/pages/:name", bearerAuth(apiToken), h.Page.Set)

	return app
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		msg = e.Message
	}

	return c.Status(code).JSON(fiber.Map{"detail": msg})
}
