package location

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("location not found")

type Repository interface {
	Create(l *Location) error
	List() ([]*Location, error)
	Read(id uuid.UUID) (*Location, error)
	Delete(id uuid.UUID) error
}

// ensure Repository interface is satisfied
var _ Repository = (*MemoryRepository)(nil)

type MemoryRepository struct {
	sync.RWMutex
	data map[uuid.UUID]*Location
}

func NewMemoryRepository() *MemoryRepository {
	repo := MemoryRepository{
		data: make(map[uuid.UUID]*Location),
	}
	return &repo
}

func (repo *MemoryRepository) Create(l *Location) error {
	repo.Lock()
	defer repo.Unlock()

	repo.data[l.ID()] = l
	return nil
}

func (repo *MemoryRepository) List() ([]*Location, error) {
	repo.RLock()
	defer repo.RUnlock()

	var ls []*Location
	for _, l := range repo.data {
		ls = append(ls, l)
	}

	return ls, nil
}

func (repo *MemoryRepository) Read(id uuid.UUID) (*Location, error) {
	repo.RLock()
	defer repo.RUnlock()

	l, ok := repo.data[id]
	if !ok {
		return nil, ErrNotFound
	}

	return l, nil
}

func (repo *MemoryRepository) Delete(id uuid.UUID) error {
	repo.Lock()
	defer repo.Unlock()

	_, ok := repo.data[id]
	if !ok {
		return ErrNotFound
	}

	delete(repo.data, id)
	return nil
}
