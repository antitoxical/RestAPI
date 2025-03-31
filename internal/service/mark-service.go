package service

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/entity"
	"RESTAPI/internal/repository"
	"errors"
)

type MarkService struct {
	repo *repository.MarkRepository
}

func NewMarkService(repo *repository.MarkRepository) *MarkService {
	return &MarkService{repo: repo}
}

func (s *MarkService) Create(req dto.MarkRequestTo) (*dto.MarkResponseTo, error) {
	mark := &entity.Mark{
		Name: req.Name,
	}
	err := s.repo.Create(mark)
	if err != nil {
		return nil, err
	}
	return &dto.MarkResponseTo{
		ID:   mark.ID,
		Name: mark.Name,
	}, nil
}

func (s *MarkService) GetById(id int64) (*dto.MarkResponseTo, error) {
	mark, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("mark not found")
	}
	return &dto.MarkResponseTo{
		ID:   mark.ID,
		Name: mark.Name,
	}, nil
}

func (s *MarkService) Update(req dto.MarkUpdateRequestTo) (*dto.MarkResponseTo, error) {
	mark := &entity.Mark{
		Name: req.Name,
		ID:   req.ID,
	}
	err := s.repo.Update(mark)
	if err != nil {
		return nil, errors.New("failed to update mark")
	}
	return &dto.MarkResponseTo{
		ID:   mark.ID,
		Name: mark.Name,
	}, nil
}

func (s *MarkService) Delete(id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *MarkService) GetAll() ([]*dto.MarkResponseTo, error) {
	marks, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	response := make([]*dto.MarkResponseTo, len(marks))
	for i, mark := range marks {
		response[i] = &dto.MarkResponseTo{
			ID:   mark.ID,
			Name: mark.Name,
		}
	}
	return response, nil
}
