package incomeservices

import (
	"coaching-app-backend/dto"
	"coaching-app-backend/internal/repository/incomerepositories"

	"github.com/sirupsen/logrus"
)

type IncomeService interface {
	GetDailyIncome(req dto.GetDailyIncomeRequest) (*dto.DailyIncomeResponse, error)
}

type incomeService struct {
	incomeRepo incomerepositories.IncomeRepository
}

func NewIncomeService(incomeRepo incomerepositories.IncomeRepository) IncomeService {

	return &incomeService{
		incomeRepo: incomeRepo,
	}
}
func (s *incomeService) GetDailyIncome(request dto.GetDailyIncomeRequest) (*dto.DailyIncomeResponse, error) {

	response, err := s.incomeRepo.GetDailyIncome(request)
	if err != nil {

		logrus.Errorf("GetDailyIncome@ServiceError=%v", err)

		return nil, err
	}

	return response, nil
}
