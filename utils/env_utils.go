package utils

import (
	"os"
)

// Load environment variables
func MustGetenv(key string) string {
	val := os.Getenv(key)
	// if val == "" {
	// 	fmt.Println("Environment variable is required but not set", key)
	// 	logrus.Fatal(fmt.Sprintf("Environment variable %s is required but not set", key))
	// }
	return val
}
