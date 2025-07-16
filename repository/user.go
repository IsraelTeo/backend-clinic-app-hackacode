package repository

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user operations
type UserRepository interface {
	GetByID(ID uint) (*model.User, error)
	Register(user *model.User) (*model.User, error)
	ExistsByEmail(email string) (bool, error)
}

// userRepository implements the UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// Dependency injection for the user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// GetByID retrieves a user by its ID
func (r *userRepository) GetByID(ID uint) (*model.User, error) {
	var user model.User

	if err := r.db.First(&user, ID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Register creates a new user in the database
func (r *userRepository) Register(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// ExistsByEmail checks if a user exists by their email
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64

	err := r.db.Model(&model.User{}).
		Where("email = ?", email).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
