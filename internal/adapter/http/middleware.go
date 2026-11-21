package http

import (
	"crypto/subtle"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func bearerAuth(token string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if token == "" {
			return fiber.NewError(fiber.StatusForbidden, "write API is disabled")
		}

		got, ok := strings.CutPrefix(c.Get(fiber.HeaderAuthorization), "Bearer ")
		if !ok || subtle.ConstantTimeCompare([]byte(got), []byte(token)) != 1 {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		return c.Next()
	}
}
