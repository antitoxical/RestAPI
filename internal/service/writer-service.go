package service

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/entity"
	"RESTAPI/internal/storage"
	"errors"
)

type WriterService struct {
	repo *storage.WriterStorage
}

func NewWriterService(repo *storage.WriterStorage) *WriterService {
	return &WriterService{repo: repo}
}

func (s *WriterService) Create(req dto.WriterRequestTo) (*dto.WriterResponseTo, error) {
	writer := &entity.Writer{
		Login:     req.Login,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	id, err := s.repo.Create(writer)
	if err != nil {
		return nil, err
	}
	return &dto.WriterResponseTo{
		ID:        id,
		Login:     writer.Login,
		FirstName: writer.FirstName,
		LastName:  writer.LastName,
	}, nil
}

func (s *WriterService) GetById(id int64) (*dto.WriterResponseTo, error) {
	writer, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("writer not found")
	}
	return &dto.WriterResponseTo{
		ID:        writer.ID,
		Login:     writer.Login,
		FirstName: writer.FirstName,
		LastName:  writer.LastName,
	}, nil
}

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

func (s *WriterService) Delete(id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

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
