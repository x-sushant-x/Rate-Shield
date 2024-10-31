package main

import (
	"gofiberapp/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(middleware.RateLimiter())

	app.Get("/api/v1/resource", func(c *fiber.Ctx) error {
		return c.SendString("This is a protected resource")
	})

	app.Listen(":3000")
}
