package repository

import (
	"errors"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/memorydb"
)

// ensure ItineraryRepository interface is satisfied
var _ ItineraryRepository = (*MemoryItineraryRepository)(nil)

type ItineraryRepository interface {
	Create(itinerary *domain.Itinerary) error
	List() ([]*domain.Itinerary, error)
	Read(id uuid.UUID) (*domain.Itinerary, error)
	Delete(itinerary *domain.Itinerary) error
}

type MemoryItineraryRepository struct {
	db *memorydb.MemoryDB[*domain.Itinerary]
}

func NewMemoryItineraryRepository() *MemoryItineraryRepository {
	repo := MemoryItineraryRepository{
		db: memorydb.New[*domain.Itinerary](),
	}
	return &repo
}

func (repo *MemoryItineraryRepository) Create(itinerary *domain.Itinerary) error {
	return repo.db.Create(itinerary)
}

func (repo *MemoryItineraryRepository) List() ([]*domain.Itinerary, error) {
	return repo.db.List()
}

func (repo *MemoryItineraryRepository) Read(id uuid.UUID) (*domain.Itinerary, error) {
	itinerary, err := repo.db.Read(id)
	if err != nil {
		switch {
		case errors.Is(err, memorydb.ErrNotFound):
			return nil, ErrNotExist
		default:
			return nil, err
		}
	}

	return itinerary, nil
}

func (repo *MemoryItineraryRepository) Delete(itinerary *domain.Itinerary) error {
	err := itinerary.CheckDelete()
	if err != nil {
		return err
	}

	err = repo.db.Delete(itinerary.ID())
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
