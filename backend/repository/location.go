package repository

import (
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
	return repo.db.Read(id)
}

func (repo *MemoryLocationRepository) Delete(id uuid.UUID) error {
	return repo.db.Delete(id)
}
