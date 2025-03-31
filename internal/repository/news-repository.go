package repository

import (
	"RESTAPI/internal/entity"

	"gorm.io/gorm"
)

type NewsRepository struct {
	BaseRepository *BaseRepository[entity.News]
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{
		BaseRepository: NewBaseRepository[entity.News](db),
	}
}

// Create создает новость
func (r *NewsRepository) Create(news *entity.News) error {
	return r.BaseRepository.Create(news)
}

// GetById получает новость по ID
func (r *NewsRepository) GetById(id int64) (entity.News, error) {
	return r.BaseRepository.GetById(id)
}

// Update обновляет новость
func (r *NewsRepository) Update(news *entity.News) error {
	return r.BaseRepository.Update(news)
}

// Delete удаляет новость по ID
func (r *NewsRepository) Delete(id int64) error {
	return r.BaseRepository.Delete(id)
}

// GetAll возвращает все новости
func (r *NewsRepository) GetAll() ([]entity.News, error) {
	var news []entity.News
	result := r.BaseRepository.db.Find(&news)
	if result.Error != nil {
		return nil, result.Error
	}
	return news, nil
}
