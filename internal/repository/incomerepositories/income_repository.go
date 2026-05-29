package incomerepositories

import (
	"coaching-app-backend/dto"
	database "coaching-app-backend/internal/storage/db"
	"coaching-app-backend/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IncomeRepository interface {
	GetDailyIncome(req dto.GetDailyIncomeRequest) (*dto.DailyIncomeResponse, error)
}

type incomeRepository struct {
	db *gorm.DB
}

func NewIncomeRepository(db *gorm.DB) IncomeRepository {

	return &incomeRepository{
		db: db,
	}
}

func (r *incomeRepository) GetDailyIncome(req dto.GetDailyIncomeRequest) (*dto.DailyIncomeResponse, error) {

	var response dto.DailyIncomeResponse

	coachID, err := strconv.Atoi(req.CoachID)
	if err != nil {
		logrus.Errorf("GetDailyIncome@InvalidCoachID=%v", err)
		return nil, fmt.Errorf("invalid coach_id")
	}

	var coachCount int64
	err = database.COACHINGDB.
		Model(&models.Coach{}).
		Where("id = ? AND deleted_at IS NULL", coachID).
		Count(&coachCount).Error
	if err != nil {
		logrus.Errorf("GetDailyIncome@CoachCheckError=%v", err)
		return nil, err
	}
	if coachCount == 0 {
		return nil, fmt.Errorf("coach not found")
	}

	days := 7
	if req.Period != "" && req.Period != "7d" {
		periodStr := strings.TrimSuffix(req.Period, "d")
		if parsed, parseErr := strconv.Atoi(periodStr); parseErr == nil && parsed > 0 {
			days = parsed
		}
	}

	today := time.Now()
	currentTo := today
	currentFrom := today.AddDate(0, 0, -(days - 1))

	previousTo := currentFrom.AddDate(0, 0, -1)
	previousFrom := previousTo.AddDate(0, 0, -(days - 1))

	var currentRows []dto.IncomePoint
	err = database.COACHINGDB.
		Model(&models.IncomeTransactions{}).
		Select("DATE(transaction_date) as date, SUM(amount) as amount").
		Where("coach_id = ?", coachID).
		Where("status = ?", "completed").
		Where("DATE(transaction_date) BETWEEN ? AND ?",
			currentFrom.Format("2006-01-02"),
			currentTo.Format("2006-01-02"),
		).
		Group("DATE(transaction_date)").
		Order("DATE(transaction_date) ASC").
		Scan(&currentRows).Error
	if err != nil {
		logrus.Errorf("GetDailyIncome@CurrentDataError=%v", err)
		return nil, err
	}

	var previousRows []dto.IncomePoint
	err = database.COACHINGDB.
		Model(&models.IncomeTransactions{}).
		Select("DATE(transaction_date) as date, SUM(amount) as amount").
		Where("coach_id = ?", coachID).
		Where("status = ?", "completed").
		Where("DATE(transaction_date) BETWEEN ? AND ?",
			previousFrom.Format("2006-01-02"),
			previousTo.Format("2006-01-02"),
		).
		Group("DATE(transaction_date)").
		Order("DATE(transaction_date) ASC").
		Scan(&previousRows).Error
	if err != nil {
		logrus.Errorf("GetDailyIncome@PreviousDataError=%v", err)
		return nil, err
	}

	var totalCurrent, totalPrevious float64
	for _, item := range currentRows {
		totalCurrent += item.Amount
	}
	for _, item := range previousRows {
		totalPrevious += item.Amount
	}

	deltaAmount := totalCurrent - totalPrevious

	var deltaPercent *float64
	if totalPrevious > 0 {
		dp := (deltaAmount / totalPrevious) * 100
		deltaPercent = &dp
	}

	trend := "unchanged"
	if deltaAmount > 0 {
		trend = "increased"
	} else if deltaAmount < 0 {
		trend = "decreased"
	}

	response.Summary = dto.Summary{
		TotalIncomeCurrent:  totalCurrent,
		TotalIncomePrevious: totalPrevious,
		DeltaAmount:         deltaAmount,
		DeltaPercent:        deltaPercent,
		Trend:               trend,
	}

	if currentRows == nil {
		currentRows = []dto.IncomePoint{}
	}
	if previousRows == nil {
		previousRows = []dto.IncomePoint{}
	}
	response.Chart.Current = currentRows
	response.Chart.Previous = previousRows

	response.Meta = dto.Meta{
		CurrentRange: dto.DateRange{
			From: currentFrom.Format("2006-01-02"),
			To:   currentTo.Format("2006-01-02"),
		},
		PreviousRange: dto.DateRange{
			From: previousFrom.Format("2006-01-02"),
			To:   previousTo.Format("2006-01-02"),
		},
	}

	return &response, nil
}
