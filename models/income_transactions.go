package models

import (
	"time"
)

// CoachIncome maps to `income_transactions` table
type IncomeTransactions struct {
	ID              uint       `gorm:"id" json:"id"`
	UUID            string     `gorm:"uuid"  json:"uuid"`
	CoachID         uint       `gorm:"coach_id"           json:"coach_id"`
	Amount          float64    `gorm:"amount" json:"amount"`
	Status          string     `gorm:"status'" json:"status"`
	TransactionDate time.Time  `gorm:"transaction_date"       json:"transaction_date"`
	CreatedAt       time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;type:timestamp;default:NULL" json:"deleted_at"`
}

func (IncomeTransactions) TableName() string {
	return "income_transactions"
}
