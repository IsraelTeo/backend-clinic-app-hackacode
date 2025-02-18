package repository

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	GetAllServicesByID(ID []uint) ([]model.Service, error)
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) GetAllServicesByID(ID []uint) ([]model.Service, error) {
	var services []model.Service
	if err := r.db.Where("id IN ?", ID).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}
