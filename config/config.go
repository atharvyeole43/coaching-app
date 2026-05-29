package config

import (
	"fmt"
	"sync"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
}

var (
	instance *EnvConfig
	mu       sync.RWMutex
)

// LoadEnvVariables initializes the config singleton (write lock)
func LoadEnvVariables() {
	mu.Lock()
	defer mu.Unlock()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}
}

// Get returns the singleton config instance (read lock)
func GetEnv() *EnvConfig {
	mu.RLock()
	defer mu.RUnlock()
	if instance == nil {
		panic("Config not loaded. Call config.LoadEnvVariables() in main.")
	}
	return instance
}
