package app

import (
	"coaching-app-backend/internal/controller/incomecontrollers"
	"coaching-app-backend/internal/repository/incomerepositories"
	"coaching-app-backend/internal/service/incomeservices"
)

type App struct {
	AppController *Controllers
}

type Controllers struct {
	IncomeControllers incomecontrollers.Controllers
}

type Services struct {
	IncomeServices incomeservices.IncomeServices
}

type Repositories struct {
	IncomeRepositories incomerepositories.IncomeRepositories
}

func InitApp() *App {

	repos := initRepositories()

	services := initServices(repos)

	controllers := initControllers(services)

	return &App{
		AppController: controllers,
	}
}

func initRepositories() *Repositories {

	return &Repositories{IncomeRepositories: incomerepositories.NewIncomeRepositories()}
}

func initServices(repos *Repositories) *Services {

	return &Services{

		IncomeServices: incomeservices.NewIncomeServices(repos.IncomeRepositories),
	}
}

func initControllers(services *Services) *Controllers {

	return &Controllers{

		IncomeControllers: incomecontrollers.NewControllers(services.IncomeServices.IncomeService),
	}
}
