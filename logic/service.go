package logic

import (
	"fmt"
	"log"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
)

type ServiceLogic interface {
	GetServiceByID(ID uint) (*model.Service, error)
	GetAllServices() ([]model.Service, error)
	CreateService(service *model.Service) error
	UpdateService(ID uint, service *model.Service) (*model.Service, error)
	DeleteService(ID uint) error
}

type serviceLogic struct {
	repository repository.Repository[model.Service]
}

func NewServiceLogic(repository repository.Repository[model.Service]) ServiceLogic {
	return &serviceLogic{repository: repository}
}

func (l *serviceLogic) GetServiceByID(ID uint) (*model.Service, error) {
	service, err := l.repository.GetByID(ID)
	if err != nil {
		log.Printf("service: Error fetching customer with ID %d: %v", ID, err)
		return nil, fmt.Errorf("failed to fetch medical service with ID %d: %w", ID, err)
	}

	return service, nil
}

func (l *serviceLogic) GetAllServices() ([]model.Service, error) {
	services, err := l.repository.GetAll()
	if err != nil {
		log.Printf("service: Error fetching medical services: %v", err)
		return nil, fmt.Errorf("failed to fetch customers: %w", err)
	}

	if len(services) == 0 {
		log.Println("service: No services found")
		return []model.Service{}, fmt.Errorf("failed to persist medical service: %w", err)
	}
	return services, nil
}

func (l *serviceLogic) CreateService(service *model.Service) error {
	if err := l.repository.Create(service); err != nil {
		log.Printf("service: Error saving medical service: %v", err)
		return fmt.Errorf("failed to persist medical service: %w", err)
	}

	return nil
}

func (l *serviceLogic) UpdateService(ID uint, service *model.Service) (*model.Service, error) {
	serviceUpdate, err := l.GetServiceByID(ID)
	if err != nil {
		log.Printf("service: Error fetching customer with ID %d: %v to update", ID, err)
		return nil, fmt.Errorf("failed to fetch medical service with ID %d: %w to update", ID, err)
	}

	serviceUpdate.Name = service.Name
	serviceUpdate.Description = service.Description
	serviceUpdate.Price = service.Price

	medicalServiceUpdated, err := l.repository.Update(serviceUpdate)
	if err != nil {
		log.Printf("service: Error updating medical service with ID %d: %v", ID, err)
		return nil, fmt.Errorf("failed to update medical service with ID %d: %w", ID, err)
	}

	return medicalServiceUpdated, nil
}

func (l *serviceLogic) DeleteService(ID uint) error {
	if err := l.repository.Delete(ID); err != nil {
		log.Printf("Error deleting customer with ID %d: %v", ID, err)
		return fmt.Errorf("service: failed to delete customer with ID %d: %w", ID, err)
	}

	return nil
}
