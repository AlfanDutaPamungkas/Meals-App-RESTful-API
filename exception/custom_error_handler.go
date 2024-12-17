package exception

import "github.com/gofiber/fiber/v2"

func ErrorHandler(statusCode int, message string, err error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(statusCode).JSON(fiber.Map{
			"code":   statusCode,
			"status": message,
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}
}
