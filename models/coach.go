package models

import (
	"time"
)

// Coach maps to `coaches` table
type Coach struct {
	ID        uint       `gorm:"id" json:"id"`
	UUID      string     `gorm:"uuid"    json:"uuid"`
	FullName  string     `gorm:"full_name" json:"full_name"`
	Email     string     `gorm:"email" json:"email"`
	Status    string     `gorm:"status'" json:"status"`
	CreatedAt time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"deleted_at" json:"deleted_at"`
}

func (Coach) TableName() string {
	return "coaches"
}
