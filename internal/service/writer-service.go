// Projects/RESTAPI/internal/service/writer-service.go
package service

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/entity"
	"RESTAPI/internal/repository"
	"fmt"
)

type WriterService struct {
	repo *repository.WriterRepository
}

func NewWriterService(repo *repository.WriterRepository) *WriterService {
	return &WriterService{repo: repo}
}

// Create creates a new writer
func (s *WriterService) Create(req dto.WriterRequestTo) (*dto.WriterResponseTo, error) {
	// Check if the login already exists
	existingWriter, err := s.repo.GetByLogin(req.Login)
	if err == nil && existingWriter != nil {
		return nil, fmt.Errorf("login_already_exists")
	}

	writer := &entity.Writer{
		Login:     req.Login,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	err = s.repo.Create(writer)
	if err != nil {
		return nil, err
	}

	return &dto.WriterResponseTo{
		ID:        writer.ID,
		Login:     writer.Login,
		FirstName: writer.FirstName,
		LastName:  writer.LastName,
	}, nil
}

// GetById gets a writer by ID
func (s *WriterService) GetById(id int64) (*dto.WriterResponseTo, error) {
	if id < 0 {
		return nil, fmt.Errorf("invalid ID: %d", id)
	}

	writer, err := s.repo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("writer not found")
	}

	return &dto.WriterResponseTo{
		ID:        writer.ID,
		Login:     writer.Login,
		FirstName: writer.FirstName,
		LastName:  writer.LastName,
	}, nil
}

// Update updates a writer
func (s *WriterService) Update(req dto.WriterUpdateRequestTo) (*dto.WriterResponseTo, error) {
	writer := &entity.Writer{
		Login:     req.Login,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		ID:        req.ID,
	}

	err := s.repo.Update(writer)
	if err != nil {
		return nil, err
	}

	return &dto.WriterResponseTo{
		ID:        writer.ID,
		Login:     writer.Login,
		FirstName: writer.FirstName,
		LastName:  writer.LastName,
	}, nil
}

// Delete deletes a writer
func (s *WriterService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// GetAll returns all writers
func (s *WriterService) GetAll() ([]*dto.WriterResponseTo, error) {
	writers, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	response := make([]*dto.WriterResponseTo, len(writers))
	for i, writer := range writers {
		response[i] = &dto.WriterResponseTo{
			ID:        writer.ID,
			Login:     writer.Login,
			FirstName: writer.FirstName,
			LastName:  writer.LastName,
		}
	}
	return response, nil
}
