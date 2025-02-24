package repository

import (
	"gorm.io/gorm"
)

type Repository[T any] interface {
	GetByID(ID uint) (*T, error)
	GetAll(limit, offset int) ([]T, error)
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

	err := r.db.First(entity, ID).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *repository[T]) GetAll(limit, offset int) ([]T, error) {
	var entities []T

	query := r.db.
		Limit(limit).
		Offset(offset)

	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *repository[T]) Create(entity *T) error {
	err := r.db.Create(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository[T]) Update(entity *T) error {
	err := r.db.Save(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository[T]) Delete(ID uint) error {
	err := r.db.Delete(new(T), ID).Error
	if err != nil {
		return err
	}

	return nil
}
