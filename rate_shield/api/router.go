package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func StartServer() {
	app := fiber.New()

	rateLimitHandler := RateLimitHandler{}

	app.Get("/rate-limiter/check", rateLimitHandler.CheckRateLimit)

	go func() {
		err := app.Listen(":8080")
		if err != nil {
			log.Err(err).Msg("unable to start server")
		}
	}()
}
