package storage

import (
	"RESTAPI/internal/entity"
	"errors"
	"sync"
)

type MessageStorage struct {
	data map[int64]*entity.Message
	mu   sync.Mutex
}

func NewMessageStorage() *MessageStorage {
	return &MessageStorage{
		data: make(map[int64]*entity.Message),
	}
}

func (s *MessageStorage) Create(message *entity.Message) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := int64(len(s.data) + 1)
	message.ID = id
	s.data[id] = message
	return id, nil
}

func (s *MessageStorage) GetById(id int64) (*entity.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	message, exists := s.data[id]
	if !exists {
		return nil, errors.New("message not found")
	}
	return message, nil
}

func (s *MessageStorage) Update(message *entity.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[message.ID]; !exists {
		return errors.New("message not found")
	}
	s.data[message.ID] = message
	return nil
}

func (s *MessageStorage) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[id]; !exists {
		return errors.New("message not found")
	}
	delete(s.data, id)
	return nil
}

func (s *MessageStorage) GetAll() ([]*entity.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	messages := make([]*entity.Message, 0, len(s.data))
	for _, message := range s.data {
		messages = append(messages, message)
	}
	return messages, nil
}
