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

// Add this method to your MarkRepository

// GetByName returns marks with the specified name
func (r *MarkRepository) GetByName(name string) ([]entity.Mark, error) {
	var marks []entity.Mark
	result := r.BaseRepository.db.Where("name = ?", name).Find(&marks)
	if result.Error != nil {
		return nil, result.Error
	}
	return marks, nil
}

// DeleteOrphaned deletes marks that are not associated with any news
func (r *MarkRepository) DeleteOrphaned() error {
	// This SQL finds and deletes marks that don't have relationships in the join table
	return r.BaseRepository.db.Exec(`
        DELETE FROM tbl_mark 
        WHERE id NOT IN (
            SELECT DISTINCT mark_id FROM news_mark
        )
    `).Error
}

// DeleteByName deletes a mark by its name
func (r *MarkRepository) DeleteByName(name string) error {
	return r.BaseRepository.db.Where("name = ?", name).Delete(&entity.Mark{}).Error
}

// DeleteMarks deletes marks by their names
func (r *MarkRepository) DeleteMarks(names []string) error {
	if len(names) == 0 {
		return nil
	}
	return r.BaseRepository.db.Where("name IN ?", names).Delete(&entity.Mark{}).Error
}
