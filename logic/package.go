package logic

import (
	"log"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/calculation"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type PackageLogic interface {
	GetPackageByID(ID uint) (*model.Package, error)
	GetAllPackages() ([]model.Package, error)
	CreatePackage(packageServices *model.CreatePackageRequest) error
	UpdatePackage(ID uint, packageServices *model.CreatePackageRequest) error
	DeletePackage(ID uint) error
}

type packageLogic struct {
	repositoryPkg  repository.Repository[model.Package]
	repositoryServ repository.Repository[model.Service]
}

func NewPackageLogic(
	repositoryPkg repository.Repository[model.Package],
	repositoryServ repository.Repository[model.Service],
) PackageLogic {
	return &packageLogic{
		repositoryPkg:  repositoryPkg,
		repositoryServ: repositoryServ,
	}
}

func (l *packageLogic) GetPackageByID(ID uint) (*model.Package, error) {
	packageService, err := l.repositoryPkg.GetByID(ID)
	if err != nil {
		log.Printf("package-logic: Error fetching package with ID %d: %v", ID, err)
		return nil, response.ErrorPackageNotFound
	}

	return packageService, nil
}

func (l *packageLogic) GetAllPackages() ([]model.Package, error) {
	packageServices, err := l.repositoryPkg.GetAll()
	if err != nil {
		log.Printf("package-logib: Error fetching packages: %v", err)
		return nil, response.ErrorPackageNotFound
	}

	if len(packageServices) == 0 {
		log.Println("package: No packages found")
		return []model.Package{}, response.ErrorListPackagesEmpty
	}

	return packageServices, nil
}

func (l *packageLogic) CreatePackage(pkg *model.CreatePackageRequest) error {
	services, err := l.repositoryServ.GetAll()
	if err != nil {
		log.Printf("package-logic: Error fetching services: %v", err)
		return response.ErrorFetchingServices
	}

	if len(services) == 0 {
		log.Println("package: No services found")
		return response.ErrorListServicesEmpty
	}

	selectedServices := []model.Service{}

	for _, serviceID := range pkg.ServiceIDs {
		serviceFound, err := l.repositoryServ.GetByID(serviceID)
		if err != nil {
			log.Printf("package: Error fetching service with ID %d: %v", serviceID, err)
			return response.ErrorServiceNotFound
		}

		selectedServices = append(selectedServices, *serviceFound)
	}

	finalPkgPrice := calculation.TotalServicePackageAmount(selectedServices)

	pkgCreated := model.Package{
		Name:     pkg.Name,
		Services: selectedServices,
		Price:    finalPkgPrice.FinalPrice,
	}

	if err := l.repositoryPkg.Create(&pkgCreated); err != nil {
		log.Printf("package: Error saving package: %v", err)
		return response.ErrorToCreatedPackage
	}

	return nil
}

func (l *packageLogic) UpdatePackage(ID uint, packageServices *model.CreatePackageRequest) error {
	existingPackage, err := l.GetPackageByID(ID)
	if err != nil {
		log.Printf("package: Error fetching package with ID %d: %v to update", ID, err)
		return response.ErrorPackageNotFound
	}

	services, err := l.repositoryServ.GetAll()
	if err != nil {
		log.Printf("package: Error fetching services: %v", err)
		return response.ErrorFetchingServices
	}

	if len(services) == 0 {
		log.Println("package: No services found")
		return response.ErrorListServicesEmpty
	}

	selectedServices := []model.Service{}

	for _, serviceID := range packageServices.ServiceIDs {
		serviceFound, err := l.repositoryServ.GetByID(serviceID)
		if err != nil {
			log.Printf("package: Error fetching service with ID %d: %v", serviceID, err)
			return response.ErrorServiceNotFound
		}
		selectedServices = append(selectedServices, *serviceFound)
	}

	finalPkgPrice := calculation.TotalServicePackageAmount(selectedServices)

	existingPackage.Name = packageServices.Name
	existingPackage.Services = selectedServices
	existingPackage.Price = finalPkgPrice.FinalPrice

	if err = l.repositoryPkg.Update(existingPackage); err != nil {
		log.Printf("package: Error updating package with ID %d: %v", ID, err)
		return response.ErrorToUpdatedPackage
	}

	return nil
}

func (l *packageLogic) DeletePackage(ID uint) error {
	if _, err := l.repositoryPkg.GetByID(ID); err != nil {
		log.Printf("package-logic: Error fetching package with ID %d: %v", ID, err)
		return response.ErrorPackageNotFound
	}

	if err := l.repositoryPkg.Delete(ID); err != nil {
		log.Printf("package: Error deleting package with ID %d: %v", ID, err)
		return response.ErrorToDeletedPackage
	}

	return nil
}
