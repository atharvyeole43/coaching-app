package middleware

import (
	"coaching-app-backend/utils"
	"context"

	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		// Use a goroutine to execute the main handler
		done := make(chan bool)
		go func() {
			c.Next()
			done <- true
		}()
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				utils.GatewayTimeoutAbortWithJSON(c, "Request Timeout")
			}
		case <-done:
			// Request completed before timeout
		}
	}
}

// Middleware function to check if the request is coming from Postman
func PostmanCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		choiceKey := c.GetHeader("x-choice-key")
		if strings.Contains(strings.ToLower(userAgent), "postman") && choiceKey == "" {
			utils.UnauthorizedAbortWithJSON(c, "Unauthorized access from Postman")
			return
		}
		c.Next()
	}
}

func GetBoolEnv(key string, defaultVal bool) bool {
	val := strings.ToLower(os.Getenv(key))
	if val == "true" {
		return true
	} else if val == "false" {
		return false
	}
	return defaultVal
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
