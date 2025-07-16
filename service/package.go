package logic

import (
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/calculation"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/repository"
	"github.com/IsraelTeo/clinic-backend-hackacode-app/response"
)

type PackageLogic interface {
	GetPackageByID(ID uint) (*model.Package, error)
	GetAllPackages(limit, offset int) ([]model.Package, error)
	CreatePackage(packageServices *model.CreatePackageRequest) error
	UpdatePackage(ID uint, packageServices *model.CreatePackageRequest) error
	DeletePackage(ID uint) error
}

type packageLogic struct {
	repositoryPkg      repository.Repository[model.Package]
	repositoryPkgMain  repository.PackageRepository
	repositoryServ     repository.Repository[model.Service]
	repositoryServMain repository.ServiceRepository
}

func NewPackageLogic(
	repositoryPkg repository.Repository[model.Package],
	repositoryPkgMain repository.PackageRepository,
	repositoryServ repository.Repository[model.Service],
	repositoryServMain repository.ServiceRepository,
) PackageLogic {
	return &packageLogic{
		repositoryPkg:      repositoryPkg,
		repositoryPkgMain:  repositoryPkgMain,
		repositoryServ:     repositoryServ,
		repositoryServMain: repositoryServMain,
	}
}

func (l *packageLogic) GetPackageByID(ID uint) (*model.Package, error) {
	packageService, err := l.repositoryPkgMain.GetByID(ID)
	if err != nil {
		log.Printf("package-logic: Error fetching package with ID %d: %v", ID, err)
		return nil, response.ErrorPackageNotFound
	}

	return packageService, nil
}

func (l *packageLogic) GetAllPackages(limit, offset int) ([]model.Package, error) {
	packageServices, err := l.repositoryPkgMain.GetAll(limit, offset)
	if err != nil {
		log.Printf("package-logic: Error fetching packages: %v", err)
		return nil, response.ErrorPackageNotFound
	}

	return packageServices, nil
}

func (l *packageLogic) CreatePackage(pkg *model.CreatePackageRequest) error {
	services, err := l.repositoryServMain.GetAll()
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

	err = l.repositoryPkg.Create(&pkgCreated)
	if err != nil {
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

	err = l.repositoryPkgMain.ClearServices(ID)
	if err != nil {
		log.Printf("package: Error clearing services for package ID %d: %v", ID, err)
		return response.ErrorClearingServices
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

	err = l.repositoryPkg.Update(existingPackage)
	if err != nil {
		log.Printf("package: Error updating package with ID %d: %v", ID, err)
		return response.ErrorToUpdatedPackage
	}

	return nil
}

func (l *packageLogic) DeletePackage(ID uint) error {
	_, err := l.repositoryPkg.GetByID(ID)
	if err != nil {
		log.Printf("package-logic: Error fetching package with ID %d: %v", ID, err)
		return response.ErrorPackageNotFound
	}

	err = l.repositoryPkgMain.Delete(ID)
	if err != nil {
		log.Printf("package: Error deleting package with ID %d: %v", ID, err)
		return response.ErrorToDeletedPackage
	}

	return nil
}
