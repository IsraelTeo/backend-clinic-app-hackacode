package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"size:50;not null"`
	Password string `gorm:"size:50;not null"`
}

/*type UserRegister struct {
	Email    string `json:"email" gorm:"size:50;not null" validate:"required, email"`
	Password string `json:"password" gorm:"size:100;not null" validate:"required"`
	IsAdmin  bool   `json:"is_admin"`
}*/

type UserLogin struct {
	Email    string `json:"email" gorm:"size:50;not null" validate:"required, email"`
	Password string `json:"password" gorm:"size:100;not null" validate:"required"`
}
