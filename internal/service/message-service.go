package service

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/entity"
	"RESTAPI/internal/storage"
	"errors"
)

type MessageService struct {
	repo *storage.MessageStorage
}

func NewMessageService(repo *storage.MessageStorage) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) Create(req dto.MessageRequestTo) (*dto.MessageResponseTo, error) {
	message := &entity.Message{
		NewsID:  req.NewsID,
		Content: req.Content,
	}
	id, err := s.repo.Create(message)
	if err != nil {
		return nil, err
	}
	return &dto.MessageResponseTo{
		ID:      id,
		NewsID:  message.NewsID,
		Content: message.Content,
	}, nil
}

func (s *MessageService) GetById(id int64) (*dto.MessageResponseTo, error) {
	message, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("message not found")
	}
	return &dto.MessageResponseTo{
		ID:      message.ID,
		NewsID:  message.NewsID,
		Content: message.Content,
	}, nil
}

func (s *MessageService) Update(req dto.MessageUpdateRequestTo) (*dto.MessageResponseTo, error) {
	message := &entity.Message{
		ID:      req.ID,
		NewsID:  req.NewsID,
		Content: req.Content,
	}
	err := s.repo.Update(message)
	if err != nil {
		return nil, errors.New("failed to update message")
	}
	return &dto.MessageResponseTo{
		ID:      message.ID,
		NewsID:  message.NewsID,
		Content: message.Content,
	}, nil
}

func (s *MessageService) Delete(id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *MessageService) GetAll() ([]*dto.MessageResponseTo, error) {
	messages, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	response := make([]*dto.MessageResponseTo, len(messages))
	for i, message := range messages {
		response[i] = &dto.MessageResponseTo{
			ID:      message.ID,
			NewsID:  message.NewsID,
			Content: message.Content,
		}
	}
	return response, nil
}
