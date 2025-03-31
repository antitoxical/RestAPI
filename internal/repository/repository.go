// Projects/RESTAPI/internal/repository/repository.go
package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

// Create creates a new record and populates its ID
func (r *BaseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// GetById gets a record by ID
func (r *BaseRepository[T]) GetById(id int64) (T, error) {
	var result T
	if err := r.db.First(&result, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, fmt.Errorf("record not found")
		}
		return result, err
	}
	return result, nil
}

// Update updates an existing record
func (r *BaseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete deletes a record by ID
func (r *BaseRepository[T]) Delete(id int64) error {
	return r.db.Delete(new(T), id).Error
}

// List returns a list of records with filtering, sorting and pagination
func (r *BaseRepository[T]) List(page, pageSize int, filter map[string]interface{}, sort string) ([]T, int64, error) {
	var entities []T
	var total int64

	query := r.db.Model(new(T)).Where(filter)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting and pagination
	if err := query.Order(sort).Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
