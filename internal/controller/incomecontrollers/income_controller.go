package incomecontrollers

import (
	"coaching-app-backend/dto"
	"coaching-app-backend/internal/service/incomeservices"
	"coaching-app-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type IncomeController interface {
	GetDailyIncome(c *gin.Context)
}

type incomeController struct {
	incomeService incomeservices.IncomeService
}

func NewIncomeController(incomeService incomeservices.IncomeService) IncomeController {

	return &incomeController{
		incomeService: incomeService,
	}
}
func (ic *incomeController) GetDailyIncome(c *gin.Context) {

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("GetDailyIncome@panicInfo:", panicInfo)
			utils.InternalServerErrorWithMessage(c, "Unexpected error occurred")
			return
		}
	}()

	var request dto.GetDailyIncomeRequest
	request.CoachID = c.Query("coach_id")
	request.Period = c.DefaultQuery("period", "7d")

	validationResp := utils.ValidateRequest(c, request)
	if validationResp != nil {
		utils.ValidationResponse(c, validationResp.(string))
		return
	}

	response, err := ic.incomeService.GetDailyIncome(request)
	if err != nil {
		logrus.Error("GetDailyIncome@Error:", err)

		if err.Error() == "coach not found" {
			utils.NotFoundResponse(c, "Coach not found")
			return
		}

		if err.Error() == "invalid coach_id" {
			utils.ValidationResponse(c, "coach_id must be a valid number")
			return
		}

		utils.InternalServerErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessResponse(c, "Daily income fetched successfully", response)
}
