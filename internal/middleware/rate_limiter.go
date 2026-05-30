package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupRateLimiter(app *fiber.App) {

	app.Use(limiter.New(limiter.Config{

		// Max requests
		Max: 30,

		// Expiration window
		Expiration: 1 * time.Second,

		// Response when limit exceeded
		LimitReached: func(c *fiber.Ctx) error {

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Too many requests",
			})
		},
	}))
}
