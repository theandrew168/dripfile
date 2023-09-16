package itinerary

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/memorydb"
)

type Repository interface {
	Create(l *Itinerary) error
	List() ([]*Itinerary, error)
	Read(id uuid.UUID) (*Itinerary, error)
	Delete(id uuid.UUID) error
}

// ensure Repository interface is satisfied
var _ Repository = (*memorydb.MemoryDB[*Itinerary])(nil)

type MemoryRepository struct {
	*memorydb.MemoryDB[*Itinerary]
}

func NewMemoryRepository() *memorydb.MemoryDB[*Itinerary] {
	return memorydb.New[*Itinerary]()
}
