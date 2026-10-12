package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type Schedule struct{}

func NewSchedule() *Schedule {
	return &Schedule{}
}

func (h *Schedule) Index(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
