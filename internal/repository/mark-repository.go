package repository

import (
	"RESTAPI/internal/entity"

	"gorm.io/gorm"
)

type MarkRepository struct {
	BaseRepository *BaseRepository[entity.Mark]
}

func NewMarkRepository(db *gorm.DB) *MarkRepository {
	return &MarkRepository{
		BaseRepository: NewBaseRepository[entity.Mark](db),
	}
}

// Create создает метку
func (r *MarkRepository) Create(mark *entity.Mark) error {
	return r.BaseRepository.Create(mark)
}

// GetById получает метку по ID
func (r *MarkRepository) GetById(id int64) (entity.Mark, error) {
	return r.BaseRepository.GetById(id)
}

// Update обновляет метку
func (r *MarkRepository) Update(mark *entity.Mark) error {
	return r.BaseRepository.Update(mark)
}

// Delete удаляет метку по ID
func (r *MarkRepository) Delete(id int64) error {
	return r.BaseRepository.Delete(id)
}

// GetAll возвращает все метки
func (r *MarkRepository) GetAll() ([]entity.Mark, error) {
	var marks []entity.Mark
	result := r.BaseRepository.db.Find(&marks)
	if result.Error != nil {
		return nil, result.Error
	}
	return marks, nil
}
