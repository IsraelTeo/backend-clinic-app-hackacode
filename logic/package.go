package logic

import (
	"log"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type PackageLogic interface {
	GetPackageByID(ID uint) (*model.Package, error)
	GetAllPackages() ([]model.Package, error)
	CreatePackage(packageServices *model.CreatePackageRequest) error
	UpdatePackage(ID uint, packageServices *model.Package) error
	DeletePackage(ID uint) error
}

type packageLogic struct {
	repositoryPkg  repository.Repository[model.Package]
	repositoryServ repository.Repository[model.Service]
}

func NewPackageLogic(repositoryPkg repository.Repository[model.Package], repositoryServ repository.Repository[model.Service]) PackageLogic {
	return &packageLogic{repositoryPkg: repositoryPkg, repositoryServ: repositoryServ}
}

func (l *packageLogic) GetPackageByID(ID uint) (*model.Package, error) {
	packageService, err := l.repositoryPkg.GetByID(ID)
	if err != nil {
		log.Printf("package: Error fetching package with ID %d: %v", ID, err)
		return nil, response.ErrorPackageNotFound
	}
	return packageService, nil
}

func (l *packageLogic) GetAllPackages() ([]model.Package, error) {
	packageServices, err := l.repositoryPkg.GetAll()
	if err != nil {
		log.Printf("package: Error fetching packages: %v", err)
		return nil, response.ErrorPackageNotFound
	}

	if len(packageServices) == 0 {
		log.Println("package: No packages found")
		return []model.Package{}, response.ErrorListPackagesEmpty
	}
	return packageServices, nil
}

func (l *packageLogic) CreatePackage(pkg *model.CreatePackageRequest) error {
	//Obtiene todos los servicios
	services, err := l.repositoryServ.GetAll()
	if err != nil {
		log.Printf("package: Error fetching services: %v", err)
		return response.ErrorFetchingServices
	}

	//Verifica que existan servicios
	if len(services) == 0 {
		log.Println("package: No services found")
		return response.ErrorListServicesEmpty
	}

	//Obtiene los servicios por ID y los asigna una nueva lista
	for _, service := range services {
		serviceFound, err := l.repositoryServ.GetByID(service.ID)
		if err != nil {
			log.Printf("package: Error fetching service with ID %d: %v", service.ID, err)
			return response.ErrorServiceNotFound
		}
		services = append(services, *serviceFound)
	}

	//Calcula el total de los servicios
	var totalAmount float64
	for _, service := range services {
		totalAmount += service.Price
	}

	//Crea el paquete los servicios encontrados y el monto total
	pgkCreated := model.Package{
		Name:     pkg.Name,
		Services: services,
		Price:    totalAmount,
	}

	//Se guarda el paquete
	if err := l.repositoryPkg.Create(&pgkCreated); err != nil {
		log.Printf("package: Error saving package: %v", err)
		return response.ErrorToCreatedPackage
	}

	return nil
}

func (l *packageLogic) UpdatePackage(ID uint, packageServices *model.Package) error {
	//Encuentra el paquete para actualizar
	existingPackage, err := l.GetPackageByID(ID)
	if err != nil {
		log.Printf("package: Error fetching package with ID %d: %v to update", ID, err)
		return response.ErrorPackageNotFound
	}

	//Encuentra los servicios
	services, err := l.repositoryServ.GetAll()
	if err != nil {
		log.Printf("package: Error fetching services: %v", err)
		return response.ErrorFetchingServices
	}

	//Verifica que existan servicios
	if len(services) == 0 {
		log.Println("package: No services found")
		return response.ErrorListServicesEmpty
	}

	//Obtiene los servicios por ID y los asigna una nueva lista
	for _, service := range services {
		serviceFound, err := l.repositoryServ.GetByID(service.ID)
		if err != nil {
			log.Printf("package: Error fetching service with ID %d: %v", service.ID, err)
			return response.ErrorServiceNotFound
		}
		services = append(services, *serviceFound)
	}

	var totalAmount float64
	for _, service := range services {
		totalAmount += service.Price
	}

	existingPackage.Name = packageServices.Name
	existingPackage.Price = totalAmount
	existingPackage.Services = services

	if err = l.repositoryPkg.Update(existingPackage); err != nil {
		log.Printf("package: Error updating package with ID %d: %v", ID, err)
		return response.ErrorToUpdatedPackage
	}

	return nil
}

func (l *packageLogic) DeletePackage(ID uint) error {
	if err := l.repositoryPkg.Delete(ID); err != nil {
		log.Printf("package: Error deleting package with ID %d: %v", ID, err)
		return response.ErrorToDeletedPackage
	}

	return nil
}
