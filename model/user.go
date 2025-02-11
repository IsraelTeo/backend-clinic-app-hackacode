package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"size:50;not null"`
	Password string `gorm:"size:50;not null"`
}


