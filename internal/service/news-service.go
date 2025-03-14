package service

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/entity"
	"RESTAPI/internal/storage"
	"errors"
	"time"
)

type NewsService struct {
	repo *storage.NewsStorage
}

func NewNewsService(repo *storage.NewsStorage) *NewsService {
	return &NewsService{repo: repo}
}

func (s *NewsService) Create(req dto.NewsRequestTo) (*dto.NewsResponseTo, error) {
	news := &entity.News{
		WriterID: req.WriterID,
		Title:    req.Title,
		Content:  req.Content,
		Created:  time.Now(),
		Modified: time.Now(),
	}
	id, err := s.repo.Create(news)
	if err != nil {
		return nil, err
	}
	return &dto.NewsResponseTo{
		ID:       id,
		WriterID: news.WriterID,
		Title:    news.Title,
		Content:  news.Content,
		Created:  news.Created,
		Modified: news.Modified,
	}, nil
}

func (s *NewsService) GetById(id int64) (*dto.NewsResponseTo, error) {
	news, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("news not found")
	}
	return &dto.NewsResponseTo{
		ID:       news.ID,
		WriterID: news.WriterID,
		Title:    news.Title,
		Content:  news.Content,
		Created:  news.Created,
		Modified: news.Modified,
	}, nil
}

func (s *NewsService) Update(req dto.NewsUpdateRequestTo) (*dto.NewsResponseTo, error) {
	news := &entity.News{
		WriterID: req.WriterID,
		Title:    req.Title,
		Content:  req.Content,
		ID:       req.ID,
	}
	err := s.repo.Update(news)
	if err != nil {
		return nil, errors.New("failed to update news")
	}
	return &dto.NewsResponseTo{
		ID:       news.ID,
		WriterID: news.WriterID,
		Title:    news.Title,
		Content:  news.Content,
		Created:  news.Created,
		Modified: news.Modified,
	}, nil
}

func (s *NewsService) Delete(id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *NewsService) GetAll() ([]*dto.NewsResponseTo, error) {
	newsList, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	response := make([]*dto.NewsResponseTo, len(newsList))
	for i, news := range newsList {
		response[i] = &dto.NewsResponseTo{
			ID:       news.ID,
			WriterID: news.WriterID,
			Title:    news.Title,
			Content:  news.Content,
			Created:  news.Created,
			Modified: news.Modified,
		}
	}
	return response, nil
}
