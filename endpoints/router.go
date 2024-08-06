package endpoints

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/x-sushant-x/RateShield/limiter"
)

func StartTestingRouter() {
	app := fiber.New()

	app.Use(checkRateLimit)

	app.Get("/", getReq)

	go func() {
		log.Fatal(app.Listen(":3000"))
	}()
}

func getReq(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"message": "success",
	})
}

func checkRateLimit(c *fiber.Ctx) error {
	if valid := limiter.RateLimiter.CheckLimit(c.IP()); !valid {
		return c.Status(http.StatusTooManyRequests).JSON(map[string]string{
			"error": "Too many requests",
		})
	}

	return c.Next()
}
