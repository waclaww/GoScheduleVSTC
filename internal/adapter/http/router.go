package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func NewRouter(h Handlers) *fiber.App {
	app := fiber.New(fiber.Config{})

	app.Use(recover.New())
	app.Use(requestid.New())

	app.Get("/", h.General.Index)

	return app
}
