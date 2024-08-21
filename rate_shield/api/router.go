package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/service"
)

func StartServer() {
	app := fiber.New()

	rateLimitHandler := RateLimitHandler{}

	app.Get("/rate-limiter/check", rateLimitHandler.CheckRateLimit)

	rulesGroup := app.Group("/rules")
	{
		svc := service.RulesServiceRedis{}
		h := NewRulesAPIHandler(svc)
		rulesGroup.Get("/all", h.ListAllRules)

		rulesGroup.Post("/add", h.CreateOrUpdateRule)
		rulesGroup.Delete("/delete", h.DeleteRule)
	}

	go func() {
		err := app.Listen(":8080")
		if err != nil {
			log.Err(err).Msg("unable to start server")
		}
	}()
}
