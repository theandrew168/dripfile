package repository

import (
	"errors"

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
	itinerary, err := repo.db.Read(id)
	if err != nil {
		switch {
		case errors.Is(err, memorydb.ErrNotFound):
			return model.Itinerary{}, ErrNotExist
		default:
			return model.Itinerary{}, err
		}
	}

	return itinerary, nil
}

func (repo *MemoryItineraryRepository) Delete(id uuid.UUID) error {
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
