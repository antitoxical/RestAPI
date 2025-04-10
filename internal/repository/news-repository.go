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
// Delete deletes a news article by ID and removes mark associations
func (r *NewsRepository) Delete(id int64) error {
	// First, get the news with its marks
	news, err := r.GetById(id)
	if err != nil {
		return err
	}

	// Begin a transaction
	tx := r.BaseRepository.db.Begin()

	if err := tx.Exec("DELETE FROM news_mark WHERE news_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Then delete the news
	if err := tx.Exec("DELETE FROM tbl_news WHERE id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Remove the associations in the join table
	if err := tx.Model(&news).Association("Marks").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Delete the news article
	if err := tx.Delete(&news).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
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
