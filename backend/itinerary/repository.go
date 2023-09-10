package itinerary

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("itinerary not found")

type Repository interface {
	Create(l *Itinerary) error
	List() ([]*Itinerary, error)
	Read(id uuid.UUID) (*Itinerary, error)
	Delete(id uuid.UUID) error
}

// ensure Repository interface is satisfied
var _ Repository = (*MemoryRepository)(nil)

type MemoryRepository struct {
	sync.RWMutex
	data map[uuid.UUID]*Itinerary
}

func NewMemoryRepository() *MemoryRepository {
	repo := MemoryRepository{
		data: make(map[uuid.UUID]*Itinerary),
	}
	return &repo
}

func (repo *MemoryRepository) Create(l *Itinerary) error {
	repo.Lock()
	defer repo.Unlock()

	repo.data[l.ID()] = l
	return nil
}

func (repo *MemoryRepository) List() ([]*Itinerary, error) {
	repo.RLock()
	defer repo.RUnlock()

	var is []*Itinerary
	for _, i := range repo.data {
		is = append(is, i)
	}

	return is, nil
}

func (repo *MemoryRepository) Read(id uuid.UUID) (*Itinerary, error) {
	repo.RLock()
	defer repo.RUnlock()

	i, ok := repo.data[id]
	if !ok {
		return nil, ErrNotFound
	}

	return i, nil
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
