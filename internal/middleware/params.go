package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ShouldHaveUrlParam(paramName string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		value := c.Params(paramName)
		if value == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":      false,
				"message": fmt.Sprintf("url should have param %s\n", paramName),
			})
		}

		return c.Next()
	}
}
