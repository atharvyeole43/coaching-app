package routes

import (
	"coaching-app-backend/app"
	"coaching-app-backend/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	fibercors "github.com/gofiber/fiber/v2/middleware/cors"
)

func RegisterRoutes(r *fiber.App, c *app.Controllers) {

	// CORS Middleware
	r.Use(fibercors.New(fibercors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "*",
		AllowCredentials: false,
		MaxAge:           int((12 * time.Hour).Seconds()),
	}))

	// Rate Limiter Middleware
	middleware.SetupRateLimiter(r)

	// API Version Group
	api := r.Group("/api/v1")

	// Income Routes
	income := api.Group("/income")

	income.Get("/daily", middleware.Timeout(30*time.Second), c.IncomeControllers.IncomeController.GetDailyIncome)
}
