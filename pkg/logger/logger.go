package logger

import (
	"coaching-app-backend/constant"
	"fmt"
	"io"

	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Initialize the logger
func Init() {
	// Prepare log file name based on the current date
	now := time.Now()
	logFileName := now.Format(constant.DateFormat) + ".log"

	// Open log file for writing (create it if doesn't exist)
	logDir := path.Join("./storage/logs")
	err := os.MkdirAll(logDir, os.ModePerm) // Create directories if they don't exist
	if err != nil {
		logrus.Error(fmt.Sprintf("Error creating log directory: %v", err))
		return
	}

	// Open log file for writing (create it if doesn't exist)
	logFilePath := path.Join(logDir, logFileName)
	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logrus.Error(fmt.Sprintf("Error opening log file: %v", err))
		return
	}

	// Set log output to the log file
	logrus.SetOutput(file)

	// Configure logrus formatter
	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableHTMLEscape: true,
		PrettyPrint:       true,
		TimestampFormat:   constant.DateFormat + " 15:04:05",
	})

	// Set gin's default writer to also write to the log file
	gin.DefaultWriter = io.MultiWriter(file)

	// Set log level and report caller
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)

	// Log the file name for debugging purposes
	logrus.Debug("Log file: ", logFileName)
}
