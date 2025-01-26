package repository

import (
	"gorm.io/gorm"
)

type Repository[T any] interface {
	GetByID(ID uint) (*T, error)
	GetAll() ([]T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(ID uint) error
}

type repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) Repository[T] {
	return &repository[T]{db: db}
}

func (r *repository[T]) GetByID(ID uint) (*T, error) {
	entity := new(T)
	if err := r.db.First(entity, ID).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *repository[T]) GetAll() ([]T, error) {
	var entities []T
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository[T]) Create(entity *T) error {
	if err := r.db.Create(entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository[T]) Update(entity *T) error {
	if err := r.db.Save(entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository[T]) Delete(ID uint) error {
	if err := r.db.Delete(new(T), ID).Error; err != nil {
		return err
	}
	return nil
}
