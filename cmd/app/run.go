package app

import (
	application "coaching-app-backend/app"
	"coaching-app-backend/config"
	"coaching-app-backend/internal/routes"
	dbstore "coaching-app-backend/internal/storage/db"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
)

func Run() {

	fmt.Println("Before db connections Initialized")
	logrus.Info("Before db connections Initialized")

	// DB Connections
	dbstore.InitAllDbConnections()

	fmt.Println("All DB Connections Initialized")
	logrus.Info("All DB Connections Initialized")

	// Load ENV
	config.LoadEnvVariables()

	// Initialize Controllers
	controllers := application.InitApp()

	// Create Fiber App
	r := fiber.New()

	// Middleware
	r.Use(logger.New())
	r.Use(recover.New())
	r.Use(cors.New())

	// Register Routes
	routes.RegisterRoutes(r, controllers.AppController)

	fmt.Println("Coaching App Backend Locked and Loaded")
	logrus.Info("Coaching App Backend Locked and Loaded")

	port := os.Getenv("APP_PORT")

	logrus.Infof("Starting HTTP server on port %s", port)
	log.Print("Starting HTTP server on port: ", port)

	// Start Server
	if err := r.Listen(":" + port); err != nil {
		logrus.Fatal("Error starting the server: ", err)
	}
}
