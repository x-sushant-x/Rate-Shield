package endpoints

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func StartTestingRouter() {
	app := fiber.New()

	app.Use(checkRateLimit)

	app.Get("/api/v1/get-data", getReq)

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
	apiEndpoint := "http://127.0.0.1:8080/rate-limiter/check"
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("ip", c.IP())
	req.Header.Add("endpoint", c.Path())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	code := resp.StatusCode

	switch code {
	case 200:
		return c.Next()
	case 429:
		return sendTooManyReq(c)
	case 500:
		return sendInternalServerError(c)
	}

	return c.Next()
}

func sendTooManyReq(c *fiber.Ctx) error {
	return c.Status(http.StatusTooManyRequests).JSON(map[string]string{
		"error": "Too many requests",
	})
}

func sendInternalServerError(c *fiber.Ctx) error {
	return c.Status(http.StatusInternalServerError).JSON(map[string]string{
		"error": "Internal server error",
	})
}
