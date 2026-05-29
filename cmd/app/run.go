package app

import (
	"coaching-app-backend/app"
	"coaching-app-backend/config"
	"coaching-app-backend/internal/routes"
	dbstore "coaching-app-backend/internal/storage/db"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run() {

	fmt.Println("Before db connections Initiallized")
	logrus.Info("Before db connections Initiallized")

	dbstore.InitAllDbConnections()
	fmt.Println("All Db Connections Initiallized")
	logrus.Info("All Db Connections Initiallized")
	config.LoadEnvVariables()
	//utils.InitEmailWorkerPool()

	controllers := app.InitApp()

	r := gin.Default()
	routes.RegisterRoutes(r, controllers.AppController)
	fmt.Println("coaching app Backend Lock and Loaded")
	logrus.Info("coaching app Backend Lock and Loaded")

	port := os.Getenv("APP_PORT")
	logrus.Infof("Starting HTTP server on port %s", port)
	log.Print("Starting HTTP server on port: ", port)

	if err := r.Run(":" + port); err != nil {
		logrus.Fatal("Error starting the server: ", err)
	}

}
