package model

import "time"

// Doctor represents a doctor in the system.
type Doctor struct {
	Person
	Especialty string    `json:"especialty" validate:"required,max=50"`
	Days       string    `json:"days" validate:"required"`
	StartTime  time.Time `json:"start_time" validate:"required"`
	EndTime    time.Time `json:"end_time" validate:"required"`
	Salary     float64   `json:"salary" validate:"required,numeric"`
	CreatedAt  time.Time `gorm:"not null"`
}
