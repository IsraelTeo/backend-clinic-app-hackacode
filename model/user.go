package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserName  string    `gorm:"size:50;not null"`
	Email     string    `gorm:"size:50;not null"`
	Password  string    `gorm:"size:50;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
