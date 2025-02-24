package repository

import (
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	GetAll() ([]model.Appointment, error)
	Delete(ID uint) error
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) GetAll() ([]model.Appointment, error) {
	var appointments []model.Appointment

	err := r.db.Find(&appointments).Error
	if err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *serviceRepository) GetAllServicesByID(ID []uint) ([]model.Service, error) {
	var services []model.Service

	err := r.db.
		Where("id IN ?", ID).
		Find(&services).
		Error
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (r *serviceRepository) Delete(ID uint) error {
	err := r.db.Exec("DELETE FROM package_services WHERE service_id = ?", ID).Error
	if err != nil {
		return err
	}

	return r.db.Delete(&model.Service{}, ID).Error
}
