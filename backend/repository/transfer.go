package repository

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/memorydb"
	"github.com/theandrew168/dripfile/backend/model"
)

// ensure TransferRepository interface is satisfied
var _ TransferRepository = (*MemoryTransferRepository)(nil)

type TransferRepository interface {
	Create(transfer model.Transfer) error
	List() ([]model.Transfer, error)
	Read(id uuid.UUID) (model.Transfer, error)
	Delete(id uuid.UUID) error
}

type MemoryTransferRepository struct {
	db *memorydb.MemoryDB[model.Transfer]
}

func NewMemoryTransferRepository() *MemoryTransferRepository {
	repo := MemoryTransferRepository{
		db: memorydb.New[model.Transfer](),
	}
	return &repo
}

func (repo *MemoryTransferRepository) Create(transfer model.Transfer) error {
	return repo.db.Create(transfer)
}

func (repo *MemoryTransferRepository) List() ([]model.Transfer, error) {
	return repo.db.List()
}

func (repo *MemoryTransferRepository) Read(id uuid.UUID) (model.Transfer, error) {
	return repo.db.Read(id)
}

func (repo *MemoryTransferRepository) Delete(id uuid.UUID) error {
	return repo.db.Delete(id)
}
