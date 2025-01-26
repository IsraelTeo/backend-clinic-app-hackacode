package model

import "time"

type Service struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"-"`
	Name        string    `json:"name" validate:"required,max=50"`
	Description string    `json:"description" validate:"required,max=250"`
	Price       float64   `json:"price" validate:"min=0,numeric"`
}
