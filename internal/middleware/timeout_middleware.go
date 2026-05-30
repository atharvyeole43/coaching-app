package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	timeoutmiddleware "github.com/gofiber/fiber/v2/middleware/timeout"
)

func Timeout(timeout time.Duration) fiber.Handler {

	return timeoutmiddleware.New(func(c *fiber.Ctx) error {

		return c.Next()

	}, timeout)
}

func getLogFile() *os.File {
	currentDate := time.Now().Format("2006-01-02")
	logFileName := "log/buidco_accounting_request" + currentDate + ".log"
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	return logFile
}
