package repository

import (
	"errors"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/memorydb"
	"github.com/theandrew168/dripfile/backend/model"
)

// ensure LocationRepository interface is satisfied
var _ LocationRepository = (*MemoryLocationRepository)(nil)

type LocationRepository interface {
	Create(location model.Location) error
	List() ([]model.Location, error)
	Read(id uuid.UUID) (model.Location, error)
	Delete(id uuid.UUID) error
}

type MemoryLocationRepository struct {
	db *memorydb.MemoryDB[model.Location]
}

func NewMemoryLocationRepository() *MemoryLocationRepository {
	repo := MemoryLocationRepository{
		db: memorydb.New[model.Location](),
	}
	return &repo
}

func (repo *MemoryLocationRepository) Create(location model.Location) error {
	return repo.db.Create(location)
}

func (repo *MemoryLocationRepository) List() ([]model.Location, error) {
	return repo.db.List()
}

func (repo *MemoryLocationRepository) Read(id uuid.UUID) (model.Location, error) {
	location, err := repo.db.Read(id)
	if err != nil {
		switch {
		case errors.Is(err, memorydb.ErrNotFound):
			return model.Location{}, ErrNotExist
		default:
			return model.Location{}, err
		}
	}

	return location, nil
}

func (repo *MemoryLocationRepository) Delete(id uuid.UUID) error {
	err := repo.db.Delete(id)
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
