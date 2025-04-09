// Projects/RESTAPI/internal/repository/writer-repository.go
package repository

import (
	"RESTAPI/internal/entity"

	"gorm.io/gorm"
)

type WriterRepository struct {
	BaseRepository *BaseRepository[entity.Writer]
}

func NewWriterRepository(db *gorm.DB) *WriterRepository {
	return &WriterRepository{
		BaseRepository: NewBaseRepository[entity.Writer](db),
	}
}

// Create creates a new writer and populates its ID
func (r *WriterRepository) Create(writer *entity.Writer) error {
	return r.BaseRepository.Create(writer)
}

// GetById gets a writer by ID
func (r *WriterRepository) GetById(id int64) (entity.Writer, error) {
	return r.BaseRepository.GetById(id)
}

// Update updates an existing writer
func (r *WriterRepository) Update(writer *entity.Writer) error {
	return r.BaseRepository.Update(writer)
}

// Delete deletes a writer by ID
func (r *WriterRepository) Delete(id int64) error {
	return r.BaseRepository.Delete(id)
}

// GetAll returns all writers
func (r *WriterRepository) GetAll() ([]entity.Writer, error) {
	var writers []entity.Writer
	result := r.BaseRepository.db.Find(&writers)
	if result.Error != nil {
		return nil, result.Error
	}
	return writers, nil
}

// This would be in your repository/writer-repository.go file
func (r *WriterRepository) GetByLogin(login string) (*entity.Writer, error) {
	var writer entity.Writer
	result := r.BaseRepository.db.Where("login = ?", login).First(&writer)
	if result.Error != nil {
		return nil, result.Error
	}
	return &writer, nil
}
