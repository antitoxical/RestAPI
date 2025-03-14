package storage

import (
	"RESTAPI/internal/entity"
	"errors"
	"sync"
)

type WriterStorage struct {
	data map[int64]*entity.Writer
	mu   sync.Mutex
}

func NewWriterStorage() *WriterStorage {
	return &WriterStorage{
		data: make(map[int64]*entity.Writer),
	}
}

func (s *WriterStorage) Create(writer *entity.Writer) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := int64(len(s.data) + 1)
	writer.ID = id
	s.data[id] = writer
	return id, nil
}

func (s *WriterStorage) GetById(id int64) (*entity.Writer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	writer, exists := s.data[id]
	if !exists {
		return nil, errors.New("writer not found")
	}
	return writer, nil
}

func (s *WriterStorage) Update(writer *entity.Writer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[writer.ID]; !exists {
		return errors.New("writer not found")
	}
	s.data[writer.ID] = writer
	return nil
}

func (s *WriterStorage) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[id]; !exists {
		return errors.New("writer not found")
	}
	delete(s.data, id)
	return nil
}

func (s *WriterStorage) GetAll() ([]*entity.Writer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	writers := make([]*entity.Writer, 0, len(s.data))
	for _, writer := range s.data {
		writers = append(writers, writer)
	}
	return writers, nil
}
