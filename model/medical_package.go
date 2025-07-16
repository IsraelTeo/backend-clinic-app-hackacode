package model

import "time"

// MedicalPackage represents a medical package in the system.
type MedicalPackage struct {
	ID        uint             `gorm:"primaryKey;autoIncrement"`
	Name      string           `gorm:"size:50;not null"`
	Services  []MedicalService `gorm:"many2many:package_services"`
	Price     float64          `gorm:"not null"`
	CreatedAt time.Time        `gorm:"not null"`
}
