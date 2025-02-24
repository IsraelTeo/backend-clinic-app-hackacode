package logic

import (
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type ServiceLogic interface {
	GetServiceByID(ID uint) (*model.Service, error)
	GetAllServices(limit, offset int) ([]model.Service, error)
	CreateService(service *model.Service) error
	UpdateService(ID uint, service *model.Service) error
	DeleteService(ID uint) error
}

type serviceLogic struct {
	repository        repository.Repository[model.Service]
	repositoryService repository.ServiceRepository
}

func NewServiceLogic(repository repository.Repository[model.Service], repositoryService repository.ServiceRepository) ServiceLogic {
	return &serviceLogic{repository: repository, repositoryService: repositoryService}
}

func (l *serviceLogic) GetServiceByID(ID uint) (*model.Service, error) {
	service, err := l.repository.GetByID(ID)
	if err != nil {
		log.Printf("service-logic: Error fetching service with ID %d: %v", ID, err)
		return nil, response.ErrorServiceNotFound
	}

	return service, nil
}

func (l *serviceLogic) GetAllServices(limit, offset int) ([]model.Service, error) {
	services, err := l.repository.GetAll(limit, offset)
	if err != nil {
		log.Printf("service-logic: Error fetching services: %v", err)
		return nil, response.ErrorServiceNotFound
	}

	return services, nil
}

func (l *serviceLogic) CreateService(service *model.Service) error {
	err := l.repository.Create(service)
	if err != nil {
		log.Printf("service-logic: Error saving medical service: %v", err)
		return response.ErrorToCreatedService
	}

	return nil
}

func (l *serviceLogic) UpdateService(ID uint, service *model.Service) error {
	serviceUpdate, err := l.GetServiceByID(ID)
	if err != nil {
		log.Printf("service-logic: Error fetching customer with ID %d: %v to update", ID, err)
		return response.ErrorServiceNotFound
	}

	serviceUpdate.Name = service.Name
	serviceUpdate.Description = service.Description
	serviceUpdate.Price = service.Price

	err = l.repository.Update(serviceUpdate)
	if err != nil {
		log.Printf("service-logic: Error updating medical service with ID %d: %v", ID, err)
		return response.ErrorToUpdatedService
	}

	return nil
}

func (l *serviceLogic) DeleteService(ID uint) error {
	_, err := l.repository.GetByID(ID)
	if err != nil {
		log.Printf("service-logic: Error fetching service with ID %d: %v", ID, err)
		return response.ErrorServiceNotFound
	}

	err = l.repositoryService.Delete(ID)
	if err != nil {
		log.Printf("services-logic: Error deleting customer with ID %d: %v", ID, err)
		return response.ErrorToDeletedService
	}

	return nil
}
