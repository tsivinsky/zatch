package middleware

import "github.com/gofiber/fiber/v2"

func ShouldHaveShortIdParam(c *fiber.Ctx) error {
	shortId := c.Params("shortId")
	if shortId == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":      false,
			"message": "shortId parameter required",
		})
	}

	return c.Next()
}
