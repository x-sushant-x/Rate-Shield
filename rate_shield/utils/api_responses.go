package utils

import "github.com/gofiber/fiber/v2"

func SendInternalError(c *fiber.Ctx) error {
	return c.Status(500).JSON(map[string]string{
		"status": "fail",
		"error":  "Internal Server Error",
	})
}

func SendBadRequestError(c *fiber.Ctx) error {
	return c.Status(500).JSON(map[string]string{
		"status": "fail",
		"error":  "Invalid Request Body",
	})
}
