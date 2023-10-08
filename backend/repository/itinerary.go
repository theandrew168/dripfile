package repository

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/memorydb"
	"github.com/theandrew168/dripfile/backend/model"
)

// ensure ItineraryRepository interface is satisfied
var _ ItineraryRepository = (*MemoryItineraryRepository)(nil)

type ItineraryRepository interface {
	Create(itinerary model.Itinerary) error
	List() ([]model.Itinerary, error)
	Read(id uuid.UUID) (model.Itinerary, error)
	Delete(id uuid.UUID) error
}

type MemoryItineraryRepository struct {
	db *memorydb.MemoryDB[model.Itinerary]
}

func NewMemoryItineraryRepository() *MemoryItineraryRepository {
	repo := MemoryItineraryRepository{
		db: memorydb.New[model.Itinerary](),
	}
	return &repo
}

func (repo *MemoryItineraryRepository) Create(itinerary model.Itinerary) error {
	return repo.db.Create(itinerary)
}

func (repo *MemoryItineraryRepository) List() ([]model.Itinerary, error) {
	return repo.db.List()
}

func (repo *MemoryItineraryRepository) Read(id uuid.UUID) (model.Itinerary, error) {
	return repo.db.Read(id)
}

func (repo *MemoryItineraryRepository) Delete(id uuid.UUID) error {
	return repo.db.Delete(id)
}
