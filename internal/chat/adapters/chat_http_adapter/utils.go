package chathttpadapter

import "github.com/gofiber/fiber/v2"

func FormatErrorResponse(c *fiber.Ctx, statusCode int, message string, details interface{}) error {

	return c.Status(statusCode).JSON(fiber.Map{
		"error":   message,
		"details": details,
	})
}
