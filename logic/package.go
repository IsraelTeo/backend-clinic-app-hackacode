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
	CreatePackage(packageServices *model.Package) error
	UpdatePackage(ID uint, packageServices *model.Package) error
	DeletePackage(ID uint) error
}

type packageLogic struct {
	repository repository.Repository[model.Package]
}

func NewPackageLogic(repository repository.Repository[model.Package]) PackageLogic {
	return &packageLogic{repository: repository}
}

func (l *packageLogic) GetPackageByID(ID uint) (*model.Package, error) {
	packageService, err := l.repository.GetByID(ID)
	if err != nil {
		log.Printf("package: Error fetching package with ID %d: %v", ID, err)
		return nil, response.ErrorPackageNotFound
	}
	return packageService, nil
}

func (l *packageLogic) GetAllPackages() ([]model.Package, error) {
	packageServices, err := l.repository.GetAll()
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

func (l *packageLogic) CreatePackage(packageServices *model.Package) error {
	if err := l.repository.Create(packageServices); err != nil {
		log.Printf("package: Error saving package: %v", err)
		return response.ErrorToCreatedPackage
	}

	return nil
}

func (l *packageLogic) UpdatePackage(ID uint, packageServices *model.Package) error {
	existingPackage, err := l.GetPackageByID(ID)
	if err != nil {
		log.Printf("package: Error fetching package with ID %d: %v to update", ID, err)
		return response.ErrorPackageNotFound
	}

	existingPackage.Name = packageServices.Name
	existingPackage.Price = packageServices.Price
	existingPackage.Services = packageServices.Services

	if err = l.repository.Update(existingPackage); err != nil {
		log.Printf("package: Error updating package with ID %d: %v", ID, err)
		return response.ErrorToUpdatedPackage
	}

	return nil
}

func (l *packageLogic) DeletePackage(ID uint) error {
	if err := l.repository.Delete(ID); err != nil {
		log.Printf("package: Error deleting package with ID %d: %v", ID, err)
		return response.ErrorToDeletedPackage
	}

	return nil
}
