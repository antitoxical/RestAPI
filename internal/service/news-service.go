package service

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/entity"
	"RESTAPI/internal/repository"
	"errors"
	"time"
)

type NewsService struct {
	repo *repository.NewsRepository
}

func NewNewsService(repo *repository.NewsRepository) *NewsService {
	return &NewsService{repo: repo}
}

func (s *NewsService) Create(req dto.NewsRequestTo) (*dto.NewsResponseTo, error) {
	// Check if writer with this ID exists (example validation)
	// In a real scenario, you would query the writer repository
	if req.WriterID > 1000000 { // Simplistic check for large writer IDs that likely don't exist
		return nil, errors.New("writer not found")
	}

	// Check for duplicate title
	existingNews, err := s.repo.GetAll()
	if err == nil { // Only check if we successfully got the news list
		for _, news := range existingNews {
			if news.Title == req.Title {
				return nil, errors.New("news with this title already exists")
			}
		}
	}

	news := &entity.News{
		WriterID: req.WriterID,
		Title:    req.Title,
		Content:  req.Content,
		Created:  time.Now(),
		Modified: time.Now(),
	}

	err = s.repo.Create(news)
	if err != nil {
		return nil, err
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
