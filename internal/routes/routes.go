package routes

import (
	"coaching-app-backend/app"
	"coaching-app-backend/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *app.Controllers) {

	// Apply CORS middleware globally
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	middleware.SetupRateLimiter(r)

	api := r.Group("/api/v1")
	{
		income := api.Group("/income")
		{
			income.GET("/daily", middleware.Timeout(30), c.IncomeControllers.IncomeController.GetDailyIncome)
		}
	}

}	
