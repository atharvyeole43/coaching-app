package dbstore

import (
	"coaching-app-backend/config"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
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
		logrus.Error("COACHING  DB connections is not initialized.")
	}

	return DataBaseConfig{
		COACHINGDB: COACHINGDB,
	}
}
func DBConnectionPool() {

	COACHINGDbConfig, err := config.COACHINGDBConfig()
	if err != nil {
		logrus.Fatal("Failed to load Insurance DB configuration: ", err)
		return
	}

	dbTcp := fmt.Sprintf("@tcp(%s:%s)/", COACHINGDbConfig.Host, COACHINGDbConfig.Port)
	dsn := fmt.Sprintf(
		"%s:%s%s%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FKolkata",
		COACHINGDbConfig.User,
		COACHINGDbConfig.Password,
		dbTcp,
		COACHINGDbConfig.Name,
	)

	// Open a DB connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("COACHINGDbConnectionPool Error: ", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatal("Failed to get sql.DB from gorm.DB: ", err)
		return
	}

	sqlDB.SetMaxIdleConns(COACHINGDbConfig.MaxIdleConn)
	sqlDB.SetMaxOpenConns(COACHINGDbConfig.MaxOpenConn)
	sqlDB.SetConnMaxIdleTime(COACHINGDbConfig.MaxIdleTime)
	sqlDB.SetConnMaxLifetime(COACHINGDbConfig.MaxLifeTime)

	// Log the connection pool stats
	stats := sqlDB.Stats()
	logrus.Infof("@COACHINGBPool MYSQL Max Open Connections: %d", stats.MaxOpenConnections)
	logrus.Infof("@COACHINGDBPool MYSQL Open Connections: %d", stats.OpenConnections)
	logrus.Infof("@COACHINGDBPool MYSQL InUse Connections: %d", stats.InUse)
	logrus.Infof("@COACHINGDBPool MYSQL Idle Connections: %d", stats.Idle)

	COACHINGDB = db
	logrus.Info("DB Connection Initiated: ")
}
