package storage

import (
	"RESTAPI/internal/entity"
	"errors"
	"sync"
)

type MarkStorage struct {
	data map[int64]*entity.Mark
	mu   sync.Mutex
}

func NewMarkStorage() *MarkStorage {
	return &MarkStorage{
		data: make(map[int64]*entity.Mark),
	}
}

func (s *MarkStorage) Create(mark *entity.Mark) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := int64(len(s.data) + 1)
	mark.ID = id
	s.data[id] = mark
	return id, nil
}

func (s *MarkStorage) GetById(id int64) (*entity.Mark, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	mark, exists := s.data[id]
	if !exists {
		return nil, errors.New("mark not found")
	}
	return mark, nil
}

func (s *MarkStorage) Update(mark *entity.Mark) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[mark.ID]; !exists {
		return errors.New("mark not found")
	}
	s.data[mark.ID] = mark
	return nil
}

func (s *MarkStorage) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[id]; !exists {
		return errors.New("mark not found")
	}
	delete(s.data, id)
	return nil
}

func (s *MarkStorage) GetAll() ([]*entity.Mark, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	marks := make([]*entity.Mark, 0, len(s.data))
	for _, mark := range s.data {
		marks = append(marks, mark)
	}
	return marks, nil
}
