package incomeservices

import "coaching-app-backend/internal/repository/incomerepositories"

type IncomeServices struct {
	IncomeService IncomeService
}

func NewIncomeServices(
	incomeRepos incomerepositories.IncomeRepositories,
) IncomeServices {

	return IncomeServices{
		IncomeService: NewIncomeService(
			incomeRepos.IncomeRepository,
		),
	}
}
