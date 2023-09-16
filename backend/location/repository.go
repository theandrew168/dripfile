package location

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/memorydb"
)

type Repository interface {
	Create(l *Location) error
	List() ([]*Location, error)
	Read(id uuid.UUID) (*Location, error)
	Delete(id uuid.UUID) error
}

// ensure Repository interface is satisfied
var _ Repository = (*memorydb.MemoryDB[*Location])(nil)

func NewMemoryRepository() *memorydb.MemoryDB[*Location] {
	return memorydb.New[*Location]()
}
