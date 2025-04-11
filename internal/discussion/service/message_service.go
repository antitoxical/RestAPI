package service

import (
	"context"
	"RESTAPI/internal/discussion/model"
	"RESTAPI/internal/discussion/repository"
)

// MessageService handles business logic for messages
type MessageService struct {
	repo repository.MessageRepository
}

// NewMessageService creates a new MessageService
func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

// CreateMessage creates a new message
func (s *MessageService) CreateMessage(ctx context.Context, message *model.Message) error {
	return s.repo.Create(ctx, message)
}

// GetMessage retrieves a message by ID
func (s *MessageService) GetMessage(ctx context.Context, id int64) (*model.Message, error) {
	return s.repo.FindByID(ctx, id)
}

// GetMessagesByNewsID retrieves all messages for a news item
func (s *MessageService) GetMessagesByNewsID(ctx context.Context, newsID int64) ([]*model.Message, error) {
	return s.repo.FindByNewsID(ctx, newsID)
}

// UpdateMessage updates an existing message
func (s *MessageService) UpdateMessage(ctx context.Context, message *model.Message) error {
	return s.repo.Update(ctx, message)
}

// DeleteMessage deletes a message by ID
func (s *MessageService) DeleteMessage(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// GetAllMessages retrieves all messages
func (s *MessageService) GetAllMessages(ctx context.Context) ([]*model.Message, error) {
	return s.repo.FindAll(ctx)
} 