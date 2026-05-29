package incomerepositories

import dbstore "coaching-app-backend/internal/storage/db"

type IncomeRepositories struct {
	IncomeRepository IncomeRepository
}

func NewIncomeRepositories() IncomeRepositories {

	return IncomeRepositories{
		IncomeRepository: NewIncomeRepository(
			dbstore.COACHINGDB,
		),
	}
}
