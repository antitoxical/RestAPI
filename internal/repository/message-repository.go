package repository

import (
	"RESTAPI/internal/entity"

	"gorm.io/gorm"
)

type MessageRepository struct {
	BaseRepository *BaseRepository[entity.Message]
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		BaseRepository: NewBaseRepository[entity.Message](db),
	}
}

// Create создает новое сообщение
func (r *MessageRepository) Create(message *entity.Message) error {
	return r.BaseRepository.Create(message)
}

// GetById получает сообщение по ID
func (r *MessageRepository) GetById(id int64) (entity.Message, error) {
	return r.BaseRepository.GetById(id)
}

// Update обновляет существующее сообщение
func (r *MessageRepository) Update(message *entity.Message) error {
	return r.BaseRepository.Update(message)
}

// Delete удаляет сообщение по ID
func (r *MessageRepository) Delete(id int64) error {
	return r.BaseRepository.Delete(id)
}

// GetAll возвращает все сообщения
func (r *MessageRepository) GetAll() ([]entity.Message, error) {
	var messages []entity.Message
	result := r.BaseRepository.db.Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}
