package middleware

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RateLimiter() fiber.Handler {
	fmt.Println("Here")
	return func(c *fiber.Ctx) error {

		ip := c.IP()         // Client IP
		endpoint := c.Path() // Requested endpoint

		req, err := http.NewRequest("GET", "http://127.0.0.1:8080/check-limit", nil)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		// Set headers for Rate Shield
		req.Header.Set("ip", ip)
		req.Header.Set("endpoint", endpoint)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		defer resp.Body.Close()

		// Handle Rate Shield response
		switch resp.StatusCode {
		case http.StatusOK:
			return c.Next() // Allow request to proceed
		case http.StatusTooManyRequests:
			return c.Status(fiber.StatusTooManyRequests).SendString("Rate limit exceeded")
		default:
			return fiber.ErrInternalServerError
		}

	}
}
