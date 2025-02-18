package repository

import (
	"fmt"
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

type PackageRepository interface {
	GetByID(ID uint) (*model.Package, error)
	GetAll() ([]model.Package, error)
}

type packageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) GetByID(ID uint) (*model.Package, error) {
	log.Printf("Getting package with ID: %d", ID)
	pkg := &model.Package{}

	err := r.db.Preload("Services").First(pkg, "id = ?", ID).Error
	if err != nil {
		return nil, err
	}

	log.Printf("resultado: %v", pkg)
	return pkg, nil
}

func (r *packageRepository) GetAll() ([]model.Package, error) {
	var packages []model.Package

	err := r.db.Preload("Services").Find(&packages).Error
	if err != nil {
		return nil, err
	}

	fmt.Printf("resultado: %v", packages)
	return packages, nil
}
