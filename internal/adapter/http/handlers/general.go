package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type General struct{}

func NewGeneral() *General {
	return &General{}
}

func (h *General) Index(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
