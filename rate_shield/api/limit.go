package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x-sushant-x/RateShield/limiter"
)

type RateLimitHandler struct{}

func (h RateLimitHandler) CheckRateLimit(c *fiber.Ctx) error {
	ip := c.Get("ip")
	endpoint := c.Get("endpoint")

	code := limiter.RateLimiter.CheckLimit(ip, endpoint)

	switch code {
	case 200:
		return c.SendStatus(200)
	case 429:
		return c.SendStatus(429)
	default:
		return c.SendStatus(500)
	}
}
