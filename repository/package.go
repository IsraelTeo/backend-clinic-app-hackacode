package repository

import (
	"log"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

type PackageRepository interface {
	GetByID(id uint) (*model.Package, error)
	GetAll() ([]model.Package, error)
}

type packageRepository struct {
	db *gorm.DB
}

// Nueva función que crea una nueva instancia del repositorio
func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{db: db}
}

// Implementación del método GetByID
func (r *packageRepository) GetByID(id uint) (*model.Package, error) {
	pkg := &model.Package{}

	// Pre-cargar los servicios asociados al paquete
	err := r.db.Preload("Services").First(pkg, "id = ?", id).Error
	if err != nil {
		log.Printf("package: Error fetching package with ID %d: %v", id, err)
		return nil, err
	}

	return pkg, nil
}

// Implementación del método GetAll
func (r *packageRepository) GetAll() ([]model.Package, error) {
	var packages []model.Package

	// Obtener todos los paquetes, con los servicios pre-cargados
	err := r.db.Preload("Services").Find(&packages).Error
	if err != nil {
		log.Printf("package: Error fetching all packages: %v", err)
		return nil, err
	}

	return packages, nil
}
