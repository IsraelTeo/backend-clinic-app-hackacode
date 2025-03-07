package repository

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ID uint) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	user := model.User{}

	err := r.db.
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByID(ID uint) (*model.User, error) {
	user := model.User{}

	err := r.db.First(user, ID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
