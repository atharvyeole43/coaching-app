package incomecontrollers

import "coaching-app-backend/internal/service/incomeservices"

type Controllers struct {
	IncomeController IncomeController
}

func NewControllers(
	incomeService incomeservices.IncomeService,
) Controllers {

	return Controllers{
		IncomeController: NewIncomeController(
			incomeService,
		),
	}
}
