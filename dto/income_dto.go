package dto

import "time"

type GetDailyIncomeRequest struct {
	CoachID string `form:"coach_id" validate:"required"`
	Period  string `form:"period"`
}

type IncomePoint struct {
	Date   CustomDate `json:"date"`
	Amount float64    `json:"amount"`
}

type Summary struct {
	TotalIncomeCurrent  float64  `json:"total_income_current"`
	TotalIncomePrevious float64  `json:"total_income_previous"`
	DeltaAmount         float64  `json:"delta_amount"`
	DeltaPercent        *float64 `json:"delta_percent"`
	Trend               string   `json:"trend"`
}

type DateRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Meta struct {
	CurrentRange  DateRange `json:"current_range"`
	PreviousRange DateRange `json:"previous_range"`
}

type ChartData struct {
	Current  []IncomePoint `json:"current"`
	Previous []IncomePoint `json:"previous"`
}

type DailyIncomeResponse struct {
	Summary Summary   `json:"summary"`
	Chart   ChartData `json:"chart"`
	Meta    Meta      `json:"meta"`
}
type CustomDate time.Time

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	t := time.Time(cd)
	formatted := t.Format("2006-01-02")
	return []byte(`"` + formatted + `"`), nil
}
