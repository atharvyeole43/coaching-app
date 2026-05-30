package dbstore

import (
	"coaching-app-backend/config"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	COACHINGDB *gorm.DB
)

type DataBaseConfig struct {
	COACHINGDB *gorm.DB
}

func InitAllDbConnections() {
	DBConnectionPool()
}

func GetDatabaseConfig() DataBaseConfig {

	if COACHINGDB == nil {
		logrus.Error("COACHING DB connection is not initialized.")
	}

	return DataBaseConfig{
		COACHINGDB: COACHINGDB,
	}
}

func DBConnectionPool() {

	COACHINGDbConfig, err := config.COACHINGDBConfig()
	if err != nil {
		logrus.Fatal("Failed to load DB configuration: ", err)
		return
	}

	// PostgreSQL DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		COACHINGDbConfig.Host,
		COACHINGDbConfig.User,
		COACHINGDbConfig.Password,
		COACHINGDbConfig.Name,
		COACHINGDbConfig.Port,
	)

	// Open PostgreSQL connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("PostgreSQL Connection Error: ", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatal("Failed to get sql.DB from gorm.DB: ", err)
		return
	}

	// Connection Pool Settings
	sqlDB.SetMaxIdleConns(COACHINGDbConfig.MaxIdleConn)
	sqlDB.SetMaxOpenConns(COACHINGDbConfig.MaxOpenConn)
	sqlDB.SetConnMaxIdleTime(COACHINGDbConfig.MaxIdleTime)
	sqlDB.SetConnMaxLifetime(COACHINGDbConfig.MaxLifeTime)

	// Stats
	stats := sqlDB.Stats()

	logrus.Infof("@COACHINGDBPool POSTGRES Max Open Connections: %d", stats.MaxOpenConnections)
	logrus.Infof("@COACHINGDBPool POSTGRES Open Connections: %d", stats.OpenConnections)
	logrus.Infof("@COACHINGDBPool POSTGRES InUse Connections: %d", stats.InUse)
	logrus.Infof("@COACHINGDBPool POSTGRES Idle Connections: %d", stats.Idle)

	COACHINGDB = db

	logrus.Info("PostgreSQL DB Connection Initiated")
}
