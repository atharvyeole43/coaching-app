package incomecontrollers

import (
	"coaching-app-backend/dto"
	"coaching-app-backend/internal/service/incomeservices"
	"coaching-app-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type IncomeController interface {
	GetDailyIncome(c *fiber.Ctx) error
}

type incomeController struct {
	incomeService incomeservices.IncomeService
}

func NewIncomeController(incomeService incomeservices.IncomeService) IncomeController {
	return &incomeController{
		incomeService: incomeService,
	}
}

func (ic *incomeController) GetDailyIncome(c *fiber.Ctx) error {

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.WithField("panic_info", panicInfo).Error("GetDailyIncome panic recovered")
			_ = utils.InternalServerErrorWithMessage(c, "Unexpected error occurred")
		}
	}()

	request := dto.GetDailyIncomeRequest{
		CoachID: c.Query("coach_id"),
		Period:  c.Query("period", "7d"),
	}

	if validationResp := utils.ValidateRequest(c, request); validationResp != nil {

		return utils.ValidationResponse(c, validationResp.(string))
	}

	response, err := ic.incomeService.GetDailyIncome(request)
	if err != nil {

		logrus.WithError(err).Error("GetDailyIncome failed")

		switch err.Error() {

		case "coach not found":
			return utils.NotFoundResponse(c, "Coach not found")

		case "invalid coach_id":
			return utils.ValidationResponse(c, "coach_id must be a valid number")

		default:
			return utils.InternalServerErrorWithMessage(c, "Failed to fetch daily income")
		}
	}

	return utils.SuccessResponse(c, "Daily income fetched successfully", response)
}
