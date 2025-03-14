package storage

type Repository[T any] interface {
	Create(entity T) (int64, error)
	GetById(id int64) (T, error)
	Update(entity T) error
	Delete(id int64) error
	GetAll() ([]T, error)
}
