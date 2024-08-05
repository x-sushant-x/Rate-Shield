package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", getHome)
	log.Fatal(app.Listen(":3000"))
}

func getHome(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"status": "successful",
	})
}
