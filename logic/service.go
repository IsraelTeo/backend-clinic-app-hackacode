package logic

import (
	"log"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type ServiceLogic interface {
	GetServiceByID(ID uint) (*model.Service, error)
	GetAllServices() ([]model.Service, error)
	CreateService(service *model.Service) error
	UpdateService(ID uint, service *model.Service) error
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
		log.Printf("service: Error fetching service with ID %d: %v", ID, err)
		return nil, response.ErrorServiceNotFound
	}
	return service, nil
}

func (l *serviceLogic) GetAllServices() ([]model.Service, error) {
	services, err := l.repository.GetAll()
	if err != nil {
		log.Printf("service: Error fetching services: %v", err)
		return nil, response.ErrorServiceNotFound
	}

	if len(services) == 0 {
		log.Println("service: No services found")
		return []model.Service{}, response.ErrorListServicesEmpty
	}
	return services, nil
}

func (l *serviceLogic) CreateService(service *model.Service) error {
	if err := l.repository.Create(service); err != nil {
		log.Printf("service: Error saving medical service: %v", err)
		return response.ErrorToCreated
	}

	return nil
}

func (l *serviceLogic) UpdateService(ID uint, service *model.Service) error {
	serviceUpdate, err := l.GetServiceByID(ID)
	if err != nil {
		log.Printf("service: Error fetching customer with ID %d: %v to update", ID, err)
		return response.ErrorServiceNotFound
	}

	serviceUpdate.Name = service.Name
	serviceUpdate.Description = service.Description
	serviceUpdate.Price = service.Price

	if err = l.repository.Update(serviceUpdate); err != nil {
		log.Printf("service: Error updating medical service with ID %d: %v", ID, err)
		return response.ErrorToUpdated
	}

	return nil
}

func (l *serviceLogic) DeleteService(ID uint) error {
	if err := l.repository.Delete(ID); err != nil {
		log.Printf("Error deleting customer with ID %d: %v", ID, err)
		return response.ErrorToDeleted
	}

	return nil
}
