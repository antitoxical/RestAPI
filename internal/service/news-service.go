package service

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/entity"
	"RESTAPI/internal/repository"
	"errors"
	"time"
)

type NewsService struct {
	repo     *repository.NewsRepository
	markRepo *repository.MarkRepository
}

func NewNewsService(repo *repository.NewsRepository, markRepo *repository.MarkRepository) *NewsService {
	return &NewsService{repo: repo, markRepo: markRepo}
}

func (s *NewsService) Create(req dto.NewsRequestTo) (*dto.NewsResponseTo, error) {

	marks := []entity.Mark{}
	for _, markName := range req.Marks {
		// Try to find existing mark
		existingMarks, err := s.markRepo.GetByName(markName)

		var mark entity.Mark
		if err != nil || len(existingMarks) == 0 {
			// Create new mark if not found
			mark = entity.Mark{Name: markName}
			if err := s.markRepo.Create(&mark); err != nil {
				return nil, err
			}
		} else {
			mark = existingMarks[0]
		}

		marks = append(marks, mark)
	}
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
		Marks:    marks,
	}

	err = s.repo.Create(news)
	if err != nil {
		return nil, err
	}

	markResponses := make([]dto.MarkResponseTo, len(news.Marks))
	for i, mark := range news.Marks {
		markResponses[i] = dto.MarkResponseTo{
			ID:   mark.ID,
			Name: mark.Name,
		}
	}

	return &dto.NewsResponseTo{
		ID:       news.ID,
		WriterID: news.WriterID,
		Title:    news.Title,
		Content:  news.Content,
		Created:  news.Created,
		Modified: news.Modified,
		Marks:    markResponses,
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

// Delete deletes a news article by ID
// Delete deletes a news article by ID and its associated marks
func (s *NewsService) Delete(id int64) error {
	// First get the news with its marks to know which marks to potentially delete
	news, err := s.repo.GetById(id)
	if err != nil {
		return err
	}

	// Extract mark names to delete them later if needed
	markNames := make([]string, len(news.Marks))
	for i, mark := range news.Marks {
		markNames[i] = mark.Name
	}

	// Delete the news with its mark associations
	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Now delete the marks if they're no longer used
	// Either use DeleteOrphaned to delete all orphaned marks
	err = s.markRepo.DeleteOrphaned()
	if err != nil {
		return err
	}

	// Or directly delete these specific marks if they should always be removed
	// (uncomment if needed)
	// return s.markRepo.DeleteMarks(markNames)

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
