package main

import (
	"coaching-app-backend/cmd/app"
	"coaching-app-backend/pkg/logger"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		logrus.Error("Error loading .env file")
	}
	fmt.Println("test")

	// Initialize logger
	logger.Init()

	// Run the application
	app.Run()
}
