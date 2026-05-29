package config

import (
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type DBConfig struct {
	Driver      string
	Name        string
	User        string
	Password    string
	Host        string
	Port        string
	MaxIdleConn int
	MaxOpenConn int
	MaxIdleTime time.Duration
	MaxLifeTime time.Duration
}

func loadPoolConfig(cfg *DBConfig, prefix string) {
	var err error
	if cfg.MaxIdleConn, err = strconv.Atoi(os.Getenv(prefix + "_DB_MAX_IDLE_CONN")); err != nil {
		logrus.Fatalf("%s Invalid DB_MAX_IDLE_CONN: %v", prefix, err)
	}
	if cfg.MaxOpenConn, err = strconv.Atoi(os.Getenv(prefix + "_DB_MAX_OPEN_CONN")); err != nil {
		logrus.Fatalf("%s Invalid DB_MAX_OPEN_CONN: %v", prefix, err)
	}
	if maxIdleTime, err := strconv.Atoi(os.Getenv(prefix + "_DB_MAX_IDLE_TIME")); err == nil {
		cfg.MaxIdleTime = time.Duration(maxIdleTime) * time.Second
	} else {
		logrus.Fatalf("%s Invalid DB_MAX_IDLE_TIME: %v", prefix, err)
	}
	if maxLifeTime, err := strconv.Atoi(os.Getenv(prefix + "_DB_MAX_LIFE_TIME")); err == nil {
		cfg.MaxLifeTime = time.Duration(maxLifeTime) * time.Second
	} else {
		logrus.Fatalf("%s Invalid DB_MAX_LIFE_TIME: %v", prefix, err)
	}
}

func COACHINGDBConfig() (*DBConfig, error) {
	cfg := &DBConfig{
		Driver:   os.Getenv("DB_DRIVER"),
		Name:     os.Getenv("DB_DATABASE_NAME"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
	}
	loadPoolConfig(cfg, "COACHING_APP")
	return cfg, nil
}
