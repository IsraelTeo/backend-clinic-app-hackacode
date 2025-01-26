package repository

import (
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	GetByID(ID uint) (*model.Service, error)
	GetAll() ([]model.Service, error)
	Create(service *model.Service) error
	Update(service *model.Service) (*model.Service, error)
	Delete(ID uint) error
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) GetByID(ID uint) (*model.Service, error) {
	service := model.Service{}
	if err := r.db.First(&service, ID).Error; err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *serviceRepository) GetAll() ([]model.Service, error) {
	var services []model.Service
	if err := r.db.Find(&services).Error; err != nil {
		return nil, err
	}

	return services, nil
}

func (r *serviceRepository) Create(service *model.Service) error {
	if err := r.db.Create(service).Error; err != nil {
		return err
	}

	return nil
}

func (r *serviceRepository) Update(service *model.Service) (*model.Service, error) {
	if err := r.db.Save(service).Error; err != nil {
		return nil, err
	}

	return service, nil
}
func (r *serviceRepository) Delete(ID uint) error {
	if err := r.db.Delete(&model.Service{}, ID).Error; err != nil {
		return err
	}

	return nil
}
