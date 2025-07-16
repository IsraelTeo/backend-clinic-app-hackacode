package model

import "time"

// Patient represents a patient in the system.
type Patient struct {
	Person
	Insurance bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}
