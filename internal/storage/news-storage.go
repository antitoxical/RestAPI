package storage

import (
	"RESTAPI/internal/entity"
	"errors"
	"sync"
	"time"
)

type NewsStorage struct {
	data map[int64]*entity.News
	mu   sync.Mutex
}

func NewNewsStorage() *NewsStorage {
	return &NewsStorage{
		data: make(map[int64]*entity.News),
	}
}

func (s *NewsStorage) Create(news *entity.News) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := int64(len(s.data) + 1)
	news.ID = id
	news.Created = time.Now()
	news.Modified = time.Now()
	s.data[id] = news
	return id, nil
}

func (s *NewsStorage) GetById(id int64) (*entity.News, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	news, exists := s.data[id]
	if !exists {
		return nil, errors.New("news not found")
	}
	return news, nil
}

func (s *NewsStorage) Update(news *entity.News) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[news.ID]; !exists {
		return errors.New("news not found")
	}
	news.Modified = time.Now()
	s.data[news.ID] = news
	return nil
}

func (s *NewsStorage) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[id]; !exists {
		return errors.New("news not found")
	}
	delete(s.data, id)
	return nil
}

func (s *NewsStorage) GetAll() ([]*entity.News, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newsList := make([]*entity.News, 0, len(s.data))
	for _, news := range s.data {
		newsList = append(newsList, news)
	}
	return newsList, nil
}
