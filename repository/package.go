package repository

import (
	"log"

	"github.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"gorm.io/gorm"
)

type PackageRepository interface {
	GetByID(ID uint) (*model.Package, error)
	GetAll(limit, offset int) ([]model.Package, error)
	ClearServices(packageID uint) error
	Delete(ID uint) error
}

type packageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) GetByID(ID uint) (*model.Package, error) {
	pkg := &model.Package{}

	err := r.db.
		Preload("Services").
		First(pkg, "id = ?", ID).
		Error
	if err != nil {
		return nil, err
	}

	log.Printf("resultado: %v", pkg)
	return pkg, nil
}

func (r *packageRepository) GetAll(limit, offset int) ([]model.Package, error) {
	var packages []model.Package

	query := r.db.Preload("Services")
	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&packages).Error
	if err != nil {
		return nil, err
	}

	return packages, nil
}

func (r *packageRepository) ClearServices(packageID uint) error {
	return r.db.Exec("DELETE FROM package_services WHERE package_id = ?", packageID).Error
}

func (r *packageRepository) Delete(ID uint) error {
	err := r.ClearServices(ID)
	if err != nil {
		return err
	}

	return r.db.Delete(&model.Package{}, ID).Error
}
