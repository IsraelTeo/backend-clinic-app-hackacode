package model

import "time"

// MedicalService represents a medical service in the system.
type MedicalService struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	Name        string  `gorm:"size:50;not null"`
	Description string  `gorm:"size:250;not null"`
	Price       float64 `gorm:"not null"`
	CreatedAt   time.Time
}
