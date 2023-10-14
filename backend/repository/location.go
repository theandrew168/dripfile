package repository

import (
	"errors"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/memorydb"
)

// ensure LocationRepository interface is satisfied
var _ LocationRepository = (*MemoryLocationRepository)(nil)

type LocationRepository interface {
	Create(location *domain.Location) error
	List() ([]*domain.Location, error)
	Read(id uuid.UUID) (*domain.Location, error)
	Delete(location *domain.Location) error
}

type MemoryLocationRepository struct {
	db *memorydb.MemoryDB[*domain.Location]
}

func NewMemoryLocationRepository() *MemoryLocationRepository {
	repo := MemoryLocationRepository{
		db: memorydb.New[*domain.Location](),
	}
	return &repo
}

func (repo *MemoryLocationRepository) Create(location *domain.Location) error {
	return repo.db.Create(location)
}

func (repo *MemoryLocationRepository) List() ([]*domain.Location, error) {
	return repo.db.List()
}

func (repo *MemoryLocationRepository) Read(id uuid.UUID) (*domain.Location, error) {
	location, err := repo.db.Read(id)
	if err != nil {
		switch {
		case errors.Is(err, memorydb.ErrNotFound):
			return nil, ErrNotExist
		default:
			return nil, err
		}
	}

	return location, nil
}

func (repo *MemoryLocationRepository) Delete(location *domain.Location) error {
	err := location.CheckDelete()
	if err != nil {
		return err
	}

	err = repo.db.Delete(location.ID())
	if err != nil {
		switch {
		case errors.Is(err, memorydb.ErrNotFound):
			return ErrNotExist
		default:
			return err
		}
	}

	return nil
}
