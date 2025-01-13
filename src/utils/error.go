package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HandleError(c *fiber.Ctx, err error) error {
	if err != nil && (strings.Contains(err.Error(), "invalid response format") ||
		strings.Contains(err.Error(), "invalid token format")) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Session expired",
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}
